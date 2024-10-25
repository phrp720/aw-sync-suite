package settings

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
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
	UserID           SettingsKey = "userId"
	Cron             SettingsKey = "cron"
	MinData          SettingsKey = "min-data"
	AsService        SettingsKey = "service"
	Standalone       SettingsKey = "standalone"
)
const configFile = "config.yaml"

type Settings struct {
	AWUrl            string   `yaml:"aw-url"`
	PrometheusUrl    string   `yaml:"prometheus-url"`
	ExcludedWatchers []string `yaml:"excluded-watchers"`
	UserID           string   `yaml:"userId"`
	Cron             string   `yaml:"cron"`
	MinData          string   `yaml:"min-data"`
	AsService        bool     `yaml:"service"`
	Standalone       bool     `yaml:"standalone"`
}

func InitSettings() *Settings {
	settings := loadYAMLConfig(configFile)
	loadEnvVariables(&settings)
	loadFlags(&settings)
	validateSettings(&settings)
	printSettings(&settings)
	return &settings
}

func loadYAMLConfig(filename string) Settings {
	file, err := os.Open(filename)
	var settings Settings

	if err != nil {
		log.Print("No config.yaml file found. Proceeding with environment variables and flags.")
	} else {

		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&settings); err != nil {
			log.Fatalf("Failed to decode settings file: %v", err)
		}
	}

	return settings
}

func loadEnvVariables(settings *Settings) {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found. Loading environment variables from the system.")
	}

	if value, exists := os.LookupEnv("ACTIVITY_WATCH_URL"); exists {
		settings.AWUrl = value
	}
	if value, exists := os.LookupEnv("PROMETHEUS_URL"); exists {
		settings.PrometheusUrl = value
	}
	if value, exists := os.LookupEnv("EXCLUDED_WATCHERS"); exists {
		settings.ExcludedWatchers = strings.Split(value, ",")
	}
	if value, exists := os.LookupEnv("USER_ID"); exists {
		settings.UserID = value
	}
	if value, exists := os.LookupEnv("CRON"); exists {
		settings.Cron = value
	}
	if value, exists := os.LookupEnv("MIN_DATA"); exists {
		settings.MinData = value
	}
	if value, exists := os.LookupEnv("SERVICE"); exists {
		settings.AsService, _ = strconv.ParseBool(value)
	}
	if value, exists := os.LookupEnv("STANDALONE"); exists {
		settings.Standalone, _ = strconv.ParseBool(value)
	}
}

func loadFlags(settings *Settings) {
	flag.StringVar(&settings.AWUrl, string(AWUrl), settings.AWUrl, "Activity Watch URL")
	flag.StringVar(&settings.PrometheusUrl, string(PrometheusUrl), settings.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&settings.UserID, string(UserID), settings.UserID, "User")
	flag.StringVar(&settings.Cron, string(Cron), settings.Cron, "Cron expression")
	flag.StringVar(&settings.MinData, string(MinData), settings.MinData, "Minimum data")
	flag.BoolVar(&settings.AsService, string(AsService), settings.AsService, "Run as service")
	flag.BoolVar(&settings.Standalone, string(Standalone), settings.Standalone, "Run in standalone mode")
	flag.Parse()
}

func validateSettings(settings *Settings) {
	if settings.AWUrl == "" {
		log.Fatal("Activity Watch URL is mandatory")
	}
	if settings.PrometheusUrl == "" {
		log.Fatal("Prometheus URL is mandatory")
	}
	if settings.Cron == "" {
		log.Print("Cron expression is empty, setting it to default value: */5 * * * * (every 5 minutes)")
		settings.Cron = "@every 5m"
	}
}

func printSettings(settings *Settings) {
	log.Print("Current Settings:")

	// Create a map of settings for easier iteration
	settingsMap := map[SettingsKey]string{
		AWUrl:            settings.AWUrl,
		PrometheusUrl:    settings.PrometheusUrl,
		ExcludedWatchers: strings.Join(settings.ExcludedWatchers, ", "),
		UserID:           settings.UserID,
		Cron:             settings.Cron,
		MinData:          settings.MinData,
		AsService:        fmt.Sprintf("%t", settings.AsService),
		Standalone:       fmt.Sprintf("%t", settings.Standalone),
	}

	maxKeyLength := 0
	maxValueLength := 0
	for key, value := range settingsMap {
		if len(key) > maxKeyLength {
			maxKeyLength = len(key)
		}
		if len(value) > maxValueLength {
			maxValueLength = len(value)
		}
	}

	borderLength := maxKeyLength + maxValueLength + 7
	border := strings.Repeat("-", borderLength)
	fmt.Println(border)
	for key, value := range settingsMap {
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLength, key, maxValueLength, value)
	}
	fmt.Println(border)
}
