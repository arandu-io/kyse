package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse"
	"github.com/arandu-io/kyse/components"
)

// TestTheCallersClassIsAddedAndNotSubstituted is the decision this type is
// built on: a caller who wants one more unit of padding does not have to know,
// and repeat, every class the component draws itself with.
func TestTheCallersClassIsAddedAndNotSubstituted(t *testing.T) {
	c := components.ComponentProps{Class: "w-full"}

	if got, want := c.PartClass("root", "btn"), "btn w-full"; got != want {
		t.Fatalf("PartClass = %q, want %q", got, want)
	}
	if got, want := c.PartClass("root", ""), "w-full"; got != want {
		t.Fatalf("a component with no class of its own = %q, want %q", got, want)
	}
	if got, want := (components.ComponentProps{}).PartClass("root", "btn"), "btn"; got != want {
		t.Fatalf("a caller who added nothing = %q, want %q", got, want)
	}
}

// TestAClassIsNotWrittenTwice covers the caller who repeats what the component
// already draws with. Two rules under one name is the same page and an
// unreadable attribute.
func TestAClassIsNotWrittenTwice(t *testing.T) {
	c := components.ComponentProps{Class: "btn p-8 p-8"}

	if got, want := c.PartClass("root", "btn card"), "btn card p-8"; got != want {
		t.Fatalf("PartClass = %q, want %q", got, want)
	}
}

// TestClassIsTheShorthandForTheRootPart pins that the two spellings are one
// thing, and which wins when a caller writes both.
//
// Neither answer is obviously right there, so the more deliberate form is the
// one that survives -- and the test exists so it stays that way rather than
// depending on the order two fields happen to be read in.
func TestClassIsTheShorthandForTheRootPart(t *testing.T) {
	short := components.ComponentProps{Class: "w-full"}
	long := components.ComponentProps{Parts: components.Parts{"root": {Class: "w-full"}}}

	if a, b := short.PartClass("root", "btn"), long.PartClass("root", "btn"); a != b {
		t.Fatalf("the shorthand gives %q and the long form gives %q", a, b)
	}

	both := components.ComponentProps{
		Class: "shorthand",
		Parts: components.Parts{"root": {Class: "explicit"}},
	}
	if got, want := both.PartClass("root", "btn"), "btn explicit"; got != want {
		t.Fatalf("with both written, PartClass = %q, want %q", got, want)
	}
}

// TestAPartReachesOnlyItsOwnElement is the promise the part vocabulary makes.
// A class written for the label may not land on the button.
func TestAPartReachesOnlyItsOwnElement(t *testing.T) {
	c := components.ComponentProps{
		Parts: components.Parts{
			"label":  {Class: "text-xs"},
			"button": {Class: "rounded-xl"},
		},
	}

	if got, want := c.PartClass("label", "label"), "label text-xs"; got != want {
		t.Fatalf("the label part = %q, want %q", got, want)
	}
	if got, want := c.PartClass("button", "btn"), "btn rounded-xl"; got != want {
		t.Fatalf("the button part = %q, want %q", got, want)
	}
	if got, want := c.PartClass("group", "button-group"), "button-group"; got != want {
		t.Fatalf("a part the caller said nothing about = %q, want %q", got, want)
	}
}

// TestAttrsReachThePartTheyWereWrittenFor is PartClass's assertion for the
// other half: attributes are addressed by part too, and the root shorthand
// works the same way.
func TestAttrsReachThePartTheyWereWrittenFor(t *testing.T) {
	c := components.ComponentProps{
		Attrs: components.Attrs{"data-testid": "actions"},
		Parts: components.Parts{"button": {Attrs: components.Attrs{"data-action": "archive"}}},
	}

	if got := c.PartAttrs("root")["data-testid"]; got != "actions" {
		t.Fatalf(`PartAttrs("root")["data-testid"] = %q, want "actions"`, got)
	}
	if got := c.PartAttrs("button")["data-action"]; got != "archive" {
		t.Fatalf(`PartAttrs("button")["data-action"] = %q, want "archive"`, got)
	}
	if got := c.PartAttrs("label"); len(got) != 0 {
		t.Fatalf(`PartAttrs("label") = %v, want nothing`, got)
	}
}

// TestTheRootGathersEverythingWrittenOnIt is the shape a component's markup
// depends on: one bag, so the root is three lines rather than six conditionals.
//
// The events are the reason it has to be a bag at all. data-kyse-on-click is a
// name assembled from data, and markup cannot write one -- the compiler reads
// the name off the source to choose the escape, and there is no source to read.
func TestTheRootGathersEverythingWrittenOnIt(t *testing.T) {
	got := components.ComponentProps{
		Attrs:    components.Attrs{"data-testid": "actions"},
		Theme:    "admin",
		Behavior: components.Behavior{Name: "message-actions", Props: map[string]any{"confirm": true}},
		Events:   components.Events{"click": "archive-message"},
	}.RootAttrs()

	for name, want := range map[string]string{
		"data-testid":        "actions",
		"data-theme":         "admin",
		"data-kyse-behavior": "message-actions",
		"data-kyse-props":    `{"confirm":true}`,
		"data-kyse-on-click": "archive-message",
	} {
		if got[name] != want {
			t.Errorf("RootAttrs()[%q] = %q, want %q", name, got[name], want)
		}
	}

	if bare := (components.ComponentProps{}).RootAttrs(); len(bare) != 0 {
		t.Errorf("a component nobody added anything to carries %v", bare)
	}
}

// TestTheComponentsOwnValueWinsTheCollision: a caller writing
// data-kyse-behavior by hand is asking for the field beside it, and the field
// is what the component reads.
func TestTheComponentsOwnValueWinsTheCollision(t *testing.T) {
	got := components.ComponentProps{
		Attrs:    components.Attrs{"data-kyse-behavior": "written-by-hand"},
		Behavior: components.Behavior{Name: "declared"},
	}.RootAttrs()

	if got["data-kyse-behavior"] != "declared" {
		t.Fatalf("data-kyse-behavior = %q, want %q", got["data-kyse-behavior"], "declared")
	}
}

// TestTheRootCarriesTheScopedStylesClass is what attaches a scoped block to the
// element it was written for. A component should not have to remember it, which
// is why the root has a method of its own.
func TestTheRootCarriesTheScopedStylesClass(t *testing.T) {
	block := kyse.CSS("& { gap: 6px; }")
	c := components.ComponentProps{Class: "w-full", Style: block}

	got := c.RootClass("card")
	for _, want := range []string{"card", "w-full", block.Class()} {
		if !strings.Contains(got, want) {
			t.Fatalf("RootClass = %q, want it to carry %q", got, want)
		}
	}

	if got := (components.ComponentProps{}).RootClass("card"); got != "card" {
		t.Fatalf("with no scoped block, RootClass = %q, want %q", got, "card")
	}
}

// TestAnEmptyBlockHasNoClass is the case that would otherwise put a class on
// every element of every component: the zero CSS hashes to something, and the
// something would be written.
func TestAnEmptyBlockHasNoClass(t *testing.T) {
	if got := kyse.CSS("").Class(); got != "" {
		t.Fatalf("the zero block has the class %q, want none", got)
	}
	if got := kyse.CSS("& { gap: 6px; }").Class(); got == "" {
		t.Fatal("a block that was written has no class")
	}
}

// TestBehaviorPropsAreJSONOrNothing holds the two ends: props become the JSON
// the behaviour reads, and a set that will not encode drops whole rather than
// half.
//
// Half a set is the failure worth naming: the behaviour mounts, reads the props
// it expects, finds one missing, and does something almost right.
func TestBehaviorPropsAreJSONOrNothing(t *testing.T) {
	b := components.Behavior{Name: "message-actions", Props: map[string]any{"confirm": true}}
	if got, want := b.PropsJSON(), `{"confirm":true}`; got != want {
		t.Fatalf("PropsJSON = %q, want %q", got, want)
	}

	if got := (components.Behavior{Props: map[string]any{"confirm": true}}).PropsJSON(); got != "" {
		t.Fatalf("props with no behaviour to receive them = %q, want nothing", got)
	}
	if got := (components.Behavior{Name: "x"}).PropsJSON(); got != "" {
		t.Fatalf("a behaviour with no props = %q, want nothing", got)
	}

	unencodable := components.Behavior{Name: "x", Props: map[string]any{"ch": make(chan int)}}
	if got := unencodable.PropsJSON(); got != "" {
		t.Fatalf("a set that will not encode came back as %q, want nothing", got)
	}
}

// TestEventNamesAreSorted is what keeps a component's markup from rendering two
// different documents for one set of props. A Go map has no order of its own.
func TestEventNamesAreSorted(t *testing.T) {
	e := components.Events{"keydown": "c", "click": "a", "input": "b"}

	got := e.EventNames()
	want := []string{"click", "input", "keydown"}
	if len(got) != len(want) {
		t.Fatalf("EventNames = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("EventNames = %v, want %v", got, want)
		}
	}

	if got := (components.Events{}).EventNames(); got != nil {
		t.Fatalf("an empty set = %v, want nil", got)
	}
}

// TestTheThemeOnAnInstanceBeatsTheGlobal is the one rule of the cascade that
// lives in Go; every other level of it is decided by CSS.
//
// Configure is not called here, on purpose: it panics on a second call, so a
// test that configured the process would decide what every other test in this
// package sees. The zero value is the state a project that configured nothing
// is in, and it is the one this can assert against.
func TestTheThemeOnAnInstanceBeatsTheGlobal(t *testing.T) {
	if got := (components.ComponentProps{}).ThemeName(); got != "" {
		t.Fatalf("with nothing configured, ThemeName = %q, want none", got)
	}
	if got := (components.ComponentProps{Theme: "admin"}).ThemeName(); got != "admin" {
		t.Fatalf("ThemeName = %q, want %q", got, "admin")
	}
}
