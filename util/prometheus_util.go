package util

import (
	"fmt"
	"net/http"
	"time"
)

// PromHealthCheck checks the health of Prometheus and the internet connection
func PromHealthCheck(prometheusUrl string) bool {
	fmt.Print("Prometheus health and Internet Connection check  started\n")
	url := fmt.Sprintf("%s/graph", prometheusUrl)
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error checking URL: %v\n", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Prometheus is up and running and connection to the internet is available")
		return true
	} else {
		fmt.Printf("Prometheus returned status code: %d\n", resp.StatusCode)
	}
	return false
}
