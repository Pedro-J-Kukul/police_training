package data

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

// Test-specific cleanup helper for training sessions
type trainingSessionTestHelper struct {
	t          *testing.T
	createdIDs []int64
	db         *sql.DB
}

func newTrainingSessionTestHelper(t *testing.T, db *sql.DB) *trainingSessionTestHelper {
	return &trainingSessionTestHelper{
		t:          t,
		createdIDs: make([]int64, 0),
		db:         db,
	}
}

func (h *trainingSessionTestHelper) addSessionID(id int64) {
	h.createdIDs = append(h.createdIDs, id)
	h.t.Logf("Tracking training session ID for cleanup: %d", id)
}

func (h *trainingSessionTestHelper) cleanup() {
	h.t.Helper()
	h.t.Logf("Cleaning up %d training session records", len(h.createdIDs))

	for _, id := range h.createdIDs {
		_, err := h.db.Exec("DELETE FROM training_sessions WHERE id = $1", id)
		if err != nil {
			h.t.Logf("Warning: Failed to cleanup training session ID %d: %v", id, err)
		} else {
			h.t.Logf("Successfully deleted training session ID: %d", id)
		}
	}
	h.createdIDs = nil
}

func TestValidateTrainingSession(t *testing.T) {
	sessionDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		session *TrainingSession
		wantErr bool
	}{
		{
			name: "Valid training session",
			session: &TrainingSession{
				FacilitatorID:    1,
				WorkshopID:       1,
				FormationID:      1,
				RegionID:         1,
				SessionDate:      sessionDate,
				StartTime:        startTime,
				EndTime:          endTime,
				TrainingStatusID: 1,
			},
			wantErr: false,
		},
		{
			name: "Missing facilitator ID",
			session: &TrainingSession{
				WorkshopID:       1,
				FormationID:      1,
				RegionID:         1,
				SessionDate:      sessionDate,
				StartTime:        startTime,
				EndTime:          endTime,
				TrainingStatusID: 1,
			},
			wantErr: true,
		},
		{
			name: "End time before start time",
			session: &TrainingSession{
				FacilitatorID:    1,
				WorkshopID:       1,
				FormationID:      1,
				RegionID:         1,
				SessionDate:      sessionDate,
				StartTime:        endTime,   // Wrong way around
				EndTime:          startTime, // Wrong way around
				TrainingStatusID: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			ValidateTrainingSession(v, tt.session)

			hasErrors := !v.IsEmpty()
			if hasErrors != tt.wantErr {
				t.Errorf("ValidateTrainingSession() error = %v, wantErr %v", hasErrors, tt.wantErr)
				if hasErrors {
					t.Logf("Validation errors: %v", v.Errors)
				}
			}
		})
	}
}

func TestInsertGetUpdateDeleteTrainingSession(t *testing.T) {
	db := setupTestDB(t)
	helper := newTrainingSessionTestHelper(t, db)
	defer helper.cleanup()

	sessionModel := TrainingSessionModel{DB: db}

	// Use existing seed data
	facilitatorID := int64(1) // Assuming user ID 1 exists
	workshopID := int64(1)    // Assuming workshop ID 1 exists
	formationID := int64(1)   // Assuming formation ID 1 exists
	regionID := int64(1)      // Assuming region ID 1 exists
	statusID := int64(1)      // Assuming status ID 1 exists

	sessionDate := time.Date(2025, 12, 25, 0, 0, 0, 0, time.UTC)
	startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
	endTime := time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)
	location := "Test Location"
	maxCapacity := 30

	session := &TrainingSession{
		FacilitatorID:    facilitatorID,
		WorkshopID:       workshopID,
		FormationID:      formationID,
		RegionID:         regionID,
		SessionDate:      sessionDate,
		StartTime:        startTime,
		EndTime:          endTime,
		Location:         &location,
		MaxCapacity:      &maxCapacity,
		TrainingStatusID: statusID,
	}

	// Test Insert
	err := sessionModel.Insert(session)
	if err != nil {
		t.Fatalf("Insert() failed: %v", err)
	}
	helper.addSessionID(session.ID)

	if session.ID == 0 {
		t.Error("Expected session ID to be set after insert")
	}

	// Test Get
	fetched, err := sessionModel.Get(session.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if fetched.FacilitatorID != session.FacilitatorID {
		t.Errorf("Expected facilitator ID %v, got %v", session.FacilitatorID, fetched.FacilitatorID)
	}

	// Test Update
	newCapacity := 50
	session.MaxCapacity = &newCapacity
	session.TrainingStatusID = statusID

	err = sessionModel.Update(session)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	// Verify update
	updated, err := sessionModel.Get(session.ID)
	if err != nil {
		t.Fatalf("Get() after update failed: %v", err)
	}

	if *updated.MaxCapacity != newCapacity {
		t.Errorf("Expected updated max capacity %v, got %v", newCapacity, *updated.MaxCapacity)
	}

	// Test Delete
	err = sessionModel.Delete(session.ID)
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	// Remove from tracking since we deleted it
	for i, id := range helper.createdIDs {
		if id == session.ID {
			helper.createdIDs = append(helper.createdIDs[:i], helper.createdIDs[i+1:]...)
			break
		}
	}

	// Verify deletion
	_, err = sessionModel.Get(session.ID)
	if err == nil {
		t.Error("Expected session to be deleted")
	}
	if err != ErrRecordNotFound {
		t.Errorf("Expected ErrRecordNotFound, got %v", err)
	}
}
