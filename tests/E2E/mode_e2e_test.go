//go:build e2e

package e2e

import (
	"fmt"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// modeList is a short list inside a narrow box, the way a dashboard card holds
// one: two lines, each a camera and a count, drawn in the mode given.
func modeList(id string, mode components.DataTableMode) string {
	rows := []components.TableRow{
		{Key: "1", Cells: []components.TableCell{{Text: "Camera 1"}, {Text: "12"}}},
		{Key: "2", Cells: []components.TableCell{{Text: "Camera 2"}, {Text: "34"}}},
	}
	return string(components.DataTable(components.DataTableProps{
		ID: id, Label: id, URL: "/mode", Complete: true, Mode: mode, RowTemplate: true,
		Columns: []components.TableColumn{
			{Label: "Camera", Key: "camera"},
			{Label: "Last hour", Key: "hour", Align: "end"},
		},
		Rows: rows,
	}))
}

// modeBoxes is the three modes, each in a box 400px wide, and the same list in
// a box 900px wide in the mode that makes cards at every width. The box is a
// class the stylesheet is compiled with from this file.
func modeBoxes() string {
	box := func(width, id string, mode components.DataTableMode) string {
		return fmt.Sprintf(`<section class="mb-6 %s border p-2"><h2 class="mb-2 text-sm font-medium">%s</h2>%s</section>`, width, id, modeList(id, mode))
	}
	return box("w-[400px]", "m-auto", components.DataTableModeAuto) +
		box("w-[400px]", "m-table", components.DataTableModeTable) +
		box("w-[400px]", "m-cards", components.DataTableModeCards) +
		box("w-[900px]", "m-wide-cards", components.DataTableModeCards)
}

// live puts a row into the list the way a script does as a read arrives: it
// clones the row the server drew in a template and writes the cells.
func live(t *testing.T, page interface{ Evaluate(string, any) }, id string, cells ...string) {
	t.Helper()
	var out any
	page.Evaluate(fmt.Sprintf(`(function () {
		var root = document.getElementById(%q), tpl = root.querySelector('template[data-part=row-template]');
		var tr = tpl.content.firstElementChild.cloneNode(true), values = %s;
		tr.querySelectorAll('td').forEach(function (td, i) { td.textContent = values[i]; });
		var body = root.querySelector('tbody');
		body.insertBefore(tr, body.firstChild);
		return 'ok';
	})()`, id, jsList(cells)), &out)
}

func jsList(values []string) string {
	quoted := make([]string, len(values))
	for i, v := range values {
		quoted[i] = fmt.Sprintf("%q", v)
	}
	return "[" + strings.Join(quoted, ",") + "]"
}

// The window decides and not the box: a box of 400px on a desktop window is a
// table in every mode but cards, with no mark on the list, and on a phone window
// the default is cards and the table mode keeps the table.
func TestTheWindowDecidesAndTheTableModeKeepsTheTableOnAPhone(t *testing.T) {
	for _, width := range []int{1280, 390} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			app := serve(t, true)
			page := open(t, app, width)
			page.Visit("/mode")

			requireStyle(t, style(t, page, "#m-table table", "", "display"), "table", "the table in the table mode")
			requireStyle(t, style(t, page, "#m-table tbody tr", "", "display"), "table-row", "a row in the table mode")
			requireStyle(t, style(t, page, "#m-table tbody td", "", "display"), "table-cell", "a cell in the table mode")
			requireStyle(t, style(t, page, "#m-table tbody td", "::before", "content"), "none", "the label in the table mode")
			requireStyle(t, style(t, page, "#m-table thead tr", "", "display"), "table-row", "the head row in the table mode")
			requireStyle(t, style(t, page, "#m-table [data-part=wrapper]", "", "border-top-width"), "1px", "the table's own border")
			requireStyle(t, eval(t, page, `String((function () {
				var tds = document.querySelectorAll('#m-table tbody tr:first-child td'), a = tds[0].getBoundingClientRect(), b = tds[1].getBoundingClientRect();
				return Math.abs(a.top - b.top) < 2 && b.left >= a.right - 1;
			})())`), "true", "the cells on one line")

			// The default mode in the same 400px box.
			want := map[int][3]string{1280: {"table", "table-row", "table-cell"}, 390: {"block", "flex", "flex"}}[width]
			requireStyle(t, style(t, page, "#m-auto table", "", "display"), want[0], "the table in the auto mode")
			requireStyle(t, style(t, page, "#m-auto tbody tr", "", "display"), want[1], "a row in the auto mode")
			requireStyle(t, style(t, page, "#m-auto tbody td", "", "display"), want[2], "a cell in the auto mode")
			shot(t, page, fmt.Sprintf("mode-%d", width))
		})
	}
}

// The cards mode draws cards at every width: on a desktop window, in a box of
// 400px and in one of 900px.
func TestACardsModeKeepsTheCardsInAWideBox(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1280)
	page.Visit("/mode")
	for _, id := range []string{"m-cards", "m-wide-cards"} {
		requireStyle(t, style(t, page, "#"+id+" table", "", "display"), "block", "the table of "+id)
		requireStyle(t, style(t, page, "#"+id+" tbody tr", "", "display"), "flex", "a row of "+id)
		requireStyle(t, style(t, page, "#"+id+" tbody td[data-column=hour]", "::before", "content"), `"Last hour"`, "the label of "+id)
	}
	requireStyle(t, eval(t, page, `String(document.getElementById('m-wide-cards').getBoundingClientRect().width > 800)`), "true", "the wide box being wide")
}

// A row that arrives live is the row the server draws: it clones the template,
// so it has the same classes, the same labels and the same boxes as the rows
// before it, in each of the modes.
func TestARowThatArrivesLiveIsDrawnLikeTheOthers(t *testing.T) {
	for _, id := range []string{"m-auto", "m-table", "m-cards", "m-wide-cards"} {
		for _, width := range []int{1280, 390} {
			t.Run(fmt.Sprintf("%s-%d", id, width), func(t *testing.T) {
				app := serve(t, true)
				page := open(t, app, width)
				page.Visit("/mode")
				live(t, page, id, "Camera 3", "56")

				requireStyle(t, eval(t, page, `document.querySelectorAll('#`+id+` tbody tr:not([data-empty-row])').length`), "3", "the rows after one arrived")
				// The attributes of the live row and of the one under it, element by element.
				requireStyle(t, eval(t, page, `(function () {
				function shape(tr) {
					return tr.className + '|' + Array.from(tr.children).map(function (td) {
						return td.tagName + '.' + td.className + '[' + (td.dataset.label || '') + '/' + (td.dataset.column || '') + '/' + td.getAttribute('data-part') + ']';
					}).join(',');
				}
				var rows = document.querySelectorAll('#`+id+` tbody tr');
				return String(shape(rows[0]) === shape(rows[1]));
			})()`), "true", "the live row's classes and labels against a server row's")
				// And the same boxes: display, width, and the label in front of the value.
				requireStyle(t, eval(t, page, `(function () {
				var rows = document.querySelectorAll('#`+id+` tbody tr'), out = [];
				[0, 1].forEach(function (n) {
					var box = rows[n].getBoundingClientRect(), td = rows[n].children[1];
					out.push([getComputedStyle(rows[n]).display, Math.round(box.width), Math.round(box.height), getComputedStyle(td).display, getComputedStyle(td, '::before').content].join('/'));
				});
				return String(out[0] === out[1]) + ' ' + out.join(' | ');
			})()`)[:4], "true", "the live row's boxes against a server row's")
				shot(t, page, fmt.Sprintf("live-%s-%d", id, width))
			})
		}
	}
}

// The template is inert: it adds no row to the list nor a cell to the order.
func TestTheRowTemplateIsNotARow(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1280)
	page.Visit("/mode")
	requireStyle(t, eval(t, page, `document.querySelectorAll('#m-table tbody tr:not([data-empty-row])').length`), "2", "the rows of the list")
	requireStyle(t, eval(t, page, `String(document.querySelector('#m-table template[data-part=row-template]').getBoundingClientRect().height)`), "0", "the height of the template")
}
