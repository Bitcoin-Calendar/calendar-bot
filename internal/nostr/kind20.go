package nostr

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"calendar-bot/internal/models"
	"github.com/nbd-wtf/go-nostr"
	"github.com/rs/zerolog/log" // For logging within image validation
)

// --- Image Validation ---

// ImageValidator provides methods to validate image URLs for NIP-68 events.
type ImageValidator struct{}

// NewImageValidator creates a new ImageValidator.
func NewImageValidator() *ImageValidator {
	return &ImageValidator{}
}

var supportedImageFormats = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".gif":  "image/gif",
	".webp": "image/webp",
	".avif": "image/avif",
	".apng": "image/apng",
}

// IsValidImageURL checks if the URL points to a supported image format based on extension.
func (iv *ImageValidator) IsValidImageURL(imageURL string) bool {
	u, err := url.Parse(imageURL)
	if err != nil {
		log.Warn().Err(err).Str("url", imageURL).Msg("Failed to parse image URL")
		return false
	}
	ext := strings.ToLower(fileExtension(u.Path))
	_, supported := supportedImageFormats[ext]
	if !supported {
		log.Debug().Str("url", imageURL).Str("extension", ext).Msg("Unsupported image extension")
	}
	return supported
}

// GetMediaType returns the IANA media type for a given image URL based on its extension.
// Returns an empty string if the format is not supported or URL is invalid.
func (iv *ImageValidator) GetMediaType(imageURL string) string {
	u, err := url.Parse(imageURL)
	if err != nil {
		log.Warn().Err(err).Str("url", imageURL).Msg("Failed to parse image URL for media type")
		return ""
	}
	ext := strings.ToLower(fileExtension(u.Path))
	mediaType, ok := supportedImageFormats[ext]
	if !ok {
		log.Debug().Str("url", imageURL).Str("extension", ext).Msg("Cannot determine media type for unsupported extension")
		return ""
	}
	return mediaType
}

// ValidateImageAccessibility checks if the image URL is accessible via a HEAD request.
// This is a basic check and doesn't download the image.
func (iv *ImageValidator) ValidateImageAccessibility(imageURL string) error {
	client := http.Client{Timeout: 10 * time.Second} // Short timeout for HEAD request
	req, err := http.NewRequest("HEAD", imageURL, nil)
	if err != nil {
		log.Warn().Err(err).Str("url", imageURL).Msg("Failed to create HEAD request for image accessibility validation")
		return fmt.Errorf("failed to create HEAD request for %s: %w", imageURL, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Warn().Err(err).Str("url", imageURL).Msg("Failed to perform HEAD request for image accessibility validation")
		return fmt.Errorf("failed to access image URL %s: %w", imageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Warn().Str("url", imageURL).Int("statusCode", resp.StatusCode).Msg("Image not accessible or returned non-OK status")
		return fmt.Errorf("image at %s returned status %d", imageURL, resp.StatusCode)
	}
	log.Debug().Str("url", imageURL).Msg("Image accessibility check passed (HEAD request successful)")
	return nil
}

// helper to get file extension
func fileExtension(path string) string {
	parts := strings.Split(path, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// --- Kind 20 Event Structure & Creation ---

// Kind20EventData represents the data needed to create a NIP-68 picture event.
// Note: `Tags` here are additional non-standard tags. Standard NIP-68 tags like `imeta`, `title`, `m` are handled by methods.
// `Date` here is the original event date string (YYYY-MM-DD) for the `d` tag.
// The struct name was changed to Kind20EventData to avoid conflict with nostr.Event type.
type Kind20EventData struct {
	Title       string   // From APIEvent.Title
	Description string   // From APIEvent.Description, used for summary tag
	ImageURL    string   // From APIEvent.Media
	MediaType   string   // Determined by ImageValidator
	Hashtags    []string // From APIEvent.Tags (parsed)
	References  []string // From APIEvent.References (parsed), for `r` tags
	EventDate   string   // YYYY-MM-DD for `d` tag, from APIEvent.Date
	// Potentially add other fields like Blurhash if we implement that.
}

// ToNostrEvent converts Kind20EventData into a nostr.Event (Kind 20).
// It builds the required NIP-68 tags.
func (k20 *Kind20EventData) ToNostrEvent() (nostr.Event, error) {
	if k20.ImageURL == "" || k20.MediaType == "" {
		return nostr.Event{}, fmt.Errorf("imageURL and mediaType are required for Kind 20 event")
	}

	allTags := nostr.Tags{}

	// Required NIP-68 tags
	allTags = append(allTags, nostr.Tag{"title", k20.Title})
	allTags = append(allTags, nostr.Tag{"imeta", "url " + k20.ImageURL})
	// Add SHA256 hash of the image URL as per NIP-68 recommendation for "imeta"
	hash := sha256.Sum256([]byte(k20.ImageURL))
	allTags = append(allTags, nostr.Tag{"imeta", "x " + hex.EncodeToString(hash[:])})

	// Optional NIP-68 tags
	if k20.Description != "" {
		allTags = append(allTags, nostr.Tag{"summary", k20.Description})
	}
	
	// Media type tag
	allTags = append(allTags, nostr.Tag{"m", k20.MediaType})

	// Preserve existing hashtags (`t` tags)
	for _, ht := range k20.Hashtags {
		if ht != "" {
			allTags = append(allTags, nostr.Tag{"t", strings.ToLower(ht)})
		}
	}

	// Preserve existing references (`r` tags)
	for _, ref := range k20.References {
		if ref != "" {
			allTags = append(allTags, nostr.Tag{"r", ref})
		}
	}

	// Date identifier tag (`d` tag)
	if k20.EventDate != "" {
		allTags = append(allTags, nostr.Tag{"d", k20.EventDate})
	}

	// Add any other specific NIP-68 tags if needed (e.g., thumb, blurhash, dimensions)
	// For now, sticking to the basic requirements plus image hash.

	// Assuming Title and Description will always be present for qualifying events.
	content := fmt.Sprintf("%s\n\n%s", k20.Title, k20.Description)

	ev := nostr.Event{
		CreatedAt: nostr.Now(),
		Kind:      20, // NIP-68 Picture Event Kind
		Tags:      allTags,
		Content:   content,
	}
	return ev, nil
}

// CreateKind20NostrEvent prepares and returns a Kind 20 Nostr event if the API event qualifies.
// It uses ImageValidator for image checks.
// Returns the event, a boolean indicating if it qualified, and an error if creation failed.
func CreateKind20NostrEvent(
	apiEvent models.APIEvent, 
	processedTags []string, 
	processedReferences []string, 
	validator *ImageValidator,
) (event nostr.Event, qualified bool, err error) {

	if apiEvent.Media == "" {
		log.Debug().Uint("apiEventID", apiEvent.ID).Msg("Kind 20: Skipped, no media URL.")
		return nostr.Event{}, false, nil
	}

	if !validator.IsValidImageURL(apiEvent.Media) {
		log.Warn().Uint("apiEventID", apiEvent.ID).Str("mediaURL", apiEvent.Media).Msg("Kind 20: Skipped, invalid or unsupported image format.")
		// Metric for invalid image format should be incremented by caller if needed
		return nostr.Event{}, false, nil // Not an error, just not qualified
	}

	mediaType := validator.GetMediaType(apiEvent.Media)
	if mediaType == "" { // Should be redundant if IsValidImageURL passed, but good check
		log.Warn().Uint("apiEventID", apiEvent.ID).Str("mediaURL", apiEvent.Media).Msg("Kind 20: Skipped, could not determine media type.")
		return nostr.Event{}, false, nil
	}

	// Optional: Validate image accessibility. This makes an external HTTP call.
	// if err := validator.ValidateImageAccessibility(apiEvent.Media); err != nil {
	// 	log.Warn().Err(err).Uint("apiEventID", apiEvent.ID).Str("mediaURL", apiEvent.Media).Msg("Kind 20: Skipped, image not accessible.")
	// 	// Metric for image accessibility failure should be incremented by caller
	// 	return nostr.Event{}, false, nil // Not an error, just not qualified
	// }

	k20Data := Kind20EventData{
		Title:       apiEvent.Title,
		Description: apiEvent.Description, // Used for summary tag
		ImageURL:    apiEvent.Media,
		MediaType:   mediaType,
		Hashtags:    processedTags,
		References:  processedReferences,
		EventDate:   apiEvent.Date.Format("2006-01-02"),
	}

	nostrEv, err := k20Data.ToNostrEvent()
	if err != nil {
		log.Error().Err(err).Uint("apiEventID", apiEvent.ID).Msg("Kind 20: Failed to build Nostr event from Kind20EventData.")
		return nostr.Event{}, false, fmt.Errorf("failed to create kind 20 nostr event: %w", err)
	}

	log.Info().Uint("apiEventID", apiEvent.ID).Str("nostrEventID", nostrEv.ID).Msg("Kind 20: Event created and qualified for publishing.")
	return nostrEv, true, nil
} 