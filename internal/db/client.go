package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Client(src string) *sql.DB {
	separator := "?"
	if strings.Contains(src, "?") {
		separator = "&"
	}

	dsn := src + separator + "_foreign_keys=on"

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
