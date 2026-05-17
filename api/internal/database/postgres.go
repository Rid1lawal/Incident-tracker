package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresConnection() (*pgxpool.Pool, error) {

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		dbname,
	)

	var pool *pgxpool.Pool
	var err error

	for i := range 10 {

		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}

		config.MaxConns = 10
		config.MinConns = 2
		config.MaxConnLifetime = time.Hour

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err == nil {
			err = pool.Ping(context.Background())
			if err == nil {
				log.Println("database connected")
				return pool, nil
			}
		}

		log.Printf("DB not ready, retrying... attempt %d\n", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("failed to connect to postgres after retries: %w", err)
}
