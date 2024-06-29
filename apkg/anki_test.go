package apkg

import (
	"context"
	"testing"

	"gorm.io/gorm"
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

func TestFindCardByFrontMatchesWholeFrontField(t *testing.T) {
	ctx := context.Background()
	storage, err := CreateInMemStorage(ctx, func(db *gorm.DB) error {
		return AutoMigrateModels(db)
	})
	if err != nil {
		t.Fatalf("CreateInMemStorage(context.Background(), AutoMigrateModels) error = %v, want nil", err)
	}
	t.Cleanup(func() {
		if err := storage.Close(); err != nil {
			t.Errorf("storage.Close() error = %v, want nil", err)
		}
	})

	noteService := &NoteService{DB: storage.DB}
	cardService := &CardService{DB: storage.DB}

	prefixNote, err := noteService.CreateNote(ctx, "cat", "small animal", NoteWithNID(101), NoteWithGUID("prefixnote"))
	if err != nil {
		t.Fatalf("CreateNote(context.Background(), %q, %q) error = %v, want nil", "cat", "small animal", err)
	}
	prefixCard, err := cardService.CreateCard(201, prefixNote)
	if err != nil {
		t.Fatalf("CreateCard(%d, prefixNote) error = %v, want nil", 201, err)
	}

	collisionNote, err := noteService.CreateNote(ctx, "catfish", "larger animal", NoteWithNID(102), NoteWithGUID("collision"))
	if err != nil {
		t.Fatalf("CreateNote(context.Background(), %q, %q) error = %v, want nil", "catfish", "larger animal", err)
	}
	if _, err := cardService.CreateCard(202, collisionNote); err != nil {
		t.Fatalf("CreateCard(%d, collisionNote) error = %v, want nil", 202, err)
	}

	notes, cards, err := cardService.FindCardByFront("cat")
	if err != nil {
		t.Fatalf("FindCardByFront(%q) error = %v, want nil", "cat", err)
	}
	if len(notes) != 1 {
		t.Fatalf("FindCardByFront(%q) returned %d notes, want 1", "cat", len(notes))
	}
	if got, want := notes[0].ID, prefixNote.ID; got != want {
		t.Errorf("FindCardByFront(%q) note ID = %d, want %d", "cat", got, want)
	}
	if len(cards) != 1 {
		t.Fatalf("FindCardByFront(%q) returned %d cards, want 1", "cat", len(cards))
	}
	if got, want := cards[0].ID, prefixCard.ID; got != want {
		t.Errorf("FindCardByFront(%q) card ID = %d, want %d", "cat", got, want)
	}
}
