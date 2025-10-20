// FileName: internal/data/training_sessions.go
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
// TrainingSession Declarations
/************************************************************************************************************/

// TrainingSession struct to represent a training session in the system
// A training sessions is the class instance of a specific workshop
type TrainingSession struct {
	ID               int64     `json:"id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	FacilitatorID    int64     `json:"facilitator_id"`
	WorkshopID       int64     `json:"workshop_id"`
	SessionDate      time.Time `json:"session_date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Location         *string   `json:"location,omitempty"`     // Nullable field
	MaxCapacity      *int64    `json:"max_capacity,omitempty"` // Nullable field
	TrainingStatusID int64     `json:"training_status_id"`
	Notes            *string   `json:"notes,omitempty"` // Nullable field
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TrainingSessionModel struct to interact with the training_sessions table in the database
type TrainingSessionModel struct {
	DB *sql.DB
}

// ValidateTrainingSession ensures training session data is valid.
func ValidateTrainingSession(v *validator.Validator, session *TrainingSession) {
	v.Check(session.FormationID > 0, "formation_id", "must be provided")
	v.Check(session.RegionID > 0, "region_id", "must be provided")
	v.Check(session.FacilitatorID > 0, "facilitator_id", "must be provided")
	v.Check(session.WorkshopID > 0, "workshop_id", "must be provided")
	v.Check(!session.SessionDate.IsZero(), "session_date", "must be provided")
	v.Check(!session.StartTime.IsZero(), "start_time", "must be provided")
	v.Check(!session.EndTime.IsZero(), "end_time", "must be provided")
	v.Check(session.TrainingStatusID > 0, "training_status_id", "must be provided")
	v.Check(len(*session.Location) <= 1000, "location", "must not exceed 1000 characters")
	v.Check(len(*session.Notes) <= 2000, "notes", "must not exceed 2000 characters")
	if session.MaxCapacity != nil {
		v.Check(*session.MaxCapacity > 0, "max_capacity", "must be a positive integer")
	}
	v.Check(session.EndTime.After(session.StartTime), "end_time", "must be after start_time")
}

// Insert creates a new training session.
func (m *TrainingSessionModel) Insert(session *TrainingSession) error {
	query := `
		INSERT INTO training_sessions (formation_id, region_id, facilitator_id, workshop_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	args := []any{
		session.FormationID,
		session.RegionID,
		session.FacilitatorID,
		session.WorkshopID,
		session.SessionDate,
		session.StartTime,
		session.EndTime,
		session.Location,
		session.MaxCapacity,
		session.TrainingStatusID,
		session.Notes,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		switch {
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Get retrieves a training session by id.
func (m *TrainingSessionModel) Get(id int64) (*TrainingSession, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, formation_id, region_id, facilitator_id, workshop_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes, created_at, updated_at
		FROM training_sessions
		WHERE id = $1`

	var session TrainingSession

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.FormationID,
		&session.RegionID,
		&session.FacilitatorID,
		&session.WorkshopID,
		&session.SessionDate,
		&session.StartTime,
		&session.EndTime,
		&session.Location,
		&session.MaxCapacity,
		&session.TrainingStatusID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &session, nil
}

// Update modifies an existing training session.
func (m *TrainingSessionModel) Update(session *TrainingSession) error {
	query := `
		UPDATE training_sessions
		SET formation_id = $1, region_id = $2, facilitator_id = $3, workshop_id = $4, session_date = $5, start_time = $6, end_time = $7, location = $8, max_capacity = $9, training_status_id = $10, notes = $11, updated_at = NOW()
		WHERE id = $12
		RETURNING updated_at`

	args := []any{
		session.FormationID,
		session.RegionID,
		session.FacilitatorID,
		session.WorkshopID,
		session.SessionDate,
		session.StartTime,
		session.EndTime,
		session.Location,
		session.MaxCapacity,
		session.TrainingStatusID,
		session.Notes,
		session.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&session.UpdatedAt)
	if err != nil {
		switch {
		case isForeignKeyViolation(err):
			return ErrForeignKeyViolation
		case isDuplicateKeyViolation(err):
			return ErrDuplicateValue
		default:
			return err
		}
	}

	return nil
}

// Delete removes a training session by id.
func (m *TrainingSessionModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM training_sessions
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

// GetAll retrievs all sessions with filtering and pagination.
func (m *TrainingSessionModel) GetAll(formation_id, region_id, facilitator_id, Workshop_id, training_status_id int64, location, notes string, start_time, end_time, date time.Time, filters Filters) ([]*TrainingSession, MetaData, error) {

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, formation_id, region_id, facilitator_id, workshop_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes, created_at, updated_at
		FROM training_sessions
		WHERE (formation_id = $1 OR $1 = 0)
		AND (region_id = $2 OR $2 = 0)
		AND (facilitator_id = $3 OR $3 = 0)
		AND (workshop_id = $4 OR $4 = 0)
		AND (training_status_id = $5 OR $5 = 0)
		AND (to_tsvector('simple', location) @@ plainto_tsquery('simple', $6) OR $6 = '')
		AND (to_tsvector('simple', notes) @@ plainto_tsquery('simple', $7) OR $7 = '')
		AND (start_time >= COALESCE(NULLIF($8::time, '0001-01-01 00:00:00'), start_time))
		AND (end_time <= COALESCE(NULLIF($9::time, '0001-01-01 00:00:00'), end_time))
		AND (session_date = COALESCE(NULLIF($10::date, '0001-01-01'), session_date))
		ORDER BY %s %s, id ASC
		LIMIT $11 OFFSET $12`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, formation_id, region_id, facilitator_id, Workshop_id, training_status_id, location, notes, start_time, end_time, date, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	totalRecords := 0
	sessions := []*TrainingSession{}

	for rows.Next() {
		var session TrainingSession

		err := rows.Scan(
			&totalRecords,
			&session.ID,
			&session.FormationID,
			&session.RegionID,
			&session.FacilitatorID,
			&session.WorkshopID,
			&session.SessionDate,
			&session.StartTime,
			&session.EndTime,
			&session.Location,
			&session.MaxCapacity,
			&session.TrainingStatusID,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, MetaData{}, err
		}

		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	meta := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return sessions, meta, nil
}
