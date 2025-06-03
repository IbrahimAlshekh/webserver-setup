package laravel

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Setup sets up the Laravel application
func Setup(config *config.Config) error {
	utils.PrintHeader("Setting Up Laravel Application")

	// Configure Git and SSH for deployment
	if err := configureGit(); err != nil {
		return err
	}

	// Clone the repository
	if err := cloneRepository(config); err != nil {
		return err
	}

	// Install Composer dependencies
	if err := installDependencies(config); err != nil {
		return err
	}

	// Configure Laravel environment
	if err := configureEnvironment(config); err != nil {
		return err
	}

	// Configure Supervisor for Laravel Queue
	if err := configureSupervisor(config); err != nil {
		return err
	}

	utils.PrintHeader("Laravel Application Setup Complete")
	utils.PrintStatus("Laravel application has been set up successfully at " + config.WebRoot)
	utils.PrintStatus("You can now access your application at http://" + config.Domain)
	utils.PrintWarning("Remember to set up SSL certificate for HTTPS access")

	return nil
}

// configureGit configures Git and SSH for deployment
func configureGit() error {
	utils.PrintHeader("Configuring Git for Deployment")
	utils.PrintStatus("Setting up SSH for Git...")

	// Create an SSH directory if it doesn't exist
	err := utils.RunCommand("mkdir", "-p", "~/.ssh")
	if err != nil {
		return err
	}

	// Set proper permissions for SSH directory
	err = utils.RunCommand("chmod", "700", "~/.ssh")
	if err != nil {
		return err
	}

	// Add GitHub to known hosts to prevent SSH prompts
	err = utils.RunCommand("ssh-keyscan", "-H", "github.com", ">>", "~/.ssh/known_hosts")
	if err != nil {
		return err
	}

	utils.PrintWarning("Please add your SSH public key to GitHub before proceeding")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Do you want to generate new ssh key? (y/n)")
	generateSSHKey, _ := reader.ReadString('\n')
	generateSSHKey = strings.TrimSpace(generateSSHKey)

	if generateSSHKey == "y" || generateSSHKey == "Y" {
		// Generate a new SSH key with a strong algorithm
		err = utils.RunInteractiveCommand("ssh-keygen", "-t", "ed25519", "-C", "deployment@"+os.Getenv("USER"))
		if err != nil {
			return err
		}
		utils.PrintStatus("SSH keygen is generated")
	}

	utils.PrintWarning("Add the public key (~/.ssh/id_ed25519.pub) to your GitHub account")

	fmt.Print("Print the public key? (y/n)")
	printPublicKey, _ := reader.ReadString('\n')
	printPublicKey = strings.TrimSpace(printPublicKey)

	if printPublicKey == "y" || printPublicKey == "Y" {
		// Read and display the public key
		pubKeyBytes, err := os.ReadFile(os.Getenv("HOME") + "/.ssh/id_ed25519.pub")
		if err != nil {
			return err
		}
		utils.PrintInformation(string(pubKeyBytes))
	}

	fmt.Print("Press Enter when you've added your SSH key to GitHub..")
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	return nil
}

// cloneRepository clones the Laravel repository
func cloneRepository(config *config.Config) error {
	utils.PrintHeader("Cloning Laravel Repository")
	utils.PrintStatus("Cloning repository to " + config.WebRoot + "...")

	// Remove the existing directory if it exists
	if _, err := os.Stat(config.WebRoot); err == nil {
		utils.PrintWarning("Directory " + config.WebRoot + " already exists. Removing...")
		err = utils.RunCommand("sudo", "rm", "-rf", config.WebRoot)
		if err != nil {
			return err
		}
	}

	// Clone the repository
	if config.RepoURL == "" {
		utils.PrintError("Repository URL cannot be empty")
		return fmt.Errorf("repository URL cannot be empty")
	}

	err := utils.RunCommand("sudo", "git", "clone", config.RepoURL, config.WebRoot)
	if err != nil {
		return err
	}

	// Set proper ownership and permissions
	utils.PrintStatus("Setting proper ownership and permissions...")
	err = utils.RunCommand("sudo", "chown", "-R", os.Getenv("USER")+":"+config.WebUser, config.WebRoot)
	if err != nil {
		return err
	}

	// Set directory permissions
	err = utils.RunCommand("sudo", "chmod", "-R", "755", config.WebRoot)
	if err != nil {
		return err
	}

	// Set storage directory permissions (needs to be writable by web server)
	err = utils.RunCommand("sudo", "chmod", "-R", "775", config.WebRoot+"/storage")
	if err != nil {
		return err
	}

	// Create bootstrap/cache directory if it doesn't exist
	err = utils.RunCommand("sudo", "mkdir", "-p", config.WebRoot+"/bootstrap/cache")
	if err != nil {
		return err
	}

	// Set bootstrap/cache directory permissions (needs to be writable by web server)
	err = utils.RunCommand("sudo", "chmod", "-R", "775", config.WebRoot+"/bootstrap/cache")
	if err != nil {
		return err
	}

	return nil
}

// installDependencies installs Composer dependencies
func installDependencies(config *config.Config) error {
	utils.PrintHeader("Installing Composer Dependencies")
	utils.PrintStatus("Installing Composer dependencies...")

	// Change to web root directory
	err := os.Chdir(config.WebRoot)
	if err != nil {
		return err
	}

	// Install Composer dependencies with optimizations for production
	err = utils.RunCommand("composer", "install", "--no-dev", "--optimize-autoloader")
	if err != nil {
		return err
	}

	return nil
}

// configureEnvironment configures the Laravel environment
func configureEnvironment(config *config.Config) error {
	utils.PrintHeader("Configuring Laravel Environment")
	utils.PrintStatus("Setting up .env file...")

	// Copy .env.example to .env if it exists
	if _, err := os.Stat(".env.example"); err == nil {
		err = utils.RunCommand("cp", ".env.example", ".env")
		if err != nil {
			return err
		}
	} else {
		utils.PrintWarning("No .env.example file found. Creating empty .env file...")
		err = utils.RunCommand("touch", ".env")
		if err != nil {
			return err
		}
	}

	// Update the .env file with database credentials
	err := utils.RunCommand("sed", "-i", "s/DB_DATABASE=laravel/DB_DATABASE="+config.DBName+"/", ".env")
	if err != nil {
		return err
	}

	err = utils.RunCommand("sed", "-i", "s/DB_USERNAME=root/DB_USERNAME="+config.DBUser+"/", ".env")
	if err != nil {
		return err
	}

	err = utils.RunCommand("sed", "-i", "s|^DB_PASSWORD=.*|DB_PASSWORD=\""+config.DBPassword+"\"|", ".env")
	if err != nil {
		return err
	}

	// Configure Redis for caching, sessions, and queue
	err = utils.RunCommand("sed", "-i", "s/CACHE_DRIVER=file/CACHE_DRIVER=redis/", ".env")
	if err != nil {
		return err
	}

	err = utils.RunCommand("sed", "-i", "s/SESSION_DRIVER=file/SESSION_DRIVER=redis/", ".env")
	if err != nil {
		return err
	}

	err = utils.RunCommand("sed", "-i", "s/QUEUE_CONNECTION=sync/QUEUE_CONNECTION=redis/", ".env")
	if err != nil {
		return err
	}

	// Generate an application key
	utils.PrintStatus("Generating application key...")
	err = utils.RunCommand("php", "artisan", "key:generate")
	if err != nil {
		return err
	}

	// Run migrations
	utils.PrintStatus("Running database migrations...")
	err = utils.RunCommand("php", "artisan", "migrate", "--force")
	if err != nil {
		return err
	}

	return nil
}

// configureSupervisor configures Supervisor for Laravel Queue
func configureSupervisor(config *config.Config) error {
	utils.PrintHeader("Configuring Supervisor for Laravel Queue")
	utils.PrintStatus("Setting up Supervisor for Laravel queue workers...")

	// Generate Supervisor configuration
	supervisorConfig := templates.GetSupervisorConfig(config.WebRoot, config.WebUser)

	// Write Supervisor configuration to file
	err := os.WriteFile("laravel-worker.conf", []byte(supervisorConfig), 0644)
	if err != nil {
		return err
	}

	// Move Supervisor configuration to conf.d directory
	err = utils.RunCommand("sudo", "mv", "laravel-worker.conf", "/etc/supervisor/conf.d/laravel-worker.conf")
	if err != nil {
		return err
	}

	// Reload Supervisor configuration
	err = utils.RunCommand("sudo", "supervisorctl", "reread")
	if err != nil {
		return err
	}

	// Update Supervisor to apply changes
	err = utils.RunCommand("sudo", "supervisorctl", "update")
	if err != nil {
		return err
	}

	// Start Laravel workers
	err = utils.RunCommand("sudo", "supervisorctl", "start", "laravel-worker:*")
	if err != nil {
		return err
	}

	return nil
}
