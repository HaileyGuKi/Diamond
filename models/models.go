// models/models.go
package models

type Match struct {
	UtcDate  string `json:"utcDate"`
	HomeTeam struct {
		Name string `json:"name"`
		Logo string `json:"crest"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Name string `json:"name"`
		Logo string `json:"crest"`
	} `json:"awayTeam"`
	Score struct {
		FullTime struct {
			Home int `json:"home"`
			Away int `json:"away"`
		} `json:"fullTime"`
	} `json:"score"`
	Season struct {
		CurrentMatchday int `json:"currentMatchday"`
	} `json:"season"`
	Competition struct {
		Type string `json:"type"`
	} `json:"competition"`
}

type ApiResponse struct {
	Matches []Match `json:"matches"`
}
