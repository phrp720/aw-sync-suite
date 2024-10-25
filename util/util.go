package util

import (
	"github.com/robfig/cron"
	"log"
)

func ValidateCronExpr(cronExpr string) string {
	_, err := cron.ParseStandard(cronExpr)
	if err != nil {
		log.Fatalf("Invalid cron expression: %v", err)
	}
	return cronExpr
}
