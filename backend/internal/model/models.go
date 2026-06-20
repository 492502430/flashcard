package model

import "time"

// User represents a WeChat mini-program user.
type User struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OpenID     string    `json:"openid" gorm:"uniqueIndex;not null"`
	Nickname   string    `json:"nickname"`
	AvatarURL  string    `json:"avatar_url"`
	InviteCode string    `json:"invite_code" gorm:"uniqueIndex;size:12"`
	TokensUsed int       `json:"tokens_used" gorm:"default:0"`
	InvitedBy  string    `json:"invited_by" gorm:"size:12"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Deck is a collection of flashcards created by a user.
type Deck struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id" gorm:"index;not null"`
	Title      string    `json:"title" gorm:"not null;size:200"`
	CardCount  int       `json:"card_count" gorm:"default:0"`
	Source     string    `json:"source"` // "text" or "pdf"
	SourceName string    `json:"source_name" gorm:"default:''"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Cards []Card `json:"cards,omitempty" gorm:"foreignKey:DeckID"`
}

// Card is a single flashcard with FSRS scheduling state.
type Card struct {
	ID           string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	DeckID       string    `json:"deck_id" gorm:"index;not null"`
	Question     string    `json:"question" gorm:"type:text;not null"`
	Answer       string    `json:"answer" gorm:"type:text;not null"`
	TagsJSON     string    `json:"-" gorm:"column:tags;type:text"` // stored as JSON array
	DocumentName string    `json:"document_name" gorm:"default:''"`
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

// CardFeedback tracks user-reported issues with cards.
type CardFeedback struct {
	ID        string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CardID    string    `json:"card_id" gorm:"index;not null"`
	UserID    string    `json:"user_id" gorm:"index;not null"`
	Type      string    `json:"type" gorm:"not null"` // "content_error", "answer_too_brief", "question_unclear"
	CreatedAt time.Time `json:"created_at"`
}

// Achievement represents a milestone badge earned by a user.
type Achievement struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID     string    `json:"user_id" gorm:"index;not null"`
	Key        string    `json:"key" gorm:"not null"` // "first_review", "cards_10", "cards_50", "cards_100", "streak_7", "streak_30"
	EarnedAt   time.Time `json:"earned_at"`
	NotifiedAt time.Time `json:"notified_at"` // null until user sees toast
}

// AchievementDef defines a milestone's metadata (not stored in DB).
type AchievementDef struct {
	Key         string `json:"key"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"` // CSS class or icon name
}

// AchievementDefinitions lists all possible achievements.
var AchievementDefinitions = []AchievementDef{
	{Key: "first_review", Title: "初次记忆", Description: "完成第一次复习", Icon: "star"},
	{Key: "cards_10", Title: "十卡入门", Description: "累计复习10张卡片", Icon: "diamond"},
	{Key: "cards_50", Title: "勤学不辍", Description: "累计复习50张卡片", Icon: "fire"},
	{Key: "cards_100", Title: "百卡达人", Description: "累计复习100张卡片", Icon: "trophy"},
	{Key: "streak_7", Title: "七日坚持", Description: "连续7天复习", Icon: "streak"},
	{Key: "streak_30", Title: "月度之星", Description: "连续30天复习", Icon: "crown"},
}
