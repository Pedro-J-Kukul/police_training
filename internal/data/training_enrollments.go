// FileName: internal/data/training_enrollments.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/************************************************************************************************************/
// TrainingEnrollment Declarations
/************************************************************************************************************/

// TrainingEnrollment struct to represent a training enrollment in the system
type TrainingEnrollment struct {
	ID                 int64      `json:"id"`
	OfficerID          int64      `json:"officer_id"`
	SessionID          int64      `json:"session_id"`
	EnrollmentStatusID int64      `json:"enrollment_status_id"`
	AttendanceStatusID *int64     `json:"attendance_status_id,omitempty"`
	ProgressStatusID   int64      `json:"progress_status_id"`
	CompletionDate     *time.Time `json:"completion_date,omitempty"`
	CertificateIssued  bool       `json:"certificate_issued"`
	CertificateNumber  *string    `json:"certificate_number,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

// TrainingEnrollmentModel struct to interact with the training_enrollments table in the database
type TrainingEnrollmentModel struct {
	DB *sql.DB
}

// ValidateTrainingEnrollment ensures training enrollment data is valid.
func ValidateTrainingEnrollment(v *validator.Validator, enrollment *TrainingEnrollment) {
	v.Check(enrollment.OfficerID > 0, "officer_id", "must be provided")
	v.Check(enrollment.SessionID > 0, "session_id", "must be provided")
	v.Check(enrollment.EnrollmentStatusID > 0, "enrollment_status_id", "must be provided")
	v.Check(enrollment.ProgressStatusID > 0, "progress_status_id", "must be provided")

	if enrollment.CertificateNumber != nil {
		v.Check(len(*enrollment.CertificateNumber) <= 100, "certificate_number", "must not exceed 100 characters")
	}

	// If certificate is issued, completion date should be set
	if enrollment.CertificateIssued {
		v.Check(enrollment.CompletionDate != nil, "completion_date", "must be provided when certificate is issued")
	}
}

// Insert creates a new training enrollment.
func (m *TrainingEnrollmentModel) Insert(enrollment *TrainingEnrollment) error {
	query := `
		INSERT INTO training_enrollments (officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		enrollment.OfficerID,
		enrollment.SessionID,
		enrollment.EnrollmentStatusID,
		enrollment.AttendanceStatusID,
		enrollment.ProgressStatusID,
		enrollment.CompletionDate,
		enrollment.CertificateIssued,
		enrollment.CertificateNumber,
	).Scan(&enrollment.ID, &enrollment.CreatedAt, &enrollment.UpdatedAt); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a training enrollment by id.
func (m *TrainingEnrollmentModel) Get(id int64) (*TrainingEnrollment, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number, created_at, updated_at
		FROM training_enrollments
		WHERE id = $1`

	var enrollment TrainingEnrollment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&enrollment.ID,
		&enrollment.OfficerID,
		&enrollment.SessionID,
		&enrollment.EnrollmentStatusID,
		&enrollment.AttendanceStatusID,
		&enrollment.ProgressStatusID,
		&enrollment.CompletionDate,
		&enrollment.CertificateIssued,
		&enrollment.CertificateNumber,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &enrollment, nil
}

// GetAll returns training enrollments filtered by various criteria.
func (m *TrainingEnrollmentModel) GetAll(officerID, sessionID, enrollmentStatusID, attendanceStatusID, progressStatusID *int64, filters Filters) ([]*TrainingEnrollment, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "created_at"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number, created_at, updated_at
		FROM training_enrollments
		WHERE ($1 = 0 OR officer_id = $1)
		AND ($2 = 0 OR session_id = $2)
		AND ($3 = 0 OR enrollment_status_id = $3)
		AND ($4 = 0 OR attendance_status_id = $4)
		AND ($5 = 0 OR progress_status_id = $5)
		ORDER BY %s %s, id ASC
		LIMIT $6 OFFSET $7`, filters.sortColumn(), filters.sortDirection())

	officerArg := int64(0)
	if officerID != nil {
		officerArg = *officerID
	}

	sessionArg := int64(0)
	if sessionID != nil {
		sessionArg = *sessionID
	}

	enrollmentStatusArg := int64(0)
	if enrollmentStatusID != nil {
		enrollmentStatusArg = *enrollmentStatusID
	}

	attendanceStatusArg := int64(0)
	if attendanceStatusID != nil {
		attendanceStatusArg = *attendanceStatusID
	}

	progressStatusArg := int64(0)
	if progressStatusID != nil {
		progressStatusArg = *progressStatusID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, officerArg, sessionArg, enrollmentStatusArg, attendanceStatusArg, progressStatusArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		enrollments  []*TrainingEnrollment
		totalRecords int
	)

	for rows.Next() {
		var enrollment TrainingEnrollment
		if err := rows.Scan(
			&totalRecords,
			&enrollment.ID,
			&enrollment.OfficerID,
			&enrollment.SessionID,
			&enrollment.EnrollmentStatusID,
			&enrollment.AttendanceStatusID,
			&enrollment.ProgressStatusID,
			&enrollment.CompletionDate,
			&enrollment.CertificateIssued,
			&enrollment.CertificateNumber,
			&enrollment.CreatedAt,
			&enrollment.UpdatedAt,
		); err != nil {
			return nil, MetaData{}, err
		}
		enrollments = append(enrollments, &enrollment)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return enrollments, metadata, nil
}

// Update modifies an existing training enrollment.
func (m *TrainingEnrollmentModel) Update(enrollment *TrainingEnrollment) error {
	query := `
		UPDATE training_enrollments
		SET officer_id = $1, session_id = $2, enrollment_status_id = $3, attendance_status_id = $4, progress_status_id = $5, completion_date = $6, certificate_issued = $7, certificate_number = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		enrollment.OfficerID,
		enrollment.SessionID,
		enrollment.EnrollmentStatusID,
		enrollment.AttendanceStatusID,
		enrollment.ProgressStatusID,
		enrollment.CompletionDate,
		enrollment.CertificateIssued,
		enrollment.CertificateNumber,
		enrollment.ID,
	).Scan(&enrollment.UpdatedAt); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Delete removes a training enrollment from the database
func (m *TrainingEnrollmentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM training_enrollments WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// GetByOfficer retrieves enrollments by officer ID
func (m *TrainingEnrollmentModel) GetByOfficer(officerID int64, filters Filters) ([]*TrainingEnrollment, MetaData, error) {
	return m.GetAll(&officerID, nil, nil, nil, nil, filters)
}

// GetBySession retrieves enrollments by session ID
func (m *TrainingEnrollmentModel) GetBySession(sessionID int64, filters Filters) ([]*TrainingEnrollment, MetaData, error) {
	return m.GetAll(nil, &sessionID, nil, nil, nil, filters)
}

// GetByOfficerAndSession retrieves a specific enrollment by officer and session
func (m *TrainingEnrollmentModel) GetByOfficerAndSession(officerID, sessionID int64) (*TrainingEnrollment, error) {
	query := `
		SELECT id, officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number, created_at, updated_at
		FROM training_enrollments
		WHERE officer_id = $1 AND session_id = $2`

	var enrollment TrainingEnrollment

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, officerID, sessionID).Scan(
		&enrollment.ID,
		&enrollment.OfficerID,
		&enrollment.SessionID,
		&enrollment.EnrollmentStatusID,
		&enrollment.AttendanceStatusID,
		&enrollment.ProgressStatusID,
		&enrollment.CompletionDate,
		&enrollment.CertificateIssued,
		&enrollment.CertificateNumber,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &enrollment, nil
}

// IssueCertificate marks an enrollment as having a certificate issued
func (m *TrainingEnrollmentModel) IssueCertificate(enrollmentID int64, certificateNumber string, completionDate time.Time) error {
	query := `
		UPDATE training_enrollments
		SET certificate_issued = true, certificate_number = $1, completion_date = $2, updated_at = NOW()
		WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, certificateNumber, completionDate, enrollmentID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
