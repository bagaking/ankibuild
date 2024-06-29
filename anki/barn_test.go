package anki

import (
	"context"
	"testing"
)

func TestParseTomlContentParsesBarnSettingsCardsAndRuntime(t *testing.T) {
	content := []byte(`
title = "RuntimeExampleBarn"
runtime = true
tags = ["complexity"]
content_fmt = "markdown"

[[q_a]]
question = "this is a runtime example"
answer = "# runtime ans"
tags = ["algorithms"]

[q_a.runtime]
cid = 101
nid = 202
guid = "stored-guid"
`)

	got, err := ParseTomlContent(context.Background(), content)
	if err != nil {
		t.Fatalf("ParseTomlContent(runtime barn TOML) error = %v, want nil", err)
	}

	if got == nil {
		t.Fatal("ParseTomlContent(runtime barn TOML) = nil, want parsed Barn")
	}
	if got.Title != "RuntimeExampleBarn" {
		t.Errorf("ParseTomlContent(runtime barn TOML).Title = %q, want %q", got.Title, "RuntimeExampleBarn")
	}
	if !got.RuntimeEnabled {
		t.Errorf("ParseTomlContent(runtime barn TOML).RuntimeEnabled = %t, want true", got.RuntimeEnabled)
	}
	if got.ContentFormatter != "markdown" {
		t.Errorf("ParseTomlContent(runtime barn TOML).ContentFormatter = %q, want %q", got.ContentFormatter, "markdown")
	}
	if len(got.Tags) != 1 || got.Tags[0] != "complexity" {
		t.Errorf("ParseTomlContent(runtime barn TOML).Tags = %v, want %v", got.Tags, []string{"complexity"})
	}
	if len(got.QnAs) != 1 {
		t.Fatalf("ParseTomlContent(runtime barn TOML).QnAs length = %d, want 1", len(got.QnAs))
	}

	card := got.QnAs[0]
	if card.Question != "this is a runtime example" {
		t.Errorf("ParseTomlContent(runtime barn TOML).QnAs[0].Question = %q, want %q", card.Question, "this is a runtime example")
	}
	if card.Answer != "# runtime ans" {
		t.Errorf("ParseTomlContent(runtime barn TOML).QnAs[0].Answer = %q, want %q", card.Answer, "# runtime ans")
	}
	if len(card.Tags) != 1 || card.Tags[0] != "algorithms" {
		t.Errorf("ParseTomlContent(runtime barn TOML).QnAs[0].Tags = %v, want %v", card.Tags, []string{"algorithms"})
	}
	if card.Runtime == nil {
		t.Fatal("ParseTomlContent(runtime barn TOML).QnAs[0].Runtime = nil, want parsed runtime")
	}
	if got, want := card.GetCardID(), 101; got != want {
		t.Errorf("QnACard.GetCardID() after ParseTomlContent(runtime barn TOML) = %d, want %d", got, want)
	}
	if got, want := card.GetNoteID(), 202; got != want {
		t.Errorf("QnACard.GetNoteID() after ParseTomlContent(runtime barn TOML) = %d, want %d", got, want)
	}
	if got, want := card.GetNoteGUID(), "stored-guid"; got != want {
		t.Errorf("QnACard.GetNoteGUID() after ParseTomlContent(runtime barn TOML) = %q, want %q", got, want)
	}
}

func TestParseTomlContentReturnsErrorForInvalidToml(t *testing.T) {
	content := []byte(`
title = "Broken"
[[q_a]]
question = "missing quote
`)

	if got, err := ParseTomlContent(context.Background(), content); err == nil {
		t.Fatalf("ParseTomlContent(invalid TOML) = %#v, nil error; want non-nil error", got)
	}
}
