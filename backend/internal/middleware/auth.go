package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/492502430/flashcard/backend/internal/handler"
)

// UserIDCtxKey is the context key for user ID (exported for handlers).
const UserIDCtxKey = "user_id"

// AuthRequired validates JWT token and sets user_id in context.
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := handler.ParseToken(token)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDCtxKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
