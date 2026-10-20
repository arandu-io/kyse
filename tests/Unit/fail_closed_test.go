package unit

import (
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A value the escapers refuse stops the component, and what it returns then is
// nothing. A prefix would end inside the refused attribute with its quote never
// closed, and the browser would read the page after it as attributes of the
// still-open element -- so the whole component is dropped, never half of it.
//
// That is the answer to a value that arrived as data. An attribute name the
// caller wrote is not drawn as nothing: it stops the render with a
// *RefusedAttribute -- see TestAnAttrCannotCarryAScript.
func TestARefusedValueDrawsNothingRatherThanHalfAComponent(t *testing.T) {
	cases := []struct {
		name string
		html string
	}{
		{"Link with a script address", string(components.Link(components.LinkProps{Label: "Docs", URL: "javascript:alert(1)"}))},
	}
	for _, tc := range cases {
		if tc.html != "" {
			t.Errorf("%s returned %q, want nothing", tc.name, tc.html)
		}
	}
	if got := components.Link(components.LinkProps{Label: "Docs", URL: "/docs"}); got == "" {
		t.Error("a safe Link drew nothing")
	}
}
