//go:build kyse

package components

@go
// LinkProps is a link, with the two things a link gets wrong on its own.
//
// The first is the rel: a link that opens a new tab hands the opened page a
// reference back to this one through window.opener, and that reference can
// navigate it. Writing rel by hand is remembering to write it every time, and
// the times it is forgotten are the ones nobody sees. Here External sets both
// target and rel together, because one without the other is the bug.
//
// The second is the underline. A link that is only a colour is a link nobody
// with a colour deficiency can find, so the underline is the default and
// removing it is a named variant rather than an omission.
//
// It publishes root.
type LinkProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label is the text of the link.
	Label string
	// URL is where it goes. Empty draws the text without an anchor, which is
	// how the current page appears in a list of links to its siblings.
	URL string
	// Variant is "muted" for a link that should recede, or "hover" for one
	// underlined only under the pointer. Empty is underlined always.
	Variant string
	// External opens the link in a new tab and sets the rel that keeps the
	// opened page from reaching back into this one.
	External bool
	// Current marks this as the link to the page being shown, which is what a
	// screen reader announces to say "you are here".
	Current bool
	// Icon is drawn after the label: an arrow, a mark saying the link leaves.
	Icon template.HTML
}

// Target is the browsing context the link opens in, and is empty for a link
// that stays.
func (p LinkProps) Target() string {
	if p.External {
		return "_blank"
	}
	return ""
}

// Rel is what the opened page is allowed to know about this one. It is set
// only when a new tab is opened, because that is the only case where the
// opened page has a reference to reach back through.
func (p LinkProps) Rel() string {
	if p.External {
		return "noopener noreferrer"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p LinkProps) PartNames() []string { return []string{"root"} }
@endgo

@if(.URL != "")
	<a
		data-part="root"
		class="{{ .RootClass("link") }}"
		@attributes(.RootAttrs())
		href="{{ .URL }}"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Target() != "")
			target="{{ .Target() }}"
			rel="{{ .Rel() }}"
			data-external="true"
		@endif
		@if(.Current)
			aria-current="page"
		@endif
	>{{ .Label }}@if(.Icon != ""){!! .Icon !!}@endif</a>
@endif
@if(.URL == "")
	<span
		data-part="root"
		class="{{ .RootClass("link") }}"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Current)
			aria-current="page"
		@endif
		@attributes(.RootAttrs())
	>{{ .Label }}@if(.Icon != ""){!! .Icon !!}@endif</span>
@endif
