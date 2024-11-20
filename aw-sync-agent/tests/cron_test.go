package tests

import (
	"aw-sync-agent/cron"
	"aw-sync-agent/util"
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
	cron.Add(c, "@every 1s", func() {})
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

func TestValidateCronExpr(t *testing.T) {
	validExpr := "@every 1s"
	invalidExpr := "invalid"

	if util.ValidateCronExpr(validExpr) != validExpr {
		t.Fatalf("expected valid cron expression, got %s", validExpr)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic for invalid cron expression")
		}
	}()
	util.ValidateCronExpr(invalidExpr)
}
