package util

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// PromHealthCheck checks the health of Prometheus and the internet connection
func PromHealthCheck(prometheusUrl string) bool {
	url := fmt.Sprintf("%s/graph", prometheusUrl)
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		log.Printf("Prometheus returned status code: %d\n", resp.StatusCode)
	}
	return false
}
