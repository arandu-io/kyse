//go:build kyse

package components

import "strings"

@go
// DialogProps is a confirmation: the thing about to happen, and the two ways
// out of it.
//
// It is a <dialog>, so the browser handles the backdrop, the focus trap and
// Escape. That is four behaviours nobody has to write, and the four that
// hand-built modals get wrong.
//
// None of the four is true unless it was opened with showModal(). One opened by
// setting the open attribute is not modal: the page behind it stays reachable,
// focus walks out of it, and Escape does nothing. And a modal dialog never
// closes on a click outside itself, which is the behaviour a confirmation
// wants and is not something this has to arrange.
//
// It confirms rather than composes. A dialog that could hold anything would
// have to take markup as a string, which is where escaping stops being
// guaranteed -- and "are you sure" is what a dialog is for nine times in ten.
type DialogProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what a button targets to open this: `onclick="ID.showModal()"`.
	ID string
	// Title is the question.
	Title string
	// Message is the consequence, spelled out. A destructive dialog that does
	// not say what is lost is a dialog people dismiss without reading.
	Message string

	// ConfirmLabel is the button that goes through with it. Empty means
	// "Confirm".
	ConfirmLabel string
	// ConfirmVariant styles it. Use "destructive" for anything that deletes.
	ConfirmVariant string
	// Action and Method are the form the confirm button submits. PUT, PATCH and
	// DELETE are carried by a POST with a hidden _method field, because browsers
	// cannot submit those methods directly. A dialog that deletes over GET is a
	// dialog a crawler can fire.
	Action string
	Method string

	// CancelLabel is the way out. Empty means "Cancel".
	CancelLabel string

	// Alert draws this as an alert dialog instead: the message is announced
	// along with the title the moment it opens, and the cursor starts on the
	// way out rather than on the way through. Use it for what cannot be undone.
	//
	// It is a field rather than a component of its own because the two differ
	// in the role, in what is announced and in where focus lands, and in
	// nothing else -- the question, the consequence and the pair of buttons are
	// the same markup. Two components would be two ways to ask "are you sure",
	// and within a year they would answer it differently.
	Alert bool

	// Token is the CSRF token. It is passed in rather than read from the page,
	// because a component does not receive the page.
	Token string
}

// Surface is the stylesheet this is drawn on: the alert one is wider, centres
// its heading on a narrow screen and has room for a figure above it.
func (p DialogProps) Surface() string {
	if p.Alert {
		return "alert-dialog"
	}
	return "dialog"
}

// DescribedBy names the consequence, and is empty when there is none. An
// aria-describedby pointing at an element that was never drawn describes
// nothing, which is worse than describing nothing at all: it reads as a
// description that failed to arrive.
func (p DialogProps) DescribedBy() string {
	if p.Message == "" {
		return ""
	}
	return p.ID + "-description"
}

// Confirm is the label of the button that goes through with it.
func (p DialogProps) Confirm() string {
	if p.ConfirmLabel == "" {
		return "Confirm"
	}
	return p.ConfirmLabel
}

// Cancel is the label of the way out.
func (p DialogProps) Cancel() string {
	if p.CancelLabel == "" {
		return "Cancel"
	}
	return p.CancelLabel
}

// FormMethod is the method a browser submits. It defaults to post because
// everything a dialog confirms changes something, and uses post as the carrier
// for PUT, PATCH and DELETE.
func (p DialogProps) FormMethod() string {
	if p.Method == "" {
		return "post"
	}
	if p.MethodOverride() != "" {
		return "post"
	}
	return strings.ToLower(p.Method)
}

// MethodOverride is the method carried in the hidden _method field, or empty
// when the browser can submit the requested method directly.
func (p DialogProps) MethodOverride() string {
	method := strings.ToUpper(p.Method)
	switch method {
	case "PUT", "PATCH", "DELETE":
		return method
	default:
		return ""
	}
}
// PartNames are the parts this component publishes.
func (p DialogProps) PartNames() []string {
	return []string{"root", "content", "header", "title", "message", "footer", "cancel", "confirm"}
}
@endgo

<dialog
	data-part="root"
	id="{{ .ID }}"
	class="{{ .RootClass(.Surface()) }}"
	aria-labelledby="{{ .ID }}-title"
	@attributes(.RootAttrs())
	@if(.Alert)
		role="alertdialog"
	@endif
	@if(.DescribedBy() != "")
		aria-describedby="{{ .DescribedBy() }}"
	@endif
>
	<article
		data-part="content"
		@if(.PartClass("content") != "")
			class="{{ .PartClass("content") }}"
		@endif
		@attributes(.PartAttrs("content"))
	>
		<header
			data-part="header"
			@if(.PartClass("header") != "")
				class="{{ .PartClass("header") }}"
			@endif
			@attributes(.PartAttrs("header"))
		>
			<h2
				data-part="title"
				@if(.PartClass("title") != "")
					class="{{ .PartClass("title") }}"
				@endif
				id="{{ .ID }}-title"
				@attributes(.PartAttrs("title"))
			>{{ .Title }}</h2>
			@if(.Message != "")
				<p
					data-part="message"
					id="{{ .ID }}-description"
					class="{{ .PartClass("message", "text-muted-foreground text-sm") }}"
					@attributes(.PartAttrs("message"))
				>{{ .Message }}</p>
			@endif
		</header>

		<footer
			data-part="footer"
			class="{{ .PartClass("footer", "flex justify-end gap-2") }}"
			@attributes(.PartAttrs("footer"))
		>
			<form method="dialog">
				<button
					data-part="cancel"
					type="submit"
					class="{{ .PartClass("cancel", "btn") }}"
					data-variant="outline"
					@attributes(.PartAttrs("cancel"))
					@if(.Alert)
						autofocus
					@endif
				>{{ .Cancel() }}</button>
			</form>
			<form method="{{ .FormMethod() }}" action="{{ .Action }}">
				@if(.MethodOverride() != "")
					<input type="hidden" name="_method" value="{{ .MethodOverride() }}">
				@endif
				<input type="hidden" name="_token" value="{{ .Token }}">
				<button
					data-part="confirm"
					type="submit"
					class="{{ .PartClass("confirm", "btn") }}"
					data-variant="{{ .ConfirmVariant }}"
					@attributes(.PartAttrs("confirm"))
				>{{ .Confirm() }}</button>
			</form>
		</footer>
	</article>
</dialog>
