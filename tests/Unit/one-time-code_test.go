package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// describedBy reads the aria-describedby the group carries, and reports whether
// it carries one at all.
var describedBy = regexp.MustCompile(`aria-describedby="([^"]*)"`)

// TestAOneTimeCodeSaysWhyItWasRefused is the rejected attempt: the code came
// back wrong, and the field has to say so in words as well as in red.
//
// The squares used to be marked invalid and nothing else, so a refused code was
// a row of red boxes with no sentence under them and nothing for a screen reader
// to announce.
func TestAOneTimeCodeSaysWhyItWasRefused(t *testing.T) {
	refused := page{errs: map[string]string{"email_code": "The code is wrong or has expired."}}

	html := string(components.OneTimeCode(components.OneTimeCodeProps{
		Name: "email_code", Label: "Email code", Page: refused,
	}))

	if !strings.Contains(html, `id="email_code-error"`) ||
		!strings.Contains(html, ">The code is wrong or has expired.</p>") {
		t.Errorf("the rejection was not drawn under its id:\n%s", html)
	}
	if !strings.Contains(html, `data-part="message"`) {
		t.Errorf("the rejection is not the message part a caller can reach:\n%s", html)
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Errorf("a refused code is not marked invalid:\n%s", html)
	}
	if got := describedBy.FindStringSubmatch(html); got == nil || got[1] != "email_code-error" {
		t.Errorf("the rejection is shown and not announced: aria-describedby is %v\n%s", got, html)
	}
	if n := strings.Count(html, "aria-describedby="); n != 1 {
		t.Errorf("aria-describedby is written %d times, want once, on the group:\n%s", n, html)
	}
}

// TestAOneTimeCodeKeepsWhatTheCallerPointedAt: a caller who drew a sentence of
// their own -- where the code was sent -- still has it announced when the code
// is refused, after the rejection rather than instead of it.
func TestAOneTimeCodeKeepsWhatTheCallerPointedAt(t *testing.T) {
	props := components.OneTimeCodeProps{
		Name: "email_code", ID: "verify-code", Label: "Email code",
		DescribedBy: "sent-to",
		Page:        page{errs: map[string]string{"email_code": "Required."}},
	}

	html := string(components.OneTimeCode(props))
	if got := describedBy.FindStringSubmatch(html); got == nil || got[1] != "verify-code-error sent-to" {
		t.Errorf("aria-describedby is %v, want the rejection and then the caller's id\n%s", got, html)
	}
	// The id hangs off the square's id, as it does on the other controls that
	// take an ID, so two code fields with different ids never share a message.
	if !strings.Contains(html, `id="verify-code-error"`) {
		t.Errorf("the rejection is not drawn under the id the group points at:\n%s", html)
	}

	props.Page = page{}
	html = string(components.OneTimeCode(props))
	if got := describedBy.FindStringSubmatch(html); got == nil || got[1] != "sent-to" {
		t.Errorf("an accepted field lost the caller's description: aria-describedby is %v\n%s", got, html)
	}
}

// TestAnAcceptedOneTimeCodePointsAtNothing: no rejection, no sentence, and no
// attribute naming an element that is not on the page. A dangling
// aria-describedby is read as nothing at best, and at worst names whatever else
// on the screen happens to carry that id.
func TestAnAcceptedOneTimeCodePointsAtNothing(t *testing.T) {
	for name, state := range map[string]components.Page{
		"no page":         nil,
		"another's error": page{errs: map[string]string{"password": "Too short."}},
	} {
		html := string(components.OneTimeCode(components.OneTimeCodeProps{
			Name: "email_code", Label: "Email code", Page: state,
		}))

		for _, unwanted := range []string{`data-part="message"`, "-error", "aria-describedby", "aria-invalid", "Too short."} {
			if strings.Contains(html, unwanted) {
				t.Errorf("%s: an accepted field wrote %s:\n%s", name, unwanted, html)
			}
		}
	}
}

// TestARefusedOneTimeCodeEscapesItsMessage: the rejection is text the server
// wrote, often around what somebody typed, and it is written as text.
func TestARefusedOneTimeCodeEscapesItsMessage(t *testing.T) {
	html := string(components.OneTimeCode(components.OneTimeCodeProps{
		Name: "email_code",
		Page: page{errs: map[string]string{"email_code": "<script>alert(1)</script>"}},
	}))
	if strings.Contains(html, "<script>") {
		t.Errorf("the rejection was written as markup:\n%s", html)
	}
}

// partTags matches the start tag of every element drawn as one part.
func partTags(html, part string) []string {
	return regexp.MustCompile(`<[a-z]+\s[^>]*data-part="`+part+`"[^>]*>`).FindAllString(html, -1)
}

// TestAOneTimeCodeWritesTheAttributesOfItsParts: the squares and the hidden
// field are published parts, so the attributes a caller writes for them arrive,
// on every square and on the field that submits.
//
// They used to be dropped without a word. The class written for the same part
// arrived, so the call looked half honoured, and a test id or a data attribute
// a script reads from the squares was simply not there.
func TestAOneTimeCodeWritesTheAttributesOfItsParts(t *testing.T) {
	html := string(components.OneTimeCode(components.OneTimeCodeProps{
		Name: "email_code", Label: "Email code", Length: 4,
		ComponentProps: components.ComponentProps{Parts: components.Parts{
			"square": {Attrs: components.Attrs{"data-testid": "code-square"}},
			"raw":    {Attrs: components.Attrs{"data-testid": "code-raw"}},
		}},
	}))

	squares := partTags(html, "square")
	if len(squares) != 4 {
		t.Fatalf("drew %d squares, want 4:\n%s", len(squares), html)
	}
	for at, square := range squares {
		if !strings.Contains(square, `data-testid="code-square"`) {
			t.Errorf("square %d does not carry the attributes written for the square part: %s", at, square)
		}
	}

	raw := partTags(html, "raw")
	if len(raw) != 1 {
		t.Fatalf("drew %d hidden fields, want 1:\n%s", len(raw), html)
	}
	if !strings.Contains(raw[0], `data-testid="code-raw"`) {
		t.Errorf("the field that submits does not carry the attributes written for the raw part: %s", raw[0])
	}

	// Each part's attributes land on that part and on nothing else.
	if n := strings.Count(html, `data-testid="code-square"`); n != 4 {
		t.Errorf("the square attributes were written %d times, want once per square:\n%s", n, html)
	}
	if n := strings.Count(html, `data-testid="code-raw"`); n != 1 {
		t.Errorf("the raw attributes were written %d times, want once:\n%s", n, html)
	}
}
