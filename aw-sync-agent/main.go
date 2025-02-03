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
	"github.com/phrp720/aw-sync-agent-plugins/plugins"
	"log"
	"os"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")

	log.Print("Initializing configurations...")
	Configs := settings.InitConfigurations()
	Configs.Settings.UserID = util.GetUserID(Configs.Settings.UserID)

	// Here abstract init of the plugins. In init we will load the plugins,load their configs if exists and do the checks like below.

	Plugins := plugins.Select(plugins.Initialize(), Configs.Settings.Plugins)
	if len(Plugins) == 0 {
		log.Print("No plugins are loaded.")
	} else {
		util.PrintPlugins(Plugins)
		if Configs.Settings.PluginsStrictOrder {
			Plugins = util.SortPlugins(Configs.Settings.Plugins, Plugins)
		}
		for _, plugin := range Plugins {
			plugin.Initialize()
		}
	}

	// If the -testConfig flag is set, test the configurations and filters and exit
	if Configs.Settings.TestConfigs {
		log.Print("Testing Settings and Plugins configuration is finished. Exiting...")
		os.Exit(1)
	}

	// If -immediate flag is set, run the sync routine and exit
	if Configs.Settings.Immediate {
		synchronizer.SyncRoutine(*Configs, Plugins)()
		os.Exit(0)
	}

	// Validate the cron expression and create a scheduler
	scheduler := util.ValidateCronExpr(Configs.Settings.Cron)
	// If the -service flag is set, creates and starts the service and exit
	if Configs.Settings.AsService {

		if util.IsWindows() {
			service.CreateWindowsService(*Configs, Plugins)
		} else if util.IsLinux() {
			service.CreateLinuxService(*Configs, Plugins)
		}
		os.Exit(0)

	}

	log.Print("Setting up Sync Cronjob...")
	c := cron.Init()
	cron.Add(c, scheduler, synchronizer.SyncRoutine(*Configs, Plugins))
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
