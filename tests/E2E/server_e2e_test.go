//go:build e2e

// Package e2e drives the data table in a real browser, against a server that
// answers the way an application's listing does: it parses the state out of
// the query, narrows a list with it and draws the component again.
//
//	go test -tags e2e ./tests/E2E/...
//
// The page is served under the policy an application serves: script-src and
// style-src 'self', so an inline style or an inline handler fails the test as
// a violation. The stylesheet is compiled with the standalone tailwindcss
// binary `aru view:build` downloads, from the component sources, so what is
// measured and photographed is what the component draws; a machine without the
// binary skips the tests that need pixels and runs the rest.
package e2e

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/framework/test/browser/cdp"
	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

var drivers = browser.Registry{"cdp": cdp.Launch}

// invoice is one row of the list the server holds.
type invoice struct {
	ID       int
	Number   string
	Customer string
	Status   string
	Total    int
}

var customers = []string{"Ada Lovelace", "Grace Hopper", "Alan Turing", "Edsger Dijkstra", "Barbara Liskov"}

func invoices() []invoice {
	statuses := []string{"paid", "open", "paid", "void", "open"}
	out := make([]invoice, 60)
	for i := range out {
		out[i] = invoice{
			ID: i + 1, Number: fmt.Sprintf("INV-%03d", i+1),
			Customer: customers[i%len(customers)], Status: statuses[(i/2)%len(statuses)],
			Total: 100 + (i*37)%900,
		}
	}
	return out
}

// request is one request the server saw.
type request struct {
	Path   string
	Query  string
	Htmx   bool
	Method string
	Form   url.Values
}

// app is the server under test and the log of what reached it.
type app struct {
	*httptest.Server
	mu  sync.Mutex
	log []request
	// delay is how long the listing takes to answer, for a test that needs a
	// request still in flight when the next thing happens.
	delay time.Duration
}

func (a *app) requests(path string, htmxOnly bool) []request {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []request
	for _, r := range a.log {
		if r.Path == path && (!htmxOnly || r.Htmx) {
			out = append(out, r)
		}
	}
	return out
}

func (a *app) slow(delay time.Duration) {
	a.mu.Lock()
	a.delay = delay
	a.mu.Unlock()
}

func (a *app) count(path string, htmxOnly bool) int { return len(a.requests(path, htmxOnly)) }

func (a *app) posted() []request {
	a.mu.Lock()
	defer a.mu.Unlock()
	var out []request
	for _, r := range a.log {
		if r.Method == http.MethodPost {
			out = append(out, r)
		}
	}
	return out
}

// listing narrows the invoices by the query, the way a handler built on the
// state does: a whitelist for every value it reads.
func listing(query url.Values) (components.DataTableValues, []invoice, int, map[string]int) {
	allowed := func(value string, among ...string) bool {
		for _, one := range among {
			if one == value {
				return true
			}
		}
		return false
	}
	state := components.DataTableValues{DefaultSize: 10, Facets: map[string][]string{}}
	state.Q = strings.TrimSpace(query.Get("q"))
	for _, status := range query["status"] {
		if allowed(status, "paid", "open", "void") {
			state.Facets["status"] = append(state.Facets["status"], status)
		}
	}
	if key := query.Get("sort"); allowed(key, "number", "customer", "status", "total") {
		state.Sort = key
		if query.Get("dir") == "desc" {
			state.Dir = "desc"
		}
	}
	state.Size = 10
	if size, _ := strconv.Atoi(query.Get("per_page")); allowed(strconv.Itoa(size), "10", "20", "50") {
		state.Size = size
	}
	for _, key := range query["hide"] {
		if allowed(key, "customer", "status", "total") {
			state.Hide = append(state.Hide, key)
		}
	}
	state.CurrentPage, _ = strconv.Atoi(query.Get("page"))

	needle := strings.ToLower(state.Q)
	search := func(one invoice) bool {
		return needle == "" || strings.Contains(strings.ToLower(one.Number+" "+one.Customer), needle)
	}
	counts := map[string]int{}
	var rows []invoice
	chosen := state.Chosen("status")
	for _, one := range invoices() {
		if !search(one) {
			continue
		}
		counts[one.Status]++
		if len(chosen) > 0 {
			match := false
			for _, status := range chosen {
				match = match || status == one.Status
			}
			if !match {
				continue
			}
		}
		rows = append(rows, one)
	}
	if state.Sort != "" {
		sign := 1
		if state.SortDir() == "desc" {
			sign = -1
		}
		sort.SliceStable(rows, func(a, b int) bool {
			var order int
			if state.Sort == "total" {
				order = rows[a].Total - rows[b].Total
			} else {
				order = strings.Compare(keyOf(rows[a], state.Sort), keyOf(rows[b], state.Sort))
			}
			if order == 0 {
				return rows[a].ID < rows[b].ID
			}
			return order*sign < 0
		})
	}
	total := len(rows)
	pages := max(1, (total+state.Size-1)/state.Size)
	state.CurrentPage = min(max(1, state.CurrentPage), pages)
	from := (state.CurrentPage - 1) * state.Size
	to := min(total, from+state.Size)
	return state, rows[from:to], total, counts
}

func keyOf(one invoice, key string) string {
	switch key {
	case "customer":
		return one.Customer
	case "status":
		return one.Status
	case "total":
		return strconv.Itoa(one.Total)
	}
	return one.Number
}

func money(total int) string { return fmt.Sprintf("$%d.00", total) }

// table draws the invoices list the way a controller would, from the query.
func table(path string, query url.Values, menu bool) string {
	state, page, total, counts := listing(query)
	rows := make([]components.TableRow, len(page))
	for i, one := range page {
		rows[i] = components.TableRow{
			Key: strconv.Itoa(one.ID), Label: "Select " + one.Number,
			Cells: []components.TableCell{{Text: one.Number}, {Text: one.Customer}, {Text: one.Status}, {Text: money(one.Total)}},
			Actions: []components.MenuItem{
				{Label: "Actions", Heading: true},
				{Label: "Copy invoice ID"},
				{Separator: true},
				{Label: "View customer", URL: "/customers/" + strconv.Itoa(one.ID)},
				{Label: "View invoice details", URL: path + "/" + strconv.Itoa(one.ID)},
			},
		}
	}
	options := []components.SelectOption{{Label: "Paid", Value: "paid"}, {Label: "Open", Value: "open"}, {Label: "Void", Value: "void"}}
	return string(components.DataTable(components.DataTableProps{
		ID: "invoices", Label: "Invoices", URL: path,
		Searchable: true, FullPager: true, HeaderMenu: menu, PageSizes: []int{10, 20, 50},
		State:     state,
		Paginator: pagination.NewLengthAwarePage(len(page), total, state.Size, state.CurrentPage, pagination.Options{Path: path}),
		Filters:   []components.DataTableFilter{{Key: "status", Label: "Status", Options: options, Counts: counts}},
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true},
			{Label: "Customer", Key: "customer", Sortable: true, Hideable: true},
			{Label: "Status", Key: "status", Sortable: true, Hideable: true},
			{Label: "Total", Key: "total", Sortable: true, Hideable: true, Align: "end"},
		},
		Rows:       rows,
		SelectName: "ids", BulkAction: "/invoices/bulk",
		BulkActions: []components.ButtonProps{
			{Label: "Archive", Type: "submit", Variant: "outline", Size: "sm", Name: "do", Value: "archive"},
			{Label: "Delete", Type: "submit", Variant: "destructive", Size: "sm", Name: "do", Value: "delete"},
		},
		Empty: components.EmptyProps{Title: "No invoices yet", Message: "Invoices you send appear here."},
	}))
}

// small draws a list that holds all its rows.
func small() string {
	rows := make([]components.TableRow, 0, 12)
	for i, one := range invoices()[:12] {
		rows = append(rows, components.TableRow{
			Key: strconv.Itoa(one.ID), Label: "Select " + one.Number,
			Facets: map[string]string{"status": one.Status},
			Cells: []components.TableCell{
				{Text: one.Number}, {Text: one.Customer}, {Text: one.Status},
				{Text: money(one.Total), SortValue: strconv.Itoa(one.Total)},
			},
		})
		_ = i
	}
	options := []components.SelectOption{{Label: "Paid", Value: "paid"}, {Label: "Open", Value: "open"}, {Label: "Void", Value: "void"}}
	return string(components.DataTable(components.DataTableProps{
		ID: "small", Label: "Small invoices", URL: "/small", Complete: true, Searchable: true, FullPager: true,
		PageSize: 5, PageSizes: []int{5, 10},
		State:   components.DataTableValues{DefaultSize: 5},
		Filters: []components.DataTableFilter{{Key: "status", Label: "Status", Options: options}},
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true},
			{Label: "Customer", Key: "customer", Sortable: true, Hideable: true},
			{Label: "Status", Key: "status", Sortable: true},
			{Label: "Total", Key: "total", Sortable: true, Align: "end"},
		},
		Rows: rows,
	}))
}

func layout(title, body string) string {
	return fmt.Sprintf(`<!doctype html><html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1"><title>%s</title>
<meta name="htmx-config" content='{"includeIndicatorStyles":false,"allowEval":false}'>
<link rel="stylesheet" href="/app.css">
<script src="%s"></script><script src="%s"></script></head>
<body hx-boost="true" class="bg-background text-foreground antialiased">
<nav class="flex gap-4 border-b p-4 text-sm"><a href="/invoices">Invoices</a><a href="/menu">Menu header</a><a href="/small">Small list</a><a href="/other">Other page</a></nav>
<main class="mx-auto max-w-5xl p-6"><h1 class="mb-2 text-2xl font-semibold tracking-tight">%s</h1>%s</main>
</body></html>`, title, view.Asset("htmx.min.js"), view.Asset("ui.js"), title, body)
}

// fragment is whether the request wants the region alone: an htmx request that
// is not the history restoring a page.
func fragment(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true" && r.Header.Get("HX-History-Restore-Request") != "true"
}

var (
	cssOnce sync.Once
	cssBody []byte
	cssErr  error
)

// stylesheet compiles the components' classes with the tailwindcss binary.
func stylesheet() ([]byte, error) {
	cssOnce.Do(func() {
		binary := os.Getenv("KYSE_E2E_TAILWIND")
		if binary == "" {
			home, _ := os.UserHomeDir()
			matches, _ := filepath.Glob(filepath.Join(home, ".arandu", "bin", "tailwindcss-v*"))
			if len(matches) > 0 {
				binary = matches[len(matches)-1]
			}
		}
		if binary == "" {
			cssErr = fmt.Errorf("no tailwindcss binary: set KYSE_E2E_TAILWIND, or run aru view:build once so it is cached in ~/.arandu/bin")
			return
		}
		_, file, _, _ := runtime.Caller(0)
		root := filepath.Join(filepath.Dir(file), "..", "..")
		dir, err := os.MkdirTemp("", "kyse-e2e-css")
		if err != nil {
			cssErr = err
			return
		}
		input := filepath.Join(dir, "in.css")
		source := fmt.Sprintf(`@import "tailwindcss" source(none);
@source %q;
@source %q;
@custom-variant dark (&:is(.dark *));
@theme inline {
  --color-background: var(--background);
  --color-foreground: var(--foreground);
  --color-card: var(--card);
  --color-popover: var(--popover);
  --color-popover-foreground: var(--popover-foreground);
  --color-primary: var(--primary);
  --color-primary-foreground: var(--primary-foreground);
  --color-secondary: var(--secondary);
  --color-secondary-foreground: var(--secondary-foreground);
  --color-muted: var(--muted);
  --color-muted-foreground: var(--muted-foreground);
  --color-accent: var(--accent);
  --color-accent-foreground: var(--accent-foreground);
  --color-destructive: var(--destructive);
  --color-border: var(--border);
  --color-input: var(--input);
  --color-ring: var(--ring);
}
:root {
  --radius: 0.625rem; --radius-md: calc(var(--radius) - 2px);
  --background: oklch(1 0 0); --foreground: oklch(0.145 0 0); --card: oklch(1 0 0);
  --popover: oklch(1 0 0); --popover-foreground: oklch(0.145 0 0);
  --primary: oklch(0.205 0 0); --primary-foreground: oklch(0.985 0 0);
  --secondary: oklch(0.97 0 0); --secondary-foreground: oklch(0.205 0 0);
  --muted: oklch(0.97 0 0); --muted-foreground: oklch(0.556 0 0);
  --accent: oklch(0.97 0 0); --accent-foreground: oklch(0.205 0 0);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.922 0 0); --input: oklch(0.922 0 0); --ring: oklch(0.708 0 0);
}
.dark {
  --background: oklch(0.145 0 0); --foreground: oklch(0.985 0 0); --card: oklch(0.205 0 0);
  --popover: oklch(0.205 0 0); --popover-foreground: oklch(0.985 0 0);
  --primary: oklch(0.922 0 0); --primary-foreground: oklch(0.205 0 0);
  --secondary: oklch(0.269 0 0); --secondary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.269 0 0); --muted-foreground: oklch(0.708 0 0);
  --accent: oklch(0.269 0 0); --accent-foreground: oklch(0.985 0 0);
  --destructive: oklch(0.704 0.191 22.216);
  --border: oklch(1 0 0 / 10%%); --input: oklch(1 0 0 / 15%%); --ring: oklch(0.556 0 0);
}
@layer base {
  *, ::after, ::before { border-color: var(--border); }
  body { background: var(--background); color: var(--foreground); }
}
`, filepath.Join(root, "components")+"/**/*.go", filepath.Dir(file)+"/**/*.go")
		if err := os.WriteFile(input, []byte(source), 0o644); err != nil {
			cssErr = err
			return
		}
		output := filepath.Join(dir, "out.css")
		if out, err := exec.Command(binary, "-i", input, "-o", output).CombinedOutput(); err != nil {
			cssErr = fmt.Errorf("tailwindcss: %v\n%s", err, out)
			return
		}
		cssBody, cssErr = os.ReadFile(output)
	})
	return cssBody, cssErr
}

// serve starts the server under test. A test that needs pixels asks for the
// stylesheet; one that does not still gets a working page.
func serve(t *testing.T, needStyles bool) *app {
	t.Helper()
	css, err := stylesheet()
	if err != nil && needStyles {
		t.Skipf("the stylesheet cannot be compiled here: %v", err)
	}
	a := &app{}
	mux := http.NewServeMux()
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(css)
	})
	page := func(path, title string, menu bool) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			a.record(r)
			a.mu.Lock()
			wait := a.delay
			a.mu.Unlock()
			time.Sleep(wait)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Vary", "HX-Request, HX-History-Restore-Request")
			w.Header().Set("Cache-Control", "no-store")
			body := table(path, r.URL.Query(), menu)
			if fragment(r) {
				fmt.Fprint(w, body)
				return
			}
			fmt.Fprint(w, layout(title, body))
		}
	}
	mux.HandleFunc("GET /invoices", page("/invoices", "Invoices", false))
	mux.HandleFunc("GET /menu", page("/menu", "Invoices", true))
	mux.HandleFunc("GET /small", func(w http.ResponseWriter, r *http.Request) {
		a.record(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, layout("Small invoices", small()))
	})
	mux.HandleFunc("GET /other", func(w http.ResponseWriter, r *http.Request) {
		a.record(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, layout("Other page", `<p>Somewhere else.</p>`))
	})
	mux.HandleFunc("POST /invoices/bulk", func(w http.ResponseWriter, r *http.Request) {
		a.record(r)
		http.Redirect(w, r, "/invoices", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /customers/{id}", func(w http.ResponseWriter, r *http.Request) {
		a.record(r)
		fmt.Fprint(w, layout("Customer", `<p>A customer.</p>`))
	})
	a.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(a.Close)
	return a
}

func (a *app) record(r *http.Request) {
	r.ParseForm()
	a.mu.Lock()
	defer a.mu.Unlock()
	a.log = append(a.log, request{
		Path: r.URL.Path, Query: r.URL.RawQuery, Htmx: r.Header.Get("HX-Request") == "true",
		Method: r.Method, Form: r.PostForm,
	})
}

// open starts a browser on the server at the width given.
func open(t *testing.T, a *app, width int) *browser.Page {
	t.Helper()
	height := 900
	if width < 600 {
		height = 844
	}
	return browser.Open(t, drivers, browser.BaseURL(a.URL), browser.WithViewport(browser.Viewport{Width: width, Height: height}))
}

// eval reads a string out of the page.
func eval(t *testing.T, page *browser.Page, expression string) string {
	t.Helper()
	var out any
	page.Evaluate(expression, &out)
	return fmt.Sprint(out)
}

// addressQuery is the query string of the address the page is on.
func addressQuery(t *testing.T, page *browser.Page) url.Values {
	t.Helper()
	values, err := url.ParseQuery(strings.TrimPrefix(eval(t, page, `location.search`), "?"))
	if err != nil {
		t.Fatalf("the address does not parse: %v", err)
	}
	return values
}

// shot saves the viewport as a PNG when KYSE_E2E_SHOTS names a directory.
func shot(t *testing.T, page *browser.Page, name string) {
	t.Helper()
	dir := os.Getenv("KYSE_E2E_SHOTS")
	if dir == "" {
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	png, err := page.Driver().Screenshot(context.Background())
	if err != nil {
		t.Fatalf("screenshot %s: %v", name, err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".png"), png, 0o644); err != nil {
		t.Fatal(err)
	}
}
