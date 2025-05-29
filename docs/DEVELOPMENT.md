# Development Guide

This guide is intended for developers who want to contribute to or extend the Bitcoin Calendar Bot.

## Project Structure

```
nostr-calendar-bot/
├── main.go              # Main application logic, event fetching from API, processing, and Nostr publishing
├── .env-example         # Example environment variables file
├── .env                 # Your local environment variables (not in repo, gitignored)
├── Dockerfile           # Defines the Docker image for the bot
├── docker-compose.yml   # Defines services for different bot instances (EN, RU, test, etc.)
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── README.md            # Project overview
├── docs/                # Documentation
│   ├── INSTALLATION.md  # Installation guide (Docker-focused)
│   ├── USAGE.md         # Usage guide (Docker-focused)
│   ├── DEVELOPMENT.md   # This development guide
│   |── CONTRIBUTING.md  # Contribution guidelines
│   |── PROJECT_STRUCTURE.md  # Overview of project structure
|   └── ROADMAP.md       # Development roadmap
└── .gitignore           # Git ignore file
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

### Main Components

1.  **Event Fetching & Processing** (`main.go`)
    -   API client logic to fetch events for the current day and specified language.
    -   JSON parsing of API responses.
    -   Nostr event creation and signing.
    -   Relay communication.

2.  **Logging System** (`main.go`)
    -   Structured JSON logging using Zerolog.
    -   Log level control via environment variables.
    -   Log rotation and management (handled by Lumberjack).

### Key Functions in `main.go`

-   `main()`: Entry point, handles configuration (reading env vars like `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE`, and the private key env var name from args), sets up logging, and orchestrates event fetching and publishing.
-   `fetchEventsFromAPI()`: Fetches events from the configured API for a given date and language.
-   `publishEvent()`: Prepares and publishes a single event to Nostr relays.
-   `getCurrentDirectory()`: Gets the current working directory (mostly for debug logging).
-   `logEnvironmentVariables()`: Logs non-sensitive environment variables (for debug).

## Working with the Bitcoin Historical Events API

The bot interacts with an external API (specified by `BOT_API_ENDPOINT`) to get events. Refer to `APIDocumentation.md` for details on the API's expected requests and responses.
-   The bot constructs requests like: `GET /api/events?month=MM&day=DD&lang=LL`
-   It expects a JSON response containing an array of events and pagination info.
-   The `APIEvent` struct in `main.go` should match the structure of events returned by the API.

## Working with Nostr

The bot uses the [go-nostr](https://github.com/nbd-wtf/go-nostr) library:

1.  **Connecting to Relays:** `nostr.RelayConnect(ctx, relayURL)`
2.  **Creating Events:** Populating `nostr.Event{...}` with content, kind, tags.
3.  **Signing Events:** `nostrEvent.Sign(privateKey)`
4.  **Publishing Events:** `relay.Publish(ctx, nostrEvent)`

## Logging System

Uses [zerolog](https://github.com/rs/zerolog) and [lumberjack.v2](https://gopkg.in/natefinch/lumberjack.v2).
-   Configure log level with `LOG_LEVEL` (`debug`, `info`, `warn`, `error`).
-   Enable console logging with `CONSOLE_LOG=true`.
-   Log directory inside the container is set by `LOG_DIR` (defaults to `/app/logs`), which is volume-mounted to `./logs` on the host.

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
