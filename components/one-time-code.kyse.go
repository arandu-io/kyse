//go:build kyse

package components

import (
	"fmt"
	"strings"
)

@go
// OneTimeCodeBehavior is the name the client behaviour is registered under with
// arandu.ui.define, and the name this component writes.
const OneTimeCodeBehavior = "one-time-code"

// OneTimeCodeProps is the box somebody types a one-time code into: one square per
// digit, filled left to right.
//
// It draws one text input carrying the whole code under the name the caller
// gave, which is what submits, and one square per character beside it, hidden.
// So the server reads one field called what the caller named, and never has to
// join six of them.
//
// # Without the script
//
// The squares are an enhancement. The client behaviour shows them and turns the
// text input into the hidden field they keep in step, handing it the id a Label
// points at. Until it runs -- scripting off, the script blocked or not loaded
// yet -- the text input is the field: labelled, asking for the code by
// autocomplete, held to the code's length and alphabet by maxlength and
// pattern, and posted by the form as it stands. Squares that nothing joins
// would post nothing under the name the server reads, so they are drawn hidden
// and only the behaviour shows them.
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
// mask behaviours keep. The data-code-complete the behaviour writes on the root
// says every square is filled, which is a count and not a verdict.
//
// # Pasting
//
// A code pasted into any square fills all of them, because that is what somebody
// copying six digits out of a message means. A code longer than the field is cut
// to length rather than refused.
//
// # When the code was refused
//
// The rejection Page holds for Name is drawn under the squares, marks the group
// invalid and is named first in the group's aria-describedby, so it is
// announced as well as shown. A refused code is the one state of this field
// somebody cannot see for themselves, and a row of red squares with nothing
// saying why reads as a field that is broken.
//
// It publishes root, square, raw and message.
type OneTimeCodeProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, carried by the input that submits.
	Name string
	// ID is the id a Label points at. Empty uses Name. The text input that
	// submits carries it, and the behaviour hands it to the first square when
	// it shows them.
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
	// DescribedBy is the id of an element the caller drew that explains this
	// field. The rejection is not it: this component draws that one and names
	// it ahead of this id.
	DescribedBy string
	// Alphanumeric takes letters as well as digits. The default is digits
	// alone, and it is what sets the on-screen keyboard to a number pad.
	Alphanumeric bool
	// Disabled takes the field out of the form and out of the tab order.
	Disabled bool
	// Autofocus puts the cursor in the field on load: the text input, and the
	// first square once the behaviour shows them.
	Autofocus bool
}

// ElementID is the id the text input that submits carries, and the one a Label
// points at.
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

// Pattern is the shape the text input holds a code to without the script:
// exactly Size characters, digits alone unless Alphanumeric. It is a count and
// an alphabet, the same filter the behaviour applies to what is typed in a
// square, and not a verdict on the code.
func (p OneTimeCodeProps) Pattern() string {
	if p.Alphanumeric {
		return fmt.Sprintf("[0-9A-Za-z]{%d}", p.Size())
	}
	return fmt.Sprintf("[0-9]{%d}", p.Size())
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

// MessageID is the id the rejection is drawn under. It hangs off ElementID, so
// two fields given different IDs never answer to the same message.
func (p OneTimeCodeProps) MessageID() string {
	return p.ElementID() + "-error"
}

// Describers names everything that explains the group, in the order it should
// be read: the rejection when there is one, then the element the caller pointed
// at with DescribedBy. Empty when there is neither, so no attribute points at
// an element that is not on the page.
func (p OneTimeCodeProps) Describers() string {
	var ids []string
	if p.Message() != "" {
		ids = append(ids, p.MessageID())
	}
	if p.DescribedBy != "" {
		ids = append(ids, p.DescribedBy)
	}
	return strings.Join(ids, " ")
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
func (p OneTimeCodeProps) PartNames() []string { return []string{"root", "square", "raw", "message"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("otp") }}"
	role="group"
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@if(.Describers() != "")
		aria-describedby="{{ .Describers() }}"
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
			hidden
			inputmode="{{ .InputMode() }}"
			maxlength="1"
			value="{{ .CharAt(index) }}"
			autocomplete="{{ .AutocompleteAt(index) }}"
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("square"))
		>
	@endforeach
	{{-- The field that submits, and the one a person types into when the
	     behaviour has not run. The behaviour makes it hidden and shows the
	     squares; nothing about the form changes when it does. --}}
	<input
		data-part="raw"
		class="{{ .PartClass("raw", "input") }}"
		type="text"
		id="{{ .ElementID() }}"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		inputmode="{{ .InputMode() }}"
		autocomplete="one-time-code"
		maxlength="{{ .Size() }}"
		pattern="{{ .Pattern() }}"
		@if(.Disabled)
			disabled
		@endif
		@if(.Autofocus)
			autofocus
		@endif
		@attributes(.PartAttrs("raw"))
	>
	{{-- Inside the group rather than after it: the root is the outermost
	     element, and the sentence belongs to the group it explains. --}}
	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .MessageID() }}"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
</div>
