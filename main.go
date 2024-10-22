package main

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"log"
)

func main() {

	// Initialize the settings
	Settings := settings.InitSettings()
	// Check if the agent should run as a service
	if *Settings[settings.AsService] == "true" {
		log.Println("Running as a service")
		// Add your code to run the agent as a service here
	} else {
		log.Println("Running as a regular application")
		// Pass the map to the synchronizer.Start function
		if !util.PromHealthCheck(*Settings[settings.PrometheusUrl]) {
			log.Fatal("Prometheus is not reachable or you don't have internet connection")
		}
		err := synchronizer.Start(Settings)
		if err != nil {
			panic(err) // handle if something wrong happens
		}
	}
}
