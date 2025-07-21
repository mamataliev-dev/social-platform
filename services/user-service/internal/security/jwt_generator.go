package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/mamataliev-dev/social-platform/services/user-service/internal/dto/domain"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/errs"
	"github.com/mamataliev-dev/social-platform/services/user-service/internal/model"
)

const refreshTokenLength = 32

type JWTGenerator struct {
	SecretKey     []byte
	TokenLifetime time.Duration
}

func NewJWTGenerator(secretKey []byte, lifetime time.Duration) *JWTGenerator {
	return &JWTGenerator{SecretKey: secretKey, TokenLifetime: lifetime}
}

func (g *JWTGenerator) CreateTokenPair(user domain.CreateTokenPairInput) (model.TokenPair, error) {
	// Access Token
	claims := jwt.MapClaims{
		"sub":      user.UserID,
		"nickname": user.Nickname,
		"exp":      time.Now().Add(g.TokenLifetime).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(g.SecretKey)
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("%w", errs.ErrTokenSigningFailed)
	}

	// Refresh Token
	refreshToken, err := g.GenerateRefreshToken()
	if err != nil {
		return model.TokenPair{}, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (g *JWTGenerator) GenerateRefreshToken() (string, error) {
	bytes := make([]byte, refreshTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("could not generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
