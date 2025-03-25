package utils

import (
	"context"
	"math/rand"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

// ContextKey is a custom type to avoid collisions
type ContextKey string

const UserIDKey ContextKey = "user_id"

const DefaultLimit = 30
const DefaultOffset = 0

// EncryptPassword hashes the given password using bcrypt
func EncryptPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ComparePassword verifies if a given password matches the stored hashed password
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func GenerateReadableID(length int) string {
	rand.Seed(time.Now().UnixNano()) // Ensure randomness on each execution
	id := make([]byte, length)
	for i := range id {
		id[i] = charset[rand.Intn(len(charset))]
	}
	return string(id)
}

func GetCurrentUser(ctx context.Context) string {
	userId, ok := ctx.Value(UserIDKey).(string)
	if ok {
		return userId
	}
	log.Infof("failed to get userid")
	return ""
}
