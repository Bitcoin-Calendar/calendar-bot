# Contributing to Bitcoin Calendar

Thank you for considering contributing to the Bitcoin Calendar Bot project! This document outlines the process for contributing and provides guidelines to help you get started.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone. Please be considerate and respectful in your communications and actions.

## Ways to Contribute

There are many ways to contribute:

1.  **Code Contributions**: Implement new features, fix bugs in the bot, or improve existing functionality.
2.  **Documentation**: Improve or translate documentation for the bot.
3.  **Testing**: Test the bot in different environments, report issues, or suggest improvements to the testing process.
4.  **Bug Reports**: Report bugs or suggest improvements for the bot via GitHub Issues.
5.  **API Event Contributions**: If you want to add or correct historical Bitcoin events, please refer to the contribution guidelines for the Bitcoin Historical Events API repository (if separate and applicable), as this bot consumes data from that API.

## Getting Started

### Setting Up Your Development Environment

1.  Follow the setup instructions in the [Development Guide](DEVELOPMENT.md), which focuses on a Docker-based environment.
2.  Familiarize yourself with the project structure and code architecture as described in [PROJECT_STRUCTURE.md](PROJECT_STRUCTURE.md).

### Finding Issues to Work On

-   Check the [GitHub Issues](https://github.com/Bitcoin-Calendar/calendar-bot/issues) for open tasks related to the bot.
-   Look for issues tagged with "good-first-issue" if you're new to the project.
-   Review the [Roadmap](ROADMAP.md) for planned features for the bot.

## Contributing Process

### For Code Contributions to the Bot

1.  **Fork the Repository**: Create your own fork of the `Bitcoin-Calendar/calendar-bot` repository on GitHub.

2.  **Create a Feature Branch**:
    ```bash
    git checkout -b feature/your-feature-name
    ```

3.  **Make Your Changes**:
    -   Follow the coding standards and practices outlined in the [Development Guide](DEVELOPMENT.md).
    -   Keep changes focused on a single issue or feature for the bot.
    -   Write clear, descriptive commit messages.

4.  **Test Your Changes**:
    -   Ensure that your bot code changes work as expected, preferably using the Docker Compose test services (`nostr-bot-en-test`, `nostr-bot-ru-test`).
    -   Run any existing tests and add new ones if necessary for Go code.

5.  **Update Documentation**:
    -   Update the relevant documentation files in the `docs/` directory to reflect your changes to the bot's functionality or setup.
    -   Add comments to your code when necessary.

6.  **Commit and Push Your Changes**:
    ```bash
    git add .
    git commit -m "Add feature/fix: Description of changes to the bot"
    git push origin feature/your-feature-name
    ```

7.  **Create a Pull Request**:
    -   Go to your fork on GitHub and create a pull request to the main `Bitcoin-Calendar/calendar-bot` repository.
    -   Provide a clear title and description for your pull request.
    -   Reference any related issues using the # symbol (e.g., "Fixes #123").

8.  **Address Review Feedback**:
    -   Be responsive to feedback and make requested changes.
    -   Keep the discussion focused and professional.

### For Event Database Contributions

As this bot fetches events from an external API, contributions to the event data itself (adding new historical events, correcting existing ones, or adding translations for events) should be directed to the maintainers or repository of that specific Bitcoin Historical Events API. Please consult the documentation or contact the owners of that API for their contribution process.

If the API data is managed within this project's ecosystem but in a different component (e.g., a separate API server repository), contributions should be made there according to its guidelines.

## Guidelines for Quality Contributions (Bot Code)

### Code Quality

-   Follow Go best practices and coding standards.
-   Write clean, readable, and maintainable code.
-   Handle errors properly and gracefully.
-   Maintain backward compatibility when possible.
-   Add appropriate comments and documentation to the code.

### Documentation (Bot)

-   Write clear, concise, and accurate documentation for any changes to the bot's setup or operation.
-   Include examples when useful.
-   Keep documentation up-to-date with code changes.

### Testing (Bot)

-   Test new features or bug fixes, ideally using the Docker test services.
-   Ensure all Go tests pass before submitting pull requests if applicable.
-   Consider edge cases and error conditions.

### Commit Messages

-   Write descriptive commit messages that explain what and why (not how).
-   Use the present tense ("Add feature" not "Added feature").
-   Reference issues where appropriate.

## Pull Request Review Process

1.  **Initial Review**: Maintainers will review your pull request for quality and relevance to the bot project.
2.  **Feedback**: You may receive feedback or requests for changes.
3.  **Iteration**: Address feedback and make necessary changes.
4.  **Acceptance**: Once approved, your pull request will be merged into the bot repository.

## Adding New Languages (for Bot Posts)

Adding support for the bot to post in a new language (e.g., Spanish) primarily depends on the **API providing event data in that language**.

If the API supports a new language (e.g., via a `lang=es` parameter):

1.  **Update `docker-compose.yml`**: 
    *   Add new services for the new language, similar to `nostr-bot-en` and `nostr-bot-ru` (and their test counterparts).
    *   Example for Spanish (`es`):
        ```yaml
        # In docker-compose.yml
        nostr-bot-es:
          # ... (copy structure from nostr-bot-en)
          container_name: nostr-bot-es
          command: ["./nostr_bot", "NOSTR_PRIVATE_KEY_ES"]
          environment:
            - BOT_PROCESSING_LANGUAGE=es # Key change here
            # ... other env vars like LOG_LEVEL
        
        nostr-bot-es-test:
          # ... (copy structure from nostr-bot-en-test)
          container_name: nostr-bot-es-test
          command: ["./nostr_bot", "NOSTR_PRIVATE_KEY_EST"]
          environment:
            - BOT_PROCESSING_LANGUAGE=es # Key change here
            # ... other env vars like LOG_LEVEL, CONSOLE_LOG
        ```
2.  **Update `.env-example` and `.env`**: 
    *   Add new Nostr private key variables (e.g., `NOSTR_PRIVATE_KEY_ES`, `NOSTR_PRIVATE_KEY_EST`).
3.  **Update Documentation**: 
    *   Mention the new language support in `README.md` and `docs/USAGE.md`.
    *   Explain how to configure and run the bot for the new language.
4.  **Testing**: 
    *   Thoroughly test the new language service (e.g., `docker-compose run --rm nostr-bot-es-test`) to ensure it correctly fetches and posts events in the new language.

The bot's Go code itself (`main.go`) is already designed to handle different languages by passing the `botLanguage` variable to the API request. The main work is configuration and ensuring the API supports the new language.

## Attribution

Contributors will be acknowledged in:
-   The project's GitHub contributors list.
-   Release notes when applicable.
-   Special recognition for significant contributions.

## Questions?

If you have questions or need help regarding contributions to the bot itself, you can:
-   Open an issue on GitHub with your question.
-   Contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) via Nostr.

Thank you for contributing to the Bitcoin Calendar project! 