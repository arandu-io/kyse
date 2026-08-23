//go:build kyse

package components

@go
// ButtonGroupProps is several buttons welded into one control: the inner
// corners squared off, the borders shared, the whole thing announced as one
// group with one name.
//
// # Why it takes buttons rather than markup
//
// Buttons are ButtonProps, drawn by the same function a button anywhere else
// is drawn by. A group that took markup would be a second place a button is
// written, and the two would drift the first time a variant is added.
//
// # What the name is for
//
// Label is the group's own, and it is not decoration. Three icon buttons in a
// row are announced one after another with nothing saying what they belong to,
// and "Archive, Report, Snooze" is a list of three unrelated commands until
// something says they are what can be done with this message.
type ButtonGroupProps struct {
	// Label is what assistive technology calls the group. Empty leaves the
	// group unnamed, which is right only when the buttons are words that
	// already say what they are together.
	Label string
	// Buttons are the members, left to right, or top to bottom when the group
	// is vertical. Fewer than two is a button, and draws as one inside a group
	// that adds nothing.
	Buttons []ButtonProps
	// Vertical stacks them instead. The stylesheet squares off the shared
	// edges either way; this is which pair of edges.
	Vertical bool
	// Separated draws a rule between the members. It is for the group whose
	// members do different things -- an action and the menu of its variants --
	// where a shared border reads as one wide button.
	Separated bool
}

// Orientation is what the stylesheet is told, and it is written only when the
// group is vertical: the horizontal rules are the ones that apply when the
// attribute is absent.
func (p ButtonGroupProps) Orientation() string {
	if p.Vertical {
		return "vertical"
	}
	return ""
}
@endgo

<div class="button-group" role="group"
	@if(.Label != "")
		aria-label="{{ .Label }}"
	@endif
	@if(.Orientation() != "")
		data-orientation="{{ .Orientation() }}"
	@endif
>
	@for(at := 0; at < len(.Buttons); at++)
		@if(.Separated && at > 0)
			<hr role="separator">
		@endif
		{!! Button(.Buttons[at]) !!}
	@endfor
</div>
