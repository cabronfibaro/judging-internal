package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"judginginternal/internal/challonge"
	"judginginternal/internal/models/tournament"
	"math/rand/v2"
	"strings"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type CreateTournamentRequest struct {
	Name                     string   `json:"name"`
	Size                     int      `json:"size"`
	StartingDate             string   `json:"startingDate"`
	Track                    string   `json:"track"`
	Entrants                 []string `json:"entrants"`
	RandomizeQualifyingOrder bool     `json:"randomizeQualifyingOrder"`
	QualifyingOrder          []string `json:"qualifyingOrder"`
	CriteriaConfigID         int      `json:"criteriaConfigID"`
}

type BracketMatch struct {
	MatchID                string   `json:"matchID"`
	ParticipantName        *string  `json:"participantName"`
	Status                 string   `json:"status"`
	Round                  int      `json:"round"`
	WinnerID               *string  `json:"winnerId"`
	ChallongeParticipantID *string  `json:"challongeParticipantID"`
	ParticipantID          *string  `json:"participantID"`
	Score                  *float64 `json:"score"`
	BothRunFinished        bool     `json:"bothRunFinished"`
}

type TournamentRepository interface {
	CreateTournament(context.Context, CreateTournamentRequest) (*tournament.Tournament, error)
	GetTournaments() ([]tournament.Tournament, error)
	GetTournamentById(string) (tournament.Tournament, error)
	DeleteTournamentById(context.Context, string) error
	UpdateTournament(context.Context, string, UpdateTournamentRequest) error
	UpdateTournamentState(context.Context, string, string) error
}

type SQLiteTournamentRepository struct {
	DB              *sql.DB
	ChallongeAPIKey string
}

func (r *SQLiteTournamentRepository) CreateTournament(ctx context.Context, req CreateTournamentRequest) (*tournament.Tournament, error) {
	name := strings.TrimSpace(req.Name)
	track := strings.TrimSpace(req.Track)

	names := make([]string, 0, len(req.Entrants))
	for _, entrant := range req.Entrants {
		if entrant = strings.TrimSpace(entrant); entrant != "" {
			names = append(names, entrant)
		}
	}

	if name == "" || track == "" || len(names) == 0 {
		return nil, fmt.Errorf("name, track, and at least one entrant are required")
	}

	if req.Size <= 0 || req.Size&(req.Size-1) != 0 {
		return nil, fmt.Errorf("size must be a positive power of two")
	}

	startingDate, err := time.Parse(time.RFC3339, req.StartingDate)
	if err != nil {
		startingDate, err = time.Parse("2006-01-02T15:04", req.StartingDate)
		if err != nil {
			return nil, fmt.Errorf("invalid starting date: %w", err)
		}
	}

	challonge, err := challonge.CreateTournament(
		ctx,
		r.ChallongeAPIKey,
		name,
		req.Size,
		startingDate,
	)
	if err != nil {
		return nil, err
	}

	challongeID := challonge.Data.ID
	challongeURL := challonge.Data.Attributes.URL

	// randomize qualifying order if enabled
	if req.RandomizeQualifyingOrder {
		rand.Shuffle(len(names), func(i, j int) {
			names[i], names[j] = names[j], names[i]
		})
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	now := time.Now().Unix()

	trackID, err := gonanoid.New(15)
	if err != nil {
		return nil, err
	}

	// track table
	_, err = tx.ExecContext(ctx, `
INSERT INTO tracks (id, name, created_at, updated_at)
VALUES (?, ?, ?, ?)
ON CONFLICT(name) DO NOTHING
`, trackID, track, now, now)
	if err != nil {
		return nil, err
	}

	err = tx.QueryRowContext(ctx,
		`SELECT id FROM tracks WHERE name = ?`,
		track,
	).Scan(&trackID)
	if err != nil {
		return nil, err
	}

	tournamentID, err := gonanoid.New(15)
	if err != nil {
		return nil, err
	}

	var criteriaConfig string
	err = tx.QueryRowContext(ctx, `
SELECT config from criteria_configs WHERE id = ?
`, req.CriteriaConfigID).Scan(&criteriaConfig)
	if err != nil {
		return nil, err
	}

	// tournaments table
	if _, err := tx.ExecContext(ctx, `
INSERT INTO tournaments (
    id, name, starting_date, track_id, qualifying_order,
    size, randomize_qualifying_order, status, created_at, updated_at, criteria_config
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`,
		tournamentID,
		name,
		startingDate.Unix(),
		trackID,
		"",
		req.Size,
		req.RandomizeQualifyingOrder,
		"qualifying",
		now,
		now,
		criteriaConfig,
	); err != nil {
		return nil, err
	}

	participantIDs := make([]string, 0, len(names))

	for _, participantName := range names {
		participantID, err := gonanoid.New(15)
		if err != nil {
			return nil, err
		}

		if _, err := tx.ExecContext(ctx, `
INSERT INTO participants (
	id, name, tournament_id, created_at, updated_at
) VALUES (?, ?, ?, ?, ?)
		`, participantID, participantName, tournamentID, now, now); err != nil {
			return nil, err
		}

		participantIDs = append(participantIDs, participantID)

		entrantID, err := gonanoid.New(15)
		if err != nil {
			return nil, err
		}

		if _, err := tx.ExecContext(ctx, `
INSERT INTO entrants (
	id, tournament_id, participant_id, status, score, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
		`, entrantID, tournamentID, participantID, "pending", 0, now, now); err != nil {
			return nil, err
		}
	}

	qualifyingOrder := strings.Join(participantIDs, ", ")

	if _, err := tx.ExecContext(ctx, `
UPDATE tournaments
SET qualifying_order = ?, updated_at = ?
WHERE id = ?
	`, qualifyingOrder, now, tournamentID); err != nil {
		return nil, err
	}

	bracketID, err := gonanoid.New(15)
	if err != nil {
		return nil, err
	}

	if _, err := tx.ExecContext(ctx, `
INSERT INTO brackets (
	id, tournament_id, challonge_id, url, status, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?)
	`, bracketID, tournamentID, challongeID, challongeURL, "not_started", now, now); err != nil {
		return nil, err
	}

	// TODO: should remove challonge tournament if transaction fails
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &tournament.Tournament{
		ID:              tournamentID,
		Name:            name,
		StartingDate:    startingDate.Format(time.RFC3339),
		Track:           &track,
		QualifyingOrder: qualifyingOrder,
		Status:          "qualifying",
		URL:             challongeURL,
		Size:            req.Size,
	}, nil
}

func (r *SQLiteTournamentRepository) GetTournaments() ([]tournament.Tournament, error) {
	var tournaments []tournament.Tournament
	rows, err := r.DB.Query(`
SELECT 
	t1.id, t1.name, t1.starting_date, t2.name as track, t1.qualifying_order, t1.size, t1.status, t3.url 
FROM tournaments t1
LEFT JOIN tracks t2
	ON t1.track_id = t2.id
LEFT JOIN brackets t3
	ON t1.id = t3.tournament_id
	`)
	if err != nil {
		return tournaments, err
	}
	defer rows.Close()

	for rows.Next() {
		var t tournament.Tournament

		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.StartingDate,
			&t.Track,
			&t.QualifyingOrder,
			&t.Size,
			&t.Status,
			&t.URL,
		); err != nil {
			return nil, err
		}

		tournaments = append(tournaments, t)
	}

	if err := rows.Err(); err != nil {
		return tournaments, err
	}

	progressRows, err := r.DB.Query(`
SELECT
    b.tournament_id,
    COUNT(DISTINCT n.match_id) AS total_matches,
    COUNT(DISTINCT CASE
        WHEN n.status = 'complete' THEN n.match_id
    END) AS completed_matches
FROM brackets b
JOIN nodes n ON n.bracket_id = b.id
GROUP BY b.tournament_id
`)
	if err != nil {
		return nil, err
	}
	defer progressRows.Close()

	totalMatchesByTournament := make(map[string]int)
	completedMatchesByTournament := make(map[string]int)

	for progressRows.Next() {
		var tournamentID string
		var totalMatches, completedMatches int

		if err := progressRows.Scan(&tournamentID, &totalMatches, &completedMatches); err != nil {
			return nil, err
		}

		totalMatchesByTournament[tournamentID] = totalMatches
		completedMatchesByTournament[tournamentID] = completedMatches
	}

	if err := progressRows.Err(); err != nil {
		return nil, err
	}

	for i := range tournaments {
		tournamentID := tournaments[i].ID

		if total, ok := totalMatchesByTournament[tournamentID]; ok && total > 0 {
			progress := 0.0
			if total > 0 {
				progress = (float64(completedMatchesByTournament[tournamentID]) / float64(total)) * 100
			}

			if progress < 0 {
				progress = 0
			} else if progress > 100 {
				progress = 100
			}

			tournaments[i].Progress = progress
		} else {
			tournaments[i].Progress = 0
		}
	}

	return tournaments, nil
}

func (r *SQLiteTournamentRepository) GetTournamentById(id string) (tournament.Tournament, error) {
	var t tournament.Tournament

	err := r.DB.QueryRow(`
SELECT
	t1.id,
	t1.name,
	t1.starting_date,
	t2.name AS track,
	t1.qualifying_order,
	t1.size,
	t1.status,
	t3.url,
	t1.criteria_config
FROM tournaments t1
LEFT JOIN tracks t2
	ON t1.track_id = t2.id
LEFT JOIN brackets t3
	ON t1.id = t3.tournament_id
WHERE t1.id = ?
	`, id).Scan(
		&t.ID,
		&t.Name,
		&t.StartingDate,
		&t.Track,
		&t.QualifyingOrder,
		&t.Size,
		&t.Status,
		&t.URL,
		&t.CriteriaConfig,
	)
	if err != nil {
		return tournament.Tournament{}, err
	}

	return t, nil
}

func (r *SQLiteTournamentRepository) DeleteTournamentById(
	ctx context.Context,
	id string,
) error {
	var challongeID string

	err := r.DB.QueryRowContext(ctx, `
		SELECT challonge_id
		FROM brackets
		WHERE tournament_id = ?
		LIMIT 1
	`, id).Scan(&challongeID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("bracket not found for this tournament")
		}
		return err
	}

	result, err := r.DB.ExecContext(ctx, `
		DELETE FROM tournaments
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return sql.ErrNoRows
	}

	if err := challonge.DeleteTournament(ctx, r.ChallongeAPIKey, challongeID); err != nil {
		return err
	}

	return nil
}

type UpdateTournamentRequest struct {
	Name            *string   `json:"name"`
	Size            *int      `json:"size"`
	StartingDate    *string   `json:"startingDate"`
	Track           *string   `json:"track"`
	QualifyingOrder *[]string `json:"qualifyingOrder"`
}

func (r *SQLiteTournamentRepository) UpdateTournament(
	ctx context.Context,
	id string,
	payload UpdateTournamentRequest,
) error {
	var set []string
	var args []any

	if payload.Name != nil {
		set = append(set, "name = ?")
		args = append(args, *payload.Name)
	}

	if payload.Size != nil {
		set = append(set, "size = ?")
		args = append(args, *payload.Size)
	}

	if payload.StartingDate != nil {
		set = append(set, "starting_date = ?")
		args = append(args, *payload.StartingDate)
	}

	if payload.Track != nil {
		set = append(set, "track = ?")
		args = append(args, *payload.Track)
	}

	if payload.QualifyingOrder != nil {
		set = append(set, "qualifying_order = ?")
		args = append(args, strings.Join(*payload.QualifyingOrder, ", "))
	}

	if len(set) == 0 {
		return nil
	}

	query := fmt.Sprintf(
		"UPDATE tournaments SET %s WHERE id = ?",
		strings.Join(set, ", "),
	)

	args = append(args, id)

	result, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update tournament: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("tournament %q not found", id)
	}

	return nil
}

func (r *SQLiteTournamentRepository) UpdateTournamentState(ctx context.Context, id, state string) error {
	validStates := map[string]bool{
		"process_checkin":      true,
		"abort_checkin":        true,
		"start_group_stage":    true,
		"finalize_group_stage": true,
		"reset_group_stage":    true,
		"start":                true,
		"finalize":             true,
		"reset":                true,
		"open_predictions":     true,
	}

	if !validStates[state] {
		return fmt.Errorf("state %q is invalid", state)
	}

	var challongeUrl string
	err := r.DB.QueryRowContext(ctx, `
		SELECT url
		FROM brackets
		WHERE tournament_id = ?
		LIMIT 1
	`, id).Scan(&challongeUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("bracket not found")
		}
		return err
	}

	if err := challonge.UpdateTournamentState(
		ctx,
		r.ChallongeAPIKey,
		challongeUrl,
		state,
	); err != nil {
		return err
	}

	if state == "start" {
		result, err := r.DB.ExecContext(ctx, `
			UPDATE brackets
			SET status = 'started'
			WHERE tournament_id = ?
		`, id)
		if err != nil {
			return err
		}

		updated, err := result.RowsAffected()
		if err != nil {
			return err
		}
		if updated == 0 {
			return errors.New("cannot update bracket status to started")
		}
	}

	return nil
}
