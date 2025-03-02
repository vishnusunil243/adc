package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("your-secret-key")

// GenerateJWTToken creates a JWT token with a given userID and expiry time
func GenerateJWTToken(userID string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(), // Token expiration time
		"iat":     time.Now().Unix(),               // Issued at time
	}

	// Create the token with claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// VerifyJWTToken validates a JWT token and extracts claims and returns user_id
func VerifyJwtToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil {
		return "", err
	}

	// Extract claims and get user_id
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, exists := claims["user_id"].(string); exists {
			return userID, nil
		}
		return "", errors.New("user_id not found in token")
	}

	return "", jwt.ErrSignatureInvalid
}
