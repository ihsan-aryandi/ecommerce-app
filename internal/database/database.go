package database

import (
	"database/sql"
	"ecommerce-app/internal/config"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func CreateConnection(cfg *config.Config) *sql.DB {
	dbConfig := cfg.Database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbConfig.Host, dbConfig.Username, dbConfig.Password, dbConfig.Name, dbConfig.Port,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open DB connection: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed to ping DB connection: %v", err)
	}

	return db
}
