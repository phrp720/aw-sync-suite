package errors

// WriteError returned if HTTP call is finished with response status code, but it was not successful.
type WriteError struct {
	Err  error
	Code int
}

func (e *WriteError) Error() string {
	return e.Err.Error()
}

func (e *WriteError) StatusCode() int {
	return e.Code
}
