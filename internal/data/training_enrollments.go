// FileName: internal/data/training_enrollments.go
package data

import (
	"context"
	"database/sql"
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
	v.Check(len(*enrollment.CertificateNumber) <= 100, "certificate_number", "must not exceed 100 characters")
	v.Check(enrollment.CertificateIssued || enrollment.CertificateNumber == nil, "certificate_number", "must be nil if certificate is not issued")
}

// Insert creates a new training enrollment.
func (m *TrainingEnrollmentModel) Insert(enrollment *TrainingEnrollment) error {
	query := `
		INSERT INTO training_enrollments (officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	args := []any{
		enrollment.OfficerID,
		enrollment.SessionID,
		enrollment.EnrollmentStatusID,
		enrollment.AttendanceStatusID,
		enrollment.ProgressStatusID,
		enrollment.CompletionDate,
		enrollment.CertificateIssued,
		enrollment.CertificateNumber,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&enrollment.ID, &enrollment.CreatedAt, &enrollment.UpdatedAt); err != nil {
		return err
	}

	return nil
}

// Get retruns a training enrollment by ID.
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
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &enrollment, nil
}

// Update modifies an existing training enrollment.
func (m *TrainingEnrollmentModel) Update(enrollment *TrainingEnrollment) error {
	query := `
		UPDATE training_enrollments
		SET officer_id = $1, session_id = $2, enrollment_status_id = $3, attendance_status_id = $4, progress_status_id = $5, completion_date = $6, certificate_issued = $7, certificate_number = $8, updated_at = NOW()
		WHERE id = $9
		RETURNING updated_at`

	args := []any{
		enrollment.OfficerID,
		enrollment.SessionID,
		enrollment.EnrollmentStatusID,
		enrollment.AttendanceStatusID,
		enrollment.ProgressStatusID,
		enrollment.CompletionDate,
		enrollment.CertificateIssued,
		enrollment.CertificateNumber,
		enrollment.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&enrollment.UpdatedAt)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return ErrRecordNotFound
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		default:
			return err
		}
	}

	return nil
}

// Delete removes a training enrollment by ID.
func (m *TrainingEnrollmentModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM training_enrollments
		WHERE id = $1`

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

// GetAll returns all training enrollments.
func (m *TrainingEnrollmentModel) GetAll(officerID, sessionID, enrollmentStatusID, progressStatusID int64, certificate_issued bool, completion_date time.Time, certificate_number string, filters Filters) ([]*TrainingEnrollment, MetaData, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, officer_id, session_id, enrollment_status_id, attendance_status_id, progress_status_id, completion_date, certificate_issued, certificate_number, created_at, updated_at
		FROM training_enrollments
		WHERE (officer_id = $1 OR $1 = 0)
		AND (session_id = $2 OR $2 = 0)
		AND (enrollment_status_id = $3 OR $3 = 0)
		AND (progress_status_id = $4 OR $4 = 0)
		AND (certificate_issued = COALESCE(NULLIF($5::boolean, false), certificate_issued))
		AND (completion_date = COALESCE(NULLIF($6::date, '0001-01-01'), completion_date))
		AND (to_tsvector('simple', certificate_number) @@ plainto_tsquery('simple', $7) OR $7 = '')
		ORDER BY %s %s, id ASC
		LIMIT $8 OFFSET $9`, filters.sortDirection(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, officerID, sessionID, enrollmentStatusID, progressStatusID, certificate_issued, completion_date, certificate_issued, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	enrollments := []*TrainingEnrollment{}
	totalRecords := 0

	for rows.Next() {
		var enrollment TrainingEnrollment
		err := rows.Scan(
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
		)
		if err != nil {
			return nil, MetaData{}, err
		}
		enrollments = append(enrollments, &enrollment)
	}
	if err = rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	meta := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return enrollments, meta, nil
}
