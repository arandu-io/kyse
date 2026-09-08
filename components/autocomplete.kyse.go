//go:build kyse

package components

import "strconv"

@go
// AutocompleteProps is a text box with suggestions the browser draws itself.
//
// It is an input and a datalist, and the list is drawn, positioned, filtered
// and dismissed by the browser. That is the whole difference from a
// [Combobox]: this one suggests, and a combobox chooses. What is typed here is
// the value whether or not it is in the list -- a town nobody has heard of, a
// tag that does not exist yet -- and what a combobox submits has to be one of
// the options it was given.
//
// So the two are not the same control with a flag between them. A combobox
// refuses a value off its list and needs a URL to have a list at all; this
// accepts anything and can have its suggestions in the page.
//
// The suggestions can also come from the server, which is what SearchURL is:
// the datalist is swapped, not the box, so what is typed is never disturbed by
// an answer arriving.
//
// It publishes root, label, input, message and hint.
type AutocompleteProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above the box.
	Label string
	// Value is what is in the box when nothing was rejected.
	Value string
	// Options are the suggestions. They are offered, never enforced.
	Options []AutocompleteOption
	// SearchURL replaces the suggestions as somebody types, by swapping the
	// list. Empty leaves the suggestions as they were drawn, which is the
	// right answer for a set that fits on the page.
	SearchURL string
	// Delay is how long typing pauses before the suggestions are fetched, in
	// milliseconds. Zero waits 200, which is long enough that a word typed at
	// speed is one request and short enough that nobody waits for it.
	Delay int
	// Placeholder is the grey text in the empty box.
	Placeholder string
	// Hint is the sentence under it while nothing is wrong.
	Hint string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
	// MaxLength is the ceiling the browser enforces while typing. Zero leaves
	// it open.
	MaxLength int
}

// AutocompleteOption is one suggestion.
type AutocompleteOption struct {
	// Value is what is put in the box when the suggestion is taken, and what
	// the form then submits.
	Value string
	// Label is the description shown beside it: the airport for a code, the
	// full name for a handle. Empty shows the value alone.
	Label string
}

// Current is what is in the box: what came back rejected, or what the caller
// set.
func (p AutocompleteProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// ListID is the id of the datalist, which the box points at with list.
func (p AutocompleteProps) ListID() string { return p.Name + "-suggestions" }

// Trigger is when the suggestions are fetched: after typing has paused, and
// only when what is in the box actually changed. Without the changed guard,
// an arrow key or a modifier would fetch the same list again.
func (p AutocompleteProps) Trigger() string {
	delay := p.Delay
	if delay <= 0 {
		delay = 200
	}
	return "input changed delay:" + strconv.Itoa(delay) + "ms"
}

// Message is the rejection for this field, and empty when there is none.
func (p AutocompleteProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the box.
func (p AutocompleteProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p AutocompleteProps) PartNames() []string {
	return []string{"root", "label", "input", "list", "message", "hint"}
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

	{{-- autocomplete="off" turns off the browser's own history for this box.
	     The two lists would otherwise be drawn one over the other, and the one
	     underneath is whatever was typed into a box of the same name on some
	     other site. --}}
	<input
		data-part="input"
		class="{{ .PartClass("input", "input") }}"
		type="text"
		id="{{ .Name }}"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		list="{{ .ListID() }}"
		autocomplete="off"
		@attributes(.PartAttrs("input"))
		@if(.SearchURL != "")
			hx-get="{{ .SearchURL }}"
			hx-trigger="{{ .Trigger() }}"
			hx-target="#{{ .ListID() }}"
			hx-swap="innerHTML"
			hx-sync="this:replace"
		@endif
		@if(.Placeholder != "")
			placeholder="{{ .Placeholder }}"
		@endif
		@if(.MaxLength > 0)
			maxlength="{{ .MaxLength }}"
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

	<datalist
		data-part="list"
		@if(.PartClass("list") != "")
			class="{{ .PartClass("list") }}"
		@endif
		id="{{ .ListID() }}"
		@attributes(.PartAttrs("list"))
	>
		@foreach(.Options as option)
			<option
				value="{{ option.Value }}"
				@if(option.Label != "")
					label="{{ option.Label }}"
				@endif
			></option>
		@endforeach
	</datalist>

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
