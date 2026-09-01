//go:build kyse

package components

@go
// SeparatorProps is the line between two things.
//
// It is an <hr>, which already carries the separator role, plus the orientation
// -- because a rule turned on its side is still a separator to a person looking
// at it and is announced as a horizontal one otherwise.
//
// # Which of the two it is
//
// A separator that stands between two groups of content is part of the
// structure, and is announced: it is what says the settings above and the
// danger below are not the same list. A separator that is only a line -- the
// one inside a toolbar, between two buttons that were already visibly apart --
// is announced as a break that does not exist, and is one more thing read out
// before the thing somebody came for. That is what Decorative is for.
//
// # Where its look comes from
//
// The classes here draw it standing on its own. Inside a button group, a field
// or an input group the stylesheet sizes the rule itself, off the role, and
// these classes would win over it: a separator in one of those is written as
// <hr role="separator"> directly, with no classes at all.
// It publishes one part, root, which is the rule itself.
type SeparatorProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Orientation is "horizontal" or "vertical". Empty is horizontal.
	//
	// A vertical one takes its height from the row it sits in, so it needs a
	// parent that has one -- a flex row -- and draws nothing in a block
	// container, where there is no height to stretch to.
	Orientation string
	// Decorative takes it out of the accessibility tree: the line is drawn and
	// nothing is announced. Use it when the separation is already obvious
	// without it and the rule is there to look right.
	Decorative bool
}

// Direction is the orientation, defaulting to horizontal.
func (p SeparatorProps) Direction() string {
	if p.Orientation == "vertical" {
		return "vertical"
	}
	return "horizontal"
}

// PartNames are the parts this component publishes.
func (p SeparatorProps) PartNames() []string { return []string{"root"} }
@endgo

<hr
	data-part="root"
	@if(.Decorative)
		role="none" aria-hidden="true"
	@else
		role="separator" aria-orientation="{{ .Direction() }}"
	@endif
	@if(.Direction() == "vertical")
		class="{{ .RootClass("bg-border w-px shrink-0 self-stretch border-0") }}"
	@else
		class="{{ .RootClass("bg-border h-px w-full shrink-0 border-0") }}"
	@endif
	@attributes(.RootAttrs())
>
