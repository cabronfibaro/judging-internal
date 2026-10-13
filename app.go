package main

import (
	"context"
	"database/sql"
	"fmt"
	"judginginternal/internal/config"
	"judginginternal/internal/db"
	"judginginternal/internal/db/repository"
	"judginginternal/internal/models/entrant"
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
	matchesRepo    repository.MatchesRepository
	bracketRepo    repository.BracketRepository
	entrantsRepo   repository.EntrantsRepository
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

	a.matchesRepo = &repository.SQLiteMatchesRepository{
		DB:              client,
		ChallongeAPIKey: a.cfg.ChallongeApiKey,
	}

	a.bracketRepo = &repository.SQLiteBracketRepository{
		DB:              client,
		ChallongeAPIKey: a.cfg.ChallongeApiKey,
	}

	a.entrantsRepo = &repository.SQLiteEntrantsRepository{
		DB:              client,
		ChallongeAPIKey: a.cfg.ChallongeApiKey,
	}

	runtime.LogInfo(a.ctx, "Database connection established and repositories initialized")
}

func (a *App) GetTournaments() ([]tournament.Tournament, error) {
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

func (a *App) GetTournament(id string) (tournament.Tournament, error) {
	t, err := a.tournamentRepo.GetTournamentById(id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to get tournament: %v", err)
		return tournament.Tournament{}, err
	}

	return t, nil
}

func (a *App) CreateTournament(req repository.CreateTournamentRequest) (tournament.Tournament, error) {
	t, err := a.tournamentRepo.CreateTournament(a.ctx, req)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to get tournaments: %v", err)
		return tournament.Tournament{}, err
	}

	return *t, err
}

func (a *App) UpdateTournament(id string, req repository.UpdateTournamentRequest) error {
	return a.tournamentRepo.UpdateTournament(a.ctx, id, req)
}

func (a *App) UpdateTournamentState(id, state string) error {
	return a.tournamentRepo.UpdateTournamentState(a.ctx, id, state)
}

func (a *App) DeleteTournament(id, status string) error {
	// only update the tournament state if tournament already been started
	if status != "qualifying" {
		if err := a.tournamentRepo.UpdateTournamentState(a.ctx, id, "reset"); err != nil {
			runtime.LogError(a.ctx, err.Error())
			return err
		}
	}

	return a.tournamentRepo.DeleteTournamentById(a.ctx, id)
}

func (a *App) SyncTournament(id string) error {
	return a.bracketRepo.Sync(a.ctx, id)
}

func (a *App) GetMatch(matchID string) (repository.Match, error) {
	match, err := a.matchesRepo.GetMatch(a.ctx, matchID)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return repository.Match{}, err
	}

	return match, nil
}

func (a *App) GetMatches(tournamentID string) ([][][]repository.Match, error) {
	matches, err := a.matchesRepo.GetMatches(a.ctx, tournamentID)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return nil, err
	}
	return matches, nil
}

func (a *App) UpdateMatchSingleParticipant(req repository.UpdateMatchSingleParticipantParams) error {
	_, err := a.matchesRepo.UpdateMatchSingleParticipant(a.ctx, req)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return err
	}
	return nil
}

func (a *App) GetActiveMatches(tournamentID string) ([]repository.ActiveBattle, error) {
	matches, err := a.matchesRepo.GetActiveMatches(a.ctx, tournamentID)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return nil, err
	}
	return matches, nil
}

func (a *App) GetEntrants(id string) ([]entrant.Entrant, error) {
	entrants, err := a.entrantsRepo.GetEntrants(a.ctx, id)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to get tournament: %v", err)
		return nil, err
	}

	return entrants, nil
}

func (a *App) UpdateEntrant(req repository.UpdateEntrantParams) error {
	return a.entrantsRepo.UpdateEntrant(a.ctx, req)
}

func (a *App) ConcludeQualifers(tournamentId string) error {
	err := a.entrantsRepo.ConcludeQualifiers(a.ctx, tournamentId)
	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to conclude qualifers: %v", err)
		return err
	}
	return nil
}

func (a *App) ConcludeBattle(tournamentID, matchID string) error {
	err := a.matchesRepo.ConcludeBattle(
		a.ctx,
		tournamentID,
		matchID,
	)

	if err != nil {
		runtime.LogErrorf(a.ctx, "Failed to conclude battle: %v", err)
		return err
	}

	return nil
}

func (a *App) ResetBattle(tournamentID, matchID string) error {
	if err := a.matchesRepo.ResetBattle(a.ctx, tournamentID, matchID); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to reset battle: %v", err)
		return err
	}

	return nil
}

func (a *App) ChangeLeader(tournamentID, matchID string) error {
	if err := a.matchesRepo.ChangeLeader(a.ctx, tournamentID, matchID); err != nil {
		runtime.LogErrorf(a.ctx, "Failed to change leader: %v", err)
		return err
	}

	return nil
}
