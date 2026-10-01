package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A form that submits with GET puts every field in the address, and an address
// is kept by the history, the server's access log and any proxy in between. A
// CSRF token there is a token handed to all of them, so a form only carries it
// when it posts.

func TestDialogSendsTheTokenOnlyWhenItPosts(t *testing.T) {
	cases := []struct {
		method string
		posts  bool
	}{
		{"", true},
		{"post", true},
		{"POST", true},
		{"delete", true},
		{"PUT", true},
		{"patch", true},
		{"get", false},
		{"GET", false},
		{"dialog", false},
		{"head", false},
	}
	for _, tc := range cases {
		html := string(components.Dialog(components.DialogProps{
			ID: "confirm", Title: "Delete?", Action: "/invoices/7", Method: tc.method, Token: "secret-token",
		}))
		hasToken := strings.Contains(html, `name="_token"`) || strings.Contains(html, "secret-token")
		if hasToken != tc.posts {
			t.Errorf("Method %q: token written = %v, want %v:\n%s", tc.method, hasToken, tc.posts, html)
		}
		if got := (components.DialogProps{Method: tc.method}).SendsToken(); got != tc.posts {
			t.Errorf("Method %q: SendsToken() = %v, want %v", tc.method, got, tc.posts)
		}
	}
}

func TestTableBulkFormSendsTheTokenOnlyWhenItPosts(t *testing.T) {
	cases := []struct {
		method string
		posts  bool
	}{
		{"", true},
		{"post", true},
		{"POST", true},
		{"get", false},
		{"GET", false},
		// A method a form cannot send is submitted by the browser as GET.
		{"delete", false},
	}
	for _, tc := range cases {
		props := components.TableProps{
			ID: "invoices", SelectName: "ids", BulkMethod: tc.method, Token: "secret-token",
			BulkActions: []components.ButtonProps{{Label: "Archive", Type: "submit"}},
			Columns:     []components.TableColumn{{Label: "Number"}},
			Rows:        []components.TableRow{{Key: "1", Cells: []components.TableCell{{Text: "2026-114"}}}},
		}
		html := string(components.Table(props))
		if !strings.Contains(html, `id="invoices-bulk"`) && !strings.Contains(html, "table-bulk") {
			t.Fatalf("no bulk form was drawn:\n%s", html)
		}
		hasToken := strings.Contains(html, `name="_token"`) || strings.Contains(html, "secret-token")
		if hasToken != tc.posts {
			t.Errorf("BulkMethod %q: token written = %v, want %v:\n%s", tc.method, hasToken, tc.posts, html)
		}
	}
}
