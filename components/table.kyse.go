//go:build kyse

package components

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

// FirstRow is the index of the row holding the tab stop on a navigable table.
func (p TableProps) FirstRow() int {
	if len(p.Rows) == 0 {
		return -1
	}
	return 0
}

// PartNames are the parts this component publishes.
func (p TableProps) PartNames() []string {
	return []string{
		"root", "table", "caption", "head", "header-cell", "row", "cell",
		"select-all", "select", "bulk",
	}
}
@endgo

@if(len(.Rows) > 0)
	{{-- A form only when there is something to submit. Wrapping every table in
	     one would nest a form inside whatever form the page already has, which
	     no browser accepts and which silently drops the inner one. --}}
	<div
		data-part="root"
		class="{{ .RootClass("table-container") }}"
		@attributes(.RootAttrs())
		@if(.Selectable())
			data-selectable="true"
		@endif
		@if(.Navigable)
			data-navigable="true"
		@endif
	>
		@if(.Selectable() && len(.BulkActions) > 0)
			<form
				data-part="bulk"
				class="{{ .PartClass("bulk", "table-bulk") }}"
				id="{{ .SelectName }}-bulk"
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
						>{{ column.Label }}</th>
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
									form="{{ .SelectName }}-bulk"
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
