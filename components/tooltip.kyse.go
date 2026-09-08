//go:build kyse

package components

@go
// TooltipProps is a short label that appears beside a control on hover and on
// focus.
//
// It draws its own trigger, and that is not a limitation worked around. A
// tooltip is bound to one control -- the tip has to name that control, the
// control has to point at the tip, and the pair has to be a single hover
// target or the tip disappears while the pointer travels to it. Written as a
// wrapper around whatever somebody passes, none of the three is guaranteed;
// written as the pair, all three are.
//
// On focus as well as on hover, always. A tooltip that appears only under a
// pointer is a tooltip nobody using a keyboard ever sees, and it is usually
// the only place a name for an icon-only button was written down.
//
// What it is not: it is not a place for a sentence, a link or anything to
// click. The tip is pointer-events-none, so nothing in it can be reached, and
// a tip long enough to want to be read is a description that belongs under the
// control.
//
// It publishes root, trigger and content.
type TooltipProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the tip and the trigger use to point at each other. Two
	// tooltips on one page need two.
	ID string
	// Text is the tip.
	Text string
	// Label is the text on the trigger. Empty draws only the icon, and then
	// the tip is the trigger's accessible name as well -- which is the
	// icon-only button case and the reason most tooltips exist.
	Label string
	// Icon is drawn on the trigger before the label.
	Icon template.HTML
	// URL makes the trigger a link. Empty makes it a button.
	URL string
	// Side is where the tip sits: "top", "bottom", "left" or "right". Empty is
	// above.
	Side string
	// Variant and Size style the trigger. See ButtonProps.
	Variant string
	Size    string
	// Disabled draws the trigger unavailable. The tip still appears, because
	// the reason a control is unavailable is exactly what somebody hovering it
	// is asking.
	Disabled bool
}

// ContentID is the id of the tip, which the trigger points at.
func (p TooltipProps) ContentID() string { return p.ID + "-tip" }

// Named is whether the trigger has visible text. Without it the tip is the
// only name the control has.
func (p TooltipProps) Named() bool { return p.Label != "" }

// PartNames are the parts this component publishes.
func (p TooltipProps) PartNames() []string {
	return []string{"root", "trigger", "content"}
}
@endgo

{{-- The tip is described-by and not labelled-by when the trigger has text of
     its own: a label would replace the name, and the name is the text. When
     the trigger has no text there is nothing to replace, so the tip becomes
     the name -- which is the only way an icon-only control gets one from
     here. --}}
<span
	data-part="root"
	class="{{ .RootClass("tooltip") }}"
	@if(.Side != "")
		data-side="{{ .Side }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(.URL != "")
		<a
			data-part="trigger"
			class="{{ .PartClass("trigger", "btn") }}"
			href="{{ .URL }}"
			@if(.Named())
				aria-describedby="{{ .ContentID() }}"
			@endif
			@if(!.Named())
				aria-labelledby="{{ .ContentID() }}"
			@endif
			@if(.Variant != "")
				data-variant="{{ .Variant }}"
			@endif
			@if(.Size != "")
				data-size="{{ .Size }}"
			@endif
			@attributes(.PartAttrs("trigger"))
		>@if(.Icon != ""){!! .Icon !!}@endif{{ .Label }}</a>
	@endif
	@if(.URL == "")
		<button
			data-part="trigger"
			class="{{ .PartClass("trigger", "btn") }}"
			type="button"
			@if(.Named())
				aria-describedby="{{ .ContentID() }}"
			@endif
			@if(!.Named())
				aria-labelledby="{{ .ContentID() }}"
			@endif
			@if(.Disabled)
				disabled
			@endif
			@if(.Variant != "")
				data-variant="{{ .Variant }}"
			@endif
			@if(.Size != "")
				data-size="{{ .Size }}"
			@endif
			@attributes(.PartAttrs("trigger"))
		>@if(.Icon != ""){!! .Icon !!}@endif{{ .Label }}</button>
	@endif

	<span
		data-part="content"
		class="{{ .PartClass("content", "tooltip-content") }}"
		id="{{ .ContentID() }}"
		role="tooltip"
		@if(.Side != "")
			data-side="{{ .Side }}"
		@endif
		@attributes(.PartAttrs("content"))
	>{{ .Text }}</span>
</span>
