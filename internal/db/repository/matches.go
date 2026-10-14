package repository

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Node struct {
	ID                     string   `json:"id"`
	BracketID              string   `json:"bracketId"`
	ParticipantID          *string  `json:"participantId"`
	MatchID                string   `json:"matchId"`
	WinnerID               *string  `json:"winnerId"`
	Status                 string   `json:"status"`
	Round                  int      `json:"round"`
	Score                  *float64 `json:"score"`
	ParticipantName        *string  `json:"participant_name,omitempty"`
	ChallongeParticipantID *string  `json:"challongeParticipantId,omitempty"`
	BothRunFinished        bool     `json:"bothRunFinished,omitempty"`
}

type Match struct {
	MatchID                string   `json:"matchId"`
	ParticipantName        *string  `json:"participant_name"`
	Status                 string   `json:"status"`
	Round                  int      `json:"round"`
	WinnerID               *string  `json:"winnerId"`
	ChallongeParticipantID *string  `json:"challongeParticipantId"`
	ParticipantID          *string  `json:"participantId"`
	Score                  *float64 `json:"score"`
	BothRunFinished        bool     `json:"bothRunFinished"`
}

type RoundGroup struct {
	Round        int
	Matches      [][]Match
	MatchIndexes map[string]int
}

type UpdateMatchSingleParticipantParams struct {
	TournamentID  string
	MatchID       string
	ParticipantID string
	Score         *float64
	WinnerID      *string
	Status        *string
	AppendScore   bool
}

type MatchesRepository interface {
	GetMatch(ctx context.Context, matchID string) (Match, error)
	GetMatches(ctx context.Context, tournamentID string) ([][][]Match, error)
	UpdateMatchSingleParticipant(ctx context.Context, params UpdateMatchSingleParticipantParams) (Node, error)
	GetNodeByID(ctx context.Context, id string) (Node, error)
	GetActiveMatches(ctx context.Context, tournamentID string) ([]ActiveBattle, error)
	ChangeLeader(ctx context.Context, tournamentID, matchID string) error
	ConcludeBattle(ctx context.Context, tournamentID string, matchID string) error
	ResetBattle(ctx context.Context, tournamentID string, matchID string) error
}

type SQLiteMatchesRepository struct {
	DB              *sql.DB
	ChallongeAPIKey string
}

func (r *SQLiteMatchesRepository) GetMatch(
	ctx context.Context,
	matchID string,
) (Match, error) {
	const query = `
		SELECT
			n.match_id,
			p.name,
			n.status,
			n.round,
			n.winner_id,
			p.challonge_participant_id,
			p.id,
			n.score,
			COALESCE(
				(
					SELECT
						COUNT(*) = 2
						AND SUM(
							CASE
								WHEN ab.lead_status = -1 THEN 1
								ELSE 0
							END
						) = 2
					FROM active_battles ab
					WHERE ab.match_id = n.match_id
					GROUP BY ab.match_id
				),
				1
			) AS both_run_finished
		FROM nodes n
		LEFT JOIN participants p ON p.id = n.participant_id
		WHERE n.match_id = ?
		LIMIT 1
	`

	var (
		match           Match
		participantName sql.NullString
		winnerID        sql.NullString
		challongeID     sql.NullString
		participantID   sql.NullString
		score           sql.NullFloat64
		bothRunFinished int
	)

	err := r.DB.QueryRowContext(ctx, query, matchID).Scan(
		&match.MatchID,
		&participantName,
		&match.Status,
		&match.Round,
		&winnerID,
		&challongeID,
		&participantID,
		&score,
		&bothRunFinished,
	)
	if err != nil {
		return Match{}, fmt.Errorf("get match: %w", err)
	}

	if participantName.Valid {
		match.ParticipantName = &participantName.String
	}
	if winnerID.Valid {
		match.WinnerID = &winnerID.String
	}
	if challongeID.Valid {
		match.ChallongeParticipantID = &challongeID.String
	}
	if participantID.Valid {
		match.ParticipantID = &participantID.String
	}
	if score.Valid {
		match.Score = &score.Float64
	}

	match.BothRunFinished = bothRunFinished != 0

	return match, nil
}

func (r *SQLiteMatchesRepository) GetMatches(
	ctx context.Context,
	tournamentID string,
) ([][][]Match, error) {
	const query = `
		SELECT
			n.match_id,
			p.name,
			n.status,
			n.round,
			n.winner_id,
			p.challonge_participant_id,
			p.id,
			n.score,
			COALESCE(
				(
					SELECT
						COUNT(*) = 2
						AND SUM(
							CASE
								WHEN ab.lead_status = -1 THEN 1
								ELSE 0
							END
						) = 2
					FROM active_battles ab
					WHERE ab.match_id = n.match_id
					GROUP BY ab.match_id
				),
				1
			) AS both_run_finished
		FROM nodes n
		INNER JOIN brackets b ON n.bracket_id = b.id
		LEFT JOIN participants p ON p.id = n.participant_id
		WHERE b.tournament_id = ?
	`

	rows, err := r.DB.QueryContext(ctx, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("query bracket matches: %w", err)
	}
	defer rows.Close()

	groups := make([]RoundGroup, 0)
	roundIndexes := make(map[int]int)

	for rows.Next() {
		var (
			match           Match
			participantName sql.NullString
			winnerID        sql.NullString
			challongeID     sql.NullString
			participantID   sql.NullString
			score           sql.NullFloat64
			bothRunFinished int
		)

		err := rows.Scan(
			&match.MatchID,
			&participantName,
			&match.Status,
			&match.Round,
			&winnerID,
			&challongeID,
			&participantID,
			&score,
			&bothRunFinished,
		)
		if err != nil {
			return nil, fmt.Errorf("scan bracket match: %w", err)
		}

		if participantName.Valid {
			match.ParticipantName = &participantName.String
		}
		if winnerID.Valid {
			match.WinnerID = &winnerID.String
		}
		if challongeID.Valid {
			match.ChallongeParticipantID = &challongeID.String
		}
		if participantID.Valid {
			match.ParticipantID = &participantID.String
		}
		if score.Valid {
			match.Score = &score.Float64
		}
		match.BothRunFinished = bothRunFinished != 0

		// preserve the order in which rounds first appear in the query results.
		roundIndex, exists := roundIndexes[match.Round]
		if !exists {
			roundIndex = len(groups)
			roundIndexes[match.Round] = roundIndex

			groups = append(groups, RoundGroup{
				Round:        match.Round,
				Matches:      make([][]Match, 0),
				MatchIndexes: make(map[string]int),
			})
		}

		group := &groups[roundIndex]

		// preserve the order in which matches first appear within each round.
		matchIndex, exists := group.MatchIndexes[match.MatchID]
		if !exists {
			matchIndex = len(group.Matches)
			group.MatchIndexes[match.MatchID] = matchIndex
			group.Matches = append(group.Matches, make([]Match, 0))
		}

		group.Matches[matchIndex] = append(
			group.Matches[matchIndex],
			match,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bracket matches: %w", err)
	}

	groupedMatches := make([][][]Match, 0, len(groups))
	for _, group := range groups {
		groupedMatches = append(groupedMatches, group.Matches)
	}

	return groupedMatches, nil
}

func (r *SQLiteMatchesRepository) UpdateMatchSingleParticipant(
	ctx context.Context,
	params UpdateMatchSingleParticipantParams,
) (Node, error) {
	var node Node

	err := r.DB.QueryRowContext(ctx, `
		SELECT n.id, n.score
		FROM nodes n
		INNER JOIN brackets b ON b.id = n.bracket_id
		WHERE b.tournament_id = ?
			AND n.match_id = ?
			AND n.participant_id = ?
		LIMIT 1
	`, params.TournamentID, params.MatchID, params.ParticipantID).Scan(
		&node.ID,
		&node.Score,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Node{}, fmt.Errorf("participant not found with id %s", params.ParticipantID)
		}
		return Node{}, fmt.Errorf("find participant node: %w", err)
	}

	var updates []string
	var args []any

	if params.Score != nil {
		updatedScore := *params.Score

		if params.AppendScore && node.Score != nil {
			updatedScore += *node.Score
		}

		updates = append(updates, "score = ?")
		args = append(args, updatedScore)
	}

	if params.WinnerID != nil {
		updates = append(updates, "winner_id = ?")
		args = append(args, *params.WinnerID)
	}

	if params.Status != nil {
		updates = append(updates, "status = ?")
		args = append(args, *params.Status)
	}

	if len(updates) > 0 {
		args = append(args, node.ID)

		query := "UPDATE nodes SET " +
			strings.Join(updates, ", ") +
			" WHERE id = ?"

		if _, err := r.DB.ExecContext(ctx, query, args...); err != nil {
			return Node{}, fmt.Errorf("update participant node: %w", err)
		}
	}

	return r.GetNodeByID(ctx, node.ID)
}

func (r *SQLiteMatchesRepository) GetNodeByID(
	ctx context.Context,
	id string,
) (Node, error) {
	var node Node
	var participantID, winnerID, participantName, challongeParticipantID sql.NullString
	var score sql.NullFloat64

	err := r.DB.QueryRowContext(ctx, `
		SELECT
			n.id,
			n.bracket_id,
			n.participant_id,
			n.match_id,
			n.winner_id,
			n.status,
			n.round,
			n.score,
			p.name,
			p.challonge_participant_id
		FROM nodes n
		LEFT JOIN participants p ON p.id = n.participant_id
		WHERE n.id = ?
	`, id).Scan(
		&node.ID,
		&node.BracketID,
		&participantID,
		&node.MatchID,
		&winnerID,
		&node.Status,
		&node.Round,
		&score,
		&participantName,
		&challongeParticipantID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Node{}, fmt.Errorf("node not found: %w", err)
		}
		return Node{}, fmt.Errorf("get node by ID: %w", err)
	}

	if participantID.Valid {
		node.ParticipantID = &participantID.String
	}
	if winnerID.Valid {
		node.WinnerID = &winnerID.String
	}
	if score.Valid {
		node.Score = &score.Float64
	}
	if participantName.Valid {
		node.ParticipantName = &participantName.String
	}
	if challongeParticipantID.Valid {
		node.ChallongeParticipantID = &challongeParticipantID.String
	}

	return node, nil
}

func (r *SQLiteMatchesRepository) GetActiveMatches(
	ctx context.Context,
	tournamentID string,
) ([]ActiveBattle, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT
			ab.id,
			ab.match_id,
			ab.participant_id,
			ab.lead_status
		FROM active_battles ab
		INNER JOIN nodes n
			ON n.match_id = ab.match_id
			AND n.participant_id = ab.participant_id
		INNER JOIN brackets b ON b.id = n.bracket_id
		WHERE b.tournament_id = ?
	`, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("query active matches: %w", err)
	}
	defer rows.Close()

	battles := make([]ActiveBattle, 0)

	for rows.Next() {
		var battle ActiveBattle

		if err := rows.Scan(
			&battle.ID,
			&battle.MatchID,
			&battle.ParticipantID,
			&battle.LeadStatus,
		); err != nil {
			return nil, fmt.Errorf("scan active battle: %w", err)
		}

		battles = append(battles, battle)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate active matches: %w", err)
	}

	return battles, nil
}

func (r *SQLiteMatchesRepository) ChangeLeader(
	ctx context.Context,
	tournamentID string,
	matchID string,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin change leader transaction: %w", err)
	}
	defer tx.Rollback()

	const query = `
		SELECT
			ab.id,
			ab.match_id,
			ab.participant_id,
			ab.lead_status
		FROM active_battles ab
		INNER JOIN nodes n
			ON n.match_id = ab.match_id
			AND n.participant_id = ab.participant_id
		INNER JOIN brackets b
			ON b.id = n.bracket_id
		WHERE b.tournament_id = ?
			AND ab.match_id = ?
	`

	rows, err := tx.QueryContext(ctx, query, tournamentID, matchID)
	if err != nil {
		return fmt.Errorf("query active battles: %w", err)
	}

	var battles []ActiveBattle

	for rows.Next() {
		var battle ActiveBattle

		if err := rows.Scan(
			&battle.ID,
			&battle.MatchID,
			&battle.ParticipantID,
			&battle.LeadStatus,
		); err != nil {
			rows.Close()
			return fmt.Errorf("scan active battle: %w", err)
		}

		battles = append(battles, battle)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate active battles: %w", err)
	}
	rows.Close()

	if len(battles) != 2 {
		return fmt.Errorf(
			"active battle must contain exactly two participants, got %d",
			len(battles),
		)
	}

	var leader, waiting, finished *ActiveBattle

	for i := range battles {
		switch battles[i].LeadStatus {
		case 1:
			leader = &battles[i]
		case 0:
			waiting = &battles[i]
		case -1:
			finished = &battles[i]
		}
	}

	switch {
	case leader != nil && waiting != nil:
		if _, err := tx.ExecContext(
			ctx,
			`UPDATE active_battles SET lead_status = -1 WHERE id = ?`,
			leader.ID,
		); err != nil {
			return fmt.Errorf("finish current leader: %w", err)
		}

		if _, err := tx.ExecContext(
			ctx,
			`UPDATE active_battles SET lead_status = 1 WHERE id = ?`,
			waiting.ID,
		); err != nil {
			return fmt.Errorf("change leader: %w", err)
		}

	case leader != nil && finished != nil:
		if _, err := tx.ExecContext(
			ctx,
			`UPDATE active_battles SET lead_status = -1 WHERE id = ?`,
			leader.ID,
		); err != nil {
			return fmt.Errorf("finish current leader: %w", err)
		}

	case leader == nil && waiting == nil:
		// both participants are already finished. Nothing to change.

	default:
		return fmt.Errorf(
			"invalid active battle lead statuses: %d, %d",
			battles[0].LeadStatus,
			battles[1].LeadStatus,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit change leader transaction: %w", err)
	}

	for i := range battles {
		switch battles[i].LeadStatus {
		case 1:
			battles[i].LeadStatus = -1
		case 0:
			battles[i].LeadStatus = 1
		}
	}

	return nil
}

func (r *SQLiteMatchesRepository) ConcludeBattle(
	ctx context.Context,
	tournamentID string,
	matchID string,
) error {
	const matchQuery = `
		SELECT
			n.participant_id,
			n.score,
			p.challonge_participant_id
		FROM nodes n
		INNER JOIN participants p ON p.id = n.participant_id
		WHERE n.match_id = ?
			AND p.tournament_id = ?
	`

	rows, err := r.DB.QueryContext(ctx, matchQuery, matchID, tournamentID)
	if err != nil {
		return fmt.Errorf("query match participants: %w", err)
	}

	type matchParticipant struct {
		participantID          string
		score                  float64
		challongeParticipantID string
	}

	var participants []matchParticipant

	for rows.Next() {
		var (
			p           matchParticipant
			score       sql.NullFloat64
			challongeID sql.NullString
		)

		if err := rows.Scan(
			&p.participantID,
			&score,
			&challongeID,
		); err != nil {
			rows.Close()
			return fmt.Errorf("scan match participant: %w", err)
		}

		if score.Valid {
			p.score = score.Float64
		}

		if challongeID.Valid {
			p.challongeParticipantID = challongeID.String
		}

		participants = append(participants, p)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate match participants: %w", err)
	}
	rows.Close()

	if len(participants) != 2 {
		return fmt.Errorf(
			"match must contain exactly two participants, got %d",
			len(participants),
		)
	}

	p1, p2 := participants[0], participants[1]

	if p1.challongeParticipantID == "" ||
		p2.challongeParticipantID == "" {
		return fmt.Errorf(
			"one or more participants do not have a Challonge participant ID",
		)
	}

	if p1.score == p2.score {
		return fmt.Errorf(
			"match cannot be concluded because the scores are tied: %v, %v",
			p1.score,
			p2.score,
		)
	}

	winner, loser := p1, p2
	if p2.score > p1.score {
		winner, loser = p2, p1
	}

	const bracketQuery = `
		SELECT url
		FROM brackets
		WHERE tournament_id = ?
		LIMIT 1
	`

	var bracketURL string
	err = r.DB.QueryRowContext(ctx, bracketQuery, tournamentID).Scan(&bracketURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("bracket not found")
		}
		return fmt.Errorf("query bracket: %w", err)
	}

	requestBody := struct {
		Data struct {
			Type       string `json:"type"`
			Attributes struct {
				Match []struct {
					ParticipantID string `json:"participant_id"`
					ScoreSet      string `json:"score_set"`
					Rank          int    `json:"rank"`
					Advancing     *bool  `json:"advancing,omitempty"`
				} `json:"match"`
				Tie bool `json:"tie"`
			} `json:"attributes"`
		} `json:"data"`
	}{}

	requestBody.Data.Type = "match"
	requestBody.Data.Attributes.Tie = false

	advancing := true

	requestBody.Data.Attributes.Match = []struct {
		ParticipantID string `json:"participant_id"`
		ScoreSet      string `json:"score_set"`
		Rank          int    `json:"rank"`
		Advancing     *bool  `json:"advancing,omitempty"`
	}{
		{
			ParticipantID: winner.challongeParticipantID,
			ScoreSet:      strconv.FormatFloat(winner.score, 'f', -1, 64),
			Rank:          1,
			Advancing:     &advancing,
		},
		{
			ParticipantID: loser.challongeParticipantID,
			ScoreSet:      strconv.FormatFloat(loser.score, 'f', -1, 64),
			Rank:          2,
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal Challonge request: %w", err)
	}

	endpoint := fmt.Sprintf(
		"https://api.challonge.com/v2.1/tournaments/%s/matches/%s.json",
		url.PathEscape(bracketURL),
		url.PathEscape(matchID),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		endpoint,
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create Challonge request: %w", err)
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", r.ChallongeAPIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("update Challonge match: %w", err)
	}
	defer resp.Body.Close()

	challongeResponse, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read Challonge response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"Challonge returned status %d: %s",
			resp.StatusCode,
			string(challongeResponse),
		)
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin conclude battle transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`UPDATE nodes SET winner_id = ? WHERE match_id = ?`,
		winner.challongeParticipantID,
		matchID,
	)
	if err != nil {
		return fmt.Errorf("update local match winner: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`DELETE FROM active_battles WHERE match_id = ?`,
		matchID,
	)
	if err != nil {
		return fmt.Errorf("delete active battles: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit conclude battle: %w", err)
	}

	return nil
}

func (r *SQLiteMatchesRepository) ResetBattle(
	ctx context.Context,
	tournamentID string,
	matchID string,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin reset battle transaction: %w", err)
	}
	defer tx.Rollback()

	const bracketQuery = `
		SELECT id
		FROM brackets
		WHERE tournament_id = ?
		LIMIT 1
	`

	var bracketID string
	if err := tx.QueryRowContext(ctx, bracketQuery, tournamentID).Scan(&bracketID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("bracket not found")
		}
		return fmt.Errorf("query bracket: %w", err)
	}

	const battlesQuery = `
		SELECT ab.id, ab.participant_id, ab.rerun
		FROM active_battles ab
		INNER JOIN nodes n
			ON n.match_id = ab.match_id
			AND n.bracket_id = ?
			AND n.participant_id = ab.participant_id
		WHERE ab.match_id = ?
	`

	rows, err := tx.QueryContext(ctx, battlesQuery, bracketID, matchID)
	if err != nil {
		return fmt.Errorf("query active battles: %w", err)
	}

	type battle struct {
		id            string
		participantID string
		rerun         int
	}

	var battles []battle
	for rows.Next() {
		var b battle
		if err := rows.Scan(&b.id, &b.participantID, &b.rerun); err != nil {
			rows.Close()
			return fmt.Errorf("scan active battle: %w", err)
		}
		battles = append(battles, b)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("iterate active battles: %w", err)
	}
	rows.Close()

	if len(battles) != 2 {
		return fmt.Errorf(
			"active battle must contain exactly two participants, got %d",
			len(battles),
		)
	}

	type entrantScore struct {
		score sql.NullFloat64
	}

	getScore := func(participantID string) (float64, error) {
		const query = `
			SELECT score
			FROM entrants
			WHERE tournament_id = ?
				AND participant_id = ?
			LIMIT 1
		`

		var score sql.NullFloat64
		err := tx.QueryRowContext(
			ctx,
			query,
			tournamentID,
			participantID,
		).Scan(&score)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return 0, fmt.Errorf("entrant data not found for participant %q", participantID)
			}
			return 0, fmt.Errorf("query entrant score: %w", err)
		}

		if !score.Valid {
			return 0, nil
		}

		return score.Float64, nil
	}

	p1Score, err := getScore(battles[0].participantID)
	if err != nil {
		return err
	}

	p2Score, err := getScore(battles[1].participantID)
	if err != nil {
		return err
	}

	rerun := battles[0].rerun + 1

	higherScoreLeads := rerun%2 == 0
	p1Leads := higherScoreLeads

	if p1Score != p2Score {
		if higherScoreLeads {
			p1Leads = p1Score > p2Score
		} else {
			p1Leads = p1Score < p2Score
		}
	}

	const resetScoresQuery = `
		UPDATE nodes
		SET score = 0
		WHERE bracket_id = ?
			AND match_id = ?
			AND participant_id IN (?, ?)
	`

	if _, err := tx.ExecContext(
		ctx,
		resetScoresQuery,
		bracketID,
		matchID,
		battles[0].participantID,
		battles[1].participantID,
	); err != nil {
		return fmt.Errorf("reset match scores: %w", err)
	}

	const updateBattleQuery = `
		UPDATE active_battles
		SET rerun = ?, lead_status = ?
		WHERE id = ?
	`

	p1LeadStatus, p2LeadStatus := 1, 0
	if !p1Leads {
		p1LeadStatus, p2LeadStatus = 0, 1
	}

	if _, err := tx.ExecContext(
		ctx,
		updateBattleQuery,
		rerun,
		p1LeadStatus,
		battles[0].id,
	); err != nil {
		return fmt.Errorf("update first active battle: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		updateBattleQuery,
		rerun,
		p2LeadStatus,
		battles[1].id,
	); err != nil {
		return fmt.Errorf("update second active battle: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit reset battle: %w", err)
	}

	return nil
}
