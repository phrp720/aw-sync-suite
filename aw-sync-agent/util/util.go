package util

import (
	internalErrors "aw-sync-agent/errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/phrp720/aw-sync-agent-plugins/models"
	"github.com/robfig/cron"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// ValidateCronExpr validates the cron expression
func ValidateCronExpr(cronExpr string) string {
	_, err := cron.ParseStandard(cronExpr)
	internalErrors.HandleFatal("Invalid cron expression: ", err)
	return cronExpr
}

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsMac() bool {
	return runtime.GOOS == "darwin"
}

// CopyBinary copies the given binary to the specified path
func CopyBinary(appPath string, binaryName string) {
	// Create the directory if it doesn't exist
	dir := filepath.Dir(appPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal(err)
	}

	// Copy the agent binary to the specified path
	src, err := os.Open(binaryName)
	if err != nil {
		log.Fatalf("Failed to open source binary: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(appPath)
	internalErrors.HandleFatal("Failed to create destination binary: ", err)
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("Error: Failed to copy binary: %v", err)
	}

	// Set the executable permissions for everyone
	if err := os.Chmod(appPath, 0755); err != nil {
		log.Fatalf("Error: Failed to set executable permissions: %v", err)
	}
}

// GetUserID fetches or generates the user ID
func GetUserID(userID string) string {
	if userID != "" {
		return userID
	}

	hostname, err := os.Hostname()
	if err == nil {
		return hostname
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return "unknown"
	}

	return id.String()
}
func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// GetRandomUUID generates a random UUID
func GetRandomUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		return "unknown"
	}
	return id.String()
}

// CreateUniqueID creates a unique ID for an event using the event ID and a random UUID suffix
func CreateUniqueID(eventID string) string {
	if eventID == "" {
		return GetRandomUUID()
	}
	return eventID + "_" + GetRandomUUID()
}

// Contains checks if a slice contains a given string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PrintPlugins prints the categories in the List in a dashboard format
func PrintPlugins(Plugins []models.Plugin) {
	log.Print("Plugins:")
	pluginNames := make([]string, 0)
	for _, plugin := range Plugins {
		pluginNames = append(pluginNames, plugin.RawName())
	}
	filtersPluginsMap := map[string]string{
		"Plugins found": strconv.Itoa(len(Plugins)),
		"Plugins":       strings.Join(pluginNames, ", "),
	}

	maxKeyLength := 0
	maxValueLength := 0
	for key, value := range filtersPluginsMap {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
		if len(value) > maxValueLength {
			maxValueLength = len(value)
		}
	}

	borderLength := maxKeyLength + maxValueLength + 7
	border := strings.Repeat("-", borderLength)
	fmt.Println(border)
	for key, value := range filtersPluginsMap {
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLength, key, maxValueLength, value)
	}
	fmt.Println(border)
}
