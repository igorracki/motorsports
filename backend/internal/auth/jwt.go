package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager interface {
	GenerateToken(userID string, duration time.Duration) (string, time.Time, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type tokenManager struct {
	secret []byte
}

func NewTokenManager(secret string) (TokenManager, error) {
	if secret == "" {
		return nil, fmt.Errorf("JWT secret cannot be empty")
	}
	return &tokenManager{secret: []byte(secret)}, nil
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (m *tokenManager) GenerateToken(userID string, duration time.Duration) (string, time.Time, error) {
	expiresAt := time.Now().Add(duration)
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	return tokenString, expiresAt, err
}

func (m *tokenManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
