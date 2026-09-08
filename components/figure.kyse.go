//go:build kyse

package components

@go
// FigureProps is a picture with the sentence that says what it is.
//
// The caption is a <figcaption>, and that is the point of the element: the
// caption becomes the figure's accessible name, so the picture and the words
// about it are one thing rather than an image followed by a paragraph that
// happens to be underneath.
//
// The alternative text and the caption are both here and are not the same
// text. The caption says what the picture is for; the alternative text says
// what is in it, to somebody who cannot see it. A caption repeated as the
// alternative text is a picture described twice and described not at all.
//
// It publishes root, media, caption and credit.
type FigureProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ImageURL is the picture.
	ImageURL string
	// Alt is what is in the picture, for somebody who cannot see it. It is
	// empty only for a picture that carries nothing the caption does not
	// already say, which is the case a decorative image is.
	Alt string
	// Caption is the sentence about it.
	Caption string
	// Credit is the smaller line under the caption: a photographer, a source,
	// a licence.
	Credit string
	// CaptionOnTop draws the caption before the picture. The element requires
	// the caption to be the first or the last child, so those are the two
	// places it can go.
	CaptionOnTop bool
}

// PartNames are the parts this component publishes.
func (p FigureProps) PartNames() []string {
	return []string{"root", "media", "caption", "credit"}
}
@endgo

<figure
	data-part="root"
	class="{{ .RootClass("figure") }}"
	@attributes(.RootAttrs())
>
	@if(.CaptionOnTop && .Caption != "")
		<figcaption
			data-part="caption"
			class="{{ .PartClass("caption", "figure-caption") }}"
			@attributes(.PartAttrs("caption"))
		>
			{{ .Caption }}
			@if(.Credit != "")
				<span
					data-part="credit"
					class="{{ .PartClass("credit", "figure-credit") }}"
					@attributes(.PartAttrs("credit"))
				>{{ .Credit }}</span>
			@endif
		</figcaption>
	@endif
	@if(.ImageURL != "")
		<img
			data-part="media"
			class="{{ .PartClass("media", "figure-media") }}"
			src="{{ .ImageURL }}"
			alt="{{ .Alt }}"
			@attributes(.PartAttrs("media"))
		>
	@endif
	@if(!.CaptionOnTop && .Caption != "")
		<figcaption
			data-part="caption"
			class="{{ .PartClass("caption", "figure-caption") }}"
			@attributes(.PartAttrs("caption"))
		>
			{{ .Caption }}
			@if(.Credit != "")
				<span
					data-part="credit"
					class="{{ .PartClass("credit", "figure-credit") }}"
					@attributes(.PartAttrs("credit"))
				>{{ .Credit }}</span>
			@endif
		</figcaption>
	@endif
</figure>
