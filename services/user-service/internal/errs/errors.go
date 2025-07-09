package errs

import "errors"

var (
	ErrInternal = errors.New("internal error")

	ErrUserNotFound = errors.New("user not found")

	ErrEmailTaken    = errors.New("email already taken")
	ErrNicknameTaken = errors.New("nickname already taken")

	ErrInvalidPassword    = errors.New("invalid password")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidEmail       = errors.New("invalid email address")

	ErrUsernameIsRequired = errors.New("username is required")
	ErrEmailIsRequired    = errors.New("email is required")
	ErrPasswordIsRequired = errors.New("password is required")
)
