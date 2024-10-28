package util

import (
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
	if err != nil {
		log.Fatalf("Invalid cron expression: %v", err)
	}
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
	if err != nil {
		log.Fatalf("Failed to create destination binary: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		log.Fatalf("Failed to copy binary: %v", err)
	}

	// Set the executable permissions for everyone
	if err := os.Chmod(appPath, 0755); err != nil {
		log.Fatalf("Failed to set executable permissions: %v", err)
	}
}
