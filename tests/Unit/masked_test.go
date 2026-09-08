package unit_test

import (
	"os"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// The masked field, which shows one value and submits another.
//
// The visible input carries the formatted value and no field name; the hidden
// one beside it carries the name and the raw value. So what reaches the server
// is what a program stores, and nothing there has to know a mask existed.

func TestTheVisibleBoxShowsTheFormattedValueAndTheHiddenOneCarriesItRaw(t *testing.T) {
	t.Parallel()

	html := string(components.Masked(components.MaskedProps{
		Name:    "postcode",
		Pattern: "00000-000",
		Value:   "01310100",
	}))

	if !strings.Contains(html, `value="01310-100"`) {
		t.Errorf("the visible box does not show the formatted value:\n%s", html)
	}
	if !strings.Contains(html, `value="01310100"`) {
		t.Errorf("the hidden input does not carry the raw value:\n%s", html)
	}
	if !strings.Contains(html, `name="postcode"`) {
		t.Error("nothing carries the field name")
	}
	if !strings.Contains(html, `name="postcode_display"`) {
		t.Error("the visible box has no name of its own, so no browser will remember it")
	}
	if !strings.Contains(html, `type="hidden"`) {
		t.Error("there is no hidden input, so the formatted value is what submits")
	}
}

// TestTheBehaviourIsHandedThePatternItRendersWith is the property the whole
// design rests on: one declaration, read by both sides.
func TestTheBehaviourIsHandedThePatternItRendersWith(t *testing.T) {
	t.Parallel()

	html := string(components.Masked(components.MaskedProps{
		Name:    "postcode",
		Pattern: "00000-000",
	}))

	if !strings.Contains(html, `data-kyse-behavior="mask"`) {
		t.Fatalf("the field names no behaviour, so nothing formats what is typed:\n%s", html)
	}
	if !strings.Contains(html, "00000-000") {
		t.Error("the pattern the server rendered with is not handed to the behaviour")
	}
}

// TestAFieldWithTwoShapesHandsOverBoth covers the alternatives: one field that
// takes either of two patterns.
func TestAFieldWithTwoShapesHandsOverBoth(t *testing.T) {
	t.Parallel()

	// The application's own patterns. None of them ships with the framework:
	// a national document is a rule before it is a shape.
	html := string(components.Masked(components.MaskedProps{
		Name:     "document",
		Patterns: []string{"000.000.000-00", "00.000.000/0000-00"},
		Value:    "12345678900",
	}))

	if !strings.Contains(html, `value="123.456.789-00"`) {
		t.Errorf("eleven digits were not formatted with the shorter pattern:\n%s", html)
	}
	for _, pattern := range []string{"000.000.000-00", "00.000.000/0000-00"} {
		if !strings.Contains(html, pattern) {
			t.Errorf("the behaviour was not handed %q, so the field cannot grow into it", pattern)
		}
	}
}

func TestTheLongerShapeIsUsedOnceTheValueNeedsIt(t *testing.T) {
	t.Parallel()

	html := string(components.Masked(components.MaskedProps{
		Name:     "document",
		Patterns: []string{"000.000.000-00", "00.000.000/0000-00"},
		Value:    "12345678000190",
	}))

	if !strings.Contains(html, `value="12.345.678/0001-90"`) {
		t.Errorf("fourteen digits were not formatted with the longer pattern:\n%s", html)
	}
}

// TestTheFormattedValueIsRightInTheFirstFrame is why the component formats at
// all instead of leaving it to the behaviour.
//
// A box that shows the raw value until a script runs is a box that flickers on
// every page load, and on a slow connection it flickers for a while.
func TestTheFormattedValueIsRightInTheFirstFrame(t *testing.T) {
	t.Parallel()

	html := string(components.Masked(components.MaskedProps{
		Name: "phone", Pattern: "(00) 00000-0000", Value: "11999998888",
	}))

	if !strings.Contains(html, `value="(11) 99999-8888"`) {
		t.Errorf("the server rendered the raw value into the visible box:\n%s", html)
	}
}

// TestAPartialValueLeavesNoDanglingSeparator holds the small thing that makes
// typing feel right: a separator with nothing behind it moves the caret for no
// reason.
func TestAPartialValueLeavesNoDanglingSeparator(t *testing.T) {
	t.Parallel()

	for _, c := range []struct{ value, want string }{
		{"123", "123"},
		{"1234", "123.4"},
		{"123456", "123.456"},
		{"1234567", "123.456.7"},
	} {
		html := string(components.Masked(components.MaskedProps{
			Name: "document", Pattern: "000.000.000-00", Value: c.value,
		}))
		if !strings.Contains(html, `value="`+c.want+`"`) {
			t.Errorf("%q was formatted as something other than %q:\n%s", c.value, c.want, html)
		}
	}
}

// TestPartNamesReachTheHiddenInput keeps the promise the parts make: a caller
// writing a form that reads its own fields has to be able to name the one that
// submits.
func TestPartNamesReachTheHiddenInput(t *testing.T) {
	t.Parallel()

	names := components.MaskedProps{}.PartNames()
	var found bool
	for _, name := range names {
		if name == "raw" {
			found = true
		}
	}
	if !found {
		t.Errorf("PartNames = %v, and the hidden input is unreachable", names)
	}
}

// TestBothSidesOfTheMaskAgree is the claim the whole design rests on.
//
// The dictionary is written twice -- once in Go, which renders the first frame,
// and once in JavaScript, which formats every keystroke after it -- because the
// two cannot import from each other. Two tables that disagree is a box whose
// value changes the moment somebody types into it, which reads as the page
// having corrected them.
//
// So the tokens are compared: every one the component accepts is one the
// behaviour accepts, and the other way round.
func TestBothSidesOfTheMaskAgree(t *testing.T) {
	t.Parallel()

	script := servedScript(t)

	// The table in ui.js, as it is written there.
	for _, token := range []string{"'0'", "'9'", "'#'", "'A'", "'S'"} {
		if !strings.Contains(script, token+": { accepts:") {
			t.Errorf("the behaviour has no entry for %s, and the component does", token)
		}
	}

	// And nothing beyond them: a sixth token on one side is a pattern that
	// formats differently depending on which side read it.
	start := strings.Index(script, "var maskTokens = {")
	if start < 0 {
		t.Fatal("ui.js declares no mask dictionary")
	}
	end := strings.Index(script[start:], "\n\t};")
	if end < 0 {
		t.Fatal("the mask dictionary in ui.js is not closed where this test can see it")
	}
	if got := strings.Count(script[start:start+end], "accepts:"); got != 5 {
		t.Errorf("the behaviour declares %d tokens and the component declares 5", got)
	}
}

// TestTheFrameworkNamesNoNationalPattern keeps the decision that made the
// country masks come back out.
//
// A national document is a rule before it is a shape -- eleven digits are not a
// CPF until the check digits agree -- so a framework that named the shape would
// be implying it checks the rule. An application names its own, which costs it
// one line, and this stops the list from growing back one convenience at a time.
func TestTheFrameworkNamesNoNationalPattern(t *testing.T) {
	t.Parallel()

	source, err := os.ReadFile("../../components/masked.kyse.go")
	if err != nil {
		t.Fatalf("reading the component: %v", err)
	}
	for _, named := range []string{
		"CPF Pattern", "CNPJ Pattern", "CEP Pattern", "PhoneBR", "PlateBR",
	} {
		if strings.Contains(string(source), named) {
			t.Errorf("the component names %s: a national shape belongs to the application that owns the rule behind it", named)
		}
	}
}
