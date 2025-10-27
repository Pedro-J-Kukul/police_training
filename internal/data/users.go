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
	ID            int64     `json:"id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Gender        string    `json:"gender"`
	Password      Password  `json:"-"`
	IsActivated   bool      `json:"is_activated"`
	IsFacilitator bool      `json:"is_facilitator"`
	IsOfficer     bool      `json:"is_officer"`
	IsDeleted     bool      `json:"is_deleted"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
	query := `
        INSERT INTO users (first_name, last_name, gender, email, password_hash, is_activated, is_facilitator, is_officer)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at, version`

	// Arguments for the SQL query - FIXED ORDER
	args := []any{
		user.FirstName,     // $1 - first_name
		user.LastName,      // $2 - last_name
		user.Gender,        // $3 - gender
		user.Email,         // $4 - email
		user.Password.hash, // $5 - password_hash
		user.IsActivated,   // $6 - is_activated
		user.IsFacilitator, // $7 - is_facilitator
		user.IsOfficer,     // $8 - is_officer
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
		SELECT id, first_name, last_name, email, gender, password_hash, is_activated, is_facilitator, is_officer, is_deleted, created_at, updated_at, version
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
		&user.IsActivated,
		&user.IsFacilitator,
		&user.IsOfficer,
		&user.IsDeleted,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Version,
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
		SET first_name = $1, last_name = $2, email = $3, gender = $4, 
		    password_hash = $5, is_activated = $6, is_facilitator = $7, 
		    is_officer = $8, updated_at = now(), version = version + 1
		WHERE id = $9 AND version = $10
		RETURNING updated_at, version`

	// Arguments for the SQL query
	args := []any{
		user.FirstName,     // $1
		user.LastName,      // $2
		user.Email,         // $3
		user.Gender,        // $4
		user.Password.hash, // $5
		user.IsActivated,   // $6
		user.IsFacilitator, // $7
		user.IsOfficer,     // $8
		user.ID,            // $9
		user.Version,       // $10
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

// UpdatePassword updates a user's password securely.
func (m *UserModel) UpdatePassword(id int64, newPlaintext string) error {
	var password Password
	if err := password.Set(newPlaintext); err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password_hash = $1, updated_at = now(), version = version + 1
		WHERE id = $2
		RETURNING updated_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var updatedAt time.Time
	var version int
	err := m.DB.QueryRowContext(ctx, query, password.hash, id).Scan(&updatedAt, &version)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// SoftDelete marks a user as deleted without removing their record.
func (m *UserModel) SoftDelete(id int64) error {
	query := `
		UPDATE users
		SET is_deleted = TRUE, updated_at = now(), version = version + 1
		WHERE id = $1 AND is_deleted = FALSE
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var returnedID int64
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

// Restore revokes the soft deletion of a user.
func (m *UserModel) Restore(id int64) error {
	query := `
		UPDATE users
		SET is_deleted = FALSE, updated_at = now(), version = version + 1
		WHERE id = $1 AND is_deleted = TRUE
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var returnedID int64
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&returnedID)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRecordNotFound
		}
		return err
	}

	return nil
}

// HardDelete permanently removes a user from the database.
func (m *UserModel) HardDelete(id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// Get retrieves a user by their ID
func (m *UserModel) Get(id int64) (*User, error) {
	// SQL query to select a user by ID
	query := `
		SELECT id, first_name, last_name, email, gender, password_hash, is_activated, is_facilitator, is_officer, is_deleted, created_at, updated_at, version
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
		&user.IsActivated,
		&user.IsFacilitator,
		&user.IsOfficer,
		&user.IsDeleted,
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
func (m *UserModel) GetAll(fname, lname, email, gender string, activated *bool, facilitator *bool, officer *bool, deleted *bool, filters Filters) ([]*User, MetaData, error) {
	// SQL query to select users with filtering, sorting, and pagination
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, first_name, last_name, email, gender, is_activated, is_facilitator, is_officer, is_deleted, created_at, updated_at, version
		FROM users
		WHERE (to_tsvector('simple', first_name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', last_name) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (to_tsvector('simple', email) @@ plainto_tsquery('simple', $3) OR $3 = '')
		AND (to_tsvector('simple', gender) @@ plainto_tsquery('simple', $4) OR $4 = '')
		AND (is_activated::boolean = $5 OR $5 IS NULL)
		AND (is_facilitator::boolean = $6 OR $6 IS NULL)
		AND (is_officer::boolean = $7 OR $7 IS NULL)
		AND (is_deleted::boolean = $8 OR $8 IS NULL)
		ORDER BY %s %s, id ASC
		LIMIT $9 OFFSET $10`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Context with a timeout for the database operation
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	// execute the query
	rows, err := m.DB.QueryContext(ctx, query, fname, lname, email, gender, activated, facilitator, officer, deleted, filters.limit(), filters.offset())
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
			&user.IsActivated,
			&user.IsFacilitator,
			&user.IsOfficer,
			&user.IsDeleted,
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
		SELECT u.id, u.first_name, u.last_name, u.email, u.gender, u.password_hash, u.is_activated, u.is_facilitator, u.is_officer, u.is_deleted, u.created_at, u.updated_at, u.version
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
		&user.IsActivated,
		&user.IsFacilitator,
		&user.IsOfficer,
		&user.IsDeleted,
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
