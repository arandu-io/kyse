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
		"name=\"scope\" value=\"platform\"",
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
