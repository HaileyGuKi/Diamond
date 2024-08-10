package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"barcelona-watch/utils"
)

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramPhotoMessage struct {
	ChatID  string `json:"chat_id"`
	Photo   string `json:"photo"`             // URL or path to the photo
	Caption string `json:"caption,omitempty"` // Optional caption for the image
}

const (
	maxRetries    = 5
	retryDelay    = 2 * time.Second
	rateLimitBase = 2
)

func SendToTelegram(message string) {
	botToken := getEnv("TELEGRAM_BOT_TOKEN")
	channelID := getEnv("TELEGRAM_CHANNEL_ID")
	proxyURL := "socks5://0.0.0.0:8086"

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	telegramMessage := TelegramMessage{
		ChatID: channelID,
		Text:   message,
	}

	jsonData, err := json.Marshal(telegramMessage)
	if err != nil {
		utils.HandleErr("Error marshalling Telegram message", err)
		return
	}

	client := utils.CreateHTTPClient(proxyURL)
	retryCount := 0

	for {
		err := sendRequest(client, apiURL, jsonData, &retryCount)
		if err != nil {
			if retryCount >= maxRetries {
				utils.HandleErr("Failed to send message to Telegram after retries", err)
				return
			}
			time.Sleep(retryDelay) // Wait before retrying
			continue
		}
		break // Exit the loop if request was successful
	}
}

func SendPhotoToTelegram(photoPath, caption string) {
	botToken := getEnv("TELEGRAM_BOT_TOKEN")
	channelID := getEnv("TELEGRAM_CHANNEL_ID")
	proxyURL := "socks5://0.0.0.0:8086"

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken)

	// Open the image file to send it
	file, err := os.Open(photoPath)
	if err != nil {
		utils.HandleErr("Error opening photo file", err)
		return
	}
	defer file.Close()

	// Create a buffer to hold the multipart form data
	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	// Create the form file field for the photo
	part, err := writer.CreateFormFile("photo", photoPath)
	if err != nil {
		utils.HandleErr("Error creating form file for photo", err)
		return
	}

	// Copy the photo file into the form file field
	_, err = io.Copy(part, file)
	if err != nil {
		utils.HandleErr("Error copying photo file", err)
		return
	}

	// Add the chat_id and caption fields
	_ = writer.WriteField("chat_id", channelID)
	_ = writer.WriteField("caption", caption)

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		utils.HandleErr("Error closing form writer", err)
		return
	}

	// Create the HTTP client with a proxy if needed
	client := utils.CreateHTTPClient(proxyURL)
	req, err := http.NewRequest("POST", apiURL, form)
	if err != nil {
		utils.HandleErr("Error creating new request", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		utils.HandleErr("Error sending request to Telegram", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.HandleErr(fmt.Sprintf("Failed to send photo, status code: %d", resp.StatusCode), err)
		return
	}

	fmt.Println("Photo sent successfully")
}

func sendRequest(client *http.Client, apiURL string, jsonData []byte, retryCount *int) error {
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("Network timeout, retrying...")
			(*retryCount)++
			return err // Return the error to trigger a retry
		}
		return err // Return the error for non-retryable errors
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil // Success, no need to retry
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := time.Duration(rateLimitBase<<*retryCount) * time.Second
		fmt.Printf("Rate limit exceeded, retrying after %v...", retryAfter)
		(*retryCount)++
		time.Sleep(retryAfter)
		return fmt.Errorf("rate limit exceeded")
	}

	return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("environment variable %s not set", key)
	}
	return value
}

func ValidateProxyURL(proxyURL string) error {
	parsedURL, err := url.Parse(proxyURL)
	utils.HandleErr("Error", err)

	switch parsedURL.Scheme {
	case "http", "https", "socks5":
	default:
		return fmt.Errorf("unsupported proxy scheme: %s", parsedURL.Scheme)
	}

	if parsedURL.Hostname() == "" {
		return fmt.Errorf("missing hostname or IP address in proxy URL")
	}

	return nil
}
