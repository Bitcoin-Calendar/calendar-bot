package main

import (
    "context"
    "encoding/csv"
    "fmt"
    "os"
    "strings"
    "time"
    "github.com/nbd-wtf/go-nostr"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gopkg.in/natefinch/lumberjack.v2"
    "github.com/joho/godotenv"
)

func main() {
    if len(os.Args) < 3 {
        log.Fatal().Msg("Usage: nostr_bot <csv_file_path> <env_var_for_private_key>")
    }

    csvFilePath := os.Args[1]
    envVarForPrivateKey := os.Args[2]

    err := godotenv.Load()
    if err != nil {
        log.Fatal().Msg("Error loading env file")
    }

    privateKey := os.Getenv(envVarForPrivateKey)
    if privateKey == "" {
        log.Fatal().Msgf("Environment variable %s is not set", envVarForPrivateKey)
    }

    // Set up log rotation for a single log file
    logFile := &lumberjack.Logger{
        Filename:   "nostr_bot.log", // Single log file for all activities
        MaxSize:    10, // megabytes
        MaxBackups: 3,
        MaxAge:     28, // days
        Compress:   true,
    }

    // Configure zerolog to log in human-readable time format
    zerolog.TimeFieldFormat = time.RFC3339
    log.Logger = zerolog.New(logFile).With().
        Timestamp().
        Str("service", "nostr-calendar-bot").
        Logger()

    // Set global log level
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    
    // Enable debug level if explicitly requested
    if os.Getenv("DEBUG") == "true" {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
        log.Debug().Msg("Debug logging enabled")
    }

    // Load today's date
    today := time.Now().Format("01-02")
    log.Info().
        Str("date", today).
        Msg("Starting bot execution")

    processCSV(csvFilePath, privateKey, today)
}

func getLogFileName(filePath string) string {
    if strings.Contains(filePath, "en") {
        return "nostr_bot_en.log"
    } else if strings.Contains(filePath, "ru") {
        return "nostr_bot_ru.log"
    }
    return "nostr_bot_unknown.log"
}

func getLanguageFromCSV(filePath string) string {
    if strings.Contains(filePath, "en") {
        return "English"
    } else if strings.Contains(filePath, "ru") {
        return "Russian"
    }
    return "Unknown"
}

func processCSV(filePath string, privateKey string, today string) {
    // Open the CSV file
    file, err := os.Open(filePath)
    if err != nil {
        log.Fatal().
            Err(err).
            Str("file", filePath).
            Msg("Error opening CSV file")
    }
    defer file.Close()

    // Read the CSV file
    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        log.Fatal().
            Err(err).
            Str("file", filePath).
            Msg("Error reading CSV file")
    }

    log.Debug().
        Int("totalRecords", len(records)-1).
        Msg("CSV records loaded")

    // Prepare Nostr relays
    relays := []string{
        "wss://relay.damus.io",
        "wss://nostr.oxtr.dev",
        "wss://relay.nostr.band",
        "wss://a.nos.lol",
        "wss://relay.primal.net",
    }

    log.Debug().
        Strs("relays", relays).
        Msg("Configured relays")

    // Track metrics
    eventsPosted := 0
    eventsSkipped := 0
    eventsFailed := 0

    // Iterate through the records
    for i, record := range records[1:] { // Skip header
        eventDate := record[0]
        eventMonthDay := eventDate[5:]
        if eventMonthDay == today {
            title := record[1]
            description := record[2]

            eventLog := log.With().
                Str("eventDate", eventDate).
                Str("eventTitle", title).
                Int("recordIndex", i+1).
                Logger()

            eventLog.Info().Msg("Processing matching event for today")

            // Construct the message with line breaks
            descriptionParts := strings.Split(description, "|")
            messageParts := []string{title}
            messageParts = append(messageParts, descriptionParts...) // Add all parts of the description

            message := strings.Join(messageParts, "\n\n") // Join with double line breaks

            // Create a unique set of tags for this event
            eventId := fmt.Sprintf("calendar-event-%s-%d", eventDate, i)
            tags := nostr.Tags{
                nostr.Tag{"calendar", "historical"},
                nostr.Tag{"date", today},
            }

            nostrEvent := nostr.Event{
                CreatedAt: nostr.Now(),
                Kind:      nostr.KindTextNote,
                Tags:      tags,
                Content:   message,
            }

            // Sign the event
            if err := nostrEvent.Sign(privateKey); err != nil {
                eventLog.Error().
                    Err(err).
                    Msg("Failed to sign event")
                eventsFailed++
                continue
            }

            eventLog.Debug().
                Str("eventID", eventId).
                Str("pubkey", nostrEvent.PubKey).
                Msg("Event signed and ready for publishing")

            // Track successful relays
            successfulRelays := 0

            // Publish the event to relays
            for _, relayURL := range relays {
                relayLog := eventLog.With().Str("relay", relayURL).Logger()
                
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                relay, err := nostr.RelayConnect(ctx, relayURL)
                
                if err != nil {
                    relayLog.Error().
                        Err(err).
                        Msg("Failed to connect to relay")
                    cancel()
                    continue
                }

                // Use defer in a function to ensure each relay is closed properly
                func(r *nostr.Relay) {
                    defer r.Close()
                    defer cancel()
                    
                    relayLog.Debug().
                        Str("eventID", nostrEvent.ID).
                        Str("eventContent", nostrEvent.Content).
                        Interface("tags", nostrEvent.Tags).
                        Msg("Preparing to publish event")
                    
                    publishStartTime := time.Now()
                    publishErr := r.Publish(ctx, nostrEvent)
                    publishDuration := time.Since(publishStartTime)
                    
                    if publishErr != nil {
                        relayLog.Error().
                            Err(publishErr).
                            Dur("duration", publishDuration).
                            Msg("Failed to publish event to relay")
                    } else {
                        relayLog.Info().
                            Dur("duration", publishDuration).
                            Msg("Successfully posted event to relay")
                        successfulRelays++
                    }
                }(relay)
            }

            if successfulRelays > 0 {
                eventLog.Info().
                    Int("successfulRelays", successfulRelays).
                    Int("totalRelays", len(relays)).
                    Msg("Event published successfully to some relays")
                eventsPosted++

                // Wait 30 minutes before posting the next event
                log.Info().
                    Int("waitMinutes", 30).
                    Msg("Waiting before posting next event")
                time.Sleep(30 * time.Minute)
            } else {
                eventLog.Warn().
                    Msg("Failed to publish event to any relay")
                eventsFailed++
            }
        } else {
            eventsSkipped++
            log.Debug().
                Str("eventDate", eventDate).
                Str("eventTitle", record[1]).
                Int("recordIndex", i+1).
                Msg("Skipped event - not matching today's date")
        }
    }

    log.Info().
        Int("eventsPosted", eventsPosted).
        Int("eventsSkipped", eventsSkipped).
        Int("eventsFailed", eventsFailed).
        Msg("Bot execution finished")
}