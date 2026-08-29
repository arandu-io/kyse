package icons

import (
	"html/template"
	"strings"
)

// Props is what an icon is drawn from.
//
// One field, and it is the one an icon gets wrong. Size and colour come from
// the stylesheet -- see the package comment -- and the geometry is the function
// you called, so what is left to decide is whether this icon is something a
// screen reader should read out.
type Props struct {
	// Label is what assistive technology says instead of the picture.
	//
	// Empty means decorative, which is the common case: an icon beside the word
	// "Delete" is drawn for the eye, and announcing "trash, Delete" reads the
	// button twice. Set it when the icon is alone -- an icon-only button with no
	// label is announced as "button", and nothing on the screen says which one.
	Label string
}

// draw writes the two paths of a duotone icon into the svg every icon shares.
//
// The svg is written here rather than stored per icon so that the 1512
// generated functions hold path geometry and nothing else. That is what makes
// the guarantee cheap to state: an icon cannot carry a script, an event
// handler or a third element, because no icon carries markup at all.
//
// Two paths and one colour: the shade is the same currentColor at a fifth of
// its strength, so an icon takes the colour of the text around it exactly as it
// did at one weight. The opacity is written here rather than carried by each
// icon, for the same reason the svg is.
func draw(p Props, shade, d string) template.HTML {
	var b strings.Builder
	b.Grow(len(shade) + len(d) + 160)

	b.WriteString(`<svg viewBox="0 0 256 256" width="1em" height="1em" fill="currentColor" `)
	if p.Label == "" {
		b.WriteString(`aria-hidden="true"`)
	} else {
		// role="img" as well as the label: without a role the element is a
		// group, and a group with a name is a landmark some screen readers
		// announce and others walk past.
		b.WriteString(`role="img" aria-label="`)
		b.WriteString(template.HTMLEscapeString(p.Label))
		b.WriteString(`"`)
	}
	b.WriteString(`>`)
	// An empty shade draws no element. Upstream leaves one out where there is
	// nothing to fill, and a path with no d is a path a validator complains
	// about and a reader cannot see.
	if shade != "" {
		b.WriteString(`<path d="`)
		b.WriteString(shade)
		b.WriteString(`" opacity=".2"/>`)
	}
	b.WriteString(`<path d="`)
	b.WriteString(d)
	b.WriteString(`"/></svg>`)

	return template.HTML(b.String())
}
