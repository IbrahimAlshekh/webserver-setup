package utils

import (
	"os"
	"os/exec"
	"strings"
)

// RunCommand executes a shell command and returns the error if any
// Streams command output to stdout and stderr for real-time feedback
func RunCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunCommandWithOutput executes a shell command and returns the output and error
// Useful when you need to capture the output for processing
func RunCommandWithOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.Output()
	return strings.TrimSpace(string(output)), err
}

// GenerateRandomPassword generates a random password using OpenSSL
// Falls back to a default password if OpenSSL fails
func GenerateRandomPassword() string {
	output, err := RunCommandWithOutput("openssl", "rand", "-base64", "16")
	if err != nil {
		// If OpenSSL fails, return a default password
		// This is not ideal but prevents the script from failing
		return "defaultpassword"
	}
	return output
}

// CheckNotRoot checks if the script is run as root
// Returns true if not running as root, false otherwise
func CheckNotRoot() bool {
	output, err := RunCommandWithOutput("id", "-u")
	if err != nil {
		PrintError("Failed to check user ID")
		return false
	}

	if output == "0" {
		PrintError("This script should not be run as root for security reasons")
		PrintWarning("Please create a regular user first:")
		PrintWarning("  adduser username")
		PrintWarning("  usermod -aG sudo username")
		PrintWarning("  su - username")
		PrintWarning("Then run this script as that user")
		return false
	}
	return true
}

// CheckSudoPrivileges checks if the user has sudo privileges
// Returns true if the user has sudo privileges, false otherwise
func CheckSudoPrivileges() bool {
	cmd := exec.Command("sudo", "-n", "true")
	err := cmd.Run()
	if err != nil {
		PrintError("This user doesn't have sudo privileges")
		PrintWarning("Please add this user to the sudo group:")
		PrintWarning("  sudo usermod -aG sudo $USER")
		return false
	}
	return true
}