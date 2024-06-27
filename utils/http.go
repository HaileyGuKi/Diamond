package utils

import (
	"net/http"
	"net/url"
	"time"
)

func CreateHTTPClient(proxyURL string) *http.Client {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			HandleErr("Error parsing proxy URL", err)
			// In case of an error, return the client without proxy settings
			return client
		}
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
	}

	return client
}
