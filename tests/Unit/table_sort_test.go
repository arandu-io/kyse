package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A header says which way its column is ordered, to the eye and to a screen
// reader, with no catalogue to ask: the arrow follows the direction and the name
// is English until a catalogue says otherwise.
func TestATableHeaderDrawsItsDirectionAndNamesItselfInEnglish(t *testing.T) {
	arrow := func(dir string) string {
		key := "n"
		if dir == "" {
			key = ""
		}
		html := string(components.Table(components.TableProps{
			ID: "t", SortKey: key, SortDir: dir, SortURL: "/t?sort={sort}&dir={dir}",
			Columns: []components.TableColumn{{Label: "Number", Key: "n", Sortable: true}},
			Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "1"}}}},
		}))
		link := html[strings.Index(html, `data-part="sort"`):]
		return link[:strings.Index(link, "</a>")]
	}
	for dir, want := range map[string]string{
		"asc":  `aria-label="Number, sorted ascending"`,
		"desc": `aria-label="Number, sorted descending"`,
		"":     `aria-label="Number, not sorted"`,
	} {
		if got := arrow(dir); !strings.Contains(got, want) {
			t.Errorf("dir %q: the header is not named %s:\n%s", dir, want, got)
		}
	}
	// Three different arrows, one for each state.
	asc, desc, none := arrow("asc"), arrow("desc"), arrow("")
	if asc == desc || asc == none || desc == none {
		t.Error("two sort states draw the same arrow")
	}
}
