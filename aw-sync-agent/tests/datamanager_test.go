package tests

import (
	"aw-sync-agent/aw"
	"aw-sync-agent/datamanager"
	"aw-sync-agent/filter"
	"aw-sync-agent/prometheus"
	"testing"
	"time"
)

func TestAggregateData(t *testing.T) {
	events := []aw.Event{
		{
			Timestamp: time.Now().Add(-2 * time.Hour),
			Duration:  30,
			Data:      map[string]interface{}{"activity": "coding"},
		},
		{
			Timestamp: time.Now().Add(-1 * time.Hour),
			Duration:  45,
			Data:      map[string]interface{}{"activity": "meeting"},
		},
	}

	watcher := "test-watcher"
	userID := "test-user"
	includeHostName := false
	filters := []filter.Filter{}

	timeSeries := datamanager.AggregateData(events, watcher, userID, includeHostName, filters)

	if len(timeSeries) != 1 {
		t.Errorf("Expected 1 time series, got %d", len(timeSeries))
	}

	expectedLabels := []prometheus.Label{
		{Name: "__name__", Value: "test_watcher"},
		{Name: "user", Value: "test-user"},
		{Name: "activity", Value: "coding"},
	}

	for i, label := range expectedLabels {
		if timeSeries[0].Labels[i] != label {
			t.Errorf("Expected label %v, got %v", label, timeSeries[0].Labels[i])
		}
	}

	if timeSeries[0].Sample.Value != 30 {
		t.Errorf("Expected sample value 30, got %v", timeSeries[0].Sample.Value)
	}
}
