package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func Client(src string) *sql.DB {
	db, err := sql.Open("sqlite3", src)
	if err != nil {
		panic(err)
	}

	return db
}
