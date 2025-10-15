// FileName: internal/data/enrollment_statuses.go
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
// EnrollmentStatus Declarations
/************************************************************************************************************/

// EnrollmentStatus struct to represent an enrollment status in the system
type EnrollmentStatus struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// EnrollmentStatusModel struct to interact with the enrollment_statuses table in the database
type EnrollmentStatusModel struct {
	DB *sql.DB
}

// ValidateEnrollmentStatus ensures status data is valid.
func ValidateEnrollmentStatus(v *validator.Validator, status *EnrollmentStatus) {
	v.Check(status.Status != "", "status", "must be provided")
	v.Check(len(status.Status) <= 150, "status", "must not exceed 150 characters")
}

// Insert creates a new enrollment status.
func (m EnrollmentStatusModel) Insert(status *EnrollmentStatus) error {
	query := `
		INSERT INTO enrollment_statuses (status, created_at)
		VALUES ($1, $2)
		RETURNING id, created_at`

	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, status.Status, now).Scan(&status.ID, &status.CreatedAt); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves an enrollment status by id.
func (m EnrollmentStatusModel) Get(id int64) (*EnrollmentStatus, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, status, created_at FROM enrollment_statuses WHERE id = $1`

	var status EnrollmentStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&status.ID, &status.Status, &status.CreatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &status, nil
}

// GetAll returns enrollment statuses filtered by name.
func (m EnrollmentStatusModel) GetAll(name string, filters Filters) ([]*EnrollmentStatus, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "status"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, status, created_at
		FROM enrollment_statuses
		WHERE ($1 = '' OR status ILIKE $1)
		ORDER BY %s %s, id ASC
		LIMIT $2 OFFSET $3`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		statuses     []*EnrollmentStatus
		totalRecords int
	)

	for rows.Next() {
		var status EnrollmentStatus
		if err := rows.Scan(&totalRecords, &status.ID, &status.Status, &status.CreatedAt); err != nil {
			return nil, MetaData{}, err
		}
		statuses = append(statuses, &status)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return statuses, metadata, nil
}

// Update modifies an existing enrollment status.
func (m EnrollmentStatusModel) Update(status *EnrollmentStatus) error {
	query := `
		UPDATE enrollment_statuses
		SET status = $1
		WHERE id = $2
		RETURNING status`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, status.Status, status.ID).Scan(&status.Status); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}
