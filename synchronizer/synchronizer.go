package synchronizer

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/settings"
	"aw-sync-agent/system_error"
	"aw-sync-agent/util"
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
	userID := util.GetUserID(Config.Settings.UserID)
	if err != nil {
		return err
	}
	for watcher, data := range scrapedData {
		log.Print("Pushing data for ", watcher, " ...")
		aggregatedData := datamanager.AggregateData(data, watcher, userID, Config.Filters) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, Config.Settings.PrometheusUrl, Config.Settings.PrometheusSecretKey, aggregatedData, watcher)
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
func SyncRoutine(Config settings.Configuration) func() {
	return func() {
		if !util.PromHealthCheck(Config.Settings.PrometheusUrl, Config.Settings.PrometheusSecretKey) {
			log.Print("Something went wrong with Prometheus or Internet connection is lost!")
		} else {
			err := Start(Config)
			system_error.HandleNormal("", err)
		}

	}
}
