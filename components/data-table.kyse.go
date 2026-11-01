//go:build kyse

package components

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/icons"
)

@go
// DataTableProps is a list somebody works in: searched, filtered by facet,
// ordered, paged, shown with the columns they want, and acted on in bulk.
//
// It is drawn as shadcn/ui draws a data table -- a toolbar over a bordered
// table over a footer -- and works the way a server-rendered list has to: every
// control is an address, the state is in the query string, and the answer is
// this component drawn again.
//
// # Where the state lives
//
// In [DataTableState], one value that says what is searched, how the list is
// ordered, which facets are chosen, which page and how many rows, and which
// columns are hidden. Every address the table draws is the state's own answer
// to "what is the query after this step?", so a link can never carry a value
// twice and a step can never forget the others. The search form carries the
// same state as fields, so a search keeps the order, the size and the hidden
// columns, and a new search sends the reader back to the first page.
//
// A list read from a database passes the [pagination.ListingResult] its query
// returned as Result and nothing else about the state: the state, the page with
// its total, the count of every facet value and the page sizes come out of it.
//
// # Where the work happens
//
// On the server. A table ordered in the browser is ordered only within the page
// that was fetched, and a filter over the rows in hand hides matches that were
// never sent. So the endpoint answers with this component, htmx swaps the
// region, and the address bar says what is on screen: a filtered, sorted page
// three can be linked to, bookmarked, reloaded and gone back to.
//
// One form owns the search and the facets, and carries the order, the page size
// and the hidden columns as fields. It makes one request however the change
// came: a word typed and paused on, Enter, a facet ticked. Every header, page
// step, page size and column switch is a link to the address of its own step.
//
// Without JavaScript the same addresses load the page.
//
// # Optional parts
//
// Searchable draws the search box, Filters the faceted filters with their
// counts, HeaderMenu the Asc, Desc and Hide menu on a column header, PageSizes
// the rows-per-page chooser, FullPager the first, previous, next and last
// buttons with "Page X of Y", and a row's Actions the menu at the end of the
// row. A column that is Hideable is listed in the columns menu.
//
// # Controls on the search's line
//
// Controls draws a text, a choice or a period beside the search, in the same
// form, so one request carries all of them and a new order, a page or a size
// keeps them. See [DataTableControl]. An application that wrote its own toolbar
// of filters and spliced them into the markup draws them with this instead.
//
// # A list that holds all its rows
//
// Complete says the rows are the whole list. The browser then searches, orders,
// filters and pages them itself and the server is not asked. Every control is
// still a link or a form, so the page works with no script.
//
// The browser searches a cell by its own words, without regard to case or to the
// marks on a letter. What a cell holds that is not the row's words -- a [Dialog]
// drawn in it, the text of an action, anything with data-search="false" -- is
// left out of the search and of the order.
//
// # On a narrow screen
//
// Where the region is narrower than 600px -- the same breakpoint Ayra's DataTable
// draws cards under -- each row is a card and nothing scrolls sideways. The
// column marked Primary on [TableColumn] is the card's title and with none marked
// it is the first column the card shows; the others are a line each, the name of
// the column and then the value; a column marked CardHidden is left out; a row's
// menu is along the foot of its card. The head is left as a bar of the controls
// that order the list, so ordering is still there, and the search, the facets and
// the pager are as they were.
//
// It is the same markup as the table, drawn under a container query on the
// region, so the width that counts is the one the table is given and not the
// window's, and there is no script, no inline style and nothing to switch on:
// a table at 600px or more is the table it was. The name in front of a value is
// the cell's data-label, drawn by CSS, so the text of a cell is still only its
// value and a screen reader reads the name with it.
//
// It publishes root, toolbar, query, search, reset, facet, facet-chosen, facet-search, filter,
// control, control-input, columns, toggle, bulk, wrapper, table, caption, head, header-cell, sort,
// header-menu, row, cell, select-all, select, row-menu, footer, count, sizes,
// pagination, page-of, first, previous, next and last.
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
	// it, so there is one address and one handler rather than a dozen.
	//
	// A query string or fragment on it is dropped, never carried: every address
	// and every hidden field this draws is built from the state alone. URL is
	// often the request's own address, and a query carried from it would put
	// whatever parameters somebody wrote into a link on every control of the
	// page.
	URL string

	// State is the listing's state: the search, the order, the facets, the page
	// size and the hidden columns, and the address of every step from them.
	// Nil is a list in its first state, with the conventional parameter names.
	State DataTableState
	// Result is what a listing returned: [pagination.ListingResult], the state
	// that was read, the page with its total and the count of every facet value.
	// Set, it fills State (through [DataTableStateFrom]), Paginator, the Counts
	// and Param of every filter and the page sizes of the listing, each of which
	// a field written beside it still overrides. A service hands the result on
	// and the table is drawn from it.
	Result *pagination.ListingResult
	// Paginator is the page the query returned: the *pagination.LengthAwarePage,
	// *pagination.Page or *pagination.CursorPage that came back beside the rows.
	// It is where the table reads the page being shown, how many there are and
	// how many rows the whole list has. A page that knows no total draws the
	// two steps and neither "Page X of Y" nor a total, because neither can be
	// true. A complete list has no use for one: there the rows are the list.
	Paginator Paginator

	// Caption is the sentence saying what the list is, drawn under the table
	// and used as its name by assistive technology. Empty draws none, and the
	// region's own Label is then the only thing naming it.
	Caption string
	// Columns are the headers. A column with a Key can be ordered by and
	// switched off; see TableColumn. A column's Hidden is only where a list
	// without a State starts: with one, State.Hidden decides.
	Columns []TableColumn
	// Rows are the lines of the page being shown, or the whole list when
	// Complete says so.
	Rows []TableRow
	// Complete says the rows above are the whole list and not a page of it.
	//
	// It is a fact about the data and not a preference. When the page holds
	// every row, the browser can search, order, filter and page correctly
	// because it has everything -- and it does, so a click costs nothing and
	// the server is not asked. When the page holds a window onto something
	// longer, the browser cannot: ordering the rows in hand orders only those,
	// and filtering them hides matches that were never sent. Then every control
	// is an address and the server answers.
	//
	// Saying true about a page of a longer list is the bug this field is named
	// to prevent: it reads as rows in the wrong order to everyone except
	// whoever wrote it, and nothing crashes.
	//
	// A complete list still keeps its addresses. The headers are still links
	// and the search is still a form, so the page works with no script -- the
	// behaviour takes the click before the link is followed, and the links ask
	// htmx to leave them alone so a boosted page does not follow them either.
	Complete bool
	// PageSize is how many rows a page holds when State says nothing: the size
	// of a complete list's window, and the size the chooser starts on. Zero
	// shows a complete list whole.
	PageSize int
	// PageSizes are the counts the "Rows per page" chooser offers. Empty draws
	// no chooser, which is right for a list whose length nobody argues with.
	PageSizes []int
	// FullPager draws the first, previous, next and last buttons with "Page X
	// of Y" beside them, instead of the two text buttons.
	FullPager bool
	// HeaderMenu draws a sortable header as a button that opens a menu with Asc,
	// Desc and Hide, instead of a button that cycles the order.
	HeaderMenu bool

	// Lines are the sentences this draws, so none of them is written in a
	// script and all of them come from the application's own catalogue.
	Lines DataTableLines
	// Empty is what stands in when the list has nothing at all. A list that
	// has rows and whose search or facets matched none of them says "No
	// results." in the table instead: nothing matched is not nothing exists.
	Empty EmptyProps

	// Searchable draws the search box. The text travels under the state's
	// search parameter.
	Searchable bool
	// SearchPlaceholder is the grey text in the empty box. Empty says "Filter
	// {label}...".
	SearchPlaceholder string
	// SearchLabel is the box's name. Empty says kyse.data_table.search from
	// PageData with Label in place of "{label}", or Label with "Search" before
	// it.
	SearchLabel string
	// Filters are the faceted filters, drawn beside the search. Their chosen
	// values are State.Chosen(Key).
	Filters []DataTableFilter
	// Controls are a text, a choice and a period drawn on the search's line,
	// inside the query form, after the facets. See [DataTableControl].
	Controls []DataTableControl

	// Navigable turns the table into a grid: one tab stop, and the arrow keys
	// moving from cell to cell. See TableProps.Navigable.
	Navigable bool

	// SelectName turns on per-row selection and is the field the checkboxes
	// submit under. The selection survives a search, a new order, a facet and a
	// page: the browser keeps the keys ticked, and the bulk bar posts all of
	// them.
	SelectName string
	// BulkActions are the buttons of the bar that appears once something is
	// ticked. A button posts the ticked keys with its own Name and Value, so
	// the server can tell which was pressed.
	BulkActions []ButtonProps
	// BulkAction is where the bar posts. Empty draws no bar: a form that posts
	// to the list's own endpoint would be answered by a handler that does not
	// know the verb.
	BulkAction string
	// BulkMethod is how the bar submits. Empty is post.
	BulkMethod string
	// PageData is the screen's own view.Page. The bulk form posts the CSRF
	// token it carries and every sentence left empty is read from it, so a list
	// drawn from a view writes PageData: .Page and never copies the token or
	// the words by hand. It is not called Page because Page is the number of
	// the page being shown.
	PageData Page
	// Token is the CSRF token, for a list with no PageData to ask. Set, it wins
	// over the page's.
	Token string

	// ColumnsLabel and SelectedLabel are the two sentences that belong to the
	// table rather than to the list. Everything else this says is in Lines --
	// two fields for one control is exactly the defect that puts a table into
	// two languages. Empty says kyse.table.columns and kyse.table.selected
	// from PageData, or the English.
	ColumnsLabel  string
	SelectedLabel string
}

// DataTableFilter is one faceted filter, drawn beside the search.
//
// It is typed rather than a slot of markup, because a slot would be a place to
// put anything, and this component's whole claim is that every control's state
// survives every other control -- which it can only promise about state it
// knows the shape of.
type DataTableFilter struct {
	// Key is the parameter the chosen values travel under, and the name State
	// knows the facet by.
	Key string
	// Label names the control.
	Label string
	// Options are what can be chosen.
	Options []SelectOption
	// Param is the query parameter the chosen values travel in, which is Key
	// unless the facet of a listing names another. Empty means Key.
	Param string
	// Counts are how many rows each option matches, by option value, the way the
	// query counted them: ignoring this facet's own choice and respecting the
	// others. Nil draws no counts.
	Counts map[string]int
}

// The kinds of control a data table draws on the search's line.
const (
	// DataTableText is a short text, asked for after a pause like the search.
	DataTableText = "text"
	// DataTableSelect is a choice among options, asked for when it changes.
	DataTableSelect = "select"
	// DataTablePeriod is two days chosen in one field, asked for when the second
	// is pressed, a shortcut is or Clear is.
	DataTablePeriod = "period"
)

// DataTableControl is a control drawn on the search's line, in the one form
// that owns the search and the facets: a text, a choice or a period.
//
// It is a facet that is not a list of checkboxes. Its value is in the state
// under Key -- and, for a period, its end under ToKey -- so it travels in every
// address the table draws (a new order, a page, a size keep it) and Reset drops
// it with the search. The query that reads the state declares those keys as the
// facets they are: a key the listing did not declare is dropped from the state,
// as in any request.
//
// On a table the server answers, one request carries the search, the facets and
// every control. On a complete table the browser narrows the rows with them
// against each row's Facets: a choice must equal the row's value for Key, a text
// must be found in it without regard to case or accents, and a period must
// hold it, a day as "2026-10-27" or a timestamp that begins with one.
type DataTableControl struct {
	// Key is the query parameter the value travels in, and for a period the
	// start of it.
	Key string
	// ToKey is the query parameter of the end of a period.
	ToKey string
	// Kind is DataTableText, DataTableSelect or DataTablePeriod. Empty is a text.
	Kind string
	// Label names the control, drawn beside it. A period names its field with it.
	Label string
	// Placeholder is the grey text of a text, and the first option of a choice,
	// the one that narrows nothing ("All").
	Placeholder string
	// Options are what a choice offers.
	Options []SelectOption
	// Value is what the control holds when the state says nothing about Key, and
	// ToValue the end of a period.
	Value   string
	ToValue string
	// Period is the picker of a period, for the words and shortcuts only: month
	// and weekday names, date format, labels, Presets. The names, id, label and
	// the two days are the control's.
	Period DateRangePickerProps
}

// Type is the kind of the control, a text when none is written.
func (c DataTableControl) Type() string {
	switch c.Kind {
	case DataTableSelect, DataTablePeriod:
		return c.Kind
	}
	return DataTableText
}

// Name is the name of the filter's boxes: the parameter its values travel in.
func (f DataTableFilter) Name() string {
	if f.Param != "" {
		return f.Param
	}
	return f.Key
}

// DataTableLines are the sentences a data table says.
//
// They are a struct and not a scatter of Label fields because they belong
// together: a table says all of them or none, and a project that translated
// four of the twenty-nine has a table that speaks two languages. Every one
// carries its placeholders, so a language that puts the number last can.
//
// Empty asks the page: a data table given PageData reads each sentence left
// empty from the application's catalogue, under the key the package doc lists
// for it, in the language the page is served in -- the only place that knows
// it. What the catalogue does not carry says the English underneath, which is
// a floor and not a default worth shipping.
type DataTableLines struct {
	// Filter is the placeholder of the search box: "Filter {label}...".
	Filter string
	// NoResult is what the body says when a search or a facet matched nothing:
	// "No results.". It is not the same sentence as an empty list: the list has
	// rows, and none of them is the answer.
	NoResult string
	// Rows is the footer of a list with no selection, in one form for every
	// count: "{total} row(s).". It is what a project wrote before the footer
	// had a singular and a plural, and it still wins over the catalogue: a
	// sentence given here is the sentence drawn for a total of one and for every
	// other, unless RowsOne or RowsMany says otherwise for its own count.
	Rows string
	// RowsOne is the footer when the list has exactly one row: "{total} row.".
	// RowsMany is the footer for every other total, zero included: "{total}
	// rows.". A language with no singular for the count -- or one that wants the
	// two to differ in more than a suffix -- writes the whole sentence in each.
	RowsOne  string
	RowsMany string
	// SelectedCount is the footer of a selectable list whose total is unknown:
	// "{n} row(s) selected.".
	SelectedCount string
	// RowsPerPage labels the size chooser: "Rows per page".
	RowsPerPage string
	// PageOf is where the reader is, beside the pager: "Page {n} of {total}".
	PageOf string
	// First, Previous, Next and Last name the four steps of the full pager for
	// a screen reader. They are names and not text, because the controls are
	// arrows.
	First    string
	Previous string
	Next     string
	Last     string
	// Back and Forward are the words on the two buttons of the plain pager:
	// "Previous" and "Next".
	Back    string
	Forward string
	// SelectAll names the header checkbox, SelectRow a row's checkbox when the
	// row has no Label of its own, and Bulk the bar of actions.
	SelectAll string
	SelectRow string
	Bulk      string
	// Apply names the submit control that makes the facets work when
	// JavaScript is unavailable. Empty uses "Apply filters".
	Apply string
	// Ascending, Descending and Unsorted are what a sortable header is called
	// in each of its three states: "{column}, sorted ascending". A header that
	// only draws an arrow is a header a screen reader reads as a link.
	Ascending  string
	Descending string
	Unsorted   string
	// Reset clears the search and every facet. ClearFilters clears one facet.
	Reset        string
	ClearFilters string
	// NoFacetResults is what a facet's own search says when nothing matched:
	// "No results found.". FacetSelected is the badge that stands for many
	// chosen values: "{n} selected".
	NoFacetResults string
	FacetSelected  string
	// OpenMenu names the button of a row's menu for a screen reader, Actions
	// the empty header over it, and ToggleColumns the columns menu.
	OpenMenu      string
	Actions       string
	ToggleColumns string
	// Asc, Desc and Hide are the three entries of a header menu.
	Asc  string
	Desc string
	Hide string
}

// Filled is the lines with every sentence the caller left empty read from the
// catalogue of the page, and what neither said filled with the English.
//
// Filling them here rather than in the script is the whole point: the script
// composes no sentence, so there is one place a translation lands and it is a
// file the application already has.
func (lines DataTableLines) Filled(page Page) DataTableLines {
	for _, line := range []struct {
		field   *string
		key     string
		english string
	}{
		{&lines.Filter, "data_table.filter", "Filter {label}..."},
		{&lines.NoResult, "data_table.no_result", "No results."},
		{&lines.SelectedCount, "data_table.selected_count", "{n} row(s) selected."},
		{&lines.RowsPerPage, "data_table.rows_per_page", "Rows per page"},
		{&lines.PageOf, "data_table.page_of", "Page {n} of {total}"},
		{&lines.First, "data_table.first", "Go to first page"},
		{&lines.Previous, "data_table.previous", "Go to previous page"},
		{&lines.Next, "data_table.next", "Go to next page"},
		{&lines.Last, "data_table.last", "Go to last page"},
		{&lines.Back, "data_table.back", "Previous"},
		{&lines.Forward, "data_table.forward", "Next"},
		{&lines.SelectAll, "table.select_all", "Select all rows"},
		{&lines.SelectRow, "data_table.select_row", "Select row"},
		{&lines.Bulk, "table.bulk", "Bulk actions"},
		{&lines.Apply, "data_table.apply", "Apply filters"},
		{&lines.Ascending, "table.ascending", "{column}, sorted ascending"},
		{&lines.Descending, "table.descending", "{column}, sorted descending"},
		{&lines.Unsorted, "table.unsorted", "{column}, not sorted"},
		{&lines.Reset, "data_table.reset", "Reset"},
		{&lines.ClearFilters, "data_table.clear_filters", "Clear filters"},
		{&lines.NoFacetResults, "data_table.no_facet_results", "No results found."},
		{&lines.FacetSelected, "data_table.facet_selected", "{n} selected"},
		{&lines.OpenMenu, "data_table.open_menu", "Open menu"},
		{&lines.Actions, "data_table.actions", "Actions"},
		{&lines.ToggleColumns, "data_table.toggle_columns", "Toggle columns"},
		{&lines.Asc, "data_table.asc", "Asc"},
		{&lines.Desc, "data_table.desc", "Desc"},
		{&lines.Hide, "data_table.hide", "Hide"},
	} {
		if *line.field == "" {
			*line.field = say(page, line.key, line.english)
		}
	}
	lines.fillRows(page)
	return lines
}

// fillRows settles the three forms of the footer. The order of who wins is the
// caller's own fields, then the catalogue, then the English:
//
//   - RowsOne or RowsMany, written on the props, is the sentence for its count.
//   - Rows, written on the props, is the sentence for whichever count has none
//     of its own, so a project that passed it before the split draws exactly what
//     it drew.
//   - Then the catalogue: kyse.data_table.rows_one and rows_many, and for a
//     count with neither, the older kyse.data_table.rows.
//   - Then "{total} row." and "{total} rows.".
//
// Rows is left holding the plural, because that is the form the browser writes
// into the footer of a list it filters (see Sentences), and a script that
// predates the singular reads only that.
func (l *DataTableLines) fillRows(page Page) {
	one, many := l.RowsOne, l.RowsMany
	if l.Rows != "" {
		if one == "" {
			one = l.Rows
		}
		if many == "" {
			many = l.Rows
		}
	}
	if one == "" {
		one = sayFirst(page, "{total} row.", "data_table.rows_one", "data_table.rows")
	}
	if many == "" {
		many = sayFirst(page, "{total} rows.", "data_table.rows_many", "data_table.rows")
	}
	l.RowsOne, l.RowsMany, l.Rows = one, many, many
}

// RowsFor is the footer for a list of total rows: RowsOne for exactly one,
// RowsMany for every other count. It reads the lines as they are, so a caller
// that wants the catalogue's words calls it on [DataTableLines.Filled].
func (l DataTableLines) RowsFor(total int) string {
	sentence := l.RowsMany
	if total == 1 {
		sentence = l.RowsOne
	}
	if sentence == "" {
		sentence = l.Rows
	}
	return strings.ReplaceAll(sentence, "{total}", strconv.Itoa(total))
}

// Sentences are the sentences the browser rewrites or reads, by name, for the
// behaviour's props.
func (l DataTableLines) Sentences() map[string]any {
	return map[string]any{
		"pageOf":        l.PageOf,
		"selectedCount": l.SelectedCount,
		"rows":          l.Rows,
		"rowsOne":       l.RowsOne,
		"ascending":     l.Ascending,
		"descending":    l.Descending,
		"unsorted":      l.Unsorted,
		"facetSelected": l.FacetSelected,
	}
}

// lines is Lines filled for this table's page.
func (p DataTableProps) lines() DataTableLines { return p.Lines.Filled(p.PageData) }

// DataTableQueryField is one hidden field of the query form: a part of the
// state no visible control of the form owns.
type DataTableQueryField struct {
	Name  string
	Value string
}

// state is State, else the state the Result was read with, else the first state
// of an unfiltered list.
func (p DataTableProps) state() DataTableState {
	switch {
	case p.State != nil:
		return p.State
	case p.Result != nil:
		return DataTableStateFrom(p.Result.State)
	}
	return DataTableValues{}
}

// stateless is whether nothing says what state the list is in, which is when a
// column's own Hidden is where it starts.
func (p DataTableProps) stateless() bool { return p.State == nil && p.Result == nil }

// paginator is Paginator, else the page of the Result.
func (p DataTableProps) paginator() Paginator {
	if p.Paginator == nil && p.Result != nil {
		return p.Result.Page
	}
	return p.Paginator
}

// Facets are the filters as they are drawn: each with the parameter and the
// counts the Result knows for it when the filter says none, and the values the
// facet lists when the filter names no options.
func (p DataTableProps) Facets() []DataTableFilter {
	if p.Result == nil || p.Result.State.Listing() == nil {
		return p.Filters
	}
	listing := p.Result.State.Listing()
	out := make([]DataTableFilter, len(p.Filters))
	for i, filter := range p.Filters {
		if spec, ok := listing.Facet(filter.Key); ok {
			if filter.Param == "" {
				filter.Param = spec.Param
			}
			if len(filter.Options) == 0 {
				for _, value := range spec.Values {
					filter.Options = append(filter.Options, SelectOption{Label: value, Value: value})
				}
			}
			if filter.Counts == nil {
				filter.Counts = map[string]int{}
				for _, count := range p.Result.FacetCounts(filter.Key) {
					filter.Counts[count.Value] = count.Count
				}
			}
		}
		out[i] = filter
	}
	return out
}

// ListedSizes are the page sizes the chooser offers: PageSizes, else the sizes
// the listing of the Result allows.
func (p DataTableProps) ListedSizes() []int {
	if len(p.PageSizes) == 0 && p.Result != nil && p.Result.State.Listing() != nil {
		return p.Result.State.Listing().PageSizes()
	}
	return p.PageSizes
}

// names are the parameters the state travels under.
func (p DataTableProps) names() DataTableParams { return p.state().Params().Filled() }

// Action is URL without its query or fragment: the path every address starts
// from, and the action of the form.
func (p DataTableProps) Action() string {
	base, err := url.Parse(p.URL)
	if err != nil {
		base = &url.URL{Path: p.URL}
	}
	base.RawQuery = ""
	base.ForceQuery = false
	base.Fragment = ""
	return base.String()
}

// Href is the address of the list after one step, built from the state's own
// canonical query and nothing else.
func (p DataTableProps) Href(change DataTableChange) string {
	query := p.state().Values(change).Encode()
	if query == "" {
		return p.Action()
	}
	return p.Action() + "?" + query
}

// SelfHref is the address of the list as it is now, which the behaviour writes
// into the address bar after a swap so the bar is canonical whatever the form
// serialised.
func (p DataTableProps) SelfHref() string { return p.Href(DataTableChange{}) }

// Target is what every control swaps: this region, by id.
func (p DataTableProps) Target() string { return "#" + p.ID }

// QueryFormID is the one form the search and the facets submit through.
// SearchID, SizesID and BulkID are the ids of the controls that name or sit
// apart from it.
func (p DataTableProps) QueryFormID() string { return p.ID + "-query" }
func (p DataTableProps) SearchID() string    { return p.ID + "-search" }
func (p DataTableProps) SizesID() string     { return p.ID + "-sizes" }
func (p DataTableProps) BulkID() string      { return p.ID + "-bulk" }

// QueryAction is the action of the query form.
func (p DataTableProps) QueryAction() string { return p.Action() }

// Query is the text searched for.
func (p DataTableProps) Query() string { return p.state().Search() }

// SearchName is the parameter the search travels under.
func (p DataTableProps) SearchName() string { return p.names().Search }

// SizeHref is the address that sets how many rows a page holds.
func (p DataTableProps) SizeHref(size int) string { return p.Href(DataTableChange{PerPage: size}) }

// QueryFields are the hidden fields of the query form: every part of the state
// that no control of the form owns, once each: the order, the page size and the
// hidden columns. The search box owns the search, the facet checkboxes their
// facets, and the page is dropped on purpose -- a new search has no page four.
func (p DataTableProps) QueryFields() []DataTableQueryField {
	names := p.names()
	owned := map[string]bool{names.Search: true, names.Page: true}
	for _, filter := range p.Facets() {
		owned[filter.Name()] = true
	}
	for _, key := range p.controlKeys() {
		owned[key] = true
	}
	values := p.state().Values(DataTableChange{})
	keys := make([]string, 0, len(values))
	for key := range values {
		if !owned[key] {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)

	var fields []DataTableQueryField
	seen := map[DataTableQueryField]bool{}
	for _, key := range keys {
		for _, value := range values[key] {
			field := DataTableQueryField{Name: key, Value: value}
			if !seen[field] {
				seen[field] = true
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// SearchTitle is what the search box is called. Empty uses the region's own
// name, because a box labelled only "Search" on a page with two lists names
// neither.
func (p DataTableProps) SearchTitle() string {
	if p.SearchLabel != "" {
		return p.SearchLabel
	}
	return strings.TrimSpace(strings.ReplaceAll(say(p.PageData, "data_table.search", "Search {label}"), "{label}", p.Label))
}

// Placeholder is the grey text in the empty search box.
func (p DataTableProps) Placeholder() string {
	if p.SearchPlaceholder != "" {
		return p.SearchPlaceholder
	}
	label := strings.ToLower(p.Label)
	sentence := strings.ReplaceAll(p.lines().Filter, "{label}", label)
	return strings.TrimSpace(strings.ReplaceAll(sentence, "  ", " "))
}

// Faceted is whether any facet is drawn.
func (p DataTableProps) Faceted() bool { return len(p.Filters) > 0 }

// Controlled is whether any control is drawn on the search's line.
func (p DataTableProps) Controlled() bool { return len(p.Controls) > 0 }

// Narrowable is whether anything on the search's line narrows the list: the
// search, a facet or a control.
func (p DataTableProps) Narrowable() bool { return p.Searchable || p.Faceted() || p.Controlled() }

// ControlID is the id of the control's field.
func (p DataTableProps) ControlID(control DataTableControl) string {
	return p.ID + "-control-" + control.Key
}

// controlValue is what the state holds under a key, and the fallback when it
// holds nothing.
func (p DataTableProps) controlValue(key, fallback string) string {
	if chosen := p.state().Chosen(key); len(chosen) > 0 {
		return chosen[0]
	}
	return fallback
}

// ControlFrom is the value a control shows: the state's for its Key, and the
// control's own Value when the state has none. A period's start.
func (p DataTableProps) ControlFrom(control DataTableControl) string {
	return p.controlValue(control.Key, control.Value)
}

// ControlTo is the end of a period.
func (p DataTableProps) ControlTo(control DataTableControl) string {
	if control.ToKey == "" {
		return ""
	}
	return p.controlValue(control.ToKey, control.ToValue)
}

// ControlChosen is whether an option of a choice is the one in force.
func (p DataTableProps) ControlChosen(control DataTableControl, option SelectOption) bool {
	return option.Value != "" && option.Value == p.ControlFrom(control)
}

// ControlPeriod is the picker a period control draws: the caller's words, with
// the names, the id, the label and the two days of the control.
func (p DataTableProps) ControlPeriod(control DataTableControl) DateRangePickerProps {
	picker := control.Period
	picker.Name = control.Key
	picker.FromName = control.Key
	picker.ToName = control.ToKey
	picker.ID = p.ControlID(control)
	picker.Label = control.Label
	picker.From = p.ControlFrom(control)
	picker.To = p.ControlTo(control)
	return picker
}

// controlKeys are the query parameters the controls own.
func (p DataTableProps) controlKeys() []string {
	var keys []string
	for _, control := range p.Controls {
		if control.Key != "" {
			keys = append(keys, control.Key)
		}
		if control.Type() == DataTablePeriod && control.ToKey != "" {
			keys = append(keys, control.ToKey)
		}
	}
	return keys
}

// Chosen are the values chosen for one facet.
func (p DataTableProps) Chosen(filter DataTableFilter) []string {
	return p.state().Chosen(filter.Key)
}

// OptionChosen is whether one option of a facet is in force.
func (p DataTableProps) OptionChosen(filter DataTableFilter, option SelectOption) bool {
	return contains(p.Chosen(filter), option.Value)
}

// Active is whether the search or any facet narrows the list, which is when
// Reset is drawn and when an empty page means "nothing matched".
func (p DataTableProps) Active() bool {
	if p.state().Search() != "" {
		return true
	}
	for _, filter := range p.Facets() {
		if len(p.Chosen(filter)) > 0 {
			return true
		}
	}
	for _, key := range p.controlKeys() {
		if len(p.state().Chosen(key)) > 0 {
			return true
		}
	}
	return false
}

// FacetBadges is how many badges the facet button shows: one per chosen value
// while there are two at most.
func (p DataTableProps) FacetBadges(filter DataTableFilter) bool { return len(p.Chosen(filter)) <= 2 }

// FacetCounted is whether the facet draws the number of rows beside each option.
func (p DataTableProps) FacetCounted(filter DataTableFilter) bool { return filter.Counts != nil }

// FacetCount is the number beside one option.
func (p DataTableProps) FacetCount(filter DataTableFilter, option SelectOption) string {
	return strconv.Itoa(filter.Counts[option.Value])
}

// FacetSelectedText is the badge standing for the chosen values of a facet.
func (p DataTableProps) FacetSelectedText(filter DataTableFilter) string {
	return strings.ReplaceAll(p.lines().FacetSelected, "{n}", strconv.Itoa(len(p.Chosen(filter))))
}

// FacetHref is the address that switches one value of a facet on or off.
func (p DataTableProps) FacetHref(filter DataTableFilter, option SelectOption) string {
	return p.Href(DataTableChange{Facet: filter.Key, Value: option.Value})
}

// ClearFacetHref is the address that drops every value of one facet.
func (p DataTableProps) ClearFacetHref(filter DataTableFilter) string {
	return p.Href(DataTableChange{ClearFacet: filter.Key})
}

// ResetHref is the address that drops the search and every facet.
func (p DataTableProps) ResetHref() string { return p.Href(DataTableChange{Reset: true}) }

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

// ColumnsName is what the menu of columns is called.
func (p DataTableProps) ColumnsName() string {
	if p.ColumnsLabel != "" {
		return p.ColumnsLabel
	}
	return say(p.PageData, "table.columns", "Columns")
}

// Off is whether a column is switched off: by the state when there is one, and
// by the column's own Hidden when there is none. A column that cannot be
// switched off is never off.
func (p DataTableProps) Off(column TableColumn) bool {
	if column.Key == "" || !column.Hideable {
		return false
	}
	if p.stateless() {
		return column.Hidden
	}
	return contains(p.state().Hidden(), column.Key)
}

// OffAt is whether the column at position i is switched off.
func (p DataTableProps) OffAt(i int) bool { return i < len(p.Columns) && p.Off(p.Columns[i]) }

// ColumnKey is the key of the column at position i, empty where a row carries
// more cells than the table has columns.
func (p DataTableProps) ColumnKey(i int) string {
	if i < len(p.Columns) {
		return p.Columns[i].Key
	}
	return ""
}

// AlignClass is the alignment of the column at position i, and the default
// where a row carries more cells than the table has columns.
func (p DataTableProps) AlignClass(i int) string {
	if i < len(p.Columns) {
		return p.Columns[i].AlignClass()
	}
	return TableColumn{}.AlignClass()
}

// ColumnHref is the address that switches one column off when it is shown and
// back on when it is hidden.
func (p DataTableProps) ColumnHref(column TableColumn) string {
	return p.Href(DataTableChange{Hide: column.Key})
}

// HideHref is the address that switches a column off.
func (p DataTableProps) HideHref(column TableColumn) string { return p.ColumnHref(column) }

// Sorted is the direction the list is ordered in by this column, and empty
// when it is not ordered by it.
func (p DataTableProps) Sorted(column TableColumn) string {
	if column.Key == "" || column.Key != p.state().SortKey() {
		return ""
	}
	if p.state().SortDir() == "desc" {
		return "desc"
	}
	return "asc"
}

// Sortable is whether the header of a column is a control that orders by it.
func (p DataTableProps) Sortable(column TableColumn) bool { return column.Sortable && column.Key != "" }

// Order is what aria-sort says for one column: the direction when it is the
// one in force, and "none" on every other sortable column. A column that cannot
// be sorted carries nothing.
func (p DataTableProps) Order(column TableColumn) string {
	if !p.Sortable(column) {
		return ""
	}
	switch p.Sorted(column) {
	case "asc":
		return "ascending"
	case "desc":
		return "descending"
	}
	return "none"
}

// NextSort is the direction a click on this header asks for. The cycle is
// ascending, descending, ascending: an ordered list never goes back to no
// order, as the reference it follows does not.
func (p DataTableProps) NextSort(column TableColumn) string {
	if p.Sorted(column) == "asc" {
		return "desc"
	}
	return "asc"
}

// SortHref is where one header points: the next step of its cycle.
func (p DataTableProps) SortHref(column TableColumn) string {
	return p.Href(DataTableChange{Sort: column.Key, Dir: p.NextSort(column)})
}

// AscHref and DescHref are the addresses of the two entries of a header menu.
func (p DataTableProps) AscHref(column TableColumn) string {
	return p.Href(DataTableChange{Sort: column.Key, Dir: "asc"})
}
func (p DataTableProps) DescHref(column TableColumn) string {
	return p.Href(DataTableChange{Sort: column.Key, Dir: "desc"})
}

// SortName is what a sortable header is called, in the state it is in.
func (p DataTableProps) SortName(column TableColumn) string {
	lines := p.lines()
	sentence := lines.Unsorted
	switch p.Sorted(column) {
	case "asc":
		sentence = lines.Ascending
	case "desc":
		sentence = lines.Descending
	}
	return strings.ReplaceAll(sentence, "{column}", column.Label)
}

// Selectable is whether the rows carry a checkbox, and the footer a count.
func (p DataTableProps) Selectable() bool { return p.SelectName != "" }

// Bulky is whether a bar of bulk actions is drawn: there are actions, there is
// somewhere to post them, and there is a selection to post.
func (p DataTableProps) Bulky() bool {
	return p.Selectable() && len(p.BulkActions) > 0 && p.BulkAction != ""
}

// PostMethod is how the bulk bar submits, and it is post unless the caller says
// otherwise -- a bulk operation changes something.
func (p DataTableProps) PostMethod() string {
	if p.BulkMethod != "" {
		return p.BulkMethod
	}
	return "post"
}

// SendsToken is whether the bulk form carries a token: only when there is one
// and the form posts. A form submitted any other way puts its fields in the
// address, which the history, the access log and every proxy between keep.
func (p DataTableProps) SendsToken() bool {
	return p.FormToken() != "" && strings.EqualFold(p.PostMethod(), "post")
}

// FormToken is the token the bulk form posts: Token when it is set, and the
// page's otherwise.
func (p DataTableProps) FormToken() string { return tokenFor(p.PageData, p.Token) }

// BulkName is what the bar of actions is called.
func (p DataTableProps) BulkName() string { return p.lines().Bulk }

// Ticked is how many rows are drawn as selected.
func (p DataTableProps) Ticked() int {
	count := 0
	for _, row := range p.Rows {
		if row.Selected {
			count++
		}
	}
	return count
}

// RowName is what one row's checkbox is called: what the caller wrote, or the
// sentence for a row.
func (p DataTableProps) RowName(row TableRow) string {
	if row.Label != "" {
		return row.Label
	}
	return p.lines().SelectRow
}

// SelectedSentence is the count's sentence with its placeholders still in it:
// SelectedLabel, or the page's kyse.table.selected, or the English.
func (p DataTableProps) SelectedSentence() string {
	if p.SelectedLabel != "" {
		return p.SelectedLabel
	}
	return say(p.PageData, "table.selected", "{n} of {total} row(s) selected.")
}

// Known is whether the size of the whole list is known: always for a complete
// list, and for a windowed one when its page counts.
func (p DataTableProps) Known() bool {
	if p.Complete {
		return true
	}
	_, ok := counted(p.paginator())
	return ok
}

// Total is how many rows the whole list has, narrowed by the search and the
// facets: the count the query ran. Zero when it is not known.
func (p DataTableProps) Total() int {
	if p.Complete {
		return len(p.Rows)
	}
	if page, ok := counted(p.paginator()); ok {
		return page.Total()
	}
	return 0
}

// FooterText is the sentence on the left of the footer: how many rows are
// selected of the whole list, or how many rows there are.
func (p DataTableProps) FooterText() string {
	total := strconv.Itoa(p.Total())
	if !p.Selectable() {
		if !p.Known() {
			return ""
		}
		return p.lines().RowsFor(p.Total())
	}
	if !p.Known() {
		return strings.ReplaceAll(p.lines().SelectedCount, "{n}", strconv.Itoa(p.Ticked()))
	}
	return strings.NewReplacer("{n}", strconv.Itoa(p.Ticked()), "{total}", total).Replace(p.SelectedSentence())
}

// Nothing is whether the list has nothing at all, which is the caller's Empty
// to say. A list that has no rows on this page, or none that match, is not
// nothing: it says "No results." in the table.
func (p DataTableProps) Nothing() bool {
	if len(p.Rows) > 0 || p.Active() {
		return false
	}
	if page, ok := counted(p.paginator()); ok && page.Total() > 0 {
		return false
	}
	return true
}

// Span is how many columns a full-width row spans.
func (p DataTableProps) Span() int {
	span := len(p.Columns)
	if p.Selectable() {
		span++
	}
	if p.Acting() {
		span++
	}
	return span
}

// Acting is whether any row carries a menu of actions.
func (p DataTableProps) Acting() bool {
	for _, row := range p.Rows {
		if len(row.Actions) > 0 {
			return true
		}
	}
	return false
}

// Role is what the table announces as: a grid when it is navigated cell by
// cell, a treegrid when those rows nest, and nothing at all otherwise.
func (p DataTableProps) Role() string {
	if !p.Navigable {
		return ""
	}
	for _, row := range p.Rows {
		if row.Level > 0 {
			return "treegrid"
		}
	}
	return "grid"
}

// RowAttrs are the attributes of a row: its facet values, as the attributes
// the browser filters a complete list by, and then the caller's, which win on a
// name both write.
func (p DataTableProps) RowAttrs(row TableRow) map[string]string {
	caller := p.PartAttrs("row")
	if len(row.Facets) == 0 {
		return caller
	}
	out := make(map[string]string, len(row.Facets)+len(caller))
	for key, value := range row.Facets {
		out["data-facet-"+key] = value
	}
	for key, value := range caller {
		out[key] = value
	}
	return out
}

// The classes that draw the table as a list of cards. Every one of them is the
// same utility written under a container query on the region -- "@max-[600px]/dt:"
// is "when the region is narrower than 600px" -- so above the breakpoint none of
// them applies and the table is the table it was. They are written out whole
// because Tailwind compiles the class names it finds as text: a class put
// together at run time from a prefix and a name would never reach the stylesheet.
//
// The region is the container, not the window, because the width that matters is
// the one the table is given: a table in the content column of a page with its
// sidebar open needs cards on a window much wider than a phone.
const (
	// cardRegion makes the region the container the queries ask.
	cardRegion = "@container/dt"
	// cardWrapper, cardScroller and cardTable stop being a bordered, scrolling
	// table: the cards carry their own border and nothing scrolls sideways.
	cardWrapper  = "@max-[600px]/dt:overflow-visible @max-[600px]/dt:rounded-none @max-[600px]/dt:border-0"
	cardScroller = "@max-[600px]/dt:overflow-visible"
	cardTable    = "@max-[600px]/dt:block"
	cardCaption  = "@max-[600px]/dt:block"
	// cardHead is the head turned into a bar of the controls that order the list,
	// and cardHeadRow the row it holds them in.
	cardHead    = "@max-[600px]/dt:block"
	cardHeadRow = "@max-[600px]/dt:flex @max-[600px]/dt:flex-wrap @max-[600px]/dt:items-center @max-[600px]/dt:gap-2 @max-[600px]/dt:border-b-0 @max-[600px]/dt:pb-3"
	// cardHeadCell keeps a header that orders the list, cardHeadSelect the
	// checkbox that ticks every row, and cardHeadOff drops a header that only
	// names a column, which each card names for itself.
	cardHeadCell   = "@max-[600px]/dt:block @max-[600px]/dt:h-auto @max-[600px]/dt:p-0"
	cardHeadSelect = "@max-[600px]/dt:block @max-[600px]/dt:h-auto @max-[600px]/dt:w-auto @max-[600px]/dt:p-0 @max-[600px]/dt:pe-1"
	cardHeadOff    = "@max-[600px]/dt:hidden"
	// cardBody stacks the cards. The rule that takes the border off the last row
	// of a table is put back, because the last card has four sides.
	cardBody = "@max-[600px]/dt:flex @max-[600px]/dt:flex-col @max-[600px]/dt:gap-3 @max-[600px]/dt:[&_tr:last-child]:border"
	// cardRow is one card.
	cardRow = "@max-[600px]/dt:relative @max-[600px]/dt:flex @max-[600px]/dt:flex-col @max-[600px]/dt:gap-1.5 @max-[600px]/dt:rounded-lg @max-[600px]/dt:border @max-[600px]/dt:bg-card @max-[600px]/dt:p-3"
	// cardCell is a line of a card: the column's name in front, drawn from the
	// cell's own data-label so the text of the cell is still only its value, and
	// the value after it.
	cardCell = "@max-[600px]/dt:flex @max-[600px]/dt:items-baseline @max-[600px]/dt:gap-x-3 @max-[600px]/dt:p-0 @max-[600px]/dt:text-start @max-[600px]/dt:whitespace-normal @max-[600px]/dt:before:w-2/5 @max-[600px]/dt:before:shrink-0 @max-[600px]/dt:before:text-xs @max-[600px]/dt:before:text-muted-foreground @max-[600px]/dt:before:content-[attr(data-label)]"
	// cardNoLabel is a cell with no column name to put in front of it.
	cardNoLabel = "@max-[600px]/dt:before:hidden"
	// cardTitle is the cell of the primary column, first and larger, with no name
	// in front; cardTitleSelect leaves room for the checkbox.
	cardTitle       = "@max-[600px]/dt:order-first @max-[600px]/dt:text-base @max-[600px]/dt:font-medium @max-[600px]/dt:before:hidden"
	cardTitleSelect = "@max-[600px]/dt:ps-8"
	// cardLeft is a column left out of the card.
	cardLeft = "@max-[600px]/dt:hidden"
	// cardSelect is the row's checkbox, at the top corner of its card.
	cardSelect = "@max-[600px]/dt:absolute @max-[600px]/dt:start-3 @max-[600px]/dt:top-3 @max-[600px]/dt:p-0"
	// cardActionCell is a cell of the column marked Actions: the foot of the card,
	// at its end and under a rule, with no name in front of the controls.
	cardActionCell = "@max-[600px]/dt:order-last @max-[600px]/dt:flex @max-[600px]/dt:items-center @max-[600px]/dt:justify-end @max-[600px]/dt:gap-x-3 @max-[600px]/dt:border-t @max-[600px]/dt:p-0 @max-[600px]/dt:pt-2 @max-[600px]/dt:text-start @max-[600px]/dt:whitespace-normal @max-[600px]/dt:before:hidden"
	// cardActions is the row's menu, along the foot of the card and in sight.
	cardActions = "@max-[600px]/dt:order-last @max-[600px]/dt:flex @max-[600px]/dt:justify-end @max-[600px]/dt:border-t @max-[600px]/dt:p-0 @max-[600px]/dt:pt-2"
	// cardFull is a row that spans the table, the one that says nothing matched.
	cardFull = "@max-[600px]/dt:block"
)

// PrimaryColumn is the position of the column that titles a card: the first one
// marked Primary, and with none marked the first column that is neither left out
// of the card nor switched off. It is -1 for a table with no column that can.
func (p DataTableProps) PrimaryColumn() int {
	for i, column := range p.Columns {
		if column.Primary {
			return i
		}
	}
	for i, column := range p.Columns {
		if !column.CardHidden && !p.Off(column) {
			return i
		}
	}
	return -1
}

// CardHiddenAt is whether column i is left out of a card. The title never is.
func (p DataTableProps) CardHiddenAt(i int) bool {
	return i < len(p.Columns) && p.Columns[i].CardHidden && i != p.PrimaryColumn()
}

// LabelAt is the name of column i, which a card writes in front of its value.
// It is empty where a row carries more cells than the table has columns.
func (p DataTableProps) LabelAt(i int) string {
	if i < len(p.Columns) {
		return p.Columns[i].Label
	}
	return ""
}

// ActionsAt is whether column i holds the row's own controls (Actions).
func (p DataTableProps) ActionsAt(i int) bool {
	return i < len(p.Columns) && p.Columns[i].Actions
}

// CardCellFor is the classes cell i draws with once the region is narrower than
// the breakpoint, given the cell itself: a cell with nothing in it is left out of
// the card, so a name is never drawn in front of nothing and a row with no action
// has no foot. The title of the card is drawn even when it is empty.
func (p DataTableProps) CardCellFor(i int, cell TableCell) string {
	if cell.Empty() && i != p.PrimaryColumn() {
		return cardLeft
	}
	return p.CardCell(i)
}

// CardCell is the classes cell i draws with once the region is narrower than the
// breakpoint, which are none above it. The cell's own content is not known here:
// CardCellFor leaves an empty one out.
func (p DataTableProps) CardCell(i int) string {
	switch {
	case p.CardHiddenAt(i):
		return cardLeft
	case p.ActionsAt(i) && i != p.PrimaryColumn():
		return cardActionCell
	case i == p.PrimaryColumn():
		if p.Selectable() {
			return joinClasses(cardCell, cardTitle, cardTitleSelect)
		}
		return joinClasses(cardCell, cardTitle)
	case p.LabelAt(i) == "":
		return joinClasses(cardCell, cardNoLabel)
	}
	return cardCell
}

// CardHead is the classes the header of a column draws with in a card layout:
// kept when it orders the list, dropped when it only names the column.
func (p DataTableProps) CardHead(column TableColumn) string {
	if p.Sortable(column) {
		return cardHeadCell
	}
	return cardHeadOff
}

// CardActions is the classes a row's last cell draws with in a card layout: the
// foot of the card when the row has a menu, and nothing at all when it has none.
func (p DataTableProps) CardActions(row TableRow) string {
	if len(row.Actions) == 0 {
		return cardLeft
	}
	return cardActions
}

// dataTableLine is one line of the body as it is drawn: the row, whether it
// lies outside the window of a complete list, and the id of its menu.
type dataTableLine struct {
	TableRow
	// Outside marks a row of a complete list that is not on the page being
	// shown, so the first paint is already the page and a page with no script
	// does not contradict its own footer.
	Outside bool
	// MenuID is the id of the row's menu, which is what lets it stay open when
	// the region is swapped.
	MenuID string
}

// Body is the rows as they are drawn.
func (p DataTableProps) Body() []dataTableLine {
	lines := make([]dataTableLine, len(p.Rows))
	first := (p.At()-1)*p.Size() + 1
	for at, row := range p.Rows {
		line := dataTableLine{TableRow: row, Outside: row.Hidden}
		if p.Complete && p.Size() > 0 {
			line.Outside = row.Hidden || at+1 < first || at+1 >= first+p.Size()
		}
		line.MenuID = p.ID + "-row-" + strconv.Itoa(at)
		if row.Key != "" {
			line.MenuID = p.ID + "-row-" + row.Key
		}
		lines[at] = line
	}
	return lines
}

// Size is how many rows a page holds: the state's, else PageSize, else the
// page the query returned.
func (p DataTableProps) Size() int {
	if size := p.state().PageSize(); size > 0 {
		return size
	}
	if p.PageSize > 0 {
		return p.PageSize
	}
	if page, ok := counted(p.paginator()); ok {
		if sized, ok := page.(interface{ PerPage() int }); ok {
			return sized.PerPage()
		}
	}
	return 0
}

// Sizes are the counts the chooser offers, with the one in force among them.
func (p DataTableProps) Sizes() []int {
	if len(p.ListedSizes()) == 0 {
		return nil
	}
	sizes := append([]int{}, p.ListedSizes()...)
	if size := p.Size(); size > 0 {
		found := false
		for _, one := range sizes {
			found = found || one == size
		}
		if !found {
			sizes = append(sizes, size)
			sort.Ints(sizes)
		}
	}
	return sizes
}

// At is the page being shown.
func (p DataTableProps) At() int {
	if !p.Complete {
		if page, ok := presentPaginator(p.paginator()).(numberedPage); ok {
			return max(1, page.CurrentPage())
		}
	}
	return max(1, p.state().Page())
}

// Sheets is how many pages there are. A complete list counts its own: the
// server sent every row, so how many pages that makes is arithmetic.
func (p DataTableProps) Sheets() int {
	if p.Complete {
		if p.Size() == 0 || len(p.Rows) == 0 {
			return 1
		}
		return (len(p.Rows) + p.Size() - 1) / p.Size()
	}
	if page, ok := counted(p.paginator()); ok {
		return max(1, page.LastPage())
	}
	return 0
}

// Paged is whether there is more than one page to move between.
func (p DataTableProps) Paged() bool {
	if page := p.uncountedPage(); page != nil {
		return page.HasPages()
	}
	return p.Sheets() > 1
}

// uncountedPage is the paginator of a windowed list when it knows no total --
// a simple page or a cursor page -- and nil otherwise.
func (p DataTableProps) uncountedPage() Paginator {
	paginator := presentPaginator(p.paginator())
	if p.Complete || paginator == nil {
		return nil
	}
	if _, ok := paginator.(countedPage); ok {
		return nil
	}
	return paginator
}

// Stepped is whether the pager is drawn: there is more than one page, or the
// list is complete and the browser may filter it down to one and back.
func (p DataTableProps) Stepped() bool { return p.Paged() || p.Complete }

// PageHref is the address of a page, kept inside the list.
func (p DataTableProps) PageHref(page int) string {
	return p.Href(DataTableChange{Page: max(1, min(page, max(1, p.Sheets())))})
}

// HasPrevious and HasNext say whether a step exists.
func (p DataTableProps) HasPrevious() bool {
	if page := p.uncountedPage(); page != nil {
		return p.steps().previous != ""
	}
	return p.At() > 1
}
func (p DataTableProps) HasNext() bool {
	if page := p.uncountedPage(); page != nil {
		return p.steps().next != ""
	}
	return p.At() < p.Sheets()
}

// FirstHref, PreviousHref, NextHref and LastHref are the four steps.
func (p DataTableProps) FirstHref() string    { return p.PageHref(1) }
func (p DataTableProps) LastHref() string     { return p.PageHref(p.Sheets()) }
func (p DataTableProps) PreviousHref() string {
	if p.uncountedPage() != nil {
		return p.steps().previous
	}
	return p.PageHref(p.At() - 1)
}
func (p DataTableProps) NextHref() string {
	if p.uncountedPage() != nil {
		return p.steps().next
	}
	return p.PageHref(p.At() + 1)
}

// steps are the two steps of a page that knows no total, each built from the
// table's own state like every other control, and empty where there is nowhere
// to go.
func (p DataTableProps) steps() steps {
	var out steps
	switch page := p.uncountedPage().(type) {
	case *pagination.CursorPage:
		cursor := func(at *pagination.Cursor) string {
			signer := page.GetOptions().Signer
			if at == nil || signer == nil {
				return ""
			}
			token := signer.Encode(*at)
			if token == "" {
				return ""
			}
			return p.Href(DataTableChange{CursorName: page.GetCursorName(), Cursor: token})
		}
		out.previous = cursor(page.PreviousCursor())
		out.next = cursor(page.NextCursor())
	case numberedPage:
		at := page.CurrentPage()
		if at > 1 {
			out.previous = p.Href(DataTableChange{Page: at - 1})
		}
		if page.HasMorePages() {
			out.next = p.Href(DataTableChange{Page: at + 1})
		}
	case nil:
	default:
		// A page of a kind this does not know the parameters of has only its
		// own addresses to offer.
		out.previous = page.PreviousPageURL()
		out.next = page.NextPageURL()
	}
	return out
}

// PageOfText is where the reader is, beside the pager: "Page 2 of 7". It is
// empty for a page that knows no total, because "Page 3 of" with nothing after
// it says less than nothing.
func (p DataTableProps) PageOfText() string {
	if !p.Known() {
		return ""
	}
	return strings.NewReplacer("{n}", strconv.Itoa(p.At()), "{total}", strconv.Itoa(p.Sheets())).Replace(p.lines().PageOf)
}

// PagerName is what the pager landmark is called.
func (p DataTableProps) PagerName() string {
	return PaginationProps{Label: p.Label, PageData: p.PageData}.LandmarkName()
}

// Chooses is whether the rows-per-page chooser is drawn.
func (p DataTableProps) Chooses() bool { return len(p.ListedSizes()) > 0 }

// Bridge is the behaviour's props: what the browser needs and cannot read from
// the markup.
func (p DataTableProps) Bridge() map[string]any {
	props := p.lines().Sentences()
	props["selected"] = p.SelectedSentence()
	props["complete"] = p.Complete
	props["pageSize"] = p.Size()
	props["selectName"] = p.SelectName
	return props
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own, and the facts the
// browser reads off the region.
func (p DataTableProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: TableBehavior, Props: p.Bridge()}
	}
	attrs := p.ComponentProps.RootAttrs()
	attrs["data-total"] = strconv.Itoa(p.Total())
	attrs["data-url"] = p.SelfHref()
	return attrs
}

// PartNames are the parts this component publishes.
func (p DataTableProps) PartNames() []string {
	return []string{
		"root", "toolbar", "query", "search", "reset", "facet", "facet-chosen", "facet-search", "filter",
		"control", "control-input", "columns", "toggle", "bulk", "wrapper", "table", "caption", "head",
		"header-cell", "sort", "header-menu", "row", "cell", "select-all",
		"select", "row-menu", "footer", "count", "sizes", "pagination", "page-of",
		"first", "previous", "next", "last",
	}
}
@endgo

{{-- One region, one id, one target. Every control swaps this element, which is
     why the search can carry the order and the order can carry the search: they
     are all rebuilt from the same state on the way back. --}}
<div
	data-part="root"
	class="{{ .RootClass("group/dt w-full", cardRegion) }}"
	id="{{ .ID }}"
	role="region"
	@if(.Complete)
		data-complete="true"
	@endif
	@if(.Navigable)
		data-navigable="true"
	@endif
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(.Nothing())
		{!! Empty(.Empty) !!}
	@else
	<div
		data-part="toolbar"
		class="{{ .PartClass("toolbar", "flex items-start gap-2 py-4") }}"
		@attributes(.PartAttrs("toolbar"))
	>
		{{-- The search and the facets are one form. With no script
		     it is a native GET; with htmx it is the single request every one of
		     them makes. The order and the hidden columns ride in it as fields,
		     once each, from the state. --}}
		<form
			id="{{ .QueryFormID() }}"
			data-part="query"
			class="{{ .PartClass("query", "flex min-w-0 flex-1 flex-wrap items-center gap-2") }}"
			action="{{ .QueryAction() }}"
			method="get"
			@if(.Searchable)
				role="search"
			@endif
			@if(.Complete)
				hx-boost="false"
			@endif
			@if(!.Complete)
				hx-get="{{ .QueryAction() }}"
				hx-trigger="submit, change"
				hx-target="{{ .Target() }}"
				hx-swap="outerHTML"
				hx-select="{{ .Target() }}"
				hx-push-url="true"
				hx-sync="closest [data-part=root]:replace"
				hx-indicator="{{ .Target() }}"
			@endif
			@attributes(.PartAttrs("query"))
		>
			@foreach(.QueryFields() as field)
				<input type="hidden" name="{{ field.Name }}" value="{{ field.Value }}">
			@endforeach
			@if(.Searchable)
				<label class="sr-only" for="{{ .SearchID() }}">{{ .SearchTitle() }}</label>
				@if(.Faceted())
					<input
						data-part="search"
						class="{{ .PartClass("search", "h-8 w-[150px] lg:w-[250px] rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30") }}"
						type="text"
						role="searchbox"
						id="{{ .SearchID() }}"
						name="{{ .SearchName() }}"
						value="{{ .Query() }}"
						autocomplete="off"
						placeholder="{{ .Placeholder() }}"
						@attributes(.PartAttrs("search"))
					>
				@endif
				@if(!.Faceted())
					<input
						data-part="search"
						class="{{ .PartClass("search", "h-8 w-full min-w-0 max-w-sm rounded-lg border border-input bg-transparent px-2.5 py-1 text-base transition-colors outline-none placeholder:text-muted-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 md:text-sm dark:bg-input/30") }}"
						type="text"
						role="searchbox"
						id="{{ .SearchID() }}"
						name="{{ .SearchName() }}"
						value="{{ .Query() }}"
						autocomplete="off"
						placeholder="{{ .Placeholder() }}"
						@attributes(.PartAttrs("search"))
					>
				@endif
			@endif

			@foreach(.Facets() as filter)
				{{-- A facet opens a popover of checkboxes, as shadcn's faceted
				     filter: the title, and what is chosen beside it. The boxes are
				     native and in the form, so ticking one is a change of the form.
				     Every option is drawn with its count; the search inside the
				     popover only hides the ones that do not match. --}}
				<details
					id="{{ .ID }}-facet-{{ filter.Key }}"
					data-part="facet"
					data-menu
					class="{{ .PartClass("facet", "group/menu relative") }}"
					@attributes(.PartAttrs("facet"))
				>
					<summary class="{{ "inline-flex cursor-pointer list-none [&::-webkit-details-marker]:hidden group/button shrink-0 items-center justify-center rounded-lg border border-dashed border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none h-8 gap-1.5 px-2.5 hover:bg-muted hover:text-foreground group-open/menu:bg-muted focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50" }}">
						{!! icons.PlusCircle(icons.Props{}) !!}
						<span>{{ filter.Label }}</span>
						<span
							class="{{ "mx-2 h-4 w-px shrink-0 bg-border" }}"
							role="separator"
							aria-orientation="vertical"
							@if(len(.Chosen(filter)) == 0)
								hidden
							@endif
						></span>
						<span
							data-part="facet-chosen"
							class="{{ .PartClass("facet-chosen", "inline-flex items-center gap-1") }}"
							@if(len(.Chosen(filter)) == 0)
								hidden
							@endif
							@attributes(.PartAttrs("facet-chosen"))
						>
							<span
								data-facet-count
								class="rounded-sm bg-secondary px-1 text-xs font-normal text-secondary-foreground lg:hidden"
							>{{ len(.Chosen(filter)) }}</span>
							<span class="hidden gap-1 lg:flex">
								<span
									data-facet-many
									class="rounded-sm bg-secondary px-1 text-xs font-normal text-secondary-foreground"
									@if(.FacetBadges(filter))
										hidden
									@endif
								>{{ .FacetSelectedText(filter) }}</span>
								@foreach(filter.Options as option)
									<span
										data-facet-badge="{{ option.Value }}"
										class="rounded-sm bg-secondary px-1 text-xs font-normal text-secondary-foreground"
										@if(!.FacetBadges(filter) || !.OptionChosen(filter, option))
											hidden
										@endif
									>{{ option.Label }}</span>
								@endforeach
							</span>
						</span>
					</summary>
					<div
						data-menu-panel
						data-align="start"
						class="{{ "fixed z-50 w-[200px] overflow-hidden rounded-lg bg-popover p-0 text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none" }}"
					>
						<div class="flex h-9 items-center gap-2 border-b px-3">
							<span class="shrink-0 opacity-50 [&_svg]:size-4" aria-hidden="true">{!! icons.MagnifyingGlass(icons.Props{}) !!}</span>
							<input
								data-part="facet-search"
								class="{{ .PartClass("facet-search", "flex h-9 w-full min-w-0 bg-transparent py-3 text-sm outline-hidden placeholder:text-muted-foreground") }}"
								type="text"
								autocomplete="off"
								placeholder="{{ filter.Label }}"
								aria-label="{{ filter.Label }}"
								@attributes(.PartAttrs("facet-search"))
							>
						</div>
						<div class="max-h-[300px] overflow-x-hidden overflow-y-auto p-1" role="group" aria-label="{{ filter.Label }}">
							@foreach(filter.Options as option)
								<label
									data-facet-item
									data-text="{{ option.Label }}"
									class="relative flex cursor-default items-center gap-2 rounded-sm px-2 py-1.5 text-sm select-none hover:bg-accent hover:text-accent-foreground has-focus-visible:bg-accent has-focus-visible:text-accent-foreground has-disabled:opacity-50"
								>
									<input
										data-part="filter"
										class="{{ .PartClass("filter", "peer relative block size-4 shrink-0 cursor-pointer appearance-none rounded-[4px] border border-input outline-none transition-colors after:absolute after:top-[42%] after:left-1/2 after:hidden after:h-[9px] after:w-[5px] after:-translate-x-1/2 after:-translate-y-1/2 after:rotate-45 after:border-r-2 after:border-b-2 after:border-current after:content-[''] checked:border-primary checked:bg-primary checked:text-primary-foreground checked:after:block disabled:cursor-not-allowed dark:bg-input/30 dark:checked:bg-primary") }}"
										type="checkbox"
										name="{{ filter.Name() }}"
										value="{{ option.Value }}"
										@if(.OptionChosen(filter, option))
											checked
										@endif
										@if(option.Disabled)
											disabled
										@endif
										@attributes(.PartAttrs("filter"))
									>
									<span class="min-w-0 flex-1 truncate">{{ option.Label }}</span>
									@if(.FacetCounted(filter))
										<span class="ml-auto flex size-4 items-center justify-center font-mono text-xs text-muted-foreground">{{ .FacetCount(filter, option) }}</span>
									@endif
								</label>
							@endforeach
							<div data-facet-empty hidden class="py-6 text-center text-sm">{{ .lines().NoFacetResults }}</div>
						</div>
						@if(len(.Chosen(filter)) > 0)
							<div class="border-t p-1" data-facet-clear>
								<a
									class="flex cursor-pointer items-center justify-center rounded-sm px-2 py-1.5 text-center text-sm select-none hover:bg-accent hover:text-accent-foreground"
									href="{{ .ClearFacetHref(filter) }}"
									@if(.Complete)
										hx-boost="false"
									@endif
									@if(!.Complete)
										hx-get="{{ .ClearFacetHref(filter) }}"
										hx-target="{{ .Target() }}"
										hx-swap="outerHTML"
										hx-select="{{ .Target() }}"
										hx-push-url="true"
										hx-sync="closest [data-part=root]:replace"
										hx-indicator="{{ .Target() }}"
									@endif
									data-facet-clear-link
								>{{ .lines().ClearFilters }}</a>
							</div>
						@endif
					</div>
				</details>
			@endforeach

			@foreach(.Controls as control)
				{{-- A control is a facet that is not a list of boxes. It sits in the
				     query form with the search, so its value is a field of the one
				     request, and its key is in the state, so every address the
				     table draws carries it. --}}
				<div
					data-part="control"
					data-control="{{ control.Type() }}"
					class="{{ .PartClass("control", "flex min-w-0 items-center gap-2 max-sm:w-full [&_.field]:flex [&_.field]:flex-row [&_.field]:items-center [&_.field]:gap-2") }}"
					@attributes(.PartAttrs("control"))
				>
					@if(control.Type() == "period")
						{!! DateRangePicker(.ControlPeriod(control)) !!}
					@endif
					@if(control.Type() == "select")
						<label class="{{ "text-sm font-medium whitespace-nowrap text-muted-foreground" }}" for="{{ .ControlID(control) }}">{{ control.Label }}</label>
						<select
							data-part="control-input"
							class="{{ .PartClass("control-input", "select h-8 w-44 max-sm:w-full py-0 text-sm") }}"
							id="{{ .ControlID(control) }}"
							name="{{ control.Key }}"
							@attributes(.PartAttrs("control-input"))
						>
							<option value="">{{ control.Placeholder }}</option>
							@foreach(control.Options as option)
								<option
									value="{{ option.Value }}"
									@if(.ControlChosen(control, option))
										selected
									@endif
									@if(option.Disabled)
										disabled
									@endif
								>{{ option.Label }}</option>
							@endforeach
						</select>
					@endif
					@if(control.Type() == "text")
						<label class="{{ "text-sm font-medium whitespace-nowrap text-muted-foreground" }}" for="{{ .ControlID(control) }}">{{ control.Label }}</label>
						<input
							data-part="control-input"
							class="{{ .PartClass("control-input", "input h-8 w-44 max-sm:w-full py-0 text-sm") }}"
							type="search"
							id="{{ .ControlID(control) }}"
							name="{{ control.Key }}"
							value="{{ .ControlFrom(control) }}"
							placeholder="{{ control.Placeholder }}"
							autocomplete="off"
							enterkeyhint="search"
							@attributes(.PartAttrs("control-input"))
						>
					@endif
				</div>
			@endforeach

			@if(.Active() && .Narrowable())
				<a
					data-part="reset"
					class="{{ .PartClass("reset", "inline-flex group/button shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none h-8 gap-1.5 px-2.5 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:hover:bg-muted/50") }}"
					href="{{ .ResetHref() }}"
					@if(.Complete)
						hx-boost="false"
					@endif
					@if(!.Complete)
						hx-get="{{ .ResetHref() }}"
						hx-target="{{ .Target() }}"
						hx-swap="outerHTML"
						hx-select="{{ .Target() }}"
						hx-push-url="true"
						hx-sync="closest [data-part=root]:replace"
						hx-indicator="{{ .Target() }}"
					@endif
					@attributes(.PartAttrs("reset"))
				>{{ .lines().Reset }}{!! icons.X(icons.Props{}) !!}</a>
			@endif
			@if(.Narrowable())
				<noscript><button type="submit" class="{{ "inline-flex group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background text-sm font-medium h-7 px-2.5" }}">{{ .lines().Apply }}</button></noscript>
			@endif
		</form>

		@if(.Hideable())
			{{-- Columns are switched on and off from a menu that stays open, as
			     shadcn's: each line is the address of the state with that column
			     toggled, so the choice lives in the address and survives the swap,
			     the reload and the history. --}}
			<details
				id="{{ .ID }}-columns"
				data-part="columns"
				data-menu
				class="{{ .PartClass("columns", "group/menu relative ms-auto") }}"
				@attributes(.PartAttrs("columns"))
			>
				<summary class="{{ "inline-flex cursor-pointer list-none [&::-webkit-details-marker]:hidden group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none h-8 gap-1.5 px-2.5 hover:bg-muted hover:text-foreground group-open/menu:bg-muted group-open/menu:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50" }}">{{ .ColumnsName() }}{!! icons.CaretDown(icons.Props{}) !!}</summary>
				<div
					data-menu-panel
					data-align="end"
					role="menu"
					aria-label="{{ .lines().ToggleColumns }}"
					class="{{ "fixed z-50 w-44 min-w-32 overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none" }}"
				>
					@foreach(.Columns as column)
						@if(column.Hideable && column.Key != "")
							<a
								data-part="toggle"
								data-column-toggle="{{ column.Key }}"
								class="{{ .PartClass("toggle", "relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pe-8 ps-1.5 text-sm capitalize outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground") }}"
								role="menuitemcheckbox"
								aria-checked="{{ !.Off(column) }}"
								href="{{ .ColumnHref(column) }}"
								@if(.Complete)
									hx-boost="false"
								@endif
								@if(!.Complete)
									hx-get="{{ .ColumnHref(column) }}"
									hx-target="{{ .Target() }}"
									hx-swap="outerHTML"
									hx-select="{{ .Target() }}"
									hx-push-url="true"
									hx-sync="closest [data-part=root]:replace"
									hx-indicator="{{ .Target() }}"
								@endif
								@attributes(.PartAttrs("toggle"))
							>
								{{ column.Label }}
								<span
									class="pointer-events-none absolute end-2 flex items-center justify-center [&_svg]:size-4 [&_path[opacity]]:hidden"
									data-indicator
									@if(.Off(column))
										hidden
									@endif
								>{!! icons.Check(icons.Props{}) !!}</span>
							</a>
						@endif
					@endforeach
				</div>
			</details>
		@endif
	</div>

	@if(.Bulky())
		{{-- The bar of bulk actions. It is shown once something is ticked, and the
		     keys ticked on other pages travel in it as hidden fields the behaviour
		     keeps, so one press posts the whole selection. --}}
		<form
			id="{{ .BulkID() }}"
			data-part="bulk"
			class="{{ .PartClass("bulk", "mb-4 flex flex-wrap items-center gap-2 rounded-lg bg-accent px-3 py-2 text-accent-foreground") }}"
			method="{{ .PostMethod() }}"
			action="{{ .BulkAction }}"
			aria-label="{{ .BulkName() }}"
			@if(.Ticked() == 0)
				hidden
			@endif
			@if(.Complete)
				hx-boost="false"
			@endif
			@attributes(.PartAttrs("bulk"))
		>
			@if(.SendsToken())
				<input type="hidden" name="_token" value="{{ .FormToken() }}">
			@endif
			@foreach(.BulkActions as action)
				{!! Button(action) !!}
			@endforeach
		</form>
	@endif

	<div
		data-part="wrapper"
		class="{{ .PartClass("wrapper", "overflow-hidden rounded-md border transition-opacity group-[.htmx-request]/dt:opacity-60", cardWrapper) }}"
		@attributes(.PartAttrs("wrapper"))
	>
		<div class="{{ "relative w-full overflow-x-auto " + cardScroller }}">
			<table
				data-part="table"
				class="{{ .PartClass("table", "w-full caption-bottom text-sm", cardTable) }}"
				@if(.Role() != "")
					role="{{ .Role() }}"
				@endif
				@attributes(.PartAttrs("table"))
			>
				@if(.Caption != "")
					<caption
						data-part="caption"
						class="{{ .PartClass("caption", "mt-4 text-sm text-muted-foreground", cardCaption) }}"
						@attributes(.PartAttrs("caption"))
					>{{ .Caption }}</caption>
				@endif
				<thead
					data-part="head"
					class="{{ .PartClass("head", "[&_tr]:border-b", cardHead) }}"
					@attributes(.PartAttrs("head"))
				>
					<tr class="{{ "border-b transition-colors hover:bg-muted/50 has-[details[open]]:bg-muted/50 data-[state=selected]:bg-muted " + cardHeadRow }}">
						@if(.Selectable())
							<th
								data-part="header-cell"
								scope="col"
								class="{{ .PartClass("header-cell", "h-10 px-2 text-start align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pe-0 [&:has([data-select])]:pe-0 w-0", cardHeadSelect) }}"
								data-select
								@attributes(.PartAttrs("header-cell"))
							>
								<input
									data-part="select-all"
									data-select-all
									class="{{ .PartClass("select-all", "peer relative block size-4 shrink-0 cursor-pointer appearance-none rounded-[4px] border border-input outline-none transition-colors after:absolute after:top-[42%] after:left-1/2 after:hidden after:h-[9px] after:w-[5px] after:-translate-x-1/2 after:-translate-y-1/2 after:rotate-45 after:border-r-2 after:border-b-2 after:border-current after:content-[''] checked:border-primary checked:bg-primary checked:text-primary-foreground checked:after:block indeterminate:after:block focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-input/30 dark:checked:bg-primary") }}"
									type="checkbox"
									aria-label="{{ .lines().SelectAll }}"
									@attributes(.PartAttrs("select-all"))
								>
							</th>
						@endif
						@foreach(.Columns as column)
							<th
								data-part="header-cell"
								scope="col"
								class="{{ .PartClass("header-cell", "h-10 px-2 align-middle font-medium whitespace-nowrap text-foreground [&:has([role=checkbox])]:pe-0 [&:has([data-select])]:pe-0", column.AlignClass(), .CardHead(column)) }}"
								@if(column.Key != "")
									data-column="{{ column.Key }}"
								@endif
								@if(.Off(column))
									hidden
								@endif
								@if(.Order(column) != "")
									aria-sort="{{ .Order(column) }}"
								@endif
								@attributes(.PartAttrs("header-cell"))
							>
								@if(.Sortable(column) && !.HeaderMenu)
									<a
										data-part="sort"
										class="{{ .PartClass("sort", "inline-flex group/button shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none h-8 gap-1.5 px-2.5 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:hover:bg-muted/50") }}"
										href="{{ .SortHref(column) }}"
										aria-label="{{ .SortName(column) }}"
										data-sort-key="{{ column.Key }}"
										@if(.Complete)
											hx-boost="false"
										@endif
										@if(!.Complete)
											hx-get="{{ .SortHref(column) }}"
											hx-target="{{ .Target() }}"
											hx-swap="outerHTML"
											hx-select="{{ .Target() }}"
											hx-push-url="true"
											hx-sync="closest [data-part=root]:replace"
											hx-indicator="{{ .Target() }}"
										@endif
										@attributes(.PartAttrs("sort"))
									>
										{{ column.Label }}
										<span
											data-sort-icon="none"
											@if(.Sorted(column) != "")
												hidden
											@endif
										>{!! icons.ArrowsDownUp(icons.Props{}) !!}</span>
										<span
											data-sort-icon="asc"
											@if(.Sorted(column) != "asc")
												hidden
											@endif
										>{!! icons.ArrowUp(icons.Props{}) !!}</span>
										<span
											data-sort-icon="desc"
											@if(.Sorted(column) != "desc")
												hidden
											@endif
										>{!! icons.ArrowDown(icons.Props{}) !!}</span>
									</a>
								@endif
								@if(.Sortable(column) && .HeaderMenu)
									<div class="flex items-center gap-2">
										<details
											id="{{ .ID }}-head-{{ column.Key }}"
											data-part="header-menu"
											data-menu
											data-closes
											class="{{ .PartClass("header-menu", "group/menu relative") }}"
											@attributes(.PartAttrs("header-menu"))
										>
											<summary
												class="{{ "-ms-3 inline-flex cursor-pointer list-none [&::-webkit-details-marker]:hidden group/button shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none h-8 gap-1 px-3 hover:bg-muted hover:text-foreground group-open/menu:bg-accent focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:hover:bg-muted/50" }}"
												aria-label="{{ .SortName(column) }}"
											>
												<span>{{ column.Label }}</span>
												<span
													data-sort-icon="none"
													@if(.Sorted(column) != "")
														hidden
													@endif
												>{!! icons.CaretUpDown(icons.Props{}) !!}</span>
												<span
													data-sort-icon="asc"
													@if(.Sorted(column) != "asc")
														hidden
													@endif
												>{!! icons.ArrowUp(icons.Props{}) !!}</span>
												<span
													data-sort-icon="desc"
													@if(.Sorted(column) != "desc")
														hidden
													@endif
												>{!! icons.ArrowDown(icons.Props{}) !!}</span>
											</summary>
											<div
												data-menu-panel
												data-align="start"
												role="menu"
												class="{{ "fixed z-50 min-w-32 overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none" }}"
											>
												<a
													role="menuitem"
													data-sort-key="{{ column.Key }}"
													data-dir="asc"
													class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4" }}"
													href="{{ .AscHref(column) }}"
													@if(.Complete)
														hx-boost="false"
													@endif
													@if(!.Complete)
														hx-get="{{ .AscHref(column) }}"
														hx-target="{{ .Target() }}"
														hx-swap="outerHTML"
														hx-select="{{ .Target() }}"
														hx-push-url="true"
														hx-sync="closest [data-part=root]:replace"
														hx-indicator="{{ .Target() }}"
													@endif
												>{!! icons.ArrowUp(icons.Props{}) !!}{{ .lines().Asc }}</a>
												<a
													role="menuitem"
													data-sort-key="{{ column.Key }}"
													data-dir="desc"
													class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4" }}"
													href="{{ .DescHref(column) }}"
													@if(.Complete)
														hx-boost="false"
													@endif
													@if(!.Complete)
														hx-get="{{ .DescHref(column) }}"
														hx-target="{{ .Target() }}"
														hx-swap="outerHTML"
														hx-select="{{ .Target() }}"
														hx-push-url="true"
														hx-sync="closest [data-part=root]:replace"
														hx-indicator="{{ .Target() }}"
													@endif
												>{!! icons.ArrowDown(icons.Props{}) !!}{{ .lines().Desc }}</a>
												@if(column.Hideable)
													<div class="-mx-1 my-1 h-px bg-border" role="separator"></div>
													<a
														role="menuitem"
														data-column-hide="{{ column.Key }}"
														class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4" }}"
														href="{{ .HideHref(column) }}"
														@if(.Complete)
															hx-boost="false"
														@endif
														@if(!.Complete)
															hx-get="{{ .HideHref(column) }}"
															hx-target="{{ .Target() }}"
															hx-swap="outerHTML"
															hx-select="{{ .Target() }}"
															hx-push-url="true"
															hx-sync="closest [data-part=root]:replace"
															hx-indicator="{{ .Target() }}"
														@endif
													>{!! icons.EyeSlash(icons.Props{}) !!}{{ .lines().Hide }}</a>
												@endif
											</div>
										</details>
									</div>
								@endif
								@if(!.Sortable(column))
									{{ column.Label }}
								@endif
							</th>
						@endforeach
						@if(.Acting())
							<th
								class="{{ "h-10 px-2 text-start align-middle font-medium whitespace-nowrap text-foreground " + cardHeadOff }}"
								scope="col"
							><span class="sr-only">{{ .lines().Actions }}</span></th>
						@endif
					</tr>
				</thead>
				<tbody class="{{ "[&_tr:last-child]:border-0 " + cardBody }}">
					@foreach(.Body() as row)
						<tr
							data-part="row"
							class="{{ .PartClass("row", "border-b transition-colors hover:bg-muted/50 has-[details[open]]:bg-muted/50 data-[state=selected]:bg-muted", cardRow) }}"
							@if(row.Selected)
								data-state="selected"
							@endif
							@if(row.Outside)
								hidden
							@endif
							@if(row.Key != "")
								data-key="{{ row.Key }}"
							@endif
							@if(row.Level > 0)
								aria-level="{{ row.Level }}"
								data-level="{{ row.Level }}"
							@endif
							@if(row.Branch && row.Expanded)
								aria-expanded="true"
							@endif
							@if(row.Branch && !row.Expanded)
								aria-expanded="false"
							@endif
							@attributes(.RowAttrs(row.TableRow))
						>
							@if(.Selectable())
								<td
									data-part="cell"
									class="{{ .PartClass("cell", "p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pe-0 [&:has([data-select])]:pe-0", cardSelect) }}"
									data-select
									@attributes(.PartAttrs("cell"))
								>
									<input
										data-part="select"
										class="{{ .PartClass("select", "peer relative block size-4 shrink-0 cursor-pointer appearance-none rounded-[4px] border border-input outline-none transition-colors after:absolute after:top-[42%] after:left-1/2 after:hidden after:h-[9px] after:w-[5px] after:-translate-x-1/2 after:-translate-y-1/2 after:rotate-45 after:border-r-2 after:border-b-2 after:border-current after:content-[''] checked:border-primary checked:bg-primary checked:text-primary-foreground checked:after:block focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-input/30 dark:checked:bg-primary") }}"
										type="checkbox"
										value="{{ row.Key }}"
										@if(.Bulky())
											name="{{ .SelectName }}"
											form="{{ .BulkID() }}"
										@endif
										aria-label="{{ .RowName(row.TableRow) }}"
										@if(row.Selected)
											checked
										@endif
										@attributes(.PartAttrs("select"))
									>
								</td>
							@endif
							@for(i := 0; i < len(row.Cells); i++)
								<td
									data-part="cell"
									class="{{ .PartClass("cell", "p-2 align-middle whitespace-nowrap [&:has([role=checkbox])]:pe-0 [&:has([data-select])]:pe-0", .AlignClass(i), .CardCellFor(i, row.Cells[i])) }}"
									@if(.ColumnKey(i) != "")
										data-column="{{ .ColumnKey(i) }}"
									@endif
									@if(.LabelAt(i) != "")
										data-label="{{ .LabelAt(i) }}"
									@endif
									@if(.OffAt(i))
										hidden
									@endif
									@if(row.Cells[i].SortValue != "")
										data-sort-value="{{ row.Cells[i].SortValue }}"
									@endif
									@attributes(.PartAttrs("cell"))
								>
									@if(row.Cells[i].HTML != "")
										{!! row.Cells[i].HTML !!}
									@endif
									@if(row.Cells[i].HTML == "")
										{{ row.Cells[i].Text }}
									@endif
								</td>
							@endfor
							@if(.Acting())
								<td class="{{ "p-2 align-middle whitespace-nowrap " + .CardActions(row.TableRow) }}">
									@if(len(row.Actions) > 0)
										<details
											id="{{ row.MenuID }}"
											data-part="row-menu"
											data-menu
											data-closes
											class="{{ .PartClass("row-menu", "group/menu relative") }}"
											@attributes(.PartAttrs("row-menu"))
										>
											<summary class="{{ "inline-flex cursor-pointer list-none [&::-webkit-details-marker]:hidden group/button shrink-0 items-center justify-center rounded-[min(var(--radius-md),10px)] border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none size-6 hover:bg-muted hover:text-foreground group-open/menu:bg-muted group-open/menu:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:hover:bg-muted/50" }}">
												<span class="sr-only">{{ .lines().OpenMenu }}</span>
												{!! icons.DotsThree(icons.Props{}) !!}
											</summary>
											<div
												data-menu-panel
												data-align="end"
												role="menu"
												class="{{ "fixed z-50 w-44 min-w-32 overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none" }}"
											>
												@for(n := 0; n < len(row.Actions); n++)
													@if(row.Actions[n].Separator)
														<div class="-mx-1 my-1 h-px bg-border" role="separator"></div>
													@elseif(row.Actions[n].Heading)
														<div class="px-1.5 py-1 text-xs font-medium text-muted-foreground" role="heading" aria-level="3">{{ row.Actions[n].Label }}</div>
													@elseif(row.Actions[n].URL != "")
														<a
															role="menuitem"
															class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4" }}"
															href="{{ row.Actions[n].URL }}"
															@if(row.Actions[n].Variant != "")
																data-variant="{{ row.Actions[n].Variant }}"
															@endif
															@if(row.Actions[n].Disabled)
																aria-disabled="true"
															@endif
														>
															{{ row.Actions[n].Label }}
															@if(row.Actions[n].Shortcut != "")
																<span class="ml-auto text-xs tracking-widest text-muted-foreground">{{ row.Actions[n].Shortcut }}</span>
															@endif
														</a>
													@else
														<button
															type="button"
															role="menuitem"
															class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md px-1.5 py-1 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[variant=destructive]:text-destructive [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4" }}"
															@if(row.Actions[n].Variant != "")
																data-variant="{{ row.Actions[n].Variant }}"
															@endif
															@if(row.Actions[n].Disabled)
																aria-disabled="true"
															@endif
															@if(row.Actions[n].HxPost != "")
																hx-post="{{ row.Actions[n].HxPost }}"
															@endif
															@if(row.Actions[n].HxGet != "")
																hx-get="{{ row.Actions[n].HxGet }}"
															@endif
															@if(row.Actions[n].HxTarget != "")
																hx-target="{{ row.Actions[n].HxTarget }}"
															@endif
															@if(row.Actions[n].HxConfirm != "")
																hx-confirm="{{ row.Actions[n].HxConfirm }}"
															@endif
														>
															{{ row.Actions[n].Label }}
															@if(row.Actions[n].Shortcut != "")
																<span class="ml-auto text-xs tracking-widest text-muted-foreground">{{ row.Actions[n].Shortcut }}</span>
															@endif
														</button>
													@endif
												@endfor
											</div>
										</details>
									@endif
								</td>
							@endif
						</tr>
					@endforeach
					@if(len(.Rows) == 0)
						<tr class="{{ "border-b transition-colors hover:bg-muted/50 " + cardFull }}">
							<td colspan="{{ .Span() }}" class="{{ "h-24 p-2 text-center " + cardFull }}">{{ .lines().NoResult }}</td>
						</tr>
					@endif
					@if(.Complete && len(.Rows) > 0)
						{{-- Drawn hidden, and shown by the browser when a search or a
						     facet leaves no row: a script that builds a styled row
						     would carry classes the stylesheet was never compiled
						     with. --}}
						<tr data-empty-row hidden class="{{ "border-b transition-colors hover:bg-muted/50 " + cardFull }}">
							<td colspan="{{ .Span() }}" class="{{ "h-24 p-2 text-center " + cardFull }}">{{ .lines().NoResult }}</td>
						</tr>
					@endif
				</tbody>
			</table>
		</div>
	</div>

	{{-- The footer: how many rows are selected of the whole list on the left,
	     and on the right the page size, where the reader is, and the steps. The
	     count is a live region so a page change is announced. --}}
	@if(.FooterText() != "" || .Stepped() || .Chooses())
		<div
			data-part="footer"
			class="{{ .PartClass("footer", "flex flex-wrap items-center justify-end gap-x-2 gap-y-2 py-4") }}"
			@attributes(.PartAttrs("footer"))
		>
			<div
				data-part="count"
				class="{{ .PartClass("count", "min-w-[10rem] flex-1 text-sm text-muted-foreground") }}"
				role="status"
				aria-live="polite"
				data-selected-template="{{ .SelectedSentence() }}"
				@attributes(.PartAttrs("count"))
			>{{ .FooterText() }}</div>

			<div class="flex flex-wrap items-center gap-x-6 gap-y-2 lg:gap-x-8">
				@if(.Chooses())
					{{-- The size is a menu of addresses, not a field of the form: the
					     state decides which page keeps the first row on screen, and
					     only a link built from the state can say so. --}}
					<div class="flex items-center space-x-2">
						<p id="{{ .SizesID() }}-label" class="text-sm font-medium">{{ .lines().RowsPerPage }}</p>
						<details
							id="{{ .SizesID() }}"
							data-part="sizes"
							data-menu
							data-closes
							class="{{ .PartClass("sizes", "group/menu relative") }}"
							@attributes(.PartAttrs("sizes"))
						>
							<summary
								class="{{ "inline-flex cursor-pointer list-none [&::-webkit-details-marker]:hidden h-8 w-[70px] items-center justify-between gap-1 rounded-md border border-input bg-transparent px-2 py-1 text-sm whitespace-nowrap shadow-xs outline-none select-none focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 group-open/menu:bg-muted [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 [&_svg]:opacity-50 dark:bg-input/30" }}"
								aria-labelledby="{{ .SizesID() }}-label {{ .SizesID() }}-value"
							>
								<span id="{{ .SizesID() }}-value" data-size-value>{{ .Size() }}</span>
								{!! icons.CaretUpDown(icons.Props{}) !!}
							</summary>
							<div
								data-menu-panel
								data-align="start"
								role="menu"
								class="{{ "fixed z-50 min-w-[70px] overflow-x-hidden overflow-y-auto rounded-lg bg-popover p-1 text-popover-foreground shadow-md ring-1 ring-foreground/10 outline-none" }}"
							>
								@foreach(.Sizes() as size)
									<a
										role="menuitemradio"
										aria-checked="{{ size == .Size() }}"
										data-size="{{ size }}"
										class="{{ "relative flex w-full cursor-default items-center gap-1.5 rounded-md py-1 pe-8 ps-1.5 text-sm outline-hidden select-none hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground" }}"
										href="{{ .SizeHref(size) }}"
										@if(.Complete)
											hx-boost="false"
										@endif
										@if(!.Complete)
											hx-get="{{ .SizeHref(size) }}"
											hx-target="{{ .Target() }}"
											hx-swap="outerHTML"
											hx-select="{{ .Target() }}"
											hx-push-url="true"
											hx-sync="closest [data-part=root]:replace"
											hx-indicator="{{ .Target() }}"
										@endif
									>
										{{ size }}
										<span
											class="pointer-events-none absolute end-2 flex items-center justify-center [&_svg]:size-4 [&_path[opacity]]:hidden"
											data-indicator
											@if(size != .Size())
												hidden
											@endif
										>{!! icons.Check(icons.Props{}) !!}</span>
									</a>
								@endforeach
							</div>
						</details>
					</div>
				@endif
				@if(.Stepped())
					@if(.FullPager && .PageOfText() != "")
						<div
							data-part="page-of"
							class="{{ .PartClass("page-of", "flex w-[100px] items-center justify-center text-sm font-medium") }}"
							data-page-of-template="{{ .lines().PageOf }}"
							@attributes(.PartAttrs("page-of"))
						>{{ .PageOfText() }}</div>
					@endif
					<nav
						data-part="pagination"
						class="{{ .PartClass("pagination", "flex items-center space-x-2") }}"
						aria-label="{{ .PagerName() }}"
						@if(!.Paged() && !.Complete)
							hidden
						@endif
						@attributes(.PartAttrs("pagination"))
					>
						@if(.FullPager && .Known())
							<a
								data-part="first"
								data-step="first"
								class="{{ .PartClass("first", "hidden lg:flex group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none size-8 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
								aria-label="{{ .lines().First }}"
								@if(.HasPrevious() || .Complete)
									href="{{ .FirstHref() }}"
								@endif
								@if(!.HasPrevious())
									aria-disabled="true"
									tabindex="-1"
								@endif
								@if(.Complete)
									hx-boost="false"
								@endif
								@if(!.Complete && .HasPrevious())
									hx-get="{{ .FirstHref() }}"
									hx-target="{{ .Target() }}"
									hx-swap="outerHTML"
									hx-select="{{ .Target() }}"
									hx-push-url="true"
									hx-sync="closest [data-part=root]:replace"
									hx-indicator="{{ .Target() }}"
								@endif
								@attributes(.PartAttrs("first"))
							>{!! icons.CaretDoubleLeft(icons.Props{}) !!}</a>
						@endif
						<a
							data-part="previous"
							data-step="previous"
							@if(.FullPager)
								class="{{ .PartClass("previous", "inline-flex group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none size-8 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
							@endif
							@if(!.FullPager)
								class="{{ .PartClass("previous", "inline-flex group/button shrink-0 items-center justify-center rounded-[min(var(--radius-md),12px)] border border-border bg-background bg-clip-padding text-[0.8rem] font-medium whitespace-nowrap transition-all outline-none select-none h-7 gap-1 px-2.5 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-3.5 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
							@endif
							aria-label="{{ .lines().Previous }}"
							@if(.HasPrevious() || .Complete)
								href="{{ .PreviousHref() }}"
							@endif
							@if(!.HasPrevious())
								aria-disabled="true"
								tabindex="-1"
							@endif
							@if(.Complete)
								hx-boost="false"
							@endif
							@if(!.Complete && .HasPrevious())
								hx-get="{{ .PreviousHref() }}"
								hx-target="{{ .Target() }}"
								hx-swap="outerHTML"
								hx-select="{{ .Target() }}"
								hx-push-url="true"
								hx-sync="closest [data-part=root]:replace"
								hx-indicator="{{ .Target() }}"
							@endif
							@attributes(.PartAttrs("previous"))
						>
							@if(.FullPager)
								{!! icons.CaretLeft(icons.Props{}) !!}
							@endif
							@if(!.FullPager)
								{{ .lines().Back }}
							@endif
						</a>
						<a
							data-part="next"
							data-step="next"
							@if(.FullPager)
								class="{{ .PartClass("next", "inline-flex group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none size-8 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
							@endif
							@if(!.FullPager)
								class="{{ .PartClass("next", "inline-flex group/button shrink-0 items-center justify-center rounded-[min(var(--radius-md),12px)] border border-border bg-background bg-clip-padding text-[0.8rem] font-medium whitespace-nowrap transition-all outline-none select-none h-7 gap-1 px-2.5 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-3.5 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
							@endif
							aria-label="{{ .lines().Next }}"
							@if(.HasNext() || .Complete)
								href="{{ .NextHref() }}"
							@endif
							@if(!.HasNext())
								aria-disabled="true"
								tabindex="-1"
							@endif
							@if(.Complete)
								hx-boost="false"
							@endif
							@if(!.Complete && .HasNext())
								hx-get="{{ .NextHref() }}"
								hx-target="{{ .Target() }}"
								hx-swap="outerHTML"
								hx-select="{{ .Target() }}"
								hx-push-url="true"
								hx-sync="closest [data-part=root]:replace"
								hx-indicator="{{ .Target() }}"
							@endif
							@attributes(.PartAttrs("next"))
						>
							@if(.FullPager)
								{!! icons.CaretRight(icons.Props{}) !!}
							@endif
							@if(!.FullPager)
								{{ .lines().Forward }}
							@endif
						</a>
						@if(.FullPager && .Known())
							<a
								data-part="last"
								data-step="last"
								class="{{ .PartClass("last", "hidden lg:flex group/button shrink-0 items-center justify-center rounded-lg border border-border bg-background bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none size-8 hover:bg-muted hover:text-foreground focus-visible:border-ring focus-visible:ring-3 focus-visible:ring-ring/50 aria-disabled:pointer-events-none aria-disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0 [&_svg:not([class*='size-'])]:size-4 dark:border-input dark:bg-input/30 dark:hover:bg-input/50") }}"
								aria-label="{{ .lines().Last }}"
								@if(.HasNext() || .Complete)
									href="{{ .LastHref() }}"
								@endif
								@if(!.HasNext())
									aria-disabled="true"
									tabindex="-1"
								@endif
								@if(.Complete)
									hx-boost="false"
								@endif
								@if(!.Complete && .HasNext())
									hx-get="{{ .LastHref() }}"
									hx-target="{{ .Target() }}"
									hx-swap="outerHTML"
									hx-select="{{ .Target() }}"
									hx-push-url="true"
									hx-sync="closest [data-part=root]:replace"
									hx-indicator="{{ .Target() }}"
								@endif
								@attributes(.PartAttrs("last"))
							>{!! icons.CaretDoubleRight(icons.Props{}) !!}</a>
						@endif
					</nav>
				@endif
			</div>
		</div>
	@endif
	@endif
</div>
