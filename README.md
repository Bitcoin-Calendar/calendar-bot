[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-purple)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)


# Bitcoin Calendar Bot

## Overview

The Bitcoin Calendar Bot is a Go-based application that reads events from CSV files and publishes them to Nostr relays. This bot is designed to automate the posting of calendar events, making it easier to share important dates and information about Bitcoin history. The bot currently supports English and Russian versions, with plans of adding more languages.

## Features

- Reads historical events from CSV files
- Automatically publishes events to multiple Nostr relays when the date matches the current day
- Formats posts with proper spacing for readability
- Implements structured JSON logging for easy debugging and monitoring
- Configurable to run as a cron job for daily automation
- Waits between posts to avoid flooding relays
- Handles multiple paragraphs and links properly

## Requirements

- Go 1.18 or higher
- Access to Nostr relays
- CSV files formatted with event data
- Private keys for the Nostr accounts

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

3. Configure your CSV files in the same directory as the executable. The CSV should have this format:
   ```
   "date","title","description"
   "YYYY-MM-DD","Event Title","Event Description |https://link.to/resource"
   ```
> Use the `|` character when wanting to introduce a line break. 

4. Set Up Environment Variables:
   - Create a `.env` file in the project directory to store your environment variables. This file should not be committed to version control.
   - Add the following lines to your `.env` file to set your Nostr private keys:
     ```
     NOSTR_PRIVATE_KEY_EN=your_english_private_key_hex_here
     NOSTR_PRIVATE_KEY_RU=your_russian_private_key_hex_here
     ```

> You can use this outsid of the bitcoin calendar bot context for posting events from different CSV files at specific times on specific dates.

5. Load Environment Variables:
   - The application uses the `github.com/joho/godotenv` package to load environment variables from the `.env` file. Ensure this package is installed and properly configured in your project.

6. Set up Cron Jobs:
   - Schedule the bot to run at specific times for each language version using cron jobs.

   Example Cron Job Entries:
   - For English events at 12 PM UTC:
     ```bash
     00 12 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN >> /path/to/nostr_bot_en.log 2>&1
     ```
   - For Russian events at 4 AM UTC:
     ```bash
     00 04 * * * cd /path/to/nostr-calendar-bot && ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU >> /path/to/nostr_bot_ru.log 2>&1
     ```

## Environment Variables

- NOSTR_PRIVATE_KEY_EN: This environment variable stores your Nostr private key for the English account.
- NOSTR_PRIVATE_KEY_RU: This environment variable stores your Nostr private key for the Russian account.
- **Keep Your `.env` File Secure**: Ensure that your `.env` file is included in your `.gitignore` to prevent it from being committed to your repository.

## Logging

The bot uses structured logging in JSON format, which allows for better parsing and analysis. Logs are rotated using the `lumberjack` package to manage log file sizes. The log files contain information about:

- Event processing status
- Relay connection success/failures
- Publishing attempts and results
- Performance metrics

## Nostr Stack

The bot utilizes the Nostr protocol to publish events to relays. It leverages the `github.com/nbd-wtf/go-nostr` package for interacting with Nostr relays, ensuring seamless integration and communication.

## Future Improvements

We are planning the following enhancements to the Calendar Bot:

### Multilingual Support
- Add support for additional languages beyond English and Russian
- Create dedicated Nostr accounts for each language

### Enhanced Event Content
- Improve event descriptions with relevant hashtags for better discoverability
- Ensure all events have media files (images) for visual appeal

### Community Contributions
- Provide documentation and tools for users to add events in their local language
- Create a submission process for community-contributed events

### Extended Integrations
- Automate cross-posting to Telegram channels
- Add support for calendar subscription formats (iCal, etc.)

### Expanded Nostr Capabilities
- Look into implementing support for more event kinds
- Explore zap-splitting for supporting content contributors

## License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/bitcoin21ideas/nostr-calendar-bot/blob/main/LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any suggestions or improvements. If you'd like to contribute to any of the future improvements listed above, please reach out to coordinate efforts.

## Contact

For any questions or feedback, please contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) via Nostr.

## Support Bitcoin Calendar

Support Bitcoin calendar via [Coinos](https://coinos.io/)
- Биткоин календарь RU ⚡️ `bitcalru@coinos.io` 🔗 [Страничка на Coinos](https://coinos.io/bitcalru).
- Bitcoin Calendar EN ⚡️ `bitcal@coinos.io` 🔗 [Coinos page](https://coinos.io/bitcal).
