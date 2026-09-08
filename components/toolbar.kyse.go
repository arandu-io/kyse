//go:build kyse

package components

@go
// ToolbarBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const ToolbarBehavior = "toolbar"

// ToolbarProps is a row of controls that is one tab stop.
//
// That is the whole point of it, and it is what separates a toolbar from a row
// of buttons. Twelve buttons in a row are twelve tab stops, and a page with
// three such rows costs somebody using a keyboard thirty-six presses to get
// past. A toolbar is one stop, and the arrow keys move inside it -- which is
// the same bargain a radio group, a menu and a tab list all make.
//
// The controls are values, not markup: a toolbar that could hold anything
// would have to take a string of HTML, which is where escaping stops being
// guaranteed by construction.
//
// It publishes root, group, control, toggle and separator.
type ToolbarProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label names the toolbar. Without it the row announces as a toolbar with
	// no subject, and a page with two of them has two of those.
	Label string
	// Items are the controls, separators and groups, in the order they read.
	Items []ToolbarItem
	// Orientation is "vertical" for a column. Empty is a row, and the arrow
	// keys follow whichever it is -- left and right along a row, up and down
	// down a column.
	Orientation string
}

// ToolbarItem is one control in the row.
//
// Separator and Pressed are what a control is instead of a plain button, and
// they are read in that order: a separator carries nothing, a control with
// Toggle set is a two-state button, and everything else acts once.
type ToolbarItem struct {
	// Label is the text on it.
	Label string
	// Icon is drawn before the label. A control with an icon and no label
	// takes its name from the label anyway, as an aria-label.
	Icon template.HTML
	// URL makes it a link. Empty makes it a button.
	URL string
	// Toggle draws it as a two-state control, with Pressed saying which state.
	// It is the bold button, the pinned filter, the muted channel.
	Toggle bool
	// Pressed is the state of a toggle, and means nothing without one.
	Pressed bool
	// Variant is "destructive" for the control that removes something.
	Variant string
	// Disabled draws it unavailable, and the arrow keys pass over it.
	Disabled bool
	// Separator draws a rule between groups instead of a control.
	Separator bool
	// IconOnly draws the label as the accessible name rather than as text.
	IconOnly bool

	// The HTMX attributes, written only when they carry something. An empty
	// hx-post is an attribute HTMX acts on: it would post to the current URL.
	HxPost   string
	HxGet    string
	HxTarget string
	HxSwap   string
}

// Vertical is whether the arrow keys run up and down rather than left and
// right.
func (p ToolbarProps) Vertical() bool { return p.Orientation == "vertical" }

// Direction is what the element declares its orientation to be, which is what
// tells assistive technology which arrow keys to promise.
func (p ToolbarProps) Direction() string {
	if p.Vertical() {
		return "vertical"
	}
	return "horizontal"
}

// First is the index of the control that holds the tab stop: the first one
// that can take it. A toolbar whose first control is disabled would otherwise
// have a tab stop nobody can land on.
func (p ToolbarProps) First() int {
	for at, item := range p.Items {
		if !item.Separator && !item.Disabled {
			return at
		}
	}
	return -1
}

// Stop is whether the control at this index is the one holding the tab stop.
func (p ToolbarProps) Stop(at int) bool { return at == p.First() }

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p ToolbarProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: ToolbarBehavior, Props: map[string]any{"vertical": p.Vertical()}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p ToolbarProps) PartNames() []string {
	return []string{"root", "control", "separator"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("toolbar") }}"
	role="toolbar"
	aria-label="{{ .Label }}"
	aria-orientation="{{ .Direction() }}"
	data-orientation="{{ .Direction() }}"
	@attributes(.RootAttrs())
>
	@for(at := 0; at < len(.Items); at++)
		@if(.Items[at].Separator)
			<hr
				data-part="separator"
				@if(.PartClass("separator") != "")
					class="{{ .PartClass("separator") }}"
				@endif
				role="separator"
				@attributes(.PartAttrs("separator"))
			>
		@elseif(.Items[at].URL != "")
			<a
				data-part="control"
				class="{{ .PartClass("control", "btn") }}"
				href="{{ .Items[at].URL }}"
				data-variant="ghost"
				data-size="sm"
				@attributes(.PartAttrs("control"))
				@if(.Stop(at))
					tabindex="0"
				@endif
				@if(!.Stop(at))
					tabindex="-1"
				@endif
				@if(.Items[at].IconOnly)
					aria-label="{{ .Items[at].Label }}"
				@endif
			>
				@if(.Items[at].Icon != "")
					{!! .Items[at].Icon !!}
				@endif
				@if(!.Items[at].IconOnly)
					{{ .Items[at].Label }}
				@endif
			</a>
		@else
			<button
				data-part="control"
				class="{{ .PartClass("control", "btn") }}"
				type="button"
				data-size="sm"
				@attributes(.PartAttrs("control"))
				@if(.Items[at].Variant != "")
					data-variant="{{ .Items[at].Variant }}"
				@endif
				@if(.Items[at].Variant == "")
					data-variant="ghost"
				@endif
				@if(.Stop(at))
					tabindex="0"
				@endif
				@if(!.Stop(at))
					tabindex="-1"
				@endif
				@if(.Items[at].Toggle && .Items[at].Pressed)
					aria-pressed="true"
				@endif
				@if(.Items[at].Toggle && !.Items[at].Pressed)
					aria-pressed="false"
				@endif
				@if(.Items[at].Disabled)
					disabled
				@endif
				@if(.Items[at].IconOnly)
					aria-label="{{ .Items[at].Label }}"
				@endif
				@if(.Items[at].HxPost != "")
					hx-post="{{ .Items[at].HxPost }}"
				@endif
				@if(.Items[at].HxGet != "")
					hx-get="{{ .Items[at].HxGet }}"
				@endif
				@if(.Items[at].HxTarget != "")
					hx-target="{{ .Items[at].HxTarget }}"
				@endif
				@if(.Items[at].HxSwap != "")
					hx-swap="{{ .Items[at].HxSwap }}"
				@endif
			>
				@if(.Items[at].Icon != "")
					{!! .Items[at].Icon !!}
				@endif
				@if(!.Items[at].IconOnly)
					{{ .Items[at].Label }}
				@endif
			</button>
		@endif
	@endfor
</div>
