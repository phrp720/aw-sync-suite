package service

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/util"
	"github.com/phrp720/service-builder/systemd"
	"log"
	"os/user"
	"path/filepath"
)

const (
	configFile  = "config.yaml"
	binaryName  = "aw-sync-agent"
	serviceName = "aw-sync-agent.service"

	configPath = "/opt/aw/" + configFile
	appPath    = "/opt/aw/" + binaryName
)

// CreateLinuxService creates a Linux service using the service-builder library github.com/phrp720/service-builder
func CreateLinuxService(sett settings.Settings) {

	// Copies  the aw-sync-agent executable to /opt/aw/ path
	util.CopyBinary(appPath, binaryName)

	// Create the config file that will be used for the service(Based on the settings) and loads it  to /opt/aw/ path
	err := settings.CreateConfigFile(sett, configPath)
	if err != nil {
		log.Fatal("Failed to create config file")
	}

	// Get the current user
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

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
		WantedBy("multi-user.target").
		Build()

	// Generate enables and starts the service
	err = systemd.GenerateAndStart(service, serviceName, true)

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Running as a Linux service...")

}
