package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bagaking/ankibuild/anki"
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

func TestWalkTomlFilesReturnsTraversalErrors(t *testing.T) {
	tempDir := t.TempDir()
	t.Chdir(tempDir)

	blockedDir := filepath.Join(tempDir, "blocked")
	if err := os.Mkdir(blockedDir, 0o755); err != nil {
		t.Fatalf("os.Mkdir(%q) error = %v, want nil", blockedDir, err)
	}
	if err := os.Chmod(blockedDir, 0); err != nil {
		t.Fatalf("os.Chmod(%q, 0) error = %v, want nil", blockedDir, err)
	}
	t.Cleanup(func() {
		if err := os.Chmod(blockedDir, 0o755); err != nil {
			t.Errorf("os.Chmod(%q, 0755) error = %v, want nil", blockedDir, err)
		}
	})

	if _, err := os.ReadDir(blockedDir); err == nil {
		t.Skip("filesystem permissions do not block directory reads")
	}

	called := false
	err := WalkTomlFiles(context.Background(), func(ctx context.Context, confK anki.Barn, pth, outDir, fileName string) error {
		called = true
		return nil
	})
	if err == nil {
		t.Fatal("WalkTomlFiles(context.Background(), processor) error = nil, want traversal error")
	}
	if called {
		t.Fatal("WalkTomlFiles(context.Background(), processor) called processor, want no calls")
	}
}
