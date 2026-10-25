package unit

import (
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/components"
)

// listing is a data table in the state the diagnosis combined: a search, an
// order, a facet, a page, a page size and a hidden column all at once, which is
// the state where a control that forgets one of them shows.
func listing() components.DataTableProps {
	return components.DataTableProps{
		ID: "invoices", Label: "Invoices", URL: "/invoices?role=admin&_method=DELETE#top",
		Searchable: true, FullPager: true, HeaderMenu: false, PageSizes: []int{10, 20, 50},
		State: components.DataTableValues{
			Q: "ada", Sort: "total", Dir: "desc", CurrentPage: 2, Size: 20, DefaultSize: 10,
			Facets: map[string][]string{"status": {"paid", "paid", "open"}}, Hide: []string{"number", "number"},
		},
		Paginator: pagination.NewLengthAwarePage(20, 137, 20, 2, pagination.Options{Path: "/invoices"}),
		Filters: []components.DataTableFilter{{
			Key: "status", Label: "Status",
			Options: []components.SelectOption{{Label: "Paid", Value: "paid"}, {Label: "Open", Value: "open"}, {Label: "Void", Value: "void"}},
			Counts:  map[string]int{"paid": 3, "open": 5, "void": 0},
		}},
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true, Hideable: true},
			{Label: "Total", Key: "total", Sortable: true, Align: "end"},
		},
		Rows:       twentyRows(),
		SelectName: "ids", BulkAction: "/invoices/bulk",
		BulkActions: []components.ButtonProps{
			{Label: "Archive", Type: "submit", Name: "do", Value: "archive"},
			{Label: "Delete", Type: "submit", Name: "do", Value: "delete"},
		},
	}
}

func unescaped(html string) string { return strings.ReplaceAll(html, "&amp;", "&") }

func attrValues(html, attribute string) []string {
	var out []string
	for _, match := range regexp.MustCompile(`\s`+attribute+`="([^"]*)"`).FindAllStringSubmatch(html, -1) {
		out = append(out, unescaped(match[1]))
	}
	return out
}

func parsed(t *testing.T, address string) url.Values {
	t.Helper()
	_, query, _ := strings.Cut(address, "?")
	values, err := url.ParseQuery(query)
	if err != nil {
		t.Fatalf("%s does not parse: %v", address, err)
	}
	return values
}

// The state encodes itself once: a value written twice is written once, in a
// stable order, and what is a default is left out.
func TestTheStateEncodesCanonically(t *testing.T) {
	state := components.DataTableValues{
		Q: "ada", Sort: "total", Dir: "desc", CurrentPage: 1, Size: 10, DefaultSize: 10,
		Facets: map[string][]string{"status": {"paid", "open", "paid", ""}, "empty": nil},
		Hide:   []string{"b", "a", "b"},
	}
	if got, want := state.Values(components.DataTableChange{}).Encode(), "dir=desc&hide=a&hide=b&q=ada&sort=total&status=open&status=paid"; got != want {
		t.Errorf("Values() = %q, want %q", got, want)
	}
	if got := (components.DataTableValues{}).Values(components.DataTableChange{}).Encode(); got != "" {
		t.Errorf("an unfiltered list encodes %q, want nothing", got)
	}
}

// A facet ticked and then unticked is the state it started from, and never a
// state with the value twice.
func TestAFacetToggledOnAndOffReturnsToWhereItStarted(t *testing.T) {
	start := components.DataTableValues{Q: "ada", Facets: map[string][]string{"status": {"paid"}}}
	on := start.Values(components.DataTableChange{Facet: "status", Value: "open"})
	if got := on["status"]; strings.Join(got, ",") != "open,paid" {
		t.Fatalf("ticking open gave status=%v, want open,paid", got)
	}
	again := components.DataTableValues{Q: "ada", Facets: map[string][]string{"status": on["status"]}}
	off := again.Values(components.DataTableChange{Facet: "status", Value: "open"})
	if got, want := off.Encode(), start.Values(components.DataTableChange{}).Encode(); got != want {
		t.Errorf("unticking gave %q, want the starting %q", got, want)
	}
	last := components.DataTableValues{Facets: map[string][]string{"status": {"paid"}}}
	if got := last.Values(components.DataTableChange{Facet: "status", Value: "paid"}); len(got) != 0 {
		t.Errorf("unticking the last value left %v", got)
	}
}

// What a step moves besides itself: a new search, facet or order goes back to
// the first page, a new size keeps the first row on screen, and the columns
// leave the page where it is.
func TestEveryStepMovesTheStateItShould(t *testing.T) {
	state := components.DataTableValues{Q: "ada", Sort: "total", Dir: "asc", CurrentPage: 3, Size: 25, DefaultSize: 10,
		Facets: map[string][]string{"status": {"paid"}}}
	for name, c := range map[string]struct {
		change components.DataTableChange
		want   url.Values
	}{
		"sort goes to page one": {components.DataTableChange{Sort: "number", Dir: "desc"},
			url.Values{"q": {"ada"}, "sort": {"number"}, "dir": {"desc"}, "status": {"paid"}, "per_page": {"25"}}},
		"a facet goes to page one": {components.DataTableChange{Facet: "status", Value: "open"},
			url.Values{"q": {"ada"}, "sort": {"total"}, "dir": {"asc"}, "status": {"open", "paid"}, "per_page": {"25"}}},
		"reset drops the search and the facets": {components.DataTableChange{Reset: true},
			url.Values{"sort": {"total"}, "dir": {"asc"}, "per_page": {"25"}}},
		"a size keeps the first row": {components.DataTableChange{PerPage: 10},
			url.Values{"q": {"ada"}, "sort": {"total"}, "dir": {"asc"}, "status": {"paid"}, "page": {"6"}}},
		"a column keeps the page": {components.DataTableChange{Hide: "number"},
			url.Values{"q": {"ada"}, "sort": {"total"}, "dir": {"asc"}, "status": {"paid"}, "per_page": {"25"}, "page": {"3"}, "hide": {"number"}}},
		"a page is a page": {components.DataTableChange{Page: 5},
			url.Values{"q": {"ada"}, "sort": {"total"}, "dir": {"asc"}, "status": {"paid"}, "per_page": {"25"}, "page": {"5"}}},
	} {
		if got, want := state.Values(c.change).Encode(), c.want.Encode(); got != want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, want)
		}
	}
	// The reference's own example: 25 a page, page two, down to 10 a page,
	// lands on page three, which holds row 26.
	two := components.DataTableValues{CurrentPage: 2, Size: 25, DefaultSize: 25}
	if got := two.Values(components.DataTableChange{PerPage: 10}).Get("page"); got != "3" {
		t.Errorf("25 to 10 on page 2 landed on page %q, want 3", got)
	}
}

// Every address the table draws is the state's own answer, once, and the htmx
// attribute on it is the same address.
func TestEveryControlOfACombinedStateIsAnExactAddress(t *testing.T) {
	html := unescaped(string(components.DataTable(listing())))

	for name, want := range map[string]string{
		"sort number":  "/invoices?dir=asc&hide=number&per_page=20&q=ada&sort=number&status=open&status=paid",
		"sort total":   "/invoices?dir=asc&hide=number&per_page=20&q=ada&sort=total&status=open&status=paid",
		"previous":     "/invoices?dir=desc&hide=number&per_page=20&q=ada&sort=total&status=open&status=paid",
		"next":         "/invoices?dir=desc&hide=number&page=3&per_page=20&q=ada&sort=total&status=open&status=paid",
		"last":         "/invoices?dir=desc&hide=number&page=7&per_page=20&q=ada&sort=total&status=open&status=paid",
		"reset":        "/invoices?dir=desc&hide=number&per_page=20&sort=total",
		"clear status": "/invoices?dir=desc&hide=number&per_page=20&q=ada&sort=total",
		"show number":  "/invoices?dir=desc&page=2&per_page=20&q=ada&sort=total&status=open&status=paid",
	} {
		if !strings.Contains(html, `href="`+want+`"`) {
			t.Errorf("%s: no link to %s in\n%s", name, want, html)
		}
	}
	if strings.Contains(html, "role=admin") || strings.Contains(html, "_method") || strings.Contains(html, "#top") {
		t.Errorf("a parameter of URL's own query was carried:\n%s", html)
	}

	// Sorting by the column in force asks for the other direction, and never
	// for no order.
	if !strings.Contains(html, `href="/invoices?dir=asc&hide=number&per_page=20&q=ada&sort=total&status=open&status=paid"`) {
		t.Errorf("a header sorted descending does not ask for ascending:\n%s", html)
	}

	for _, address := range hrefs(html) {
		values := parsed(t, address)
		for name, list := range values {
			seen := map[string]bool{}
			for _, value := range list {
				if seen[value] {
					t.Errorf("%s=%s is written twice in %s", name, value, address)
				}
				seen[value] = true
			}
		}
	}

	// Every link that asks htmx for the swap asks for the address it points at.
	tags := regexp.MustCompile(`<a\b[^>]*hx-get="[^"]*"[^>]*>`).FindAllString(html, -1)
	if len(tags) < 8 {
		t.Fatalf("expected the headers, the steps, the columns, reset and clear to swap, found %d", len(tags))
	}
	for _, tag := range tags {
		href := regexp.MustCompile(`href="([^"]*)"`).FindStringSubmatch(tag)
		get := regexp.MustCompile(`hx-get="([^"]*)"`).FindStringSubmatch(tag)
		if href == nil || href[1] != get[1] {
			t.Errorf("a link's hx-get is not its href: %s", tag)
		}
	}
}

// The form is the one owner of the search, the facets and the size, and carries
// the order and the hidden columns once each.
func TestTheFormCarriesTheStateOnceAndMakesOneRequest(t *testing.T) {
	html := unescaped(string(components.DataTable(listing())))

	forms := regexp.MustCompile(`(?s)<form\b[^>]*>`).FindAllString(html, -1)
	var query string
	for _, form := range forms {
		if strings.Contains(form, `id="invoices-query"`) {
			query = form
		}
	}
	if query == "" {
		t.Fatalf("no query form:\n%s", html)
	}
	for _, want := range []string{
		`hx-get="/invoices"`, `hx-push-url="true"`, `hx-target="#invoices"`, `hx-swap="outerHTML"`,
		`hx-select="#invoices"`, `hx-trigger="submit, change"`,
	} {
		if !strings.Contains(query, want) {
			t.Errorf("the form is missing %s: %s", want, query)
		}
	}
	if strings.Count(html, `hx-trigger=`) != 1 {
		t.Errorf("more than one element listens for a trigger; the state has two owners:\n%s", html)
	}

	hidden := regexp.MustCompile(`<input type="hidden" name="([^"]*)" value="([^"]*)">`).FindAllStringSubmatch(html, -1)
	got := map[string]int{}
	for _, field := range hidden {
		if field[1] == "_token" {
			continue
		}
		got[field[1]+"="+field[2]]++
	}
	for field, count := range got {
		if count != 1 {
			t.Errorf("hidden field %s appears %d times", field, count)
		}
	}
	for _, want := range []string{"dir=desc", "hide=number", "sort=total", "per_page=20"} {
		if got[want] != 1 {
			t.Errorf("the form does not carry %s once: %v", want, got)
		}
	}
	for _, owned := range []string{"q=", "page=", "status="} {
		for field := range got {
			if strings.HasPrefix(field, owned) {
				t.Errorf("the form carries %s as a hidden field, but a control owns it", field)
			}
		}
	}
	// The size is a menu of addresses: the state keeps the first row on screen.
	if !strings.Contains(html, `href="/invoices?dir=desc&hide=number&page=3&q=ada&sort=total&status=open&status=paid"`) {
		t.Errorf("ten a page from page two of twenty does not land on page three, which holds row 21:\n%s", html)
	}
	if !strings.Contains(html, `href="/invoices?dir=desc&hide=number&per_page=50&q=ada&sort=total&status=open&status=paid"`) {
		t.Errorf("fifty a page from page two of twenty does not land on page one:\n%s", html)
	}
}

func TestACompleteListNeverAsksTheServerAndKeepsCorrectAddresses(t *testing.T) {
	props := listing()
	props.Complete, props.Paginator, props.PageSize = true, nil, 4
	props.State = components.DataTableValues{DefaultSize: 4}
	html := unescaped(string(components.DataTable(props)))

	if strings.Contains(html, "hx-get") || strings.Contains(html, "hx-trigger") || strings.Contains(html, "hx-push-url") {
		t.Errorf("a complete list asks the server:\n%s", html)
	}
	for _, link := range regexp.MustCompile(`<a\b[^>]*href="[^"]*"[^>]*>`).FindAllString(html, -1) {
		if !strings.Contains(link, `hx-boost="false"`) {
			t.Errorf("a link of a complete list may be boosted into a request: %s", link)
		}
	}
	if !strings.Contains(html, `hx-boost="false"`) || !strings.Contains(html, `data-complete="true"`) {
		t.Errorf("a complete list is not marked, nor kept out of the boost:\n%s", html)
	}
	for _, address := range hrefs(html) {
		if strings.Contains(address, "page=-") || strings.Contains(address, "page=0") {
			t.Errorf("an address leaves the list: %s", address)
		}
	}
	// 20 rows, 4 to a page: the first page, so the first step goes nowhere and
	// the next goes to page two.
	if !regexp.MustCompile(`data-step="previous"[^>]*aria-disabled="true"`).MatchString(html) &&
		!regexp.MustCompile(`aria-disabled="true"[^>]*data-step="previous"`).MatchString(html) {
		t.Errorf("the previous step is not disabled on the first page:\n%s", html)
	}
	if !strings.Contains(html, `href="/invoices?page=2"`) {
		t.Errorf("the next step does not go to page two:\n%s", html)
	}
	if strings.Count(html, "<tr") < 22 {
		t.Errorf("a complete list must draw every row, hidden outside the window, found %d", strings.Count(html, "<tr"))
	}
	if !strings.Contains(html, "data-empty-row") {
		t.Errorf("a complete list draws no row for a search that matches nothing")
	}
}

func TestThePagerIsDisabledAtTheEnds(t *testing.T) {
	first := components.DataTableProps{ID: "p", URL: "/p", FullPager: true,
		Columns:   []components.TableColumn{{Label: "N"}},
		Rows:      twentyRows(),
		Paginator: pagination.NewLengthAwarePage(20, 60, 20, 1, pagination.Options{})}
	html := unescaped(string(components.DataTable(first)))
	for _, step := range []string{"first", "previous"} {
		tag := regexp.MustCompile(`<a\b[^>]*data-step="` + step + `"[^>]*>`).FindString(html)
		if !strings.Contains(tag, `aria-disabled="true"`) || strings.Contains(tag, "href=") || strings.Contains(tag, "hx-get") {
			t.Errorf("%s on the first page must be disabled and go nowhere: %s", step, tag)
		}
	}
	for _, step := range []string{"next", "last"} {
		tag := regexp.MustCompile(`<a\b[^>]*data-step="` + step + `"[^>]*>`).FindString(html)
		if strings.Contains(tag, `aria-disabled="true"`) || !strings.Contains(tag, "href=") {
			t.Errorf("%s on the first page must be a link: %s", step, tag)
		}
	}
	if !strings.Contains(html, "Page 1 of 3") {
		t.Errorf("no 'Page 1 of 3':\n%s", html)
	}
	for _, step := range []string{"first", "last"} {
		tag := regexp.MustCompile(`<a\b[^>]*data-step="` + step + `"[^>]*>`).FindString(html)
		if !strings.Contains(tag, "hidden") || !strings.Contains(tag, "lg:flex") {
			t.Errorf("%s is not hidden below the large breakpoint: %s", step, tag)
		}
	}
}

// A list with nothing in it says the caller's Empty; a list that has rows and
// matched none says so in the table.
func TestNothingMatchedIsNotNothingExists(t *testing.T) {
	base := components.DataTableProps{ID: "n", URL: "/n", Label: "Notes", Searchable: true,
		Columns: []components.TableColumn{{Label: "Title"}},
		Empty:   components.EmptyProps{Title: "No notes yet", Message: "Write the first."}}

	html := string(components.DataTable(base))
	if !strings.Contains(html, "No notes yet") || strings.Contains(html, "No results.") {
		t.Errorf("an empty list must say the caller's Empty:\n%s", html)
	}

	filtered := base
	filtered.State = components.DataTableValues{Q: "zzz"}
	html = string(components.DataTable(filtered))
	if !strings.Contains(html, "No results.") || strings.Contains(html, "No notes yet") {
		t.Errorf("a search with no match must say No results.:\n%s", html)
	}
	if !strings.Contains(html, `colspan="1"`) || !strings.Contains(html, "h-24") || !strings.Contains(html, `data-part="reset"`) {
		t.Errorf("the empty row is not the reference's, or Reset is missing:\n%s", html)
	}
}

func TestHiddenColumnsComeFromTheStateAndOnlyHideableOnesAreListed(t *testing.T) {
	props := components.DataTableProps{ID: "h", URL: "/h", Label: "H",
		State: components.DataTableValues{Hide: []string{"status", "id"}},
		Columns: []components.TableColumn{
			{Label: "Id", Key: "id"},
			{Label: "Status", Key: "status", Hideable: true},
			{Label: "Total", Key: "total", Hideable: true},
		},
		Rows: []components.TableRow{{Cells: []components.TableCell{{Text: "1"}, {Text: "open"}, {Text: "9"}}}}}
	html := string(components.DataTable(props))

	if !regexp.MustCompile(`data-column="status"\s+hidden`).MatchString(html) {
		t.Errorf("a column the state hides is drawn:\n%s", html)
	}
	if regexp.MustCompile(`data-column="id"\s+hidden`).MatchString(html) {
		t.Errorf("a column that cannot be hidden was hidden by the state")
	}
	menu := html[strings.Index(html, `data-part="columns"`):]
	menu = menu[:strings.Index(menu, "</details>")]
	if strings.Contains(menu, ">\n\t\t\t\t\t\t\t\tId\n") || strings.Count(menu, `data-column-toggle=`) != 2 {
		t.Errorf("the menu must list the two hideable columns:\n%s", menu)
	}
	if !strings.Contains(menu, "capitalize") || !strings.Contains(menu, `aria-checked="false"`) || !strings.Contains(menu, `aria-checked="true"`) {
		t.Errorf("the menu does not show which columns are on:\n%s", menu)
	}
	// A list without a state still starts a column off by its own Hidden.
	props.State = nil
	props.Columns[2].Hidden = true
	if html := string(components.DataTable(props)); !regexp.MustCompile(`data-column="total"\s+hidden`).MatchString(html) {
		t.Errorf("a column's own Hidden is not where a stateless list starts:\n%s", html)
	}
}

func TestTheBulkBarNeedsAnExplicitDestination(t *testing.T) {
	props := components.DataTableProps{ID: "b", URL: "/b", Label: "B", SelectName: "ids",
		BulkActions: []components.ButtonProps{{Label: "Archive", Type: "submit", Name: "do", Value: "archive"}},
		Columns:     []components.TableColumn{{Label: "N"}},
		Rows:        []components.TableRow{{Key: "1", Cells: []components.TableCell{{Text: "x"}}}}}

	html := string(components.DataTable(props))
	if strings.Contains(html, `data-part="bulk"`) || strings.Contains(html, "<form id=\"b-bulk\"") || strings.Contains(html, `form="b-bulk"`) {
		t.Errorf("a bulk bar with nowhere to post was drawn:\n%s", html)
	}

	props.BulkAction, props.Token = "/b/bulk", "tok"
	html = string(components.DataTable(props))
	for _, want := range []string{`action="/b/bulk"`, `method="post"`, `name="do"`, `value="archive"`, `name="_token" value="tok"`, `form="b-bulk"`, `name="ids"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s:\n%s", want, html)
		}
	}
	if !regexp.MustCompile(`id="b-bulk"[^>]*hidden`).MatchString(html) {
		t.Errorf("the bar is shown with nothing ticked:\n%s", html)
	}
	props.Rows[0].Selected = true
	if html = string(components.DataTable(props)); regexp.MustCompile(`id="b-bulk"[^>]*hidden`).MatchString(html) {
		t.Errorf("the bar is hidden with a row ticked")
	}
}

func TestTheSortIconAndNameFollowTheDirection(t *testing.T) {
	for name, c := range map[string]struct {
		dir  string
		icon string
		sort string
	}{"none": {"", "none", "Number, not sorted"}, "asc": {"asc", "asc", "Number, sorted ascending"}, "desc": {"desc", "desc", "Number, sorted descending"}} {
		props := components.DataTableProps{ID: "s", URL: "/s", Label: "S",
			Columns: []components.TableColumn{{Label: "Number", Key: "number", Sortable: true}},
			Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "1"}}}}}
		if c.dir != "" {
			props.State = components.DataTableValues{Sort: "number", Dir: c.dir}
		}
		html := string(components.DataTable(props))
		if !strings.Contains(html, `aria-label="`+c.sort+`"`) {
			t.Errorf("%s: the header is not called %q with no catalogue:\n%s", name, c.sort, html)
		}
		for _, icon := range []string{"none", "asc", "desc"} {
			tag := regexp.MustCompile(`<span\s+data-sort-icon="` + icon + `"[^>]*>`).FindString(html)
			if tag == "" {
				t.Errorf("%s: no %s icon", name, icon)
			}
			if (icon == c.icon) == strings.Contains(tag, "hidden") {
				t.Errorf("%s: the %s icon visibility is wrong: %s", name, icon, tag)
			}
		}
	}
}

func TestTheHeaderMenuOffersAscDescAndHide(t *testing.T) {
	props := components.DataTableProps{ID: "m", URL: "/m", Label: "M", HeaderMenu: true,
		State:   components.DataTableValues{Sort: "number", Dir: "asc"},
		Columns: []components.TableColumn{{Label: "Number", Key: "number", Sortable: true, Hideable: true}},
		Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "1"}}}}}
	html := unescaped(string(components.DataTable(props)))
	for _, want := range []string{
		`href="/m?dir=asc&sort=number"`, `href="/m?dir=desc&sort=number"`, `href="/m?hide=number&dir=asc&sort=number"`,
		">Asc<", ">Desc<", ">Hide<", `role="menuitem"`,
	} {
		if want == `href="/m?hide=number&dir=asc&sort=number"` {
			want = `href="/m?dir=asc&hide=number&sort=number"`
		}
		if !strings.Contains(html, want) {
			t.Errorf("missing %s:\n%s", want, html)
		}
	}
}

func TestTheFooterCountsTheSelectionAgainstTheWholeList(t *testing.T) {
	html := string(components.DataTable(listing()))
	if !strings.Contains(html, "0 of 137 row(s) selected.") {
		t.Errorf("the footer does not count against the query's total:\n%s", html)
	}
	for _, want := range []string{`data-total="137"`, `data-url="/invoices?dir=desc`} {
		if !strings.Contains(unescaped(html), want) {
			t.Errorf("the region does not say %s", want)
		}
	}
	plain := listing()
	plain.SelectName, plain.BulkActions, plain.BulkAction = "", nil, ""
	if html := string(components.DataTable(plain)); !strings.Contains(html, "137 row(s).") {
		t.Errorf("an unselectable list does not count its rows:\n%s", html)
	}
}

func TestAFacetDrawsItsCountsItsSearchAndItsClearLink(t *testing.T) {
	html := unescaped(string(components.DataTable(listing())))
	for _, want := range []string{"border-dashed", "font-mono text-xs", "data-facet-empty", "Clear filters", "Reset", "2 selected", "lg:hidden"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q:\n%s", want, html)
		}
	}
	// Two values chosen show two badges; three would show the count.
	state := components.DataTableValues{Facets: map[string][]string{"status": {"paid", "open", "void"}}}
	props := listing()
	props.State = state
	if html := string(components.DataTable(props)); !regexp.MustCompile(`data-facet-many\s+class="[^"]*"\s*>3 selected`).MatchString(html) {
		t.Errorf("three chosen values do not show '3 selected':\n%s", html)
	}
}

func TestEveryTextTheTableDrawsIsEscaped(t *testing.T) {
	bad := `<script>alert(1)</script>`
	props := listing()
	props.Label, props.Caption = bad, bad
	props.Filters[0].Label = bad
	props.Filters[0].Options[0].Label = bad
	props.Columns[0].Label = bad
	props.Rows[0].Label = bad
	props.Rows[0].Cells[0].Text = bad
	props.Rows[0].Actions = []components.MenuItem{{Label: bad, URL: "/x"}}
	props.State = components.DataTableValues{Q: bad, Facets: map[string][]string{"status": {bad}}}
	if html := string(components.DataTable(props)); strings.Contains(html, "<script>alert(1)") {
		t.Errorf("a prop reached the page unescaped:\n%s", html)
	}
}

// Every control that replaces the region takes the region out of the answer by
// the same selector. An endpoint that answers the whole page -- which is what
// it answers without htmx -- then still swaps one list for one list, instead of
// drawing the page inside itself.
func TestEveryRegionSwapSelectsTheRegionFromTheAnswer(t *testing.T) {
	props := listing()
	props.HeaderMenu = true
	html := string(components.DataTable(props))

	swaps := regexp.MustCompile(`<[a-z]+\b[^>]*hx-target="#invoices"[^>]*>`).FindAllString(html, -1)
	if len(swaps) < 12 {
		t.Fatalf("expected the form, the headers, the facets, the columns, the sizes and the steps to swap the region, found %d", len(swaps))
	}
	for _, tag := range swaps {
		for _, want := range []string{`hx-swap="outerHTML"`, `hx-select="#invoices"`, `hx-push-url="true"`, `hx-sync="closest [data-part=root]:replace"`} {
			if !strings.Contains(tag, want) {
				t.Errorf("a control that swaps the region is missing %s:\n%s", want, tag)
			}
		}
	}
	// And nothing asks for a swap without a target: the closest ancestor would
	// answer, and that is a caller's own element.
	for _, tag := range regexp.MustCompile(`<[a-z]+\b[^>]*hx-get="[^"]*"[^>]*>`).FindAllString(html, -1) {
		if !strings.Contains(tag, `hx-target="#invoices"`) {
			t.Errorf("a control asks the server without saying what it swaps:\n%s", tag)
		}
	}
}
