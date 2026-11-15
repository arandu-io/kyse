//go:build e2e

package e2e

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// requestsApp serves the components that ask the server for something on
// their own -- a deferred region, a list that grows, a row that goes away, a
// search -- under the policy an application serves, and keeps the log of what
// reached it. Each answers the way the native rules of the request library
// read it: a fragment for a request that carries YPY-Request, an empty 200 for
// what is removed, and a status the library does not swap for a refusal.
type requestsApp struct {
	*httptest.Server
	mu  sync.Mutex
	log []string
}

func (a *requestsApp) saw(prefix string) int {
	a.mu.Lock()
	defer a.mu.Unlock()
	n := 0
	for _, line := range a.log {
		if strings.HasPrefix(line, prefix) {
			n++
		}
	}
	return n
}

func serveRequests(t *testing.T) *requestsApp {
	t.Helper()
	a := &requestsApp{}
	mux := http.NewServeMux()
	record := func(r *http.Request) {
		a.mu.Lock()
		defer a.mu.Unlock()
		a.log = append(a.log, fmt.Sprintf("%s %s?%s ypy=%s target=%s", r.Method, r.URL.Path, r.URL.RawQuery,
			r.Header.Get("YPY-Request"), r.Header.Get("YPY-Target")))
	}
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		body := string(components.LazyLoad(components.LazyLoadProps{URL: "/lazy", Label: "Summary", Message: "Loading"})) +
			`<ul id="feed"><li>first</li></ul>` +
			string(components.LoadMore(components.LoadMoreProps{URL: "/more?page=2"})) +
			string(components.EditInPlace(components.EditInPlaceProps{ID: "t422", Name: "title", Label: "Title 422", Value: "Old 422", EditURL: "/title/edit?as=422"})) +
			string(components.EditInPlace(components.EditInPlaceProps{ID: "t500", Name: "title", Label: "Title 500", Value: "Old 500", EditURL: "/title/edit?as=500"})) +
			`<table><tbody>` +
			`<tr id="row-ok"><td>ok</td><td>` + string(components.DeleteRow(components.DeleteRowProps{URL: "/rows/ok", Description: "Delete ok"})) + `</td></tr>` +
			`<tr id="row-204"><td>204</td><td>` + string(components.DeleteRow(components.DeleteRowProps{URL: "/rows/204", Description: "Delete 204"})) + `</td></tr>` +
			`<tr id="row-500"><td>500</td><td>` + string(components.DeleteRow(components.DeleteRowProps{URL: "/rows/500", Description: "Delete 500"})) + `</td></tr>` +
			`</tbody></table>` +
			string(components.FileUpload(components.FileUploadProps{Name: "doc", Label: "Document", UploadURL: "/upload"})) +
			`<ul id="doc-files"></ul>` +
			string(components.ActiveSearch(components.ActiveSearchProps{Name: "q", Label: "Search", URL: "/search", Delay: 50})) +
			`<div id="q-results"></div>` +
			`<div id="spacer"></div>` +
			string(components.LazyLoad(components.LazyLoadProps{URL: "/lazy-far", Label: "Far", Message: "Loading far", Deferred: true}))
		fmt.Fprint(w, layout("Requests", body))
	})
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		fmt.Fprint(w, `#spacer{height:2400px}`)
	})
	mux.HandleFunc("GET /lazy", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		fmt.Fprint(w, `<section id="summary">The summary</section>`)
	})
	mux.HandleFunc("GET /lazy-far", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		fmt.Fprint(w, `<section id="far">The far panel</section>`)
	})
	mux.HandleFunc("GET /more", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		fmt.Fprint(w, `<li>second</li>`)
	})
	mux.HandleFunc("DELETE /rows/{id}", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		switch r.PathValue("id") {
		case "204":
			w.WriteHeader(http.StatusNoContent)
		case "500":
			http.Error(w, "boom", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /title/edit", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		kind := r.URL.Query().Get("as")
		fmt.Fprint(w, components.EditInPlace(components.EditInPlaceProps{ID: "t" + kind, Name: "title", Label: "Title " + kind, Value: "Old " + kind, Editing: true, SaveURL: "/title/" + kind}))
	})
	mux.HandleFunc("PUT /title/{kind}", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		switch r.PathValue("kind") {
		case "422":
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, components.EditInPlace(components.EditInPlaceProps{ID: "t422", Name: "title", Label: "Title 422", Value: "x", Editing: true, SaveURL: "/title/422", Message: "Too short"}))
		case "500":
			http.Error(w, "boom", http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("POST /upload", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		if err := r.ParseMultipartForm(1 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, header, err := r.FormFile("doc")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, "<li>uploaded "+template.HTMLEscapeString(header.Filename)+"</li>")
	})
	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		record(r)
		fmt.Fprint(w, template.HTMLEscapeString("results for "+r.URL.Query().Get("q")))
	})
	a.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		w.Header().Set("Vary", "YPY-Request")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(a.Close)
	return a
}

// A region that fetches itself once the page has drawn is replaced by what the
// server answers, and fetches once.
func TestALazyRegionIsReplacedOnceByItsAnswer(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Visit("/")
	page.ExpectText("The summary")
	if n := app.saw("GET /lazy?"); n != 1 {
		t.Errorf("the region asked %d times, want once", n)
	}
	if got := eval(t, page, `document.querySelectorAll('[aria-busy=true]').length`); got != "1" {
		t.Errorf("%s regions are still busy, want only the one that has not been scrolled to", got)
	}
}

// A region that waits to be revealed does not fetch until it is scrolled to.
func TestADeferredRegionFetchesWhenItIsScrolledTo(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Visit("/")
	page.Evaluate(`window.scrollTo(0, 0)`, nil)
	page.ExpectText("The summary")
	if n := app.saw("GET /lazy-far"); n != 0 {
		t.Fatalf("the far region asked %d times before it was in view", n)
	}
	page.Evaluate(`document.querySelector('[aria-label=Far]').scrollIntoView()`, nil)
	page.ExpectText("The far panel")
}

// The control that grows a list swaps itself for the next rows and the next
// control.
func TestLoadMoreAppendsTheRowsAndLeavesNoControlBehind(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Visit("/")
	page.ClickRole("button", "Load more")
	page.ExpectText("second")
	if got := eval(t, page, `document.querySelectorAll('button[ypy-get="/more?page=2"]').length`); got != "0" {
		t.Errorf("the control that asked is still on the page: %s", got)
	}
}

// What is deleted leaves the page when the server answers with nothing; a
// refusal and a 204 leave it where it was, which is the library's native rule
// for both.
func TestADeletedRowLeavesOnlyWhenTheServerAnswersWithNothing(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Allow("/rows/500")
	page.Visit("/")
	for _, name := range []string{"Delete 500", "Delete 204", "Delete ok"} {
		page.ClickRole("button", name)
		page.Click(browser.CSS("dialog[data-kyse-confirm] [data-part=confirm]", "the button of the question"))
	}
	// The row goes after the 200 ms the swap waits, so it is waited for.
	page.Evaluate(`new Promise(function (resolve) {
		var until = Date.now() + 3000;
		(function poll() { if (!document.getElementById('row-ok') || Date.now() > until) resolve(true); else setTimeout(poll, 50); })();
	})`, nil)
	has := func(id string) string { return eval(t, page, `String(!!document.getElementById('`+id+`'))`) }
	if has("row-ok") != "false" {
		t.Error("the row answered with an empty 200 is still on the page")
	}
	if has("row-500") != "true" {
		t.Error("the row whose delete was refused with a 500 left the page")
	}
	if has("row-204") != "true" {
		t.Error("the row answered with a 204 left the page: 204 is the status that swaps nothing")
	}
}

// Typing asks once the pause has passed, and the answer goes into the region
// the search names.
func TestTheActiveSearchAnswersIntoItsRegion(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Visit("/")
	page.Fill("Search", "ada")
	page.ExpectText("results for ada")
}

// serverOf is the browser's way to the app.
func serverOf(a *requestsApp) *app { return &app{Server: a.Server} }

// A save the server refuses as invalid draws the form again with its message;
// a save that fails on the server leaves the form as it was and draws no error
// page into it.
func TestASaveRefusedAsInvalidShowsItsMessageAndAFailedOneChangesNothing(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Allow("/title/422")
	page.Allow("/title/500")
	page.Visit("/")

	page.Click(browser.CSS("#t422 [data-part=edit]", "the edit button of the 422 editor"))
	page.Click(browser.CSS("#t422 [data-part=save]", "the save button"))
	page.ExpectText("Too short")
	if got := eval(t, page, `String(!!document.querySelector('form#t422 input'))`); got != "true" {
		t.Errorf("the form with the message is not on the page: %s", eval(t, page, `document.getElementById('t422').outerHTML`))
	}

	page.Click(browser.CSS("#t500 [data-part=edit]", "the edit button of the 500 editor"))
	page.Click(browser.CSS("#t500 [data-part=save]", "the save button"))
	page.Evaluate(`new Promise(function (resolve) { setTimeout(resolve, 600) })`, nil)
	if got := eval(t, page, `String(!!document.querySelector('form#t500 input'))`); got != "true" {
		t.Errorf("the form was replaced by the answer to a failed save: %s", eval(t, page, `document.body.innerHTML`))
	}
	page.ExpectNoText("boom")
}

// A file chosen is sent as the multipart body the server reads.
func TestAChosenFileIsSentAsMultipart(t *testing.T) {
	app := serveRequests(t)
	page := open(t, serverOf(app), 1280)
	page.Visit("/")
	page.Evaluate(`(function () {
		var input = document.querySelector('input[type=file]');
		var data = new DataTransfer();
		data.items.add(new File(['hello'], 'hello.txt', { type: 'text/plain' }));
		input.files = data.files;
		input.dispatchEvent(new Event('change', { bubbles: true }));
	})()`, nil)
	page.ExpectText("uploaded hello.txt")
}
