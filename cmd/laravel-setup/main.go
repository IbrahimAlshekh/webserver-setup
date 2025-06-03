package main

import (
	"bufio"
	"flag"
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

// getSkipStatus returns a string indicating whether a step will be skipped
func getSkipStatus(skip bool) string {
	if skip {
		return " (will be skipped)"
	}
	return ""
}

func main() {
	// Define command-line flags
	cleanupFlag := flag.Bool("cleanup", false, "Clean up temporary files created during the setup process")
	configPathFlag := flag.String("config-path", "", "Path to the configuration file (default: ~/config.toml)")

	// Module selection flags
	skipSystemUpdateFlag := flag.Bool("skip-system-update", false, "Skip system update step")
	skipEssentialsFlag := flag.Bool("skip-essentials", false, "Skip installing essential packages")
	skipPHPFlag := flag.Bool("skip-php", false, "Skip PHP installation")
	skipMySQLFlag := flag.Bool("skip-mysql", false, "Skip MySQL installation")
	skipNginxFlag := flag.Bool("skip-nginx", false, "Skip Nginx installation")
	skipSecurityFlag := flag.Bool("skip-security", false, "Skip security configuration")
	skipLaravelFlag := flag.Bool("skip-laravel", false, "Skip Laravel setup")
	skipServicesFlag := flag.Bool("skip-services", false, "Skip services configuration")

	flag.Parse()

	// Check if a cleanup flag is set
	if *cleanupFlag {
		err := utils.CleanupTempFiles()
		if err != nil {
			utils.PrintError("Failed to clean up temporary files: " + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}

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
	cfg, err := config.InitConfig(*configPathFlag)
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
	utils.PrintStatus("1. System update" + getSkipStatus(*skipSystemUpdateFlag))
	utils.PrintStatus("2. Installing essential packages" + getSkipStatus(*skipEssentialsFlag))
	utils.PrintStatus("3. Installing PHP 8.3 and extensions" + getSkipStatus(*skipPHPFlag))
	utils.PrintStatus("4. Installing and configuring MySQL" + getSkipStatus(*skipMySQLFlag))
	utils.PrintStatus("5. Installing and configuring Nginx" + getSkipStatus(*skipNginxFlag))
	utils.PrintStatus("6. Configuring security (firewall, fail2ban, SSH)" + getSkipStatus(*skipSecurityFlag))
	utils.PrintStatus("7. Setting up Laravel application" + getSkipStatus(*skipLaravelFlag))
	utils.PrintStatus("8. Configuring and starting services" + getSkipStatus(*skipServicesFlag))
	utils.PrintStatus("")
	utils.PrintWarning("This process may take some time. Please be patient.")
	utils.PrintWarning("You will be prompted for input at certain stages.")
	utils.PrintStatus("")
	utils.PrintStatus("Note: You can skip any step by using the corresponding command-line flag:")
	utils.PrintStatus("  --skip-system-update, --skip-essentials, --skip-php, --skip-mysql,")
	utils.PrintStatus("  --skip-nginx, --skip-security, --skip-laravel, --skip-services")
	utils.PrintStatus("")

	fmt.Print("Press Enter to begin the setup process...")
	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		return
	}

	// Apply skip flags from config file if not overridden by command-line flags
	skipSystemUpdate := *skipSystemUpdateFlag || cfg.SkipSystemUpdate
	skipEssentials := *skipEssentialsFlag || cfg.SkipEssentials
	skipPHP := *skipPHPFlag || cfg.SkipPHP
	skipMySQL := *skipMySQLFlag || cfg.SkipMySQL
	skipNginx := *skipNginxFlag || cfg.SkipNginx
	skipSecurity := *skipSecurityFlag || cfg.SkipSecurity
	skipLaravel := *skipLaravelFlag || cfg.SkipLaravel
	skipServices := *skipServicesFlag || cfg.SkipServices

	// Run each step of the setup process, skipping those that the user has opted to skip
	if !skipSystemUpdate {
		runConfigStep("System Update", system.Update, nil)
	} else {
		utils.PrintStatus("Skipping System Update step as requested")
	}

	if !skipEssentials {
		runConfigStep("Install Essentials", system.InstallEssentials, nil)
	} else {
		utils.PrintStatus("Skipping Install Essentials step as requested")
	}

	if !skipPHP {
		runConfigStep("Install PHP", php.Install, nil)
	} else {
		utils.PrintStatus("Skipping PHP Installation step as requested")
	}

	if !skipMySQL {
		runConfigStep("Install MySQL", mysql.Install, cfg)
	} else {
		utils.PrintStatus("Skipping MySQL Installation step as requested")
	}

	if !skipNginx {
		runConfigStep("Install Nginx", nginx.Install, cfg)
	} else {
		utils.PrintStatus("Skipping Nginx Installation step as requested")
	}

	if !skipSecurity {
		runConfigStep("Configure Security", security.Configure, cfg)
	} else {
		utils.PrintStatus("Skipping Security Configuration step as requested")
	}

	if !skipLaravel {
		runConfigStep("Setup Laravel", laravel.Setup, cfg)
	} else {
		utils.PrintStatus("Skipping Laravel Setup step as requested")
	}

	if !skipServices {
		runConfigStep("Configure Services", services.Configure, cfg)
	} else {
		utils.PrintStatus("Skipping Services Configuration step as requested")
	}

	// Final message
	utils.PrintHeader("Setup Complete!")
	utils.PrintStatus("Laravel production server has been successfully set up")
	utils.PrintStatus("Server information has been saved to: /home/" + os.Getenv("USER") + "/server_info.txt")
	utils.PrintStatus("MySQL credentials have been saved to: /home/" + os.Getenv("USER") + "/mysql_credentials.txt")
	utils.PrintWarning("Remember to:")
	utils.PrintWarning("1. Point your domain DNS to this server")
	utils.PrintWarning("2. Set up SSL certificate if you haven't already")
	utils.PrintWarning("3. Change SSH port in your SSH client to: " + cfg.SSHPort)
	utils.PrintStatus("")
	utils.PrintStatus("You can clean up temporary files by running: laravel-setup -cleanup")

	fmt.Printf("%sYour Laravel production server is ready!%s\n", utils.ColorGreen, utils.ColorReset)
}
