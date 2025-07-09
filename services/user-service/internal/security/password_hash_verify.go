package security

import (
	"log/slog"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		slog.Error("bcrypt error", "err", err)
		return errs.ErrInvalidPassword
	}
	return nil
}
