---
name: qweather
description: QWeather API for real-time weather and forecasts
metadata: {"clawdbot":{"emoji":"🌤️"}}
---

# qweather

Use the QWeather (和风天气) API to get weather data including current conditions, daily forecasts, and city search.


## Commands

### Get Current Weather

Query real-time weather conditions for any location:

```bash
# By city name
qweather now --location "北京"
qweather now --location "Beijing"

# By location ID (get from city search)
qweather now --location "101010100"

# By coordinates (latitude,longitude)
qweather now --location "116.41,39.92"

# JSON output
qweather now --location "北京" --format json

# Table output
qweather now --location "北京" --format table
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
qweather forecast --location "北京" --days 3

# 7-day forecast with table format
qweather forecast --location "上海" --days 7 --format table

# 15-day forecast with JSON output
qweather forecast --location "101010100" --days 15 --format json

# 30-day forecast
qweather forecast --location "116.41,39.92" --days 30
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
qweather search --query "北京"
qweather search --query "beijing"

# Table format
qweather search --query "london" --format table

# JSON format for programmatic use
qweather search --query "paris" --format json
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
qweather now -l "北京" -f json
qweather forecast -l "上海" -d 7 -f table
qweather search -q "tokyo" -f text
```

## Examples

### Complete Workflow

```bash
# 1. Search for a city to get its location ID
qweather search --query "杭州"

# Output shows: ID: 101210101

# 2. Get current qweather using the location ID
qweather now --location "101210101"

# 3. Get 7-day forecast
qweather forecast --location "101210101" --days 7 --format table

# 4. Export to JSON for further processing
qweather now --location "101210101" --format json > weather.json
```

### Using with Scripts

```bash
#!/bin/bash
# Get weather and check if it's raining

WEATHER=$(qweather now --location "北京" --format json)
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
- All times are in ISO 8601 format with timezone offset
- Temperature is in Celsius by default
