package utils

import (
	"fmt"
)

// Color codes for terminal output
const (
	ColorRed    = "\033[0;31m"
	ColorGreen  = "\033[0;32m"
	ColorYellow = "\033[1;33m"
	ColorBlue   = "\033[0;34m"
	ColorReset  = "\033[0m"
)

// PrintStatus prints a status message with green color
func PrintStatus(message string) {
	fmt.Printf("%s[INFO]%s %s\n", ColorGreen, ColorReset, message)
}

// PrintWarning prints a warning message with yellow color
func PrintWarning(message string) {
	fmt.Printf("%s[WARNING]%s %s\n", ColorYellow, ColorReset, message)
}

// PrintError prints an error message with red color
func PrintError(message string) {
	fmt.Printf("%s[ERROR]%s %s\n", ColorRed, ColorReset, message)
}

// PrintHeader prints a header with blue color
func PrintHeader(message string) {
	fmt.Printf("%s================================%s\n", ColorBlue, ColorReset)
	fmt.Printf("%s%s%s\n", ColorBlue, message, ColorReset)
	fmt.Printf("%s================================%s\n", ColorBlue, ColorReset)
}

// PrintInformation prints information with green color
func PrintInformation(message string) {
	fmt.Printf("%s================================%s\n", ColorGreen, ColorReset)
	fmt.Printf("%s%s%s\n", ColorGreen, message, ColorReset)
	fmt.Printf("%s================================%s\n", ColorGreen, ColorReset)
}
