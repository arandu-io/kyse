//go:build kyse

package components

@go
// SkeletonProps is the grey rectangle standing where something has not
// arrived yet.
//
// It is one rectangle. A card that is waiting is four of these in the shape of
// the card it will become, written where the card is written, because what
// makes the wait readable is the arrangement -- and an arrangement decided
// inside a component is one the screen it serves cannot change.
//
// It has no colour and no animation of its own: both come from the stylesheet,
// which is what keeps every rectangle on a page pulsing together rather than
// each on its own clock.
type SkeletonProps struct {
	// Shape is what is being waited for: "line" for a line of text, "block" for
	// a picture or a chart, "circle" for an avatar. Empty means "line".
	Shape string
	// Width is how much of the space across it fills: "full", "3/4", "2/3",
	// "1/2" or "1/4". Empty means "full".
	//
	// A paragraph whose last line runs the full width does not read as a
	// paragraph, which is the whole reason this field exists. A circle ignores
	// it, being as wide as it is tall.
	Width string
	// Label is what assistive technology says while the rectangle is on screen.
	//
	// Empty is the ordinary case and marks it hidden from assistive technology:
	// a screen of grey boxes announced one at a time says nothing that the one
	// announcement for the region does not say better. Give the label to one
	// rectangle per region and leave the rest silent.
	Label string
}

// Geometry is the shape and the reach, as the classes that draw them.
//
// One answer rather than one per field, because the two decide together: a
// circle has no width to set, being as wide as it is tall, and a shape with a
// width appended separately leaves the class attribute holding an empty half.
//
// The names are written out rather than built from Shape and Width, because a
// class name assembled at run time is a name that never appears in the source
// the stylesheet is compiled from, and the rule for it is never emitted.
func (p SkeletonProps) Geometry() string {
	if p.Shape == "circle" {
		return "size-10 shrink-0 rounded-full"
	}

	shape := "h-4"
	if p.Shape == "block" {
		shape = "aspect-video"
	}

	switch p.Width {
	case "3/4":
		return shape + " w-3/4"
	case "2/3":
		return shape + " w-2/3"
	case "1/2":
		return shape + " w-1/2"
	case "1/4":
		return shape + " w-1/4"
	default:
		return shape + " w-full"
	}
}
@endgo

<div
	class="skeleton {{ .Geometry() }}"
	@if(.Label != "")
		role="status"
		aria-busy="true"
		aria-label="{{ .Label }}"
	@endif
	@if(.Label == "")
		aria-hidden="true"
	@endif
></div>
