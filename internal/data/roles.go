// Filename: internal/data/roles.go
package data

import "database/sql"

/************************************************************************************************************/
// Role Declarations
/************************************************************************************************************/

// Role struct to represent a role in the system
type Role struct {
	ID   int64  `json:"id"`
	Role string `json:"role"`
}

// RoleModel struct to interact with the roles table in the database
type RoleModel struct {
	DB *sql.DB
}

// Roles type to represent a list of roles
type Roles []string
