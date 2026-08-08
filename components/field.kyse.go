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
// Error being set does three things at once: it draws the message, marks the
// input invalid for the stylesheet, and points aria-describedby at the message
// so a screen reader announces it. That is why it is one field and not three.
type FieldProps struct {
	// Name is the form field name, and the id everything else is hung off.
	Name string
	// Label is the text above the input.
	Label string
	// Type is the input type: "text", "email", "password", "date", "url".
	// Empty means "text".
	Type string
	// Value is what the input starts with. It is what makes a rejected form come
	// back filled in rather than blank.
	Value string
	// Placeholder is the grey text inside an empty input. It is not a label:
	// it disappears as soon as somebody types.
	Placeholder string
	// Hint is the sentence under the input that is always there.
	Hint string
	// Error is what came back from validation. Empty means the field is fine.
	Error string
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

// DescribedBy names the element that explains this input, so the error and the
// hint are announced rather than only shown.
func (p FieldProps) DescribedBy() string {
	if p.Error != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
@endgo

<div class="field">
	<label class="label" for="{{ .Name }}">{{ .Label }}</label>
	<input
		class="input"
		type="{{ .InputType() }}"
		id="{{ .Name }}"
		name="{{ .Name }}"
		value="{{ .Value }}"
		@if(.Placeholder != "")
			placeholder="{{ .Placeholder }}"
		@endif
		@if(.Autocomplete != "")
			autocomplete="{{ .Autocomplete }}"
		@endif
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Error != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			required
		@endif
		@if(.Autofocus)
			autofocus
		@endif
	>
	@if(.Error != "")
		<p id="{{ .Name }}-error" class="text-destructive text-sm">{{ .Error }}</p>
	@endif
	@if(.Error == "")
		@if(.Hint != "")
			<p id="{{ .Name }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
		@endif
	@endif
</div>
