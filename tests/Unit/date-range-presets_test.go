package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// Polar's nine named ranges, as a picker offers them: five count back from
// today, and four end in the past or begin on a calendar boundary. "All time" is
// the field's own Clear.
func TestAPickerOffersPolarsNamedRanges(t *testing.T) {
	html := string(components.DateRangePicker(components.DateRangePickerProps{
		Name: "p", Label: "Period",
		Presets: []components.DateRangePreset{
			{Label: "Today", Days: 1},
			{Label: "Yesterday", Days: 1, EndsAgo: 1},
			{Label: "This week", Unit: "week"},
			{Label: "This month", Unit: "month"},
			{Label: "Last month", Unit: "month", Back: 1},
			{Label: "Last 3 months", Months: 3},
			{Label: "This year", Unit: "year"},
			{Label: "Last year", Unit: "year", Back: 1},
			{Label: "Last 7 days", Days: 7},
		},
	}))
	for label, attrs := range map[string][]string{
		"Today":         {`data-range-days="1"`},
		"Yesterday":     {`data-range-days="1"`, `data-range-ends-ago="1"`},
		"This week":     {`data-range-unit="week"`},
		"This month":    {`data-range-unit="month"`},
		"Last month":    {`data-range-unit="month"`, `data-range-back="1"`},
		"Last 3 months": {`data-range-months="3"`},
		"This year":     {`data-range-unit="year"`},
		"Last year":     {`data-range-unit="year"`, `data-range-back="1"`},
		"Last 7 days":   {`data-range-days="7"`},
	} {
		tag := presetTag(t, html, label)
		for _, want := range attrs {
			if !strings.Contains(tag, want) {
				t.Errorf("%s: %s lacks %s", label, tag, want)
			}
		}
	}
	// A preset says its range one way: a calendar unit carries no day count.
	if tag := presetTag(t, html, "Last month"); strings.Contains(tag, "data-range-days") || strings.Contains(tag, "data-range-months") {
		t.Errorf("a calendar month also carries a count of days: %s", tag)
	}
	// Today is exactly what it always drew: a count of days and nothing about the end.
	if tag := presetTag(t, html, "Today"); strings.Contains(tag, "ends-ago") || strings.Contains(tag, "range-unit") {
		t.Errorf("a preset that sets only Days drew more than it did: %s", tag)
	}
}

// presetTag is the opening tag of the button that carries label.
func presetTag(t *testing.T, html, label string) string {
	t.Helper()
	end := strings.Index(html, ">"+label+"</button>")
	if end < 0 {
		t.Fatalf("no shortcut %q in:\n%s", label, html)
	}
	start := strings.LastIndex(html[:end], "<button")
	return html[start:end]
}

// The caller's attributes for the preset part stay on the button, and the ones
// that say its range win on a name both write.
func TestTheRangeOfAPresetCannotBeOverwrittenByAPartAttribute(t *testing.T) {
	html := string(components.DateRangePicker(components.DateRangePickerProps{
		Name: "p", Label: "Period",
		ComponentProps: components.ComponentProps{Parts: components.Parts{"preset": {Attrs: components.Attrs{"data-range-days": "99", "data-testid": "shortcut"}}}},
		Presets:        []components.DateRangePreset{{Label: "Yesterday", Days: 1, EndsAgo: 1}},
	}))
	tag := presetTag(t, html, "Yesterday")
	if !strings.Contains(tag, `data-testid="shortcut"`) || !strings.Contains(tag, `data-range-days="1"`) || strings.Contains(tag, `"99"`) {
		t.Fatalf("the shortcut is %s", tag)
	}
}
