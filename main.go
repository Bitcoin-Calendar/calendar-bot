package main

import (
    "context"
    "encoding/csv"
    "fmt"
    "os"
    "runtime"
    "strings"
    "time"
    "github.com/nbd-wtf/go-nostr"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "gopkg.in/natefinch/lumberjack.v2"
    "github.com/joho/godotenv"
)

// getCurrentDirectory gets the current working directory
func getCurrentDirectory() string {
    dir, err := os.Getwd()
    if err != nil {
        return "unknown"
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

    // Use language-specific logs based on CSV file
    baseLogFileName := getLogFileName(csvFilePath)
    
    // Check for custom log directory
    logDir := os.Getenv("LOG_DIR")
    var logFilePath string
    if logDir != "" {
        // Ensure log directory exists
        if err := os.MkdirAll(logDir, 0755); err != nil {
            log.Fatal().Err(err).Str("directory", logDir).Msg("Failed to create log directory")
        }
        logFilePath = fmt.Sprintf("%s/%s", logDir, baseLogFileName)
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

    // Get language for logging context
    language := getLanguageFromCSV(csvFilePath)

    // Configure zerolog to log in human-readable time format
    zerolog.TimeFieldFormat = time.RFC3339
    
    // Add environment info to logs
    hostName, _ := os.Hostname()
    log.Logger = zerolog.New(logFile).With().
        Timestamp().
        Str("service", "nostr-calendar-bot").
        Str("language", language).
        Str("version", "1.0.0").
        Str("host", hostName).
        Logger()

    // Set global log level
    logLevel := os.Getenv("LOG_LEVEL")
    switch strings.ToLower(logLevel) {
    case "debug":
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    case "info":
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    case "warn":
        zerolog.SetGlobalLevel(zerolog.WarnLevel)
    case "error":
        zerolog.SetGlobalLevel(zerolog.ErrorLevel)
    default:
        // Use INFO as default
        zerolog.SetGlobalLevel(zerolog.InfoLevel)
    }
    
    // Enable debug level if explicitly requested (for backward compatibility)
    debugMode := os.Getenv("DEBUG") == "true"
    if debugMode {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
        
        // Log system information
        log.Debug().
            Str("os", runtime.GOOS).
            Str("arch", runtime.GOARCH).
            Str("goVersion", runtime.Version()).
            Int("cpus", runtime.NumCPU()).
            Str("workingDir", getCurrentDirectory()).
            Str("csvFile", csvFilePath).
            Str("envVarForPrivateKey", envVarForPrivateKey).
            Msg("Debug mode enabled - System information")
            
        // Log environment variables (excluding sensitive ones)
        logEnvironmentVariables()
        
        log.Debug().Msg("Debug logging enabled")
    }

    // Near the zerolog setup in main()
    consoleOutput := os.Getenv("CONSOLE_LOG") == "true"
    if consoleOutput {
        // Pretty console output
        consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
        multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
        log.Logger = zerolog.New(multi).With().
            Timestamp().
            Str("service", "nostr-calendar-bot").
            Str("language", language).
            Str("version", "1.0.0").
            Str("host", hostName).
            Logger()
    }

    // Load today's date
    today := time.Now().Format("01-02")
    log.Info().
        Str("date", today).
        Msg("Starting bot execution")

    // Initialize metrics collector
    metrics := NewMetricsCollector()
    
    // Process the CSV file and collect metrics
    processCSV(csvFilePath, privateKey, today, metrics)
    
    // Log metrics summary at the end
    metrics.LogSummary()
    
    // Export metrics to a file
    metricsFilePath := fmt.Sprintf("metrics_%s_%s.json", language, time.Now().Format("2006-01-02"))
    if err := metrics.ExportMetrics(metricsFilePath); err != nil {
        log.Error().Err(err).Str("file", metricsFilePath).Msg("Failed to export metrics")
    } else {
        log.Info().Str("file", metricsFilePath).Msg("Metrics exported successfully")
    }
}

func getLogFileName(filePath string) string {
    // Extract the base filename
    baseName := filePath
    // Find the last slash to get just the filename
    if lastSlash := strings.LastIndex(filePath, "/"); lastSlash != -1 {
        baseName = filePath[lastSlash+1:]
    }
    
    // Check for specific language patterns in the filename
    if strings.Contains(baseName, "_en.") || strings.Contains(baseName, "_en_") {
        return "nostr_bot_en.log"
    } else if strings.Contains(baseName, "_ru.") || strings.Contains(baseName, "_ru_") {
        return "nostr_bot_ru.log"
    }
    
    // For backwards compatibility, check the entire path if no match is found
    if strings.Contains(filePath, "_en") {
        return "nostr_bot_en.log"
    } else if strings.Contains(filePath, "_ru") {
        return "nostr_bot_ru.log"
    }
    
    return "nostr_bot_unknown.log"
}

func getLanguageFromCSV(filePath string) string {
    // Extract the base filename
    baseName := filePath
    // Find the last slash to get just the filename
    if lastSlash := strings.LastIndex(filePath, "/"); lastSlash != -1 {
        baseName = filePath[lastSlash+1:]
    }
    
    // Check for specific language patterns in the filename
    if strings.Contains(baseName, "_en.") || strings.Contains(baseName, "_en_") {
        return "English"
    } else if strings.Contains(baseName, "_ru.") || strings.Contains(baseName, "_ru_") {
        return "Russian"
    }
    
    // For backwards compatibility, check the entire path if no match is found
    if strings.Contains(filePath, "_en") {
        return "English"
    } else if strings.Contains(filePath, "_ru") {
        return "Russian"
    }
    
    return "Unknown"
}

func processCSV(filePath string, privateKey string, today string, metrics *MetricsCollector) {
    // Log file information before opening
    fileInfo, err := os.Stat(filePath)
    if err != nil {
        log.Error().
            Err(err).
            Str("file", filePath).
            Msg("Error checking CSV file")
    } else {
        log.Debug().
            Str("file", filePath).
            Int64("size", fileInfo.Size()).
            Str("modTime", fileInfo.ModTime().Format(time.RFC3339)).
            Msg("CSV file information")
    }

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
            
            // Generate a request ID for tracking this event across logs
            requestID := fmt.Sprintf("%s-%d-%d", today, i, time.Now().UnixNano())
            
            eventLog = eventLog.With().
                Str("requestID", requestID).
                Logger()
                
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
                metrics.EventsFailed++
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
                
                // Log the connection attempt
                relayLog.Debug().Msg("Attempting to connect to relay")
                
                connectStartTime := time.Now()
                ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                relay, err := nostr.RelayConnect(ctx, relayURL)
                connectDuration := time.Since(connectStartTime)
                
                if err != nil {
                    relayLog.Error().
                        Err(err).
                        Dur("connectionAttemptDuration", connectDuration).
                        Msg("Failed to connect to relay")
                    metrics.RecordRelayFailure(relayURL)
                    cancel()
                    continue
                }
                
                relayLog.Debug().
                    Dur("connectionDuration", connectDuration).
                    Msg("Successfully connected to relay")

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
                        metrics.RecordRelayFailure(relayURL)
                    } else {
                        relayLog.Info().
                            Dur("duration", publishDuration).
                            Msg("Successfully posted event to relay")
                        metrics.RecordRelaySuccess(relayURL, publishDuration)
                        successfulRelays++
                    }
                }(relay)
            }

            if successfulRelays > 0 {
                eventLog.Info().
                    Int("successfulRelays", successfulRelays).
                    Int("totalRelays", len(relays)).
                    Msg("Event published successfully to some relays")
                metrics.EventsPosted++

                // Wait 30 minutes before posting the next event
                log.Info().
                    Int("waitMinutes", 30).
                    Msg("Waiting before posting next event")
                time.Sleep(30 * time.Minute)
            } else {
                eventLog.Warn().
                    Msg("Failed to publish event to any relay")
                metrics.EventsFailed++
            }
        } else {
            metrics.EventsSkipped++
            log.Debug().
                Str("eventDate", eventDate).
                Str("eventTitle", record[1]).
                Int("recordIndex", i+1).
                Msg("Skipped event - not matching today's date")
        }
    }
}