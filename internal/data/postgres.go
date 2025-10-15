package data

import (
	"errors"
	"strings"

	"github.com/lib/pq"
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
