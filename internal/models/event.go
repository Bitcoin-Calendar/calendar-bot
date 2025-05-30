package models

import "time"

// APIEvent struct to match the expected data structure from the API
// This struct is also used by the api client, but defined here as a shared model.
type APIEvent struct {
	ID          uint      `json:"ID"`
	Date        time.Time `json:"Date"`
	Title       string    `json:"Title"`
	Description string    `json:"Description"`
	Tags        string    `json:"Tags"`       // JSON array string
	Media       string    `json:"Media"`      // URL
	References  string    `json:"References"` // JSON array string
	Hashtags    []string  `json:"hashtags"`
	Olas        bool      `json:"olas"` // Assuming 'olas' is a boolean based on previous discussions
}

// APIResponseWrapper represents the full structure of the API response,
// including the list of events and pagination information.
type APIResponseWrapper struct {
	Events     []APIEvent  `json:"events"`
	Pagination interface{} `json:"pagination"` // Using interface{} for now as pagination details aren't used
}

// Constants for event types, if needed elsewhere 