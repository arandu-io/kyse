package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// An entry of a menu asks "are you sure" by opening a Dialog drawn beside the
// menu. Before MenuItem.Dialog the only way to ask from a menu was HxConfirm,
// which htmx turns into the browser's own confirm box -- and a menu entry had
// no way to name a Kyse Dialog at all.
//
// The entry has to name the dialog the way DialogTrigger does, because the
// view layer's script reads that attribute and nothing else; it must not carry
// a request of its own, or the click would send what the dialog was meant to
// ask about; and it announces that it opens a dialog.
func TestAMenuEntryOpensADialogAndSendsNothing(t *testing.T) {
	html := string(components.DropdownMenu(components.DropdownMenuProps{
		ID: "actions", Label: "Actions",
		Items: []components.MenuItem{
			{Label: "Archive", Variant: "destructive", Dialog: "archive-meter"},
			{Label: "Rename", HxGet: "/rename"},
		},
	}))

	entry := between(t, html, `id="actions-item-0"`, `</button>`)
	for _, want := range []string{`data-dialog-open="archive-meter"`, `aria-haspopup="dialog"`, `data-variant="destructive"`, `type="button"`} {
		if !strings.Contains(entry, want) {
			t.Errorf("the entry lacks %s:\n%s", want, entry)
		}
	}
	for _, refused := range []string{"hx-post", "hx-get", "hx-confirm", "href="} {
		if strings.Contains(entry, refused) {
			t.Errorf("an entry that opens a dialog also carries %s:\n%s", refused, entry)
		}
	}

	other := between(t, html, `id="actions-item-1"`, `</button>`)
	if strings.Contains(other, "data-dialog-open") || strings.Contains(other, "aria-haspopup") {
		t.Errorf("an entry with no Dialog claims to open one:\n%s", other)
	}
}

// A link that opens a dialog is left out of a boosted layout, as DialogTrigger
// does: htmx answers the click on the link itself, before the document the
// dialog script listens on hears it, and the page would navigate instead.
func TestALinkEntryThatOpensADialogIsNotBoosted(t *testing.T) {
	html := string(components.DropdownMenu(components.DropdownMenuProps{
		ID: "actions", Label: "Actions",
		Items: []components.MenuItem{{Label: "Archive", URL: "/meters/1/archive", Dialog: "archive-meter"}},
	}))
	entry := between(t, html, `id="actions-item-0"`, `</a>`)
	for _, want := range []string{`href="/meters/1/archive"`, `data-dialog-open="archive-meter"`, `hx-boost="false"`} {
		if !strings.Contains(entry, want) {
			t.Errorf("the link entry lacks %s:\n%s", want, entry)
		}
	}
}

func between(t *testing.T, html, from, to string) string {
	t.Helper()
	start := strings.Index(html, from)
	if start < 0 {
		t.Fatalf("the markup has no %s:\n%s", from, html)
	}
	// The id is written after the tag's first attributes, so the entry starts at the tag.
	start = strings.LastIndex(html[:start], "<")
	end := strings.Index(html[start:], to)
	if end < 0 {
		t.Fatalf("the markup has no %s after %s:\n%s", to, from, html)
	}
	return html[start : start+end]
}
