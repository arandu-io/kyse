//go:build kyse

package components

import "strings"

@go
// RadioGroupProps is one answer chosen from a few, drawn as the whole set at
// once.
//
// The group is the component, and not the single radio, because a radio on its
// own cannot be got right: what makes the choice exclusive is the shared name,
// what names the set is a label no single control owns, and what tells a screen
// reader it is looking at a set is a role on the element around them. Anything
// drawing one radio at a time leaves all three to be repeated by hand.
//
// Field draws none of it -- its row is one box somebody types into. Select is
// the same choice for a longer list: a group is worth its space up to about five
// options, and past that it is a wall.
type RadioGroupProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name shared by every radio, which is what makes
	// the choice exclusive, and what Page is asked about.
	Name string
	// Label names the whole set. It is not a label element, because a label
	// element names one control and this names several.
	Label string
	// Value is the option selected when nothing was rejected: the stored record
	// on an edit form, and empty on a create form, where nothing is selected.
	// What was chosen on a rejected attempt takes precedence -- see Current.
	Value string
	// Hint is the sentence under the set that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this set asks for its
	// message and for what was chosen. Nil draws no message and keeps Value.
	Page Page
	// Options are the choices, top to bottom. A set with one option is a
	// question with one answer, and a radio nobody can clear once it is picked.
	Options []RadioOption
	// Required refuses the form until one option is chosen.
	Required bool
}

// RadioOption is one choice in the set.
type RadioOption struct {
	// Label is the text beside the radio.
	Label string
	// Value is what this choice submits, and what Value and Current are
	// compared against.
	Value string
	// Hint is the sentence under the label saying what the choice means. Empty
	// draws none.
	Hint string
	// Disabled takes this one choice out of the form and out of the tab order,
	// leaving the rest of the set alone.
	Disabled bool
}

// OptionID is the id of the radio carrying a value, and what its label points
// at.
//
// It is derived from the group name and the value rather than asked for, so that
// a set of options is written as a list of what they say and what they submit.
// The value is folded to the characters an id may hold, which means two values
// differing only in the characters that are folded come out as one id -- the
// price of deriving rather than being told.
func (p RadioGroupProps) OptionID(value string) string {
	folded := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9':
			return r
		case r == '-', r == '_':
			return r
		default:
			return '-'
		}
	}, value)
	return p.Name + "-" + folded
}

// Message is what validation left for this set, or empty.
func (p RadioGroupProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is the value drawn selected: what was chosen on the rejected attempt,
// and Value when there was none.
func (p RadioGroupProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// DescribedBy names the element that explains this set, so the message and the
// hint are announced rather than only shown.
func (p RadioGroupProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
// PartNames are the parts this component publishes.
func (p RadioGroupProps) PartNames() []string {
	return []string{"root", "label", "group", "option", "input", "option-label", "message", "hint"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("field") }}"
	@attributes(.RootAttrs())
	@if(.Message() != "")
		data-invalid="true"
	@endif
>
	<span
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		id="{{ .Name }}-label"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</span>
	<div
		data-part="group"
		@if(.PartClass("group") != "")
			class="{{ .PartClass("group") }}"
		@endif
		role="radiogroup"
		aria-labelledby="{{ .Name }}-label"
		@attributes(.PartAttrs("group"))
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			aria-required="true"
		@endif
	>
		@foreach(.Options as option)
			<div
				data-part="option"
				class="{{ .PartClass("option", "field") }}"
				data-orientation="horizontal"
				@attributes(.PartAttrs("option"))
				@if(option.Disabled)
					data-disabled="true"
				@endif
			>
				<input
					data-part="input"
					class="{{ .PartClass("input", "input") }}"
					type="radio"
					id="{{ .OptionID(option.Value) }}"
					name="{{ .Name }}"
					value="{{ option.Value }}"
					@attributes(.PartAttrs("input"))
					@if(option.Value == .Current())
						checked
					@endif
					@if(option.Hint != "")
						aria-describedby="{{ .OptionID(option.Value) }}-hint"
					@endif
					@if(.Required)
						required
					@endif
					@if(option.Disabled)
						disabled
					@endif
				>
				<section>
					<label
						data-part="option-label"
						class="{{ .PartClass("option-label", "label") }}"
						for="{{ .OptionID(option.Value) }}"
						@attributes(.PartAttrs("option-label"))
					>{{ option.Label }}</label>
					@if(option.Hint != "")
						<p id="{{ .OptionID(option.Value) }}-hint" class="text-muted-foreground text-sm">{{ option.Hint }}</p>
					@endif
				</section>
			</div>
		@endforeach
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
