//go:build kyse

package components

@go
// CardProps is one item in a list: a title, a sentence, and where it goes.
//
// It does not take content. A component that wrapped arbitrary markup would
// have to receive it as a string of HTML, and a string of HTML is the one thing
// this view layer refuses to accept -- it is exactly where escaping stops being
// guaranteed by construction. So a card takes the fields it draws, and markup
// that needs a frame around it uses `class="card"` directly.
type CardProps struct {
	// Title is the heading, and the link text when Href is set.
	Title string
	// Description is the sentence under it.
	Description string
	// Href makes the title a link. Empty draws a plain card.
	Href string
	// Meta is the small grey line at the bottom: a date, an author, a count.
	Meta string
	// Badge draws a badge beside the title. Empty draws none.
	Badge string
	// BadgeVariant styles it. See BadgeProps.
	BadgeVariant string
}
@endgo

<article class="card p-5">
	<header class="flex items-start justify-between gap-3">
		<h3 class="font-semibold tracking-tight">
			@if(.Href != "")
				<a class="hover:underline" href="{{ .Href }}">{{ .Title }}</a>
			@endif
			@if(.Href == "")
				{{ .Title }}
			@endif
		</h3>
		@if(.Badge != "")
			{!! Badge(BadgeProps{Label: .Badge, Variant: .BadgeVariant}) !!}
		@endif
	</header>

	@if(.Description != "")
		<p class="text-muted-foreground mt-2 text-sm">{{ .Description }}</p>
	@endif
	@if(.Meta != "")
		<p class="text-muted-foreground mt-3 text-xs">{{ .Meta }}</p>
	@endif
</article>
