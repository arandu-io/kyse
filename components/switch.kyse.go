//go:build kyse

package components

@go
// SwitchProps is one setting that is on or off, drawn as a track with a thumb
// and named from the left.
//
// It is the same element as Checkbox with role="switch" on it, and it is a
// separate component rather than a flag on that one because the two are read
// aloud differently -- "on" and "off" against "ticked" and "not ticked" -- and
// because they are used at different moments. A checkbox is an answer given
// while filling a form in. A switch is a setting somebody throws, and the row
// reads as a setting: the name on the left, the control on the right, the way
// the rest of the row does.
//
// Field draws neither. Its row is a name over a box somebody types into.
//
// aria-checked is deliberately not written. The state belongs to the input, and
// the browser reports it from there; an attribute written once at render would
// still say "off" after somebody turned the switch on, which is worse than
// having said nothing.
type SwitchProps struct {
	// Name is the form field name, and what Page is asked about.
	Name string
	// ID is the id the markup carries, and the id the label points at. Empty
	// uses Name -- which a list of switches sharing one name has to set, since
	// otherwise every label in the list points at the first switch.
	ID string
	// Label is the text naming the setting.
	Label string
	// Value is what a switch that is on submits. Empty submits "on", which is
	// what a browser sends for a checkbox with no value of its own.
	Value string
	// Hint is the sentence under the label saying what the setting does.
	Hint string
	// Page is the screen's own view.Page, which is what this switch asks for its
	// message and for whether it came back on. Nil draws no message and keeps
	// On.
	Page Page
	// On is whether the switch starts on when nothing was rejected: the stored
	// record on an edit form, and the default on a new one.
	On bool
	// Disabled takes the switch out of the form and out of the tab order.
	Disabled bool
}

// ElementID is the id the switch carries: ID when one was given, and Name when
// none was.
func (p SwitchProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// SubmittedValue is what a switch that is on sends: Value when it was given, and
// "on" when it was not.
func (p SwitchProps) SubmittedValue() string {
	if p.Value == "" {
		return "on"
	}
	return p.Value
}

// Message is what validation left for this switch, or empty.
func (p SwitchProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current reports whether the switch is drawn on.
//
// A switch that is on submits its value and one that is off submits nothing at
// all, so what comes back from a rejected attempt answers only half the
// question: the value being there means it was on, and it being absent means
// either that it was off or that there was no attempt. Absent therefore falls
// back to On, which is right on the form nobody submitted and wrong in one case
// -- a switch somebody turned off on an attempt that was then rejected comes
// back on. Telling the two apart needs a second field on the wire, and this
// component does not invent one.
func (p SwitchProps) Current() bool {
	if p.Page == nil {
		return p.On
	}
	fallback := ""
	if p.On {
		fallback = p.SubmittedValue()
	}
	return p.Page.OldOr(p.Name, fallback) == p.SubmittedValue()
}

// DescribedBy names the element that explains this switch, so the message and
// the hint are announced rather than only shown.
func (p SwitchProps) DescribedBy() string {
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
	data-orientation="horizontal"
	@if(.Message() != "")
		data-invalid="true"
	@endif
	@if(.Disabled)
		data-disabled="true"
	@endif
>
	<section>
		<label class="label" for="{{ .ElementID() }}">{{ .Label }}</label>
		@if(.Message() != "")
			<p id="{{ .ElementID() }}-error" class="text-destructive text-sm">{{ .Message() }}</p>
		@endif
		@if(.Message() == "")
			@if(.Hint != "")
				<p id="{{ .ElementID() }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
			@endif
		@endif
	</section>
	<input
		class="input"
		type="checkbox"
		role="switch"
		id="{{ .ElementID() }}"
		name="{{ .Name }}"
		@if(.Value != "")
			value="{{ .Value }}"
		@endif
		@if(.Current())
			checked
		@endif
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Disabled)
			disabled
		@endif
	>
</div>
