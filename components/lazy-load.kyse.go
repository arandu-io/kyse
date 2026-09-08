//go:build kyse

package components

@go
// LazyLoadProps is a region that fetches its own contents after the page has
// drawn.
//
// It is for the part of a page that is slow and is not what the page is for: a
// panel of aggregates, a chart over a year, a third-party embed. Waiting for
// it means nothing renders until the slowest query finishes; deferring it
// means the page arrives and the panel fills in.
//
// The space is reserved before the answer comes back. Without that, the
// content arriving pushes everything under it down -- which is the layout
// shift a reader experiences as the page moving while they are reading it. So
// a placeholder of the right height is drawn and swapped for the real thing.
//
// It publishes root and placeholder.
type LazyLoadProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// URL is what answers with the contents.
	URL string
	// Label names the region while it is empty, and is what a screen reader
	// says is loading. Without it the region announces as a busy nothing.
	Label string
	// Deferred waits until the region is scrolled towards instead of fetching
	// as soon as the page has loaded. It is right for something far down and
	// wrong for something above the fold, which would then load only after the
	// reader has already looked at where it should be.
	Deferred bool
	// Message is the line shown in the placeholder: "Loading the summary".
	// Empty draws the placeholder with no words, which is right when the
	// placeholder is shaped like what is coming.
	Message string
}

// Trigger is when the fetch happens: once the page has loaded, or once the
// region is approached. Once either way -- a region that re-fetched on every
// scroll past would replace what a reader is looking at.
func (p LazyLoadProps) Trigger() string {
	if p.Deferred {
		return "revealed once"
	}
	return "load once"
}

// PartNames are the parts this component publishes.
func (p LazyLoadProps) PartNames() []string { return []string{"root", "placeholder"} }
@endgo

{{-- aria-busy is true while the region is empty and goes away with the swap,
     because the swap replaces the whole element. That is the honest sequence:
     the region says it is working, and then it is not that region any more. --}}
<div
	data-part="root"
	class="{{ .RootClass("lazy-load") }}"
	role="region"
	aria-label="{{ .Label }}"
	aria-busy="true"
	@attributes(.RootAttrs())
	hx-get="{{ .URL }}"
	hx-trigger="{{ .Trigger() }}"
	hx-target="this"
	hx-swap="outerHTML"
>
	<div
		data-part="placeholder"
		class="{{ .PartClass("placeholder", "skeleton") }}"
		@attributes(.PartAttrs("placeholder"))
	>{{ .Message }}</div>
</div>
