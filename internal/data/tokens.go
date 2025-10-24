// FileName: internal/data/tokens.go
package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"time"

	"github.com/Pedro-J-Kukul/police_training/internal/validator"
)

/*************************************************************************************************************/
// Declarations
/*************************************************************************************************************/
const (
	ScopeActivation     = "activation"     // Scope for account activation tokens
	ScopeAuthentication = "authentication" // Scope for authentication tokens
	ScopePasswordReset  = "password_reset" // Scope for password reset tokens
)

// Define our token
type Token struct {
	Plaintext string    `json:"token"`  // Plaintext token value
	Hash      []byte    `json:"-"`      // Hashed token value (not exposed in JSON)
	UserID    int64     `json:"-"`      // ID of the user the token belongs to
	Expiry    time.Time `json:"expiry"` // Expiry time of the token
	Scope     string    `json:"-"`      // Scope of the token (not exposed in JSON)
}

// TokenModel struct wraps a database connection pool
type TokenModel struct {
	DB *sql.DB // Database connection pool
}

/*************************************************************************************************************/
// helper functions
/*************************************************************************************************************/

// generateToken generates a new token for a user with a specific scope and expiry duration
func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Create a new token with the provided userID, scope, and calculated expiry time
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Generate a random plaintext token
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, err // Return error if random byte generation fails
	}
	token.Plaintext = base64.RawURLEncoding.EncodeToString(randomBytes) // Encode to URL-safe base64
	hash := sha256.Sum256([]byte(token.Plaintext))                      // Hash the plaintext token

	token.Hash = hash[:] // Set the hash field to the hashed value
	return token, nil    // Return the generated token
}

// ValidateTokenPlaintext checks if a token is valid based on its plaintext value
func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")                // Token must be provided
	v.Check(len(tokenPlaintext) == 22, "token", "must be 22 characters long") // Token must match expected length
}

/*************************************************************************************************************/
// Methods
/*************************************************************************************************************/

// New creates a new token for a user and stores it in the database
func (m *TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Generate a new token
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err // Return error if token generation fails
	}

	err = m.Insert(token) // Insert the token into the database
	return token, err     // Return the token and any insertion error
}

// Insert adds a new token to the database
func (m *TokenModel) Insert(token *Token) error {
	// SQL query to insert a new token into the tokens table
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	// Prepare the arguments for the query
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	_, err := m.DB.ExecContext(ctx, query, args...) // Execute the insert query
	return err                                      // Return any error that occurred during execution
}

// DeleteAllForUser removes all tokens for a specific user and scope from the database
func (m *TokenModel) DeleteAllForUser(scope string, userID int64) error {
	// SQL query to delete tokens for a specific user and scope
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // Create a context with a 3-second timeout
	defer cancel()                                                          // Ensure the context is cancelled to free resources

	_, err := m.DB.ExecContext(ctx, query, scope, userID) // Execute the delete query
	return err                                            // Return any error that occurred during execution
}
