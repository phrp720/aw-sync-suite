package checkpoint

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var checkpointFile = "checkpoint.json"

// Read reads the checkpoint for the given watcher
func Read(watcher string) *time.Time {

	// Read the existing data from the file | If the file does not exist, creates it
	file, err := os.OpenFile(checkpointFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error: Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		// If the file is empty or invalid, initialize an empty map
		checkpoints = make(map[string]string)
	}

	// Get the checkpoint for the given watcher
	timestampStr, ok := checkpoints[watcher]
	if !ok {
		return nil
	}
	timestamp, err := time.Parse(time.RFC3339Nano, timestampStr)
	if err != nil {
		log.Fatalf("Error: Failed to parse timestamp: %v", err)
	}
	return &timestamp
}

// Update updates the checkpoint for the given watcher
func Update(watcher string, timestamp time.Time) {

	// Read the existing data from the file | If the file does not exist, creates it
	file, err := os.OpenFile(checkpointFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error: Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		// If the file is empty or invalid, initialize an empty map
		checkpoints = make(map[string]string)
	}

	// Update the checkpoint for the given watcher
	checkpoints[watcher] = timestamp.Format("2006-01-02T15:04:05.000000-07:00")

	// Write the updated data back to the file
	file.Seek(0, 0) // Move the file pointer to the beginning
	encoder := json.NewEncoder(file)
	if err = encoder.Encode(checkpoints); err != nil {
		log.Fatalf("Error: Failed to write to checkpoint file: %v", err)
	}

	// Truncate the file to the current size (in case the new content is smaller)
	tmpFile, err := file.Seek(0, 1)
	if err != nil {
		log.Fatalf("Error: Failed to get the current size of the checkpoint file: %v", err)
	}
	if err = file.Truncate(tmpFile); err != nil {
		log.Fatalf("Error: Failed to truncate checkpoint file: %v", err)
	}
}

// GetCheckpointFile returns the checkpoint file name
func GetCheckpointFile() string {
	return checkpointFile
}

// SetCheckpointFile sets the checkpoint file name
func SetCheckpointFile(file string) {
	checkpointFile = file
}
