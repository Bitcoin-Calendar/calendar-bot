# Nostr Calendar Bot

## Overview

The Nostr Calendar Bot is a Go-based application that reads events from a CSV file and publishes them to Nostr relays. This bot is designed to automate the posting of calendar events, making it easier to share important dates and information.

## Features

- Reads events from a CSV file.
- Publishes events to multiple Nostr relays.
- Structured logging for easy debugging and monitoring.
- Configurable to run as a cron job for daily automation.

## Requirements

- Go 1.16 or higher
- Access to Nostr relays
- A CSV file formatted with event data

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/nostr-calendar-bot.git
   cd nostr-calendar-bot
   ```

2. Build the application:
   ```bash
   go build -o nostr_bot main.go
   ```

3. Configure your `events.csv` file in the same directory as the executable.

4. Set up a cron job to run the bot daily:
   ```bash
   crontab -e
   ```
   Add the following line to schedule the bot:
   ```bash
   00 04 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot >> /path/to/nostr_bot.log 2>&1
   ```

## Logging

The bot uses structured logging in JSON format, which allows for better parsing and analysis. Logs are rotated using the `lumberjack` package to manage log file sizes.

## License

See the [LICENSE](/calendar-bot/LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements.

## Contact

For any questions or feedback, please contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) via Nostr.
