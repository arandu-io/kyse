package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// Every control of a paged list that replaces the region takes the region out
// of the answer by the same selector. An endpoint that answers the whole page
// -- which is what it answers without HTMX -- then still swaps one list for
// one list, instead of drawing the page inside itself.
func TestEveryRegionSwapSelectsTheRegionFromTheAnswer(t *testing.T) {
	rows := make([]components.TableRow, 0, 2)
	for _, email := range []string{"ada@example.test", "grace@example.test"} {
		rows = append(rows, components.TableRow{Cells: []components.TableCell{{Text: email}}})
	}
	html := string(components.DataTable(components.DataTableProps{
		ID: "users", Label: "Users", URL: "/users",
		SearchName: "q", Total: 60, Page: 2, Pages: 3,
		Filters: []components.DataTableFilter{{
			Key: "status", Label: "Status",
			Options: []components.SelectOption{{Label: "Active", Value: "active"}},
		}},
		Columns: []components.TableColumn{{Label: "Email", Key: "email", Sortable: true}},
		Rows:    rows,
	}))

	swaps := regexp.MustCompile(`<[a-z]+[^>]*hx-target="#users"[^>]*>`).FindAllString(html, -1)
	if len(swaps) < 5 {
		t.Fatalf("expected the search, the facet, the header and the pages to swap the region, found %d:\n%s", len(swaps), html)
	}
	for _, tag := range swaps {
		if strings.Contains(tag, `hx-swap="outerHTML"`) && !strings.Contains(tag, `hx-select="#users"`) {
			t.Errorf("a control replaces the region without selecting it from the answer:\n%s", tag)
		}
	}
}
