package datamanager

import (
	"aw-sync-agent/aw"
	"aw-sync-agent/checkpoint"
	"aw-sync-agent/filter"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/util"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// ScrapeData scrapes the data from the local ActivityWatch instance via the aw Client
func ScrapeData(awUrl string, excludedWatchers []string) (aw.WatcherNameToEventsMap, error) {
	log.Print("Fetching buckets  ...\n")
	buckets, err := aw.GetBuckets(awUrl)
	if err != nil {
		return nil, fmt.Errorf("Error fetching buckets: %v", err)
	}

	log.Print("Buckets fetched successfully")
	log.Print("Total buckets fetched: ", len(buckets))
	util.RemoveExcludedWatchers(buckets, excludedWatchers)
	eventsMap := make(aw.WatcherNameToEventsMap)
	for name, bucket := range buckets {
		log.Print("Fetching events from ", bucket.Client, " ...")
		startPoint := checkpoint.Read(bucket.Client)

		//endPoint := time.Now().AddDate(0, 0, -1) // Set end date to one day before the current date
		events, err := aw.GetEvents(awUrl, name, startPoint, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("error fetching events for bucket %s: %v", bucket.Client, err)
		}
		eventsMap[bucket.Client] = events
	}
	log.Print("Events fetched successfully")
	return eventsMap, nil
}

// AggregateData aggregates the data
// This is going to be called with events for each watcher separately
func AggregateData(events []aw.Event, watcher string, userID string, includeHostName bool, filters []filter.Filter) []prometheus.TimeSeries {

	// Sort events by timestamp. Older to newer.
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// Remove the newest event because it might be incomplete.
	if len(events) > 0 {
		events = events[:len(events)-1]
	}

	var timeSeriesList []prometheus.TimeSeries

	//Apply the filters
	var watcherFilters []filter.Filter
	if watcher != "aw-watcher-afk" {
		watcherFilters = filter.GetMatchingFilters(filters, watcher)
		// Sort watcherFilters so filters with a Category take priority
		sort.Slice(watcherFilters, func(i, j int) bool {
			return watcherFilters[i].Category != "" && watcherFilters[j].Category == ""
		})
	}

	var dropEvent bool
	for _, event := range events {

		//Apply the filters
		if watcher != "aw-watcher-afk" {
			event.Data["category"] = "Other" //Default category
			event.Data, dropEvent = filter.Apply(event.Data, watcherFilters)
		}

		// Drop the event if it matches the filter
		if dropEvent {
			continue
		}
		var labels []prometheus.Label
		labels = append(labels, prometheus.Label{
			Name:  "__name__",
			Value: strings.ReplaceAll(watcher, "-", "_"),
		})
		// Unique ID for each event to avoid duplicate errors of timestamp seconds
		labels = append(labels, prometheus.Label{
			Name:  "unique_id",
			Value: util.GetRandomUUID(),
		})
		//Event ID created from activityWatch
		labels = append(labels, prometheus.Label{
			Name:  "aw_id",
			Value: strconv.Itoa(event.ID),
		})
		labels = append(labels, prometheus.Label{
			Name:  "user",
			Value: userID,
		})
		var hostValue string

		if includeHostName {
			hostValue = util.GetHostname()
		} else {
			hostValue = "Unknown"
		}
		labels = append(labels, prometheus.Label{
			Name:  "host",
			Value: hostValue,
		})
		// Add the data as labels
		for key, value := range event.Data {
			labels = append(labels, prometheus.Label{
				Name:  key,
				Value: fmt.Sprintf("%v", value),
			})
		}
		sample := prometheus.Sample{
			Value: event.Duration,
			Time:  event.Timestamp,
		}

		timeSeries := prometheus.TimeSeries{
			Labels: labels,
			Sample: sample,
		}

		timeSeriesList = append(timeSeriesList, timeSeries)
	}

	return timeSeriesList
}

// PushData pushes  data to the server via the Prometheus Client
func PushData(client *prometheus.Client, prometheusUrl string, prometheusSecretKey string, timeseries []prometheus.TimeSeries, watcher string) error {
	const chunkSize = 20

	log.Print("Pushing data for [", watcher, "] ...")

	for i := 0; i < len(timeseries); i += chunkSize {
		if !util.PromHealthCheck(prometheusUrl, prometheusSecretKey) {
			return errors.New("prometheus is not reachable or Internet connection is lost. Data will be pushed when health is recovered")
		}
		end := i + chunkSize
		if end > len(timeseries) {
			end = len(timeseries)
		}
		chunk := timeseries[i:end]
		_, err := client.Write(context.Background(), prometheusSecretKey, &prometheus.WriteRequest{TimeSeries: chunk})
		if err != nil {
			log.Printf("Error pushing data: %v", err)
			return err
		}

		checkpoint.Update(watcher, chunk[len(chunk)-1].Sample.Time)
		log.Printf("Pushed %d time series records", len(chunk))

	}
	if len(timeseries) == 0 {
		log.Print("No data to push for [", watcher, "]")
	} else {
		log.Print("Data pushed successfully for [", watcher, "]")
	}

	return nil

}
