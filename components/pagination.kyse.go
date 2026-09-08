//go:build kyse

package components

import (
	"strconv"
	"strings"
)

@go
// PaginationProps is the trail of page numbers under a list.
//
// It is a nav landmark holding a list of links, so it is reachable as
// navigation and every entry is a real address -- which means the middle of a
// long list can be linked to, bookmarked, opened in a new tab and gone back
// to. A pager written as buttons that fetch has none of those.
//
// The current page is a link that points at itself carrying aria-current, and
// not a disabled one. Removing it from the tab order would move the keyboard
// past the one entry that says where the reader is.
//
// The window of numbers is computed here, with the first and the last always
// present and an ellipsis where numbers were left out. Ten thousand pages draw
// as nine entries.
//
// It publishes root, list, item, link, previous, next and ellipsis.
type PaginationProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is the page being shown, counting from one.
	Page int
	// Pages is how many there are. Zero or one draws nothing at all: a pager
	// over a single page is chrome with no function.
	Pages int
	// URL is the address of a page, with the number written as "{page}":
	// "/invoices?page={page}". It is a template and not a base so that the
	// number can sit anywhere in it, including in the path.
	URL string
	// Around is how many numbers to keep either side of the current one. Zero
	// keeps one, which with the first and the last is seven entries.
	Around int
	// Label names the landmark. Empty says "Pagination", which is what
	// separates it from the site navigation for somebody listing the
	// landmarks.
	Label string
	// PreviousLabel and NextLabel are the two ends. Empty says "Previous" and
	// "Next".
	PreviousLabel string
	NextLabel     string

	// The HTMX attributes, written on every entry, for a pager that swaps the
	// list instead of loading a document. The links stay real addresses either
	// way -- HTMX takes the click and the URL is still there for everything
	// else.
	HxTarget string
	HxSwap   string
}

// PaginationEntry is one thing drawn in the trail: a number, or the gap where
// numbers were left out.
type PaginationEntry struct {
	// Number is the page, and zero on an ellipsis.
	Number int
	// Ellipsis is whether this is the gap rather than a page.
	Ellipsis bool
	// Current is whether this is the page being shown.
	Current bool
}

// Entries are the numbers to draw, with the gaps already in place.
//
// The first and the last are always there, because a reader who is on page 400
// still has to be able to get to the beginning without knowing the address.
func (p PaginationProps) Entries() []PaginationEntry {
	if p.Pages < 2 {
		return nil
	}
	around := p.Around
	if around <= 0 {
		around = 1
	}
	current := p.Page
	if current < 1 {
		current = 1
	}
	if current > p.Pages {
		current = p.Pages
	}

	keep := map[int]bool{1: true, p.Pages: true}
	for at := current - around; at <= current+around; at++ {
		if at >= 1 && at <= p.Pages {
			keep[at] = true
		}
	}

	entries := make([]PaginationEntry, 0, len(keep)+2)
	previous := 0
	for at := 1; at <= p.Pages; at++ {
		if !keep[at] {
			continue
		}
		if previous != 0 && at-previous > 1 {
			entries = append(entries, PaginationEntry{Ellipsis: true})
		}
		entries = append(entries, PaginationEntry{Number: at, Current: at == current})
		previous = at
	}
	return entries
}

// Href is the address of one page.
func (p PaginationProps) Href(number int) string {
	return strings.ReplaceAll(p.URL, "{page}", strconv.Itoa(number))
}

// HasPrevious and HasNext say whether each end is reachable.
func (p PaginationProps) HasPrevious() bool { return p.Page > 1 }

// HasNext is whether there is a page after this one.
func (p PaginationProps) HasNext() bool { return p.Page < p.Pages }

// PreviousHref and NextHref are the two ends' addresses.
func (p PaginationProps) PreviousHref() string { return p.Href(p.Page - 1) }

// NextHref is the address of the page after this one.
func (p PaginationProps) NextHref() string { return p.Href(p.Page + 1) }

// LandmarkName is what the nav is called.
func (p PaginationProps) LandmarkName() string {
	if p.Label != "" {
		return p.Label
	}
	return "Pagination"
}

// PreviousName and NextName are what the two ends read.
func (p PaginationProps) PreviousName() string {
	if p.PreviousLabel != "" {
		return p.PreviousLabel
	}
	return "Previous"
}

// NextName is what the forward end reads.
func (p PaginationProps) NextName() string {
	if p.NextLabel != "" {
		return p.NextLabel
	}
	return "Next"
}

// PageName is what one number reads to a screen reader. The digit alone is
// read as a digit, which in a row of them says nothing about what it does.
func (p PaginationProps) PageName(number int) string {
	return "Page " + strconv.Itoa(number)
}

// PartNames are the parts this component publishes.
func (p PaginationProps) PartNames() []string {
	return []string{"root", "list", "item", "link", "previous", "next", "ellipsis"}
}
@endgo

@if(.Pages > 1)
	<nav
		data-part="root"
		class="{{ .RootClass("pagination") }}"
		aria-label="{{ .LandmarkName() }}"
		@attributes(.RootAttrs())
	>
		<ul
			data-part="list"
			@if(.PartClass("list") != "")
				class="{{ .PartClass("list") }}"
			@endif
			@attributes(.PartAttrs("list"))
		>
			@if(.HasPrevious())
				<li
					data-part="item"
					@if(.PartClass("item") != "")
						class="{{ .PartClass("item") }}"
					@endif
					@attributes(.PartAttrs("item"))
				>
					<a
						data-part="previous"
						@if(.PartClass("previous") != "")
							class="{{ .PartClass("previous") }}"
						@endif
						@attributes(.PartAttrs("previous"))
						href="{{ .PreviousHref() }}"
						rel="prev"
						@if(.HxTarget != "")
							hx-get="{{ .PreviousHref() }}"
							hx-target="{{ .HxTarget }}"
						@endif
						@if(.HxSwap != "")
							hx-swap="{{ .HxSwap }}"
						@endif
					>{{ .PreviousName() }}</a>
				</li>
			@endif

			@foreach(.Entries() as entry)
				<li
					data-part="item"
					@if(.PartClass("item") != "")
						class="{{ .PartClass("item") }}"
					@endif
					@attributes(.PartAttrs("item"))
				>
					@if(entry.Ellipsis)
						<span
							data-part="ellipsis"
							class="{{ .PartClass("ellipsis", "pagination-ellipsis") }}"
							aria-hidden="true"
							@attributes(.PartAttrs("ellipsis"))
						>&hellip;</span>
					@endif
					@if(!entry.Ellipsis)
						<a
							data-part="link"
							@if(.PartClass("link") != "")
								class="{{ .PartClass("link") }}"
							@endif
							@attributes(.PartAttrs("link"))
							href="{{ .Href(entry.Number) }}"
							aria-label="{{ .PageName(entry.Number) }}"
							@if(entry.Current)
								aria-current="page"
							@endif
							@if(.HxTarget != "")
								hx-get="{{ .Href(entry.Number) }}"
								hx-target="{{ .HxTarget }}"
							@endif
							@if(.HxSwap != "")
								hx-swap="{{ .HxSwap }}"
							@endif
						>{{ entry.Number }}</a>
					@endif
				</li>
			@endforeach

			@if(.HasNext())
				<li
					data-part="item"
					@if(.PartClass("item") != "")
						class="{{ .PartClass("item") }}"
					@endif
					@attributes(.PartAttrs("item"))
				>
					<a
						data-part="next"
						@if(.PartClass("next") != "")
							class="{{ .PartClass("next") }}"
						@endif
						@attributes(.PartAttrs("next"))
						href="{{ .NextHref() }}"
						rel="next"
						@if(.HxTarget != "")
							hx-get="{{ .NextHref() }}"
							hx-target="{{ .HxTarget }}"
						@endif
						@if(.HxSwap != "")
							hx-swap="{{ .HxSwap }}"
						@endif
					>{{ .NextName() }}</a>
				</li>
			@endif
		</ul>
	</nav>
@endif
