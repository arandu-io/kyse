//go:build kyse

package components

@go
// ResponsiveImageProps is one picture in several forms, with the browser
// choosing.
//
// The choosing is the point, and it is not a decision a script can make well:
// the browser knows the pixel density, the viewport, the layout width at the
// moment the image is discovered, and which formats it can decode -- and it
// knows all of that before layout, which is why the image starts downloading
// before any script has run.
//
// Three things travel together here and are usually separated to the image's
// cost. The width and height are always written, because without them the
// space is not reserved and everything under the image jumps when it arrives.
// The sizes attribute says how wide the image will be drawn, without which the
// browser assumes the full viewport and downloads the largest source for a
// thumbnail. And loading is eager for the picture at the top of the page and
// lazy for the rest -- a lazy hero is a hero that starts loading after the
// page has decided it is on screen.
//
// It publishes root, source and image.
type ResponsiveImageProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// URL is the fallback, and the one every browser can read: a JPEG or a
	// PNG. It is what the img element points at.
	URL string
	// Alt is what the picture shows, for somebody who cannot see it. Empty is
	// a decorative picture, which is a claim and not an omission.
	Alt string
	// Width and Height are the intrinsic size in pixels. They are not the
	// drawn size -- CSS decides that -- they are the ratio, and writing them
	// is what reserves the space.
	Width  int
	Height int
	// Sources are the alternatives the browser picks from, in the order it
	// tries them: the most preferred format first.
	Sources []ImageSource
	// SrcSet is the fallback's own set of widths: "cover-800.jpg 800w,
	// cover-1600.jpg 1600w". Empty offers the one file.
	SrcSet string
	// Sizes is how wide the image is drawn, as media conditions:
	// "(min-width: 60rem) 40rem, 100vw". Empty means the full viewport, which
	// is right for a hero and wrong for everything else.
	Sizes string
	// Eager loads the image immediately instead of when it approaches the
	// viewport. It is for the one picture that is on the screen when the page
	// opens, and for nothing else.
	Eager bool
	// Caption draws the picture inside a figure with the sentence under it.
	Caption string
}

// ImageSource is one alternative the browser may pick.
type ImageSource struct {
	// SrcSet is the file or files, with their widths or densities.
	SrcSet string
	// Type is the MIME type: "image/avif", "image/webp". It is what lets the
	// browser skip a source it cannot decode without downloading it.
	Type string
	// Media is a condition the source applies under: "(max-width: 40rem)",
	// "(prefers-color-scheme: dark)". Empty applies always.
	Media string
}

// Loading is when the browser starts fetching.
func (p ResponsiveImageProps) Loading() string {
	if p.Eager {
		return "eager"
	}
	return "lazy"
}

// Priority is the fetch priority. The one picture that is on the screen at the
// start is worth asking for ahead of everything else; the rest are not.
func (p ResponsiveImageProps) Priority() string {
	if p.Eager {
		return "high"
	}
	return "auto"
}

// PartNames are the parts this component publishes.
func (p ResponsiveImageProps) PartNames() []string {
	return []string{"root", "image", "caption"}
}
@endgo

@if(.Caption != "")
	<figure
		data-part="root"
		class="{{ .RootClass("figure") }}"
		@attributes(.RootAttrs())
	>
		<picture>
			@foreach(.Sources as source)
				<source
					srcset="{{ source.SrcSet }}"
					@if(source.Type != "")
						type="{{ source.Type }}"
					@endif
					@if(source.Media != "")
						media="{{ source.Media }}"
					@endif
					@if(.Sizes != "")
						sizes="{{ .Sizes }}"
					@endif
				>
			@endforeach
			<img
				data-part="image"
				class="{{ .PartClass("image", "figure-media") }}"
				src="{{ .URL }}"
				alt="{{ .Alt }}"
				loading="{{ .Loading() }}"
				decoding="async"
				fetchpriority="{{ .Priority() }}"
				@attributes(.PartAttrs("image"))
				@if(.Width > 0)
					width="{{ .Width }}"
				@endif
				@if(.Height > 0)
					height="{{ .Height }}"
				@endif
				@if(.SrcSet != "")
					srcset="{{ .SrcSet }}"
				@endif
				@if(.Sizes != "")
					sizes="{{ .Sizes }}"
				@endif
			>
		</picture>
		<figcaption
			data-part="caption"
			class="{{ .PartClass("caption", "figure-caption") }}"
			@attributes(.PartAttrs("caption"))
		>{{ .Caption }}</figcaption>
	</figure>
@endif
@if(.Caption == "")
	<picture
		data-part="root"
		@if(.RootClass() != "")
			class="{{ .RootClass() }}"
		@endif
		@attributes(.RootAttrs())
	>
		@foreach(.Sources as source)
			<source
				srcset="{{ source.SrcSet }}"
				@if(source.Type != "")
					type="{{ source.Type }}"
				@endif
				@if(source.Media != "")
					media="{{ source.Media }}"
				@endif
				@if(.Sizes != "")
					sizes="{{ .Sizes }}"
				@endif
			>
		@endforeach
		<img
			data-part="image"
			class="{{ .PartClass("image", "responsive-image") }}"
			src="{{ .URL }}"
			alt="{{ .Alt }}"
			loading="{{ .Loading() }}"
			decoding="async"
			fetchpriority="{{ .Priority() }}"
			@attributes(.PartAttrs("image"))
			@if(.Width > 0)
				width="{{ .Width }}"
			@endif
			@if(.Height > 0)
				height="{{ .Height }}"
			@endif
			@if(.SrcSet != "")
				srcset="{{ .SrcSet }}"
			@endif
			@if(.Sizes != "")
				sizes="{{ .Sizes }}"
			@endif
		>
	</picture>
@endif
