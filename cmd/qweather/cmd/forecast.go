package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
	"github.com/pangu-studio/awesome-skills/internal/config"
	"github.com/pangu-studio/awesome-skills/internal/output"
	"github.com/spf13/cobra"
)

var (
	forecastLocation string
	forecastCity     string
	forecastDays     int
)

// forecastCmd represents the forecast command
var forecastCmd = &cobra.Command{
	Use:   "forecast",
	Short: "Get weather forecast",
	Long: `Get daily weather forecast for a specified location.

Location can be specified using either --location or --city:
  - Location ID (e.g., "101010100")
  - Coordinates (e.g., "116.41,39.92")
  - City name (use --city, e.g., "北京", "Shanghai")

Forecast days can be: 3, 7, 10, 15, or 30

Examples:
  qweather forecast --location "101010100" --days 7
  qweather forecast --location "116.41,39.92" --days 15 --format json
  qweather forecast --city "北京" --days 3
  qweather forecast --city "shanghai" --days 7 --format table`,
	RunE: runForecast,
}

func init() {
	rootCmd.AddCommand(forecastCmd)
	forecastCmd.Flags().StringVarP(&forecastLocation, "location", "l", "", "Location ID or coordinates (required, mutually exclusive with --city)")
	forecastCmd.Flags().StringVarP(&forecastCity, "city", "c", "", "City name (auto-resolve to location ID, required, mutually exclusive with --location)")
	forecastCmd.Flags().IntVarP(&forecastDays, "days", "d", 3, "Forecast days: 3, 7, 10, 15, or 30")
	forecastCmd.MarkFlagsMutuallyExclusive("location", "city")
}

func runForecast(cmd *cobra.Command, args []string) error {
	// Validate flags
	if forecastLocation == "" && forecastCity == "" {
		return fmt.Errorf("either --location or --city is required (mutually exclusive)")
	}

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

	// Resolve location: if using --city, search for city ID first
	location := forecastLocation
	if forecastCity != "" {
		searchData, err := client.SearchCity(ctx, forecastCity)
		if err != nil {
			printError(fmt.Errorf("failed to search city %q: %w", forecastCity, err))
			return err
		}
		if len(searchData.Location) == 0 {
			errMsg := fmt.Sprintf("no location found for city: %q\n", forecastCity)
			errMsg += "Tip: try a different spelling or use location ID directly"
			printError(fmt.Errorf("%s", errMsg))
			return fmt.Errorf("city not found: %s", forecastCity)
		}
		location = searchData.Location[0].ID
		if verboseFlag {
			fmt.Printf("Resolved %s to location ID: %s\n", forecastCity, location)
		}
	}

	// Get weather forecast
	forecastData, err := client.GetDailyForecast(ctx, location, forecastDays)
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
