package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// An action that worked says so with Context.With, redirects, and the page it
// lands on carries the sentence in its flash. The toast reads it from there, so
// a layout draws the confirmation with one call on every page.

var _ components.FlashPage = view.Page{}

func TestAToastDrawsWhatThePageFlashed(t *testing.T) {
	page := view.Page{Messages: map[string]string{"status": "Invoice saved."}}
	html := string(components.Toast(components.ToastProps{Page: page}))
	if !strings.Contains(html, ">Invoice saved.</h2>") {
		t.Errorf("the flashed status was not drawn:\n%s", html)
	}
	if !strings.Contains(html, `role="status"`) {
		t.Errorf("the flashed status is not announced:\n%s", html)
	}
}

func TestAToastOnAPageWithNothingFlashedDrawsNothing(t *testing.T) {
	for name, page := range map[string]components.Page{
		"nil page":         nil,
		"field-only page":  page{},
		"nothing flashed":  view.Page{},
		"another key only": view.Page{Messages: map[string]string{"resent": "Sent again."}},
		"empty status":     view.Page{Messages: map[string]string{"status": ""}},
	} {
		if html := strings.TrimSpace(string(components.Toast(components.ToastProps{Page: page}))); html != "" {
			t.Errorf("%s: drew a toast with nothing to say:\n%s", name, html)
		}
	}
}

func TestAToastReadsTheKeyItIsGiven(t *testing.T) {
	page := view.Page{Messages: map[string]string{"status": "Saved.", "resent": "Sent again."}}
	html := string(components.Toast(components.ToastProps{Page: page, FlashKey: "resent"}))
	if !strings.Contains(html, ">Sent again.</h2>") || strings.Contains(html, "Saved.") {
		t.Errorf("FlashKey did not choose the sentence:\n%s", html)
	}
}

func TestATitleWrittenOnTheToastWinsOverTheFlash(t *testing.T) {
	page := view.Page{Messages: map[string]string{"status": "Saved."}}
	html := string(components.Toast(components.ToastProps{Page: page, Title: "Archived."}))
	if !strings.Contains(html, ">Archived.</h2>") || strings.Contains(html, "Saved.") {
		t.Errorf("the flash replaced the title the caller wrote:\n%s", html)
	}
}

func TestAFlashedSentenceIsEscaped(t *testing.T) {
	page := view.Page{Messages: map[string]string{"status": `<script>alert(1)</script>`}}
	html := string(components.Toast(components.ToastProps{Page: page}))
	if strings.Contains(html, "<script>") || !strings.Contains(html, "&lt;script&gt;") {
		t.Errorf("a flashed sentence was written as markup:\n%s", html)
	}
}

// The toast's doc used to call itself the whole flash mechanism, with "no
// redirect carrying a string" -- the opposite of how a form that worked is
// answered. It names Context.With now, and the sentence that contradicted it
// stays out.
func TestTheToastDocDescribesTheFlash(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "components", "toast.kyse.go"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(body)
	if strings.Contains(doc, "no redirect carrying a string") {
		t.Error("the toast doc still says a redirect carries no sentence")
	}
	if !strings.Contains(doc, "Context.With") {
		t.Error("the toast doc does not name Context.With, where the sentence comes from")
	}
}
