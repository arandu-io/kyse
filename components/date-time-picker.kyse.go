//go:build kyse

package components

import (
	"time"

	"github.com/arandu-io/kyse/icons"
)

@go
// DateTimePickerProps is a date, a time, or both, in the field the browser
// already draws a calendar for.
//
// There is no calendar written here, and that is the decision. A date field is
// where a locale is least negotiable: the order of day and month, the first
// day of the week, the names of the months, the calendar system itself. The
// browser knows the reader's; a script does not, and one that guesses is wrong
// in exactly the places where being wrong is expensive.
//
// The value crossing the wire is never what is drawn. A date field submits
// "2026-09-07" whatever the reader sees, and a datetime-local field submits
// "2026-09-07T14:30" with no zone at all -- so a moment that has to be exact
// is either stored with the zone alongside it or asked for in UTC.
//
// It publishes root, label, input, message and hint.
type DateTimePickerProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id the label points at.
	Name string
	// Label is the text above it.
	Label string
	// Kind is what is being asked for: "date", "time", "datetime", "month" or
	// "week". Empty asks for a date.
	Kind string
	// Value is what is shown, in the form the chosen kind submits:
	// "2026-09-07" for a date, "14:30" for a time, "2026-09-07T14:30" for
	// both, "2026-09" for a month, "2026-W37" for a week. A value in any other
	// shape is ignored by the browser and the field comes back empty.
	Value string
	// Min and Max are the range, in the same form as Value. A booking that
	// cannot be in the past sets Min and is then refused by the browser as
	// well as by the server.
	Min string
	Max string
	// Step is the granularity, in seconds for a time and in days for a date.
	// "60" is whole minutes, which is the default for a time field; "900" is
	// quarter hours. Empty leaves the browser's.
	Step string
	// Hint is the sentence under it while nothing is wrong. It is where the
	// zone belongs, because the field carries none: "Times are UTC".
	Hint string
	// Required marks the field.
	Required bool
	// Disabled draws it unavailable.
	Disabled bool
	// ReadOnly shows the value without letting it change. It still submits,
	// which is what separates it from Disabled.
	ReadOnly bool

	// Calendar draws a month grid beside the box instead of leaving the
	// choosing to the browser's own picker.
	//
	// The native popup is chrome: no rule reaches inside it, in any engine,
	// so a project with a design cannot have the picker look like the rest of
	// its screens. This draws one that can.
	//
	// What does not change is where the value lives. The input is still the
	// field, still submits, and still enforces Min, Max, Step and Required --
	// the calendar writes into it. So the native picker remains the whole
	// answer before any script has run, which is why this is drawn beside the
	// box rather than in place of it.
	//
	// It has no effect on a time, a month or a week: there is no grid of
	// those, and the browser's own control for them is not the one that is
	// ugly.
	Calendar bool
	// MonthNames, WeekdayNames and FirstDay are handed to that grid. See
	// CalendarProps, where the reason they are data rather than a guess is
	// written out.
	MonthNames   []string
	WeekdayNames []string
	FirstDay     time.Weekday
	// OpenLabel names the control that opens the grid. Empty says "Choose a
	// date", which a screen reader needs because the control is an icon.
	OpenLabel string
	// TodayLabel and ClearLabel are the two controls under the grid.
	TodayLabel string
	ClearLabel string
}

// InputType is the element type for the kind asked for. An unknown kind is a
// date, because a field that fell back to text would accept anything and
// submit it.
func (p DateTimePickerProps) InputType() string {
	switch p.Kind {
	case "time":
		return "time"
	case "datetime":
		return "datetime-local"
	case "month":
		return "month"
	case "week":
		return "week"
	default:
		return "date"
	}
}

// Current is the value shown: what came back rejected, or what the caller set.
func (p DateTimePickerProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// Message is the rejection for this field, and empty when there is none.
func (p DateTimePickerProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the field.
func (p DateTimePickerProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// Griddable is whether a month grid means anything for what is being asked
// for. A time, a month and a week have no grid of days to draw.
func (p DateTimePickerProps) Griddable() bool {
	return p.Calendar && (p.Kind == "" || p.Kind == "date" || p.Kind == "datetime")
}

// Day is the date half of the value, which is what the grid selects. A
// datetime carries a time after it, and the grid neither reads nor writes
// that -- it replaces the ten characters in front.
func (p DateTimePickerProps) Day() string {
	current := p.Current()
	if len(current) >= 10 {
		return current[:10]
	}
	return current
}

// DayMin and DayMax are the bounds as the grid reads them, which is the date
// half of whatever the field was given.
func (p DateTimePickerProps) DayMin() string { return firstTen(p.Min) }

// DayMax is the ceiling as the grid reads it.
func (p DateTimePickerProps) DayMax() string { return firstTen(p.Max) }

// firstTen is the date half of a value that may carry a time.
func firstTen(value string) string {
	if len(value) >= 10 {
		return value[:10]
	}
	return value
}

// OpenName is what the control that opens the grid is called.
func (p DateTimePickerProps) OpenName() string {
	if p.OpenLabel != "" {
		return p.OpenLabel
	}
	return "Choose a date"
}

// CalendarID and PanelID are the ids the grid and the panel hang off.
func (p DateTimePickerProps) CalendarID() string { return p.Name + "-calendar" }

// PanelID is the id of the element that opens.
func (p DateTimePickerProps) PanelID() string { return p.Name + "-calendar-panel" }

// Grid is the calendar this picker draws, built from the picker's own fields
// so the two cannot disagree about the day, the bounds or the language.
func (p DateTimePickerProps) Grid() CalendarProps {
	return CalendarProps{
		// The caller's parts go through, so a name published here reaches the
		// element it names whether that element is drawn by this component or
		// by the one it composes. A composition that swallowed them would
		// publish names nothing could reach.
		ComponentProps: ComponentProps{Parts: p.Parts},
		ID:           p.CalendarID(),
		Label:        p.OpenName(),
		Value:        p.Day(),
		Min:          p.DayMin(),
		Max:          p.DayMax(),
		MonthNames:   p.MonthNames,
		WeekdayNames: p.WeekdayNames,
		FirstDay:     p.FirstDay,
		TodayLabel:   p.TodayLabel,
		ClearLabel:   p.ClearLabel,
		Target:       p.Name,
	}
}

// PartNames are the parts this component publishes.
func (p DateTimePickerProps) PartNames() []string {
	// The calendar's names are published here as well, because this draws one
	// and the parts a caller can reach are the parts on the page -- not the
	// parts this file happens to write itself.
	return append(
		[]string{"root", "label", "group", "input", "open", "panel", "message", "hint"},
		CalendarProps{}.PartNames()[1:]...,
	)
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("field") }}"
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		for="{{ .Name }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	{{-- The group is only drawn when there is a grid to sit beside the box.
	     A field with one child in a flex row is a wrapper that exists for
	     nothing, and it would change the markup of every picker that does not
	     ask for a calendar. --}}
	@if(.Griddable())
		<div
			data-part="group"
			class="{{ .PartClass("group", "date-picker") }}"
			@attributes(.PartAttrs("group"))
		>
	@endif

	<input
		data-part="input"
		class="{{ .PartClass("input", "input") }}"
		type="{{ .InputType() }}"
		id="{{ .Name }}"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		@attributes(.PartAttrs("input"))
		@if(.Min != "")
			min="{{ .Min }}"
		@endif
		@if(.Max != "")
			max="{{ .Max }}"
		@endif
		@if(.Step != "")
			step="{{ .Step }}"
		@endif
		@if(.DescribedBy() != "")
			aria-describedby="{{ .DescribedBy() }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			required
		@endif
		@if(.Disabled)
			disabled
		@endif
		@if(.ReadOnly)
			readonly
		@endif
	>

	@if(.Griddable())
			{{-- The panel is a popover the shipped behaviour already opens and
			     light-dismisses, so nothing here decides when it is on screen.
			     The grid inside it is the Calendar component, built from this
			     picker's own fields -- the two cannot disagree about the day,
			     the bounds or the language, because there is one source for
			     each. --}}
			<button
				data-part="open"
				class="{{ .PartClass("open", "btn") }}"
				type="button"
				data-variant="ghost"
				data-size="sm"
				data-align="end"
				id="{{ .PanelID() }}-trigger"
				aria-controls="{{ .PanelID() }}"
				aria-expanded="false"
				aria-label="{{ .OpenName() }}"
				@if(.Disabled)
					disabled
				@endif
				@attributes(.PartAttrs("open"))
			>{!! icons.CalendarBlank(icons.Props{}) !!}</button>

			<div
				data-part="panel"
				class="{{ .PartClass("panel", "date-picker-panel") }}"
				id="{{ .PanelID() }}"
				data-popover
				aria-hidden="true"
				aria-labelledby="{{ .PanelID() }}-trigger"
				@attributes(.PartAttrs("panel"))
			>{!! Calendar(.Grid()) !!}</div>
		</div>
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
