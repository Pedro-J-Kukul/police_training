// FileName: internal/data/workshops.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// Workshop Declarations
/************************************************************************************************************/

// Workshop struct to represent a workshop in the system
type Workshop struct {
	ID             int64     `json:"id"`
	WorkshopName   string    `json:"workshop_name"`
	CategoryID     int64     `json:"category_id"`
	TrainingTypeID int64     `json:"training_type_id"`
	CreditHours    int       `json:"credit_hours"`
	Description    *string   `json:"description,omitempty"`
	Objectives     *string   `json:"objectives,omitempty"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// WorkshopModel struct to interact with the workshops table in the database
type WorkshopModel struct {
	DB *sql.DB
}
