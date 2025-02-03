package settings

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"io"
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
	PrometheusSecretKey SettingsKey = "prometheusAuth"
	AsService           SettingsKey = "service"
	Immediate           SettingsKey = "immediate"
	IncludeHostname     SettingsKey = "includeHostname"
	TestConfigs         SettingsKey = "testConfig"
	Plugins             SettingsKey = "plugins"
	PluginsStrictOrder  SettingsKey = "pluginsStrictOrder"
)
const configFile = "aw-sync-settings.yaml"
const configDir = "./config"
const configAbs = configDir + "/" + configFile

type Setts struct {
	AWUrl               string   `yaml:"awUrl"`
	PrometheusUrl       string   `yaml:"prometheusUrl"`
	PrometheusSecretKey string   `yaml:"prometheusAuth"`
	ExcludedWatchers    []string `yaml:"excludedWatchers"`
	UserID              string   `yaml:"userId"`
	IncludeHostname     bool     `yaml:"includeHostname"`
	Cron                string   `yaml:"cron"`
	Plugins             []string `yaml:"plugins"`
	PluginsStrictOrder  bool     `yaml:"pluginsStrictOrder"`
	AsService           bool     `yaml:"-"`
	Immediate           bool     `yaml:"-"`
	TestConfigs         bool     `yaml:"-"` // TestConfigs is a flag to test the configurations/filters
}

// Configuration struct
type Configuration struct {
	Settings Setts `yaml:"Settings"`
}

// InitConfigurations initializes the settings
func InitConfigurations() *Configuration {
	settings := LoadYAMLConfig(configAbs)
	loadEnvVariables(&settings)
	loadFlags(&settings)
	ValidateSettings(&settings)
	printSettings(&settings)
	return &settings
}

// LoadYAMLConfig  Load the YAML config file
func LoadYAMLConfig(filename string) Configuration {
	file, err := os.Open(filename)
	var config Configuration

	if err != nil {
		log.Printf("No %s file found. Proceeding with environment variables and flags.", configFile)
	} else {
		log.Printf("Loading settings from %s file.", configFile)
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err = decoder.Decode(&config); err != nil && err != io.EOF {
			log.Fatalf("Error: Failed to decode settings file: %v", err)
		}

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
		config.Settings.ExcludedWatchers = strings.Split(value, "|")
	}
	if value, exists := os.LookupEnv("PLUGINS"); exists {
		config.Settings.Plugins = strings.Split(value, "|")
	}
	if value, exists := os.LookupEnv("PLUGINS_STRICT_ORDER"); exists {
		config.Settings.PluginsStrictOrder = value == "true"
	}
	if value, exists := os.LookupEnv("USER_ID"); exists {
		config.Settings.UserID = value
	}
	if value, exists := os.LookupEnv("CRON"); exists {
		config.Settings.Cron = value
	}
	if value, exists := os.LookupEnv("PROMETHEUS_AUTH"); exists {
		config.Settings.PrometheusSecretKey = value
	}
	if value, exists := os.LookupEnv("INCLUDE_HOSTNAME"); exists {
		config.Settings.IncludeHostname = value == "true"
	}

}

// Load the flags
func loadFlags(config *Configuration) {
	flag.StringVar(&config.Settings.AWUrl, string(AWUrl), config.Settings.AWUrl, "Activity Watch URL")
	flag.StringVar(&config.Settings.PrometheusUrl, string(PrometheusUrl), config.Settings.PrometheusUrl, "Prometheus URL")
	flag.StringVar(&config.Settings.UserID, string(UserID), config.Settings.UserID, "User Identification")
	flag.BoolVar(&config.Settings.IncludeHostname, string(IncludeHostname), config.Settings.IncludeHostname, "Include hostname in the metrics")

	var excludedWatchers StringSliceFlag
	flag.Var(&excludedWatchers, string(ExcludedWatchers), "Excluded watchers")

	var plugins StringSliceFlag
	flag.Var(&plugins, string(Plugins), "Plugins to load")

	flag.BoolVar(&config.Settings.PluginsStrictOrder, string(PluginsStrictOrder), config.Settings.PluginsStrictOrder, "Plugins strict order")

	flag.StringVar(&config.Settings.Cron, string(Cron), config.Settings.Cron, "Cron expression")
	flag.StringVar(&config.Settings.PrometheusSecretKey, string(PrometheusSecretKey), config.Settings.PrometheusSecretKey, "Prometheus Secret Key")
	flag.BoolVar(&config.Settings.AsService, string(AsService), config.Settings.AsService, "Run as service")
	flag.BoolVar(&config.Settings.Immediate, string(Immediate), config.Settings.Immediate, "Run a sync immediately")
	flag.BoolVar(&config.Settings.TestConfigs, string(TestConfigs), config.Settings.TestConfigs, "Test the configurations/filters")

	flag.Parse()
	if len(excludedWatchers) > 0 {
		config.Settings.ExcludedWatchers = excludedWatchers
	}
	if len(plugins) > 0 {
		config.Settings.Plugins = plugins

	}

	log.Print("Loading settings from flags.")
}

// ValidateSettings Validate the settings
func ValidateSettings(config *Configuration) {
	if config.Settings.AWUrl == "" {
		log.Print("Warning: Activity Watch URL is not defined. Defaulting to http://localhost:5600")
		config.Settings.AWUrl = "http://localhost:5600"
	}
	if config.Settings.PrometheusUrl == "" {
		log.Fatal("Error: Prometheus URL is mandatory")
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
		ExcludedWatchers:    strings.Join(config.Settings.ExcludedWatchers, ","),
		UserID:              config.Settings.UserID,
		IncludeHostname:     fmt.Sprintf("%t", config.Settings.IncludeHostname),
		Plugins:             strings.Join(config.Settings.Plugins, ","),
		PluginsStrictOrder:  fmt.Sprintf("%t", config.Settings.PluginsStrictOrder),
		Cron:                config.Settings.Cron,
	}

	// Define the order of the settings
	order := []SettingsKey{
		AWUrl,
		PrometheusUrl,
		PrometheusSecretKey,
		ExcludedWatchers,
		UserID,
		IncludeHostname,
		Plugins,
		PluginsStrictOrder,
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
