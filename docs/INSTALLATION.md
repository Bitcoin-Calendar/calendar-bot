# Installation Guide

This guide provides detailed instructions for installing and setting up the Bitcoin Calendar Bot.

## Prerequisites

- Go 1.18 or higher
- Git
- Access to Nostr relays
- Nostr private keys for posting (one for each language)

## Step-by-Step Installation

### 1. Clone the Repository

```bash
git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
cd nostr-calendar-bot
```

### 2. Install Dependencies

The project uses Go modules to manage dependencies. Run the following command to install all required dependencies:

```bash
go mod download
```

### 3. Build the Application

Build the application using the Go compiler:

```bash
go build -o nostr_bot main.go metrics.go
```

### 4. Configure Environment Variables

Create a `.env` file in the project directory to store your environment variables:

```
# Required - Nostr private keys for each language
NOSTR_PRIVATE_KEY_EN=your_english_private_key_hex_here
NOSTR_PRIVATE_KEY_RU=your_russian_private_key_hex_here

# Optional - Logging configuration
LOG_LEVEL=info      # Options: debug, info, warn, error (default: info)
LOG_DIR=/path/to/logs  # Custom log directory (default: current directory)
CONSOLE_LOG=true    # Set to true to also output logs to console (default: false)
DEBUG=true          # Legacy option to enable debug mode (use LOG_LEVEL=debug instead)
```

> **Important:** Ensure that your `.env` file is included in your `.gitignore` to prevent it from being committed to your repository.

### 5. Prepare CSV Event Files

Ensure your CSV files are properly formatted and placed in the project directory:

- `events_en.csv` - Events in English
- `events_ru.csv` - Events in Russian

For details on the CSV format, see the [CSV Format Guide](CSV_FORMAT.md).

## Setting Up for Production

### Automated Execution with Cron

To schedule the bot to run automatically at specific times, set up cron jobs:

1. Edit your crontab:
   ```bash
   crontab -e
   ```

2. Add entries for each language. For example:
   ```bash
   # For English events at 12 PM UTC
   00 12 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
   
   # For Russian events at 4 AM UTC
   00 04 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
   ```

### Docker Installation (Optional)

For containerized deployment:

1. Build the Docker image:
   ```bash
   docker build -t bitcoin-calendar-bot .
   ```

2. Run the container:
   ```bash
   docker run -v /path/to/your/env:/app/.env -v /path/to/your/csv:/app/data bitcoin-calendar-bot events_en.csv NOSTR_PRIVATE_KEY_EN
   ```

## Troubleshooting

### Common Issues

1. **Permission Denied**
   
   If you encounter permission issues when executing the bot:
   ```bash
   chmod +x nostr_bot
   ```

2. **Missing Environment Variables**
   
   Ensure your `.env` file is in the same directory as the executable and contains all required variables.

3. **CSV File Not Found**
   
   Verify that your CSV files are in the correct location and properly named.

4. **Build Errors**
   
   Make sure to include all Go files in your build command:
   ```bash
   go build -o nostr_bot main.go metrics.go
   ```

## Next Steps

Once you've completed the installation, refer to the [Usage Guide](USAGE.md) for information on how to use the bot effectively. 

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
