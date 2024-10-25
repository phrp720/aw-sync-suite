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

// Define constants for each setting name(These values are the flags and the environment variables)
const (
	AWUrl            SettingsKey = "awUrl"
	PrometheusUrl    SettingsKey = "prometheusUrl"
	ExcludedWatchers SettingsKey = "excludedWatchers"
	UserID           SettingsKey = "userID"
	Cron             SettingsKey = "cron"
	MinData          SettingsKey = "minData"
	AsService        SettingsKey = "service"
	Standalone       SettingsKey = "standalone"
)

func InitSettings() map[SettingsKey]*string {
	// These are the settings that contains the Env Variables/Flags
	Settings := map[SettingsKey]*string{
		AWUrl:            InitFlag("ACTIVITY_WATCH_URL", string(AWUrl), true),
		PrometheusUrl:    InitFlag("PROMETHEUS_URL", string(PrometheusUrl), true),
		ExcludedWatchers: InitFlag("EXCLUDED_WATCHERS", string(ExcludedWatchers), false),
		UserID:           InitFlag("USER_ID", string(UserID), false),
		Cron:             InitFlag("CRON", string(Cron), false),
		MinData:          InitFlag("MIN_DATA", string(MinData), false),
		AsService:        InitBooleanFlag(string(AsService)),
		Standalone:       InitBooleanFlag(string(Standalone)),
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

func InitBooleanFlag(flagName string) *string {
	// Define the command-line flag and sets the default value to the environment variable value
	flagValue := flag.Bool(flagName, false, "")
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

	}
	SetBooleanSetting(settings, AsService)
	SetBooleanSetting(settings, Standalone)

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

func FlagExists(flg string) bool {
	// Check if the -service flag was set
	Exists := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == flg {
			Exists = true
		}
	})
	return Exists
}

func SetBooleanSetting(settings map[SettingsKey]*string, key SettingsKey) {
	if FlagExists(string(key)) {
		srv := "true"
		settings[key] = &srv
	} else {
		srv := "false"
		settings[key] = &srv
	}
}

func IsStandalone(standalone string) bool {
	return standalone == "true"
}

func IsService(service string) bool {
	return service == "true"
}
