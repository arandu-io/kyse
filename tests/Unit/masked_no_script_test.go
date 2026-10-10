package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/hesape/mask"
	"github.com/arandu-io/kyse/components"
)

// TestAMaskedFieldPostsWhatWasTypedWithoutScript: with no script running the
// form sends what somebody typed in the box, under the name the caller gave,
// and nothing else under it.
//
// It used to send the value the page was drawn with. The box that took the
// typing was named Name_display and the field with the caller's name was hidden
// and filled only by the behaviour, so the server read an empty field on a new
// form and the rejected value back on a corrected one.
func TestAMaskedFieldPostsWhatWasTypedWithoutScript(t *testing.T) {
	for name, c := range map[string]struct {
		props components.MaskedProps
		typed string
	}{
		"raw digits":     {components.MaskedProps{Name: "postcode", Pattern: "00000-000"}, "01310100"},
		"as formatted":   {components.MaskedProps{Name: "postcode", Pattern: "00000-000"}, "01310-100"},
		"over a value":   {components.MaskedProps{Name: "postcode", Pattern: "00000-000", Value: "20040002"}, "01310100"},
		"refused before": {components.MaskedProps{Name: "postcode", Pattern: "00000-000", Page: page{old: map[string]string{"postcode": "2004"}, errs: map[string]string{"postcode": "Too short."}}}, "01310-100"},
		"two shapes":     {components.MaskedProps{Name: "postcode", Patterns: []string{"00000-000", "00000-0000"}}, "013101000"},
	} {
		markup := string(components.Masked(c.props))

		posted := postWithoutScript(markup, c.typed)
		if got := posted["postcode"]; len(got) != 1 || got[0] != c.typed {
			t.Errorf("%s: without script the form posts postcode=%q, want exactly [%q]\n%s", name, got, c.typed, markup)
		}
		if len(posted) != 1 {
			t.Errorf("%s: without script the form posts %v, want the one field and nothing beside it", name, posted)
		}
	}
}

// TestWhatArrivesWithoutScriptUnmasksToWhatArrivesWithIt: the value posted
// without script is as typed, and the server that stores the field unmasks it
// with the same pattern. That answers the raw value the behaviour would have
// sent, whichever way somebody typed it.
func TestWhatArrivesWithoutScriptUnmasksToWhatArrivesWithIt(t *testing.T) {
	pattern := "00000-000"
	markup := string(components.Masked(components.MaskedProps{Name: "postcode", Pattern: pattern}))
	for _, typed := range []string{"01310100", "01310-100", "01310 100"} {
		got := postWithoutScript(markup, typed)["postcode"]
		if len(got) != 1 {
			t.Fatalf("%q: posted %v, want one value", typed, got)
		}
		if raw := mask.Pattern(pattern).Unmask(got[0]); raw != "01310100" {
			t.Errorf("%q posted as %q unmasks to %q, want the raw 01310100", typed, got[0], raw)
		}
	}
}

// TestTheMaskedBoxIsTheFieldUntilTheBehaviourRuns: the visible box carries the
// field's name and the id a Label points at, the hidden input waits disabled,
// and the behaviour is handed the name the box takes when it hands the field's
// name over.
func TestTheMaskedBoxIsTheFieldUntilTheBehaviourRuns(t *testing.T) {
	markup := string(components.Masked(components.MaskedProps{
		Name: "postcode", Pattern: "00000-000", Value: "01310100",
	}))

	var box, raw input
	for _, in := range inputs(markup) {
		switch in["data-part"] {
		case "input":
			box = in
		case "raw":
			raw = in
		}
	}
	if box["name"] != "postcode" || box["id"] != "postcode" || box["value"] != "01310-100" || box.has("disabled") {
		t.Errorf("the visible box is not the field without script: %v", box)
	}
	if raw["name"] != "postcode" || raw["type"] != "hidden" || raw["value"] != "01310100" || !raw.has("disabled") {
		t.Errorf("the hidden input does not wait disabled with the raw value for the behaviour: %v", raw)
	}
	if !strings.Contains(markup, "&#34;displayName&#34;:&#34;postcode_display&#34;") {
		t.Errorf("the behaviour is not handed the name the box takes once the hidden input submits:\n%s", markup)
	}
}

// TestADisabledMaskedFieldPostsNothing: Disabled takes the field out of the
// form, and with the hidden input drawn disabled too that is now true of both.
// The hidden input used to post the drawn value from a field shown disabled.
func TestADisabledMaskedFieldPostsNothing(t *testing.T) {
	markup := string(components.Masked(components.MaskedProps{
		Name: "postcode", Pattern: "00000-000", Value: "01310100", Disabled: true,
	}))
	if posted := postWithoutScript(markup, "99999999"); len(posted) != 0 {
		t.Errorf("a disabled field posts %v, want nothing", posted)
	}
}
