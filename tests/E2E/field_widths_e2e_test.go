//go:build e2e

package e2e

import (
	"encoding/json"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// fieldsForm is a form of the kinds of field that put a control or a label of
// their own under a .field: a choice from a list, a file area, a text box. Its
// card is 28rem wide, the width an application gives a form.
func fieldsForm() string {
	return `<form id="form" class="flex max-w-md flex-col gap-4">` +
		string(components.Field(components.FieldProps{Name: "city", Label: "City"})) +
		string(components.Select(components.SelectProps{
			Name: "status", Label: "Status", Placeholder: "Choose",
			Options: []components.SelectOption{{Value: "a", Label: "Open"}, {Value: "b", Label: "Paid"}},
		})) +
		string(components.Select(components.SelectProps{
			Name: "tags", Label: "Tags", Multiple: true, Size: 3,
			Options: []components.SelectOption{{Value: "a", Label: "One"}, {Value: "b", Label: "Two"}},
		})) +
		string(components.FileUpload(components.FileUploadProps{Name: "photos", Label: "Photos"})) +
		`</form>`
}

// The width of what a field holds is the width of the field, whatever the
// stylesheet says about a label or a native select on its own: Basecoat gives
// both width:fit-content, in a layer the field's own width:100% shares and
// loses to by specificity, so a select in a field was as wide as its longest
// option and the dropzone as wide as its prompt. The components say width in
// the markup, in the utilities layer, which wins.
func TestWhatAFieldHoldsTakesTheWidthOfTheField(t *testing.T) {
	for _, width := range []int{1280, 390} {
		page := open(t, serve(t, true), width)
		page.Visit("/fields")
		var got struct {
			Form, Text, Select, Multiple, Dropzone, Label float64
		}
		out := eval(t, page, `(function () {
			var w = function (sel) { return document.querySelector(sel).getBoundingClientRect().width; };
			return JSON.stringify({Form: w('#form'), Text: w('#city'), Select: w('#status'), Multiple: w('#tags'), Dropzone: w('[data-part=dropzone]'), Label: w('label[for=status]')});
		})()`)
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatal(err)
		}
		for name, v := range map[string]float64{"text box": got.Text, "select": got.Select, "multiple select": got.Multiple, "dropzone": got.Dropzone} {
			if v < got.Form-1 || v > got.Form+1 {
				t.Errorf("at %dpx the %s is %.0fpx wide, want the %.0fpx of the form", width, name, v, got.Form)
			}
		}
		if got.Label >= got.Form-1 {
			t.Errorf("at %dpx the label of a field fills the row (%.0fpx of %.0f): a label stays as wide as its text", width, got.Label, got.Form)
		}
		shot(t, page, "fields-"+map[int]string{1280: "1280", 390: "390"}[width])
	}
}

// The search box stays inside the card that holds it, however narrow the card:
// it keeps the width it is given while there is room and gives way when there
// is not, instead of pushing out of the card.
func TestTheSearchBoxStaysInsideItsCard(t *testing.T) {
	app := serve(t, true)
	for _, c := range []struct {
		path  string
		width int
	}{{"/invoices", 390}, {"/invoices", 320}, {"/cards", 390}, {"/narrow", 390}} {
		page := open(t, app, c.width)
		page.Visit(c.path)
		out := eval(t, page, `(function () {
			var search = document.querySelector('[data-part=search]').getBoundingClientRect();
			var card = (document.getElementById('card') || document.querySelector('[data-part=root]')).getBoundingClientRect();
			return JSON.stringify({left: search.left >= card.left - 0.5, right: search.right <= card.right + 0.5, search: search.right, card: card.right});
		})()`)
		var got struct {
			Left, Right  bool
			Search, Card float64
		}
		if err := json.Unmarshal([]byte(out), &got); err != nil {
			t.Fatal(err)
		}
		if !got.Left || !got.Right {
			t.Errorf("%s at %dpx: the search box ends at %.0fpx, the card at %.0fpx", c.path, c.width, got.Search, got.Card)
		}
	}
}
