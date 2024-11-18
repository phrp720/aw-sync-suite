package checkpoint

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const checkpointFile = "checkpoint.json"

// Read reads the checkpoint for the given watcher
func Read(watcher string) *time.Time {

	// Read the existing data from the file | If the file does not exist, creates it
	file, err := os.OpenFile(checkpointFile, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]time.Time
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		// If the file is empty or invalid, initialize an empty map
		checkpoints = make(map[string]time.Time)
	}

	// Get the checkpoint for the given watcher
	timestamp, ok := checkpoints[watcher]
	if !ok {
		return nil
	}

	return &timestamp
}

// Update updates the checkpoint for the given watcher
func Update(watcher string, timestamp time.Time) {

	// Read the existing data from the file | If the file does not exist, creates it
	file, err := os.OpenFile(checkpointFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Failed to open checkpoint file: %v", err)
	}
	defer file.Close()

	var checkpoints map[string]time.Time
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&checkpoints); err != nil {
		// If the file is empty or invalid, initialize an empty map
		checkpoints = make(map[string]time.Time)
	}

	// Update the checkpoint for the given watcher
	checkpoints[watcher] = timestamp

	// Write the updated data back to the file
	file.Seek(0, 0) // Move the file pointer to the beginning
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(checkpoints); err != nil {
		log.Fatalf("Failed to write to checkpoint file: %v", err)
	}

	// Truncate the file to the current size (in case the new content is smaller)
	tmpFile, err := file.Seek(0, 1)
	if err != nil {
		log.Fatalf("Failed to get the current size of the checkpoint file: %v", err)
	}
	if err = file.Truncate(tmpFile); err != nil {
		log.Fatalf("Failed to truncate checkpoint file: %v", err)
	}
}
