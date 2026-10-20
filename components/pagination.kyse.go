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
// It draws the page a query returned: Paginator takes the
// *pagination.LengthAwarePage, *pagination.Page or *pagination.CursorPage as it
// came back, and the numbers, the window and every address are that page's.
// A page that knows its total draws its numbers, with the first and the last
// always present and an ellipsis where the window left numbers out; a page that
// does not draws the two steps.
//
// It publishes root, list, item, link, previous, next and ellipsis.
type PaginationProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Paginator is the page the query returned. Set, it is what the trail is
	// drawn from -- the current page, the last, the window and every address
	// -- and Page, Pages, URL and Around are not read. A page with one page
	// draws nothing at all: a pager over a single page is chrome with no
	// function.
	Paginator Paginator

	// Page is the page being shown, counting from one.
	//
	// Deprecated: set Paginator to the page the query returned. Page, Pages
	// and URL are that page copied out into three fields.
	Page int
	// Pages is how many there are. Zero or one draws nothing at all.
	//
	// Deprecated: set Paginator.
	Pages int
	// URL is the address of a page, with the number written as "{page}":
	// "/invoices?page={page}".
	//
	// Deprecated: set Paginator, whose addresses keep the query the request
	// arrived with.
	URL string
	// Around is how many numbers to keep either side of the current one. Zero
	// keeps one, which with the first and the last is seven entries.
	//
	// Deprecated: set Paginator, and OnEachSide on the page.
	Around int
	// Label names the landmark. Empty says "Pagination", which is what
	// separates it from the site navigation for somebody listing the
	// landmarks.
	Label string
	// PreviousLabel and NextLabel are the two ends. Empty says "Previous" and
	// "Next".
	PreviousLabel string
	NextLabel     string
	// PageLabel names one numbered control, with the number written as "{n}":
	// "Page {n}". Empty says that in English, which is a floor and not a
	// default worth shipping.
	PageLabel string

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
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginatorEntries(paginator)
	}
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

// paginatorEntries are the numbers a page that knows its total says to draw, in
// the order it lists them, and none for a page that does not.
func paginatorEntries(paginator Paginator) []PaginationEntry {
	page, ok := counted(paginator)
	if !ok || !page.HasPages() {
		return nil
	}
	links := page.Links()
	entries := make([]PaginationEntry, 0, len(links))
	for _, link := range links {
		if link.Page == 0 {
			entries = append(entries, PaginationEntry{Ellipsis: true})
			continue
		}
		entries = append(entries, PaginationEntry{Number: link.Page, Current: link.Active})
	}
	return entries
}

// Shows is whether there is a trail to draw: more than one page.
func (p PaginationProps) Shows() bool {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginator.HasPages()
	}
	return p.Pages > 1
}

// Href is the address of one page.
func (p PaginationProps) Href(number int) string {
	if page, ok := counted(p.Paginator); ok {
		return page.URL(number)
	}
	return strings.ReplaceAll(p.URL, "{page}", strconv.Itoa(number))
}

// HasPrevious and HasNext say whether each end is reachable.
func (p PaginationProps) HasPrevious() bool {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginator.PreviousPageURL() != ""
	}
	return p.Page > 1
}

// HasNext is whether there is a page after this one.
func (p PaginationProps) HasNext() bool {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginator.NextPageURL() != ""
	}
	return p.Page < p.Pages
}

// PreviousHref and NextHref are the two ends' addresses.
func (p PaginationProps) PreviousHref() string {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginator.PreviousPageURL()
	}
	return p.Href(p.Page - 1)
}

// NextHref is the address of the page after this one.
func (p PaginationProps) NextHref() string {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		return paginator.NextPageURL()
	}
	return p.Href(p.Page + 1)
}

// PreviousNumber and NextNumber are the pages the two ends lead to, and zero
// for a page that has no number -- a cursor names a row, not a page.
func (p PaginationProps) PreviousNumber() int {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		if page, ok := paginator.(numberedPage); ok {
			return page.CurrentPage() - 1
		}
		return 0
	}
	return p.Page - 1
}

// NextNumber is the page the forward end leads to, or zero when it has none.
func (p PaginationProps) NextNumber() int {
	if paginator := presentPaginator(p.Paginator); paginator != nil {
		if page, ok := paginator.(numberedPage); ok {
			return page.CurrentPage() + 1
		}
		return 0
	}
	return p.Page + 1
}

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
	sentence := p.PageLabel
	if sentence == "" {
		sentence = "Page {n}"
	}
	return strings.ReplaceAll(sentence, "{n}", strconv.Itoa(number))
}

// PartNames are the parts this component publishes.
func (p PaginationProps) PartNames() []string {
	return []string{"root", "list", "item", "link", "previous", "next", "ellipsis"}
}
@endgo

@if(.Shows())
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
						href="{{ .PreviousHref() }}"
						rel="prev"
						@if(.PreviousNumber() > 0)
							data-page="{{ .PreviousNumber() }}"
						@endif
						@if(.HxTarget != "")
							hx-get="{{ .PreviousHref() }}"
							hx-target="{{ .HxTarget }}"
						@endif
						@if(.HxSwap != "")
							hx-swap="{{ .HxSwap }}"
						@endif
						@if(.HxSwap == "outerHTML" && .HxTarget != "")
							hx-select="{{ .HxTarget }}"
						@endif
						@attributes(.PartAttrs("previous"))
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
							href="{{ .Href(entry.Number) }}"
							data-page="{{ entry.Number }}"
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
							@if(.HxSwap == "outerHTML" && .HxTarget != "")
								hx-select="{{ .HxTarget }}"
							@endif
							@attributes(.PartAttrs("link"))
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
						href="{{ .NextHref() }}"
						rel="next"
						@if(.NextNumber() > 0)
							data-page="{{ .NextNumber() }}"
						@endif
						@if(.HxTarget != "")
							hx-get="{{ .NextHref() }}"
							hx-target="{{ .HxTarget }}"
						@endif
						@if(.HxSwap != "")
							hx-swap="{{ .HxSwap }}"
						@endif
						@if(.HxSwap == "outerHTML" && .HxTarget != "")
							hx-select="{{ .HxTarget }}"
						@endif
						@attributes(.PartAttrs("next"))
					>{{ .NextName() }}</a>
				</li>
			@endif
		</ul>
	</nav>
@endif
