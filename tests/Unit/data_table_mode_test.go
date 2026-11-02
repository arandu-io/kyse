package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// The mode of a data table: auto is the card layout under 600px that every list
// drew before the mode existed, table keeps the table at every width with no
// card class on any element, and cards keeps the cards. They are three sets of
// classes on the same markup.

func moded(mode components.DataTableMode, template bool) string {
	list := phoneList(phoneColumns, true, "")
	list.Mode = mode
	list.RowTemplate = template
	return string(components.DataTable(list))
}

const (
	underQuery  = "max-[600px]:"
	alwaysQuery = "min-[0px]:"
)

func TestAutoIsTheDefaultAndAnUnknownModeIsAuto(t *testing.T) {
	zero := moded("", false)
	if !strings.Contains(zero, underQuery+"flex") {
		t.Fatal("the zero mode is not the card layout under 600px")
	}
	for _, mode := range []components.DataTableMode{components.DataTableModeAuto, "auto", "never", "Table"} {
		if got := moded(mode, false); got != zero {
			t.Errorf("mode %q draws something else than the default", mode)
		}
	}
}

func TestTheTableModeDrawsNoCardClassAndNoContainer(t *testing.T) {
	html := moded(components.DataTableModeTable, false)
	for _, banned := range []string{underQuery, alwaysQuery, "@container", "@max-[", "@min-["} {
		if strings.Contains(html, banned) {
			t.Errorf("the table mode draws %q", banned)
		}
	}
	// It is still the same markup: the table, a row per row, and each cell naming
	// its column so a script and a screen reader have it.
	for _, kept := range []string{"<table", `data-part="row"`, `data-label="Customer"`, `data-column="customer"`, "group/dt"} {
		if !strings.Contains(html, kept) {
			t.Errorf("the table mode lost %q", kept)
		}
	}
}

func TestTheCardsModeIsTheAutoClassesUnderAQueryThatIsAlwaysTrue(t *testing.T) {
	auto := moded(components.DataTableModeAuto, true)
	cards := moded(components.DataTableModeCards, true)
	if strings.Contains(cards, underQuery) {
		t.Error("the cards mode asks the width of the region")
	}
	if want := strings.ReplaceAll(auto, underQuery, alwaysQuery); cards != want {
		t.Errorf("the cards mode is not the auto classes under %s\n got: %s\nwant: %s", alwaysQuery, cards, want)
	}
}

func TestTheModeIsNotInTheMarkupOutsideTheClasses(t *testing.T) {
	// Strip every class attribute: what is left is the same in the three modes.
	class := regexp.MustCompile(`class="[^"]*"`)
	base := class.ReplaceAllString(moded(components.DataTableModeAuto, true), "")
	for _, mode := range []components.DataTableMode{components.DataTableModeTable, components.DataTableModeCards} {
		if got := class.ReplaceAllString(moded(mode, true), ""); got != base {
			t.Errorf("mode %q changes more than the classes", mode)
		}
	}
}

func TestNoRowTemplateUnlessAsked(t *testing.T) {
	if strings.Contains(moded(components.DataTableModeAuto, false), "row-template") {
		t.Error("a template is drawn that nobody asked for")
	}
}

// rowTemplate is the cells of the template: the first row of it.
func rowTemplate(t *testing.T, html string) (row string, cells []string) {
	t.Helper()
	start := strings.Index(html, `<template data-part="row-template">`)
	if start < 0 {
		t.Fatalf("no row template in\n%s", html)
	}
	inner := html[start:]
	inner = inner[:strings.Index(inner, "</template>")]
	return squeeze(tags(inner, "tr")[0]), squeezeAll(tags(inner, "td"))
}

// squeeze is a tag with its whitespace made one space: the template is indented
// where it stands, and the indentation is not what is compared.
func squeeze(tag string) string { return strings.Join(strings.Fields(tag), " ") }

func squeezeAll(tags []string) []string {
	out := make([]string, len(tags))
	for i, tag := range tags {
		out[i] = squeeze(tag)
	}
	return out
}

// The template is the row the server draws: its row and each of its cells carry
// the same class attribute as a real row's, in every mode.
func TestTheRowTemplateCarriesTheClassesOfARealRow(t *testing.T) {
	for _, mode := range []components.DataTableMode{components.DataTableModeAuto, components.DataTableModeTable, components.DataTableModeCards} {
		html := moded(mode, true)
		row, cells := rowTemplate(t, html)

		body := html[strings.Index(html, "<tbody"):strings.Index(html, "<template")]
		realRow := squeeze(tags(body, "tr")[0])
		classOf := func(tag string) string { return regexp.MustCompile(`class="[^"]*"`).FindString(tag) }
		if classOf(row) != classOf(realRow) {
			t.Errorf("mode %q: the template's row is\n%s\nand a real row is\n%s", mode, row, realRow)
		}
		// Row 1 of phoneList has a value in every column but the last, so the
		// cells it draws in front of a value are the cells the template draws.
		var real []string
		for _, td := range tags(body, "td") {
			if strings.Contains(td, "data-column") || strings.Contains(td, "data-label") {
				real = append(real, squeeze(td))
			}
		}
		if len(cells) != len(phoneColumns) {
			t.Fatalf("mode %q: the template has %d cells, want one per column (%d)", mode, len(cells), len(phoneColumns))
		}
		for i := 0; i < len(phoneColumns)-1; i++ {
			if cells[i] != real[i] {
				t.Errorf("mode %q: cell %d of the template is\n%s\nand a real cell is\n%s", mode, i, cells[i], real[i])
			}
		}
	}
}

func TestTheRowTemplateIsOutsideTheBody(t *testing.T) {
	html := moded(components.DataTableModeTable, true)
	body := html[strings.Index(html, "<tbody"):strings.Index(html, "</tbody>")]
	if strings.Contains(body, "<template") {
		t.Error("the template is inside the body, where the list's own scripts would count it as a row")
	}
}
