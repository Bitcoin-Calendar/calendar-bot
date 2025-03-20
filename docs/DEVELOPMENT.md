# Development Guide

This guide is intended for developers who want to contribute to or extend the Bitcoin Calendar Bot.

## Project Structure

```
nostr-calendar-bot/
├── main.go              # Main application logic, event processing
├── metrics.go           # Metrics collection and reporting
├── events_en.csv        # English events data
├── events_ru.csv        # Russian events data
├── .env                 # Environment variables (not in repo)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── README.md            # Project overview
├── docs/                # Documentation
│   ├── INSTALLATION.md  # Installation guide
│   ├── USAGE.md         # Usage guide
│   ├── CSV_FORMAT.md    # CSV format documentation
│   ├── DEVELOPMENT.md   # This development guide
│   |── CONTRIBUTING.md  # Contribution guidelines
│   |── PROJECT_STRUCTURE.md  # Overview of project structure
|   └── ROADMAP.md       # Development roadmap
└── .gitignore           # Git ignore file
```

## Development Environment Setup

1. **Install Go:**
   Ensure you have Go 1.18 or higher installed. Download from [golang.org](https://golang.org/dl/).

2. **Clone the Repository:**
   ```bash
   git clone https://github.com/bitcoin21ideas/nostr-calendar-bot.git
   cd nostr-calendar-bot
   ```

3. **Install Dependencies:**
   ```bash
   go mod download
   ```

4. **Set Up Environment Variables:**
   Create a `.env` file with test keys for development:
   ```
   NOSTR_PRIVATE_KEY_EN=your_test_private_key_here
   NOSTR_PRIVATE_KEY_RU=your_test_private_key_here
   LOG_LEVEL=debug
   CONSOLE_LOG=true
   ```

5. **Create Test CSV Files:**
   Use small CSV files with test data for development.

## Code Architecture

### Main Components

1. **Event Processing** (`main.go`)
   - CSV file parsing
   - Date matching logic
   - Nostr event creation and signing
   - Relay communication

2. **Metrics Collection** (`metrics.go`)
   - Tracking success/failure metrics
   - Relay performance monitoring
   - Metrics reporting and export

3. **Logging System** (`main.go`)
   - Structured JSON logging
   - Log level control
   - Log rotation and management
   - Language-specific log files

### Key Functions

- `main()`: Entry point, handles configuration and orchestration
- `processCSV()`: Processes CSV files and posts matching events
- `getLogFileName()`: Determines the log file name based on language
- `getLanguageFromCSV()`: Extracts language information from CSV file path
- `getCurrentDirectory()`: Gets the current working directory
- `logEnvironmentVariables()`: Logs non-sensitive environment variables

## Working with Nostr

The bot uses the [go-nostr](https://github.com/nbd-wtf/go-nostr) library to interact with Nostr relays:

1. **Connecting to Relays:**
   ```go
   relay, err := nostr.RelayConnect(ctx, relayURL)
   ```

2. **Creating Events:**
   ```go
   nostrEvent := nostr.Event{
       CreatedAt: nostr.Now(),
       Kind:      nostr.KindTextNote,
       Tags:      tags,
       Content:   message,
   }
   ```

3. **Signing Events:**
   ```go
   nostrEvent.Sign(privateKey)
   ```

4. **Publishing Events:**
   ```go
   relay.Publish(ctx, nostrEvent)
   ```

## Logging System

The bot uses [zerolog](https://github.com/rs/zerolog) for structured logging:

1. **Setting Log Level:**
   ```go
   zerolog.SetGlobalLevel(zerolog.DebugLevel)
   ```

2. **Creating Contextual Logs:**
   ```go
   eventLog := log.With().
       Str("eventDate", eventDate).
       Str("eventTitle", title).
       Int("recordIndex", i+1).
       Logger()
   ```

3. **Writing Logs:**
   ```go
   log.Info().Str("key", "value").Msg("Log message")
   log.Debug().Interface("data", object).Msg("Debug info")
   log.Error().Err(err).Msg("Error occurred")
   ```

## Testing

### Running Tests

```bash
go test ./...
```

### Manual Testing

For manual testing:

1. Create a small test CSV file with today's date:
   ```csv
   "date","title","description"
   "YYYY-MM-DD","Test Event","This is a test event |With multiple lines"
   ```
   (Replace YYYY-MM-DD with today's date)

2. Run the bot with debug logging:
   ```bash
   LOG_LEVEL=debug CONSOLE_LOG=true ./nostr_bot test_events.csv NOSTR_PRIVATE_KEY_EN
   ```

3. Check the console output and log files for expected behavior.

## Building for Production

```bash
go build -o nostr_bot main.go metrics.go
```

## Adding Features

When adding new features:

1. **Create a Branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Follow Go Best Practices:**
   - Use meaningful variable and function names
   - Add comments for complex logic
   - Handle errors properly
   - Write unit tests when possible

3. **Update Documentation:**
   - Update relevant documentation files in the `docs/` directory
   - Add comments for new functions and parameters

4. **Submit a Pull Request:**
   - Provide a clear description of the changes
   - Reference any related issues
   - Ensure tests pass

## Logging Enhancements

The bot offers enhanced logging functionality:

- **Log Level Control**: Setting via the `LOG_LEVEL` environment variable
- **Custom Log Directory**: Setting via the `LOG_DIR` environment variable 
- **Event ID Generation**: Unique IDs for tracking events across logs
- **Request ID Tracking**: For following event processing through logs

When extending the logging system, maintain these capabilities while ensuring logs remain structured and useful.

## Future Development

See the [Roadmap](ROADMAP.md) for planned future development directions.

## Code Style and Standards

- Follow standard Go code formatting (use `gofmt`)
- Follow Go best practices from [Effective Go](https://golang.org/doc/effective_go)
- Write clear comments, especially for complex logic
- Maintain backward compatibility when possible
- Handle errors explicitly

## Questions and Support

For development questions, please:
1. Check existing documentation
2. Review code comments
3. Contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) on Nostr

For more details on contributing, see [CONTRIBUTING.md](CONTRIBUTING.md). 

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
