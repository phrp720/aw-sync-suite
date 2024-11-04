package main

import (
	"aw-sync-agent/cron"
	"aw-sync-agent/service"
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"log"
	"os"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")
	log.Print("Initializing settings...")
	Settings := settings.InitSettings()
	scheduler := util.ValidateCronExpr(Settings.Cron)

	if Settings.AsService {

		if util.IsWindows() {
			service.CreateWindowsService(*Settings)
		} else if util.IsLinux() {
			service.CreateLinuxService(*Settings)
		}
		os.Exit(0)

	}

	log.Print("Setting up Sync Cronjob...")
	c := cron.Init()
	cron.Add(c, scheduler, synchronizer.SyncRoutine(*Settings))
	cron.Start(c)

	log.Print("Agent Started Successfully")

	// Keep the main program running
	select {}
}
