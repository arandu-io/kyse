//go:build e2e

package e2e

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// serveFocus serves the components whose first control has to hold the cursor
// when the markup that draws it arrives by a swap, which is the one case the
// browser's own autofocus does not cover: it runs once for a document, and the
// request library moves the focus into swapped content only for what asks.
func serveFocus(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	html := func(w http.ResponseWriter, body string) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, body)
	}
	editor := func(extra components.EditInPlaceProps) string {
		extra.ID, extra.Name, extra.Label = "title", "title", "Title"
		return string(components.EditInPlace(extra))
	}
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		body := `<a href="/signin">Sign in</a> <input id="notes" aria-label="Notes"> ` +
			editor(components.EditInPlaceProps{Value: "Old", EditURL: "/title/edit"}) +
			editor2() +
			`<button type="button" ypy-get="/frag/input" ypy-target="#slot" ypy-swap="innerHTML">Load input</button>` +
			`<button type="button" ypy-get="/frag/password" ypy-target="#slot" ypy-swap="innerHTML">Load password</button>` +
			`<button type="button" ypy-get="/frag/dialog" ypy-target="#slot" ypy-swap="innerHTML">Load dialog</button>` +
			`<button type="button" ypy-get="/frag/form-dialog" ypy-target="#slot" ypy-swap="innerHTML">Load form dialog</button>` +
			`<div id="slot"></div>`
		html(w, layout("Focus", body))
	})
	mux.HandleFunc("GET /signin", func(w http.ResponseWriter, r *http.Request) {
		html(w, layout("Sign in", string(components.Input(components.InputProps{Name: "email", ID: "email", AriaLabel: "Email", Autofocus: true}))))
	})
	mux.HandleFunc("GET /title/edit", func(w http.ResponseWriter, r *http.Request) {
		html(w, editor(components.EditInPlaceProps{Value: "Old", Editing: true, SaveURL: "/title/save"}))
	})
	mux.HandleFunc("GET /title2/edit", func(w http.ResponseWriter, r *http.Request) {
		html(w, editor2Editing(""))
	})
	mux.HandleFunc("PUT /title2/save", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		html(w, editor2Editing("Too short"))
	})
	mux.HandleFunc("GET /frag/input", func(w http.ResponseWriter, r *http.Request) {
		html(w, string(components.Input(components.InputProps{Name: "name", ID: "name", AriaLabel: "Name", Autofocus: true})))
	})
	mux.HandleFunc("GET /frag/password", func(w http.ResponseWriter, r *http.Request) {
		html(w, string(components.Password(components.PasswordProps{Name: "secret", ID: "secret", Label: "Secret", Autofocus: true})))
	})
	mux.HandleFunc("GET /frag/dialog", func(w http.ResponseWriter, r *http.Request) {
		html(w, string(components.DialogTrigger(components.DialogTriggerProps{Dialog: "ask", Label: "Open ask"}))+
			string(components.Dialog(components.DialogProps{ID: "ask", Title: "Delete it?", Alert: true, Action: "/x"})))
	})
	mux.HandleFunc("GET /frag/form-dialog", func(w http.ResponseWriter, r *http.Request) {
		html(w, string(components.FormDialog(components.FormDialogProps{
			ID: "ask-form", Title: "Delete it?", Alert: true, Open: true, Action: "/x",
			Body: nil,
		})))
	})
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		w.Header().Set("Vary", "YPY-Request")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func editor2() string {
	return string(components.EditInPlace(components.EditInPlaceProps{ID: "title2", Name: "title2", Label: "Title two", Value: "Old two", EditURL: "/title2/edit"}))
}

func editor2Editing(message string) string {
	return string(components.EditInPlace(components.EditInPlaceProps{ID: "title2", Name: "title2", Label: "Title two", Value: "x", Editing: true, SaveURL: "/title2/save", Message: message}))
}

// expectFocus waits for the element the selector names to hold the focus: the
// request library moves it after the swap settles, not with it.
func expectFocus(t *testing.T, page *browser.Page, selector, what string) {
	t.Helper()
	var ok any
	page.Evaluate(`new Promise(function (resolve) {
		var until = Date.now() + 3000;
		(function poll() {
			var e = document.activeElement;
			if ((e && e.matches(`+quote(selector)+`)) || Date.now() > until) resolve(!!(e && e.matches(`+quote(selector)+`)));
			else setTimeout(poll, 25);
		})();
	})`, &ok)
	if fmt.Sprint(ok) != "true" {
		t.Errorf("the focus is not on %s; it is on %s", what, eval(t, page, `document.activeElement ? document.activeElement.outerHTML.slice(0, 160) : 'nothing'`))
	}
}

func quote(s string) string { return fmt.Sprintf("%q", s) }

// Pressing Edit swaps the field in and the cursor is in it, so the person types
// without a second click.
func TestEditInPlaceOpensWithTheCursorInTheField(t *testing.T) {
	srv := serveFocus(t)
	page := open(t, &app{Server: srv}, 1280)
	page.Visit("/")
	// The cursor is somewhere else when the answer comes, which is when the
	// browser's own autofocus does nothing: it takes the field only when the
	// page holds no cursor at all.
	page.Evaluate(`(function () {
		document.querySelector('#title [data-part=edit]').click();
		document.getElementById('notes').focus();
	})()`, nil)
	page.ExpectVisible(browser.CSS("#title input", "the field"))
	expectFocus(t, page, "#title-input", "the field of the editor")
}

// A save the server refuses draws the form again, and the cursor goes back to
// the field it refused, which the button that was pressed no longer holds.
func TestEditInPlaceRefusedSaveKeepsTheCursorInTheField(t *testing.T) {
	srv := serveFocus(t)
	page := open(t, &app{Server: srv}, 1280)
	page.Allow("/title2/save")
	page.Visit("/")
	page.Click(browser.CSS("#title2 [data-part=edit]", "the edit button"))
	page.Evaluate(`(function () {
		document.querySelector('#title2 [data-part=save]').click();
		document.getElementById('notes').focus();
	})()`, nil)
	page.ExpectText("Too short")
	expectFocus(t, page, "#title2-input", "the field of the refused editor")
}

// A field drawn with Autofocus by a fragment takes the cursor, like one drawn
// with the page.
func TestAnAutofocusFieldThatArrivesBySwapTakesTheCursor(t *testing.T) {
	srv := serveFocus(t)
	page := open(t, &app{Server: srv}, 1280)
	page.Visit("/")
	page.ClickRole("button", "Load input")
	page.ExpectVisible(browser.CSS("#slot input", "the field"))
	expectFocus(t, page, "#name", "the field that arrived")

	page.ClickRole("button", "Load password")
	page.ExpectVisible(browser.CSS("#slot input[name=secret]", "the password field"))
	expectFocus(t, page, "#slot input[name=secret]", "the password field that arrived")
}

// A page with an Autofocus field keeps the cursor in it when it is reached by a
// boosted link: the heading the navigation would focus does not take it.
func TestAnAutofocusFieldOfABoostedPageKeepsTheCursor(t *testing.T) {
	srv := serveFocus(t)
	page := open(t, &app{Server: srv}, 1280)
	page.Visit("/")
	page.ClickRole("link", "Sign in")
	page.ExpectVisible(browser.CSS("#email", "the field"))
	expectFocus(t, page, "#email", "the field of the page")
}

// A dialog that arrives by a swap puts the cursor on its cancel button when it
// is an alert, whichever way it is opened.
func TestADialogThatArrivesBySwapFocusesItsCancelButton(t *testing.T) {
	srv := serveFocus(t)
	page := open(t, &app{Server: srv}, 1280)
	page.Visit("/")
	page.ClickRole("button", "Load dialog")
	page.ClickRole("button", "Open ask")
	page.ExpectModal("Delete it?")
	expectFocus(t, page, "dialog#ask [data-part=cancel]", "the cancel button of the dialog")

	page.Visit("/")
	page.ClickRole("button", "Load form dialog")
	page.ExpectModal("Delete it?")
	expectFocus(t, page, "dialog#ask-form [data-part=cancel]", "the cancel button of the form dialog")
}
