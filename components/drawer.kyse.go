//go:build kyse

package components

@go
// DrawerProps is a panel that comes in from an edge of the screen, carrying the
// places to go from here.
//
// It is a <dialog>, and that is where its behaviour comes from rather than from
// anything written here. Opened with showModal(), the rest of the document is
// inert, focus cannot leave the panel while it is open, Escape closes it, and
// closing puts focus back on whatever opened it. Four behaviours, none of them
// written, and the four that hand-built panels get wrong.
//
// # What it does not do
//
// It does not open itself, and the way it is opened decides whether any of the
// above is true. A trigger calls showModal() on it. A panel whose open
// attribute is set instead is not modal at all: the page behind it stays
// reachable, focus walks straight out of it, and Escape does nothing.
//
// It does not close when the backdrop is clicked. A modal dialog never does on
// its own, and the click handler that would add it is a script this library
// does not ship.
//
// It slides in and it does not slide out. The stylesheet starts the entry
// transition by itself; the exit needs an attribute set at the moment of
// closing and removed when the transition ends, and that is a script as well.
// The panel disappears rather than leaving.
//
// It carries destinations rather than arbitrary content, for the reason a card
// does: content would have to arrive as a string of HTML, and a string of HTML
// is where escaping stops being guaranteed by construction. A panel that has to
// hold a form is written as <dialog class="drawer"> directly.
type DrawerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what a trigger targets to open this: `onclick="ID.showModal()"`.
	ID string
	// Side is the edge it comes in from: "left", "right", "top" or "bottom".
	// Empty is the bottom, because that is what the stylesheet draws when
	// nothing says otherwise.
	Side string
	// Title is the heading, and the name assistive technology announces when
	// the panel opens.
	Title string
	// Description is the sentence under the title, and may be empty.
	Description string
	// Links are the destinations, drawn in the order they are given. None
	// draws the header and the footer with nothing between them.
	Links []DrawerLink
	// CloseLabel is the button that shuts it. Empty means "Close".
	CloseLabel string
}

// DrawerLink is one destination in a drawer.
type DrawerLink struct {
	// Label is the text.
	Label string
	// Href is where it goes.
	Href string
	// Current marks this as the page already being shown. It is what tells
	// somebody who cannot see the highlight where in the list they are.
	Current bool
}

// Close is the label of the button that shuts the panel.
func (p DrawerProps) Close() string {
	if p.CloseLabel == "" {
		return "Close"
	}
	return p.CloseLabel
}

// DescribedBy names the sentence under the title, and is empty when there is
// none. An aria-describedby pointing at an element that was never drawn
// describes nothing, and reads as a panel whose description failed to load.
func (p DrawerProps) DescribedBy() string {
	if p.Description == "" {
		return ""
	}
	return p.ID + "-description"
}
// PartNames are the parts this component publishes.
func (p DrawerProps) PartNames() []string {
	return []string{"root", "content", "header", "title", "description", "nav", "link", "footer", "close"}
}
@endgo

<dialog
	data-part="root"
	id="{{ .ID }}"
	class="{{ .RootClass("drawer") }}"
	aria-labelledby="{{ .ID }}-title"
	@attributes(.RootAttrs())
	@if(.Side != "")
		data-side="{{ .Side }}"
	@endif
	@if(.DescribedBy() != "")
		aria-describedby="{{ .DescribedBy() }}"
	@endif
>
	<article
		data-part="content"
		@if(.PartClass("content") != "")
			class="{{ .PartClass("content") }}"
		@endif
		@attributes(.PartAttrs("content"))
	>
		<header
			data-part="header"
			@if(.PartClass("header") != "")
				class="{{ .PartClass("header") }}"
			@endif
			@attributes(.PartAttrs("header"))
		>
			<h2
				data-part="title"
				@if(.PartClass("title") != "")
					class="{{ .PartClass("title") }}"
				@endif
				id="{{ .ID }}-title"
				@attributes(.PartAttrs("title"))
			>{{ .Title }}</h2>
			@if(.Description != "")
				<p
					data-part="description"
					@if(.PartClass("description") != "")
						class="{{ .PartClass("description") }}"
					@endif
					id="{{ .ID }}-description"
					@attributes(.PartAttrs("description"))
				>{{ .Description }}</p>
			@endif
		</header>

		@if(len(.Links) > 0)
			<section class="p-4">
				<nav
					data-part="nav"
					class="{{ .PartClass("nav", "grid gap-1") }}"
					@attributes(.PartAttrs("nav"))
				>
					@foreach(.Links as link)
						<a
							data-part="link"
							class="{{ .PartClass("link", "btn justify-start") }}"
							data-variant="ghost"
							href="{{ link.Href }}"
							@attributes(.PartAttrs("link"))
							@if(link.Current)
								aria-current="page"
							@endif
						>{{ link.Label }}</a>
					@endforeach
				</nav>
			</section>
		@endif

		<footer
			data-part="footer"
			@if(.PartClass("footer") != "")
				class="{{ .PartClass("footer") }}"
			@endif
			@attributes(.PartAttrs("footer"))
		>
			<form method="dialog">
				<button
					data-part="close"
					type="submit"
					class="{{ .PartClass("close", "btn") }}"
					data-variant="outline"
					@attributes(.PartAttrs("close"))
				>{{ .Close() }}</button>
			</form>
		</footer>
	</article>
</dialog>
