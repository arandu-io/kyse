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
	// Name is the form field name, and the id the label points at.
	Name string
	// Label is the text above the box.
	Label string
	// Value is what the box starts with.
	Value string
	// Placeholder is the grey text inside an empty box.
	Placeholder string
	// Hint is the sentence under the box.
	Hint string
	// Error is what came back from validation.
	Error string
	// Rows is the height in lines. Zero draws the browser default.
	Rows int
	// Required marks the box required.
	Required bool
}

// DescribedBy names the element that explains this box.
func (p TextareaProps) DescribedBy() string {
	if p.Error != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
@endgo

<div class="field">
	<label class="label" for="{{ .Name }}">{{ .Label }}</label>
	<textarea
		class="textarea"
		id="{{ .Name }}"
		name="{{ .Name }}"
		@if(.Rows > 0)
			rows="{{ .Rows }}"
		@endif
		@if(.Placeholder != "")
			placeholder="{{ .Placeholder }}"
		@endif
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Error != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			required
		@endif
	>{{ .Value }}</textarea>
	@if(.Error != "")
		<p id="{{ .Name }}-error" class="text-destructive text-sm">{{ .Error }}</p>
	@endif
	@if(.Error == "")
		@if(.Hint != "")
			<p id="{{ .Name }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
		@endif
	@endif
</div>
