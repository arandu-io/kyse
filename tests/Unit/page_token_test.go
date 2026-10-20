package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// A form that posts carries the session's token, and the page every screen
// embeds already holds it. These hold that each component that posts asks the
// page, so a view writes Page: .Page and never copies the token into a second
// field -- which one project did seventy-eight times.

var _ components.TokenPage = view.Page{}

// tokenForms renders every component that posts, given the page and the token
// written on its props.
func tokenForms(page components.Page, token string) map[string]string {
	bulk := []components.ButtonProps{{Label: "Archive", Type: "submit"}}
	columns := []components.TableColumn{{Label: "Number"}}
	rows := []components.TableRow{{Key: "1", Cells: []components.TableCell{{Text: "2026-114"}}}}
	return map[string]string{
		"Dialog": string(components.Dialog(components.DialogProps{
			ID: "delete", Title: "Delete?", Action: "/invoices/7", Method: "delete", Page: page, Token: token,
		})),
		"FormDialog": string(components.FormDialog(components.FormDialogProps{
			ID: "rename", Title: "Rename", Action: "/invoices/7", Method: "patch", Page: page, Token: token,
		})),
		"Table": string(components.Table(components.TableProps{
			ID: "invoices", SelectName: "ids", BulkActions: bulk, Columns: columns, Rows: rows, Page: page, Token: token,
		})),
		"DataTable": string(components.DataTable(components.DataTableProps{
			ID: "invoices", URL: "/invoices", SelectName: "ids", BulkActions: bulk, Columns: columns, Rows: rows,
			PageData: page, Token: token,
		})),
	}
}

func TestAFormThatPostsCarriesThePagesToken(t *testing.T) {
	page := view.Page{Token: "token-from-the-page"}
	for name, html := range tokenForms(page, "") {
		if !strings.Contains(html, `name="_token" value="token-from-the-page"`) {
			t.Errorf("%s did not post the page's token:\n%s", name, html)
		}
	}
}

func TestATokenWrittenOnThePropsWinsOverThePages(t *testing.T) {
	page := view.Page{Token: "token-from-the-page"}
	for name, html := range tokenForms(page, "token-from-the-props") {
		if !strings.Contains(html, `name="_token" value="token-from-the-props"`) {
			t.Errorf("%s did not post the token written on its props:\n%s", name, html)
		}
		if strings.Contains(html, "token-from-the-page") {
			t.Errorf("%s posted the page's token over the one it was given:\n%s", name, html)
		}
	}
}

// A page that answers only the field questions has no token to give, and the
// forms post what they posted before there was a page to ask.
func TestAPageWithNoTokenChangesNothing(t *testing.T) {
	before := tokenForms(nil, "token-from-the-props")
	after := tokenForms(page{}, "token-from-the-props")
	for name := range before {
		if before[name] != after[name] {
			t.Errorf("%s drew differently for a page with no token:\n%s\n---\n%s", name, before[name], after[name])
		}
	}
}
