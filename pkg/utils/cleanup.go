package utils

import (
	"os"
)

// CleanupTempFiles removes all temporary files created during the setup process
func CleanupTempFiles() error {
	PrintHeader("Cleaning up temporary files")

	// List of temporary files that might be created during the setup process
	tempFiles := []string{
		"mysql_config.sql",     // Created in mysql/install.go
		"nginx_site.conf",      // Created in nginx/install.go
		"opcache.ini",          // Created in php/install.go
		"laravel-worker.conf",  // Created in laravel/setup.go
		"fail2ban_custom.conf", // Created in security/configure.go
		"ssh_security.conf",    // Created in security/configure.go
	}

	// Remove each temporary file if it exists
	for _, file := range tempFiles {
		if _, err := os.Stat(file); err == nil {
			PrintStatus("Removing temporary file: " + file)
			err := os.Remove(file)
			if err != nil {
				PrintError("Failed to remove temporary file: " + file)
				return err
			}
		}
	}

	PrintStatus("All temporary files have been cleaned up")
	return nil
}
