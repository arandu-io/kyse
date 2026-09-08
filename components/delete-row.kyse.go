//go:build kyse

package components

@go
// DeleteRowProps is the button that removes one row, in place.
//
// It is a DELETE and not a POST, because the request says what it does and a
// crawler, a proxy and a browser's own retry all read that. It asks first, and
// it asks with the browser's own confirmation rather than with a dialog drawn
// here: a page that draws one has to also decide where it opens, what closes
// it and what happens when the row underneath is swapped away while it is
// open.
//
// The row is removed after the server says it is gone, never before. An
// optimistic removal is a row that vanishes and a record that did not, and the
// person who did it has already moved on.
//
// It publishes root.
type DeleteRowProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// URL is what is deleted.
	URL string
	// Target is what is removed when it answers: the row, usually as
	// "closest tr" or the id of the element. Empty removes the button's own
	// closest list item or row.
	Target string
	// Label is the text on the button. Empty says "Delete".
	Label string
	// Confirm is the question asked before the request. Empty asks "Delete
	// this?" -- and it is never skipped, because a delete with no confirmation
	// beside a row is a delete that happens on a mis-click.
	Confirm string
	// Description names what is being deleted, for a button drawn once per row
	// where "Delete" alone is the same name repeated down the page: "Delete
	// invoice 2026-114".
	Description string
	// Variant and Size style it. Empty is destructive, because that is what
	// this is.
	Variant string
	Size    string
	// Icon is drawn before the label.
	Icon template.HTML
	// IconOnly draws the name as the accessible name rather than as text.
	IconOnly bool
	// Disabled draws it unavailable.
	Disabled bool
}

// Text is the label, or the word used when none was given.
func (p DeleteRowProps) Text() string {
	if p.Label != "" {
		return p.Label
	}
	return "Delete"
}

// Name is what the button is called to a screen reader: the description when
// there is one, and the label otherwise.
func (p DeleteRowProps) Name() string {
	if p.Description != "" {
		return p.Description
	}
	return p.Text()
}

// Question is what is asked before the request.
func (p DeleteRowProps) Question() string {
	if p.Confirm != "" {
		return p.Confirm
	}
	return "Delete this?"
}

// Removes is what is taken out of the page when the server answers.
func (p DeleteRowProps) Removes() string {
	if p.Target != "" {
		return p.Target
	}
	return "closest tr, closest li"
}

// Look is the variant, and it is destructive unless the caller says otherwise.
func (p DeleteRowProps) Look() string {
	if p.Variant != "" {
		return p.Variant
	}
	return "destructive"
}

// PartNames are the parts this component publishes.
func (p DeleteRowProps) PartNames() []string { return []string{"root"} }
@endgo

{{-- hx-swap outerHTML with the target set to the row is what removes it: the
     endpoint answers with nothing, and nothing replacing the row is the row
     gone. A 204 answers the same way and is the honest status for it. --}}
<button
	data-part="root"
	class="{{ .RootClass("btn") }}"
	type="button"
	data-variant="{{ .Look() }}"
	@attributes(.RootAttrs())
	hx-delete="{{ .URL }}"
	hx-confirm="{{ .Question() }}"
	hx-target="{{ .Removes() }}"
	hx-swap="outerHTML swap:200ms"
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
	@if(.Description != "")
		aria-label="{{ .Name() }}"
	@endif
	@if(.IconOnly && .Description == "")
		aria-label="{{ .Text() }}"
	@endif
	@if(.Disabled)
		disabled
	@endif
>@if(.Icon != ""){!! .Icon !!}@endif@if(!.IconOnly){{ .Text() }}@endif</button>
