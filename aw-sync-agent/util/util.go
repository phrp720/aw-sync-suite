package util

import (
	"aw-sync-agent/errors"
	"github.com/google/uuid"

	"github.com/robfig/cron"
	"io"
	"log"
	"os"
	"path/filepath"

	"runtime"
)

// ValidateCronExpr validates the cron expression
func ValidateCronExpr(cronExpr string) string {
	_, err := cron.ParseStandard(cronExpr)
	errors.HandleFatal("Invalid cron expression: ", err)
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
	errors.HandleFatal("Failed to create destination binary: ", err)
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
