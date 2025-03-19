package utils

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken creates a JWT with the given claims using the SECRET_KEY.
func GenerateToken(claims jwt.MapClaims) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "mysecret"
	}
	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenJWT.SignedString([]byte(secretKey))
}
