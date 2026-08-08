//go:build kyse

package components

@go
// DialogProps is a confirmation: the thing about to happen, and the two ways
// out of it.
//
// It is a <dialog>, so the browser handles the backdrop, the focus trap and
// Escape. That is four behaviours nobody has to write, and the four that
// hand-built modals get wrong.
//
// It confirms rather than composes. A dialog that could hold anything would
// have to take markup as a string, which is where escaping stops being
// guaranteed -- and "are you sure" is what a dialog is for nine times in ten.
type DialogProps struct {
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
	// Action and Method are the form the confirm button submits. A dialog that
	// deletes over GET is a dialog a crawler can fire.
	Action string
	Method string

	// CancelLabel is the way out. Empty means "Cancel".
	CancelLabel string

	// Token is the CSRF token. It is passed in rather than read from the page,
	// because a component does not receive the page.
	Token string
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

// FormMethod defaults to post, because everything a dialog confirms changes
// something.
func (p DialogProps) FormMethod() string {
	if p.Method == "" {
		return "post"
	}
	return p.Method
}
@endgo

<dialog id="{{ .ID }}" class="dialog" aria-labelledby="{{ .ID }}-title">
	<article>
		<header>
			<h2 id="{{ .ID }}-title">{{ .Title }}</h2>
			@if(.Message != "")
				<p class="text-muted-foreground text-sm">{{ .Message }}</p>
			@endif
		</header>

		<footer class="flex justify-end gap-2">
			<form method="dialog">
				<button type="submit" class="btn" data-variant="outline">{{ .Cancel() }}</button>
			</form>
			<form method="{{ .FormMethod() }}" action="{{ .Action }}">
				<input type="hidden" name="_csrf" value="{{ .Token }}">
				<button type="submit" class="btn" data-variant="{{ .ConfirmVariant }}">{{ .Confirm() }}</button>
			</form>
		</footer>
	</article>
</dialog>
