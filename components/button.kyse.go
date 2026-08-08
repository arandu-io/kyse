//go:build kyse

package components

@go
// ButtonProps is what a button is drawn from.
//
// Variant and Size are the shadcn attribute API, unchanged: "default",
// "secondary", "outline", "ghost", "destructive", "link" for the first, and
// "default", "xs", "sm", "lg", "icon" for the second. The stylesheet reads them
// off the element, so the markup carries two short attributes instead of the
// forty utility classes the same button costs when it is written by hand.
//
// An empty Variant or Size means the default, and the stylesheet says so: the
// rules are written for `.btn:not([data-variant])` as well as for the named one.
type ButtonProps struct {
	// Label is the text on the button.
	Label string
	// Variant and Size select the look. Empty is the default of each.
	Variant string
	Size    string
	// Type is the HTML type: "button", "submit" or "reset". Empty means
	// "button", because a button inside a form that does not say otherwise
	// submits it, and that surprise is worth one line to remove.
	Type string
	// Disabled draws the button as unavailable and stops it answering.
	Disabled bool

	// The HTMX attributes, written only when they carry something. An empty
	// hx-post is an attribute HTMX acts on: it would post to the current URL.
	HxPost    string
	HxGet     string
	HxTarget  string
	HxSwap    string
	HxConfirm string
}

// ButtonType is the type attribute, defaulting to a button that does nothing on
// its own.
func (p ButtonProps) ButtonType() string {
	if p.Type == "" {
		return "button"
	}
	return p.Type
}
@endgo

<button
	type="{{ .ButtonType() }}"
	class="btn"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
	@if(.Disabled)
		disabled
	@endif
	@if(.HxPost != "")
		hx-post="{{ .HxPost }}"
	@endif
	@if(.HxGet != "")
		hx-get="{{ .HxGet }}"
	@endif
	@if(.HxTarget != "")
		hx-target="{{ .HxTarget }}"
	@endif
	@if(.HxSwap != "")
		hx-swap="{{ .HxSwap }}"
	@endif
	@if(.HxConfirm != "")
		hx-confirm="{{ .HxConfirm }}"
	@endif
>{{ .Label }}</button>
