package handler

import "net/http"

// GetCheckin returns review counts for the past 30 days (one row per day).
// Used by the frontend calendar heatmap.
func (h *Handler) GetCheckin(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	type dayStat struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}

	var stats []dayStat
	h.DB.Raw(`
		SELECT created_at::date::text AS date, COUNT(*) AS count
		FROM review_records
		WHERE user_id = ? AND created_at >= CURRENT_DATE - INTERVAL '30 days'
		GROUP BY created_at::date
		ORDER BY date
	`, userID).Scan(&stats)

	if stats == nil {
		stats = []dayStat{}
	}

	writeJSON(w, 200, stats)
}
