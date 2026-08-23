//go:build kyse

package components

@go
// CheckboxProps is one box that is either ticked or not, with its label beside
// it rather than above it.
//
// Field cannot draw this and should not be asked to. Its row is a name over a
// box somebody types into, and it carries a string; a checkbox is a box with the
// name to its right, and what it carries is a yes or a no. Branching Field on
// which of the two it was would be two components sharing one name.
//
// Use it for consent and for one answer of many on a form somebody submits. Use
// Switch for a setting that takes effect when it is thrown.
type CheckboxProps struct {
	// Name is the form field name, and what Page is asked about.
	Name string
	// ID is the id the markup carries, and the id the label points at. Empty
	// uses Name -- which a list of boxes sharing one name has to set, since
	// otherwise every label in the list points at the first box.
	ID string
	// Label is the text beside the box.
	Label string
	// Value is what a ticked box submits. Empty submits "on", which is what a
	// browser sends for a box with no value of its own.
	Value string
	// Hint is the sentence under the label that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this box asks for its
	// message and for whether it came back ticked. Nil draws no message and
	// keeps Checked.
	Page Page
	// Checked is whether the box starts ticked when nothing was rejected: the
	// stored record on an edit form, and usually false on a create form.
	Checked bool
	// Required refuses the form until the box is ticked, which is what a consent
	// box wants and what a preference must not have.
	Required bool
	// Disabled takes the box out of the form and out of the tab order.
	Disabled bool
}

// ElementID is the id the box carries: ID when one was given, and Name when
// none was.
func (p CheckboxProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// SubmittedValue is what a ticked box sends: Value when it was given, and "on"
// when it was not.
func (p CheckboxProps) SubmittedValue() string {
	if p.Value == "" {
		return "on"
	}
	return p.Value
}

// Message is what validation left for this box, or empty.
func (p CheckboxProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current reports whether the box is drawn ticked.
//
// A ticked box submits its value and an unticked box submits nothing at all,
// so what comes back from a rejected attempt answers only half the question: the
// value being there means the box was ticked, and it being absent means either
// that it was unticked or that there was no attempt. Absent therefore falls
// back to Checked, which is right on the form nobody submitted and wrong in one
// case -- a box somebody unticked on an attempt that was then rejected comes
// back ticked. Telling the two apart needs a second field on the wire, and this
// component does not invent one.
func (p CheckboxProps) Current() bool {
	if p.Page == nil {
		return p.Checked
	}
	fallback := ""
	if p.Checked {
		fallback = p.SubmittedValue()
	}
	return p.Page.OldOr(p.Name, fallback) == p.SubmittedValue()
}

// DescribedBy names the element that explains this box, so the message and the
// hint are announced rather than only shown.
func (p CheckboxProps) DescribedBy() string {
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
	<input
		class="input"
		type="checkbox"
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
		@if(.Required)
			required
		@endif
		@if(.Disabled)
			disabled
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
</div>
