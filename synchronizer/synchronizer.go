package synchronizer

import (
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/settings"
	"fmt"
	"log"
	"strings"
)

// Start starts the synchronization process of data with prometheus
func Start(Settings map[settings.SettingsKey]*string) error {

	log.Print("==================================================================")
	log.Print("Starting synchronization process...\n")
	log.Print("==================================================================")

	prometheusClient := prometheus.NewClient(fmt.Sprintf("%s%s", *Settings[settings.PrometheusUrl], "/api/v1/write"))
	scrapedData, err := datamanager.ScrapeData(*Settings[settings.AWUrl], *Settings[settings.ExcludedWatchers])
	if err != nil {
		return err
	}
	for watcher, data := range scrapedData {
		log.Print("Pushing data for ", watcher, " ...")
		aggregatedData := datamanager.AggregateData(data, strings.ReplaceAll(watcher, "-", "_")) //metric names must not have '-'
		err = datamanager.PushData(prometheusClient, *Settings[settings.PrometheusUrl], aggregatedData, watcher)
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
