[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)

# Bitcoin Calendar Bot

### Archiving and relaying every Bitcoin milestone 

![bitcal-logo-wide](https://i.nostr.build/dOwtfOe0dvsriH7K.png)

## Overview

The Bitcoin Calendar Bot is a Go-based application that reads events from CSV files and publishes them to Nostr relays. This bot automates the posting of calendar events about Bitcoin history. The bot currently supports English and Russian versions, with plans to add more languages.

## Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
   cd nostr-calendar-bot
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

4. **Run the bot**
   ```bash
   # For English events
   ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
   
   # For Russian events
   ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
   ```

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

Support Bitcoin calendar via [Coinos](https://coinos.io/)
- Bitcoin Calendar EN ⚡️ `bitcal@coinos.io` 🔗 [Coinos page](https://coinos.io/bitcal)
- Биткоин календарь RU ⚡️ `bitcalru@coinos.io` 🔗 [Страничка на Coinos](https://coinos.io/bitcalru)

...or 👇

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
