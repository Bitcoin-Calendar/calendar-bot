# Project Structure

This document provides an overview of the Bitcoin Calendar Bot project structure to help you navigate and understand the codebase.

## Directory Structure

```
nostr-calendar-bot/
├── main.go              # Main application logic, event processing, and logging
├── metrics.go           # Metrics collection and reporting functionality
├── events_en.csv        # English events database
├── events_ru.csv        # Russian events database
├── .env                 # Environment variables (not in repository)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── README.md            # Project overview and quick start
├── LICENSE.txt          # MIT license
├── docs/                # Documentation
│   ├── INSTALLATION.md  # Installation guide
│   ├── USAGE.md         # Usage guide
│   ├── CSV_FORMAT.md    # CSV format documentation
│   ├── DEVELOPMENT.md   # Development guide
│   ├── CONTRIBUTING.md  # Contribution guidelines
│   |── PROJECT_STRUCTURE.md  # This file
|   └── ROADMAP.md       # Development roadmap
└── .gitignore           # Git ignore file
```

## Key Files

### `main.go`

Contains the core application logic:

- Application entry point and command-line processing
- CSV file parsing and event processing
- Nostr event creation and publishing
- Logging system configuration and management
- Environment variable handling
- Language detection and log file selection

### `metrics.go`

Contains metrics collection functionality:

- Tracking of events posted, skipped, and failed
- Relay performance metrics
- Success/failure statistics
- JSON export of metrics

### Event CSV Files

- `events_en.csv`: Database of Bitcoin historical events in English
- `events_ru.csv`: Database of Bitcoin historical events in Russian

Each CSV file follows the format documented in [CSV_FORMAT.md](CSV_FORMAT.md).

### Configuration Files

- `.env`: Contains environment variables including private keys (not in repository)
- `go.mod` and `go.sum`: Go module definitions and dependency management

## Log Files

The bot generates the following log files:

- `nostr_bot_en.log`: Logs for English events processing
- `nostr_bot_ru.log`: Logs for Russian events processing
- `nostr_bot_unknown.log`: Logs for events with unrecognized language
- `metrics_<language>_<date>.json`: Metrics export file for each run

## Documentation Structure

The documentation is organized into separate files for clarity:

- `README.md`: Project overview and quick start guide
- `docs/INSTALLATION.md`: Detailed installation instructions
- `docs/USAGE.md`: Guide on how to use the bot
- `docs/CSV_FORMAT.md`: Specification for CSV file format
- `docs/DEVELOPMENT.md`: Guidelines for developers
- `docs/CONTRIBUTING.md`: Contribution process and guidelines
- `docs/PROJECT_STRUCTURE.md`: This file, explaining the codebase structure
- `docs/ROADMAP.md`: Future development plans


## Code Organization

### Main Application Flow

1. **Initialization** (in `main()` function):
   - Parse command-line arguments
   - Load environment variables
   - Configure logging
   - Initialize metrics collector

2. **Event Processing** (in `processCSV()` function):
   - Open and read CSV file
   - Compare event dates with today's date
   - Format matching events
   - Post events to Nostr relays

3. **Metrics Collection** (throughout the code):
   - Track events processed
   - Monitor relay performance
   - Export metrics at the end of execution

### Important Functions

#### In `main.go`:

- `main()`: Entry point and orchestration
- `processCSV()`: Main event processing logic
- `getLogFileName()`: Determines log file based on language
- `getLanguageFromCSV()`: Extracts language from CSV path
- `getCurrentDirectory()`: Gets working directory
- `logEnvironmentVariables()`: Logs non-sensitive variables

#### In `metrics.go`:

- `NewMetricsCollector()`: Creates a new metrics collector
- `RecordRelaySuccess()`: Records successful relay operations
- `RecordRelayFailure()`: Records failed relay operations
- `LogSummary()`: Outputs metrics summary to logs
- `ExportMetrics()`: Exports metrics to JSON file

## Dependencies

The project relies on several external Go packages:

- `github.com/nbd-wtf/go-nostr`: Nostr protocol implementation | *h/t [@fiatjaf](https://njump.me/fiatjaf.com)*
- `github.com/rs/zerolog`: Structured logging library
- `gopkg.in/natefinch/lumberjack.v2`: Log rotation
- `github.com/joho/godotenv`: Environment variable loading

## Build and Output

When the application is built, it produces a single executable file:

- `nostr_bot`: The main executable binary

This binary accepts command-line arguments for the CSV file and environment variable name containing the private key.

## Next Steps

- To contribute to the project, see [CONTRIBUTING.md](CONTRIBUTING.md)
- For development guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md)
- To understand how to use the bot, see [USAGE.md](USAGE.md) 

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
