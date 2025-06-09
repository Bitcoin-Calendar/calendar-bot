package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	APIEndpoint         string
	APIKey              string
	PrivateKey          string
	ProcessingLanguage  string
	LogDir              string
	LogLevel            string
	ConsoleLog          bool
	Debug               bool
	NostrRelays         []string
	EnvVarForPrivateKey string // To store the name of the env var holding the private key
}

// Validate checks the configuration for any errors.
// Placeholder for now.
func (c *Config) Validate() error {
	// TODO: Implement validation logic
	if c.APIEndpoint == "" {
		return fmt.Errorf("APIEndpoint is required")
	}
	if c.APIKey == "" {
		return fmt.Errorf("APIKey is required")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("PrivateKey is required")
	}
	if c.ProcessingLanguage == "" {
		return fmt.Errorf("ProcessingLanguage is required")
	}
	if c.ProcessingLanguage != "en" {
		return fmt.Errorf("Invalid BOT_PROCESSING_LANGUAGE '%s'. Must be 'en'", c.ProcessingLanguage)
	}
	if len(c.NostrRelays) == 0 {
		return fmt.Errorf("NostrRelays are required")
	}
	// Add more validation as needed
	return nil
}

// loadConfig loads configuration from environment variables and command-line arguments.
func LoadConfig(envVarForPrivateKeyName string) (*Config, error) {
	cfg := &Config{
		EnvVarForPrivateKey: envVarForPrivateKeyName,
	}

	// Attempt to load .env file, but don't make it fatal if it doesn't exist
	// This will be logged later if it fails, after logger is set up.
	_ = godotenv.Load() // Error is handled by a log message later if logger is set up.

	cfg.APIEndpoint = os.Getenv("BOT_API_ENDPOINT")
	// No need to check for empty here, Validate() will do it.

	cfg.APIKey = os.Getenv("BOT_API_KEY")
	// No need to check for empty here, Validate() will do it.

	cfg.PrivateKey = os.Getenv(cfg.EnvVarForPrivateKey)
	// No need to check for empty here, Validate() will do it.

	cfg.ProcessingLanguage = os.Getenv("BOT_PROCESSING_LANGUAGE")
	// No need to check for empty or invalid value here, Validate() will do it.

	cfg.LogDir = os.Getenv("BOT_LOG_DIR")
	if cfg.LogDir == "" {
		cfg.LogDir = "logs" // Default log directory
	}

	cfg.LogLevel = os.Getenv("BOT_LOG_LEVEL")
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info" // Default log level
	}

	consoleLog := os.Getenv("BOT_CONSOLE_LOG")
	if consoleLog == "true" {
		cfg.ConsoleLog = true
	}

	debugEnv := os.Getenv("BOT_DEBUG")
	if debugEnv == "true" {
		cfg.Debug = true
	}
	
	nostrRelaysEnv := os.Getenv("NOSTR_RELAYS")
	if nostrRelaysEnv != "" {
		// Split by comma again
		cfg.NostrRelays = strings.Split(nostrRelaysEnv, ",") 
		// Trim whitespace from each relay URL, and filter out empty strings if any result from split
		validRelays := make([]string, 0, len(cfg.NostrRelays))
		for _, relay := range cfg.NostrRelays {
			trimmedRelay := strings.TrimSpace(relay)
			if trimmedRelay != "" {
				validRelays = append(validRelays, trimmedRelay)
			}
		}
		cfg.NostrRelays = validRelays
	}
	// No need to check for empty here, Validate() will do it.


	// Perform validation after loading all values
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
} 