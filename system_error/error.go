package system_error

import (
	"fmt"
	"log"
)

// EnvVarError Defines custom errors types
type EnvVarError struct {
	VarName string
}

func (e *EnvVarError) Error() string {
	return fmt.Sprintf("Environment variable %s is not set or is empty", e.VarName)
}

func HandleNormal(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}
func HandleFatal(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
