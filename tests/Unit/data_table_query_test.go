package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

func TestDataTableSearchAndFacetsShareOneNativeQueryForm(t *testing.T) {
	html := string(components.DataTable(components.DataTableProps{
		ID: "users", Label: "Users", URL: "/users?scope=platform",
		SearchName: "q", Query: "ada", SortKey: "email", SortDir: "desc",
		Filters: []components.DataTableFilter{{
			Key: "status", Label: "Status", Selected: []string{"active"},
			Options: []components.SelectOption{{Label: "Active", Value: "active"}, {Label: "Suspended", Value: "suspended"}},
		}},
		Columns: []components.TableColumn{{Label: "Email", Key: "email", Sortable: true}},
		Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "ada@example.test"}}}},
	}))

	for _, want := range []string{
		"id=\"users-query\"",
		"action=\"/users\"",
		"name=\"sort\" value=\"email\"",
		"name=\"dir\" value=\"desc\"",
		"name=\"q\"",
		"name=\"status\"",
		"hx-include=\"closest form\"",
		"Apply filters",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("query form is missing %q:\n%s", want, html)
		}
	}
	if strings.Contains(html, `name="scope"`) || strings.Contains(html, "platform") {
		t.Fatalf("a parameter of URL's own query was carried:\n%s", html)
	}
	if strings.Count(html, "<form") != 1 {
		t.Fatalf("search and facets are not one native form:\n%s", html)
	}
	if strings.Index(html, "name=\"q\"") > strings.Index(html, "</form>") ||
		strings.Index(html, "name=\"status\"") > strings.Index(html, "</form>") {
		t.Fatalf("a query control sits outside the shared form:\n%s", html)
	}
}

func TestDataTableServerRangeDoesNotInvertPastTheLastPage(t *testing.T) {
	html := string(components.DataTable(components.DataTableProps{
		ID: "users", Label: "Users", URL: "/users",
		Page: 3, Pages: 2, Total: 50, PageSize: 25,
		Columns: []components.TableColumn{{Label: "Email"}},
		Rows:    nil,
	}))
	if strings.Contains(html, "51 to 50") {
		t.Fatalf("a stale page rendered an inverted range:\n%s", html)
	}
}

// TestDataTableDoesNotCarryTheQueryOfURL holds every address and hidden field
// the table draws to its own state. URL is often the request's address, and a
// parameter copied from it -- a spoofed method, a role, a tenant -- would be
// resubmitted by every control on the page.
func TestDataTableDoesNotCarryTheQueryOfURL(t *testing.T) {
	props := components.DataTableProps{
		ID: "users", Label: "Users", URL: "/users?role=admin&_method=DELETE&tenant=2#top",
		SearchName: "q", Query: "ada", SortKey: "email", SortDir: "desc", Page: 2, Pages: 3,
		Filters: []components.DataTableFilter{{
			Key: "status", Label: "Status", Selected: []string{"active"},
			Options: []components.SelectOption{{Label: "Active", Value: "active"}},
		}},
		Columns: []components.TableColumn{{Label: "Email", Key: "email", Sortable: true}},
		Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "ada@example.test"}}}},
	}
	html := string(components.DataTable(props))

	for _, leaked := range []string{`name="role"`, "admin", "_method", "DELETE", "tenant", "#top"} {
		if strings.Contains(html, leaked) {
			t.Errorf("the table carried %q from URL's query:\n%s", leaked, html)
		}
	}
	for _, address := range []string{props.SortURL(), props.PageURL(), props.SearchAddress(), props.QueryAction()} {
		if !strings.HasPrefix(address, "/users") || strings.Contains(address, "role=") ||
			strings.Contains(address, "_method") || strings.Contains(address, "tenant") || strings.Contains(address, "#") {
			t.Errorf("address %q carries URL's own query", address)
		}
	}
	if got := props.PageURL(); got != "/users?dir=desc&q=ada&sort=email&status=active&page={page}" {
		t.Errorf("PageURL() = %q, want only the table's own state", got)
	}

	fields := props.QueryFields()
	want := []components.DataTableQueryField{{Name: "dir", Value: "desc"}, {Name: "sort", Value: "email"}}
	if len(fields) != len(want) || fields[0] != want[0] || fields[1] != want[1] {
		t.Errorf("QueryFields() = %+v, want %+v", fields, want)
	}
	if hidden := strings.Count(html, `type="hidden"`); hidden != 2 {
		t.Errorf("the query form has %d hidden fields, want the order alone:\n%s", hidden, html)
	}
	if got := (components.DataTableProps{URL: "/users?role=admin"}).QueryFields(); len(got) != 0 {
		t.Errorf("an unsorted table drew hidden fields %+v", got)
	}
}
