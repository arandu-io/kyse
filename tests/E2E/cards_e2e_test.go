//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"testing"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/kyse/components"
)

// cardsTable is a list that marks its columns for the card layout: the customer
// titles the card, the number is left out of it, and the status draws a badge.
func cardsTable(total int) string {
	rows := make([]components.TableRow, 0, 6)
	for _, one := range invoices()[:min(6, total)] {
		rows = append(rows, components.TableRow{
			Key: strconv.Itoa(one.ID), Label: "Select " + one.Number,
			Cells: []components.TableCell{
				{Text: one.Number}, {Text: one.Customer},
				{HTML: components.Badge(components.BadgeProps{Label: one.Status})},
				{Text: money(one.Total)}, {Text: ""},
			},
			Actions: []components.MenuItem{{Label: "View invoice", URL: "/customers/" + strconv.Itoa(one.ID)}, {Label: "Void", Variant: "destructive"}},
		})
	}
	return string(components.DataTable(components.DataTableProps{
		ID: "cards", Label: "Invoices", URL: "/cards", Searchable: true, FullPager: true,
		State:     components.DataTableValues{DefaultSize: 10},
		Paginator: pagination.NewLengthAwarePage(len(rows), total, 10, 1, pagination.Options{Path: "/cards"}),
		Columns: []components.TableColumn{
			{Label: "Number", Key: "number", Sortable: true, CardHidden: true},
			{Label: "Customer", Key: "customer", Sortable: true, Primary: true},
			{Label: "Status", Key: "status", Sortable: true},
			{Label: "Total", Key: "total", Sortable: true, Align: "end"},
			{Label: "", Key: ""},
		},
		Rows: rows,
	}))
}

// footTable is a list whose last column holds the row's own controls, the way an
// application draws "Edit" and "Remove" as icon buttons in a cell: the first and
// third rows have some, the second and fourth have none and their cell is empty.
func footTable() string {
	rows := make([]components.TableRow, 0, 4)
	for n, name := range []string{"Ada Lovelace", "Grace Hopper", "Alan Turing", "Edsger Dijkstra"} {
		var controls components.TableCell
		if n%2 == 0 {
			controls.HTML = template.HTML(`<div class="flex justify-end gap-0.5" data-row-actions>` +
				string(components.Button(components.ButtonProps{Label: "Edit", Variant: "ghost", Size: "sm"})) +
				string(components.Button(components.ButtonProps{Label: "Remove", Variant: "ghost", Size: "sm"})) + `</div>`)
		}
		rows = append(rows, components.TableRow{
			Key:   strconv.Itoa(n + 1),
			Cells: []components.TableCell{{Text: name}, {Text: []string{"Free", "Pro", "Pro", "Free"}[n]}, {Text: ""}, controls},
		})
	}
	return string(components.DataTable(components.DataTableProps{
		ID: "foot", Label: "Accounts", URL: "/foot",
		Paginator: pagination.NewLengthAwarePage(len(rows), len(rows), 10, 1, pagination.Options{Path: "/foot"}),
		Columns: []components.TableColumn{
			{Label: "Account", Key: "account", Primary: true},
			{Label: "Plan", Key: "plan"},
			{Label: "Note", Key: "note"},
			{Label: "Actions", Align: "end", Actions: true},
		},
		Rows: rows,
	}))
}

// style reads one computed property of the first element a selector finds, or
// of its pseudo-element.
func style(t *testing.T, page interface {
	Evaluate(string, any)
}, selector, pseudo, property string) string {
	t.Helper()
	var out any
	page.Evaluate(fmt.Sprintf(`(function () { var el = document.querySelector(%q); return el ? getComputedStyle(el, %q).getPropertyValue(%q) : 'missing'; })()`, selector, pseudo, property), &out)
	return fmt.Sprint(out)
}

func requireStyle(t *testing.T, got, want, what string) {
	t.Helper()
	if got != want {
		t.Errorf("%s is %q, want %q", what, got, want)
	}
}

// Under 600px every row is a card: the table and its rows stop being table
// boxes, a cell carries the name of its column, the primary column leads and
// the cells marked CardHidden are gone.
func TestARowIsACardOnAPhone(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 390)
	page.Visit("/cards")

	requireStyle(t, style(t, page, "#cards table", "", "display"), "block", "the table")
	requireStyle(t, style(t, page, "#cards tbody tr", "", "display"), "flex", "a row")
	requireStyle(t, style(t, page, "#cards tbody tr:last-child", "", "border-top-width"), "1px", "the border of the last card")
	// The title is the Customer, marked Primary, drawn first and with no name before it.
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=customer]", "", "order"), "-9999", "the primary cell's order")
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=customer]", "::before", "display"), "none", "the title's label")
	// Number is CardHidden.
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=number]", "", "display"), "none", "the CardHidden cell")
	// The other cells say their column first, from the data-label, and keep their text.
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=status]", "::before", "content"), `"Status"`, "the label of the status cell")
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=total]", "::before", "content"), `"Total"`, "the label of the total cell")
	requireStyle(t, eval(t, page, `document.querySelector('#cards tbody tr td[data-column=total]').textContent.trim()`), "$100.00", "the text of the total cell")
	// A cell with no column name draws no label.
	// The header keeps only what orders the list.
	requireStyle(t, style(t, page, "#cards thead tr", "", "display"), "flex", "the head row")
	if got := eval(t, page, `Array.from(document.querySelectorAll('#cards thead th')).filter(function (th) { return getComputedStyle(th).display !== 'none'; }).map(function (th) { return th.dataset.column; }).join(',')`); got != "number,customer,status,total" {
		t.Errorf("the head keeps %q, want the four headers that order the list", got)
	}
	// Nothing scrolls sideways and the menu of the row is in view, under its values.
	requireStyle(t, eval(t, page, `String(document.documentElement.scrollWidth <= document.documentElement.clientWidth)`), "true", "the page fitting its width")
	requireStyle(t, eval(t, page, `String((function () {
		var row = document.querySelector('#cards tbody tr'), menu = row.querySelector('[data-part=row-menu]').getBoundingClientRect(), box = row.getBoundingClientRect(), last = row.querySelector('td[data-column=total]').getBoundingClientRect();
		return menu.left >= 0 && menu.right <= innerWidth && menu.top >= last.bottom && menu.bottom <= box.bottom;
	})())`), "true", "the row menu sitting inside the card, under the values")
	shot(t, page, "cards-390-light")
}

// The cell of the column marked Actions is the foot of the card with no name in
// front of the controls, at the end of the card. A row with nothing in it draws
// no cell, so no "Actions" stands alone and no rule closes an empty foot; the
// label of any other column with no value is left out the same way.
func TestTheActionsOfACardAreItsFootAndAnEmptyOneDrawsNothing(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 390)
	page.Visit("/foot")

	const cell = `#foot tbody tr:nth-child(%d) td[data-label=Actions]`
	first := fmt.Sprintf(cell, 1)
	requireStyle(t, style(t, page, first, "", "display"), "flex", "the actions cell of a row that has some")
	requireStyle(t, style(t, page, first, "::before", "display"), "none", "the label of the actions cell")
	requireStyle(t, style(t, page, first, "", "justify-content"), "flex-end", "where the actions sit")
	requireStyle(t, style(t, page, first, "", "border-top-width"), "1px", "the rule over the actions")
	requireStyle(t, eval(t, page, `String((function () {
		var td = document.querySelector('`+first+`'), box = td.getBoundingClientRect(), group = td.querySelector('[data-row-actions]').getBoundingClientRect(), card = td.closest('tr').getBoundingClientRect();
		return Math.abs(group.right - box.right) < 2 && box.right <= card.right && group.left > box.left + box.width / 2;
	})())`), "true", "the controls at the end of the card, not after a label")
	// Rows two and four have none: nothing is drawn for them, no label and no rule.
	for _, n := range []int{2, 4} {
		requireStyle(t, style(t, page, fmt.Sprintf(cell, n), "", "display"), "none", fmt.Sprintf("the actions cell of row %d, which has none", n))
	}
	// Nor does the empty "Note" of any card draw its name in front of nothing.
	requireStyle(t, style(t, page, `#foot tbody tr:nth-child(1) td[data-column=note]`, "", "display"), "none", "an empty note")
	// No card draws the name of the actions column.
	requireStyle(t, eval(t, page, `String(Array.from(document.querySelectorAll('#foot tbody td[data-label=Actions]')).some(function (td) {
		return getComputedStyle(td).display !== 'none' && getComputedStyle(td, '::before').display !== 'none';
	}))`), "false", "a card drawing the name of the actions column")
	requireStyle(t, eval(t, page, `String(document.documentElement.scrollWidth <= document.documentElement.clientWidth)`), "true", "the page fitting its width")
	shot(t, page, "foot-390-light")
}

// A table above the breakpoint keeps the column of actions as a column, its
// header and its empty cells too.
func TestTheActionsColumnIsAColumnOnADesktop(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1280)
	page.Visit("/foot")
	requireStyle(t, style(t, page, "#foot tbody tr:nth-child(2) td[data-label=Actions]", "", "display"), "table-cell", "the empty actions cell")
	requireStyle(t, style(t, page, "#foot tbody tr:nth-child(1) td[data-label=Actions]", "", "border-top-width"), "0px", "a rule over the actions of a table row")
	requireStyle(t, style(t, page, "#foot thead th:last-child", "", "display"), "table-cell", "the header of the actions column")
	shot(t, page, "foot-1280-light")
}

// Above the breakpoint nothing changed: it is a table, with its labels unspoken.
func TestARowIsARowOnADesktop(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1280)
	page.Visit("/cards")

	requireStyle(t, style(t, page, "#cards table", "", "display"), "table", "the table")
	requireStyle(t, style(t, page, "#cards tbody tr", "", "display"), "table-row", "a row")
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=number]", "", "display"), "table-cell", "the CardHidden cell")
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=status]", "::before", "content"), "none", "the label of the status cell")
	requireStyle(t, style(t, page, "#cards tbody tr td[data-column=customer]", "", "order"), "0", "the primary cell's order")
	shot(t, page, "cards-1280-light")
}

// The popover of a row's menu is placed from the viewport with position:fixed,
// and the region is a container for the card layout: the container must not
// become its containing block.
func TestTheRowMenuOpensBesideItsButtonInsideTheContainer(t *testing.T) {
	for _, width := range []int{390, 1280} {
		t.Run(strconv.Itoa(width), func(t *testing.T) {
			app := serve(t, true)
			page := open(t, app, width)
			page.Visit("/cards")
			page.Evaluate(`document.querySelector('#cards tbody tr [data-part=row-menu] > summary').click()`, nil)
			page.Evaluate(`new Promise(function (resolve) { setTimeout(resolve, 200); })`, nil)
			got := eval(t, page, `(function () {
				var menu = document.querySelector('#cards tbody tr [data-part=row-menu]'), button = menu.querySelector('summary').getBoundingClientRect(), panel = menu.querySelector('[data-menu-panel]').getBoundingClientRect();
				return JSON.stringify({open: menu.open, below: Math.round(panel.top - button.bottom), inside: panel.left >= 0 && panel.right <= innerWidth && panel.bottom <= innerHeight, near: Math.abs(panel.right - button.right) < 40 || Math.abs(panel.left - button.left) < 40});
			})()`)
			var state struct {
				Open   bool `json:"open"`
				Below  int  `json:"below"`
				Inside bool `json:"inside"`
				Near   bool `json:"near"`
			}
			if err := json.Unmarshal([]byte(got), &state); err != nil {
				t.Fatal(err)
			}
			if !state.Open || !state.Inside || !state.Near || state.Below < -2 || state.Below > 24 {
				t.Fatalf("the menu opened at %s", got)
			}
			shot(t, page, "cards-"+strconv.Itoa(width)+"-menu")
		})
	}
}

// A list with one row says "1 row.", and the older sentence is the same words
// for any count.
func TestTheFooterSaysOneRowInTheSingular(t *testing.T) {
	app := serve(t, false)
	page := open(t, app, 390)
	page.Visit("/cards?total=1")
	page.ExpectText("1 row.")
	page.Visit("/cards")
	page.ExpectText("60 rows.")
}

// The sort bar still orders a list of cards, and the card layout survives the
// swap the order asks for.
func TestTheCardsSurviveASwap(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 390)
	page.Visit("/invoices")
	page.Evaluate(`document.querySelector('#invoices thead [data-sort-key=total]').click()`, nil)
	page.Evaluate(`new Promise(function (resolve) { setTimeout(resolve, 500); })`, nil)
	requireStyle(t, eval(t, page, `location.search`), "?dir=asc&sort=total", "the address")
	requireStyle(t, style(t, page, "#invoices tbody tr", "", "display"), "flex", "a row after the swap")
	shot(t, page, "invoices-390-sorted")
}

// The photographs of the card layout and of the table above it, in both
// schemes, written when KYSE_E2E_SHOTS names a directory. They are looked at by
// a person; nothing here asserts a pixel.
func TestPhotographTheCardsAndTheTable(t *testing.T) {
	if os.Getenv("KYSE_E2E_SHOTS") == "" {
		t.Skip("KYSE_E2E_SHOTS names no directory")
	}
	app := serve(t, true)
	for _, size := range []struct{ width, height int }{{390, 1300}, {1280, 800}} {
		for _, scheme := range []string{"light", "dark"} {
			for _, path := range []string{"cards", "invoices", "foot"} {
				page := browser.Open(t, drivers, browser.BaseURL(app.URL), browser.WithViewport(browser.Viewport{Width: size.width, Height: size.height}))
				page.Visit("/" + path)
				if scheme == "dark" {
					page.Evaluate(`document.documentElement.classList.add('dark')`, nil)
					// The colours transition; a frame taken at once is half way.
					page.Evaluate(`new Promise(function (resolve) { setTimeout(resolve, 600); })`, nil)
				}
				shot(t, page, fmt.Sprintf("%s-%d-%s", path, size.width, scheme))
			}
		}
	}
}

// wideTable is a list whose rows do not fit 600 pixels, over three pages: the
// case a table scrolls sideways in, which a phone no longer does.
func wideTable(page int) string {
	cells := func(n int) []components.TableCell {
		out := make([]components.TableCell, 6)
		for i := range out {
			out[i] = components.TableCell{Text: fmt.Sprintf("a value long enough to push a column wide %d-%d", n, i)}
		}
		return out
	}
	rows := make([]components.TableRow, 0, 3)
	for i := 0; i < 3; i++ {
		rows = append(rows, components.TableRow{Key: strconv.Itoa(i), Cells: cells(page*10 + i)})
	}
	columns := make([]components.TableColumn, 6)
	for i := range columns {
		columns[i] = components.TableColumn{Label: fmt.Sprintf("Column %d", i), Key: fmt.Sprintf("c%d", i)}
	}
	return string(components.DataTable(components.DataTableProps{
		ID: "wide", Label: "Wide", URL: "/wide", FullPager: true,
		State:     components.DataTableValues{DefaultSize: 3, CurrentPage: page},
		Paginator: pagination.NewLengthAwarePage(3, 9, 3, page, pagination.Options{Path: "/wide"}),
		Columns:   columns, Rows: rows,
	}))
}

// The width that counts is the region's, not the window's: the page keeps 24
// pixels of padding on each side, so a 640 pixel window gives the table 592 and
// it is cards, and a 700 pixel one gives it 652 and it is a table.
func TestTheRegionAndNotTheWindowDecidesBetweenCardsAndTable(t *testing.T) {
	for width, want := range map[int]string{640: "block", 700: "table"} {
		app := serve(t, true)
		page := open(t, app, width)
		page.Visit("/cards")
		requireStyle(t, style(t, page, "#cards table", "", "display"), want, fmt.Sprintf("the table in a %d pixel window", width))
	}
}
