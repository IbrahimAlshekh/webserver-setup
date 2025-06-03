package mysql

import (
	"os"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Install installs and configures MySQL for Laravel
func Install(config *config.Config) error {
	utils.PrintHeader("Installing MySQL 8.0")
	utils.PrintStatus("Installing MySQL server and client...")

	// Install MySQL server and client
	err := utils.RunCommand("sudo", "apt", "install", "-y", "mysql-server", "mysql-client")
	if err != nil {
		return err
	}

	utils.PrintStatus("MySQL installed successfully")

	// Secure MySQL installation
	utils.PrintHeader("Securing MySQL Installation")
	utils.PrintWarning("Please set a strong root password when prompted")
	err = utils.RunInteractiveCommand("sudo", "mysql_secure_installation")
	if err != nil {
		return err
	}

	// Configure MySQL for Laravel
	utils.PrintHeader("Configuring MySQL for Laravel")
	utils.PrintStatus("Configuring MySQL database and user...")
	utils.PrintStatus("Creating database: " + config.DBName)
	utils.PrintStatus("Creating user: " + config.DBUser)

	// Create MySQL configuration script
	mysqlConfig := templates.GetMySQLConfig(
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBRootPassword,
	)

	// Write MySQL configuration to file
	err = os.WriteFile("mysql_config.sql", []byte(mysqlConfig), 0644)
	if err != nil {
		return err
	}

	// Apply MySQL configuration
	err = utils.RunCommandWithFileInput("mysql_config.sql", "sudo", "mysql")
	if err != nil {
		return err
	}

	utils.PrintStatus("MySQL configured successfully")

	// Save credentials securely
	credentialsContent := templates.GetMySQLCredentialsContent(
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBRootPassword,
	)

	// Write credentials to file with restricted permissions
	err = os.WriteFile("/home/"+os.Getenv("USER")+"/mysql_credentials.txt", []byte(credentialsContent), 0600)
	if err != nil {
		return err
	}

	utils.PrintStatus("MySQL credentials saved to ~/mysql_credentials.txt")

	// Configure Redis for caching
	if err := configureRedis(); err != nil {
		return err
	}

	return nil
}

// configureRedis configures Redis for caching
// Redis is commonly used with Laravel for caching, sessions, and queue
func configureRedis() error {
	utils.PrintHeader("Configuring Redis")
	utils.PrintStatus("Optimizing Redis configuration...")

	// Set maximum memory to prevent Redis from using all available memory
	err := utils.RunCommand("sudo", "sed", "-i", "s/# maxmemory <bytes>/maxmemory 256mb/", "/etc/redis/redis.conf")
	if err != nil {
		return err
	}

	// Set eviction policy to remove least recently used keys when memory is full
	err = utils.RunCommand("sudo", "sed", "-i", "s/# maxmemory-policy noeviction/maxmemory-policy allkeys-lru/", "/etc/redis/redis.conf")
	if err != nil {
		return err
	}

	// Enable Redis to start on boot
	err = utils.RunCommand("sudo", "systemctl", "enable", "redis-server")
	if err != nil {
		return err
	}

	// Restart Redis to apply changes
	err = utils.RunCommand("sudo", "systemctl", "restart", "redis-server")
	if err != nil {
		return err
	}

	utils.PrintStatus("Redis configured successfully")
	return nil
}
