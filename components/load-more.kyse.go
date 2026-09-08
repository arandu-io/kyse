//go:build kyse

package components

@go
// LoadMoreProps is the end of a list that fetches the next page into it.
//
// It replaces itself. The endpoint answers with the next rows followed by a
// fresh copy of this control pointing at the page after, so there is no
// counter kept anywhere and no state to get out of step -- the page that comes
// back carries the knowledge of where it is.
//
// It is a button by default and can watch the scroll instead. The button is
// the default deliberately: an infinite list has no end to reach, so the
// footer under it is unreachable, and somebody who scrolls past what they
// wanted cannot get back to it. Auto is right for a feed and wrong for
// anything with something after it.
//
// It publishes root, button and indicator.
type LoadMoreProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// URL is the next page.
	URL string
	// Label is the text on the button. Empty says "Load more".
	Label string
	// Auto fetches when the control scrolls into view instead of when it is
	// pressed. The button is still drawn, because a fetch that fails leaves it
	// as the way to try again.
	Auto bool
	// Margin is how far ahead of the viewport the automatic fetch starts, as a
	// CSS length: "200px" begins loading before the reader arrives. Empty
	// starts when it is visible.
	Margin string
	// Variant and Size style the button. Empty is the outline one, because
	// this is not the page's primary action.
	Variant string
	Size    string
	// Exhausted draws nothing at all, for the end of the list. It is a field
	// rather than the caller omitting the component, so the loop that draws
	// the page does not have to branch.
	Exhausted bool
}

// Text is the label, or the words used when none was given.
func (p LoadMoreProps) Text() string {
	if p.Label != "" {
		return p.Label
	}
	return "Load more"
}

// Look is the variant, and it is outline unless the caller says otherwise.
func (p LoadMoreProps) Look() string {
	if p.Variant != "" {
		return p.Variant
	}
	return "outline"
}

// Trigger is when the fetch happens: a press, or the control coming into view
// once. Once, because a control that fetched every time it was scrolled past
// would fetch the same page again on the way back up.
func (p LoadMoreProps) Trigger() string {
	if !p.Auto {
		return "click"
	}
	if p.Margin != "" {
		return "intersect once threshold:0.1, click"
	}
	return "revealed once, click"
}

// PartNames are the parts this component publishes.
func (p LoadMoreProps) PartNames() []string { return []string{"root", "indicator"} }
@endgo

@if(!.Exhausted)
	{{-- outerHTML on itself: what comes back is the next rows and the next
	     control, in place of this one. Nothing counts pages here because
	     nothing here could be told when the count is wrong. --}}
	<button
		data-part="root"
		class="{{ .RootClass("btn") }}"
		type="button"
		data-variant="{{ .Look() }}"
		@attributes(.RootAttrs())
		hx-get="{{ .URL }}"
		hx-trigger="{{ .Trigger() }}"
		hx-target="this"
		hx-swap="outerHTML"
		@if(.Size != "")
			data-size="{{ .Size }}"
		@endif
	>
		{{ .Text() }}
		<span
			data-part="indicator"
			class="{{ .PartClass("indicator", "spinner") }}"
			aria-hidden="true"
			@attributes(.PartAttrs("indicator"))
		></span>
	</button>
@endif
