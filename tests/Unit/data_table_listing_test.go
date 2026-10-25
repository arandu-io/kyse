package unit

import (
	"net/url"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/components"
)

func invoiceListing(t *testing.T) *pagination.Listing {
	t.Helper()
	listing, err := pagination.NewListing(pagination.ListingSpec{
		Search:      []string{"number", "customer"},
		Sorts:       map[string]string{"number": "number", "customer": "customer", "total": "total"},
		DefaultSort: "number",
		TieBreak:    "id",
		Facets:      []pagination.FacetSpec{{Key: "status", Param: "s", Column: "status", Values: []string{"paid", "open", "void"}}},
		PageSizes:   []int{10, 25, 50},
		Hideable:    []string{"customer", "total"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return listing
}

func parsedState(t *testing.T, raw string) pagination.ListingState {
	t.Helper()
	values, err := url.ParseQuery(raw)
	if err != nil {
		t.Fatal(err)
	}
	return invoiceListing(t).Parse(values)
}

// The adapter adds no rule of its own: every step is the listing state's own
// method, and the query it answers is the state's own Values.
func TestTheListingAdapterAnswersWithTheListingStatesOwnSteps(t *testing.T) {
	state := parsedState(t, "q=ada&s=paid&sort=total&dir=desc&page=3&per_page=25&hide=customer")
	adapted := components.DataTableStateFrom(state)

	for name, c := range map[string]struct {
		change components.DataTableChange
		want   pagination.ListingState
	}{
		"no step":                  {components.DataTableChange{}, state},
		"a new order":              {components.DataTableChange{Sort: "number", Dir: "asc"}, state.WithSort("number", pagination.SortAsc)},
		"the other direction":      {components.DataTableChange{Sort: "total", Dir: "asc"}, state.WithSort("total", pagination.SortAsc)},
		"a facet value on":         {components.DataTableChange{Facet: "status", Value: "open"}, state.ToggleFacet("status", "open")},
		"a facet value off":        {components.DataTableChange{Facet: "status", Value: "paid"}, state.ToggleFacet("status", "paid")},
		"a facet cleared":          {components.DataTableChange{ClearFacet: "status"}, state.ClearFacet("status")},
		"the filters reset":        {components.DataTableChange{Reset: true}, state.ClearFilters()},
		"a size, same row":         {components.DataTableChange{PerPage: 10}, state.WithPageSize(10)},
		"a column on":              {components.DataTableChange{Hide: "customer"}, state.ToggleHidden("customer")},
		"a column off":             {components.DataTableChange{Hide: "total"}, state.ToggleHidden("total")},
		"a page":                   {components.DataTableChange{Page: 5}, state.WithPage(5)},
		"a size not offered":       {components.DataTableChange{PerPage: 7}, state},
		"a sort not declared":      {components.DataTableChange{Sort: "password", Dir: "asc"}, state},
		"a column not hideable":    {components.DataTableChange{Hide: "number"}, state},
		"a facet value not listed": {components.DataTableChange{Facet: "status", Value: "refunded"}, state},
		"a cursor":                 {components.DataTableChange{CursorName: "cursor", Cursor: "abc"}, state},
	} {
		if got, want := adapted.Values(c.change).Encode(), c.want.Encode(); got != want {
			t.Errorf("%s:\n got %s\nwant %s", name, got, want)
		}
	}

	// Row 51 opened page three at 25 a page; at 10 a page it is on page six.
	if got := adapted.Values(components.DataTableChange{PerPage: 10}).Get("page"); got != "6" {
		t.Errorf("a size of 10 from page 3 of 25 landed on page %q, want 6", got)
	}
	// What the table reads of it.
	if adapted.Search() != "ada" || adapted.SortKey() != "total" || adapted.SortDir() != "desc" ||
		adapted.Page() != 3 || adapted.PageSize() != 25 ||
		strings.Join(adapted.Chosen("status"), ",") != "paid" || strings.Join(adapted.Hidden(), ",") != "customer" {
		t.Errorf("the adapter reads the state wrongly: %+v", adapted)
	}
	// A listing's own names reach the table.
	if params := adapted.Params().Filled(); params.Search != "q" || params.PerPage != "per_page" || params.Hide != "hide" {
		t.Errorf("Params = %+v", params)
	}
}

// A state that no listing read draws as a list nobody has touched.
func TestAStateNoListingReadIsAnUntouchedList(t *testing.T) {
	adapted := components.DataTableStateFrom(pagination.ListingState{})
	if adapted.Search() != "" || adapted.SortKey() != "" || adapted.SortDir() != "" || adapted.Page() != 1 {
		t.Errorf("the zero state reads %+v", adapted)
	}
	if got := adapted.Values(components.DataTableChange{Page: 4}).Encode(); got != "" {
		t.Errorf("the zero state wrote %q", got)
	}
}

// A service hands the result on and the table is drawn from it: the state, the
// page with its total, the counts of each facet value, the facet's own
// parameter and the sizes the listing allows.
func TestAResultFillsTheStateThePageTheCountsAndTheSizes(t *testing.T) {
	state := parsedState(t, "s=paid&s=open&page=2&per_page=25&hide=total")
	counts := map[string][]pagination.FacetCount{"status": {
		{Value: "paid", Count: 24, Selected: true}, {Value: "open", Count: 7, Selected: true}, {Value: "void", Count: 0},
	}}
	result := pagination.NewListingResult(state, pagination.NewLengthAwarePage(6, 31, 25, 2, state.PageOptions("/invoices")), counts)

	props := components.DataTableProps{
		ID: "invoices", Label: "Invoices", URL: "/invoices", Searchable: true, FullPager: true,
		Result:  result,
		Filters: []components.DataTableFilter{{Key: "status", Label: "Status"}},
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true},
			{Label: "Total", Key: "total", Sortable: true, Hideable: true},
		},
		Rows: []components.TableRow{{Key: "1", Cells: []components.TableCell{{Text: "INV-001"}, {Text: "10"}}}},
	}
	html := strings.Join(strings.Fields(strings.ReplaceAll(string(components.DataTable(props)), "&amp;", "&")), " ")
	for _, want := range []string{
		"Page 2 of 2", "31 row(s).",
		`name="s" value="paid"`, `name="s" value="open"`, `name="s" value="void"`,
		`>24</span>`, `>7</span>`,
		`data-column="total"`, `hidden`,
		`data-size="10"`, `data-size="25"`, `data-size="50"`,
		`href="/invoices?hide=total&per_page=25&s=paid&s=open"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s:\n%s", want, html)
		}
	}
	// The options came from the facet's own values, once each.
	if strings.Count(html, `name="s" value="paid"`) != 1 {
		t.Errorf("a facet option is drawn more than once")
	}
	// Nothing is built from a parameter that is not the listing's.
	if strings.Contains(html, "status=") {
		t.Errorf("an address used the facet's key where its parameter is s:\n%s", html)
	}
}
