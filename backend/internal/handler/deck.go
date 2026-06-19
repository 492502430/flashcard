package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/492502430/flashcard/backend/internal/ai"
	"github.com/go-chi/chi/v5"
)

// CreateDeckRequest is the payload for creating a new deck.
type CreateDeckRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// CreateDeck creates a new deck. If text is provided, AI generates cards.
func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		writeError(w, 400, "title is required (max 200 chars)")
		return
	}

	// 1. Create deck
	var deck struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CardCount int    `json:"card_count"`
	}

	err := h.DB.Raw(`
		INSERT INTO decks (user_id, title) VALUES (?, ?)
		RETURNING id, title, card_count
	`, userID, req.Title).Scan(&deck).Error

	if err != nil {
		writeError(w, 500, "failed to create deck")
		return
	}

	// 2. If text provided, generate cards via AI
	if req.Text != "" {
		go h.generateCardsAsync(deck.ID, userID, req.Text)
	}

	writeJSON(w, 201, deck)
}

// generateCardsAsync calls the AI service and inserts cards into the deck.
func (h *Handler) generateCardsAsync(deckID, userID, text string) {
	aiClient := ai.NewClient("http://localhost:8001")

	result, err := aiClient.GenerateCards(text, deckID)
	if err != nil {
		log.Printf("AI generation failed for deck %s: %v", deckID, err)
		return
	}

	// Insert each card
	for _, card := range result.Cards {
		tags, _ := json.Marshal(card.Tags)
		h.DB.Exec(`
			INSERT INTO cards (deck_id, question, answer, tags, next_review_at)
			VALUES (?, ?, ?, ?, NOW())
		`, deckID, card.Question, card.Answer, string(tags))
	}

	// Update card count
	h.DB.Exec(`UPDATE decks SET card_count = ?, updated_at = NOW() WHERE id = ?`,
		result.Count, deckID)

	// Increment user's tokens_used
	if result.TokensUsed > 0 {
		h.DB.Exec(`UPDATE users SET tokens_used = tokens_used + ? WHERE id = ?`,
			result.TokensUsed, userID)
	}

	log.Printf("AI generated %d cards for deck %s (tokens: %d)", result.Count, deckID, result.TokensUsed)
}

// ListDecks returns all decks for the current user.
func (h *Handler) ListDecks(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var decks []struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CardCount int    `json:"card_count"`
		CreatedAt string `json:"created_at"`
	}

	h.DB.Raw(`
		SELECT id, title, card_count, created_at
		FROM decks WHERE user_id = ? ORDER BY created_at DESC
	`, userID).Scan(&decks)

	writeJSON(w, 200, decks)
}

// GetDeck returns a single deck with its cards.
func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	deckID := chi.URLParam(r, "id")

	var deck struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CardCount int    `json:"card_count"`
	}

	err := h.DB.Raw(`
		SELECT id, title, card_count FROM decks WHERE id = ? AND user_id = ?
	`, deckID, userID).Scan(&deck).Error

	if err != nil || deck.ID == "" {
		writeError(w, 404, "deck not found")
		return
	}

	var cards []struct {
		ID       string `json:"id"`
		Question string `json:"question"`
		Answer   string `json:"answer"`
		State    string `json:"state"`
	}
	h.DB.Raw(`SELECT id, question, answer, state FROM cards WHERE deck_id = ? ORDER BY created_at`, deckID).Scan(&cards)

	writeJSON(w, 200, map[string]interface{}{"deck": deck, "cards": cards})
}

// DeleteDeck deletes a deck and its cards.
func (h *Handler) DeleteDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	deckID := chi.URLParam(r, "id")

	result := h.DB.Exec(`DELETE FROM decks WHERE id = ? AND user_id = ?`, deckID, userID)
	if result.RowsAffected == 0 {
		writeError(w, 404, "deck not found")
		return
	}
	writeJSON(w, 200, map[string]bool{"deleted": true})
}
