//go:build kyse

package components

@go
// ProgressProps is how far along something is, drawn as a bar.
//
// It reports and does not drive: the bar is redrawn when the fragment holding
// it is replaced, and a bar that animated itself would be a bar telling a story
// about work nothing measured.
//
// Two states, and both are the same component. A value that is known fills the
// bar to it and says so with aria-valuenow. A value that is not known fills the
// bar and pulses, and writes no aria-valuenow at all -- which is precisely how
// an unknown share is spelled, and not an omission.
// It publishes root and fill, fill being the part of the bar that is filled.
type ProgressProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label is what is progressing: "Upload", "Import", "Restore". It becomes
	// the accessible name, and a bar without one is announced as a percentage
	// of nothing in particular.
	Label string
	// Value is how far along it is, in the same unit as Max. A value outside
	// the range is held at the nearest end rather than drawn outside the bar.
	Value int
	// Max is the value that means finished. Zero and below mean 100, so Value
	// alone reads as a percentage.
	Max int
	// Indeterminate says the share done is not known -- work has started and
	// nothing counts the steps. It draws a full pulsing bar and no value.
	Indeterminate bool
}

// Ceiling is the value that means finished.
func (p ProgressProps) Ceiling() int {
	if p.Max <= 0 {
		return 100
	}
	return p.Max
}

// Now is Value held inside the bar's own range.
func (p ProgressProps) Now() int {
	if p.Value < 0 {
		return 0
	}
	if ceiling := p.Ceiling(); p.Value > ceiling {
		return ceiling
	}
	return p.Value
}

// Percent is how much of the bar is filled, from 0 to 100.
func (p ProgressProps) Percent() int {
	return p.Now() * 100 / p.Ceiling()
}

// progressWidths are the twenty-one widths the bar is drawn at, one every five
// percent.
//
// They are written out rather than built from the number, because a class name
// assembled at run time never appears in the source the stylesheet is compiled
// from, and the rule for it is never emitted.
var progressWidths = [21]string{
	"w-0", "w-[5%]", "w-[10%]", "w-[15%]", "w-[20%]",
	"w-[25%]", "w-[30%]", "w-[35%]", "w-[40%]", "w-[45%]",
	"w-[50%]", "w-[55%]", "w-[60%]", "w-[65%]", "w-[70%]",
	"w-[75%]", "w-[80%]", "w-[85%]", "w-[90%]", "w-[95%]",
	"w-full",
}

// WidthClass is how far the bar is filled, as the class that draws it.
//
// A class and not a style attribute: the policy this renders under allows no
// inline style, so a width written into one is dropped by the browser and the
// bar never moves. The cost is that it fills in steps of five percent, which is
// one pixel on a bar twenty pixels wide. What is announced is the exact value.
func (p ProgressProps) WidthClass() string {
	return progressWidths[(p.Percent()+2)/5]
}

// PartNames are the parts this component publishes.
func (p ProgressProps) PartNames() []string { return []string{"root", "fill"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("progress") }}"
	role="progressbar"
	@attributes(.RootAttrs())
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	aria-valuemin="0"
	aria-valuemax="{{ .Ceiling() }}"
	@if(!.Indeterminate)
		aria-valuenow="{{ .Now() }}"
	@endif
>
	@if(.Indeterminate)
		<span
			data-part="fill"
			class="{{ .PartClass("fill", "w-full animate-pulse") }}"
			@attributes(.PartAttrs("fill"))
		></span>
	@else
		<span
			data-part="fill"
			class="{{ .PartClass("fill", .WidthClass()) }}"
			@attributes(.PartAttrs("fill"))
		></span>
	@endif
</div>
