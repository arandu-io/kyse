//go:build kyse

package components

@go
// FieldProps is one labelled input, with its error and its hint.
//
// It is the component that pays for itself fastest. A form field written out is
// a label, an input, a conditional paragraph for the error and the two ARIA
// attributes that tie them together -- eleven lines that are identical on every
// field of every form, and the two attributes are the ones people forget.
//
// A message being there does three things at once: it draws the sentence, marks
// the input invalid for the stylesheet, and points aria-describedby at the
// sentence so a screen reader announces it. That is why it is one answer and not
// three.
//
// # Where the message and the typed value come from
//
// From Page, by Name, and not from two more props beside it. See the Page
// interface for the failure that shape prevents: a name written three times in
// one call is a name two of the three can disagree with, silently, on the screen
// where a missing message is the complaint.
type FieldProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, the id everything else is hung off, and what
	// Page is asked about.
	Name string
	// Label is the text above the input.
	Label string
	// Type is the input type: "text", "email", "password", "date", "url".
	// Empty means "text".
	Type string
	// Value is what the input starts with when nothing was rejected: the stored
	// record on an edit form, and empty on a create form. What was typed on a
	// rejected attempt takes precedence over it -- see Current.
	Value string
	// Placeholder is the grey text inside an empty input. It is not a label:
	// it disappears as soon as somebody types.
	Placeholder string
	// Hint is the sentence under the input that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this input asks for its
	// message and for what was typed. Nil draws no message and keeps Value.
	Page Page
	// Required marks the input required, in the markup and to the browser.
	Required bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
	// Autocomplete is the browser hint: "email", "current-password",
	// "new-password", "off".
	Autocomplete string
}

// InputType is the type attribute, defaulting to text.
func (p FieldProps) InputType() string {
	if p.Type == "" {
		return "text"
	}
	return p.Type
}

// Message is what validation left for this input, or empty.
func (p FieldProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is what the input is drawn with: what was typed on the rejected
// attempt, and Value when there was none.
//
// A password is never among what was typed -- the flash drops the value and
// keeps the message -- so this answers empty for one, which is the intended
// behaviour and not an omission: the box comes back blank and the sentence under
// it says why.
func (p FieldProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// DescribedBy names the element that explains this input, so the error and the
// hint are announced rather than only shown.
func (p FieldProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
// PartNames are the parts this component publishes.
func (p FieldProps) PartNames() []string { return []string{"root", "label", "input", "message", "hint"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("field") }}"
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		for="{{ .Name }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>
	<input
		data-part="input"
		class="{{ .PartClass("input", "input") }}"
		type="{{ .InputType() }}"
		id="{{ .Name }}"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		@attributes(.PartAttrs("input"))
		@if(.Placeholder != "")
			placeholder="{{ .Placeholder }}"
		@endif
		@if(.Autocomplete != "")
			autocomplete="{{ .Autocomplete }}"
		@endif
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			required
		@endif
		@if(.Autofocus)
			autofocus
		@endif
	>
	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .Name }}-error"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p
				data-part="hint"
				id="{{ .Name }}-hint"
				class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("hint"))
			>{{ .Hint }}</p>
		@endif
	@endif
</div>
