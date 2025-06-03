package main

import (
	"bufio"
	"fmt"
	"os"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/laravel"
	"laravel-setup/pkg/mysql"
	"laravel-setup/pkg/nginx"
	"laravel-setup/pkg/php"
	"laravel-setup/pkg/security"
	"laravel-setup/pkg/services"
	"laravel-setup/pkg/system"
	"laravel-setup/pkg/utils"
)

// runConfigStep runs a function that takes a config parameter and checks its exit status
func runConfigStep(name string, fn func(*config.Config) error, cfg *config.Config) {
	utils.PrintHeader("Running " + name)

	err := fn(cfg)
	if err != nil {
		utils.PrintError("Step " + name + " failed: " + err.Error())
		os.Exit(1)
	}
}

func main() {
	// Print a welcome message
	utils.PrintHeader("Laravel Production Server Setup")

	// Check if running as root and if the user has sudo privileges
	// These are security checks to ensure the script is run correctly
	if !utils.CheckNotRoot() {
		os.Exit(1)
	}

	if !utils.CheckSudoPrivileges() {
		os.Exit(1)
	}

	// Initialize configuration
	cfg, err := config.InitConfig()
	if err != nil {
		utils.PrintError("Failed to initialize configuration: " + err.Error())
		os.Exit(1)
	}

	utils.PrintStatus("Setting up server for domain: " + cfg.Domain)
	utils.PrintStatus("Running as user: " + os.Getenv("USER"))

	// Main setup process
	utils.PrintHeader("Starting Laravel Server Setup Process")
	utils.PrintStatus("This script will set up a complete Laravel production server")
	utils.PrintStatus("The setup process is divided into several steps:")
	utils.PrintStatus("1. System update")
	utils.PrintStatus("2. Installing essential packages")
	utils.PrintStatus("3. Installing PHP 8.3 and extensions")
	utils.PrintStatus("4. Installing and configuring MySQL")
	utils.PrintStatus("5. Installing and configuring Nginx")
	utils.PrintStatus("6. Configuring security (firewall, fail2ban, SSH)")
	utils.PrintStatus("7. Setting up Laravel application")
	utils.PrintStatus("8. Configuring and starting services")
	utils.PrintStatus("")
	utils.PrintWarning("This process may take some time. Please be patient.")
	utils.PrintWarning("You will be prompted for input at certain stages.")
	utils.PrintStatus("")

	fmt.Print("Press Enter to begin the setup process...")
	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		return
	}

	// Run each step of the setup process
	runConfigStep("System Update", system.Update, nil)
	runConfigStep("Install Essentials", system.InstallEssentials, nil)
	runConfigStep("Install PHP", php.Install, nil)
	runConfigStep("Install MySQL", mysql.Install, cfg)
	runConfigStep("Install Nginx", nginx.Install, cfg)
	runConfigStep("Configure Security", security.Configure, cfg)
	runConfigStep("Setup Laravel", laravel.Setup, cfg)
	runConfigStep("Configure Services", services.Configure, cfg)

	// Final message
	utils.PrintHeader("Setup Complete!")
	utils.PrintStatus("Laravel production server has been successfully set up")
	utils.PrintStatus("Server information has been saved to: /home/" + os.Getenv("USER") + "/server_info.txt")
	utils.PrintStatus("MySQL credentials have been saved to: /home/" + os.Getenv("USER") + "/mysql_credentials.txt")
	utils.PrintWarning("Remember to:")
	utils.PrintWarning("1. Point your domain DNS to this server")
	utils.PrintWarning("2. Set up SSL certificate if you haven't already")
	utils.PrintWarning("3. Change SSH port in your SSH client to: " + cfg.SSHPort)

	fmt.Printf("%sYour Laravel production server is ready!%s\n", utils.ColorGreen, utils.ColorReset)
}
