package synchronizer

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"fmt"
	"log"
	"strings"
)

// Start starts the synchronization process of data with prometheus
func Start(awUrl string, prometheusUrl string) error {
	log.Print("Starting synchronization process...\n")
	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", prometheusUrl, "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(awUrl)
	if err != nil {
		return err
	}
	for watcher, data := range scrapedData {
		log.Print("Pushing data for ", watcher, " ...")
		aggregatedData := datamanager.AggregateData(data, strings.ReplaceAll(watcher, "-", "_")) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, prometheusUrl, aggregatedData)
		if err != nil {
			return err
		}
		log.Print("Data pushed successfully for ", watcher, "\n")
	}
	log.Print("Synchronization process finished successfully\n")
	return nil
}
