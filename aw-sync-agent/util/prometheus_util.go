package util

import (
	internalErrors "aw-sync-agent/errors"
	"aw-sync-agent/prometheus"
	"fmt"
	"log"
	"net/http"
	"regexp"
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
		internalErrors.HandleNormal("Failed to create Prometheus health-check request: ", err)
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

func AddMetricLabel(labels *[]prometheus.Label, key string, value string) {
	if value != "" {
		*labels = append(*labels, prometheus.Label{
			Name:  key,
			Value: value,
		})
	}
}

// SanitizeLabelName ensures the label name conforms to Prometheus naming conventions
func SanitizeLabelName(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9_]`)
	return re.ReplaceAllString(name, "_")
}
