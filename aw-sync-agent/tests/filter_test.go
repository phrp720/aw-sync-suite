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
			FilterName: "CategoryInvalidFilter",
			Target: []filter.Target{
				{Key: "key3", Value: regexp.MustCompile("value3")},
			},
			Category: "Email",
			Drop:     true,
			Enable:   true,
		},
		{
			FilterName: "CategoryValidFilter",
			Target: []filter.Target{
				{Key: "key3", Value: regexp.MustCompile("value3")},
			},
			Category: "Email",
			Enable:   true,
		},
		{
			FilterName: "InvalidFilter",
			Target: []filter.Target{
				{Key: "", Value: regexp.MustCompile("value2")},
			},
			Enable: true,
		},
	}

	validFilters, total, invalid, disabled := filter.ValidateFilters(filters)
	if len(validFilters) != 2 || total != 5 || invalid != 2 || disabled != 1 {
		t.Errorf("expected 4 total filters,1 valid, 1 invalid, 1 disabled; got %d valid, %d total, %d invalid, %d disabled", len(validFilters), total, invalid, disabled)
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

	data := map[string]interface{}{"key1": "value12"}
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
				{Key: "key1", Value: regexp.MustCompile("2value1")},
			},
			RegexReplace: []filter.RegexReplace{
				{Key: "key1", Expression: regexp.MustCompile("value1"), Value: "newValue1"},
			},
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "2value1"}
	expected := map[string]interface{}{"key1": "2newValue1"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["key1"] != expected["key1"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithPlainAndRegex tests the Apply function with both Plain and Regex Replace in one filter
func TestApplyWithPlainAndRegex(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "Filter with both Plain and Regex Replace",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			PlainReplace: []filter.PlainReplace{
				{Key: "key1", Value: "leValue1"},
			},
			RegexReplace: []filter.RegexReplace{
				{Key: "key1", Expression: regexp.MustCompile("Value1"), Value: "changedValue1"},
			},
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1", "key2": "value2"}
	expected := map[string]interface{}{"key1": "lechangedValue1"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["key1"] != expected["key1"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithMultipleFilters tests the Apply function with multiple filters
func TestApplyWithDropAndFilters(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "Just a Filter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			PlainReplace: []filter.PlainReplace{
				{Key: "key1", Value: "newValue1"},
			},
			Enable: true,
		},
		{
			FilterName: "Just a Drop filter",
			Target: []filter.Target{
				{Key: "key2", Value: regexp.MustCompile("value2")},
			},
			Drop:   true,
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1", "key2": "value2"}
	expected := map[string]interface{}{}

	result, dropped := filter.Apply(data, filters)
	if !dropped {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithDisabledFilter tests the Apply function with a disabled filter
func TestApplyWithDisabledFilter(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "DisabledFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			PlainReplace: []filter.PlainReplace{
				{Key: "key1", Value: "newValue1"},
			},
			Enable: false,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	expected := map[string]interface{}{"key1": "value1"}
	filters, _, _, _ = filter.ValidateFilters(filters)
	result, dropped := filter.Apply(data, filters)
	if dropped || result["key1"] != expected["key1"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithCategory tests the Apply function with a category filter
func TestApplyWithCategory(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "CategoryFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			Category: "Email",
			Enable:   true,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	expected := map[string]interface{}{"key1": "value1", "category": "Email"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["category"] != expected["category"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}

// TestApplyWithDefaultCategory tests the Apply function with the default category when no category is set
func TestApplyWithDefaultCategory(t *testing.T) {
	filters := []filter.Filter{
		{
			FilterName: "NoCategoryFilter",
			Target: []filter.Target{
				{Key: "key1", Value: regexp.MustCompile("value1")},
			},
			Enable: true,
		},
	}

	data := map[string]interface{}{"key1": "value1"}
	expected := map[string]interface{}{"key1": "value1", "category": "Other"}

	result, dropped := filter.Apply(data, filters)
	if dropped || result["category"] != expected["category"] {
		t.Errorf("expected %v, got %v, dropped: %v", expected, result, dropped)
	}
}
