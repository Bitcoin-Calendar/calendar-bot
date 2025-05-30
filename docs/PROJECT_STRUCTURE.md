# Project Structure

This document provides an overview of the Bitcoin Calendar Bot project structure to help you navigate and understand the codebase.

## Directory Structure

```
nostr-calendar-bot/
├── main.go              # Application entry point, orchestrates internal modules
├── internal/            # Internal application logic, not intended for external import
│   ├── api/             # Client for interacting with the Bitcoin Calendar events API
│   │   └── client.go
│   ├── config/          # Configuration loading and validation
│   │   └── config.go
│   ├── logging/         # Logging setup and management
│   │   └── setup.go
│   ├── metrics/         # Metrics collection
│   │   └── collector.go
│   ├── models/          # Shared data structures (e.g., APIEvent)
│   │   └── event.go
│   └── nostr/           # Nostr event creation and publishing
│       ├── publisher.go   # Core Nostr event publishing logic
│       ├── kind1.go       # Kind 1 (text) event creation
│       └── kind20.go      # Kind 20 (NIP-68 picture) event creation & image validation
├── Dockerfile           # Defines the Docker image for building and running the bot
├── docker-compose.yml   # Defines Docker Compose services for different bot instances (EN, RU, tests)
├── .env-example         # Example environment variables file
├── .env                 # Your local environment variables (not in repository, gitignored)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── README.md            # Project overview and quick start
├── LICENSE.txt          # MIT license
├── docs/                # Documentation
│   ├── INSTALLATION.md  # Installation guide (Docker-focused)
│   ├── USAGE.md         # Usage guide (Docker-focused)
│   ├── DEVELOPMENT.md   # Development guide
│   ├── CONTRIBUTING.md  # Contribution guidelines
│   |── PROJECT_STRUCTURE.md  # This file
|   └── ROADMAP.md       # Development roadmap
└── .gitignore           # Git ignore file
```

(Note: Event CSV files such as `events_en.csv` or `events_ru.csv` are no longer used by the bot as events are fetched dynamically from an API.)

## Key Files and Directories

### `main.go`

The main application entry point. Its primary responsibilities now include:
-   Parsing command-line arguments (e.g., for Nostr private key environment variable name).
-   Initializing and orchestrating the various internal modules:
    -   Loading configuration (`internal/config`).
    -   Setting up logging (`internal/logging`).
    -   Initializing the API client (`internal/api`).
    -   Initializing the metrics collector (`internal/metrics`).
    -   Initializing the Nostr event publisher and handlers (`internal/nostr`).
-   Managing the main application loop: fetching events via the API client, processing them, and triggering Nostr publications for Kind 1 and Kind 20 events.
-   Handling graceful shutdown and signal processing (if implemented).

### `internal/` Directory

This directory houses the core logic of the application, organized into distinct packages:

-   **`internal/config`**: Manages application configuration. It loads settings from environment variables and `.env` files, validates them, and provides a `Config` struct to the rest of the application.
-   **`internal/api`**: Contains the `Client` for interacting with the external Bitcoin Calendar events API. It handles request construction, sending HTTP requests, parsing responses, and includes retry logic.
-   **`internal/logging`**: Responsible for setting up the global logger (using `zerolog`). It configures log levels, output (console/file), and log rotation (using `lumberjack`).
-   **`internal/metrics`**: Defines the `Collector` for tracking various application metrics, such as the number of events fetched, successfully published (Kind 1 and Kind 20), or failed. It includes methods to increment counters and log summaries.
-   **`internal/models`**: Contains shared data structures used throughout the application, such as `APIEvent` (representing an event from the API) and `APIResponseWrapper` (for handling the API's response structure).
-   **`internal/nostr`**: Encapsulates all logic related to Nostr.
    -   `publisher.go`: Implements `EventPublisher` which handles the actual signing and publishing of `nostr.Event` objects to multiple relays, including connection management and timeouts.
    -   `kind1.go`: Contains `CreateKind1NostrEvent` for constructing Kind 1 (text-based) Nostr events from `APIEvent` data.
    -   `kind20.go`: Contains `CreateKind20NostrEvent` for constructing NIP-68 Kind 20 (picture-based) Nostr events. This includes logic for image URL validation (`ImageValidator`), media type checking, and assembling the specific tags required by NIP-68.

### `Dockerfile`

Defines the multi-stage Docker build process. It compiles the Go application in a builder stage and then copies the binary to a lightweight Alpine image for the final runtime environment. It also sets up a non-root user and necessary directories.

### `docker-compose.yml`

Defines the services for running different instances of the bot (e.g., `nostr-bot-en`, `nostr-bot-ru`, and their test counterparts). Each service is configured with specific environment variables (like `BOT_PROCESSING_LANGUAGE`) and the command to run the bot with the correct Nostr private key argument. It also manages volume mounts for persistent logs.

### Configuration Files

-   `.env-example`: A template for the `.env` file.
-   `.env`: Stores your actual environment variables, including API credentials and Nostr private keys. This file is loaded by `docker-compose` and also by `godotenv` if running `main.go` directly.
-   `go.mod` and `go.sum`: Go module definitions and dependency management.

## Log Files

The bot generates a single primary log file, `nostr_bot.log`, within the log directory configured inside the container (`/app/logs` by default, mapped to `./logs` on the host via `docker-compose.yml`).

-   Log entries include the language being processed for that run, allowing differentiation of activities if multiple language bots log to the same file over time (though separate host directories can be configured in `docker-compose.yml` if strict separation is needed).

## Documentation Structure

The documentation is organized into separate files for clarity:

-   `README.md`: Project overview and quick start guide.
-   `docs/INSTALLATION.md`: Detailed installation instructions (Docker-focused).
-   `docs/USAGE.md`: Guide on how to use the bot (Docker-focused).
-   `docs/DEVELOPMENT.md`: Guidelines for developers.
-   `docs/CONTRIBUTING.md`: Contribution process and guidelines.
-   `docs/PROJECT_STRUCTURE.md`: This file, explaining the codebase structure.
-   `docs/ROADMAP.md`: Future development plans.

## Code Organization (Deprecated - Refer to `internal/` directory structure)

The primary application flow is orchestrated by `main.go`, which utilizes the various packages within the `internal/` directory.

1.  **Initialization (in `main.go`)**:
    *   Load configuration using `config.LoadConfig()`.
    *   Set up logger using `logging.Setup()`.
    *   Initialize API client using `api.NewClient()`.
    *   Initialize metrics collector using `metrics.NewCollector()`.
    *   Initialize Nostr publisher using `nostr.NewEventPublisher()`.
    *   Initialize image validator using `nostr.NewImageValidator()`.
2.  **Main Loop (in `main.go`)**:
    *   Determine current date and configured language.
    *   Fetch events using the API client's `FetchEvents()` method.
    *   For each fetched `APIEvent`:
        *   Attempt to create and publish a Kind 1 event using `nostr.CreateKind1NostrEvent()` and the `EventPublisher`. Update Kind 1 metrics.
        *   If Kind 1 was successful, attempt to create and publish a Kind 20 event (if applicable, based on `APIEvent.Media`) using `nostr.CreateKind20NostrEvent()` and the `EventPublisher`. Update Kind 20 metrics.
        *   Implement a wait period if necessary (e.g., after successful Kind 1 publication).
3.  **Metrics Summary (in `main.go`)**:
    *   Log a summary of collected metrics at the end of the run using `metricsCollector.LogSummary()`.

## Dependencies

The project relies on several external Go packages:

-   `github.com/nbd-wtf/go-nostr`: Nostr protocol implementation.
-   `github.com/rs/zerolog`: Structured logging library.
-   `gopkg.in/natefinch/lumberjack.v2`: Log rotation.
-   `github.com/joho/godotenv`: Loading environment variables from `.env` files.

## Build and Output

The application is built into a Docker image using the `Dockerfile`. The `docker-compose.yml` file then uses this image to run the bot services.

-   The Go binary `nostr_bot` is created inside the Docker image.
-   This binary accepts one command-line argument: the *name* of the environment variable that holds the Nostr private key for the current run.

## Next Steps

-   To contribute to the project, see [CONTRIBUTING.md](CONTRIBUTING.md).
-   For development guidelines, see [DEVELOPMENT.md](DEVELOPMENT.md).
-   To understand how to use the bot, see [USAGE.md](USAGE.md).

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
