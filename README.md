[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)

# Bitcoin Calendar Bot

### Archiving and relaying every Bitcoin milestone 

![bitcal-structure](https://haven.bitcoin-calendar.org/4d4f81403e39c0c4a454a35cb6913a9420539c9665cb5240fdaff4e27b1e5176.webp)

## Overview

The Bitcoin Calendar Bot is a Go-based application that fetches historical Bitcoin event data from the Bitcoin Historical Events API and publishes these events to Nostr relays. This bot automates the posting of calendar events about Bitcoin history.

Key functionalities include:
*   Fetching events from a configurable API endpoint.
*   Publishing events as Nostr Kind 1 (text-based) notes.
*   **NEW**: Publishing events as Nostr NIP-68 Kind 20 (picture-based) notes for events with associated images.
*   Support for English, configurable at runtime.
*   Metrics collection for monitoring event posting success and failures.

The project has recently undergone a significant refactoring to improve modularity and maintainability. Core functionalities such as configuration management, API interaction, logging, metrics collection, and Nostr event publishing have been moved into dedicated packages within an `internal` directory.

## Getting Started

The recommended way to run the bot is with Docker and Docker Compose. For detailed setup instructions, please see the **[Installation Guide](docs/INSTALLATION.md)**.

A typical workflow involves:
1.  Cloning the repository.
2.  Configuring your API and Nostr keys in a `.env` file.
3.  Building the Docker image with `docker-compose build`.
4.  Running the test bot with `docker-compose run --rm nostr-bot-en-test` to verify your setup.
5.  Setting up a cron job to run the production bot (`nostr-bot-en`) daily.

## Project Structure

The main application logic is organized within the `internal` directory:
*   `internal/config`: Handles loading and validation of application configuration.
*   `internal/api`: Contains the client for interacting with the Bitcoin Calendar events API.
*   `internal/logging`: Manages logger setup and configuration.
*   `internal/metrics`: Provides a collector for tracking application-specific metrics.
*   `internal/models`: Defines shared data structures like `APIEvent`.
*   `internal/nostr`: Handles all Nostr-related operations, including event creation and publishing.

The `main.go` file serves as the entry point, orchestrating these components.

## Documentation

- [Installation Guide](docs/INSTALLATION.md) - Detailed installation and setup instructions
- [Usage Guide](docs/USAGE.md) - How to use the bot effectively and available configuration options
- [Development Guide](docs/DEVELOPMENT.md) - Information for developers
- [Roadmap](docs/ROADMAP.md) - Future development plans

## License

This project is licensed under the MIT License. See the [LICENSE.txt](LICENSE.txt) file for details.

## Contributing

Contributions are welcome! Please see our [Contributing Guide](docs/CONTRIBUTING.md) for details.

## Support

Support Bitcoin Calendar via [Coinos](https://coinos.io/)
- Bitcoin Calendar EN ⚡️ `bitcal@coinos.io` 🔗 [Coinos page](https://coinos.io/bitcal)

Support Bitcoin Calendar on [Geyser](https://geyser.fund/project/bitcoincalendar)

...or 👇

[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)