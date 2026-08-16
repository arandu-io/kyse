//go:build kyse

package components

@go
// StatCardProps is a table of numbers with a heading over it: how many of a
// thing there are, one row per whatever the things are grouped by.
//
// It exists because a count nobody can look at is a count nobody checks. The
// application already holds the numbers -- a counter, a query, a length -- and
// what was missing was the shape that draws them without every screen inventing
// its own table.
//
// # What it deliberately does not do
//
// It does not fetch. There is no URL on it and no polling: the caller reads its
// own numbers, with its own authority, and hands them over. A component that
// fetched would be a component deciding who may see what.
//
// It keeps no history. No peak, no average, no window, no chart. Those are what
// separate a panel from an APM product, and building an APM product is a
// decision this project took the other way (ADR 0055).
//
// It does not format. Values arrive as strings that are already what they
// should read as, because a thousands separator is a local decision and a
// component that made it would make it wrong in half the places it is used.
type StatCardProps struct {
	// Title is the heading: what is being counted.
	Title string
	// Meta is the small line under it -- when the numbers were read, what
	// period they cover. Empty draws none.
	Meta string
	// Columns are the headers over the numbers, left to right. A row with more
	// values than there are columns draws them anyway; the extra ones sit under
	// no header, which is visible and therefore fixable.
	Columns []string
	// Rows are the lines. Empty draws Empty instead of an empty table.
	Rows []StatRow
	// Empty is what stands in when there are no rows. A blank card reads as a
	// card that failed rather than one with nothing to show.
	Empty EmptyProps
}

// StatRow is one line: what it is, and the numbers about it.
type StatRow struct {
	// Label is what the numbers are about -- a tenant, a queue, a host.
	Label string
	// Values are the numbers, in the order Columns names them, already
	// formatted.
	Values []string
}
@endgo

<article class="card p-5">
	<header>
		<h3 class="font-semibold tracking-tight">{{ .Title }}</h3>
		@if(.Meta != "")
			<p class="text-muted-foreground mt-1 text-xs">{{ .Meta }}</p>
		@endif
	</header>

	@forelse(.Rows as row)
		<table class="mt-4 w-full text-sm">
			<thead>
				<tr class="text-muted-foreground text-left text-xs">
					<th class="font-medium">Name</th>
					@foreach(.Columns as column)
						<th class="font-medium text-right">{{ column }}</th>
					@endforeach
				</tr>
			</thead>
			<tbody class="divide-y">
				@foreach(.Rows as row)
					<tr>
						<td class="py-2">{{ row.Label }}</td>
						@foreach(row.Values as value)
							<td class="py-2 text-right tabular-nums">{{ value }}</td>
						@endforeach
					</tr>
				@endforeach
			</tbody>
		</table>
	@empty
		<div class="mt-4">{!! Empty(.Empty) !!}</div>
	@endforelse
</article>
