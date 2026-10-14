package tournament

type Tournament struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	StartingDate    string  `json:"starting_date"`
	Track           *string `json:"track"`
	QualifyingOrder string  `json:"qualifying_order"`
	Status          string  `json:"status"`
	URL             string  `json:"url"`
	Size            int     `json:"size"`
	Progress        float64 `json:"progress"`
}
