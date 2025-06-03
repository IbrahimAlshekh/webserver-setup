package system

import (
	"laravel-setup/pkg/config"
	"laravel-setup/pkg/utils"
)

// Update updates the system packages
// This ensures the system has the latest security patches and package information
func Update(_ *config.Config) error {
	utils.PrintStatus("Updating system packages...")

	// Update package lists
	err := utils.RunCommand("sudo", "apt", "update")
	if err != nil {
		return err
	}

	// Upgrade installed packages
	err = utils.RunCommand("sudo", "apt", "upgrade", "-y")
	if err != nil {
		return err
	}

	utils.PrintStatus("System update completed successfully")
	return nil
}
