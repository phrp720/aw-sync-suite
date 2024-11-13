package filter

import (
	"fmt"
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
