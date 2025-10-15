// FileName: internal/data/errors.go
package data

import "errors"

// Predefined errors for common scenarios
var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrEditConflict       = errors.New("edit conflict")
	ErrDuplicateEmail     = errors.New("duplicate email")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrNoMatch            = errors.New("no matching records found")
)
