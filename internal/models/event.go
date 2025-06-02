package models

import (
	"encoding/json"
	"time"
)

// APIEvent struct to match the expected data structure from the API.
// The Media field will store multiple URLs.
type APIEvent struct {
	ID          uint      `json:"ID"`
	Date        time.Time `json:"Date"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Tags        string    `json:"Tags"`       // JSON array string from API
	Media       []string  `json:"Media"`      // Parsed from JSON array string from API
	References  string    `json:"References"` // JSON array string from API
	Hashtags    []string  `json:"hashtags"`   // Parsed from Tags by the bot
	Olas        bool      `json:"olas"`
}

// Custom unmarshalling logic for APIEvent
// This is necessary because the 'Media', 'Tags', and 'References' fields
// come as JSON strings from the API, but we want to use them as structured types.
type apiEventRaw struct {
	ID          uint      `json:"ID"`
	Date        time.Time `json:"Date"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Tags        string    `json:"Tags"`
	Media       string    `json:"Media"` // Media comes as a string (potentially a JSON array string)
	References  string    `json:"References"`
	// Hashtags are processed later by the bot from Tags, not directly from this initial parse.
	Olas bool `json:"olas"`
}

func (ae *APIEvent) UnmarshalJSON(data []byte) error {
	var raw apiEventRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	ae.ID = raw.ID
	ae.Date = raw.Date
	ae.Title = raw.Title
	ae.Description = raw.Description
	ae.Tags = raw.Tags             // Keep as string, bot parses this into ae.Hashtags
	ae.References = raw.References // Keep as string, bot parses this
	ae.Olas = raw.Olas

	// Unmarshal Media string into []string
	if raw.Media != "" && raw.Media != "[]" {
		if err := json.Unmarshal([]byte(raw.Media), &ae.Media); err != nil {
			// If it's not a valid JSON array, treat it as a single URL in a slice
			// This provides backward compatibility if some entries are single URLs
			// and not JSON arrays. Or you could return an error here.
			// For now, let's log a warning or handle as a single item.
			// If it's not already a JSON array string, we'll assume it's a single URL.
			// However, the goal is to store comma-separated or JSON array in DB.
			// Let's assume for now API will provide a valid JSON array string for Media.
			// If not, this part needs robust error handling or decision.
			// For now, strict unmarshalling:
			return err // Or, ae.Media = []string{raw.Media} if you want to be lenient
		}
	} else {
		ae.Media = []string{} // Ensure it's an empty slice, not nil
	}

	// Hashtags are typically derived from Tags string later in the bot's processing logic
	// So, ae.Hashtags is not populated here from raw.Hashtags (which isn't in apiEventRaw)

	return nil
}

// APIResponseWrapper represents the full structure of the API response,
// including the list of events and pagination information.
type APIResponseWrapper struct {
	Events     []APIEvent  `json:"events"`
	Pagination interface{} `json:"pagination"` // Using interface{} for now as pagination details aren't used
}

// Constants for event types, if needed elsewhere
