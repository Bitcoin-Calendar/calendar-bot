[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)

# Bitcoin Calendar Bot

### Archiving and relaying every Bitcoin milestone 

![bitcal-logo-wide](https://haven.bitcoin-calendar.org/9db5c5d32dff9f016bda92280deb4e14e137856704499827f5f0e6d83d7cb326.webp)

## Overview

The Bitcoin Calendar Bot is a Go-based application that reads events from CSV files and publishes them to Nostr relays. This bot automates the posting of calendar events about Bitcoin history. The bot supports English and Russian versions.

## Quick Start with Docker (Recommended)

This is the recommended way to run the Bitcoin Calendar Bot. It uses Docker and Docker Compose to manage the bot instances for different languages.

1.  **Prerequisites**:
    *   Ensure you have [Docker](https://docs.docker.com/get-docker/) and [Docker Compose](https://docs.docker.com/compose/install/) installed.

2.  **Clone the repository**:
    ```bash
    git clone https://github.com/Bitcoin-Calendar/calendar-bot.git
    cd calendar-bot
    ```

3.  **Set up environment variables**:
    *   Copy the example environment file:
        ```bash
        cp .env-example .env
        ```
    *   Edit the `.env` file and add your Nostr private keys for production and testing:
        ```env
        # Production Keys
        NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex"
        NOSTR_PRIVATE_KEY_RU="your_russian_specific_private_key_hex"

        # Test Keys
        NOSTR_PRIVATE_KEY_ENT="your_english_TEST_private_key_hex"
        NOSTR_PRIVATE_KEY_RUT="your_russian_TEST_private_key_hex"
        ```
    *   You can also adjust `LOG_LEVEL`, `CONSOLE_LOG`, etc., in the `.env` file to override defaults set in `docker-compose.yml`.

4.  **Build the Docker image**:
    You need to build the image that all services will use. You can do this by building one of the services (it will build the shared image):
    ```bash
    docker-compose build nostr-bot-en 
    # Or: docker-compose build (to build images for all services defined)
    ```

5.  **Running Test Instances Manually**:
    To test your setup without posting to your main Nostr accounts, you can run the test services defined in `docker-compose.yml` (e.g., `nostr-bot-en-test`). These use specific test keys from your `.env` file. For detailed commands and information on test instance behavior (logging, shared files), please see the [Running Test Instances section in our Usage Guide](docs/USAGE.md#running-test-instances).

6.  **Set up Cron Jobs for Production Bots**:
    The production bots (`nostr-bot-en`, `nostr-bot-ru`) are designed for scheduled runs using your system's cron. For detailed instructions on setting up cron jobs with `docker-compose`, including how Docker handles environment variables from the `.env` file for these services, please refer to the [Automated Operation section in our Usage Guide](docs/USAGE.md#automated-operation).

7.  **Viewing Logs**:
    The bot generates log files for its operations. For test runs, output is also sent to the console by default. Detailed information on log file locations, naming, and rotation can be found in the [Log Files section of our Usage Guide](docs/USAGE.md#log-files).

8.  **Building after code changes**:
    If you update the bot's Go code, you need to rebuild the Docker image before your cron jobs or manual test runs pick up the changes:
    ```bash
    cd /path/to/your/calendar-bot 
    docker-compose build
    ```

The `docker-compose.yml` file defines production services (`nostr-bot-en`, `nostr-bot-ru`) for scheduled runs and test services (`nostr-bot-en-test`, `nostr-bot-ru-test`) for manual testing. All services use the same built image but different command parameters and environment variable expectations for keys.

<details>
<summary>Legacy Manual Setup (Without Docker)</summary>

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
   NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex"
   NOSTR_PRIVATE_KEY_RU="your_russian_specific_private_key_hex"
   # Ensure these variable names (NOSTR_PRIVATE_KEY_EN, NOSTR_PRIVATE_KEY_RU)
   # match what you use in the run command.
   ```
   Or copy the `.env-example` file and adjust its contents as needed. (Note: The `.env-example` has been updated for Docker usage, so ensure your manual setup aligns with how you use the private key variables).
   For the manual run commands below, the bot expects the *name* of the environment variable (e.g., `NOSTR_PRIVATE_KEY_EN`) as a command-line argument. Ensure this variable, containing your actual private key, is set in your shell's environment *before* running the command. You can often achieve this by sourcing the `.env` file (e.g., `set -a; source .env; set +a` for bash/zsh if your `.env` file uses `KEY=VALUE` format) or by other means appropriate for your shell.

4. **Run the bot**
   ```bash
   # For English events
   LOG_DIR=./logs ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN

   # For Russian events
   LOG_DIR=./logs ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
   ```

### Automated Daily Posting with Cron

To ensure the bot posts events daily, you can set up a cron job on your host system. The bot is designed to run, process events for the current day, and then exit.
Cron jobs run with a minimal environment, so you must ensure that necessary environment variables (like `NOSTR_PRIVATE_KEY_EN` and `NOSTR_PRIVATE_KEY_RU`) are available within the cron execution scope. This might involve defining them directly in your crontab, sourcing a script that exports them, or other methods depending on your system.

Example cron entries:
```cron
# Ensure the path to your project directory is correct.
# Run Russian bot at 04:00 AM daily
00 04 * * * cd /path/to/your/calendar-bot && LOG_DIR=./logs ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU

# Run English bot at 15:00 PM daily
00 15 * * * cd /path/to/your/calendar-bot && LOG_DIR=./logs ./nostr_bot events_en.csv NOSTR_PRIVATE_KEY_EN
```
Find more details and examples of how to setup a cron job, including managing environment variables, in our Usage Guide.

</details>

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