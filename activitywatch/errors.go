package activitywatch

import "fmt"

// Define custom error types
type EnvVarError struct {
	VarName string
}

func (e *EnvVarError) Error() string {
	return fmt.Sprintf("Environment variable %s is not set or is empty", e.VarName)
}

type HTTPError struct {
	URL string
	Err error
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("failed to get %s: %v", e.URL, e.Err)
}

type DecodeError struct {
	Err error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("failed to decode response: %v", e.Err)
}
