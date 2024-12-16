package tests

import (
	"aw-sync-agent/cron"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	c := cron.Init()
	if c == nil {
		t.Fatal("expected cron instance, got nil")
	}
}

func TestAdd(t *testing.T) {
	c := cron.Init()
	cron.Add(c, "* * * * *", func() {})
	if len(c.Entries()) != 1 {
		t.Fatalf("expected 1 cron entry, got %d", len(c.Entries()))
	}
}

func TestStartAndStop(t *testing.T) {
	c := cron.Init()
	cron.Add(c, "@every 1s", func() {})
	cron.Start(c)
	time.Sleep(2 * time.Second)
	cron.Stop(c)
	if len(c.Entries()) != 1 {
		t.Fatalf("expected 1 cron entry, got %d", len(c.Entries()))
	}
}
