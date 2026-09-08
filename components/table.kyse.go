//go:build kyse

package components

import (
	"strconv"
	"strings"
)

@go
// TableProps is a grid of rows under a row of headers.
//
// The head and the body are drawn from two slices that line up by position and
// by nothing else: the third cell of every row is the third column, whatever
// that column is called. A row carrying more cells than the table has columns
// draws them anyway, under no header, which is visible on the screen and
// therefore fixable.
//
// # What it does not do
//
// It does not sort, page or filter. Each of those is a request the server
// answers, and a component that made them would be a component deciding what a
// caller may ask for; the table draws the rows it was handed.
//
// It does not format. A number arrives as the string it should read as, because
// a thousands separator and a currency are decided where the locale is known,
// and a component that chose would choose wrong in half the places it is drawn.
type TableProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Caption is the sentence saying what the table lists. It is drawn under
	// the table and is the table's name to assistive technology, which is what
	// makes a page of three tables navigable. Empty draws none.
	Caption string
	// Columns are the headers, left to right. Each one also decides where the
	// cells beneath it sit.
	Columns []TableColumn
	// Rows are the lines. Empty draws Empty rather than a head with nothing
	// under it.
	Rows []TableRow
	// Empty is what stands in when there are no rows. A table with headers and
	// no body reads as one that failed to load rather than one with nothing to
	// list.
	Empty EmptyProps

	// SelectName turns on per-row selection and is the field the checkboxes
	// submit under. Empty draws no checkboxes at all.
	//
	// The rows carry their own key, and the bulk bar is revealed by CSS from
	// the fact that something is checked -- so nothing counts selections and
	// nothing has to be told when one is swapped away.
	SelectName string
	// SelectAllLabel names the header checkbox. Empty says "Select all rows",
	// which is the name it needs: a checkbox in a header cell with no name is
	// read as a checkbox in a header cell.
	SelectAllLabel string
	// BulkActions are the controls revealed once something is selected. They
	// submit the checked rows, so each one is a form button and not a link.
	BulkActions []ButtonProps
	// BulkLabel names the bar the actions sit in. Empty says "Bulk actions".
	BulkLabel string
	// BulkAction is where the bar submits, and BulkMethod is how. Empty posts
	// to the current URL, which is right when the list and its bulk operations
	// are one endpoint.
	BulkAction string
	BulkMethod string

	// Navigable turns the table into a grid: one tab stop for the whole thing,
	// and the arrow keys moving from cell to cell inside it.
	//
	// It is for a table people work in rather than read -- a spreadsheet, a
	// roster, a matrix of settings. It is the wrong answer for a list of
	// records, where the tab key reaching each link in turn is what somebody
	// expects, and where a grid takes that away.
	//
	// A row that carries a Level makes it a treegrid instead, and the arrow
	// keys then open and close as well as move.
	Navigable bool
	// SortKey and SortDir are the ordering in force: the column's Key, and
	// "asc" or "desc". They are what the server sorted by, drawn back onto the
	// header so the table says how it is ordered rather than leaving somebody
	// to infer it.
	SortKey string
	SortDir string
	// SortURL is the address a sort link goes to, with the column written as
	// "{sort}" and the direction as "{dir}": "/invoices?sort={sort}&dir={dir}".
	// Empty draws no sort links even on a sortable column, because a link with
	// nowhere to go is a header that looks like a control and is not one.
	SortURL string
	// ColumnsLabel names the menu that switches columns off. Empty says
	// "Columns". The menu is drawn only when some column is Hideable.
	ColumnsLabel string
	// SelectedLabel is the line saying how many rows are chosen, with the
	// count written as "{n}" and the total as "{total}": "{n} of {total} rows
	// selected". Empty draws that sentence in English.
	SelectedLabel string

	// The HTMX attributes for a sort link or a column menu, so the table can
	// swap itself instead of loading a document. The links stay real
	// addresses either way.
	HxTarget string
	HxSwap   string

	// Composed says this table is drawn inside another component, which has
	// already mounted the behaviour and owns the region.
	//
	// Without it the bridge is written twice -- once on the wrapper and once
	// here -- so the count runs twice per click and the live region is
	// rewritten twice, which a screen reader reads out twice.
	Composed bool
	// ID is what the bulk form and its checkboxes hang their ids off. Empty
	// falls back to SelectName, which two tables on one page share -- and then
	// each checkbox points at whichever form the browser found first.
	ID string
	// AscendingLabel, DescendingLabel and UnsortedLabel are what a sortable
	// header is called in each of its three states: "{column}, sorted
	// ascending". A header that only draws an arrow is a header a screen
	// reader reads as a link.
	AscendingLabel  string
	DescendingLabel string
	UnsortedLabel   string
	// Token is the CSRF token for the bulk form. A form that changes something
	// and carries none is refused by the framework, which is the intended
	// outcome and a confusing one to debug -- so it is a field here rather
	// than something to remember.
	Token string
}

// TableColumn is one header, and where everything under it sits.
type TableColumn struct {
	// Label is the header text.
	Label string
	// Key names the column: what a sort link sends, and what a visibility
	// toggle switches. Empty leaves the column unsortable and unhideable,
	// which is right for a column of controls.
	Key string
	// Sortable draws the header as the control that orders by it.
	//
	// The ordering happens on the server. A table sorted in the browser is
	// sorted only within the page that was fetched, so page two of a list
	// sorted by name holds whatever the server thought page two was -- and the
	// bug shows up as rows that seem to be in the wrong order to everyone
	// except whoever wrote it.
	Sortable bool
	// Hideable lets the column be switched off from the columns menu.
	Hideable bool
	// Hidden starts it switched off.
	Hidden bool
	// Align is "start", "center" or "end". Empty means "start".
	//
	// A column of numbers is "end": the digits then line up on the units, which
	// is what makes two amounts comparable by looking at them.
	Align string
}

// AlignClass is the alignment this column draws with.
func (c TableColumn) AlignClass() string {
	switch c.Align {
	case "center":
		return "text-center"
	case "end":
		return "text-end"
	default:
		return "text-start"
	}
}

// TableRow is one line of the body.
type TableRow struct {
	// Cells are the values, in the order Columns names them.
	Cells []TableCell
	// Key is what a checked row submits: the id of the record. It is required
	// for selection and ignored without it.
	Key string
	// Label names the row's checkbox: "Select invoice 2026-114". Without it
	// every checkbox down the column has the same name, which is the same as
	// having none.
	Label string
	// Selected draws the row checked.
	Selected bool
	// Level is how deep the row sits in a hierarchy, counting from one. Any
	// row carrying it makes the table a treegrid, where a row can be a branch
	// that opens.
	Level int
	// Expanded is whether a branch row is open, and means nothing on a row
	// with no children.
	Expanded bool
	// Hidden keeps the row out of the page without taking it off it, for a
	// list the browser pages: the server draws every row and marks the ones
	// outside the window, so the first paint is already right.
	//
	// Drawing them all visible and letting a script hide them afterwards is a
	// first frame that lies, and a page with no script that contradicts its
	// own footer.
	Hidden bool
	// Branch is whether the row has children under it. It is separate from
	// Expanded because a leaf has no expanded state at all -- which is what
	// the attribute being absent means, and is different from it being false.
	Branch bool
}

// TableCell is one value, as text or as markup.
type TableCell struct {
	// Text is what the cell says. It is escaped on the way out, and it is what
	// almost every cell is.
	Text string
	// SortValue is what a browser compares this cell by, when the browser is
	// the one doing the ordering.
	//
	// It exists because what a cell reads as and what it sorts as are not the
	// same string, and this component refuses to format: "1.240,00" sorts
	// before "860,00" as text and after it as money, and "8 de setembro"
	// sorts nowhere at all. So the server hands over the comparable form --
	// "1240.00", "2026-09-08" -- and keeps the readable one in Text.
	//
	// A value that parses as a number is compared as one; everything else is
	// compared as text, which for an ISO date is the same order. Empty falls
	// back to the cell's own text, which is right for a name and wrong for
	// everything with a unit in it.
	SortValue string
	// HTML is markup, drawn in place of Text whenever it is set.
	//
	// Nothing here escapes it, so what goes in is what the page gets, and the
	// caller is what makes that safe: pass what another component returned --
	// a badge, a button, an icon -- and never a string assembled around a value
	// somebody typed.
	HTML template.HTML
}

// AlignClass is the alignment of the column at position i, and the default
// where a row carries more cells than the table has columns.
func (p TableProps) AlignClass(i int) string {
	if i < len(p.Columns) {
		return p.Columns[i].AlignClass()
	}
	return TableColumn{}.AlignClass()
}
// Selectable is whether rows carry checkboxes.
func (p TableProps) Selectable() bool { return p.SelectName != "" }

// SelectAllName is what the header checkbox is called.
func (p TableProps) SelectAllName() string {
	if p.SelectAllLabel != "" {
		return p.SelectAllLabel
	}
	return "Select all rows"
}

// BulkName is what the bar of actions is called.
func (p TableProps) BulkName() string {
	if p.BulkLabel != "" {
		return p.BulkLabel
	}
	return "Bulk actions"
}

// PostMethod is how the bulk bar submits, and it is post unless the caller
// says otherwise -- a bulk operation changes something.
func (p TableProps) PostMethod() string {
	if p.BulkMethod != "" {
		return p.BulkMethod
	}
	return "post"
}

// RowName is what one row's checkbox is called: what the caller wrote, or the
// row's key, which is at least distinct.
func (p TableProps) RowName(row TableRow) string {
	if row.Label != "" {
		return row.Label
	}
	return "Select " + row.Key
}

// Hierarchical is whether any row declares a depth, which is what turns the
// grid into a treegrid.
func (p TableProps) Hierarchical() bool {
	for _, row := range p.Rows {
		if row.Level > 0 {
			return true
		}
	}
	return false
}

// Role is what the table announces as: a grid when it is navigated cell by
// cell, a treegrid when those rows nest, and nothing at all otherwise -- a
// plain table already has the right role and stating it again is noise.
func (p TableProps) Role() string {
	if !p.Navigable {
		return ""
	}
	if p.Hierarchical() {
		return "treegrid"
	}
	return "grid"
}

// Sortable is whether any column offers to order the table, which is what
// decides whether the header cells are controls at all.
func (p TableProps) Sortable() bool {
	if p.SortURL == "" {
		return false
	}
	for _, column := range p.Columns {
		if column.Sortable && column.Key != "" {
			return true
		}
	}
	return false
}

// SortedBy is whether the table is ordered by this column.
func (p TableProps) SortedBy(column TableColumn) bool {
	return column.Key != "" && column.Key == p.SortKey
}

// Order is what aria-sort says for one column: the direction when it is the
// one in force, and "none" on every other sortable column.
//
// A column that cannot be sorted carries nothing. aria-sort of "none" is a
// claim that the column could be sorted and is not, which on a column of
// avatars is a promise nobody can keep.
func (p TableProps) Order(column TableColumn) string {
	if !column.Sortable || column.Key == "" || !p.Sortable() {
		return ""
	}
	if !p.SortedBy(column) {
		return "none"
	}
	if p.SortDir == "desc" {
		return "descending"
	}
	return "ascending"
}

// NextDir is the direction a click on this header asks for: the other one when
// the table is already ordered by it, and ascending when it is not.
func (p TableProps) NextDir(column TableColumn) string {
	if p.SortedBy(column) && p.SortDir != "desc" {
		return "desc"
	}
	return "asc"
}

// SortHref is where one header points.
func (p TableProps) SortHref(column TableColumn) string {
	href := strings.ReplaceAll(p.SortURL, "{sort}", column.Key)
	return strings.ReplaceAll(href, "{dir}", p.NextDir(column))
}

// Hideable is whether any column can be switched off, which is what decides
// whether the columns menu is drawn.
func (p TableProps) Hideable() bool {
	for _, column := range p.Columns {
		if column.Hideable && column.Key != "" {
			return true
		}
	}
	return false
}

// ColumnsName is what the menu of columns is called.
func (p TableProps) ColumnsName() string {
	if p.ColumnsLabel != "" {
		return p.ColumnsLabel
	}
	return "Columns"
}

// Chosen is how many rows are drawn as selected.
func (p TableProps) Chosen() int {
	count := 0
	for _, row := range p.Rows {
		if row.Selected {
			count++
		}
	}
	return count
}

// SelectedText is the line saying how many rows are chosen.
//
// It is drawn by the server from what the server sent, and kept in step by the
// behaviour as boxes are ticked. Both write the same sentence, because both
// take it from here.
func (p TableProps) SelectedText() string {
	sentence := p.SelectedLabel
	if sentence == "" {
		sentence = "{n} of {total} rows selected"
	}
	sentence = strings.ReplaceAll(sentence, "{n}", strconv.Itoa(p.Chosen()))
	return strings.ReplaceAll(sentence, "{total}", strconv.Itoa(len(p.Rows)))
}

// ColumnKey is the key of the column at position i, and empty where a row
// carries more cells than the table has columns.
func (p TableProps) ColumnKey(i int) string {
	if i < len(p.Columns) {
		return p.Columns[i].Key
	}
	return ""
}

// ColumnHidden is whether the column at position i starts switched off.
func (p TableProps) ColumnHidden(i int) bool {
	return i < len(p.Columns) && p.Columns[i].Hidden
}

// TableBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const TableBehavior = "table"

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the table has anything for it to do.
//
// Sorting is not on that list: it is a link the server answers, and a link
// needs nothing mounted. What the behaviour does is the three things only the
// browser can -- keep the count in step as boxes are ticked, switch a column
// off, and move a cell at a time on a navigable table.
func (p TableProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" && !p.Composed && (p.Selectable() || p.Hideable() || p.Navigable) {
		p.Behavior = Behavior{Name: TableBehavior, Props: map[string]any{
			"selected": p.SelectedLabel,
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// FormID is the id of the bulk form, and what a checkbox points at.
func (p TableProps) FormID() string {
	if p.ID != "" {
		return p.ID + "-bulk"
	}
	return p.SelectName + "-bulk"
}

// Submits is whether a checkbox names a form at all. Without bulk actions
// there is no form, and an attribute naming an element that does not exist
// takes the checkbox out of every form instead of putting it in one.
func (p TableProps) Submits() bool { return p.Selectable() && len(p.BulkActions) > 0 }

// SortName is what a sortable header is called, in the state it is in.
func (p TableProps) SortName(column TableColumn) string {
	sentence := p.UnsortedLabel
	switch p.Order(column) {
	case "ascending":
		sentence = p.AscendingLabel
	case "descending":
		sentence = p.DescendingLabel
	}
	if sentence == "" {
		return ""
	}
	return strings.ReplaceAll(sentence, "{column}", column.Label)
}

// PartNames are the parts this component publishes.
func (p TableProps) PartNames() []string {
	return []string{
		"root", "table", "caption", "head", "header-cell", "row", "cell",
		"select-all", "select", "bulk", "count", "sort", "columns", "toggle",
	}
}
@endgo

@if(len(.Rows) > 0)
	{{-- A form only when there is something to submit. Wrapping every table in
	     one would nest a form inside whatever form the page already has, which
	     no browser accepts and which silently drops the inner one. --}}
	{{-- A composed table writes no root part: the component that draws it owns
	     the region, and two elements answering to the same name is a name that
	     reaches whichever the caller did not mean. --}}
	<div
		@if(!.Composed)
			data-part="root"
		@endif
		class="{{ .RootClass("table-container") }}"
		@attributes(.RootAttrs())
		@if(.Selectable())
			data-selectable="true"
		@endif
		@if(.Navigable)
			data-navigable="true"
		@endif
	>
		{{-- The count is drawn from what the server sent and kept in step by the
		     behaviour as boxes are ticked. Both write the same sentence, from
		     SelectedText, so the two cannot say different things. --}}
		@if(.Selectable())
			<p
				data-part="count"
				class="{{ .PartClass("count", "table-count") }}"
				role="status"
				aria-live="polite"
				data-selected-template="{{ .SelectedLabel }}"
				@attributes(.PartAttrs("count"))
			>{{ .SelectedText() }}</p>
		@endif

		@if(.Hideable())
			{{-- Checkboxes and not menu items: switching a column off is a
			     two-state choice that stays chosen, which is what a checkbox
			     is and what a menu item is not. --}}
			<details
				data-part="columns"
				class="{{ .PartClass("columns", "table-columns") }}"
				@attributes(.PartAttrs("columns"))
			>
				<summary>{{ .ColumnsName() }}</summary>
				<div>
					@foreach(.Columns as column)
						@if(column.Hideable && column.Key != "")
							<label
								data-part="toggle"
								@if(.PartClass("toggle") != "")
									class="{{ .PartClass("toggle") }}"
								@endif
								@attributes(.PartAttrs("toggle"))
							>
								<input
									type="checkbox"
									data-column-toggle="{{ column.Key }}"
									@if(!column.Hidden)
										checked
									@endif
								>
								{{ column.Label }}
							</label>
						@endif
					@endforeach
				</div>
			</details>
		@endif

		@if(.Selectable() && len(.BulkActions) > 0)
			<form
				data-part="bulk"
				class="{{ .PartClass("bulk", "table-bulk") }}"
				id="{{ .FormID() }}"
				method="{{ .PostMethod() }}"
				aria-label="{{ .BulkName() }}"
				@attributes(.PartAttrs("bulk"))
				@if(.BulkAction != "")
					action="{{ .BulkAction }}"
				@endif
			>
				@if(.Token != "")
					<input type="hidden" name="_token" value="{{ .Token }}">
				@endif
				@foreach(.BulkActions as action)
					{!! Button(action) !!}
				@endforeach
			</form>
		@endif

		<table
			data-part="table"
			class="{{ .PartClass("table", "table") }}"
			@attributes(.PartAttrs("table"))
			@if(.Role() != "")
				role="{{ .Role() }}"
			@endif
		>
			@if(.Caption != "")
				<caption
					data-part="caption"
					@if(.PartClass("caption") != "")
						class="{{ .PartClass("caption") }}"
					@endif
					@attributes(.PartAttrs("caption"))
				>{{ .Caption }}</caption>
			@endif
			<thead
				data-part="head"
				@if(.PartClass("head") != "")
					class="{{ .PartClass("head") }}"
				@endif
				@attributes(.PartAttrs("head"))
			>
				<tr>
					@if(.Selectable())
						<th
							data-part="header-cell"
							scope="col"
							class="{{ .PartClass("header-cell", "w-0") }}"
							@attributes(.PartAttrs("header-cell"))
						>
							<input
								data-part="select-all"
								@if(.PartClass("select-all") != "")
									class="{{ .PartClass("select-all") }}"
								@endif
								type="checkbox"
								data-select-all
								aria-label="{{ .SelectAllName() }}"
								@attributes(.PartAttrs("select-all"))
							>
						</th>
					@endif
					@foreach(.Columns as column)
						<th
							data-part="header-cell"
							scope="col"
							class="{{ .PartClass("header-cell", column.AlignClass()) }}"
							@attributes(.PartAttrs("header-cell"))
							@if(column.Key != "")
								data-column="{{ column.Key }}"
							@endif
							@if(column.Hidden)
								hidden
							@endif
							@if(.Order(column) != "")
								aria-sort="{{ .Order(column) }}"
							@endif
						>
							@if(column.Sortable && .Sortable() && column.Key != "")
								<a
									data-part="sort"
									class="{{ .PartClass("sort", "table-sort") }}"
									href="{{ .SortHref(column) }}"
									@attributes(.PartAttrs("sort"))
									@if(.SortName(column) != "")
										aria-label="{{ .SortName(column) }}"
									@endif
									@if(.HxTarget != "")
										hx-get="{{ .SortHref(column) }}"
										hx-target="{{ .HxTarget }}"
									@endif
									@if(.HxSwap != "")
										hx-swap="{{ .HxSwap }}"
									@endif
								>{{ column.Label }}</a>
							@endif
							@if(!column.Sortable || !.Sortable() || column.Key == "")
								{{ column.Label }}
							@endif
						</th>
					@endforeach
				</tr>
			</thead>
			<tbody>
				@foreach(.Rows as row)
					<tr
						data-part="row"
						@if(.PartClass("row") != "")
							class="{{ .PartClass("row") }}"
						@endif
						@attributes(.PartAttrs("row"))
						@if(row.Hidden)
							hidden
						@endif
						@if(row.Level > 0)
							aria-level="{{ row.Level }}"
							data-level="{{ row.Level }}"
						@endif
						@if(row.Branch && row.Expanded)
							aria-expanded="true"
						@endif
						@if(row.Branch && !row.Expanded)
							aria-expanded="false"
						@endif
						@if(.Selectable() && row.Selected)
							aria-selected="true"
						@endif
					>
						@if(.Selectable())
							<td
								data-part="cell"
								class="{{ .PartClass("cell", "w-0") }}"
								@attributes(.PartAttrs("cell"))
							>
								<input
									data-part="select"
									@if(.PartClass("select") != "")
										class="{{ .PartClass("select") }}"
									@endif
									type="checkbox"
									name="{{ .SelectName }}"
									value="{{ row.Key }}"
									@if(.Submits())
										form="{{ .FormID() }}"
									@endif
									@if(row.Hidden)
										disabled
									@endif
									aria-label="{{ .RowName(row) }}"
									@attributes(.PartAttrs("select"))
									@if(row.Selected)
										checked
									@endif
								>
							</td>
						@endif
						@for(i := 0; i < len(row.Cells); i++)
							<td
								data-part="cell"
								class="{{ .PartClass("cell", .AlignClass(i)) }}"
								@attributes(.PartAttrs("cell"))
								@if(.ColumnKey(i) != "")
									data-column="{{ .ColumnKey(i) }}"
								@endif
								@if(.ColumnHidden(i))
									hidden
								@endif
								@if(row.Cells[i].SortValue != "")
									data-sort-value="{{ row.Cells[i].SortValue }}"
								@endif
							>
								@if(row.Cells[i].HTML != "")
									{!! row.Cells[i].HTML !!}
								@else
									{{ row.Cells[i].Text }}
								@endif
							</td>
						@endfor
					</tr>
				@endforeach
			</tbody>
		</table>
	</div>
@else
	{!! Empty(.Empty) !!}
@endif
