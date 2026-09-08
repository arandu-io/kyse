//go:build kyse

package components

@go
// SkipLinkProps is the first thing in the tab order, and it is there for the
// people who reach the page by tabbing.
//
// Without it, every visit that starts at the keyboard begins by tabbing past
// the whole navigation to get to the page. With a header of thirty links that
// is thirty presses on every page, and it is the same thirty every time.
//
// It is hidden until it takes focus, so it costs a pointer user nothing and
// appears the moment it is the thing that would be used.
//
// It publishes root.
type SkipLinkProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label is what it says. Empty says "Skip to main content".
	Label string
	// TargetID is the id it jumps to, without the hash. Empty jumps to "main".
	//
	// The element it names has to be able to take focus for the jump to move
	// the keyboard and not only the scroll, which for a <main> means tabindex
	// of minus one on it.
	TargetID string
}

// Text is the label, or the sentence used when none was given.
func (p SkipLinkProps) Text() string {
	if p.Label != "" {
		return p.Label
	}
	return "Skip to main content"
}

// Href is the fragment it jumps to.
func (p SkipLinkProps) Href() string {
	if p.TargetID != "" {
		return "#" + p.TargetID
	}
	return "#main"
}

// PartNames are the parts this component publishes.
func (p SkipLinkProps) PartNames() []string { return []string{"root"} }
@endgo

<a
	data-part="root"
	class="{{ .RootClass("skip-link") }}"
	@attributes(.RootAttrs())
	href="{{ .Href() }}"
>{{ .Text() }}</a>
