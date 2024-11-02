package cron

import (
	"github.com/robfig/cron"
	"log"
)

// Init initializes a new cron job
func Init() *cron.Cron {
	c := cron.New()
	return c
}

// Add adds a new function to the cron job
func Add(c *cron.Cron, scheduler string, fun func()) {
	err := c.AddFunc(scheduler, fun)
	if err != nil {
		log.Fatal(err)
	}
}

// Start starts the cron job
func Start(c *cron.Cron) {
	c.Start()
}

// Stop stops the cron job
func Stop(c *cron.Cron) {
	c.Stop()
}
