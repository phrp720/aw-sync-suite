package service

import (
	internalErrors "aw-sync-agent/errors"
	"aw-sync-agent/settings"
	"aw-sync-agent/util"
	"github.com/phrp720/go-service-builder/systemd"
	"log"
	"os/user"
	"path/filepath"
)

const (
	LinuxConfig     = "aw-sync-settings.yaml"
	LinuxExecutable = "aw-sync-agent"
	LinuxService    = "aw-sync-agent.service"
)

// CreateLinuxService creates a Linux service using the service-builder library github.com/phrp720/service-builder
func CreateLinuxService(config settings.Configuration) {
	// Get the current user
	currentUser, err := user.Current()
	internalErrors.HandleFatal("Failed to get current user: ", err)

	// Define paths dynamically based on the user's home directory
	homeDir := currentUser.HomeDir
	configPath := filepath.Join(homeDir, ".config", "aw", "config", LinuxConfig)
	appPath := filepath.Join(homeDir, ".config", "aw", LinuxExecutable)
	serviceFilePath := filepath.Join(homeDir, ".config", "systemd", "user", LinuxService)

	// Stop the old service if it is running
	systemd.StopService(LinuxExecutable, false)

	// Copies the aw-sync-agent executable to the user's config path
	util.CopyBinary(appPath, LinuxExecutable)

	// Create the config file that will be used for the service (Based on the settings) and loads it to the user's config path
	err = settings.CreateConfigFile(config, configPath)
	internalErrors.HandleFatal("Failed to create config file: ", err)

	// Get the working directory
	workingDirectory := filepath.Dir(appPath)

	builder := systemd.NewServiceBuilder()
	//Note https://unix.stackexchange.com/questions/438064/failed-to-determine-supplementary-groups-operation-not-permitted
	service := builder.
		// Unit
		Description("ActivityWatch Sync Agent").
		After("network.target").
		// Service
		ExecStart(appPath).
		Restart("always").
		WorkingDirectory(workingDirectory).
		RestartSec("5").
		// Install
		WantedBy("default.target").
		Build()

	// Generate, enable, and start the service
	err = systemd.CreateService(service, serviceFilePath)
	internalErrors.HandleFatal("Failed to create service: ", err)

	err = systemd.StartService(LinuxExecutable, false)
	internalErrors.HandleFatal("Failed to start service: ", err)

	log.Print("Running as a Linux service...")
}
