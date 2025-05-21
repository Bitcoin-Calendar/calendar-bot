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

If you specify a custom `LOG_DIR`, logs will be stored in that directory with the same naming convention.

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

## Automated Operation

For regular posting, it's recommended to set up cron jobs to run the bot automatically:

```bash
# For English events at 12 PM UTC
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
