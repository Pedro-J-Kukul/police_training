// FileName: internal/data/enrollment_statuses.go
package data

import (
	"database/sql"
	"time"
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
