package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds all the settings for the Laravel setup
type Config struct {
	Domain         string
	RepoURL        string
	DBName         string
	DBUser         string
	DBPassword     string
	DBRootPassword string
	WebUser        string
	SSHPort        string
	WebRoot        string
	ScriptDir      string
	// Skip flags
	SkipSystemUpdate bool
	SkipEssentials   bool
	SkipPHP          bool
	SkipMySQL        bool
	SkipNginx        bool
	SkipSecurity     bool
	SkipLaravel      bool
	SkipServices     bool
}

// NewConfig initializes a new configuration with default values
func NewConfig() *Config {
	return &Config{
		DBName:  "production_db",
		DBUser:  "db_user",
		SSHPort: "2222",
		WebUser: "www-data",
	}
}

// LoadConfigFromFile loads configuration from a TOML file
func LoadConfigFromFile(configPath string) (*Config, error) {
	config := NewConfig()

	// Check if the file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Return default config if file doesn't exist
	}

	// Read and parse the TOML file
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetDefaultConfigPath returns the default path for the config file in the user's home directory
func GetDefaultConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "config.toml"), nil
}
