package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Open() (*sql.DB, error) {

	connString := "postgres://postgres:postgres@localhost:5432/notesdb"

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	log.Println("Database connection established")
	if err != nil {
		return nil, err
	}

	return db, nil
}
