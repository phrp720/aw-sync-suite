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

func CreateLinuxService(sett settings.Settings) {

	// Copy the agent executable to the specified path
	util.CopyBinary(appPath, binaryName)

	// Create the config file that will be used for the service(Based on the settings)
	err := settings.CreateConfigFile(sett, configPath)
	if err != nil {
		log.Fatal("Failed to create config file")
	}

	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}

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

	err = systemd.GenerateAndStart(service, serviceName, true)

	if err != nil {
		log.Fatal(err)
	}
	log.Print("Running as a service In Linux...")

}
