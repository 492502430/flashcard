package handler

import (
	"encoding/json"
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

// LoginRequest is the payload from WeChat mini-program wx.login.
type LoginRequest struct {
	Code string `json:"code"`
}

// LoginResponse returns JWT token and user info.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

type UserInfo struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// WxLogin handles WeChat mini-program login.
func (h *Handler) WxLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Code == "" {
		writeError(w, 400, "code is required")
		return
	}

	// TODO: Call WeChat code2session API to get real openid
	openid := req.Code

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

	// Generate invite code for new users (empty invite_code means first login)
	if user.InviteCode == "" {
		code := generateInviteCode()
		h.DB.Exec(`UPDATE users SET invite_code = ? WHERE id = ?`, code, user.ID)
		user.InviteCode = code
	}

	token, err := GenerateToken(user.ID)
	if err != nil {
		writeError(w, 500, "failed to generate token")
		return
	}

	writeJSON(w, 200, LoginResponse{
		Token: token,
		User: UserInfo{ID: user.ID, Nickname: user.Nickname},
	})
}
