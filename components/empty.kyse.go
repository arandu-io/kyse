//go:build kyse

package components

@go
// EmptyProps is what a list draws when it has nothing in it.
//
// It exists because the alternative -- a blank area -- reads as a page that
// failed rather than a page with nothing to show, and the difference costs a
// heading and a sentence. The action is what turns it from an apology into the
// obvious next step.
type EmptyProps struct {
	// Title is the heading: what is not here.
	Title string
	// Message is the sentence under it.
	Message string
	// ActionLabel and ActionURL draw a button. Both empty draws none.
	ActionLabel string
	ActionURL   string
}
@endgo

<div class="empty">
	<header>
		<h2>{{ .Title }}</h2>
		@if(.Message != "")
			<p>{{ .Message }}</p>
		@endif
	</header>
	@if(.ActionURL != "")
		<section>
			<a class="btn" data-variant="outline" href="{{ .ActionURL }}">{{ .ActionLabel }}</a>
		</section>
	@endif
</div>
