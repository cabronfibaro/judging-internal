package main

import (
	"embed"
	"judginginternal/internal/config"
	"judginginternal/internal/db"
	"log"
	"os"
	"path"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	workspace, err := initializeWorkspace()
	if err != nil {
		log.Fatalf("could not initialize workspace folder: %v", err)
	}

	cfg, err := config.New(workspace)
	if err != nil {
		log.Fatalf("could not load config file: %v", err)
	}

	if err := db.RunMigrations(workspace); err != nil {
		log.Fatalf("could not run database migrations: %v", err)
	}

	app := NewApp(cfg)

	err = wails.Run(&options.App{
		Title:  "judging-internal",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func initializeWorkspace() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	p := path.Join(dir, ".judging-internal")
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(p, os.ModePerm); err != nil {
			return "", nil
		}
	} else if err != nil {
		return "", nil
	}
	return p, nil
}
