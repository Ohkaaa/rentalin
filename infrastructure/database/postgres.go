package database

import (
	"context"
	"fmt"
	"log"
	"rentalin/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPosgres(cfg *config.Config) *pgxpool.Pool {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse db config:", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal("unable to create connection pool:", err)
	}

	if err := dbpool.Ping(ctx); err != nil {
		log.Fatal("database unreachable:", err)
	}

	log.Println("database connected")

	return dbpool
}
