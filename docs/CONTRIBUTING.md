# Contributing to Bitcoin Calendar

Thank you for considering contributing to the Bitcoin Calendar! This document outlines the process for contributing to the project and provides guidelines to help you get started.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for everyone. Please be considerate and respectful in your communications and actions.

## Ways to Contribute

There are many ways to contribute to the Bitcoin Calendar:

1. **Adding Historical Events**: Expand the event database with accurate Bitcoin historical events
2. **Code Contributions**: Implement new features or fix bugs
3. **Documentation**: Improve or translate documentation
4. **Translations**: Help translate the bot to support more languages
5. **Testing**: Test the bot in different environments and report issues
6. **Bug Reports**: Report bugs or suggest improvements

## Getting Started

### Setting Up Your Development Environment

1. Follow the setup instructions in the [Development Guide](DEVELOPMENT.md)
2. Familiarize yourself with the project structure and code architecture

### Finding Issues to Work On

- Check the [GitHub Issues](https://github.com/bitcoin21ideas/nostr-calendar-bot/issues) for open tasks
- Look for issues tagged with "good-first-issue" if you're new to the project
- Review the [Roadmap](ROADMAP.md) for planned features

## Contributing Process

### For Code Contributions

1. **Fork the Repository**:
   - Create your own fork of the repository on GitHub

2. **Create a Feature Branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Your Changes**:
   - Follow the coding standards and practices outlined in the [Development Guide](DEVELOPMENT.md)
   - Keep changes focused on a single issue or feature
   - Write clear, descriptive commit messages

4. **Test Your Changes**:
   - Ensure that your changes work as expected
   - Run any existing tests and add new ones if necessary

5. **Update Documentation**:
   - Update the relevant documentation to reflect your changes
   - Add comments to your code when necessary

6. **Commit and Push Your Changes**:
   ```bash
   git add .
   git commit -m "Add feature/fix: Description of changes"
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request**:
   - Go to your fork on GitHub and create a pull request to the main repository
   - Provide a clear title and description for your pull request
   - Reference any related issues using the # symbol (e.g., "Fixes #123")

8. **Address Review Feedback**:
   - Be responsive to feedback and make requested changes
   - Keep the discussion focused and professional

### For Event Database Contributions

1. **Review Existing Events**:
   - Check the CSV files to understand the format and existing events
   - Follow the guidelines in the [CSV Format Guide](CSV_FORMAT.md)

2. **Research Historical Events**:
   - Ensure historical accuracy of dates and details
   - Include reliable sources for verification
   - Follow the established format for event descriptions

3. **Submit Your Contributions**:
   - Follow the same fork and pull request process outlined above
   - Include links to sources in your pull request description

## Guidelines for Quality Contributions

### Code Quality

- Follow Go best practices and coding standards
- Write clean, readable, and maintainable code
- Handle errors properly and gracefully
- Maintain backward compatibility when possible
- Add appropriate comments and documentation

### Documentation

- Write clear, concise, and accurate documentation
- Include examples when useful
- Keep documentation up-to-date with code changes

### Testing

- Write tests for new features or bug fixes
- Ensure all tests pass before submitting pull requests
- Consider edge cases and error conditions

### Commit Messages

- Write descriptive commit messages that explain what and why (not how)
- Use the present tense ("Add feature" not "Added feature")
- Reference issues where appropriate

## Pull Request Review Process

1. **Initial Review**: Maintainers will review your pull request for quality and relevance
2. **Feedback**: You may receive feedback or requests for changes
3. **Iteration**: Address feedback and make necessary changes
4. **Acceptance**: Once approved, your pull request will be merged

## Adding New Languages

If you want to add support for a new language:

1. Create a new CSV file for events in that language (e.g., `events_es.csv` for Spanish)
2. Follow the format guidelines in the [CSV Format Guide](CSV_FORMAT.md)
3. Update documentation to mention the new language support
4. Test the bot with the new language file

## Attribution

Contributors will be acknowledged in:
- The project's GitHub contributors list
- Release notes when applicable
- Special recognition for significant contributions

## Questions?

If you have questions or need help, you can:
- Open an issue on GitHub with your question
- Contact [Tony](https://njump.me/npub10awzknjg5r5lajnr53438ndcyjylgqsrnrtq5grs495v42qc6awsj45ys7) via Nostr.

Thank you for contributing to the Bitcoin Calendar project! 