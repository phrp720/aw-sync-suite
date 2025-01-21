package settings

import "strings"

// StringSliceFlag is a custom flag type that holds a slice of strings
type StringSliceFlag []string

// String returns the string representation of the StringSlice
func (s *StringSliceFlag) String() string {
	return strings.Join(*s, "|")
}

// Set appends a new value to the StringSlice
func (s *StringSliceFlag) Set(value string) error {
	*s = strings.Split(value, "|")
	return nil
}
