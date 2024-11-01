package service

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/system_error"
	"aw-sync-agent/util"
	"github.com/phrp720/service-builder/nssm"
	"log"
	"os"
	"path/filepath"
)

const (
	WinConfig     = "config.yaml"
	WinExecutable = "aw-sync-agent.exe"
	WinFolder     = "AwSyncAgent"
	WinService    = "aw-sync-agent"
)

func CreateWindowsService(sett settings.Settings) {

	// Construct paths relative to the user's home directory
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	system_error.HandleFatal("Failed to get user home directory: ", err)

	windowsRootPath := filepath.Join(homeDir, WinFolder)
	windowsAppPath := filepath.Join(windowsRootPath, WinExecutable)
	windowsConfigPath := filepath.Join(windowsRootPath, WinConfig)

	util.CopyBinary(windowsRootPath, WinExecutable)

	// Create the config file that will be used for the service(Based on the settings) and loads it  to /opt/aw/ path
	err = settings.CreateConfigFile(sett, windowsConfigPath)
	system_error.HandleFatal("Failed to create config file: ", err)

	err = nssm.StartNssm(windowsRootPath)
	system_error.HandleFatal("", err)

	builder := nssm.NewServiceBuilder()
	service := builder.ServiceName(WinService).
		AppDirectory(windowsRootPath).
		DisplayName("ActivityWatch Sync Agent").
		Application(windowsAppPath).
		Build()

	nssm.RemoveService(service.ServiceName)

	err = nssm.CreateService(service)
	system_error.HandleFatal("Failed to create service: ", err)

	err = nssm.StartService(service.ServiceName)
	system_error.HandleFatal("Failed to start service: ", err)

	log.Print("Running as a service In Windows...")

}
