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
//
// It publishes root, header, title, description and meta.
type CardProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
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

// PartNames are the parts this component publishes.
func (p CardProps) PartNames() []string {
	return []string{"root", "header", "title", "description", "meta"}
}
@endgo

<article
	data-part="root"
	class="{{ .RootClass("card p-5") }}"
	@attributes(.RootAttrs())
>
	<header
		data-part="header"
		class="{{ .PartClass("header", "flex items-start justify-between gap-3") }}"
		@attributes(.PartAttrs("header"))
	>
		<h3
			data-part="title"
			class="{{ .PartClass("title", "font-semibold tracking-tight") }}"
			@attributes(.PartAttrs("title"))
		>
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
		<p
			data-part="description"
			class="{{ .PartClass("description", "text-muted-foreground mt-2 text-sm") }}"
			@attributes(.PartAttrs("description"))
		>{{ .Description }}</p>
	@endif
	@if(.Meta != "")
		<p
			data-part="meta"
			class="{{ .PartClass("meta", "text-muted-foreground mt-3 text-xs") }}"
			@attributes(.PartAttrs("meta"))
		>{{ .Meta }}</p>
	@endif
</article>
