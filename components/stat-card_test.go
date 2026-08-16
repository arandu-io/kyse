package components_test

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
