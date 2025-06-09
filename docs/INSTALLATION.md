# Bitcoin Calendar Bot Installation Guide

This guide provides detailed instructions for installing and setting up the Bitcoin Calendar Bot using Docker and Docker Compose, which is the recommended method.

## Recommended: Docker Installation

### 1. Prerequisites

- Ensure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.
- Git

### 2. Clone the Repository
```bash
git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
cd calendar-bot
```

### 3. Configure Environment Variables

Copy the example `.env` file and customize it with your details.

```bash
cp .env-example .env
```

Open the `.env` file and add your API configuration and Nostr private keys:

```env
# --- API Configuration (Required) ---
BOT_API_ENDPOINT="http://your_api_vps_ip:port/api" # Replace with your API's base URL
BOT_API_KEY="your_secret_api_key"             # Replace with your API key

# --- Nostr Private Keys (Required) ---
# Used by the production `nostr-bot-en` service
NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex"

# Used by the `nostr-bot-en-test` service for manual runs
NOSTR_PRIVATE_KEY_ENT="your_english_TEST_private_key_hex"

# --- Optional: Logging Configuration ---
# LOG_LEVEL="debug" # Or "info", "warn", "error"
# CONSOLE_LOG="true"  # Or "false"
```
The `BOT_PROCESSING_LANGUAGE` is set to `en` within the `docker-compose.yml` file and does not need to be set here.

### 4. Build the Docker Image
All services defined in `docker-compose.yml` use the same base image, which you must build first.

```bash
docker-compose build
```

### 5. Verify Setup (Run a Test Bot)

To ensure your configuration is correct, run the test bot service. This will attempt to fetch events and simulate posting without affecting your production keys.

```bash
docker-compose run --rm nostr-bot-en-test
```
Check the console output for any errors. The bot should fetch events for the current day and attempt to post them.

### 6. Automated Operation (Cron Job)

Once you have confirmed the test bot is working, you can set up a cron job to run the production service (`nostr-bot-en`) automatically.

Refer to the [Automated Operation section in the Usage Guide](USAGE.md#automated-operation) for detailed instructions on configuring cron jobs with `docker-compose run`.

### 7. Rebuilding After Code Changes
If you modify the Go source code, you must rebuild the Docker image for the changes to take effect.
```bash
docker-compose build
```

---

## Manual Installation (Without Docker - Deprecated)

Setting up and running the bot manually is not recommended. The Docker-based approach is strongly preferred for its simplicity and reliability.

If you must run manually:

1.  **Install Go**: Ensure you have Go 1.18+ installed.
2.  **Build from source**: `go build -o nostr_bot main.go`
3.  **Set Environment Variables**: You must manually set `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE` (to `en`), and the environment variable holding your Nostr private key (e.g., `MY_NOSTR_KEY="actual_hex_key"`) in your shell environment *before* running the bot.
4.  **Run the bot**: `./nostr_bot MY_NOSTR_KEY` (where `MY_NOSTR_KEY` is the *name* of the environment variable holding the private key).

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
