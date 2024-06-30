package apkg

import (
	"archive/zip"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestExportToAPKGWritesCollectionAtArchiveRoot(t *testing.T) {
	deckDir := t.TempDir()
	deck, err := CreateDeck(context.Background(), deckDir)
	if err != nil {
		t.Fatalf("CreateDeck(context.Background(), %q) error = %v, want nil", deckDir, err)
	}

	exportPath := filepath.Join(t.TempDir(), "deck.apkg")
	if err := deck.ExportToAPKG(exportPath); err != nil {
		t.Fatalf("Deck.ExportToAPKG(%q) error = %v, want nil", exportPath, err)
	}

	reader, err := zip.OpenReader(exportPath)
	if err != nil {
		t.Fatalf("zip.OpenReader(%q) error = %v, want nil", exportPath, err)
	}
	t.Cleanup(func() {
		if err := reader.Close(); err != nil {
			t.Errorf("zip.Reader.Close() error = %v, want nil", err)
		}
	})

	tests := []struct {
		name string
		want string
	}{
		{
			name: "single collection payload",
			want: ApkgDBName,
		},
	}

	if got, want := len(reader.File), len(tests); got != want {
		t.Fatalf("zip.OpenReader(%q) file count = %d, want %d", exportPath, got, want)
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := reader.File[i]
			if got := file.Name; got != tt.want {
				t.Errorf("zip.OpenReader(%q).File[%d].Name = %q, want %q", exportPath, i, got, tt.want)
			}
			if got := file.UncompressedSize64; got == 0 {
				t.Errorf("zip.OpenReader(%q).File[%d].UncompressedSize64 = %d, want non-zero", exportPath, i, got)
			}
		})
	}

	collectionPath := filepath.Join(deckDir, ApkgDBName)
	if _, err := os.Stat(collectionPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("os.Stat(%q) error = %v, want os.ErrNotExist after export", collectionPath, err)
	}
}
