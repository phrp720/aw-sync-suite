package datamanager

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/checkpoint"
	"aw-sync-agent/prometheus"
	"context"
	"errors"
	"fmt"
	"github.com/phrp720/aw-sync-agent-plugins/models"
	"log"
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
func AggregateData(Plugins []models.Plugin, events []activitywatch.Event, watcher string, userID string, includeHostName bool) []prometheus.TimeSeries {

	events = activitywatch.SortAndTrimEvents(events)
	var timeSeriesList []prometheus.TimeSeries
	var unmarshaledEvents models.Events
	for _, plugin := range Plugins {
		unmarshaledEvents = plugin.Execute(activitywatch.ToPluginEvent(events), watcher, userID, includeHostName)
	}
	if len(Plugins) > 0 {
		events = activitywatch.ToAwEvent(unmarshaledEvents)
	}
	for _, event := range events {
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
