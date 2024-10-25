package main

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
	"aw-sync-agent/util"
	"github.com/robfig/cron"
	"log"
)

func main() {

	log.Print("Starting ActivityWatch Sync Agent...")

	log.Print("Initializing settings...")
	Settings := settings.InitSettings()

	log.Print("Setting up Sync Cronjob...")
	scheduler := util.ValidateCronExpr(*Settings[settings.Cron])
	sync := func() {
		if !util.PromHealthCheck(*Settings[settings.PrometheusUrl]) {
			log.Fatal("Prometheus is not reachable or you don't have internet connection")
		}
		err := synchronizer.Start(Settings)
		if err != nil {
			panic(err) // handle if something wrong happens
		}
	}
	c := cron.New()

	err := c.AddFunc(scheduler, sync)
	if err != nil {
		log.Print(err)
	}
	c.Start()
	log.Print("Agent Started Successfully")

	// Keep the main program running
	select {}
}
