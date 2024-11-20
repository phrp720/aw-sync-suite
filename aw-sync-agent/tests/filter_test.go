package tests

import (
	"aw-sync-agent/filter"
	"regexp"
	"testing"
)

// TestValidateFilters tests the ValidateFilters function
func TestValidateFilters(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "CorrectFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			Enable: true,
		},
		{
			FilterName: "DisabledFilter",
			Target: []filter.Target{
				{Key: "key2", Value: regexp.MustCompile("value2")},
			},
			Enable: false,
		},
		{
			FilterName: "InvalidFilter",
			Target: []filter.Target{
				{Key: "", Value: regexp.MustCompile("value2")},
			},
			Enable: false,
		},
	}

	validFilters, total, invalid, disabled := filter.ValidateFilters(filters)
	if len(validFilters) != 1 || total != 3 || invalid != 1 || disabled != 1 {
		t.Errorf("expected 1 valid filter, 3 total, 1 invalid, 1 disabled; got %d valid, %d total, %d invalid, %d disabled", len(validFilters), total, invalid, disabled)
	}
}

// TestGetMatchingFilters tests the GetMatchingFilters function
func TestGetMatchingFilters(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "Filter_1",
			Watchers:   []string{"watcher1"},
			Enable:     true,
		},
		{
			FilterName: "Filter_2",
			Watchers:   []string{"watcher2"},
			Enable:     true,
		},
	}

	matchingFilters := filter.GetMatchingFilters(filters, "watcher1")
	if len(matchingFilters) != 1 || matchingFilters[0].FilterName != "Filter_1" {
		t.Errorf("expected 1 matching filter with name 'Filter_1', got %d matching filters with name '%s'", len(matchingFilters), matchingFilters[0].FilterName)
	}
}

// TestApplyWithDrop tests the Apply function with a filter that should drop the data
func TestApplyWithDrop(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "DropFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			Drop:   true,
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	result, dropped := filter.Apply(data, filters)
	if result != nil || !dropped {
		t.Errorf("expected data to be dropped, got %v, dropped: %v", result, dropped)
	}
}

// TestApplyWithDrop tests the Apply function with a filter that should not drop the data
func TestApplyWithPlainReplace(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "PlainReplaceFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			PlainReplace: []filter.PlainReplace{
				{Key: "key1", Value: "newValue1"},
			},
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	expected := map[string]interface{}{"key1": "newValue1"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["key1"] != expected["key1"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithRegexReplace tests the Apply function with a regex replace filter
func TestApplyWithRegexReplace(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "RegexReplaceFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			RegexReplace: []filter.RegexReplace{
				{Key: "key1", Expression: regexp.MustCompile("val"), Value: "newValue1"},
			},
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	expected := map[string]interface{}{"key1": "newValue1"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["key1"] != expected["key1"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}
