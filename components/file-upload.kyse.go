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
// # What is chosen
//
// A field that leaves its files to the form lists them as soon as they are
// picked or dropped: a row per file with a thumbnail for an image, its name and
// its size, and a count above the rows that is announced politely. Without it a
// person who picks a photo sees the screen not change and takes it that nothing
// happened. The list is drawn by the "file-upload" behaviour of the view runtime
// from what the input holds, out of the row this component writes in a
// template, so every class and every sentence is the component's and nothing
// is built from strings in the script. The thumbnail is a data: URL, which is
// what a policy of img-src 'self' data: allows.
//
// MaxFiles and MaxBytes make the list a warning as well. A pick over either is
// not refused -- the files stay in the input -- but the count says what is wrong
// before the submit, the rows over the size are marked, and the browser holds
// the submit with its own message. The server still checks both, because a
// request that skipped the page carries none of this.
//
// A field with an UploadURL posts each file as it is chosen and the server's
// answer fills its list, so it draws none of this.
//
// The sentences are the page's own words: kyse.file_upload.chosen_one,
// chosen_many, too_many, too_big_one and too_big_many in the application's
// catalogue, or the fields of the same name here, which win.
//
// It publishes root, label, dropzone, input, prompt, summary, list, file,
// thumbnail, file-name, file-size, hint and message. The last four are the row
// the behaviour draws once per file.
type FileUploadProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejection comes from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// ID is the id the markup carries, and the id the label, the message and
	// the hint hang off. Empty uses Name. Two forms on one page that ask for the
	// same Name -- two dialogs that each want a reason -- need one each, or the
	// second label points at the first box.
	ID string
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
	// MaxFiles is how many files one submit takes. Zero sets no limit. It is
	// what the page warns about; the server enforces it.
	MaxFiles int
	// MaxBytes is how large one file may be, in bytes. Zero sets no limit. It
	// is what the page warns about; the server enforces it.
	MaxBytes int64
	// ChosenOne is the count for a single file. Empty says
	// kyse.file_upload.chosen_one from the Page, or "1 file chosen".
	ChosenOne string
	// ChosenMany is the count for several files, with {n} where the number
	// goes. Empty says kyse.file_upload.chosen_many from the Page, or "{n}
	// files chosen".
	ChosenMany string
	// TooMany is the warning for more files than MaxFiles, with {max} where
	// the limit goes. Empty says kyse.file_upload.too_many from the Page, or
	// "Send up to {max} files at a time: choose again."
	TooMany string
	// TooBigOne is the warning for one file over MaxBytes, with {size} where
	// the limit goes, written in the language of the page. Empty says
	// kyse.file_upload.too_big_one from the Page, or "One file is over {size}:
	// replace it before sending."
	TooBigOne string
	// TooBigMany is the warning for several files over MaxBytes, with {n} and
	// {size}. Empty says kyse.file_upload.too_big_many from the Page, or "{n}
	// files are over {size}: replace them before sending."
	TooBigMany string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
}

// Previews is whether the files are listed by the page: they are, unless the
// field posts each one as it is chosen and lets the server's answer fill the
// list.
func (p FileUploadProps) Previews() bool { return p.UploadURL == "" }

// BehaviorProps are what the behaviour is given: the limits, and the sentences
// it writes the count and the warnings with, each the one a caller wrote, or
// else the page's line, or else English. Nil when the field lists nothing.
func (p FileUploadProps) BehaviorProps() map[string]any {
	if !p.Previews() {
		return nil
	}
	written := func(own, said string) string {
		if own != "" {
			return own
		}
		return said
	}
	props := map[string]any{
		"chosenOne":  written(p.ChosenOne, say(p.Page, "file_upload.chosen_one", "1 file chosen")),
		"chosenMany": written(p.ChosenMany, say(p.Page, "file_upload.chosen_many", "{n} files chosen")),
		"tooMany":    written(p.TooMany, say(p.Page, "file_upload.too_many", "Send up to {max} files at a time: choose again.")),
		"tooBigOne":  written(p.TooBigOne, say(p.Page, "file_upload.too_big_one", "One file is over {size}: replace it before sending.")),
		"tooBigMany": written(p.TooBigMany, say(p.Page, "file_upload.too_big_many", "{n} files are over {size}: replace them before sending.")),
	}
	if p.MaxFiles > 0 {
		props["maxFiles"] = p.MaxFiles
	}
	if p.MaxBytes > 0 {
		props["maxBytes"] = p.MaxBytes
	}
	return props
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
	return p.ElementID() + "-files"
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
		return p.ElementID() + "-error"
	}
	if p.Hint != "" {
		return p.ElementID() + "-hint"
	}
	return ""
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p FileUploadProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: FileUploadBehavior, Props: p.BehaviorProps()}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p FileUploadProps) PartNames() []string {
	return []string{"root", "label", "dropzone", "input", "prompt", "summary", "list", "file", "thumbnail", "file-name", "file-size", "hint", "message"}
}

// ElementID is the id the markup carries: ID, or Name when ID is empty.
func (p FileUploadProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
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
		for="{{ .ElementID() }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	{{-- The area is a label for the input, so clicking anywhere in it opens the
	     picker with no handler at all, and the input itself stays where a
	     keyboard and a screen reader find it. --}}
	<label
		data-part="dropzone"
		class="{{ .PartClass("dropzone", "file-upload w-full") }}"
		for="{{ .ElementID() }}"
		data-dropzone
		@attributes(.PartAttrs("dropzone"))
	>
		<input
			data-part="input"
			class="{{ .PartClass("input", "sr-only") }}"
			type="file"
			id="{{ .ElementID() }}"
			name="{{ .FieldName() }}"
			@if(.Accept != "")
				accept="{{ .Accept }}"
			@endif
			@if(.Multiple)
				multiple
			@endif
			@if(.UploadURL != "")
				ypy-post="{{ .UploadURL }}"
				ypy-encoding="multipart/form-data"
				ypy-trigger="change"
				ypy-target="#{{ .Region() }}"
				ypy-swap="beforeend"
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
			@attributes(.PartAttrs("input"))
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

	{{-- What was chosen, said in one sentence the screen reader announces
	     politely. It is on the page before anything is chosen and empty, and an
	     empty one is taken out of sight and not out of the page: a live region
	     that arrives together with its words is often not announced. --}}
	@if(.Previews())
		<p
			data-part="summary"
			class="{{ .PartClass("summary", "text-foreground text-sm font-medium empty:sr-only data-[over=true]:text-destructive") }}"
			role="status"
			aria-live="polite"
			aria-atomic="true"
			@attributes(.PartAttrs("summary"))
		></p>
	@endif

	{{-- The list is on the page before anything is chosen. When the server
	     fills it, as an upload as it is chosen does, it is polite: files
	     appearing is not an interruption. When the page draws it, the summary
	     above is what is announced, and a list that read out every row again on
	     every pick would only be noise. --}}
	<ul
		data-part="list"
		class="{{ .PartClass("list", "file-upload-list") }}"
		id="{{ .Region() }}"
		@if(!.Previews())
			aria-live="polite"
		@endif
		@attributes(.PartAttrs("list"))
	></ul>

	{{-- The row the behaviour copies once per file. It is a template so that it
	     is not on the page, and so its classes are in the source the stylesheet
	     is built from. The rules that colour a row over the limit read the
	     data-over the behaviour writes on it. --}}
	@if(.Previews())
		<template data-file-template>
			<li
				data-part="file"
				class="{{ .PartClass("file", "group flex min-w-0 items-center gap-3") }}"
				@attributes(.PartAttrs("file"))
			>
				<img
					data-part="thumbnail"
					class="{{ .PartClass("thumbnail", "bg-muted size-12 flex-none rounded-md object-cover") }}"
					alt=""
					width="48"
					height="48"
					hidden
					@attributes(.PartAttrs("thumbnail"))
				>
				<span
					data-part="file-name"
					class="{{ .PartClass("file-name", "text-foreground min-w-0 flex-1 truncate group-data-[over=true]:text-destructive") }}"
					@attributes(.PartAttrs("file-name"))
				></span>
				<span
					data-part="file-size"
					class="{{ .PartClass("file-size", "text-muted-foreground flex-none tabular-nums group-data-[over=true]:text-destructive") }}"
					@attributes(.PartAttrs("file-size"))
				></span>
			</li>
		</template>
	@endif

	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .ElementID() }}-error"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p
				data-part="hint"
				id="{{ .ElementID() }}-hint"
				class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("hint"))
			>{{ .Hint }}</p>
		@endif
	@endif
</div>
