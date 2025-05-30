package nostr

import (
	"context"
	"time"

	"calendar-bot/internal/metrics"
	"calendar-bot/internal/models" // For APIEvent type

	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog"
)

// EventPublisher handles the creation and publishing of Nostr events.
// It will manage connections to relays and orchestrate different event kinds.
type EventPublisher struct {
	relays         []string
	privateKey     string
	metrics        *metrics.Collector
	logger         zerolog.Logger
	defaultWaitTime time.Duration // Time to wait after a successful publish batch
}

// NewEventPublisher creates a new EventPublisher.
func NewEventPublisher(relays []string, privateKey string, metrics *metrics.Collector, logger zerolog.Logger) *EventPublisher {
	return &EventPublisher{
		relays:         relays,
		privateKey:     privateKey,
		metrics:        metrics,
		logger:         logger.With().Str("component", "EventPublisher").Logger(),
		defaultWaitTime: 30 * time.Minute, // Default from existing logic
	}
}

// DefaultWaitTime returns the default wait time for the EventPublisher
func (ep *EventPublisher) DefaultWaitTime() time.Duration {
	return ep.defaultWaitTime
}

// PublishEvent orchestrates the publishing of an API event to Nostr.
// This will eventually handle both Kind 1 and Kind 20 events.
// For now, it will contain the generic relay publishing logic.
// The actual Nostr event creation will be delegated.
// Returns: successful_publish_count, error (error is primarily for signing issues)
func (ep *EventPublisher) PublishEvent(apiEvent models.APIEvent, nostrEv nostr.Event, eventType string) (int, error) {
	eventSpecificLogger := ep.logger.With().Uint("apiEventID", apiEvent.ID).Str("nostrEventID", nostrEv.ID).Str("eventType", eventType).Logger()
	eventSpecificLogger.Info().Msg("Preparing to publish event to Nostr relays")

	if err := nostrEv.Sign(ep.privateKey); err != nil {
		eventSpecificLogger.Error().Err(err).Msg("Failed to sign Nostr event")
		return 0, err
	}

	eventSpecificLogger.Debug().Str("pubkey", nostrEv.PubKey).Int("tagCount", len(nostrEv.Tags)).Msg("Event signed and ready for publishing")

	successfulRelayPublishes := 0
	for _, relayURL := range ep.relays {
		relayLog := eventSpecificLogger.With().Str("relayURL", relayURL).Logger()
		relayLog.Debug().Msg("Attempting to connect to relay")

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // Connection timeout
		relayConn, err := nostr.RelayConnect(ctx, relayURL)
		if err != nil {
			relayLog.Warn().Err(err).Msg("Failed to connect to relay")
			ep.metrics.RecordRelayFailure(relayURL)
			cancel()
			continue
		}

		relayLog.Debug().Msg("Successfully connected to relay. Preparing to publish.")

		publishCtx, publishCancel := context.WithTimeout(context.Background(), 25*time.Second) // Publish operation timeout
		err = relayConn.Publish(publishCtx, nostrEv)

		if err != nil {
			relayLog.Warn().Err(err).Msg("Failed to publish event to relay")
			ep.metrics.RecordRelayFailure(relayURL)
		} else {
			relayLog.Info().Msg("Event successfully published to relay")
			ep.metrics.RecordRelaySuccess(relayURL, 0) 
			successfulRelayPublishes++
		}

		publishCancel()
		relayConn.Close()
		cancel()
	}

	if successfulRelayPublishes > 0 {
		eventSpecificLogger.Info().Int("successfulRelaysCount", successfulRelayPublishes).Int("totalRelaysAttempted", len(ep.relays)).Msg("Event publishing process completed for one or more relays.")
	} else {
		eventSpecificLogger.Warn().Msg("Event was not successfully published to any of the configured relays.")
	}

	return successfulRelayPublishes, nil
} 