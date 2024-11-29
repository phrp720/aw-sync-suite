package synchronizer

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/settings"
	"aw-sync-agent/system_error"
	util2 "aw-sync-agent/util"
	"fmt"
	"log"
)

// Start starts the synchronization process of data with prometheus
func Start(Config settings.Configuration) error {

	log.Print("==================================================================")
	log.Print("Starting synchronization process...\n")
	log.Print("==================================================================")

	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", Config.Settings.PrometheusUrl, "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(Config.Settings.AWUrl, Config.Settings.ExcludedWatchers)
	userID := util2.GetUserID(Config.Settings.UserID)
	if err != nil {
		return err
	}

	for watcher, data := range scrapedData {
		log.Print("------------------------------------------------------------------")
		log.Print("Processing data for watcher: ", watcher)
		log.Print("------------------------------------------------------------------")

		log.Print("Aggregating data for watcher: [", watcher, "] ...")
		aggregatedData := datamanager.AggregateData(data, watcher, userID, Config.Settings.IncludeHostname, Config.Filters) //metric names must not have '-'
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
func SyncRoutine(Config settings.Configuration) func() {
	return func() {
		if !util2.PromHealthCheck(Config.Settings.PrometheusUrl, Config.Settings.PrometheusSecretKey) {
			log.Print("Something went wrong with Prometheus or Internet connection is lost!")
		} else {
			err := Start(Config)
			system_error.HandleNormal("", err)
		}

	}
}
