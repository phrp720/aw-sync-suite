package main

import (
	"aw-sync-agent/cron"
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"log"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")

	log.Print("Initializing settings...")
	Settings := settings.InitSettings()

	log.Print("Setting up Sync Cronjob...")
	scheduler := util.ValidateCronExpr(*Settings[settings.Cron])

	c := cron.InitCron()
	cron.Add(c, scheduler, synchronizer.SyncRoutine(Settings))
	cron.Start(c)

	log.Print("Agent Started Successfully")

	// Keep the main program running
	select {}
}
