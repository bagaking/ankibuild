package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuildAPKGsFromTomlRejectsEmptyCardsWithoutWritingAPKG(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	configPath := "empty.apkg.toml"
	if err := os.WriteFile(configPath, []byte(`title = "EmptyDeck"`), 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) error = %v, want nil", configPath, err)
	}

	err := BuildAPKGsFromToml(context.Background())
	if err == nil {
		t.Fatal("BuildAPKGsFromToml(context.Background()) error = nil, want non-nil error")
	}
	if got, want := err.Error(), "empty.apkg.toml has no q_a cards"; !strings.Contains(got, want) {
		t.Errorf("BuildAPKGsFromToml(context.Background()) error = %q, want containing %q", got, want)
	}

	apkgPath := filepath.Join(tempDir, "EmptyDeck.apkg")
	if _, statErr := os.Stat(apkgPath); !errors.Is(statErr, os.ErrNotExist) {
		t.Errorf("os.Stat(%q) error = %v, want os.ErrNotExist", apkgPath, statErr)
	}
}
