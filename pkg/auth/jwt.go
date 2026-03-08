package auth

import (
	"crypto/sha256"
	"fmt"
	"os"
	"time"
	"github.com/golang-jwt/jwt/v4"
)

var secretKey = []byte("default-secret-key")

type Claims struct {
	PasswordHash string `json:"password_hash"`
	jwt.RegisteredClaims
}

func GenerateToken(password string) (string, error) {
	passwordHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))

	claims := &Claims{
		PasswordHash: passwordHash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString, password string) bool {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return false
	}

	currentHash := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	return claims.PasswordHash == currentHash
}
