package fsrs

import (
	"math"
	"time"
)

// Card represents a flashcard with FSRS scheduling state.
type Card struct {
	State      string  // new, learning, review
	Stability  float64 // days until next review
	Difficulty float64 // 0..1, how hard the card is
}

// Schedule computes the next review time and new stability for a card
// based on the user's rating.
// rating: 1=forgot, 2=hard, 3=good, 4=easy
func Schedule(card *Card, rating int) (time.Time, float64) {
	base := map[int]float64{1: 0.5, 2: 2.0, 3: 7.0, 4: 15.0}

	if card.Stability == 0 {
		card.Stability = base[rating]
	} else if rating == 1 {
		// Forgot — reset stability to a low value
		card.Stability = math.Max(card.Stability*0.2, 0.5)
	} else {
		difficultyFactor := 1.0 + float64(4-rating)*0.15
		card.Stability = card.Stability * difficultyFactor
	}

	// Update difficulty based on rating
	if rating <= 2 {
		card.Difficulty = math.Min(1.0, card.Difficulty+0.1)
	} else {
		card.Difficulty = math.Max(0.1, card.Difficulty-0.05)
	}

	card.State = "review"
	return time.Now().Add(time.Duration(card.Stability*24) * time.Hour), card.Stability
}
