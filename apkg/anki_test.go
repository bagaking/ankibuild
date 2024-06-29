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
	cardService, noteService, cleanup := newTestCardAndNoteServices(t)
	defer cleanup()

	ctx := context.Background()

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

func TestFindCardByFrontTreatsLikeWildcardsAsLiterals(t *testing.T) {
	tests := []struct {
		name           string
		front          string
		collisionFront string
	}{
		{
			name:           "percent",
			front:          "100%",
			collisionFront: "100x",
		},
		{
			name:           "underscore",
			front:          "card_1",
			collisionFront: "cardx1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cardService, noteService, cleanup := newTestCardAndNoteServices(t)
			defer cleanup()

			ctx := context.Background()

			wantNote, err := noteService.CreateNote(ctx, tt.front, "literal front", NoteWithNID(301), NoteWithGUID("literal"))
			if err != nil {
				t.Fatalf("CreateNote(context.Background(), %q, %q) error = %v, want nil", tt.front, "literal front", err)
			}
			wantCard, err := cardService.CreateCard(401, wantNote)
			if err != nil {
				t.Fatalf("CreateCard(%d, wantNote) error = %v, want nil", 401, err)
			}

			collisionNote, err := noteService.CreateNote(ctx, tt.collisionFront, "wildcard collision", NoteWithNID(302), NoteWithGUID("collision"))
			if err != nil {
				t.Fatalf("CreateNote(context.Background(), %q, %q) error = %v, want nil", tt.collisionFront, "wildcard collision", err)
			}
			if _, err := cardService.CreateCard(402, collisionNote); err != nil {
				t.Fatalf("CreateCard(%d, collisionNote) error = %v, want nil", 402, err)
			}

			notes, cards, err := cardService.FindCardByFront(tt.front)
			if err != nil {
				t.Fatalf("FindCardByFront(%q) error = %v, want nil", tt.front, err)
			}
			if len(notes) != 1 {
				t.Fatalf("FindCardByFront(%q) returned %d notes, want 1", tt.front, len(notes))
			}
			if got, want := notes[0].ID, wantNote.ID; got != want {
				t.Errorf("FindCardByFront(%q) note ID = %d, want %d", tt.front, got, want)
			}
			if len(cards) != 1 {
				t.Fatalf("FindCardByFront(%q) returned %d cards, want 1", tt.front, len(cards))
			}
			if got, want := cards[0].ID, wantCard.ID; got != want {
				t.Errorf("FindCardByFront(%q) card ID = %d, want %d", tt.front, got, want)
			}
		})
	}
}

func newTestCardAndNoteServices(t *testing.T) (*CardService, *NoteService, func()) {
	t.Helper()

	ctx := context.Background()
	storage, err := CreateInMemStorage(ctx, func(db *gorm.DB) error {
		return AutoMigrateModels(db)
	})
	if err != nil {
		t.Fatalf("CreateInMemStorage(context.Background(), AutoMigrateModels) error = %v, want nil", err)
	}

	cleanup := func() {
		if err := storage.Close(); err != nil {
			t.Errorf("storage.Close() error = %v, want nil", err)
		}
	}

	return &CardService{DB: storage.DB}, &NoteService{DB: storage.DB}, cleanup
}
