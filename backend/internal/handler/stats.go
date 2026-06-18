package handler

import "net/http"

// GetStats returns daily review counts for the past 7 days.
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	type dayStat struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var stats []dayStat
	h.DB.Raw(`
		SELECT created_at::date::text AS date, COUNT(*) AS count
		FROM review_records
		WHERE user_id = ? AND created_at >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY created_at::date
		ORDER BY date
	`, userID).Scan(&stats)

	if stats == nil {
		stats = []dayStat{}
	}

	writeJSON(w, 200, stats)
}
