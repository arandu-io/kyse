//go:build kyse

package components

import (
	"html/template"

	"github.com/arandu-io/kyse/icons"
)

@go
// AccordionProps is a stack of headings, each with a section folded behind it.
//
// It is details and summary, which is the disclosure the browser already has.
// The open state, the click, Enter and Space, the announcement to a screen
// reader, and find-in-page expanding a closed section are five behaviours
// nobody writes here -- and they are the five a hand-built accordion drops one
// of. There is no role and no aria-expanded on it for the same reason: the
// elements carry both already, and writing them again is how they come to
// disagree with the state.
//
// One section open at a time is an attribute rather than a script. Two details
// elements that share a name are exclusive to each other, and the browser
// closes the one that was open when the next is unfolded.
type AccordionProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is the name the sections share, which is what makes opening one close
	// the last. Two accordions on a page need two, and an empty ID leaves every
	// section independent whatever Multiple says.
	ID string
	// Items are the sections, top to bottom.
	Items []AccordionItem
	// Multiple lets more than one section stand open at once.
	Multiple bool
}

// AccordionItem is one section: the heading that is always visible, and what is
// behind it.
type AccordionItem struct {
	// Label is the heading, and what is clicked.
	Label string
	// Content is what the section holds. It is markup rather than text, so it
	// comes from another component or from a view -- and both of those escaped
	// whatever a person wrote on the way in.
	Content template.HTML
	// Open draws this section already unfolded. When the sections are exclusive
	// the browser keeps the first of several.
	Open bool
	// Disabled makes the section neither clickable nor reachable by keyboard.
	// The stylesheet stops the pointer and the heading leaves the tab order,
	// which together are the two ways a section is opened. It is said on the
	// heading rather than on the section, because the heading is the control.
	Disabled bool
}

// Group is the name the sections share, or empty when each opens on its own.
func (p AccordionProps) Group() string {
	if p.Multiple {
		return ""
	}
	return p.ID
}
// PartNames are the parts this component publishes.
func (p AccordionProps) PartNames() []string { return []string{"root", "item", "trigger", "panel"} }
@endgo

<section
	data-part="root"
	class="{{ .RootClass("accordion") }}"
	@attributes(.RootAttrs())
	@if(.ID != "")
		id="{{ .ID }}"
	@endif
>
	@foreach(.Items as item)
		<details
			data-part="item"
			@if(.PartClass("item") != "")
				class="{{ .PartClass("item") }}"
			@endif
			@attributes(.PartAttrs("item"))
			@if(.Group() != "")
				name="{{ .Group() }}"
			@endif
			@if(item.Open)
				open
			@endif
		>
			<summary
				data-part="trigger"
				@if(.PartClass("trigger") != "")
					class="{{ .PartClass("trigger") }}"
				@endif
				@attributes(.PartAttrs("trigger"))
				@if(item.Disabled)
					aria-disabled="true"
					tabindex="-1"
				@endif
			>{{ item.Label }}{!! icons.CaretDown(icons.Props{}) !!}</summary>
			<section
				data-part="panel"
				@if(.PartClass("panel") != "")
					class="{{ .PartClass("panel") }}"
				@endif
				@attributes(.PartAttrs("panel"))
			>{!! item.Content !!}</section>
		</details>
	@endforeach
</section>
