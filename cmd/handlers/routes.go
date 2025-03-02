package handlers

import (
	"github.com/labstack/echo/v4"
	user_handlers "main.go/cmd/handlers/user"
	"main.go/cmd/middlewares"
)

func RegisterHandlers(e *echo.Echo) {
	registerUserHandlers(e)
}

func registerUserHandlers(e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.POST("/signup/", user_handlers.NewUserHandler().Signup)
	userGroup.Use(middlewares.JWTMiddleware)
	userGroup.POST("/login/", user_handlers.NewUserHandler().Login)
	userGroup.GET("/list/", user_handlers.NewUserHandler().List)
	userGroup.GET("/:id/", user_handlers.NewUserHandler().Get)
}
