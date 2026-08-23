//go:build kyse

package components

import "html/template"

@go
// PopoverProps is a button and the panel it shows beside itself.
//
// It is the plain form of the pair: the trigger says what it controls and
// whether the panel is open, and the panel says whether it is hidden. Escape
// closes it and gives focus back to the trigger, a click outside closes it, and
// opening one closes any other -- so two of them on a page cannot both stand
// open over each other.
//
// Where it sits is Side and Align, and both are the stylesheet's: the panel is
// positioned against the trigger by CSS rather than by measuring anything, so
// it lands correctly on the first paint and stays there while the page scrolls.
//
// It is not a dialog. There is no backdrop, nothing behind it is made inert,
// and focus is free to leave -- which is what separates a panel with a filter
// in it from something that has to be answered before the page continues.
type PopoverProps struct {
	// ID is what the trigger and the panel hang their ids off. Two popovers on
	// a page need two.
	ID string
	// Label is the text on the trigger.
	Label string
	// Variant and Size style the trigger. Empty is the default of each.
	Variant string
	Size    string
	// Title is the heading inside the panel. Empty draws no header.
	Title string
	// Description is the sentence under the heading.
	Description string
	// Content is what the panel holds under the header. It is markup, so it
	// comes from another component or from a view.
	Content template.HTML
	// Side is which way the panel opens: "top", "bottom", "left", "right", or
	// the logical "inline-start" and "inline-end", which turn around in a
	// right-to-left document. Empty opens downward.
	Side string
	// Align is where the panel sits along that side: "start", "end" or
	// "center". Empty aligns it to the start.
	Align string
}

// TriggerID is the id of the button that opens the panel.
func (p PopoverProps) TriggerID() string { return p.ID + "-trigger" }

// PanelID is the id of the panel it opens.
func (p PopoverProps) PanelID() string { return p.ID + "-popover" }
@endgo

<div class="popover" id="{{ .ID }}">
	<button
		type="button"
		class="btn"
		id="{{ .TriggerID() }}"
		aria-controls="{{ .PanelID() }}"
		aria-expanded="false"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
	>{{ .Label }}</button>

	<div
		id="{{ .PanelID() }}"
		data-popover
		aria-hidden="true"
		@if(.Side != "")
			data-side="{{ .Side }}"
		@endif
		@if(.Align != "")
			data-align="{{ .Align }}"
		@endif
	>
		@if(.Title != "")
			<header>
				<h3>{{ .Title }}</h3>
				@if(.Description != "")
					<p>{{ .Description }}</p>
				@endif
			</header>
		@endif
		{!! .Content !!}
	</div>
</div>
