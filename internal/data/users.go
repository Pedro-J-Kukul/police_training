// FileName: internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

/************************************************************************************************************/
// User declarations
/************************************************************************************************************/

// Password stores the hashed password and optional plaintext (used for validation during write operations).
type Password struct {
	hash      []byte
	plaintext *string
}

// User represents an application user.
type User struct {
	ID          int64     `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Gender      string    `json:"gender"`
	Password    Password  `json:"-"`
	Activated   bool      `json:"activated"`
	Facilitator bool      `json:"facilitator"`
	Version     int       `json:"version"`
	IsOfficer   bool      `json:"is_officer"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserModel wraps the database connection pool for user CRUD operations.
type UserModel struct {
	DB *sql.DB
}

// AnonymousUser is a sentinel anonymous user instance.
var AnonymousUser = &User{}

/************************************************************************************************************/
// Password helpers
/************************************************************************************************************/

// Set hashes the supplied plaintext password.
func (p *Password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.hash = hash
	p.plaintext = &plaintext
	return nil
}

// Matches verifies that the supplied plaintext password matches the stored hash.
func (p *Password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

/************************************************************************************************************/
// User Validation
/************************************************************************************************************/

// IsAnonymous checks if the user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser // Return true if the user is the anonymous user
}

// ValidateEmail checks if the email is valid
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")                                      // Check if email is not empty
	v.Check(len(email) <= 254, "email", "must not be more than 254 bytes long")            // Check if email length is within limit
	v.Check(v.Matches(email, validator.EmailRX), "email", "must be a valid email address") // Check if email matches the regex
}

// ValidatePasswordPlaintext checks if the plaintext password is valid
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")                                                              // Check if password is not empty
	v.Check(len(password) >= 8, "password", "must be at least 8 characters long")                                        // Check if password length is at least 8 characters
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters long")                                 // Check if password length is within limit
	v.Check(v.Matches(password, validator.PasswordNumberRX), "password", "must contain at least one number")             // Check if password contains at least one number
	v.Check(v.Matches(password, validator.PasswordUpperRX), "password", "must contain at least one uppercase letter")    // Check if password contains at least one uppercase letter
	v.Check(v.Matches(password, validator.PasswordLowerRX), "password", "must contain at least one lowercase letter")    // Check if password contains at least one lowercase letter
	v.Check(v.Matches(password, validator.PasswordSpecialRX), "password", "must contain at least one special character") // Check if password contains at least one special character
}

// ValidateUser checks if the user struct is valid
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.FirstName != "", "first_name", "must be provided")                              // Check if first name is not empty
	v.Check(len(user.FirstName) <= 50, "first_name", "must not be more than 50 characters long") // Check if first name length is within limit
	v.Check(user.LastName != "", "last_name", "must be provided")                                // Check if last name is not empty
	v.Check(len(user.LastName) <= 50, "last_name", "must not be more than 50 characters long")   // Check if last name length is within limit
	ValidateEmail(v, user.Email)                                                                 // Validate the email
	if user.Password.plaintext != nil {                                                          // If plaintext password is set, validate it
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}
	if user.Password.hash == nil { // Check if the password hash is set
		panic("missing password hash for user")
	}
	v.Check(user.Gender != "", "gender", "must be provided")                     // Check if gender is not empty
	v.Check(v.Permitted(user.Gender, "m", "f"), "gender", "must be 'm', or 'f'") // Check if gender is one of the permitted values
	v.Check(len(user.Gender) == 1, "gender", "must only be 1 character long")
	if user.IsOfficer {
		v.Check(user.IsOfficer, "is_officer", "must be true or false") // Check if IsOfficer is true or false
	}
}

/************************************************************************************************************/
// User Model Methods
/************************************************************************************************************/

// Insert adds a new user to the database
func (m *UserModel) Insert(user *User) error {
	// SQL query to insert a new user
	query := `
		INSERT INTO users (first_name, last_name, gender, email, password_hash, is_activated, is_facilitator, is_officer)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at, version`

	// Arguments for the SQL query
	args := []any{
		user.FirstName,
		user.LastName,
		user.Gender,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.Facilitator,
		user.IsOfficer,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// Execute the query and scan the returned values into the user struct
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail // Email already exists
		default:
			return err // Some other error occurred
		}
	}

	return nil // Everything went well
}

// GetByEmail retrieves a user by their email address
func (m *UserModel) GetByEmail(email string) (*User, error) {
	// SQL query to select a user by email
	query := `
		SELECT id, first_name, last_name, email, gender, password_hash, is_activated, is_facilitator, created_at, updated_at, version, is_officer
		FROM users
		WHERE email = $1`

	var user User // Variable to hold the user data

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// Execute the query and scan the result into the user struct
	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Gender,
		&user.Password.hash,
		&user.Activated,
		&user.Facilitator,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
		&user.IsOfficer,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound // No user found with the given email
		default:
			return nil, err // Some other error occurred
		}
	}
	return &user, nil // Return the user data
}

// Update modifies an existing user in the database
func (m *UserModel) Update(user *User) error {
	// SQL query to update a user
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, email = $3, gender = $4, password_hash = $5, is_activated = $6, is_facilitator = $7, updated_at = now(), version = version + 1, is_officer = $10
		WHERE id = $8 AND version = $9
		RETURNING updated_at, version`

	// Arguments for the SQL query
	args := []any{
		user.FirstName,
		user.LastName,
		user.Email,
		user.Gender,
		user.Password.hash,
		user.Activated,
		user.Facilitator,
		user.ID,
		user.Version,
		user.IsOfficer,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// Execute the query and scan the returned values into the user struct
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.UpdatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail // Email already exists
		case err == sql.ErrNoRows:
			return ErrEditConflict // Edit conflict occurred
		default:
			return err // Some other error occurred
		}
	}

	return nil // Everything went well
}

// Get retrieves a user by their ID
func (m *UserModel) Get(id int64) (*User, error) {
	// SQL query to select a user by ID
	query := `
		SELECT id, first_name, last_name, email, gender, password_hash, is_activated, is_facilitator, created_at, updated_at, version
		FROM users
		WHERE id = $1`

	var user User // Variable to hold the user data

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// Execute the query and scan the result into the user struct
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Gender,
		&user.Password.hash,
		&user.Activated,
		&user.Facilitator,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound // No user found with the given ID
		default:
			return nil, err // Some other error occurred
		}
	}
	return &user, nil // Return the user data
}

// GetAll retrieves all users from the database
func (m *UserModel) GetAll(fname, lname, email, gender string, activated *bool, facilitator *bool, filters Filters) ([]*User, MetaData, error) {
	// SQL query to select users with filtering, sorting, and pagination
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, first_name, last_name, email, gender, is_activated, is_facilitator, created_at, updated_at, version
		FROM users
		WHERE (to_tsvector('simple', first_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', last_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
		AND (gender = $4 OR $4 = '')
		AND (is_activated = $5 OR $5 IS NULL)
		AND (is_facilitator = $6 OR $6 IS NULL)
		ORDER BY %s %s, id ASC
		LIMIT $7 OFFSET $8`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// execute the query
	rows, err := m.DB.QueryContext(ctx, query, fname, lname, email, gender, activated, facilitator, filters.limit(), filters.offset())
	if err != nil {
		return nil, MetaData{}, err // Return any error encountered while executing the query
	}
	defer rows.Close() // Ensure the rows are closed after reading

	totalRecords := 0  // Variable to hold the total number of records
	users := []*User{} // Slice to hold the retrieved users
	for rows.Next() {  // Iterate over the rows
		var user User     // Variable to hold the user data
		err := rows.Scan( // Scan the row into the user struct
			&totalRecords,
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Gender,
			&user.Activated,
			&user.Facilitator,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Version,
		)
		if err != nil {
			return nil, MetaData{}, err // Return any error encountered while scanning the row
		}
		users = append(users, &user) // Add the user to the slice
	}
	if err := rows.Err(); err != nil {
		return nil, MetaData{}, err // Return any error encountered while iterating over the rows
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize) // Calculate pagination metadata
	return users, metadata, nil
}

/*************************************************************************************************************/
// Tokens
/*************************************************************************************************************/

// GetForToken retrieves a user based on a token scope and plaintext token
func (m *UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	// SQL query to select a user based on a token
	query := `
		SELECT u.id, u.first_name, u.last_name, u.email, u.gender, u.password_hash, u.is_activated, u.is_facilitator, u.created_at, u.updated_at, u.version
		FROM users u
		INNER JOIN tokens t ON u.id = t.user_id
		WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3`

	// Arguments for the SQL query
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	args := []any{
		tokenHash[:], // Convert array to slice for SQL compatibility
		tokenScope,   // Token scope
		time.Now(),   // Current time to check for expiry
	}
	var user User // Variable to hold the user data

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// Execute the query and scan the result into the user struct
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Gender,
		&user.Password.hash,
		&user.Activated,
		&user.Facilitator,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err == sql.ErrNoRows:
			return nil, ErrRecordNotFound // No user found with the given token
		default:
			return nil, err // Some other error occurred
		}
	}
	return &user, nil // Return the user data
}

func (m *UserModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, id)
	return err
}

// Check if a user has all required permissions through their roles
func (m *UserModel) IsAllowedTo(id int64, requiredPermissions ...string) (bool, error) {
	// First, check which superuser permissions the user has
	superuserQuery := `
		SELECT DISTINCT p.code
		FROM roles_users ru
		INNER JOIN roles r ON ru.role_id = r.id
		INNER JOIN roles_permissions rp ON r.id = rp.role_id
		INNER JOIN permissions p ON rp.permission_id = p.id
		WHERE ru.user_id = $1 AND p.code IN ('CAN_READ', 'CAN_CREATE', 'CAN_DELETE', 'CAN_MODIFY')`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, superuserQuery, id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	// Create boolean map of superuser permissions
	superuserPermissions := map[string]bool{
		"CAN_READ":   false,
		"CAN_CREATE": false,
		"CAN_DELETE": false,
		"CAN_MODIFY": false,
	}

	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return false, err
		}
		superuserPermissions[permission] = true
	}
	if err := rows.Err(); err != nil {
		return false, err
	}

	// Print superuser permissions for debugging
	fmt.Printf("Superuser Permissions: %v\n", superuserPermissions)

	// Filter out scoped permissions if user has corresponding superuser permission
	var filteredPermissions []string
	for _, perm := range requiredPermissions {
		// Skip scoped permissions if user has corresponding superuser permission
		switch {
		case perm == "CAN_READ" || (len(perm) > 9 && perm[:9] == "CAN_READ_"):
			IsAllowedToDoThis(superuserPermissions, perm, "CAN_READ", &filteredPermissions)

		case perm == "CAN_CREATE" || (len(perm) > 11 && perm[:11] == "CAN_CREATE_"):
			IsAllowedToDoThis(superuserPermissions, perm, "CAN_CREATE", &filteredPermissions)

		case perm == "CAN_DELETE" || (len(perm) > 11 && perm[:11] == "CAN_DELETE_"):
			IsAllowedToDoThis(superuserPermissions, perm, "CAN_DELETE", &filteredPermissions)

		case perm == "CAN_MODIFY" || (len(perm) > 11 && perm[:11] == "CAN_MODIFY_"):
			IsAllowedToDoThis(superuserPermissions, perm, "CAN_MODIFY", &filteredPermissions)

		default:
			print("should almost never trigger")
			filteredPermissions = append(filteredPermissions, perm)
		}
	}

	// If all permissions were filtered out (user has all required superuser permissions), return true
	if len(filteredPermissions) == 0 {
		return true, nil
	}

	// Print out evaluated permissions for debugging
	fmt.Printf("Evaluated Permissions: %v\n", filteredPermissions)

	// SQL query to get remaining permissions the user has through their roles
	query := `
		SELECT COUNT(DISTINCT p.code)
		FROM roles_users ru
		INNER JOIN roles r ON ru.role_id = r.id
		INNER JOIN roles_permissions rp ON r.id = rp.role_id
		INNER JOIN permissions p ON rp.permission_id = p.id
		WHERE ru.user_id = $1 AND p.code = ANY($2)`

	var count int
	err = m.DB.QueryRowContext(ctx, query, id, pq.Array(filteredPermissions)).Scan(&count)
	if err != nil {
		return false, err
	}

	// Return true only if the user has ALL remaining required permissions
	return count == len(filteredPermissions), nil
}

func IsAllowedToDoThis(superuserPermissions map[string]bool, perm string, target string, filteredPermissions *[]string) {

	// Superuser overrides
	if superuserPermissions[target] && perm == target {
		*filteredPermissions = append(*filteredPermissions, perm)
	}

	// Scoped permissions
	if !superuserPermissions[target] && perm != target {
		*filteredPermissions = append(*filteredPermissions, perm)
	}
}
