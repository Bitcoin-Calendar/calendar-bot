# Bitcoin Calendar - Detailed Roadmap

This document outlines the development roadmap for the Bitcoin Calendar project, organized by quarter with specific goals and milestones.

## Q2 2025 (April-June)

`Version: 1.1.0`

### Database Enhancement
- **English Event Population**
  - Target: Expand from current March & April events to 500+ events across the year (prepare English version based on the extensive Russian version)
  - Add missing historical milestones
  - Standardize formatting and citation requirements

- **Content Quality Improvements**
  - Add media files for all events (images, charts)
  - Include additional backup archive links to prevent link rot
  - Implement consistent formatting standards

### Infrastructure Improvements
- **Monitoring and Alerting**
  - Implement health check system
  - Set up automated alerts for failed posts

- **Dockerization of Calendar Bot**
  - [x] Create Dockerfile for the Go application
  - [x] Implement Docker Compose setup for easy deployment and management (production and test instances)
  - [x] Updated documentation for Docker-based setup

### Data Redundancy Improvements
- **Deeper Nostr Integration**
  - [x] Spin up Haven relay
  - [x] Spin up Blossom Server

## Q3 2025 (July-September)

`Version: 1.2.0`

### Content Enhancement
- **Metadata Optimization**
  - Implement consistent hashtag strategy
  - Create event categories and tagging system

- **Content Expansion**
  - Create thematic collections of related events
  - Develop "This Week in Bitcoin History" compilations and weekly statistics reports

## Q4 2025 (October-December)
`Version: 2.0.0`

### Web Presence Development
- **Interactive Website**
  - Develop responsive frontend for browsing events
  - Create search functionality
  - Add filtering by categories, tags, and dates

### API Development
- **Public API**
  - Create RESTful API for accessing events
  - Implement rate limiting and authentication
  - Develop comprehensive API documentation

- **Integration Options**
  - Implement webhooks for event notifications
  - Develop widgets for embedding calendar events
  - Create RSS feeds for various categories

## Q1 2026 (January-March)

`Version: 2.1.0`

### Mobile Experience
- **Progressive Web App**
  - Develop PWA for mobile access
  - Implement offline functionality
  - Create notification system for upcoming events
  - Implement timeline visualization

## Q2 2026 (April-June)

`Version: 2.2.0`

### Community Features
- **User Accounts**
  - Implement Nostr login for personalized experience
  - Create user profiles with contribution history
  - Develop preferences and notification settings

- **Contribution System**
  - Add event submission portal
  - Implement community voting on submissions
  - Create contributor recognition system

## Long-term Vision

### Bitcoin Archive
- **Ensuring no Bitcoin event is lost**
  - Setup ArchiveBox to locally host all the web pages and documents Bitcoin Calendar refers to
  - Integrate archived pages across the Bitcoin Calendar ecosystem – the bot, the website - to ensure the project does not have to rely on any third-party services 

### Educational Resources
- **Learning Paths**
  - Develop curated learning paths based on historical events and [21ideas](https://21ideas.org/en/) Bitcoin knowledge base
  - Create downloadable educational materials to complement calendar format

### Integration Ecosystem
- **Calendar Integration**
  - Calendar subscription service (iCal, Google Calendar)
  - Event reminder system
  - Personal calendar integration

### Advanced Features
- **Content Exploration**
  - Implement search by concepts and relationships
  - Develop visualization of historical connections
  - Create AI-assisted navigation of Bitcoin history

### Content System Evolution
- **Submission System**
  - Design curator workflow for approving submissions
  - Implement verification process for historical accuracy
  - Create citation standards and guidelines


[![⚡️zapmeacoffee](https://img.shields.io/badge/⚡️zap_-me_a_coffee-violet?style=plastic)](https://zapmeacoffee.com/npub1tcalvjvswjh5rwhr3gywmfjzghthexjpddzvlxre9wxfqz4euqys0309hn)
