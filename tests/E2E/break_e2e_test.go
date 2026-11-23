//go:build e2e

package e2e

import (
	"fmt"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// breakList is a list of six columns, as a listing of a dashboard is, in the
// auto mode and with CardBelow as given.
func breakList(id string, below components.DataTableBreak) string {
	columns := []components.TableColumn{
		{Label: "Customer", Key: "customer"}, {Label: "Plan", Key: "plan"}, {Label: "State", Key: "state"},
		{Label: "Created", Key: "created"}, {Label: "Amount", Key: "amount", Align: "end"}, {Label: "Fee", Key: "fee", Align: "end"},
	}
	var rows []components.TableRow
	for _, n := range []string{"1", "2"} {
		rows = append(rows, components.TableRow{Key: n, Cells: []components.TableCell{
			{Text: "Customer " + n}, {Text: "Pro"}, {Text: "Paid"}, {Text: "2026-11-0" + n}, {Text: "100,00"}, {Text: "3,50"},
		}})
	}
	return string(components.DataTable(components.DataTableProps{
		ID: id, Label: id, URL: "/break", Complete: true, CardBelow: below, Columns: columns, Rows: rows,
	}))
}

// breakBoxes is the same list in boxes of three widths on one page, for each
// break, and the list beside a menu of 256px, the way the dashboard of an
// application holds it. A box is a class the stylesheet is compiled with from
// this file: its width is the content, as p-2 and the border take 18px of it.
func breakBoxes() string {
	box := func(width, id string, below components.DataTableBreak) string {
		return fmt.Sprintf(`<section class="mb-6 %s border p-2">%s</section>`, width, breakList(id, below))
	}
	var out string
	for _, c := range []struct {
		id    string
		below components.DataTableBreak
	}{{"b-screen", components.DataTableBreakScreen}, {"b-block", components.DataTableBreakBlock}, {"b-wide", components.DataTableBreakBlockWide}} {
		out += box("w-[600px]", c.id+"-600", c.below) + box("w-[680px]", c.id+"-680", c.below) + box("w-[800px]", c.id+"-800", c.below)
	}
	return out
}

// breakBeside is the lists in the content of a page that has a menu of 256px open
// beside it: the menu keeps its width and the content takes what is left, less
// the 16px of padding on each side.
func breakBeside() string {
	return `<div class="flex"><aside class="w-64 shrink-0 border-e p-4">Menu</aside><main class="min-w-0 flex-1 p-4">` +
		breakList("s-screen", components.DataTableBreakScreen) +
		breakList("s-block", components.DataTableBreakBlock) +
		breakList("s-wide", components.DataTableBreakBlockWide) + `</main></div>`
}

// whatItDraws reads the layout of the list: the table or the cards.
func whatItDraws(t *testing.T, page interface{ Evaluate(string, any) }, id string) string {
	t.Helper()
	switch got := style(t, page, "#"+id+" tbody tr", "", "display"); got {
	case "table-row":
		return "table"
	case "flex":
		return "cards"
	default:
		return got
	}
}

// The block decides and not the screen: on a window of 1280px the same list is
// cards in a box under the break and a table in a box over it, at 640px for
// DataTableBreakBlock and at 720px for DataTableBreakBlockWide, and the default
// is a table in all three because the screen is wide.
func TestTheBlockDecidesWhereCardBelowAsksIt(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 1280)
	page.Visit("/break")
	want := map[string][3]string{
		"b-screen": {"table", "table", "table"},
		"b-block":  {"cards", "table", "table"},
		"b-wide":   {"cards", "cards", "table"},
	}
	for name, states := range want {
		for i, width := range []string{"600", "680", "800"} {
			id := name + "-" + width
			requireStyle(t, whatItDraws(t, page, id), states[i], fmt.Sprintf("%s (a box of %spx)", name, width))
		}
	}
	// A card is a card: the label of a column in front of its value, and nothing
	// scrolls sideways.
	requireStyle(t, style(t, page, "#b-wide-600 tbody td[data-column=amount]", "::before", "content"), `"Amount"`, "the label of a cell in a card")
	requireStyle(t, eval(t, page, `String(document.documentElement.scrollWidth <= window.innerWidth)`), "true", "the page not scrolling sideways")
	shot(t, page, "break-1280")
}

// On a phone window the default is cards, and the block breaks still ask the box
// and not the window: a box of 800px is a table in it even though the window is
// 390px, which is the question the break exists to put.
func TestOnAPhoneTheDefaultFollowsTheScreenAndTheBlockBreaksTheBox(t *testing.T) {
	app := serve(t, true)
	page := open(t, app, 390)
	page.Visit("/break")
	requireStyle(t, whatItDraws(t, page, "b-screen-800"), "cards", "the default on a phone")
	requireStyle(t, whatItDraws(t, page, "b-block-800"), "table", "CardBelow block in a box of 800px")
	requireStyle(t, whatItDraws(t, page, "b-block-600"), "cards", "CardBelow block in a box of 600px")
}

// Beside a menu: a list in the content next to a menu 256px wide is a table when
// the window leaves it room and cards when it does not, which the screen alone
// cannot tell.
func TestAListBesideAMenuFollowsTheRoomItHas(t *testing.T) {
	for _, c := range []struct {
		window int
		screen string
		block  string
		wide   string
	}{
		// content = window - 256 - 2 * 16px of padding: 480, 736 and 992px
		{768, "table", "cards", "cards"},
		{1024, "table", "table", "table"},
		{1280, "table", "table", "table"},
	} {
		t.Run(fmt.Sprint(c.window), func(t *testing.T) {
			app := serve(t, true)
			page := open(t, app, c.window)
			page.Visit("/beside")
			requireStyle(t, whatItDraws(t, page, "s-screen"), c.screen, "the default")
			requireStyle(t, whatItDraws(t, page, "s-block"), c.block, "CardBelow block")
			requireStyle(t, whatItDraws(t, page, "s-wide"), c.wide, "CardBelow block-wide")
			requireStyle(t, eval(t, page, `String(document.documentElement.scrollWidth <= window.innerWidth)`), "true", "the page not scrolling sideways")
			shot(t, page, fmt.Sprintf("beside-%d", c.window))
		})
	}
}
