package security

import (
	"golang.org/x/crypto/bcrypt"
	"log/slog"

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
