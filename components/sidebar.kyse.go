//go:build kyse

package components

import (
	"html/template"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
//
// # Full and slim
//
// Mode is a second way of folding the column on a wide screen, and the one a
// panel usually wants: "slim" keeps it on the page as a strip of icons instead
// of taking it away. Each entry shows its Icon alone, its label becomes the
// link's name and a tip beside it on hover and on keyboard focus, and a group's
// heading becomes a line between groups. The page beside it keeps the strip's
// width, and the current entry stays marked.
//
// Mode is drawn by the server, so the page is right on the first paint with no
// script and no storage at all. The mode is written on the root as
// data-sidebar-mode, and switching it later is setting that attribute to the
// other value: the stylesheet follows the attribute, so whatever switches it
// -- a behaviour, or a choice remembered on the device -- needs to touch nothing
// else. On a narrow screen the column is the overlay in either mode, with its
// labels.
//
// The look of the slim strip travels with the component, as utilities, so a
// project's stylesheet needs no rule for it. --sidebar-slim-width sets the
// strip's width, and 3rem when the stylesheet does not define it.
//
// It publishes root, nav, header, group, group-label, list, item, link,
// tooltip and footer. The tooltip is drawn only when Mode is set.
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
	// a workspace switcher, an account control. Empty draws neither. In the
	// slim strip they are as wide as the strip and clip what does not fit, so
	// what goes there is drawn to fit it or hidden through Parts.
	Header template.HTML
	Footer template.HTML
	// Collapsed serves the page with the column already folded away.
	Collapsed bool
	// MobileOpen serves a narrow screen with the overlay already open. It is
	// separate from Collapsed because the two are different screens: the column
	// is furniture on a wide one and a modal on a narrow one.
	MobileOpen bool
	// Mode is how the column is drawn on a wide screen: "full", with the
	// labels, or "slim", as a strip of icons. Empty draws the column as it was
	// drawn before modes existed -- no data-sidebar-mode, and nothing that
	// answers to one -- so a page that never asks for a mode is unchanged. Any
	// other value is "full".
	Mode string
}

// SidebarGroup is one titled run of entries.
type SidebarGroup struct {
	// Label is the heading over the group. Empty draws the entries with no
	// heading and leaves the group unnamed rather than naming it nothing. In
	// the slim strip the heading is drawn as a line, and the group keeps the
	// label as its name.
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
	// Icon is drawn before the label, and is all the slim strip shows of the
	// entry. Pass what an icon function returned with no Label of its own --
	// icons.House(icons.Props{}) -- because the entry's label already names
	// the link, and a labelled icon would read it twice. Empty draws no icon,
	// and the slim strip draws the label's first letter in its place.
	Icon template.HTML
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

// Switchable is whether the column was given a mode, and so draws what the
// slim strip needs.
func (p SidebarProps) Switchable() bool { return p.Mode != "" }

// ModeName is the mode written on the root: "slim" or "full". Empty when no
// mode was asked for.
func (p SidebarProps) ModeName() string {
	switch {
	case p.Mode == "":
		return ""
	case p.Mode == "slim":
		return "slim"
	default:
		return "full"
	}
}

// Slim is whether the page is served with the strip.
func (p SidebarProps) Slim() bool { return p.ModeName() == "slim" }

// The classes below are written out whole, each in one string, because the
// stylesheet is compiled from the text of this file: a class assembled from
// pieces is a class no rule is emitted for.
//
// Every one that changes how the column looks waits on
// data-sidebar-mode="slim" and on a wide screen, so the full column and the
// overlay on a narrow screen look exactly as the stylesheet draws them without
// a mode. The transitions are the exception: they are shortened in every mode,
// because the column that moves between the two is the same column.

// RootModeClass is what the root adds when the column has a mode: the margin
// the strip leaves the page beside it, on the side it sits, unless the column
// is folded away; and the page's transition, shortened and taken away under
// reduced motion.
func (p SidebarProps) RootModeClass() string {
	if !p.Switchable() {
		return ""
	}
	if p.Edge() == "right" {
		return "[&+*]:duration-200 motion-reduce:[&+*]:transition-none md:data-[sidebar-mode=slim]:[&:not([aria-hidden=true])+*]:mr-[var(--sidebar-slim-width,3rem)]"
	}
	return "[&+*]:duration-200 motion-reduce:[&+*]:transition-none md:data-[sidebar-mode=slim]:[&:not([aria-hidden=true])+*]:ml-[var(--sidebar-slim-width,3rem)]"
}

// NavModeClass narrows the column to the strip, and moves its width with the
// same transition that slides it.
func (p SidebarProps) NavModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "transition-[transform,translate,width] duration-200 motion-reduce:transition-none md:in-data-[sidebar-mode=slim]:w-[var(--sidebar-slim-width,3rem)]"
}

// EdgeModeClass keeps what the header and the footer hold inside the strip.
func (p SidebarProps) EdgeModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "md:in-data-[sidebar-mode=slim]:overflow-hidden"
}

// GroupLabelModeClass draws a group's heading as a line in the strip. The
// text stays in the heading, clipped, so it is still read out.
func (p SidebarProps) GroupLabelModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "md:in-data-[sidebar-mode=slim]:h-0 md:in-data-[sidebar-mode=slim]:overflow-hidden md:in-data-[sidebar-mode=slim]:rounded-none md:in-data-[sidebar-mode=slim]:border-t md:in-data-[sidebar-mode=slim]:border-sidebar-border md:in-data-[sidebar-mode=slim]:p-0 md:in-data-[sidebar-mode=slim]:mb-2"
}

// LinkModeClass marks the link as the element the tip beside it watches.
func (p SidebarProps) LinkModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "peer"
}

// LabelModeClass hides an entry's text from the eye in the strip and keeps it
// as the link's name.
func (p SidebarProps) LabelModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "md:in-data-[sidebar-mode=slim]:sr-only"
}

// InitialClass draws the stand-in for a missing icon, in the strip only.
func (p SidebarProps) InitialClass() string {
	return "hidden size-4 shrink-0 items-center justify-center text-xs font-semibold md:in-data-[sidebar-mode=slim]:flex"
}

// SectionModeClass lets the list overflow the strip instead of scrolling it.
//
// A scrolling box clips whatever it holds, and the tip beside an entry sits
// outside the strip, so a list that scrolled would cut every tip off at the
// strip's edge. A fixed tip escapes the clip and does not follow the list when
// it scrolls -- measured, it stays where the entry was before the scroll. So
// the strip does not scroll, which is what the icon sidebar it is modelled on
// does: a strip holds the handful of places a panel goes, and a list long
// enough to scroll is one that wants its labels.
func (p SidebarProps) SectionModeClass() string {
	if !p.Switchable() {
		return ""
	}
	return "md:in-data-[sidebar-mode=slim]:overflow-visible"
}

// TooltipClass is the tip beside an entry of the strip: the tooltip's own
// look, beside the entry on the side away from the edge, shown while the link
// is hovered or has keyboard focus, and only in the strip.
func (p SidebarProps) TooltipClass() string {
	if p.Edge() == "right" {
		return "tooltip-content hidden md:in-data-[sidebar-mode=slim]:block top-1/2 bottom-auto left-auto right-[calc(100%+0.75rem)] -translate-y-1/2 peer-hover:opacity-100 peer-focus-visible:opacity-100"
	}
	return "tooltip-content hidden md:in-data-[sidebar-mode=slim]:block top-1/2 bottom-auto left-[calc(100%+0.75rem)] -translate-y-1/2 peer-hover:opacity-100 peer-focus-visible:opacity-100"
}

// Plain is whether an entry is drawn exactly as it was before icons and
// modes existed: no mode on the column and no icon on the entry.
func (p SidebarProps) Plain(i SidebarItem) bool { return !p.Switchable() && i.Icon == "" }

// Initial is the letter the strip draws for an entry with no icon: the first
// letter of its label, upper-cased.
func (i SidebarItem) Initial() string {
	label := strings.TrimSpace(i.Label)
	if label == "" {
		return ""
	}
	r, _ := utf8.DecodeRuneInString(label)
	return string(unicode.ToUpper(r))
}

// PartNames are the parts this component publishes.
func (p SidebarProps) PartNames() []string {
	return []string{"root", "nav", "header", "group", "group-label", "list", "item", "link", "tooltip", "footer"}
}
@endgo

{{-- The tip is aria-hidden because it says nothing new: the link is already
     named by its label, which in the strip is the aria-label as well as the
     hidden text. A tip that was also a description would read the name twice.

     The slim classes are written only when the column has a mode, so a page
     that never asked for one is sent the markup it was sent before. --}}
<aside
	data-part="root"
	class="{{ .RootClass("sidebar", .RootModeClass()) }}"
	id="{{ .ID }}"
	data-side="{{ .Edge() }}"
	@if(.Switchable())
		data-sidebar-mode="{{ .ModeName() }}"
	@endif
	@if(.Collapsed)
		data-initial-open="false"
	@endif
	@if(.MobileOpen)
		data-initial-mobile-open="true"
	@endif
	@attributes(.RootAttrs())
>
	<nav
		data-part="nav"
		@if(.PartClass("nav", .NavModeClass()) != "")
			class="{{ .PartClass("nav", .NavModeClass()) }}"
		@endif
		aria-label="{{ .Landmark() }}"
		@attributes(.PartAttrs("nav"))
	>
		@if(.Header != "")
			<header
				data-part="header"
				@if(.PartClass("header", .EdgeModeClass()) != "")
					class="{{ .PartClass("header", .EdgeModeClass()) }}"
				@endif
				@attributes(.PartAttrs("header"))
			>{!! .Header !!}</header>
		@endif

		@if(.Switchable())
		<section class="{{ .SectionModeClass() }}">
		@else
		<section>
		@endif
			@foreach(.Groups as group)
				<div
					data-part="group"
					@if(.PartClass("group") != "")
						class="{{ .PartClass("group") }}"
					@endif
					role="group"
					@if(group.Label != "")
						aria-label="{{ group.Label }}"
					@endif
					@attributes(.PartAttrs("group"))
				>
					@if(group.Label != "")
						<h3
							data-part="group-label"
							@if(.PartClass("group-label", .GroupLabelModeClass()) != "")
								class="{{ .PartClass("group-label", .GroupLabelModeClass()) }}"
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
								@if(.Plain(item))
								<a
									data-part="link"
									@if(.PartClass("link") != "")
										class="{{ .PartClass("link") }}"
									@endif
									href="{{ item.URL }}"
									@if(item.Current)
										aria-current="page"
									@endif
									@if(item.Disabled)
										aria-disabled="true"
									@endif
									@attributes(.PartAttrs("link"))
								><span>{{ item.Label }}</span></a>
								@else
								<a
									data-part="link"
									@if(.PartClass("link", .LinkModeClass()) != "")
										class="{{ .PartClass("link", .LinkModeClass()) }}"
									@endif
									href="{{ item.URL }}"
									@if(.Slim())
										aria-label="{{ item.Label }}"
									@endif
									@if(item.Current)
										aria-current="page"
									@endif
									@if(item.Disabled)
										aria-disabled="true"
									@endif
									@attributes(.PartAttrs("link"))
								>
									@if(item.Icon != "")
										{!! item.Icon !!}
									@endif
									@if(.Switchable() && item.Icon == "")
										<span class="{{ .InitialClass() }}" aria-hidden="true">{{ item.Initial() }}</span>
									@endif
									<span
										@if(.LabelModeClass() != "")
											class="{{ .LabelModeClass() }}"
										@endif
									>{{ item.Label }}</span>
								</a>
								@endif
								@if(.Switchable())
									<span
										data-part="tooltip"
										class="{{ .PartClass("tooltip", .TooltipClass()) }}"
										aria-hidden="true"
										@attributes(.PartAttrs("tooltip"))
									>{{ item.Label }}</span>
								@endif
							</li>
						@endforeach
					</ul>
				</div>
			@endforeach
		</section>

		@if(.Footer != "")
			<footer
				data-part="footer"
				@if(.PartClass("footer", .EdgeModeClass()) != "")
					class="{{ .PartClass("footer", .EdgeModeClass()) }}"
				@endif
				@attributes(.PartAttrs("footer"))
			>{!! .Footer !!}</footer>
		@endif
	</nav>
</aside>
