package unit

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/arandu-io/kyse/components"
)

func rangeMarkup(p components.DateRangePickerProps) string {
	return string(components.DateRangePicker(p))
}

// Day cells of the markup, by date and state: the marks a day carries are what
// the stylesheet paints the band from.
var rangeDay = regexp.MustCompile(`<td\b[^>]*data-date="([0-9-]+)"[^>]*>`)

func dayTag(t *testing.T, html, date string) string {
	t.Helper()
	for _, m := range rangeDay.FindAllStringSubmatch(html, -1) {
		if m[1] == date && !strings.Contains(m[0], `data-outside="true"`) {
			return m[0]
		}
	}
	t.Fatalf("no cell for %s", date)
	return ""
}

// A period is one field and two values: the hidden inputs carry the ISO days under
// the names the caller gave, so a controller reads what it always read.
func TestAPeriodSubmitsTwoIsoDaysUnderTheNamesItWasGiven(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{
		Name: "period", FromName: "de", ToName: "ate", Label: "Período",
		From: "2026-09-15", To: "2026-10-07", DateFormat: "02/01/2006",
	})
	for _, want := range []string{
		`type="hidden"`, `name="de"`, `value="2026-09-15"`, `name="ate"`, `value="2026-10-07"`,
		`15/09/2026 – 07/10/2026`, `for="period-trigger"`, `id="period-trigger"`,
		`data-kyse-behavior="date-range"`, `aria-haspopup="dialog"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the period field draws %q", want)
		}
	}
	for _, not := range []string{`name="period_from"`, `name="period_to"`, `data-empty="true"`} {
		if strings.Contains(html, not) {
			t.Errorf("the period field draws %q", not)
		}
	}

	byDefault := rangeMarkup(components.DateRangePickerProps{Name: "period", Label: "Period", From: "2026-09-15", To: "2026-10-07"})
	for _, want := range []string{`name="period_from"`, `name="period_to"`, `2026-09-15 – 2026-10-07`} {
		if !strings.Contains(byDefault, want) {
			t.Errorf("by default the field draws %q", want)
		}
	}
}

// With no period the field says so and submits nothing: a period is both ends or
// neither, so one end alone is the placeholder too.
func TestAnEmptyPeriodShowsThePlaceholder(t *testing.T) {
	for name, props := range map[string]components.DateRangePickerProps{
		"none":      {Name: "p", Label: "P", Placeholder: "Escolher período"},
		"only from": {Name: "p", Label: "P", Placeholder: "Escolher período", From: "2026-09-15"},
		"bad day":   {Name: "p", Label: "P", Placeholder: "Escolher período", From: "15/09/2026", To: "07/10/2026"},
	} {
		html := rangeMarkup(props)
		if !strings.Contains(html, `data-empty="true"`) || !strings.Contains(html, `>Escolher período</span>`) {
			t.Errorf("%s: the field shows the placeholder", name)
		}
		if strings.Contains(html, `data-range=`) {
			t.Errorf("%s: no day is inside a period", name)
		}
	}
}

// The period is drawn as a band: the first day, the last, and the ones between,
// across the two months; a period of one day is a single mark.
func TestThePeriodIsDrawnAsABandWithItsEndsMarked(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{Name: "p", Label: "P", From: "2026-09-15", To: "2026-10-07"})
	for date, want := range map[string]string{
		"2026-09-15": "start", "2026-09-16": "middle", "2026-09-30": "middle", "2026-10-01": "middle", "2026-10-07": "end",
	} {
		if tag := dayTag(t, html, date); !strings.Contains(tag, `data-range="`+want+`"`) || !strings.Contains(tag, `aria-selected="true"`) {
			t.Errorf("%s is %q: %s", date, want, tag)
		}
	}
	for _, date := range []string{"2026-09-14", "2026-10-08"} {
		if tag := dayTag(t, html, date); strings.Contains(tag, "data-range=") || !strings.Contains(tag, `aria-selected="false"`) {
			t.Errorf("%s is outside the period: %s", date, tag)
		}
	}

	if tag := dayTag(t, rangeMarkup(components.DateRangePickerProps{Name: "p", From: "2026-09-15", To: "2026-09-15"}), "2026-09-15"); !strings.Contains(tag, `data-range="single"`) {
		t.Errorf("a one-day period is a single mark: %s", tag)
	}
	// An end before the start is the same period read backwards.
	if tag := dayTag(t, rangeMarkup(components.DateRangePickerProps{Name: "p", From: "2026-10-07", To: "2026-09-15"}), "2026-09-15"); !strings.Contains(tag, `data-range="start"`) {
		t.Errorf("a period read backwards is drawn forwards: %s", tag)
	}
}

// Two months always, the first the one the period starts in; each with six
// rows, the sixth hidden when it holds no day of the month. A day of a
// neighbouring month is drawn empty, because the pair already holds it.
func TestTwoMonthsAreDrawnAndAnEmptyRowIsHidden(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{
		Name: "p", Label: "P", From: "2026-09-15", To: "2026-10-07",
		MonthNames:   []string{"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho", "Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro"},
		WeekdayNames: []string{"Dom", "Seg", "Ter", "Qua", "Qui", "Sex", "Sáb"},
	})
	if n := strings.Count(html, `data-part="month"`); n != 2 {
		t.Fatalf("the picker draws %d months, want 2", n)
	}
	for _, want := range []string{"Setembro 2026", "Outubro 2026", `data-month="2026-09"`, ">Dom</th>", ">Sáb</th>"} {
		if !strings.Contains(html, want) {
			t.Errorf("the picker draws %q", want)
		}
	}
	if n := strings.Count(html, `data-part="week"`); n != 12 {
		t.Errorf("the picker draws %d rows, want 6 a month", n)
	}
	// September 2026 starts on a Tuesday and ends on a Wednesday: five rows. The
	// sixth row of each month is hidden when it is all outside.
	if n := len(regexp.MustCompile(`<tr\s+data-part="week"\s+hidden`).FindAllString(html, -1)); n < 1 {
		t.Error("a row with no day of the month is hidden")
	}
	if tag := rangeDay.FindAllString(html, -1); len(tag) != 84 {
		t.Errorf("the picker draws %d cells, want 84", len(tag))
	}
	// Outside cells carry no number: they would show the same day twice.
	if regexp.MustCompile(`data-outside="true"[^>]*>\s*[0-9]`).MatchString(html) {
		t.Error("a day of a neighbouring month is drawn with a number")
	}
}

// Min and Max bound what can be chosen: a day outside is drawn and disabled.
func TestDaysOutsideTheBoundsAreDisabled(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{Name: "p", Month: "2026-09", Min: "2026-09-10", Max: "2026-10-20"})
	if tag := dayTag(t, html, "2026-09-09"); !strings.Contains(tag, `aria-disabled="true"`) {
		t.Errorf("a day before Min is disabled: %s", tag)
	}
	if tag := dayTag(t, html, "2026-10-21"); !strings.Contains(tag, `aria-disabled="true"`) {
		t.Errorf("a day after Max is disabled: %s", tag)
	}
	if tag := dayTag(t, html, "2026-09-10"); strings.Contains(tag, `aria-disabled`) {
		t.Errorf("Min itself can be chosen: %s", tag)
	}
	for _, want := range []string{`"min":"2026-09-10"`, `"max":"2026-10-20"`} {
		if !strings.Contains(html, strings.ReplaceAll(want, `"`, `&#34;`)) && !strings.Contains(html, want) {
			t.Errorf("the behaviour is told %s", want)
		}
	}
}

// The grid is one tab stop: exactly one cell can be reached by tabbing.
func TestTheGridHoldsOneTabStop(t *testing.T) {
	for name, props := range map[string]components.DateRangePickerProps{
		"a period":   {Name: "p", From: "2026-09-15", To: "2026-10-07"},
		"none":       {Name: "p", Month: "2026-09"},
		"this month": {Name: "p", Month: time.Now().UTC().Format("2006-01")},
	} {
		html := rangeMarkup(props)
		stops := 0
		for _, m := range rangeDay.FindAllString(html, -1) {
			if strings.Contains(m, `tabindex="0"`) {
				stops++
			}
		}
		if stops != 1 {
			t.Errorf("%s: %d tab stops, want 1", name, stops)
		}
	}
}

// A rejected value comes back to the field: the message of either end, the days
// that were typed, and the two attributes that tie them to the trigger.
func TestARejectedPeriodShowsItsMessageAndItsOldDays(t *testing.T) {
	state := page{
		errs: map[string]string{"ate": "O fim vem antes do início."},
		old:  map[string]string{"de": "2026-10-07", "ate": "2026-09-15"},
	}
	html := rangeMarkup(components.DateRangePickerProps{Name: "p", FromName: "de", ToName: "ate", Label: "P", Page: state, From: "2020-01-01", To: "2020-01-02"})
	for _, want := range []string{
		`O fim vem antes do início.`, `aria-invalid="true"`, `aria-describedby="p-error"`, `id="p-error"`,
		`value="2026-10-07"`, `value="2026-09-15"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the rejected field draws %q", want)
		}
	}
	if strings.Contains(html, "2020-01-01") {
		t.Error("what came back beats what the caller set")
	}
}

// Shortcuts are drawn with the days they span, and none are drawn unless asked.
func TestShortcutsCarryTheirSpan(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{Name: "p", Presets: []components.DateRangePreset{{Label: "Hoje", Days: 1}, {Label: "7 dias", Days: 7}}})
	for _, want := range []string{`data-range-days="1"`, `>Hoje</button>`, `data-range-days="7"`, `>7 dias</button>`} {
		if !strings.Contains(html, want) {
			t.Errorf("the shortcuts draw %q", want)
		}
	}
	if strings.Contains(rangeMarkup(components.DateRangePickerProps{Name: "p"}), `data-range-days`) {
		t.Error("no shortcut is drawn unless one is asked for")
	}
}

// Nothing here writes a script, a handler or a style: the behaviour is a name the
// view layer's script answers to, and the panel starts closed.
func TestThePickerWritesNoScriptNoHandlerNoStyle(t *testing.T) {
	html := rangeMarkup(components.DateRangePickerProps{Name: "p", Label: "P", From: "2026-09-15", To: "2026-10-07", Presets: []components.DateRangePreset{{Label: "Hoje", Days: 1}}})
	for _, not := range []string{"<script", " style=", "onclick", "x-on:", "@click", "hx-on"} {
		if strings.Contains(html, not) {
			t.Errorf("the picker writes %q", not)
		}
	}
	if !regexp.MustCompile(`data-part="panel"[^>]*\shidden`).MatchString(html) {
		t.Error("the panel starts closed")
	}
	if !strings.Contains(html, `aria-expanded="false"`) {
		t.Error("the trigger says the panel is closed")
	}
}

// A data table publishes previous, next and footer for its pager, and a period
// field sits in its toolbar: none of the picker's parts, nor the calendar's, may
// carry those names, or a lookup by part finds the wrong element.
func TestNoPartIsCalledWhatADataTableCallsItsPager(t *testing.T) {
	taken := map[string]bool{"previous": true, "next": true, "footer": true}
	names := append(components.DateRangePickerProps{}.PartNames(), components.CalendarProps{}.PartNames()...)
	names = append(names, components.DateTimePickerProps{}.PartNames()...)
	for _, name := range names {
		if taken[name] {
			t.Errorf("a part is called %q, which the data table's pager owns", name)
		}
	}
}

// The one-day picker's group is a popover: the shipped behaviour binds to that
// class, and a group without it draws a calendar nothing opens.
func TestTheDayPickersGroupIsAPopover(t *testing.T) {
	html := string(components.DateTimePicker(components.DateTimePickerProps{Name: "d", Label: "D", Calendar: true}))
	if !strings.Contains(html, `class="date-picker popover"`) {
		t.Errorf("the group of a calendar picker is a popover: %s", regexp.MustCompile(`data-part="group"[^>]*`).FindString(html))
	}
	if !strings.Contains(html, `data-part="month-previous"`) || !strings.Contains(html, `data-part="calendar-footer"`) {
		t.Error("the calendar's arrows and footer are not named like the pager's")
	}
}

// The one-day picker names its arrows in the reader's language: the grid it
// builds takes the words from the picker, so an application does not have to
// rewrite the markup afterwards.
func TestTheDayPickerNamesItsArrows(t *testing.T) {
	html := string(components.DateTimePicker(components.DateTimePickerProps{
		Name: "d", Label: "D", Calendar: true, PreviousLabel: "Mês anterior", NextLabel: "Próximo mês",
	}))
	for _, want := range []string{`aria-label="Mês anterior"`, `aria-label="Próximo mês"`} {
		if !strings.Contains(html, want) {
			t.Errorf("the picker draws %q", want)
		}
	}
	if strings.Contains(html, "Previous month") {
		t.Error("the picker says its arrows in English when it was told otherwise")
	}
}
