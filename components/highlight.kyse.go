//go:build kyse

package components

import (
	"html/template"
	"strings"

	"github.com/arandu-io/hesape/view"
)

@go
// HighlightProps is a line of text with the part somebody searched for marked
// in it.
//
// The marking is <mark>, which is the element that means "relevant to what is
// going on right now" rather than "important" -- and that is exactly what a
// search hit is. A span with a yellow background says nothing to anyone not
// looking at it.
//
// The match is found here, on the server, because the results are rendered
// here. A script that marked them in the browser would have to run again after
// every swap, and would be marking text it did not fetch.
//
// It publishes root and match.
type HighlightProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Text is the line, as it is stored.
	Text string
	// Query is what to mark in it. Empty marks nothing and draws the line
	// unchanged, which is what a results page does before anyone has typed.
	Query string
	// CaseSensitive matches the query exactly as written. It is off by
	// default, because a person searching for "ada" means the same thing as
	// one searching for "Ada".
	CaseSensitive bool
}

// HighlightSegment is one piece of the line, and whether it matched.
type HighlightSegment struct {
	// Text is the piece.
	Text string
	// Match is true when this piece is what was searched for.
	Match bool
}

// Segments is the text cut at the matches: each piece, and whether that piece
// is one of them.
func (p HighlightProps) Segments() []HighlightSegment {
	if p.Text == "" {
		return nil
	}
	if p.Query == "" {
		return []HighlightSegment{{Text: p.Text}}
	}
	haystack, needle := p.Text, p.Query
	if !p.CaseSensitive {
		haystack, needle = strings.ToLower(haystack), strings.ToLower(needle)
	}
	segments := make([]HighlightSegment, 0, 3)
	for cursor := 0; cursor < len(p.Text); {
		at := strings.Index(haystack[cursor:], needle)
		if at < 0 {
			segments = append(segments, HighlightSegment{Text: p.Text[cursor:]})
			break
		}
		start := cursor + at
		if start > cursor {
			segments = append(segments, HighlightSegment{Text: p.Text[cursor:start]})
		}
		end := start + len(needle)
		segments = append(segments, HighlightSegment{Text: p.Text[start:end], Match: true})
		cursor = end
	}
	return segments
}

// Marked is the line with the matches wrapped, assembled here rather than in
// the template.
//
// The template is written with a line per element, and a line per element puts
// a newline between them. Everywhere else that is invisible, because the
// elements are blocks; here the pieces are halves of one word, and a newline
// between them is a space inside it -- "Lovelace" drawn as "Love lace". So the
// splice is made where no newline exists to be introduced.
//
// Every piece is escaped as it is joined, which is the same guarantee the
// template gives and the reason this returns assembled HTML rather than taking
// any.
func (p HighlightProps) Marked() template.HTML {
	segments := p.Segments()
	if len(segments) == 0 {
		return ""
	}
	open := `<mark data-part="match"`
	if class := p.PartClass("match", "highlight-match"); class != "" {
		open += ` class="` + view.TextAttr(class) + `"`
	}
	if attrs, err := view.Attributes(p.PartAttrs("match")); err == nil {
		open += attrs
	}
	open += ">"

	var b strings.Builder
	for _, segment := range segments {
		escaped := template.HTMLEscapeString(view.Text(segment.Text))
		if !segment.Match {
			b.WriteString(escaped)
			continue
		}
		b.WriteString(open)
		b.WriteString(escaped)
		b.WriteString("</mark>")
	}
	return template.HTML(b.String())
}

// PartNames are the parts this component publishes.
func (p HighlightProps) PartNames() []string { return []string{"root", "match"} }
@endgo

<span
	data-part="root"
	class="{{ .RootClass("highlight") }}"
	@attributes(.RootAttrs())
>{!! .Marked() !!}</span>
