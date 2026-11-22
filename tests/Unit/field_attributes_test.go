package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A phone's keyboard, the browser's own check and the box's length are what a
// field says about itself, so they are props of the field and not an attribute
// a caller has to know to pass through Parts.
func TestAFieldWritesTheAttributesItHasPropsFor(t *testing.T) {
	html := string(components.Field(components.FieldProps{
		Name: "zip", Label: "ZIP", Autocomplete: "postal-code",
		InputMode: "numeric", Pattern: `\d{5}`, MaxLength: 5, MinLength: 5,
	}))
	for _, want := range []string{`autocomplete="postal-code"`, `inputmode="numeric"`, `pattern="\d{5}"`, `maxlength="5"`, `minlength="5"`} {
		if !strings.Contains(html, want) {
			t.Errorf("the input does not carry %s:\n%s", want, html)
		}
	}
}

func TestAFieldWritesNoneOfThemUnlessAskedTo(t *testing.T) {
	html := string(components.Field(components.FieldProps{Name: "city", Label: "City"}))
	for _, name := range []string{"inputmode", "pattern", "maxlength", "minlength", "autocomplete"} {
		if strings.Contains(html, name+"=") {
			t.Errorf("an input nobody asked anything of writes %s:\n%s", name, html)
		}
	}
}

// What a caller passes through Parts still works, for an attribute that has no
// prop and for one that has: the prop empty leaves the name to the caller.
func TestAFieldStillTakesTheseThroughParts(t *testing.T) {
	html := string(components.Field(components.FieldProps{
		Name: "zip", Label: "ZIP",
		ComponentProps: components.ComponentProps{Parts: map[string]components.PartProps{
			"input": {Attrs: components.Attrs{"inputmode": "numeric", "enterkeyhint": "next"}},
		}},
	}))
	for _, want := range []string{`inputmode="numeric"`, `enterkeyhint="next"`} {
		if !strings.Contains(html, want) {
			t.Errorf("the input lost %s from Parts:\n%s", want, html)
		}
	}
}

// A value a caller wrote cannot leave the attribute it is in.
func TestAFieldEscapesWhatItWritesInTheNewAttributes(t *testing.T) {
	html := string(components.Field(components.FieldProps{
		Name: "x", Label: "X", InputMode: `"><script>alert(1)</script>`, Pattern: `"><script>alert(2)</script>`,
	}))
	if strings.Contains(html, "<script>") {
		t.Errorf("a prop escaped its attribute:\n%s", html)
	}
}

// The select and the dropzone say their own width, in the utilities layer,
// because the stylesheet gives a label and a native select width:fit-content in
// the layer below the field's own width:100% and by a more specific selector. A
// caller who wants another width still has Parts, with the important marker
// against the utility.
func TestWhatAFieldHoldsSaysItsWidth(t *testing.T) {
	selectHTML := string(components.Select(components.SelectProps{Name: "s", Label: "S", Options: []components.SelectOption{{Value: "a", Label: "A"}}}))
	if !strings.Contains(selectHTML, `class="select w-full"`) {
		t.Errorf("the select does not say it fills the field:\n%s", selectHTML)
	}
	upload := string(components.FileUpload(components.FileUploadProps{Name: "f", Label: "F"}))
	if !strings.Contains(upload, `class="file-upload w-full"`) {
		t.Errorf("the dropzone does not say it fills the field:\n%s", upload)
	}
	narrower := string(components.Select(components.SelectProps{
		Name: "s", Label: "S",
		ComponentProps: components.ComponentProps{Parts: map[string]components.PartProps{"input": {Class: "w-48!"}}},
	}))
	if !strings.Contains(narrower, `class="select w-full w-48!"`) {
		t.Errorf("a caller's width is not added after the component's:\n%s", narrower)
	}
}
