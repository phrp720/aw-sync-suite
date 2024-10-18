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
	MinData          SettingsKey = "minData"
	AsService        SettingsKey = "asService"
)

func InitSettings() map[SettingsKey]*string {
	// These are the settings that contains the Env Variables/Flags
	Settings := map[SettingsKey]*string{
		AWUrl:            InitFlag("ACTIVITY_WATCH_URL", "awUrl", true),
		PrometheusUrl:    InitFlag("PROMETHEUS_URL", "prometheusUrl", true),
		ExcludedWatchers: InitFlag("EXCLUDED_WATCHERS", "excludedWatchers", false),
		UserID:           InitFlag("USER_ID", "userID", false),
		Cron:             InitFlag("CRON", "cron", false),
		MinData:          InitFlag("MIN_DATA", "minData", false),
		AsService:        InitFlag("AS_SERVICE", "asService", false),
	}
	flag.Parse()
	for key, value := range Settings {
		CheckSettingValue(*value, key == AWUrl || key == PrometheusUrl)
	}
	return Settings
}
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

func InitFlag(envVarName string, flagName string, isMandatory bool) *string {
	// Define the command-line flag
	envValue, _ := GetEnvVar(envVarName, isMandatory)
	// Define the command-line flag and sets the default value to the environment variable value
	flagValue := flag.String(flagName, envValue, envVarName+" (env Variable)")

	return flagValue
}

func CheckSettingValue(value string, isMandatory bool) {
	if value == "" && isMandatory {
		log.Fatalf("The %s is mandatory", value)
	}
}
