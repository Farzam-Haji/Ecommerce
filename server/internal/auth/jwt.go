package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const TokenExpiresAfter = 72 * time.Hour

var ErrInvalidToken = errors.New("invalid token")
var ErrExpiredToken = errors.New("token has expired")



type CustomClaims struct {
	UserID 		int64
	jwt.RegisteredClaims
}

type JWTManager struct {
	Secret 	[]byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		Secret: []byte(secret),
	}
}

func (m *JWTManager) GenerateToken(userID int64) (string, error) {
	claim := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(TokenExpiresAfter),
			),
			IssuedAt: jwt.NewNumericDate(
				time.Now(),
			),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(m.Secret)
}

func (m *JWTManager) ValidateToken(singedToken string) (*CustomClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(singedToken, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return m.Secret, nil
	})
	if err != nil {
		switch{
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrExpiredToken
		default:
			return nil, fmt.Errorf("%w: %w", ErrInvalidToken, err)
		}
	}

	if claims, ok := parsedToken.Claims.(*CustomClaims); ok && parsedToken.Valid {
		return claims, nil
	}
	
	return nil, ErrInvalidToken
}