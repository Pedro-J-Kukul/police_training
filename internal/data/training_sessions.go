// FileName: internal/data/training_sessions.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// TrainingSession Declarations
/************************************************************************************************************/

// TrainingSession struct to represent a training session in the system
type TrainingSession struct {
	ID               int64     `json:"id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	FacilitatorID    int64     `json:"facilitator_id"`
	WorkshopID       int64     `json:"workshop_id"`
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
