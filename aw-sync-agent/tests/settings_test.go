package tests

import (
	"aw-sync-agent/settings"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	config := settings.LoadYAMLConfig("test_config.yaml")

	expectedConfig := settings.Configuration{
		Settings: settings.Setts{
			AWUrl:               "http://localhost:5600",
			PrometheusUrl:       "http://localhost:80",
			ExcludedWatchers:    []string{"watcher1", "watcher2"},
			UserID:              "Tester",
			PrometheusSecretKey: "secretKey",
			IncludeHostname:     true,
		},
	}

	if config.Settings.AWUrl != expectedConfig.Settings.AWUrl {
		t.Errorf("expected AWUrl to be %s, got %s", expectedConfig.Settings.AWUrl, config.Settings.AWUrl)
	}

	if config.Settings.PrometheusUrl != expectedConfig.Settings.PrometheusUrl {
		t.Errorf("expected PrometheusUrl to be %s, got %s", expectedConfig.Settings.PrometheusUrl, config.Settings.PrometheusUrl)
	}

	if config.Settings.UserID != expectedConfig.Settings.UserID {
		t.Errorf("expected UserID to be %s, got %s", expectedConfig.Settings.UserID, config.Settings.UserID)
	}

	if config.Settings.PrometheusSecretKey != expectedConfig.Settings.PrometheusSecretKey {
		t.Errorf("expected PrometheusSecretKey to be %s, got %s", expectedConfig.Settings.PrometheusSecretKey, config.Settings.PrometheusSecretKey)
	}

	if config.Settings.IncludeHostname != expectedConfig.Settings.IncludeHostname {
		t.Errorf("expected IncludeHostname to be %v, got %v", expectedConfig.Settings.IncludeHostname, config.Settings.IncludeHostname)
	}
}
