package challonge

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type ChallongeTournamentResponse struct {
	Data struct {
		ID         string `json:"id"`
		Attributes struct {
			URL string `json:"url"`
		} `json:"attributes"`
	} `json:"data"`
}

func CreateTournament(
	ctx context.Context,
	apiKey string,
	name string,
	size int,
	startingDate time.Time,
) (*ChallongeTournamentResponse, error) {
	const tournamentsURL = "https://api.challonge.com/v2.1/tournaments.json"
	payload := map[string]any{
		"data": map[string]any{
			"type": "tournaments",
			"attributes": map[string]any{
				"name":            name,
				"url":             fmt.Sprintf("%s_%d", strings.ReplaceAll(strings.ToLower(name), " ", "_"), time.Now().UnixMilli()),
				"tournament_type": "single elimination",
				"game_name":       "CarX Drift Racing",
				"private":         false,
				"starts_at":       startingDate.UTC().Format(time.RFC3339Nano),
				"description":     fmt.Sprintf("Top %d head-to-head tandem battle knockout bracket.", size),
				"match_options": map[string]any{
					"consolation_matches_target_rank": 3,
					"accept_attachments":              false,
				},
				"registration_options": map[string]any{
					"open_signup": false,
				},
				"seeding_options": map[string]any{
					"hide_seeds":          false,
					"sequential_pairings": false,
				},
				"station_options": map[string]any{
					"auto_assign": false,
					"only_start_matches_with_assigned_stations": false,
				},
				"group_stage_enabled": false,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		tournamentsURL,
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create Challonge tournament: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf(
			"Challonge returned HTTP %d: %s",
			resp.StatusCode,
			string(respBody),
		)
	}

	var result ChallongeTournamentResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode Challonge response: %w", err)
	}

	return &result, nil
}

func DeleteTournament(ctx context.Context, apiKey string, tournamentID string) error {
	url := fmt.Sprintf(
		"https://api.challonge.com/v2.1/tournaments/%s.json",
		tournamentID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete Challonge tournament: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("[DeleteTournament] Challonge tournament with ID/URL %s not found", tournamentID)
		}

		return fmt.Errorf("Challonge returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func UpdateTournamentState(
	ctx context.Context,
	apiKey string,
	tournamentID string,
	state string,
) error {
	payload := map[string]any{
		"data": map[string]any{
			"type": "TournamentState",
			"attributes": map[string]any{
				"state": state,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal state payload: %w", err)
	}

	url := fmt.Sprintf(
		"https://api.challonge.com/v2.1/tournaments/%s/change_state.json",
		tournamentID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("change Challonge tournament state: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusNotFound {
			return fmt.Errorf("[UpdateTournamentState] Challonge with ID/URL %s not found", tournamentID)
		}
		return fmt.Errorf("Challonge returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

type ParticipantReference struct {
	Data struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"data"`
}

type ParticipantScore struct {
	ParticipantID int   `json:"participant_id"`
	Scores        []int `json:"scores"`
}

type MatchTimestamps struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type MatchRelationships struct {
	Player1 ParticipantReference `json:"player1"`
	Player2 ParticipantReference `json:"player2"`
}

type MatchAttributes struct {
	State               string             `json:"state"`
	Round               int                `json:"round"`
	Identifier          string             `json:"identifier"`
	SuggestedPlayOrder  int                `json:"suggested_play_order"`
	Scores              string             `json:"scores"`
	ScoreInSets         [][2]int           `json:"score_in_sets"`
	PointsByParticipant []ParticipantScore `json:"points_by_participant"`
	Timestamps          MatchTimestamps    `json:"timestamps"`
	WinnerID            int                `json:"winner_id"`
	Relationships       MatchRelationships `json:"relationships"`
}

type MatchData struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Attributes MatchAttributes `json:"attributes"`
}

type MatchResponse struct {
	Data []MatchData `json:"data"`
}

func GetMatches(
	ctx context.Context,
	challongeTournamentID string,
	apiKey string,
) ([]MatchData, error) {
	url := fmt.Sprintf(
		"https://api.challonge.com/v2.1/tournaments/%s/matches.json",
		challongeTournamentID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization-Type", "v1")
	req.Header.Set("Authorization", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get Challonge matches: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		if resp.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("[GetMatches] Challonge with tournament ID/URL %s not found", challongeTournamentID)
		}
		return nil, fmt.Errorf("Challonge returned HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	var result MatchResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("decode Challonge matches: %w", err)
	}

	return result.Data, nil
}
