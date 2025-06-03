package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"laravel-setup/pkg/utils"
)

// InitConfig initializes the configuration with user input
func InitConfig() (*Config, error) {
	// Start with default configuration
	config := NewConfig()

	// Get script directory
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	config.ScriptDir = filepath.Dir(ex)

	// Get domain from user input
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the Domain for your Laravel project: ")
	domain, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	config.Domain = strings.TrimSpace(domain)

	// Get repository URL from user input
	fmt.Print("Enter the Git repository URL for your Laravel project: ")
	repoURL, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	config.RepoURL = strings.TrimSpace(repoURL)

	// Generate random passwords for database
	config.DBPassword = utils.GenerateRandomPassword()
	config.DBRootPassword = utils.GenerateRandomPassword()

	// Set web root based on domain
	config.WebRoot = "/var/www/" + config.Domain

	return config, nil
}