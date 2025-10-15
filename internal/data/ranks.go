// FileName: internal/data/ranks.go
package data

import "database/sql"

/************************************************************************************************************/
// Rank Declarations
/************************************************************************************************************/

// Rank struct to represent a rank in the system
type Rank struct {
	ID                            int64  `json:"id"`
	Rank                          string `json:"rank"`
	Code                          string `json:"code"`
	AnnualTrainingHoursRequired   int    `json:"annual_training_hours_required"`
}

// RankModel struct to interact with the ranks table in the database
type RankModel struct {
	DB *sql.DB
}