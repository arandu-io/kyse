//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

@go
// SplitButtonProps is the action people take, welded to the ones they take
// less often.
//
// It exists because a menu of five equal entries makes the common one cost the
// same as the rare ones: open, read, aim, click. Splitting it means the common
// one is a click, and the other four are still there.
//
// The two halves are two buttons and not one, because they do different things
// and a screen reader has to be able to say so. One button with a menu on part
// of its surface is a control nobody can operate without a pointer.
//
// It publishes root, action, trigger, panel, menu, item and shortcut.
type SplitButtonProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the trigger, the panel and the menu hang their ids off. Two
	// split buttons on one page need two.
	ID string
	// Label is the text on the acting half.
	Label string
	// URL makes that half a link. Empty makes it a button.
	URL string
	// Type is the HTML type of the acting half when it is a button:
	// "button", "submit" or "reset". Empty submits nothing.
	Type string
	// Variant and Size style both halves together, because two halves of one
	// control styled apart is two controls.
	Variant string
	Size    string
	// Disabled draws both halves as unavailable.
	Disabled bool
	// MenuLabel is what the second half is called to a screen reader: "More
	// send options". Without it the control announces as an unnamed button,
	// and the caret is not a name.
	MenuLabel string
	// Items are the entries of the menu, in the order they read.
	Items []MenuItem

	// The HTMX attributes of the acting half, written only when they carry
	// something. An empty hx-post is an attribute HTMX acts on: it would post
	// to the current URL.
	HxPost   string
	HxGet    string
	HxTarget string
	HxSwap   string
}

// ButtonType is the type attribute of the acting half, defaulting to a button
// that does nothing on its own.
func (p SplitButtonProps) ButtonType() string {
	if p.Type == "" {
		return "button"
	}
	return p.Type
}

// TriggerID, PanelID and MenuID are the ids the three elements need to point at
// one another.
func (p SplitButtonProps) TriggerID() string { return p.ID + "-trigger" }

// PanelID is the id of the element that opens.
func (p SplitButtonProps) PanelID() string { return p.ID + "-panel" }

// MenuID is the id of the menu inside it.
func (p SplitButtonProps) MenuID() string { return p.ID + "-menu" }

// TriggerName is what the second half is called, and it is never empty: a
// control with no name is a control a screen reader reads as "button".
func (p SplitButtonProps) TriggerName() string {
	if p.MenuLabel != "" {
		return p.MenuLabel
	}
	return "More options"
}

// PartNames are the parts this component publishes.
func (p SplitButtonProps) PartNames() []string {
	return []string{"root", "action", "trigger", "panel", "menu", "item", "shortcut"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("split-button") }}"
	id="{{ .ID }}"
	@attributes(.RootAttrs())
>
	@if(.URL != "")
		<a
			data-part="action"
			class="{{ .PartClass("action", "btn") }}"
			href="{{ .URL }}"
			@if(.Variant != "")
				data-variant="{{ .Variant }}"
			@endif
			@if(.Size != "")
				data-size="{{ .Size }}"
			@endif
			@attributes(.PartAttrs("action"))
		>{{ .Label }}</a>
	@endif
	@if(.URL == "")
		<button
			data-part="action"
			class="{{ .PartClass("action", "btn") }}"
			type="{{ .ButtonType() }}"
			@if(.Disabled)
				disabled
			@endif
			@if(.Variant != "")
				data-variant="{{ .Variant }}"
			@endif
			@if(.Size != "")
				data-size="{{ .Size }}"
			@endif
			@if(.HxPost != "")
				hx-post="{{ .HxPost }}"
			@endif
			@if(.HxGet != "")
				hx-get="{{ .HxGet }}"
			@endif
			@if(.HxTarget != "")
				hx-target="{{ .HxTarget }}"
			@endif
			@if(.HxSwap != "")
				hx-swap="{{ .HxSwap }}"
			@endif
			@attributes(.PartAttrs("action"))
		>{{ .Label }}</button>
	@endif

	<button
		data-part="trigger"
		class="{{ .PartClass("trigger", "btn") }}"
		type="button"
		id="{{ .TriggerID() }}"
		aria-haspopup="menu"
		aria-controls="{{ .MenuID() }}"
		aria-expanded="false"
		aria-label="{{ .TriggerName() }}"
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
	>{!! icons.CaretDown(icons.Props{}) !!}</button>

	<div
		data-part="panel"
		@if(.PartClass("panel") != "")
			class="{{ .PartClass("panel") }}"
		@endif
		id="{{ .PanelID() }}"
		data-popover
		aria-hidden="true"
		@attributes(.PartAttrs("panel"))
	>
		<div
			data-part="menu"
			@if(.PartClass("menu") != "")
				class="{{ .PartClass("menu") }}"
			@endif
			role="menu"
			id="{{ .MenuID() }}"
			aria-labelledby="{{ .TriggerID() }}"
			@attributes(.PartAttrs("menu"))
		>
			@for(at := 0; at < len(.Items); at++)
				@if(.Items[at].Separator)
					<hr role="separator">
				@elseif(.Items[at].Heading)
					<div role="presentation">{{ .Items[at].Label }}</div>
				@elseif(.Items[at].URL != "")
					<a
						data-part="item"
						@if(.PartClass("item") != "")
							class="{{ .PartClass("item") }}"
						@endif
						role="menuitem"
						href="{{ .Items[at].URL }}"
						@if(.Items[at].Variant != "")
							data-variant="{{ .Items[at].Variant }}"
						@endif
						@attributes(.PartAttrs("item"))
					>
						{{ .Items[at].Label }}
						@if(.Items[at].Shortcut != "")
							<kbd
								data-part="shortcut"
								class="{{ .PartClass("shortcut", "kbd") }}"
								@attributes(.PartAttrs("shortcut"))
							>{{ .Items[at].Shortcut }}</kbd>
						@endif
					</a>
				@else
					<button
						data-part="item"
						@if(.PartClass("item") != "")
							class="{{ .PartClass("item") }}"
						@endif
						role="menuitem"
						type="button"
						@if(.Items[at].Disabled)
							disabled
						@endif
						@if(.Items[at].Variant != "")
							data-variant="{{ .Items[at].Variant }}"
						@endif
						@attributes(.PartAttrs("item"))
					>
						{{ .Items[at].Label }}
						@if(.Items[at].Shortcut != "")
							<kbd
								data-part="shortcut"
								class="{{ .PartClass("shortcut", "kbd") }}"
								@attributes(.PartAttrs("shortcut"))
							>{{ .Items[at].Shortcut }}</kbd>
						@endif
					</button>
				@endif
			@endfor
		</div>
	</div>
</div>
