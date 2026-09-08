//go:build kyse

package components

@go
// TextareaProps is a labelled textarea, with its error.
//
// It is FieldProps without the type and with rows, kept separate because a
// textarea is a different element with a different way of carrying its value --
// between the tags rather than in an attribute -- and one component branching on
// which element to draw would be two components sharing a name.
type TextareaProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, and the id the label points at.
	Name string
	// Label is the text above the box.
	Label string
	// Value is what the box starts with when nothing was rejected. What was
	// typed on a rejected attempt takes precedence -- see Current.
	Value string
	// Placeholder is the grey text inside an empty box.
	Placeholder string
	// Hint is the sentence under the box.
	Hint string
	// Page is the screen's own view.Page, which is what this box asks for its
	// message and for what was typed. Nil draws no message and keeps Value.
	Page Page
	// Rows is the height in lines. Zero draws the browser default.
	Rows int
	// Required marks the box required.
	Required bool
	// Autosize grows and shrinks the box to fit what is in it, with Rows as
	// the floor.
	//
	// It is one CSS declaration -- field-sizing: content -- and not the loop
	// that reads scrollHeight and writes a height back on every keystroke.
	// That loop is what this replaces: it forces layout twice per character,
	// and it is wrong the first time a font loads late or the box is drawn
	// while hidden. A browser without the declaration draws a fixed box of
	// Rows lines, which is the box that was there before.
	//
	// There is no ceiling here, and that is deliberate: a maximum height is a
	// length, and a length written as a field would be a length this component
	// has to turn into CSS at run time -- which is the one thing a class
	// cannot be built from. A caller who wants one adds it where every other
	// class goes: Parts["input"].Class, with "max-h-64" or whatever the design
	// says.
	Autosize bool
}

// Message is what validation left for this box, or empty.
func (p TextareaProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is what the box is drawn with: what was typed on the rejected
// attempt, and Value when there was none.
func (p TextareaProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// DescribedBy names the element that explains this box.
func (p TextareaProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
// PartNames are the parts this component publishes.
func (p TextareaProps) PartNames() []string { return []string{"root", "label", "input", "message", "hint"} }
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
	<textarea
		data-part="input"
		class="{{ .PartClass("input", "textarea") }}"
		id="{{ .Name }}"
		name="{{ .Name }}"
		@attributes(.PartAttrs("input"))
		@if(.Autosize)
			data-autosize="true"
		@endif
		@if(.Rows > 0)
			rows="{{ .Rows }}"
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
	>{{ .Current() }}</textarea>
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
