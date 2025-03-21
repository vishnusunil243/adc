package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"main.go/cmd/handlers"
	"main.go/common/cfg"
	"main.go/common/database"
	"main.go/internal/models"
)

func main() {
	cfg := cfg.LoadConfig()
	db := database.InitDb(cfg.DBDSN)
	if err := models.Migrate(db); err != nil {
		log.Fatalf("error migrating tables: %v", err.Error())
	}
	e := echo.New()
	v1 := e.Group("/api/v1")
	handlers.RegisterHandlers(v1)

	// Start the server
	e.Logger.Fatal(e.Start(":8001"))
}
