package filter

import (
	"aw-sync-agent/util"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Filter struct
type Filter struct {
	FilterName   string         `yaml:"filter-name"`   // FilterName is the name of the filter
	Watchers     []string       `yaml:"watchers"`      // Watchers is the list of watchers to be filtered
	Target       []Target       `yaml:"target"`        // Target is the key-value pair to be matched
	PlainReplace []PlainReplace `yaml:"plain-replace"` // Replace is the key-value pair to be replaced
	RegexReplace []RegexReplace `yaml:"regex-replace"` // Replace is the key-value pair to be replaced
	Enable       bool           `yaml:"enable"`        // Enabled is the flag to enable or disable the filter
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
func ValidateFilters(filters []Filter) ([]Filter, int, int, int) {
	validFilters := []Filter{}
	invalid := 0
	disabled := 0
	total := len(filters)
	for _, filter := range filters {
		if filter.Enable {
			var targetList []Target
			for _, target := range filter.Target {
				if strings.TrimSpace(target.Key) != "" && target.Value != nil {
					targetList = append(targetList, target)
				}
			}
			if len(targetList) != 0 { // Check if the filter has at least one valid target
				validFilters = append(validFilters, filter)
			} else {
				invalid++
			}
		} else {
			disabled++
		}
	}
	return validFilters, total, invalid, disabled
}

// Apply applies the filters to the data
func Apply(data map[string]interface{}, filters []Filter) (map[string]interface{}, bool) {
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
			// Drop the event if the filter matches
			if filter.Drop {
				return nil, true
			}
			// Apply replacements
			data = Replace(data, filter.PlainReplace, filter.RegexReplace)

		}
	}
	return data, false
}

// Replace replaces the values in the data
func Replace(data map[string]interface{}, plain []PlainReplace, regex []RegexReplace) map[string]interface{} {

	// Apply replacements

	// Plain replacements
	for _, replace := range plain {
		if _, exists := data[replace.Key]; exists {
			data[replace.Key] = replace.Value
		}
	}

	// Regex replacements
	for _, replace := range regex {
		if value, exists := data[replace.Key]; exists {
			// Check if the value matches the target's regex
			if replace.Expression.MatchString(fmt.Sprintf("%v", value)) {
				// Replace the value with the formatted string
				data[replace.Key] = replace.Expression.ReplaceAllString(fmt.Sprintf("%v", value), replace.Value)
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
		log.Print("Filters applied for [", watcher, "]: ", len(matchingFilters))
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
