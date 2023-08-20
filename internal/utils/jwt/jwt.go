package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/Xacor/gophermart/internal/entity"
	"github.com/golang-jwt/jwt/v5"
)

const tokenExp = time.Hour * 3

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

var (
	ErrEmptySignKey  = errors.New("empty sign key")
	ErrInvalidToken  = errors.New("invalid token")
	ErrInvalidClaims = errors.New("couldn't parse claims")
	ErrTokenExpired  = errors.New("token expired")
)

func BuildToken(user entity.User, key string) (string, error) {
	if len(key) == 0 {
		return "", ErrEmptySignKey
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: user.ID,
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString, key string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(key), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	if !claims.ExpiresAt.After(time.Now()) {
		return nil, ErrTokenExpired
	}

	return claims, nil
}
