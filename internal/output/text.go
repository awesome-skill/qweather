package output

import (
	"fmt"
	"io"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
)

// TextFormatter formats output as human-readable text
type TextFormatter struct{}

// FormatWeatherNow formats current weather as text
func (f *TextFormatter) FormatWeatherNow(data *qweather.WeatherNowResponse, w io.Writer) error {
	now := data.Now

	fmt.Fprintf(w, "Current Weather\n")
	fmt.Fprintf(w, "===============\n\n")
	fmt.Fprintf(w, "Observation Time: %s\n", now.ObsTime)
	fmt.Fprintf(w, "Condition:        %s\n", now.Text)
	fmt.Fprintf(w, "Temperature:      %s°C\n", now.Temp)
	fmt.Fprintf(w, "Feels Like:       %s°C\n", now.FeelsLike)
	fmt.Fprintf(w, "Humidity:         %s%%\n", now.Humidity)
	fmt.Fprintf(w, "Wind:             %s %s (%s km/h)\n", now.WindDir, now.WindScale, now.WindSpeed)
	fmt.Fprintf(w, "Pressure:         %s hPa\n", now.Pressure)
	fmt.Fprintf(w, "Visibility:       %s km\n", now.Vis)

	if now.Precip != "" && now.Precip != "0.0" {
		fmt.Fprintf(w, "Precipitation:    %s mm\n", now.Precip)
	}

	fmt.Fprintf(w, "\nUpdate Time: %s\n", data.UpdateTime)

	return nil
}

// FormatWeatherDaily formats daily forecast as text
func (f *TextFormatter) FormatWeatherDaily(data *qweather.WeatherDailyResponse, w io.Writer) error {
	fmt.Fprintf(w, "Weather Forecast\n")
	fmt.Fprintf(w, "================\n\n")

	for i, day := range data.Daily {
		if i > 0 {
			fmt.Fprintf(w, "\n")
		}

		fmt.Fprintf(w, "Date:             %s\n", day.FxDate)
		fmt.Fprintf(w, "Daytime:          %s (%s°C - %s°C)\n", day.TextDay, day.TempMin, day.TempMax)
		fmt.Fprintf(w, "Night:            %s\n", day.TextNight)
		fmt.Fprintf(w, "Wind (Day):       %s %s (%s km/h)\n", day.WindDirDay, day.WindScaleDay, day.WindSpeedDay)
		fmt.Fprintf(w, "Wind (Night):     %s %s (%s km/h)\n", day.WindDirNight, day.WindScaleNight, day.WindSpeedNight)
		fmt.Fprintf(w, "Humidity:         %s%%\n", day.Humidity)
		fmt.Fprintf(w, "UV Index:         %s\n", day.UvIndex)

		if day.Sunrise != "" {
			fmt.Fprintf(w, "Sunrise/Sunset:   %s / %s\n", day.Sunrise, day.Sunset)
		}

		if i < len(data.Daily)-1 {
			fmt.Fprintf(w, "---\n")
		}
	}

	fmt.Fprintf(w, "\nUpdate Time: %s\n", data.UpdateTime)

	return nil
}

// FormatCitySearch formats city search results as text
func (f *TextFormatter) FormatCitySearch(data *qweather.CitySearchResponse, w io.Writer) error {
	fmt.Fprintf(w, "City Search Results\n")
	fmt.Fprintf(w, "===================\n\n")

	if len(data.Location) == 0 {
		fmt.Fprintf(w, "No cities found.\n")
		return nil
	}

	for i, loc := range data.Location {
		fmt.Fprintf(w, "%d. %s\n", i+1, loc.Name)
		fmt.Fprintf(w, "   ID:       %s\n", loc.ID)
		fmt.Fprintf(w, "   Location: %s, %s, %s\n", loc.Adm2, loc.Adm1, loc.Country)
		fmt.Fprintf(w, "   Coords:   %s, %s\n", loc.Lat, loc.Lon)
		fmt.Fprintf(w, "   Timezone: %s (UTC%s)\n", loc.Tz, loc.UtcOffset)

		if i < len(data.Location)-1 {
			fmt.Fprintf(w, "\n")
		}
	}

	return nil
}
