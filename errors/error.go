package errors

import "fmt"

// EnvVarError Defines custom errors types
type EnvVarError struct {
	VarName string
}

func (e *EnvVarError) Error() string {
	return fmt.Sprintf("Environment variable %s is not set or is empty", e.VarName)
}
