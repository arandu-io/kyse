package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// Every input that names its box after Name takes an ID instead, because a page
// with two dialogs that each ask for a reason drew two boxes with id="reason":
// the second label pointed at the first box, and so did the second message.
//
// Each row is rendered with Name "reason" and ID "cancel-reason". The box has
// to carry the ID, submit under the Name, and leave no id built from the Name
// anywhere in the markup -- a label, a message or a list still pointing at
// "reason" is the same defect one element over.
func TestEveryInputTakesItsOwnID(t *testing.T) {
	const name, id = "reason", "cancel-reason"
	page := refusedPage{errors: map[string]string{name: "Say why."}, old: map[string]string{}}

	for _, c := range []struct {
		name string
		html string
	}{
		{"Field", string(components.Field(components.FieldProps{Name: name, ID: id, Label: "Reason", Hint: "Optional.", Page: page}))},
		{"Textarea", string(components.Textarea(components.TextareaProps{Name: name, ID: id, Label: "Reason", Hint: "Optional.", Page: page}))},
		{"RadioGroup", string(components.RadioGroup(components.RadioGroupProps{Name: name, ID: id, Label: "Reason", Page: page,
			Options: []components.RadioOption{{Value: "price", Label: "Price"}}}))},
		{"ColorPicker", string(components.ColorPicker(components.ColorPickerProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"Password", string(components.Password(components.PasswordProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"Rating", string(components.Rating(components.RatingProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"RangeSlider", string(components.RangeSlider(components.RangeSliderProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"Combobox", string(components.Combobox(components.ComboboxProps{Name: name, ID: id, Label: "Reason", SearchURL: "/r", Page: page}))},
		{"Autocomplete", string(components.Autocomplete(components.AutocompleteProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"FileUpload", string(components.FileUpload(components.FileUploadProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"DateTimePicker", string(components.DateTimePicker(components.DateTimePickerProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"InputGroup", string(components.InputGroup(components.InputGroupProps{Name: name, ID: id, Label: "Reason", Start: "R$", Page: page}))},
		{"NumberInput", string(components.NumberInput(components.NumberInputProps{Name: name, ID: id, Label: "Reason", Page: page}))},
		{"SegmentedControl", string(components.SegmentedControl(components.SegmentedControlProps{Name: name, ID: id, Label: "Reason", Page: page,
			Options: []components.SegmentedOption{{Value: "price", Label: "Price"}}}))},
	} {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(c.html, `name="`+name) {
				t.Errorf("it no longer submits under its Name:\n%s", c.html)
			}
			if !strings.Contains(c.html, `"`+id) {
				t.Errorf("the ID is not on the markup:\n%s", c.html)
			}
			for _, attr := range []string{`id="` + name, `for="` + name, `aria-describedby="` + name, `aria-labelledby="` + name, `aria-controls="` + name, `list="` + name} {
				if strings.Contains(c.html, attr) {
					t.Errorf("an id is still built from the Name (%s...):\n%s", attr, c.html)
				}
			}
		})
	}
}

// Two dialogs on one list each ask for a reason. One was refused; the other
// was not even opened. Only the fields of the refused one draw the message and
// the typed value.
func TestOnlyTheRefusedDialogsFieldsAnswer(t *testing.T) {
	page := refusedPage{
		errors: map[string]string{"reason": "Say why."},
		old:    map[string]string{components.FormDialogField: "cancel-subscription", "reason": "Too expensive"},
	}

	refused := string(components.Field(components.FieldProps{
		Name: "reason", ID: "cancel-subscription-reason", Label: "Reason",
		Page: components.DialogPage(page, "cancel-subscription"),
	}))
	other := string(components.Field(components.FieldProps{
		Name: "reason", ID: "refund-order-reason", Label: "Reason", Value: "stored",
		Page: components.DialogPage(page, "refund-order"),
	}))

	for _, want := range []string{"Say why.", `value="Too expensive"`, `aria-invalid="true"`} {
		if !strings.Contains(refused, want) {
			t.Errorf("the refused dialog's field lost %q:\n%s", want, refused)
		}
	}
	for _, unwanted := range []string{"Say why.", "Too expensive", `aria-invalid`} {
		if strings.Contains(other, unwanted) {
			t.Errorf("a field of a dialog nobody submitted drew %q:\n%s", unwanted, other)
		}
	}
	if !strings.Contains(other, `value="stored"`) {
		t.Errorf("a field of a dialog nobody submitted lost the value it was given:\n%s", other)
	}

	if got := components.DialogPage(nil, "x").OldOr("reason", "kept"); got != "kept" {
		t.Errorf("a dialog page over no page answered %q, want the fallback", got)
	}
}

func TestAHiddenInputCarriesTheRefusedValue(t *testing.T) {
	plain := string(components.HiddenInput(components.HiddenInputProps{Name: "subscription", Value: "sub_1"}))
	if !strings.Contains(plain, `type="hidden"`) || !strings.Contains(plain, `name="subscription"`) || !strings.Contains(plain, `value="sub_1"`) {
		t.Errorf("not a hidden input carrying its value:\n%s", plain)
	}
	if strings.Contains(plain, "id=") || strings.Contains(plain, "class=") {
		t.Errorf("a hidden input wrote an id or a class nobody asked for:\n%s", plain)
	}

	page := refusedPage{old: map[string]string{"subscription": "sub_2"}}
	back := string(components.HiddenInput(components.HiddenInputProps{Name: "subscription", Value: "sub_1", Page: page}))
	if !strings.Contains(back, `value="sub_2"`) {
		t.Errorf("the refused attempt's value did not come back:\n%s", back)
	}
}

// A layout that boosts its body boosts every link in it, and htmx answers the
// click on the link before the dialog script hears it: the page navigated to
// the form instead of opening the dialog.
func TestALinkTriggerIsNotBoosted(t *testing.T) {
	link := string(components.DialogTrigger(components.DialogTriggerProps{Dialog: "d", Label: "Open", Href: "/d"}))
	if !strings.Contains(link, `hx-boost="false"`) {
		t.Errorf("the link trigger is boosted by a boosted layout:\n%s", link)
	}
}
