package service

import (
	"aw-sync-agent/aw-sync-agent/settings"
	"aw-sync-agent/aw-sync-agent/system_error"
	"aw-sync-agent/aw-sync-agent/util"
	"github.com/phrp720/go-service-builder/systemd"
	"log"
	"os/user"
	"path/filepath"
)

const (
	LinuxConfig     = "aw-sync-agent.yaml"
	LinuxExecutable = "aw-sync-agent"
	LinuxService    = "aw-sync-agent.service"
)

// CreateLinuxService creates a Linux service using the service-builder library github.com/phrp720/service-builder
func CreateLinuxService(config settings.Configuration) {
	// Get the current user
	currentUser, err := user.Current()
	system_error.HandleFatal("Failed to get current user: ", err)

	// Define paths dynamically based on the user's home directory
	homeDir := currentUser.HomeDir
	configPath := filepath.Join(homeDir, ".config", "aw", LinuxConfig)
	appPath := filepath.Join(homeDir, ".config", "aw", LinuxExecutable)
	serviceFilePath := filepath.Join(homeDir, ".config", "systemd", "user", LinuxService)

	// Copies the aw-sync-agent executable to the user's config path
	util.CopyBinary(appPath, LinuxExecutable)

	// Create the config file that will be used for the service (Based on the settings) and loads it to the user's config path
	err = settings.CreateConfigFile(config, configPath)
	system_error.HandleFatal("Failed to create config file: ", err)

	// Get the working directory
	workingDirectory := filepath.Dir(appPath)

	builder := systemd.NewServiceBuilder()
	service := builder.
		// Unit
		Description("ActivityWatch Sync Agent").
		After("network.target").
		// Service
		ExecStart(appPath).
		Restart("always").
		User(currentUser.Username).
		Group(currentUser.Username).
		WorkingDirectory(workingDirectory).
		RestartSec("5").
		// Install
		WantedBy("default.target").
		Build()

	// Generate, enable, and start the service
	err = systemd.CreateService(service, serviceFilePath)
	system_error.HandleFatal("Failed to create service: ", err)

	err = systemd.StartService(LinuxExecutable, false)
	system_error.HandleFatal("Failed to start service: ", err)

	log.Print("Running as a Linux service...")
}
