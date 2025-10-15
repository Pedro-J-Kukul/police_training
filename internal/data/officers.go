// FileName: internal/data/officers.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// Officer Declarations
/************************************************************************************************************/

// Officer struct to represent an officer in the system
type Officer struct {
	ID               int64     `json:"id"`
	RegulationNumber string    `json:"regulation_number"`
	PostingID        int64     `json:"posting_id"`
	RankID           int64     `json:"rank_id"`
	FormationID      int64     `json:"formation_id"`
	RegionID         int64     `json:"region_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// OfficerModel struct to interact with the officers table in the database
type OfficerModel struct {
	DB *sql.DB
}
