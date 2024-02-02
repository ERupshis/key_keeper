package jwtgenerator

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims struct that keeps standard jwtgenerator claims plus custom UserID.
type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

// JWTGenerator generator itself.
type JwtGenerator struct {
	jwtKey   string
	tokenExp int
}

// NewJWTGenerator creates JWT tokens generator.
func NewJWTGenerator(jwtKey string, tokenExp int) (*JwtGenerator, error) {
	if jwtKey == "" {
		return nil, fmt.Errorf("jwt key is not set")
	}

	return &JwtGenerator{
		jwtKey:   jwtKey,
		tokenExp: tokenExp,
	}, nil
}

// BuildJWTString creates token and returns it as string.
func (j *JwtGenerator) BuildJWTString(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.tokenExp) * time.Hour)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(j.jwtKey))
	if err != nil {
		return "", fmt.Errorf("generate jwt token: %w", err)
	}

	return tokenString, nil
}

// GetUserID gets token in string format, parse it and returns userID.
func (j *JwtGenerator) GetUserID(tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(j.jwtKey), nil
		})
	if err != nil {
		return -1, fmt.Errorf("parse jwt token: %w", err)
	}

	if !token.Valid {
		return -1, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
