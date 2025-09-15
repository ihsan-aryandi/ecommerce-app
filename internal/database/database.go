package database

import (
	"context"
	"database/sql"
	"ecommerce-app/internal/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func CreateConnection(cfg *config.Config) *goqu.Database {
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

	log.Println("Database connected")
	return goqu.New("postgres", db)
}

func CreateRDBConnection(cfg *config.Config) *redis.Client {
	redisCfg := cfg.Redis

	db, _ := strconv.Atoi(redisCfg.DB)
	poolSize, _ := strconv.Atoi(redisCfg.PoolSize)
	minIdleConnections, _ := strconv.Atoi(redisCfg.MinIdleConnections)

	addr := fmt.Sprintf("%s:%s", redisCfg.Host, redisCfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     redisCfg.Password,
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConnections,
	})

	// Ping
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to Redis: %s", err)
	}

	log.Println("Redis connected")
	return rdb
}
