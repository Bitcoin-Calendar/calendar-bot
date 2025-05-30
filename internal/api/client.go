package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"calendar-bot/internal/models" // Import the new models package
	"github.com/rs/zerolog/log"
)

const (
	defaultRetryAttempts = 3
	defaultRetryDelay    = 5 * time.Second
)

// Client is a client for interacting with the Bitcoin Calendar API.
// It will handle request construction, sending, and response parsing.
// It will also implement retry logic for transient errors.
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Retries    int
	RetryDelay time.Duration
}

// NewClient creates a new API client.
func NewClient(baseURL string, apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second}, // Sensible default timeout
		Retries:    defaultRetryAttempts,
		RetryDelay: defaultRetryDelay,
	}
}

// FetchEvents retrieves events for a specific month, day, and language from the API.
// It includes retry logic for transient errors.
func (c *Client) FetchEvents(month string, day string, language string) ([]models.APIEvent, error) {
	var lastErr error

	// Construct the API URL
	// Corrected: Use /events endpoint with query parameters.
	// c.BaseURL should be like http://host:port/api
	url := fmt.Sprintf("%s/events?month=%s&day=%s&lang=%s", c.BaseURL, month, day, language)
	log.Debug().Str("url", url).Msg("Constructed API URL for fetching events by date and language")

	for i := 0; i < c.Retries; i++ {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create API request: %w", err)
		}
		req.Header.Set("X-API-Key", c.APIKey)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to send API request on attempt %d: %w", i+1, err)
			log.Warn().Err(lastErr).Int("attempt", i+1).Int("maxRetries", c.Retries).Msg("API request failed, retrying...")
			time.Sleep(c.RetryDelay)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			lastErr = fmt.Errorf("API request failed with status code %d on attempt %d: %s", resp.StatusCode, i+1, string(bodyBytes))
			log.Warn().Err(lastErr).Int("attempt", i+1).Int("maxRetries", c.Retries).Int("statusCode", resp.StatusCode).Msg("API request non-OK status, retrying...")
			// For certain status codes (e.g., 4xx client errors), retrying might not be useful.
			// However, the current logic retries on any non-200. This can be refined.
			time.Sleep(c.RetryDelay)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// This error is less likely to be transient, but we'll retry based on current loop structure
			lastErr = fmt.Errorf("failed to read API response body on attempt %d: %w", i+1, err)
			log.Warn().Err(lastErr).Int("attempt", i+1).Int("maxRetries", c.Retries).Msg("Failed to read API response body, retrying...")
			time.Sleep(c.RetryDelay)
			continue
		}

		var apiResponse models.APIResponseWrapper // Use the wrapper struct
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal API response body: %w", err)
		}

		return apiResponse.Events, nil // Return the slice of events from the wrapper
	}

	// If loop finishes, all retries failed
	return nil, fmt.Errorf("failed to fetch events from API after %d attempts: %w", c.Retries, lastErr)
} 