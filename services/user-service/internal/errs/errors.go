package errs

import "errors"

var (
	ErrInternal  = errors.New("domain error")
	ErrDBFailure = errors.New("database failure")

	ErrMissingMetadata         = errors.New("missing metadata")
	ErrUnexpectedSigningMethod = errors.New("unexpected JWT signing method")

	ErrMissingAuthToken   = errors.New("authorization token is not supplied")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenSigningFailed = errors.New("jwt signing failed")
	ErrTokenNotFound      = errors.New("token not found")

	ErrUserNotFound = errors.New("user not found")

	ErrEmailTaken    = errors.New("email already taken")
	ErrNicknameTaken = errors.New("nickname already taken")

	ErrInvalidPassword = errors.New("invalid password")
	ErrHashingFailed   = errors.New("hashing failed")

	ErrInvalidArgument = errors.New("invalid argument")
)
