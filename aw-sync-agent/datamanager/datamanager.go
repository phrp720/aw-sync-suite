package datamanager

import (
	aw2 "aw-sync-agent/aw"
	"aw-sync-agent/checkpoint"
	"aw-sync-agent/filter"
	prometheus2 "aw-sync-agent/prometheus"
	util2 "aw-sync-agent/util"
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

// ScrapeData scrapes the data from the local ActivityWatch instance via the aw Client
func ScrapeData(awUrl string, excludedWatchers []string) (aw2.WatcherNameToEventsMap, error) {
	log.Print("Fetching buckets  ...\n")
	buckets, err := aw2.GetBuckets(awUrl)
	if err != nil {
		return nil, fmt.Errorf("Error fetching buckets: %v", err)
	}

	log.Print("Buckets fetched successfully")
	log.Print("Total buckets fetched: ", len(buckets))
	util2.RemoveExcludedWatchers(buckets, excludedWatchers)
	eventsMap := make(aw2.WatcherNameToEventsMap)
	for name, bucket := range buckets {
		log.Print("Fetching events from ", bucket.Client, " ...")
		startPoint := checkpoint.Read(bucket.Client)
		//endPoint := time.Now().AddDate(0, 0, -1) // Set end date to one day before the current date
		events, err := aw2.GetEvents(awUrl, name, startPoint, nil, nil)
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
func AggregateData(events []aw2.Event, watcher string, userID string, filters []filter.Filter) []prometheus2.TimeSeries {

	// Sort events by timestamp. Older to newer.
	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	// Remove the newest event because it might be incomplete.
	if len(events) > 0 {
		events = events[:len(events)-1]
	}

	var timeSeriesList []prometheus2.TimeSeries

	watcherFilters := filter.GetMatchingFilters(filters, watcher)
	for _, event := range events {

		//Apply the filters
		event.Data = filter.Apply(event.Data, watcherFilters)

		var labels []prometheus2.Label
		labels = append(labels, prometheus2.Label{
			Name:  "__name__",
			Value: strings.ReplaceAll(watcher, "-", "_"),
		})
		labels = append(labels, prometheus2.Label{
			Name:  "user",
			Value: userID,
		})
		for key, value := range event.Data {
			labels = append(labels, prometheus2.Label{
				Name:  key,
				Value: fmt.Sprintf("%v", value),
			})
		}
		sample := prometheus2.Sample{
			Value: event.Duration,
			Time:  event.Timestamp,
		}

		timeSeries := prometheus2.TimeSeries{
			Labels: labels, // Add more labels as needed
			Sample: sample,
		}

		timeSeriesList = append(timeSeriesList, timeSeries)
	}

	return timeSeriesList
}

// PushData pushes  data to the server via the Prometheus Client
func PushData(client *prometheus2.Client, prometheusUrl string, prometheusSecretKey string, timeseries []prometheus2.TimeSeries, watcher string) error {
	const chunkSize = 20
	for i := 0; i < len(timeseries); i += chunkSize {
		if !util2.PromHealthCheck(prometheusUrl, prometheusSecretKey) {
			return errors.New("prometheus is not reachable or Internet connection is lost. Data will be pushed when health is recovered")
		}
		end := i + chunkSize
		if end > len(timeseries) {
			end = len(timeseries)
		}
		chunk := timeseries[i:end]
		_, err := client.Write(context.Background(), prometheusSecretKey, &prometheus2.WriteRequest{TimeSeries: chunk})
		if err != nil {
			log.Printf("Error pushing data: %v", err)
			return err
		}
		checkpoint.Update(watcher, chunk[len(chunk)-1].Sample.Time)
		log.Printf("Pushed %d time series records", len(chunk))
	}
	return nil

}
