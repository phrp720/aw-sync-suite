package util

import (
	"aw-sync-agent/aw"
	internalErrors "aw-sync-agent/errors"
	"log"
	"net/http"
	"time"
)

// RemoveExcludedWatchers removes the excluded watchers from the buckets
func RemoveExcludedWatchers(buckets aw.Watchers, excludedWatchers []string) aw.Watchers {
	if len(excludedWatchers) > 0 {
		for _, excludedWatcher := range excludedWatchers {
			for id, bucket := range buckets {
				if bucket.Client == excludedWatcher {
					delete(buckets, id)
				}
			}
		}
	}
	log.Print("Buckets excluded: ", len(excludedWatchers), excludedWatchers, "\n")
	return buckets
}

func ActivityWatchHealthCheck(activityWatchUrl string) bool {

	resp, err := getRequest(activityWatchUrl)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	} else {
		log.Printf("ActivityWatch returned status code: %d\n", resp.StatusCode)
	}
	return false
}

// GetRequest makes a request to the given URL
func getRequest(url string) (*http.Response, error) {

	client := http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		internalErrors.HandleNormal("Failed to create ActivityWatch health-check request: ", err)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
