package main

import (
	"log"
	"net/http"

	"github.com/492502430/flashcard/backend/internal/config"
	"github.com/492502430/flashcard/backend/internal/handler"
	"github.com/492502430/flashcard/backend/internal/middleware"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate
	db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		openid TEXT UNIQUE NOT NULL,
		nickname TEXT DEFAULT '',
		avatar_url TEXT DEFAULT '',
		invite_code TEXT UNIQUE,
		tokens_used INT DEFAULT 0,
		invited_by TEXT DEFAULT '',
		created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS decks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		title TEXT NOT NULL, card_count INT DEFAULT 0,
		source TEXT DEFAULT 'text', created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS cards (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		deck_id UUID NOT NULL REFERENCES decks(id) ON DELETE CASCADE,
		question TEXT NOT NULL, answer TEXT NOT NULL, tags TEXT DEFAULT '[]',
		stability FLOAT DEFAULT 0, difficulty FLOAT DEFAULT 0.5,
		next_review_at TIMESTAMPTZ DEFAULT NOW(), review_count INT DEFAULT 0,
		state TEXT DEFAULT 'new', created_at TIMESTAMPTZ DEFAULT NOW(),
		updated_at TIMESTAMPTZ DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS review_records (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		rating INT NOT NULL, stability FLOAT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS card_feedbacks (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		card_id UUID NOT NULL REFERENCES cards(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type TEXT NOT NULL,
		created_at TIMESTAMPTZ DEFAULT NOW())`)
	db.Exec(`CREATE TABLE IF NOT EXISTS achievements (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		key TEXT NOT NULL,
		earned_at TIMESTAMPTZ DEFAULT NOW(),
		notified_at TIMESTAMPTZ,
		UNIQUE(user_id, key))`)
	// Migration: add columns for existing databases
	db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS invite_code TEXT UNIQUE`)
	db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS tokens_used INT DEFAULT 0`)
	db.Exec(`ALTER TABLE users ADD COLUMN IF NOT EXISTS invited_by TEXT DEFAULT ''`)
	db.Exec(`ALTER TABLE decks ADD COLUMN IF NOT EXISTS source_name TEXT DEFAULT ''`)
	db.Exec(`ALTER TABLE cards ADD COLUMN IF NOT EXISTS document_name TEXT DEFAULT ''`)
	log.Println("Database migrated")

	h := handler.New(db, cfg.WxAppID, cfg.WxAppSecret)
	r := chi.NewRouter()
	r.Use(chimw.Logger, chimw.Recoverer, middleware.RequestLog)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Admin routes (protected by X-Admin-Token header)
	adminAuth := handler.AdminAuth(cfg.AdminPassword)
	r.Group(func(r chi.Router) {
		r.Use(adminAuth)
		r.Get("/admin", h.AdminPage)
		r.Get("/api/admin/stats", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
				return
			}
			h.AdminStats(w, r)
		})
	})

	// API documentation (no auth)
	r.Get("/api/docs", h.DocsPage)

	// Template decks (listing is public, import requires auth)
	r.Get("/api/templates", h.ListTemplates)
	r.Get("/api/templates/{id}", h.GetTemplate)

	r.Post("/api/auth/login", h.WxLogin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequired)
		r.Post("/api/decks", h.CreateDeck)
		r.Get("/api/decks", h.ListDecks)
		r.Get("/api/decks/{id}", h.GetDeck)
		r.Put("/api/decks/{id}", h.UpdateDeck)
		r.Post("/api/decks/{id}/rename", h.UpdateDeck)
		r.Get("/api/decks/{id}/review", h.GetDeckReview)
		r.Delete("/api/decks/{id}", h.DeleteDeck)
		r.Delete("/api/cards/{id}", h.DeleteCard)
		r.Put("/api/cards/{id}", h.UpdateCard)
		r.Post("/api/cards/optimize", h.OptimizeCards)
		r.Get("/api/cards/search", h.SearchCards)
		r.Get("/api/export", h.ExportAll)
		r.Get("/api/stats", h.GetStats)
		r.Post("/api/upload", h.Upload)
		r.Get("/api/review/today", h.GetDueCards)
		r.Post("/api/review", h.SubmitReview)
		r.Post("/api/cards/{id}/feedback", h.SubmitCardFeedback)
		r.Get("/api/achievements", h.GetAchievements)
		r.Get("/api/checkin", h.GetCheckin)
		r.Get("/api/invite/my-code", h.GetMyInviteCode)
		r.Get("/api/invite/stats", h.GetInviteStats)
		r.Post("/api/templates/{id}/import", h.ImportTemplate)
	})

	log.Printf("Server starting on :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
