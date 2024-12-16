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
	// 0 is the seconds.
	//This happens to remove the SECONDS options from the cron expression(default in github.com/robfig/cron ) so it conforms to the standard. https://en.wikipedia.org/wiki/Cron
	err := c.AddFunc("0 "+scheduler, fun)
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
