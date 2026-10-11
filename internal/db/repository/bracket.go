package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"judginginternal/internal/challonge"
	"math/rand"
	"strconv"
	"strings"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

type BracketRepository interface {
	GetBracket(ctx context.Context, tournamentID string) (Bracket, error)
	Sync(ctx context.Context, tournamentID string) (SyncResult, error)
}

type SQLiteBracketRepository struct {
	DB              *sql.DB
	ChallongeAPIKey string
}

type Bracket struct {
	ID                    string `json:"id"`
	TournamentID          string `json:"tournamentId"`
	ChallongeTournamentID string `json:"challongeTournamentId"`
	Status                string `json:"status"`
}

type SyncResult struct {
	Matches       [][][]BracketMatch `json:"matches"`
	ActiveBattles []ActiveBattle     `json:"activeBattles"`
}

type ActiveBattle struct {
	ID            string `json:"id"`
	MatchID       string `json:"matchId"`
	ParticipantID string `json:"participantId"`
	LeadStatus    int    `json:"leadStatus"`
}

type existingNode struct {
	ID      string
	MatchID string
}

func (r *SQLiteBracketRepository) GetBracket(
	ctx context.Context,
	tournamentID string,
) (Bracket, error) {
	var b Bracket

	err := r.DB.QueryRowContext(ctx, `
		SELECT id, tournament_id, url, status
		FROM brackets
		WHERE tournament_id = ?
		LIMIT 1
	`, tournamentID).Scan(
		&b.ID,
		&b.TournamentID,
		&b.ChallongeTournamentID,
		&b.Status,
	)
	if err != nil {
		return Bracket{}, err
	}

	return b, nil
}

func (r *SQLiteBracketRepository) Sync(
	ctx context.Context,
	tournamentID string,
) (SyncResult, error) {
	bracket, err := r.GetBracket(ctx, tournamentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return SyncResult{}, errors.New("bracket not found")
		}
		return SyncResult{}, err
	}

	matches, err := challonge.GetMatches(
		ctx,
		bracket.ChallongeTournamentID,
		r.ChallongeAPIKey,
	)
	if err != nil {
		return SyncResult{}, err
	}

	if err := r.syncNodes(ctx, bracket, tournamentID, matches); err != nil {
		return SyncResult{}, err
	}

	groupedMatches, err := r.getBracketMatches(ctx, tournamentID)
	if err != nil {
		return SyncResult{}, err
	}

	if err := r.createActiveBattles(ctx, tournamentID, groupedMatches); err != nil {
		return SyncResult{}, err
	}

	battles, err := r.getActiveBattles(ctx, tournamentID)
	if err != nil {
		return SyncResult{}, err
	}

	return SyncResult{
		Matches:       groupedMatches,
		ActiveBattles: battles,
	}, nil
}

func (r *SQLiteBracketRepository) syncNodes(
	ctx context.Context,
	bracket Bracket,
	tournamentID string,
	matches []challonge.MatchData,
) error {
	participantIDs := make(map[string]struct{})

	for _, match := range matches {
		for _, entry := range match.Attributes.PointsByParticipant {
			participantIDs[strconv.Itoa(entry.ParticipantID)] = struct{}{}
		}
	}

	participantMap := make(map[string]string)

	if len(participantIDs) > 0 {
		args := []any{tournamentID}
		placeholders := make([]string, 0, len(participantIDs))

		for id := range participantIDs {
			placeholders = append(placeholders, "?")
			args = append(args, id)
		}

		query := `
			SELECT id, challonge_participant_id
			FROM participants
			WHERE tournament_id = ?
			  AND challonge_participant_id IN (` +
			strings.Join(placeholders, ",") + `)`

		rows, err := r.DB.QueryContext(ctx, query, args...)
		if err != nil {
			return err
		}

		for rows.Next() {
			var localID string
			var challongeID sql.NullString

			if err := rows.Scan(&localID, &challongeID); err != nil {
				rows.Close()
				return err
			}

			if challongeID.Valid {
				participantMap[challongeID.String] = localID
			}
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return err
		}

		if err := rows.Close(); err != nil {
			return err
		}
	}

	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, match_id
		FROM nodes
		WHERE bracket_id = ?
	`, bracket.ID)
	if err != nil {
		return err
	}

	existingNodes := make(map[string][]existingNode)

	for rows.Next() {
		var n existingNode

		if err := rows.Scan(&n.ID, &n.MatchID); err != nil {
			rows.Close()
			return err
		}

		existingNodes[n.MatchID] = append(existingNodes[n.MatchID], n)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}

	if err := rows.Close(); err != nil {
		return err
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, match := range matches {
		attrs := match.Attributes
		matchID := match.ID

		matchNodes := existingNodes[matchID]

		var winnerID any
		if attrs.WinnerID != 0 {
			winnerID = strconv.Itoa(attrs.WinnerID)
		}

		for index, entry := range attrs.PointsByParticipant {
			var participantID any

			if entry.ParticipantID != 0 {
				localID, ok := participantMap[strconv.Itoa(entry.ParticipantID)]
				if !ok {
					return fmt.Errorf(
						"local participant not found for Challonge participant %d",
						entry.ParticipantID,
					)
				}
				participantID = localID
			}

			if index < len(matchNodes) {
				// update the existing node at the same index.
				// leave its score unchanged, matching the original.
				_, err := tx.ExecContext(ctx, `
					UPDATE nodes
					SET participant_id = ?,
					    status = ?,
					    round = ?,
					    winner_id = ?
					WHERE id = ?
				`,
					participantID,
					attrs.State,
					attrs.Round,
					winnerID,
					matchNodes[index].ID,
				)
				if err != nil {
					return err
				}
			} else {
				// insert a node if there was no existing node at this index.
				nodeID, err := gonanoid.New(15)
				if err != nil {
					return err
				}

				_, err = tx.ExecContext(ctx, `
					INSERT INTO nodes (
						id,
						bracket_id,
						participant_id,
						match_id,
						status,
						round,
						winner_id
					)
					VALUES (?, ?, ?, ?, ?, ?, ?)
				`,
					nodeID,
					bracket.ID,
					participantID,
					matchID,
					attrs.State,
					attrs.Round,
					winnerID,
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return tx.Commit()
}

func (r *SQLiteBracketRepository) createActiveBattles(
	ctx context.Context,
	tournamentID string,
	groupedMatches [][][]BracketMatch,
) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, round := range groupedMatches {
		for _, match := range round {
			if len(match) != 2 {
				// matches without exactly two nodes are skipped.
				continue
			}

			p1 := match[0]
			p2 := match[1]

			if p1.ParticipantID == nil || p2.ParticipantID == nil {
				continue
			}

			if p1.ChallongeParticipantID == nil ||
				p2.ChallongeParticipantID == nil {
				continue
			}

			if p1.Status == "complete" || p2.Status == "complete" {
				continue
			}

			var p1Score, p2Score float64

			err := tx.QueryRowContext(ctx, `
				SELECT score
				FROM entrants
				WHERE tournament_id = ?
				  AND participant_id = ?
				LIMIT 1
			`, tournamentID, *p1.ParticipantID).Scan(&p1Score)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				return err
			}

			err = tx.QueryRowContext(ctx, `
				SELECT score
				FROM entrants
				WHERE tournament_id = ?
				  AND participant_id = ?
				LIMIT 1
			`, tournamentID, *p2.ParticipantID).Scan(&p2Score)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				return err
			}

			var existingCount int

			err = tx.QueryRowContext(ctx, `
				SELECT COUNT(*)
				FROM active_battles
				WHERE match_id = ?
			`, p1.MatchID).Scan(&existingCount)
			if err != nil {
				return err
			}

			if existingCount >= 2 {
				continue
			}

			var p1Leads bool

			if p1Score == p2Score {
				// randomly assign the lead when qualifying scores are tied.
				p1Leads = rand.Intn(2) == 0
			} else {
				p1Leads = p1Score > p2Score
			}

			p2Leads := !p1Leads

			p1BattleID, err := gonanoid.New(15)
			if err != nil {
				return err
			}

			p2BattleID, err := gonanoid.New(15)
			if err != nil {
				return err
			}

			// intentionally insert both rows whenever fewer than two exist,
			// matching the original behavior.
			_, err = tx.ExecContext(ctx, `
				INSERT INTO active_battles (
					id, match_id, participant_id, lead_status
				)
				VALUES (?, ?, ?, ?), (?, ?, ?, ?)
			`,
				p1BattleID, p1.MatchID, *p1.ParticipantID, boolToInt(p1Leads),
				p2BattleID, p2.MatchID, *p2.ParticipantID, boolToInt(p2Leads),
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (r *SQLiteBracketRepository) getBracketMatches(
	ctx context.Context,
	tournamentID string,
) ([][][]BracketMatch, error) {
	rows, err := r.DB.QueryContext(ctx, `
		SELECT
			n.match_id,
			p.name,
			n.status,
			n.round,
			n.winner_id,
			p.challonge_participant_id,
			p.id,
			n.score,
			COALESCE(bs.complete, 1)
		FROM nodes n
		INNER JOIN brackets b
			ON n.bracket_id = b.id
		LEFT JOIN (
			SELECT
				match_id,
				CASE
					WHEN COUNT(*) = 2
						AND SUM(
							CASE WHEN lead_status = -1 THEN 1 ELSE 0 END
						) = 2
					THEN 1
					ELSE 0
				END AS complete
			FROM active_battles
			GROUP BY match_id
		) bs ON bs.match_id = n.match_id
		LEFT JOIN participants p
			ON p.id = n.participant_id
		WHERE b.tournament_id = ?
	`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Preserve the order in which rounds and matches are first encountered.
	var groupedMatches [][][]BracketMatch
	roundIndexes := make(map[int]int)
	matchIndexes := make(map[int]map[string]int)

	for rows.Next() {
		var m BracketMatch
		var name, winnerID, challongeParticipantID, participantID sql.NullString
		var score sql.NullFloat64
		var complete int

		if err := rows.Scan(
			&m.MatchID,
			&name,
			&m.Status,
			&m.Round,
			&winnerID,
			&challongeParticipantID,
			&participantID,
			&score,
			&complete,
		); err != nil {
			return nil, err
		}

		if name.Valid {
			m.ParticipantName = &name.String
		}
		if winnerID.Valid {
			m.WinnerID = &winnerID.String
		}
		if challongeParticipantID.Valid {
			m.ChallongeParticipantID = &challongeParticipantID.String
		}
		if participantID.Valid {
			m.ParticipantID = &participantID.String
		}
		if score.Valid {
			m.Score = &score.Float64
		}

		m.BothRunFinished = complete == 1

		roundIndex, ok := roundIndexes[m.Round]
		if !ok {
			roundIndex = len(groupedMatches)
			roundIndexes[m.Round] = roundIndex
			groupedMatches = append(groupedMatches, [][]BracketMatch{})
			matchIndexes[m.Round] = make(map[string]int)
		}

		matchIndex, ok := matchIndexes[m.Round][m.MatchID]
		if !ok {
			matchIndex = len(groupedMatches[roundIndex])
			matchIndexes[m.Round][m.MatchID] = matchIndex
			groupedMatches[roundIndex] = append(
				groupedMatches[roundIndex],
				[]BracketMatch{},
			)
		}

		groupedMatches[roundIndex][matchIndex] = append(
			groupedMatches[roundIndex][matchIndex],
			m,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groupedMatches, nil
}

func (r *SQLiteBracketRepository) getActiveBattles(
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
		INNER JOIN brackets b
			ON b.id = n.bracket_id
		WHERE b.tournament_id = ?
	`, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var battles []ActiveBattle

	for rows.Next() {
		var battle ActiveBattle

		if err := rows.Scan(
			&battle.ID,
			&battle.MatchID,
			&battle.ParticipantID,
			&battle.LeadStatus,
		); err != nil {
			return nil, err
		}

		battles = append(battles, battle)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return battles, nil
}

func boolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
