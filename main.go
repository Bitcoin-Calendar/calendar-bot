package main

import (
    "context"
	"encoding/json"
    "fmt"
	"io/ioutil"
	"net/http"
    "os"
    "runtime"
    "strings"
    "time"

	"github.com/joho/godotenv"
    "github.com/nbd-wtf/go-nostr"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gopkg.in/natefinch/lumberjack.v2"
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

// APIEvent struct to match the expected data structure from the API
type APIEvent struct {
	ID          uint      `json:"ID"`
	Date        time.Time `json:"Date"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Tags        string    `json:"Tags"`       // JSON array string
	Media       string    `json:"Media"`      // URL
	References  string    `json:"References"` // JSON array string
}

// getCurrentDirectory gets the current working directory
func getCurrentDirectory() string {
    dir, err := os.Getwd()
    if err != nil {
		// If there's an error, log it and return a placeholder
		// Using the global logger might not be initialized yet if this is called very early
		// For now, assume logger is available or this function is called after setup
		log.Error().Err(err).Msg("Failed to get current working directory")
		return "unknown_dir"
    }
    return dir
}

// logEnvironmentVariables logs non-sensitive environment variables
func logEnvironmentVariables() {
    envVars := make(map[string]string)
    for _, env := range os.Environ() {
        parts := strings.SplitN(env, "=", 2)
        if len(parts) == 2 {
            // Skip sensitive environment variables
            if !strings.Contains(strings.ToLower(parts[0]), "key") && 
               !strings.Contains(strings.ToLower(parts[0]), "secret") &&
               !strings.Contains(strings.ToLower(parts[0]), "password") &&
               !strings.Contains(strings.ToLower(parts[0]), "token") {
                envVars[parts[0]] = parts[1]
            }
        }
    }
    log.Debug().Interface("environment", envVars).Msg("Environment variables")
}

func main() {
	// Expect only <env_var_for_private_key> as a command-line argument
	if len(os.Args) < 2 {
		// Before logger is initialized, print to stderr
		fmt.Fprintln(os.Stderr, "Usage: calendar-bot <env_var_for_private_key>")
		os.Exit(1)
	}

	envVarForPrivateKeyName := os.Args[1]

	// Load configuration
	cfg, err := loadConfig(envVarForPrivateKeyName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Setup Logging using the configuration
	// Note: setupLogger itself will handle .env loading message if it uses the global logger
	// which gets configured by setupLogger.
	setupLogger(cfg)

	// Log the processing language after logger is configured
	log.Info().Str("language", cfg.ProcessingLanguage).Msg("Bot configured to process events for language.")

	if cfg.Debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel) // Ensure debug level is set if cfg.Debug is true
		log.Debug().Msg("Debug logging enabled.")
        log.Debug().
            Str("os", runtime.GOOS).
            Str("arch", runtime.GOARCH).
            Str("goVersion", runtime.Version()).
            Int("cpus", runtime.NumCPU()).
            Str("workingDir", getCurrentDirectory()).
            Str("envVarForPrivateKey", cfg.EnvVarForPrivateKey). // Log the name of the env var used
			Msg("System information")
        logEnvironmentVariables() // This function could also take cfg if it needs specific values
	}

	today := time.Now().Format("01-02") // Format is "MM-DD"
	log.Info().Str("date", today).Msg("Starting bot execution. Fetching events from API.")

    metrics := NewMetricsCollector()
    
	// Extract month and day from today's date string
	parts := strings.Split(today, "-")
	var currentMonth, currentDay string
	if len(parts) == 2 {
		currentMonth = parts[0]
		currentDay = parts[1]
		log.Debug().Str("month", currentMonth).Str("day", currentDay).Msg("Extracted month and day for API query")
    } else {
		log.Error().Str("todayFormat", today).Msg("Failed to parse month and day from today's date format. Cannot query API by date.")
		// Decide if this is fatal or if it should attempt to fetch all events as a fallback
		// For now, exiting as specific date filtering is the new primary logic.
		os.Exit(1)
	}

	apiEvents, err := fetchEventsFromAPI(cfg.APIEndpoint, cfg.APIKey, currentMonth, currentDay, cfg.ProcessingLanguage)
	if err != nil {
		log.Error().Err(err).Msg("Fatal: Failed to fetch events from API. Bot will exit.")
		metrics.LogSummary() // Log whatever metrics might have been gathered (likely none in this path)
        // Export metrics even on critical failure before exiting
        metricsDir := "metrics"
        if mkDirErr := os.MkdirAll(metricsDir, 0755); mkDirErr != nil {
            log.Error().Err(mkDirErr).Str("directory", metricsDir).Msg("Failed to create metrics directory for error export")
        }
        metricsFilePath := fmt.Sprintf("%s/metrics_error_%s.json", metricsDir, time.Now().Format("2006-01-02_15-04-05"))
        if exportErr := metrics.ExportMetrics(metricsFilePath); exportErr != nil {
            log.Error().Err(exportErr).Str("file", metricsFilePath).Msg("Failed to export metrics during error shutdown")
    } else {
            log.Info().Str("file", metricsFilePath).Msg("Metrics exported successfully during error shutdown")
    }
		os.Exit(1) // Exit if API fetch fails critically
	}
	log.Info().Int("eventsFetchedCount", len(apiEvents)).Msg("Successfully fetched events from API.")

    // Relays are now from cfg.NostrRelays
	log.Debug().Strs("relays", cfg.NostrRelays).Msg("Configured relays for publishing from NOSTR_RELAYS env var")

	eventsToPublishToday := 0
	for _, event := range apiEvents {
		if event.Date.Format("01-02") == today {
			eventsToPublishToday++
			log.Info().Str("eventTitle", event.Title).Str("eventAPIDate", event.Date.Format("2006-01-02")).Uint("eventID", event.ID).Msg("Processing matching API event for today")

			var currentEventAPITags []string
			if event.Tags != "" && event.Tags != "[]" { // Check for empty or empty JSON array string
				if err := json.Unmarshal([]byte(event.Tags), &currentEventAPITags); err != nil {
					log.Warn().Err(err).Uint("eventID", event.ID).Str("tagsString", event.Tags).Msg("Failed to unmarshal event Tags. Proceeding with no API tags for this event.")
					// Proceed with empty currentEventAPITags
				}
			}

			var currentEventAPIReferences []string
			if event.References != "" && event.References != "[]" { // Check for empty or empty JSON array string
				if err := json.Unmarshal([]byte(event.References), &currentEventAPIReferences); err != nil {
					log.Warn().Err(err).Uint("eventID", event.ID).Str("referencesString", event.References).Msg("Failed to unmarshal event References. Proceeding with no API references for this event.")
					// Proceed with empty currentEventAPIReferences
				}
			}

			requestID := fmt.Sprintf("api-event-%d-%s-%d", event.ID, today, time.Now().UnixNano())
			eventSpecificLogger := log.With().Str("requestID", requestID).Uint("eventID", event.ID).Logger()

			publishEvent(
				cfg.PrivateKey,
				event.Date.Format("2006-01-02"),
				event.Title,
				event.Description,
				currentEventAPITags,
				event.Media,
				currentEventAPIReferences,
				metrics,
				cfg.NostrRelays, // Use relays from config
				eventSpecificLogger,
			)
			// The 30-minute wait is inside publishEvent, triggered after a successful multi-relay publish attempt.
		} else {
			metrics.EventsSkipped++
			log.Debug().Uint("eventID", event.ID).Str("eventTitle", event.Title).Str("eventAPIDate", event.Date.Format("2006-01-02")).Msg("Skipped API event: Date does not match today.")
		}
	}

	if eventsToPublishToday == 0 {
		log.Info().Msg("No events found for today's date after checking all fetched events.")
	}

	log.Info().Msg("Bot execution finished for today.")
	metrics.LogSummary()

	metricsDir := "metrics"
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		log.Error().Err(err).Str("directory", metricsDir).Msg("Failed to create metrics directory for final export")
	}
	metricsFilePath := fmt.Sprintf("%s/metrics_run_%s.json", metricsDir, time.Now().Format("2006-01-02_15-04-05"))
	if err := metrics.ExportMetrics(metricsFilePath); err != nil {
		log.Error().Err(err).Str("file", metricsFilePath).Msg("Failed to export metrics at end of run")
	} else {
		log.Info().Str("file", metricsFilePath).Msg("Metrics exported successfully at end of run")
	}
}

// loadConfig loads configuration from environment variables and command-line arguments.
func loadConfig(envVarForPrivateKeyName string) (*Config, error) {
	cfg := &Config{
		EnvVarForPrivateKey: envVarForPrivateKeyName,
	}

	// Attempt to load .env file, but don't make it fatal if it doesn't exist
	// This will be logged later if it fails, after logger is set up.
	_ = godotenv.Load() // Error is handled by a log message later if logger is set up.

	cfg.APIEndpoint = os.Getenv("BOT_API_ENDPOINT")
	if cfg.APIEndpoint == "" {
		return nil, fmt.Errorf("error: Environment variable BOT_API_ENDPOINT is not set")
	}

	cfg.APIKey = os.Getenv("BOT_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("error: Environment variable BOT_API_KEY is not set")
	}

	cfg.PrivateKey = os.Getenv(cfg.EnvVarForPrivateKey)
	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("error: Environment variable %s (for private key) is not set", cfg.EnvVarForPrivateKey)
	}

	cfg.ProcessingLanguage = os.Getenv("BOT_PROCESSING_LANGUAGE")
	if cfg.ProcessingLanguage == "" {
		return nil, fmt.Errorf("error: Environment variable BOT_PROCESSING_LANGUAGE is not set. Must be 'en' or 'ru'")
	}
	if cfg.ProcessingLanguage != "en" && cfg.ProcessingLanguage != "ru" {
		return nil, fmt.Errorf("error: Invalid BOT_PROCESSING_LANGUAGE '%s'. Must be 'en' or 'ru'", cfg.ProcessingLanguage)
	}

	cfg.LogDir = os.Getenv("LOG_DIR") // Optional, defaults to current directory if empty

	cfg.LogLevel = strings.ToLower(os.Getenv("LOG_LEVEL"))
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info" // Default log level
	}

	cfg.ConsoleLog = os.Getenv("CONSOLE_LOG") == "true"
	cfg.Debug = os.Getenv("DEBUG") == "true"

	relayEnvVar := os.Getenv("NOSTR_RELAYS")
	if relayEnvVar == "" {
		return nil, fmt.Errorf("error: Environment variable NOSTR_RELAYS is not set. Please provide a comma-separated list of relay URLs")
	}
	relays := strings.Split(relayEnvVar, ",")
	parsedRelays := make([]string, 0, len(relays))
	for _, r := range relays {
		trimmed := strings.TrimSpace(r)
		if trimmed != "" {
			parsedRelays = append(parsedRelays, trimmed)
		}
	}
	if len(parsedRelays) == 0 {
		return nil, fmt.Errorf("error: NOSTR_RELAYS environment variable is set but contains no valid relay URLs after parsing")
	}
	cfg.NostrRelays = parsedRelays

	return cfg, nil
}

func fetchEventsFromAPI(apiBaseURL string, apiKey string, month string, day string, language string) ([]APIEvent, error) {
	// Construct URL with month, day, and language query parameters
	requestURL := fmt.Sprintf("%s/events?month=%s&day=%s&lang=%s", apiBaseURL, month, day, language)
	log.Debug().Str("url", requestURL).Msg("Fetching events from API for specific month, day, and language")

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create API request: %w", err)
	}

	req.Header.Set("X-API-KEY", apiKey)
	req.Header.Set("Accept", "application/json") // Good practice to set Accept header

	client := &http.Client{Timeout: time.Second * 30} // 30-second timeout for the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute API request to %s: %w", requestURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := ioutil.ReadAll(resp.Body) // Use ioutil for older Go versions, or io.ReadAll for Go 1.16+
		if readErr != nil {
			return nil, fmt.Errorf("API request to %s failed with status %s, and failed to read error response body: %w", requestURL, resp.Status, readErr)
		}
		return nil, fmt.Errorf("API request to %s failed with status %s: %s", requestURL, resp.Status, string(bodyBytes))
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response body from %s: %w", requestURL, err)
	}

	// Define struct to match the *full* API response, including pagination if present
	var apiResponse struct {
		Events     []APIEvent `json:"events"`
		Pagination struct {
			Total       int `json:"total"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			LastPage    int `json:"last_page"`
		} `json:"pagination"` // Assuming pagination structure, adjust if different
	}

	if err := json.Unmarshal(bodyBytes, &apiResponse); err != nil {
		// Log the body for debugging if unmarshal fails
		log.Error().Err(err).Str("responseBody", string(bodyBytes)).Msg("Failed to unmarshal API response JSON")
		return nil, fmt.Errorf("failed to unmarshal API response from %s: %w. Body: %s", requestURL, err, string(bodyBytes))
	}

	log.Debug().Int("eventsReceived", len(apiResponse.Events)).Interface("paginationDetails", apiResponse.Pagination).Msg("Successfully unmarshalled API response")
	return apiResponse.Events, nil
}

func publishEvent(sk string, eventDateStr string, title string, description string, apiTags []string, mediaURL string, apiReferences []string, metrics *MetricsCollector, relays []string, eventLogger zerolog.Logger) {
	eventLogger.Info().Msg("Preparing to publish event")

	messageParts := []string{title, "\n", description} // Start with title and description

	if mediaURL != "" {
		messageParts = append(messageParts, fmt.Sprintf("\n\n%s", mediaURL))
	}
	if len(apiReferences) > 0 {
		messageParts = append(messageParts, "\n") // Keep a newline for separation if there was no media, or add to existing newlines
		for _, ref := range apiReferences {
			// Each reference on a new line, ensure no double newlines if mediaURL was also present and added \n\n
			// If mediaURL was empty, the first reference will follow the \n added above.
			// If mediaURL was present, it added \n\n, so this \n should ensure references start clearly.
			// This might result in triple newline if description, media, and references are all present.
			// Let's refine the newline logic carefully.

			// New logic: Add a newline separator only if there's something before it (description or media)
			// and then add each reference.
			// This will be handled by the Join and ReplaceAll logic later more generally.
			messageParts = append(messageParts, fmt.Sprintf("\n- %s", ref))
		}
	}
	message := strings.Join(messageParts, "") // Join parts, newlines are already included or handled by description
    // Ensure double newlines between sections if not already present
    message = strings.ReplaceAll(message, "\n\n\n", "\n\n") // Consolidate triple newlines to double
    message = strings.ReplaceAll(message, "\n \n", "\n\n") // Clean up potential space between newlines
    
    // Trim leading/trailing newlines that might result from conditional appends
    message = strings.TrimSpace(message)
    // After trimming, if the message was only newlines, it will be empty.
    // We need to re-ensure the core structure: title

//description then optional sections.

    // Rebuild message more deterministically for newlines
    var finalMessageBuilder strings.Builder
    finalMessageBuilder.WriteString(title)
    finalMessageBuilder.WriteString("\n\n") // Ensure separation after title
    finalMessageBuilder.WriteString(description)

    if mediaURL != "" {
        finalMessageBuilder.WriteString("\n\n")
        finalMessageBuilder.WriteString(mediaURL)
    }

    if len(apiReferences) > 0 {
        finalMessageBuilder.WriteString("\n\n") // Separator for references section
        for i, ref := range apiReferences {
            if i > 0 {
                finalMessageBuilder.WriteString("\n") // Newline between references
            }
            finalMessageBuilder.WriteString(fmt.Sprintf("- %s", ref))
        }
    }
    message = finalMessageBuilder.String()


	// Nostr tags
	defaultTags := []string{"bitcoin", "history", "onthisday", "calendar", "btc"} // Example default tags
	allEventTags := nostr.Tags{}
	for _, t := range defaultTags {
		allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(t)})
	}
	for _, apiTag := range apiTags {
		if apiTag != "" { // Avoid empty tags
			allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(apiTag)})
		}
	}
	allEventTags = append(allEventTags, nostr.Tag{"d", eventDateStr}) // Date of the event, for identification

	ev := nostr.Event{
                CreatedAt: nostr.Now(),
                Kind:      nostr.KindTextNote,
		Tags:      allEventTags,
                Content:   message,
            }

	if err := ev.Sign(sk); err != nil {
		eventLogger.Error().Err(err).Msg("Failed to sign Nostr event")
                metrics.EventsFailed++
		return
            }

	eventLogger.Debug().Str("eventID", ev.ID).Str("pubkey", ev.PubKey).Int("tagCount", len(ev.Tags)).Msg("Event signed and ready for publishing")

	successfulRelayPublishes := 0
            for _, relayURL := range relays {
		relayLog := eventLogger.With().Str("relayURL", relayURL).Logger()
                relayLog.Debug().Msg("Attempting to connect to relay")
                
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // Shortened timeout for connect
		relayConn, err := nostr.RelayConnect(ctx, relayURL)
                if err != nil {
			relayLog.Warn().Err(err).Msg("Failed to connect to relay")
                    metrics.RecordRelayFailure(relayURL)
			cancel() // Important to cancel context on error
                    continue
                }
		// Ensure cancel is called regardless of what happens next in this iteration
        // defer cancel() // This would be for the function scope, need per-loop management

		relayLog.Debug().Msg("Successfully connected to relay. Preparing to publish.")

		publishCtx, publishCancel := context.WithTimeout(context.Background(), 25*time.Second) // Timeout for publish operation
		err = relayConn.Publish(publishCtx, ev) // Corrected: Assume Publish returns a single error value based on compiler message
        

		if err != nil {
			relayLog.Warn().Err(err).Msg("Failed to publish event to relay") // Corrected: Removed status as it's not available with single error return
                        metrics.RecordRelayFailure(relayURL)
                    } else {
			// If Publish returns only error, a nil error implies success.
			// The switch statement for publishStatus is removed as the constants are undefined.
			relayLog.Info().Msg("Event successfully published to relay (inferred from nil error)")
			metrics.RecordRelaySuccess(relayURL, 0) // Duration can be added if measured
			successfulRelayPublishes++
                    }
        publishCancel() // Cancel the publish context
        relayConn.Close() // Close the connection to the relay
        cancel() // Cancel the connection context
            }

	if successfulRelayPublishes > 0 {
		eventLogger.Info().Int("successfulRelaysCount", successfulRelayPublishes).Int("totalRelaysAttempted", len(relays)).Msg("Event publishing process completed for one or more relays.")
                metrics.EventsPosted++
		log.Info().Int("waitMinutes", 30).Msg("Waiting 30 minutes after successful publish before processing next event for today (if any)...")
                time.Sleep(30 * time.Minute)
            } else {
		eventLogger.Warn().Msg("Event was not successfully published to any of the configured relays.")
                metrics.EventsFailed++
            }
}

// MetricsCollector stores metrics about bot operation
type MetricsCollector struct {
	EventsPosted      int                    `json:"eventsPosted"`
	EventsSkipped     int                    `json:"eventsSkipped"`
	EventsFailed      int                    `json:"eventsFailed"`
	RelaySuccesses    map[string]int         `json:"relaySuccesses"`
	RelayFailures     map[string]int         `json:"relayFailures"`
	RelaySuccessTimes map[string][]time.Duration `json:"relaySuccessTimes,omitempty"` // omitempty if not always populated
}

// NewMetricsCollector initializes a new MetricsCollector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		RelaySuccesses:    make(map[string]int),
		RelayFailures:     make(map[string]int),
		RelaySuccessTimes: make(map[string][]time.Duration),
	}
}

// RecordRelaySuccess records a successful relay publish
func (mc *MetricsCollector) RecordRelaySuccess(relayURL string, duration time.Duration) {
	mc.RelaySuccesses[relayURL]++
	if duration > 0 { // Only record times if a valid duration is provided
		mc.RelaySuccessTimes[relayURL] = append(mc.RelaySuccessTimes[relayURL], duration)
	}
}

// RecordRelayFailure records a failed relay publish
func (mc *MetricsCollector) RecordRelayFailure(relayURL string) {
	mc.RelayFailures[relayURL]++
}

// LogSummary logs a summary of collected metrics using the global logger
func (mc *MetricsCollector) LogSummary() {
	log.Info().
		Int("eventsPostedTotal", mc.EventsPosted).
		Int("eventsSkippedTotal", mc.EventsSkipped).
		Int("eventsFailedTotal", mc.EventsFailed).
		Interface("relaySuccessesPerRelay", mc.RelaySuccesses).
		Interface("relayFailuresPerRelay", mc.RelayFailures).
		Msg("Run Metrics Summary")

	// Example of logging average times, can be expanded
	for relay, times := range mc.RelaySuccessTimes {
		if len(times) > 0 {
			var totalDuration time.Duration
			for _, t := range times {
				totalDuration += t
			}
			avgDuration := totalDuration / time.Duration(len(times))
			log.Debug().Str("relayURL", relay).Dur("averageSuccessfulPublishTimeMs", avgDuration).Int("successfulPublishCount", len(times)).Msg("Relay Performance Detail")
		}
	}
}

// ExportMetrics saves the collected metrics to a JSON file
func (mc *MetricsCollector) ExportMetrics(filePath string) error {
	data, err := json.MarshalIndent(mc, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metrics to JSON: %w", err)
	}
	err = ioutil.WriteFile(filePath, data, 0644) // Standard file permissions
	if err != nil {
		return fmt.Errorf("failed to write metrics JSON to file %s: %w", filePath, err)
	}
	return nil
}

// setupLogger configures the global zerolog logger based on environment variables.
// func setupLogger() {
// Modified to accept Config
func setupLogger(cfg *Config) {
	// Setup Logging (after essential env vars are checked)
	baseLogFileName := "nostr_bot.log"
    // logDir := os.Getenv("LOG_DIR") // Use cfg.LogDir
	var logFilePath string
    if cfg.LogDir != "" {
        if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			// Use fmt.Fprintf for critical errors before logger might be fully set up
			// or if writing to log file itself fails.
			fmt.Fprintf(os.Stderr, "Failed to create log directory %s: %v\n", cfg.LogDir, err)
			// Exiting here because logging is fundamental. If the directory can't be made,
			// subsequent log writes will likely fail or misbehave.
			os.Exit(1) 
        }
        logFilePath = fmt.Sprintf("%s/%s", cfg.LogDir, baseLogFileName)
    } else {
        logFilePath = baseLogFileName
    }
    
    logFile := &lumberjack.Logger{
        Filename:   logFilePath,
        MaxSize:    10, // megabytes
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    }

    zerolog.TimeFieldFormat = time.RFC3339
    hostName, _ := os.Hostname()
	// Default logger to file
    log.Logger = zerolog.New(logFile).With().
        Timestamp().
        Str("service", "nostr-calendar-bot").
		Str("version", "1.1.0"). // Updated version example
        Str("host", hostName).
        Logger()

	// logLevelStr := os.Getenv("LOG_LEVEL") // Use cfg.LogLevel
	switch cfg.LogLevel { // Already lowercased in loadConfig
    case "debug":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    case "info":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    case "warn":
        zerolog.SetGlobalLevel(zerolog.WarnLevel)
    case "error":
        zerolog.SetGlobalLevel(zerolog.ErrorLevel)
    default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // Default to Info (already handled in loadConfig too)
    }
    
	// If CONSOLE_LOG is true, add console writer to the logger
	// if os.Getenv("CONSOLE_LOG") == "true" { // Use cfg.ConsoleLog
	if cfg.ConsoleLog {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		multi := zerolog.MultiLevelWriter(logFile, consoleWriter) // Log to both file and console
		log.Logger = zerolog.New(multi).With().
			Timestamp().
			Str("service", "nostr-calendar-bot").
			Str("version", "1.1.0").
			Str("host", hostName).
			Logger()
		log.Info().Msg("Console logging enabled.") // Log this message to confirm console output
	}
    
    // This message was previously printed before full logger setup in some cases.
    if godotenv.Load() != nil { // Check .env loading status again, log if needed
        log.Info().Msg("No .env file found or error loading it. Relying on system environment variables.")
    }
}