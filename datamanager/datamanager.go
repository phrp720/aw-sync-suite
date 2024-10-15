package datamanager

import (
	"aw-sync-agent/aw"
	"aw-sync-agent/util"
	"log"
)

// ScrapeData scrapes the data from the local ActivityWatch instance via the aw Client
func ScrapeData(awUrl string) (aw.WatcherNameToEventsMap, error) {
	log.Print("Fetching buckets  ...\n")
	buckets, err := aw.GetBuckets(awUrl)
	if err != nil {
		log.Printf("Error getting buckets: %s", err)
		return nil, err
	}

	log.Print("Buckets fetched successfully")
	log.Print("Total buckets fetched: ", len(buckets))
	util.RemoveExcludedWatchers(buckets)
	eventsMap := make(aw.WatcherNameToEventsMap)
	for name, bucket := range buckets {
		log.Print("Getting events for bucket ", bucket.Client, " ...")
		events, err := aw.GetEvents(awUrl, name, nil, nil, nil)
		if err != nil {
			log.Printf("Error getting events for bucket %s: %v", bucket.Client, err)
			return nil, err
		}
		eventsMap[bucket.ID] = events
	}
	log.Print("Events fetched successfully\n")
	return eventsMap, nil
}

// AggregateData aggregates the data
func AggregateData() {

}

// PushData pushes the data to the server via the Prometheus Client
func PushData(prometheusUrl string, eventsMap aw.WatcherNameToEventsMap) {

}
