package middleware

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTProvider struct {
	secret []byte
}

func NewJWTProvider(secret string) *JWTProvider {
	return &JWTProvider{
		secret: []byte(secret),
	}
}

func (p *JWTProvider) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(jwt.NewNumericDate(time.Now()).Add(24 * time.Hour)),
	})
	signed, err := token.SignedString(p.secret)
	if err != nil {
		return "", err
	}

	return signed, nil
}

func (p *JWTProvider) ValidateJWS(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return p.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("could not parse token: %w", err)
	}
	return token, nil
}
