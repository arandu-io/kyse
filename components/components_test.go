package components_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// The component library shipped eleven components and proved nothing.
//
// What is checked here is what a component library can get wrong without
// anybody noticing until it is in somebody's page: markup that escapes nothing,
// an attribute written empty, and a generated file that no longer matches the
// source it claims to come from.

// TestEveryComponentRenders is the smoke test, and it is a table so that adding
// a component without adding a line here is visible.
func TestEveryComponentRenders(t *testing.T) {
	for _, c := range []struct {
		name string
		html string
	}{
		{"Alert", string(components.Alert(components.AlertProps{Title: "Saved"}))},
		{"Avatar", string(components.Avatar(components.AvatarProps{Name: "Ada Lovelace"}))},
		{"Badge", string(components.Badge(components.BadgeProps{Label: "draft"}))},
		{"Button", string(components.Button(components.ButtonProps{Label: "Save"}))},
		{"Card", string(components.Card(components.CardProps{Title: "A post"}))},
		{"Dialog", string(components.Dialog(components.DialogProps{ID: "d", Title: "Sure?"}))},
		{"Empty", string(components.Empty(components.EmptyProps{Title: "Nothing here"}))},
		{"Field", string(components.Field(components.FieldProps{Name: "email", Label: "Email"}))},
		{"Textarea", string(components.Textarea(components.TextareaProps{Name: "body", Label: "Body"}))},
		{"ThemeToggle", string(components.ThemeToggle())},
		{"Toast", string(components.Toast(components.ToastProps{Title: "Saved"}))},
	} {
		t.Run(c.name, func(t *testing.T) {
			if strings.TrimSpace(c.html) == "" {
				t.Fatal("rendered nothing")
			}
			if !strings.Contains(c.html, "<") {
				t.Fatalf("rendered no markup: %q", c.html)
			}
			// Every component carries a Basecoat class, because that is where
			// its appearance comes from. One that renders bare markup is one
			// that will look unstyled in a project and correct in review.
			if !strings.Contains(c.html, `class="`) {
				t.Errorf("rendered no class attribute:\n%s", c.html)
			}
		})
	}
}

// TestAPropIsEscaped is the one that matters.
//
// A component is interpolated with {!! !!}, which does not escape -- and that is
// correct only because everything inside the component was escaped by the view
// compiler on the way in. If that ever stops being true, every component becomes
// an injection point, and it becomes untrue silently: the markup still renders.
func TestAPropIsEscaped(t *testing.T) {
	const payload = `<script>alert(1)</script>`

	for name, html := range map[string]string{
		"Button": string(components.Button(components.ButtonProps{Label: payload})),
		"Badge":  string(components.Badge(components.BadgeProps{Label: payload})),
		"Alert":  string(components.Alert(components.AlertProps{Title: payload})),
		"Card":   string(components.Card(components.CardProps{Title: payload})),
		"Field":  string(components.Field(components.FieldProps{Name: "x", Label: payload})),
		"Empty":  string(components.Empty(components.EmptyProps{Title: payload})),
		"Toast":  string(components.Toast(components.ToastProps{Title: payload})),
	} {
		if strings.Contains(html, payload) {
			t.Errorf("%s wrote a script tag through unescaped:\n%s", name, html)
		}
		if !strings.Contains(html, "&lt;script&gt;") {
			t.Errorf("%s did not escape the payload at all:\n%s", name, html)
		}
	}
}

// TestAnEmptyOptionalAttributeIsNotWritten.
//
// `hx-post=""` is not the absence of hx-post: HTMX acts on it and posts to the
// current URL. The same for data-variant, which the stylesheet matches on.
func TestAnEmptyOptionalAttributeIsNotWritten(t *testing.T) {
	plain := string(components.Button(components.ButtonProps{Label: "Save"}))

	for _, attr := range []string{"data-variant", "data-size", "hx-post", "hx-get", "hx-target", "disabled"} {
		if strings.Contains(plain, attr) {
			t.Errorf("a button with no %s wrote it anyway:\n%s", attr, plain)
		}
	}

	// And when they are set, they are there.
	full := string(components.Button(components.ButtonProps{
		Label: "Save", Variant: "destructive", Size: "sm", HxPost: "/x", Disabled: true,
	}))
	for _, want := range []string{`data-variant="destructive"`, `data-size="sm"`, `hx-post="/x"`, "disabled"} {
		if !strings.Contains(full, want) {
			t.Errorf("%s is missing:\n%s", want, full)
		}
	}
}

// TestAToastWithNothingToSayDrawsNothing: a page includes it unconditionally, so
// the empty case has to be empty rather than an empty box.
func TestAToastWithNothingToSayDrawsNothing(t *testing.T) {
	if got := strings.TrimSpace(string(components.Toast(components.ToastProps{}))); got != "" {
		t.Errorf("a toast with no title drew %q", got)
	}
}

// TestInitialsWalkRunes: a name starting with É is one letter to a person and
// two bytes to Go, and slicing bytes prints half a character.
func TestInitialsWalkRunes(t *testing.T) {
	for name, want := range map[string]string{
		"Ada Lovelace":    "AL",
		"Émile Borel":     "ÉB",
		"Prince":          "P",
		"Ada B. Lovelace": "AL",
	} {
		if got := (components.AvatarProps{Name: name}).Initials(); got != want {
			t.Errorf("Initials(%q) = %q, want %q", name, got, want)
		}
	}
}

// TestTheGeneratedFilesMatchTheirSources.
//
// The .go files are committed, because a module whose generated files are
// missing is one `go get` cannot use. That makes a stale one shippable: the
// source says one thing, the published package does another, and nothing in the
// build notices.
//
// This does not recompile -- that would need the view compiler, which is in aru
// and would make this module depend on the CLI. It checks the weaker property
// that catches the real mistake: every source has an output, every output has a
// source, and the output names the source it came from.
//
// A missing output is caught earlier and harder: the package stops compiling,
// because the function it declared is gone. What only this test catches is the
// other direction -- a .go left behind by a component that was renamed, which
// compiles perfectly and ships a component nobody maintains.
func TestTheGeneratedFilesMatchTheirSources(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}

	sources, generated := map[string]bool{}, map[string]bool{}
	for _, e := range entries {
		switch name := e.Name(); {
		case strings.HasSuffix(name, ".kyse.go"):
			sources[strings.TrimSuffix(name, ".kyse.go")] = true
		case name == "doc.go" || name == "components_test.go":
		case strings.HasSuffix(name, ".go"):
			generated[strings.TrimSuffix(name, ".go")] = true
		}
	}

	if len(sources) == 0 {
		t.Fatal("no .kyse.go sources: this test is proving nothing")
	}

	for name := range sources {
		if !generated[name] {
			t.Errorf("%s.kyse.go has no compiled %s.go: run `aru view:build` and commit it", name, name)
			continue
		}
		body, err := os.ReadFile(filepath.Join(".", name+".go"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), name+".kyse.go") {
			t.Errorf("%s.go does not name the source it came from", name)
		}
	}
	for name := range generated {
		if !sources[name] {
			t.Errorf("%s.go has no source: it is left over from a component that was renamed or removed", name)
		}
	}
}

// page is a stand-in for framework/view.Page, so this module can be tested
// without importing the framework it is deliberately independent of.
//
// It answers the two questions components.Page asks, in the same shape view.Page
// does: the first message for a name, and what was typed falling back to what
// the caller passed.
type page struct {
	errs map[string]string
	old  map[string]string
}

func (p page) FieldError(name string) string { return p.errs[name] }

func (p page) OldOr(name, fallback string) string {
	if v, typed := p.old[name]; typed {
		return v
	}
	return fallback
}

// TestAFieldAsksThePageByItsOwnName is the property the Page prop exists for.
//
// The name is written once. A field that took its message as a separate prop
// took the name three times in one call, and nothing checked that the three
// agreed -- so a box labelled "Confirm password" could draw the message of
// "password", or none, and look correct in review.
func TestAFieldAsksThePageByItsOwnName(t *testing.T) {
	state := page{
		errs: map[string]string{"password": "must be at least 12 characters"},
		old:  map[string]string{"email": "ada@example.com"},
	}

	password := string(components.Field(components.FieldProps{
		Name: "password", Label: "Password", Type: "password", Page: state,
	}))
	if !strings.Contains(password, "must be at least 12 characters") {
		t.Errorf("the message for this field was not drawn:\n%s", password)
	}
	if !strings.Contains(password, `aria-invalid="true"`) {
		t.Errorf("a field with a message is not marked invalid:\n%s", password)
	}
	if !strings.Contains(password, `aria-describedby="password-error"`) {
		t.Errorf("the message is shown and not announced:\n%s", password)
	}

	// The other field of the same form was accepted, and asks the same page.
	email := string(components.Field(components.FieldProps{
		Name: "email", Label: "Email", Type: "email", Page: state,
	}))
	if strings.Contains(email, "must be at least") {
		t.Errorf("a field drew another field's message:\n%s", email)
	}
	if !strings.Contains(email, `value="ada@example.com"`) {
		t.Errorf("what was typed did not come back:\n%s", email)
	}
}

// TestAPasswordComesBackEmptyAndSaysWhy.
//
// The flash carries the message for a password and never the value, so the box
// is blank and the sentence under it is not. A blank box with nothing under it
// is the failure this whole path was built to end.
func TestAPasswordComesBackEmptyAndSaysWhy(t *testing.T) {
	html := string(components.Field(components.FieldProps{
		Name: "password", Label: "Password", Type: "password",
		Page: page{errs: map[string]string{"password": "must be at least 12 characters"}},
	}))

	if !strings.Contains(html, `value=""`) {
		t.Errorf("the password box came back carrying a value:\n%s", html)
	}
	if !strings.Contains(html, "must be at least 12 characters") {
		t.Errorf("the empty box does not say why it was rejected:\n%s", html)
	}
}

// TestWhatWasTypedBeatsTheStoredValue is the edit form: the box starts at the
// row and comes back carrying the change somebody was in the middle of making.
func TestWhatWasTypedBeatsTheStoredValue(t *testing.T) {
	stored := components.FieldProps{Name: "title", Label: "Title", Value: "Stored title"}

	if html := string(components.Field(stored)); !strings.Contains(html, `value="Stored title"`) {
		t.Errorf("a field with no page behind it lost its value:\n%s", html)
	}

	stored.Page = page{old: map[string]string{"title": "What I was typing"}}
	if html := string(components.Field(stored)); !strings.Contains(html, `value="What I was typing"`) {
		t.Errorf("the rejected edit reverted to the stored row:\n%s", html)
	}
}

// TestATextareaAsksThePageToo: the long-text half of a form is the same
// question, and a component that answered it differently would be the second way
// to draw a message.
func TestATextareaAsksThePageToo(t *testing.T) {
	html := string(components.Textarea(components.TextareaProps{
		Name: "body", Label: "Body",
		Page: page{
			errs: map[string]string{"body": "is required"},
			old:  map[string]string{"body": "half a paragraph"},
		},
	}))

	if !strings.Contains(html, "is required") {
		t.Errorf("the message was not drawn:\n%s", html)
	}
	if !strings.Contains(html, "half a paragraph") {
		t.Errorf("what was typed did not come back:\n%s", html)
	}
}
