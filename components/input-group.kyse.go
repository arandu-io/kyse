//go:build kyse

package components

import "strings"

@go
// InputGroupProps is one labelled input with something written inside the box
// beside it: the currency it is in, the unit it is measured in, the domain a
// name is under, the count of what was found.
//
// It is a field whose border is drawn around the group rather than around the
// control. That is the whole difference, and it is the reason the addon has to
// be inside: a symbol placed next to a bordered box is a symbol that sits
// outside it, and every attempt to fake it with padding puts the two out of
// line at the first change of font size.
//
// # The addons are read out
//
// Start and End are named and pointed at by aria-describedby, so a box worth
// "USD" is announced as a box worth USD. An addon that is only drawn is an
// addon half the people filling in the form never learn about, and it is
// usually the one carrying the unit the number has to be in.
//
// # Where the message and the typed value come from
//
// From Page, by Name, in the same way a field asks. See the Page interface.
type InputGroupProps struct {
	// Name is the form field name, the id everything else is hung off, and
	// what Page is asked about.
	Name string
	// Label is the text above the group.
	Label string
	// Type is the input type: "text", "email", "number", "search", "url".
	// Empty means "text".
	Type string
	// Value is what the input starts with when nothing was rejected: the
	// stored record on an edit form, and empty on a create form.
	Value string
	// Placeholder is the grey text inside an empty input.
	Placeholder string
	// Start is what is written inside the box before the control: a currency
	// symbol, a protocol, an at sign. Empty draws none.
	Start string
	// End is what is written inside the box after the control: a unit, a
	// domain, a count of what was found. Empty draws none.
	End string
	// Shortcut is the key combination that puts the cursor in this box, drawn
	// as a key cap at the end. It is shown and never bound: what the keys do
	// is the application's.
	Shortcut string
	// Hint is the sentence under the group that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this input asks for
	// its message and for what was typed. Nil draws no message and keeps
	// Value.
	Page Page
	// Required marks the input required, in the markup and to the browser.
	Required bool
	// Disabled draws the whole group as unavailable and stops it answering.
	Disabled bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
	// Autocomplete is the browser hint: "email", "current-password",
	// "new-password", "off".
	Autocomplete string
}

// InputType is the type attribute, defaulting to text.
func (p InputGroupProps) InputType() string {
	if p.Type == "" {
		return "text"
	}
	return p.Type
}

// StartID names the addon before the control.
func (p InputGroupProps) StartID() string {
	return p.Name + "-start"
}

// EndID names the addon after the control.
func (p InputGroupProps) EndID() string {
	return p.Name + "-end"
}

// Message is what validation left for this input, or empty.
func (p InputGroupProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is what the input is drawn with: what was typed on the rejected
// attempt, and Value when there was none.
func (p InputGroupProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// DescribedBy names everything that explains this input, in the order it
// should be read: what is written inside the box, then the message or the
// hint.
//
// It is a list rather than one id, because an input in a group has more than
// one thing attached to it and aria-describedby is the attribute that takes
// several. Naming only the message would drop the unit; naming only the unit
// would drop the complaint.
func (p InputGroupProps) DescribedBy() string {
	var ids []string
	if p.Start != "" {
		ids = append(ids, p.StartID())
	}
	if p.End != "" {
		ids = append(ids, p.EndID())
	}
	switch {
	case p.Message() != "":
		ids = append(ids, p.Name+"-error")
	case p.Hint != "":
		ids = append(ids, p.Name+"-hint")
	}
	return strings.Join(ids, " ")
}
@endgo

<div class="field">
	<label class="label" for="{{ .Name }}">{{ .Label }}</label>

	<div class="input-group">
		<input
			type="{{ .InputType() }}"
			id="{{ .Name }}"
			name="{{ .Name }}"
			value="{{ .Current() }}"
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
			@if(.Disabled)
				disabled
			@endif
			@if(.Autofocus)
				autofocus
			@endif
		>

		@if(.Start != "")
			<span id="{{ .StartID() }}" data-align="start">{{ .Start }}</span>
		@endif
		@if(.End != "")
			<span id="{{ .EndID() }}" data-align="end">{{ .End }}</span>
		@endif
		@if(.Shortcut != "")
			<span data-align="end"><kbd>{{ .Shortcut }}</kbd></span>
		@endif
	</div>

	@if(.Message() != "")
		<p id="{{ .Name }}-error" class="text-destructive text-sm">{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p id="{{ .Name }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
		@endif
	@endif
</div>
