# Usage Guide

This guide explains how to use the Bitcoin Calendar Bot effectively, including available configuration options and command examples.

## Basic Usage

The basic syntax for running the Bitcoin Calendar Bot is:

```bash
./nostr_bot <csv_file_path> <env_var_for_private_key>
```

### Parameters:

- `<csv_file_path>`: Path to the CSV file containing the events to be posted
- `<env_var_for_private_key>`: Name of the environment variable containing the Nostr private key

### Examples:

```bash
# Run the bot with English events
./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

# Run the bot with Russian events
./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU

# Run the bot with a custom CSV file
./nostr_bot custom_events.csv CUSTOM_NOSTR_KEY
```

## Configuration Options

The Bitcoin Calendar Bot can be configured using environment variables, either set in your shell or in the `.env` file.

### Required Environment Variables

These must be set in your `.env` file:

- `NOSTR_PRIVATE_KEY_EN`: Private key for posting English events
- `NOSTR_PRIVATE_KEY_RU`: Private key for posting Russian events

### Logging Configuration

Control the logging behavior with these optional environment variables:

| Variable | Description | Options | Default |
|----------|-------------|---------|---------|
| `LOG_LEVEL` | Sets the logging verbosity | `debug`, `info`, `warn`, `error` | `info` |
| `LOG_DIR` | Custom directory for log files | Any valid directory path | Current directory |
| `CONSOLE_LOG` | Output logs to console in addition to files | `true`, `false` | `false` |
| `DEBUG` | Legacy option to enable debug mode | `true`, `false` | `false` |

### Examples:

```bash
# Run with debug logging
LOG_LEVEL=debug ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

# Run with console output
CONSOLE_LOG=true ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

# Run with custom log directory
LOG_DIR=/var/log/nostr-calendar-bot ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

# Combine multiple options
LOG_LEVEL=debug CONSOLE_LOG=true LOG_DIR=/var/log/nostr-calendar-bot ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
```

## Log Files

The bot automatically creates log files based on the language of the CSV file:

- English events: `nostr_bot_en.log` 
- Russian events: `nostr_bot_ru.log`
- Other (unidentified) languages: `nostr_bot_unknown.log`

If you specify a custom `LOG_DIR` (e.g., via environment variable or Docker volume mount), logs will be stored in that directory with the same naming convention. 

When using Docker, the logs are typically mounted to a directory on your host machine (e.g., `./logs` in the provided `docker-compose.yml`). You can view them by navigating to this directory on your host system.
For test services run with `docker-compose run --rm nostr-bot-en-test`, logs are also output to the console by default due to `CONSOLE_LOG=true` in their service definition.

### Log Rotation

Log files are automatically rotated when they reach 10MB, with the following settings:
- Maximum of 3 backup files are kept
- Log files older than 28 days are deleted
- Rotated log files are compressed to save space

## Event Processing

The bot processes the CSV file to find events that match the current date:

1. It scans the CSV file for events matching today's date (MM-DD format)
2. For each matching event, it:
   - Generates a unique request ID for tracking
   - Formats the event content with proper line breaks
   - Signs the event with your private key
   - Publishes the event to multiple Nostr relays
   - Waits 30 minutes before posting the next event (if multiple events match today's date)

## Metrics Collection

The bot collects and reports metrics about its operation:

- Number of events posted
- Number of events skipped
- Number of events that failed to post
- Success/failure statistics for each relay
- Performance metrics (connection and publishing times)

These metrics are logged and exported to a JSON file named `metrics_<language>_<date>.json` at the end of each run.

## Running Test Instances

Before setting up automated cron jobs for production, it's highly recommended to test your setup using the dedicated test services. These services use the test Nostr private keys defined in your `.env` file (`NOSTR_PRIVATE_KEY_ENT`, `NOSTR_PRIVATE_KEY_RUT`) to avoid posting to your main accounts during testing.

To run the test instances:

```bash
# Test the English bot with test keys
docker-compose run --rm nostr-bot-en-test

# Test the Russian bot with test keys
docker-compose run --rm nostr-bot-ru-test
```

Key characteristics of test instances:

*   **Verbose Logging**: They default to more verbose logging settings (`LOG_LEVEL=debug`, `CONSOLE_LOG=true`) as defined in `docker-compose.yml`. This helps in debugging and verifying behavior.
*   **Shared Event Files**: They use the same `events_en.csv` and `events_ru.csv` files as the production bots.
*   **Logs and Metrics**: By default, logs and metrics from test runs will go into the same `./logs` and `./metrics` directories on your host machine. If you prefer to keep test outputs separate, you can modify the `volumes` in `docker-compose.yml` for the test services (e.g., change `./logs:/app/logs` to `./logs-test:/app/logs`).

## Automated Operation

For regular posting, it's recommended to set up cron jobs to run the bot automatically. The bot is designed to run, process events for the current day, and then exit.

### With Docker (Recommended)

If you are using the Docker setup, you will use `docker-compose run --rm <service-name>` in your cron jobs. This starts the specified production service, which runs its defined command and then exits. The `--rm` flag ensures the container is removed after completion.

Modify your crontab (e.g., by running `crontab -e`) to include entries similar to the following. Ensure the path to your project directory (where `docker-compose.yml` is located) is correct.
When you use `docker-compose run`, it automatically looks for and loads environment variables from a `.env` file located in the same directory as your `docker-compose.yml` file, making them available to the container. This means your `NOSTR_PRIVATE_KEY_EN`, `NOSTR_PRIVATE_KEY_RU`, etc., defined in `.env` will be accessible to the bot services.

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

### Manual Setup (Without Docker)

If you are running the bot manually (not using Docker), you can set up cron jobs as follows:

Important: Cron jobs run with a minimal environment and typically do not automatically source your shell's profile files (like `.bashrc` or `.zshrc`) or the project's `.env` file. You must ensure that the required environment variables (e.g., `NOSTR_PRIVATE_KEY_EN`, `LOG_DIR` if you set it outside the command) are available when the cron job executes. Common methods include:
*   Defining the variables directly within the crontab itself (e.g., `NOSTR_PRIVATE_KEY_EN=yourkey 00 12 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN`).
*   Creating a small wrapper script that first sources a file exporting the variables (or exports them directly) and then calls the bot. For example, `00 12 * * * /path/to/run_english_bot.sh`, where `run_english_bot.sh` might contain:
    ```bash
    #!/bin/bash
    # Adjust path to your .env file or define variables directly
    # export $(grep -v '^#' /path/to/your/calendar-bot/.env | xargs) # Example for .env sourcing
    cd /path/to/nostr-calendar-bot
    LOG_DIR=./logs ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
    ```
    Ensure the wrapper script is executable (`chmod +x /path/to/run_english_bot.sh`).
*   Some systems allow prefixing the command with `env VAR=value` in the crontab, but this can vary.

Make sure to replace `/path/to/nostr-calendar-bot` with the actual absolute path to your project and adjust variable definitions as needed.

```bash
# For English events at 12 PM UTC
# Example assuming variables are set within the cron environment or a wrapper script:
00 12 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

# For Russian events at 4 AM UTC
00 04 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
```

## Troubleshooting

If you encounter issues, check the log files for detailed error messages and information.

Enable debug logging for more verbose output:

```bash
LOG_LEVEL=debug ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
```

Common issues and their solutions are documented in the [Installation Guide](INSTALLATION.md#troubleshooting). 

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
