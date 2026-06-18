package fsrs

import (
	"testing"
	"time"
)

func TestScheduleNewCard_Good(t *testing.T) {
	card := Card{State: "new", Stability: 0, Difficulty: 0.5}

	nextReview, stability := Schedule(&card, 3) // rating=good

	if stability < 5 || stability > 10 {
		t.Errorf("expected stability 5-10 days for rating=good, got %.1f", stability)
	}

	expectedMin := time.Now().Add(5 * 24 * time.Hour)
	if nextReview.Before(expectedMin) {
		t.Errorf("next review too soon: %v (expected >= %v)", nextReview, expectedMin)
	}

	if card.State != "review" {
		t.Errorf("expected state 'review', got '%s'", card.State)
	}
}

func TestScheduleNewCard_Easy(t *testing.T) {
	card := Card{State: "new", Stability: 0, Difficulty: 0.5}

	_, stability := Schedule(&card, 4) // rating=easy

	if stability < 12 || stability > 18 {
		t.Errorf("expected stability 12-18 days for rating=easy, got %.1f", stability)
	}
}

func TestScheduleNewCard_Forgot(t *testing.T) {
	card := Card{State: "new", Stability: 0, Difficulty: 0.5}

	_, stability := Schedule(&card, 1) // rating=forgot

	if stability > 1 {
		t.Errorf("forgotten new card should have very low stability, got %.1f", stability)
	}
}

func TestScheduleReview_Forgotten(t *testing.T) {
	card := Card{State: "review", Stability: 30, Difficulty: 0.5}

	_, stability := Schedule(&card, 1) // rating=forgot

	// Forgotten card: stability should drop significantly
	if stability > 8 {
		t.Errorf("forgotten card stability should reset low, got %.1f", stability)
	}
}

func TestScheduleReview_Good_StabilityIncreases(t *testing.T) {
	card := Card{State: "review", Stability: 7, Difficulty: 0.5}

	_, stability := Schedule(&card, 3) // rating=good

	if stability <= 7 {
		t.Errorf("stability should increase from 7, got %.1f", stability)
	}
}

func TestScheduleReview_Hard_IncreasesDifficulty(t *testing.T) {
	card := Card{State: "review", Stability: 10, Difficulty: 0.5}

	Schedule(&card, 2) // rating=hard

	if card.Difficulty <= 0.5 {
		t.Errorf("difficulty should increase for hard rating, got %.3f", card.Difficulty)
	}
}

func TestScheduleReview_Easy_DecreasesDifficulty(t *testing.T) {
	card := Card{State: "review", Stability: 10, Difficulty: 0.5}

	Schedule(&card, 4) // rating=easy

	if card.Difficulty >= 0.5 {
		t.Errorf("difficulty should decrease for easy rating, got %.3f", card.Difficulty)
	}
}

func TestSchedule_DifficultyClamped(t *testing.T) {
	// Difficulty should never exceed 1.0
	card := Card{State: "review", Stability: 10, Difficulty: 0.95}
	for i := 0; i < 10; i++ {
		Schedule(&card, 1) // forgot repeatedly
	}
	if card.Difficulty > 1.0 {
		t.Errorf("difficulty should be clamped at 1.0, got %.3f", card.Difficulty)
	}

	// Difficulty should never go below 0.1
	card2 := Card{State: "review", Stability: 10, Difficulty: 0.15}
	for i := 0; i < 10; i++ {
		Schedule(&card2, 4) // easy repeatedly
	}
	if card2.Difficulty < 0.1 {
		t.Errorf("difficulty should not go below 0.1, got %.3f", card2.Difficulty)
	}
}
