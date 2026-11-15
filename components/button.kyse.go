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
//
// It publishes one part, root, which is the button itself.
type ButtonProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
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
	// Name and Value are what a submit button sends with its form, which is
	// how one form with several actions tells the server which was pressed.
	// Both are written only when Name carries something.
	Name  string
	Value string

	// The YPY attributes, written only when they carry something. An empty
	// ypy-post is an attribute YPY acts on: it would post to the current URL.
	//
	// YpyTarget and YpySwap ask for a piece of the page back rather than a
	// document, so what answers the request has to be a view that draws that
	// piece and nothing around it -- one that extends no layout. Aimed at a
	// view that extends one, the reply is a whole document: the swap drops its
	// head, puts the body's children inside the target, and the page ends up
	// holding a second copy of its own header and navigation.
	//
	// The status the answer carries decides whether the swap happens at all.
	// A rejected form is answered 422, and 422 is not swapped unless the
	// document says it is, so a button that narrows a swap and never seems to
	// receive one is asking a question the page has not been configured to
	// accept.
	//
	// YpyConfirm is the question asked before the request. The view layer's
	// script draws it as a centred confirmation dialog (hesape v0.103.0 and
	// later) and never as the browser's own box; the button's Variant is the
	// variant of the button that goes through with it, so a destructive button
	// asks with a destructive one.
	YpyPost    string
	YpyGet     string
	YpyTarget  string
	YpySwap    string
	YpyConfirm string
}

// ButtonType is the type attribute, defaulting to a button that does nothing on
// its own.
func (p ButtonProps) ButtonType() string {
	if p.Type == "" {
		return "button"
	}
	return p.Type
}

// PartNames are the parts this component publishes.
func (p ButtonProps) PartNames() []string { return []string{"root"} }
@endgo

<button
	data-part="root"
	type="{{ .ButtonType() }}"
	class="{{ .RootClass("btn") }}"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
	@if(.Disabled)
		disabled
	@endif
	@if(.Name != "")
		name="{{ .Name }}"
		value="{{ .Value }}"
	@endif
	@if(.YpyPost != "")
		ypy-post="{{ .YpyPost }}"
	@endif
	@if(.YpyGet != "")
		ypy-get="{{ .YpyGet }}"
	@endif
	@if(.YpyTarget != "")
		ypy-target="{{ .YpyTarget }}"
	@endif
	@if(.YpySwap != "")
		ypy-swap="{{ .YpySwap }}"
	@endif
	@if(.YpyConfirm != "")
		ypy-confirm="{{ .YpyConfirm }}"
	@endif
	@attributes(.RootAttrs())
>{{ .Label }}</button>
