package qweather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	// Arrange
	apiKey := "test-api-key"
	baseURL := "https://api.test.com"

	// Act
	client := NewClient(apiKey, baseURL)

	// Assert
	assert.NotNil(t, client)
	assert.Equal(t, apiKey, client.APIKey)
	assert.Equal(t, baseURL, client.BaseURL)
	assert.NotNil(t, client.HTTPClient)
	assert.Equal(t, 30*time.Second, client.HTTPClient.Timeout)
}

func TestNewClient_DefaultBaseURL(t *testing.T) {
	// Arrange
	apiKey := "test-api-key"

	// Act
	client := NewClient(apiKey, "")

	// Assert
	assert.NotNil(t, client)
	assert.Equal(t, "https://devapi.qweather.com", client.BaseURL)
}

func TestNewClient_AddHTTPS(t *testing.T) {
	// Arrange
	apiKey := "test-api-key"
	baseURL := "api.test.com"

	// Act
	client := NewClient(apiKey, baseURL)

	// Assert
	assert.Equal(t, "https://"+baseURL, client.BaseURL)
}

func TestGetNowWeather_Success(t *testing.T) {
	// Arrange
	expectedResponse := WeatherNowResponse{
		Code: "200",
		Now: WeatherNow{
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/v7/weather/now", r.URL.Path)
		assert.Equal(t, "101010100", r.URL.Query().Get("location"))
		assert.Equal(t, "test-key", r.Header.Get("X-QW-Api-Key"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.GetNowWeather(context.Background(), "101010100")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "200", response.Code)
	assert.Equal(t, "15", response.Now.Temp)
	assert.Equal(t, "Sunny", response.Now.Text)
}

func TestGetNowWeather_APIErrorCode(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WeatherNowResponse{
			Code: "401",
		})
	}))
	defer server.Close()

	client := NewClient("invalid-key", server.URL)

	// Act
	response, err := client.GetNowWeather(context.Background(), "101010100")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API returned error code")
	assert.Nil(t, response)
}

func TestGetNowWeather_HTTPError(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.GetNowWeather(context.Background(), "101010100")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API request failed")
	assert.Nil(t, response)
}

func TestGetDailyForecast_Success(t *testing.T) {
	// Arrange
	expectedResponse := WeatherDailyResponse{
		Code: "200",
		Daily: []WeatherDay{
			{
				FxDate:         "2024-02-25",
				Sunrise:        "06:30",
				Sunset:         "18:00",
				Moonrise:       "19:00",
				Moonset:        "07:00",
				MoonPhase:      "Waxing Gibbous",
				MoonPhaseIcon:  "803",
				TempMax:        "20",
				TempMin:        "10",
				IconDay:        "100",
				TextDay:        "Sunny",
				IconNight:      "150",
				TextNight:      "Clear",
				Wind360Day:     "90",
				WindDirDay:     "East",
				WindScaleDay:   "2",
				WindSpeedDay:   "10",
				Wind360Night:   "270",
				WindDirNight:   "West",
				WindScaleNight: "1",
				WindSpeedNight: "5",
				Humidity:       "60",
				Precip:         "0.0",
				Pressure:       "1013",
				Vis:            "10",
				Cloud:          "20",
				UvIndex:        "5",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/v7/weather/3d", r.URL.Path)
		assert.Equal(t, "101010100", r.URL.Query().Get("location"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.GetDailyForecast(context.Background(), "101010100", 3)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "200", response.Code)
	assert.Len(t, response.Daily, 1)
	assert.Equal(t, "2024-02-25", response.Daily[0].FxDate)
	assert.Equal(t, "20", response.Daily[0].TempMax)
}

func TestGetDailyForecast_ValidDays(t *testing.T) {
	// Test all valid days values
	validDays := []int{3, 7, 10, 15, 30}

	for _, days := range validDays {
		t.Run(fmt.Sprintf("%d days", days), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := fmt.Sprintf("/v7/weather/%dd", days)
				assert.Equal(t, expectedPath, r.URL.Path)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(WeatherDailyResponse{
					Code:  "200",
					Daily: []WeatherDay{},
				})
			}))
			defer server.Close()

			client := NewClient("test-key", server.URL)
			_, err := client.GetDailyForecast(context.Background(), "101010100", days)
			assert.NoError(t, err)
		})
	}
}

func TestGetDailyForecast_InvalidDays(t *testing.T) {
	// Arrange
	client := NewClient("test-key", "https://api.test.com")

	testCases := []struct {
		name string
		days int
	}{
		{"zero days", 0},
		{"negative days", -1},
		{"invalid days", 5},
		{"too many days", 31},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			response, err := client.GetDailyForecast(context.Background(), "101010100", tc.days)

			// Assert
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid days parameter")
			assert.Nil(t, response)
		})
	}
}

func TestGetDailyForecast_APIErrorCode(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WeatherDailyResponse{
			Code: "404",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.GetDailyForecast(context.Background(), "101010100", 3)

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API returned error code")
	assert.Nil(t, response)
}

func TestSearchCity_Success(t *testing.T) {
	// Arrange
	expectedResponse := CitySearchResponse{
		Code: "200",
		Location: []Location{
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
				IsDst:     "0",
				Type:      "city",
				Rank:      "10",
				FxLink:    "http://hfx.link/abc123",
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		assert.Equal(t, "/geo/v2/city/lookup", r.URL.Path)
		assert.Equal(t, "beijing", r.URL.Query().Get("location"))

		// Return mock response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedResponse)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.SearchCity(context.Background(), "beijing")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "200", response.Code)
	assert.Len(t, response.Location, 1)
	assert.Equal(t, "Beijing", response.Location[0].Name)
	assert.Equal(t, "101010100", response.Location[0].ID)
}

func TestSearchCity_NoResults(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CitySearchResponse{
			Code:     "200",
			Location: []Location{},
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.SearchCity(context.Background(), "nonexistentcity12345")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "200", response.Code)
	assert.Empty(t, response.Location)
}

func TestSearchCity_APIErrorCode(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CitySearchResponse{
			Code: "400",
		})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Act
	response, err := client.SearchCity(context.Background(), "")

	// Assert
	require.Error(t, err)
	assert.Contains(t, err.Error(), "API returned error code")
	assert.Nil(t, response)
}

func TestDoRequest_ContextCancellation(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(WeatherNowResponse{Code: "200"})
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL)

	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Act
	response, err := client.GetNowWeather(ctx, "101010100")

	// Assert
	require.Error(t, err)
	assert.Nil(t, response)
}
