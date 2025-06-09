package nostr

import (
	"strings"

	"calendar-bot/internal/models"

	"github.com/nbd-wtf/go-nostr"
)

// CreateKind1NostrEvent creates a Nostr kind 1 text event from an APIEvent.
func CreateKind1NostrEvent(apiEvent models.APIEvent, processedTags []string, processedReferences []string) (nostr.Event, error) {
	var finalMessageBuilder strings.Builder
	finalMessageBuilder.WriteString(apiEvent.Title)
	finalMessageBuilder.WriteString("\n\n")
	finalMessageBuilder.WriteString(apiEvent.Description)

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

	defaultTags := []string{"bitcoin", "history", "onthisday", "calendar", "bitcoincalendar", "bitcoinhistory", "autopost"}
	allEventTags := nostr.Tags{}

	for _, t := range defaultTags {
		allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(t)})
	}
	for _, apiTag := range processedTags {
		if apiTag != "" {
			allEventTags = append(allEventTags, nostr.Tag{"t", strings.ToLower(apiTag)})
		}
	}
	allEventTags = append(allEventTags, nostr.Tag{"d", apiEvent.Date.Format("2006-01-02")})

	ev := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      nostr.KindTextNote,
		Tags:      allEventTags,
		Content:   message,
	}

	// The event is not signed here; the EventPublisher handles signing.
	return ev, nil
}
