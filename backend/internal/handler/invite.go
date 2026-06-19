package handler

import (
	"crypto/rand"
	"math/big"
	"net/http"
)

// generateInviteCode creates a random 8-character alphanumeric code.
func generateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	code := make([]byte, 8)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[n.Int64()]
	}
	return string(code)
}

// GetMyInviteCode returns the current user's invite code.
func (h *Handler) GetMyInviteCode(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var result struct {
		InviteCode string `json:"invite_code"`
	}
	err := h.DB.Raw(`SELECT invite_code FROM users WHERE id = ?`, userID).Scan(&result).Error
	if err != nil || result.InviteCode == "" {
		writeError(w, 500, "invite code not found")
		return
	}

	writeJSON(w, 200, map[string]string{
		"invite_code": result.InviteCode,
	})
}

// GetInviteStats returns how many people used this user's invite code.
func (h *Handler) GetInviteStats(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var inviteCode string
	h.DB.Raw(`SELECT invite_code FROM users WHERE id = ?`, userID).Scan(&inviteCode)
	if inviteCode == "" {
		writeJSON(w, 200, map[string]int{"invited_count": 0})
		return
	}

	var count int
	h.DB.Raw(`SELECT COUNT(*) FROM users WHERE invited_by = ?`, inviteCode).Scan(&count)

	writeJSON(w, 200, map[string]int{
		"invited_count": count,
	})
}
