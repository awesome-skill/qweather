package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	formatFlag  string
	verboseFlag bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "weather",
	Short: "Weather CLI powered by QWeather API",
	Long: `A command-line tool to query weather information using QWeather API.

Supports current weather, weather forecast, and city search.
Requires QWeather API key to be configured.`,
	Version: "0.1.0",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&formatFlag, "format", "f", "text", "Output format: text, json, table")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Verbose output")
}

// printError prints error message to stderr
func printError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}
