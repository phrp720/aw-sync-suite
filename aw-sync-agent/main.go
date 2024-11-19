package main

import (
	"aw-sync-agent/cron"
	"aw-sync-agent/filter"
	service2 "aw-sync-agent/service"
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

	log.Print("Validating filters...")
	var totalFilters, invalidFilters, disabledFilters int
	Configs.Filters, totalFilters, invalidFilters, disabledFilters = filter.ValidateFilters(Configs.Filters)
	log.Print("| Total filters: ", totalFilters, "| Disabled Filters: ", disabledFilters, " | Valid filters: ", totalFilters-invalidFilters, " | Invalid filters: ", invalidFilters, " |")

	//	filter.PrintFilters(Configs.Filters)
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
			service2.CreateWindowsService(*Configs)
		} else if util.IsLinux() {
			service2.CreateLinuxService(*Configs)
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
