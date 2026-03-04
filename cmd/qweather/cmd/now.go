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
	nowLocation string
	nowCity     string
)

// nowCmd represents the now command
var nowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get current weather",
	Long: `Get current weather conditions for a specified location.

Location can be specified using either --location or --city:
  - Location ID (e.g., "101010100")
  - Coordinates (e.g., "116.41,39.92")
  - City name (use --city, e.g., "北京", "Beijing")

Examples:
  qweather now --location "101010100"
  qweather now --location "116.41,39.92"
  qweather now --city "北京"
  qweather now --city "beijing" --format json`,
	RunE: runNow,
}

func init() {
	rootCmd.AddCommand(nowCmd)
	nowCmd.Flags().StringVarP(&nowLocation, "location", "l", "", "Location ID or coordinates (required, mutually exclusive with --city)")
	nowCmd.Flags().StringVarP(&nowCity, "city", "c", "", "City name (auto-resolve to location ID, required, mutually exclusive with --location)")
	nowCmd.MarkFlagsMutuallyExclusive("location", "city")
}

func runNow(cmd *cobra.Command, args []string) error {
	// Validate flags
	if nowLocation == "" && nowCity == "" {
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
	location := nowLocation
	if nowCity != "" {
		searchData, err := client.SearchCity(ctx, nowCity)
		if err != nil {
			printError(fmt.Errorf("failed to search city %q: %w", nowCity, err))
			return err
		}
		if len(searchData.Location) == 0 {
			errMsg := fmt.Sprintf("no location found for city: %q\n", nowCity)
			errMsg += "Tip: try a different spelling or use location ID directly"
			printError(fmt.Errorf("%s", errMsg))
			return fmt.Errorf("city not found: %s", nowCity)
		}
		location = searchData.Location[0].ID
		if verboseFlag {
			fmt.Printf("Resolved %s to location ID: %s\n", nowCity, location)
		}
	}

	// Get current weather
	weatherData, err := client.GetNowWeather(ctx, location)
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
