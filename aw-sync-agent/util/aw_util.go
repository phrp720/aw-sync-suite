package util

import (
	"aw-sync-agent/aw"
	"log"
)

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
