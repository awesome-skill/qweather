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
	forecastLocation string
	forecastDays     int
)

// forecastCmd represents the forecast command
var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get weather forecast",
	Long: `Get daily weather forecast for a specified location.

Location can be:
  - City name (e.g., "北京", "Beijing")
  - Location ID (e.g., "101010100")
  - Coordinates (e.g., "116.41,39.92")

Forecast days can be: 3, 7, 10, 15, or 30

Examples:
  weather forecast --location "北京" --days 3
  weather forecast --location "101010100" --days 7 --format table
  weather forecast --location "116.41,39.92" --days 15 --format json`,
	RunE: runForecast,
}

func init() {
	rootCmd.AddCommand(forecastCmd)
	forecastCmd.Flags().StringVarP(&forecastLocation, "location", "l", "", "Location to query (required)")
	forecastCmd.Flags().IntVarP(&forecastDays, "days", "d", 3, "Forecast days: 3, 7, 10, 15, or 30")
	forecastCmd.MarkFlagRequired("location")
}

func runForecast(cmd *cobra.Command, args []string) error {
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

	// Get weather forecast
	forecastData, err := client.GetDailyForecast(ctx, forecastLocation, forecastDays)
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

	if err := formatter.FormatWeatherDaily(forecastData, os.Stdout); err != nil {
		printError(err)
		return err
	}

	return nil
}
