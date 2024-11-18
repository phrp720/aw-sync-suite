package aw

import "time"

// Watcher represents a bucket in the aw database
type Watcher struct {
	ID          string                 `json:"id"`
	Created     time.Time              `json:"created"`
	Name        *string                `json:"name"`
	Type        string                 `json:"type"`
	Client      string                 `json:"client"`
	Hostname    string                 `json:"hostname"`
	Data        map[string]interface{} `json:"data"`
	LastUpdated time.Time              `json:"last_updated"`
}

type Watchers map[string]Watcher

// Event represents an event in the aw database
type Event struct {
	ID        int                    `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Duration  float64                `json:"duration"`
	Data      map[string]interface{} `json:"data"`
}

type Events []Event

type WatcherNameToEventsMap map[string]Events
