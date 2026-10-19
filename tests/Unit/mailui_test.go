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

// TestTheMessageSaysItsLanguage: lang is what a screen reader picks its voice
// by, and a message in Portuguese announced as English is read with the wrong
// phonetics from the first word.
func TestTheMessageSaysItsLanguage(t *testing.T) {
	out := string(mailui.Layout(mailui.LayoutProps{Lang: "pt-BR", Heading: "Confirme seu e-mail"}))
	if !strings.Contains(out, `<html lang="pt-BR">`) {
		t.Errorf("the language given is not on the message:\n%s", out)
	}
}

// TestALanguageLeftOutIsEnglish keeps every caller written before Lang existed
// drawing what it drew.
func TestALanguageLeftOutIsEnglish(t *testing.T) {
	out := string(mailui.Layout(mailui.LayoutProps{Heading: "Confirm"}))
	if !strings.Contains(out, `<html lang="en">`) {
		t.Errorf("a message with no language is not English:\n%s", out)
	}
}

// TestALanguageThatIsNotATagIsRefused: the value lands in an attribute, and a
// lang that is not a tag is one no reader can use.
func TestALanguageThatIsNotATagIsRefused(t *testing.T) {
	for _, lang := range []string{
		`"><script>alert(1)</script>`,
		"pt BR",
		"pt--BR",
		"-pt",
		"p",
		"pt_BR",
		"português",
		"en-abcdefghi",
	} {
		out := string(mailui.Layout(mailui.LayoutProps{Lang: lang}))
		if !strings.Contains(out, `<html lang="en">`) || strings.Contains(out, "<script>") {
			t.Errorf("the language %q was not replaced by the floor:\n%s", lang, out)
		}
	}
	for _, lang := range []string{"en", "pt", "pt-BR", "es-419", "zh-Hant-TW"} {
		out := string(mailui.Layout(mailui.LayoutProps{Lang: lang}))
		if !strings.Contains(out, `<html lang="`+lang+`">`) {
			t.Errorf("the valid language %q was refused", lang)
		}
	}
}

// TestTheFallbackSaysTheSentenceItIsGiven: a message in Portuguese with an
// English line above the address reads as two messages stitched together.
func TestTheFallbackSaysTheSentenceItIsGiven(t *testing.T) {
	const pt = "Se o botão não abrir, copie este endereço:"
	out := string(mailui.FallbackWith(mailui.FallbackProps{Href: "https://x/verify?t=1&u=2", Text: pt}))

	if !strings.Contains(out, "Se o botão não abrir, copie este endereço:") {
		t.Errorf("the sentence given is not in the fallback:\n%s", out)
	}
	if strings.Contains(out, "If the button") {
		t.Errorf("the English floor is still drawn beside the sentence given:\n%s", out)
	}
	if !strings.Contains(out, "https://x/verify?t=1&amp;u=2") {
		t.Errorf("the address is not written out:\n%s", out)
	}
}

// TestTheFallbackFloorIsTheEnglishItAlwaysWas: Fallback and FallbackWith with
// no sentence draw the same bytes Fallback drew before the sentence was a
// field.
func TestTheFallbackFloorIsTheEnglishItAlwaysWas(t *testing.T) {
	const href = "https://x/verify"
	want := `<p style="margin:0;font-size:13px;line-height:1.6;color:#71717a">If the button does not open, copy this address:<br><span style="word-break:break-all">https://x/verify</span></p>`

	if got := string(mailui.Fallback(href)); got != want {
		t.Errorf("Fallback changed:\n got %s\nwant %s", got, want)
	}
	if got := string(mailui.FallbackWith(mailui.FallbackProps{Href: href})); got != want {
		t.Errorf("FallbackWith with no text is not Fallback:\n got %s\nwant %s", got, want)
	}
	if mailui.FallbackWith(mailui.FallbackProps{Text: "Copie:"}) != "" {
		t.Error("a fallback with no address drew a sentence introducing nothing")
	}
}

// TestTheFallbackSentenceIsEscaped: it comes from a catalogue, and a catalogue
// is a file somebody edits.
func TestTheFallbackSentenceIsEscaped(t *testing.T) {
	out := string(mailui.FallbackWith(mailui.FallbackProps{Href: "https://x", Text: `"><script>alert(1)</script>`}))
	if strings.Contains(out, "<script>") {
		t.Errorf("the fallback sentence was not escaped:\n%s", out)
	}
}
