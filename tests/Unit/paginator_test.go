package unit

import (
	"net/url"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/encryption"
	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/components"
)

// A list's repository returns the rows and the page they are a page of. The
// pager and the data table draw from that page as it came back, so nobody
// copies its current page, its last page and its total into three fields --
// which is what every list did, through an adapter of its own.

var (
	_ components.Paginator = (*pagination.LengthAwarePage)(nil)
	_ components.Paginator = (*pagination.Page)(nil)
	_ components.Paginator = (*pagination.CursorPage)(nil)
)

// listQuery is the query a list page arrived with: a filter the reader chose,
// and a parameter somebody else wrote into the link.
var listQuery = url.Values{"status": {"open"}, "role": {"admin"}}

func signedCursorOptions(path string) pagination.Options {
	signer := pagination.NewCursorSigner(encryption.NewSigner([]byte("0123456789abcdef0123456789abcdef")), 0)
	return pagination.Options{Path: path, Query: listQuery, Signer: signer}
}

// cursorPage is the first page of a keyset list that read one row more than it
// shows, so it has a next page and no previous one.
func cursorPage() *pagination.CursorPage {
	return pagination.NewCursorPage(21, 20, nil, func(i int) map[string]string {
		return map[string]string{"id": strconv.Itoa(i + 1)}
	}, signedCursorOptions("/invoices"))
}

var hrefPattern = regexp.MustCompile(`href="([^"]*)"`)

// hrefs are the addresses a rendering links to, unescaped.
func hrefs(html string) []string {
	var out []string
	for _, match := range hrefPattern.FindAllStringSubmatch(html, -1) {
		out = append(out, strings.ReplaceAll(match[1], "&amp;", "&"))
	}
	return out
}

func TestAPagerDrawsTheLengthAwarePageItWasGiven(t *testing.T) {
	page := pagination.NewLengthAwarePage(20, 137, 20, 3, pagination.Options{Path: "/invoices", Query: url.Values{"status": {"open"}}})
	props := components.PaginationProps{Paginator: page}
	html := string(components.Pagination(props))

	for _, want := range []string{
		`href="/invoices?page=2&amp;status=open"`,
		`href="/invoices?page=4&amp;status=open"`,
		`href="/invoices?page=7&amp;status=open"`,
		`rel="prev"`, `rel="next"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %s:\n%s", want, html)
		}
	}
	if !regexp.MustCompile(`href="/invoices\?page=3&amp;status=open"[^>]*aria-current="page"`).MatchString(html) {
		t.Errorf("page 3 is not marked as the page being read:\n%s", html)
	}

	// The window is the page's own, not one computed again here.
	var want, got []int
	for _, link := range page.Links() {
		want = append(want, link.Page)
	}
	for _, entry := range props.Entries() {
		got = append(got, entry.Number)
	}
	if !slices.Equal(got, want) {
		t.Errorf("Entries = %v, want the page's own window %v", got, want)
	}
}

func TestAPagerOverOnePageDrawsNothing(t *testing.T) {
	var none *pagination.LengthAwarePage
	for name, paginator := range map[string]components.Paginator{
		"one page":          pagination.NewLengthAwarePage(5, 5, 20, 1, pagination.Options{}),
		"no more rows":      pagination.NewPage(5, 20, 1, pagination.Options{}),
		"a nil page":        none,
		"no paginator, nil": nil,
	} {
		if html := strings.TrimSpace(string(components.Pagination(components.PaginationProps{Paginator: paginator}))); html != "" {
			t.Errorf("%s: drew a pager:\n%s", name, html)
		}
	}
}

func TestAPagerOverAPageWithNoTotalDrawsTheTwoSteps(t *testing.T) {
	for name, paginator := range map[string]components.Paginator{
		"simple": pagination.NewPage(21, 20, 2, pagination.Options{Path: "/invoices"}),
		"cursor": cursorPage(),
	} {
		html := string(components.Pagination(components.PaginationProps{Paginator: paginator}))
		if strings.Contains(html, `data-part="link"`) {
			t.Errorf("%s: drew numbers for a page that has none:\n%s", name, html)
		}
		if paginator.NextPageURL() == "" || !strings.Contains(strings.ReplaceAll(html, "&amp;", "&"), `href="`+paginator.NextPageURL()+`"`) {
			t.Errorf("%s: the next step is not the page's own address %q:\n%s", name, paginator.NextPageURL(), html)
		}
		if previous := paginator.PreviousPageURL(); previous != "" && !strings.Contains(strings.ReplaceAll(html, "&amp;", "&"), `href="`+previous+`"`) {
			t.Errorf("%s: the previous step is not the page's own address %q:\n%s", name, previous, html)
		}
	}

	cursor := string(components.Pagination(components.PaginationProps{Paginator: cursorPage()}))
	if strings.Contains(cursor, `rel="prev"`) || strings.Contains(cursor, "data-page=") {
		t.Errorf("a first cursor page drew a previous step or a page number:\n%s", cursor)
	}
}

// twentyRows is one page of a windowed list.
func twentyRows() []components.TableRow {
	rows := make([]components.TableRow, 20)
	for i := range rows {
		key := strconv.Itoa(41 + i)
		rows[i] = components.TableRow{Key: key, Cells: []components.TableCell{{Text: "invoice " + key}}}
	}
	return rows
}

func windowedList(paginator components.Paginator) components.DataTableProps {
	return components.DataTableProps{
		ID: "invoices", Label: "Invoices", URL: "/invoices",
		Columns:    []components.TableColumn{{Label: "Number", Key: "number"}},
		Rows:       twentyRows(),
		Searchable: true, FullPager: true,
		State:     components.DataTableValues{Q: "ada"},
		Paginator: paginator,
	}
}

// The page the query returned says where the reader is and how long the list
// is, and the table draws the numbers from it.
func TestADataTableDrawsTheLengthAwarePage(t *testing.T) {
	page := pagination.NewLengthAwarePage(20, 137, 20, 3, pagination.Options{Path: "/invoices", Query: listQuery})
	given := string(components.DataTable(windowedList(page)))
	for _, want := range []string{"Page 3 of 7", "137 row(s)."} {
		if !strings.Contains(given, want) {
			t.Errorf("missing %q:\n%s", want, given)
		}
	}
}

// The table builds its addresses from its own state, as it does for every other
// control: the search travels, and a parameter that came in on the request does
// not, even though the page carries it.
func TestADataTableKeepsItsOwnAddressesOverThePages(t *testing.T) {
	for name, paginator := range map[string]components.Paginator{
		"length aware": pagination.NewLengthAwarePage(20, 137, 20, 3, pagination.Options{Path: "/invoices", Query: listQuery}),
		"simple":       pagination.NewPage(21, 20, 3, pagination.Options{Path: "/invoices", Query: listQuery}),
		"cursor":       cursorPage(),
	} {
		html := string(components.DataTable(windowedList(paginator)))
		for _, address := range hrefs(html) {
			if strings.Contains(address, "role=admin") {
				t.Errorf("%s: an address carried a parameter the table does not own: %s", name, address)
			}
		}
		if !strings.Contains(strings.ReplaceAll(html, "&amp;", "&"), "q=ada") {
			t.Errorf("%s: the search did not travel on the steps:\n%s", name, html)
		}
	}
}

func TestADataTableWritesThePagesParameterName(t *testing.T) {
	page := pagination.NewLengthAwarePage(20, 137, 20, 3, pagination.Options{Path: "/invoices", PageName: "sheet"})
	props := windowedList(page)
	props.State = components.DataTableValues{Q: "ada", Names: components.DataTableParams{Page: "sheet"}}
	html := strings.ReplaceAll(string(components.DataTable(props)), "&amp;", "&")
	if !strings.Contains(html, "sheet=4") || !strings.Contains(html, "sheet=2") {
		t.Errorf("the steps do not use the state's parameter:\n%s", html)
	}
	if strings.Contains(html, "page=") {
		t.Errorf("an address used the conventional name instead of the state's:\n%s", html)
	}
}

func TestADataTableOverAPageWithNoTotalDrawsTwoStepsAndNoCount(t *testing.T) {
	simple := strings.ReplaceAll(string(components.DataTable(windowedList(pagination.NewPage(21, 20, 3, pagination.Options{})))), "&amp;", "&")
	if !strings.Contains(simple, `href="/invoices?page=2&q=ada"`) || !strings.Contains(simple, `href="/invoices?page=4&q=ada"`) {
		t.Errorf("a simple page did not step to the pages either side:\n%s", simple)
	}

	cursor := strings.ReplaceAll(string(components.DataTable(windowedList(cursorPage()))), "&amp;", "&")
	next := regexp.MustCompile(`data-part="next"[^>]*href="([^"]*)"`).FindStringSubmatch(cursor)
	if next == nil || !strings.Contains(next[1], "cursor=") || !strings.Contains(next[1], "q=ada") {
		t.Errorf("a cursor page's next step does not carry the signed cursor and the search:\n%s", cursor)
	}

	for name, html := range map[string]string{"simple": simple, "cursor": cursor} {
		if strings.Contains(html, `data-part="page-of"`) || strings.Contains(html, `data-part="first"`) || strings.Contains(html, `data-part="last"`) {
			t.Errorf("%s: drew a count or a jump for a page that knows no total:\n%s", name, html)
		}
		if !strings.Contains(html, `data-part="pagination"`) {
			t.Errorf("%s: drew no steps:\n%s", name, html)
		}
	}
}
