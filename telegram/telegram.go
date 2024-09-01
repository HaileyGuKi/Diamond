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

	"barcelona-watch/global"
	"barcelona-watch/utils"
)

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

const (
	maxRetries    = 5
	retryDelay    = 2 * time.Second
	rateLimitBase = 2
)

func SendToTelegram(message string) {
	botToken := utils.GetEnv("TELEGRAM_BOT_TOKEN")
	channelID := utils.GetEnv("TELEGRAM_CHANNEL_ID")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	telegramMessage := TelegramMessage{
		ChatID: channelID,
		Text:   message,
	}

	jsonData, err := json.Marshal(telegramMessage)
	utils.HandleErr("Error marshalling Telegram message", err)

	client := utils.CreateHTTPClient(global.ProxyURL)
	retryCount := 0

	for {
		err := sendRequest(client, apiURL, jsonData, &retryCount)
		if err != nil {
			if retryCount >= maxRetries {
				utils.HandleErr("Failed to send message to Telegram after retries", err)
			}
			time.Sleep(retryDelay) // Wait before retrying
			continue
		}
		break // Exit the loop if request was successful
	}
}

func SendPhotoToTelegram(photoPath, caption string) {
	botToken := utils.GetEnv("TELEGRAM_BOT_TOKEN")
	channelID := utils.GetEnv("TELEGRAM_CHANNEL_ID")
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", botToken)

	// Open the image file to send it
	file, err := os.Open(photoPath)
	utils.HandleErr("Error opening photo file", err)
	defer file.Close()

	// Create a buffer to hold the multipart form data
	form := &bytes.Buffer{}
	writer := multipart.NewWriter(form)

	// Create the form file field for the photo
	part, err := writer.CreateFormFile("photo", photoPath)
	utils.HandleErr("Error creating form file for photo", err)

	// Copy the photo file into the form file field
	_, err = io.Copy(part, file)
	utils.HandleErr("Error copying photo file", err)

	// Add the chat_id and caption fields
	_ = writer.WriteField("chat_id", channelID)
	_ = writer.WriteField("caption", caption)

	// Close the writer to finalize the form data
	err = writer.Close()
	utils.HandleErr("Error closing form writer", err)

	// Create the HTTP client with a proxy if needed
	client := utils.CreateHTTPClient(global.ProxyURL)
	req, err := http.NewRequest("POST", apiURL, form)
	utils.HandleErr("Error creating new request", err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	resp, err := client.Do(req)
	utils.HandleErr("Error sending request to Telegram", err)
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		utils.HandleErr(fmt.Sprintf("Failed to send photo, status code: %d", resp.StatusCode), err)
	}
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
