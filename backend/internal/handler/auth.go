package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("flashcard-dev-secret-change-in-production")

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{},
		func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil })
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

type LoginRequest struct {
	Code       string `json:"code"`
	InviteCode string `json:"invite_code,omitempty"`
}

// WxLogin handles WeChat mini-program login.
func (h *Handler) WxLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		writeError(w, 400, "code is required")
		return
	}

	openid := req.Code

	// Try real WeChat code2session if credentials configured
	if h.WxAppID != "" && h.WxAppSecret != "" {
		if realOpenid, _, err := h.Code2Session(req.Code); err == nil {
			openid = realOpenid
		} else {
			log.Printf("WxLogin: code2session failed (%v), using code as openid", err)
		}
	}

	var user struct {
		ID         string
		OpenID     string
		Nickname   string
		InviteCode string
	}

	err := h.DB.Raw(`
		INSERT INTO users (openid) VALUES (?)
		ON CONFLICT (openid) DO UPDATE SET openid=excluded.openid
		RETURNING id, openid, COALESCE(nickname, '') as nickname, COALESCE(invite_code, '') as invite_code
	`, openid).Scan(&user).Error

	if err != nil {
		writeError(w, 500, "failed to create user")
		return
	}

	if user.InviteCode == "" {
		code := generateInviteCode()
		h.DB.Exec(`UPDATE users SET invite_code = ? WHERE id = ?`, code, user.ID)
		user.InviteCode = code
	}

	if req.InviteCode != "" && req.InviteCode != user.InviteCode {
		h.DB.Exec(`UPDATE users SET invited_by = ? WHERE id = ? AND COALESCE(invited_by, '') = ''`,
			req.InviteCode, user.ID)
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		writeError(w, 500, "failed to generate token")
		return
	}

	writeJSON(w, 200, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"nickname": user.Nickname,
		},
	})
}


// UpdateProfile updates the user's nickname and avatar.
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)
	var req struct {
		Nickname  string `json:"nickname"`
		AvatarURL string `json:"avatar_url"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	h.DB.Exec(`UPDATE users SET nickname = ?, avatar_url = ? WHERE id = ?`, req.Nickname, req.AvatarURL, userID)
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
