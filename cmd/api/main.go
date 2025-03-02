package main

import (
	"github.com/labstack/echo/v4"
	"main.go/cmd/handlers"
	"main.go/common/cfg"
	"main.go/common/database"
)

func main() {
	cfg := cfg.NewConfig()
	database.InitDb(cfg.DBDSN)
	e := echo.New()
	handlers.RegisterHandlers(e)

	// Start the server
	e.Logger.Fatal(e.Start(":8001"))
}
