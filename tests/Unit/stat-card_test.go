package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

func TestStatCardRendersRowsAndEmpty(t *testing.T) {
	full := string(components.StatCard(components.StatCardProps{
		Title:   "Connections",
		Meta:    "read at 21:44",
		Columns: []string{"Open", "Channels"},
		Rows: []components.StatRow{
			{Label: "acme", Values: []string{"12", "3"}},
			{Label: "globex", Values: []string{"0", "0"}},
		},
		Empty: components.EmptyProps{Title: "No sockets"},
	}))
	for _, want := range []string{"Connections", "read at 21:44", "Open", "Channels", "acme", "12", "globex"} {
		if !strings.Contains(full, want) {
			t.Fatalf("the card did not draw %q:\n%s", want, full)
		}
	}
	if strings.Contains(full, "No sockets") {
		t.Fatal("the empty state was drawn beside rows")
	}

	bare := string(components.StatCard(components.StatCardProps{
		Title: "Connections",
		Empty: components.EmptyProps{Title: "No sockets", Message: "Nothing is connected."},
	}))
	if !strings.Contains(bare, "No sockets") {
		t.Fatalf("a card with no rows did not draw its empty state:\n%s", bare)
	}
	if strings.Contains(bare, "<table") {
		t.Fatal("a card with no rows drew an empty table")
	}
}

func TestStatCardEscapesWhatItIsGiven(t *testing.T) {
	got := string(components.StatCard(components.StatCardProps{
		Title: `<script>alert(1)</script>`,
		Rows:  []components.StatRow{{Label: `<img onerror=alert(1)>`, Values: []string{`"&`}}},
	}))
	if strings.Contains(got, "<script>") || strings.Contains(got, "<img onerror") {
		t.Fatalf("markup survived into the output:\n%s", got)
	}
}

// TestStatCardDrawsOneTable counts instead of searching, because the defect it
// guards against is invisible to strings.Contains.
//
// A loop written around the whole table rather than around its rows renders one
// complete table per row, each holding the full row set. Every label is present,
// every column header is present, nothing is missing and nothing extra is
// named -- so an assertion that only looks for substrings passes on five copies
// of the card as happily as on one.
func TestStatCardDrawsOneTable(t *testing.T) {
	rows := make([]components.StatRow, 0, 5)
	for _, label := range []string{"tenant-1", "tenant-2", "tenant-3", "tenant-4", "tenant-5"} {
		rows = append(rows, components.StatRow{Label: label, Values: []string{"1", "2"}})
	}

	got := string(components.StatCard(components.StatCardProps{
		Title:   "Connections",
		Columns: []string{"Open", "Channels"},
		Rows:    rows,
		Empty:   components.EmptyProps{Title: "No sockets"},
	}))

	for _, c := range []struct {
		needle string
		want   int
	}{
		{"<table", 1},
		{"</table>", 1},
		{"<thead", 1},
		// One header row over the rows themselves.
		{"<tr", 1 + len(rows)},
		{"tenant-1<", 1},
		{"tenant-5<", 1},
	} {
		if n := strings.Count(got, c.needle); n != c.want {
			t.Errorf("%s appears %d time(s), want %d:\n%s", c.needle, n, c.want, got)
		}
	}
}

// TestStatCardWithNoRowsDrawsTheEmptyComponent pins the other half of the same
// branch: what a caller sees when there is nothing to draw.
//
// The card does not render a headed, bodyless table there. It renders the Empty
// component, whole, which is asserted by rendering that component separately and
// looking for its output inside the card -- so the two cannot drift apart
// without this failing.
func TestStatCardWithNoRowsDrawsTheEmptyComponent(t *testing.T) {
	empty := components.EmptyProps{Title: "No sockets", Message: "Nothing is connected."}

	got := string(components.StatCard(components.StatCardProps{
		Title: "Connections",
		Empty: empty,
	}))

	if n := strings.Count(got, "<table"); n != 0 {
		t.Errorf("a card with no rows drew %d table(s):\n%s", n, got)
	}
	if want := strings.TrimSpace(string(components.Empty(empty))); !strings.Contains(got, want) {
		t.Errorf("the card did not draw the Empty component:\nwant to find:\n%s\ngot:\n%s", want, got)
	}
	if !strings.Contains(got, `<div class="mt-4">`) {
		t.Errorf("the empty state lost its spacing wrapper:\n%s", got)
	}
}
