package filter

import (
	"aw-sync-agent/system_error"
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	match_app_id     = "match_app_id"
	match_app_name   = "match_app_name"
	replace_app_id   = "replace_app_id"
	replace_app_name = "replace_app_name"
)

// Filter struct
type Filter struct {
	Match_app_id     *regexp.Regexp
	Match_app_name   *regexp.Regexp
	Replace_app_id   string
	Replace_app_name string
}

// List of filters
type List struct {
	Filters []Filter
}

// LoadFilters loads filters from the specified file and returns a List
func LoadFilters(filename string) *List {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open filters file: %v", err)
	}
	defer file.Close()

	var filters []Filter
	scanner := bufio.NewScanner(file)
	var currentFilter Filter

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "[Filter]" {
			// Add the previous filter if it has a valid match field
			if currentFilter.Match_app_id != nil || currentFilter.Match_app_name != nil {
				filters = append(filters, currentFilter)
			}
			currentFilter = Filter{}
			continue
		}

		// Split key-value pairs by "="
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

		// Parse values based on the key
		switch key {
		case match_app_id:
			re, err := regexp.Compile(value)
			system_error.HandleFatal("invalid regex for match_app_id: ", err)
			currentFilter.Match_app_id = re

		case match_app_name:
			re, err := regexp.Compile(value)
			system_error.HandleFatal("invalid regex for match_app_name: ", err)
			currentFilter.Match_app_name = re

		case replace_app_id:
			if currentFilter.Match_app_id != nil || currentFilter.Match_app_name != nil {
				currentFilter.Replace_app_id = value
			}

		case replace_app_name:
			if currentFilter.Match_app_id != nil || currentFilter.Match_app_name != nil {
				currentFilter.Replace_app_name = value
			}
		}
	}

	// Add the last filter if it has a valid match field
	if currentFilter.Match_app_id != nil || currentFilter.Match_app_name != nil {
		filters = append(filters, currentFilter)
	}

	if err := scanner.Err(); err != nil {
		system_error.HandleFatal("Failed to read filters file: ", err)
	}

	return &List{Filters: filters}
}

// Apply applies the filters in the  to the app_id and app_name
func (f *Filter) Apply(app_id string, app_name string) (string, string) {
	if f.Match_app_id.MatchString(app_id) {
		app_id = f.Replace_app_id
	}
	if f.Match_app_name.MatchString(app_name) {
		app_name = f.Replace_app_name
	}
	return app_id, app_name
}

// Apply applies the filters in the List to the app_id and app_name
func (l *List) Apply(app_id string, app_name string) (string, string) {
	for _, f := range l.Filters {
		app_id, app_name = f.Apply(app_id, app_name)
	}
	return app_id, app_name
}

// PrintFilters prints the filters in the List
func PrintFilters(list List) {
	for i, filter := range list.Filters {
		fmt.Printf("Filter %d:\n", i+1)
		fmt.Printf("  Match_app_id: %s\n", filter.Match_app_id.String())
		fmt.Printf("  Match_app_name: %s\n", filter.Match_app_name.String())
		fmt.Printf("  Replace_app_id: %s\n", filter.Replace_app_id)
		fmt.Printf("  Replace_app_name: %s\n", filter.Replace_app_name)
	}
}
