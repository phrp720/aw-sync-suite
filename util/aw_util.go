package util

import (
	"aw-sync-agent/aw"
	"log"
	"strings"
)

// getExcludedWatchers gets the excluded watchers from the environment variable
func getExcludedWatchers(excludedWatchers string) []string {
	excluded := strings.Split(excludedWatchers, "|")
	if len(excluded) == 1 && excluded[0] == "" {
		return []string{}
	}
	return excluded
}

// RemoveExcludedWatchers removes the excluded watchers from the buckets
func RemoveExcludedWatchers(buckets aw.Watchers, excludedWatchers []string) aw.Watchers {
	if len(excludedWatchers) > 0 {
		for _, excludedWatcher := range excludedWatchers {
			for id, bucket := range buckets {
				if bucket.Client == excludedWatcher {
					delete(buckets, id)
				}
			}
		}
	}
	log.Print("Buckets excluded: ", len(excludedWatchers), excludedWatchers, "\n")
	return buckets
}
