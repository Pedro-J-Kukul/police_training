// FileName: internal/data/training_categories.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// TrainingCategory Declarations
/************************************************************************************************************/

// TrainingCategory struct to represent a training category in the system
type TrainingCategory struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// TrainingCategoryModel struct to interact with the training_categories table in the database
type TrainingCategoryModel struct {
	DB *sql.DB
}