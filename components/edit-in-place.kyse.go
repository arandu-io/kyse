//go:build kyse

package components

@go
// EditInPlaceProps is a value shown as text that becomes a field when it is
// edited.
//
// It is two fragments and one endpoint pair, and that is the whole design. The
// reading view is drawn here with a button that fetches the editing view; the
// editing view saves and answers with the reading view again. Cancel re-fetches
// the reading view rather than reverting anything in the browser, so what is
// on the screen after a cancel is what the server has -- which is the point of
// cancelling.
//
// Nothing is kept anywhere between the two. A component that remembered the
// original value would be remembering it in the one place that cannot be told
// when somebody else changes it.
//
// It publishes root, value, edit, form, input, save and cancel.
type EditInPlaceProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the region hangs its id off, and what the endpoints target.
	ID string
	// Name is the field name when this is drawn as the editing view.
	Name string
	// Label names the value: "Display name". It is the accessible name of the
	// edit button and of the field, so it is needed in both views.
	Label string
	// Value is what is stored, drawn as text in the reading view and as the
	// field's value in the editing one.
	Value string
	// Placeholder is what the reading view says when the value is empty:
	// "Not set". Without it an empty value draws an empty line with nothing
	// to click.
	Placeholder string
	// EditURL answers with the editing view. It is fetched when the button is
	// pressed.
	EditURL string
	// SaveURL takes the new value and answers with the reading view.
	SaveURL string
	// CancelURL answers with the reading view unchanged. Empty uses EditURL's
	// sibling -- there is none, so a component without it draws no cancel, and
	// escape is then the only way out.
	CancelURL string
	// Editing draws the editing view instead of the reading one. It is what
	// the endpoint sets when it answers the edit request.
	Editing bool
	// Multiline draws a box rather than a line, for a value that is a
	// paragraph.
	Multiline bool
	// EditLabel, SaveLabel and CancelLabel are the three buttons. Empty says
	// "Edit", "Save" and "Cancel".
	EditLabel   string
	SaveLabel   string
	CancelLabel string
	// Message is a rejection from the save, drawn under the field.
	Message string
}

// Shown is what the reading view draws: the value, or the words standing in
// for an empty one.
func (p EditInPlaceProps) Shown() string {
	if p.Value != "" {
		return p.Value
	}
	return p.Placeholder
}

// EditName is what the edit button is called. "Edit" repeated down a page of
// fields is the same name on every one of them.
func (p EditInPlaceProps) EditName() string {
	if p.EditLabel != "" {
		return p.EditLabel
	}
	if p.Label != "" {
		return "Edit " + p.Label
	}
	return "Edit"
}

// SaveName is the text on the save button.
func (p EditInPlaceProps) SaveName() string {
	if p.SaveLabel != "" {
		return p.SaveLabel
	}
	return "Save"
}

// CancelName is the text on the cancel button.
func (p EditInPlaceProps) CancelName() string {
	if p.CancelLabel != "" {
		return p.CancelLabel
	}
	return "Cancel"
}

// PartNames are the parts this component publishes.
func (p EditInPlaceProps) PartNames() []string {
	return []string{"root", "value", "edit", "label", "input", "save", "cancel", "message"}
}
@endgo

@if(!.Editing)
	<div
		data-part="root"
		class="{{ .RootClass("edit-in-place") }}"
		id="{{ .ID }}"
		@attributes(.RootAttrs())
	>
		<span
			data-part="value"
			@if(.PartClass("value") != "")
				class="{{ .PartClass("value") }}"
			@endif
			@attributes(.PartAttrs("value"))
		>{{ .Shown() }}</span>

		<button
			data-part="edit"
			class="{{ .PartClass("edit", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			aria-label="{{ .EditName() }}"
			@attributes(.PartAttrs("edit"))
			hx-get="{{ .EditURL }}"
			hx-target="#{{ .ID }}"
			hx-swap="outerHTML"
		>{{ .EditName() }}</button>
	</div>
@endif

@if(.Editing)
	{{-- autofocus, because the whole interaction is "press edit and type": a
	     field that opens without the caret in it costs a click nobody expects
	     to have to make. --}}
	<form
		data-part="root"
		class="{{ .RootClass("edit-in-place") }}"
		id="{{ .ID }}"
		@attributes(.RootAttrs())
		hx-put="{{ .SaveURL }}"
		hx-target="#{{ .ID }}"
		hx-swap="outerHTML"
	>
		<label
			data-part="label"
			class="{{ .PartClass("label", "sr-only") }}"
			for="{{ .ID }}-input"
			@attributes(.PartAttrs("label"))
		>{{ .Label }}</label>

		@if(.Multiline)
			<textarea
				data-part="input"
				class="{{ .PartClass("input", "textarea") }}"
				id="{{ .ID }}-input"
				name="{{ .Name }}"
				autofocus
				@attributes(.PartAttrs("input"))
				@if(.Message != "")
					aria-invalid="true"
					aria-describedby="{{ .ID }}-error"
				@endif
			>{{ .Value }}</textarea>
		@endif
		@if(!.Multiline)
			<input
				data-part="input"
				class="{{ .PartClass("input", "input") }}"
				type="text"
				id="{{ .ID }}-input"
				name="{{ .Name }}"
				value="{{ .Value }}"
				autofocus
				@attributes(.PartAttrs("input"))
				@if(.Message != "")
					aria-invalid="true"
					aria-describedby="{{ .ID }}-error"
				@endif
			>
		@endif

		<button
			data-part="save"
			class="{{ .PartClass("save", "btn") }}"
			type="submit"
			data-size="sm"
			@attributes(.PartAttrs("save"))
		>{{ .SaveName() }}</button>

		@if(.CancelURL != "")
			<button
				data-part="cancel"
				class="{{ .PartClass("cancel", "btn") }}"
				type="button"
				data-variant="ghost"
				data-size="sm"
				@attributes(.PartAttrs("cancel"))
				hx-get="{{ .CancelURL }}"
				hx-target="#{{ .ID }}"
				hx-swap="outerHTML"
			>{{ .CancelName() }}</button>
		@endif

		@if(.Message != "")
			<p
				data-part="message"
				id="{{ .ID }}-error"
				class="{{ .PartClass("message", "text-destructive text-sm") }}"
				@attributes(.PartAttrs("message"))
			>{{ .Message }}</p>
		@endif
	</form>
@endif
