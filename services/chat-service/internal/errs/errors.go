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
)
