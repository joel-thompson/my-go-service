package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	serverURL string
	format    string
	verbose   bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mycli",
	Short: "CLI tool for testing my-go-service API",
	Long: `A command line interface for testing and interacting with the my-go-service API.
	
This tool provides convenient commands for testing all API endpoints
without having to write curl commands manually.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "http://localhost:8080", "API server URL")
	rootCmd.PersistentFlags().StringVar(&format, "format", "pretty", "Output format (pretty|json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(helloCmd)
	rootCmd.AddCommand(itemsCmd)
}

// Helper function to handle verbose output
func verboseLog(message string) {
	if verbose {
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", message)
	}
}
