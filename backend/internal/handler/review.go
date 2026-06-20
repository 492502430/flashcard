package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/492502430/flashcard/backend/internal/fsrs"
)

// GetDueCards returns cards due for review today.
func (h *Handler) GetDueCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var cards []struct {
		ID           string    `json:"id"`
		DeckID       string    `json:"deck_id"`
		DeckTitle    string    `json:"deck_title"`
		Question     string    `json:"question"`
		Answer       string    `json:"answer"`
		State        string    `json:"state"`
		Stability    float64   `json:"stability"`
		Difficulty   float64   `json:"difficulty"`
		ReviewCount  int       `json:"review_count"`
		NextReviewAt time.Time `json:"next_review_at"`
	}

	h.DB.Raw(`
		SELECT c.id, c.deck_id, d.title AS deck_title, c.question, c.answer, c.state,
			c.stability, c.difficulty, c.review_count, c.next_review_at
		FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE d.user_id = ? AND c.next_review_at <= ?
		ORDER BY c.next_review_at ASC LIMIT 50
	`, userID, time.Now()).Scan(&cards)

	// Count today's reviews
	var reviewedToday int
	h.DB.Raw(`
		SELECT COUNT(*) FROM review_records
		WHERE user_id = ?
			AND timezone('Asia/Shanghai', created_at)::date = timezone('Asia/Shanghai', NOW())::date
	`, userID).Scan(&reviewedToday)

	// Count all submitted reviews for profile and lifetime stats.
	var reviewedTotal int
	h.DB.Raw(`
		SELECT COUNT(*) FROM review_records
		WHERE user_id = ?
	`, userID).Scan(&reviewedTotal)

	// Count streak (consecutive days with reviews)
	var streak int
	h.DB.Raw(`
		WITH RECURSIVE dates AS (
			SELECT DISTINCT timezone('Asia/Shanghai', created_at)::date AS review_date
			FROM review_records WHERE user_id = ?
		), streak AS (
			SELECT review_date, 1 AS n FROM dates WHERE review_date = timezone('Asia/Shanghai', NOW())::date
			UNION ALL
			SELECT d.review_date, s.n + 1
			FROM dates d JOIN streak s ON d.review_date = s.review_date - 1
		)
		SELECT COALESCE(MAX(n), 0) FROM streak
	`, userID).Scan(&streak)

	writeJSON(w, 200, map[string]interface{}{
		"cards":          cards,
		"total":          len(cards),
		"reviewed_today": reviewedToday,
		"reviewed_total": reviewedTotal,
		"streak":         streak,
	})
}

// SubmitReviewRequest is the payload for submitting a review rating.
type SubmitReviewRequest struct {
	CardID string `json:"card_id"`
	Rating int    `json:"rating"`
}

// SubmitReview processes a review rating and updates card scheduling.
func (h *Handler) SubmitReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req SubmitReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CardID == "" || req.Rating < 1 || req.Rating > 4 {
		writeError(w, 400, "card_id and rating (1-4) are required")
		return
	}

	var card struct {
		ID         string
		State      string
		Stability  float64
		Difficulty float64
	}

	err := h.DB.Raw(`
		SELECT c.id, c.state, c.stability, c.difficulty
		FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE c.id = ? AND d.user_id = ?
	`, req.CardID, userID).Scan(&card).Error

	if err != nil || card.ID == "" {
		writeError(w, 404, "card not found")
		return
	}

	fsrscard := fsrs.Card{State: card.State, Stability: card.Stability, Difficulty: card.Difficulty}
	nextReview, stability := fsrs.Schedule(&fsrscard, req.Rating)

	h.DB.Exec(`
		UPDATE cards SET state='review', stability=?, difficulty=?, next_review_at=?,
			review_count=review_count+1, updated_at=NOW() WHERE id=?
	`, stability, fsrscard.Difficulty, nextReview, req.CardID)

	h.DB.Exec(`INSERT INTO review_records (card_id, user_id, rating, stability) VALUES (?, ?, ?, ?)`,
		req.CardID, userID, req.Rating, stability)

	// Check for newly earned achievements
	newAchievements := h.CheckAndAwardAchievements(userID)

	writeJSON(w, 200, map[string]interface{}{
		"next_review":      nextReview,
		"new_stability":    stability,
		"new_achievements": newAchievements,
	})
}

// GetDeckReview returns due cards for a specific deck.
func (h *Handler) GetDeckReview(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	deckID := chi.URLParam(r, "id")

	var cards []struct {
		ID           string    `json:"id"`
		DeckID       string    `json:"deck_id"`
		DeckTitle    string    `json:"deck_title"`
		Question     string    `json:"question"`
		Answer       string    `json:"answer"`
		State        string    `json:"state"`
		Stability    float64   `json:"stability"`
		Difficulty   float64   `json:"difficulty"`
		ReviewCount  int       `json:"review_count"`
		NextReviewAt time.Time `json:"next_review_at"`
	}

	h.DB.Raw(`
		SELECT c.id, c.deck_id, d.title AS deck_title, c.question, c.answer, c.state,
			c.stability, c.difficulty, c.review_count, c.next_review_at
		FROM cards c
		JOIN decks d ON c.deck_id = d.id
		WHERE d.id = ? AND d.user_id = ? AND c.next_review_at <= ?
		ORDER BY c.next_review_at ASC LIMIT 50
	`, deckID, userID, time.Now()).Scan(&cards)

	writeJSON(w, 200, map[string]interface{}{
		"cards": cards,
		"total": len(cards),
	})
}
