package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A question about one row is drawn in that row's cell, and carries the row.
func TestADialogPostsTheRowItWasDrawnFor(t *testing.T) {
	html := string(components.Dialog(components.DialogProps{
		ID: "archive-114", Title: "Archive 2026-114?", Message: "It stays in the history.",
		Action: "/invoices/114/archive", ConfirmLabel: "Archive", ConfirmVariant: "destructive", Token: "csrf",
		Fields:      []components.DialogField{{Name: "return", Value: "/invoices?page=2"}, {Name: "id", Value: "114"}},
		ConfirmName: "do", ConfirmValue: "archive",
	}))
	for _, want := range []string{
		`<input type="hidden" name="return" value="/invoices?page=2">`, `<input type="hidden" name="id" value="114">`,
		`<input type="hidden" name="_token" value="csrf">`, `name="do"`, `value="archive"`, `action="/invoices/114/archive"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the dialog lacks %s:\n%s", want, html)
		}
	}
	// The words of the question are not the row's: a complete table leaves them out of its search.
	if !strings.Contains(html, `data-search="false"`) {
		t.Errorf("the dialog does not say its words are not searched:\n%s", html)
	}
}

// A dialog for the bulk bar of a table submits the form the table declared.
func TestADialogCanConfirmAFormDeclaredElsewhere(t *testing.T) {
	html := string(components.Dialog(components.DialogProps{
		ID: "archive-many", Title: "Archive the selected?", ConfirmLabel: "Archive", FormID: "invoices-bulk",
		ConfirmName: "do", ConfirmValue: "archive", Action: "/ignored",
	}))
	if !strings.Contains(html, `form="invoices-bulk"`) || !strings.Contains(html, `name="do"`) || !strings.Contains(html, `value="archive"`) {
		t.Errorf("the confirm button does not submit the table's form:\n%s", html)
	}
	if strings.Contains(html, `/ignored`) || strings.Count(html, "<form") != 1 {
		t.Errorf("a dialog confirming another form drew a form of its own:\n%s", html)
	}
}
