//go:build kyse

package components

@go
// FileUploadBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const FileUploadBehavior = "file-upload"

// FileUploadProps is a file field with a drop area over it.
//
// The field is a real input of type file, and the drop area is drawn around
// it rather than instead of it. That is what keeps the whole thing working: a
// keyboard opens the picker, a screen reader announces a file input, the form
// submits it as multipart, and dropping a file is an addition on top for
// people using a pointer.
//
// A drop area built as a div with handlers is none of those. It is unreachable
// by keyboard, announces as nothing, and on a phone -- where files are picked
// and never dropped -- it does nothing at all.
//
// The accept list and the size are stated here and enforced on the server
// regardless. Everything a browser checks is a courtesy to whoever is
// uploading, and a request that bypassed the page carries none of it.
//
// It publishes root, label, dropzone, input, prompt, list, hint and message.
type FileUploadProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejection comes from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above the area.
	Label string
	// Accept is what the picker offers, as extensions or MIME types:
	// ".pdf,image/*". Empty offers everything, which on a phone means the
	// whole file system rather than the camera roll.
	Accept string
	// Multiple takes more than one file.
	Multiple bool
	// Prompt is the sentence inside the area. Empty says what the area does,
	// which is the one line that has to be there.
	Prompt string
	// Hint is the line under it, and is where the limits belong: "PDF up to
	// 10 MB".
	Hint string
	// MaxSizeLabel is the size written into the area itself, for the limit
	// that matters most and is read least.
	MaxSizeLabel string
	// UploadURL posts each file as it is chosen instead of waiting for the
	// form. Empty leaves the files to the form the field is in.
	UploadURL string
	// ListID is the id of the region the uploaded files are listed in.
	// Empty means Name with "-files" after it.
	ListID string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
}

// PromptText is the sentence inside the area.
func (p FileUploadProps) PromptText() string {
	if p.Prompt != "" {
		return p.Prompt
	}
	if p.Multiple {
		return "Drop files here, or choose them"
	}
	return "Drop a file here, or choose one"
}

// Region is the id of the list of chosen files.
func (p FileUploadProps) Region() string {
	if p.ListID != "" {
		return p.ListID
	}
	return p.Name + "-files"
}

// FieldName is what the control submits under. A multiple file input submits
// once per file, and the bracket is what tells a form parser to keep them as a
// list rather than to overwrite with the last one.
func (p FileUploadProps) FieldName() string {
	if p.Multiple {
		return p.Name + "[]"
	}
	return p.Name
}

// Message is the rejection for this field, and empty when there is none.
func (p FileUploadProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the field.
func (p FileUploadProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p FileUploadProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: FileUploadBehavior}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p FileUploadProps) PartNames() []string {
	return []string{"root", "label", "dropzone", "input", "prompt", "list", "hint", "message"}
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

	{{-- The area is a label for the input, so clicking anywhere in it opens the
	     picker with no handler at all, and the input itself stays where a
	     keyboard and a screen reader find it. --}}
	<label
		data-part="dropzone"
		class="{{ .PartClass("dropzone", "file-upload") }}"
		for="{{ .Name }}"
		data-dropzone
		@attributes(.PartAttrs("dropzone"))
	>
		<input
			data-part="input"
			class="{{ .PartClass("input", "sr-only") }}"
			type="file"
			id="{{ .Name }}"
			name="{{ .FieldName() }}"
			@attributes(.PartAttrs("input"))
			@if(.Accept != "")
				accept="{{ .Accept }}"
			@endif
			@if(.Multiple)
				multiple
			@endif
			@if(.UploadURL != "")
				hx-post="{{ .UploadURL }}"
				hx-encoding="multipart/form-data"
				hx-trigger="change"
				hx-target="#{{ .Region() }}"
				hx-swap="beforeend"
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

		<span
			data-part="prompt"
			class="{{ .PartClass("prompt", "file-upload-prompt") }}"
			@attributes(.PartAttrs("prompt"))
		>
			{{ .PromptText() }}
			@if(.MaxSizeLabel != "")
				<small>{{ .MaxSizeLabel }}</small>
			@endif
		</span>
	</label>

	{{-- The list is on the page before anything is chosen, and it is polite:
	     files appearing is not an interruption. --}}
	<ul
		data-part="list"
		class="{{ .PartClass("list", "file-upload-list") }}"
		id="{{ .Region() }}"
		aria-live="polite"
		@attributes(.PartAttrs("list"))
	></ul>

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
