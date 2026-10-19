//go:build kyse

package components

import (
	"html/template"
	"strings"

	"github.com/arandu-io/kyse/icons"
)

@go
// FormDialogBehavior is the name the client behaviour is registered under with
// arandu.ui.define. It is written only on a dialog that has to open by itself;
// see FormDialogProps.Opened.
const FormDialogBehavior = "form-dialog"

// FormDialogField is the hidden field a form dialog posts its own ID in.
//
// It is what lets the page that a refused form is sent back to tell which of
// its dialogs was refused: the framework carries what was typed back to that
// page, this field among it, and the dialog whose ID it names opens itself.
// The handler receives it with the rest of the form and has nothing to do with
// it.
const FormDialogField = "_dialog"

// FormDialogProps is a form in a dialog over the page it belongs to: "New
// discount" over the list of discounts, "Refund order" over the order.
//
// It is a <dialog>, opened with showModal() by the script the framework
// serves (Hesape's ui.js, v0.66.0 or later), so the backdrop, the focus trap
// and Escape come from the element, and focus goes back to the trigger when it
// closes. It is opened by a DialogTrigger, which is a link to the page that
// holds the same form, so a browser that never runs the script follows the
// link and gets the form as a page.
//
// # Submitting it, and being refused
//
// The form posts like any other form -- no fragment, no script -- to the same
// action the full page posts to. Accepted, the handler redirects as it always
// does. Refused, the framework sends the person back where they posted from,
// which is the page this dialog sits on, carrying the messages and what was
// typed. The page is drawn again, the fields inside Body find their own
// messages through the Page they were given, and the dialog opens itself: it
// posts its ID in FormDialogField, the ID comes back with the rest of the typed
// values, and Opened answers true for the dialog it names. Nothing is added to
// the handler, the redirect or the address.
//
// A dialog that opens itself is drawn with the open attribute, so its form is
// readable on a page where the script never ran, and the script turns it into
// a modal and puts the cursor on the first refused field.
//
// # Two dialogs on one page
//
// A page answers the fields by name, so two dialogs that each ask for a reason
// both draw the message of the one that was refused, and the reason typed in
// it. The fields of a dialog are given DialogPage(page, ID) instead of the
// page, which answers only while that dialog was the one refused, and an ID of
// their own, so the two boxes called reason are two ids:
//
//	components.Field(components.FieldProps{
//		Name: "reason", ID: "cancel-subscription-reason", Label: "Reason",
//		Page: components.DialogPage(.Page, "cancel-subscription"),
//	})
//
// A value the form carries without asking -- the subscription being cancelled
// -- is a HiddenInput in Body, given the same DialogPage.
//
// # Leaving it
//
// Cancel and the close button submit a second, empty form whose method is
// "dialog", which the browser answers by closing the dialog -- with or without
// the script, and without validating or sending anything. They are not the
// form's own buttons, so pressing Enter in a field submits the form rather
// than closing it.
//
// # Confirming something that cannot be undone
//
// Alert and SubmitVariant "destructive" are the confirmation with a form in
// it: "Cancel subscription" with the choice of when, "Delete" with the name
// typed again. A confirmation with nothing to fill in is Dialog.
//
// It publishes root, content, header, title, description, body, footer,
// cancel, submit and close.
type FormDialogProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what a DialogTrigger names to open this, and what the title, the
	// description and the form hang their ids off. Two on a page need two.
	ID string
	// Title is the heading, and the name assistive technology announces when
	// the dialog opens.
	Title string
	// Description is the sentence under the title, and may be empty.
	Description string
	// Body is what the form holds, in order: each one the output of another
	// component -- a Field, a Select, a Password -- or of a view. It is markup,
	// so it never holds a string assembled around a value somebody typed; that
	// is what makes writing it unescaped safe.
	Body []template.HTML

	// Action and Method are where the form posts and how. PUT, PATCH and
	// DELETE are carried by a POST with a hidden _method field, because
	// browsers cannot submit those methods directly. Use the action the full
	// page posts to, so a refusal lands where the person came from either way.
	Action string
	Method string
	// Token is the CSRF token, written only when the form posts.
	Token string
	// Multipart encodes the form as multipart/form-data, for a form with a
	// file in it.
	Multipart bool

	// SubmitLabel is the button that sends the form. Empty means "Save".
	SubmitLabel string
	// SubmitVariant styles it. Use "destructive" for one that deletes or
	// cancels something.
	SubmitVariant string
	// CancelLabel is the way out. Empty means "Cancel".
	CancelLabel string
	// CloseLabel names the close button in the corner, which is an icon and has
	// no text of its own. Empty means "Close".
	CloseLabel string

	// Size is how wide it gets from the small breakpoint up: "sm", "md" or
	// "lg". Empty is "md". Below the small breakpoint the dialog takes the
	// whole screen, because a phone has no room around a form to leave
	// showing; an alert stays a card, because a question is short.
	Size string
	// Alert draws this as an alert dialog: the role tells assistive technology
	// the question has to be answered, and the surface is the alert one.
	Alert bool

	// Page is where the dialog asks whether the form that was just refused was
	// this one. A nil Page never opens it by itself.
	Page Page
	// Open draws it open regardless, for a page that decides by itself -- one
	// whose handler sends the person back with a sentence rather than with
	// what they typed.
	Open bool
}

// Opened is whether the dialog opens by itself when the page is drawn: when
// Open says so, or when the form the page was sent back with was this one.
func (p FormDialogProps) Opened() bool {
	if p.Open {
		return true
	}
	if p.Page == nil || p.ID == "" {
		return false
	}
	return p.Page.OldOr(FormDialogField, "") == p.ID
}

// RootAttrs are the outermost element's attributes, with the client behaviour
// filled in when the dialog opens by itself and the caller named no behaviour
// of their own. A dialog that waits for its trigger needs none: the trigger is
// read by the same script without one.
func (p FormDialogProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" && p.Opened() {
		p.Behavior = Behavior{Name: FormDialogBehavior, Props: map[string]any{"open": true}}
	}
	return p.ComponentProps.RootAttrs()
}

// Surface is the stylesheet this is drawn on.
func (p FormDialogProps) Surface() string {
	if p.Alert {
		return "alert-dialog"
	}
	return "dialog"
}

// SizeName is the size in force, written as data-size on a dialog that is not
// an alert. The alert surface reads data-size="sm" as a different layout,
// with its two buttons side by side in a grid, and the width is already
// WidthClass's.
func (p FormDialogProps) SizeName() string {
	switch p.Size {
	case "sm", "lg":
		return p.Size
	default:
		return "md"
	}
}

// WidthClass is the width from the small breakpoint up. Each branch is
// written out whole, because a class assembled at run time is a name no rule
// was compiled for.
func (p FormDialogProps) WidthClass() string {
	switch p.SizeName() {
	case "sm":
		return "sm:max-w-sm"
	case "lg":
		return "sm:max-w-2xl"
	default:
		return "sm:max-w-lg"
	}
}

// SheetClass takes the whole screen below the small breakpoint, and is empty
// for an alert.
func (p FormDialogProps) SheetClass() string {
	if p.Alert {
		return ""
	}
	return "max-sm:inset-0 max-sm:translate-x-0 max-sm:translate-y-0 max-sm:h-dvh max-sm:max-h-none max-sm:w-full max-sm:max-w-none max-sm:rounded-none max-sm:ring-0"
}

// FormID is the id of the form, for a control drawn outside it that submits
// it with the form attribute.
func (p FormDialogProps) FormID() string { return p.ID + "-form" }

// DismissID is the id of the empty form Cancel and the close button submit.
func (p FormDialogProps) DismissID() string { return p.ID + "-dismiss" }

// DescribedBy names the description, and is empty when there is none.
func (p FormDialogProps) DescribedBy() string {
	if p.Description == "" {
		return ""
	}
	return p.ID + "-description"
}

// Submit is the label of the button that sends the form.
func (p FormDialogProps) Submit() string {
	if p.SubmitLabel == "" {
		return "Save"
	}
	return p.SubmitLabel
}

// Cancel is the label of the way out.
func (p FormDialogProps) Cancel() string {
	if p.CancelLabel == "" {
		return "Cancel"
	}
	return p.CancelLabel
}

// Close is the accessible name of the close button.
func (p FormDialogProps) Close() string {
	if p.CloseLabel == "" {
		return "Close"
	}
	return p.CloseLabel
}

// FormMethod is the method a browser submits: post unless the caller asked
// for get, and post as the carrier for PUT, PATCH and DELETE.
func (p FormDialogProps) FormMethod() string {
	if strings.EqualFold(p.Method, "get") {
		return "get"
	}
	return "post"
}

// MethodOverride is the method carried in the hidden _method field, or empty
// when the browser can submit the requested method directly.
func (p FormDialogProps) MethodOverride() string {
	method := strings.ToUpper(p.Method)
	switch method {
	case "PUT", "PATCH", "DELETE":
		return method
	default:
		return ""
	}
}

// SendsToken is whether the form carries the CSRF token: only when it posts,
// because a form submitted with get puts its fields in the address.
func (p FormDialogProps) SendsToken() bool { return p.FormMethod() == "post" }

// Encoding is the enctype, and is empty for an ordinary form.
func (p FormDialogProps) Encoding() string {
	if p.Multipart {
		return "multipart/form-data"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p FormDialogProps) PartNames() []string {
	return []string{"root", "content", "header", "title", "description", "body", "footer", "cancel", "submit", "close"}
}
@endgo

<dialog
	data-part="root"
	id="{{ .ID }}"
	class="{{ .RootClass(.Surface()) }}"
	aria-labelledby="{{ .ID }}-title"
	@if(.Alert)
		role="alertdialog"
	@else
		data-size="{{ .SizeName() }}"
	@endif
	@if(.DescribedBy() != "")
		aria-describedby="{{ .DescribedBy() }}"
	@endif
	@if(.Opened())
		open
	@endif
	@attributes(.RootAttrs())
>
	{{-- The form is the panel itself rather than something inside it, because
	     the stylesheet places the header, the section and the footer as
	     children of the panel, and a form wrapped around them would take them
	     out of its reach. --}}
	<form
		data-part="content"
		id="{{ .FormID() }}"
		method="{{ .FormMethod() }}"
		action="{{ .Action }}"
		@if(.Encoding() != "")
			enctype="{{ .Encoding() }}"
		@endif
		class="{{ .PartClass("content", .WidthClass(), .SheetClass()) }}"
		@attributes(.PartAttrs("content"))
	>
		<input type="hidden" name="_dialog" value="{{ .ID }}">
		@if(.MethodOverride() != "")
			<input type="hidden" name="_method" value="{{ .MethodOverride() }}">
		@endif
		@if(.SendsToken())
			<input type="hidden" name="_token" value="{{ .Token }}">
		@endif

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
			@if(.Description != "")
				<p
					data-part="description"
					@if(.PartClass("description") != "")
						class="{{ .PartClass("description") }}"
					@endif
					id="{{ .ID }}-description"
					@attributes(.PartAttrs("description"))
				>{{ .Description }}</p>
			@endif
		</header>

		<section
			data-part="body"
			class="{{ .PartClass("body", "-m-1 grid min-h-0 content-start gap-4 overflow-y-auto p-1") }}"
			@attributes(.PartAttrs("body"))
		>
			@foreach(.Body as block)
				{!! block !!}
			@endforeach
		</section>

		{{-- Submit is first in the document and drawn last: the first submit
		     button a form owns is the one Enter presses, and that has to be the
		     one that sends it. Cancel is owned by the dismiss form, so it is
		     never that button wherever it stands. --}}
		<footer
			data-part="footer"
			class="{{ .PartClass("footer", "flex-col sm:flex-row-reverse sm:justify-start") }}"
			@attributes(.PartAttrs("footer"))
		>
			<button
				data-part="submit"
				type="submit"
				class="{{ .PartClass("submit", "btn") }}"
				@if(.SubmitVariant != "")
					data-variant="{{ .SubmitVariant }}"
				@endif
				@attributes(.PartAttrs("submit"))
			>{{ .Submit() }}</button>
			<button
				data-part="cancel"
				type="submit"
				form="{{ .DismissID() }}"
				class="{{ .PartClass("cancel", "btn") }}"
				data-variant="outline"
				@if(.Alert)
					autofocus
				@endif
				@attributes(.PartAttrs("cancel"))
			>{{ .Cancel() }}</button>
		</footer>

		@if(!.Alert)
			<button
				data-part="close"
				type="submit"
				form="{{ .DismissID() }}"
				@if(.PartClass("close") != "")
					class="{{ .PartClass("close") }}"
				@endif
				aria-label="{{ .Close() }}"
				@attributes(.PartAttrs("close"))
			>{!! icons.X(icons.Props{}) !!}</button>
		@endif
	</form>
	<form method="dialog" id="{{ .DismissID() }}" hidden></form>
</dialog>
