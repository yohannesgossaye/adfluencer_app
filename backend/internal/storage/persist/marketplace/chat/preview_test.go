package chat

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMessagePreview(t *testing.T) {
	amharic := strings.Repeat("ሰላም እንዴት ነህ ", 10) // 120 runes, 3 bytes per letter
	emojiAtCut := strings.Repeat("a", 96) + "😀😀 and more text after the emoji"

	cases := []struct {
		name      string
		body      string
		wantRunes int
		truncated bool
	}{
		{"short english", "hello", 5, false},
		{"exactly 100 runes", strings.Repeat("x", 100), 100, false},
		{"long english", strings.Repeat("x", 150), 100, true},
		{"long amharic", amharic, 100, true},
		{"emoji at the cut point", emojiAtCut, 100, true},
		{"short amharic", "ሰላም", 3, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := messagePreview(tc.body)
			if !utf8.ValidString(got) {
				t.Fatalf("preview is not valid UTF-8: %q", got)
			}
			if n := utf8.RuneCountInString(got); n != tc.wantRunes {
				t.Fatalf("preview has %d runes, want %d", n, tc.wantRunes)
			}
			if tc.truncated != strings.HasSuffix(got, "...") {
				t.Fatalf("truncated = %v, preview %q", !tc.truncated, got)
			}
		})
	}
}
