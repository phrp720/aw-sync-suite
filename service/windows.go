package service

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/util"
	"github.com/phrp720/service-builder/nssm"
	"log"
)

const (
	windows_configFile  = "config.yaml"
	windows_binary      = "aw-sync-agent.exe"
	windows_serviceName = "aw-sync-agent"

	windows_rooPath    = "C:\\AwSyncAgent\\"
	windows_configPath = "C:\\AwSyncAgent\\" + windows_configFile
	windows_appPath    = "C:\\AwSyncAgent\\" + windows_binary
)

func CreateWindowsService(sett settings.Settings) {
	//Here we will handle the windows and windows service creation
	//We will use the nssm for windows and the systemd for windows
	//We will create a service that will run the agent as a service and with -service flag we will pass all the data to the excutable
	//os.Exit(0)
	// Copies  the aw-sync-agent executable to /opt/aw/ path
	util.CopyBinary(windows_appPath, windows_binary)

	// Create the config file that will be used for the service(Based on the settings) and loads it  to /opt/aw/ path
	err := settings.CreateConfigFile(sett, windows_configPath)
	if err != nil {
		log.Fatal("Failed to create config file")
	}

	err = nssm.InitNssm(windows_rooPath)
	if err != nil {
		log.Fatal(err)
	}
	builder := nssm.NewServiceBuilder()
	service := builder.ServiceName(windows_serviceName).
		AppDirectory("C:\\AwSyncAgent\\").
		DisplayName("ActivityWatch Sync Agent").
		Application("C:\\AwSyncAgent\\aw-sync-agent.exe").
		Build()
	nssm.RemoveService(service.ServiceName)
	err = nssm.CreateService(service)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	err = nssm.StartService(service.ServiceName)
	if err != nil {
		log.Fatalf("Failed to start service: %v", err)
	}

	log.Print("Running as a service In Windows...")

}
