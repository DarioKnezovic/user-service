package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Claims represents the JWT claims structure.
type Claims struct {
	UserID uint   `json:"id"`
	Email  string `json:"email"`
	jwt.StandardClaims
}

// GenerateJWT generates a new JWT token.
func GenerateJWT(userID uint, email string, secretKey []byte, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token: %w", err)
	}

	return signedToken, nil
}

func VerifyJWT(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid JWT token")
	}

	// Check if the token has expired
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("JWT token has expired")
	}

	return claims, nil
}
