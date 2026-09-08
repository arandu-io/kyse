package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// TestNoComponentLeaksADirective is the gate that was missing on the day nine
// components shipped with `@if(...)@endif` printed on the page.
//
// The compiler recognises a directive only at the start of a line. Written
// after anything else -- `>@if(.Icon != "")` -- it is not a mistake it reports:
// it is text, and it is copied into the output verbatim. So a component
// compiles, renders, passes every other test here, and draws
// `@if(.OffIcon != "")@endifLike 12` in a browser.
//
// Nothing caught it. The parts test asks which elements are reachable, the
// render test asks whether markup came out at all, and neither reads what the
// markup says. This one does, over every component and every state the table
// already builds -- so a directive that survives into the output fails here
// rather than on a page.
func TestNoComponentLeaksADirective(t *testing.T) {
	// The whole set the compiler answers to. A directive it does not know is
	// text everywhere, which is the same defect with a different spelling.
	directives := []string{
		"@if(", "@elseif(", "@else", "@endif",
		"@foreach(", "@endforeach", "@for(", "@endfor",
		"@attributes(", "@csrf", "@extends(", "@section(", "@endsection",
		"@go", "@endgo",
	}

	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			for _, got := range c.render(components.ComponentProps{}) {
				for _, directive := range directives {
					if !strings.Contains(got, directive) {
						continue
					}
					t.Errorf("the rendered output carries %q, so the compiler read it as text:\n%s",
						directive, excerpt(got, directive))
				}
			}
		})
	}
}

// excerpt is the line the directive is on, which is what says where to look.
func excerpt(markup, directive string) string {
	for _, line := range strings.Split(markup, "\n") {
		if strings.Contains(line, directive) {
			return strings.TrimSpace(line)
		}
	}
	return markup
}
