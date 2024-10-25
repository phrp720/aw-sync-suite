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
	if settings.IsService(*Settings[settings.AsService]) {
		//Here we will handle the windows and linux service creation
		//We will use the nssm for windows and the systemd for linux
		//We will create a service that will run the agent as a service and with -service flag we will pass all the data to the excutable
		if util.IsWindows() {
			service.CreateWindowsService()
		} else if util.IsLinux() {
			service.CreateLinuxService()
		}
		os.Exit(0)
	}
	if settings.IsStandalone(*Settings[settings.Standalone]) {
		log.Print("Running as a service...")
		log.Print("Setting up Sync Cronjob...")
		scheduler := util.ValidateCronExpr(*Settings[settings.Cron])
		print(scheduler)
		c := cron.Init()
		cron.Add(c, "@every 5s", synchronizer.SyncRoutine(Settings))
		cron.Start(c)

		log.Print("Agent Started Successfully")

		// Keep the main program running
		select {}
	}
	log.Print("If you want to run the agent as a service, please run the agent with the -service flag")
	log.Print("If you want to run the agent as a standalone process, please run the agent with the -standalone flag")
}
