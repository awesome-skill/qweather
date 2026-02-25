package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
)

// Format defines the output format type
type Format string

const (
	// FormatText outputs in human-readable text format
	FormatText Format = "text"
	// FormatJSON outputs in JSON format
	FormatJSON Format = "json"
	// FormatTable outputs in table format
	FormatTable Format = "table"
)

// Formatter formats data for output
type Formatter interface {
	FormatWeatherNow(data *qweather.WeatherNowResponse, w io.Writer) error
	FormatWeatherDaily(data *qweather.WeatherDailyResponse, w io.Writer) error
	FormatCitySearch(data *qweather.CitySearchResponse, w io.Writer) error
}

// NewFormatter creates a new formatter based on the format type
func NewFormatter(format string) (Formatter, error) {
	switch Format(format) {
	case FormatText:
		return &TextFormatter{}, nil
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatTable:
		return &TableFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: text, json, table)", format)
	}
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// FormatWeatherNow formats current weather as JSON
func (f *JSONFormatter) FormatWeatherNow(data *qweather.WeatherNowResponse, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FormatWeatherDaily formats daily forecast as JSON
func (f *JSONFormatter) FormatWeatherDaily(data *qweather.WeatherDailyResponse, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FormatCitySearch formats city search results as JSON
func (f *JSONFormatter) FormatCitySearch(data *qweather.CitySearchResponse, w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
