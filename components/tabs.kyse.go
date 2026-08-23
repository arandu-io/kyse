//go:build kyse

package components

import "html/template"

@go
// TabsProps is one region of a page with several views of itself and a row of
// buttons that choose between them.
//
// The row is a tablist, each button is a tab that names the panel it controls,
// and each panel names the tab that labels it. That pair of references is what
// the keyboard and the screen reader both run on: the arrow keys move along the
// row, Home and End reach its ends, and only the selected tab is in the tab
// order, so one press of Tab leaves the row rather than walking through every
// button in it.
//
// Every panel is rendered, and the ones not being shown carry the hidden
// attribute. That is what makes the first paint correct before any script has
// run, and it is why a panel is a value rather than an address: switching tabs
// is not a request.
type TabsProps struct {
	// ID is what every tab and panel hangs its id off. Two tab sets on one page
	// need two, or the second set's panels answer to the first set's tabs.
	ID string
	// Tabs are the buttons and their panels, in the order they are drawn.
	Tabs []Tab
	// Active is the ID of the tab drawn selected. Empty selects the first one
	// that is not disabled.
	Active string
	// Label names the row for a screen reader, which announces it as a set
	// rather than as loose buttons.
	Label string
	// Variant is "line" for the underlined row. Empty is the default.
	Variant string
	// Vertical stacks the row down the side instead of across the top. It also
	// moves the arrow keys from left and right to up and down, which is why it
	// is a property of the set and not a class the caller adds.
	Vertical bool
}

// Tab is one button and the panel it shows.
type Tab struct {
	// ID identifies the tab within the set. It is part of two element ids, so
	// it has to be unique among the tabs, and nobody sees it.
	ID string
	// Label is the text on the button.
	Label string
	// Panel is what the tab shows. It is markup, so it comes from another
	// component or from a view.
	Panel template.HTML
	// Disabled draws the tab unusable and takes it out of the arrow-key order.
	Disabled bool
}

// TabID is the id of a tab's button.
func (p TabsProps) TabID(t Tab) string { return p.ID + "-tab-" + t.ID }

// PanelID is the id of the panel a tab controls.
func (p TabsProps) PanelID(t Tab) string { return p.ID + "-panel-" + t.ID }

// IsActive reports whether a tab is the one being shown.
//
// With no Active named it is the first tab that is not disabled, rather than
// the first tab: a set whose first entry is unavailable would otherwise open
// showing a panel nothing can select again.
func (p TabsProps) IsActive(t Tab) bool {
	if p.Active != "" {
		return t.ID == p.Active
	}
	for _, candidate := range p.Tabs {
		if candidate.Disabled {
			continue
		}
		return candidate.ID == t.ID
	}
	return false
}

// Selected is the aria-selected value of a tab.
func (p TabsProps) Selected(t Tab) string {
	if p.IsActive(t) {
		return "true"
	}
	return "false"
}

// TabIndex is a tab's place in the keyboard order.
//
// The selected tab is the single entry point into the row and the arrow keys
// reach the rest, which is what keeps a set of nine tabs from being nine stops
// on the way to the panel.
func (p TabsProps) TabIndex(t Tab) string {
	if p.IsActive(t) {
		return "0"
	}
	return "-1"
}

// Folded reports whether a panel starts hidden.
func (p TabsProps) Folded(t Tab) bool { return !p.IsActive(t) }

// Orientation is which way the row runs, which is also what decides the pair of
// arrow keys that moves along it.
func (p TabsProps) Orientation() string {
	if p.Vertical {
		return "vertical"
	}
	return "horizontal"
}
@endgo

<div class="tabs" id="{{ .ID }}">
	<div
		role="tablist"
		aria-orientation="{{ .Orientation() }}"
		@if(.Label != "")
			aria-label="{{ .Label }}"
		@endif
		@if(.Variant != "")
			data-variant="{{ .Variant }}"
		@endif
	>
		@foreach(.Tabs as tab)
			<button
				type="button"
				role="tab"
				id="{{ .TabID(tab) }}"
				aria-controls="{{ .PanelID(tab) }}"
				aria-selected="{{ .Selected(tab) }}"
				tabindex="{{ .TabIndex(tab) }}"
				@if(tab.Disabled)
					aria-disabled="true"
				@endif
			>{{ tab.Label }}</button>
		@endforeach
	</div>

	@foreach(.Tabs as tab)
		<div
			role="tabpanel"
			id="{{ .PanelID(tab) }}"
			aria-labelledby="{{ .TabID(tab) }}"
			tabindex="0"
			@if(.Folded(tab))
				hidden
			@endif
		>{!! tab.Panel !!}</div>
	@endforeach
</div>
