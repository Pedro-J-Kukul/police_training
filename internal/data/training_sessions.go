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
type TrainingSession struct {
	ID               int64     `json:"id"`
	FacilitatorID    int64     `json:"facilitator_id"`
	WorkshopID       int64     `json:"workshop_id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	SessionDate      time.Time `json:"session_date"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	Location         *string   `json:"location,omitempty"`
	MaxCapacity      *int      `json:"max_capacity,omitempty"`
	TrainingStatusID int64     `json:"training_status_id"`
	Notes            *string   `json:"notes,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TrainingSessionModel struct to interact with the training_sessions table in the database
type TrainingSessionModel struct {
	DB *sql.DB
}

// ValidateTrainingSession ensures training session data is valid.
func ValidateTrainingSession(v *validator.Validator, session *TrainingSession) {
	v.Check(session.FacilitatorID > 0, "facilitator_id", "must be provided")
	v.Check(session.WorkshopID > 0, "workshop_id", "must be provided")
	v.Check(session.FormationID > 0, "formation_id", "must be provided")
	v.Check(session.RegionID > 0, "region_id", "must be provided")
	v.Check(!session.SessionDate.IsZero(), "session_date", "must be provided")
	v.Check(!session.StartTime.IsZero(), "start_time", "must be provided")
	v.Check(!session.EndTime.IsZero(), "end_time", "must be provided")
	v.Check(session.EndTime.After(session.StartTime), "end_time", "must be after start time")
	v.Check(session.TrainingStatusID > 0, "training_status_id", "must be provided")

	if session.MaxCapacity != nil {
		v.Check(*session.MaxCapacity > 0, "max_capacity", "must be greater than zero")
	}

	if session.Location != nil {
		v.Check(len(*session.Location) <= 255, "location", "must not exceed 255 characters")
	}
}

// Insert creates a new training session.
func (m *TrainingSessionModel) Insert(session *TrainingSession) error {
	query := `
		INSERT INTO training_sessions (facilitator_id, workshop_id, formation_id, region_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		session.FacilitatorID,
		session.WorkshopID,
		session.FormationID,
		session.RegionID,
		session.SessionDate,
		session.StartTime,
		session.EndTime,
		session.Location,
		session.MaxCapacity,
		session.TrainingStatusID,
		session.Notes,
	).Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt); err != nil {
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

// Get retrieves a training session by id.
func (m *TrainingSessionModel) Get(id int64) (*TrainingSession, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, facilitator_id, workshop_id, formation_id, region_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes, created_at, updated_at
		FROM training_sessions
		WHERE id = $1`

	var session TrainingSession

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&session.ID,
		&session.FacilitatorID,
		&session.WorkshopID,
		&session.FormationID,
		&session.RegionID,
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

// GetAll returns training sessions filtered by various criteria.
func (m *TrainingSessionModel) GetAll(facilitatorID, workshopID, formationID, regionID, statusID *int64, sessionDate *time.Time, filters Filters) ([]*TrainingSession, MetaData, error) {
	if filters.Sort == "" {
		filters.Sort = "session_date"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, facilitator_id, workshop_id, formation_id, region_id, session_date, start_time, end_time, location, max_capacity, training_status_id, notes, created_at, updated_at
		FROM training_sessions
		WHERE ($1 = 0 OR facilitator_id = $1)
		AND ($2 = 0 OR workshop_id = $2)
		AND ($3 = 0 OR formation_id = $3)
		AND ($4 = 0 OR region_id = $4)
		AND ($5 = 0 OR training_status_id = $5)
		AND ($6::date IS NULL OR session_date = $6::date)
		ORDER BY %s %s, id ASC
		LIMIT $7 OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	facilitatorArg := int64(0)
	if facilitatorID != nil {
		facilitatorArg = *facilitatorID
	}

	workshopArg := int64(0)
	if workshopID != nil {
		workshopArg = *workshopID
	}

	formationArg := int64(0)
	if formationID != nil {
		formationArg = *formationID
	}

	regionArg := int64(0)
	if regionID != nil {
		regionArg = *regionID
	}

	statusArg := int64(0)
	if statusID != nil {
		statusArg = *statusID
	}

	var dateArg interface{} = nil
	if sessionDate != nil {
		dateArg = *sessionDate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, facilitatorArg, workshopArg, formationArg, regionArg, statusArg, dateArg, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err
	}
	defer rows.Close()

	var (
		sessions     []*TrainingSession
		totalRecords int
	)

	for rows.Next() {
		var session TrainingSession
		if err := rows.Scan(
			&totalRecords,
			&session.ID,
			&session.FacilitatorID,
			&session.WorkshopID,
			&session.FormationID,
			&session.RegionID,
			&session.SessionDate,
			&session.StartTime,
			&session.EndTime,
			&session.Location,
			&session.MaxCapacity,
			&session.TrainingStatusID,
			&session.Notes,
			&session.CreatedAt,
			&session.UpdatedAt,
		); err != nil {
			return nil, MetaData{}, err
		}
		sessions = append(sessions, &session)
	}

	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return sessions, metadata, nil
}

// Update modifies an existing training session.
func (m *TrainingSessionModel) Update(session *TrainingSession) error {
	query := `
		UPDATE training_sessions
		SET facilitator_id = $1, workshop_id = $2, formation_id = $3, region_id = $4, session_date = $5, start_time = $6, end_time = $7, location = $8, max_capacity = $9, training_status_id = $10, notes = $11, updated_at = NOW()
		WHERE id = $12
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := m.DB.QueryRowContext(ctx, query,
		session.FacilitatorID,
		session.WorkshopID,
		session.FormationID,
		session.RegionID,
		session.SessionDate,
		session.StartTime,
		session.EndTime,
		session.Location,
		session.MaxCapacity,
		session.TrainingStatusID,
		session.Notes,
		session.ID,
	).Scan(&session.UpdatedAt); err != nil {
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

// Delete removes a training session from the database
func (m *TrainingSessionModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM training_sessions WHERE id = $1`

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

// GetByFacilitator retrieves training sessions by facilitator ID
func (m *TrainingSessionModel) GetByFacilitator(facilitatorID int64, filters Filters) ([]*TrainingSession, MetaData, error) {
	return m.GetAll(&facilitatorID, nil, nil, nil, nil, nil, filters)
}

// GetByWorkshop retrieves training sessions by workshop ID
func (m *TrainingSessionModel) GetByWorkshop(workshopID int64, filters Filters) ([]*TrainingSession, MetaData, error) {
	return m.GetAll(nil, &workshopID, nil, nil, nil, nil, filters)
}

// GetByDate retrieves training sessions by session date
func (m *TrainingSessionModel) GetByDate(sessionDate time.Time, filters Filters) ([]*TrainingSession, MetaData, error) {
	return m.GetAll(nil, nil, nil, nil, nil, &sessionDate, filters)
}
