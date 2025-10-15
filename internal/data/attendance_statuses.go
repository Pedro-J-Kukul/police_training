// FileName: internal/data/attendance_statuses.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// AttendanceStatus Declarations
/************************************************************************************************************/

// AttendanceStatus struct to represent an attendance status in the system
type AttendanceStatus struct {
	ID              int64     `json:"id"`
	Status          string    `json:"status"`
	CountsAsPresent bool      `json:"counts_as_present"`
	CreatedAt       time.Time `json:"created_at"`
}

// AttendanceStatusModel struct to interact with the attendance_statuses table in the database
type AttendanceStatusModel struct {
	DB *sql.DB
}