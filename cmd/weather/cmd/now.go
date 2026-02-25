package cmd

import (
	"context"
	"os"
	"time"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
	"github.com/pangu-studio/awesome-skills/internal/config"
	"github.com/pangu-studio/awesome-skills/internal/output"
	"github.com/spf13/cobra"
)

var (
	nowLocation string
)

// nowCmd represents the now command
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get current weather",
	Long: `Get current weather conditions for a specified location.

Location can be:
  - City name (e.g., "北京", "Beijing")
  - Location ID (e.g., "101010100")
  - Coordinates (e.g., "116.41,39.92")

Examples:
  weather now --location "北京"
  weather now --location "101010100"
  weather now --location "116.41,39.92" --format json`,
	RunE: runNow,
}

func init() {
	rootCmd.AddCommand(nowCmd)
	nowCmd.Flags().StringVarP(&nowLocation, "location", "l", "", "Location to query (required)")
	nowCmd.MarkFlagRequired("location")
}

func runNow(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		printError(err)
		return err
	}

	// Create API client
	client := qweather.NewClient(cfg.QWeather.APIKey, cfg.QWeather.APIHost)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get current weather
	weatherData, err := client.GetNowWeather(ctx, nowLocation)
	if err != nil {
		printError(err)
		return err
	}

	// Format and print output
	formatter, err := output.NewFormatter(formatFlag)
	if err != nil {
		printError(err)
		return err
	}

	if err := formatter.FormatWeatherNow(weatherData, os.Stdout); err != nil {
		printError(err)
		return err
	}

	return nil
}
