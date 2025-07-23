// Package errs defines domain, database, authentication, and validation errors
// for the chat-service. It centralizes error values for consistent handling and
// supports the Single Responsibility and Open/Closed principles.
package errs

import "errors"

var (
	// ErrInternal indicates a generic domain error.
	ErrInternal = errors.New("domain error")
	// ErrDBFailure indicates a database operation failure.
	ErrDBFailure = errors.New("database failure")
)
