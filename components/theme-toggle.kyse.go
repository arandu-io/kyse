//go:build kyse

package components

@go
// ThemeAccents are the accent colours the toggle offers, in the order they are
// read.
//
// The list is here because the swatches are markup and the server draws markup:
// six buttons, each carrying the name it sets. The script that applies a theme
// keeps a list of its own and checks a name against it before anything reaches
// the document, so a name that is on one list and not the other sets nothing --
// it draws a swatch that does not answer, never an unchecked value applied to
// the page.
var ThemeAccents = []string{"slate", "blue", "green", "amber", "rose", "violet"}
@endgo

{{-- The theme is client state -- which colours somebody prefers on this device
     -- so the server neither sets it nor is told about it. It lives in
     localStorage and on the html element.

     It takes no props, and still earns its place: the markup is a button, a menu
     of six swatches and the ARIA that ties them together, and repeating that on
     every layout is how two of them drift apart.

     Nothing here is drawn twice and chosen between at render time. The glyph and
     the word are both in the document, and the stylesheet shows the one that
     matches the html element -- which the theme script sets before the body is
     parsed, so the first paint is already right and there is no frame to hide.
     The attributes below carry data and never an expression: which accent is in
     force is compared against the html element by the client script, and the
     swatch colour is a custom property the server names. --}}

<div class="dropdown-menu">
	<button type="button" class="btn" data-variant="ghost" data-size="icon"
		aria-label="Change the theme" aria-haspopup="menu" aria-expanded="false">
		<span aria-hidden="true" data-theme-glyph="light">○</span>
		<span aria-hidden="true" data-theme-glyph="dark">◑</span>
	</button>

	<div data-popover aria-hidden="true" class="w-44">
		<div role="menu">
			{{-- The word names what the click will do, so it is the opposite of
			     the state it is shown in. --}}
			<button type="button" role="menuitem" data-theme-dark>
				<span data-theme-glyph="light">Dark</span>
				<span data-theme-glyph="dark">Light</span>
			</button>

			<hr role="separator">

			@for(at := 0; at < len(ThemeAccents); at++)
				<button type="button" role="menuitem" data-theme-accent="{{ ThemeAccents[at] }}">
					{{-- The colour is a class and not a style attribute. The policy an
					     Arandu application serves its pages under is style-src 'self',
					     which refuses a style attribute as surely as it refuses an
					     inline script -- the swatch would be a colourless dot, and the
					     console would say so in a sentence about inline styles that
					     reads like a different problem. The class is defined once
					     beside the rest of the theme rules. --}}
					<span class="size-3 rounded-full border swatch-{{ ThemeAccents[at] }}" data-theme-swatch></span>
					<span>{{ ThemeAccents[at] }}</span>
				</button>
			@endfor
		</div>
	</div>
</div>
