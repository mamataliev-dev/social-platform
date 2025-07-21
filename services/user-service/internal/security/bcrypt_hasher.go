// Package security provides bcrypt-based password hashing and verification for
// the user-service. It supports Dependency Inversion and Single Responsibility principles.
package security

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
)

type BcryptHasher struct{}

func (BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (BcryptHasher) VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		slog.Error("bcrypt error", "err", err)
		return errs.ErrInvalidPassword
	}
	return nil
}
