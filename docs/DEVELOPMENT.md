# Development Guide

This guide is intended for developers who want to contribute to or extend the Bitcoin Calendar Bot.

## Project Structure

```
nostr-calendar-bot/
├── main.go              # Application entry point, orchestrates internal modules
├── internal/            # Internal application logic
│   ├── api/             # API client (client.go)
│   ├── config/          # Configuration (config.go)
│   ├── logging/         # Logging setup (setup.go)
│   ├── metrics/         # Metrics collection (collector.go)
│   ├── models/          # Shared data models (event.go)
│   └── nostr/           # Nostr logic (publisher.go, kind1.go, kind20.go)
├── Dockerfile
├── docker-compose.yml
├── .env-example
├── go.mod
├── go.sum
├── README.md
├── docs/
│   ├── INSTALLATION.md
│   ├── USAGE.md
│   ├── DEVELOPMENT.md   # This development guide
│   |── CONTRIBUTING.md
│   |── PROJECT_STRUCTURE.md
|   └── ROADMAP.md
└── .gitignore
```

(Note: Event CSV files like `events_en.csv` are no longer part of the core bot structure as events are fetched from an API.)

## Development Environment Setup

1.  **Install Go:**
    Ensure you have Go 1.20 or higher installed (as per `go.mod` or `Dockerfile` if specified, e.g., golang:1.24-alpine).
    Download from [golang.org](https://golang.org/dl/).

2.  **Install Docker & Docker Compose:**
    Required for building and running the bot as intended. See [Docker docs](https://docs.docker.com/get-docker/) and [Docker Compose docs](https://docs.docker.com/compose/install/).

3.  **Clone the Repository:**
    ```bash
    git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
    cd calendar-bot
    ```

4.  **Set Up Environment Variables for Local Development (Optional, if running outside Docker manually):**
    If you plan to run `go run main.go` directly for quick local tests (not the recommended way for full simulation), create a `.env` file:
    ```env
    BOT_API_ENDPOINT="http://your_api_vps_ip_or_localhost:port/api"
    BOT_API_KEY="your_api_key"
    BOT_PROCESSING_LANGUAGE="en" # or "ru", for direct Go runs
    
    NOSTR_PRIVATE_KEY_ENT="your_test_private_key_hex" # For testing
    # Add other keys if your manual test needs them

    LOG_LEVEL="debug"
    CONSOLE_LOG="true"
    ```
    For Docker-based development, the `.env` file is used by `docker-compose` as described in `INSTALLATION.md`.

5.  **Install Dependencies (for local Go tools/LSP):**
    ```bash
    go mod download
    ```

## Code Architecture

The application is now structured with a main orchestrator (`main.go`) and several specialized packages within the `internal/` directory.

### Core Components in `internal/`

1.  **`internal/config` (`config.go`)**: 
    *   Handles loading and validation of all configuration parameters from environment variables and `.env` files (e.g., API details, Nostr keys, logging settings).
    *   Provides a `Config` struct to the rest of the application.

2.  **`internal/api` (`client.go`)**: 
    *   Provides an API `Client` to interact with the Bitcoin Historical Events API.
    *   Manages request creation, response parsing (including the `APIResponseWrapper` for the expected JSON structure), and retry mechanisms.

3.  **`internal/logging` (`setup.go`)**: 
    *   Configures the global `zerolog` logger.
    *   Sets log level, console output, and file logging with rotation (`lumberjack`).

4.  **`internal/metrics` (`collector.go`)**: 
    *   Defines a `Collector` to track application metrics.
    *   Includes counters for events fetched, Kind 1 events posted/failed, Kind 20 events posted/failed/skipped, and image validation failures.
    *   Provides a `LogSummary` method to output collected metrics.

5.  **`internal/models` (`event.go`)**: 
    *   Defines shared data structures like `APIEvent` (matching the API's event structure) and `APIResponseWrapper`.

6.  **`internal/nostr`**: This package contains all Nostr-related functionality.
    *   **`publisher.go`**: Defines `EventPublisher` responsible for:
        *   Connecting to Nostr relays.
        *   Signing `nostr.Event` objects using the configured private key.
        *   Publishing events to relays with timeout and retry logic for individual relays.
        *   Updating metrics for relay-specific publish attempts.
    *   **`kind1.go`**: Contains `CreateKind1NostrEvent`, which takes an `APIEvent` and constructs a `nostr.Event` of Kind 1 (text note). It formats the event content and tags.
    *   **`kind20.go`**: Contains `CreateKind20NostrEvent` for NIP-68 Kind 20 (picture) events.
        *   Includes `ImageValidator` with methods to check image URL validity (extension-based), get media type, and optionally validate accessibility.
        *   Defines `Kind20EventData` to hold necessary data for a Kind 20 event.
        *   The `ToNostrEvent()` method on `Kind20EventData` assembles the `nostr.Event` with all required NIP-68 tags (`title`, `imeta` (URL and hash), `m`, `summary`, `t`, `r`, `d`) and content.

### Orchestration in `main.go`

`main.go` initializes and uses the above components:
-   Loads config using `config.LoadConfig()`.
-   Sets up logging via `logging.Setup()`.
-   Creates instances of `api.Client`, `metrics.Collector`, `nostr.EventPublisher`, and `nostr.ImageValidator`.
-   The main loop fetches events using the API client.
-   For each event:
    -   It calls `nostr.CreateKind1NostrEvent` and then `eventPublisher.PublishEvent`.
    -   It then calls `nostr.CreateKind20NostrEvent` (if applicable) and `eventPublisher.PublishEvent`.
    -   Updates metrics accordingly using the `metrics.Collector`.
-   Logs a metrics summary at the end.

### Key Functions (Illustrative - refer to specific packages)

-   `main()` in `main.go`: Entry point, initialization, main processing loop.
-   `config.LoadConfig()` in `internal/config/config.go`: Loads and validates configuration.
-   `api.Client.FetchEvents()` in `internal/api/client.go`: Fetches events from the API.
-   `nostr.EventPublisher.PublishEvent()` in `internal/nostr/publisher.go`: Publishes a generic `nostr.Event`.
-   `nostr.CreateKind1NostrEvent()` in `internal/nostr/kind1.go`: Creates a Kind 1 event.
-   `nostr.CreateKind20NostrEvent()` in `internal/nostr/kind20.go`: Creates a Kind 20 event.
-   `metrics.Collector.LogSummary()` in `internal/metrics/collector.go`: Logs collected metrics.

## Working with the Bitcoin Historical Events API

The `internal/api/client.go` module handles interaction with the external API (specified by `Config.API.Endpoint`).
-   It constructs requests like: `GET /api/events?month=MM&day=DD&lang=LL`
-   It expects a JSON response structured as `{"events": [...], "pagination": ...}` which is unmarshalled into `models.APIResponseWrapper`.
-   The `models.APIEvent` struct should match the structure of individual events within the `events` array.

## Working with Nostr

Nostr interactions are primarily handled within the `internal/nostr/` package, using the [go-nostr](https://github.com/nbd-wtf/go-nostr) library.

1.  **Connecting to Relays**: Done within `EventPublisher` (`nostr.RelayConnect`).
2.  **Creating Events**: 
    *   Kind 1: `CreateKind1NostrEvent` in `kind1.go` populates `nostr.Event`.
    *   Kind 20: `Kind20EventData.ToNostrEvent` in `kind20.go` populates `nostr.Event` with NIP-68 tags.
3.  **Signing Events**: Done by `EventPublisher` (`nostrEvent.Sign`).
4.  **Publishing Events**: Done by `EventPublisher` (`relay.Publish`).

## Logging System

Managed by `internal/logging/setup.go`, using [zerolog](https://github.com/rs/zerolog) and [lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2).
-   Configuration (log level, console, directory) is sourced from `config.Config`.

## Testing

### Using Docker Compose (Recommended for most testing)

This is the best way to test the bot in an environment that closely mirrors production.

1.  Ensure your `.env` file is configured with `BOT_API_ENDPOINT`, `BOT_API_KEY`, and the necessary `NOSTR_PRIVATE_KEY_ENT` / `NOSTR_PRIVATE_KEY_RUT`.
2.  Build the image if you made code changes: `docker-compose build`
3.  Run a test service:
    ```bash
    docker-compose run --rm nostr-bot-en-test 
    # or
    docker-compose run --rm nostr-bot-ru-test
    ```
4.  Check console output and log files (`./logs/nostr_bot.log` on host) for expected behavior.

### Manual Go Runs (for quick iterative development)

For very quick local tests of specific functions or small changes without rebuilding the Docker image:

1.  Ensure relevant environment variables are set in your shell or `.env` file (and that `godotenv.Load()` is active).
    ```bash
    export BOT_API_ENDPOINT="http://localhost:3000/api" # Or your test API
    export BOT_API_KEY="your_key"
    export BOT_PROCESSING_LANGUAGE="en"
    export NOSTR_PRIVATE_KEY_ENT="your_hex_key"
    export LOG_LEVEL="debug"
    export CONSOLE_LOG="true"
    ```
2.  Run the bot from the project root:
    ```bash
    go run main.go NOSTR_PRIVATE_KEY_ENT
    ```

## Building for Production (via Docker)

The `Dockerfile` handles building the Go application into a lightweight Alpine image.
To build:
```bash
docker-compose build
```
This command reads the `Dockerfile` and `docker-compose.yml` to build the image named (by default) based on the directory name (e.g., `calendar-bot_nostr-bot-en` or a shared image name if configured so).

The resulting image is then used by all services defined in `docker-compose.yml`.

## Adding Features

When adding new features:

1.  **Create a Branch:** `git checkout -b feature/your-feature-name`
2.  **Follow Go Best Practices:** Clean code, comments, error handling.
3.  **Update Documentation:** Reflect changes in relevant `docs/` files.
4.  **Test Thoroughly:** Preferably using the Docker Compose test services.
5.  **Submit a Pull Request:** Clear description, reference issues.

## Code Style and Standards

-   Follow standard Go code formatting (`gofmt`).
-   Follow Go best practices from [Effective Go](https://golang.org/doc/effective_go).
-   Write clear comments, especially for complex logic.
-   Handle errors explicitly and provide context.

## Questions and Support

For development questions, please:
1.  Check existing documentation.
2.  Review code comments.
3.  Contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) on Nostr.

For more details on contributing, see [CONTRIBUTING.md](CONTRIBUTING.md).

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
