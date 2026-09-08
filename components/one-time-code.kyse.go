//go:build kyse

package components

@go
// OneTimeCodeBehavior is the name the client behaviour is registered under with
// arandu.ui.define, and the name this component writes.
const OneTimeCodeBehavior = "one-time-code"

// OneTimeCodeProps is the box somebody types a one-time code into: one square per
// digit, filled left to right.
//
// It draws one visible input per digit and a hidden input carrying the whole
// code, which is what submits. So the server reads one field called what the
// caller named, and never has to join six of them.
//
// # Why squares and not one box
//
// A code arrives in a message and is read a few characters at a time. Six
// squares show how many are expected before anything is typed, show which one
// is being answered, and let somebody fix the third without retyping the rest.
// One text box shows none of that, and a code half typed into it looks the same
// as a code fully typed into it.
//
// # What it does not do
//
// It does not decide. Whether the code is the right one is the server's answer,
// and this component is not consulted -- the same division the password and
// mask behaviours keep. data-otp-complete says every square is filled, which is
// a count and not a verdict.
//
// # Pasting
//
// A code pasted into any square fills all of them, because that is what somebody
// copying six digits out of a message means. A code longer than the field is cut
// to length rather than refused.
type OneTimeCodeProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, carried by the hidden input that submits.
	Name string
	// ID is the id of the first square, so a Label can point at it. Empty uses
	// Name.
	ID string
	// Length is how many characters the code has. Zero means six, which is what
	// nearly every code is.
	Length int
	// Value is the code this starts with, which is empty on an ordinary load
	// and is what was typed on a rejected attempt.
	Value string
	// Page is the screen's own view.Page, asked for what was typed and for
	// whether it was rejected.
	Page Page
	// Label is the accessible name of the group. A code field with no label is
	// announced as a row of unnamed boxes.
	Label string
	// DescribedBy is the id of the element that explains this field.
	DescribedBy string
	// Alphanumeric takes letters as well as digits. The default is digits
	// alone, and it is what sets the on-screen keyboard to a number pad.
	Alphanumeric bool
	// Disabled takes the field out of the form and out of the tab order.
	Disabled bool
	// Autofocus puts the cursor in the first square on load.
	Autofocus bool
}

// ElementID is the id the first square carries.
func (p OneTimeCodeProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// Size is how many squares this field draws.
//
// Six when the caller named none, because six is what nearly every code is, and
// a field that drew none until told would be a field that renders as nothing.
func (p OneTimeCodeProps) Size() int {
	if p.Length > 0 {
		return p.Length
	}
	return 6
}

// Squares is one entry per character, so the view can walk them.
func (p OneTimeCodeProps) Squares() []int {
	out := make([]int, 0, p.Size())
	for i := 0; i < p.Size(); i++ {
		out = append(out, i)
	}
	return out
}

// Current is the code this field is drawn with: what was typed on a rejected
// attempt, and Value when there was none.
func (p OneTimeCodeProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// CharAt is the character in one square, and empty past the end of the value.
func (p OneTimeCodeProps) CharAt(at int) string {
	runes := []rune(p.Current())
	if at < 0 || at >= len(runes) {
		return ""
	}
	return string(runes[at])
}

// InputMode is the on-screen keyboard this field asks for.
func (p OneTimeCodeProps) InputMode() string {
	if p.Alphanumeric {
		return "text"
	}
	return "numeric"
}

// AutocompleteAt is the autofill hint for one square.
//
// Only the first asks for the code, because a browser filling every square with
// the whole code is what happens when they all ask. The rest are told not to
// autofill at all, and the behaviour spreads what arrives in the first across
// them.
func (p OneTimeCodeProps) AutocompleteAt(at int) string {
	if at == 0 {
		return "one-time-code"
	}
	return "off"
}

// Message is the rejection for this field, or empty.
func (p OneTimeCodeProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p OneTimeCodeProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: OneTimeCodeBehavior, Props: map[string]any{
			"length":       p.Size(),
			"alphanumeric": p.Alphanumeric,
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p OneTimeCodeProps) PartNames() []string { return []string{"root", "square", "raw"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("otp") }}"
	role="group"
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@if(.DescribedBy != "")
		aria-describedby="{{ .DescribedBy }}"
	@endif
	@if(.Message() != "")
		aria-invalid="true"
	@endif
	@attributes(.RootAttrs())
>
	@foreach(.Squares() as index)
		<input
			data-part="square"
			data-code-index="{{ index }}"
			class="{{ .PartClass("square", "otp-square") }}"
			type="text"
			@if(index == 0)
				id="{{ .ElementID() }}"
			@endif
			inputmode="{{ .InputMode() }}"
			maxlength="1"
			value="{{ .CharAt(index) }}"
			autocomplete="{{ .AutocompleteAt(index) }}"
			@if(.Disabled)
				disabled
			@endif
			@if(.Autofocus && index == 0)
				autofocus
			@endif
		>
	@endforeach
	<input
		data-part="raw"
		class="{{ .PartClass("raw", "") }}"
		type="hidden"
		name="{{ .Name }}"
		value="{{ .Current() }}"
	>
</div>
