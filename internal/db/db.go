package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open() (*sql.DB, error) {

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is not set")
	}
	// Open a connection to the PostgreSQL database using the pgx driver
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}
	// Ping the database to ensure the connection is established
	err = db.Ping()
	log.Println("Database connection established")
	if err != nil {
		return nil, err
	}

	return db, nil
}
