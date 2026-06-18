package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client calls the Python AI service for card generation.
type Client struct {
	baseURL string
	client  *http.Client
}

// Card represents a generated flashcard.
type Card struct {
	Question string   `json:"q"`
	Answer   string   `json:"a"`
	Tags     []string `json:"tags"`
}

// GenerateResponse is the AI service response.
type GenerateResponse struct {
	DeckID string `json:"deck_id"`
	Cards  []Card `json:"cards"`
	Count  int    `json:"count"`
}

// NewClient creates an AI service client.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

// GenerateCards calls the AI service to generate flashcards from text.
func (c *Client) GenerateCards(text, deckID string) (*GenerateResponse, error) {
	body := map[string]string{
		"text":    text,
		"deck_id": deckID,
	}
	payload, _ := json.Marshal(body)

	resp, err := c.client.Post(
		c.baseURL+"/generate",
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return nil, fmt.Errorf("ai service error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ai service returned %d", resp.StatusCode)
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &result, nil
}
