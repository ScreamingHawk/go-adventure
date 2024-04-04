package narrator_test

import (
	"testing"

	"github.com/ScreamingHawk/go-adventure/narrator"
)

func TestDecodeChatResponse(t *testing.T) {
	// Create a test response
	response := `{"plot": "You are a brave knight on a quest to rescue the captured princess from the dragon's castle.","choices": ["Approach the castle quietly", "Challenge the dragon to a duel"],"hidden": "Remember to stay focused and use your courage and cleverness to succeed in your quest."}`
	chatResponse, err := narrator.DecodeChatResponse(response)
	if err != nil {
		t.Fatalf("failed to decode chat response: %v", err)
	}
	if chatResponse.Plot != "You are a brave knight on a quest to rescue the captured princess from the dragon's castle." {
		t.Fatalf("expected plot to be %q, got %q", "You are a brave knight on a quest to rescue the captured princess from the dragon's castle.", chatResponse.Plot)
	}
	if len(chatResponse.Choices) != 2 {
		t.Fatalf("expected 2 choices, got %d", len(chatResponse.Choices))
	}
	if chatResponse.Choices[0] != "Approach the castle quietly" {
		t.Fatalf("expected first choice to be %q, got %q", "Approach the castle quietly", chatResponse.Choices[0])
	}
	if chatResponse.Choices[1] != "Challenge the dragon to a duel" {
		t.Fatalf("expected second choice to be %q, got %q", "Challenge the dragon to a duel", chatResponse.Choices[1])
	}
	if chatResponse.Hidden != "Remember to stay focused and use your courage and cleverness to succeed in your quest." {
		t.Fatalf("expected hidden to be %q, got %q", "Remember to stay focused and use your courage and cleverness to succeed in your quest.", chatResponse.Hidden)
	}
}
