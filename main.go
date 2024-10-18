package main

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
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
		err := synchronizer.Start(Settings)
		if err != nil {
			panic(err) // handle if something wrong happens
		}
	}
}
