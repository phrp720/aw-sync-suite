package filter

import (
	"fmt"
	"regexp"
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

func ValidateFilters(filters []Filter) []Filter {
	validFilters := []Filter{}
	for _, filter := range filters {
		if (filter.Replace_app_id != "" || filter.Replace_app_name != "") &&
			(filter.Match_app_id == nil && filter.Match_app_name == nil) {
			// Ignore this filter as it does not meet the condition
			continue
		}
		validFilters = append(validFilters, filter)
	}
	return validFilters
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
func PrintFilters(filters []Filter) {
	for i, filter := range filters {
		fmt.Printf("Filter %d:\n", i+1)
		fmt.Printf("  Match_app_id: %s\n", filter.Match_app_id.String())
		fmt.Printf("  Match_app_name: %s\n", filter.Match_app_name.String())
		fmt.Printf("  Replace_app_id: %s\n", filter.Replace_app_id)
		fmt.Printf("  Replace_app_name: %s\n", filter.Replace_app_name)
	}
}
