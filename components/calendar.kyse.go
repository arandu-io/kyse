//go:build kyse

package components

import (
	"strconv"
	"time"

	"github.com/arandu-io/kyse/icons"
)

@go
// CalendarBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const CalendarBehavior = "calendar"

// CalendarProps is a month, drawn as a grid, with the days that can be chosen.
//
// # Why the names arrive as data
//
// A calendar is where a locale is least negotiable: the order of the columns,
// the day the week starts on, what the months are called. A component that
// guessed would be wrong in exactly the places where being wrong is expensive,
// and Go's time package cannot localise a month name at all.
//
// So it takes them. MonthNames and WeekdayNames come from the application's own
// catalogue, which is the one thing on the page that already knows the reader's
// language. Empty draws English and says so, rather than pretending.
//
// That is the same bargain [RatingProps.OptionLabels] makes, and it is what
// separates this from a script reading the browser's locale: the calendar and
// the rest of the page then agree, because they are translated from one file.
//
// # What the client does, and does not
//
// Moving between months is done in the browser, with the names it was handed
// -- so a month change costs nothing and cannot disagree with the first month
// the server drew. Nothing here fetches, and nothing computes a name.
//
// It publishes root, header, title, previous, next, grid, weekday, week, day
// and footer.
type CalendarProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the grid and the controls hang their ids off. Two calendars
	// on one page need two.
	ID string
	// Label names the calendar for a screen reader: "Choose a date".
	Label string
	// Month is the month drawn, as "2026-09". Empty draws the month of Value,
	// and the current month when there is no value either.
	Month string
	// Value is the chosen day, as "2026-09-08". Empty chooses nothing.
	Value string
	// Min and Max are the range that can be chosen, in the same form. A day
	// outside it is drawn and cannot be picked, rather than hidden -- a
	// calendar with holes in it is a calendar nobody can read.
	Min string
	Max string
	// MonthNames are the twelve, from January. Fewer than twelve draws
	// English, because a partial list is a list that would name some months
	// and number others.
	MonthNames []string
	// WeekdayNames are the seven, from Sunday, as they head the columns: "Su",
	// "Mo". They are short by convention and this does not shorten them --
	// what fits is a decision about the language.
	WeekdayNames []string
	// FirstDay is the day the week starts on. Zero is Sunday, which is the
	// zero value of time.Weekday and is wrong for most of the world -- so a
	// calendar drawn without it is a calendar that has not been told.
	FirstDay time.Weekday
	// TodayLabel and ClearLabel are the two controls under the grid. Empty
	// says "Today" and "Clear"; neither is drawn when the label is a single
	// space, which is how a calendar asks for no footer at all.
	TodayLabel string
	ClearLabel string
	// PreviousLabel and NextLabel name the two arrows. Empty says "Previous
	// month" and "Next month".
	PreviousLabel string
	NextLabel     string
	// Target is the id of the input this calendar writes into. Empty writes
	// nowhere and leaves the day as a value the caller reads off the grid,
	// which is what a calendar used for display rather than for choosing is.
	Target string
}

// CalendarDay is one cell of the grid.
type CalendarDay struct {
	// Date is the day, as "2026-09-08". It is what a cell carries and what is
	// written into the target.
	Date string
	// Number is the day of the month, which is what is drawn.
	Number int
	// Outside is whether the day belongs to a neighbouring month. Those cells
	// are drawn rather than left blank, so the grid is always six full weeks
	// and the calendar does not change height when the month does.
	Outside bool
	// Today marks the day the page was drawn on.
	Today bool
	// Selected marks the chosen day.
	Selected bool
	// Disabled is a day outside Min and Max.
	Disabled bool
}

// englishMonths and englishWeekdays are what is drawn when the caller named
// nothing. They are not a default anybody should ship: they are what makes a
// calendar readable while somebody is still wiring the catalogue up.
var englishMonths = []string{
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

var englishWeekdays = []string{"Su", "Mo", "Tu", "We", "Th", "Fr", "Sa"}

// Months are the twelve names, the caller's or English.
func (p CalendarProps) Months() []string {
	if len(p.MonthNames) == 12 {
		return p.MonthNames
	}
	return englishMonths
}

// Weekdays are the seven column heads, from Sunday, the caller's or English.
func (p CalendarProps) Weekdays() []string {
	if len(p.WeekdayNames) == 7 {
		return p.WeekdayNames
	}
	return englishWeekdays
}

// Columns are the column heads in the order they are drawn, which is the seven
// names rotated to start on FirstDay.
func (p CalendarProps) Columns() []string {
	names := p.Weekdays()
	start := int(p.FirstDay) % 7
	out := make([]string, 0, 7)
	for at := 0; at < 7; at++ {
		out = append(out, names[(start+at)%7])
	}
	return out
}

// Visible is the month drawn: the one asked for, the one the value is in, or
// the one the page was drawn in.
func (p CalendarProps) Visible() time.Time {
	if month, err := time.Parse("2006-01", p.Month); err == nil {
		return month
	}
	if day, err := time.Parse("2006-01-02", p.Value); err == nil {
		return time.Date(day.Year(), day.Month(), 1, 0, 0, 0, 0, time.UTC)
	}
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
}

// Title is the month and year as they head the grid.
func (p CalendarProps) Title() string {
	visible := p.Visible()
	return p.Months()[int(visible.Month())-1] + " " + strconv.Itoa(visible.Year())
}

// MonthValue is the drawn month as "2026-09", which is what the behaviour
// steps from.
func (p CalendarProps) MonthValue() string { return p.Visible().Format("2006-01") }

// Weeks are the six rows of the grid, each of seven days.
//
// Six always, and full always. A grid that grew a row in a long month and lost
// it in a short one would move everything under it every time the month
// changed, and a month drawn with blanks at the ends is a month whose first
// row reads as a week that started on a Wednesday.
func (p CalendarProps) Weeks() [][]CalendarDay {
	visible := p.Visible()
	offset := (int(visible.Weekday()) - int(p.FirstDay) + 7) % 7
	start := visible.AddDate(0, 0, -offset)

	today := time.Now().UTC().Format("2006-01-02")
	weeks := make([][]CalendarDay, 0, 6)
	for week := 0; week < 6; week++ {
		row := make([]CalendarDay, 0, 7)
		for column := 0; column < 7; column++ {
			day := start.AddDate(0, 0, week*7+column)
			date := day.Format("2006-01-02")
			row = append(row, CalendarDay{
				Date:     date,
				Number:   day.Day(),
				Outside:  day.Month() != visible.Month(),
				Today:    date == today,
				Selected: date == p.Value,
				Disabled: p.outOfRange(date),
			})
		}
		weeks = append(weeks, row)
	}
	return weeks
}

// outOfRange is whether a day falls outside Min and Max. The comparison is on
// the string, which is exact for this format and is why the format is fixed.
func (p CalendarProps) outOfRange(date string) bool {
	if p.Min != "" && date < p.Min {
		return true
	}
	return p.Max != "" && date > p.Max
}

// GridID is the id of the grid, which the controls point at.
func (p CalendarProps) GridID() string { return p.ID + "-grid" }

// TitleID is the id of the heading, which names the grid.
func (p CalendarProps) TitleID() string { return p.ID + "-title" }

// PreviousName and NextName are what the two arrows are called.
func (p CalendarProps) PreviousName() string {
	if p.PreviousLabel != "" {
		return p.PreviousLabel
	}
	return "Previous month"
}

// NextName is what the forward arrow is called.
func (p CalendarProps) NextName() string {
	if p.NextLabel != "" {
		return p.NextLabel
	}
	return "Next month"
}

// TodayText and ClearText are the two controls under the grid. A single space
// draws neither, which is how a calendar asks for no footer.
func (p CalendarProps) TodayText() string {
	if p.TodayLabel != "" {
		return p.TodayLabel
	}
	return "Today"
}

// ClearText is the control that unsets the day.
func (p CalendarProps) ClearText() string {
	if p.ClearLabel != "" {
		return p.ClearLabel
	}
	return "Clear"
}

// HasFooter is whether the two controls are drawn at all.
func (p CalendarProps) HasFooter() bool {
	return p.TodayText() != " " && p.ClearText() != " "
}

// Stop is whether this day holds the tab stop: the chosen one, or today, or
// the first of the month. A grid is one tab stop and the arrow keys move
// inside it, so exactly one cell can be reached by tabbing.
func (p CalendarProps) Stop(day CalendarDay) bool {
	if p.Value != "" {
		return day.Date == p.Value
	}
	visible := p.Visible()
	if today := time.Now().UTC(); today.Year() == visible.Year() && today.Month() == visible.Month() {
		return day.Date == today.Format("2006-01-02")
	}
	return !day.Outside && day.Number == 1
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
//
// The names go over with it. The behaviour redraws a month without asking the
// server and without computing a name, so what it draws and what the server
// drew are the same words.
func (p CalendarProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: CalendarBehavior, Props: map[string]any{
			"months":   p.Months(),
			"weekdays": p.Columns(),
			"firstDay": int(p.FirstDay),
			"min":      p.Min,
			"max":      p.Max,
			"target":   p.Target,
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p CalendarProps) PartNames() []string {
	return []string{
		"root", "header", "title", "previous", "next",
		"grid", "weekday", "week", "day", "footer",
	}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("calendar") }}"
	id="{{ .ID }}"
	data-month="{{ .MonthValue() }}"
	@if(.Value != "")
		data-value="{{ .Value }}"
	@endif
	@attributes(.RootAttrs())
>
	<header
		data-part="header"
		class="{{ .PartClass("header", "calendar-header") }}"
		@attributes(.PartAttrs("header"))
	>
		<button
			data-part="previous"
			class="{{ .PartClass("previous", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			data-calendar-step="-1"
			aria-label="{{ .PreviousName() }}"
			aria-controls="{{ .GridID() }}"
			@attributes(.PartAttrs("previous"))
		>{!! icons.CaretLeft(icons.Props{}) !!}</button>

		<h2
			data-part="title"
			class="{{ .PartClass("title", "calendar-title") }}"
			id="{{ .TitleID() }}"
			aria-live="polite"
			@attributes(.PartAttrs("title"))
		>{{ .Title() }}</h2>

		<button
			data-part="next"
			class="{{ .PartClass("next", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			data-calendar-step="1"
			aria-label="{{ .NextName() }}"
			aria-controls="{{ .GridID() }}"
			@attributes(.PartAttrs("next"))
		>{!! icons.CaretRight(icons.Props{}) !!}</button>
	</header>

	{{-- A table, because a month is tabular data: the column a day sits in is
	     what says which weekday it is, and a screen reader announces that from
	     the column head rather than from anything written on the cell. --}}
	<table
		data-part="grid"
		class="{{ .PartClass("grid", "calendar-grid") }}"
		id="{{ .GridID() }}"
		role="grid"
		aria-labelledby="{{ .TitleID() }}"
		@if(.Label != "")
			aria-label="{{ .Label }}"
		@endif
		@attributes(.PartAttrs("grid"))
	>
		<thead>
			<tr>
				@foreach(.Columns() as column)
					<th
						data-part="weekday"
						class="{{ .PartClass("weekday", "calendar-weekday") }}"
						scope="col"
						@attributes(.PartAttrs("weekday"))
					>{{ column }}</th>
				@endforeach
			</tr>
		</thead>
		<tbody>
			@foreach(.Weeks() as week)
				<tr
					data-part="week"
					@if(.PartClass("week") != "")
						class="{{ .PartClass("week") }}"
					@endif
					@attributes(.PartAttrs("week"))
				>
					@foreach(week as day)
						<td
							data-part="day"
							class="{{ .PartClass("day", "calendar-day") }}"
							role="gridcell"
							data-date="{{ day.Date }}"
							@attributes(.PartAttrs("day"))
							@if(day.Outside)
								data-outside="true"
							@endif
							@if(day.Today)
								data-today="true"
							@endif
							@if(day.Selected)
								aria-selected="true"
							@endif
							@if(!day.Selected)
								aria-selected="false"
							@endif
							@if(day.Disabled)
								aria-disabled="true"
							@endif
							@if(.Stop(day))
								tabindex="0"
							@endif
							@if(!.Stop(day))
								tabindex="-1"
							@endif
						>{{ day.Number }}</td>
					@endforeach
				</tr>
			@endforeach
		</tbody>
	</table>

	@if(.HasFooter())
		<footer
			data-part="footer"
			class="{{ .PartClass("footer", "calendar-footer") }}"
			@attributes(.PartAttrs("footer"))
		>
			<button type="button" class="btn" data-variant="ghost" data-size="sm" data-calendar-clear>{{ .ClearText() }}</button>
			<button type="button" class="btn" data-variant="ghost" data-size="sm" data-calendar-today>{{ .TodayText() }}</button>
		</footer>
	@endif
</div>
