[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)

# Bitcoin Calendar Bot

### Archiving and relaying every Bitcoin milestone 

![bitcal-structure](https://haven.bitcoin-calendar.org/4d4f81403e39c0c4a454a35cb6913a9420539c9665cb5240fdaff4e27b1e5176.webp)

## Overview

The Bitcoin Calendar Bot is a Go-based application that fetches historical Bitcoin event data from the Bitcoin Historical Events API and publishes these events to Nostr relays. This bot automates the posting of calendar events about Bitcoin history.

Key functionalities include:
*   Fetching events from a configurable API endpoint.
*   Publishing events as Nostr Kind 1 (text-based) notes.
*   **NEW**: Publishing events as Nostr NIP-68 Kind 20 (picture-based) notes for events with associated images.
*   Support for multiple languages (e.g., English, Russian), configurable at runtime.
*   Metrics collection for monitoring event posting success and failures.

The project has recently undergone a significant refactoring to improve modularity and maintainability. Core functionalities such as configuration management, API interaction, logging, metrics collection, and Nostr event publishing have been moved into dedicated packages within an `internal` directory.

## Quick Start with Docker (Recommended)

This is the recommended way to run the Bitcoin Calendar Bot. It uses Docker and Docker Compose to manage the bot instances for different languages.

1.  **Prerequisites**:
    *   Ensure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.

2.  **Clone the repository**:
    ```bash
    git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
    cd calendar-bot
    ```

3.  **Set up environment variables**:
    *   Copy the example environment file:
        ```bash
        cp .env-example .env
        ```
    *   Edit the `.env` file and add your API configuration and Nostr private keys for production and testing:
        ```env
        # --- API Configuration (Required) ---
        BOT_API_ENDPOINT="http://your_api_vps_ip:port/api" # Replace with your API's base URL
        BOT_API_KEY="your_secret_api_key"             # Replace with your API key

        # --- Nostr Private Keys (Required) ---
        # Production Keys
        NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex"
        NOSTR_PRIVATE_KEY_RU="your_russian_specific_private_key_hex"

        # Test Keys
        NOSTR_PRIVATE_KEY_ENT="your_english_TEST_private_key_hex"
        NOSTR_PRIVATE_KEY_RUT="your_russian_TEST_private_key_hex"
        
        # --- Optional: Logging Configuration ---
        # LOG_LEVEL="debug" # Or "info", "warn", "error"
        # CONSOLE_LOG="true"  # Or "false"
        # LOG_DIR="/app/logs" # Path inside the container, usually mapped to a host volume
        ```
    *   The `BOT_PROCESSING_LANGUAGE` is set per-service in `docker-compose.yml` and should not be in the `.env` file.

4.  **Build the Docker image**:
    You need to build the image that all services will use.
    ```bash
    docker-compose build
    # Or, to build a specific service's image if needed:
    # docker-compose build nostr-bot-en 
    ```

5.  **Running Test Instances Manually**:
    To test your setup without posting to your main Nostr accounts, you can run the test services defined in `docker-compose.yml` (e.g., `nostr-bot-en-test`, `nostr-bot-ru-test`). These use specific test keys from your `.env` file and are configured for the respective language. For detailed commands, please see the [Running Test Instances section in our Usage Guide](docs/USAGE.md#running-test-instances).
    Example:
    ```bash
    docker-compose run --rm nostr-bot-en-test
    docker-compose run --rm nostr-bot-ru-test
    ```

6.  **Set up Cron Jobs for Production Bots**:
    The production bots (`nostr-bot-en`, `nostr-bot-ru`) are designed for scheduled runs using your system's cron. Each service is pre-configured in `docker-compose.yml` to run for a specific language. For detailed instructions on setting up cron jobs with `docker-compose`, please refer to the [Automated Operation section in our Usage Guide](docs/USAGE.md#automated-operation).

7.  **Viewing Logs**:
    The bot generates log files for its operations. For test runs, output is also sent to the console by default. Detailed information on log file locations and rotation can be found in the [Log Files section of our Usage Guide](docs/USAGE.md#log-files).

8.  **Building after code changes**:
    If you update the bot's Go code, you need to rebuild the Docker image before your cron jobs or manual test runs pick up the changes:
    ```bash
    cd /path/to/your/calendar-bot 
    docker-compose build
    ```

The `docker-compose.yml` file defines production services (`nostr-bot-en`, `nostr-bot-ru`) for scheduled runs and test services (`nostr-bot-en-test`, `nostr-bot-ru-test`) for manual testing. All services use the same built image but are configured with different environment variables (like `BOT_PROCESSING_LANGUAGE`) and command parameters (for Nostr private key environment variable names) to control their behavior.

The main application logic is now organized within the `internal` directory:
*   `internal/config`: Handles loading and validation of application configuration from environment variables and a `.env` file.
*   `internal/api`: Contains the client for interacting with the Bitcoin Calendar events API, including event fetching and retry logic.
*   `internal/logging`: Manages logger setup and configuration (using zerolog), supporting console and file-based logging.
*   `internal/metrics`: Provides a collector for tracking application-specific metrics, such as the number of events processed and published (for both Kind 1 and Kind 20).
*   `internal/models`: Defines shared data structures used across the application, like `APIEvent` and `APIResponseWrapper`.
*   `internal/nostr`: Handles all Nostr-related operations.
    *   `publisher.go`: Core event publishing logic to relays.
    *   `kind1.go`: Logic for creating Kind 1 (text) Nostr events.
    *   `kind20.go`: Logic for creating NIP-68 Kind 20 (picture) Nostr events, including image validation.

The `main.go` file now serves as the entry point, orchestrating these components.

<details>
<summary>Legacy Manual Setup (Without Docker - Deprecated)</summary>

This method is no longer recommended as the primary way to run the bot due to the API-driven nature and Docker-first approach for managing configurations like `BOT_PROCESSING_LANGUAGE`.

1. **Clone the repository**
   ```bash
   git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
   cd calendar-bot
   ```

2. **Build the application**
   ```bash
   go build -o nostr_bot main.go
   ```

3. **Set up environment variables**
   Create a `.env` file in the project directory or ensure these are set in your environment:
   ```env
   # --- API Configuration (Required) ---
   BOT_API_ENDPOINT="http://your_api_vps_ip:port/api" 
   BOT_API_KEY="your_secret_api_key"

   # --- Language for this instance (Required for manual run) ---
   BOT_PROCESSING_LANGUAGE="en" # or "ru"

   # --- Nostr Private Key (Required) ---
   # The bot takes the NAME of the env var holding the key as a command argument.
   # So, if you pass NOSTR_KEY_MY_ACCOUNT as an argument, set it here:
   NOSTR_KEY_MY_ACCOUNT="your_private_key_hex"
   ```

4. **Run the bot**
   The bot fetches events based on the current date, its API configuration, and the `BOT_PROCESSING_LANGUAGE` env var.
   ```bash
   # Example for an English bot instance:
   # Ensure BOT_API_ENDPOINT, BOT_API_KEY, BOT_PROCESSING_LANGUAGE="en", 
   # and NOSTR_PRIVATE_KEY_FOR_EN (or your chosen name) are set in the environment.
   LOG_DIR=./logs ./nostr_bot NOSTR_PRIVATE_KEY_FOR_EN 
   # Replace NOSTR_PRIVATE_KEY_FOR_EN with the env var name for the specific key you want to use.
   ```

### Automated Daily Posting with Cron (Manual Setup - Deprecated)

Cron jobs run with a minimal environment. Ensure `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE`, and the specific `NOSTR_PRIVATE_KEY_...` variable are available.

Example cron entries (adapt as needed):
```cron
# Ensure the path to your project directory is correct.

# Run English bot instance at 04:00 AM daily
00 04 * * * cd /path/to/your/calendar-bot && BOT_PROCESSING_LANGUAGE="en" LOG_DIR=./logs ./nostr_bot NOSTR_PRIVATE_KEY_EN

# Run Russian bot instance at 05:00 AM daily
00 05 * * * cd /path/to/your/calendar-bot && BOT_PROCESSING_LANGUAGE="ru" LOG_DIR=./logs ./nostr_bot NOSTR_PRIVATE_KEY_RU
```
Consider using the Docker setup for cron jobs for easier management of environment variables per language.

</details>

## Documentation

- [Installation Guide](docs/INSTALLATION.md) - Detailed installation and setup instructions
- [Usage Guide](docs/USAGE.md) - How to use the bot effectively and available configuration options
- [Development Guide](docs/DEVELOPMENT.md) - Information for developers
- [Roadmap](docs/ROADMAP.md) - Future development plans

## License

This project is licensed under the MIT License. See the [LICENSE.txt](LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

## Support

Support Bitcoin Calendar via [Coinos](https://coinos.io/)
- Bitcoin Calendar EN ⚡️ `bitcal@coinos.io` 🔗 [Coinos page](https://coinos.io/bitcal)
- Биткоин календарь RU ⚡️ `bitcalru@coinos.io` 🔗 [Страничка на Coinos](https://coinos.io/bitcalru)

Support Bitcoin Calendar on [Geyser](https://geyser.fund/project/bitcoincalendar)

...or 👇

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)