package filter

import (
	"aw-sync-agent/util"
	"fmt"
	"log"
	"regexp"
)

// Filter struct
type Filter struct {
	Watchers []string       `yaml:"watchers"`
	Key      string         `yaml:"key"`
	Value    *regexp.Regexp `yaml:"value"`
	Replace  []Replace      `yaml:"replace"`
}

type Replace struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// ValidateFilters validates the filters in the List
func ValidateFilters(filters []Filter) ([]Filter, int, int) {
	validFilters := []Filter{}
	invalid := 0
	total := len(filters)
	for _, filter := range filters {
		if filter.Key == "" || filter.Value == nil {
			invalid++
			continue
		}
		validFilters = append(validFilters, filter)
	}
	return validFilters, total, invalid
}

func Apply(data map[string]interface{}, filters []Filter) map[string]interface{} {

	for _, filter := range filters {
		// Check if the data contains the key to be filtered
		if value, ok := data[filter.Key]; ok {
			// Check if the value matches the filter's regex
			if filter.Value.MatchString(fmt.Sprintf("%v", value)) {
				// Apply replacements
				for _, replace := range filter.Replace {
					if _, exists := data[replace.Key]; exists {
						data[replace.Key] = replace.Value
					}
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
	log.Print(watcher, " filters applied : ", len(matchingFilters))
	return matchingFilters
}

// PrintFilters prints the filters in the List
func PrintFilters(filters []Filter) {
	for i, filter := range filters {
		fmt.Printf("Filter %d:\n", i+1)
		fmt.Printf("  Watchers: %v\n", filter.Watchers)
		fmt.Printf("  Key: %s\n", filter.Key)
		fmt.Printf("  Value: %s\n", filter.Value)
		for j, replace := range filter.Replace {
			fmt.Printf("  Replace %d:\n", j+1)
			fmt.Printf("    Key: %s\n", replace.Key)
			fmt.Printf("    Value: %s\n", replace.Value)
		}
	}
}
