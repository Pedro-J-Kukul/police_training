// Filename: internal/data/roles.go
package data

import (
	"context"
	"database/sql"
	"slices"
	"time"
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

// Include - CHeck if a role exists in the Roles slice
func (r Roles) Include(role string) bool {
	return slices.Contains(r, role)
}

/*************************************************************************************************************/
// Methods
/*************************************************************************************************************/

// GetAllForUser - Retrieve all roles associated with a specific user
func (m RoleModel) GetAllForUser(userID int64) (Roles, error) {
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
func (m RoleModel) AssignToUser(userID int64, roles ...string) error {

	// Query to insert roles for a user
	query := `
		INSERT INTO roles (user_id, role_id)
		SELECT $1, r.id
		FROM roles r
		WHERE r.role = ANY($2)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set a 3-second timeout
	defer cancel()                                                          // Ensure the context is canceled to free resources

	// Execute the query
	_, err := m.DB.ExecContext(ctx, query, userID, roles)
	if err != nil {
		return err
	}

	return nil // Return nil error on successful assignment
}
