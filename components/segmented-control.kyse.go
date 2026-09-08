//go:build kyse

package components

@go
// SegmentedControlProps is a small closed set picked from a joined bar: a
// value, not a panel.
//
// The difference from tabs is what is being chosen. Tabs pick which panel is
// shown, and the panel is on the page. This picks a value -- a period, a unit,
// a sort order -- and the page reacts by fetching something. So it is a radio
// group and not a tab list, and a screen reader announces it as a choice
// rather than as navigation.
//
// It is radio inputs with labels over them, so arrow keys, the single tab
// stop, the checked state and the submitted value are the browser's.
//
// It publishes root, label, group, option and segment.
type SegmentedControlProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Page is where a rejected value comes from.
	Page Page
	// Name is the field name, and the id every option hangs off.
	Name string
	// Label names the group. Without it the bar announces as loose radios.
	Label string
	// LabelHidden keeps the name for a screen reader and draws nothing, for a
	// bar whose meaning the surrounding page already gives.
	LabelHidden bool
	// Options are the choices, in the order they read.
	Options []SegmentedOption
	// Value is the one chosen. Empty chooses the first, because a segmented
	// control with nothing chosen is a control in a state it cannot be put
	// back into.
	Value string
	// Size is "sm" or "lg". Empty is the default.
	Size string

	// The HTMX attributes, written only when they carry something. A segmented
	// control usually fetches on change, and HxTrigger defaults to that.
	HxGet     string
	HxPost    string
	HxTarget  string
	HxSwap    string
	HxTrigger string
}

// SegmentedOption is one segment of the bar.
type SegmentedOption struct {
	// Label is the text on it.
	Label string
	// Value is what choosing it submits.
	Value string
	// Disabled draws it unavailable, and the arrow keys pass over it.
	Disabled bool
}

// Current is the chosen value: what came back rejected, what the caller set,
// or the first option.
func (p SegmentedControlProps) Current() string {
	chosen := p.Value
	if p.Page != nil {
		chosen = p.Page.OldOr(p.Name, p.Value)
	}
	if chosen == "" && len(p.Options) > 0 {
		return p.Options[0].Value
	}
	return chosen
}

// Chosen is whether one option is the chosen one.
func (p SegmentedControlProps) Chosen(option SegmentedOption) bool {
	return option.Value == p.Current()
}

// OptionID is the id of one option, which its label points at.
func (p SegmentedControlProps) OptionID(option SegmentedOption) string {
	return p.Name + "-" + option.Value
}

// Trigger is when the request is made. A segmented control that fetched on
// anything but change would fetch on the arrow key that is passing over an
// option on the way to another.
func (p SegmentedControlProps) Trigger() string {
	if p.HxTrigger != "" {
		return p.HxTrigger
	}
	return "change"
}

// Fetches is whether this bar asks the server for anything at all.
func (p SegmentedControlProps) Fetches() bool {
	return p.HxGet != "" || p.HxPost != ""
}

// PartNames are the parts this component publishes.
func (p SegmentedControlProps) PartNames() []string {
	return []string{"root", "label", "group", "option", "segment"}
}
@endgo

<fieldset
	data-part="root"
	class="{{ .RootClass("field") }}"
	@attributes(.RootAttrs())
>
	<legend
		data-part="label"
		@if(.LabelHidden)
			class="{{ .PartClass("label", "sr-only") }}"
		@endif
		@if(!.LabelHidden)
			class="{{ .PartClass("label", "label") }}"
		@endif
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</legend>

	<div
		data-part="group"
		class="{{ .PartClass("group", "segmented-control") }}"
		@attributes(.PartAttrs("group"))
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
		@if(.HxGet != "")
			hx-get="{{ .HxGet }}"
		@endif
		@if(.HxPost != "")
			hx-post="{{ .HxPost }}"
		@endif
		@if(.Fetches())
			hx-trigger="{{ .Trigger() }}"
		@endif
		@if(.HxTarget != "")
			hx-target="{{ .HxTarget }}"
		@endif
		@if(.HxSwap != "")
			hx-swap="{{ .HxSwap }}"
		@endif
	>
		@foreach(.Options as option)
			<input
				data-part="option"
				class="{{ .PartClass("option", "sr-only") }}"
				type="radio"
				id="{{ .OptionID(option) }}"
				name="{{ .Name }}"
				value="{{ option.Value }}"
				@attributes(.PartAttrs("option"))
				@if(.Chosen(option))
					checked
				@endif
				@if(option.Disabled)
					disabled
				@endif
			>
			<label
				data-part="segment"
				@if(.PartClass("segment") != "")
					class="{{ .PartClass("segment") }}"
				@endif
				for="{{ .OptionID(option) }}"
				@attributes(.PartAttrs("segment"))
			>{{ option.Label }}</label>
		@endforeach
	</div>
</fieldset>
