//go:build kyse

package components

import (
	"math"
	"strconv"
)

@go
// MeterProps is a measurement inside a range that is known: disk used of disk
// there, seats taken of seats sold, score out of the score possible.
//
// It is not a progress bar, and the difference decides which one to draw. A
// progress bar is a task moving towards being finished, and it ends. A meter
// is a level that goes up and down and has no end -- disk usage at sixty per
// cent is not sixty per cent finished.
//
// The three thresholds are what makes it a meter rather than a bar: low, high
// and optimum tell the browser which end of the range is the good end, so the
// same seventy per cent draws as healthy for a score and as a warning for a
// disk.
//
// It publishes root.
type MeterProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Value is the measurement.
	Value float64
	// Min is the floor of the range. Zero is the floor and is also the default.
	Min float64
	// Max is the ceiling. Zero means one, which is what an unset ceiling means
	// to the element itself, so a fraction can be passed with nothing else set.
	Max float64
	// Low and High are where the range stops being low and starts being high.
	// Both zero leaves the range undivided.
	Low  float64
	High float64
	// Optimum is where in the range the good value is. It is what decides
	// whether the segment under the value draws as good or as a warning: a
	// disk sets it near Min, a battery near Max.
	Optimum float64
	// Text is what the element reads for a browser that does not draw meters,
	// and it is what a person is told the number means. Empty writes the value
	// and the ceiling with a slash between them, which is right for a count
	// and wrong for anything carrying a unit.
	Text string
	// Title is the unit, which is the one place the element itself has for it:
	// "gigabytes", "seats".
	Title string
	// Label names the measurement. Without it the meter is a bar with no
	// subject.
	Label string
}

// Ceiling is the top of the range, and is one when nothing set it, because a
// meter whose ceiling is zero has no range at all.
func (p MeterProps) Ceiling() float64 {
	if p.Max == 0 {
		return 1
	}
	return p.Max
}

// Reading is what the element says: what the caller wrote, or the value over
// the ceiling.
func (p MeterProps) Reading() string {
	if p.Text != "" {
		return p.Text
	}
	return formatNumber(p.Value) + " / " + formatNumber(p.Ceiling())
}

// formatNumber writes a float the way a person would: without the decimal
// point when there is nothing after it.
func formatNumber(v float64) string {
	if v == math.Trunc(v) && math.Abs(v) < 1e15 {
		return strconv.FormatInt(int64(v), 10)
	}
	return strconv.FormatFloat(v, 'f', -1, 64)
}

// PartNames are the parts this component publishes.
func (p MeterProps) PartNames() []string { return []string{"root"} }
@endgo

<meter
	data-part="root"
	class="{{ .RootClass("meter") }}"
	value="{{ .Value }}"
	min="{{ .Min }}"
	max="{{ .Ceiling() }}"
	@if(.Low != 0)
		low="{{ .Low }}"
	@endif
	@if(.High != 0)
		high="{{ .High }}"
	@endif
	@if(.Optimum != 0)
		optimum="{{ .Optimum }}"
	@endif
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@if(.Title != "")
		title="{{ .Title }}"
	@endif
	@attributes(.RootAttrs())
>{{ .Reading() }}</meter>
