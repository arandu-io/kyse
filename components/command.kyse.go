//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/kyse/icons"
)

@go
// CommandProps is the palette: everywhere a screen can go, in one list, with a
// box that narrows it.
//
// # Where the filtering happens
//
// In the browser, and nowhere else: a palette holds the places one application
// has, which is a list the page can carry whole, and asking a server between
// two keystrokes to re-sort thirty lines somebody is looking straight at is a
// round trip they would feel.
//
// The list is drawn once and every line stays in the document. Narrowing marks
// the lines that do not match hidden, which is also what empties a group whose
// lines all went, and what puts EmptyText on the screen when none are left.
// Nothing is fetched, so nothing arrives late and nothing arrives out of order.
//
// # What a line is
//
// A link. It is reachable by tab, opens in a new tab on the modifier somebody
// already uses, and goes where it says it goes with the palette's own keyboard
// switched off entirely. A line with no URL is drawn without one, which is an
// anchor the browser will not focus or follow, and is marked disabled so it is
// announced the way it looks.
type CommandProps struct {
	// ID is the id everything inside is hung off. Two palettes on one page
	// need two of them.
	ID string
	// Label is what assistive technology calls the palette as a whole.
	Label string
	// Placeholder is the grey text inside the empty search box.
	Placeholder string
	// EmptyText is the sentence drawn in place of the list when nothing
	// matched. Empty takes the stylesheet's own wording.
	EmptyText string
	// Groups are the lines, in the order they are read, under the headings
	// they belong to. A group whose lines all fail the search draws neither
	// its heading nor a gap.
	Groups []CommandGroup
}

// CommandGroup is one heading and the lines under it.
type CommandGroup struct {
	// Heading is the small line over the group. Empty draws none, which is
	// right for a palette with one group and wrong for a palette with three.
	Heading string
	// Items are the lines.
	Items []CommandItem
}

// CommandItem is one line of the palette.
type CommandItem struct {
	// Label is what the line reads as, and what the search matches first.
	Label string
	// Keywords are the other words that should find this line: what somebody
	// calls it when they do not call it what it is called. They are matched
	// and never shown.
	Keywords string
	// Shortcut is the key combination drawn at the end of the line. It is
	// shown and never bound: what the keys do is the application's, and a
	// component that bound them would take a combination the page already
	// uses.
	Shortcut string
	// URL is where the line goes. Empty draws the line as unavailable.
	URL string
	// Disabled draws the line as unavailable while leaving it in the list, so
	// that what exists but cannot be reached is still visible.
	Disabled bool
}

// Available reports whether this line can be followed.
func (i CommandItem) Available() bool {
	return i.URL != "" && !i.Disabled
}

// Search is what this line is matched against: what it reads as, and the words
// somebody might look for it under.
func (i CommandItem) Search() string {
	if i.Keywords == "" {
		return i.Label
	}
	return i.Label + " " + i.Keywords
}

// InputID names the search box.
func (p CommandProps) InputID() string {
	return p.ID + "-input"
}

// MenuID names the list, which is what the search box says it controls.
func (p CommandProps) MenuID() string {
	return p.ID + "-menu"
}

// HeadingID names a group's heading, so the group is announced by it.
func (p CommandProps) HeadingID(group int) string {
	return p.ID + "-heading-" + strconv.Itoa(group)
}

// ItemID names one line.
//
// It is numbered rather than named after the line, because
// aria-activedescendant holds an id and an id with a space in it is read as
// two.
func (p CommandProps) ItemID(group, item int) string {
	return p.ID + "-item-" + strconv.Itoa(group) + "-" + strconv.Itoa(item)
}
@endgo

<div class="command" id="{{ .ID }}"
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	x-data="{
		q: '',
		active: null,
		lines() { return Array.from(this.$refs.menu.querySelectorAll('[role=menuitem]')).filter(l =&gt; l.getAttribute('aria-hidden') !== 'true' &amp;&amp; l.getAttribute('aria-disabled') !== 'true') },
		hidden(line) {
			const q = this.q.trim().toLowerCase()
			if (!q) return false
			return !line.dataset.search.toLowerCase().includes(q)
		},
		move(step) {
			const all = this.lines()
			if (!all.length) return
			const at = all.findIndex(l =&gt; l.id === this.active)
			const to = at &lt; 0 ? (step &gt; 0 ? 0 : all.length - 1) : (at + step + all.length) % all.length
			this.active = all[to].id
			all[to].scrollIntoView({ block: 'nearest' })
		},
		run() {
			const line = this.active ? document.getElementById(this.active) : null
			if (line) line.click()
		}
	}"
>
	<header>
		{!! icons.MagnifyingGlass(icons.Props{}) !!}
		<input
			type="text"
			role="combobox"
			id="{{ .InputID() }}"
			autocomplete="off"
			autocorrect="off"
			spellcheck="false"
			aria-autocomplete="list"
			aria-expanded="true"
			aria-controls="{{ .MenuID() }}"
			x-model="q"
			x-bind:aria-activedescendant="active"
			x-on:keydown.arrow-down.prevent="move(1)"
			x-on:keydown.arrow-up.prevent="move(-1)"
			x-on:keydown.enter.prevent="run()"
			x-on:keydown.escape="q = ''; active = null"
			@if(.Placeholder != "")
				placeholder="{{ .Placeholder }}"
			@endif
		>
	</header>

	<div role="menu" id="{{ .MenuID() }}" aria-orientation="vertical" x-ref="menu"
		@if(.EmptyText != "")
			data-empty="{{ .EmptyText }}"
		@endif
	>
		@for(g := 0; g < len(.Groups); g++)
			<div role="group"
				@if(.Groups[g].Heading != "")
					aria-labelledby="{{ .HeadingID(g) }}"
				@endif
			>
				@if(.Groups[g].Heading != "")
					<span role="heading" id="{{ .HeadingID(g) }}">{{ .Groups[g].Heading }}</span>
				@endif

				@for(i := 0; i < len(.Groups[g].Items); i++)
					<a
						role="menuitem"
						id="{{ .ItemID(g, i) }}"
						data-search="{{ .Groups[g].Items[i].Search() }}"
						x-bind:aria-hidden="hidden($el)"
						x-bind:class="{ active: active === $el.id }"
						@if(.Groups[g].Items[i].Available())
							href="{{ .Groups[g].Items[i].URL }}"
						@endif
						@if(!.Groups[g].Items[i].Available())
							aria-disabled="true"
						@endif
					>
						<span>{{ .Groups[g].Items[i].Label }}</span>
						@if(.Groups[g].Items[i].Shortcut != "")
							<span data-shortcut>{{ .Groups[g].Items[i].Shortcut }}</span>
						@endif
					</a>
				@endfor
			</div>
		@endfor
	</div>
</div>
