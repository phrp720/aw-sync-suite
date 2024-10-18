package settings

import (
	"aw-sync-agent/errors"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// SettingsKey is a custom type for settings keys
type SettingsKey string

// Define constants for each setting name
const (
	AWUrl            SettingsKey = "awUrl"
	PrometheusUrl    SettingsKey = "prometheusUrl"
	ExcludedWatchers SettingsKey = "excludedWatchers"
	UserID           SettingsKey = "userID"
	Cron             SettingsKey = "cron"
)

func GetEnvVar(variable string, mandatory bool) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	if os.Getenv(variable) == "" && mandatory {
		return "", &errors.EnvVarError{VarName: variable}
	} else if os.Getenv(variable) == "" && !mandatory {
		return "", nil
	}
	return os.Getenv(variable), nil
}

func GetFlagOrEnvVar(envVarName string, flagName string, isMandatory bool) string {
	// Define the command-line flag
	envValue, _ := GetEnvVar(envVarName, isMandatory)
	// Define the command-line flag and sets the default value to the environment variable value
	flagValue := flag.String(flagName, envValue, envVarName+" (command-line flag)")

	// Parse command-line flags
	flag.Parse()

	// Get the value from the command-line flag
	value := *flagValue
	if value == "" && isMandatory {
		log.Fatalf("The %s is mandatory", flagName)
	}
	return value
}
