# Nostr Calendar Bot

## Overview

The Nostr Calendar Bot is a Go-based application that reads events from a CSV file and publishes them to Nostr relays. This bot is designed to automate the posting of calendar events, making it easier to share important dates and information about Bitcoin and cryptocurrency history.

## Features

- Reads historical events from a CSV file
- Automatically publishes events to multiple Nostr relays when the date matches the current day
- Formats posts with proper spacing for readability
- Implements structured JSON logging for easy debugging and monitoring
- Configurable to run as a cron job for daily automation
- Waits between posts to avoid flooding relays
- Handles multiple paragraphs and links properly

## Requirements

- Go 1.18 or higher
- Access to Nostr relays
- A CSV file formatted with event data
- Private key for the Nostr account

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/bitcoin21ideas/nostr-calendar-bot.git
   cd nostr-calendar-bot
   ```

2. Build the application:
   ```bash
   go build -o nostr_bot main.go
   ```

3. Configure your `events.csv` file in the same directory as the executable. The CSV should have this format:
   ```
   "date","title","description"
   "YYYY-MM-DD","Event Title","Event Description |https://link.to/resource"
   ```

   Add `|` where you want to add a line break (does not work after the media link if it is fetched as preview by Nostr clients).

4. Set up a cron job to run the bot daily:
   ```bash
   crontab -e
   ```
   Add the following line to schedule the bot to run at 4:00 AM UTC:
   ```bash
   00 04 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot >> /path/to/nostr_bot.log 2>&1
   ```

## Logging

The bot uses structured logging in JSON format, which allows for better parsing and analysis. Logs are rotated using the `lumberjack` package to manage log file sizes. The log file (`nostr_bot.log`) contains information about:

- Event processing status
- Relay connection success/failures
- Publishing attempts and results
- Performance metrics

## Future Improvements

We are planning the following enhancements to the Calendar Bot:

### Multilingual Support
- Add support for English, Spanish, and German versions of events
- Create dedicated Nostr accounts for each language

### Enhanced Event Content
- Improve event descriptions with relevant hashtags for better discoverability
- Ensure all events have media files (images) for visual appeal

### Community Contributions
- Provide documentation and tools for users to add events in their local language
- Create a submission process for community-contributed events

### Extended Integrations
- Automate cross-posting to Telegram channels
- Explore integration with Twitter/X and other social platforms
- Add support for calendar subscription formats (iCal, etc.)

### Expanded Nostr Capabilities
- Look into implementing support for kind:20 events
- Explore zap-splitting for supporting content contributors

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/bitcoin21ideas/nostr-calendar-bot/blob/main/LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements. If you'd like to contribute to any of the future improvements listed above, please reach out to coordinate efforts.

## Contact

For any questions or feedback, please contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) via Nostr.
