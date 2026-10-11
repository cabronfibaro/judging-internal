package db

import (
	"database/sql"
	"embed"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func RunMigrations(workspace string) error {
	if err := os.MkdirAll(workspace, 0755); err != nil {
		return err
	}

	migrationFS, err := fs.Sub(migrations, "migrations")
	if err != nil {
		return err
	}

	sourceDriver, err := iofs.New(migrationFS, ".")
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", filepath.Join(workspace, "sqlite.db"))
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		"sqlite3",
		driver,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
