//go:build kyse

package components

@go
// LabelProps is the text that names a control, tied to it by id.
//
// Field draws its own, so a labelled form row needs nothing from here. This is
// the other half of Input: when the caller lays the row out, the name is a
// separate element, and the two halves have to agree on one id.
//
// The required marker is drawn from a flag rather than written into Text,
// because it is decoration and has to be kept out of the accessible name: read
// aloud it turns "Email" into "Email star", and it is not what announces the
// constraint either -- the required attribute on the control is.
// It publishes root and marker, marker being the asterisk beside the text.
type LabelProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// For is the id of the control this names, which is the control's id and not
	// its name. The two are the same only where nothing set them apart.
	For string
	// Text is what the label says.
	Text string
	// Required draws the marker beside the text. It is a mark on the page and
	// nothing more: the control carries the required attribute.
	Required bool
}

// PartNames are the parts this component publishes.
func (p LabelProps) PartNames() []string { return []string{"root", "marker"} }
@endgo

<label
	data-part="root"
	class="{{ .RootClass("label") }}"
	for="{{ .For }}"
	@attributes(.RootAttrs())
>
	{{ .Text }}
	@if(.Required)
		<span
			data-part="marker"
			class="{{ .PartClass("marker", "text-destructive") }}"
			aria-hidden="true"
			@attributes(.PartAttrs("marker"))
		>*</span>
	@endif
</label>
