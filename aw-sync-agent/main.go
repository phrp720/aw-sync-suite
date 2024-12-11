package main

import (
	"aw-sync-agent/cron"
	"aw-sync-agent/filter"
	"aw-sync-agent/service"
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"log"
	"os"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")

	log.Print("Initializing configurations...")
	Configs := settings.InitConfigurations()

	if Configs.Filters != nil {
		log.Print("Validating Filters...")
		var totalFilters, invalidFilters, disabledFilters int
		Configs.Filters, totalFilters, invalidFilters, disabledFilters = filter.ValidateFilters(Configs.Filters)
		filter.PrintFilters(totalFilters, invalidFilters, disabledFilters)

		log.Print("Extracting Categories from Filters...")
		categories := filter.GetCategories(Configs.Filters)
		if len(categories) > 0 {
			filter.PrintCategories(categories)
		} else {
			log.Print("No Categories found.")
		}
	}
	if Configs.Settings.TestConfigs {
		log.Print("Test Settings and Filters is finished. Exiting...")
		os.Exit(1)
	}

	// If immediate flag is set, run the sync routine and exit
	if Configs.Settings.Immediate {
		synchronizer.SyncRoutine(*Configs)()
		os.Exit(0)
	}

	// Validate the cron expression and create a scheduler
	scheduler := util.ValidateCronExpr(Configs.Settings.Cron)

	// If the -service flag is set, creates and starts the service and exit
	if Configs.Settings.AsService {

		if util.IsWindows() {
			service.CreateWindowsService(*Configs)
		} else if util.IsLinux() {
			service.CreateLinuxService(*Configs)
		}
		os.Exit(0)

	}

	log.Print("Setting up Sync Cronjob...")
	c := cron.Init()
	cron.Add(c, scheduler, synchronizer.SyncRoutine(*Configs))
	cron.Start(c)

	log.Print("Agent Started Successfully")

	// Keep the main program running
	select {}
}
