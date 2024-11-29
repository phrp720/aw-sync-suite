package tests

import (
	"aw-sync-agent/checkpoint"
	"os"
	"testing"
	"time"
)

const testCheckpointFile = "test_checkpoint.json"

func setup() {
	checkpoint.SetCheckpointFile(testCheckpointFile)
}

func teardown() {
	os.Remove(testCheckpointFile)
}

func TestReadAndUpdate(t *testing.T) {
	setup()
	defer teardown()

	watcher := "testWatcher"
	timestamp := time.Now()

	// Test Update
	checkpoint.Update(watcher, timestamp)

	// Test Read
	readTimestamp := checkpoint.Read(watcher)
	if readTimestamp == nil {
		t.Fatalf("expected timestamp, got nil")
	}
	if (readTimestamp.Format("2006-01-02T15:04:05.000000-07:00")) != (timestamp.Format("2006-01-02T15:04:05.000000-07:00")) {
		t.Errorf("expected %v, got %v", timestamp, readTimestamp)
	}
}
