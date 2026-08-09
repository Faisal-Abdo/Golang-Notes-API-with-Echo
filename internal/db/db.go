package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open(connString string) (*sql.DB, error) {
	// Open a connection to the PostgreSQL database using the pgx driver
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	// Ping the database to ensure the connection is established
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database connection established")

	return db, nil
}
