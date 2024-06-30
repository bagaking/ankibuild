package anki

import (
	"context"
	"fmt"
	"slices"
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

func TestParseTomlContentHandlesTOMLBoundaries(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    *Barn
		wantErr bool
	}{
		{
			name: "unknown fields ignored but known fields parse",
			content: []byte(`
title = "KnownTitle"
unexpected_root = "ignored"
runtime = true

[[q_a]]
question = "known question"
answer = "known answer"
unexpected_card = "ignored"
`),
			want: &Barn{
				BarnSetting: BarnSetting{
					Title:          "KnownTitle",
					RuntimeEnabled: true,
				},
				QnAs: []QnACard{
					{
						Question: "known question",
						Answer:   "known answer",
					},
				},
			},
		},
		{
			name: "multiple q_a cards",
			content: []byte(`
title = "TwoCards"

[[q_a]]
question = "first question"
answer = "first answer"
tags = ["first"]

[[q_a]]
question = "second question"
answer = "second answer"
content_fmt = "markdown"
`),
			want: &Barn{
				BarnSetting: BarnSetting{Title: "TwoCards"},
				QnAs: []QnACard{
					{
						Meta:     Meta{Tags: []string{"first"}},
						Question: "first question",
						Answer:   "first answer",
					},
					{
						Meta:     Meta{ContentFormatter: "markdown"},
						Question: "second question",
						Answer:   "second answer",
					},
				},
			},
		},
		{
			name: "missing runtime table leaves Runtime nil",
			content: []byte(`
[[q_a]]
question = "no runtime"
answer = "no runtime answer"
`),
			want: &Barn{
				QnAs: []QnACard{
					{
						Question: "no runtime",
						Answer:   "no runtime answer",
					},
				},
			},
		},
		{
			name: "runtime false still parses per-card runtime data",
			content: []byte(`
runtime = false

[[q_a]]
question = "stored runtime"
answer = "kept even when barn runtime writes are disabled"

[q_a.runtime]
cid = 11
nid = 22
guid = "note-guid"
`),
			want: &Barn{
				BarnSetting: BarnSetting{RuntimeEnabled: false},
				QnAs: []QnACard{
					{
						Question: "stored runtime",
						Answer:   "kept even when barn runtime writes are disabled",
						Runtime: &Runtime{
							CardID:   11,
							NoteID:   22,
							NoteGUID: "note-guid",
						},
					},
				},
			},
		},
		{
			name:    "invalid type returns error",
			content: []byte(`runtime = "not a bool"`),
			wantErr: true,
		},
		{
			name:    "empty content parses as zero-value Barn",
			content: []byte{},
			want:    &Barn{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTomlContent(context.Background(), tt.content)
			if gotErr := err != nil; gotErr != tt.wantErr {
				t.Fatalf("ParseTomlContent(%q) error = %v, want error presence = %t", tt.name, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			assertBarnEqual(t, tt.want, got)
		})
	}
}

func assertBarnEqual(t *testing.T, want, got *Barn) {
	t.Helper()
	if want == nil || got == nil {
		if want != got {
			t.Fatalf("Barn = %#v, want %#v", got, want)
		}
		return
	}

	if got.Title != want.Title {
		t.Errorf("Barn.Title = %q, want %q", got.Title, want.Title)
	}
	if got.RuntimeEnabled != want.RuntimeEnabled {
		t.Errorf("Barn.RuntimeEnabled = %t, want %t", got.RuntimeEnabled, want.RuntimeEnabled)
	}
	assertMetaEqual(t, "Barn.Meta", want.Meta, got.Meta)

	if (got.QnAs == nil) != (want.QnAs == nil) {
		t.Fatalf("Barn.QnAs nil = %t, want %t", got.QnAs == nil, want.QnAs == nil)
	}
	if len(got.QnAs) != len(want.QnAs) {
		t.Fatalf("Barn.QnAs length = %d, want %d", len(got.QnAs), len(want.QnAs))
	}
	for i := range want.QnAs {
		path := fmt.Sprintf("Barn.QnAs[%d]", i)
		assertQnACardEqual(t, path, want.QnAs[i], got.QnAs[i])
	}
}

func assertQnACardEqual(t *testing.T, path string, want, got QnACard) {
	t.Helper()
	assertMetaEqual(t, path+".Meta", want.Meta, got.Meta)
	if got.Question != want.Question {
		t.Errorf("%s.Question = %q, want %q", path, got.Question, want.Question)
	}
	if got.Answer != want.Answer {
		t.Errorf("%s.Answer = %q, want %q", path, got.Answer, want.Answer)
	}
	assertRuntimeEqual(t, path+".Runtime", want.Runtime, got.Runtime)
}

func assertMetaEqual(t *testing.T, path string, want, got Meta) {
	t.Helper()
	if (got.Tags == nil) != (want.Tags == nil) {
		t.Errorf("%s.Tags nil = %t, want %t", path, got.Tags == nil, want.Tags == nil)
		return
	}
	if !slices.Equal(got.Tags, want.Tags) {
		t.Errorf("%s.Tags = %v, want %v", path, got.Tags, want.Tags)
	}
	if got.ContentFormatter != want.ContentFormatter {
		t.Errorf("%s.ContentFormatter = %q, want %q", path, got.ContentFormatter, want.ContentFormatter)
	}
}

func assertRuntimeEqual(t *testing.T, path string, want, got *Runtime) {
	t.Helper()
	if want == nil || got == nil {
		if want != got {
			t.Errorf("%s = %#v, want %#v", path, got, want)
		}
		return
	}
	if got.CardID != want.CardID {
		t.Errorf("%s.CardID = %d, want %d", path, got.CardID, want.CardID)
	}
	if got.NoteID != want.NoteID {
		t.Errorf("%s.NoteID = %d, want %d", path, got.NoteID, want.NoteID)
	}
	if got.NoteGUID != want.NoteGUID {
		t.Errorf("%s.NoteGUID = %q, want %q", path, got.NoteGUID, want.NoteGUID)
	}
}
