package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// DeleteCard deletes a single card (must belong to user's deck).
func (h *Handler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	cardID := chi.URLParam(r, "id")

	result := h.DB.Exec(`
		DELETE FROM cards
		USING decks
		WHERE cards.id = ? AND cards.deck_id = decks.id AND decks.user_id = ?
	`, cardID, userID)

	if result.RowsAffected == 0 {
		writeError(w, 404, "card not found")
		return
	}

	// Decrement deck card_count
	h.DB.Exec(`UPDATE decks SET card_count = GREATEST(card_count - 1, 0),
		updated_at = NOW()
		WHERE id IN (SELECT deck_id FROM cards WHERE id = ?)`, cardID)

	writeJSON(w, 200, map[string]string{"status": "deleted"})
}

// updateCardRequest is the request body for updating a card.
type updateCardRequest struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// UpdateCard updates a single card's question and answer (must belong to user's deck).
func (h *Handler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	cardID := chi.URLParam(r, "id")

	var req updateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 400, "invalid request body")
		return
	}

	if req.Question == "" || req.Answer == "" {
		writeError(w, 400, "question and answer are required")
		return
	}

	// Update card with ownership check via deck join
	result := h.DB.Exec(`
		UPDATE cards
		SET question = ?, answer = ?, updated_at = ?
		FROM decks
		WHERE cards.id = ? AND cards.deck_id = decks.id AND decks.user_id = ?
	`, req.Question, req.Answer, time.Now(), cardID, userID)

	if result.RowsAffected == 0 {
		writeError(w, 404, "card not found")
		return
	}

	// Fetch and return the updated card
	var card struct {
		ID           string    `json:"id"`
		DeckID       string    `json:"deck_id"`
		Question     string    `json:"question"`
		Answer       string    `json:"answer"`
		Stability    float64   `json:"stability"`
		Difficulty   float64   `json:"difficulty"`
		NextReviewAt time.Time `json:"next_review_at"`
		ReviewCount  int       `json:"review_count"`
		State        string    `json:"state"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}
	if err := h.DB.Raw(`
		SELECT id, deck_id, question, answer, stability, difficulty,
		       next_review_at, review_count, state, created_at, updated_at
		FROM cards WHERE id = ?
	`, cardID).Scan(&card).Error; err != nil {
		writeError(w, 500, "failed to fetch updated card")
		return
	}

	writeJSON(w, 200, card)
}

// SearchCards searches cards by keyword across all of a user's decks.
func (h *Handler) SearchCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	q := r.URL.Query().Get("q")
	if q == "" {
		writeJSON(w, 200, []interface{}{})
		return
	}

	pattern := fmt.Sprintf("%%%s%%", q)

	type cardResult struct {
		ID        string `json:"id"`
		DeckID    string `json:"deck_id"`
		DeckTitle string `json:"deck_title"`
		Question  string `json:"question"`
		Answer    string `json:"answer"`
		State     string `json:"state"`
	}

	var cards []cardResult
	h.DB.Raw(`
		SELECT cards.id, cards.deck_id, decks.title AS deck_title,
			cards.question, cards.answer, cards.state
		FROM cards
		JOIN decks ON cards.deck_id = decks.id
		WHERE decks.user_id = $1
			AND (cards.question ILIKE $2 OR cards.answer ILIKE $2)
		ORDER BY cards.updated_at DESC
		LIMIT 30
	`, userID, pattern).Scan(&cards)

	writeJSON(w, 200, cards)
}

// ExportAll returns all user data (decks + cards) as JSON.
func (h *Handler) ExportAll(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	type exportCard struct {
		Question string `json:"question"`
		Answer   string `json:"answer"`
		State    string `json:"state"`
	}
	type exportDeck struct {
		Title string       `json:"title"`
		Cards []exportCard `json:"cards"`
	}

	var decks []struct {
		ID    string
		Title string
	}
	h.DB.Raw(`SELECT id, title FROM decks WHERE user_id = ? ORDER BY created_at`, userID).Scan(&decks)

	result := make([]exportDeck, 0, len(decks))
	for _, d := range decks {
		var cards []exportCard
		h.DB.Raw(`SELECT question, answer, state FROM cards WHERE deck_id = ? ORDER BY created_at`, d.ID).Scan(&cards)
		result = append(result, exportDeck{Title: d.Title, Cards: cards})
	}

	writeJSON(w, 200, result)
}
