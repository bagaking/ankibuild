package apkg

import (
	"context"
	"testing"
)

func TestDeckColReturnsInitializedCollection(t *testing.T) {
	deck, err := CreateDeck(context.Background(), t.TempDir())
	if err != nil {
		t.Fatalf("CreateDeck(context.Background(), tempDir) error = %v, want nil", err)
	}
	t.Cleanup(func() {
		if err := deck.Close(); err != nil {
			t.Errorf("deck.Close() error = %v, want nil", err)
		}
	})

	col := deck.Col()
	if col == nil {
		t.Fatal("Deck.Col() = nil, want initialized collection")
	}
	if col.ID == 0 {
		t.Errorf("Deck.Col().ID = %d, want non-zero ID", col.ID)
	}
	if col.Models == "" {
		t.Errorf("Deck.Col().Models = %q, want non-empty models JSON", col.Models)
	}
	if col.Decks == "" {
		t.Errorf("Deck.Col().Decks = %q, want non-empty decks JSON", col.Decks)
	}
}
