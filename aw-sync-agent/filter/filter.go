package filter

import (
	"aw-sync-agent/util"
	"fmt"
	"log"
	"regexp"
)

// Filter struct
type Filter struct {
	FilterName   string         `yaml:"filter-name"`   // FilterName is the name of the filter
	Watchers     []string       `yaml:"watchers"`      // Watchers is the list of watchers to be filtered
	Target       []Target       `yaml:"target"`        // Target is the key-value pair to be matched
	PlainReplace []PlainReplace `yaml:"plain-replace"` // Replace is the key-value pair to be replaced
	RegexReplace []RegexReplace `yaml:"regex-replace"` // Replace is the key-value pair to be replaced
	Enabled      bool           `yaml:"enabled"`       // Enabled is the flag to enable or disable the filter
	Drop         bool           `yaml:"drop"`          // Drop is the flag to drop the event if the filter matches
}

type Target struct {
	Key   string         `yaml:"key"`
	Value *regexp.Regexp `yaml:"value"`
}

type PlainReplace struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type RegexReplace struct {
	Key        string         `yaml:"key"`
	Expression *regexp.Regexp `yaml:"from"`
	Value      string         `yaml:"value"`
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
			for _, replace := range filter.PlainReplace {
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
		for j, replace := range filter.PlainReplace {
			fmt.Printf("  Plain String Replace %d:\n", j+1)
			fmt.Printf("    Key: %s\n", replace.Key)
			fmt.Printf("    Value: %s\n", replace.Value)
		}
		for m, replace := range filter.RegexReplace {
			fmt.Printf("  Regex Value Replace %d:\n", m+1)
			fmt.Printf("    Key: %s\n", replace.Key)
			fmt.Printf("    Value: %s\n", replace.Value)
			fmt.Printf("    Expression: %s\n", replace.Expression)
		}
	}
}
