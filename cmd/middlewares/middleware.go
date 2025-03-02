package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"main.go/common/utils"
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

		// Set user ID in context
		ctx := context.WithValue(c.Request().Context(), UserIDKey, userID)
		c.SetRequest(c.Request().WithContext(ctx))

		// Call the next handler
		return next(c)
	}
}
