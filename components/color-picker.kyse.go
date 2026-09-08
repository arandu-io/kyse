//go:build kyse

package components

@go
// ColorPickerProps is a colour, chosen with the picker the operating system
// already has.
//
// The whole control is an input of type color. The wheel, the eyedropper, the
// recent swatches and the keyboard path into all of them belong to the
// platform, and a picker written here would be a worse copy of each on every
// platform at once.
//
// What is here is the frame around it and the readout beside it, because an
// input of type color draws as a small unlabelled swatch and says nothing
// about which colour it holds.
//
// The value is always a seven-character hex string with the hash, because that
// is the only form the element accepts and the only one it returns. A named
// colour or an rgb() written into it is silently ignored, and the box comes
// back black.
//
// It publishes root, label, group, input, readout, message and hint.
type ColorPickerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above it.
	Label string
	// Value is the colour, as "#1d4ed8". Empty is black, because the element
	// has no empty state -- it always holds a colour.
	Value string
	// ShowValue draws the hex beside the swatch, which is what makes the
	// chosen colour something a person can write down, paste or compare.
	ShowValue bool
	// Swatches are the colours offered before the full picker is opened, as
	// hex strings. The browser draws them in its own picker, and a browser
	// that does not support the list ignores it.
	Swatches []string
	// Hint is the sentence under it while nothing is wrong.
	Hint string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
}

// Current is the colour shown: what came back rejected, what the caller set,
// or black.
func (p ColorPickerProps) Current() string {
	chosen := p.Value
	if p.Page != nil {
		chosen = p.Page.OldOr(p.Name, p.Value)
	}
	if chosen == "" {
		return "#000000"
	}
	return chosen
}

// ListID is the id of the swatch list, and empty when there is none.
func (p ColorPickerProps) ListID() string {
	if len(p.Swatches) == 0 {
		return ""
	}
	return p.Name + "-swatches"
}

// Message is the rejection for this field, and empty when there is none.
func (p ColorPickerProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the field.
func (p ColorPickerProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p ColorPickerProps) PartNames() []string {
	return []string{"root", "label", "group", "input", "readout", "message", "hint"}
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
		class="{{ .PartClass("group", "color-picker") }}"
		@attributes(.PartAttrs("group"))
	>
		<input
			data-part="input"
			@if(.PartClass("input") != "")
				class="{{ .PartClass("input") }}"
			@endif
			type="color"
			id="{{ .Name }}"
			name="{{ .Name }}"
			value="{{ .Current() }}"
			@attributes(.PartAttrs("input"))
			@if(.ListID() != "")
				list="{{ .ListID() }}"
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

		@if(.ShowValue)
			<output
				data-part="readout"
				class="{{ .PartClass("readout", "text-muted-foreground font-mono text-sm uppercase") }}"
				for="{{ .Name }}"
				aria-live="polite"
				@attributes(.PartAttrs("readout"))
			>{{ .Current() }}</output>
		@endif
	</div>

	@if(.ListID() != "")
		<datalist id="{{ .ListID() }}">
			@foreach(.Swatches as swatch)
				<option value="{{ swatch }}"></option>
			@endforeach
		</datalist>
	@endif

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
