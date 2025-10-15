// FileName: internal/data/training_types.go
package data

import "database/sql"

/************************************************************************************************************/
// TrainingType Declarations
/************************************************************************************************************/

// TrainingType struct to represent a training type in the system
type TrainingType struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

// TrainingTypeModel struct to interact with the training_types table in the database
type TrainingTypeModel struct {
	DB *sql.DB
}