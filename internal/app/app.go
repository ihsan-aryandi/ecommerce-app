package app

import (
	"ecommerce-app/internal/config"
	"fmt"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/joho/godotenv"
)

func Run() error {
	// Load .env file if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on actual environment variables")
	}

	// Configuration
	cfg := config.Load()

	// Set Goqu default dialect
	goqu.RegisterDialect("postgres", nil)

	// Router
	router := InitializeApp()

	// Run router
	log.Printf("Server started on port %s\n", cfg.ServerAddr)
	return router.Run(fmt.Sprintf(":%s", cfg.ServerAddr))
}
