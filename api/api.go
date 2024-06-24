// api/api.go
package api

import (
	"barcelona-watch/models"
	"barcelona-watch/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// makeRequest performs an HTTP request and returns the response body.
func makeRequest(apiKey, status string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.football-data.org/v4/teams/81/matches?status=%s&limit=1", status)
	req, err := http.NewRequest("GET", url, nil)
	utils.HandleErr("error creating request", err)

	req.Header.Set("X-Auth-Token", apiKey)
	resp, err := client.Do(req)
	utils.HandleErr("error sending request", err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	utils.HandleErr("error reading response body", err)

	return body, nil
}

// parseResponse parses the JSON response into ApiResponse struct.
func parseResponse(body []byte) (models.ApiResponse, error) {
	var apiResponse models.ApiResponse
	err := json.Unmarshal(body, &apiResponse)
	utils.HandleErr("error parsing JSON", err)
	return apiResponse, nil
}

// CheckFinishedMatches handles the logic for checking the last finished match.
func CheckFinishedMatches(apiKey string) {
	body, err := makeRequest(apiKey, "FINISHED")
	utils.HandleErr("Error", err)

	apiResponse, err := parseResponse(body)
	utils.HandleErr("Error", err)

	if len(apiResponse.Matches) == 0 {
		fmt.Println("No finished matches found.")
		return
	}

	// Handle the result
	match := apiResponse.Matches[0]
	matchTime, err := utils.ParseTime(match.UtcDate)
	utils.HandleErr("Error parsing match date", err)

	// Check if the match was yesterday
	if utils.IsYesterday(matchTime) {
		fmt.Printf("%s %d - %d %s\n",
			match.AwayTeam.Name,
			match.Score.FullTime.Away,
			match.Score.FullTime.Home,
			match.HomeTeam.Name,
		)
	} else {
		fmt.Println("No match played yesterday.")
	}
}

// CheckScheduledMatches handles the logic for checking the next scheduled match.
func CheckScheduledMatches(apiKey string) {
	body, err := makeRequest(apiKey, "SCHEDULED")
	utils.HandleErr("Error", err)

	apiResponse, err := parseResponse(body)
	utils.HandleErr("Error", err)

	if len(apiResponse.Matches) == 0 {
		fmt.Println("No upcoming matches found.")
		return
	}

	match := apiResponse.Matches[0]
	matchTime, err := utils.ParseTime(match.UtcDate)
	utils.HandleErr("Error parsing match date", err)

	// Calculate days remaining until the next match using time.Until
	daysUntilMatch := int(time.Until(matchTime).Hours() / 24)

	// Check if the match is today
	if daysUntilMatch != 0 {
		// Print the days remaining until the match
		fmt.Printf("[Days Until Match: %d] - %s vs %s - %s\n",
			daysUntilMatch,
			match.HomeTeam.Name,
			match.AwayTeam.Name,
			matchTime.Format("2006-01-02 15:04"))
	} else {
		fmt.Printf("[MatchDay] - %s vs %s - %s\n",
			match.HomeTeam.Name,
			match.AwayTeam.Name,
			matchTime.Format("2006-01-02 15:04"))
	}
}
