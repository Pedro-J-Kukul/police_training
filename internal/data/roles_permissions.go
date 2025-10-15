// FileName: internal/data/roles_permissions.go
package data

import "database/sql"

/************************************************************************************************************/
// RolePermission Declarations
/************************************************************************************************************/

// RolePermission struct to represent the many-to-many relationship between roles and permissions
type RolePermission struct {
	PermissionID int64 `json:"permission_id"`
	RoleID       int64 `json:"role_id"`
}

// RolePermissionModel struct to interact with the roles_permissions table in the database
type RolePermissionModel struct {
	DB *sql.DB
}
