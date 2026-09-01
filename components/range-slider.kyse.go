//go:build kyse

package components

import "strconv"

@go
// RangeSliderProps is a number chosen by dragging: one thumb on one track,
// between a floor and a ceiling.
//
// # Why the input is the slider
//
// It is an input of type range, so the thumb, the drag, the arrow keys, Home,
// End, Page Up and Page Down, the role of slider and the three values a screen
// reader reads out all come from the browser. None of that is written here and
// none of it can be got wrong here. What is written here is the label, the
// bounds, and the coloured part of the track.
//
// # What the client is for
//
// Two things, and both are what somebody is looking at while they drag: how
// much of the track is filled, and the number under the thumb. Both are drawn
// by the server first, so a slider is correct before anything runs and stays
// correct if nothing does; while a thumb is moving they are kept in step with
// it, because a number that only catches up on release is a number nobody can
// aim with.
//
// The value that is submitted is the input's own, and the browser sends it
// whether or not any of that happened.
type RangeSliderProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, the id everything else is hung off, and
	// what Page is asked about.
	Name string
	// Label is the text above the track.
	Label string
	// Value is where the thumb starts when nothing was rejected.
	Value int
	// Min is the floor. It defaults to zero, which is also what an unset Min
	// means, so a slider that starts below zero has to say so.
	Min int
	// Max is the ceiling. Zero means a hundred, because a slider whose floor
	// and ceiling are both zero has nowhere to go.
	Max int
	// Step is how far one press of an arrow key moves the thumb. Zero means
	// one.
	Step int
	// Unit is what the number is in, drawn after it: per cent, seconds,
	// pixels. Empty draws the number alone.
	Unit string
	// ShowValue draws the current number under the track. Without it the only
	// way to read a slider is to guess from the thumb, which is fine for
	// loudness and wrong for anything somebody has to report.
	ShowValue bool
	// Hint is the sentence under the track that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this slider asks for
	// its message and for where the thumb was left. Nil draws no message and
	// keeps Value.
	Page Page
	// Disabled draws the slider as unavailable and stops it answering.
	Disabled bool
}

// Floor is the lowest value the thumb reaches.
func (p RangeSliderProps) Floor() int {
	return p.Min
}

// Ceiling is the highest value the thumb reaches.
func (p RangeSliderProps) Ceiling() int {
	if p.Max == 0 {
		return 100
	}
	return p.Max
}

// Tick is how far one press of an arrow key moves the thumb.
func (p RangeSliderProps) Tick() int {
	if p.Step == 0 {
		return 1
	}
	return p.Step
}

// Message is what validation left for this slider, or empty.
func (p RangeSliderProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// submitted is what the form left for this slider: where the thumb was on the
// rejected attempt, and Value when there was none. It is text, because that is
// what a form carries and what Page answers with.
func (p RangeSliderProps) submitted() string {
	if p.Page == nil {
		return strconv.Itoa(p.Value)
	}
	return p.Page.OldOr(p.Name, strconv.Itoa(p.Value))
}

// At is where the thumb is, as a number inside the bounds.
//
// Anything that is not a number, and anything outside the bounds, lands on the
// floor rather than throwing. What comes back from a rejected attempt is
// whatever was posted, and a slider drawn carrying a word is a slider whose
// thumb the browser puts in one place and whose track is coloured to another.
func (p RangeSliderProps) At() int {
	value, err := strconv.Atoi(p.submitted())
	if err != nil || value < p.Floor() || value > p.Ceiling() {
		return p.Floor()
	}
	return value
}

// Current is where the thumb is drawn, which is At written out.
func (p RangeSliderProps) Current() string {
	return strconv.Itoa(p.At())
}

// sliderFills are the twenty-one positions the filled half of the track is
// drawn at, one every five percent.
//
// They are written out rather than built from the number, because a class name
// assembled at run time never appears in the source the stylesheet is compiled
// from, and the rule for it is never emitted.
var sliderFills = [21]string{
	"slider-fill-0", "slider-fill-5", "slider-fill-10", "slider-fill-15", "slider-fill-20",
	"slider-fill-25", "slider-fill-30", "slider-fill-35", "slider-fill-40", "slider-fill-45",
	"slider-fill-50", "slider-fill-55", "slider-fill-60", "slider-fill-65", "slider-fill-70",
	"slider-fill-75", "slider-fill-80", "slider-fill-85", "slider-fill-90", "slider-fill-95",
	"slider-fill-100",
}

// FillClass is how much of the track is coloured, as the class that draws it,
// so the slider is painted right on the first frame rather than on the first
// drag.
//
// A class and not a style attribute: the policy this renders under allows no
// inline style, so a value written into one is dropped by the browser and the
// track paints unfilled. The cost is that it starts in steps of five percent
// until the pointer takes over. What is announced, and what is submitted, is
// the exact value.
func (p RangeSliderProps) FillClass() string {
	span := p.Ceiling() - p.Floor()
	if span <= 0 {
		return sliderFills[0]
	}
	return sliderFills[((p.At()-p.Floor())*100/span+2)/5]
}

// DescribedBy names the element that explains this slider, so the error and
// the hint are announced rather than only shown.
func (p RangeSliderProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
// PartNames are the parts this component publishes.
func (p RangeSliderProps) PartNames() []string {
	return []string{"root", "label", "input", "value", "message", "hint"}
}
@endgo

{{-- data-slider is how the track finds the number that belongs to it, and it is
     the whole of what the client is told. There is nothing to seed: the fill and
     the number below are both written by the server, so the first frame is
     already right. --}}
<div
	data-part="root"
	class="{{ .RootClass("field") }}"
	data-slider
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		for="{{ .Name }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	{{-- One class attribute, not two. It carried "input" and the fill class in
	     separate attributes, and HTML keeps the first occurrence of a repeated
	     attribute and drops the rest -- so the fill was written into the markup
	     and never reached the page. --}}
	<input
		data-part="input"
		class="{{ .PartClass("input", "input", .FillClass()) }}"
		type="range"
		id="{{ .Name }}"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		min="{{ .Floor() }}"
		max="{{ .Ceiling() }}"
		step="{{ .Tick() }}"
		data-slider-track
		@attributes(.PartAttrs("input"))
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Disabled)
			disabled
		@endif
	>

	@if(.ShowValue)
		<p
			data-part="value"
			class="{{ .PartClass("value", "text-muted-foreground text-sm") }}"
			@attributes(.PartAttrs("value"))
		>
			<output for="{{ .Name }}" data-slider-output>{{ .Current() }}</output>
			@if(.Unit != "")
				{{ .Unit }}
			@endif
		</p>
	@endif

	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .Name }}-error"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p
				data-part="hint"
				id="{{ .Name }}-hint"
				class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("hint"))
			>{{ .Hint }}</p>
		@endif
	@endif
</div>
