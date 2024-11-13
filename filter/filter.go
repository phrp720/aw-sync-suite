package filter

import (
	"aw-sync-agent/util"
	"fmt"
	"log"
	"regexp"
)

// Filter struct
type Filter struct {
	Watchers []string  `yaml:"watchers"` // Watchers is the list of watchers to be filtered
	Target   []Target  `yaml:"target"`   // Target is the key-value pair to be matched
	Replace  []Replace `yaml:"replace"`  // Replace is the key-value pair to be replaced
}

type Target struct {
	Key   string         `yaml:"key"`
	Value *regexp.Regexp `yaml:"value"`
}

type Replace struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// ValidateFilters validates the filters in the List
func ValidateFilters(filters []Filter) ([]Filter, int, int) {
	validFilters := []Filter{}
	var targetList []Target
	invalid := 0
	total := len(filters)
	for _, filter := range filters {
		for _, target := range filter.Target {
			if target.Key != "" && target.Value != nil {
				targetList = append(targetList, target)
			}
		}
		if len(targetList) != 0 { // Check if the filter has at least one valid target
			validFilters = append(validFilters, filter)
		} else {
			invalid++
		}
	}
	return validFilters, total, invalid
}

func Apply(data map[string]interface{}, filters []Filter) map[string]interface{} {
	for _, filter := range filters {
		allMatch := true
		for _, target := range filter.Target {
			// Check if the data contains the key to be filtered
			if value, ok := data[target.Key]; ok {
				// Check if the value matches the target's regex
				if !target.Value.MatchString(fmt.Sprintf("%v", value)) {
					allMatch = false
					break
				}
			} else {
				allMatch = false
				break
			}
		}

		if allMatch {
			// Apply replacements
			for _, replace := range filter.Replace {
				if _, exists := data[replace.Key]; exists {
					data[replace.Key] = replace.Value
				}
			}
		}
	}
	return data
}

// GetMatchingFilters returns filters that match the given watcher or have an empty watcher list
func GetMatchingFilters(filters []Filter, watcher string) []Filter {
	var matchingFilters []Filter
	for _, filter := range filters {
		if len(filter.Watchers) == 0 || util.Contains(filter.Watchers, watcher) {
			matchingFilters = append(matchingFilters, filter)
		}
	}
	if len(matchingFilters) > 0 {
		log.Print(watcher, " filters applied : ", len(matchingFilters))
	}

	return matchingFilters
}

// PrintFilters prints the filters in the List
func PrintFilters(filters []Filter) {
	for i, filter := range filters {
		fmt.Printf("Filter %d:\n", i+1)
		fmt.Printf("  Watchers: %v\n", filter.Watchers)
		for k, target := range filter.Target {
			fmt.Printf("  Target %d:\n", k+1)
			fmt.Printf("    Key: %s\n", target.Key)
			fmt.Printf("    Value: %s\n", target.Value)
		}
		for j, replace := range filter.Replace {
			fmt.Printf("  Replace %d:\n", j+1)
			fmt.Printf("    Key: %s\n", replace.Key)
			fmt.Printf("    Value: %s\n", replace.Value)
		}
	}
}
