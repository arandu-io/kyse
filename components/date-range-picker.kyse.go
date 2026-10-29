//go:build kyse

package components

import (
	"time"

	"github.com/arandu-io/kyse/icons"
)

@go
// DateRangePickerBehavior is the name the client behaviour is registered under
// with arandu.ui.define.
const DateRangePickerBehavior = "date-range"

// DateRangePickerProps is a period chosen in one field: "15/09/2026 –
// 07/10/2026", with a calendar that takes the first day pressed as the start
// and the second as the end.
//
// # Why one field and two values
//
// A period is one thing to the person choosing it and two columns to the
// database. The field says it as one sentence; the form submits it as two
// plain ISO days, FromName and ToName, in hidden inputs the server draws. So
// the controller reads what it always read, a day at a time, and nothing in the
// request says a picker was involved.
//
// # Where the words and the dates come from
//
// Month names, column heads, the first day of the week and the shape the field
// writes a date in (DateFormat) arrive as data, for the reason
// [CalendarProps] gives: a calendar is where a locale is least negotiable, and
// the script that moves between months computes numbers and never a name.
// Empty draws English, and ISO days in the field.
//
// # What a request is
//
// The second day pressed, a shortcut and Clear are the change. Both inputs are
// written first and one change event follows, so a form that asks the server on
// change asks once and never with half a period. Closing the panel by Escape or
// by pressing outside writes nothing.
//
// # Where the panel sits
//
// It is placed from the field -- under it, from its start edge, turned over
// when there is no room below, kept inside the window -- and is not clipped by
// a dialog or a card around it. On a narrow screen it is a sheet at the bottom
// and shows one month; the arrows move one month either way on every width.
//
// It publishes root, label, group, trigger, value, from, to, backdrop, panel,
// presets, preset, months, month-previous, month-next, month, month-title,
// month-grid, weekday, week, day, actions, status, clear, message and hint.
// None of them is called previous, next or footer: a data table publishes those
// for its pager, and a period field inside its toolbar must not be found when
// the pager is looked up by part.
type DateRangePickerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is what the field is called when FromName and ToName say nothing:
	// the two days submit as Name_from and Name_to.
	Name string
	// FromName and ToName are the names the two days submit under. A list that
	// already reads "de" and "ate" says so here and keeps its addresses.
	FromName string
	ToName   string
	// ID is the id the markup hangs off. Empty uses Name. Two fields on one
	// page need one each.
	ID string
	// Label is the text above the field.
	Label string
	// From and To are the period shown, as "2026-09-15". Either empty shows the
	// placeholder: a period is both ends or neither.
	From string
	To   string
	// Min and Max are the days that can be chosen, in the same form. A day
	// outside them is drawn and cannot be pressed.
	Min string
	Max string
	// Month is the first of the two months drawn, as "2026-09", when no period
	// is held. Empty draws the month of From, and otherwise the current one.
	Month string
	// DateFormat is the shape the field writes a day in, as a Go reference
	// layout: "02/01/2006" for the Brazilian one. Empty writes "2006-01-02".
	DateFormat string
	// Placeholder is the field while it holds no period: "Choose a period".
	Placeholder string
	// PanelLabel names the panel for a screen reader. Empty uses Placeholder.
	PanelLabel string
	// Hint is the sentence under the field while nothing is wrong.
	Hint string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
	// Presets are the shortcuts above the calendar: "Today", "Last 7 days".
	// None draws none.
	Presets []DateRangePreset
	// MonthNames, WeekdayNames and FirstDay are handed to the grids. See
	// CalendarProps, where the reason they are data is written out.
	MonthNames   []string
	WeekdayNames []string
	FirstDay     time.Weekday
	// PreviousLabel and NextLabel name the two arrows. Empty says "Previous
	// month" and "Next month".
	PreviousLabel string
	NextLabel     string
	// PendingLabel is said while the start is chosen and the end is not:
	// "Choose the end date". Empty says it in English.
	PendingLabel string
	// ClearLabel is the control that empties the period. Empty says "Clear".
	ClearLabel string
}

// DateRangePreset is a shortcut: the last Days days, counting today.
type DateRangePreset struct {
	// Label is what the shortcut is called.
	Label string
	// Days is how many days the period spans, ending today: 1 is today, 7 is
	// the last week, 30 the last thirty days.
	Days int
}

// DateRangeDay is one cell of a month the picker draws.
type DateRangeDay struct {
	// Date is the day, as "2026-09-08".
	Date string
	// Number is the day of the month, which is what is drawn.
	Number int
	// Outside is a day of a neighbouring month. It is drawn empty, because the
	// two months side by side already hold it.
	Outside bool
	// Today marks the day the page was drawn on.
	Today bool
	// Disabled is a day outside Min and Max.
	Disabled bool
	// Kind is where the day stands in the period: "start", "middle", "end", or
	// "single" when the period is one day. Empty is outside it.
	Kind string
	// Stop holds the tab stop.
	Stop bool
}

// DateRangeWeek is one row of a month: seven days.
type DateRangeWeek struct {
	// Days are the seven cells of the row.
	Days []DateRangeDay
	// Empty is a row with no day of the month in it, which is the sixth of a
	// month that fits in five. It is drawn hidden, so the grid is not a row of
	// nothing taller than it has to be, and the script shows it again when
	// the month it moves to needs it.
	Empty bool
}

// DateRangeMonth is one month of the pair, as a grid of six rows.
type DateRangeMonth struct {
	// Title is the month and year as they head the grid.
	Title string
	// Weeks are the six rows of seven days.
	Weeks []DateRangeWeek
}

// ElementID is the id the markup carries: ID, or Name when ID is empty.
func (p DateRangePickerProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// FromField is the name the start submits under.
func (p DateRangePickerProps) FromField() string {
	if p.FromName != "" {
		return p.FromName
	}
	return p.Name + "_from"
}

// ToField is the name the end submits under.
func (p DateRangePickerProps) ToField() string {
	if p.ToName != "" {
		return p.ToName
	}
	return p.Name + "_to"
}

// FromValue is the start shown: what came back rejected, or what the caller
// set.
func (p DateRangePickerProps) FromValue() string {
	if p.Page == nil {
		return p.From
	}
	return p.Page.OldOr(p.FromField(), p.From)
}

// ToValue is the end shown: what came back rejected, or what the caller set.
func (p DateRangePickerProps) ToValue() string {
	if p.Page == nil {
		return p.To
	}
	return p.Page.OldOr(p.ToField(), p.To)
}

// span is the period as two days in order, and whether there is one. A start
// after the end is the same period read backwards, and is drawn forwards.
func (p DateRangePickerProps) span() (time.Time, time.Time, bool) {
	from, errFrom := time.Parse("2006-01-02", p.FromValue())
	to, errTo := time.Parse("2006-01-02", p.ToValue())
	if errFrom != nil || errTo != nil {
		return time.Time{}, time.Time{}, false
	}
	if to.Before(from) {
		from, to = to, from
	}
	return from, to, true
}

// Layout is the shape of a day in the field.
func (p DateRangePickerProps) Layout() string {
	if p.DateFormat != "" {
		return p.DateFormat
	}
	return "2006-01-02"
}

// Text is what the field says: the period in the shape asked for, or the
// placeholder.
func (p DateRangePickerProps) Text() string {
	from, to, ok := p.span()
	if !ok {
		return p.PlaceholderText()
	}
	return from.Format(p.Layout()) + " – " + to.Format(p.Layout())
}

// PlaceholderText is what the field says while it holds no period.
func (p DateRangePickerProps) PlaceholderText() string {
	if p.Placeholder != "" {
		return p.Placeholder
	}
	return "Choose a period"
}

// Empty is whether the field holds no period.
func (p DateRangePickerProps) Empty() bool {
	_, _, ok := p.span()
	return !ok
}

// PanelName is what the panel is called.
func (p DateRangePickerProps) PanelName() string {
	if p.PanelLabel != "" {
		return p.PanelLabel
	}
	return p.PlaceholderText()
}

// Message is the rejection for either end, and empty when there is none.
func (p DateRangePickerProps) Message() string {
	if p.Page == nil {
		return ""
	}
	if message := p.Page.FieldError(p.FromField()); message != "" {
		return message
	}
	return p.Page.FieldError(p.ToField())
}

// DescribedBy is the id of whatever is explaining the field.
func (p DateRangePickerProps) DescribedBy() string {
	if p.Message() != "" {
		return p.ElementID() + "-error"
	}
	if p.Hint != "" {
		return p.ElementID() + "-hint"
	}
	return ""
}

// TriggerID, PanelID and BackdropID are the ids the field, the panel and the
// sheet's backdrop carry.
func (p DateRangePickerProps) TriggerID() string { return p.ElementID() + "-trigger" }

// PanelID is the id of the element that opens.
func (p DateRangePickerProps) PanelID() string { return p.ElementID() + "-panel" }

// calendar is the grid's vocabulary, built from this picker's own fields so the
// two cannot disagree about the language or the bounds.
func (p DateRangePickerProps) calendar(month string) CalendarProps {
	return CalendarProps{
		Month:        month,
		Min:          p.Min,
		Max:          p.Max,
		MonthNames:   p.MonthNames,
		WeekdayNames: p.WeekdayNames,
		FirstDay:     p.FirstDay,
	}
}

// First is the first of the two months drawn.
func (p DateRangePickerProps) First() time.Time {
	if month, err := time.Parse("2006-01", p.Month); err == nil {
		return month
	}
	if from, _, ok := p.span(); ok {
		return time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	return p.calendar("").Visible()
}

// MonthValue is the first month as "2026-09", which is what the behaviour steps
// from.
func (p DateRangePickerProps) MonthValue() string { return p.First().Format("2006-01") }

// Columns are the column heads in the order they are drawn.
func (p DateRangePickerProps) Columns() []string { return p.calendar("").Columns() }

// stop is the day that holds the tab stop: the start, or today when it is drawn,
// or the first of the first month.
func (p DateRangePickerProps) stop(months []time.Time) string {
	candidates := []string{}
	if from, _, ok := p.span(); ok {
		candidates = append(candidates, from.Format("2006-01-02"))
	}
	candidates = append(candidates, time.Now().UTC().Format("2006-01-02"))
	for _, candidate := range candidates {
		for _, month := range months {
			if len(candidate) >= 7 && candidate[:7] == month.Format("2006-01") {
				return candidate
			}
		}
	}
	return months[0].Format("2006-01-02")
}

// Grids are the two months side by side, each of six full weeks.
//
// Two, always: the second one is hidden by the stylesheet on a narrow screen,
// so the markup does not depend on the width it is read at.
func (p DateRangePickerProps) Grids() []DateRangeMonth {
	first := p.First()
	months := []time.Time{first, first.AddDate(0, 1, 0)}
	from, to, held := p.span()
	today := time.Now().UTC().Format("2006-01-02")
	stop := p.stop(months)

	out := make([]DateRangeMonth, 0, 2)
	for _, month := range months {
		calendar := p.calendar(month.Format("2006-01"))
		weeks := make([]DateRangeWeek, 0, 6)
		for _, week := range calendar.Weeks() {
			row := make([]DateRangeDay, 0, 7)
			empty := true
			for _, cell := range week {
				if !cell.Outside {
					empty = false
				}
				day := DateRangeDay{Date: cell.Date, Number: cell.Number, Outside: cell.Outside}
				if !cell.Outside {
					day.Today = cell.Date == today
					day.Disabled = cell.Disabled
					day.Stop = cell.Date == stop
					if held {
						at, _ := time.Parse("2006-01-02", cell.Date)
						switch {
						case at.Before(from) || at.After(to):
						case from.Equal(to):
							day.Kind = "single"
						case at.Equal(from):
							day.Kind = "start"
						case at.Equal(to):
							day.Kind = "end"
						default:
							day.Kind = "middle"
						}
					}
				}
				row = append(row, day)
			}
			weeks = append(weeks, DateRangeWeek{Days: row, Empty: empty})
		}
		out = append(out, DateRangeMonth{Title: calendar.Title(), Weeks: weeks})
	}
	return out
}

// PreviousName and NextName are what the two arrows are called.
func (p DateRangePickerProps) PreviousName() string {
	if p.PreviousLabel != "" {
		return p.PreviousLabel
	}
	return "Previous month"
}

// NextName is what the forward arrow is called.
func (p DateRangePickerProps) NextName() string {
	if p.NextLabel != "" {
		return p.NextLabel
	}
	return "Next month"
}

// ClearText is the control that empties the period.
func (p DateRangePickerProps) ClearText() string {
	if p.ClearLabel != "" {
		return p.ClearLabel
	}
	return "Clear"
}

// PendingText is said while the end is still to be chosen.
func (p DateRangePickerProps) PendingText() string {
	if p.PendingLabel != "" {
		return p.PendingLabel
	}
	return "Choose the end date"
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p DateRangePickerProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: DateRangePickerBehavior, Props: map[string]any{
			"months":      p.calendar("").Months(),
			"firstDay":    int(p.FirstDay),
			"min":         p.Min,
			"max":         p.Max,
			"format":      p.Layout(),
			"placeholder": p.PlaceholderText(),
			"pending":     p.PendingText(),
			"month":       p.Month,
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p DateRangePickerProps) PartNames() []string {
	return []string{
		"root", "label", "group", "trigger", "value", "from", "to", "backdrop", "panel",
		"presets", "preset", "months", "month-previous", "month-next", "month",
		"month-title", "month-grid", "weekday", "week", "day", "actions", "status",
		"clear", "message", "hint",
	}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("field date-range") }}"
	data-month="{{ .MonthValue() }}"
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		for="{{ .TriggerID() }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	<div
		data-part="group"
		class="{{ .PartClass("group", "date-range-box") }}"
		@attributes(.PartAttrs("group"))
	>
		{{-- A button, not a text box: the period is chosen in the calendar and
		     never typed, so there is no half-written date for a form to send.
		     Without the script it is still a field that shows the period and
		     submits it, because the value is in the two hidden inputs below. --}}
		<button
			data-part="trigger"
			class="{{ .PartClass("trigger", "btn date-range-trigger") }}"
			type="button"
			data-variant="outline"
			id="{{ .TriggerID() }}"
			aria-haspopup="dialog"
			aria-expanded="false"
			aria-controls="{{ .PanelID() }}"
			@if(.DescribedBy() != "")
				aria-describedby="{{ .DescribedBy() }}"
			@endif
			@if(.Message() != "")
				aria-invalid="true"
			@endif
			@if(.Required)
				aria-required="true"
			@endif
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("trigger"))
		>
			{!! icons.CalendarBlank(icons.Props{}) !!}
			<span
				data-part="value"
				class="{{ .PartClass("value", "date-range-value") }}"
				@if(.Empty())
					data-empty="true"
				@endif
				@attributes(.PartAttrs("value"))
			>{{ .Text() }}</span>
		</button>

		<input
			data-part="from"
			type="hidden"
			@if(.PartClass("from") != "")
				class="{{ .PartClass("from") }}"
			@endif
			name="{{ .FromField() }}"
			value="{{ .FromValue() }}"
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("from"))
		>
		<input
			data-part="to"
			type="hidden"
			@if(.PartClass("to") != "")
				class="{{ .PartClass("to") }}"
			@endif
			name="{{ .ToField() }}"
			value="{{ .ToValue() }}"
			@if(.Disabled)
				disabled
			@endif
			@attributes(.PartAttrs("to"))
		>

		<div
			data-part="backdrop"
			class="{{ .PartClass("backdrop", "date-range-backdrop") }}"
			hidden
			@attributes(.PartAttrs("backdrop"))
		></div>

		<div
			data-part="panel"
			class="{{ .PartClass("panel", "date-range-panel") }}"
			id="{{ .PanelID() }}"
			role="dialog"
			aria-label="{{ .PanelName() }}"
			hidden
			@attributes(.PartAttrs("panel"))
		>
			@if(len(.Presets) > 0)
				<div
					data-part="presets"
					class="{{ .PartClass("presets", "date-range-presets") }}"
					@attributes(.PartAttrs("presets"))
				>
					@foreach(.Presets as preset)
						<button
							data-part="preset"
							class="{{ .PartClass("preset", "btn") }}"
							type="button"
							data-variant="outline"
							data-size="sm"
							data-range-days="{{ preset.Days }}"
							@attributes(.PartAttrs("preset"))
						>{{ preset.Label }}</button>
					@endforeach
				</div>
			@endif

			<div
				data-part="months"
				class="{{ .PartClass("months", "date-range-months") }}"
				@attributes(.PartAttrs("months"))
			>
				<button
					data-part="month-previous"
					class="{{ .PartClass("month-previous", "btn date-range-nav") }}"
					type="button"
					data-variant="ghost"
					data-size="sm"
					data-range-step="-1"
					aria-label="{{ .PreviousName() }}"
					@attributes(.PartAttrs("month-previous"))
				>{!! icons.CaretLeft(icons.Props{}) !!}</button>
				<button
					data-part="month-next"
					class="{{ .PartClass("month-next", "btn date-range-nav") }}"
					type="button"
					data-variant="ghost"
					data-size="sm"
					data-range-step="1"
					aria-label="{{ .NextName() }}"
					@attributes(.PartAttrs("month-next"))
				>{!! icons.CaretRight(icons.Props{}) !!}</button>

				@foreach(.Grids() as month)
					<section
						data-part="month"
						class="{{ .PartClass("month", "date-range-month") }}"
						@attributes(.PartAttrs("month"))
					>
						<h2
							data-part="month-title"
							class="{{ .PartClass("month-title", "date-range-title") }}"
							aria-live="polite"
							@attributes(.PartAttrs("month-title"))
						>{{ month.Title }}</h2>
						<table
							data-part="month-grid"
							class="{{ .PartClass("month-grid", "date-range-grid") }}"
							role="grid"
							aria-label="{{ month.Title }}"
							@attributes(.PartAttrs("month-grid"))
						>
							<thead>
								<tr>
									@foreach(.Columns() as column)
										<th
											data-part="weekday"
											class="{{ .PartClass("weekday", "date-range-weekday") }}"
											scope="col"
											@attributes(.PartAttrs("weekday"))
										>{{ column }}</th>
									@endforeach
								</tr>
							</thead>
							<tbody>
								@foreach(month.Weeks as week)
									<tr
										data-part="week"
										@if(week.Empty)
											hidden
										@endif
										@if(.PartClass("week") != "")
											class="{{ .PartClass("week") }}"
										@endif
										@attributes(.PartAttrs("week"))
									>
										@foreach(week.Days as day)
											<td
												data-part="day"
												class="{{ .PartClass("day", "date-range-day") }}"
												role="gridcell"
												data-date="{{ day.Date }}"
												@if(day.Outside)
													data-outside="true"
												@endif
												@if(day.Today)
													data-today="true"
												@endif
												@if(day.Kind != "")
													data-range="{{ day.Kind }}"
													aria-selected="true"
												@endif
												@if(day.Kind == "")
													aria-selected="false"
												@endif
												@if(day.Disabled)
													aria-disabled="true"
												@endif
												@if(day.Stop)
													tabindex="0"
												@endif
												@if(!day.Stop)
													tabindex="-1"
												@endif
												@attributes(.PartAttrs("day"))
											>
												@if(!day.Outside)
													{{ day.Number }}
												@endif
											</td>
										@endforeach
									</tr>
								@endforeach
							</tbody>
						</table>
					</section>
				@endforeach
			</div>

			<footer
				data-part="actions"
				class="{{ .PartClass("actions", "date-range-actions") }}"
				@attributes(.PartAttrs("actions"))
			>
				<p
					data-part="status"
					class="{{ .PartClass("status", "date-range-status") }}"
					role="status"
					aria-live="polite"
					@attributes(.PartAttrs("status"))
				></p>
				<button
					data-part="clear"
					class="{{ .PartClass("clear", "btn") }}"
					type="button"
					data-variant="ghost"
					data-size="sm"
					data-range-clear
					@attributes(.PartAttrs("clear"))
				>{{ .ClearText() }}</button>
			</footer>
		</div>
	</div>

	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .ElementID() }}-error"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p
				data-part="hint"
				id="{{ .ElementID() }}-hint"
				class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("hint"))
			>{{ .Hint }}</p>
		@endif
	@endif
</div>
