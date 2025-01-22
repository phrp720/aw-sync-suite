package service

import (
	internalErrors "aw-sync-agent/errors"
	"aw-sync-agent/settings"
	"aw-sync-agent/util"
	"fmt"
	"github.com/phrp720/aw-sync-agent-plugins/models"
	"github.com/phrp720/go-service-builder/nssm"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	WinConfig     = "aw-sync-settings.yaml"
	WinExecutable = "aw-sync-agent.exe"
	WinFolder     = "AwSyncAgent"
	WinService    = "aw-sync-agent"
)

func CreateWindowsService(config settings.Configuration, plugins []models.Plugin) {

	// Construct paths relative to the user's home directory
	// Get the user's home directory
	homeDir, err := os.UserHomeDir()
	internalErrors.HandleFatal("Failed to get user home directory: ", err)

	windowsRootPath := filepath.Join(homeDir, WinFolder)
	windowsAppPath := filepath.Join(windowsRootPath, WinExecutable)
	windowsConfigPath := filepath.Join(windowsRootPath, "config", WinConfig)

	err = nssm.InitNssm(windowsRootPath)
	internalErrors.HandleFatal("", err)

	// Stop and remove the service if it already exists
	StopAndRemoveService(WinService, nssm.GetNssmPath())

	util.CopyBinary(windowsAppPath, WinExecutable)

	// Create the config file that will be used for the service(Based on the settings) and loads it  to /opt/aw/ path
	err = settings.CreateConfigFile(config, windowsConfigPath)
	for _, plugin := range plugins {
		plugin.ReplicateConfig(homeDir + "/.config/aw/config")
	}

	internalErrors.HandleFatal("Failed to create config file: ", err)

	builder := nssm.NewServiceBuilder()

	service := builder.ServiceName(WinService).
		AppDirectory(windowsRootPath).
		DisplayName("ActivityWatch Sync Agent").
		Application(windowsAppPath).
		Start("SERVICE_AUTO_START").
		Build()

	err = nssm.CreateService(service)
	internalErrors.HandleFatal("Failed to create service: ", err)

	err = nssm.StartService(WinService)
	internalErrors.HandleFatal("Failed to start service: ", err)

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
