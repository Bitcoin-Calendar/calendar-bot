package main

import (
	// "context" // No longer used directly
	"encoding/json"
	"fmt"
	// "io/ioutil" // No longer used directly
	// "net/http" // No longer used directly
	"os"
	"runtime"
	"strings"
	"time"

	"calendar-bot/internal/api"
	"calendar-bot/internal/config"
	"calendar-bot/internal/logging"
	"calendar-bot/internal/metrics"
	// "calendar-bot/internal/models" // models.APIEvent is used by api, nostr packages
	"calendar-bot/internal/nostr"

	// "github.com/joho/godotenv" // Now used in logging package
	// "github.com/rs/zerolog" // zerolog/log is used, but not this directly
	"github.com/rs/zerolog/log"
)

// APIEvent struct removed, moved to internal/models/event.go

// cleanURL removes unwanted characters and formatting from a URL string.
func cleanURL(url string) string {
	// First, remove leading/trailing whitespace
	cleaned := strings.TrimSpace(url)
	// Remove potential JSON array wrapping for single items
	cleaned = strings.TrimPrefix(cleaned, "[\"")
	cleaned = strings.TrimSuffix(cleaned, "\"]")
	// Remove list-like prefixes
	cleaned = strings.TrimPrefix(cleaned, "- ")
	// Trim space again in case the prefixes left any
	cleaned = strings.TrimSpace(cleaned)
	return cleaned
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

	// Load configuration using the new config package
	cfg, err := config.LoadConfig(envVarForPrivateKeyName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}

	// Setup Logging using the new logging package
	logging.Setup(cfg)

	// Now that logger is configured, we can use it.
	log.Info().Str("language", cfg.ProcessingLanguage).Msg("Bot configured to process events for language.")

	if cfg.Debug {
		log.Debug().Msg("Debug logging enabled.")
		log.Debug().
			Str("os", runtime.GOOS).
			Str("arch", runtime.GOARCH).
			Str("goVersion", runtime.Version()).
			Int("cpus", runtime.NumCPU()).
			Str("workingDir", getCurrentDirectory()).
			Str("envVarForPrivateKey", cfg.EnvVarForPrivateKey).
			Msg("System information")
		logEnvironmentVariables()
	}

	today := time.Now().Format("01-02") // Format is "MM-DD"
	log.Info().Str("date", today).Msg("Starting bot execution. Fetching events from API.")

	// Initialize MetricsCollector from the new package
	metricsCollector := metrics.NewCollector()
	
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

	// Initialize the API client
	apiClient := api.NewClient(cfg.APIEndpoint, cfg.APIKey)

	// Initialize the EventPublisher from the nostr package
	eventPublisher := nostr.NewEventPublisher(cfg.NostrRelays, cfg.PrivateKey, metricsCollector, log.Logger)
	// Initialize the ImageValidator
	imageValidator := nostr.NewImageValidator()

	// Fetch events using the API client
	// The APIEvent type is now models.APIEvent
	apiEvents, err := apiClient.FetchEvents(currentMonth, currentDay, cfg.ProcessingLanguage)
	if err != nil {
		log.Error().Err(err).Msg("Fatal: Failed to fetch events from API. Bot will exit.")
		metricsCollector.LogSummary() // Use new metricsCollector
		metricsDir := "metrics-logs" // This is also the package name, maybe rename dir?
		if mkDirErr := os.MkdirAll(metricsDir, 0755); mkDirErr != nil {
			log.Error().Err(mkDirErr).Str("directory", metricsDir).Msg("Failed to create metrics directory for error export")
		}
		metricsFilePath := fmt.Sprintf("%s/metrics_error_%s.json", metricsDir, time.Now().Format("2006-01-02_15-04-05"))
		if exportErr := metricsCollector.ExportMetrics(metricsFilePath); exportErr != nil { // Use new metricsCollector
			log.Error().Err(exportErr).Str("file", metricsFilePath).Msg("Failed to export metrics during error shutdown")
		} else {
			log.Info().Str("file", metricsFilePath).Msg("Metrics exported successfully during error shutdown")
		}
		os.Exit(1)
	}
	log.Info().Int("eventsFetchedCount", len(apiEvents)).Msg("Successfully fetched events from API.")

	eventsToPublishToday := 0
	for _, apiEvent := range apiEvents {
		if apiEvent.Date.Format("01-02") == today {
			eventsToPublishToday++
			requestID := fmt.Sprintf("api-event-%d-%s-%d", apiEvent.ID, today, time.Now().UnixNano())
			eventSpecificLogger := log.With().Str("requestID", requestID).Uint("apiEventID", apiEvent.ID).Logger()
			eventSpecificLogger.Info().Str("eventTitle", apiEvent.Title).Msg("Processing matching API event for today")

			// Clean up media and reference URLs
			for i := range apiEvent.Media {
				apiEvent.Media[i] = cleanURL(apiEvent.Media[i])
			}
			currentEventAPIReferences := make([]string, 0, len(apiEvent.References))
			for _, ref := range apiEvent.References {
				currentEventAPIReferences = append(currentEventAPIReferences, cleanURL(ref))
			}

			// Parse tags once for both kind 1 and kind 20
			var currentEventAPITags []string
			if apiEvent.Tags != "" && apiEvent.Tags != "[]" {
				if err := json.Unmarshal([]byte(apiEvent.Tags), &currentEventAPITags); err != nil {
					eventSpecificLogger.Warn().Err(err).Str("tagsString", apiEvent.Tags).Msg("Failed to unmarshal event Tags. Proceeding with no API tags.")
				}
			}

			kind1PublishedSuccessfully := false

			// --- Publish Kind 1 Event ---
			eventSpecificLogger.Info().Msg("Attempting to publish Kind 1 event.")
			kind1NostrEvent, err := nostr.CreateKind1NostrEvent(apiEvent, currentEventAPITags, currentEventAPIReferences)
			if err != nil {
				eventSpecificLogger.Error().Err(err).Msg("Failed to create Kind 1 Nostr event object.")
				metricsCollector.Kind1EventsFailed++
			} else {
				successfulK1Publishes, pubErr := eventPublisher.PublishEvent(apiEvent, kind1NostrEvent, "kind1")
				if pubErr != nil {
					eventSpecificLogger.Error().Err(pubErr).Msg("Failed to sign Kind 1 event.")
					metricsCollector.Kind1EventsFailed++
				} else if successfulK1Publishes > 0 {
					eventSpecificLogger.Info().Int("successfulRelays", successfulK1Publishes).Msg("Kind 1 event successfully published.")
					metricsCollector.Kind1EventsPosted++
					kind1PublishedSuccessfully = true
				} else {
					eventSpecificLogger.Warn().Msg("Kind 1 event was processed but failed to publish to any relay.")
					metricsCollector.Kind1EventsFailed++
				}
			}

			// --- Publish Kind 20 Event (NIP-68) ---
			eventSpecificLogger.Info().Msg("Checking eligibility and attempting to publish Kind 20 event.")
			// The "olas" tag check is part of the feature toggle logic as per spec.
			// The feature toggle is: "media link (ending with png, jpeg, jpg, webp, avif, gif) in the media column"
			// This is handled by CreateKind20NostrEvent via ImageValidator.IsValidImageURL.
			// No explicit "olas" tag check mentioned for this stage in CreateKind20NostrEvent, so proceeding.

			kind20NostrEvent, qualified, errK20Create := nostr.CreateKind20NostrEvent(apiEvent, currentEventAPITags, currentEventAPIReferences, imageValidator)
			if errK20Create != nil {
				eventSpecificLogger.Error().Err(errK20Create).Msg("Error creating Kind 20 Nostr event object.")
				metricsCollector.Kind20EventsFailed++ // Or a new metric like Kind20CreationFailed?
			} else if qualified {
				eventSpecificLogger.Info().Msg("Event qualified for Kind 20. Attempting to publish.")
				successfulK20Publishes, pubErrK20 := eventPublisher.PublishEvent(apiEvent, kind20NostrEvent, "kind20")
				if pubErrK20 != nil {
					eventSpecificLogger.Error().Err(pubErrK20).Msg("Failed to sign Kind 20 event.")
					metricsCollector.Kind20EventsFailed++
				} else if successfulK20Publishes > 0 {
					eventSpecificLogger.Info().Int("successfulRelays", successfulK20Publishes).Msg("Kind 20 event successfully published.")
					metricsCollector.Kind20EventsPosted++
					// Note: If Kind 1 also succeeded, the wait will happen. If only Kind 20, it will also wait here.
					// The requirement is "Dual Publishing Workflow: ... Immediately publish kind 20 version (no delay between them)
					// Maintain 30-minute delay before processing next event"
					// So, the wait should be after *both* attempts for a single API event.
				} else {
					eventSpecificLogger.Warn().Msg("Kind 20 event was processed but failed to publish to any relay.")
					metricsCollector.Kind20EventsFailed++
				}
			} else {
				eventSpecificLogger.Info().Msg("Event did not qualify for Kind 20 publishing (e.g., no valid image, or other criteria).")
				metricsCollector.Kind20EventsSkipped++
			}

			// Wait 30 minutes if at least Kind 1 was successfully published.
			// As per spec: "Dual Publishing Workflow: ... Maintain 30-minute delay before processing next event"
			// This implies the wait happens after attempting both for the current API event.
			if kind1PublishedSuccessfully { // Or if any event type for this API event was successful
				log.Info().Msgf("Waiting %v after processing event ID %d before next event...", eventPublisher.DefaultWaitTime(), apiEvent.ID)
				time.Sleep(eventPublisher.DefaultWaitTime())
			}

		} else {
			metricsCollector.EventsSkipped++
			log.Debug().Uint("apiEventID", apiEvent.ID).Str("eventTitle", apiEvent.Title).Str("eventAPIDate", apiEvent.Date.Format("2006-01-02")).Msg("Skipped API event: Date does not match today.")
		}
	}

	if eventsToPublishToday == 0 {
		log.Info().Msg("No events found for today's date after checking all fetched events.")
	}

	log.Info().Msg("Bot execution finished for today.")
	metricsCollector.LogSummary() // Use new metricsCollector

	metricsDir := "metrics-logs"
	if err := os.MkdirAll(metricsDir, 0755); err != nil {
		log.Error().Err(err).Str("directory", metricsDir).Msg("Failed to create metrics directory for final export")
	}
	metricsFilePath := fmt.Sprintf("%s/metrics_run_%s.json", metricsDir, time.Now().Format("2006-01-02_15-04-05"))
	if err := metricsCollector.ExportMetrics(metricsFilePath); err != nil { // Use new metricsCollector
		log.Error().Err(err).Str("file", metricsFilePath).Msg("Failed to export metrics at end of run")
	} else {
		log.Info().Str("file", metricsFilePath).Msg("Metrics exported successfully at end of run")
	}
}