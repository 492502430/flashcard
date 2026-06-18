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
	log.Println("Database migrated")

	h := handler.New(db)
	r := chi.NewRouter()
	r.Use(chimw.Logger, chimw.Recoverer, middleware.RequestLog)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Post("/api/auth/login", h.WxLogin)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthRequired)
		r.Post("/api/decks", h.CreateDeck)
		r.Get("/api/decks", h.ListDecks)
		r.Get("/api/decks/{id}", h.GetDeck)
		r.Get("/api/decks/{id}/review", h.GetDeckReview)
		r.Delete("/api/decks/{id}", h.DeleteDeck)
		r.Post("/api/upload", h.Upload)
		r.Get("/api/review/today", h.GetDueCards)
		r.Post("/api/review", h.SubmitReview)
	})

	log.Printf("Server starting on :%s", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}
