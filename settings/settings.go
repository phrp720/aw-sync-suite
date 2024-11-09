package settings

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// SettingsKey is a custom type for settings keys
type SettingsKey string

// Define constants for each setting name
const (
	AWUrl               SettingsKey = "awUrl"
	PrometheusUrl       SettingsKey = "prometheusUrl"
	ExcludedWatchers    SettingsKey = "excludedWatchers"
	UserID              SettingsKey = "userId"
	Cron                SettingsKey = "cron"
	PrometheusSecretKey SettingsKey = "prometheus-secret-key"
	AsService           SettingsKey = "service"
	Immediate           SettingsKey = "immediate"
)
const configFile = "config.yaml"

// Settings struct
type Settings struct {
	AWUrl               string   `yaml:"aw-url"`
	PrometheusUrl       string   `yaml:"prometheus-url"`
	PrometheusSecretKey string   `yaml:"prometheus-secret-key"`
	ExcludedWatchers    []string `yaml:"excluded-watchers"`
	UserID              string   `yaml:"userId"`
	Cron                string   `yaml:"cron"`
	AsService           bool     `yaml:"-"`
	Immediate           bool     `yaml:"-"`
}

// InitSettings initializes the settings
func InitSettings() *Settings {
	settings := loadYAMLConfig(configFile)
	loadEnvVariables(&settings)
	loadFlags(&settings)
	validateSettings(&settings)
	printSettings(&settings)
	return &settings
}

// Load the YAML config file
func loadYAMLConfig(filename string) Settings {
	file, err := os.Open(filename)
	var settings Settings

	if err != nil {
		log.Print("No config.yaml file found. Proceeding with environment variables and flags.")
	} else {
		log.Print("Loading settings from config.yaml file.")
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&settings); err != nil {
			log.Fatalf("Failed to decode settings file: %v", err)
		}
		// Remove loading of SERVICE and STANDALONE from YAML config
		settings.AsService = false
	}

	return settings
}

// Load the environment variables
func loadEnvVariables(settings *Settings) {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found. Loading environment variables from the system.")
	} else {
		log.Print("Loading environment variables.")
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
	if value, exists := os.LookupEnv("PROMETHEUS_SECRET_KEY"); exists {
		settings.PrometheusSecretKey = value
	}

}

// Load the flags
func loadFlags(settings *Settings) {
	flag.StringVar(&settings.AWUrl, string(AWUrl), settings.AWUrl, "Activity Watch URL")
	flag.StringVar(&settings.PrometheusUrl, string(PrometheusUrl), settings.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&settings.UserID, string(UserID), settings.UserID, "User")
	flag.StringVar(&settings.Cron, string(Cron), settings.Cron, "Cron expression")
	flag.StringVar(&settings.PrometheusSecretKey, string(PrometheusSecretKey), settings.PrometheusSecretKey, "Prometheus Secret Key")
	flag.BoolVar(&settings.AsService, string(AsService), settings.AsService, "Run as service")
	flag.BoolVar(&settings.Immediate, string(Immediate), settings.Immediate, "Run a sync immediately")

	flag.Parse()
	log.Print("Loading settings from flags.")
}

// Validate the settings
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

// Pretty Print of the settings
func printSettings(settings *Settings) {
	log.Print("Current Settings:")

	// Create a map of settings for easier iteration
	settingsMap := map[SettingsKey]string{
		AWUrl:               settings.AWUrl,
		PrometheusUrl:       settings.PrometheusUrl,
		PrometheusSecretKey: settings.PrometheusSecretKey,
		ExcludedWatchers:    strings.Join(settings.ExcludedWatchers, ", "),
		UserID:              settings.UserID,
		Cron:                settings.Cron,
		AsService:           fmt.Sprintf("%t", settings.AsService),
	}

	// Define the order of the settings
	order := []SettingsKey{
		AWUrl,
		PrometheusUrl,
		PrometheusSecretKey,
		ExcludedWatchers,
		UserID,
		Cron,
		AsService,
	}

	maxKeyLength := 0
	maxValueLength := 0
	for _, key := range order {
		value := settingsMap[key]
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
	for _, key := range order {
		value := settingsMap[key]
		fmt.Printf("| %-*s | %-*s |\n", maxKeyLength, key, maxValueLength, value)
	}
	fmt.Println(border)
}

// CreateConfigFile creates a config file to a given path based on the settings
func CreateConfigFile(settings Settings, path string) error {

	content, err := yaml.Marshal(&settings)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)

}
