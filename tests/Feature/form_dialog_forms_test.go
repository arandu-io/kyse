package feature

import (
	"html/template"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

var formTagPattern = regexp.MustCompile(`<form\b|</form>`)

// formDepth is the deepest nesting of <form> elements in the markup. A browser
// refuses a form inside a form, so anything over one means the dialog was drawn
// as a form around the forms of its body.
func formDepth(markup string) int {
	depth, deepest := 0, 0
	for _, tag := range formTagPattern.FindAllString(markup, -1) {
		if tag == "</form>" {
			depth--
			continue
		}
		depth++
		if depth > deepest {
			deepest = depth
		}
	}
	return deepest
}

func tabbedDialog(page components.Page) string {
	return string(components.FormDialog(components.FormDialogProps{
		ID: "update", Title: "Update", Forms: true, Page: page,
		Body: []template.HTML{
			template.HTML(`<form method="post" action="/product"><input name="price"><button type="submit">Save product</button></form>`),
			template.HTML(`<form method="post" action="/discount"><input name="discount"><button type="submit">Save discount</button></form>`),
		},
	}))
}

// TestAFormDialogForFormsIsNotAFormAroundForms: a dialog that carries the
// panels of a Tabs, each with its own form and address, must not be a form
// itself, and must not draw the submit, the hidden fields or the footer of the
// one form it no longer has.
func TestAFormDialogForFormsIsNotAFormAroundForms(t *testing.T) {
	markup := tabbedDialog(nil)
	if depth := formDepth(markup); depth != 1 {
		t.Fatalf("forms nest %d deep, want 1 (the body's own forms, side by side):\n%s", depth, markup)
	}
	if !regexp.MustCompile(`<div\s+data-part="content"`).MatchString(markup) {
		t.Fatalf("the panel is not a div:\n%s", markup)
	}
	for _, drawn := range []string{`data-part="submit"`, `data-part="footer"`, `name="_dialog"`, `name="_token"`, `action=""`} {
		if strings.Contains(markup, drawn) {
			t.Fatalf("the dialog drew %s with no form of its own:\n%s", drawn, markup)
		}
	}
	for _, kept := range []string{`id="update-title"`, `data-part="body"`, `data-part="close"`, `action="/product"`, `action="/discount"`, `id="update-dismiss"`} {
		if !strings.Contains(markup, kept) {
			t.Fatalf("the dialog lost %s:\n%s", kept, markup)
		}
	}
}

// TestAFormDialogWithoutFormsStaysOneForm: the default is what it always was.
func TestAFormDialogWithoutFormsStaysOneForm(t *testing.T) {
	markup := string(components.FormDialog(components.FormDialogProps{ID: "plain", Title: "Plain", Action: "/x", Token: "t"}))
	if !regexp.MustCompile(`<form\s+data-part="content"`).MatchString(markup) {
		t.Fatalf("the dialog is no longer a form:\n%s", markup)
	}
	for _, kept := range []string{`data-part="submit"`, `data-part="footer"`, `name="_dialog"`, `name="_token"`} {
		if !strings.Contains(markup, kept) {
			t.Fatalf("the dialog lost %s:\n%s", kept, markup)
		}
	}
}

type oldPage struct{ components.Page }

func (oldPage) OldOr(key, fallback string) string {
	if key == components.FormDialogField {
		return "update"
	}
	return fallback
}

// TestAFormDialogForFormsOpensWhereItWasSentBack: a refused panel posts the
// dialog's own name, and the page that comes back opens the same dialog.
func TestAFormDialogForFormsOpensWhereItWasSentBack(t *testing.T) {
	if markup := tabbedDialog(oldPage{}); !regexp.MustCompile(`<dialog[^>]*\bopen\b`).MatchString(markup) {
		t.Fatalf("the dialog stayed closed on the page it was sent back to:\n%s", markup)
	}
}
