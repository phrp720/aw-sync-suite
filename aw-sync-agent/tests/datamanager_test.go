package tests

import (
	"aw-sync-agent/activitywatch"
	"aw-sync-agent/datamanager"
	"aw-sync-agent/prometheus"
	"testing"
	"time"
)

func TestAggregateData(t *testing.T) {
	events := []activitywatch.Event{
		{
			Timestamp: time.Now().Add(-1 * time.Hour),
			Duration:  30,
			Data:      map[string]interface{}{"app": "Meeting App", "title": "Meeting with team"},
		},
		{
			Timestamp: time.Now().Add(-3 * time.Hour),
			Duration:  45,
			Data:      map[string]interface{}{"app": "code", "title": "Coding awesome stuff"},
		},
	}

	watcher := "test-watcher"
	userID := "test-user"

	timeSeries := datamanager.AggregateData(nil, events, watcher, userID, false)

	if len(timeSeries) != 1 {
		t.Errorf("Expected 1 time series, got %d", len(timeSeries))
	}

	expectedLabels := []prometheus.Label{
		{Name: "__name__", Value: "test_watcher"},
		{Name: "user", Value: "test-user"},
		{Name: "app", Value: "code"},
		{Name: "title", Value: "Coding awesome stuff"},
	}

	for _, label := range expectedLabels {
		for _, l := range timeSeries[0].Labels {
			if l.Name == label.Name {
				if l.Value != label.Value {
					t.Errorf("Expected label %v, got %v", label.Value, l.Value)
				}
				break
			}
		}
	}

	if timeSeries[0].Sample.Value != 45 {
		t.Errorf("Expected sample value 45, got %v", timeSeries[0].Sample.Value)
	}
}
