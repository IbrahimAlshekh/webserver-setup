package nginx

import (
	"os"
	"strings"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Install installs and configures Nginx for Laravel
func Install(config *config.Config) error {
	utils.PrintHeader("Installing Nginx")
	utils.PrintStatus("Installing Nginx web server...")

	// Install Nginx
	err := utils.RunCommand("sudo", "apt", "install", "-y", "nginx")
	if err != nil {
		return err
	}

	utils.PrintStatus("Nginx installed successfully")

	// Configure Nginx for Laravel
	utils.PrintHeader("Configuring Nginx for Laravel")
	utils.PrintStatus("Setting up Nginx configuration for domain: " + config.Domain)

	// Add rate-limiting zones to nginx.conf for security
	// This helps prevent brute force and DoS attacks

	// Check if rate-limiting zones already exist to prevent duplication
	_, err = utils.RunCommandWithOutput("sudo", "grep", "limit_req_zone.*zone=login", "/etc/nginx/nginx.conf")
	if err != nil && !strings.Contains(err.Error(), "exit status 1") {
		// Real error, not just "not found"
		return err
	}

	// Only add rate-limiting zones if they don't already exist
	if err != nil {
		// Rate-limiting zones don't exist, add them
		err = utils.RunCommand("sudo", "sed", "-i", `/http {/a\\n    # Rate limiting zones\n    limit_req_zone $binary_remote_addr zone=login:10m rate=10r/m;\n    limit_req_zone $binary_remote_addr zone=api:10m rate=100r/m;`, "/etc/nginx/nginx.conf")
		if err != nil {
			return err
		}
		utils.PrintStatus("Added rate limiting zones to nginx.conf")
	} else {
		utils.PrintStatus("Rate limiting zones already exist in nginx.conf")
	}

	// Create the site configuration using the template
	nginxConfig := templates.GetNginxConfig(config.Domain, config.WebRoot)

	// Write Nginx configuration to file
	err = os.WriteFile("nginx_site.conf", []byte(nginxConfig), 0644)
	if err != nil {
		return err
	}

	// Move Nginx configuration to sites-available directory
	err = utils.RunCommand("sudo", "mv", "nginx_site.conf", "/etc/nginx/sites-available/"+config.Domain)
	if err != nil {
		return err
	}

	// Enable the site by creating a symbolic link in sites-enabled
	err = utils.RunCommand("sudo", "ln", "-sf", "/etc/nginx/sites-available/"+config.Domain, "/etc/nginx/sites-enabled/")
	if err != nil {
		return err
	}

	// Remove default site to prevent conflicts
	err = utils.RunCommand("sudo", "rm", "-f", "/etc/nginx/sites-enabled/default")
	if err != nil {
		return err
	}

	// Test Nginx configuration
	utils.PrintStatus("Testing Nginx configuration...")
	err = utils.RunCommand("sudo", "nginx", "-t")
	if err != nil {
		return err
	}

	utils.PrintStatus("Nginx configuration is valid")

	// Restart Nginx to apply changes
	err = utils.RunCommand("sudo", "systemctl", "restart", "nginx")
	if err != nil {
		return err
	}

	// Create web directory if it doesn't exist
	utils.PrintStatus("Setting up web directory...")
	err = utils.RunCommand("sudo", "mkdir", "-p", config.WebRoot)
	if err != nil {
		return err
	}

	// Set proper ownership and permissions
	// This allows the web server to access the files while maintaining security
	err = utils.RunCommand("sudo", "chown", "-R", os.Getenv("USER")+":"+config.WebUser, config.WebRoot)
	if err != nil {
		return err
	}

	err = utils.RunCommand("sudo", "chmod", "-R", "755", config.WebRoot)
	if err != nil {
		return err
	}

	utils.PrintStatus("Nginx configured successfully for " + config.Domain)

	return nil
}
