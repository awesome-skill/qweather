---
name: qweather
description: QWeather API for real-time weather and forecasts
homepage: https://dev.qweather.com
metadata: {"clawdbot":{"emoji":"🌤️"}}
---

# qweather

Use the QWeather (和风天气) API to get weather data including current conditions, daily forecasts, and city search.

## Setup

1. Register for a QWeather account at https://id.qweather.com/#/register
2. Create a project and obtain your API key from https://console.qweather.com
3. Configure the API key using one of the following methods:

**Method 1: Configuration file** (Recommended)
```bash
mkdir -p ~/.config/awesome-skills/qweather
echo "your_api_key_here" > ~/.config/awesome-skills/qweather/api_key
```

**Method 2: Environment variable**
```bash
export QWEATHER_API_KEY="your_api_key_here"
export QWEATHER_API_HOST="devapi.qweather.com"  # Optional, defaults to free tier
```

**Note:** Free tier uses `devapi.qweather.com`, paid tier uses `api.qweather.com`

## Commands

### Get Current Weather

Query real-time weather conditions for any location:

```bash
# By city name
weather now --location "北京"
weather now --location "Beijing"

# By location ID (get from city search)
weather now --location "101010100"

# By coordinates (latitude,longitude)
weather now --location "116.41,39.92"

# JSON output
weather now --location "北京" --format json

# Table output
weather now --location "北京" --format table
```

**Output includes:**
- Temperature and feels-like temperature
- Weather condition (sunny, cloudy, rainy, etc.)
- Wind direction, speed, and scale
- Humidity, pressure, visibility
- Precipitation (if any)
- Observation and update time

### Get Weather Forecast

Query daily weather forecast (3, 7, 10, 15, or 30 days):

```bash
# 3-day forecast (default)
weather forecast --location "北京" --days 3

# 7-day forecast with table format
weather forecast --location "上海" --days 7 --format table

# 15-day forecast with JSON output
weather forecast --location "101010100" --days 15 --format json

# 30-day forecast
weather forecast --location "116.41,39.92" --days 30
```

**Output includes for each day:**
- Date
- Daytime and nighttime weather conditions
- Temperature range (min/max)
- Wind information (day and night)
- Humidity and UV index
- Sunrise and sunset times (when available)

### Search Cities

Find cities by name with fuzzy matching:

```bash
# Search for cities
weather search --query "北京"
weather search --query "beijing"

# Table format
weather search --query "london" --format table

# JSON format for programmatic use
weather search --query "paris" --format json
```

**Output includes:**
- City name
- Location ID (for use in other commands)
- Administrative regions (city, province/state, country)
- Coordinates (latitude, longitude)
- Timezone information

## Output Formats

The CLI supports three output formats:

- **text** (default): Human-readable text format
- **json**: JSON format for programmatic processing
- **table**: ASCII table format for easy viewing of multiple records

Specify format with the `--format` or `-f` flag:

```bash
weather now -l "北京" -f json
weather forecast -l "上海" -d 7 -f table
weather search -q "tokyo" -f text
```

## API Information

### Request Limits

- **Free tier**: 1,000 requests per day
- **Standard tier**: Higher limits available

### Data Update Frequency

- Current weather: Updated every 10-20 minutes
- Daily forecast: Updated multiple times daily
- City database: Updated regularly

### Error Codes

If you encounter errors, check:
- API key is correctly configured
- Location parameter is valid
- Request limits not exceeded
- Internet connection is available

### Documentation

Full API documentation: https://dev.qweather.com/docs/api/

### Support

- QWeather Console: https://console.qweather.com
- API Documentation: https://dev.qweather.com/docs/
- Contact: https://www.qweather.com/contact/

## Examples

### Complete Workflow

```bash
# 1. Search for a city to get its location ID
weather search --query "杭州"

# Output shows: ID: 101210101

# 2. Get current weather using the location ID
weather now --location "101210101"

# 3. Get 7-day forecast
weather forecast --location "101210101" --days 7 --format table

# 4. Export to JSON for further processing
weather now --location "101210101" --format json > weather.json
```

### Using with Scripts

```bash
#!/bin/bash
# Get weather and check if it's raining

WEATHER=$(weather now --location "北京" --format json)
CONDITION=$(echo $WEATHER | jq -r '.now.text')

if [[ $CONDITION == *"雨"* ]]; then
    echo "It's raining! Don't forget your umbrella."
else
    echo "Weather looks good: $CONDITION"
fi
```

## Notes

- Location ID is the most reliable way to query weather (use city search to find it)
- Coordinates should be in decimal format: longitude,latitude
- Free tier API uses `devapi.qweather.com`, paid tier uses `api.qweather.com`
- All times are in ISO 8601 format with timezone offset
- Temperature is in Celsius by default
