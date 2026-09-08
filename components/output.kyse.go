//go:build kyse

package components

@go
// OutputProps is the result of a calculation, in the element the browser
// already announces when it changes.
//
// An <output> carries an implicit live region, so a value swapped into it is
// read out without anything declaring aria-live. That is the whole reason to
// use it rather than a span: a span that changes is a change nobody hears.
//
// It publishes root.
type OutputProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Value is what is shown. It is a string and not a number because what a
	// result reads as -- the currency, the unit, the number of decimals -- is
	// the caller's decision and not this element's.
	Value string
	// Name is the form field name, for a result submitted with the form.
	Name string
	// For is the space-separated ids of the fields this result was computed
	// from, which is what tells assistive technology what the number belongs
	// to.
	For string
	// Tone is "muted", "primary" or "destructive". Empty is the plain one.
	Tone string
	// Urgent announces the change by interrupting rather than waiting. It is
	// for a result that is a failure; a number that merely changed is not.
	Urgent bool
}

// Live is how insistently a change is announced.
func (p OutputProps) Live() string {
	if p.Urgent {
		return "assertive"
	}
	return "polite"
}

// PartNames are the parts this component publishes.
func (p OutputProps) PartNames() []string { return []string{"root"} }
@endgo

<output
	data-part="root"
	class="{{ .RootClass("output") }}"
	aria-live="{{ .Live() }}"
	@if(.Name != "")
		name="{{ .Name }}"
	@endif
	@if(.For != "")
		for="{{ .For }}"
	@endif
	@if(.Tone != "")
		data-tone="{{ .Tone }}"
	@endif
	@attributes(.RootAttrs())
>{{ .Value }}</output>
