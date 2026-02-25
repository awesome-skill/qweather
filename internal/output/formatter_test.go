package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/pangu-studio/awesome-skills/internal/client/qweather"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFormatter(t *testing.T) {
	testCases := []struct {
		name        string
		format      string
		expectError bool
		expectType  interface{}
	}{
		{"text format", "text", false, &TextFormatter{}},
		{"json format", "json", false, &JSONFormatter{}},
		{"table format", "table", false, &TableFormatter{}},
		{"invalid format", "invalid", true, nil},
		{"empty format", "", true, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter, err := NewFormatter(tc.format)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, formatter)
			} else {
				require.NoError(t, err)
				assert.IsType(t, tc.expectType, formatter)
			}
		})
	}
}

// Test data
func getTestWeatherNowResponse() *qweather.WeatherNowResponse {
	return &qweather.WeatherNowResponse{
		Code:       "200",
		UpdateTime: "2024-02-25T10:30+08:00",
		Now: qweather.WeatherNow{
			ObsTime:   "2024-02-25T10:00+08:00",
			Temp:      "15",
			FeelsLike: "13",
			Icon:      "100",
			Text:      "Sunny",
			Wind360:   "90",
			WindDir:   "East",
			WindScale: "2",
			WindSpeed: "10",
			Humidity:  "65",
			Precip:    "0.0",
			Pressure:  "1013",
			Vis:       "10",
			Cloud:     "20",
			Dew:       "8",
		},
	}
}

func getTestWeatherDailyResponse() *qweather.WeatherDailyResponse {
	return &qweather.WeatherDailyResponse{
		Code:       "200",
		UpdateTime: "2024-02-25T10:30+08:00",
		Daily: []qweather.WeatherDay{
			{
				FxDate:         "2024-02-25",
				Sunrise:        "06:30",
				Sunset:         "18:00",
				TempMax:        "20",
				TempMin:        "10",
				TextDay:        "Sunny",
				TextNight:      "Clear",
				WindDirDay:     "East",
				WindScaleDay:   "2",
				WindSpeedDay:   "10",
				WindDirNight:   "West",
				WindScaleNight: "1",
				WindSpeedNight: "5",
				Humidity:       "60",
				UvIndex:        "5",
			},
			{
				FxDate:         "2024-02-26",
				Sunrise:        "06:29",
				Sunset:         "18:01",
				TempMax:        "22",
				TempMin:        "12",
				TextDay:        "Cloudy",
				TextNight:      "Partly Cloudy",
				WindDirDay:     "South",
				WindScaleDay:   "3",
				WindSpeedDay:   "15",
				WindDirNight:   "South",
				WindScaleNight: "2",
				WindSpeedNight: "10",
				Humidity:       "55",
				UvIndex:        "4",
			},
		},
	}
}

func getTestCitySearchResponse() *qweather.CitySearchResponse {
	return &qweather.CitySearchResponse{
		Code: "200",
		Location: []qweather.Location{
			{
				Name:      "Beijing",
				ID:        "101010100",
				Lat:       "39.90498",
				Lon:       "116.40528",
				Adm2:      "Beijing",
				Adm1:      "Beijing",
				Country:   "China",
				Tz:        "Asia/Shanghai",
				UtcOffset: "+08:00",
			},
			{
				Name:      "Shunyi",
				ID:        "101010200",
				Lat:       "40.13",
				Lon:       "116.65",
				Adm2:      "Shunyi",
				Adm1:      "Beijing",
				Country:   "China",
				Tz:        "Asia/Shanghai",
				UtcOffset: "+08:00",
			},
		},
	}
}

// JSON Formatter Tests
func TestJSONFormatter_FormatWeatherNow(t *testing.T) {
	// Arrange
	formatter := &JSONFormatter{}
	data := getTestWeatherNowResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherNow(data, buf)

	// Assert
	require.NoError(t, err)

	// Verify it's valid JSON
	var result qweather.WeatherNowResponse
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "200", result.Code)
	assert.Equal(t, "15", result.Now.Temp)
	assert.Equal(t, "Sunny", result.Now.Text)
}

func TestJSONFormatter_FormatWeatherDaily(t *testing.T) {
	// Arrange
	formatter := &JSONFormatter{}
	data := getTestWeatherDailyResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherDaily(data, buf)

	// Assert
	require.NoError(t, err)

	var result qweather.WeatherDailyResponse
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "200", result.Code)
	assert.Len(t, result.Daily, 2)
	assert.Equal(t, "2024-02-25", result.Daily[0].FxDate)
}

func TestJSONFormatter_FormatCitySearch(t *testing.T) {
	// Arrange
	formatter := &JSONFormatter{}
	data := getTestCitySearchResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatCitySearch(data, buf)

	// Assert
	require.NoError(t, err)

	var result qweather.CitySearchResponse
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)
	assert.Equal(t, "200", result.Code)
	assert.Len(t, result.Location, 2)
	assert.Equal(t, "Beijing", result.Location[0].Name)
}

// Text Formatter Tests
func TestTextFormatter_FormatWeatherNow(t *testing.T) {
	// Arrange
	formatter := &TextFormatter{}
	data := getTestWeatherNowResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherNow(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	// Verify key information is present
	assert.Contains(t, output, "Current Weather")
	assert.Contains(t, output, "Temperature:      15°C")
	assert.Contains(t, output, "Condition:        Sunny")
	assert.Contains(t, output, "Feels Like:       13°C")
	assert.Contains(t, output, "Humidity:         65%")
	assert.Contains(t, output, "Wind:             East 2 (10 km/h)")
	assert.Contains(t, output, "Pressure:         1013 hPa")
	assert.Contains(t, output, "Visibility:       10 km")
}

func TestTextFormatter_FormatWeatherNow_WithPrecipitation(t *testing.T) {
	// Arrange
	formatter := &TextFormatter{}
	data := getTestWeatherNowResponse()
	data.Now.Precip = "5.2" // Add precipitation
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherNow(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "Precipitation:    5.2 mm")
}

func TestTextFormatter_FormatWeatherDaily(t *testing.T) {
	// Arrange
	formatter := &TextFormatter{}
	data := getTestWeatherDailyResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherDaily(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	// Verify header
	assert.Contains(t, output, "Weather Forecast")

	// Verify first day
	assert.Contains(t, output, "Date:             2024-02-25")
	assert.Contains(t, output, "Daytime:          Sunny (10°C - 20°C)")
	assert.Contains(t, output, "Night:            Clear")
	assert.Contains(t, output, "Humidity:         60%")
	assert.Contains(t, output, "UV Index:         5")
	assert.Contains(t, output, "Sunrise/Sunset:   06:30 / 18:00")

	// Verify second day
	assert.Contains(t, output, "Date:             2024-02-26")
	assert.Contains(t, output, "Daytime:          Cloudy (12°C - 22°C)")
}

func TestTextFormatter_FormatCitySearch(t *testing.T) {
	// Arrange
	formatter := &TextFormatter{}
	data := getTestCitySearchResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatCitySearch(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	assert.Contains(t, output, "City Search Results")
	assert.Contains(t, output, "1. Beijing")
	assert.Contains(t, output, "ID:       101010100")
	assert.Contains(t, output, "Location: Beijing, Beijing, China")
	assert.Contains(t, output, "Coords:   39.90498, 116.40528")
	assert.Contains(t, output, "Timezone: Asia/Shanghai (UTC+08:00)")

	assert.Contains(t, output, "2. Shunyi")
	assert.Contains(t, output, "ID:       101010200")
}

func TestTextFormatter_FormatCitySearch_NoResults(t *testing.T) {
	// Arrange
	formatter := &TextFormatter{}
	data := &qweather.CitySearchResponse{
		Code:     "200",
		Location: []qweather.Location{},
	}
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatCitySearch(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "No cities found")
}

// Table Formatter Tests
func TestTableFormatter_FormatWeatherNow(t *testing.T) {
	// Arrange
	formatter := &TableFormatter{}
	data := getTestWeatherNowResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherNow(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	// Verify table header
	assert.Contains(t, output, "Current Weather")
	assert.Contains(t, output, strings.Repeat("=", 60))

	// Verify key-value pairs are formatted correctly
	assert.Contains(t, output, "Observation Time")
	assert.Contains(t, output, "Temperature")
	assert.Contains(t, output, "15°C")
	assert.Contains(t, output, "Condition")
	assert.Contains(t, output, "Sunny")
}

func TestTableFormatter_FormatWeatherDaily(t *testing.T) {
	// Arrange
	formatter := &TableFormatter{}
	data := getTestWeatherDailyResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatWeatherDaily(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	// Verify table header
	assert.Contains(t, output, "Weather Forecast")
	assert.Contains(t, output, strings.Repeat("=", 100))
	assert.Contains(t, output, "Date")
	assert.Contains(t, output, "Day")
	assert.Contains(t, output, "Night")
	assert.Contains(t, output, "Temp(Min-Max)")

	// Verify data rows
	assert.Contains(t, output, "2024-02-25")
	assert.Contains(t, output, "Sunny")
	assert.Contains(t, output, "Clear")
	assert.Contains(t, output, "10°C - 20°C")

	assert.Contains(t, output, "2024-02-26")
	assert.Contains(t, output, "Cloudy")
	assert.Contains(t, output, "12°C - 22°C")
}

func TestTableFormatter_FormatCitySearch(t *testing.T) {
	// Arrange
	formatter := &TableFormatter{}
	data := getTestCitySearchResponse()
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatCitySearch(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()

	// Verify table structure
	assert.Contains(t, output, "City Search Results")
	assert.Contains(t, output, strings.Repeat("=", 100))
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "ID")
	assert.Contains(t, output, "Region")
	assert.Contains(t, output, "Country")

	// Verify data
	assert.Contains(t, output, "Beijing")
	assert.Contains(t, output, "101010100")
	assert.Contains(t, output, "China")
	assert.Contains(t, output, "Shunyi")
}

func TestTableFormatter_FormatCitySearch_NoResults(t *testing.T) {
	// Arrange
	formatter := &TableFormatter{}
	data := &qweather.CitySearchResponse{
		Code:     "200",
		Location: []qweather.Location{},
	}
	buf := &bytes.Buffer{}

	// Act
	err := formatter.FormatCitySearch(data, buf)

	// Assert
	require.NoError(t, err)
	output := buf.String()
	assert.Contains(t, output, "No cities found")
}

// Test format constants
func TestFormatConstants(t *testing.T) {
	assert.Equal(t, Format("text"), FormatText)
	assert.Equal(t, Format("json"), FormatJSON)
	assert.Equal(t, Format("table"), FormatTable)
}
