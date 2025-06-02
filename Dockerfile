# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Explicitly set GOROOT and GOPATH (Alpine Go image should do this, but for safety)
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH

WORKDIR /app

# Copy all source code, including go.mod, go.sum, main.go, and internal/
# This ensures that the module structure is fully present before any go commands are run in this context.
COPY . /app/

# Now, all go commands will run from the root of the copied module.
# Go should correctly interpret this as the module 'calendar-bot'
# because go.mod is at /app/go.mod.

RUN go mod download
RUN go mod tidy

# Explicitly set Go module mode (should be default but good for safety)
ENV GO111MODULE=on

# Before building, let's list the contents of /app and /app/internal to be sure
RUN ls -R /app

# Build the application. Go should resolve "calendar-bot/internal/metrics"
# relative to the module root at /app.
RUN go build -v -o /app/nostr_bot .
# Added -v for verbose output from the build command

# Stage 2: Create the final lightweight image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Install strace
RUN apk add --no-cache strace

# Copy the built binary from the builder stage
COPY --from=builder /app/nostr_bot /app/nostr_bot

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
CMD ["./nostr_bot", "NOSTR_PRIVATE_KEY_ENT"] 