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
	//
	// A query string or fragment on it is dropped, never carried: every
	// address and every hidden field this draws is built from the fields of
	// this struct alone. URL is often the request's own address, and a query
	// carried from it would put whatever parameters somebody wrote into a link
	// on every control of the page. A parameter the list needs is one of these
	// fields -- a Filter for a facet -- or part of the path.
	URL string

	// Caption is the sentence saying what the list is, drawn under the table
	// and used as its name by assistive technology. Empty draws none, and the
	// region's own Label is then the only thing naming it.
	Caption string
	// Columns are the headers. A column with a Key can be ordered by and
	// switched off; see TableColumn.
	Columns []TableColumn
	// Rows are the lines of the page being shown, or the whole list when
	// Complete says so.
	Rows []TableRow
	// Complete says the rows above are the whole list and not a page of it.
	//
	// It is a fact about the data and not a preference, and everything else
	// follows from it. When the page holds every row, the browser can search,
	// order and page correctly because it has everything -- and it does, so a
	// click costs nothing and the server is not asked. When the page holds a
	// window onto something longer, the browser cannot: ordering the rows in
	// hand orders only those, and filtering them hides matches that were never
	// sent. Then every control is an address and the server answers.
	//
	// Saying true about a page of a longer list is the bug this field is named
	// to prevent: it reads as rows in the wrong order to everyone except
	// whoever wrote it, and nothing crashes.
	//
	// A complete list still keeps its addresses. The headers are still links
	// and the search is still a form, so the page works with no script -- the
	// behaviour takes the click before the link is followed. That is what
	// makes this an enhancement rather than a second implementation.
	Complete bool
	// PageSize is how many rows a browser-paged list shows at once. Zero shows
	// them all, which is right for anything short enough not to need a pager.
	// It means nothing when the server is paging, because then the page it
	// sent is the page.
	PageSize int
	// PageSizes are the counts the reader can switch between. Empty draws no
	// chooser, which is right for a list whose length nobody argues with.
	//
	// It is a list and not a range because the numbers are a judgement about
	// the rows: ten of something wide, a hundred of something narrow.
	PageSizes []int
	// Lines are the sentences this draws, so none of them is written in a
	// script and all of them come from the application's own catalogue.
	Lines DataTableLines
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

	// Navigable turns the table into a grid: one tab stop, and the arrow keys
	// moving from cell to cell. See TableProps.Navigable.
	Navigable bool
	// BulkMethod is how the bulk bar submits. See TableProps.BulkMethod.
	BulkMethod string
	// Filters are the per-column facets, drawn beside the search. Their chosen
	// values ride on every sort and page link, which is what keeps a facet
	// from being thrown away by the next click.
	Filters []DataTableFilter
	// SortKey and SortDir are the ordering in force, and travel as "sort" and
	// "dir".
	SortKey string
	SortDir string

	// Paginator is the page the query returned, for a windowed list: the
	// *pagination.LengthAwarePage, *pagination.Page or *pagination.CursorPage
	// that came back beside the rows. Set, it is where the table reads the
	// page being shown, how many there are, how many rows the whole list has
	// and the range this page holds, and Page, Pages and Total are not read.
	//
	// The addresses stay the table's own, built from its fields like every
	// other control, with the page's parameter name -- and, for a cursor page,
	// its signed cursor. A page that knows no total draws the two steps and
	// neither the range nor "Page n of total", because neither can be true.
	// A complete list has no use for one: there the rows are the list.
	Paginator Paginator

	// Page is the page being shown, counting from one, and Pages is how many
	// there are. It travels as "page". Pages of zero or one draws no pager.
	//
	// Deprecated: set Paginator to the page the query returned. Page, Pages
	// and Total are that page copied out into three fields.
	Page  int
	Pages int
	// Total is how many rows the whole list has, which on a windowed list only
	// the server knows -- it is the count it already ran to work out Pages.
	//
	// Without it the range under the table cannot be true, so it is not drawn:
	// Rows is one page, and saying "1 to 20 of 20" under page three of a
	// hundred and thirty-seven is worse than saying nothing. A complete list
	// needs none of this, because there the rows are the list.
	//
	// Deprecated: set Paginator.
	Total int

	// SelectName turns on per-row selection and is the field the checkboxes
	// submit under. BulkActions are the controls revealed once something is
	// chosen.
	SelectName  string
	BulkActions []ButtonProps
	// BulkAction is where the bar posts. Empty posts to URL, which is the
	// endpoint that already knows the list.
	BulkAction string
	// PageData is the screen's own view.Page. The bulk form posts the CSRF
	// token it carries, so a list drawn from a view writes PageData: .Page and
	// never copies the token by hand. It is not called Page because Page is
	// the number of the page being shown.
	PageData Page
	// Token is the CSRF token, for a list with no PageData to ask. Set, it
	// wins over the page's.
	Token string

	// ColumnsLabel and SelectedLabel are the two sentences that belong to the
	// table rather than to the list. Everything else this says is in Lines --
	// two fields for one control is exactly the defect that puts a table into
	// two languages.
	ColumnsLabel  string
	SelectedLabel string
}

// DataTableLines are the sentences a data table says.
//
// They are a struct and not a scatter of Label fields because they belong
// together: a table says all of them or none, and a project that translated
// four of the seven has a table that speaks two languages. Every one carries
// its placeholders, so a language that puts the number last can.
//
// Empty says the English underneath. That is a floor and not a default worth
// shipping -- a table drawn from a catalogue says these in the language the
// page is served in, which is the only place that knows it.
type DataTableLines struct {
	// Showing is the range under the table: "{from} to {to} of {total}".
	Showing string
	// PerPage labels the chooser beside it: "per page".
	PerPage string
	// NoResult is what the body says when the search matched nothing. It is
	// not the same sentence as an empty list: the table has rows, and none of
	// them is the answer.
	NoResult string
	// Previous and Next name the two ends of the pager for a screen reader.
	// They are names and not text, because the controls are arrows.
	Previous string
	Next     string
	// Page names one numbered control: "Page {n}". The digit alone is read as
	// a digit, which in a row of them says nothing about what it does.
	Page string
	// PageOf is where the reader is, beside the pager: "Page {n} of {total}".
	PageOf string
	// Back and Forward are the words on the pager's two buttons: "Previous"
	// and "Next". Previous and Next above are what a screen reader calls them.
	Back    string
	Forward string
	// SelectAll names the header checkbox, and Bulk the bar of actions.
	SelectAll string
	Bulk      string
	// Apply names the submit control that makes server-side facets work when
	// JavaScript is unavailable. Empty uses "Apply filters".
	Apply string
	// Ascending, Descending and Unsorted are what a sortable header is called
	// in each of its three states: "{column}, sorted ascending". A header that
	// only draws an arrow is a header a screen reader reads as a link.
	Ascending  string
	Descending string
	Unsorted   string
}

// Sentences are the lines with the English floor filled in under whatever the
// caller left empty.
//
// Filling them here rather than in the script is the whole point: the script
// composes no sentence, so there is one place a translation lands and it is a
// file the application already has.
func (l DataTableLines) Sentences() map[string]any {
	fallback := func(given, floor string) string {
		if given != "" {
			return given
		}
		return floor
	}
	return map[string]any{
		"showing":    fallback(l.Showing, "{from} to {to} of {total}"),
		"perPage":    fallback(l.PerPage, "per page"),
		"noResult":   fallback(l.NoResult, "Nothing matched"),
		"previous":   fallback(l.Previous, "Previous page"),
		"next":       fallback(l.Next, "Next page"),
		"page":       fallback(l.Page, "Page {n}"),
		"pageOf":     fallback(l.PageOf, "Page {n} of {total}"),
		"back":       fallback(l.Back, "Previous"),
		"forward":    fallback(l.Forward, "Next"),
		"ascending":  fallback(l.Ascending, "{column}, sorted ascending"),
		"descending": fallback(l.Descending, "{column}, sorted descending"),
		"unsorted":   fallback(l.Unsorted, "{column}, not sorted"),
		"selectAll":  fallback(l.SelectAll, "Select all rows"),
		"bulk":       fallback(l.Bulk, "Bulk actions"),
		"apply":      fallback(l.Apply, "Apply filters"),
	}
}

// Address is URL with a query string built from everything except what the
// caller is changing.
//
// Every control goes through here, which is what stops the four from throwing
// each other away: a sort link carries the search, a page link carries both,
// and a new search drops the page rather than keeping a number that no longer
// points anywhere.
//
// # Two things it has to do that are easy to miss
//
// It drops whatever query URL already carries and builds the new one from this
// struct alone. URL is often the address the request arrived on, so a query
// seeded from it would copy any parameter somebody put in a link -- _method,
// a role, a tenant -- onto every sort header, every page and the search form.
// The state a list needs is its fields; a parameter that is not one of them
// does not travel. Parsing also means one "?" rather than a second one
// appended to a URL that had one, which every parser reads as one parameter
// whose value is the rest of the query.
//
// And it writes the templated pairs by hand, after the encoding. A value of
// "{sort}" through url.Values.Encode comes out "%7Bsort%7D", and the caller
// looking for "{sort}" then finds nothing -- every header and every page
// pointing at one dead address. That was the bug this comment exists to keep
// from coming back.
func (p DataTableProps) Address(changes map[string]string) string {
	base, err := url.Parse(p.URL)
	if err != nil {
		base = &url.URL{Path: p.URL}
	}
	values := url.Values{}

	set := func(key, value string) {
		if value == "" {
			values.Del(key)
			return
		}
		values.Set(key, value)
	}
	if p.SearchName != "" {
		set(p.SearchName, p.Query)
	}
	set("sort", p.SortKey)
	if p.SortKey != "" {
		set("dir", p.SortDir)
	}
	if at := p.at(); at > 1 {
		set(p.pageParam(), strconv.Itoa(at))
	}
	for _, filter := range p.Filters {
		values.Del(filter.Key)
		for _, chosen := range filter.Selected {
			values.Add(filter.Key, chosen)
		}
	}

	// A templated value is held out of the encoder and appended verbatim.
	templated := make([][2]string, 0, len(changes))
	for key, value := range changes {
		if value == "" {
			values.Del(key)
			continue
		}
		if strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}") {
			values.Del(key)
			templated = append(templated, [2]string{key, value})
			continue
		}
		values.Set(key, value)
	}
	sort.Slice(templated, func(a, b int) bool { return templated[a][0] < templated[b][0] })

	query := values.Encode()
	for _, pair := range templated {
		if query != "" {
			query += "&"
		}
		query += pair[0] + "=" + pair[1]
	}

	base.RawQuery = ""
	base.ForceQuery = false
	base.Fragment = ""
	address := base.String()
	if query == "" {
		return address
	}
	return address + "?" + query
}

// SortURL is the template the table builds its header links from. The page is
// dropped: an order that kept it would land somebody on page four of a list
// they have just reordered, where the rows are not the ones they were looking
// at.
func (p DataTableProps) SortURL() string {
	return p.Address(map[string]string{"sort": "{sort}", "dir": "{dir}", p.pageParam(): ""})
}

// PageURL is the template the pager builds from. It keeps the search, the
// order and every facet, because a page of a filtered, sorted list is a page
// of that list.
func (p DataTableProps) PageURL() string {
	return p.Address(map[string]string{p.pageParam(): "{page}"})
}

// DataTableFilter is one per-column facet, drawn beside the search.
//
// It is typed rather than a slot of markup, because a slot would be a place to
// put anything, and this component's whole claim is that every control's state
// survives every other control -- which it can only promise about state it
// knows the shape of.
type DataTableFilter struct {
	// Key is the parameter the chosen values travel under.
	Key string
	// Label names the control.
	Label string
	// Options are what can be chosen.
	Options []SelectOption
	// Selected are the values in force, which every link then carries.
	Selected []string
}

// Filtered is whether any facet is drawn.
func (p DataTableProps) Filtered() bool { return len(p.Filters) > 0 }

// Chosen is whether one option of a facet is in force.
func (p DataTableProps) Chosen(filter DataTableFilter, option SelectOption) bool {
	for _, value := range filter.Selected {
		if value == option.Value {
			return true
		}
	}
	return false
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

// Sheets is how many pages there are.
//
// A complete list counts its own: the server sent every row, so how many
// pages that makes is arithmetic and not something a caller should have to
// pass in beside the rows it already passed. A windowed list is told, because
// only the server knows how much it did not send.
func (p DataTableProps) Sheets() int {
	if p.Complete && p.PageSize > 0 {
		if len(p.Rows) == 0 {
			return 1
		}
		return (len(p.Rows) + p.PageSize - 1) / p.PageSize
	}
	if page, ok := p.countedPage(); ok {
		return page.LastPage()
	}
	return p.Pages
}

// Paged is whether there is more than one page to move between.
func (p DataTableProps) Paged() bool {
	if paginator := p.uncountedPage(); paginator != nil {
		return paginator.HasPages()
	}
	return p.Sheets() > 1
}

// countedPage is the paginator of a windowed list when it knows its total.
func (p DataTableProps) countedPage() (countedPage, bool) {
	if p.Complete {
		return nil, false
	}
	return counted(p.Paginator)
}

// uncountedPage is the paginator of a windowed list when it knows no total --
// a simple page or a cursor page -- and nil otherwise.
func (p DataTableProps) uncountedPage() Paginator {
	paginator := presentPaginator(p.Paginator)
	if p.Complete || paginator == nil {
		return nil
	}
	if _, ok := paginator.(countedPage); ok {
		return nil
	}
	return paginator
}

// at is the page being shown: the paginator's number when it has one, Page
// otherwise.
func (p DataTableProps) at() int {
	if !p.Complete {
		if page, ok := presentPaginator(p.Paginator).(numberedPage); ok {
			return page.CurrentPage()
		}
	}
	return p.Page
}

// pageParam is the query parameter the page number travels under: the
// paginator's, so the address the table writes is the one the query reads.
func (p DataTableProps) pageParam() string {
	if page, ok := presentPaginator(p.Paginator).(numberedPage); ok && page.GetPageName() != "" {
		return page.GetPageName()
	}
	return "page"
}

// stepAddresses are the two steps of a page that knows no total, each built
// from the table's own state like every other control, and empty where there
// is nowhere to go.
func (p DataTableProps) stepAddresses(paginator Paginator) steps {
	var out steps
	switch page := paginator.(type) {
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
			return p.Address(map[string]string{page.GetCursorName(): token, p.pageParam(): ""})
		}
		out.previous = cursor(page.PreviousCursor())
		out.next = cursor(page.NextCursor())
	case numberedPage:
		at := page.CurrentPage()
		if at > 1 {
			out.previous = p.Address(map[string]string{p.pageParam(): strconv.Itoa(at - 1)})
		}
		if page.HasMorePages() {
			out.next = p.Address(map[string]string{p.pageParam(): strconv.Itoa(at + 1)})
		}
	default:
		// A page of a kind this does not know the parameters of has only its
		// own addresses to offer.
		out.previous = paginator.PreviousPageURL()
		out.next = paginator.NextPageURL()
	}
	return out
}

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
	rows := p.Rows
	if p.Complete && p.PageSize > 0 {
		// The window is marked here so the first paint is already the page.
		from, to, _ := p.window()
		rows = make([]TableRow, 0, len(p.Rows))
		for at, row := range p.Rows {
			row.Hidden = at+1 < from || at+1 > to
			rows = append(rows, row)
		}
	}

	return TableProps{
		ComponentProps:  ComponentProps{Parts: p.Parts},
		Composed:        true,
		ID:              p.ID,
		Navigable:       p.Navigable,
		BulkMethod:      p.BulkMethod,
		SelectAllLabel:  p.Lines.SelectAll,
		BulkLabel:       p.Lines.Bulk,
		AscendingLabel:  p.Lines.Ascending,
		DescendingLabel: p.Lines.Descending,
		UnsortedLabel:   p.Lines.Unsorted,
		Caption:         p.Caption,
		Columns:         columns,
		Rows:            rows,
		Empty:           p.Empty,
		SortKey:        p.SortKey,
		SortDir:        p.SortDir,
		SortURL:        p.SortURL(),
		SelectName:     p.SelectName,
		SelectedLabel:  p.SelectedLabel,
		BulkActions:    p.BulkActions,
		BulkAction:     p.BulkAction,
		Token:          p.Token,
		Page:           p.PageData,
		HxTarget:       p.swapTarget(),
		HxSwap:         p.swapKind(),
	}
}

// swapTarget and swapKind are empty on a complete list, and that is the whole
// of what keeps one gesture from running two engines.
//
// preventDefault does not cancel HTMX: its listener sits on the link itself,
// so a header that carries hx-get fires a request the browser was told not to
// follow, and the outerHTML swap then hands back the server's region -- wiping
// the page, the size, the filter, the hidden columns and the selection the
// browser had just worked out. href, action and method stay either way, so the
// path with no script is unchanged.
func (p DataTableProps) swapTarget() string {
	if p.Complete {
		return ""
	}
	return p.Target()
}

// swapKind is how the region is replaced, and nothing when the browser is the
// one working the rows.
func (p DataTableProps) swapKind() string {
	if p.Complete {
		return ""
	}
	return "outerHTML"
}

// Pager is the pagination this draws, or the zero value when there is one page.
func (p DataTableProps) Pager() PaginationProps {
	pager := PaginationProps{
		ComponentProps: ComponentProps{Parts: p.Parts},
		Page:           p.at(),
		Pages:          p.Sheets(),
		URL:            p.PageURL(),
		Label:          p.Label,
		PreviousLabel:  p.Lines.Previous,
		NextLabel:      p.Lines.Next,
		PageLabel:      p.Lines.Page,
		HxTarget:       p.swapTarget(),
		HxSwap:         p.swapKind(),
	}
	if paginator := p.uncountedPage(); paginator != nil {
		// A page that knows no total has no numbers to template, and is drawn
		// as its two steps.
		pager.Paginator = p.stepAddresses(paginator)
	}
	return pager
}

// ColumnsName is what the menu of columns is called.
func (p DataTableProps) ColumnsName() string {
	if p.ColumnsLabel != "" {
		return p.ColumnsLabel
	}
	return "Columns"
}

// DataTableQueryField is one hidden field that preserves query state across
// a native GET form submission.
type DataTableQueryField struct {
	Name  string
	Value string
}

// QueryFormID is the one form search and facets submit through.
func (p DataTableProps) QueryFormID() string { return p.ID + "-query" }

// QueryAction is URL without its query. The state that belongs to the
// DataTable is submitted as named fields so a native GET form cannot discard it.
func (p DataTableProps) QueryAction() string {
	base, err := url.Parse(p.URL)
	if err != nil {
		return p.URL
	}
	base.RawQuery = ""
	base.ForceQuery = false
	base.Fragment = ""
	return base.String()
}


// QueryFields are the hidden fields of the query form: the order in force,
// which has no visible control of its own there. Search, page and facet values
// are submitted by their own controls.
//
// Nothing of URL's own query is carried, for the reason given on URL: a hidden
// field copied from the request's address is a parameter somebody else chose,
// submitted again by whoever next searches.
func (p DataTableProps) QueryFields() []DataTableQueryField {
	if p.SortKey == "" {
		return nil
	}
	fields := make([]DataTableQueryField, 0, 2)
	if p.SortDir != "" {
		fields = append(fields, DataTableQueryField{Name: "dir", Value: p.SortDir})
	}
	return append(fields, DataTableQueryField{Name: "sort", Value: p.SortKey})
}

// ApplyText is the no-script facet submit label.
func (p DataTableProps) ApplyText() string {
	if p.Lines.Apply != "" {
		return p.Lines.Apply
	}
	return "Apply filters"
}

// SearchTrigger is when the list is refetched as somebody types: after typing
// has paused, and only when what is in the box actually changed.
func (p DataTableProps) SearchTrigger() string {
	return "input changed delay:300ms, search"
}

// SearchAddress is where the box submits: the endpoint with the order kept and
// the page dropped, because a new search has no page four.
func (p DataTableProps) SearchAddress() string {
	return p.Address(map[string]string{p.SearchName: "", p.pageParam(): ""})
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
//
// It is the table's behaviour, because what needs mounting is the table's:
// the count, the column toggles and the roving cell. Nothing about the
// searching, the ordering or the paging needs a script -- they are addresses.
func (p DataTableProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		props := p.Lines.Sentences()
		props["selected"] = p.SelectedLabel
		props["complete"] = p.Complete
		props["pageSize"] = p.PageSize
		p.Behavior = Behavior{Name: TableBehavior, Props: props}
	}
	return p.ComponentProps.RootAttrs()
}

// Sizes are the counts the chooser offers.
func (p DataTableProps) Sizes() []int { return p.PageSizes }

// Chooses is whether the chooser is drawn: only on a complete list, because
// only there can the count change without asking the server. A windowed list
// is paged by the endpoint, and the size is the endpoint's to decide.
func (p DataTableProps) Chooses() bool { return p.Complete && len(p.PageSizes) > 0 }

// PerPageText is the label beside the chooser.
func (p DataTableProps) PerPageText() string {
	if p.Lines.PerPage != "" {
		return p.Lines.PerPage
	}
	return "per page"
}

// Showing is the range under the table, drawn by the server for the first
// page and rewritten by the behaviour after that -- from the same sentence,
// so the two cannot say it differently.
func (p DataTableProps) Showing() string {
	sentence := p.Lines.Showing
	if sentence == "" {
		sentence = "{from} to {to} of {total}"
	}
	from, to, total := p.window()
	sentence = strings.ReplaceAll(sentence, "{from}", strconv.Itoa(from))
	sentence = strings.ReplaceAll(sentence, "{to}", strconv.Itoa(to))
	return strings.ReplaceAll(sentence, "{total}", strconv.Itoa(total))
}

// window is the range the table is showing and the size of the list it is a
// window onto.
func (p DataTableProps) window() (int, int, int) {
	if page, ok := p.countedPage(); ok {
		return page.FirstItem(), page.LastItem(), page.Total()
	}
	total := p.Total
	if p.Complete {
		total = len(p.Rows)
	}
	if total == 0 {
		return 0, 0, 0
	}
	if !p.Complete && len(p.Rows) == 0 {
		return 0, 0, total
	}
	at := p.at()
	if at < 1 {
		at = 1
	}
	if p.PageSize > 0 {
		from := (at-1)*p.PageSize + 1
		to := from + p.PageSize - 1
		if to > total {
			to = total
		}
		return from, to, total
	}
	if p.Complete {
		return 1, total, total
	}
	// A windowed list with no page size still knows where it is from the rows
	// it was handed.
	from := (at-1)*len(p.Rows) + 1
	return from, from + len(p.Rows) - 1, total
}

// Ranged is whether the range can be stated truthfully. A windowed list that
// was not told the total cannot, and draws nothing rather than a number that
// is silently the page size.
func (p DataTableProps) Ranged() bool {
	if page, ok := p.countedPage(); ok {
		return page.Total() > 0
	}
	if p.uncountedPage() != nil {
		return false
	}
	return p.Complete || p.Total > 0
}

// PartNames are the parts this component publishes.
// PageOfText is where the reader is, beside the pager: "Page 2 of 7".
// It is empty for a page that knows no total, because "Page 3 of" with nothing
// after it says less than nothing.
func (p DataTableProps) PageOfText() string {
	if p.uncountedPage() != nil {
		return ""
	}
	page := p.at()
	if page < 1 {
		page = 1
	}
	return strings.NewReplacer("{n}", strconv.Itoa(page), "{total}", strconv.Itoa(p.Sheets())).
		Replace(p.Lines.Sentences()["pageOf"].(string))
}

// PageOfTemplate is the sentence the browser rewrites as the page moves.
func (p DataTableProps) PageOfTemplate() string { return p.Lines.Sentences()["pageOf"].(string) }

// PreviousPage and NextPage are the pages the two steps lead to, kept inside
// the list.
func (p DataTableProps) PreviousPage() int { return max(1, p.at()-1) }
func (p DataTableProps) NextPage() int     { return min(max(1, p.Sheets()), max(1, p.at())+1) }

// PreviousText and NextText are the words on the two steps.
func (p DataTableProps) PreviousText() string { return p.Lines.Sentences()["back"].(string) }
func (p DataTableProps) NextText() string     { return p.Lines.Sentences()["forward"].(string) }

// Selectable is whether the rows carry a checkbox, and the footer a count.
func (p DataTableProps) Selectable() bool { return p.SelectName != "" }

// SelectedText is the count of rows ticked, as the table's own count says it.
func (p DataTableProps) SelectedText() string { return p.Grid().SelectedText() }

// FacetChosen are the labels of a facet's options in force, in the order the
// options are listed, which is what the facet's button shows.
func (p DataTableProps) FacetChosen(filter DataTableFilter) []string {
	out := []string{}
	for _, option := range filter.Options {
		if p.Chosen(filter, option) {
			out = append(out, option.Label)
		}
	}
	return out
}

// FacetSummary is what the facet's button shows of its choice: every label
// while there are two at most, and the count past that, so a facet with five
// choices does not push the toolbar off the screen.
func (p DataTableProps) FacetSummary(filter DataTableFilter) []string {
	chosen := p.FacetChosen(filter)
	if len(chosen) > 2 {
		return []string{strconv.Itoa(len(chosen))}
	}
	return chosen
}

// FacetName names one facet's menu for a screen reader: its label.
func (p DataTableProps) FacetName(filter DataTableFilter) string { return filter.Label }

func (p DataTableProps) PartNames() []string {
	// The table's names and the pager's are published here as well, because
	// this draws both and the parts a caller can reach are the parts on the
	// page -- not the parts this file happens to write itself.
	names := []string{
		"root", "toolbar", "search", "indicator", "columns", "toggle",
		"pagination", "sizes", "showing", "footer", "filter",
		"facet", "facet-chosen", "page-of", "previous", "next",
	}
	names = append(names, TableProps{}.PartNames()[1:]...)

	// Deduplicated, because this composes two components that publish some of
	// the same names and a list that repeats one is a list nobody can trust to
	// be the set.
	seen := map[string]bool{}
	unique := names[:0]
	for _, name := range names {
		if seen[name] {
			continue
		}
		seen[name] = true
		unique = append(unique, name)
	}
	return unique
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
	@if(.Complete)
		data-complete="true"
	@endif
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(.Searchable() || .Hideable() || .Chooses())
		<div
			data-part="toolbar"
			class="{{ .PartClass("toolbar", "data-table-toolbar") }}"
			@attributes(.PartAttrs("toolbar"))
		>
			@if(.Searchable() || .Filtered())
				{{-- Search and facets are one native GET form. Keeping them in one
				     form is what makes the no-script path preserve the same state
				     HTMX sends, instead of turning each control into a separate
				     implementation. --}}
				<form
					class="data-table-controls"
					id="{{ .QueryFormID() }}"
					action="{{ .QueryAction() }}"
					method="get"
					@if(.Searchable())
						role="search"
					@endif
				>
					@foreach(.QueryFields() as field)
						<input type="hidden" name="{{ field.Name }}" value="{{ field.Value }}">
					@endforeach
					@if(.Searchable())
						<label class="sr-only" for="{{ .ID }}-search">{{ .SearchTitle() }}</label>
						<input
							data-part="search"
							class="{{ .PartClass("search", "input") }}"
							type="search"
							id="{{ .ID }}-search"
							name="{{ .SearchName }}"
							value="{{ .Query }}"
							autocomplete="off"
							@if(!.Complete)
								hx-get="{{ .SearchAddress() }}"
								hx-trigger="{{ .SearchTrigger() }}"
								hx-include="closest form"
								hx-target="{{ .Target() }}"
								hx-swap="outerHTML"
								hx-select="{{ .Target() }}"
								hx-sync="this:replace"
								hx-indicator="#{{ .ID }}-indicator"
							@endif
							@if(.SearchPlaceholder != "")
								placeholder="{{ .SearchPlaceholder }}"
							@endif
							@attributes(.PartAttrs("search"))
						>
						<span
							data-part="indicator"
							class="{{ .PartClass("indicator", "spinner") }}"
							id="{{ .ID }}-indicator"
							aria-hidden="true"
							@attributes(.PartAttrs("indicator"))
						></span>
					@endif

					@if(.Filtered())
						{{-- A facet is a button that opens a list of checkboxes, as
						     shadcn's faceted filter: the label, and what is chosen
						     beside it. The checkboxes are native and in the form, so
						     with no script the same choice is submitted by Apply. --}}
						@foreach(.Filters as filter)
							<details
								data-part="facet"
								class="{{ .PartClass("facet", "data-table-facet") }}"
								@attributes(.PartAttrs("facet"))
							>
								<summary>
									{!! icons.PlusCircle(icons.Props{}) !!}
									<span>{{ filter.Label }}</span>
									@if(len(.FacetSummary(filter)) > 0)
										<span
											data-part="facet-chosen"
											class="{{ .PartClass("facet-chosen", "data-table-facet-chosen") }}"
											@attributes(.PartAttrs("facet-chosen"))
										>
											@foreach(.FacetSummary(filter) as label)
												<span class="badge" data-variant="secondary">{{ label }}</span>
											@endforeach
										</span>
									@endif
								</summary>
								<div
									class="data-table-facet-menu"
									role="group"
									aria-label="{{ .FacetName(filter) }}"
									@if(!.Complete)
										hx-get="{{ .SearchAddress() }}"
										hx-trigger="change"
										hx-include="closest form"
										hx-target="{{ .Target() }}"
										hx-swap="outerHTML"
										hx-select="{{ .Target() }}"
										hx-sync="closest form:replace"
										hx-indicator="#{{ .ID }}-indicator"
									@endif
								>
									@foreach(filter.Options as option)
										<label class="data-table-facet-option">
											<input
												data-part="filter"
												class="{{ .PartClass("filter", "input") }}"
												type="checkbox"
												name="{{ filter.Key }}"
												value="{{ option.Value }}"
												@if(.Chosen(filter, option))
													checked
												@endif
												@attributes(.PartAttrs("filter"))
											>
											<span>{{ option.Label }}</span>
										</label>
									@endforeach
								</div>
							</details>
						@endforeach
						<noscript>{!! Button(ButtonProps{Label: .ApplyText(), Type: "submit", Variant: "outline", Size: "sm"}) !!}</noscript>
					@endif
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
					<summary>{{ .ColumnsName() }}{!! icons.CaretDown(icons.Props{}) !!}</summary>
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
										class="input"
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

	{{-- The range and the pager sit together, because between them they say
	     one thing: where in the list the reader is. The range is a live region
	     so a page change is announced -- rows changing under somebody is
	     otherwise silent. --}}
	@if(.Paged() || .Complete || .Searchable() || .Chooses() || .Selectable())
		<div
			data-part="footer"
			class="{{ .PartClass("footer", "data-table-footer") }}"
			@attributes(.PartAttrs("footer"))
		>
			<div class="data-table-status">
				@if(.Selectable())
					<p
						data-part="count"
						class="{{ .PartClass("count", "data-table-count") }}"
						role="status"
						aria-live="polite"
						data-selected-template="{{ .SelectedLabel }}"
						@attributes(.PartAttrs("count"))
					>{{ .SelectedText() }}</p>
				@endif
				@if(.Ranged())
					<p
						data-part="showing"
						class="{{ .PartClass("showing", "data-table-showing") }}"
						role="status"
						aria-live="polite"
						@attributes(.PartAttrs("showing"))
					>{{ .Showing() }}</p>
				@endif
			</div>

			<div class="data-table-pager">
				@if(.Chooses())
					<label class="data-table-sizes">
						<span>{{ .PerPageText() }}</span>
						<select
							data-part="sizes"
							class="{{ .PartClass("sizes", "select") }}"
							data-page-size
							@attributes(.PartAttrs("sizes"))
						>
							@foreach(.Sizes() as size)
								<option
									value="{{ size }}"
									@if(size == .PageSize)
										selected
									@endif
								>{{ size }}</option>
							@endforeach
						</select>
					</label>
				@endif
				@if(.Paged() || .Complete)
					{{-- Where the reader is, and the two steps: a numbered row of
					     pages is noise on a list that is read top to bottom. The
					     links stay addresses, so with no script they load the
					     page from the server. --}}
					<nav
						data-part="pagination"
						class="{{ .PartClass("pagination", "data-table-pagination") }}"
						aria-label="{{ .Pager().LandmarkName() }}"
						data-page-of-template="{{ .PageOfTemplate() }}"
						@if(!.Paged())
							hidden
						@endif
						@attributes(.PartAttrs("pagination"))
					>
						@if(.PageOfText() != "")
						<span
							data-part="page-of"
							class="{{ .PartClass("page-of", "data-table-page-of") }}"
							@attributes(.PartAttrs("page-of"))
						>{{ .PageOfText() }}</span>
						@endif
						<a
							data-part="previous"
							class="{{ .PartClass("previous", "btn") }}"
							data-variant="outline"
							data-size="sm"
							href="{{ .Pager().PreviousHref() }}"
							data-page="{{ .PreviousPage() }}"
							aria-label="{{ .Pager().PreviousName() }}"
							@if(!.Pager().HasPrevious())
								aria-disabled="true"
								tabindex="-1"
							@endif
							@if(!.Complete)
								hx-get="{{ .Pager().PreviousHref() }}"
								hx-target="{{ .Target() }}"
								hx-swap="outerHTML"
								hx-select="{{ .Target() }}"
							@endif
							@attributes(.PartAttrs("previous"))
						>{!! icons.CaretLeft(icons.Props{}) !!}<span>{{ .PreviousText() }}</span></a>
						<a
							data-part="next"
							class="{{ .PartClass("next", "btn") }}"
							data-variant="outline"
							data-size="sm"
							href="{{ .Pager().NextHref() }}"
							data-page="{{ .NextPage() }}"
							aria-label="{{ .Pager().NextName() }}"
							@if(!.Pager().HasNext())
								aria-disabled="true"
								tabindex="-1"
							@endif
							@if(!.Complete)
								hx-get="{{ .Pager().NextHref() }}"
								hx-target="{{ .Target() }}"
								hx-swap="outerHTML"
								hx-select="{{ .Target() }}"
							@endif
							@attributes(.PartAttrs("next"))
						><span>{{ .NextText() }}</span>{!! icons.CaretRight(icons.Props{}) !!}</a>
					</nav>
				@endif
			</div>
		</div>
	@endif
</div>
