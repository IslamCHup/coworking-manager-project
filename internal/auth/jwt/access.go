package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func accessSecret() ([]byte, error) {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_ACCESS_SECRET is not set")
	}
	return []byte(secret), nil
}

func GenerateAccessToken(userID uint) (string, error) {
	secret, err := accessSecret()
	if err != nil {
		return "", err
	}

	claims := AccessClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}

func ParseAccessToken(tokenString string) (*AccessClaims, error) {
	fmt.Println("PARSE TOKEN CALLED")
	fmt.Println("TOKEN STRING:", tokenString)
	fmt.Println("ACCESS SECRET (PARSE):", os.Getenv("JWT_ACCESS_SECRET"))

	secret, err := accessSecret()
	if err != nil {
		fmt.Println("PARSE ERROR: secret error:", err)
		return nil, err
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AccessClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)

	if err != nil {
		fmt.Println("PARSE ERROR:", err)
		return nil, err
	}

	claims, ok := token.Claims.(*AccessClaims)
	if !ok || !token.Valid {
		fmt.Println("PARSE ERROR: invalid claims")
		return nil, errors.New("invalid token")
	}

	fmt.Println("PARSE SUCCESS, USER:", claims.UserID)
	return claims, nil
}
