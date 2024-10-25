package util

import (
	"github.com/robfig/cron"
	"log"
	"runtime"
)

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
