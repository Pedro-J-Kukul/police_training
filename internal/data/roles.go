// Filename: internal/data/roles.go
package data

import (
	"context"
	"database/sql"
	"slices"
	"time"

	"github.com/lib/pq"
)

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

// Include - Check if a role exists in the Roles slice
func (r Roles) Include(role string) bool {
	return slices.Contains(r, role)
}

/*************************************************************************************************************/
// Methods
/*************************************************************************************************************/

// GetAllForUser - Retrieve all roles associated with a specific user
func (m *RoleModel) GetAllForUser(userID int64) (Roles, error) {
	query := `
		SELECT r.role
		FROM roles r
		INNER JOIN roles_users ur ON ur.role_id = r.id
		INNER JOIN users u ON ur.user_id = u.id
		WHERE ur.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set a 3-second timeout
	defer cancel()                                                          // Ensure the context is canceled to free resources

	// Execute the query
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	// Ensure the rows are closed after reading
	defer rows.Close()

	var roles Roles
	// Iterate through the result set and scan each role into the slice
	for rows.Next() {
		var role string // Temporary variable to hold the scanned role
		// Scan the role from the current row
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}
		roles = append(roles, role) // Append the scanned role to the roles slice
	}

	return roles, nil // Return the list of roles and nil error
}

// AssignToUser - Assign a list of roles to a specific user
func (m *RoleModel) AssignToUser(userID int64, roles ...string) error {
	// Query to insert roles for a user
	query := `
        INSERT INTO roles_users (user_id, role_id)
        SELECT $1, r.id
        FROM roles r
        WHERE r.role = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query - Use pq.Array() for the string slice
	_, err := m.DB.ExecContext(ctx, query, userID, pq.Array(roles))
	if err != nil {
		return err
	}

	return nil
}

// GetAllPermissionsForUser - Retrieve all permission codes associated with a specific user
func (m *RoleModel) GetAllPermissionsForUser(userID int64) (Permissions, error) {
	query := `
		SELECT DISTINCT p.code
		FROM permissions p
		INNER JOIN roles_permissions rp ON rp.permission_id = p.id
		INNER JOIN roles_users ru ON ru.role_id = rp.role_id
		WHERE ru.user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		permissions = append(permissions, code)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// HasPermission - Check if a user has a specific permission
func (m *RoleModel) HasPermission(userID int64, permissionCode string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			INNER JOIN roles_permissions rp ON rp.permission_id = p.id
			INNER JOIN roles_users ru ON ru.role_id = rp.role_id
			WHERE ru.user_id = $1 AND p.code = $2
		)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, userID, permissionCode).Scan(&exists)
	return exists, err
}

// HasAllPermissions - Check if a user has all the required permissions
func (m *RoleModel) HasAllPermissions(userID int64, permissionCodes ...string) (bool, error) {
	query := `
		SELECT COUNT(DISTINCT p.code)
		FROM permissions p
		INNER JOIN roles_permissions rp ON rp.permission_id = p.id
		INNER JOIN roles_users ru ON ru.role_id = rp.role_id
		WHERE ru.user_id = $1 AND p.code = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, query, userID, pq.Array(permissionCodes)).Scan(&count)
	if err != nil {
		return false, err
	}

	// User has all permissions if the count matches the number of required permissions
	return count == len(permissionCodes), nil
}
