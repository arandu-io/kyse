//go:build kyse

package components

import "html/template"

@go
// DialogTriggerProps is the control that opens a dialog: a FormDialog, a
// Dialog or a Drawer, named by its ID.
//
// With Href it is a link to the page that holds the same thing the dialog
// does -- the form as a page of its own. The script the framework serves opens
// the dialog on an ordinary click instead of following the link, and leaves a
// click that asks for a new tab to the browser. A page that never runs the
// script follows the link, and the form is posted and refused there like any
// other: nothing about the dialog has to work for the form to.
//
// Without Href it is a button, for a dialog that has no page behind it -- a
// confirmation, which a page without the script cannot ask anyway.
//
// It is drawn as a button either way, and publishes root and icon.
type DialogTriggerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Dialog is the ID of the dialog it opens.
	Dialog string
	// Label is the text on it.
	Label string
	// Href is the page that holds the same form, followed when the dialog
	// cannot be opened. Empty draws a button.
	Href string
	// Variant and Size style it, as on a Button. Empty is the default of each.
	Variant string
	Size    string
	// Icon is drawn before the label: the output of an icon function.
	Icon template.HTML
	// IconOnly draws the label as the accessible name rather than as text.
	IconOnly bool
}

// PartNames are the parts this component publishes.
func (p DialogTriggerProps) PartNames() []string { return []string{"root", "icon"} }
@endgo

@if(.Href != "")
	<a
		data-part="root"
		href="{{ .Href }}"
		class="{{ .RootClass("btn") }}"
		data-dialog-open="{{ .Dialog }}"
		aria-haspopup="dialog"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
		@if(.IconOnly)
			aria-label="{{ .Label }}"
		@endif
		@attributes(.RootAttrs())
	>
		@if(.Icon != "")
			<span
				data-part="icon"
				@if(.PartClass("icon") != "")
					class="{{ .PartClass("icon") }}"
				@endif
				aria-hidden="true"
				@attributes(.PartAttrs("icon"))
			>{!! .Icon !!}</span>
		@endif
		@if(!.IconOnly)
			{{ .Label }}
		@endif
	</a>
@else
	<button
		data-part="root"
		type="button"
		class="{{ .RootClass("btn") }}"
		data-dialog-open="{{ .Dialog }}"
		aria-haspopup="dialog"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
		@if(.IconOnly)
			aria-label="{{ .Label }}"
		@endif
		@attributes(.RootAttrs())
	>
		@if(.Icon != "")
			<span
				data-part="icon"
				@if(.PartClass("icon") != "")
					class="{{ .PartClass("icon") }}"
				@endif
				aria-hidden="true"
				@attributes(.PartAttrs("icon"))
			>{!! .Icon !!}</span>
		@endif
		@if(!.IconOnly)
			{{ .Label }}
		@endif
	</button>
@endif
