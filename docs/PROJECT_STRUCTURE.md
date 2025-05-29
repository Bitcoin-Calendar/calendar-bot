# Project Structure

This document provides an overview of the Bitcoin Calendar Bot project structure to help you navigate and understand the codebase.

## Directory Structure

```
nostr-calendar-bot/
├── main.go              # Main application logic: API event fetching, processing, Nostr publishing, logging
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

## Key Files

### `main.go`

Contains the core application logic:

-   Application entry point and command-line argument processing (for Nostr private key environment variable name).
-   Loading and validation of environment variables (`BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE`, etc.).
-   Fetching events from the configured API based on current date and language.
-   JSON parsing of API responses.
-   Nostr event creation and publishing to relays.
-   Logging system configuration and management using Zerolog and Lumberjack.

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

## Code Organization

### Main Application Flow (within `main.go`)

1.  **Initialization**:
    *   Parse command-line argument (Nostr private key env var name).
    *   Load environment variables (API config, language, logging settings).
    *   Configure logging (Zerolog with Lumberjack for rotation).
2.  **Event Fetching & Processing**:
    *   Call `fetchEventsFromAPI()` with current date and language to get events.
    *   Iterate through fetched events that match today's date.
    *   For each valid event, call `publishEvent()`.
3.  **Event Publishing** (`publishEvent()`):
    *   Format event content (title, description, media, references).
    *   Create Nostr event structure with tags.
    *   Sign the event.
    *   Attempt to publish to all configured relays.

### Important Functions

#### In `main.go`:

-   `main()`: Entry point and orchestration.
-   `fetchEventsFromAPI()`: Handles API communication to retrieve events.
-   `publishEvent()`: Manages the creation, signing, and relaying of a single Nostr event.

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
