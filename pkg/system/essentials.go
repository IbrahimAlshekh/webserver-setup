package system

import (
	"laravel-setup/pkg/config"
	"strings"

	"laravel-setup/pkg/utils"
)

// InstallEssentials installs essential system packages
// These packages are required for the Laravel server setup
func InstallEssentials(_ *config.Config) error {
	utils.PrintStatus("Installing essential system packages...")

	// Install essential packages
	// These packages provide core functionality for the server
	err := utils.RunCommand("sudo", "apt", "install", "-y",
		"curl", "wget", "git", "unzip", "software-properties-common",
		"apt-transport-https", "ca-certificates", "gnupg", "lsb-release",
		"ufw", "fail2ban", "htop", "tree", "vim", "supervisor",
		"redis-server", "certbot", "python3-certbot-nginx")
	if err != nil {
		return err
	}

	utils.PrintStatus("Essential packages installed successfully")

	// Install Composer (PHP dependency manager)
	if err := installComposer(); err != nil {
		return err
	}

	// Install Node.js and npm (JavaScript runtime and package manager)
	if err := installNodeJS(); err != nil {
		return err
	}

	return nil
}

// installComposer installs the Composer PHP dependency manager
func installComposer() error {
	utils.PrintHeader("Installing Composer")
	utils.PrintStatus("Downloading and installing Composer...")

	// Download Composer installer
	err := utils.RunCommand("curl", "-sS", "https://getcomposer.org/installer", "-o", "composer-setup.php")
	if err != nil {
		return err
	}

	// Run the installer
	err = utils.RunCommand("php", "composer-setup.php")
	if err != nil {
		return err
	}

	// Move composer.phar to a directory in the PATH
	err = utils.RunCommand("sudo", "mv", "composer.phar", "/usr/local/bin/composer")
	if err != nil {
		return err
	}

	// Make it executable
	err = utils.RunCommand("sudo", "chmod", "+x", "/usr/local/bin/composer")
	if err != nil {
		return err
	}

	utils.PrintStatus("Composer installed successfully")
	return nil
}

// installNodeJS installs Node.js and npm
func installNodeJS() error {
	utils.PrintHeader("Installing Node.js and npm")
	utils.PrintStatus("Adding Node.js repository and installing Node.js...")

	// Download Node.js setup script
	err := utils.RunCommand("curl", "-fsSL", "https://deb.nodesource.com/setup_20.x", "-o", "nodejs-setup.sh")
	if err != nil {
		return err
	}

	// Run the setup script
	err = utils.RunCommand("sudo", "bash", "nodejs-setup.sh")
	if err != nil {
		return err
	}

	// Install Node.js
	err = utils.RunCommand("sudo", "apt", "install", "-y", "nodejs")
	if err != nil {
		return err
	}

	// Check Node.js and npm versions
	nodeVersion, err := utils.RunCommandWithOutput("node", "-v")
	if err != nil {
		return err
	}

	npmVersion, err := utils.RunCommandWithOutput("npm", "-v")
	if err != nil {
		return err
	}

	utils.PrintStatus("Node.js and npm installed successfully")
	utils.PrintStatus("Node.js version: " + strings.TrimSpace(nodeVersion))
	utils.PrintStatus("npm version: " + strings.TrimSpace(npmVersion))

	return nil
}
