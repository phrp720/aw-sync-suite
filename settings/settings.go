package settings

import (
	"aw-sync-agent/errors"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
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
		AsService:        InitServiceFlag("service"),
	}
	flag.Parse()
	validateSettings(Settings)
	for key, value := range Settings {
		CheckSettingValue(*value, key == AWUrl || key == PrometheusUrl)
	}
	PrintSettings(Settings)
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

func InitServiceFlag(flagName string) *string {
	// Define the command-line flag and sets the default value to the environment variable value
	flagValue := flag.Bool(flagName, false, "Run as a service")
	flagValueStr := strconv.FormatBool(*flagValue)

	return &flagValueStr
}

func CheckSettingValue(value string, isMandatory bool) {
	if value == "" && isMandatory {
		log.Fatalf("The %s is mandatory", value)
	}

}
func validateSettings(settings map[SettingsKey]*string) map[SettingsKey]*string {
	for key, value := range settings {
		if key == Cron && *value == "" {
			log.Print("Cron expression is empty, setting it to default value: */5 * * * * (every 5 minutes)")

			*value = "@every 5m"
		}
		if key == AsService && *value == "" {
			*value = "false"
		}
		// Check if the -service flag was set
		if isService() {
			srv := "true"
			settings[AsService] = &srv
		}
		//
	}
	return settings
}

// PrintSettings prints the settings in a symmetric box format
func PrintSettings(settings map[SettingsKey]*string) {
	log.Print("Current Settings:")
	maxKeyLength := 0
	maxValueLength := 0
	for key, value := range settings {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
		if len(*value) > maxValueLength {
			maxValueLength = len(*value)
		}
	}
	borderLength := maxKeyLength + maxValueLength + 7
	border := strings.Repeat("-", borderLength)
	fmt.Println(border)
	for key, value := range settings {
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLength, key, maxValueLength, *value)
	}
	fmt.Println(border)
}

func isService() bool {
	// Check if the -service flag was set
	IsService := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "service" {
			IsService = true
		}
	})
	return IsService
}
