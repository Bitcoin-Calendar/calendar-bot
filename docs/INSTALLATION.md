# Installation Guide

This guide provides detailed instructions for installing and setting up the Bitcoin Calendar Bot using Docker and Docker Compose, which is the recommended method.

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed
- [Docker Compose](https://docs.docker.com/compose/install/) installed
- Git
- API Endpoint and Key for the Bitcoin Historical Events API
- Nostr private keys for posting (e.g., one for English posts, one for Russian posts, and test keys)

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
NOSTR_PRIVATE_KEY_RU="your_russian_specific_private_key_hex"

# Test Keys (used by -test services for manual runs)
NOSTR_PRIVATE_KEY_ENT="your_english_TEST_private_key_hex"
NOSTR_PRIVATE_KEY_RUT="your_russian_TEST_private_key_hex"

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

Before setting up cron jobs, test one of the pre-configured test services. These use your test Nostr keys and are set to output logs to the console.

```bash
# Test the English bot
docker-compose run --rm nostr-bot-en-test

# Test the Russian bot
docker-compose run --rm nostr-bot-ru-test
```

Check the console output for any errors. The bot should fetch events for the current day for the specified language and attempt to post them.

## Setting Up for Production (Automated Cron Jobs)

Once you have confirmed the test bots are working, you can set up cron jobs to run the production services (`nostr-bot-en` and `nostr-bot-ru`) automatically.

Refer to the [Automated Operation section in the Usage Guide](USAGE.md#automated-operation) for detailed instructions on configuring cron jobs with `docker-compose run`.

## Troubleshooting

- **Connection Errors to API**: Ensure the `BOT_API_ENDPOINT` in your `.env` file is correct and that the API server is running and accessible from where Docker is executing.
- **Authentication Errors with API**: Verify that `BOT_API_KEY` in your `.env` file is correct.
- **Nostr Posting Issues**: Double-check that the `NOSTR_PRIVATE_KEY_...` variables in your `.env` file are correct (hex format, no `nsec` prefix) and correspond to the key names used in the `command` section of your `docker-compose.yml` services.
- **Language Mismatch**: Ensure `BOT_PROCESSING_LANGUAGE` is correctly set in the `environment` section of each service in `docker-compose.yml` (`en` for English services, `ru` for Russian services).
- **Build Issues**: If `docker-compose build` fails, check the output for errors. Ensure your `Dockerfile` is correct and your Go environment (if any local Go tools are invoked indirectly) is sound.

For more general usage information, refer to the [Usage Guide](USAGE.md).

## Legacy Manual Setup (Deprecated)

Setting up and running the bot manually without Docker is deprecated due to the complexities of managing environment variables (like `BOT_PROCESSING_LANGUAGE`) per instance. The Docker-based approach is strongly recommended.

If you must proceed with a manual setup:

1.  **Install Go:** Ensure Go 1.18 or higher is installed.
2.  **Build the application:** `go build -o nostr_bot main.go`
3.  **Set Environment Variables:** You must manually set `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE` (to `en` or `ru`), and the environment variable holding your Nostr private key (e.g., `MY_NOSTR_KEY="actual_hex_key"`) in your shell environment *before* running the bot.
4.  **Run the bot:** `./nostr_bot MY_NOSTR_KEY` (where `MY_NOSTR_KEY` is the *name* of the environment variable holding the private key).

This manual method lacks the per-service configuration benefits of Docker Compose, making it harder to manage different language instances reliably.

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
