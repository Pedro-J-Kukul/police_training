// FileName: internal/data/regions.go
package data

import "database/sql"

/************************************************************************************************************/
// Region Declarations
/************************************************************************************************************/

// Region struct to represent a region in the system
type Region struct {
	ID     int64  `json:"id"`
	Region string `json:"region"`
}

// RegionModel struct to interact with the regions table in the database
type RegionModel struct {
	DB *sql.DB
}