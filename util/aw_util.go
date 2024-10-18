package util

import (
	"aw-sync-agent/aw"
	"log"
	"strings"
)

// getExcludedWatchers gets the excluded watchers from the environment variable
func getExcludedWatchers() []string {
	excluded, _ := GetEnvVar("EXCLUDED_WATCHERS", false)
	return strings.Split(excluded, "|")
}

// RemoveExcludedWatchers removes the excluded watchers from the buckets
func RemoveExcludedWatchers(buckets aw.Watchers) aw.Watchers {
	excluded := getExcludedWatchers()
	if len(excluded) > 0 {
		for _, excludedWatcher := range excluded {
			for id, bucket := range buckets {
				if bucket.Client == excludedWatcher {
					delete(buckets, id)
				}
			}
		}
	}
	log.Print("Buckets excluded: ", len(excluded), excluded, "\n")
	return buckets
}
