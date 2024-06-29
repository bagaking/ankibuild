package anki

import "testing"

func TestQnACardRuntimeGettersReturnStoredValues(t *testing.T) {
	card := QnACard{
		Runtime: &Runtime{
			CardID:   123,
			NoteID:   456,
			NoteGUID: "stored-guid",
		},
	}

	if got, want := card.GetNoteID(), 456; got != want {
		t.Errorf("QnACard.GetNoteID() = %d, want %d", got, want)
	}
	if got, want := card.GetNoteGUID(), "stored-guid"; got != want {
		t.Errorf("QnACard.GetNoteGUID() = %q, want %q", got, want)
	}
	if got, want := card.GetCardID(), 123; got != want {
		t.Errorf("QnACard.GetCardID() = %d, want %d", got, want)
	}
}

func TestQnACardRuntimeGettersReturnZeroValuesWhenRuntimeMissing(t *testing.T) {
	card := QnACard{}

	if got, want := card.GetNoteID(), 0; got != want {
		t.Errorf("QnACard.GetNoteID() = %d, want %d", got, want)
	}
	if got, want := card.GetNoteGUID(), ""; got != want {
		t.Errorf("QnACard.GetNoteGUID() = %q, want %q", got, want)
	}
	if got, want := card.GetCardID(), 0; got != want {
		t.Errorf("QnACard.GetCardID() = %d, want %d", got, want)
	}
}
