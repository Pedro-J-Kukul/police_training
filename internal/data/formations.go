// FileName: internal/data/formations.go
package data

import (
	"database/sql"
	"time"
)

/************************************************************************************************************/
// Formation Declarations
/************************************************************************************************************/

// Formation struct to represent a formation in the system
type Formation struct {
	ID        int64     `json:"id"`
	Formation string    `json:"formation"`
	RegionID  int64     `json:"region_id"`
	CreatedAt time.Time `json:"created_at"`
}

// FormationModel struct to interact with the formations table in the database
type FormationModel struct {
	DB *sql.DB
}