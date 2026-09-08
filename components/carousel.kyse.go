//go:build kyse

package components

import "strconv"

@go
// CarouselBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const CarouselBehavior = "carousel"

// CarouselProps is a row of slides that scrolls one at a time.
//
// The scrolling is CSS scroll-snap, so a drag, a trackpad swipe, a touch flick
// and the scrollbar all work before any script has run -- and the buttons are
// an addition on top rather than the only way through. A carousel whose
// movement is animated in script is a carousel that does nothing while the
// script loads and fights the finger when it does.
//
// It is the APG carousel and not a slideshow: nothing moves on its own. An
// automatic carousel takes the reading away from whoever is reading, and the
// pause button it then needs is the tell that it should not have moved. If a
// rotation is wanted, it is the caller's to start and the reader's to stop.
//
// It publishes root, viewport, slide, previous, next, tabs and tab.
type CarouselProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the slides and the dots hang their ids off. Two carousels on
	// one page need two.
	ID string
	// Label names the carousel. Without it it announces with no subject.
	Label string
	// Slides are the panels, in the order they scroll.
	Slides []CarouselSlide
	// PreviousLabel and NextLabel name the two buttons. Empty says "Previous
	// slide" and "Next slide".
	PreviousLabel string
	NextLabel     string
	// Dots draws the row of markers under it, each one a way to jump to a
	// slide. They are a tablist, so the arrow keys move along them.
	Dots bool
	// Vertical scrolls up and down instead of left and right.
	Vertical bool
}

// CarouselSlide is one panel.
type CarouselSlide struct {
	// Title is the heading on it. Empty draws none.
	Title string
	// Caption is the sentence under the heading.
	Caption string
	// ImageURL is the picture.
	ImageURL string
	// Alt is what the picture shows, for somebody who cannot see it.
	Alt string
	// URL makes the whole slide a link.
	URL string
	// Label names the slide in the dots: "Slide 1 of 4" is what it falls back
	// to, and a name is better.
	Label string
}

// SlideID is the id of one slide, which its dot points at.
func (p CarouselProps) SlideID(at int) string {
	return p.ID + "-slide-" + strconv.Itoa(at)
}

// DotID is the id of one dot, which its slide points back at.
func (p CarouselProps) DotID(at int) string {
	return p.ID + "-dot-" + strconv.Itoa(at)
}

// SlideName is what one slide is called: what the caller wrote, or its place
// in the run.
func (p CarouselProps) SlideName(slide CarouselSlide, at int) string {
	if slide.Label != "" {
		return slide.Label
	}
	return "Slide " + strconv.Itoa(at+1) + " of " + strconv.Itoa(len(p.Slides))
}

// PreviousName and NextName are what the two buttons are called.
func (p CarouselProps) PreviousName() string {
	if p.PreviousLabel != "" {
		return p.PreviousLabel
	}
	return "Previous slide"
}

// NextName is what the forward button is called.
func (p CarouselProps) NextName() string {
	if p.NextLabel != "" {
		return p.NextLabel
	}
	return "Next slide"
}

// Direction is which way it scrolls.
func (p CarouselProps) Direction() string {
	if p.Vertical {
		return "vertical"
	}
	return "horizontal"
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p CarouselProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: CarouselBehavior, Props: map[string]any{"vertical": p.Vertical}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p CarouselProps) PartNames() []string {
	return []string{"root", "viewport", "slide", "media", "title", "caption", "previous", "next", "tabs", "tab"}
}
@endgo

<section
	data-part="root"
	class="{{ .RootClass("carousel") }}"
	id="{{ .ID }}"
	role="group"
	aria-roledescription="carousel"
	aria-label="{{ .Label }}"
	data-orientation="{{ .Direction() }}"
	@attributes(.RootAttrs())
>
	{{-- The viewport is what scrolls, and it is focusable and a group so a
	     keyboard can put the caret in it and use the arrow keys -- which is
	     what a scrollable region needs to be reachable at all. --}}
	<div
		data-part="viewport"
		class="{{ .PartClass("viewport", "carousel-viewport") }}"
		role="group"
		aria-label="{{ .Label }}"
		tabindex="0"
		@attributes(.PartAttrs("viewport"))
	>
		@for(at := 0; at < len(.Slides); at++)
			<div
				data-part="slide"
				class="{{ .PartClass("slide", "carousel-slide") }}"
				id="{{ .SlideID(at) }}"
				role="group"
				aria-roledescription="slide"
				aria-label="{{ .SlideName(.Slides[at], at) }}"
				@attributes(.PartAttrs("slide"))
				@if(.Dots)
					aria-labelledby="{{ .DotID(at) }}"
				@endif
			>
				@if(.Slides[at].ImageURL != "")
					<img
						data-part="media"
						class="{{ .PartClass("media", "carousel-media") }}"
						src="{{ .Slides[at].ImageURL }}"
						alt="{{ .Slides[at].Alt }}"
						loading="lazy"
						decoding="async"
						@attributes(.PartAttrs("media"))
					>
				@endif
				@if(.Slides[at].Title != "")
					<h3
						data-part="title"
						class="{{ .PartClass("title", "font-semibold") }}"
						@attributes(.PartAttrs("title"))
					>
						@if(.Slides[at].URL != "")
							<a class="hover:underline" href="{{ .Slides[at].URL }}">{{ .Slides[at].Title }}</a>
						@endif
						@if(.Slides[at].URL == "")
							{{ .Slides[at].Title }}
						@endif
					</h3>
				@endif
				@if(.Slides[at].Caption != "")
					<p
						data-part="caption"
						class="{{ .PartClass("caption", "text-muted-foreground text-sm") }}"
						@attributes(.PartAttrs("caption"))
					>{{ .Slides[at].Caption }}</p>
				@endif
			</div>
		@endfor
	</div>

	<button
		data-part="previous"
		class="{{ .PartClass("previous", "btn") }}"
		type="button"
		data-variant="outline"
		data-size="sm"
		data-carousel-step="-1"
		aria-label="{{ .PreviousName() }}"
		aria-controls="{{ .ID }}"
		@attributes(.PartAttrs("previous"))
	></button>

	<button
		data-part="next"
		class="{{ .PartClass("next", "btn") }}"
		type="button"
		data-variant="outline"
		data-size="sm"
		data-carousel-step="1"
		aria-label="{{ .NextName() }}"
		aria-controls="{{ .ID }}"
		@attributes(.PartAttrs("next"))
	></button>

	@if(.Dots)
		{{-- The dots are links to the slides, so they work with no script at
		     all: the fragment scrolls the viewport, which is the same movement
		     the buttons make. --}}
		<div
			data-part="tabs"
			class="{{ .PartClass("tabs", "carousel-dots") }}"
			role="group"
			aria-label="{{ .Label }}"
			@attributes(.PartAttrs("tabs"))
		>
			@for(at := 0; at < len(.Slides); at++)
				<a
					data-part="tab"
					class="{{ .PartClass("tab", "carousel-dot") }}"
					id="{{ .DotID(at) }}"
					href="#{{ .SlideID(at) }}"
					aria-label="{{ .SlideName(.Slides[at], at) }}"
					@attributes(.PartAttrs("tab"))
				></a>
			@endfor
		</div>
	@endif
</section>
