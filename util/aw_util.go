package util

import (
	"aw-sync-agent/aw"
	"log"
	"strings"
)

// getExcludedWatchers gets the excluded watchers from the environment variable
func getExcludedWatchers(excludedWatchers string) []string {
	return strings.Split(excludedWatchers, "|")
}

// RemoveExcludedWatchers removes the excluded watchers from the buckets
func RemoveExcludedWatchers(buckets aw.Watchers, excludedWatchers string) aw.Watchers {
	excluded := getExcludedWatchers(excludedWatchers)
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
