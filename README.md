[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)

# Bitcoin Calendar Bot

### Archiving and relaying every Bitcoin milestone 

![bitcal-logo-wide](https://haven.bitcoin-calendar.org/9db5c5d32dff9f016bda92280deb4e14e137856704499827f5f0e6d83d7cb326.webp)

## Overview

The Bitcoin Calendar Bot is a Go-based application that reads events from CSV files and publishes them to Nostr relays. This bot automates the posting of calendar events about Bitcoin history. The bot supports English and Russian versions.

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
   cd calendar-bot
   ```

2. **Build the application**
   ```bash
   go build -o nostr_bot main.go metrics.go
   ```

3. **Set up environment variables**
   Create a `.env` file in the project directory:
   ```
   NOSTR_PRIVATE_KEY_1=your_first_private_key_hex_here
   NOSTR_PRIVATE_KEY_2=your_second_private_key_hex_here 
   ```

   Or copy the `.env-example` file and adjust its contents as needed.

4. **Run the bot**
   ```bash
   # For English events
   ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
   
   # For Russian events
   ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
   ```

### Automated Daily Posting with Cron

To ensure the bot posts events daily, you can set up a cron job. Find details and examples of how to setup a cronjob in our Usage Guide.

## Documentation

- [Installation Guide](docs/INSTALLATION.md) - Detailed installation and setup instructions
- [Usage Guide](docs/USAGE.md) - How to use the bot effectively and available configuration options
- [CSV Format](docs/CSV_FORMAT.md) - How to structure your CSV event files
- [Development Guide](docs/DEVELOPMENT.md) - Information for developers
- [Roadmap](docs/ROADMAP.md) - Future development plans

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

## Support

Support Bitcoin Calendar via [Coinos](https://coinos.io/)
- Bitcoin Calendar EN ⚡️ `bitcal@coinos.io` 🔗 [Coinos page](https://coinos.io/bitcal)
- Биткоин календарь RU ⚡️ `bitcalru@coinos.io` 🔗 [Страничка на Coinos](https://coinos.io/bitcalru)

Support Bitcoin Calendar on [Geyser](https://geyser.fund/project/bitcoincalendar)

...or 👇

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)