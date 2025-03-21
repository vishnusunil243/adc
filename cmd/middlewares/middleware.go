package middlewares

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"main.go/common/cfg"
	"main.go/common/utils"
	"main.go/internal/models"
	"main.go/internal/service/user_service"
)

// ContextKey is a custom type to avoid collisions
type ContextKey string

const UserIDKey ContextKey = "user_id"

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		// Extract token (assuming Bearer token format)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.VerifyJwtToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or Expired Token")
		}
		user, serErr := user_service.NewUserService().GetUser(c.Request().Context(), &user_service.GetUserReq{
			Id: userID,
		})
		if serErr != nil {
			return serErr
		}
		if user == nil {
			return fmt.Errorf("user not found")
		}

		// Set user ID in context
		ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		// Extract token (assuming Bearer token format)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := utils.VerifyJwtToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid or Expired Token")
		}
		user, serErr := user_service.NewUserService().GetUser(c.Request().Context(), &user_service.GetUserReq{
			Id: userID,
		})
		if serErr != nil {
			return serErr
		}
		if user == nil {
			return fmt.Errorf("user not found")
		}
		if user.UserType != models.Admin {
			return fmt.Errorf("you don't have authorisation to perform this action")
		}
		// Set user ID in context
		ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}

// BasicAuthMiddleware checks for basic authentication
func BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization Header")
		}

		// Check if the authorization header starts with "Basic "
		if !strings.HasPrefix(authHeader, "Basic ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization Format")
		}

		// Extract the base64 encoded credentials
		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Failed to Decode Credentials")
		}

		// Split the decoded credentials into username and password
		credentials := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(credentials) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Credentials Format")
		}

		username := credentials[0]
		password := credentials[1]

		// Get the expected username and password from environment variables
		expectedUsername := cfg.LoadConfig().AuthUsername
		expectedPassword := cfg.LoadConfig().AuthPassword

		// Check if username and password match
		if username != expectedUsername || password != expectedPassword {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Username or Password")
		}

		// Optional: You can set the username in the context if needed for further use in the handlers
		ctx := context.WithValue(c.Request().Context(), UserIDKey, username)
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}
