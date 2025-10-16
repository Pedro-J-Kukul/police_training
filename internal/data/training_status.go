// FileName: internal/data/training_status.go
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
// TrainingStatus Declarations
/************************************************************************************************************/

// TrainingStatus struct to represent a training status in the system
type TrainingStatus struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

// TrainingStatusModel struct to interact with the training_status table in the database
type TrainingStatusModel struct {
	DB *sql.DB
}

// ValidateTrainingStatus ensures status data is valid.
func ValidateTrainingStatus(v *validator.Validator, status *TrainingStatus) {
	v.Check(status.Status != "", "status", "must be provided")
	v.Check(len(status.Status) <= 150, "status", "must not exceed 150 characters")
}

// Insert adds a new training status.
func (m TrainingStatusModel) Insert(status *TrainingStatus) error {
	query := `
		INSERT INTO training_status (status)
		VALUES ($1)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, status.Status).Scan(&status.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a training status by id.
func (m TrainingStatusModel) Get(id int64) (*TrainingStatus, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, status FROM training_status WHERE id = $1`

	var status TrainingStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&status.ID, &status.Status)
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

// GetAll returns training statuses filtered by name.
func (m TrainingStatusModel) GetAll(name string, filters Filters) ([]*TrainingStatus, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "status"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, status
		FROM training_status
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
		statuses     []*TrainingStatus
		totalRecords int
	)

	for rows.Next() {
		var status TrainingStatus
		if err := rows.Scan(&totalRecords, &status.ID, &status.Status); err != nil {
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

// Update modifies an existing training status.
func (m TrainingStatusModel) Update(status *TrainingStatus) error {
	query := `
		UPDATE training_status
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
