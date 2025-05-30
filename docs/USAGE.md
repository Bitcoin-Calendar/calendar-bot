# Usage Guide

This guide explains how to use the Bitcoin Calendar Bot effectively, primarily focusing on the Docker Compose setup, including configuration options and command examples.

## Basic Operation Principle

The Bitcoin Calendar Bot fetches historical Bitcoin events for the current day (month and day) from a configured API endpoint. It then posts these events to Nostr relays using a specified Nostr private key. Events are posted as Kind 1 (text notes). If an event includes a valid image URL in its `Media` field, the bot will also attempt to post a NIP-68 Kind 20 (picture note) event for that image, in addition to the Kind 1 text note. The language of the events (e.g., English or Russian) is determined by the `BOT_PROCESSING_LANGUAGE` environment variable, which is pre-configured for each service in the `docker-compose.yml` file.

## Running the Bot with Docker Compose

The primary way to run the bot is by using `docker-compose run <service-name>`. The `docker-compose.yml` file defines several services:

-   `nostr-bot-en`: Production service for English posts.
-   `nostr-bot-ru`: Production service for Russian posts.
-   `nostr-bot-en-test`: Test service for English posts (uses test keys).
-   `nostr-bot-ru-test`: Test service for Russian posts (uses test keys).

### Examples:

```bash
# Run the English production bot (typically for a cron job)
docker-compose run --rm nostr-bot-en

# Run the Russian test bot manually
docker-compose run --rm nostr-bot-ru-test
```
Each service is configured with the correct `BOT_PROCESSING_LANGUAGE` and the name of the environment variable that holds its corresponding Nostr private key.

## Configuration Options

The Bitcoin Calendar Bot is configured using environment variables. Most of these are set in the `.env` file, while `BOT_PROCESSING_LANGUAGE` is set per-service in `docker-compose.yml`.

### Required Environment Variables (in `.env` file)

-   `BOT_API_ENDPOINT`: Full base URL of the Bitcoin Historical Events API (e.g., `http://your_api_ip:port/api`).
-   `BOT_API_KEY`: Your secret API key for the events API.
-   `NOSTR_PRIVATE_KEY_EN`: Hexadecimal private key for posting English events (production).
-   `NOSTR_PRIVATE_KEY_RU`: Hexadecimal private key for posting Russian events (production).
-   `NOSTR_PRIVATE_KEY_ENT`: Hexadecimal private key for posting English events (testing).
-   `NOSTR_PRIVATE_KEY_RUT`: Hexadecimal private key for posting Russian events (testing).

### Logging Configuration (can be set in `.env` or defaults in `docker-compose.yml` used)

| Variable      | Description                                     | Options                        | Default (in `docker-compose.yml` for prod/test) | 
|---------------|-------------------------------------------------|--------------------------------|-------------------------------------------------| 
| `LOG_LEVEL`   | Sets the logging verbosity.                     | `debug`, `info`, `warn`, `error` | `info` (prod), `debug` (test)                   |
| `LOG_DIR`     | Custom directory *inside the container* for log files. | Any valid directory path       | `/app/logs` (mounted to host `./logs`)          |
| `CONSOLE_LOG` | Output logs to console in addition to files.    | `true`, `false`                | `false` (prod), `true` (test)                   |

### Examples of Overriding Logging in `.env`:

To make all services (including production) log at debug level to console:
```env
# In .env file
LOG_LEVEL="debug"
CONSOLE_LOG="true"
```

## Log Files

The bot automatically creates log files named `nostr_bot.log` within the directory specified by `LOG_DIR` (inside the container). This directory is mapped to `./logs` on your host machine by default in `docker-compose.yml`.

-   The log file contains entries from all runs, regardless of language, as it's a single log file now inside the container. The language being processed is indicated in the log messages themselves.

### Log Rotation

Log files are automatically rotated when they reach 10MB, with the following settings:
-   Maximum of 3 backup files are kept.
-   Log files older than 28 days are deleted.
-   Rotated log files are compressed to save space.

## Event Processing Flow

The bot performs the following steps:

1.  Reads its configuration (API endpoint, API key, Nostr private key name, processing language, relays, etc.) using the `internal/config` module.
2.  Sets up logging using the `internal/logging` module.
3.  Initializes clients and services: API client (`internal/api`), metrics collector (`internal/metrics`), Nostr event publisher and image validator (`internal/nostr`).
4.  Fetches events for the current calendar day (month and day) from the API, for the configured language, using the API client.
5.  For each matching `APIEvent`:
    *   Generates a unique request ID for tracking (this is part of the logger context usually).
    *   **Kind 1 Event**: Creates a Kind 1 (text) Nostr event using `nostr.CreateKind1NostrEvent()`.
    *   Publishes the Kind 1 event to configured Nostr relays via `eventPublisher.PublishEvent()`. Updates Kind 1 metrics.
    *   **Kind 20 Event (if applicable)**: If the `APIEvent.Media` field contains a valid image URL, it creates a NIP-68 Kind 20 (picture) Nostr event using `nostr.CreateKind20NostrEvent()` (which includes image validation).
    *   Publishes the Kind 20 event to relays via `eventPublisher.PublishEvent()`. Updates Kind 20 metrics.
    *   If multiple events are found for the day, and the Kind 1 event for the current API event was successfully published to at least one relay, it waits 30 minutes before processing the next API event from the list.
6.  Logs a summary of collected metrics using `metricsCollector.LogSummary()`.

## Running Test Instances

Before setting up automated cron jobs for production, it's highly recommended to test your setup using the dedicated test services. These services use the test Nostr private keys defined in your `.env` file (`NOSTR_PRIVATE_KEY_ENT`, `NOSTR_PRIVATE_KEY_RUT`) to avoid posting to your main accounts during testing.

To run the test instances:

```bash
# Test the English bot with test keys
docker-compose run --rm nostr-bot-en-test

# Test the Russian bot with test keys
docker-compose run --rm nostr-bot-ru-test
```

Key characteristics of test instances as defined in `docker-compose.yml`:

*   **Language Specific**: `nostr-bot-en-test` is configured for English, `nostr-bot-ru-test` for Russian via `BOT_PROCESSING_LANGUAGE`.
*   **Test Keys**: They use `NOSTR_PRIVATE_KEY_ENT` and `NOSTR_PRIVATE_KEY_RUT` respectively.
*   **Verbose Logging**: They default to more verbose logging settings (`LOG_LEVEL=debug`, `CONSOLE_LOG=true`).
*   **Logs**: By default, logs from test runs will go into the same `./logs` directory on your host machine as production runs. If you prefer to keep test outputs separate, you can modify the `volumes` in `docker-compose.yml` for the test services (e.g., change `./logs:/app/logs` to `./logs-test:/app/logs`).

## Automated Operation (Cron Jobs)

For regular posting, set up cron jobs to run the production bot services automatically. The bot is designed to run, process events for the current day, and then exit.

### With Docker (Recommended)

Use `docker-compose run --rm <service-name>` in your cron jobs. This starts the specified production service, which runs its defined command and then exits. The `--rm` flag ensures the container is removed after completion.

Modify your crontab (e.g., by running `crontab -e`) to include entries similar to the following. Ensure the path to your project directory (where `docker-compose.yml` is located) is correct. `docker-compose run` automatically loads variables from `.env` in that directory.

```cron
# Bitcoin Calendar Bot Cron Jobs
# Ensure the path to your project directory (where docker-compose.yml is) is correct.

# Example: Run Russian bot at 04:00 AM daily (UTC)
# Replace /path/to/your/calendar-bot with the actual absolute path to your project.
00 04 * * * cd /path/to/your/calendar-bot && docker-compose run --rm nostr-bot-ru

# Example: Run English bot at 15:00 PM daily (UTC)
# Replace /path/to/your/calendar-bot with the actual absolute path to your project.
00 15 * * * cd /path/to/your/calendar-bot && docker-compose run --rm nostr-bot-en
```

Make sure to adjust the schedule and the path (`/path/to/your/calendar-bot`) to match your setup.

### Manual Setup (Deprecated)

Running the bot manually without Docker is not recommended for production or cron jobs due to the difficulty of managing distinct configurations (especially `BOT_PROCESSING_LANGUAGE`) for different language instances. The Docker setup handles this cleanly via services.

If you were to run it manually, you would need to set `BOT_API_ENDPOINT`, `BOT_API_KEY`, `BOT_PROCESSING_LANGUAGE`, and the relevant `NOSTR_PRIVATE_KEY_...` in the environment before executing:

```bash
# Example for a manual English run (variables must be pre-set in shell):
# BOT_PROCESSING_LANGUAGE="en" LOG_LEVEL=debug ./nostr_bot NOSTR_PRIVATE_KEY_EN
```

## Troubleshooting

If you encounter issues, check the log files (in `./logs` on the host by default) for detailed error messages. The console output from `docker-compose run` (especially for test services) will also show logs.

-   Ensure all required environment variables are correctly set in your `.env` file.
-   Verify API connectivity and key validity.
-   Confirm Nostr private keys are correct and in hex format.
-   Check `docker-compose.yml` for correct service configurations, especially `BOT_PROCESSING_LANGUAGE` and the private key variable names passed in the `command`.

Common issues and their solutions are also documented in the [Installation Guide](INSTALLATION.md#troubleshooting).

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
