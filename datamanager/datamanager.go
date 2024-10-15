package datamanager

import (
	"aw-sync-agent/aw"
	"fmt"
	"log"
	"os"
	"strings"
)

// ScrapeData scrapes the data from the local ActivityWatch instance via the aw Client
func ScrapeData() aw.WatcherNameToEventsMap {
	fmt.Print("Fetching buckets  ...\n")
	buckets, err := aw.GetBuckets()
	if err != nil {
		log.Fatal("Error getting buckets: ", err)

	}

	fmt.Print("Buckets fetched successfully\n", "Total buckets fetched: ", len(buckets), "\n")
	RemoveExcludedWatchers(buckets)
	eventsMap := make(aw.WatcherNameToEventsMap)
	for name, bucket := range buckets {
		fmt.Print("Getting events for bucket ", bucket.Client, " ...\n")
		events, err := aw.GetEvents(name, nil, nil, nil)
		if err != nil {
			log.Fatal("Error getting events for bucket ", bucket.Client, " : ", err)
		}
		eventsMap[bucket.ID] = events
	}
	fmt.Print("Events fetched successfully\n")
	return eventsMap
}

// AggregateData aggregates the data
func AggregateData() {

}

// PushData pushes the data to the server via the Prometheus Client
func PushData() {

}
func getExcludedWatchers() []string {
	return strings.Split(os.Getenv("EXCLUDED_WATCHERS"), ",")
}

func RemoveExcludedWatchers(buckets aw.Watchers) aw.Watchers {
	excluded := getExcludedWatchers()
	if len(excluded) > 0 {
		for _, excludedWatcher := range excluded {
			for id, bucket := range buckets {
				if bucket.Client == excludedWatcher {
					delete(buckets, id)
				}
			}
		}
	}
	fmt.Print("Buckets excluded: ", len(excluded), excluded, "\n")
	return buckets
}
