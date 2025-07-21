// Package errs defines domain, database, authentication, and validation errors
// for the user-service. It centralizes error values for consistent handling and
// supports the Single Responsibility and Open/Closed principles.
package errs

import "errors"

var (
	// ErrInternal indicates a generic domain error.
	ErrInternal = errors.New("domain error")
	// ErrDBFailure indicates a database operation failure.
	ErrDBFailure = errors.New("database failure")

	// ErrMissingMetadata indicates missing required metadata in a request.
	ErrMissingMetadata = errors.New("missing metadata")
	// ErrUnexpectedSigningMethod indicates an unexpected JWT signing method.
	ErrUnexpectedSigningMethod = errors.New("unexpected JWT signing method")

	// ErrMissingAuthToken indicates that an authorization token was not supplied.
	ErrMissingAuthToken = errors.New("authorization token is not supplied")
	// ErrInvalidToken indicates an invalid JWT token.
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenSigningFailed indicates a JWT signing failure.
	ErrTokenSigningFailed = errors.New("jwt signing failed")
	// ErrTokenNotFound indicates a missing refresh token.
	ErrTokenNotFound = errors.New("token not found")

	// ErrUserNotFound indicates that a user was not found in the database.
	ErrUserNotFound = errors.New("user not found")

	// ErrEmailTaken indicates that the email is already registered.
	ErrEmailTaken = errors.New("email already taken")
	// ErrNicknameTaken indicates that the nickname is already registered.
	ErrNicknameTaken = errors.New("nickname already taken")

	// ErrInvalidPassword indicates an invalid password attempt.
	ErrInvalidPassword = errors.New("invalid password")
	// ErrHashingFailed indicates a password hashing failure.
	ErrHashingFailed = errors.New("hashing failed")

	// ErrInvalidArgument indicates invalid input data (validation failed).
	ErrInvalidArgument = errors.New("invalid argument")
)
