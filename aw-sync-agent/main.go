package main

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/cron"
	internalErrors "aw-sync-agent/errors"
	"aw-sync-agent/prometheus"
	"aw-sync-agent/service"
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"errors"
	"log"
	"os"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")

	log.Print("Initializing configurations...")
	Configs := settings.InitConfigurations()
	Configs.Settings.UserID = util.GetUserID(Configs.Settings.UserID)

	// Here abstract init of the plugins. In init we will load the plugins,load their configs if exists and do the checks like below.

	//if Configs.Filters != nil {
	//	log.Print("Validating Filters...")
	//	var totalFilters, invalidFilters, disabledFilters int
	//	Configs.Filters, totalFilters, invalidFilters, disabledFilters = filter.ValidateFilters(Configs.Filters)
	//	filter.PrintFilters(totalFilters, invalidFilters, disabledFilters)
	//
	//	log.Print("Extracting Categories from Filters...")
	//	categories := filter.GetCategories(Configs.Filters)
	//	if len(categories) > 0 {
	//		filter.PrintCategories(categories)
	//	} else {
	//		log.Print("No Categories found.")
	//	}
	//}

	// If the -testConfig flag is set, test the configurations and filters and exit
	if Configs.Settings.TestConfigs {
		log.Print("Testing Settings and Plugins configuration is finished. Exiting...")
		os.Exit(1)
	}

	// If -immediate flag is set, run the sync routine and exit
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
	defer cron.Stop(c)

	if !prometheus.HealthCheck(Configs.Settings.PrometheusUrl, Configs.Settings.PrometheusSecretKey) {
		internalErrors.HandleNormal("Warning:", errors.New("Prometheus is not reachable or Internet connection is lost. Please fix the issue before the next synchronization."))
	}
	if !activitywatch.HealthCheck(Configs.Settings.AWUrl) {
		internalErrors.HandleNormal("Warning:", errors.New("ActivityWatch is not reachable. Please fix the issue before the next synchronization."))
	}

	log.Print("Agent Started Successfully!")

	// Keep the main program running
	select {}

}
