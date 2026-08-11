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
