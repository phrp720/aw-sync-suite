package synchronizer

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/settings"
	"aw-sync-agent/util"
	"fmt"
	"log"
	"strings"
)

// Start starts the synchronization process of data with prometheus
func Start(Settings settings.Settings) error {

	log.Print("==================================================================")
	log.Print("Starting synchronization process...\n")
	log.Print("==================================================================")

	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", Settings.PrometheusUrl, "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(Settings.AWUrl, Settings.ExcludedWatchers)
	if err != nil {
		return err
	}
	for watcher, data := range scrapedData {
		log.Print("Pushing data for ", watcher, " ...")
		aggregatedData := datamanager.AggregateData(data, strings.ReplaceAll(watcher, "-", "_")) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, Settings.PrometheusUrl, aggregatedData, watcher)
		if err != nil {
			return err
		}
		log.Print("Data pushed successfully for ", watcher, "\n")
	}

	log.Print("==================================================================")
	log.Print("Synchronization process finished successfully\n")
	log.Print("==================================================================")

	return nil
}

// SyncRoutine returns a function that init the synchronization and starts the  process
func SyncRoutine(Settings settings.Settings) func() {
	return func() {
		if !util.PromHealthCheck(Settings.PrometheusUrl) {
			log.Fatal("Prometheus is not reachable or you don't have internet connection")
		}
		err := Start(Settings)
		if err != nil {
			panic(err) // handle if something wrong happens
		}
	}
}
