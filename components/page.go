package components

// Page is the small part of a screen's page data that an input asks about
// itself: the message validation left for it, and what was typed into it on the
// attempt that was rejected.
//
// It is not a component. It is a contract, declared here so that this library
// stays a library: framework/view.Page implements both methods, and every screen
// already embeds it, so a form writes
//
//	{!! components.Field(components.FieldProps{
//		Name: "email", Label: "Email", Type: "email", Page: .Page,
//	}) !!}
//
// and the field looks itself up by the Name it was already given.
//
// # Why the input asks rather than being told
//
// The alternative is what this replaced: Value and Error passed in beside Name,
// which writes the field's name three times in one call. Nothing checks that the
// three agree, and when they do not the failure is silent -- a field labelled
// "Confirm password" drawing the message belonging to "password", or drawing
// none at all, on the screen where a missing message is the whole complaint.
// Asked this way the name is written once and the other two cannot drift from
// it.
//
// A nil Page is the ordinary case for a screen with no form state behind it, and
// it draws no message and keeps whatever Value it was given. It is not an error:
// a component that panicked on it would take down a page over a field that has
// nothing to say.
type Page interface {
	// FieldError is the first message for an input, or empty when it was
	// accepted -- and empty for every input on a page nobody was rejected on,
	// which is nearly all of them.
	FieldError(name string) string
	// OldOr is what was typed into an input on the rejected attempt, falling
	// back to what the caller passed when nothing was typed.
	//
	// The fallback is what an edit form starts at, so a rejected edit comes back
	// carrying the change somebody was in the middle of making rather than
	// reverting to the stored row.
	OldOr(name, fallback string) string
}

// DialogPage is the page as the fields of one FormDialog see it: the messages
// and the typed values of a refused form, but only when the form refused was
// that dialog's.
//
//	components.Field(components.FieldProps{
//		Name: "reason", ID: "cancel-subscription-reason", Label: "Reason",
//		Page: components.DialogPage(.Page, "cancel-subscription"),
//	})
//
// # Why the fields need it
//
// A page answers by field name, and a page with three dialogs can have three
// fields called reason. A refusal from one of them came back with a message for
// reason, and every box called reason drew it -- in the two dialogs nobody
// submitted as well as in the one that was refused, along with the value typed
// into the one. The dialog that was refused is the one whose ID came back in
// FormDialogField, so this answers for that dialog's fields and is silent for
// every other: no message, and the value the field was given.
//
// A nil page answers nothing, like a nil Page given to a field directly.
func DialogPage(page Page, dialog string) Page {
	return dialogPage{page: page, dialog: dialog}
}

// dialogPage is a Page that answers only while its dialog was the one refused.
type dialogPage struct {
	page   Page
	dialog string
}

// refused is whether the form sent back to this page was this dialog's.
func (d dialogPage) refused() bool {
	return d.page != nil && d.dialog != "" && d.page.OldOr(FormDialogField, "") == d.dialog
}

// FieldError is the page's message for name when this dialog was refused.
func (d dialogPage) FieldError(name string) string {
	if !d.refused() {
		return ""
	}
	return d.page.FieldError(name)
}

// OldOr is what was typed into name when this dialog was refused, and the
// fallback otherwise.
func (d dialogPage) OldOr(name, fallback string) string {
	if !d.refused() {
		return fallback
	}
	return d.page.OldOr(name, fallback)
}
