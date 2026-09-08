//go:build kyse

package components

@go
// StatusProps is a line the page updates and a screen reader reads without
// being asked.
//
// It is the polite half of the pair a toast completes: a toast is something
// that happened and then goes away, and this is a line that stays and changes
// -- "saved", "3 of 40 uploaded", "connection lost". Both announce; only one
// of them disappears.
//
// The region has to be in the document before the text arrives. A region
// created and filled in the same swap is generally not announced at all,
// because there was nothing there to be watching. So this is drawn empty, on
// the page, and the endpoint swaps its contents.
//
// It publishes root.
type StatusProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Message is the current line. Empty draws the region with nothing in it,
	// which is how it is put on the page to be filled later.
	Message string
	// Log makes it a running list rather than one line: each addition is read,
	// and what came before is not read again. It is right for a transcript and
	// wrong for a single value.
	Log bool
	// Tone is "success", "destructive" or "muted". Empty is the plain one.
	Tone string
	// Label names the region for a reader that lists them.
	Label string
}

// Role is what kind of region this is.
func (p StatusProps) Role() string {
	if p.Log {
		return "log"
	}
	return "status"
}

// Atomic is whether the whole region is read on every change, or only what was
// added. A single value is read whole; a running list is not, or every line
// would be read again on every line.
func (p StatusProps) Atomic() string {
	if p.Log {
		return "false"
	}
	return "true"
}

// PartNames are the parts this component publishes.
func (p StatusProps) PartNames() []string { return []string{"root"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("status") }}"
	role="{{ .Role() }}"
	aria-live="polite"
	aria-atomic="{{ .Atomic() }}"
	@if(.Tone != "")
		data-tone="{{ .Tone }}"
	@endif
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@attributes(.RootAttrs())
>{{ .Message }}</div>
