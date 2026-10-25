//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/arandu-io/framework/test/browser"
)

const (
	rowsExpression      = `Array.from(document.querySelectorAll('#invoices tbody tr[data-key]')).map(function (row) { return row.cells[1].textContent.trim(); }).join(',')`
	totalExpression     = `String(document.querySelector('#invoices').getAttribute('data-total'))`
	countExpression     = `document.querySelector('#invoices [data-part=count]').textContent.trim()`
	customersExpression = `Array.from(document.querySelectorAll('#invoices tbody tr[data-key]')).map(function (row) { return row.cells[2].textContent.trim(); }).join(',')`
	pauseExpression     = `new Promise(function (resolve) { setTimeout(resolve, 700); })`
)

func search() browser.Locator {
	return browser.Role("searchbox", "Search Invoices")
}

func facetButton() browser.Locator {
	return browser.CSS("#invoices-facet-status > summary", "a summary has no role the locators name")
}

func facetBox(value string) browser.Locator {
	return browser.CSS("#invoices-facet-status input[value="+value+"]", "the facet box is a native checkbox inside a labelled row")
}

func menuOpen(t *testing.T, page *browser.Page, id string) bool {
	t.Helper()
	return eval(t, page, `String(document.getElementById('`+id+`').open)`) == "true"
}

func requireQuery(t *testing.T, page *browser.Page, want url.Values) {
	t.Helper()
	got := query(t, page)
	for _, values := range got {
		sort.Strings(values)
	}
	for _, values := range want {
		sort.Strings(values)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("the address query is %v, want %v", got, want)
	}
}

// A facet ticked and unticked leaves the address where it started, with every
// value once, and each tick is one request and keeps the popover open.
func TestAFacetIsTickedAndUntickedWithoutMultiplying(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.Click(facetButton())

	steps := []struct {
		box  string
		on   bool
		want url.Values
	}{
		{"paid", true, url.Values{"status": {"paid"}}},
		{"open", true, url.Values{"status": {"paid", "open"}}},
		{"paid", false, url.Values{"status": {"open"}}},
		{"open", false, url.Values{}},
	}
	for _, step := range steps {
		before := app.count("/invoices", true)
		if step.on {
			page.Find(facetBox(step.box)).Check()
		} else {
			page.Find(facetBox(step.box)).Uncheck()
		}
		requireQuery(t, page, step.want)
		if got := app.count("/invoices", true) - before; got != 1 {
			t.Fatalf("ticking %s=%v made %d requests, want 1", step.box, step.on, got)
		}
		if !menuOpen(t, page, "invoices-facet-status") {
			t.Fatalf("the facet popover closed after %s=%v", step.box, step.on)
		}
		// Every request the form made carried each value once.
		last := app.requests("/invoices", true)
		values, _ := url.ParseQuery(last[len(last)-1].Query)
		for name, list := range values {
			seen := map[string]bool{}
			for _, one := range list {
				if seen[one] {
					t.Fatalf("the request carried %s=%s twice: %s", name, one, last[len(last)-1].Query)
				}
				seen[one] = true
			}
		}
	}
	if got := eval(t, page, `location.search`); got != "" {
		t.Fatalf("with every facet unticked the address is %q, want none", got)
	}
}

// A word typed and paused on is one request; a word typed and followed by
// Enter is one request; and the box keeps what was typed and the focus.
func TestTheSearchMakesOneRequestWhetherItIsPausedOnOrSubmitted(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.Click(search())

	type_ := func(text string) {
		page.Evaluate(`(function () {
			var box = document.querySelector('#invoices-search');
			var text = `+fmt.Sprintf("%q", text)+`;
			for (var i = 1; i <= text.length; i++) {
				box.value = text.slice(0, i);
				box.dispatchEvent(new Event('input', { bubbles: true }));
			}
			return true;
		})()`, nil)
	}

	before := app.count("/invoices", true)
	type_("ada")
	page.ExpectURL("/invoices?q=ada")
	page.Evaluate(pauseExpression, nil)
	if got := app.count("/invoices", true) - before; got != 1 {
		t.Fatalf("a word typed and paused on made %d requests, want 1", got)
	}
	if got := eval(t, page, `document.activeElement && document.activeElement.id`); got != "invoices-search" {
		t.Fatalf("the search box lost the focus: %q", got)
	}
	if got := eval(t, page, `document.querySelector('#invoices-search').value`); got != "ada" {
		t.Fatalf("the search box holds %q, want ada", got)
	}
	for _, name := range strings.Split(eval(t, page, customersExpression), ",") {
		if !strings.Contains(strings.ToLower(name), "ada") {
			t.Fatalf("a row that does not match the search is listed: %s", name)
		}
	}

	// Typed as a person types, and submitted before the pause is over, with the
	// server slow enough that the pause ends while the request is in flight:
	// the box's change event and the pause must both stay out of it.
	app.slow(600 * time.Millisecond)
	page.FillIn(search(), "")
	page.ExpectURL("/invoices")
	app.slow(0)
	app.mu.Lock()
	app.log = nil
	app.mu.Unlock()
	app.slow(600 * time.Millisecond)
	page.Click(search())
	driver := page.Driver()
	if err := driver.InsertText(context.Background(), "gra"); err != nil {
		t.Fatal(err)
	}
	if err := driver.Press(context.Background(), "Enter"); err != nil {
		t.Fatal(err)
	}
	page.ExpectURL("/invoices?q=gra")
	page.Evaluate(pauseExpression, nil)
	if got := app.count("/invoices", true); got != 1 {
		t.Fatalf("a word typed and submitted made %d requests, want 1: %v", got, app.requests("/invoices", true))
	}
	app.slow(0)

	// Clearing it comes back to the bare address, and not to "?q=" that the
	// form serialised.
	page.FillIn(search(), "")
	until(t, page, `location.search`, "")
}

func TestOrderingShowsTheDirectionAndKeepsTheOtherControls(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices?q=a&status=paid")

	icon := func() string {
		return eval(t, page, `Array.from(document.querySelectorAll('#invoices th[data-column=total] [data-sort-icon]')).filter(function (one) { return !one.hidden; }).map(function (one) { return one.getAttribute('data-sort-icon'); }).join(',')`)
	}
	if got := icon(); got != "none" {
		t.Fatalf("an unordered column draws the %q icon, want none", got)
	}
	page.ClickRole("link", "Total, not sorted")
	page.ExpectURL("/invoices?dir=asc&q=a&sort=total&status=paid")
	if got := icon(); got != "asc" {
		t.Fatalf("an ascending column draws the %q icon, want asc", got)
	}
	page.ClickRole("link", "Total, sorted ascending")
	page.ExpectURL("/invoices?dir=desc&q=a&sort=total&status=paid")
	if got := icon(); got != "desc" {
		t.Fatalf("a descending column draws the %q icon, want desc", got)
	}
	page.ClickRole("link", "Total, sorted descending")
	page.ExpectURL("/invoices?dir=asc&q=a&sort=total&status=paid")
	if got := eval(t, page, `document.querySelector('#invoices th[data-column=total]').getAttribute('aria-sort')`); got != "ascending" {
		t.Fatalf("aria-sort = %q, want ascending", got)
	}
}

func TestPagingAndTheSizeKeepTheFirstRowOnScreen(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.ExpectText("Page 1 of 6")
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=2")
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=3")
	page.ExpectText("Page 3 of 6")
	page.Click(browser.CSS("#invoices-sizes > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemradio", "20")
	// Row 21 opened page three; with twenty a page it is on page two.
	page.ExpectURL("/invoices?page=2&per_page=20")
	page.ExpectText("Page 2 of 3")
	page.ClickRole("link", "Go to last page")
	page.ExpectURL("/invoices?page=3&per_page=20")
	page.ExpectText("Page 3 of 3")
	if got := eval(t, page, `String(document.querySelector('#invoices [data-step=next]').hasAttribute('href'))`); got != "false" {
		t.Fatalf("next on the last page is a link")
	}
	page.ClickRole("link", "Go to first page")
	page.ExpectURL("/invoices?per_page=20")
	// A search sends the reader back to the first page and keeps the size.
	page.FillIn(search(), "grace")
	page.ExpectURL("/invoices?per_page=20&q=grace")
}

func TestAHiddenColumnSurvivesTheSwapTheReloadAndTheMenuStaysOpen(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	hidden := func() string {
		return eval(t, page, `String(document.querySelector('#invoices th[data-column=customer]').hidden)`)
	}
	page.Click(browser.CSS("#invoices-columns > summary", "a summary has no role the locators name"))
	if got := eval(t, page, `Array.from(document.querySelectorAll('#invoices-columns [data-part=toggle]')).map(function (one) { return one.textContent.trim(); }).join(',')`); got != "Customer,Status,Total" {
		t.Fatalf("the columns menu lists %q, want the three hideable columns", got)
	}
	page.ClickRole("menuitemcheckbox", "Customer")
	page.ExpectURL("/invoices?hide=customer")
	if hidden() != "true" || !menuOpen(t, page, "invoices-columns") {
		t.Fatalf("the column is hidden=%s and the menu open=%v after the toggle", hidden(), menuOpen(t, page, "invoices-columns"))
	}
	// A swap that is not the menu's keeps the column off.
	page.ClickRole("link", "Number, not sorted")
	page.ExpectURL("/invoices?dir=asc&hide=customer&sort=number")
	if hidden() != "true" {
		t.Fatal("a sort brought the hidden column back")
	}
	// And so does a reload of the address.
	page.Visit("/invoices?dir=asc&hide=customer&sort=number")
	if hidden() != "true" {
		t.Fatal("a reload brought the hidden column back")
	}
	page.Click(browser.CSS("#invoices-columns > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemcheckbox", "Customer")
	page.ExpectURL("/invoices?dir=asc&sort=number")
	if hidden() != "false" {
		t.Fatal("showing the column again left it hidden")
	}
}

func TestTheSelectionSurvivesTheSearchTheOrderThePageAndPostsEveryKey(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	bar := browser.CSS("#invoices-bulk", "the bar is a form with no role of its own")
	page.ExpectHidden(bar)

	page.Find(browser.Role("checkbox", "Select INV-001")).Check()
	page.ExpectVisible(bar)
	page.ExpectText("1 of 60 row(s) selected.")
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=2")
	page.Find(browser.Role("checkbox", "Select INV-012")).Check()
	page.ExpectText("2 of 60 row(s) selected.")

	// The header ticks the page, and the box is indeterminate while only some of it is.
	state := func() string {
		return eval(t, page, `(function () { var m = document.querySelector('#invoices [data-select-all]'); return m.checked + '/' + m.indeterminate; })()`)
	}
	if got := state(); got != "false/true" {
		t.Fatalf("the header box is %s with one of ten ticked, want false/true", got)
	}
	page.Find(browser.Role("checkbox", "Select all rows")).Check()
	page.ExpectText("11 of 60 row(s) selected.")
	if got := state(); got != "true/false" {
		t.Fatalf("the header box is %s with the page ticked, want true/false", got)
	}

	// A search and an order change the rows and not the selection.
	page.FillIn(search(), "a")
	page.ExpectURL("/invoices?q=a")
	page.ClickRole("link", "Total, not sorted")
	page.ExpectText("11 of 60 row(s) selected.")
	page.Visit("/invoices")
	// A reload is a new page and starts over: the selection belongs to the page's life.
	page.ExpectText("0 of 60 row(s) selected.")

	page.Find(browser.Role("checkbox", "Select INV-002")).Check()
	page.ClickRole("link", "Go to next page")
	page.Find(browser.Role("checkbox", "Select INV-015")).Check()
	page.ClickRole("link", "Go to previous page")
	if got := eval(t, page, `String(document.querySelector('#invoices input[aria-label="Select INV-002"]').checked)`); got != "true" {
		t.Fatalf("the box ticked on page one is %s after coming back", got)
	}
	page.ClickRole("button", "Archive")
	page.ExpectURL("/invoices")
	posted := app.posted()
	if len(posted) != 1 {
		t.Fatalf("the bar posted %d times, want once", len(posted))
	}
	got := posted[0].Form["ids"]
	sort.Strings(got)
	if want := []string{"15", "2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("the bar posted ids %v, want %v", got, want)
	}
	if posted[0].Form.Get("do") != "archive" {
		t.Fatalf("the bar posted do=%q, want the pressed action", posted[0].Form.Get("do"))
	}
}

func TestTheAddressIsTheStateAndReloadAndHistoryRestoreIt(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=2")
	page.Click(facetButton())
	page.Find(facetBox("paid")).Check()
	page.ExpectURL("/invoices?status=paid")
	want := eval(t, page, rowsExpression)

	// Reload: the server draws the same list from the same address.
	page.Visit("/invoices?status=paid")
	if got := eval(t, page, rowsExpression); got != want {
		t.Fatalf("a reload drew %q, want %q", got, want)
	}
	if got := eval(t, page, `String(document.querySelector('#invoices-facet-status input[value=paid]').checked)`); got != "true" {
		t.Fatalf("the facet is not ticked after a reload")
	}
}

// Every step the swaps made is one entry of the history, so back and forward
// walk them in order.
func TestBackAndForwardWalkTheStepsTheSwapsMade(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	entries := eval(t, page, `String(history.length)`)
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=2")
	page.Click(facetButton())
	page.Find(facetBox("paid")).Check()
	page.ExpectURL("/invoices?status=paid")
	page.FillIn(search(), "ada")
	page.ExpectURL("/invoices?q=ada&status=paid")
	if got, was := eval(t, page, `String(history.length)`), entries; got != fmt.Sprint(atoi(was)+3) {
		t.Fatalf("three steps left %s entries after %s, want one entry each", got, was)
	}
	walk := func(move string, want string) {
		page.Evaluate(`history.`+move+`()`, nil)
		page.ExpectURL(want)
	}
	walk("back", "/invoices?status=paid")
	walk("back", "/invoices?page=2")
	walk("back", "/invoices")
	page.ExpectText("Page 1 of 6")
	walk("forward", "/invoices?page=2")
	page.ExpectText("Page 2 of 6")
	walk("forward", "/invoices?status=paid")
	walk("forward", "/invoices?q=ada&status=paid")
}

func atoi(s string) int {
	n := 0
	fmt.Sscan(s, &n)
	return n
}

func TestACompleteListNeverReachesTheServerEvenUnderBoost(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/small")
	if got := app.count("/small", false); got != 1 {
		t.Fatalf("the document made %d requests, want 1", got)
	}
	first := func() string {
		return eval(t, page, `Array.from(document.querySelectorAll('#small tbody tr[data-key]')).filter(function (r) { return !r.hidden; }).map(function (r) { return r.cells[1].textContent.trim(); }).join(',')`)
	}
	start := first()
	_ = start

	page.ClickRole("link", "Total, not sorted")
	if got := eval(t, page, `document.querySelector('#small th[data-column=total]').getAttribute('aria-sort')`); got != "ascending" {
		t.Fatalf("aria-sort = %q after the click, want ascending", got)
	}
	page.ClickRole("link", "Total, sorted ascending")
	if got := eval(t, page, `document.querySelector('#small th[data-column=total]').getAttribute('aria-sort')`); got != "descending" {
		t.Fatalf("aria-sort = %q after the second click, want descending", got)
	}
	if first() == start {
		t.Fatal("the rows were not reordered")
	}

	page.ClickRole("link", "Go to next page")
	page.ExpectText("Page 2 of 3")
	page.ClickRole("link", "Go to previous page")
	page.ExpectText("Page 1 of 3")
	page.Click(browser.CSS("#small-sizes > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemradio", "10")
	page.ExpectText("Page 1 of 2")

	page.FillIn(browser.Role("searchbox", "Search Small invoices"), "zzzz")
	page.ExpectText("No results.")
	page.FillIn(browser.Role("searchbox", "Search Small invoices"), "")

	page.Click(browser.CSS("#small-columns > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitemcheckbox", "Customer")
	if got := eval(t, page, `String(document.querySelector('#small th[data-column=customer]').hidden)`); got != "true" {
		t.Fatalf("the column is not hidden in the browser: %s", got)
	}
	if got := eval(t, page, `location.search`); got != "?hide=customer" {
		t.Fatalf("the address says %q, want the hidden column kept", got)
	}

	page.Evaluate(pauseExpression, nil)
	if got := app.count("/small", false); got != 1 {
		t.Fatalf("working a complete list made %d requests, want only the document", got)
	}
}

func TestAnOpenSelectionDoesNotLeakIntoAnotherPage(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.Find(browser.Role("checkbox", "Select INV-001")).Check()
	page.ExpectText("1 of 60 row(s) selected.")
	// A boosted visit to another page and back is a new list.
	page.ClickRole("link", "Other page")
	page.ExpectText("Somewhere else.")
	page.ClickRole("link", "Invoices")
	page.ExpectText("0 of 60 row(s) selected.")
}

func TestTheHeaderMenuOrdersAndHides(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/menu")
	page.Click(browser.CSS("#invoices-head-total > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitem", "Desc")
	page.ExpectURL("/menu?dir=desc&sort=total")
	if got := eval(t, page, `Array.from(document.querySelectorAll('#invoices th[data-column=total] summary [data-sort-icon]')).filter(function (one) { return !one.hidden; }).map(function (one) { return one.getAttribute('data-sort-icon'); }).join(',')`); got != "desc" {
		t.Fatalf("the header shows the %q icon, want desc", got)
	}
	page.Click(browser.CSS("#invoices-head-total > summary", "a summary has no role the locators name"))
	page.ClickRole("menuitem", "Hide")
	page.ExpectURL("/menu?dir=desc&hide=total&sort=total")
	if got := eval(t, page, `String(document.querySelector('#invoices th[data-column=total]').hidden)`); got != "true" {
		t.Fatalf("Hide left the column drawn: %s", got)
	}
}

func TestARowMenuOpensBesideItsButtonAndHighlightsTheRow(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	row := `#invoices-row-3`
	page.Click(browser.CSS(row+" > summary", "a summary has no role the locators name"))
	if !menuOpen(t, page, "invoices-row-3") {
		t.Fatal("the row menu did not open")
	}
	if got := eval(t, page, `(function () { var p = document.querySelector('`+row+` > [data-menu-panel]'); var b = document.querySelector('`+row+` > summary').getBoundingClientRect(); var r = p.getBoundingClientRect(); return (r.top >= b.bottom - 1 || r.bottom <= b.top + 1) && r.right <= window.innerWidth; })()`); got != "true" {
		t.Fatalf("the panel is not placed beside its button: %s", got)
	}
	page.ClickRole("menuitem", "View customer")
	page.ExpectText("A customer.")
}

func TestEscapeAndAClickElsewhereCloseAMenu(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.Click(facetButton())
	if !menuOpen(t, page, "invoices-facet-status") {
		t.Fatal("the facet did not open")
	}
	page.Press("Escape")
	if menuOpen(t, page, "invoices-facet-status") {
		t.Fatal("Escape left the popover open")
	}
	page.Click(facetButton())
	page.Click(browser.CSS("#invoices-columns > summary", "a summary has no role the locators name"))
	if menuOpen(t, page, "invoices-facet-status") || !menuOpen(t, page, "invoices-columns") {
		t.Fatal("opening one menu did not close the other")
	}
	page.Click(browser.CSS("main h1", "the heading is a place that is neither a menu nor a control"))
	if menuOpen(t, page, "invoices-columns") {
		t.Fatal("a click elsewhere left the menu open")
	}
}

func TestThePageIsNeverWiderThanTheScreen(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1440)
	page.Visit("/invoices")
	page.Responsive()
	page.Visit("/invoices?status=open&status=paid&q=a&hide=customer")
	page.Responsive()
	page.Visit("/small")
	page.Responsive()
}

// until waits for an expression to read want, for the answer to a request the
// page makes by itself.
func until(t *testing.T, page *browser.Page, expression, want string) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for {
		got := eval(t, page, expression)
		if got == want {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("%s reads %q, want %q", expression, got, want)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

// screen puts the page at a width in a colour scheme, the way a person's
// device is.
func screen(t *testing.T, page *browser.Page, width int, scheme string) {
	t.Helper()
	height := map[int]int{1440: 900, 768: 1024, 390: 844}[width]
	ctx := context.Background()
	if err := page.Driver().SetViewport(ctx, browser.Viewport{Width: width, Height: height}); err != nil {
		t.Fatal(err)
	}
	if err := page.Driver().EmulateMedia(ctx, browser.Media{ColorScheme: scheme}); err != nil {
		t.Fatal(err)
	}
	page.Evaluate(`document.documentElement.classList.toggle('dark', `+fmt.Sprint(scheme == "dark")+`)`, nil)
}

// TestPhotographs writes the component in its states, at the three widths in
// both schemes, to the directory KYSE_E2E_SHOTS names, so it can be set beside
// the reference. It asserts nothing and does nothing without the variable.
func TestPhotographs(t *testing.T) {
	if os.Getenv("KYSE_E2E_SHOTS") == "" {
		t.Skip("KYSE_E2E_SHOTS names no directory")
	}
	app := serve(t, true)
	page := open(t, app, 1440)
	summary := func(id string) browser.Locator {
		return browser.CSS("#"+id+" > summary", "a summary has no role the locators name")
	}
	for _, width := range []int{1440, 768, 390} {
		for _, scheme := range []string{"light", "dark"} {
			tag := fmt.Sprintf("%d-%s", width, scheme)
			at := func(address string) {
				page.Visit(address)
				screen(t, page, width, scheme)
			}
			at("/invoices")
			shot(t, page, tag+"-1-default")

			at("/invoices?dir=desc&q=a&sort=total&status=paid")
			shot(t, page, tag+"-2-sorted-filtered")

			at("/invoices?status=paid")
			page.Click(summary("invoices-facet-status"))
			shot(t, page, tag+"-3-facet-popover")

			at("/invoices")
			page.Find(browser.Role("checkbox", "Select INV-001")).Check()
			page.Find(browser.Role("checkbox", "Select INV-002")).Check()
			shot(t, page, tag+"-4-selected")

			at("/invoices")
			page.Click(summary("invoices-columns"))
			shot(t, page, tag+"-5-columns-menu")

			at("/invoices")
			page.Click(summary("invoices-row-3"))
			shot(t, page, tag+"-6-row-menu")

			at("/menu?dir=asc&sort=customer")
			page.Click(summary("invoices-head-customer"))
			shot(t, page, tag+"-7-header-menu")

			at("/small")
			shot(t, page, tag+"-8-complete")

			at("/invoices?q=zzz")
			shot(t, page, tag+"-9-no-results")
		}
	}
}

// A swap that replaces the region keeps what the reader had under their hands:
// how far the table was scrolled sideways, and a menu that was open.
func TestASwapKeepsTheSidewaysScroll(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 390)
	page.Visit("/invoices")
	scroller := `document.querySelector('#invoices [data-part=table]').parentElement`
	page.Evaluate(`(function () { var s = `+scroller+`; s.scrollLeft = 80; return s.scrollLeft; })()`, nil)
	if got := eval(t, page, scroller+`.scrollLeft`); got != "80" {
		t.Fatalf("the table does not scroll sideways at 390 pixels: %s", got)
	}
	page.ClickRole("link", "Go to next page")
	page.ExpectURL("/invoices?page=2")
	if got := eval(t, page, scroller+`.scrollLeft`); got != "80" {
		t.Fatalf("the scroll is %s after the swap, want the 80 it had", got)
	}
}
