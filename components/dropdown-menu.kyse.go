//go:build kyse

package components

import "strconv"

@go
// DropdownMenuProps is a button and the list of things it opens.
//
// The three parts each carry the name of another: the trigger says which menu
// it controls, the menu says which trigger labels it, and the panel between
// them is what is shown and hidden. Written out by hand that is four ids and
// five attributes, and the one people leave out is aria-expanded -- which is
// the only part that says whether the menu is open at all.
//
// It is a menu rather than a list of links in a box, and the difference is what
// the keyboard does: the arrow keys walk the entries, Home and End reach the
// ends, Escape closes it and hands focus back to the trigger, and a click
// anywhere else closes it. Focus stays on the trigger the whole time and the
// entry the arrows have reached is named by aria-activedescendant, which is why
// every entry is numbered here -- an entry with no id is highlighted and not
// announced, which is the half of the behaviour nobody sees go missing.
//
// A separator and a heading are announced as what they are, so a group of five
// entries is not read as five of nine.
//
// Entries are values rather than markup. A menu is a list of choices, and one
// that could hold anything would have to take a string of HTML -- which is
// where escaping stops being guaranteed by construction.
type DropdownMenuProps struct {
	// ID is what the trigger, the panel, the menu and every entry hang their
	// ids off. Two menus on one page need two.
	ID string
	// Label is the text on the trigger.
	Label string
	// Variant and Size style the trigger. Empty is the default of each.
	Variant string
	Size    string
	// Items are the entries, separators and headings, in the order they read.
	Items []MenuItem
	// Side is which way the panel opens: "top", "bottom", "left", "right", or
	// the logical "inline-start" and "inline-end", which turn around in a
	// right-to-left document. Empty opens downward.
	Side string
	// Align is where the panel sits along that side: "start", "end" or
	// "center". Empty aligns it to the start.
	Align string
}

// MenuItem is one line of the menu.
//
// Separator and Heading are what a line is instead of an entry, and they are
// read in that order: a separator carries no text, a heading is text that
// cannot be chosen, and everything else is a choice.
type MenuItem struct {
	// Label is the text.
	Label string
	// URL makes the entry a link. Empty makes it a button, which is what an
	// entry that acts rather than navigates is.
	URL string
	// Variant is "destructive" for the entry that deletes something. It is the
	// one entry a menu should not draw like the others.
	Variant string
	// Shortcut is the key combination shown down the far side. It is a reminder
	// of a binding the page already has, and drawing it does not create one.
	Shortcut string
	// Heading draws the line as the title of the group under it.
	Heading bool
	// Separator draws the line as a rule between groups.
	Separator bool
	// Disabled draws the entry unavailable, and the arrow keys pass over it.
	Disabled bool

	// The HTMX attributes, written only when they carry something. An empty
	// hx-post is an attribute HTMX acts on: it would post to the current URL.
	HxPost    string
	HxGet     string
	HxTarget  string
	HxConfirm string
}

// TriggerID is the id of the button that opens the menu.
func (p DropdownMenuProps) TriggerID() string { return p.ID + "-trigger" }

// PanelID is the id of the box that is shown and hidden.
func (p DropdownMenuProps) PanelID() string { return p.ID + "-popover" }

// MenuID is the id of the list of entries inside it.
func (p DropdownMenuProps) MenuID() string { return p.ID + "-menu" }

// ItemID is the id of the entry at a position, which is what
// aria-activedescendant points at while the arrow keys are on it.
//
// It is numbered rather than taken from the entry, because the entry has
// nothing unique on it: two menus commonly hold two entries reading "Delete",
// and an id repeated is an announcement naming the wrong line.
func (p DropdownMenuProps) ItemID(at int) string { return p.ID + "-item-" + strconv.Itoa(at) }
@endgo

<div class="dropdown-menu" id="{{ .ID }}">
	<button
		type="button"
		class="btn"
		id="{{ .TriggerID() }}"
		aria-haspopup="menu"
		aria-controls="{{ .MenuID() }}"
		aria-expanded="false"
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
	>{{ .Label }}</button>

	<div
		id="{{ .PanelID() }}"
		data-popover
		aria-hidden="true"
		@if(.Side != "")
			data-side="{{ .Side }}"
		@endif
		@if(.Align != "")
			data-align="{{ .Align }}"
		@endif
	>
		<div role="menu" id="{{ .MenuID() }}" aria-labelledby="{{ .TriggerID() }}">
			@for(at := 0; at < len(.Items); at++)
				@if(.Items[at].Separator)
					<hr role="separator">
				@elseif(.Items[at].Heading)
					<div role="heading" id="{{ .ItemID(at) }}">{{ .Items[at].Label }}</div>
				@elseif(.Items[at].URL != "")
					<a
						role="menuitem"
						id="{{ .ItemID(at) }}"
						href="{{ .Items[at].URL }}"
						@if(.Items[at].Variant != "")
							data-variant="{{ .Items[at].Variant }}"
						@endif
						@if(.Items[at].Disabled)
							aria-disabled="true"
						@endif
					>
						{{ .Items[at].Label }}
						@if(.Items[at].Shortcut != "")
							<kbd>{{ .Items[at].Shortcut }}</kbd>
						@endif
					</a>
				@else
					<button
						type="button"
						role="menuitem"
						id="{{ .ItemID(at) }}"
						@if(.Items[at].Variant != "")
							data-variant="{{ .Items[at].Variant }}"
						@endif
						@if(.Items[at].Disabled)
							aria-disabled="true"
						@endif
						@if(.Items[at].HxPost != "")
							hx-post="{{ .Items[at].HxPost }}"
						@endif
						@if(.Items[at].HxGet != "")
							hx-get="{{ .Items[at].HxGet }}"
						@endif
						@if(.Items[at].HxTarget != "")
							hx-target="{{ .Items[at].HxTarget }}"
						@endif
						@if(.Items[at].HxConfirm != "")
							hx-confirm="{{ .Items[at].HxConfirm }}"
						@endif
					>
						{{ .Items[at].Label }}
						@if(.Items[at].Shortcut != "")
							<kbd>{{ .Items[at].Shortcut }}</kbd>
						@endif
					</button>
				@endif
			@endfor
		</div>
	</div>
</div>
