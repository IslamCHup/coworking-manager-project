package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func getAccessSecret() ([]byte, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET is not set")
	}
	return []byte(secret), nil
}

func GenerateAccessToken(userID uint) (string, error) {
	accessSecret, err := getAccessSecret()
	if err != nil {
		return "", err
	}

	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessSecret)
}

func ParseAccessToken(tokenString string) (*AccessClaims, error) {
	if tokenString == "" {
		return nil, errors.New("token is empty")
	}

	accessSecret, err := getAccessSecret()
	if err != nil {
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return accessSecret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
