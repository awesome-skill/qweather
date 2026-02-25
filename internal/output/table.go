package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
)

// TableFormatter formats output as ASCII tables
type TableFormatter struct{}

// FormatWeatherNow formats current weather as a table
func (f *TableFormatter) FormatWeatherNow(data *qweather.WeatherNowResponse, w io.Writer) error {
	now := data.Now

	fmt.Fprintf(w, "Current Weather\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 60))

	rows := [][]string{
		{"Observation Time", now.ObsTime},
		{"Condition", now.Text},
		{"Temperature", fmt.Sprintf("%s°C", now.Temp)},
		{"Feels Like", fmt.Sprintf("%s°C", now.FeelsLike)},
		{"Humidity", fmt.Sprintf("%s%%", now.Humidity)},
		{"Wind", fmt.Sprintf("%s %s (%s km/h)", now.WindDir, now.WindScale, now.WindSpeed)},
		{"Pressure", fmt.Sprintf("%s hPa", now.Pressure)},
		{"Visibility", fmt.Sprintf("%s km", now.Vis)},
	}

	if now.Precip != "" && now.Precip != "0.0" {
		rows = append(rows, []string{"Precipitation", fmt.Sprintf("%s mm", now.Precip)})
	}

	f.printTable(w, rows)
	fmt.Fprintf(w, "\nUpdate Time: %s\n", data.UpdateTime)

	return nil
}

// FormatWeatherDaily formats daily forecast as a table
func (f *TableFormatter) FormatWeatherDaily(data *qweather.WeatherDailyResponse, w io.Writer) error {
	fmt.Fprintf(w, "Weather Forecast\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 100))

	// Header
	header := []string{"Date", "Day", "Night", "Temp(Min-Max)", "Humidity", "UV"}
	f.printTableHeader(w, header)

	// Rows
	for _, day := range data.Daily {
		row := []string{
			day.FxDate,
			fmt.Sprintf("%s", day.TextDay),
			fmt.Sprintf("%s", day.TextNight),
			fmt.Sprintf("%s°C - %s°C", day.TempMin, day.TempMax),
			fmt.Sprintf("%s%%", day.Humidity),
			day.UvIndex,
		}
		f.printTableRow(w, row)
	}

	fmt.Fprintf(w, "\nUpdate Time: %s\n", data.UpdateTime)

	return nil
}

// FormatCitySearch formats city search results as a table
func (f *TableFormatter) FormatCitySearch(data *qweather.CitySearchResponse, w io.Writer) error {
	fmt.Fprintf(w, "City Search Results\n")
	fmt.Fprintf(w, "%s\n", strings.Repeat("=", 100))

	if len(data.Location) == 0 {
		fmt.Fprintf(w, "No cities found.\n")
		return nil
	}

	// Header
	header := []string{"Name", "ID", "Region", "Country", "Coordinates"}
	f.printTableHeader(w, header)

	// Rows
	for _, loc := range data.Location {
		row := []string{
			loc.Name,
			loc.ID,
			fmt.Sprintf("%s, %s", loc.Adm2, loc.Adm1),
			loc.Country,
			fmt.Sprintf("%s, %s", loc.Lat, loc.Lon),
		}
		f.printTableRow(w, row)
	}

	return nil
}

// printTable prints a simple key-value table
func (f *TableFormatter) printTable(w io.Writer, rows [][]string) {
	maxKeyLen := 0
	for _, row := range rows {
		if len(row[0]) > maxKeyLen {
			maxKeyLen = len(row[0])
		}
	}

	for _, row := range rows {
		fmt.Fprintf(w, "%-*s: %s\n", maxKeyLen, row[0], row[1])
	}
}

// printTableHeader prints table header
func (f *TableFormatter) printTableHeader(w io.Writer, columns []string) {
	// Ensure we have enough columns
	for len(columns) < 6 {
		columns = append(columns, "")
	}
	fmt.Fprintf(w, "%-12s %-15s %-15s %-18s %-10s %-8s\n", columns[0], columns[1], columns[2], columns[3], columns[4], columns[5])
	fmt.Fprintf(w, "%s\n", strings.Repeat("-", 100))
}

// printTableRow prints a table row
func (f *TableFormatter) printTableRow(w io.Writer, columns []string) {
	// Ensure we have enough columns
	for len(columns) < 6 {
		columns = append(columns, "")
	}
	fmt.Fprintf(w, "%-12s %-15s %-15s %-18s %-10s %-8s\n", columns[0], columns[1], columns[2], columns[3], columns[4], columns[5])
}
