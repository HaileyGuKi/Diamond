// models/models.go
package models

type Match struct {
	UtcDate  string `json:"utcDate"`
	HomeTeam struct {
		Name string `json:"name"`
	} `json:"homeTeam"`
	AwayTeam struct {
		Name string `json:"name"`
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
}

type ApiResponse struct {
	Matches []Match `json:"matches"`
}
