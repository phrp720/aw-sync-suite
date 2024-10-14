package remote_write

// WriteError returned if HTTP call is finished with response status code, but it was not successful.
type WriteError struct {
	err  error
	code int
}

func (e *WriteError) Error() string {
	return e.err.Error()
}

func (e *WriteError) StatusCode() int {
	return e.code
}
