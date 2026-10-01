package unit

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/arandu-io/kyse/components"
)

// TestHighlightCutsTheOriginalTextAtRuneBoundaries holds the matching to the
// text it was given. Lower-casing a string can change its length -- "Ⱥ" is two
// bytes and its lower case is three -- so offsets found in a lowered copy point
// somewhere else in the original, past its end or into the middle of a
// character.
func TestHighlightCutsTheOriginalTextAtRuneBoundaries(t *testing.T) {
	cases := []struct {
		name, text, query string
		sensitive         bool
		want              []components.HighlightSegment
	}{
		{
			name: "a lower case that grows does not run past the end",
			text: "ȺȺȺȺx", query: "x",
			want: []components.HighlightSegment{{Text: "ȺȺȺȺ"}, {Text: "x", Match: true}},
		},
		{
			name: "a lower case that grows is still matched",
			text: "aȺb", query: "ⱥ",
			want: []components.HighlightSegment{{Text: "a"}, {Text: "Ⱥ", Match: true}, {Text: "b"}},
		},
		{
			name: "a lower case that splits into two runes does not shift the rest",
			text: "İİİİx", query: "x",
			want: []components.HighlightSegment{{Text: "İİİİ"}, {Text: "x", Match: true}},
		},
		{
			name: "the Kelvin sign is matched whole, never half of it",
			text: "Kab<script>", query: "k",
			want: []components.HighlightSegment{{Text: "K", Match: true}, {Text: "ab<script>"}},
		},
		{
			name: "a lower case that shrinks does not cut the next character",
			text: "KKab<script>", query: "ab",
			want: []components.HighlightSegment{{Text: "KK"}, {Text: "ab", Match: true}, {Text: "<script>"}},
		},
		{
			name: "sharp s folds to its capital",
			text: "Straße", query: "STRAẞE",
			want: []components.HighlightSegment{{Text: "Straße", Match: true}},
		},
		{
			name: "sharp s is not expanded to two letters",
			text: "Straße", query: "strasse",
			want: []components.HighlightSegment{{Text: "Straße"}},
		},
		{
			name: "emoji around a match",
			text: "👍OK👍", query: "ok",
			want: []components.HighlightSegment{{Text: "👍"}, {Text: "OK", Match: true}, {Text: "👍"}},
		},
		{
			name: "an emoji query",
			text: "a👍b", query: "👍",
			want: []components.HighlightSegment{{Text: "a"}, {Text: "👍", Match: true}, {Text: "b"}},
		},
		{
			name: "an empty query marks nothing",
			text: "Ada Lovelace", query: "",
			want: []components.HighlightSegment{{Text: "Ada Lovelace"}},
		},
		{
			name: "overlapping matches are taken left to right without overlap",
			text: "aaa", query: "AA",
			want: []components.HighlightSegment{{Text: "aa", Match: true}, {Text: "a"}},
		},
		{
			name: "repeated matches are all marked",
			text: "aaaa", query: "aa",
			want: []components.HighlightSegment{{Text: "aa", Match: true}, {Text: "aa", Match: true}},
		},
		{
			name: "a case-sensitive query ignores other cases",
			text: "Ada ada", query: "ada", sensitive: true,
			want: []components.HighlightSegment{{Text: "Ada "}, {Text: "ada", Match: true}},
		},
		{
			name: "a query longer than the text matches nothing",
			text: "Ⱥ", query: "ⱥⱥ",
			want: []components.HighlightSegment{{Text: "Ⱥ"}},
		},
		{
			name: "invalid UTF-8 is carried through unchanged",
			text: "a\xffb", query: "B",
			want: []components.HighlightSegment{{Text: "a\xff"}, {Text: "b", Match: true}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			props := components.HighlightProps{Text: tc.text, Query: tc.query, CaseSensitive: tc.sensitive}
			var got []components.HighlightSegment
			func() {
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("Segments panicked on %q / %q: %v", tc.text, tc.query, r)
					}
				}()
				got = props.Segments()
			}()
			if len(got) != len(tc.want) {
				t.Fatalf("Segments(%q, %q) = %+v, want %+v", tc.text, tc.query, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("Segments(%q, %q) = %+v, want %+v", tc.text, tc.query, got, tc.want)
				}
			}
			assertSegmentsRebuildText(t, tc.text, got)
			// The rendered component must not panic either.
			_ = components.Highlight(props)
		})
	}
}

// TestHighlightNeverSplitsACharacter runs every query that is one character
// against text made of characters whose case forms differ in length, and
// checks that each piece is whole UTF-8 and the pieces add up to the text.
func TestHighlightNeverSplitsACharacter(t *testing.T) {
	text := "ȺⱥİiıKKkſsßẞ👍ÅåÅ<&>\"x"
	for _, query := range []string{"Ⱥ", "ⱥ", "İ", "i", "I", "ı", "K", "k", "K", "S", "ſ", "ß", "ẞ", "👍", "Å", "å", "Å", "<", "x", "X", "ks", "KS"} {
		props := components.HighlightProps{Text: text, Query: query}
		segments := func() (out []components.HighlightSegment) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("Segments panicked on query %q: %v", query, r)
				}
			}()
			return props.Segments()
		}()
		assertSegmentsRebuildText(t, text, segments)
		for _, segment := range segments {
			if segment.Match && !strings.EqualFold(segment.Text, query) {
				t.Errorf("query %q marked %q, which does not fold to it", query, segment.Text)
			}
		}
	}
}

// TestHighlightEscapesAroundAFoldedMatch checks the rendered line: the match
// is wrapped and nothing around it is written unescaped.
func TestHighlightEscapesAroundAFoldedMatch(t *testing.T) {
	html := string(components.Highlight(components.HighlightProps{Text: "KKab<script>", Query: "AB"}))
	if !strings.Contains(html, `<mark data-part="match" class="highlight-match">ab</mark>&lt;script&gt;`) {
		t.Fatalf("the match is not marked or the rest is not escaped:\n%s", html)
	}
	if !utf8.ValidString(html) {
		t.Fatalf("the rendered line is not valid UTF-8: %q", html)
	}
}

func assertSegmentsRebuildText(t *testing.T, text string, segments []components.HighlightSegment) {
	t.Helper()
	var rebuilt strings.Builder
	for _, segment := range segments {
		if segment.Text == "" {
			t.Errorf("an empty segment in %+v", segments)
		}
		if utf8.ValidString(text) && !utf8.ValidString(segment.Text) {
			t.Errorf("segment %q splits a character of %q", segment.Text, text)
		}
		rebuilt.WriteString(segment.Text)
	}
	if rebuilt.String() != text {
		t.Errorf("segments %+v do not add up to %q", segments, text)
	}
}
