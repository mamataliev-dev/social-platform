package errs

import "errors"

var (
	ErrInternal  = errors.New("internal error")
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

	ErrInvalidPassword     = errors.New("invalid password")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrMissingRequiredData = errors.New("missing or empty required fields")

	ErrHashingFailed = errors.New("hashing failed")

	ErrMissingOrInvalidRefreshToken = errors.New("missing or invalid refresh token")
	ErrMissingOrInvalidAccessToken  = errors.New("missing or invalid access token")
	ErrMissingJWTSecret             = errors.New("missing secret")
)
