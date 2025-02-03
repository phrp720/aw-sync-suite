package synchronizer

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/datamanager"
	internalErrors "aw-sync-agent/errors"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/settings"
	"fmt"
	"github.com/phrp720/aw-sync-agent-plugins/models"
	"log"
)

// Start starts the synchronization process of data with prometheus
func Start(Config settings.Configuration, Plugins []models.Plugin) error {

	log.Print("==================================================================")
	log.Print("Starting synchronization process...\n")
	log.Print("==================================================================")

	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", Config.Settings.PrometheusUrl, "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(Config.Settings.AWUrl, Config.Settings.ExcludedWatchers)
	if err != nil {
		return err
	}

	for watcher, data := range scrapedData {
		log.Print("------------------------------------------------------------------")
		log.Print("Processing data for watcher: ", watcher)
		log.Print("------------------------------------------------------------------")

		log.Print("Aggregating data for watcher: [", watcher, "] ...")
		aggregatedData := datamanager.AggregateData(Plugins, data, watcher, Config.Settings.UserID, Config.Settings.IncludeHostname) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, Config.Settings.PrometheusUrl, Config.Settings.PrometheusSecretKey, aggregatedData, watcher)
		if err != nil {
			return err
		}

	}
	log.Print("==================================================================")
	log.Print("Synchronization process finished successfully\n")
	log.Print("==================================================================")

	return nil
}

// SyncRoutine returns a function that init the synchronization and starts the  process
func SyncRoutine(Config settings.Configuration, Plugins []models.Plugin) func() {
	return func() {
		if !prometheus.HealthCheck(Config.Settings.PrometheusUrl, Config.Settings.PrometheusSecretKey) {
			log.Print("Something went wrong with Prometheus or Internet connection is lost! Data will be pushed at the next synchronization.")
		} else if !activitywatch.HealthCheck(Config.Settings.AWUrl) {
			log.Print("ActivityWatch is not reachable! Data will be pushed at the next synchronization.")
		} else {
			err := Start(Config, Plugins)
			internalErrors.HandleNormal("Error:", err)
		}

	}
}
