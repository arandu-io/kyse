//go:build kyse

package components

import "strings"

@go
// AvatarProps is who somebody is, in a circle.
//
// With no image it draws initials, which is the case that actually happens: an
// application where everybody uploaded a photograph does not exist, and a broken
// image icon beside every comment is worse than two letters.
// It publishes one part, root, which is the circle itself.
type AvatarProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the person. It is the alt text, and the initials come from it.
	Name string
	// ImageURL is their picture. Empty draws the initials.
	ImageURL string
	// Size is "sm", "lg", or empty for the default.
	Size string
}

// Initials are the first letters of the first and last word of the name.
//
// First and LAST, not the first two: "Ada B. Lovelace" is AL, and the version
// that walked forwards until it had two letters answered AB -- a middle initial
// where a surname belongs. Nobody would have noticed from looking at the
// component; the letters are the right shape either way.
//
// It walks runes rather than bytes: a name starting with É is one letter to a
// person and two bytes to Go, and slicing bytes would print half a character.
func (p AvatarProps) Initials() string {
	words := strings.Fields(p.Name)
	if len(words) == 0 {
		return ""
	}

	out := []rune{[]rune(words[0])[0]}
	if len(words) > 1 {
		out = append(out, []rune(words[len(words)-1])[0])
	}
	return strings.ToUpper(string(out))
}

// PartNames are the parts this component publishes.
func (p AvatarProps) PartNames() []string { return []string{"root"} }
@endgo

<span
	data-part="root"
	class="{{ .RootClass("avatar") }}"
	@attributes(.RootAttrs())
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
>
	@if(.ImageURL != "")
		<img src="{{ .ImageURL }}" alt="{{ .Name }}">
	@endif
	@if(.ImageURL == "")
		<span aria-hidden="true">{{ .Initials() }}</span>
		<span class="sr-only">{{ .Name }}</span>
	@endif
</span>
