//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/kyse/icons"
)

@go
// RatingProps is a score out of five, as the radio group it actually is.
//
// It is radio inputs and labels, and nothing else. That is not a limitation
// worked around -- it is what makes the arrow keys, the tab stop, the required
// check and the submitted value all belong to the browser. A rating written as
// a row of buttons has none of them and has to grow each one back.
//
// The hover preview is a sibling selector in the stylesheet, so a pointer sees
// the stars fill ahead of the one under it without a line of script. Reading
// the score back is the value of the checked radio, which is what a form
// sends.
//
// It publishes root, label, group, option, star, message and hint.
type RatingProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value and its message come from.
	Page Page
	// Name is the field name, and the id every option hangs off.
	Name string
	// Label names the group. It is a fieldset legend and not a floating line,
	// because five radios without one announce as five unrelated choices.
	Label string
	// Value is the score chosen, as the number written: "4". Empty is
	// unrated, which is different from zero.
	Value string
	// Count is how many stars there are. Zero means five.
	Count int
	// OptionLabels are what each star is called, from the first up: {"Terrible",
	// "Poor", "Fair", "Good", "Excellent"}. Fewer than Count leaves the rest
	// named by their number, which is worse and is still a name.
	OptionLabels []string
	// Hint is the sentence under it while nothing is wrong.
	Hint string
	// Required refuses to submit unrated.
	Required bool
	// ReadOnly draws the score without letting it change, for a rating being
	// shown rather than given.
	ReadOnly bool
}

// Stars is how many there are.
func (p RatingProps) Stars() int {
	if p.Count <= 0 {
		return 5
	}
	return p.Count
}

// Current is the chosen score: what came back rejected, or what the caller
// set.
func (p RatingProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// Chosen is whether the star at position n is the one chosen.
func (p RatingProps) Chosen(n int) bool {
	return p.Current() == strconv.Itoa(n)
}

// OptionName is what the star at position n is called: the word the caller
// gave, or the number it is.
func (p RatingProps) OptionName(n int) string {
	if n-1 < len(p.OptionLabels) && p.OptionLabels[n-1] != "" {
		return p.OptionLabels[n-1]
	}
	if n == 1 {
		return "1 star"
	}
	return strconv.Itoa(n) + " stars"
}

// OptionID is the id of one option, which the label points at.
func (p RatingProps) OptionID(n int) string {
	return p.Name + "-" + strconv.Itoa(n)
}

// Message is the rejection for this field, and empty when there is none.
func (p RatingProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy is the id of whatever is explaining the group.
func (p RatingProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p RatingProps) PartNames() []string {
	return []string{"root", "label", "group", "option", "star", "message", "hint"}
}
@endgo

<fieldset
	data-part="root"
	class="{{ .RootClass("field rating") }}"
	@attributes(.RootAttrs())
	@if(.DescribedBy() != "")
		aria-describedby="{{ .DescribedBy() }}"
	@endif
>
	<legend
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</legend>

	<div
		data-part="group"
		class="{{ .PartClass("group", "rating-stars") }}"
		@attributes(.PartAttrs("group"))
	>
		@for(star := 1; star <= .Stars(); star++)
			<input
				data-part="option"
				class="{{ .PartClass("option", "sr-only") }}"
				type="radio"
				id="{{ .OptionID(star) }}"
				name="{{ .Name }}"
				value="{{ star }}"
				@attributes(.PartAttrs("option"))
				@if(.Chosen(star))
					checked
				@endif
				@if(.Required)
					required
				@endif
				@if(.ReadOnly)
					disabled
				@endif
			>
			<label
				data-part="star"
				@if(.PartClass("star") != "")
					class="{{ .PartClass("star") }}"
				@endif
				for="{{ .OptionID(star) }}"
				@attributes(.PartAttrs("star"))
			><span class="sr-only">{{ .OptionName(star) }}</span>{!! icons.Star(icons.Props{}) !!}</label>
		@endfor
	</div>

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
</fieldset>
