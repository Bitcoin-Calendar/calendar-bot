# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
# The output will be in /app/nostr_bot
RUN go build -o /app/nostr_bot main.go metrics.go

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy the built binary from the builder stage
COPY --from=builder /app/nostr_bot /app/nostr_bot

# Copy the CSV event files
COPY events_en.csv /app/events_en.csv
COPY events_ru.csv /app/events_ru.csv
# Add other CSV files if needed, e.g.:
# COPY events_es.csv /app/events_es.csv

# Create an empty .env file.
# Your application tries to load .env but can also work with environment variables
# passed directly by Docker. This file prevents an error if it's not found,
# allowing Docker-provided environment variables to be used.
RUN touch /app/.env

# Create directories for logs and metrics and set ownership
RUN mkdir -p /app/logs && chown appuser:appgroup /app/logs
RUN mkdir -p /app/metrics && chown appuser:appgroup /app/metrics

# Switch to the non-root user
USER appuser

# Set default environment variables.
# These can be overridden at runtime (docker run -e ...).
ENV LOG_DIR="/app/logs"
ENV LOG_LEVEL="info"
# Example: you will need to pass the actual private key when running the container
ENV NOSTR_PRIVATE_KEY_EN="your_english_specific_private_key_hex_goes_here_at_runtime"

# Default command to run the application (e.g., for English events).
# You can override this when running the container for different languages/keys.
# For example, for Russian:
# docker run -e NOSTR_PRIVATE_KEY_RU="..." your_image_name ./nostr_bot events_ru.csv NOSTR_PRIVATE_KEY_RU
CMD ["./nostr_bot", "events_en.csv", "NOSTR_PRIVATE_KEY_ENT"] 