//go:build kyse

package components

@go
// ContainerCardProps is a card that lays itself out from its own width rather
// than from the window's.
//
// That difference is the whole component. A card that reads the viewport is a
// card that has to know where on the page it was put: the same markup in a
// narrow sidebar and in a wide column gets the same layout, and one of the two
// is wrong. A container query asks the box how wide it is, so one card is
// correct in both places and neither the page nor the caller has to say which
// it is in.
//
// It publishes root, layout, media, body, title, description and footer.
type ContainerCardProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Title is the heading.
	Title string
	// Description is the sentence under it.
	Description string
	// URL makes the title a link.
	URL string
	// ImageURL is the picture beside or above the text.
	ImageURL string
	// Alt is what the picture shows, for somebody who cannot see it. It is
	// empty only for a picture the text already accounts for.
	Alt string
	// Meta is the small line at the bottom: a date, an author, a count.
	Meta string
	// ActionLabel and ActionURL draw a control in the footer.
	ActionLabel string
	ActionURL   string
	// Wide flips the layout at a narrower width, for a card whose text is
	// short enough to sit beside a picture sooner. It is a named step rather
	// than a length, because the threshold has to exist as a literal in the
	// stylesheet for a rule to be emitted for it.
	Wide bool
}

// Heading is the title, and the link text when URL is set.
func (p ContainerCardProps) Heading() string { return p.Title }

// PartNames are the parts this component publishes.
func (p ContainerCardProps) PartNames() []string {
	return []string{"root", "layout", "media", "body", "title", "description", "meta", "footer"}
}
@endgo

<article
	data-part="root"
	class="{{ .RootClass("container-card") }}"
	@if(.Wide)
		data-break="wide"
	@endif
	@attributes(.RootAttrs())
>
	<div
		data-part="layout"
		class="{{ .PartClass("layout", "container-card-layout") }}"
		@attributes(.PartAttrs("layout"))
	>
		@if(.ImageURL != "")
			<img
				data-part="media"
				class="{{ .PartClass("media", "container-card-media") }}"
				src="{{ .ImageURL }}"
				alt="{{ .Alt }}"
				@attributes(.PartAttrs("media"))
			>
		@endif

		<div
			data-part="body"
			class="{{ .PartClass("body", "container-card-body") }}"
			@attributes(.PartAttrs("body"))
		>
			<h3
				data-part="title"
				class="{{ .PartClass("title", "font-semibold leading-none") }}"
				@attributes(.PartAttrs("title"))
			>
				@if(.URL != "")
					<a class="hover:underline" href="{{ .URL }}">{{ .Heading() }}</a>
				@endif
				@if(.URL == "")
					{{ .Heading() }}
				@endif
			</h3>

			@if(.Description != "")
				<p
					data-part="description"
					class="{{ .PartClass("description", "text-muted-foreground text-sm") }}"
					@attributes(.PartAttrs("description"))
				>{{ .Description }}</p>
			@endif

			@if(.Meta != "")
				<p
					data-part="meta"
					class="{{ .PartClass("meta", "text-muted-foreground text-xs") }}"
					@attributes(.PartAttrs("meta"))
				>{{ .Meta }}</p>
			@endif

			@if(.ActionURL != "")
				<div
					data-part="footer"
					class="{{ .PartClass("footer", "container-card-footer") }}"
					@attributes(.PartAttrs("footer"))
				>
					<a class="btn" data-variant="outline" data-size="sm" href="{{ .ActionURL }}">{{ .ActionLabel }}</a>
				</div>
			@endif
		</div>
	</div>
</article>
