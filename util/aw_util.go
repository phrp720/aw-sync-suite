package util

import (
	"aw-sync-agent/aw"
	"log"
	"strings"
)

func getExcludedWatchers() []string {
	excluded, _ := GetEnvVar("EXCLUDED_WATCHERS", false)
	return strings.Split(excluded, "|")
}

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
