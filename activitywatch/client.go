package activitywatch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// GetBuckets gets the buckets from the activitywatch database
func GetBuckets() (Watchers, error) {
	awUrl := os.Getenv("ACTIVITY_WATCH_URL")
	if awUrl == "" {
		fmt.Println("Environment variable ACTIVITY_WATCHER_URL is not set or is empty")
		os.Exit(1)
	}
	url := awUrl + "/api/0/buckets"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get buckets: %w", err)
	}
	defer resp.Body.Close()

	var buckets Watchers
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return buckets, nil
}

// GetEvents gets the events from a specific bucket
func GetEvents(bucket string, start *time.Time, end *time.Time, limit *int) (Events, error) {
	awUrl := os.Getenv("ACTIVITY_WATCH_URL")
	if awUrl == "" {
		fmt.Println("Environment variable ACTIVITY_WATCHER_URL is not set or is empty")
		os.Exit(1)
	}
	url := fmt.Sprintf("%s/api/0/buckets/%s/events", awUrl, bucket)
	url = getEventsURL(url, start, end, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer resp.Body.Close()

	var events Events
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return events, nil
}

func getEventsURL(url string, start *time.Time, end *time.Time, limit *int) string {
	// Add start and end time parameters if they are provided
	if start != nil && end != nil {
		url = fmt.Sprintf("%s?start=%s&end=%s", url, start.Format(time.RFC3339), end.Format(time.RFC3339))
	} else if start != nil {
		url = fmt.Sprintf("%s?start=%s", url, start.Format(time.RFC3339))
	} else if end != nil {
		url = fmt.Sprintf("%s?end=%s", url, end.Format(time.RFC3339))

	}
	if limit != nil {
		if start != nil || end != nil {
			url = fmt.Sprintf("%s&limit=%d", url, *limit)
		} else {
			url = fmt.Sprintf("%s?limit=%d", url, *limit)
		}
	}

	return url
}
