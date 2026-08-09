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

// draw writes one path into the svg every icon shares.
//
// The svg is written here rather than stored per icon so that the 1512
// generated functions hold path geometry and nothing else. That is what makes
// the guarantee cheap to state: an icon cannot carry a script, an event
// handler or a second element, because no icon carries markup at all.
func draw(p Props, d string) template.HTML {
	var b strings.Builder
	b.Grow(len(d) + 128)

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
	b.WriteString(`><path d="`)
	b.WriteString(d)
	b.WriteString(`"/></svg>`)

	return template.HTML(b.String())
}
