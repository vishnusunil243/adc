package handlers

import (
	"github.com/labstack/echo/v4"
	cart_handlers "main.go/cmd/handlers/cart"
	product_handlers "main.go/cmd/handlers/product"
	user_handlers "main.go/cmd/handlers/user"
	"main.go/cmd/middlewares"
)

func RegisterHandlers(e *echo.Echo) {
	registerUserHandlers(e)
	registerProductHandlers(e)
	registerCartHandlers(e)
}

func registerUserHandlers(e *echo.Echo) {
	userGroup := e.Group("/user")
	userGroup.POST("/signup/", user_handlers.NewUserHandler().Signup)
	userGroup.Use(middlewares.JWTMiddleware)
	userGroup.POST("/login/", user_handlers.NewUserHandler().Login)
	userGroup.GET("/list/", user_handlers.NewUserHandler().List)
	userGroup.GET("/:id/", user_handlers.NewUserHandler().Get)
}

func registerProductHandlers(e *echo.Echo) {
	productGroup := e.Group("/product")
	productGroup.Use(middlewares.JWTMiddleware)
	productGroup.POST("/", product_handlers.NewProductHandler().Create)
	productGroup.GET("/:id/", product_handlers.NewProductHandler().Get)
	productGroup.POST("/delete/", product_handlers.NewProductHandler().Delete)
	productGroup.GET("/list/", product_handlers.NewProductHandler().List)
	productGroup.PATCH("/:id/", product_handlers.NewProductHandler().Update)
}

func registerCartHandlers(e *echo.Echo) {
	cartGroup := e.Group("/cart")
	cartGroup.Use(middlewares.JWTMiddleware)
	cartGroup.POST("/", cart_handlers.NewCartHandler().AddToCart)
	cartGroup.PATCH("/:id/", cart_handlers.NewCartHandler().UpdateCart)
	cartGroup.DELETE("/:id/", cart_handlers.NewCartHandler().DeleteCart)
	cartGroup.GET("/:id/", cart_handlers.NewCartHandler().GetCart)
	cartGroup.GET("/list/", cart_handlers.NewCartHandler().ListCart)
}
