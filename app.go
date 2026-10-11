package main

import (
	"context"
	"database/sql"
	"fmt"
	"myproject/internal/config"
	"myproject/internal/db"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx    context.Context
	cfg    *config.Config
	client *sql.DB
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	app := &App{cfg: cfg}
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	client := db.Client(path.Join(a.cfg.Workspace, "sqlite.db"))
	if err := client.Ping(); err != nil {
		panic(fmt.Sprintf("could not return connection from database: %v", err))
	}

	a.client = client

	runtime.LogInfo(a.ctx, "Database connection established")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
