//go:build kyse

package components

@go
// DateTimePickerProps is a date, a time, or both, in the field the browser
// already draws a calendar for.
//
// There is no calendar written here, and that is the decision. A date field is
// where a locale is least negotiable: the order of day and month, the first
// day of the week, the names of the months, the calendar system itself. The
// browser knows the reader's; a script does not, and one that guesses is wrong
// in exactly the places where being wrong is expensive.
//
// The value crossing the wire is never what is drawn. A date field submits
// "2026-09-07" whatever the reader sees, and a datetime-local field submits
// "2026-09-07T14:30" with no zone at all -- so a moment that has to be exact
// is either stored with the zone alongside it or asked for in UTC.
//
// It publishes root, label, input, message and hint.
type DateTimePickerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above it.
	Label string
	// Kind is what is being asked for: "date", "time", "datetime", "month" or
	// "week". Empty asks for a date.
	Kind string
	// Value is what is shown, in the form the chosen kind submits:
	// "2026-09-07" for a date, "14:30" for a time, "2026-09-07T14:30" for
	// both, "2026-09" for a month, "2026-W37" for a week. A value in any other
	// shape is ignored by the browser and the field comes back empty.
	Value string
	// Min and Max are the range, in the same form as Value. A booking that
	// cannot be in the past sets Min and is then refused by the browser as
	// well as by the server.
	Min string
	Max string
	// Step is the granularity, in seconds for a time and in days for a date.
	// "60" is whole minutes, which is the default for a time field; "900" is
	// quarter hours. Empty leaves the browser's.
	Step string
	// Hint is the sentence under it while nothing is wrong. It is where the
	// zone belongs, because the field carries none: "Times are UTC".
	Hint string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
	// ReadOnly shows the value without letting it change. It still submits,
	// which is what separates it from Disabled.
	ReadOnly bool
}

// InputType is the element type for the kind asked for. An unknown kind is a
// date, because a field that fell back to text would accept anything and
// submit it.
func (p DateTimePickerProps) InputType() string {
	switch p.Kind {
	case "time":
		return "time"
	case "datetime":
		return "datetime-local"
	case "month":
		return "month"
	case "week":
		return "week"
	default:
		return "date"
	}
}

// Current is the value shown: what came back rejected, or what the caller set.
func (p DateTimePickerProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// Message is the rejection for this field, and empty when there is none.
func (p DateTimePickerProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the field.
func (p DateTimePickerProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p DateTimePickerProps) PartNames() []string {
	return []string{"root", "label", "input", "message", "hint"}
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

	<input
		data-part="input"
		class="{{ .PartClass("input", "input") }}"
		type="{{ .InputType() }}"
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
		@if(.ReadOnly)
			readonly
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
