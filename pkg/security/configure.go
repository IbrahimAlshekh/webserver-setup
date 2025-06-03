package security

import (
	"os"

	"laravel-setup/pkg/config"
	"laravel-setup/pkg/templates"
	"laravel-setup/pkg/utils"
)

// Configure sets up security measures for the server
func Configure(config *config.Config) error {
	// Configure firewall
	if err := configureFirewall(config); err != nil {
		return err
	}

	// Configure fail2ban
	if err := configureFail2ban(config); err != nil {
		return err
	}

	// Configure SSH
	if err := configureSSH(config); err != nil {
		return err
	}

	utils.PrintHeader("Security Configuration Complete")
	utils.PrintStatus("Firewall, fail2ban, and SSH security have been configured")

	return nil
}

// configureFirewall sets up UFW firewall with appropriate rules
func configureFirewall(config *config.Config) error {
	utils.PrintHeader("Configuring UFW Firewall")
	utils.PrintStatus("Setting up firewall rules...")

	// Set default policies
	err := utils.RunCommand("sudo", "ufw", "default", "deny", "incoming")
	if err != nil {
		return err
	}

	err = utils.RunCommand("sudo", "ufw", "default", "allow", "outgoing")
	if err != nil {
		return err
	}

	// Allow SSH on custom port
	err = utils.RunCommand("sudo", "ufw", "allow", config.SSHPort+"/tcp")
	if err != nil {
		return err
	}

	// Allow HTTP
	err = utils.RunCommand("sudo", "ufw", "allow", "80/tcp")
	if err != nil {
		return err
	}

	// Allow HTTPS
	err = utils.RunCommand("sudo", "ufw", "allow", "443/tcp")
	if err != nil {
		return err
	}

	// Enable firewall
	utils.PrintStatus("Enabling firewall...")
	err = utils.RunCommand("sudo", "ufw", "--force", "enable")
	if err != nil {
		return err
	}

	utils.PrintStatus("Firewall configured and enabled successfully")
	err = utils.RunCommand("sudo", "ufw", "status")
	if err != nil {
		return err
	}

	return nil
}

// configureFail2ban sets up fail2ban to protect against brute force attacks
func configureFail2ban(config *config.Config) error {
	utils.PrintHeader("Configuring Fail2ban")
	utils.PrintStatus("Setting up fail2ban for intrusion prevention...")

	// Copy default configuration
	err := utils.RunCommand("sudo", "cp", "/etc/fail2ban/jail.conf", "/etc/fail2ban/jail.local")
	if err != nil {
		return err
	}

	// Generate fail2ban configuration
	fail2banConfig := templates.GetFail2banConfig(config.SSHPort)

	// Write fail2ban configuration to file
	err = os.WriteFile("fail2ban_custom.conf", []byte(fail2banConfig), 0644)
	if err != nil {
		return err
	}

	// Move fail2ban configuration to jail.d directory
	err = utils.RunCommand("sudo", "mv", "fail2ban_custom.conf", "/etc/fail2ban/jail.d/custom.conf")
	if err != nil {
		return err
	}

	// Enable fail2ban to start on boot
	err = utils.RunCommand("sudo", "systemctl", "enable", "fail2ban")
	if err != nil {
		return err
	}

	// Restart fail2ban to apply changes
	err = utils.RunCommand("sudo", "systemctl", "restart", "fail2ban")
	if err != nil {
		return err
	}

	utils.PrintStatus("Fail2ban configured and started successfully")
	return nil
}

// configureSSH hardens SSH configuration for better security
func configureSSH(config *config.Config) error {
	utils.PrintHeader("Configuring SSH Security")
	utils.PrintStatus("Hardening SSH configuration...")

	// Backup original SSH configuration
	err := utils.RunCommand("sudo", "cp", "/etc/ssh/sshd_config", "/etc/ssh/sshd_config.backup")
	if err != nil {
		return err
	}

	// Generate SSH configuration
	sshConfig := templates.GetSSHConfig(config.SSHPort)

	// Write SSH configuration to file
	err = os.WriteFile("ssh_security.conf", []byte(sshConfig), 0644)
	if err != nil {
		return err
	}

	// Create sshd_config.d directory if it doesn't exist
	err = utils.RunCommand("sudo", "mkdir", "-p", "/etc/ssh/sshd_config.d")
	if err != nil {
		return err
	}

	// Move SSH configuration to sshd_config.d directory
	err = utils.RunCommand("sudo", "mv", "ssh_security.conf", "/etc/ssh/sshd_config.d/security.conf")
	if err != nil {
		return err
	}

	// Restart SSH service to apply changes
	utils.PrintStatus("Restarting SSH service to apply changes...")
	err = utils.RunCommand("sudo", "systemctl", "restart", "ssh")
	if err != nil {
		return err
	}

	utils.PrintStatus("SSH security configured successfully")
	utils.PrintWarning("SSH port has been changed to: " + config.SSHPort)
	utils.PrintWarning("Make sure to update your SSH client configuration")

	return nil
}