//go:build kyse

package components

import "strconv"

@go
// MenubarBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const MenubarBehavior = "menubar"

// MenubarProps is the row of menus along the top of an application: File,
// Edit, View.
//
// It is one tab stop for the whole row, and the arrow keys move along it --
// left and right between the menus, down into one, up and down inside it. That
// is the contract people already know from every desktop application, and it
// is the reason a menubar is not a row of dropdowns: a row of six dropdowns is
// six tab stops, and moving between them means closing one and opening the
// next by hand.
//
// Once a menu is open, moving sideways opens the next one rather than closing
// everything -- which is the behaviour that makes a menubar quick and is the
// part most reimplementations leave out.
//
// It publishes root, trigger, panel, menu, item and shortcut.
type MenubarProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what every trigger, panel and menu hangs its ids off. Two menubars
	// on one page need two.
	ID string
	// Label names the row. Without it the menubar announces with no subject.
	Label string
	// Menus are the top-level menus, in the order they read.
	Menus []MenubarMenu
}

// MenubarMenu is one menu of the row: the word on the bar, and what is under it.
type MenubarMenu struct {
	// Label is the word on the bar.
	Label string
	// Items are the entries, separators and headings under it.
	Items []MenuItem
	// Disabled draws the menu unavailable, and the arrow keys pass over it.
	Disabled bool
}

// TriggerID, PanelID and MenuID are the ids the three elements of one menu use
// to point at one another.
func (p MenubarProps) TriggerID(at int) string {
	return p.ID + "-trigger-" + strconv.Itoa(at)
}

// PanelID is the id of the element that opens under one menu.
func (p MenubarProps) PanelID(at int) string {
	return p.ID + "-panel-" + strconv.Itoa(at)
}

// MenuID is the id of the menu inside that panel.
func (p MenubarProps) MenuID(at int) string {
	return p.ID + "-menu-" + strconv.Itoa(at)
}

// First is the index of the menu holding the tab stop: the first one that can
// take it.
func (p MenubarProps) First() int {
	for at, menu := range p.Menus {
		if !menu.Disabled {
			return at
		}
	}
	return -1
}

// Stop is whether the menu at this index holds the tab stop.
func (p MenubarProps) Stop(at int) bool { return at == p.First() }

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
func (p MenubarProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: MenubarBehavior}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p MenubarProps) PartNames() []string {
	return []string{"root", "trigger", "panel", "menu", "item", "shortcut"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("menubar") }}"
	id="{{ .ID }}"
	role="menubar"
	aria-label="{{ .Label }}"
	aria-orientation="horizontal"
	@attributes(.RootAttrs())
>
	@for(at := 0; at < len(.Menus); at++)
		<div class="menubar-menu">
			<button
				data-part="trigger"
				class="{{ .PartClass("trigger", "btn") }}"
				type="button"
				role="menuitem"
				data-variant="ghost"
				data-size="sm"
				id="{{ .TriggerID(at) }}"
				aria-haspopup="menu"
				aria-controls="{{ .MenuID(at) }}"
				aria-expanded="false"
				@attributes(.PartAttrs("trigger"))
				@if(.Stop(at))
					tabindex="0"
				@endif
				@if(!.Stop(at))
					tabindex="-1"
				@endif
				@if(.Menus[at].Disabled)
					disabled
				@endif
			>{{ .Menus[at].Label }}</button>

			<div
				data-part="panel"
				@if(.PartClass("panel") != "")
					class="{{ .PartClass("panel") }}"
				@endif
				id="{{ .PanelID(at) }}"
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
					id="{{ .MenuID(at) }}"
					aria-labelledby="{{ .TriggerID(at) }}"
					@attributes(.PartAttrs("menu"))
				>
					@for(line := 0; line < len(.Menus[at].Items); line++)
						@if(.Menus[at].Items[line].Separator)
							<hr role="separator">
						@elseif(.Menus[at].Items[line].Heading)
							<div role="presentation">{{ .Menus[at].Items[line].Label }}</div>
						@elseif(.Menus[at].Items[line].URL != "")
							<a
								data-part="item"
								@if(.PartClass("item") != "")
									class="{{ .PartClass("item") }}"
								@endif
								role="menuitem"
								href="{{ .Menus[at].Items[line].URL }}"
								tabindex="-1"
								@attributes(.PartAttrs("item"))
								@if(.Menus[at].Items[line].Variant != "")
									data-variant="{{ .Menus[at].Items[line].Variant }}"
								@endif
							>
								{{ .Menus[at].Items[line].Label }}
								@if(.Menus[at].Items[line].Shortcut != "")
									<kbd
										data-part="shortcut"
										class="{{ .PartClass("shortcut", "kbd") }}"
										@attributes(.PartAttrs("shortcut"))
									>{{ .Menus[at].Items[line].Shortcut }}</kbd>
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
								tabindex="-1"
								@attributes(.PartAttrs("item"))
								@if(.Menus[at].Items[line].Disabled)
									disabled
								@endif
								@if(.Menus[at].Items[line].Variant != "")
									data-variant="{{ .Menus[at].Items[line].Variant }}"
								@endif
							>
								{{ .Menus[at].Items[line].Label }}
								@if(.Menus[at].Items[line].Shortcut != "")
									<kbd
										data-part="shortcut"
										class="{{ .PartClass("shortcut", "kbd") }}"
										@attributes(.PartAttrs("shortcut"))
									>{{ .Menus[at].Items[line].Shortcut }}</kbd>
								@endif
							</button>
						@endif
					@endfor
				</div>
			</div>
		</div>
	@endfor
</div>
