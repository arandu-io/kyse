package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// The fill is a span, and a span is inline: width and height do nothing to an
// inline box, so a bar whose fill carried only a width class drew a track with
// nothing in it. The fill has to be a block as tall as the track.
func TestProgressFillIsABlockAsTallAsTheTrack(t *testing.T) {
	for name, p := range map[string]components.ProgressProps{
		"known":         {Label: "Upload", Value: 28, Max: 100},
		"indeterminate": {Label: "Import", Indeterminate: true},
	} {
		html := string(components.Progress(p))
		i := strings.Index(html, `data-part="fill"`)
		if i < 0 {
			t.Fatalf("%s: no fill part in %s", name, html)
		}
		class := html[i:]
		class = class[strings.Index(class, `class="`):]
		class = class[:strings.Index(class[7:], `"`)+8]
		for _, want := range []string{"block", "h-full"} {
			found := false
			for _, c := range strings.Fields(strings.Trim(class[len(`class=`):], `"`)) {
				if c == want {
					found = true
				}
			}
			if !found {
				t.Errorf("%s: fill class %s lacks %q, so the fill collapses to nothing", name, class, want)
			}
		}
	}
}

// "default" is what a status enum calls the unmarked badge. No rule styles a
// data-variant of "default", so it was drawn as bare text: the same status read
// as a pill in one place and as plain words in another.
func TestBadgeDefaultVariantIsTheUnmarkedBadge(t *testing.T) {
	html := string(components.Badge(components.BadgeProps{Label: "Ativa", Variant: "default"}))
	if strings.Contains(html, "data-variant") {
		t.Errorf("a default badge wrote a data-variant no rule styles: %s", html)
	}
	html = string(components.Badge(components.BadgeProps{Label: "Falhou", Variant: "destructive"}))
	if !strings.Contains(html, `data-variant="destructive"`) {
		t.Errorf("a destructive badge lost its variant: %s", html)
	}
}

// A picture that fails to load draws its alt text and the broken-image icon,
// and an absolutely placed one draws them past the circle. The name is read from
// a span, so the picture itself is decorative and has no text to draw.
func TestAvatarPictureCarriesNoAltTextButTheNameIsStillRead(t *testing.T) {
	html := string(components.Avatar(components.AvatarProps{Name: "Marina Costa", ImageURL: "/m.png"}))
	if !strings.Contains(html, `alt=""`) {
		t.Errorf("the picture still carries the name as alt text, so a broken one draws it: %s", html)
	}
	if !strings.Contains(html, `sr-only">Marina Costa<`) {
		t.Errorf("the name is no longer read for a picture avatar: %s", html)
	}
}

// A long address or URL has no place to break, and in a grid track it pushed
// the alert wider than the screen. The title and the message break anywhere.
func TestAlertTextBreaksAnywhere(t *testing.T) {
	html := string(components.Alert(components.AlertProps{Title: "t", Message: "m"}))
	if strings.Count(html, "wrap-anywhere") != 2 {
		t.Errorf("the alert's title and message do not both break anywhere: %s", html)
	}
}

// A badge wider than its container clips instead of widening the page.
func TestBadgeNeverWiderThanItsContainer(t *testing.T) {
	html := string(components.Badge(components.BadgeProps{Label: "x"}))
	if !strings.Contains(html, "max-w-full") {
		t.Errorf("the badge can outgrow its container: %s", html)
	}
}

// The first column's header was the English word "Name" in every language.
func TestStatCardLabelHeaderIsTheCallers(t *testing.T) {
	rows := []components.StatRow{{Label: "Plano Pro", Values: []string{"31"}}}
	html := string(components.StatCard(components.StatCardProps{Title: "Receita", LabelHeader: "Produto", Columns: []string{"Pedidos"}, Rows: rows}))
	if !strings.Contains(html, ">Produto</th>") || strings.Contains(html, ">Name</th>") {
		t.Errorf("the label column does not carry the caller's header: %s", html)
	}
	html = string(components.StatCard(components.StatCardProps{Title: "Receita", Columns: []string{"Pedidos"}, Rows: rows}))
	if !strings.Contains(html, ">Name</th>") {
		t.Errorf("the default header changed: %s", html)
	}
}
