// FileName: internal/data/training_status.go
package data

import "database/sql"

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