package handlers

import (
	"github.com/labstack/echo/v4"
	address_handlers "main.go/cmd/handlers/address"
	cart_handlers "main.go/cmd/handlers/cart"
	order_handlers "main.go/cmd/handlers/order"
	product_handlers "main.go/cmd/handlers/product"
	user_handlers "main.go/cmd/handlers/user"
	"main.go/cmd/middlewares"
)

func RegisterHandlers(e *echo.Group) {
	registerUserHandlers(e)
	registerProductHandlers(e)
	registerCartHandlers(e)
	registerOrderHandlers(e)
	registerAddressHandlers(e)
}

func registerUserHandlers(e *echo.Group) {
	userGroup := e.Group("/user")
	userGroup.POST("/login/", user_handlers.NewUserHandler().Login, middlewares.BasicAuthMiddleware)
	userGroup.POST("/signup/", user_handlers.NewUserHandler().Signup, middlewares.BasicAuthMiddleware)
	userGroup.Use(middlewares.JWTMiddleware)
	userGroup.GET("/list/", user_handlers.NewUserHandler().List)
	userGroup.GET("/:id/", user_handlers.NewUserHandler().Get)
}

func registerProductHandlers(e *echo.Group) {
	productGroup := e.Group("/product")
	productGroup.Use(middlewares.JWTMiddleware)
	productGroup.POST("/", product_handlers.NewProductHandler().Create)
	productGroup.GET("/:id/", product_handlers.NewProductHandler().Get)
	productGroup.DELETE("/:id/", product_handlers.NewProductHandler().Delete)
	productGroup.GET("/list/", product_handlers.NewProductHandler().List)
	productGroup.PATCH("/:id/", product_handlers.NewProductHandler().Update)
}

func registerCartHandlers(e *echo.Group) {
	cartGroup := e.Group("/cart")
	cartGroup.Use(middlewares.JWTMiddleware)
	cartGroup.POST("/", cart_handlers.NewCartHandler().AddToCart)
	cartGroup.PATCH("/:id/", cart_handlers.NewCartHandler().UpdateCart)
	cartGroup.DELETE("/:id/", cart_handlers.NewCartHandler().DeleteCart)
	cartGroup.GET("/:id/", cart_handlers.NewCartHandler().GetCart)
	cartGroup.GET("/list/", cart_handlers.NewCartHandler().ListCart)
}

func registerOrderHandlers(e *echo.Group) {
	orderGroup := e.Group("/order")
	orderGroup.Use(middlewares.JWTMiddleware)
	orderGroup.POST("/", order_handlers.NewOrderHandler().AddOrder)
	orderGroup.GET("/list/", order_handlers.NewOrderHandler().ListOrders)
	orderGroup.GET("/:id/", order_handlers.NewOrderHandler().GetOrder)
	orderGroup.PATCH("/:id/", order_handlers.NewOrderHandler().UpdateOrder)
}

func registerAddressHandlers(e *echo.Group) {
	addressGroup := e.Group("/address")
	addressGroup.Use(middlewares.JWTMiddleware)
	addressGroup.POST("/", address_handlers.NewAddressHandler().Create)
	addressGroup.PATCH("/:id/", address_handlers.NewAddressHandler().Update)
	addressGroup.GET("/:id/", address_handlers.NewAddressHandler().Get)
	addressGroup.GET("/list/", address_handlers.NewAddressHandler().List)
	addressGroup.DELETE("/:id/", address_handlers.NewAddressHandler().Delete)
}
