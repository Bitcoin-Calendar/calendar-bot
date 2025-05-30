package logging

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"calendar-bot/internal/config" // Requires config for log settings

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Setup configures the global zerolog logger based on the provided configuration.
func Setup(cfg *config.Config) {
	// Attempt to load .env file. Log if it's loaded or if it fails.
	if err := godotenv.Load(); err == nil {
		// Use a temporary basic logger if global one isn't set yet, or log after full setup.
		// For now, assuming log.Info() will work or queue if logger not fully ready.
		log.Info().Msg("Loaded configuration from .env file.")
	} else {
		log.Debug().Err(err).Msg("No .env file found or error loading it. Using environment variables.")
	}

	level, err := zerolog.ParseLevel(strings.ToLower(cfg.LogLevel))
	if err != nil {
		// If parsing fails, log a warning and default to InfoLevel.
		// This uses the potentially uninitialized global logger. This is a common pattern.
		log.Warn().Err(err).Str("configuredLogLevel", cfg.LogLevel).Msg("Invalid log level string. Defaulting to 'info'.")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var writers []io.Writer

	if cfg.ConsoleLog {
		// Use os.Stderr for console output, as is common for logs.
		// TimeFormat is set for better readability.
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	}

	if cfg.LogDir != "" {
		// Check if log directory exists, create if not.
		if _, statErr := os.Stat(cfg.LogDir); os.IsNotExist(statErr) {
			if mkdirErr := os.MkdirAll(cfg.LogDir, 0755); mkdirErr != nil {
				// If console logging isn't already on, add it as a fallback.
				if !cfg.ConsoleLog {
					writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
				}
			}
		}

		// If LogDir is set (and ideally created), add file logger.
		// Re-check cfg.LogDir in case MkdirAll failed but we continued (e.g. with console fallback)
		// However, the current structure ensures that if MkdirAll fails and console is not on, writers might be empty before this.
		// A check on `writers` length or a successful dir creation flag might be more robust.
		if _, statErr := os.Stat(cfg.LogDir); !os.IsNotExist(statErr) { // Check if dir exists now
			writers = append(writers, &lumberjack.Logger{
				Filename:   fmt.Sprintf("%s/calendar-bot.log", cfg.LogDir),
				MaxSize:    1, // megabytes
				MaxBackups: 3,
				MaxAge:     28, // days
				Compress:   true,
			})
			log.Info().Str("filePath", fmt.Sprintf("%s/calendar-bot.log", cfg.LogDir)).Msg("File logging enabled.")
		} else if cfg.LogDir != "" { // LogDir was specified but doesn't exist (e.g. creation failed)
			log.Warn().Str("path", cfg.LogDir).Msg("Log directory specified but could not be accessed/created. File logging disabled.")
		}

	} else {
		// If LogDir is explicitly empty, ensure console logging if not already added.
		if !cfg.ConsoleLog {
			writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
			log.Info().Msg("Log directory not specified. Logging to console.")
		}
	}

	// Fallback: if no writers are configured (e.g., console log false, and logdir empty or failed to create/access), default to stderr.
	if len(writers) == 0 {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		log.Warn().Msg("No log output explicitly configured (or file log failed). Defaulting to console (stderr).")
	}

	// Create the multi-level writer and set it as the global logger.
	// Add common fields like timestamp. Service name, version, hostname can be added here too if desired.
	multiWriter := io.MultiWriter(writers...)

	log.Logger = zerolog.New(multiWriter).With().Timestamp().Logger()

	hostname, errHost := os.Hostname()
	if errHost != nil {
		// Log the error using the just-configured logger if possible, or a direct fmt.Print if logger is bad.
		// For now, we'll assume log.Logger is good enough to report this critical failure.
		log.Error().Err(errHost).Msg("CRITICAL - Failed to get hostname. Using 'unknown-host'.")
		hostname = "unknown-host"
	}

	log.Logger = log.Logger.With().Str("service", "calendar-bot").Str("host", hostname).Logger()

	log.Info().Msg("Logger initialized successfully.")
	log.Debug().Msg("Debug logging is active.") // This will only print if level is debug
} 