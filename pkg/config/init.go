package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"laravel-setup/pkg/utils"
)

// InitConfig initializes the configuration with user input and/or config file
func InitConfig(configPath string) (*Config, error) {
	var config *Config
	var err error

	// Try to load configuration from a file
	if configPath != "" {
		// Use provided a config path
		utils.PrintStatus("Loading configuration from: " + configPath)
		config, err = LoadConfigFromFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	} else {
		// Try default config a path
		defaultPath, err := GetDefaultConfigPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default config path: %w", err)
		}

		if _, err := os.Stat(defaultPath); err == nil {
			utils.PrintStatus("Loading configuration from: " + defaultPath)
			config, err = LoadConfigFromFile(defaultPath)
			if err != nil {
				return nil, fmt.Errorf("failed to load config from %s: %w", defaultPath, err)
			}
		} else {
			// No config file found, use default config
			utils.PrintStatus("No configuration file found, using default configuration")
			config = NewConfig()
		}
	}

	// Get script directory
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	config.ScriptDir = filepath.Dir(ex)

	reader := bufio.NewReader(os.Stdin)

	// Get domain from user input if not in config
	if config.Domain == "" {
		fmt.Print("Enter the Domain for your Laravel project: ")
		domain, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		config.Domain = strings.TrimSpace(domain)
	}

	// Get repository URL from user input if not in config
	if config.RepoURL == "" {
		fmt.Print("Enter the Git repository URL for your Laravel project: ")
		repoURL, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		config.RepoURL = strings.TrimSpace(repoURL)
	}

	// Generate random passwords for database if not in config
	if config.DBPassword == "" {
		config.DBPassword = utils.GenerateRandomPassword()
	}

	if config.DBRootPassword == "" {
		config.DBRootPassword = utils.GenerateRandomPassword()
	}

	// Set web root based on domain if not in config
	if config.WebRoot == "" {
		config.WebRoot = "/var/www/" + config.Domain
	}

	return config, nil
}
