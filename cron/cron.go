package cron

import (
	"github.com/robfig/cron"
	"log"
)

func InitCron() *cron.Cron {
	c := cron.New()
	return c
}

func Add(c *cron.Cron, scheduler string, fun func()) {
	err := c.AddFunc(scheduler, fun)
	if err != nil {
		log.Fatal(err)
	}
}

func Start(c *cron.Cron) {
	c.Start()
}

func Stop(c *cron.Cron) {
	c.Stop()
}
