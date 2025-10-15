// Filename: internal/data/permissions.go
package data

import "database/sql"

/************************************************************************************************************/
// Permission Declarations
/************************************************************************************************************/

// Permission struct to represent a permission in the system
type Permission struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

// PermissionModel struct to interact with the permissions table in the database
type PermissionModel struct {
	DB *sql.DB
}

// Permissions type to represent a list of permissions
type Permissions []string
