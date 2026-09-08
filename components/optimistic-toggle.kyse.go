//go:build kyse

package components

@go
// OptimisticToggleBehavior is the name the client behaviour is registered
// under with arandu.ui.define.
const OptimisticToggleBehavior = "optimistic-toggle"

// OptimisticToggleProps is a like, a follow, a pin: a two-state button that
// flips immediately and reconciles with the server afterwards.
//
// The flip is immediate because the alternative is a control that does nothing
// for as long as the network takes, which people answer by pressing it again.
// The reconciliation is what makes that safe: the endpoint answers with the
// state it actually recorded, and a failure puts the button back where it was
// and says so.
//
// It is a button with aria-pressed and not a checkbox, because it acts rather
// than collects: nothing here is submitted with a form, and the press is the
// whole transaction.
//
// Without script it is a form button that posts and gets the page back with
// the new state in it. The optimism is the enhancement; the state change is
// not.
//
// It publishes root, on, off and count.
type OptimisticToggleProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// URL is what records the change. It answers with this control redrawn in
	// the state it stored, which is what corrects an optimistic flip that was
	// wrong.
	URL string
	// Pressed is the state now.
	Pressed bool
	// Label is what the control is called when it is off: "Like". It is the
	// accessible name in both states, because a name that changes is a control
	// a screen reader can no longer find -- what changes is aria-pressed.
	Label string
	// OnLabel and OffLabel are the words drawn in each state, which may differ
	// from the name: "Following" and "Follow". Empty draws the label in both.
	OnLabel  string
	OffLabel string
	// OnIcon and OffIcon are the pictures for each state. Both are in the
	// markup and one is hidden, so flipping is showing and hiding rather than
	// building an icon.
	OnIcon  template.HTML
	OffIcon template.HTML
	// Count is the number beside it: likes, followers. Empty draws none.
	Count string
	// Variant and Size style it. See ButtonProps.
	Variant string
	Size    string
	// Disabled draws it unavailable.
	Disabled bool
}

// OnText is what is drawn while it is on.
func (p OptimisticToggleProps) OnText() string {
	if p.OnLabel != "" {
		return p.OnLabel
	}
	return p.Label
}

// OffText is what is drawn while it is off.
func (p OptimisticToggleProps) OffText() string {
	if p.OffLabel != "" {
		return p.OffLabel
	}
	return p.Label
}

// State is what aria-pressed says.
func (p OptimisticToggleProps) State() string {
	if p.Pressed {
		return "true"
	}
	return "false"
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p OptimisticToggleProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: OptimisticToggleBehavior}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p OptimisticToggleProps) PartNames() []string {
	return []string{"root", "on", "off", "count"}
}
@endgo

{{-- The answer replaces this control and nothing else, so a failed request
     leaves the page as it was and the behaviour puts the state back. --}}
<button
	data-part="root"
	class="{{ .RootClass("btn") }}"
	type="button"
	aria-pressed="{{ .State() }}"
	aria-label="{{ .Label }}"
	@attributes(.RootAttrs())
	hx-post="{{ .URL }}"
	hx-target="this"
	hx-swap="outerHTML"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
	@if(.Disabled)
		disabled
	@endif
>
	<span
		data-part="on"
		@if(.PartClass("on") != "")
			class="{{ .PartClass("on") }}"
		@endif
		data-toggle-on
		@attributes(.PartAttrs("on"))
		@if(!.Pressed)
			hidden
		@endif
	>
		@if(.OnIcon != "")
			{!! .OnIcon !!}
		@endif
		{{ .OnText() }}
	</span>
	<span
		data-part="off"
		@if(.PartClass("off") != "")
			class="{{ .PartClass("off") }}"
		@endif
		data-toggle-off
		@attributes(.PartAttrs("off"))
		@if(.Pressed)
			hidden
		@endif
	>
		@if(.OffIcon != "")
			{!! .OffIcon !!}
		@endif
		{{ .OffText() }}
	</span>
	@if(.Count != "")
		<span
			data-part="count"
			class="{{ .PartClass("count", "tabular-nums") }}"
			data-toggle-count
			@attributes(.PartAttrs("count"))
		>{{ .Count }}</span>
	@endif
</button>
