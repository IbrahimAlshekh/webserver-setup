package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Configure configures and starts all services
func Configure(config *config.Config) error {
	utils.PrintHeader("Configuring and Starting Services")

	// Enable and start Nginx
	if err := enableService("nginx"); err != nil {
		return err
	}

	// Enable and start PHP-FPM
	if err := enableService("php8.3-fpm"); err != nil {
		return err
	}

	// Enable and start MySQL
	if err := enableService("mysql"); err != nil {
		return err
	}

	// Enable and start Redis
	if err := enableService("redis-server"); err != nil {
		return err
	}

	// Enable and start Supervisor
	if err := enableService("supervisor"); err != nil {
		return err
	}

	// Setup SSL certificate
	if err := setupSSL(config); err != nil {
		return err
	}

	// Create server information file
	if err := createServerInfo(config); err != nil {
		return err
	}

	utils.PrintHeader("Services Configuration Complete")
	utils.PrintStatus("All services have been configured and started")
	utils.PrintStatus("Server information saved to: /home/" + os.Getenv("USER") + "/server_info.txt")

	return nil
}

// enableService enables and starts a service
func enableService(service string) error {
	utils.PrintStatus("Enabling and starting " + service + "...")
	
	// Enable service to start on boot
	err := utils.RunCommand("sudo", "systemctl", "enable", service)
	if err != nil {
		return err
	}

	// Restart service to apply changes
	err = utils.RunCommand("sudo", "systemctl", "restart", service)
	if err != nil {
		return err
	}

	return nil
}

// setupSSL sets up SSL certificate using Let's Encrypt
func setupSSL(config *config.Config) error {
	utils.PrintHeader("Setting up SSL Certificate")
	utils.PrintWarning("Make sure your domain DNS is pointing to this server before running SSL setup")
	
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to setup SSL certificate now? (y/n): ")
	setupSSL, _ := reader.ReadString('\n')
	setupSSL = strings.TrimSpace(setupSSL)

	if setupSSL == "y" || setupSSL == "Y" {
		// Use Certbot to obtain and install SSL certificate
		err := utils.RunCommand("sudo", "certbot", "--nginx", "-d", config.Domain, "-d", "www."+config.Domain)
		if err != nil {
			utils.PrintError("Failed to install SSL certificate")
			utils.PrintWarning("You can try again later with: sudo certbot --nginx -d " + config.Domain + " -d www." + config.Domain)
		} else {
			utils.PrintStatus("SSL certificate installed successfully")

			// Setup auto-renewal via cron job
			err = utils.RunCommand("sudo", "bash", "-c", "echo \"0 12 * * * /usr/bin/certbot renew --quiet\" | sudo crontab -")
			if err != nil {
				return err
			}

			utils.PrintStatus("SSL auto-renewal configured")
		}
	} else {
		utils.PrintWarning("SSL certificate setup skipped")
		utils.PrintWarning("You can set it up later with: sudo certbot --nginx -d " + config.Domain + " -d www." + config.Domain)
	}

	return nil
}

// createServerInfo creates a server information file
func createServerInfo(config *config.Config) error {
	utils.PrintHeader("Creating Server Information File")
	utils.PrintStatus("Saving server information to file...")

	// Generate server information content
	serverInfo := templates.GetServerInfoContent(
		config.Domain,
		config.WebRoot,
		config.DBName,
		config.DBUser,
		config.SSHPort,
		os.Getenv("USER"),
	)

	// Write server information to file with restricted permissions
	err := os.WriteFile("/home/"+os.Getenv("USER")+"/server_info.txt", []byte(serverInfo), 0600)
	if err != nil {
		return err
	}

	return nil
}