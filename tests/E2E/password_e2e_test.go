//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/validation"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// The password box is drawn with a policy that asks for every kind of character
// the validator can name, so each line of the checklist has a rule of its own
// to tick for.
func passwordBox() string {
	return string(components.Password(components.PasswordProps{
		Name: "password", Label: "Password", Hint: "Use 12 to 128 characters.",
		Policy: validation.PasswordMin(12).Max(128).Letters().MixedCase().Numbers().Symbols(),
	}))
}

// passwordApp serves the box under the policy an application serves, and the
// same box as a fragment an htmx swap brings in. The stylesheet is the one the
// other tests measure: it is built from the component sources and carries no
// Basecoat file, so the meter has no rule of the application to lean on.
func passwordApp(t *testing.T) *app {
	t.Helper()
	css, err := stylesheet()
	if err != nil {
		t.Skipf("the stylesheet cannot be compiled here: %v", err)
	}
	a := &app{}
	mux := http.NewServeMux()
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(css)
	})
	mux.HandleFunc("GET /password", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, layout("Password", `<form class="flex max-w-sm flex-col gap-4">`+passwordBox()+`</form>`))
	})
	mux.HandleFunc("GET /password/swap", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, layout("Swap", `<button type="button" class="btn" hx-get="/password/fragment" hx-target="#slot" hx-swap="innerHTML">Choose a password</button><form id="slot" class="flex max-w-sm flex-col gap-4"></form>`))
	})
	mux.HandleFunc("GET /password/fragment", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, passwordBox())
	})
	a.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(a.Close)
	return a
}

// checklist is what a person sees of the box at one moment: which line is
// ticked, how wide the bar is drawn against its track, and what is typed.
type checklist struct {
	Typed int             `json:"typed"`
	Met   map[string]bool `json:"met"`
	Fill  float64         `json:"fill"`
	Track float64         `json:"track"`
	Open  bool            `json:"open"`
}

func readChecklist(t *testing.T, page *browser.Page) checklist {
	t.Helper()
	var out checklist
	raw := eval(t, page, `JSON.stringify((function () {
		var root = document.querySelector('[data-kyse-behavior="password"]');
		var met = {};
		root.querySelectorAll('[data-requirement]').forEach(function (line) {
			met[line.getAttribute('data-requirement')] = line.getAttribute('data-met') === 'true';
		});
		var fill = root.querySelector('[data-part="fill"]').getBoundingClientRect();
		var track = root.querySelector('[data-part="meter"]').getBoundingClientRect();
		return {
			typed: Array.from(root.querySelector('input').value).length,
			met: met, fill: fill.width, track: track.width,
			open: !root.querySelector('[data-part="panel"]').hidden,
		};
	})())`)
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		t.Fatalf("the checklist does not read: %v\n%s", err, raw)
	}
	return out
}

var (
	lowerCase = regexp.MustCompile(`\p{Ll}`)
	upperCase = regexp.MustCompile(`\p{Lu}`)
)

// judged is what the policy says about a value, worked out here and not read
// from the page: the browser is asked, and this is the answer it is held to.
func judged(value string) map[string]bool {
	lower, upper := lowerCase.MatchString(value), upperCase.MatchString(value)
	return map[string]bool{
		"min":       utf8.RuneCountInString(value) >= 12,
		"max":       utf8.RuneCountInString(value) <= 128,
		"letters":   regexp.MustCompile(`\p{L}`).MatchString(value),
		"mixedCase": lower && upper,
		"numbers":   regexp.MustCompile(`\p{N}`).MatchString(value),
		"symbols":   regexp.MustCompile(`[\p{P}\p{S}\p{Z}]`).MatchString(value),
	}
}

// requireChecklist holds the page to what judged says for the value typed:
// every line ticked or not as the policy asks, and the bar as wide as the
// share of the lines ticked.
func requireChecklist(t *testing.T, page *browser.Page, value, when string) {
	t.Helper()
	got := readChecklist(t, page)
	want := judged(value)
	counted, ticked := 0, 0
	for key, ok := range want {
		if key == "max" {
			continue
		}
		counted++
		if value != "" && ok {
			ticked++
		}
		if value == "" {
			ok = false
		}
		if got.Met[key] != ok {
			t.Errorf("%s: the line %q is ticked=%v for %q, and the policy says %v", when, key, got.Met[key], value, ok)
		}
	}
	if got.Typed != utf8.RuneCountInString(value) {
		t.Fatalf("%s: %d characters are in the box, %d were typed", when, got.Typed, utf8.RuneCountInString(value))
	}
	if got.Track <= 0 {
		t.Fatalf("%s: the bar has no track to fill", when)
	}
	if share := got.Track * float64(ticked) / float64(counted); math.Abs(got.Fill-share) > 1.5 {
		t.Errorf("%s: the bar is %.1fpx of a %.1fpx track for %d of %d lines ticked, and it should be %.1fpx",
			when, got.Fill, got.Track, ticked, counted, share)
	}
}

// Each line ticks as its rule is met, one key at a time, and takes the tick
// back when the key that met it is deleted; the bar follows the lines, so what
// a person sees move is what the lines say.
func TestEachRuleTicksAsThePasswordIsTyped(t *testing.T) {
	const typed = "Abcdef1!ghij"
	app := passwordApp(t)
	page := open(t, app, 1280)
	page.Visit("/password")
	page.Click(browser.Label("Password"))
	requireChecklist(t, page, "", "before a key")
	if !readChecklist(t, page).Open {
		t.Error("the panel is shut with the cursor in the box")
	}

	for i, r := range typed {
		page.Press(string(r))
		requireChecklist(t, page, typed[:i+1], fmt.Sprintf("after typing %q", typed[:i+1]))
	}
	shot(t, page, "password-typed")
	for i := len(typed) - 1; i >= 0; i-- {
		page.Press("Backspace")
		requireChecklist(t, page, typed[:i], fmt.Sprintf("after deleting to %q", typed[:i]))
	}
	shot(t, page, "password-deleted")
}

// A value that arrives as one input, a paste or a password manager, ticks the
// same lines the keys would have, and clearing it takes them all back.
func TestAPastedPasswordTicksEveryRuleItMeets(t *testing.T) {
	app := passwordApp(t)
	page := open(t, app, 390)
	page.Visit("/password")
	page.Fill("Password", "Sunrise-2026-Over")
	requireChecklist(t, page, "Sunrise-2026-Over", "after a paste")
	page.Fill("Password", "")
	requireChecklist(t, page, "", "after clearing the box")
}

// A box an htmx swap brings in is mounted by the swap and ticks as one drawn
// with the page.
func TestAPasswordBoxBroughtByASwapTicksToo(t *testing.T) {
	app := passwordApp(t)
	page := open(t, app, 1280)
	page.Visit("/password/swap")
	page.ClickText("Choose a password")
	page.Click(browser.Label("Password"))
	for i, r := range "Zz9?" {
		page.Press(string(r))
		requireChecklist(t, page, "Zz9?"[:i+1], fmt.Sprintf("after typing %d characters into the swapped box", i+1))
	}
	if !strings.Contains(eval(t, page, `document.querySelector('[data-kyse-behavior="password"]').getAttribute('data-kyse-mounted')`), "true") {
		t.Error("the swapped box was never mounted")
	}
}
