package activitywatch

import (
	internalErrors "aw-sync-agent/errors"
	"github.com/phrp720/aw-sync-agent-plugins/models"
	"log"
	"net/http"
	"sort"
	"time"
)

// RemoveExcludedWatchers removes the excluded watchers from the buckets
func RemoveExcludedWatchers(buckets Watchers, excludedWatchers []string) Watchers {
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

func HealthCheck(activityWatchUrl string) bool {

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

func SortAndTrimEvents(events []Event) []Event {
	// Sort events by timestamp. Older to newer.
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// Remove the newest event because it might be incomplete.
	if len(events) > 0 {
		events = events[:len(events)-1]
	}

	return events
}

func ToPluginEvent(events []Event) models.Events {
	var convertedEvents models.Events
	for _, event := range events {

		convertedEvents = append(convertedEvents, models.Event{ID: event.ID, Timestamp: event.Timestamp, Duration: event.Duration, Data: event.Data})
	}
	return convertedEvents
}

func ToAwEvent(events models.Events) Events {
	var convertedEvents Events
	for _, event := range events {
		convertedEvents = append(convertedEvents, Event{ID: event.ID, Timestamp: event.Timestamp, Duration: event.Duration, Data: event.Data})
	}
	return convertedEvents
}
