//go:build e2e

package e2e

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/components"
)

// issued is the day an invoice was issued, in October 2026.
func issued(one invoice) string { return fmt.Sprintf("2026-10-%02d", 1+one.ID%28) }

// foldedCustomers have letters with marks, for the controls' text to be looked
// for without them.
func foldedName(one invoice) string {
	if one.ID%5 == 0 {
		return "Grâce Hopper"
	}
	return one.Customer
}

// controlledTable draws the invoices with a text, a choice and a period on the
// search's line, from the query, the way a controller built on the state does.
func controlledTable(query url.Values) string {
	state := components.DataTableValues{DefaultSize: 10, Facets: map[string][]string{}, Size: 10}
	for _, key := range []string{"customer", "state", "de", "ate"} {
		if value := strings.TrimSpace(query.Get(key)); value != "" {
			state.Facets[key] = []string{value}
		}
	}
	if key := query.Get("sort"); key == "total" || key == "number" {
		state.Sort = key
		if query.Get("dir") == "desc" {
			state.Dir = "desc"
		}
	}
	state.CurrentPage, _ = strconv.Atoi(query.Get("page"))

	fold := func(text string) string {
		return strings.NewReplacer("â", "a", "é", "e", "ç", "c").Replace(strings.ToLower(text))
	}
	var rows []invoice
	for _, one := range invoices() {
		if want := state.Chosen("customer"); len(want) > 0 && !strings.Contains(fold(foldedName(one)), fold(want[0])) {
			continue
		}
		if want := state.Chosen("state"); len(want) > 0 && one.Status != want[0] {
			continue
		}
		if from := state.Chosen("de"); len(from) > 0 && issued(one) < from[0] {
			continue
		}
		if to := state.Chosen("ate"); len(to) > 0 && issued(one) > to[0] {
			continue
		}
		rows = append(rows, one)
	}
	sort.SliceStable(rows, func(a, b int) bool {
		switch state.Sort {
		case "total":
			if state.SortDir() == "desc" {
				return rows[a].Total > rows[b].Total
			}
			return rows[a].Total < rows[b].Total
		case "number":
			if state.SortDir() == "desc" {
				return rows[a].ID > rows[b].ID
			}
		}
		return rows[a].ID < rows[b].ID
	})
	total := len(rows)
	pages := max(1, (total+state.Size-1)/state.Size)
	state.CurrentPage = min(max(1, state.CurrentPage), pages)
	from := (state.CurrentPage - 1) * state.Size
	page := rows[from:min(total, from+state.Size)]

	cells := make([]components.TableRow, len(page))
	for i, one := range page {
		cells[i] = components.TableRow{Key: strconv.Itoa(one.ID),
			Cells: []components.TableCell{{Text: one.Number}, {Text: foldedName(one)}, {Text: one.Status}, {Text: money(one.Total)}}}
	}
	return string(components.DataTable(components.DataTableProps{
		ID: "controlled", Label: "Controlled invoices", URL: "/controlled", FullPager: true,
		State:     state,
		Paginator: pagination.NewLengthAwarePage(len(page), total, state.Size, state.CurrentPage, pagination.Options{Path: "/controlled"}),
		Controls: []components.DataTableControl{
			{Key: "customer", Label: "Customer", Placeholder: "Name"},
			{Key: "state", Kind: components.DataTableSelect, Label: "State", Placeholder: "Any",
				Options: []components.SelectOption{{Value: "paid", Label: "Paid"}, {Value: "open", Label: "Open"}, {Value: "void", Label: "Void"}}},
			{Key: "de", ToKey: "ate", Kind: components.DataTablePeriod, Label: "Period", Period: components.DateRangePickerProps{
				Placeholder: "Pick a period",
				Presets: []components.DateRangePreset{
					{Label: "Today", Days: 1}, {Label: "Yesterday", Days: 1, EndsAgo: 1},
					{Label: "Last month", Unit: "month", Back: 1}, {Label: "Last year", Unit: "year", Back: 1},
				},
			}},
		},
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true},
			{Label: "Customer", Key: "customer"},
			{Label: "Status", Key: "status"},
			{Label: "Total", Key: "total", Sortable: true, Align: "end"},
		},
		Rows:  cells,
		Empty: components.EmptyProps{Title: "No invoices yet"},
	}))
}

func controlRows(t *testing.T, page *browser.Page) string {
	t.Helper()
	var out any
	page.Evaluate(`Array.from(document.querySelectorAll('#controlled tbody tr')).map(function (r) { return r.cells[1].textContent.trim(); }).join(',')`, &out)
	return fmt.Sprint(out)
}

// A text on the search's line is asked for after a pause, in the one request
// the form makes; a new order and a page keep it, and a name with a mark in it
// is found without the mark.
func TestATextControlIsOneRequestAndSurvivesTheOrderAndThePage(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/controlled")

	before := app.count("/controlled", true)
	page.FillIn(browser.CSS("#controlled-control-customer", "the customer text"), "grace")
	until(t, page, `String(new URL(location.href).searchParams.get('customer'))`, "grace")
	page.Evaluate(`new Promise(function (r) { setTimeout(r, 700); })`, nil)
	if got := app.count("/controlled", true) - before; got != 1 {
		t.Fatalf("a text typed and paused on made %d requests, want 1", got)
	}
	rows := controlRows(t, page)
	if rows == "" || strings.Contains(rows, "Ada") || !strings.Contains(rows, "Grâce Hopper") {
		t.Fatalf("grace finds %q, want the Graces and Grâce among them", rows)
	}
	if got := eval(t, page, `document.getElementById('controlled-control-customer').value`); got != "grace" {
		t.Fatalf("the box holds %q after the swap", got)
	}

	page.ClickRole("link", "Total, not sorted")
	until(t, page, `String(new URL(location.href).searchParams.get('sort'))`, "total")
	if got := addressQuery(t, page).Get("customer"); got != "grace" {
		t.Fatalf("a new order dropped the text: customer=%q", got)
	}
	page.ClickRole("link", "Reset")
	until(t, page, `String(new URL(location.href).searchParams.get('customer'))`, "null")
	if got := eval(t, page, `document.getElementById('controlled-control-customer').value`); got != "" {
		t.Fatalf("Reset left %q in the box", got)
	}
}

// A choice is one request when it changes, and it is the value in force when
// the region comes back.
func TestAChoiceControlIsAskedForWhenItChanges(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/controlled")
	page.Select("State", "Open")
	until(t, page, `String(new URL(location.href).searchParams.get('state'))`, "open")
	if got := eval(t, page, `document.getElementById('controlled-control-state').value`); got != "open" {
		t.Fatalf("the choice holds %q after the swap", got)
	}
	rows := eval(t, page, `Array.from(document.querySelectorAll('#controlled tbody tr')).map(function (r) { return r.cells[2].textContent.trim(); }).join(',')`)
	if strings.Contains(rows, "paid") || strings.Contains(rows, "void") || rows == "" {
		t.Fatalf("the open invoices are %q", rows)
	}
}

// A period shortcut that ends in the past is one request carrying both days.
func TestAPeriodShortcutThatEndsInThePastIsOneRequestWithBothDays(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 1440)
	page.Visit("/controlled")
	before := app.count("/controlled", true)

	page.Click(browser.CSS("#controlled-control-de-trigger", "the period field"))
	page.Click(browser.Text("Yesterday"))
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	until(t, page, `String(new URL(location.href).searchParams.get('de'))`, yesterday)
	query := addressQuery(t, page)
	if query.Get("de") != yesterday || query.Get("ate") != yesterday {
		t.Fatalf("yesterday sent de=%q ate=%q, want %s twice", query.Get("de"), query.Get("ate"), yesterday)
	}
	page.Evaluate(`new Promise(function (r) { setTimeout(r, 500); })`, nil)
	if got := app.count("/controlled", true) - before; got != 1 {
		t.Fatalf("a shortcut made %d requests, want 1", got)
	}
}
