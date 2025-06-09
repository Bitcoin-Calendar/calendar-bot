package nostr

import (
	"strings"

	"calendar-bot/internal/models"

	"github.com/nbd-wtf/go-nostr"
	// "github.com/rs/zerolog" // Logger not directly used here, but by caller
)

// CreateKind1NostrEvent creates a Nostr kind 1 text event from an APIEvent.
func CreateKind1NostrEvent(apiEvent models.APIEvent, processedTags []string, processedReferences []string) (nostr.Event, error) {
	// Reconstruct the message content (similar to old publishEvent)
	var finalMessageBuilder strings.Builder
	finalMessageBuilder.WriteString(apiEvent.Title)
	finalMessageBuilder.WriteString("\n\n")
	finalMessageBuilder.WriteString(apiEvent.Description)

	// Append all media URLs if present
	if len(apiEvent.Media) > 0 {
		finalMessageBuilder.WriteString("\n")
		for _, mediaURL := range apiEvent.Media {
			if mediaURL != "" {
				finalMessageBuilder.WriteString("\n")
				finalMessageBuilder.WriteString(mediaURL)
			}
		}
	}

	if len(processedReferences) > 0 {
		finalMessageBuilder.WriteString("\n")
		for _, ref := range processedReferences {
			if ref != "" {
				finalMessageBuilder.WriteString("\n")
				finalMessageBuilder.WriteString(ref)
			}
		}
	}
	message := finalMessageBuilder.String()

	// Nostr tags
	// Default tags could be defined elsewhere or passed in if they become configurable.
	defaultTags := []string{"bitcoin", "history", "onthisday", "calendar", "bitcoincalendar", "bitcoinhistory", "autopost"}
	allEventTags := nostr.Tags{}

	for _, t := range defaultTags {
		allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(t)})
	}
	for _, apiTag := range processedTags {
		if apiTag != "" { // Avoid empty tags
			allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(apiTag)})
		}
	}
	allEventTags = append(allEventTags, nostr.Tag{"d", apiEvent.Date.Format("2006-01-02")}) // Date of the event

	ev := nostr.Event{
		CreatedAt: nostr.Now(), // or time.Now() if nostr.Now() is specific to a context
		Kind:      nostr.KindTextNote,
		Tags:      allEventTags,
		Content:   message,
		// PubKey will be set by Sign method
	}

	// The event is not signed here. Signing is handled by the EventPublisher.
	return ev, nil
}
