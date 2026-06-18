package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// CreateDeckRequest is the payload for creating a new deck.
type CreateDeckRequest struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// CreateDeck creates a new deck.
func (h *Handler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" {
		writeError(w, 400, "title is required (max 200 chars)")
		return
	}

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

	writeJSON(w, 201, deck)
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
