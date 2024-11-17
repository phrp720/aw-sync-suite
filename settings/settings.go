package settings

import (
	"aw-sync-agent/filter"
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
const configFile = "aw-sync-agent.yaml"

type sett struct {
	AWUrl               string   `yaml:"aw-url"`
	PrometheusUrl       string   `yaml:"prometheus-url"`
	PrometheusSecretKey string   `yaml:"prometheus-secret-key"`
	ExcludedWatchers    []string `yaml:"excluded-watchers"`
	UserID              string   `yaml:"userId"`
	Cron                string   `yaml:"cron"`
	AsService           bool     `yaml:"-"`
	Immediate           bool     `yaml:"-"`
}

// Configuration struct
type Configuration struct {
	Settings sett            `yaml:"Settings"`
	Filters  []filter.Filter `yaml:"Filters"`
}

// InitConfigurations initializes the settings
func InitConfigurations() *Configuration {
	settings := loadYAMLConfig(configFile)
	loadEnvVariables(&settings)
	loadFlags(&settings)
	validateSettings(&settings)
	printSettings(&settings)
	return &settings
}

// Load the YAML config file
func loadYAMLConfig(filename string) Configuration {
	file, err := os.Open(filename)
	var config Configuration

	if err != nil {
		log.Print("No aw-sync-agent.yaml file found. Proceeding with environment variables and flags.")
	} else {
		log.Print("Loading settings from aw-sync-agent.yaml file.")
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Fatalf("Failed to decode settings file: %v", err)
		}
		// Remove loading of SERVICE and STANDALONE from YAML config
		config.Settings.AsService = false
	}

	return config
}

// Load the environment variables
func loadEnvVariables(config *Configuration) {
	if err := godotenv.Load(".env"); err != nil {
		log.Print("No .env file found. Loading environment variables from the system.")
	} else {
		log.Print("Loading environment variables.")
	}

	if value, exists := os.LookupEnv("ACTIVITY_WATCH_URL"); exists {
		config.Settings.AWUrl = value
	}
	if value, exists := os.LookupEnv("PROMETHEUS_URL"); exists {
		config.Settings.PrometheusUrl = value
	}
	if value, exists := os.LookupEnv("EXCLUDED_WATCHERS"); exists {
		config.Settings.ExcludedWatchers = strings.Split(value, ",")
	}
	if value, exists := os.LookupEnv("USER_ID"); exists {
		config.Settings.UserID = value
	}
	if value, exists := os.LookupEnv("CRON"); exists {
		config.Settings.Cron = value
	}
	if value, exists := os.LookupEnv("PROMETHEUS_SECRET_KEY"); exists {
		config.Settings.PrometheusSecretKey = value
	}

}

// Load the flags
func loadFlags(config *Configuration) {
	flag.StringVar(&config.Settings.AWUrl, string(AWUrl), config.Settings.AWUrl, "Activity Watch URL")
	flag.StringVar(&config.Settings.PrometheusUrl, string(PrometheusUrl), config.Settings.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&config.Settings.UserID, string(UserID), config.Settings.UserID, "User")
	flag.StringVar(&config.Settings.Cron, string(Cron), config.Settings.Cron, "Cron expression")
	flag.StringVar(&config.Settings.PrometheusSecretKey, string(PrometheusSecretKey), config.Settings.PrometheusSecretKey, "Prometheus Secret Key")
	flag.BoolVar(&config.Settings.AsService, string(AsService), config.Settings.AsService, "Run as service")
	flag.BoolVar(&config.Settings.Immediate, string(Immediate), config.Settings.Immediate, "Run a sync immediately")

	flag.Parse()
	log.Print("Loading settings from flags.")
}

// Validate the settings
func validateSettings(config *Configuration) {
	if config.Settings.AWUrl == "" {
		log.Print("Activity Watch URL is mandatory but it isn't defined! Setting it to default value: http://localhost:5600")
		config.Settings.AWUrl = "http://localhost:5600"
	}
	if config.Settings.PrometheusUrl == "" {
		log.Fatal("Prometheus URL is mandatory")
	}
	if config.Settings.Cron == "" {
		log.Print("Cron expression is empty, setting it to default value: */5 * * * * (every 5 minutes)")
		config.Settings.Cron = "*/5 * * * *"
	}
}

// Pretty Print of the settings
func printSettings(config *Configuration) {
	log.Print("Current Settings:")

	// Create a map of settings for easier iteration
	settingsMap := map[SettingsKey]string{
		AWUrl:               config.Settings.AWUrl,
		PrometheusUrl:       config.Settings.PrometheusUrl,
		PrometheusSecretKey: config.Settings.PrometheusSecretKey,
		ExcludedWatchers:    strings.Join(config.Settings.ExcludedWatchers, ", "),
		UserID:              config.Settings.UserID,
		Cron:                config.Settings.Cron,
	}

	// Define the order of the settings
	order := []SettingsKey{
		AWUrl,
		PrometheusUrl,
		PrometheusSecretKey,
		ExcludedWatchers,
		UserID,
		Cron,
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
func CreateConfigFile(config Configuration, path string) error {

	content, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}
