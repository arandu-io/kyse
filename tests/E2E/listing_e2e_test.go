//go:build e2e

package e2e

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/auth"
	"github.com/arandu-io/hesape/database/query"
	"github.com/arandu-io/hesape/database/query/grammars"
	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// This is the flow an application is built on, end to end and with nothing in
// the middle written for the test: a listing declared once with
// pagination.NewListing, the request's query parsed by it, the hesape query
// builder applying it to a real SQLite database under the tenant of a Grant, the
// result handed to the data table as it came back, and a browser clicking what
// the table drew.
//
// SQLite is the sqlite3 program: the module has two dependencies and a driver
// would be a third, so the connection below runs each statement the builder
// compiled through the program, bindings written in as literals, and reads the
// rows back as JSON. The statements are the SQLite grammar's own. A machine
// with no sqlite3 skips these.

// sqliteProgram runs the statements of one database file.
type sqliteProgram struct {
	binary string
	file   string
}

func newSQLite(t *testing.T) *sqliteProgram {
	t.Helper()
	binary, err := exec.LookPath("sqlite3")
	if err != nil {
		t.Skip("there is no sqlite3 program on this machine")
	}
	return &sqliteProgram{binary: binary, file: filepath.Join(t.TempDir(), "listing.db")}
}

func (s *sqliteProgram) run(sql string, json bool) ([]byte, error) {
	args := []string{}
	if json {
		args = append(args, "-json")
	}
	args = append(args, s.file, sql)
	out, err := exec.Command(s.binary, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sqlite3: %v: %s\n%s", err, out, sql)
	}
	return out, nil
}

// literal writes a binding into a statement.
func literal(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case bool:
		if v {
			return "1"
		}
		return "0"
	}
	return "'" + strings.ReplaceAll(fmt.Sprint(value), "'", "''") + "'"
}

// Select is the only door the listing uses.
func (s *sqliteProgram) Select(_ context.Context, sql string, bindings []any, _ bool) ([]query.Record, error) {
	var out strings.Builder
	next := 0
	for _, r := range sql {
		if r == '?' {
			if next >= len(bindings) {
				return nil, errors.New("more placeholders than bindings")
			}
			out.WriteString(literal(bindings[next]))
			next++
			continue
		}
		out.WriteRune(r)
	}
	if next != len(bindings) {
		return nil, errors.New("more bindings than placeholders")
	}
	raw, err := s.run(out.String(), true)
	if err != nil {
		return nil, err
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		return nil, nil
	}
	var rows []query.Record
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}
	for _, row := range rows {
		for key, value := range row {
			if number, ok := value.(float64); ok && number == float64(int64(number)) {
				row[key] = int64(number)
			}
		}
	}
	return rows, nil
}

func (s *sqliteProgram) Insert(context.Context, string, []any) (bool, error) {
	return false, errors.New("read only")
}
func (s *sqliteProgram) Update(context.Context, string, []any) (int64, error) {
	return 0, errors.New("read only")
}
func (s *sqliteProgram) Delete(context.Context, string, []any) (int64, error) {
	return 0, errors.New("read only")
}
func (s *sqliteProgram) Statement(context.Context, string, []any) (bool, error) {
	return false, errors.New("read only")
}

type passThrough struct{}

func (passThrough) ProcessSelect(_ *query.Builder, rows []query.Record) []query.Record { return rows }
func (passThrough) ProcessInsertGetID(context.Context, *query.Builder, string, []any, string) (int64, error) {
	return 0, errors.New("read only")
}

// seed writes the invoices of the list under test: the sixty of every other
// test under the tenant acme, and five more under another tenant that must
// never appear.
func (s *sqliteProgram) seed(t *testing.T) {
	t.Helper()
	var sql strings.Builder
	sql.WriteString(`create table invoices (id integer primary key, tenant_id text not null, number text, customer text, status text, total integer);`)
	for _, one := range invoices() {
		fmt.Fprintf(&sql, "insert into invoices values (%d, 'acme', %s, %s, %s, %d);",
			one.ID, literal(one.Number), literal(one.Customer), literal(one.Status), one.Total)
	}
	for i := 1; i <= 5; i++ {
		fmt.Fprintf(&sql, "insert into invoices values (%d, 'other', 'OTH-%03d', 'Ada Lovelace', 'paid', %d);", 100+i, i, 1000+i)
	}
	if _, err := s.run(sql.String(), false); err != nil {
		t.Fatal(err)
	}
}

func invoiceListing(t *testing.T) *pagination.Listing {
	t.Helper()
	listing, err := pagination.NewListing(pagination.ListingSpec{
		Search:      []string{"number", "customer"},
		Sorts:       map[string]string{"number": "number", "customer": "customer", "status": "status", "total": "total"},
		DefaultSort: "number",
		TieBreak:    "id",
		Facets: []pagination.FacetSpec{{
			Key: "status", Column: "status", Values: []string{"paid", "open", "void"},
		}},
		PageSizes: []int{10, 20, 50},
		Hideable:  []string{"customer", "status", "total"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return listing
}

// reference is what the list should hold for a query, worked out in Go from the
// rows alone: the answer the database has to agree with.
func reference(values url.Values) (numbers []string, total int, counts map[string]int) {
	needle := strings.ToLower(values.Get("q"))
	chosen := map[string]bool{}
	for _, status := range values["status"] {
		chosen[status] = true
	}
	counts = map[string]int{}
	var rows []invoice
	for _, one := range invoices() {
		if needle != "" && !strings.Contains(strings.ToLower(one.Number+" "+one.Customer), needle) {
			continue
		}
		counts[one.Status]++
		if len(chosen) > 0 && !chosen[one.Status] {
			continue
		}
		rows = append(rows, one)
	}
	key, desc := values.Get("sort"), values.Get("dir") == "desc"
	if key == "" {
		key = "number"
	}
	sort.SliceStable(rows, func(a, b int) bool {
		var order int
		if key == "total" {
			order = rows[a].Total - rows[b].Total
		} else {
			order = strings.Compare(keyOf(rows[a], key), keyOf(rows[b], key))
		}
		if order == 0 {
			return rows[a].ID < rows[b].ID
		}
		if desc {
			return order > 0
		}
		return order < 0
	})
	size, _ := strconv.Atoi(values.Get("per_page"))
	if size == 0 {
		size = 10
	}
	page, _ := strconv.Atoi(values.Get("page"))
	page = max(page, 1)
	total = len(rows)
	from := min((page-1)*size, total)
	to := min(from+size, total)
	for _, one := range rows[from:to] {
		numbers = append(numbers, one.Number)
	}
	return numbers, total, counts
}

// listed is the service, the controller and the view in one handler: the
// request is parsed by the listing, the builder applies it under the tenant,
// and the result goes to the table as it came back.
func listed(t *testing.T, program *sqliteProgram, listing *pagination.Listing) http.HandlerFunc {
	t.Helper()
	grant := auth.SystemGrant("invoice.list", "acme")
	return func(w http.ResponseWriter, r *http.Request) {
		state := listing.ParseURL(r.URL)
		builder := query.NewBuilder(program, grammars.NewSQLiteGrammar(), passThrough{}).From("invoices")
		found, result, err := builder.Listing(r.Context(), grant, state, "/listing")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		rows := make([]components.TableRow, len(found))
		for i, row := range found {
			rows[i] = components.TableRow{
				Key: fmt.Sprint(row["id"]), Label: "Select " + fmt.Sprint(row["number"]),
				Cells: []components.TableCell{
					{Text: fmt.Sprint(row["number"])}, {Text: fmt.Sprint(row["customer"])},
					{Text: fmt.Sprint(row["status"])}, {Text: money(int(row["total"].(int64)))},
				},
			}
		}
		body := string(components.DataTable(components.DataTableProps{
			ID: "invoices", Label: "Invoices", URL: "/listing",
			Searchable: true, FullPager: true,
			Result:  result,
			Filters: []components.DataTableFilter{{Key: "status", Label: "Status", Options: []components.SelectOption{{Label: "Paid", Value: "paid"}, {Label: "Open", Value: "open"}, {Label: "Void", Value: "void"}}}},
			Columns: []components.TableColumn{
				{Label: "Number", Key: "number", Sortable: true},
				{Label: "Customer", Key: "customer", Sortable: true, Hideable: true},
				{Label: "Status", Key: "status", Sortable: true, Hideable: true},
				{Label: "Total", Key: "total", Sortable: true, Hideable: true, Align: "end"},
			},
			Rows:  rows,
			Empty: components.EmptyProps{Title: "No invoices yet"},
		}))
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		if fragment(r) {
			fmt.Fprint(w, body)
			return
		}
		fmt.Fprint(w, layout("Invoices", body))
	}
}

// The database and the reference agree about every state the table asks for,
// and the other tenant's rows are never in an answer.
func TestTheBuilderAnswersEveryStateAsTheReferenceDoes(t *testing.T) {
	program := newSQLite(t)
	program.seed(t)
	listing := invoiceListing(t)
	grant := auth.SystemGrant("invoice.list", "acme")

	for _, raw := range []string{
		"", "q=ada", "q=A", "status=paid", "status=paid&status=open", "sort=total&dir=desc",
		"sort=customer&status=void&page=1", "page=3", "page=2&per_page=20", "q=inv-05&sort=total",
		"q=zzz", "q=100%25", "page=99", "sort=password&dir=sideways", "status=refunded",
	} {
		values, _ := url.ParseQuery(raw)
		builder := query.NewBuilder(program, grammars.NewSQLiteGrammar(), passThrough{}).From("invoices")
		rows, result, err := builder.Listing(context.Background(), grant, listing.Parse(values), "/listing")
		if err != nil {
			t.Fatalf("%q: %v", raw, err)
		}
		// The reference reads the state the listing read, so a value the
		// listing dropped is dropped from both.
		read := result.State.Values()
		want, total, counts := reference(read)
		var got []string
		for _, row := range rows {
			got = append(got, fmt.Sprint(row["number"]))
			if strings.HasPrefix(fmt.Sprint(row["number"]), "OTH") {
				t.Fatalf("%q: a row of another tenant was returned", raw)
			}
		}
		if wantPage := result.State.Page(); raw == "page=99" && wantPage != 6 {
			t.Errorf("%q: a page past the end was brought back to page %d, want 6", raw, wantPage)
		}
		if raw != "page=99" {
			if strings.Join(got, ",") != strings.Join(want, ",") || result.Total() != total {
				t.Errorf("%q: rows %v total %d, want %v total %d", raw, got, result.Total(), want, total)
			}
		}
		for _, count := range result.FacetCounts("status") {
			if count.Count != counts[count.Value] {
				t.Errorf("%q: %s counted %d, want %d", raw, count.Value, count.Count, counts[count.Value])
			}
		}
	}
}

// numbersOnPage are the invoice numbers the table shows, from the page.
func numbersOnPage(t *testing.T, page *browser.Page) string {
	t.Helper()
	return eval(t, page, `Array.from(document.querySelectorAll('#invoices tbody tr[data-key]')).map(function (row) { return row.querySelector('td[data-column=number]').textContent.trim(); }).join(',')`)
}

// TestTheListingFlowFromTheDatabaseToTheClick drives the flow the skeleton is
// built on: every address is the listing state's own, every row is the
// database's, and the counts beside a facet are the query's.
func TestTheListingFlowFromTheDatabaseToTheClick(t *testing.T) {
	program := newSQLite(t)
	program.seed(t)
	listing := invoiceListing(t)

	css, err := stylesheet()
	if err != nil {
		t.Logf("no stylesheet: %v", err)
	}
	var mu sync.Mutex
	var seen []string
	mux := http.NewServeMux()
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(css)
	})
	list := listed(t, program, listing)
	mux.HandleFunc("GET /listing", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		seen = append(seen, r.URL.RequestURI())
		mu.Unlock()
		list(w, r)
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(server.Close)

	page := browser.Open(t, drivers, browser.BaseURL(server.URL), browser.WithViewport(browser.Viewport{Width: 1440, Height: 900}))

	// expect is the page as the database answers the address it is at.
	expect := func(address string) {
		t.Helper()
		page.ExpectURL(address)
		values, _ := url.ParseQuery(strings.TrimPrefix(eval(t, page, `location.search`), "?"))
		wantRows, wantTotal, _ := reference(values)
		if got := numbersOnPage(t, page); got != strings.Join(wantRows, ",") {
			t.Fatalf("%s drew %s, the database holds %s", address, got, strings.Join(wantRows, ","))
		}
		if got := eval(t, page, totalExpression); got != strconv.Itoa(wantTotal) {
			t.Fatalf("%s says %s rows, the query counted %d", address, got, wantTotal)
		}
	}

	page.Visit("/listing")
	expect("/listing")
	page.ExpectText("60 row(s).")
	page.ExpectText("Page 1 of 6")

	// A facet ticked and unticked: the counts are the query's and the address
	// is the state's.
	page.Click(facetButton())
	page.ExpectText("Paid")
	if got := eval(t, page, `Array.from(document.querySelectorAll('#invoices-facet-status [data-facet-item]')).map(function (i) { return i.getAttribute('data-text') + ':' + i.lastElementChild.textContent.trim(); }).join(',')`); got != "Paid:24,Open:24,Void:12" {
		t.Fatalf("the facet counts are %s", got)
	}
	page.Find(facetBox("paid")).Check()
	expect("/listing?status=paid")
	page.ExpectText("24 row(s).")
	page.Find(facetBox("open")).Check()
	expect("/listing?status=paid&status=open")
	page.Find(facetBox("paid")).Uncheck()
	expect("/listing?status=open")
	page.Find(facetBox("open")).Uncheck()
	expect("/listing")
	page.Press("Escape")

	// An order, a page, and a size that keeps the first row on screen.
	page.ClickRole("link", "Total, not sorted")
	expect("/listing?sort=total")
	page.ClickRole("link", "Go to next page")
	expect("/listing?page=2&sort=total")
	page.ClickRole("link", "Go to next page")
	expect("/listing?page=3&sort=total")
	page.ExpectText("Page 3 of 6")
	page.Click(browser.CSS("#invoices-sizes > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemradio", "20")
	// Row 21 opened page three at ten; at twenty it is on page two.
	expect("/listing?page=2&per_page=20&sort=total")
	page.ExpectText("Page 2 of 3")

	// A search goes back to the first page and keeps the order and the size.
	page.FillIn(search(), "ada")
	until(t, page, `location.search`, "?per_page=20&q=ada&sort=total")
	expect("/listing?per_page=20&q=ada&sort=total")
	if got := numbersOnPage(t, page); strings.Contains(got, "OTH") {
		t.Fatalf("a row of another tenant was drawn: %s", got)
	}

	// A column hidden stays in the address, and the next step carries it.
	page.Click(browser.CSS("#invoices-columns > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemcheckbox", "Customer")
	page.ExpectURL("/listing?hide=customer&per_page=20&q=ada&sort=total")
	if got := eval(t, page, `String(document.querySelector('#invoices th[data-column=customer]').hidden)`); got != "true" {
		t.Fatalf("the customer column is drawn")
	}

	// And a reload of the address is the same page.
	address := eval(t, page, `location.pathname + location.search`)
	before := numbersOnPage(t, page)
	page.Visit(address)
	if got := numbersOnPage(t, page); got != before {
		t.Fatalf("a reload drew %s, want %s", got, before)
	}

	// Nothing the browser sent was a parameter the listing does not read, or a
	// value twice.
	mu.Lock()
	defer mu.Unlock()
	known := map[string]bool{"q": true, "sort": true, "dir": true, "status": true, "page": true, "per_page": true, "hide": true}
	for _, one := range seen {
		parsed, _ := url.Parse(one)
		for name, values := range parsed.Query() {
			if !known[name] {
				t.Errorf("request %s carries %s, which the listing does not read", one, name)
			}
			distinct := map[string]bool{}
			for _, value := range values {
				if distinct[value] {
					t.Errorf("request %s carries %s=%s twice", one, name, value)
				}
				distinct[value] = true
			}
		}
	}
}
