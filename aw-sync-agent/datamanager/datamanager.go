package datamanager

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/checkpoint"
	"aw-sync-agent/filter"
	"aw-sync-agent/prometheus"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
)

// ScrapeData scrapes the data from the local ActivityWatch instance via the aw Client
func ScrapeData(awUrl string, excludedWatchers []string) (activitywatch.WatcherNameToEventsMap, error) {
	if !activitywatch.HealthCheck(awUrl) {
		return nil, errors.New("activityWatch is not reachable. Data will be pushed at the next synchronization")
	}
	log.Print("Fetching buckets  ...\n")
	buckets, err := activitywatch.GetBuckets(awUrl)
	if err != nil {
		return nil, fmt.Errorf("Error fetching buckets: %v", err)
	}

	log.Print("Buckets fetched successfully")
	log.Print("Total buckets fetched: ", len(buckets))
	activitywatch.RemoveExcludedWatchers(buckets, excludedWatchers)
	eventsMap := make(activitywatch.WatcherNameToEventsMap)
	for name, bucket := range buckets {
		if !activitywatch.HealthCheck(awUrl) {
			return nil, errors.New("activityWatch is not reachable. Data will be pushed at the next synchronization")
		}
		log.Print("Fetching events from ", bucket.Client, " ...")
		startPoint := checkpoint.Read(bucket.Client)

		//endPoint := time.Now().AddDate(0, 0, -1) // Set end date to one day before the current date
		events, err := activitywatch.GetEvents(awUrl, name, startPoint, nil, nil)
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
func AggregateData(events []activitywatch.Event, watcher string, userID string, includeHostName bool, filters []filter.Filter) []prometheus.TimeSeries {

	events = activitywatch.SortAndTrimEvents(events)

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

		// Here it will be the abstract run of each plugin.We can follow strict order of execution.Each plugin will have its own function and must return Event type.

		//Apply the filters
		if watcher != "aw-watcher-afk" {
			event.Data["category"] = "Other" //Default category
			event.Data, dropEvent = filter.Apply(event.Data, watcherFilters)
		}

		// Drop the event if it matches the filter
		if dropEvent {
			continue
		}

		timeSeriesList = append(timeSeriesList, prometheus.AttachTimeSeriesPayload(event, includeHostName, watcher, userID))
	}
	return timeSeriesList
}

// PushData pushes  data to the server via the Prometheus Client
func PushData(client *prometheus.Client, prometheusUrl string, prometheusSecretKey string, timeseries []prometheus.TimeSeries, watcher string) error {
	const chunkSize = 20

	log.Print("Pushing data for [", watcher, "] ...")

	for i := 0; i < len(timeseries); i += chunkSize {
		if !prometheus.HealthCheck(prometheusUrl, prometheusSecretKey) {
			return errors.New("prometheus is not reachable or Internet connection is lost. Data will be pushed at the next synchronization")
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
