package utils

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

// TODO: Add documentation and improve naming

func IsUniqueViolation(err error) (violated bool, constraint string) {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) && pqErr.Code == "23505" {
		return true, pqErr.Constraint
	}
	return false, ""
}

func IsUserExists(err error) bool {
	if errors.Is(err, sql.ErrNoRows) {
		return false
	}
	return true
}
