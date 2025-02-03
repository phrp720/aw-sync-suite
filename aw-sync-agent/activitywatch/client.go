package activitywatch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// GetBuckets gets the buckets from the aw database
func GetBuckets(awUrl string) (Watchers, error) {

	url := awUrl + "/api/0/buckets"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buckets Watchers
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, err
	}
	return buckets, nil
}

// GetEvents gets the events from a specific bucket
func GetEvents(awUrl string, bucket string, start *time.Time, end *time.Time, limit *int) (Events, error) {
	url := fmt.Sprintf("%s/api/0/buckets/%s/events", awUrl, bucket)
	url = addQueryParams(url, start, end, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events Events
	if err = json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, err
	}

	// Filter out events that have the same timestamp as the start point
	if start != nil {
		filteredEvents := make(Events, 0, len(events))
		for _, event := range events {
			if !event.Timestamp.Equal(*start) {
				filteredEvents = append(filteredEvents, event)
			}
		}
		events = filteredEvents
	}

	return events, nil
}

// addQueryParams adds query parameters to the get Events url
func addQueryParams(BaseUrl string, start *time.Time, end *time.Time, limit *int) string {
	// Add start and end time parameters if they are provided
	if start != nil && end != nil {
		BaseUrl = fmt.Sprintf("%s?start=%s&end=%s", BaseUrl, url.QueryEscape(start.Format("2006-01-02T15:04:05.000000-07:00")), url.QueryEscape(end.Format("2006-01-02T15:04:05.000000-07:00")))
	} else if start != nil {
		BaseUrl = fmt.Sprintf("%s?start=%s", BaseUrl, url.QueryEscape(start.Format("2006-01-02T15:04:05.000000-07:00")))
	} else if end != nil {
		BaseUrl = fmt.Sprintf("%s?end=%s", BaseUrl, url.QueryEscape(end.Format("2006-01-02T15:04:05.000000-07:00")))
	}
	if limit != nil {
		if start != nil || end != nil {
			BaseUrl = fmt.Sprintf("%s&limit=%d", BaseUrl, *limit)
		} else {
			BaseUrl = fmt.Sprintf("%s?limit=%d", BaseUrl, *limit)
		}
	}
	return BaseUrl
}
