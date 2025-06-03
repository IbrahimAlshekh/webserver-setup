package php

import (
	"laravel-setup/pkg/config"
	"os"

	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Install installs PHP 8.3 and required extensions for Laravel
func Install(_ *config.Config) error {
	utils.PrintHeader("Installing PHP 8.3 and Extensions")
	utils.PrintStatus("Adding PHP repository and installing PHP 8.3 with extensions...")

	// Add a PHP repository from Ondrej (maintained PPA for latest PHP versions)
	err := utils.RunCommand("sudo", "add-apt-repository", "ppa:ondrej/php", "-y")
	if err != nil {
		return err
	}

	// Update package lists after adding the repository
	err = utils.RunCommand("sudo", "apt", "update")
	if err != nil {
		return err
	}

	// Install PHP and extensions required for Laravel
	err = utils.RunCommand("sudo", "apt", "install", "-y",
		"php8.3", "php8.3-fpm", "php8.3-mysql", "php8.3-mbstring",
		"php8.3-xml", "php8.3-bcmath", "php8.3-curl", "php8.3-gd",
		"php8.3-zip", "php8.3-intl", "php8.3-soap", "php8.3-redis",
		"php8.3-imagick", "php8.3-cli", "php8.3-common", "php8.3-opcache")
	if err != nil {
		return err
	}

	utils.PrintStatus("PHP 8.3 and extensions installed successfully")

	// Configure PHP-FPM for optimal Laravel performance
	if err := configurePHPFPM(); err != nil {
		return err
	}

	// Display a PHP version for verification
	err = utils.RunCommand("php", "-v")
	if err != nil {
		return err
	}

	return nil
}

// configurePHPFPM configures PHP-FPM for optimal Laravel performance
func configurePHPFPM() error {
	utils.PrintHeader("Configuring PHP-FPM")
	utils.PrintStatus("Optimizing PHP configuration for Laravel...")

	// Adjust PHP settings for Laravel
	// Disable path info fixing for security
	err := utils.RunCommand("sudo", "sed", "-i", "s/;cgi.fix_pathinfo=1/cgi.fix_pathinfo=0/", "/etc/php/8.3/fpm/php.ini")
	if err != nil {
		return err
	}

	// Increase upload size limit for larger file uploads
	err = utils.RunCommand("sudo", "sed", "-i", "s/upload_max_filesize = 2M/upload_max_filesize = 64M/", "/etc/php/8.3/fpm/php.ini")
	if err != nil {
		return err
	}

	// Increase post-size limit to match upload size
	err = utils.RunCommand("sudo", "sed", "-i", "s/post_max_size = 8M/post_max_size = 64M/", "/etc/php/8.3/fpm/php.ini")
	if err != nil {
		return err
	}

	// Increase execution time for longer-running scripts
	err = utils.RunCommand("sudo", "sed", "-i", "s/max_execution_time = 30/max_execution_time = 300/", "/etc/php/8.3/fpm/php.ini")
	if err != nil {
		return err
	}

	// Increase the memory limit for more complex applications
	err = utils.RunCommand("sudo", "sed", "-i", "s/memory_limit = 128M/memory_limit = 512M/", "/etc/php/8.3/fpm/php.ini")
	if err != nil {
		return err
	}

	// Configure OPcache for better performance
	utils.PrintStatus("Configuring OPcache for better performance...")

	// Write OPcache configuration to file
	err = os.WriteFile("opcache.ini", []byte(templates.OPcacheConfig), 0644)
	if err != nil {
		return err
	}

	// Move OPcache configuration to PHP configuration directory
	err = utils.RunCommand("sudo", "mv", "opcache.ini", "/etc/php/8.3/fpm/conf.d/10-opcache.ini")
	if err != nil {
		return err
	}

	// Restart PHP-FPM to apply changes
	err = utils.RunCommand("sudo", "systemctl", "restart", "php8.3-fpm")
	if err != nil {
		return err
	}

	utils.PrintStatus("PHP-FPM configured successfully")
	return nil
}
