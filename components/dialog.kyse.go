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
// # In a row of a table
//
// A question about one row is drawn in that row's cell, next to the control that
// opens it: DialogTrigger with the dialog's ID, and a Dialog carrying the row in
// Fields. A <dialog> is shown in the top layer whatever cell it is drawn in, so
// the table does not clip it. Its words are not the row's: it writes
// data-search="false", and a complete table leaves it out of what it searches
// and orders by. A table the server draws again replaces the dialog with the
// rest of the region, so a dialog opened while the region is swapped is closed
// with it; ask the question from a page the table does not redraw when the list
// updates by itself.
//
// It confirms rather than composes. A dialog that could hold anything would
// have to take markup as a string, which is where escaping stops being
// guaranteed -- and "are you sure" is what a dialog is for nine times in ten.
// The tenth, a form over the page it belongs to, is FormDialog, which takes the
// output of other components rather than a string.
type DialogProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what a button targets to open this:
	//
	//	<button data-dialog-open="ID">Delete</button>
	//
	// Not an inline onclick. The policy this framework serves under is
	// script-src 'self' with no 'unsafe-inline', and that blocks an inline
	// handler exactly as it blocks an inline script -- so the handler works
	// while somebody is building and fails when the headers go on. The
	// attribute is read by the shipped script, which looks the id up and
	// never evaluates anything.
	ID string
	// Title is the question.
	Title string
	// Message is the consequence, spelled out. A destructive dialog that does
	// not say what is lost is a dialog people dismiss without reading.
	Message string

	// ConfirmLabel is the button that goes through with it. Empty says
	// kyse.dialog.confirm from the Page, or "Confirm".
	ConfirmLabel string
	// ConfirmVariant styles it. Use "destructive" for anything that deletes.
	ConfirmVariant string
	// Action and Method are the form the confirm button submits. PUT, PATCH and
	// DELETE are carried by a POST with a hidden _method field, because browsers
	// cannot submit those methods directly. A dialog that deletes over GET is a
	// dialog a crawler can fire.
	Action string
	Method string

	// Fields are hidden fields the confirm form posts beside the token: the key
	// of the row the question is about, a return address. A dialog drawn once
	// for a row of a table carries the row in them, so the one form posts the
	// row it was opened for and nothing else.
	Fields []DialogField
	// FormID makes the confirm button submit a form that is declared elsewhere on
	// the page, by its id, instead of one of its own: the bulk bar of a data
	// table, whose selection rides in hidden fields the table keeps. The dialog
	// then draws no form to confirm with, and Action, Method and Fields are the
	// other form's.
	FormID string
	// ConfirmName and ConfirmValue are the name and value the confirm button
	// submits with, for a form with several buttons that tell the server which
	// was pressed.
	ConfirmName  string
	ConfirmValue string

	// CancelLabel is the way out. Empty says kyse.dialog.cancel from the Page,
	// or "Cancel".
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

	// Page is the screen's own view.Page. The confirm form posts the CSRF
	// token it carries, the two buttons read their words from it when their
	// labels are empty, so a dialog drawn from a view writes Page: .Page and
	// never copies the token by hand. Nil posts Token, or none.
	Page Page
	// Token is the CSRF token, for a dialog with no Page to ask. Set, it wins
	// over the page's. It is written only when the form posts; see
	// SendsToken.
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

// DialogField is a hidden field of the form a dialog confirms with.
type DialogField struct {
	Name  string
	Value string
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
		return say(p.Page, "dialog.confirm", "Confirm")
	}
	return p.ConfirmLabel
}

// Cancel is the label of the way out.
func (p DialogProps) Cancel() string {
	if p.CancelLabel == "" {
		return say(p.Page, "dialog.cancel", "Cancel")
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

// SendsToken is whether the confirm form carries the CSRF token: only when it
// posts. A form submitted any other way puts its fields in the address, and an
// address is kept by the history, the access log and every proxy between --
// so the token is left out rather than handed to all of them.
func (p DialogProps) SendsToken() bool { return p.FormMethod() == "post" }

// FormToken is the token the confirm form posts: Token when it is set, and the
// page's otherwise.
func (p DialogProps) FormToken() string { return tokenFor(p.Page, p.Token) }

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
	@if(.Alert)
		role="alertdialog"
	@endif
	@if(.DescribedBy() != "")
		aria-describedby="{{ .DescribedBy() }}"
	@endif
	data-search="false"
	@attributes(.RootAttrs())
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
					@if(.Alert)
						autofocus
					@endif
					@attributes(.PartAttrs("cancel"))
				>{{ .Cancel() }}</button>
			</form>
			@if(.FormID != "")
				<button
					data-part="confirm"
					type="submit"
					form="{{ .FormID }}"
					class="{{ .PartClass("confirm", "btn") }}"
					data-variant="{{ .ConfirmVariant }}"
					@if(.ConfirmName != "")
						name="{{ .ConfirmName }}"
						value="{{ .ConfirmValue }}"
					@endif
					@attributes(.PartAttrs("confirm"))
				>{{ .Confirm() }}</button>
			@endif
			@if(.FormID == "")
				<form method="{{ .FormMethod() }}" action="{{ .Action }}">
					@if(.MethodOverride() != "")
						<input type="hidden" name="_method" value="{{ .MethodOverride() }}">
					@endif
					@if(.SendsToken())
						<input type="hidden" name="_token" value="{{ .FormToken() }}">
					@endif
					@foreach(.Fields as field)
						<input type="hidden" name="{{ field.Name }}" value="{{ field.Value }}">
					@endforeach
					<button
						data-part="confirm"
						type="submit"
						class="{{ .PartClass("confirm", "btn") }}"
						data-variant="{{ .ConfirmVariant }}"
						@if(.ConfirmName != "")
							name="{{ .ConfirmName }}"
							value="{{ .ConfirmValue }}"
						@endif
						@attributes(.PartAttrs("confirm"))
					>{{ .Confirm() }}</button>
				</form>
			@endif
		</footer>
	</article>
</dialog>
