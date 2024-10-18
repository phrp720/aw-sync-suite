package main

import (
	"aw-sync-agent/settings"
	"aw-sync-agent/synchronizer"
)

func main() {
	// These are the settings that contains the Env Variables/Flags
	Settings := map[settings.SettingsKey]string{
		settings.AWUrl:            settings.GetFlagOrEnvVar("ACTIVITY_WATCH_URL", "awUrl", true),
		settings.PrometheusUrl:    settings.GetFlagOrEnvVar("PROMETHEUS_URL", "prometheusUrl", true),
		settings.ExcludedWatchers: settings.GetFlagOrEnvVar("EXCLUDED_WATCHERS", "excludedWatchers", false),
		settings.UserID:           settings.GetFlagOrEnvVar("USER_ID", "userID", false),
		settings.Cron:             settings.GetFlagOrEnvVar("CRON", "cron", false),
	}

	// Pass the map to the synchronizer.Start function
	err := synchronizer.Start(Settings)
	if err != nil {
		panic(err) // handle if something wrong happens
	}
}
