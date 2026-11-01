package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/translation"
	"github.com/arandu-io/kyse/components"
)

var squareName = regexp.MustCompile(`<input\b[^>]*data-part="square"[^>]*aria-label="([^"]*)"`)

func squareNames(html string) []string {
	var names []string
	for _, match := range squareName.FindAllStringSubmatch(html, -1) {
		names = append(names, match[1])
	}
	return names
}

// Every square is an input a screen reader stops on, and an input with no name
// is announced as "edit text" six times over: each carries its own, with its
// place in the code.
func TestEverySquareOfAOneTimeCodeHasAName(t *testing.T) {
	html := string(components.OneTimeCode(components.OneTimeCodeProps{Name: "email_code", Label: "Email code"}))
	got := squareNames(html)
	want := []string{
		"Email code, digit 1 of 6", "Email code, digit 2 of 6", "Email code, digit 3 of 6",
		"Email code, digit 4 of 6", "Email code, digit 5 of 6", "Email code, digit 6 of 6",
	}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("the squares are named %q, want %q\n%s", got, want, html)
	}
	if n := strings.Count(html, `data-part="square"`); n != 6 {
		t.Fatalf("%d squares are drawn", n)
	}
}

// The length is the code's, not six.
func TestTheNamesOfTheSquaresCountTheLengthOfTheCode(t *testing.T) {
	html := string(components.OneTimeCode(components.OneTimeCodeProps{Name: "pin", Label: "PIN", Length: 4}))
	got := squareNames(html)
	if len(got) != 4 || got[3] != "PIN, digit 4 of 4" {
		t.Fatalf("a code of four squares names them %q", got)
	}
}

// A group with no label has the name without a stray comma in front of it.
func TestASquareOfAnUnlabelledGroupIsNamedWithoutAComma(t *testing.T) {
	html := string(components.OneTimeCode(components.OneTimeCodeProps{Name: "pin"}))
	if got := squareNames(html); got[0] != "Digit 1 of 6" {
		t.Fatalf("the first square is named %q, want %q", got[0], "Digit 1 of 6")
	}
}

// The sentence is the caller's when the caller wrote one, and the catalogue's in
// the language of the page otherwise -- Portuguese here, as the application's
// lang/pt-BR/kyse.json carries it.
func TestTheNameOfASquareIsSaidInTheLanguageOfThePage(t *testing.T) {
	loader := translation.NewArrayLoader().AddMessages("pt-BR", "kyse", translation.Lines{
		"one_time_code.square": "{label}, dígito {n} de {total}",
	}, "")
	pt := translatedPage{translator: translation.New(loader, "en", "en"), locale: "pt-BR"}

	html := string(components.OneTimeCode(components.OneTimeCodeProps{Name: "email_code", Label: "Código do e-mail", Page: pt}))
	got := squareNames(html)
	if len(got) != 6 || got[0] != "Código do e-mail, dígito 1 de 6" || got[5] != "Código do e-mail, dígito 6 de 6" {
		t.Fatalf("in Portuguese the squares are named %q", got)
	}

	// A page that carries no line for it says the English.
	empty := translatedPage{translator: translation.New(translation.NewArrayLoader(), "en", "en"), locale: "pt-BR"}
	english := squareNames(string(components.OneTimeCode(components.OneTimeCodeProps{Name: "c", Label: "Code", Page: empty})))
	if english[2] != "Code, digit 3 of 6" {
		t.Fatalf("without a line the third square is %q", english[2])
	}

	// The sentence on the props wins over the catalogue.
	own := squareNames(string(components.OneTimeCode(components.OneTimeCodeProps{
		Name: "c", Label: "Código", Page: pt, SquareLabel: "{label}: {n}/{total}",
	})))
	if own[0] != "Código: 1/6" {
		t.Fatalf("the caller's sentence gives %q", own[0])
	}
}

// A name for a square is not something a part can be given: aria- is refused on
// every part, and the component writes the one it owns.
func TestACallerCannotRenameASquare(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("an aria-label on the square part was written instead of refused")
		}
	}()
	_ = components.OneTimeCode(components.OneTimeCodeProps{
		Name: "c", Label: "Code",
		ComponentProps: components.ComponentProps{Parts: components.Parts{"square": {Attrs: components.Attrs{"aria-label": "forged"}}}},
	})
}
