package service

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/system_error"
	"aw-sync-agent/util"
	"fmt"
	"github.com/phrp720/service-builder/nssm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	WinConfig     = "config.yaml"
	WinExecutable = "aw-sync-agent.exe"
	WinFolder     = "AwSyncAgent"
	WinService    = "aw-sync-agent"
	NssmExe       = "nssm.exe"
)

func CreateWindowsService(sett settings.Settings) {

	// Construct paths relative to the user's home directory
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	system_error.HandleFatal("Failed to get user home directory: ", err)

	windowsRootPath := filepath.Join(homeDir, WinFolder)
	windowsAppPath := filepath.Join(windowsRootPath, WinExecutable)
	windowsConfigPath := filepath.Join(windowsRootPath, WinConfig)

	err = nssm.InitNssm(windowsRootPath)
	system_error.HandleFatal("", err)

	// Stop and remove the service if it already exists
	err = StopAndRemoveService(WinService, windowsRootPath+"/"+NssmExe)
	system_error.HandleNormal("", err)

	util.CopyBinary(windowsAppPath, WinExecutable)

	// Create the config file that will be used for the service(Based on the settings) and loads it  to /opt/aw/ path
	err = settings.CreateConfigFile(sett, windowsConfigPath)
	system_error.HandleFatal("Failed to create config file: ", err)

	builder := nssm.NewServiceBuilder()

	service := builder.ServiceName(WinService).
		AppDirectory(windowsRootPath).
		DisplayName("ActivityWatch Sync Agent").
		Application(windowsAppPath).
		Build()

	err = nssm.CreateService(service)
	system_error.HandleFatal("Failed to create service: ", err)

	err = nssm.StartService(WinService)
	system_error.HandleFatal("Failed to start service: ", err)

	log.Print("Running as a service In Windows...")

}

// StopAndRemoveService stops and removes the service using nssm.exe
func StopAndRemoveService(service string, nssmPath string) error {
	// StopService stops the service using nssm.exe
	cmd := exec.Command(nssmPath, "stop", service)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop the  service: %w", err)
	}

	cmd = exec.Command(nssmPath, "remove", service, "confirm")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to remove the service: %w", err)
	}
	return nil

}
