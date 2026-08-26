package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/mailui"
)

// TestTheMessageIsTableBasedAndInlineStyled.
//
// Not preference. Outlook renders with Word, which has no flexbox and no grid;
// Gmail strips <style>; an external stylesheet is not fetched at all. A message
// that used the class names the rest of the application uses would arrive
// unstyled in most of the clients that receive it.
func TestTheMessageIsTableBasedAndInlineStyled(t *testing.T) {
	out := string(mailui.Layout(mailui.LayoutProps{
		Brand: "Arandu", Heading: "Confirm your email address",
		Body:   mailui.Paragraph("Hello Ada."),
		Footer: "You are receiving this because of an action on your account.",
	}))

	if !strings.Contains(out, "<table") {
		t.Error("the layout is not a table: Word has no flexbox")
	}
	if strings.Contains(out, "<style") || strings.Contains(out, `class="`) {
		t.Error("the message carries a stylesheet or a class, and neither reaches an inbox")
	}
	if !strings.Contains(out, "style=") {
		t.Error("nothing is styled inline")
	}
}

// TestEverythingIsEscaped.
//
// Every function here returns template.HTML, which tells the engine the value is
// already safe -- so the engine escapes nothing, and a heading is often a
// subject somebody typed.
func TestEverythingIsEscaped(t *testing.T) {
	const attack = `"><script>alert(1)</script>`

	for name, got := range map[string]string{
		"Heading":   string(mailui.Layout(mailui.LayoutProps{Heading: attack})),
		"Paragraph": string(mailui.Paragraph(attack)),
		"Button":    string(mailui.Button(mailui.ButtonProps{Label: attack, Href: "https://x"})),
		"Href":      string(mailui.Button(mailui.ButtonProps{Label: "Go", Href: attack})),
		"Panel":     string(mailui.Panel(attack)),
		"Small":     string(mailui.Small(attack)),
		"Fallback":  string(mailui.Fallback(attack)),
	} {
		if strings.Contains(got, "<script>") {
			t.Errorf("%s did not escape its input:\n%s", name, got)
		}
	}
}

// TestAButtonRefusesADangerousDestination covers the part HTML escaping cannot
// answer: a javascript URL contains no syntax that needs escaping, but the
// browser executes it when the reader follows the link.
func TestAButtonRefusesADangerousDestination(t *testing.T) {
	for _, href := range []string{
		"javascript:alert(1)",
		"JaVaScRiPt:alert(1)",
		"java\tscript:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"//evil.example/steal",
	} {
		if got := mailui.Button(mailui.ButtonProps{Label: "Open", Href: href}); got != "" {
			t.Errorf("a button rendered the refused destination %q:\n%s", href, got)
		}
	}
}

// TestAButtonKeepsAnAllowedDestination is the other half of the refusal. A
// guard that removes every link is safe from scripts and useless as mail UI.
func TestAButtonKeepsAnAllowedDestination(t *testing.T) {
	got := string(mailui.Button(mailui.ButtonProps{
		Label: "Open",
		Href:  "/account?tab=security&from=mail",
	}))

	if !strings.Contains(got, `href="/account?tab=security&amp;from=mail"`) {
		t.Errorf("an allowed destination was refused or changed incorrectly:\n%s", got)
	}
}

// TestAButtonWithNoDestinationDrawsNothing: a call to action that goes nowhere
// is worse than none, because the reader clicks it.
func TestAButtonWithNoDestinationDrawsNothing(t *testing.T) {
	if mailui.Button(mailui.ButtonProps{Label: "Confirm"}) != "" {
		t.Error("a button with no href was drawn")
	}
	if mailui.Button(mailui.ButtonProps{Href: "https://x"}) != "" {
		t.Error("a button with no label was drawn")
	}
	if mailui.Fallback("") != "" {
		t.Error("an empty fallback drew an empty address")
	}
}

// TestThePreheaderIsHiddenAndPadded.
//
// It is what the inbox list shows beside the subject. Left out, the client shows
// the first words of the body -- usually "Hello there," which tells nobody
// anything. The padding stops the client filling the rest of the preview line
// with whatever comes next in the markup.
func TestThePreheaderIsHiddenAndPadded(t *testing.T) {
	out := string(mailui.Layout(mailui.LayoutProps{
		Heading: "Confirm", Preheader: "One click and the account is yours.",
	}))

	if !strings.Contains(out, "One click and the account is yours.") {
		t.Fatal("the preheader is not in the message")
	}
	if !strings.Contains(out, "display:none") {
		t.Error("the preheader is visible in the body as well as in the list")
	}
	if !strings.Contains(out, "&zwnj;") {
		t.Error("the preheader is not padded: the client fills the preview with the markup after it")
	}
}

// TestNoRemoteImage: every client blocks them until the reader allows it, so a
// logo delivered that way is a broken-image box on the first impression.
func TestNoRemoteImage(t *testing.T) {
	out := string(mailui.Layout(mailui.LayoutProps{
		Brand: "Arandu", Heading: "Hello", Body: mailui.Paragraph("Text."),
	}))
	if strings.Contains(out, "<img") {
		t.Error("the message carries an image")
	}
	if !strings.Contains(out, "Arandu") {
		t.Error("the brand is not in the message at all")
	}
}
