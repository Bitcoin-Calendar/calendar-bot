package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/rs/zerolog/log" // Assuming global logger is okay for LogSummary
)

// Collector stores metrics about bot operation.
// It includes fields for both general operation and NIP-68 specific events.
type Collector struct {
	// Existing fields from main.go
	EventsPosted      int                    `json:"eventsPosted"` // Renamed for clarity, was Kind 1 implicitly
	EventsSkipped     int                    `json:"eventsSkipped"`
	EventsFailed      int                    `json:"eventsFailed"`  // Renamed for clarity, was Kind 1 implicitly
	RelaySuccesses    map[string]int         `json:"relaySuccesses"`
	RelayFailures     map[string]int         `json:"relayFailures"`
	RelaySuccessTimes map[string][]time.Duration `json:"relaySuccessTimes,omitempty"`

	// New NIP-68 specific metrics fields
	Kind1EventsPosted    int `json:"kind1EventsPosted"`
	Kind1EventsFailed    int `json:"kind1EventsFailed"`
	Kind20EventsPosted   int `json:"kind20EventsPosted"`
	Kind20EventsFailed   int `json:"kind20EventsFailed"`
	Kind20EventsSkipped  int `json:"kind20EventsSkipped"` // No olas tag or invalid image
	ImageValidationFails int `json:"imageValidationFails"`
}

// NewCollector initializes a new MetricsCollector.
func NewCollector() *Collector {
	return &Collector{
		RelaySuccesses:    make(map[string]int),
		RelayFailures:     make(map[string]int),
		RelaySuccessTimes: make(map[string][]time.Duration),
		// NIP-68 fields will be zero-initialized by default
	}
}

// RecordRelaySuccess records a successful relay publish.
func (mc *Collector) RecordRelaySuccess(relayURL string, duration time.Duration) {
	mc.RelaySuccesses[relayURL]++
	if duration > 0 { // Only record times if a valid duration is provided
		mc.RelaySuccessTimes[relayURL] = append(mc.RelaySuccessTimes[relayURL], duration)
	}
}

// RecordRelayFailure records a failed relay publish.
func (mc *Collector) RecordRelayFailure(relayURL string) {
	mc.RelayFailures[relayURL]++
}

// LogSummary logs a summary of collected metrics using the global logger.
// This will need to be updated to show the new NIP-68 fields.
func (mc *Collector) LogSummary() {
	log.Info().
		Int("eventsPostedTotal_DEPRECATED", mc.EventsPosted). // Mark old fields for review/removal
		Int("eventsSkippedTotal", mc.EventsSkipped).
		Int("eventsFailedTotal_DEPRECATED", mc.EventsFailed).   // Mark old fields for review/removal
		Int("kind1EventsPosted", mc.Kind1EventsPosted).
		Int("kind1EventsFailed", mc.Kind1EventsFailed).
		Int("kind20EventsPosted", mc.Kind20EventsPosted).
		Int("kind20EventsFailed", mc.Kind20EventsFailed).
		Int("kind20EventsSkipped", mc.Kind20EventsSkipped).
		Int("imageValidationFails", mc.ImageValidationFails).
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

// ExportMetrics saves the collected metrics to a JSON file.
func (mc *Collector) ExportMetrics(filePath string) error {
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