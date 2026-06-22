package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"gorm.io/gorm"
)

// Handler holds shared dependencies.
type Handler struct {
	DB          *gorm.DB
	WxAppID     string
	WxAppSecret string
}

func New(db *gorm.DB, wxAppID, wxAppSecret string) *Handler {
	return &Handler{DB: db, WxAppID: wxAppID, WxAppSecret: wxAppSecret}
}

// Code2Session calls WeChat's code2session API and returns openid and session_key.
func (h *Handler) Code2Session(code string) (openid, sessionKey string, err error) {
	apiURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		url.QueryEscape(h.WxAppID), url.QueryEscape(h.WxAppSecret), url.QueryEscape(code),
	)
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", "", fmt.Errorf("code2session request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result struct {
		OpenID     string `json:"openid"`
		SessionKey string `json:"session_key"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", fmt.Errorf("code2session parse failed: %w (body: %s)", err, string(body))
	}
	if result.ErrCode != 0 {
		return "", "", fmt.Errorf("code2session error [%d]: %s", result.ErrCode, result.ErrMsg)
	}
	return result.OpenID, result.SessionKey, nil
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
