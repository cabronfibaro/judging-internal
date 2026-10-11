package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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

type roundGroup struct {
	round        int
	matches      [][]Match
	matchIndexes map[string]int
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
	GetMatches(ctx context.Context, tournamentID string) ([][][]Match, error)
	UpdateMatchSingleParticipant(ctx context.Context, params UpdateMatchSingleParticipantParams) (Node, error)
	GetNodeByID(ctx context.Context, id string) (Node, error)
}

type SQLiteMatchesRepository struct {
	DB              *sql.DB
	ChallongeAPIKey string
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

	groups := make([]roundGroup, 0)
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

			groups = append(groups, roundGroup{
				round:        match.Round,
				matches:      make([][]Match, 0),
				matchIndexes: make(map[string]int),
			})
		}

		group := &groups[roundIndex]

		// preserve the order in which matches first appear within each round.
		matchIndex, exists := group.matchIndexes[match.MatchID]
		if !exists {
			matchIndex = len(group.matches)
			group.matchIndexes[match.MatchID] = matchIndex
			group.matches = append(group.matches, make([]Match, 0))
		}

		group.matches[matchIndex] = append(
			group.matches[matchIndex],
			match,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bracket matches: %w", err)
	}

	groupedMatches := make([][][]Match, 0, len(groups))
	for _, group := range groups {
		groupedMatches = append(groupedMatches, group.matches)
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
			return Node{}, fmt.Errorf("participant not found")
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
