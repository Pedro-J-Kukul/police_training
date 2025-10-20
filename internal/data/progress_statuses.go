// FileName: internal/data/progress_statuses.go
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
// ProgressStatus Declarations
/************************************************************************************************************/

// ProgressStatus struct to represent a progress status in the system
type ProgressStatus struct {
	ID        int64     `json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ProgressStatusModel struct to interact with the progress_statuses table in the database
type ProgressStatusModel struct {
	DB *sql.DB
}

// ValidateProgressStatus ensures progress status data is valid.
func ValidateProgressStatus(v *validator.Validator, status *ProgressStatus) {
	v.Check(status.Status != "", "status", "must be provided")
	v.Check(len(status.Status) <= 150, "status", "must not exceed 150 characters")
}

// Insert creates a new progress status.
func (m ProgressStatusModel) Insert(status *ProgressStatus) error {
	query := `
		INSERT INTO progress_statuses (status, created_at)
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

// Get retrieves a progress status by id.
func (m ProgressStatusModel) Get(id int64) (*ProgressStatus, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, status, created_at FROM progress_statuses WHERE id = $1`

	var status ProgressStatus

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

// GetAll returns progress statuses filtered by name.
func (m ProgressStatusModel) GetAll(name string, filters Filters) ([]*ProgressStatus, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "status"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, status, created_at
		FROM progress_statuses
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
		statuses     []*ProgressStatus
		totalRecords int
	)

	for rows.Next() {
		var status ProgressStatus
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

// Update modifies an existing progress status.
func (m ProgressStatusModel) Update(status *ProgressStatus) error {
	query := `
		UPDATE progress_statuses
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
