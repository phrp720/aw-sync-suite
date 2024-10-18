package util

import (
	"aw-sync-agent/errors"
	"os"
)

func GetEnvVar(variable string, mandatory bool) (string, error) {
	if os.Getenv(variable) == "" && mandatory {
		return "", &errors.EnvVarError{VarName: variable}
	} else if os.Getenv(variable) == "" && !mandatory {
		return "", nil
	}
	return os.Getenv(variable), nil
}
