package util

import (
	"aw-sync-agent/aw-sync-agent/system_error"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PromHealthCheck checks the health of Prometheus and the internet connection
func PromHealthCheck(prometheusUrl string, secretKey string) bool {
	url := fmt.Sprintf("%s/-/healthy", prometheusUrl)
	resp, err := MakeRequest(url, secretKey)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else if resp.StatusCode == http.StatusUnauthorized {
		log.Printf("Unauthorized access to Prometheus. Status Code: %d\n", resp.StatusCode)
	} else {
		log.Printf("Prometheus returned status code: %d\n", resp.StatusCode)
	}
	return false
}

// MakeRequest makes a request to the given URL
func MakeRequest(url string, secretKey string) (*http.Response, error) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		system_error.HandleNormal("Failed to create Prometheus health-check request: ", err)
	}

	// Set Bearer if exists
	if secretKey != "" {
		req.Header.Set("Authorization", "Bearer "+secretKey)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
