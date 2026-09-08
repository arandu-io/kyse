//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/kyse/icons"
)

@go
// NumberInputBehavior is the name the client behaviour is registered under
// with arandu.ui.define.
const NumberInputBehavior = "number-input"

// NumberInputProps is a number, with the two buttons that step it.
//
// The box is an input of type number, so the arrow keys, the range check and
// the numeric keypad on a phone are the browser's and not this file's. What is
// added is the pair of buttons, because the browser's own spinner is a target
// four pixels tall in most engines and absent in some.
//
// The buttons carry aria-hidden and the box keeps the label, so the control
// announces as one spin button and not as a field followed by two anonymous
// buttons. Nothing is lost by that: a screen reader user steps the value with
// the arrow keys, which is the same interaction the buttons exist to replace
// for a pointer.
//
// It publishes root, label, group, decrement, input, increment, message and
// hint.
type NumberInputProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above the box.
	Label string
	// Value is what is shown when nothing was rejected. It is a string so that
	// an empty box and a zero are different things -- a quantity nobody has
	// typed is not a quantity of none.
	Value string
	// Min and Max are the range. Both empty leaves it open, which is right for
	// a difference and wrong for a quantity.
	Min string
	Max string
	// Step is how much a press moves it. Empty steps by one; "0.01" is the
	// step of an amount of money, and "any" accepts anything.
	Step string
	// Unit is drawn after the box: "kg", "%", "days". It is not part of the
	// value and is never submitted.
	Unit string
	// Placeholder is the grey text in the empty box.
	Placeholder string
	// Hint is the sentence under it while nothing is wrong.
	Hint string
	// Required marks the field, and the browser refuses to submit without it.
	Required bool
	// Disabled draws the box and both buttons as unavailable.
	Disabled bool
	// DecrementLabel and IncrementLabel name the two buttons for a screen
	// reader that reaches them anyway. Empty says "Decrease" and "Increase".
	DecrementLabel string
	IncrementLabel string
}

// Current is the value in the box: what came back rejected, or what the caller
// set.
func (p NumberInputProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// Message is the rejection for this field, and empty when there is none.
func (p NumberInputProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the box: the rejection, or
// the hint, and never both -- a field draws one of them.
func (p NumberInputProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// StepAmount is how much a press moves the value, as a number. It is one when
// the step is empty or is "any", because "any" means the box accepts anything
// and says nothing about what a press should do.
func (p NumberInputProps) StepAmount() float64 {
	if p.Step == "" || p.Step == "any" {
		return 1
	}
	amount, err := strconv.ParseFloat(p.Step, 64)
	if err != nil || amount <= 0 {
		return 1
	}
	return amount
}

// DecrementName is what the down button is called.
func (p NumberInputProps) DecrementName() string {
	if p.DecrementLabel != "" {
		return p.DecrementLabel
	}
	return "Decrease"
}

// IncrementName is what the up button is called.
func (p NumberInputProps) IncrementName() string {
	if p.IncrementLabel != "" {
		return p.IncrementLabel
	}
	return "Increase"
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p NumberInputProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: NumberInputBehavior, Props: map[string]any{"step": p.StepAmount()}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p NumberInputProps) PartNames() []string {
	return []string{"root", "label", "group", "decrement", "input", "increment", "unit", "message", "hint"}
}
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

	<div
		data-part="group"
		class="{{ .PartClass("group", "input-group") }}"
		@attributes(.PartAttrs("group"))
	>
		<button
			data-part="decrement"
			class="{{ .PartClass("decrement", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			data-align="start"
			data-number-step="-1"
			tabindex="-1"
			aria-hidden="true"
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("decrement"))
		>{!! icons.Minus(icons.Props{}) !!}</button>

		<input
			data-part="input"
			@if(.PartClass("input") != "")
				class="{{ .PartClass("input") }}"
			@endif
			type="number"
			inputmode="decimal"
			id="{{ .Name }}"
			name="{{ .Name }}"
			value="{{ .Current() }}"
			@attributes(.PartAttrs("input"))
			@if(.Min != "")
				min="{{ .Min }}"
			@endif
			@if(.Max != "")
				max="{{ .Max }}"
			@endif
			@if(.Step != "")
				step="{{ .Step }}"
			@endif
			@if(.Placeholder != "")
				placeholder="{{ .Placeholder }}"
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
		>

		@if(.Unit != "")
			<span
				data-part="unit"
				class="{{ .PartClass("unit", "text-muted-foreground text-sm") }}"
				data-align="end"
				@attributes(.PartAttrs("unit"))
			>{{ .Unit }}</span>
		@endif

		<button
			data-part="increment"
			class="{{ .PartClass("increment", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			data-align="end"
			data-number-step="1"
			tabindex="-1"
			aria-hidden="true"
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("increment"))
		>{!! icons.Plus(icons.Props{}) !!}</button>
	</div>

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
