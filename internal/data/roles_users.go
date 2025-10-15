// FileName: internal/data/roles_users.go
package data

import "database/sql"

/************************************************************************************************************/
// RoleUser Declarations
/************************************************************************************************************/

// RoleUser struct to represent the many-to-many relationship between roles and users
type RoleUser struct {
	RoleID int64 `json:"role_id"`
	UserID int64 `json:"user_id"`
}

// RoleUserModel struct to interact with the roles_users table in the database
type RoleUserModel struct {
	DB *sql.DB
}