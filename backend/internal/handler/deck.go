package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/492502430/flashcard/backend/internal/ai"
	"github.com/go-chi/chi/v5"
)

// CreateDeckRequest is the payload for creating a new deck.
type CreateDeckRequest struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	SourceName string `json:"source_name"`
	CardCount  int    `json:"card_count"`
}

type UpdateDeckRequest struct {
	Title string `json:"title"`
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
		ID         string `json:"id"`
		Title      string `json:"title"`
		CardCount  int    `json:"card_count"`
		SourceName string `json:"source_name"`
	}

	err := h.DB.Raw(`
		INSERT INTO decks (user_id, title, source_name) VALUES (?, ?, ?)
		RETURNING id, title, card_count, source_name
	`, userID, req.Title, req.SourceName).Scan(&deck).Error

	if err != nil {
		writeError(w, 500, "failed to create deck")
		return
	}

	// 2. If text provided, generate cards via AI
	if req.Text != "" {
		sourceName := req.SourceName
		if sourceName == "" {
			sourceName = req.Title
		}
		cardCount := req.CardCount
		if cardCount <= 0 {
			cardCount = 25
		}
		go h.generateCardsAsync(deck.ID, userID, req.Text, sourceName, cardCount)
	}

	writeJSON(w, 201, deck)
}

// generateCardsAsync calls the AI service and inserts cards into the deck.
func (h *Handler) generateCardsAsync(deckID, userID, text, sourceName string, cardCount int) {
	aiClient := ai.NewClient("http://localhost:8001")

	result, err := aiClient.GenerateCards(text, deckID, cardCount)
	if err != nil {
		log.Printf("AI generation failed for deck %s: %v", deckID, err)
		return
	}

	// Insert each card
	for _, card := range result.Cards {
		tags, _ := json.Marshal(card.Tags)
		h.DB.Exec(`
			INSERT INTO cards (deck_id, question, answer, tags, document_name, next_review_at)
			VALUES (?, ?, ?, ?, ?, NOW())
		`, deckID, card.Question, card.Answer, string(tags), sourceName)
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
		ID         string `json:"id"`
		Title      string `json:"title"`
		CardCount  int    `json:"card_count"`
		SourceName string `json:"source_name"`
		CreatedAt  string `json:"created_at"`
	}

	h.DB.Raw(`
		SELECT id, title, card_count, source_name, created_at
		FROM decks WHERE user_id = ? ORDER BY created_at DESC
	`, userID).Scan(&decks)

	writeJSON(w, 200, decks)
}

// GetDeck returns a single deck with its cards.
func (h *Handler) GetDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	deckID := chi.URLParam(r, "id")

	var deck struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		CardCount  int    `json:"card_count"`
		SourceName string `json:"source_name"`
	}

	err := h.DB.Raw(`
		SELECT id, title, card_count, source_name FROM decks WHERE id = ? AND user_id = ?
	`, deckID, userID).Scan(&deck).Error

	if err != nil || deck.ID == "" {
		writeError(w, 404, "deck not found")
		return
	}

	var cards []struct {
		ID           string `json:"id"`
		Question     string `json:"question"`
		Answer       string `json:"answer"`
		State        string `json:"state"`
		Tags         string `json:"tags"`
		DocumentName string `json:"document_name"`
	}
	h.DB.Raw(`SELECT id, question, answer, state, tags, document_name FROM cards WHERE deck_id = ? ORDER BY created_at`, deckID).Scan(&cards)

	writeJSON(w, 200, map[string]interface{}{"deck": deck, "cards": cards})
}

// UpdateDeck updates deck metadata owned by the current user.
func (h *Handler) UpdateDeck(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	deckID := chi.URLParam(r, "id")

	var req UpdateDeckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		writeError(w, 400, "title is required")
		return
	}
	if len([]rune(title)) > 80 {
		writeError(w, 400, "title is too long")
		return
	}

	var deck struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		CardCount  int    `json:"card_count"`
		SourceName string `json:"source_name"`
	}

	result := h.DB.Exec(`
		UPDATE decks
		SET title = ?, updated_at = NOW()
		WHERE id = ? AND user_id = ?
	`, title, deckID, userID)
	if result.Error != nil {
		writeError(w, 500, "failed to update deck")
		return
	}
	if result.RowsAffected == 0 {
		writeError(w, 404, "deck not found")
		return
	}

	h.DB.Raw(`
		SELECT id, title, card_count, source_name
		FROM decks
		WHERE id = ? AND user_id = ?
	`, deckID, userID).Scan(&deck)

	writeJSON(w, 200, deck)
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

// OptimizeCards proxies the AI card optimization request.
func (h *Handler) OptimizeCards(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Cards []map[string]interface{} `json:"cards"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid body")
		return
	}

	body, _ := json.Marshal(map[string]interface{}{"cards": req.Cards})

	resp, err := http.Post("http://localhost:8001/optimize", "application/json", bytes.NewReader(body))
	if err != nil {
		writeError(w, 500, "AI service unreachable")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if resp.StatusCode != 200 {
		writeError(w, resp.StatusCode, fmt.Sprintf("optimize failed: %v", result))
		return
	}

	writeJSON(w, 200, result)
}
