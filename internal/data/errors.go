// FileName: internal/data/errors.go
package data

import (
	"errors"
	"strings"

	"github.com/lib/pq"
)

// Predefined errors for common scenarios
var (
	ErrRecordNotFound      = errors.New("record not found")
	ErrEditConflict        = errors.New("edit conflict")
	ErrDuplicateEmail      = errors.New("duplicate email")
	ErrDuplicateValue      = errors.New("duplicate value")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrNoMatch             = errors.New("no matching records found")
	ErrForeignKeyViolation = errors.New("constraint violation")
)

func isDuplicateKeyViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23505"
	}
	return strings.Contains(strings.ToLower(err.Error()), "duplicate key value")
}

func isForeignKeyViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return pqErr.Code == "23503"
	}
	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "violates foreign key constraint")
}
