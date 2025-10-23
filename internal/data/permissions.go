// Filename: internal/data/permissions.go
package data

import (
	"context"
	"database/sql"
	"slices"
	"time"
)

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

// Includes - Check if a specific permission code exists in the Permissions slice
func (p Permissions) Includes(code string) bool {
	return slices.Contains(p, code)
}

/*************************************************************************************************************/
// Methods
/*************************************************************************************************************/

// GetAllFor Role - Retrieve all permissions associated with a specific role
func (m PermissionModel) GetAllForRole(roleID int64) (Permissions, error) {
	query := `
		SELECT p.code
		FROM permissions p
		INNER JOIN roles_permissions rp ON rp.permission_id = p.id
		INNER JOIN roles r ON rp.role_id = r.id
		WHERE rp.role_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set a 3-second timeout
	defer cancel()                                                          // Ensure the context is canceled to free resources

	// Execute the query
	rows, err := m.DB.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}

	// Ensure the rows are closed after reading
	defer rows.Close()

	var permissions Permissions // Initialize an empty slice to hold permissions
	// Iterate through the result set and scan each permission code into the slice
	for rows.Next() {
		var code string // Temporary variable to hold the scanned code
		// Scan the code from the current row
		if err := rows.Scan(&code); err != nil {
			return nil, err
		}
		permissions = append(permissions, code) // Append the scanned code to the permissions slice
	}

	return permissions, nil // Return the list of permissions and nil error
}

// AssignToRole - Assign a list of permissions to a specific role
func (m PermissionModel) AssignToRole(roleID int64, codes ...string) error {

	// Prepare the SQL statement for inserting role-permission associations
	query := `
		INSERT INTO roles_permissions (role_id, permission_id)
		SELECT $1, p.id
		FROM permissions p
		WHERE p.code = ANY($2)
		ON CONFLICT (role_id, permission_id) DO NOTHING`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Set a 3-second timeout
	defer cancel()                                                          // Ensure the context is canceled to free resources

	// Execute the insert statement with the provided role ID and permission codes
	_, err := m.DB.ExecContext(ctx, query, roleID, codes)
	if err != nil {
		return err
	}

	return nil // Return nil error if the operation is successful
}
