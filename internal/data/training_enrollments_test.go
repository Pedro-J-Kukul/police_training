package data

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// Test helper for training enrollments
type trainingEnrollmentTestHelper struct {
	t          *testing.T
	createdIDs []int64
	db         *sql.DB
}

func newTrainingEnrollmentTestHelper(t *testing.T, db *sql.DB) *trainingEnrollmentTestHelper {
	return &trainingEnrollmentTestHelper{
		t:          t,
		createdIDs: make([]int64, 0),
		db:         db,
	}
}

func (h *trainingEnrollmentTestHelper) addEnrollmentID(id int64) {
	h.createdIDs = append(h.createdIDs, id)
	h.t.Logf("Tracking training enrollment ID for cleanup: %d", id)
}

func (h *trainingEnrollmentTestHelper) cleanup() {
	h.t.Helper()
	h.t.Logf("Cleaning up %d training enrollment records", len(h.createdIDs))

	for _, id := range h.createdIDs {
		_, err := h.db.Exec("DELETE FROM training_enrollments WHERE id = $1", id)
		if err != nil {
			h.t.Logf("Warning: Failed to cleanup training enrollment ID %d: %v", id, err)
		} else {
			h.t.Logf("Successfully deleted training enrollment ID: %d", id)
		}
	}
	h.createdIDs = nil
}

func TestValidateTrainingEnrollment(t *testing.T) {
	completionDate := time.Now().AddDate(0, 0, 7) // 7 days from now

	tests := []struct {
		name       string
		enrollment *TrainingEnrollment
		wantErr    bool
	}{
		{
			name: "Valid training enrollment",
			enrollment: &TrainingEnrollment{
				OfficerID:          1,
				SessionID:          1,
				EnrollmentStatusID: 1,
				ProgressStatusID:   1,
				CertificateIssued:  false,
			},
			wantErr: false,
		},
		{
			name: "Missing officer ID",
			enrollment: &TrainingEnrollment{
				SessionID:          1,
				EnrollmentStatusID: 1,
				ProgressStatusID:   1,
			},
			wantErr: true,
		},
		{
			name: "Certificate issued without completion date",
			enrollment: &TrainingEnrollment{
				OfficerID:          1,
				SessionID:          1,
				EnrollmentStatusID: 1,
				ProgressStatusID:   1,
				CertificateIssued:  true, // This should require completion date
			},
			wantErr: true,
		},
		{
			name: "Valid certificate with completion date",
			enrollment: &TrainingEnrollment{
				OfficerID:          1,
				SessionID:          1,
				EnrollmentStatusID: 1,
				ProgressStatusID:   1,
				CertificateIssued:  true,
				CompletionDate:     &completionDate,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateTrainingEnrollment(v, tt.enrollment)

			hasErrors := !v.IsEmpty()
			if hasErrors != tt.wantErr {
				t.Errorf("ValidateTrainingEnrollment() error = %v, wantErr %v", hasErrors, tt.wantErr)
				if hasErrors {
					t.Logf("Validation errors: %v", v.Errors)
				}
			}
		})
	}
}

func TestInsertGetUpdateDeleteTrainingEnrollment(t *testing.T) {
	db := setupTestDB(t)
	helper := newTrainingEnrollmentTestHelper(t, db)
	defer helper.cleanup()

	enrollmentModel := TrainingEnrollmentModel{DB: db}

	// Use existing seed data IDs
	officerID := int64(1)          // Assuming officer ID 1 exists
	sessionID := int64(1)          // Assuming session ID 1 exists
	enrollmentStatusID := int64(1) // Assuming status ID 1 exists
	progressStatusID := int64(1)   // Assuming status ID 1 exists

	enrollment := &TrainingEnrollment{
		OfficerID:          officerID,
		SessionID:          sessionID,
		EnrollmentStatusID: enrollmentStatusID,
		ProgressStatusID:   progressStatusID,
		CertificateIssued:  false,
	}

	// Test Insert
	err := enrollmentModel.Insert(enrollment)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}
	helper.addEnrollmentID(enrollment.ID)

	if enrollment.ID == 0 {
		t.Error("Expected enrollment ID to be set after insert")
	}

	// Test Get
	fetched, err := enrollmentModel.Get(enrollment.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if fetched.OfficerID != enrollment.OfficerID {
		t.Errorf("Expected officer ID %v, got %v", enrollment.OfficerID, fetched.OfficerID)
	}

	// Test GetByOfficerAndSession
	fetchedByOS, err := enrollmentModel.GetByOfficerAndSession(officerID, sessionID)
	if err != nil {
		t.Fatalf("GetByOfficerAndSession() failed: %v", err)
	}

	if fetchedByOS.ID != enrollment.ID {
		t.Errorf("Expected enrollment ID %v, got %v", enrollment.ID, fetchedByOS.ID)
	}

	// Test Update
	enrollment.CertificateIssued = true
	completionDate := time.Now()
	enrollment.CompletionDate = &completionDate
	certificateNumber := "CERT-001"
	enrollment.CertificateNumber = &certificateNumber

	err = enrollmentModel.Update(enrollment)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify update
	updated, err := enrollmentModel.Get(enrollment.ID)
	if err != nil {
		t.Fatalf("Get() after update failed: %v", err)
	}

	if !updated.CertificateIssued {
		t.Error("Expected certificate to be issued after update")
	}

	if updated.CertificateNumber == nil || *updated.CertificateNumber != certificateNumber {
		t.Errorf("Expected certificate number %v, got %v", certificateNumber, updated.CertificateNumber)
	}

	// Test IssueCertificate method
	newCertNumber := "CERT-002"
	newCompletionDate := time.Now().AddDate(0, 0, 1)

	err = enrollmentModel.IssueCertificate(enrollment.ID, newCertNumber, newCompletionDate)
	if err != nil {
		t.Fatalf("IssueCertificate() failed: %v", err)
	}

	// Verify certificate issuance
	updated, err = enrollmentModel.Get(enrollment.ID)
	if err != nil {
		t.Fatalf("Get() after certificate issuance failed: %v", err)
	}

	if !updated.CertificateIssued {
		t.Error("Expected certificate to be issued")
	}

	if updated.CertificateNumber == nil || *updated.CertificateNumber != newCertNumber {
		t.Errorf("Expected certificate number %v, got %v", newCertNumber, updated.CertificateNumber)
	}

	// Test Delete
	err = enrollmentModel.Delete(enrollment.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Remove from tracking since we deleted it
	for i, id := range helper.createdIDs {
		if id == enrollment.ID {
			helper.createdIDs = append(helper.createdIDs[:i], helper.createdIDs[i+1:]...)
			break
		}
	}

	// Verify deletion
	_, err = enrollmentModel.Get(enrollment.ID)
	if err == nil {
		t.Error("Expected enrollment to be deleted")
	}
	if err != ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}

func TestTrainingEnrollmentGetAll(t *testing.T) {
	db := setupTestDB(t)
	helper := newTrainingEnrollmentTestHelper(t, db)
	defer helper.cleanup()

	enrollmentModel := TrainingEnrollmentModel{DB: db}

	// Create multiple test enrollments
	numEnrollments := 3
	officerID := int64(1)
	sessionID := int64(1)

	for i := 0; i < numEnrollments; i++ {
		enrollment := &TrainingEnrollment{
			OfficerID:          officerID,
			SessionID:          sessionID + int64(i), // Different sessions
			EnrollmentStatusID: 1,
			ProgressStatusID:   1,
			CertificateIssued:  false,
		}
		err := enrollmentModel.Insert(enrollment)
		if err != nil {
			t.Fatalf("Insert() failed: %v", err)
		}
		helper.addEnrollmentID(enrollment.ID)
	}

	// Test GetAll
	filters := Filters{
		Page:         1,
		PageSize:     10,
		Sort:         "id",
		SortSafelist: []string{"id", "-id", "created_at", "-created_at"},
	}

	enrollments, metadata, err := enrollmentModel.GetAll(nil, nil, nil, nil, nil, filters)
	if err != nil {
		t.Fatalf("GetAll() failed: %v", err)
	}

	if len(enrollments) < numEnrollments {
		t.Errorf("Expected at least %d enrollments, got %d", numEnrollments, len(enrollments))
	}

	if metadata.TotalRecords < numEnrollments {
		t.Errorf("Expected total records at least %d, got %d", numEnrollments, metadata.TotalRecords)
	}

	// Test GetByOfficer
	officerEnrollments, _, err := enrollmentModel.GetByOfficer(officerID, filters)
	if err != nil {
		t.Fatalf("GetByOfficer() failed: %v", err)
	}

	if len(officerEnrollments) < numEnrollments {
		t.Errorf("Expected at least %d enrollments for officer, got %d", numEnrollments, len(officerEnrollments))
	}
}
