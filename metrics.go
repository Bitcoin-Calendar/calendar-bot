package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// MetricsCollector tracks operational metrics for the bot
type MetricsCollector struct {
	StartTime        time.Time
	EventsPosted     int
	EventsSkipped    int
	EventsFailed     int
	SuccessfulRelays map[string]int
	FailedRelays     map[string]int
	RelayLatencies   map[string][]time.Duration
	mu               sync.Mutex // For thread safety
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		StartTime:        time.Now(),
		SuccessfulRelays: make(map[string]int),
		FailedRelays:     make(map[string]int),
		RelayLatencies:   make(map[string][]time.Duration),
	}
}

// RecordRelaySuccess records a successful relay connection
func (m *MetricsCollector) RecordRelaySuccess(relay string, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.SuccessfulRelays[relay]++
	m.RelayLatencies[relay] = append(m.RelayLatencies[relay], duration)
}

// RecordRelayFailure records a failed relay connection
func (m *MetricsCollector) RecordRelayFailure(relay string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FailedRelays[relay]++
}

// GetAverageLatency calculates the average latency for a relay
func (m *MetricsCollector) GetAverageLatency(relay string) time.Duration {
	m.mu.Lock()
	defer m.mu.Unlock()

	latencies := m.RelayLatencies[relay]
	if len(latencies) == 0 {
		return 0
	}

	var total time.Duration
	for _, duration := range latencies {
		total += duration
	}
	return total / time.Duration(len(latencies))
}

// LogSummary logs a summary of the collected metrics
func (m *MetricsCollector) LogSummary() {
	m.mu.Lock()
	defer m.mu.Unlock()

	totalDuration := time.Since(m.StartTime)

	// Calculate relay performance metrics
	type relayMetric struct {
		Relay      string
		Success    int
		Failures   int
		AvgLatency time.Duration
	}

	var relayMetrics []relayMetric

	// Combine all relay keys
	relays := make(map[string]bool)
	for relay := range m.SuccessfulRelays {
		relays[relay] = true
	}
	for relay := range m.FailedRelays {
		relays[relay] = true
	}

	// Build metrics for each relay
	for relay := range relays {
		metric := relayMetric{
			Relay:    relay,
			Success:  m.SuccessfulRelays[relay],
			Failures: m.FailedRelays[relay],
		}

		if latencies, ok := m.RelayLatencies[relay]; ok && len(latencies) > 0 {
			var total time.Duration
			for _, d := range latencies {
				total += d
			}
			metric.AvgLatency = total / time.Duration(len(latencies))
		}

		relayMetrics = append(relayMetrics, metric)
	}

	// Log the summary
	log.Info().
		Dur("totalRuntime", totalDuration).
		Int("eventsPosted", m.EventsPosted).
		Int("eventsSkipped", m.EventsSkipped).
		Int("eventsFailed", m.EventsFailed).
		Interface("relayPerformance", relayMetrics).
		Msg("Bot execution summary")
}

// ExportMetrics exports metrics to a JSON file
func (m *MetricsCollector) ExportMetrics(filePath string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create metrics summary
	summary := struct {
		Timestamp     time.Time          `json:"timestamp"`
		Runtime       string             `json:"runtime"`
		EventsPosted  int                `json:"events_posted"`
		EventsSkipped int                `json:"events_skipped"`
		EventsFailed  int                `json:"events_failed"`
		RelaySuccess  map[string]int     `json:"relay_success"`
		RelayFailures map[string]int     `json:"relay_failures"`
		AvgLatencies  map[string]string  `json:"avg_latencies"`
	}{
		Timestamp:     time.Now(),
		Runtime:       time.Since(m.StartTime).String(),
		EventsPosted:  m.EventsPosted,
		EventsSkipped: m.EventsSkipped,
		EventsFailed:  m.EventsFailed,
		RelaySuccess:  m.SuccessfulRelays,
		RelayFailures: m.FailedRelays,
		AvgLatencies:  make(map[string]string),
	}

	// Format latencies as human-readable strings
	for relay, latencies := range m.RelayLatencies {
		if len(latencies) > 0 {
			var total time.Duration
			for _, d := range latencies {
				total += d
			}
			avg := total / time.Duration(len(latencies))
			summary.AvgLatencies[relay] = avg.String()
		}
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return err
	}

	// Write to file
	return os.WriteFile(filePath, data, 0644)
}
