package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"judginginternal/internal/models/entrant"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

type EntrantsRepository interface {
	GetEntrants(context.Context, string) ([]entrant.Entrant, error)
	UpdateEntrant(ctx context.Context, params UpdateEntrantParams) error
	ConcludeQualifiers(ctx context.Context, tournamentID string) error
}

type SQLiteEntrantsRepository struct {
	DB              *sql.DB
	ChallongeAPIKey string
}

func (r *SQLiteEntrantsRepository) GetEntrants(
	ctx context.Context,
	tournamentID string,
) ([]entrant.Entrant, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT
			e.id,
			e.participant_id,
			p.name,
			e.status,
			e.score
		FROM entrants e
		INNER JOIN participants p ON p.id = e.participant_id
		WHERE e.tournament_id = ?
	`, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("query tournament entrants: %w", err)
	}
	defer rows.Close()

	entrants := make([]entrant.Entrant, 0)

	for rows.Next() {
		var e entrant.Entrant
		var score sql.NullFloat64

		if err := rows.Scan(
			&e.ID,
			&e.ParticipantID,
			&e.Name,
			&e.Status,
			&score,
		); err != nil {
			return nil, fmt.Errorf("scan tournament entrant: %w", err)
		}

		if score.Valid {
			e.Score = &score.Float64
		}

		entrants = append(entrants, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tournament entrants: %w", err)
	}

	return entrants, nil
}

type UpdateEntrantParams struct {
	TournamentID string
	EntrantID    string
	Status       *string
	Score        *float64
}

func (r *SQLiteEntrantsRepository) UpdateEntrant(
	ctx context.Context,
	params UpdateEntrantParams,
) error {
	var updates []string
	var args []any

	if params.Status != nil {
		switch *params.Status {
		case "qualified", "not_qualified", "pending", "scored":
		default:
			return fmt.Errorf("invalid entrant status")
		}

		updates = append(updates, "status = ?")
		args = append(args, *params.Status)
	}

	if params.Score != nil {
		updates = append(updates, "score = ?")
		args = append(args, *params.Score)
	}

	if len(updates) == 0 {
		return fmt.Errorf("nothing to update")
	}

	args = append(args, params.EntrantID, params.TournamentID)

	query := `
		UPDATE entrants
		SET ` + strings.Join(updates, ", ") + `
		WHERE id = ? AND tournament_id = ?
		RETURNING id, tournament_id, participant_id, status, score
	`

	var e entrant.Entrant
	var score sql.NullFloat64

	err := r.DB.QueryRowContext(ctx, query, args...).Scan(
		&e.ID,
		&e.TournamentID,
		&e.ParticipantID,
		&e.Status,
		&score,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("entrant not found")
		}
		return fmt.Errorf("update entrant: %w", err)
	}

	return nil
}

func (r *SQLiteEntrantsRepository) ConcludeQualifiers(
	ctx context.Context,
	tournamentID string,
) error {
	var size int
	var challongeTournamentID sql.NullString

	err := r.DB.QueryRowContext(ctx, `
		SELECT t.size, b.url
		FROM tournaments t
		LEFT JOIN brackets b ON b.tournament_id = t.id
		WHERE t.id = ?
		LIMIT 1
	`, tournamentID).Scan(&size, &challongeTournamentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("tournament not found")
		}
		return fmt.Errorf("get tournament: %w", err)
	}

	// Load entrants and their participant names.
	rows, err := r.DB.QueryContext(ctx, `
		SELECT
			e.id,
			e.participant_id,
			e.score,
			p.name
		FROM entrants e
		LEFT JOIN participants p ON p.id = e.participant_id
		WHERE e.tournament_id = ?
	`, tournamentID)
	if err != nil {
		return fmt.Errorf("get tournament entrants: %w", err)
	}

	type qualifier struct {
		ID            string
		ParticipantID string
		Name          *string
		Score         float64
	}

	entrants := make([]qualifier, 0)

	for rows.Next() {
		var entrant qualifier
		var score sql.NullFloat64
		var name sql.NullString

		if err := rows.Scan(
			&entrant.ID,
			&entrant.ParticipantID,
			&score,
			&name,
		); err != nil {
			rows.Close()
			return fmt.Errorf("scan entrant: %w", err)
		}

		if score.Valid {
			entrant.Score = score.Float64
		}
		if name.Valid {
			entrant.Name = &name.String
		}

		entrants = append(entrants, entrant)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate entrants: %w", err)
	}
	rows.Close()

	sort.SliceStable(entrants, func(i, j int) bool {
		return entrants[i].Score > entrants[j].Score
	})

	if size < 0 {
		size = 0
	}
	if size < len(entrants) {
		entrants = entrants[:size]
	}

	if !challongeTournamentID.Valid || challongeTournamentID.String == "" {
		return fmt.Errorf("Challonge tournament ID not found")
	}

	type challongeParticipantInput struct {
		Name *string `json:"name"`
	}

	payload := struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Participants []challongeParticipantInput `json:"participants"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	payload.Data.Type = "Participants"
	payload.Data.Attributes.Participants = make(
		[]challongeParticipantInput, 0, len(entrants),
	)

	for _, entrant := range entrants {
		payload.Data.Attributes.Participants = append(
			payload.Data.Attributes.Participants,
			challongeParticipantInput{Name: entrant.Name},
		)
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode Challonge request: %w", err)
	}

	endpoint := "https://api.challonge.com/v2.1/tournaments/" +
		url.PathEscape(challongeTournamentID.String) +
		"/participants/bulk_add.json"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create Challonge request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", r.ChallongeAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request Challonge participants: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read Challonge response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Challonge returned HTTP %d: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	var challongeResponse struct {
		Data []struct {
			ID         string `json:"id"`
			Attributes struct {
				Name string `json:"name"`
			} `json:"attributes"`
		} `json:"data"`
	}

	if err := json.Unmarshal(responseBody, &challongeResponse); err != nil {
		return fmt.Errorf("decode Challonge response: %w", err)
	}

	// Match the original name-to-ID mapping behavior.
	challongeIDsByName := make(map[string]string)
	for _, participant := range challongeResponse.Data {
		challongeIDsByName[participant.Attributes.Name] = participant.ID
	}

	now := time.Now().Unix()

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin qualifying transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE entrants
		SET status = 'not_qualified', updated_at = ?
		WHERE tournament_id = ?
	`, now, tournamentID)
	if err != nil {
		return fmt.Errorf("mark entrants not qualified: %w", err)
	}

	for _, entrant := range entrants {
		if entrant.Name == nil {
			return fmt.Errorf("Challonge participant not found: <nil>")
		}

		challongeID, ok := challongeIDsByName[*entrant.Name]
		if !ok || challongeID == "" {
			return fmt.Errorf(
				"Challonge participant not found: %s",
				*entrant.Name,
			)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE entrants
			SET status = 'qualified', updated_at = ?
			WHERE id = ? AND tournament_id = ?
		`, now, entrant.ID, tournamentID)
		if err != nil {
			return fmt.Errorf("mark entrant qualified: %w", err)
		}

		_, err = tx.ExecContext(ctx, `
			UPDATE participants
			SET challonge_participant_id = ?, updated_at = ?
			WHERE id = ? AND tournament_id = ?
		`, challongeID, now, entrant.ParticipantID, tournamentID)
		if err != nil {
			return fmt.Errorf("update Challonge participant ID: %w", err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE tournaments
		SET status = 'elimination', updated_at = ?
		WHERE id = ?
	`, now, tournamentID)
	if err != nil {
		return fmt.Errorf("update tournament status: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit qualifying transaction: %w", err)
	}

	return nil
}
