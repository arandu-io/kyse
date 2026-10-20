package components

import "github.com/arandu-io/hesape/translation"

// Page is the small part of a screen's page data that an input asks about
// itself: the message validation left for it, and what was typed into it on the
// attempt that was rejected.
//
// It is not a component. It is a contract, declared here so that this library
// stays a library: hesape/view.Page implements both methods, and every screen
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

// TokenPage is a page that carries the CSRF token of the session it is drawn
// for. hesape/view.Page implements it, as CSRFToken.
//
// A component that posts asks its Page for it, so a confirmation drawn from a
// view writes Page: .Page and the form carries the token without anybody
// copying it into a second field:
//
//	{!! components.Dialog(components.DialogProps{
//		ID: "delete-invoice", Title: "Delete this invoice?",
//		Action: "/invoices/7", Method: "delete", Page: .Page,
//	}) !!}
//
// Dialog, FormDialog, Table and DataTable ask; a Token written on the props
// still wins, for the page that has no token to offer.
//
// It is apart from Page because a field asks for its message and what was
// typed, and nothing else: a page written for a test of one field is not made
// to invent a token it has no use for.
type TokenPage interface {
	// CSRFToken is the token every form of this page posts, or empty when the
	// page has none.
	CSRFToken() string
}

// FlashPage is a page that carries what the request before it said about what
// it did. hesape/view.Page implements it, as Flashed.
//
// It is the page half of Context.With: an action that worked says so under a
// key, redirects, and the page it lands on reads the sentence once. A toast
// asks for it, so a layout draws the confirmation with one call:
//
//	{!! components.Toast(components.ToastProps{Page: .Page}) !!}
//
// It is apart from Page for the reason TokenPage is.
type FlashPage interface {
	// Flashed is the sentence left under key by the request that redirected
	// here, or empty.
	Flashed(key string) string
}

// TranslatorPage is a page that can say a component's own words -- "Cancel",
// "Next page", "Close" -- in the language the page is served in.
//
// A component that draws a word nobody wrote asks for it here before it falls
// back to English. The words are read from the application's own catalogue,
// under the group "kyse": lang/pt-BR/kyse.json holding
//
//	{"dialog": {"cancel": "Cancelar"}, "pagination": {"next": "Próxima"}}
//
// answers kyse.dialog.cancel and kyse.pagination.next. The package doc lists
// every key. A label written on the props still wins, and a key the catalogue
// does not carry -- or a page that is not a TranslatorPage -- draws the English
// it drew before.
//
// It is apart from Page for the reason TokenPage is. hesape/view.Page does not
// implement it: a screen's own page data does, by returning the translator it
// was wired with and the locale of the request.
type TranslatorPage interface {
	// Translator is the application's catalogue, or nil for none.
	Translator() *translation.Translator
	// Locale is the language of the request being drawn. Empty asks the
	// translator for its own.
	Locale() string
}

// say is a component's own word for key: the page's catalogue line for
// kyse.<key> when it has one, and english otherwise.
func say(page Page, key, english string) string {
	speaker, ok := page.(TranslatorPage)
	if !ok {
		return english
	}
	translator := speaker.Translator()
	if translator == nil {
		return english
	}
	key = "kyse." + key
	locale := speaker.Locale()
	if !translator.Has(locale, key) {
		return english
	}
	return translator.Get(locale, key, nil)
}

// tokenFor is the token a form posts: the one the caller wrote, or else the
// one its page carries.
func tokenFor(page Page, given string) string {
	if given != "" {
		return given
	}
	if carrier, ok := page.(TokenPage); ok {
		return carrier.CSRFToken()
	}
	return ""
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
