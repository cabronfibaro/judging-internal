package main

import (
	"context"
	"database/sql"
	"fmt"
	"judginginternal/internal/config"
	"judginginternal/internal/db"
	"judginginternal/internal/db/repository"
	"judginginternal/internal/models/tournament"
	"path"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var client *sql.DB = nil

// App struct
type App struct {
	ctx            context.Context
	cfg            *config.Config
	tournamentRepo repository.TournamentRepository
	bracketRepo    repository.BracketRepository
}

// NewApp creates a new App application struct
func NewApp(cfg *config.Config) *App {
	return &App{
		cfg: cfg,
	}
}

// startup is called when the app starts.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	client = db.Client(path.Join(a.cfg.Workspace, "sqlite.db"))
	if err := client.Ping(); err != nil {
		panic(fmt.Sprintf("could not return connection from database: %v", err))
	}

	a.tournamentRepo = &repository.SQLiteTournamentRepository{
		DB:              client,
		ChallongeAPIKey: a.cfg.ChallongeApiKey,
	}
	a.bracketRepo = &repository.SQLiteBracketRepository{
		DB:              client,
		ChallongeAPIKey: a.cfg.ChallongeApiKey,
	}

	runtime.LogInfo(a.ctx, "Database connection established and repositories initialized")
}

func (a *App) GetTournaments() ([]tournament.Tournament, error) {
	if a.tournamentRepo == nil {
		err := fmt.Errorf("tournament repository is not initialized")
		runtime.LogError(a.ctx, err.Error())
		return []tournament.Tournament{}, err
	}

	tournaments, err := a.tournamentRepo.GetTournaments()
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to get tournaments: %v", err)
		return []tournament.Tournament{}, err
	}

	if tournaments == nil {
		return []tournament.Tournament{}, nil
	}

	return tournaments, nil
}

func (a *App) CreateTournament(req repository.CreateTournamentRequest) (tournament.Tournament, error) {
	if a.tournamentRepo == nil {
		err := fmt.Errorf("tournament repository is not initialized")
		runtime.LogError(a.ctx, err.Error())
		return tournament.Tournament{}, err
	}

	t, err := a.tournamentRepo.CreateTournament(a.ctx, req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to get tournaments: %v", err)
		return tournament.Tournament{}, err
	}

	return *t, err
}

func (a *App) DeleteTournament(id, status string) error {
	if a.tournamentRepo == nil {
		err := fmt.Errorf("tournament repository is not initialized")
		runtime.LogError(a.ctx, err.Error())
		return err
	}

	// only update the tournament state if tournament already been started
	if status != "qualifying" {
		if err := a.tournamentRepo.UpdateTournamentState(a.ctx, id, "reset"); err != nil {
			runtime.LogError(a.ctx, err.Error())
			return err
		}
	}

	return a.tournamentRepo.DeleteTournamentById(a.ctx, id)
}
