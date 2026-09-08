//go:build kyse

package components

import (
	"net/url"
	"strconv"
)

@go
// DataTableProps is a list somebody works in: searched, ordered, paged, and
// acted on in bulk.
//
// # Why it is one component and not a page somebody assembles
//
// The four controls around a table are not independent. Searching has to reset
// the page, or the reader lands on page four of two results. Ordering has to
// keep the search, or typing a name and then sorting throws the name away.
// Paging has to keep both. Every one of those is a bug somebody writes once
// per list, and the reason is that the state lives in four places on the page
// and nothing owns it.
//
// Here it lives in one struct and every link is built from all of it, so the
// four cannot disagree. That is the whole of what this adds over drawing a
// [Table], a [Pagination] and an [ActiveSearch] beside each other.
//
// # Where the work happens
//
// On the server, and the reason is not taste. A table ordered in the browser
// is ordered only within the page that was fetched: page two of a list sorted
// by name holds whatever the server thought page two was, and it reads as rows
// in the wrong order to everyone except whoever wrote it. The same is true of
// a filter -- filtering the rows in hand hides matches that were never sent.
//
// So every control is an address. The endpoint answers with this component
// redrawn, HTMX swaps the region, and the URL says what is on screen -- which
// means a filtered, sorted page three can be linked to, bookmarked and gone
// back to. None of that is true of a table whose state lives in a script.
//
// Without JavaScript the same addresses load the page. The search is a form,
// the headers and the pages are links, and the bulk bar is a form that posts.
//
// It publishes root, toolbar, search, columns, toggle, table, count, bulk and
// pagination, and the parts of the table it draws.
type DataTableProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is the id of the region every control swaps. Two lists on one page
	// need two.
	ID string
	// Label names the region: "Invoices". It is what a screen reader says the
	// list is, and what the search box is called.
	Label string
	// URL is the endpoint, without a query string. Every control appends to
	// it, so there is one address and one handler rather than four.
	URL string

	// Caption is the sentence saying what the list is, drawn under the table
	// and used as its name by assistive technology. Empty draws none, and the
	// region's own Label is then the only thing naming it.
	Caption string
	// Columns are the headers. A column with a Key can be ordered by and
	// switched off; see TableColumn.
	Columns []TableColumn
	// Rows are the lines of the page being shown.
	Rows []TableRow
	// Empty is what stands in when the list has nothing -- which on a searched
	// list means "nothing matched" and not "nothing exists", so the message is
	// the caller's to write.
	Empty EmptyProps

	// Query is what was searched for, and SearchName is the parameter it
	// travels under. Empty SearchName draws no search box.
	Query      string
	SearchName string
	// SearchPlaceholder is the grey text in the empty box.
	SearchPlaceholder string
	// SearchLabel is the box's name. Empty uses Label with "Search" before it.
	SearchLabel string

	// SortKey and SortDir are the ordering in force, and travel as "sort" and
	// "dir".
	SortKey string
	SortDir string

	// Page is the page being shown, counting from one, and Pages is how many
	// there are. It travels as "page". Pages of zero or one draws no pager.
	Page  int
	Pages int

	// SelectName turns on per-row selection and is the field the checkboxes
	// submit under. BulkActions are the controls revealed once something is
	// chosen.
	SelectName  string
	BulkActions []ButtonProps
	// BulkAction is where the bar posts. Empty posts to URL, which is the
	// endpoint that already knows the list.
	BulkAction string
	// Token is the CSRF token for that form.
	Token string

	// ColumnsLabel, SelectedLabel and the pager's two ends are the sentences
	// this draws. Each is documented on the component it belongs to.
	ColumnsLabel  string
	SelectedLabel string
	PreviousLabel string
	NextLabel     string
}

// Address is URL with a query string built from everything except what the
// caller is changing.
//
// Every control goes through here, which is what stops the four from throwing
// each other away: a sort link carries the search, a page link carries both,
// and a new search drops the page rather than keeping a number that no longer
// points anywhere.
func (p DataTableProps) Address(changes map[string]string) string {
	values := url.Values{}
	if p.SearchName != "" && p.Query != "" {
		values.Set(p.SearchName, p.Query)
	}
	if p.SortKey != "" {
		values.Set("sort", p.SortKey)
		if p.SortDir != "" {
			values.Set("dir", p.SortDir)
		}
	}
	if p.Page > 1 {
		values.Set("page", strconv.Itoa(p.Page))
	}
	for key, value := range changes {
		if value == "" {
			values.Del(key)
			continue
		}
		values.Set(key, value)
	}
	if len(values) == 0 {
		return p.URL
	}
	return p.URL + "?" + values.Encode()
}

// SortURL is the template the table builds its header links from. The page is
// dropped: an order that kept it would land somebody on page four of a list
// they have just reordered, where the rows are not the ones they were looking
// at.
func (p DataTableProps) SortURL() string {
	return p.Address(map[string]string{"sort": "{sort}", "dir": "{dir}", "page": ""})
}

// PageURL is the template the pager builds from. It keeps the search and the
// order, because a page of a filtered, sorted list is a page of that list.
func (p DataTableProps) PageURL() string {
	return p.Address(map[string]string{"page": "{page}"})
}

// Target is what every control swaps: this region, by id.
func (p DataTableProps) Target() string { return "#" + p.ID }

// Searchable is whether a search box is drawn.
func (p DataTableProps) Searchable() bool { return p.SearchName != "" }

// SearchName_ is what the search box is called. Empty uses the region's own
// name, because a box labelled only "Search" on a page with two lists names
// neither.
func (p DataTableProps) SearchTitle() string {
	if p.SearchLabel != "" {
		return p.SearchLabel
	}
	if p.Label != "" {
		return "Search " + p.Label
	}
	return "Search"
}

// Hideable is whether any column can be switched off, which is what decides
// whether the columns menu is drawn.
func (p DataTableProps) Hideable() bool {
	for _, column := range p.Columns {
		if column.Hideable && column.Key != "" {
			return true
		}
	}
	return false
}

// Paged is whether there is more than one page to move between.
func (p DataTableProps) Paged() bool { return p.Pages > 1 }

// Grid is the table this draws, built from this component's own fields so the
// two cannot disagree about the order, the selection or the rows.
//
// The count and the columns menu are turned off there and drawn here instead:
// they belong in the toolbar with the search, not stacked above the headers.
func (p DataTableProps) Grid() TableProps {
	columns := make([]TableColumn, 0, len(p.Columns))
	for _, column := range p.Columns {
		// Hideable is cleared so the table draws no menu of its own. The
		// column keeps its Key and its Hidden, which is what the toggle here
		// reaches and what the cells are marked with.
		column.Hideable = false
		columns = append(columns, column)
	}
	return TableProps{
		ComponentProps: ComponentProps{Parts: p.Parts},
		Caption:        p.Caption,
		Columns:        columns,
		Rows:           p.Rows,
		Empty:          p.Empty,
		SortKey:        p.SortKey,
		SortDir:        p.SortDir,
		SortURL:        p.SortURL(),
		SelectName:     p.SelectName,
		SelectedLabel:  p.SelectedLabel,
		BulkActions:    p.BulkActions,
		BulkAction:     p.BulkAction,
		Token:          p.Token,
		HxTarget:       p.Target(),
		HxSwap:         "outerHTML",
	}
}

// Pager is the pagination this draws, or the zero value when there is one page.
func (p DataTableProps) Pager() PaginationProps {
	return PaginationProps{
		ComponentProps: ComponentProps{Parts: p.Parts},
		Page:           p.Page,
		Pages:          p.Pages,
		URL:            p.PageURL(),
		Label:          p.Label,
		PreviousLabel:  p.PreviousLabel,
		NextLabel:      p.NextLabel,
		HxTarget:       p.Target(),
		HxSwap:         "outerHTML",
	}
}

// ColumnsName is what the menu of columns is called.
func (p DataTableProps) ColumnsName() string {
	if p.ColumnsLabel != "" {
		return p.ColumnsLabel
	}
	return "Columns"
}

// SearchTrigger is when the list is refetched as somebody types: after typing
// has paused, and only when what is in the box actually changed.
func (p DataTableProps) SearchTrigger() string {
	return "input changed delay:300ms, search"
}

// SearchAddress is where the box submits: the endpoint with the order kept and
// the page dropped, because a new search has no page four.
func (p DataTableProps) SearchAddress() string {
	return p.Address(map[string]string{p.SearchName: "", "page": ""})
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
//
// It is the table's behaviour, because what needs mounting is the table's:
// the count, the column toggles and the roving cell. Nothing about the
// searching, the ordering or the paging needs a script -- they are addresses.
func (p DataTableProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: TableBehavior, Props: map[string]any{
			"selected": p.SelectedLabel,
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p DataTableProps) PartNames() []string {
	// The table's names and the pager's are published here as well, because
	// this draws both and the parts a caller can reach are the parts on the
	// page -- not the parts this file happens to write itself.
	names := []string{"root", "toolbar", "search", "indicator", "columns", "toggle", "pagination"}
	names = append(names, TableProps{}.PartNames()[1:]...)
	return append(names, PaginationProps{}.PartNames()[1:]...)
}
@endgo

{{-- One region, one id, one target. Every control below swaps this element,
     which is why the search can carry the order and the order can carry the
     search: they are all rebuilt from the same struct on the way back. --}}
<div
	data-part="root"
	class="{{ .RootClass("data-table") }}"
	id="{{ .ID }}"
	role="region"
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(.Searchable() || .Hideable())
		<div
			data-part="toolbar"
			class="{{ .PartClass("toolbar", "data-table-toolbar") }}"
			@attributes(.PartAttrs("toolbar"))
		>
			@if(.Searchable())
				{{-- A form, so Enter submits to the same endpoint and the page
				     comes back with the results in it. HTMX takes the typing
				     before that ever happens. --}}
				<form
					class="data-table-search"
					action="{{ .SearchAddress() }}"
					method="get"
					role="search"
				>
					<label class="sr-only" for="{{ .ID }}-search">{{ .SearchTitle() }}</label>
					<input
						data-part="search"
						class="{{ .PartClass("search", "input") }}"
						type="search"
						id="{{ .ID }}-search"
						name="{{ .SearchName }}"
						value="{{ .Query }}"
						autocomplete="off"
						hx-get="{{ .SearchAddress() }}"
						hx-trigger="{{ .SearchTrigger() }}"
						hx-target="{{ .Target() }}"
						hx-swap="outerHTML"
						hx-sync="this:replace"
						hx-indicator="#{{ .ID }}-indicator"
						@attributes(.PartAttrs("search"))
						@if(.SearchPlaceholder != "")
							placeholder="{{ .SearchPlaceholder }}"
						@endif
					>
					<span
						data-part="indicator"
						class="{{ .PartClass("indicator", "spinner") }}"
						id="{{ .ID }}-indicator"
						aria-hidden="true"
						@attributes(.PartAttrs("indicator"))
					></span>
				</form>
			@endif

			@if(.Hideable())
				{{-- Checkboxes and not menu entries: switching a column off is a
				     two-state choice that stays chosen, which is what a checkbox
				     is and what a menu entry is not. --}}
				<details
					data-part="columns"
					class="{{ .PartClass("columns", "table-columns") }}"
					@attributes(.PartAttrs("columns"))
				>
					<summary>{{ .ColumnsName() }}</summary>
					<div>
						@foreach(.Columns as column)
							@if(column.Hideable && column.Key != "")
								<label
									data-part="toggle"
									@if(.PartClass("toggle") != "")
										class="{{ .PartClass("toggle") }}"
									@endif
									@attributes(.PartAttrs("toggle"))
								>
									<input
										type="checkbox"
										data-column-toggle="{{ column.Key }}"
										@if(!column.Hidden)
											checked
										@endif
									>
									{{ column.Label }}
								</label>
							@endif
						@endforeach
					</div>
				</details>
			@endif
		</div>
	@endif

	{!! Table(.Grid()) !!}

	@if(.Paged())
		<div
			data-part="pagination"
			class="{{ .PartClass("pagination", "data-table-pagination") }}"
			@attributes(.PartAttrs("pagination"))
		>{!! Pagination(.Pager()) !!}</div>
	@endif
</div>
