# Installation Guide

This guide provides detailed instructions for installing and setting up the Bitcoin Calendar Bot using Docker and Docker Compose, which is the recommended method.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed
- [Docker Compose](https://docs.docker.com/compose/install/) installed
- Git
- API Endpoint and Key for the Bitcoin Historical Events API
- Nostr private keys for posting (e.g., one for English posts, and a test key)

## Step-by-Step Installation (Docker & Docker Compose)

### 1. Clone the Repository

```bash
git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
cd calendar-bot
```

### 2. Configure Environment Variables

Create a `.env` file by copying the example. This file will store your API credentials and Nostr private keys.

```bash
cp .env-example .env
```

Now, edit the `.env` file with your actual credentials:

```env
# --- API Configuration (Required) ---
BOT_API_ENDPOINT="http://your_api_vps_ip:port/api" # Replace with your API's base URL
BOT_API_KEY="your_secret_api_key"             # Replace with your API key

# --- Nostr Private Keys (Required) ---
# Production Keys (used by cron jobs)
NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex"

# Test Keys (used by -test services for manual runs)
NOSTR_PRIVATE_KEY_ENT="your_english_TEST_private_key_hex"

# --- Optional: Logging Configuration (defaults are set in docker-compose.yml) ---
# These can be uncommented and set here to override service defaults if needed.
# LOG_LEVEL="debug"
# CONSOLE_LOG="true"
# LOG_DIR="/app/logs" # Path inside the container, usually mapped to a host volume
```

> **Important:** The `.env` file contains sensitive keys. Ensure it is listed in your `.gitignore` file (it should be by default) to prevent accidental commits.
> The `BOT_PROCESSING_LANGUAGE` is configured per-service within the `docker-compose.yml` and should not be set globally in the `.env` file.

### 3. Build the Docker Image

Build the Docker image which will be used by all services defined in `docker-compose.yml`:

```bash
docker-compose build
```

Alternatively, you can build the image for a specific service (which also builds the base image if not present):
```bash
# docker-compose build nostr-bot-en
```

### 4. Verify Setup (Run a Test Bot)

To ensure your configuration is correct, run the test bot service. This will attempt to fetch events and simulate posting without affecting your production keys.

# Test the English bot
```bash
docker-compose run --rm nostr-bot-en-test
```

### 5. Automated Operation (Cron Job)

Once you have confirmed the test bot is working, you can set up a cron job to run the production service (`nostr-bot-en`) automatically.

Refer to the [Automated Operation section in the Usage Guide](USAGE.md#automated-operation) for detailed instructions on configuring cron jobs with `docker-compose run`.

## Troubleshooting

- **Connection Errors to API**: Ensure the `BOT_API_ENDPOINT` in your `.env` file is correct and that the API server is running and accessible from where Docker is executing.
- **Authentication Errors**: Double-check your `BOT_API_KEY` and ensure the `NOSTR_PRIVATE_KEY...` variables in your `.env` match the keys expected by the `docker-compose.yml` (`en` for English services).
- **Configuration Issues**: If the bot fails to start, check the logs for messages about missing environment variables. You can view logs via `docker-compose logs nostr-bot-en`.
- **Permissions**: Ensure your user has the necessary permissions to run Docker and Docker Compose.

## Manual Installation (Without Docker)

Setting up and running the bot manually without Docker is deprecated due to the complexities of managing environment variables (like `BOT_PROCESSING_LANGUAGE`) per instance. The Docker-based approach is strongly recommended for its simplicity and reliability.

If you must run manually:

1.  **Install Go**: Ensure you have Go 1.18+ installed.
2.  **Build from source**: `go build -o nostr_bot main.go`
3.  **Set Environment Variables**: You must manually set `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE` (to `en`), and the environment variable holding your Nostr private key (e.g., `MY_NOSTR_KEY="actual_hex_key"`) in your shell environment *before* running the bot.
4.  **Run the bot**: `./nostr_bot MY_NOSTR_KEY` (where `MY_NOSTR_KEY` is the *name* of the environment variable holding the private key).

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
