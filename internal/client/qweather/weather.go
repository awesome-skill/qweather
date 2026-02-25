package qweather

import (
	"context"
	"fmt"
	"net/url"
)

// WeatherNowResponse represents the current weather API response
type WeatherNowResponse struct {
	Code       string     `json:"code"`
	UpdateTime string     `json:"updateTime"`
	FxLink     string     `json:"fxLink"`
	Now        WeatherNow `json:"now"`
	Refer      Refer      `json:"refer"`
}

// WeatherNow contains current weather data
type WeatherNow struct {
	ObsTime   string `json:"obsTime"`
	Temp      string `json:"temp"`
	FeelsLike string `json:"feelsLike"`
	Icon      string `json:"icon"`
	Text      string `json:"text"`
	Wind360   string `json:"wind360"`
	WindDir   string `json:"windDir"`
	WindScale string `json:"windScale"`
	WindSpeed string `json:"windSpeed"`
	Humidity  string `json:"humidity"`
	Precip    string `json:"precip"`
	Pressure  string `json:"pressure"`
	Vis       string `json:"vis"`
	Cloud     string `json:"cloud"`
	Dew       string `json:"dew"`
}

// WeatherDailyResponse represents the daily forecast API response
type WeatherDailyResponse struct {
	Code       string       `json:"code"`
	UpdateTime string       `json:"updateTime"`
	FxLink     string       `json:"fxLink"`
	Daily      []WeatherDay `json:"daily"`
	Refer      Refer        `json:"refer"`
}

// WeatherDay contains daily weather forecast data
type WeatherDay struct {
	FxDate         string `json:"fxDate"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	Moonrise       string `json:"moonrise"`
	Moonset        string `json:"moonset"`
	MoonPhase      string `json:"moonPhase"`
	MoonPhaseIcon  string `json:"moonPhaseIcon"`
	TempMax        string `json:"tempMax"`
	TempMin        string `json:"tempMin"`
	IconDay        string `json:"iconDay"`
	TextDay        string `json:"textDay"`
	IconNight      string `json:"iconNight"`
	TextNight      string `json:"textNight"`
	Wind360Day     string `json:"wind360Day"`
	WindDirDay     string `json:"windDirDay"`
	WindScaleDay   string `json:"windScaleDay"`
	WindSpeedDay   string `json:"windSpeedDay"`
	Wind360Night   string `json:"wind360Night"`
	WindDirNight   string `json:"windDirNight"`
	WindScaleNight string `json:"windScaleNight"`
	WindSpeedNight string `json:"windSpeedNight"`
	Humidity       string `json:"humidity"`
	Precip         string `json:"precip"`
	Pressure       string `json:"pressure"`
	Vis            string `json:"vis"`
	Cloud          string `json:"cloud"`
	UvIndex        string `json:"uvIndex"`
}

// Refer contains data source information
type Refer struct {
	Sources []string `json:"sources"`
	License []string `json:"license"`
}

// GetNowWeather retrieves current weather data for a location
func (c *Client) GetNowWeather(ctx context.Context, location string) (*WeatherNowResponse, error) {
	params := url.Values{}
	params.Set("location", location)

	var result WeatherNowResponse
	if err := c.doRequest(ctx, "/v7/weather/now", params, &result); err != nil {
		return nil, fmt.Errorf("get current weather: %w", err)
	}

	// Check API response code
	if result.Code != "200" {
		return nil, fmt.Errorf("API returned error code: %s", result.Code)
	}

	return &result, nil
}

// GetDailyForecast retrieves daily weather forecast for a location
func (c *Client) GetDailyForecast(ctx context.Context, location string, days int) (*WeatherDailyResponse, error) {
	// Validate days parameter
	validDays := map[int]bool{3: true, 7: true, 10: true, 15: true, 30: true}
	if !validDays[days] {
		return nil, fmt.Errorf("invalid days parameter: %d (valid values: 3, 7, 10, 15, 30)", days)
	}

	endpoint := fmt.Sprintf("/v7/weather/%dd", days)
	params := url.Values{}
	params.Set("location", location)

	var result WeatherDailyResponse
	if err := c.doRequest(ctx, endpoint, params, &result); err != nil {
		return nil, fmt.Errorf("get daily forecast: %w", err)
	}

	// Check API response code
	if result.Code != "200" {
		return nil, fmt.Errorf("API returned error code: %s", result.Code)
	}

	return &result, nil
}
