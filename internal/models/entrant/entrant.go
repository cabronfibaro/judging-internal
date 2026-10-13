package entrant

type Entrant struct {
	ID            string   `json:"id"`
	ParticipantID string   `json:"participantId"`
	TournamentID  string   `json:"tournamentId"`
	Name          string   `json:"name"`
	Status        string   `json:"status"`
	Score         *float64 `json:"score"`
}
