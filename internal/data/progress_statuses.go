// FileName: internal/data/progress_statuses.go
package data

import (
	"database/sql"
	"time"
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