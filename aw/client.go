package aw

import (
	"aw-sync-agent/errors"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetBuckets gets the buckets from the aw database
func GetBuckets(awUrl string) (Watchers, error) {

	url := awUrl + "/api/0/buckets"
	resp, err := http.Get(url)
	if err != nil {
		return nil, &errors.HTTPError{URL: url, Err: err}
	}
	defer resp.Body.Close()

	var buckets Watchers
	if err := json.NewDecoder(resp.Body).Decode(&buckets); err != nil {
		return nil, &errors.DecodeError{Err: err}
	}
	return buckets, nil
}

// GetEvents gets the events from a specific bucket
func GetEvents(awUrl string, bucket string, start *time.Time, end *time.Time, limit *int) (Events, error) {

	url := fmt.Sprintf("%s/api/0/buckets/%s/events", awUrl, bucket)
	url = addQueryParams(url, start, end, limit)

	resp, err := http.Get(url)
	if err != nil {
		return nil, &errors.HTTPError{URL: url, Err: err}
	}
	defer resp.Body.Close()

	var events Events
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, &errors.DecodeError{Err: err}
	}
	return events, nil
}

// addQueryParams adds query parameters to the get Events url
func addQueryParams(url string, start *time.Time, end *time.Time, limit *int) string {
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