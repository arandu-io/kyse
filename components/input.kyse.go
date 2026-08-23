//go:build kyse

package components

@go
// InputProps is the bare control: one input element, with no label around it
// and no paragraph under it.
//
// Field is the labelled form row -- label, input, message and the ARIA that ties
// the three together -- and it is what an ordinary form is written with. This
// one is for the places that row does not fit: a search box in a toolbar, a cell
// in an editable table, a box sitting beside a button in a group. There the name
// is given by an element the caller drew, or to assistive technology alone, and
// there is no room under the box for a sentence.
//
// Choosing this over Field means taking on the three things Field does without
// being asked. Label draws the name. The caller draws the message. DescribedBy
// is how the box is pointed at whichever of the two explains it. What is not
// given up is the value and the invalid state: both still come from Page, by
// Name, exactly as they do for Field.
type InputProps struct {
	// Name is the form field name, and what Page is asked about.
	Name string
	// ID is the id the markup carries, so a Label can point at it. Empty uses
	// Name, which is right until the same name is drawn twice on one screen.
	ID string
	// Type is the input type: "text", "email", "password", "search", "number",
	// "date", "url". Empty means "text".
	Type string
	// Value is what the input starts with when nothing was rejected: the stored
	// record on an edit form, and empty on a create form. What was typed on a
	// rejected attempt takes precedence over it -- see Current.
	Value string
	// Placeholder is the grey text inside an empty input. It is not a label:
	// it disappears as soon as somebody types.
	Placeholder string
	// Page is the screen's own view.Page, which is what this input asks for
	// what was typed and for whether it was rejected. Nil keeps Value and marks
	// nothing invalid.
	Page Page
	// AriaLabel is the accessible name for an input no label element points at.
	// An input with neither is announced as an unnamed edit box.
	AriaLabel string
	// DescribedBy is the id of the element that explains this input -- the
	// message, or the hint, whichever the caller drew. This component draws
	// neither, so nothing else can fill this in.
	DescribedBy string
	// Autocomplete is the browser hint: "email", "current-password",
	// "new-password", "off".
	Autocomplete string
	// Required marks the input required, in the markup and to the browser.
	Required bool
	// Disabled takes the input out of the form and out of the tab order.
	Disabled bool
	// Readonly refuses edits while keeping the value focusable, selectable and
	// submitted, which is what Disabled gives up.
	Readonly bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
}

// InputType is the type attribute, defaulting to text.
func (p InputProps) InputType() string {
	if p.Type == "" {
		return "text"
	}
	return p.Type
}

// ElementID is the id the input carries: ID when one was given, and Name when
// none was.
func (p InputProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// Message is what validation left for this input, or empty.
//
// It is not drawn here -- this component has nowhere to draw it. It decides
// aria-invalid, so the box is marked wrong for the stylesheet and for a screen
// reader even though the sentence explaining it belongs to the caller.
func (p InputProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is what the input is drawn with: what was typed on the rejected
// attempt, and Value when there was none.
//
// A password is never among what was typed, so this answers empty for one and
// the box comes back blank.
func (p InputProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}
@endgo

<input
	class="input"
	type="{{ .InputType() }}"
	id="{{ .ElementID() }}"
	name="{{ .Name }}"
	value="{{ .Current() }}"
	@if(.Placeholder != "")
		placeholder="{{ .Placeholder }}"
	@endif
	@if(.AriaLabel != "")
		aria-label="{{ .AriaLabel }}"
	@endif
	@if(.Autocomplete != "")
		autocomplete="{{ .Autocomplete }}"
	@endif
	@if(.DescribedBy != "")
		aria-describedby="{{ .DescribedBy }}"
	@endif
	@if(.Message() != "")
		aria-invalid="true"
	@endif
	@if(.Required)
		required
	@endif
	@if(.Disabled)
		disabled
	@endif
	@if(.Readonly)
		readonly
	@endif
	@if(.Autofocus)
		autofocus
	@endif
>
