package unit

import (
	"html"
	"html/template"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/hesape/translation"
	"github.com/arandu-io/kyse/components"
)

// What the footer says for one row and for the others, and what a list says on
// a narrow screen: each row a card, its column marked Primary the title, the
// columns marked CardHidden left out, the others a name and a value, the row's
// menu in sight. The card is the same markup as the table, drawn under a
// container query, so what these hold is the classes and attributes that draw it.

// footerOf is a data table of total rows over a page that counts them, with the
// lines given, and the text of its footer.
func footerOf(t *testing.T, total int, lines components.DataTableLines, page components.Page) string {
	t.Helper()
	html := string(components.DataTable(components.DataTableProps{
		ID: "n", Label: "N", URL: "/n", Lines: lines, PageData: page,
		Paginator: pagination.NewLengthAwarePage(min(total, 10), total, 10, 1, pagination.Options{Path: "/n"}),
		Columns:   []components.TableColumn{{Label: "Name", Key: "name"}},
		Rows:      []components.TableRow{{Cells: []components.TableCell{{Text: "a"}}}},
	}))
	match := regexp.MustCompile(`(?s)data-part="count".*?>([^<]*)</div>`).FindStringSubmatch(html)
	if match == nil {
		t.Fatalf("no footer in\n%s", html)
	}
	return strings.TrimSpace(match[1])
}

func TestTheFooterSaysOneRowInTheSingularByDefault(t *testing.T) {
	for total, want := range map[int]string{1: "1 row.", 0: "0 rows.", 2: "2 rows.", 137: "137 rows."} {
		if got := footerOf(t, total, components.DataTableLines{}, nil); got != want {
			t.Errorf("a total of %d says %q, want %q", total, got, want)
		}
	}
}

func TestTheFooterTakesASingularAndAPluralFromTheApplication(t *testing.T) {
	lines := components.DataTableLines{RowsOne: "{total} linha", RowsMany: "{total} linhas"}
	for total, want := range map[int]string{1: "1 linha", 0: "0 linhas", 2: "2 linhas"} {
		if got := footerOf(t, total, lines, nil); got != want {
			t.Errorf("a total of %d says %q, want %q", total, got, want)
		}
	}
}

// A sentence written before the split is still the sentence for both counts,
// and a form written beside it wins for its own.
func TestTheOlderRowsSentenceStillSpeaksForBothCounts(t *testing.T) {
	old := components.DataTableLines{Rows: "{total} linha(s)."}
	for total, want := range map[int]string{1: "1 linha(s).", 5: "5 linha(s).", 0: "0 linha(s)."} {
		if got := footerOf(t, total, old, nil); got != want {
			t.Errorf("the older sentence at %d says %q, want %q", total, got, want)
		}
	}
	mixed := components.DataTableLines{Rows: "{total} linha(s).", RowsOne: "1 linha."}
	if got := footerOf(t, 1, mixed, nil); got != "1 linha." {
		t.Errorf("the one-row sentence is %q, want the form written for it", got)
	}
	if got := footerOf(t, 3, mixed, nil); got != "3 linha(s)." {
		t.Errorf("the other counts say %q, want the older sentence", got)
	}
}

// The catalogue speaks in the same order: the split keys, then the older one,
// and a sentence written on the props beats them all.
func TestTheCatalogueSpeaksForTheFooterInBothKeys(t *testing.T) {
	older := ptPage(translation.Lines{"data_table.rows": "{total} linha(s)"})
	if got := footerOf(t, 1, components.DataTableLines{}, older); got != "1 linha(s)" {
		t.Errorf("a catalogue with the older key says %q, want it for one row", got)
	}
	if got := footerOf(t, 4, components.DataTableLines{}, older); got != "4 linha(s)" {
		t.Errorf("a catalogue with the older key says %q, want it for four", got)
	}

	split := ptPage(translation.Lines{"data_table.rows": "{total} velho", "data_table.rows_one": "{total} linha", "data_table.rows_many": "{total} linhas"})
	if got := footerOf(t, 1, components.DataTableLines{}, split); got != "1 linha" {
		t.Errorf("the split catalogue says %q for one, want 1 linha", got)
	}
	if got := footerOf(t, 4, components.DataTableLines{}, split); got != "4 linhas" {
		t.Errorf("the split catalogue says %q for four, want 4 linhas", got)
	}

	if got := footerOf(t, 1, components.DataTableLines{Rows: "{total} escrita"}, split); got != "1 escrita" {
		t.Errorf("a sentence on the props says %q, want it over the catalogue", got)
	}
}

// RowsFor is the same rule for an application that draws its own line.
func TestRowsForChoosesTheFormByTheCount(t *testing.T) {
	lines := components.DataTableLines{RowsOne: "um", RowsMany: "{total} vários"}.Filled(nil)
	if got := lines.RowsFor(1); got != "um" {
		t.Errorf("RowsFor(1) is %q", got)
	}
	if got := lines.RowsFor(0); got != "0 vários" {
		t.Errorf("RowsFor(0) is %q", got)
	}
}

// The browser rewrites the footer of a list it filters, so it is handed both
// forms: rows is the plural, which is what a script that predates the singular
// reads, and rowsOne is the new one.
func TestTheBehaviourIsHandedBothFormsOfTheFooter(t *testing.T) {
	page := html.UnescapeString(string(components.DataTable(components.DataTableProps{
		ID: "c", Label: "C", URL: "/c", Complete: true,
		Lines:   components.DataTableLines{RowsOne: "{total} linha", RowsMany: "{total} linhas"},
		Columns: []components.TableColumn{{Label: "Name"}},
		Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "a"}}}},
	})))
	for _, want := range []string{`"rows":"{total} linhas"`, `"rowsOne":"{total} linha"`} {
		if !strings.Contains(page, want) {
			t.Errorf("the behaviour props lack %s", want)
		}
	}
	old := html.UnescapeString(string(components.DataTable(components.DataTableProps{
		ID: "c", Label: "C", URL: "/c", Complete: true, Lines: components.DataTableLines{Rows: "{total} linha(s)."},
		Columns: []components.TableColumn{{Label: "Name"}},
		Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "a"}}}},
	})))
	if !strings.Contains(old, `"rows":"{total} linha(s)."`) || !strings.Contains(old, `"rowsOne":"{total} linha(s)."`) {
		t.Errorf("the older sentence is not handed on as both forms:\n%s", old)
	}
}

// phoneList is a list marked for the card layout.
func phoneList(columns []components.TableColumn, actions bool, selectName string) components.DataTableProps {
	row := components.TableRow{Key: "1", Cells: []components.TableCell{{Text: "INV-1"}, {Text: "Ada"}, {Text: "paid"}, {Text: "10.00"}, {Text: ""}}}
	if actions {
		row.Actions = []components.MenuItem{{Label: "Open", URL: "/1"}}
	}
	return components.DataTableProps{
		ID: "p", Label: "P", URL: "/p", Columns: columns, Rows: []components.TableRow{row, {Key: "2", Cells: row.Cells}},
		SelectName: selectName,
	}
}

var phoneColumns = []components.TableColumn{
	{Label: "Number", Key: "number", Sortable: true, CardHidden: true},
	{Label: "Customer", Key: "customer", Sortable: true, Primary: true},
	{Label: "Status", Key: "status"},
	{Label: "Total", Key: "total", Align: "end"},
	{Label: ""},
}

// tags is each opening tag of the element named, in order.
func tags(html, element string) []string {
	return regexp.MustCompile(`(?s)<`+element+`\b[^>]*>`).FindAllString(html, -1)
}

func withData(tag, name string) bool { return strings.Contains(tag, "data-column=\""+name+"\"") }

func TestTheScreenAndNotTheBlockIsWhatTheCardsAreAskedAgainst(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, true, "")))
	// The block is not a container: a container query would ask its width, which
	// on a desktop is under 600px inside any card of a dashboard.
	if strings.Contains(html, "@container") || strings.Contains(html, "@max-[") || strings.Contains(html, "@min-[") {
		t.Error("a card rule asks the width of the block the table is in")
	}
	// Every card rule is under the media query, so above it nothing applies.
	for _, class := range regexp.MustCompile(`[^\s"]*max-\[600px\]:[^\s"]+`).FindAllString(html, -1) {
		if !strings.HasPrefix(class, "max-[600px]:") {
			t.Errorf("a card rule outside the query: %s", class)
		}
	}
	if !strings.Contains(html, "max-[600px]:flex") {
		t.Error("no row is drawn as a card under 600px")
	}
}

func TestEachCellKeepsTheNameOfItsColumnForTheCard(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, true, "")))
	body := html[strings.Index(html, "<tbody"):]
	var status, total string
	for _, td := range tags(body, "td") {
		switch {
		case withData(td, "status"):
			status = td
		case withData(td, "total"):
			total = td
		}
	}
	for name, tag := range map[string]string{"status": status, "total": total} {
		want := `data-label="` + strings.ToUpper(name[:1]) + name[1:] + `"`
		if !strings.Contains(tag, want) {
			t.Errorf("the %s cell lacks %s: %s", name, want, tag)
		}
		if !strings.Contains(tag, "before:content-[attr(data-label)]") {
			t.Errorf("the %s cell does not draw its label from data-label: %s", name, tag)
		}
	}
	if strings.Contains(body, ">Status<") || strings.Contains(body, ">Total<") {
		t.Error("the label is text in a cell, so a script reading a cell's text would read it too")
	}
	// A cell of a column with no name has none to draw in front of its value.
	for _, td := range tags(body, "td") {
		if !strings.Contains(td, "data-column") && strings.Contains(td, "data-label=") {
			t.Errorf("a cell with no column name carries a label: %s", td)
		}
	}
}

func TestThePrimaryColumnTitlesTheCardAndCardHiddenIsLeftOut(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, false, "")))
	body := html[strings.Index(html, "<tbody"):]
	for _, td := range tags(body, "td") {
		switch {
		case withData(td, "customer"):
			if !strings.Contains(td, "max-[600px]:order-first") || !strings.Contains(td, "max-[600px]:before:hidden") {
				t.Errorf("the primary cell is not the title: %s", td)
			}
		case withData(td, "number"):
			if !strings.Contains(td, "max-[600px]:hidden") {
				t.Errorf("a CardHidden cell stays in the card: %s", td)
			}
		case withData(td, "status"):
			if strings.Contains(td, "order-first") || strings.Contains(td, "max-[600px]:hidden") {
				t.Errorf("a plain cell is a title or hidden: %s", td)
			}
		}
	}
}

// With none marked the title is the first column that the card would show.
func TestWithNoPrimaryTheFirstColumnTheCardShowsIsTheTitle(t *testing.T) {
	columns := []components.TableColumn{
		{Label: "Number", Key: "number", CardHidden: true},
		{Label: "Customer", Key: "customer"},
		{Label: "Status", Key: "status"},
	}
	p := phoneList(columns, false, "")
	if got := p.PrimaryColumn(); got != 1 {
		t.Errorf("the title is column %d, want 1: the first that is not CardHidden", got)
	}
	columns[1].Hideable, columns[1].Hidden = true, true
	if got := phoneList(columns, false, "").PrimaryColumn(); got != 2 {
		t.Errorf("a switched-off column titles the card: column %d", got)
	}
	plain := phoneList([]components.TableColumn{{Label: "A", Key: "a"}, {Label: "B", Key: "b"}}, false, "")
	if got := plain.PrimaryColumn(); got != 0 {
		t.Errorf("the default title is column %d, want the first", got)
	}
	if got := phoneList(nil, false, "").PrimaryColumn(); got != -1 {
		t.Errorf("a table with no column has title %d", got)
	}
}

// A column marked both is the title: a card with no title is the worse defect.
func TestPrimaryBeatsCardHidden(t *testing.T) {
	p := phoneList([]components.TableColumn{{Label: "A", Key: "a", Primary: true, CardHidden: true}, {Label: "B", Key: "b"}}, false, "")
	if p.CardHiddenAt(0) {
		t.Error("the title is left out of its own card")
	}
}

// The head turns into the bar that orders the list: what orders it stays, what
// only names a column goes, and the table is never wider than the screen.
func TestTheHeadKeepsOnlyTheControlsThatOrderTheList(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, true, "ids")))
	head := html[strings.Index(html, "<thead"):strings.Index(html, "</thead>")]
	for _, th := range tags(head, "th") {
		switch {
		case withData(th, "number"), withData(th, "customer"):
			if strings.Contains(th, "max-[600px]:hidden") {
				t.Errorf("a header that orders the list is dropped: %s", th)
			}
		case withData(th, "status"), withData(th, "total"):
			if !strings.Contains(th, "max-[600px]:hidden") {
				t.Errorf("a header that only names a column stays: %s", th)
			}
		}
	}
	if !strings.Contains(html, "max-[600px]:overflow-visible") {
		t.Error("the table still scrolls sideways under the breakpoint")
	}
}

func TestARowMenuStaysInTheCardAndARowWithNoneDrawsNoFoot(t *testing.T) {
	p := phoneList(phoneColumns, true, "")
	p.Rows[1].Actions = nil
	html := string(components.DataTable(p))
	body := html[strings.Index(html, "<tbody"):]
	feet := 0
	for _, td := range tags(body, "td") {
		if strings.Contains(td, "max-[600px]:border-t") {
			feet++
		}
	}
	if feet != 1 {
		t.Errorf("%d cards have a foot, want one: only the row with a menu", feet)
	}
}

func TestTheCheckboxSitsInTheCornerAndTheTitleLeavesRoomForIt(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, false, "ids")))
	body := html[strings.Index(html, "<tbody"):]
	if !strings.Contains(body, "max-[600px]:absolute") || !strings.Contains(body, "max-[600px]:ps-8") {
		t.Error("the checkbox is not in the corner of the card, or the title does not leave room for it")
	}
	without := string(components.DataTable(phoneList(phoneColumns, false, "")))
	if strings.Contains(without[strings.Index(without, "<tbody"):], "max-[600px]:ps-8") {
		t.Error("a table with no selection leaves room for a checkbox")
	}
}

// The card is classes and attributes: nothing a strict policy would refuse.
func TestTheCardLayoutNeedsNoInlineStyleOrScript(t *testing.T) {
	html := string(components.DataTable(phoneList(phoneColumns, true, "ids")))
	for _, banned := range []string{" style=", "<style", "<script", " onclick=", "x-data"} {
		if strings.Contains(html, banned) {
			t.Errorf("the markup carries %q", banned)
		}
	}
}

// The empty row that says nothing matched is a block in a card layout, or it
// would be an anonymous table as narrow as its text.
func TestTheNoResultRowIsABlockInACardLayout(t *testing.T) {
	p := phoneList(phoneColumns, false, "")
	p.Rows = nil
	p.Searchable = true
	p.State = components.DataTableValues{Q: "zzz"}
	html := string(components.DataTable(p))
	if !strings.Contains(html, "No results.") || !strings.Contains(html, "max-[600px]:block") {
		t.Errorf("the no-result row is not a block:\n%s", html)
	}
}

// actionsList is a list whose last column holds the row's own controls: the
// first row has some and the second has an empty cell.
func actionsList() components.DataTableProps {
	return components.DataTableProps{
		ID: "a", Label: "A", URL: "/a",
		Columns: []components.TableColumn{
			{Label: "Name", Key: "name", Primary: true},
			{Label: "Note", Key: "note"},
			{Label: "Actions", Align: "end", Actions: true},
		},
		Rows: []components.TableRow{
			{Key: "1", Cells: []components.TableCell{{Text: "Ada"}, {Text: ""}, {HTML: `<a href="/1">Edit</a>`}}},
			{Key: "2", Cells: []components.TableCell{{Text: "Grace"}, {Text: "kept"}, {}}},
			{Key: "3", Cells: []components.TableCell{{Text: "Alan"}, {Text: "kept"}, {HTML: " \n\t "}}},
		},
	}
}

// cellsOf is the opening tag of each cell of the body, row by row.
func cellsOf(html string) [][]string {
	body := html[strings.Index(html, "<tbody"):]
	var rows [][]string
	for _, tr := range strings.Split(body, "<tr")[1:] {
		rows = append(rows, tags(tr, "td"))
	}
	return rows
}

// The cell of the column marked Actions is the foot of the card with no name in
// front of it; one with nothing in it is left out of the card, and so is the
// label of any other cell that has no value.
func TestTheActionsCellIsTheFootOfTheCardAndAnEmptyOneIsLeftOut(t *testing.T) {
	rows := cellsOf(string(components.DataTable(actionsList())))
	if len(rows) != 3 {
		t.Fatalf("%d rows, want 3", len(rows))
	}
	foot := rows[0][2]
	for _, want := range []string{"max-[600px]:border-t", "max-[600px]:justify-end", "max-[600px]:before:hidden", "max-[600px]:order-last"} {
		if !strings.Contains(foot, want) {
			t.Errorf("the actions cell lacks %s: %s", want, foot)
		}
	}
	// Under the breakpoint only: the cell is the cell it was above it.
	for _, class := range strings.Fields(foot[strings.Index(foot, `class="`)+7:]) {
		if strings.Contains(class, "border-t") && !strings.HasPrefix(class, "max-[600px]:") {
			t.Errorf("a rule outside the card query: %s", class)
		}
	}
	// A note with nothing in it is not drawn; one with a value is.
	if !strings.Contains(rows[0][1], "max-[600px]:hidden") {
		t.Errorf("an empty note stays in the card: %s", rows[0][1])
	}
	if strings.Contains(rows[1][1], "max-[600px]:hidden") {
		t.Errorf("a note with a value is left out: %s", rows[1][1])
	}
	// An empty actions cell, or one of white space only, draws no foot and no label.
	for _, n := range []int{1, 2} {
		cell := rows[n][2]
		if !strings.Contains(cell, "max-[600px]:hidden") || strings.Contains(cell, "max-[600px]:border-t") {
			t.Errorf("row %d has no actions and draws a foot: %s", n+1, cell)
		}
	}
	// The header still names the column for the table above the breakpoint.
	page := string(components.DataTable(actionsList()))
	if !strings.Contains(page[strings.Index(page, "<thead"):strings.Index(page, "</thead>")], "Actions") {
		t.Error("the column of actions lost its header")
	}
}

// The card's title is drawn when it is empty, so a card never loses its heading
// to the rule that leaves empty cells out.
func TestAnEmptyTitleIsStillTheCardsTitle(t *testing.T) {
	p := actionsList()
	p.Rows[0].Cells[0] = components.TableCell{}
	rows := cellsOf(string(components.DataTable(p)))
	if strings.Contains(rows[0][0], "max-[600px]:hidden") || !strings.Contains(rows[0][0], "max-[600px]:order-first") {
		t.Errorf("an empty title is left out of its card: %s", rows[0][0])
	}
}

// A column marked Actions and Primary is the title: the card has to have one.
func TestAnActionsColumnThatIsAlsoThePrimaryIsStillTheTitle(t *testing.T) {
	p := actionsList()
	p.Columns[2].Primary = true
	p.Columns[0].Primary = false
	rows := cellsOf(string(components.DataTable(p)))
	if !strings.Contains(rows[0][2], "max-[600px]:order-first") || strings.Contains(rows[0][2], "max-[600px]:order-last") {
		t.Errorf("the primary column is not the title: %s", rows[0][2])
	}
}

func TestACellIsEmptyWhenItHasNeitherTextNorMarkup(t *testing.T) {
	for cell, want := range map[string]bool{"": true, "  ": true, "x": false} {
		if got := (components.TableCell{Text: cell}).Empty(); got != want {
			t.Errorf("Text %q: Empty is %t, want %t", cell, got, want)
		}
		if got := (components.TableCell{HTML: template.HTML(cell)}).Empty(); got != want {
			t.Errorf("HTML %q: Empty is %t, want %t", cell, got, want)
		}
	}
}
