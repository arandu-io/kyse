package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// CardBelow says what a data table in the auto mode measures to turn into cards:
// the screen at 600px, which is what every list drew before it existed, or the
// block the table sits in, at 640px or at 720px. The three are the same markup
// under three sets of classes.

func broken(mode components.DataTableMode, below components.DataTableBreak) string {
	list := phoneList(phoneColumns, true, "ids")
	list.Mode = mode
	list.CardBelow = below
	list.RowTemplate = true
	return string(components.DataTable(list))
}

func TestTheScreenIsTheDefaultBreakAndAnUnknownOneIsTheScreen(t *testing.T) {
	zero := broken(components.DataTableModeAuto, "")
	for _, below := range []components.DataTableBreak{components.DataTableBreakScreen, "screen", "never", "Block", "block-wide "} {
		if got := broken(components.DataTableModeAuto, below); got != zero {
			t.Errorf("CardBelow %q draws something else than the default", below)
		}
	}
	// The default is byte for byte what it was: no container on the root.
	for _, banned := range []string{"@container", "@max-[", "@min-["} {
		if strings.Contains(zero, banned) {
			t.Errorf("the default draws %q", banned)
		}
	}
}

// The block sets are the screen set with its query swapped for a container query
// on the region, so a class added to one and forgotten in another fails here.
func TestTheBlockSetsAreTheScreenSetUnderAContainerQuery(t *testing.T) {
	screen := broken(components.DataTableModeAuto, components.DataTableBreakScreen)
	cases := []struct {
		below components.DataTableBreak
		query string
	}{
		{components.DataTableBreakBlock, "@max-[640px]/dt:"},
		{components.DataTableBreakBlockWide, "@max-[720px]/dt:"},
	}
	for _, c := range cases {
		html := broken(components.DataTableModeAuto, c.below)
		if strings.Contains(html, underQuery) {
			t.Errorf("%q still asks the width of the screen", c.below)
		}
		// The root is the container, and it is the only difference outside the classes.
		rooted := strings.Replace(html, `class="group/dt w-full @container/dt`, `class="group/dt w-full`, 1)
		if rooted == html {
			t.Fatalf("%q does not make the root a container:\n%s", c.below, html[:400])
		}
		if want := strings.ReplaceAll(screen, underQuery, c.query); rooted != want {
			t.Errorf("%q is not the screen classes under %s", c.below, c.query)
		}
	}
}

// Every card rule of a block set carries the query, so above it nothing applies,
// and the markup outside the classes is the same in every set.
func TestABlockBreakChangesOnlyTheClasses(t *testing.T) {
	class := regexp.MustCompile(`class="[^"]*"`)
	base := class.ReplaceAllString(broken(components.DataTableModeAuto, ""), "")
	for _, below := range []components.DataTableBreak{components.DataTableBreakBlock, components.DataTableBreakBlockWide} {
		html := broken(components.DataTableModeAuto, below)
		if got := class.ReplaceAllString(html, ""); got != base {
			t.Errorf("%q changes more than the classes", below)
		}
		for _, name := range regexp.MustCompile(`[^\s"]*@max-\[\d+px\]/dt:[^\s"]+`).FindAllString(html, -1) {
			if !strings.HasPrefix(name, "@max-[") {
				t.Errorf("%q: a card rule outside the query: %s", below, name)
			}
		}
	}
}

// The table and cards modes do not ask anything, so CardBelow changes nothing
// there and neither of them makes a container.
func TestTheTableAndCardsModesIgnoreCardBelow(t *testing.T) {
	for _, mode := range []components.DataTableMode{components.DataTableModeTable, components.DataTableModeCards} {
		want := broken(mode, "")
		for _, below := range []components.DataTableBreak{components.DataTableBreakBlock, components.DataTableBreakBlockWide} {
			if got := broken(mode, below); got != want {
				t.Errorf("mode %q with CardBelow %q draws something else than without", mode, below)
			}
		}
		if strings.Contains(want, "@container") {
			t.Errorf("mode %q makes the root a container", mode)
		}
	}
}

// The row template is drawn like a real row in a block set, as in the others.
func TestTheRowTemplateCarriesTheBlockClasses(t *testing.T) {
	for _, below := range []components.DataTableBreak{components.DataTableBreakBlock, components.DataTableBreakBlockWide} {
		html := broken(components.DataTableModeAuto, below)
		row, _ := rowTemplate(t, html)
		body := html[strings.Index(html, "<tbody"):strings.Index(html, "<template")]
		classOf := func(tag string) string { return regexp.MustCompile(`class="[^"]*"`).FindString(tag) }
		if got, want := classOf(row), classOf(squeeze(tags(body, "tr")[0])); got != want {
			t.Errorf("%q: the template's row is %s and a real row is %s", below, got, want)
		}
	}
}

func TestTheBlockBreaksNeedNoInlineStyleOrScript(t *testing.T) {
	for _, below := range []components.DataTableBreak{components.DataTableBreakBlock, components.DataTableBreakBlockWide} {
		html := broken(components.DataTableModeAuto, below)
		for _, banned := range []string{"style=", "<script", " onclick=", " x-data"} {
			if strings.Contains(html, banned) {
				t.Errorf("%q draws %q", below, banned)
			}
		}
	}
}
