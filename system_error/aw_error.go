package system_error

import "fmt"

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
