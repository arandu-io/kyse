//go:build kyse

package components

@go
// SelectProps is one answer chosen from a list, drawn as the browser's own
// select.
//
// It is the native element and not a button opening a styled list of divs.
// There is one way to choose from a list here, and this is it: the native
// element is the one that submits with the form, that opens as the platform
// picker on a phone, that types-to-find and pages with the keyboard without
// being taught how, and that is still a working control when a script has not
// run. The styled kind buys a menu that matches the rest of the page and pays
// for it with a widget whose state lives in JavaScript -- and the two together
// would be two ways to ask one question.
//
// Field is the labelled row for a box somebody types into; this is the labelled
// row for a list. RadioGroup is the same choice with every option on the screen
// at once, which is worth its space up to about five.
type SelectProps struct {
	// Name is the form field name, and what Page is asked about.
	Name string
	// ID is the id the markup carries, and the id the label points at. Empty
	// uses Name.
	ID string
	// Label is the text above the list.
	Label string
	// Value is the option selected when nothing was rejected: the stored record
	// on an edit form, and empty on a create form. What was chosen on a rejected
	// attempt takes precedence -- see Current.
	Value string
	// Placeholder is the first line, drawn as an option submitting nothing and
	// selected while Value is empty. Empty draws no such line, and the list then
	// starts on its first real option, which is a choice nobody made.
	Placeholder string
	// Hint is the sentence under the list that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this list asks for its
	// message and for what was chosen. Nil draws no message and keeps Value.
	Page Page
	// Options are the choices, in the order they are drawn.
	Options []SelectOption
	// Required refuses the form while the placeholder is the chosen line.
	Required bool
	// Disabled takes the list out of the form and out of the tab order.
	Disabled bool
}

// SelectOption is one line of the list.
type SelectOption struct {
	// Label is what the line reads.
	Label string
	// Value is what this line submits, and what Value and Current are compared
	// against.
	Value string
	// Disabled draws the line and refuses it, which says a choice exists and is
	// not available -- what removing it from the list cannot say.
	Disabled bool
}

// ElementID is the id the list carries: ID when one was given, and Name when
// none was.
func (p SelectProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// Message is what validation left for this list, or empty.
func (p SelectProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is the value drawn selected: what was chosen on the rejected attempt,
// and Value when there was none.
func (p SelectProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// DescribedBy names the element that explains this list, so the message and the
// hint are announced rather than only shown.
func (p SelectProps) DescribedBy() string {
	if p.Message() != "" {
		return p.ElementID() + "-error"
	}
	if p.Hint != "" {
		return p.ElementID() + "-hint"
	}
	return ""
}
@endgo

<div
	class="field"
	@if(.Message() != "")
		data-invalid="true"
	@endif
	@if(.Disabled)
		data-disabled="true"
	@endif
>
	<label class="label" for="{{ .ElementID() }}">{{ .Label }}</label>
	<select
		class="select"
		id="{{ .ElementID() }}"
		name="{{ .Name }}"
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
		@if(.Placeholder != "")
			<option
				value=""
				@if(.Current() == "")
					selected
				@endif
			>{{ .Placeholder }}</option>
		@endif
		@foreach(.Options as option)
			<option
				value="{{ option.Value }}"
				@if(option.Value == .Current())
					selected
				@endif
				@if(option.Disabled)
					disabled
				@endif
			>{{ option.Label }}</option>
		@endforeach
	</select>
	@if(.Message() != "")
		<p id="{{ .ElementID() }}-error" class="text-destructive text-sm">{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p id="{{ .ElementID() }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
		@endif
	@endif
</div>
