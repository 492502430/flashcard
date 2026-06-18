package model

import "time"

// User represents a WeChat mini-program user.
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OpenID    string    `json:"openid" gorm:"uniqueIndex;not null"`
	Nickname  string    `json:"nickname"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Deck is a collection of flashcards created by a user.
type Deck struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	Title     string    `json:"title" gorm:"not null;size:200"`
	CardCount int       `json:"card_count" gorm:"default:0"`
	Source    string    `json:"source"` // "text" or "pdf"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Cards []Card `json:"cards,omitempty" gorm:"foreignKey:DeckID"`
}

// Card is a single flashcard with FSRS scheduling state.
type Card struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeckID       string    `json:"deck_id" gorm:"index;not null"`
	Question     string    `json:"question" gorm:"type:text;not null"`
	Answer       string    `json:"answer" gorm:"type:text;not null"`
	TagsJSON     string    `json:"-" gorm:"column:tags;type:text"`      // stored as JSON array
	Stability    float64   `json:"stability" gorm:"default:0"`
	Difficulty   float64   `json:"difficulty" gorm:"default:0.5"`
	NextReviewAt time.Time `json:"next_review_at"`
	ReviewCount  int       `json:"review_count" gorm:"default:0"`
	State        string    `json:"state" gorm:"default:'new'"` // new, learning, review
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ReviewRecord logs each review submission for analytics.
type ReviewRecord struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CardID    string    `json:"card_id" gorm:"index;not null"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	Rating    int       `json:"rating"` // 1-4
	Stability float64   `json:"stability"`
	CreatedAt time.Time `json:"created_at"`
}
