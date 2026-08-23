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
@endgo

@if(len(.Rows) > 0)
	<div class="table-container">
		<table class="table">
			@if(.Caption != "")
				<caption>{{ .Caption }}</caption>
			@endif
			<thead>
				<tr>
					@foreach(.Columns as column)
						<th scope="col" class="{{ column.AlignClass() }}">{{ column.Label }}</th>
					@endforeach
				</tr>
			</thead>
			<tbody>
				@foreach(.Rows as row)
					<tr>
						@for(i := 0; i < len(row.Cells); i++)
							<td class="{{ .AlignClass(i) }}">
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
