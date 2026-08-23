//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

@go
// CollapsibleProps is one heading and the paragraph it hides: a question with
// its answer, a section of a form nobody opens twice, the detail under a
// summary.
//
// It is <details> and <summary>, so the toggling, the keyboard, the focus and
// the aria-expanded a screen reader reads are the browser's. No script is
// loaded and none is needed, which is also why it keeps working on a page whose
// JavaScript failed.
//
// # This one or a group of them
//
// Use this for a disclosure that stands on its own -- one thing to open, and
// opening it says nothing about anything else on the screen.
//
// Use a group when the items belong to one list and are read as a set: a page
// of frequently asked questions, the sections of a long form. A group is
// written as `class="accordion"` around several <details>, the same way markup
// that needs a frame around it uses `class="card"`. That is not a second
// component to reach for -- it is the same markup, and the stylesheet draws the
// rules between the items only when they are siblings of each other. Several of
// these side by side would be several one-item lists, each with no rule between
// it and the next.
type CollapsibleProps struct {
	// Label is the heading: what opening it will show.
	Label string
	// Content is the text under it.
	Content string
	// Open draws it already open. Use it for the one item somebody came for;
	// a screen where everything starts open is a screen with no summary.
	Open bool
}
@endgo

<section class="accordion">
	<details
		@if(.Open)
			open
		@endif
	>
		<summary>{{ .Label }}{!! icons.CaretDown(icons.Props{}) !!}</summary>
		<section>{{ .Content }}</section>
	</details>
</section>
