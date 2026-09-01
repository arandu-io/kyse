//go:build kyse

package components

@go
// EmptyProps is what a list draws when it has nothing in it.
//
// It exists because the alternative -- a blank area -- reads as a page that
// failed rather than a page with nothing to show, and the difference costs a
// heading and a sentence. The action is what turns it from an apology into the
// obvious next step.
// It publishes root, header, title, message and action.
type EmptyProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Title is the heading: what is not here.
	Title string
	// Message is the sentence under it.
	Message string
	// ActionLabel and ActionURL draw a button. Both empty draws none.
	ActionLabel string
	ActionURL   string
}

// PartNames are the parts this component publishes.
func (p EmptyProps) PartNames() []string {
	return []string{"root", "header", "title", "message", "action"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("empty") }}"
	@attributes(.RootAttrs())
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
			@attributes(.PartAttrs("title"))
		>{{ .Title }}</h2>
		@if(.Message != "")
			<p
				data-part="message"
				@if(.PartClass("message") != "")
					class="{{ .PartClass("message") }}"
				@endif
				@attributes(.PartAttrs("message"))
			>{{ .Message }}</p>
		@endif
	</header>
	@if(.ActionURL != "")
		<section>
			<a
				data-part="action"
				class="{{ .PartClass("action", "btn") }}"
				data-variant="outline"
				href="{{ .ActionURL }}"
				@attributes(.PartAttrs("action"))
			>{{ .ActionLabel }}</a>
		</section>
	@endif
</div>
