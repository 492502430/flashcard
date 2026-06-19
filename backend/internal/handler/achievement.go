package handler

import (
	"net/http"
	"time"

	"github.com/492502430/flashcard/backend/internal/model"
)

// GetAchievements returns the user's earned achievements with metadata.
func (h *Handler) GetAchievements(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	// Fetch earned achievements for this user
	type earnedRow struct {
		Key        string     `json:"key"`
		EarnedAt   time.Time  `json:"earned_at"`
		NotifiedAt *time.Time `json:"notified_at"`
	}
	var earned []earnedRow
	h.DB.Raw(`SELECT key, earned_at, notified_at FROM achievements WHERE user_id = ? ORDER BY earned_at`, userID).Scan(&earned)

	// Build a set of earned keys for O(1) lookup
	earnedMap := make(map[string]earnedRow)
	for _, e := range earned {
		earnedMap[e.Key] = e
	}

	// Merge definitions with earned data
	type achievementResult struct {
		Key         string     `json:"key"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Icon        string     `json:"icon"`
		Earned      bool       `json:"earned"`
		EarnedAt    *time.Time `json:"earned_at"`
		NotifiedAt  *time.Time `json:"notified_at"`
	}

	results := make([]achievementResult, 0, len(model.AchievementDefinitions))
	for _, def := range model.AchievementDefinitions {
		res := achievementResult{
			Key:         def.Key,
			Title:       def.Title,
			Description: def.Description,
			Icon:        def.Icon,
		}
		if e, ok := earnedMap[def.Key]; ok {
			res.Earned = true
			res.EarnedAt = &e.EarnedAt
			res.NotifiedAt = e.NotifiedAt
		}
		results = append(results, res)
	}

	writeJSON(w, 200, map[string]interface{}{
		"achievements": results,
	})
}

// CheckAndAwardAchievements checks whether the user has reached any new milestones
// and awards them. Returns a list of newly earned achievement keys.
func (h *Handler) CheckAndAwardAchievements(userID string) []string {
	// Count total reviews
	var totalReviews int
	h.DB.Raw(`SELECT COUNT(*) FROM review_records WHERE user_id = ?`, userID).Scan(&totalReviews)

	// Count streak
	var streak int
	h.DB.Raw(`
		WITH RECURSIVE dates AS (
			SELECT DISTINCT created_at::date AS review_date
			FROM review_records WHERE user_id = ?
		), streak AS (
			SELECT review_date, 1 AS n FROM dates WHERE review_date = CURRENT_DATE
			UNION ALL
			SELECT d.review_date, s.n + 1
			FROM dates d JOIN streak s ON d.review_date = s.review_date - 1
		)
		SELECT COALESCE(MAX(n), 0) FROM streak
	`, userID).Scan(&streak)

	// Determine which achievements should be earned
	shouldEarn := make(map[string]bool)
	if totalReviews >= 1 {
		shouldEarn["first_review"] = true
	}
	if totalReviews >= 10 {
		shouldEarn["cards_10"] = true
	}
	if totalReviews >= 50 {
		shouldEarn["cards_50"] = true
	}
	if totalReviews >= 100 {
		shouldEarn["cards_100"] = true
	}
	if streak >= 7 {
		shouldEarn["streak_7"] = true
	}
	if streak >= 30 {
		shouldEarn["streak_30"] = true
	}

	// Fetch already earned
	var existing []string
	h.DB.Raw(`SELECT key FROM achievements WHERE user_id = ?`, userID).Pluck("key", &existing)
	existingSet := make(map[string]bool)
	for _, k := range existing {
		existingSet[k] = true
	}

	// Insert new ones
	var newKeys []string
	for key := range shouldEarn {
		if !existingSet[key] {
			h.DB.Exec(`INSERT INTO achievements (user_id, key) VALUES (?, ?) ON CONFLICT (user_id, key) DO NOTHING`,
				userID, key)
			newKeys = append(newKeys, key)
		}
	}

	return newKeys
}
