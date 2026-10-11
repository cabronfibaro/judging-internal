package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Client(src string) *sql.DB {
	db, err := sql.Open("sqlite3", src)
	if err != nil {
		panic(err)
	}

	// enable foreign key support for this connection
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		log.Fatalf("failed to enable foreign keys: %v", err)
	}

	return db
}
