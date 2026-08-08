//go:build kyse

package components

import "strings"

@go
// AvatarProps is who somebody is, in a circle.
//
// With no image it draws initials, which is the case that actually happens: an
// application where everybody uploaded a photograph does not exist, and a broken
// image icon beside every comment is worse than two letters.
type AvatarProps struct {
	// Name is the person. It is the alt text, and the initials come from it.
	Name string
	// ImageURL is their picture. Empty draws the initials.
	ImageURL string
	// Size is "sm", "lg", or empty for the default.
	Size string
}

// Initials are the first letters of the first and last word of the name.
//
// It walks runes rather than bytes: a name starting with É is one letter to a
// person and two bytes to Go, and slicing bytes would print half a character.
func (p AvatarProps) Initials() string {
	var out []rune
	for i, word := range strings.Fields(p.Name) {
		if i > 0 && len(out) > 0 {
			out = out[:1]
		}
		out = append(out, []rune(word)[0])
		if len(out) == 2 {
			break
		}
	}
	return strings.ToUpper(string(out))
}
@endgo

<span class="avatar"
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
