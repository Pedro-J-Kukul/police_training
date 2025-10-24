// FileName: internal/data/attendance_statuses.go
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
// AttendanceStatus Declarations
/************************************************************************************************************/

// AttendanceStatus struct to represent an attendance status in the system
type AttendanceStatus struct {
	ID              int64  `json:"id"`
	Status          string `json:"status"`
	CountsAsPresent bool   `json:"counts_as_present"`
}

// AttendanceStatusModel struct to interact with the attendance_statuses table in the database
type AttendanceStatusModel struct {
	DB *sql.DB
}

// ValidateAttendanceStatus ensures attendance status data is valid.
func ValidateAttendanceStatus(v *validator.Validator, status *AttendanceStatus) {
	v.Check(status.Status != "", "status", "must be provided")
	v.Check(len(status.Status) <= 150, "status", "must not exceed 150 characters")
}

// Insert creates a new attendance status.
func (m AttendanceStatusModel) Insert(status *AttendanceStatus) error {
	query := `
		INSERT INTO attendance_statuses (status, counts_as_present)
		VALUES ($1, $2, $3)
		RETURNING id`

	now := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, status.Status, status.CountsAsPresent, now).Scan(&status.ID); err != nil {
		switch {
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves an attendance status by id.
func (m AttendanceStatusModel) Get(id int64) (*AttendanceStatus, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, status, counts_as_present FROM attendance_statuses WHERE id = $1`

	var status AttendanceStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&status.ID, &status.Status, &status.CountsAsPresent)
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

// GetAll returns attendance statuses filtered by name or counts_as_present flag.
func (m AttendanceStatusModel) GetAll(name string, countsAsPresent *bool, filters Filters) ([]*AttendanceStatus, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "status"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, status, counts_as_present
		FROM attendance_statuses
		WHER (to_tsvector('simple', status) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND ($2::boolean IS NULL OR counts_as_present = $2)
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	var countsAsPresentArg any = nil
	if countsAsPresent != nil {
		countsAsPresentArg = *countsAsPresent
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, countsAsPresentArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		statuses     []*AttendanceStatus
		totalRecords int
	)

	for rows.Next() {
		var status AttendanceStatus
		if err := rows.Scan(&totalRecords, &status.ID, &status.Status, &status.CountsAsPresent); err != nil {
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

// Update modifies an existing attendance status.
func (m AttendanceStatusModel) Update(status *AttendanceStatus) error {
	query := `
		UPDATE attendance_statuses
		SET status = $1, counts_as_present = $2
		WHERE id = $3
		RETURNING status, counts_as_present`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query, status.Status, status.CountsAsPresent, status.ID).Scan(&status.Status, &status.CountsAsPresent); err != nil {
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

// Delete removes an attendance status by id.
func (m AttendanceStatusModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM attendance_statuses WHERE id = $1`

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

func (m AttendanceStatusModel) GetByName(name string) (*AttendanceStatus, error) {
	query := `SELECT id, status, counts_as_present FROM attendance_statuses WHERE status = $1`

	var status AttendanceStatus

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, name).Scan(&status.ID, &status.Status, &status.CountsAsPresent)
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
