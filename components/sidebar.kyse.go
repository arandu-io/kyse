//go:build kyse

package components

import "html/template"

@go
// SidebarProps is the column of navigation beside a page: groups of links, with
// the one being looked at marked.
//
// It is an aside holding one nav, and the page content is its next sibling.
// That order is not a suggestion -- the stylesheet gives the element after the
// aside the margin that leaves room for the column, so a wrapper between them
// leaves the page underneath the sidebar.
//
// The current entry carries aria-current, which is both what the stylesheet
// paints and what a screen reader announces. Marking it with a class instead
// gives the second half to nobody.
//
// # What collapsing it costs
//
// Collapsed is the state the page is served in. Closing and reopening it after
// that belongs to the button that does it, and the button cannot live here: the
// sidebar is inert while it is closed, so a control inside it is unreachable
// exactly when it is needed. ID is the handle that button uses.
//
// On a narrow screen the column becomes an overlay, and it closes itself when a
// link inside it is followed or the shade beside it is clicked -- so a page
// that never draws a button of its own is still usable there.
type SidebarProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is the handle the control that opens and closes the sidebar reaches it
	// by.
	ID string
	// Groups are the sections of the navigation, top to bottom.
	Groups []SidebarGroup
	// Label names the navigation landmark, which is what a screen reader offers
	// when it lists the ways around a page. Empty means "Sidebar".
	Label string
	// Side is "left" or "right". Empty means left, and it is physical: it names
	// a side of the viewport rather than the start of a line, so it does not
	// turn around in a right-to-left document.
	Side string
	// Header and Footer are the regions above and below the scrolling list --
	// a workspace switcher, an account control. Empty draws neither.
	Header template.HTML
	Footer template.HTML
	// Collapsed serves the page with the column already folded away.
	Collapsed bool
	// MobileOpen serves a narrow screen with the overlay already open. It is
	// separate from Collapsed because the two are different screens: the column
	// is furniture on a wide one and a modal on a narrow one.
	MobileOpen bool
}

// SidebarGroup is one titled run of entries.
type SidebarGroup struct {
	// Label is the heading over the group. Empty draws the entries with no
	// heading and leaves the group unnamed rather than naming it nothing.
	Label string
	// Items are the entries, top to bottom.
	Items []SidebarItem
}

// SidebarItem is one entry: somewhere to go, and whether that is where we are.
//
// It is a link. An entry that acts rather than navigates belongs in Header or
// Footer, where a button is what it looks like -- a list of links with one
// button in it reads as a link that behaves strangely.
type SidebarItem struct {
	// Label is the text.
	Label string
	// URL is where it goes.
	URL string
	// Current marks this as the page being looked at. At most one per sidebar,
	// because the question it answers has one answer.
	Current bool
	// Disabled draws the entry unavailable and says so.
	Disabled bool
}

// Landmark is what the navigation is announced as.
func (p SidebarProps) Landmark() string {
	if p.Label == "" {
		return "Sidebar"
	}
	return p.Label
}

// Edge is the side of the viewport the column is fixed to.
func (p SidebarProps) Edge() string {
	if p.Side == "" {
		return "left"
	}
	return p.Side
}
// PartNames are the parts this component publishes.
func (p SidebarProps) PartNames() []string {
	return []string{"root", "nav", "header", "group", "group-label", "list", "item", "link", "footer"}
}
@endgo

<aside
	data-part="root"
	class="{{ .RootClass("sidebar") }}"
	id="{{ .ID }}"
	data-side="{{ .Edge() }}"
	@attributes(.RootAttrs())
	@if(.Collapsed)
		data-initial-open="false"
	@endif
	@if(.MobileOpen)
		data-initial-mobile-open="true"
	@endif
>
	<nav
		data-part="nav"
		@if(.PartClass("nav") != "")
			class="{{ .PartClass("nav") }}"
		@endif
		aria-label="{{ .Landmark() }}"
		@attributes(.PartAttrs("nav"))
	>
		@if(.Header != "")
			<header
				data-part="header"
				@if(.PartClass("header") != "")
					class="{{ .PartClass("header") }}"
				@endif
				@attributes(.PartAttrs("header"))
			>{!! .Header !!}</header>
		@endif

		<section>
			@foreach(.Groups as group)
				<div
					data-part="group"
					@if(.PartClass("group") != "")
						class="{{ .PartClass("group") }}"
					@endif
					role="group"
					@attributes(.PartAttrs("group"))
					@if(group.Label != "")
						aria-label="{{ group.Label }}"
					@endif
				>
					@if(group.Label != "")
						<h3
							data-part="group-label"
							@if(.PartClass("group-label") != "")
								class="{{ .PartClass("group-label") }}"
							@endif
							@attributes(.PartAttrs("group-label"))
						>{{ group.Label }}</h3>
					@endif
					<ul
						data-part="list"
						@if(.PartClass("list") != "")
							class="{{ .PartClass("list") }}"
						@endif
						@attributes(.PartAttrs("list"))
					>
						@foreach(group.Items as item)
							<li
								data-part="item"
								@if(.PartClass("item") != "")
									class="{{ .PartClass("item") }}"
								@endif
								@attributes(.PartAttrs("item"))
							>
								<a
									data-part="link"
									@if(.PartClass("link") != "")
										class="{{ .PartClass("link") }}"
									@endif
									href="{{ item.URL }}"
									@attributes(.PartAttrs("link"))
									@if(item.Current)
										aria-current="page"
									@endif
									@if(item.Disabled)
										aria-disabled="true"
									@endif
								><span>{{ item.Label }}</span></a>
							</li>
						@endforeach
					</ul>
				</div>
			@endforeach
		</section>

		@if(.Footer != "")
			<footer
				data-part="footer"
				@if(.PartClass("footer") != "")
					class="{{ .PartClass("footer") }}"
				@endif
				@attributes(.PartAttrs("footer"))
			>{!! .Footer !!}</footer>
		@endif
	</nav>
</aside>
