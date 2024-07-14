// main.go
package main

import (
	"os"

	"barcelona-watch/api"
	"barcelona-watch/utils"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("/home/mohammad/Videos/go/Barcelona-watch/.env")
	utils.HandleErr("Error loading .env file", err)
}

func main() {
	// Load the API key from .env
	apiKey := os.Getenv("API_KEY")

	// Check for the last finished match
	api.CheckFinishedMatches(apiKey)

	// Check for the next scheduled match
	api.CheckScheduledMatches(apiKey)
}
