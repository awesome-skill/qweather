# Awesome Skills

A collection of skill implementations for AI agent platforms like OpenClaw. Each skill is implemented as a CLI tool that provides specific functionality through a standardized interface.

## Overview

This project provides skills in the form of command-line tools built with Go. Each skill:

- Is implemented as a standalone CLI using the Cobra framework
- Outputs data to stdout in multiple formats (text, JSON, table)
- Follows consistent patterns for configuration and error handling
- Includes detailed skill definition documentation

## Available Skills

### Weather (QWeather)

Get real-time weather data and forecasts powered by QWeather API.

**Features:**
- Current weather conditions
- Daily weather forecasts (3-30 days)
- City search with fuzzy matching
- Multiple output formats (text, JSON, table)

**Skill Definition:** [skills/weather/qweather.md](skills/weather/qweather.md)

## Quick Start

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile)
- QWeather API key (for weather skill)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/pangu-studio/awesome-skills.git
cd awesome-skills
```

2. Build all CLIs:
```bash
make build
```

Or build a specific CLI:
```bash
make build-weather
```

3. Install to system (optional):
```bash
make install
```

Binaries will be installed to `$GOPATH/bin`.

### Configuration

#### Weather Skill Setup

1. Get your QWeather API key from https://console.qweather.com

2. Configure using one of these methods:

**Option 1: Configuration file**
```bash
mkdir -p ~/.config/awesome-skills/qweather
echo "your_api_key_here" > ~/.config/awesome-skills/qweather/api_key
```

**Option 2: Environment variable**
```bash
export QWEATHER_API_KEY="your_api_key_here"
```

3. Test the configuration:
```bash
./bin/weather now --location "北京"
```

## Usage Examples

### Weather CLI

**Get current weather:**
```bash
weather now --location "北京"
weather now --location "101010100" --format json
```

**Get weather forecast:**
```bash
weather forecast --location "上海" --days 7
weather forecast --location "116.41,39.92" --days 15 --format table
```

**Search for cities:**
```bash
weather search --query "beijing"
weather search --query "杭州" --format table
```

## Project Structure

```
awesome-skills/
├── cmd/                    # CLI entry points
│   ├── weather/           # Weather CLI
│   │   ├── main.go
│   │   └── cmd/          # Cobra commands
│   └── ask-cli/          # Other CLIs
├── internal/              # Private packages
│   ├── client/           # API clients
│   │   └── qweather/    # QWeather client
│   ├── config/          # Configuration management
│   └── output/          # Output formatters
├── skills/               # Skill definitions
│   └── weather/         # Weather skill docs
├── Makefile             # Build automation
├── go.mod               # Go module definition
└── README.md            # This file
```

## Development

### Building

```bash
# Build all CLIs
make build

# Build specific CLI
make build-weather

# Install to system
make install
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-cover
```

### Code Quality

```bash
# Format code
make fmt

# Run go vet
make vet

# Run all linters
make lint
```

### Cleaning

```bash
# Clean build artifacts
make clean
```

## Adding New Skills

To add a new skill:

1. **Create CLI implementation:**
   - Create directory under `cmd/` (e.g., `cmd/translate/`)
   - Implement using Cobra framework
   - Follow existing patterns for config and output

2. **Create API client (if needed):**
   - Add client package under `internal/client/`
   - Implement API methods with proper error handling
   - Use context for timeouts

3. **Add to build system:**
   - Update `BINARIES` variable in Makefile
   - Add build target if needed

4. **Write skill definition:**
   - Create markdown file under `skills/`
   - Follow the format in `skills/weather/qweather.md`
   - Include setup, usage examples, and API information

5. **Update documentation:**
   - Add entry to this README
   - Document any configuration requirements

## Code Guidelines

This project follows the Go best practices documented in [AGENTS.md](AGENTS.md). Key points:

- Use `gofmt` for formatting
- Follow standard Go naming conventions
- Write godoc comments for public APIs
- Handle errors explicitly with context
- Write tests for core functionality
- Keep functions small and focused

## Configuration

### Environment Variables

- `QWEATHER_API_KEY`: QWeather API key
- `QWEATHER_API_HOST`: QWeather API host (default: devapi.qweather.com)

### Configuration Files

Configuration files are stored in:
- Linux/macOS: `~/.config/awesome-skills/`
- Windows: `%APPDATA%\awesome-skills\`

## Output Formats

All CLIs support multiple output formats:

- **text**: Human-readable text (default)
- **json**: JSON format for programmatic use
- **table**: ASCII table format for structured data

Specify format with `--format` or `-f` flag:
```bash
weather now -l "Beijing" -f json
weather forecast -l "Shanghai" -d 7 -f table
```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Follow the code guidelines in AGENTS.md
4. Write tests for new functionality
5. Submit a pull request

## Support

For issues and questions:
- Open an issue on GitHub
- Check skill-specific documentation in `skills/`
- Review API provider documentation

## Acknowledgments

- Weather data powered by [QWeather](https://www.qweather.com)
- CLI framework: [Cobra](https://github.com/spf13/cobra)
