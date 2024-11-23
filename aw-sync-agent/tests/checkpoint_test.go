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
	if !readTimestamp.Equal(timestamp) {
		t.Errorf("expected %v, got %v", timestamp, readTimestamp)
	}
}
