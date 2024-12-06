// main.go
package main

import (
	"flag"
	"os"

	"barcelona-watch/api"
	"barcelona-watch/global"
	"barcelona-watch/utils"

	"github.com/joho/godotenv"
)

// FOR SYSTEM
func init() {
	err := godotenv.Load("/home/mohammad/Videos/go/Barcelona-watch/.env")
	utils.HandleErr("Error loading .env file", err)
}

// FOR ACTION
// func init() {
// 	// Get the correct path to the .env file in GitHub Actions
// 	envFile := filepath.Join(os.Getenv("GITHUB_WORKSPACE"), ".env")
// 	err := godotenv.Load(envFile)
// 	utils.HandleErr("Error loading .env file", err)
// }

func flagParser() {
	flag.StringVar(&global.ProxyURL, "proxy", "", "Proxy URL to use for sending Telegram messages")
	flag.Parse()
}

func main() {
	flagParser()

	// Load the API key from .env
	apiKey := os.Getenv("API_KEY")

	// Check for the last finished match
	api.CheckFinishedMatches(apiKey)

	// Check for the next scheduled match
	api.CheckScheduledMatches(apiKey)
}
