//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

{{-- The theme is client state -- which colours somebody prefers on this device
     -- so the server neither sets it nor is told about it. It lives in
     localStorage and on the html element.

     Three options, not a switch. A switch has to start somewhere, and wherever
     it starts is wrong for half the people arriving: somebody whose device is
     dark and whose last visit left it light gets a white page, and nothing on
     the control says why. "Auto" is the state where the page follows the
     device, it is the default, and choosing is how you stop following.

     There is no accent picker. The accent of a site is its identity, and
     offering a visitor five others is offering them five other brands.

     It takes no props, and still earns its place: the markup is a button, a
     menu of three radios and the ARIA that ties them together, and repeating
     that on every layout is how two of them drift apart.

     Nothing here is drawn twice and chosen between at render time. Both trigger
     icons are in the document and the stylesheet shows the one matching the
     html element, which the theme script sets before the body is parsed -- so
     the first paint is already right and there is no frame to hide. The
     attributes carry data and never an expression: which mode is in force is
     compared against the html element by the client script, which validates a
     name against its own list before anything reaches the document. --}}

<div class="dropdown-menu">
	<button type="button" class="btn" data-variant="ghost" data-size="icon"
		aria-label="Change the theme" aria-haspopup="menu" aria-expanded="false">
		{{-- The trigger shows what is in force rather than what was chosen. In
		     auto those are different things, and the sun and the moon are what
		     the page actually looks like; a third glyph here would name the
		     choice and say nothing about the screen. --}}
		<span aria-hidden="true" class="block size-[18px]" data-theme-glyph="light">{!! icons.Sun(icons.Props{}) !!}</span>
		<span aria-hidden="true" class="block size-[18px]" data-theme-glyph="dark">{!! icons.Moon(icons.Props{}) !!}</span>
	</button>

	<div data-popover aria-hidden="true" class="w-40">
		{{-- role="menu" with menuitemradio, and not a radiogroup. The two say the
		     same thing about mutual exclusion, and only one of them is what a
		     menu is: the script that opens this popover looks for a menu inside
		     it, and a radiogroup leaves it reporting a missing element and the
		     arrow keys doing nothing. --}}
		<div role="menu" aria-label="Theme">
			<button type="button" role="menuitemradio" aria-checked="true" data-theme-mode="auto">
				<span aria-hidden="true" class="block size-4">{!! icons.Desktop(icons.Props{}) !!}</span>
				<span>Auto</span>
			</button>
			<button type="button" role="menuitemradio" aria-checked="false" data-theme-mode="light">
				<span aria-hidden="true" class="block size-4">{!! icons.Sun(icons.Props{}) !!}</span>
				<span>Light</span>
			</button>
			<button type="button" role="menuitemradio" aria-checked="false" data-theme-mode="dark">
				<span aria-hidden="true" class="block size-4">{!! icons.Moon(icons.Props{}) !!}</span>
				<span>Dark</span>
			</button>
		</div>
	</div>
</div>
