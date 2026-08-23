package unit

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// The component library shipped eleven components and proved nothing.
//
// What is checked here is what a component library can get wrong without
// anybody noticing until it is in somebody's page: markup that escapes nothing
// and an attribute written empty.
//
// One property does not fit here and lives in tests/Feature instead: that every
// .kyse.go source has a compiled .go beside it. That one reads the shipped tree
// rather than a rendered string.

// TestEveryComponentRenders is the smoke test, and it is a table.
//
// The table said of itself that adding a component without adding a line here
// was visible, and it was not: the table is hand-kept, and a missing line looks
// exactly like a component that does not exist. Twenty-one were added at once
// and it stayed green. TestEveryComponentIsInTheTable below is what makes the
// claim true -- it reads the directory rather than this list.
func TestEveryComponentRenders(t *testing.T) {
	for _, c := range []struct {
		name string
		html string
	}{
		{"Accordion", string(components.Accordion(components.AccordionProps{
			Items: []components.AccordionItem{{Label: "What is it", Content: "A framework."}},
		}))},
		{"Alert", string(components.Alert(components.AlertProps{Title: "Saved"}))},
		{"Avatar", string(components.Avatar(components.AvatarProps{Name: "Ada Lovelace"}))},
		{"Badge", string(components.Badge(components.BadgeProps{Label: "draft"}))},
		{"Button", string(components.Button(components.ButtonProps{Label: "Save"}))},
		{"ButtonGroup", string(components.ButtonGroup(components.ButtonGroupProps{
			Label:   "Message actions",
			Buttons: []components.ButtonProps{{Label: "Archive"}, {Label: "Report"}},
		}))},
		{"Breadcrumb", string(components.Breadcrumb(components.BreadcrumbProps{
			Items: []components.Crumb{{Label: "Home", URL: "/"}, {Label: "Invoices"}},
		}))},
		{"Card", string(components.Card(components.CardProps{Title: "A post"}))},
		{"Checkbox", string(components.Checkbox(components.CheckboxProps{Name: "terms", Label: "I agree"}))},
		{"Collapsible", string(components.Collapsible(components.CollapsibleProps{
			Label: "Details", Content: "The rest of it.",
		}))},
		{"Combobox", string(components.Combobox(components.ComboboxProps{
			Name: "framework", Label: "Framework", SearchURL: "/frameworks",
		}))},
		{"Command", string(components.Command(components.CommandProps{
			ID: "palette", Label: "Command palette",
			Groups: []components.CommandGroup{{Heading: "Go to", Items: []components.CommandItem{
				{Label: "Profile", URL: "/profile"},
			}}},
		}))},
		{"Dialog", string(components.Dialog(components.DialogProps{ID: "d", Title: "Sure?"}))},
		{"Drawer", string(components.Drawer(components.DrawerProps{ID: "nav", Title: "Menu"}))},
		{"DropdownMenu", string(components.DropdownMenu(components.DropdownMenuProps{
			Label: "Actions", Items: []components.MenuItem{{Label: "Rename"}},
		}))},
		{"Empty", string(components.Empty(components.EmptyProps{Title: "Nothing here"}))},
		{"Field", string(components.Field(components.FieldProps{Name: "email", Label: "Email"}))},
		{"Input", string(components.Input(components.InputProps{Name: "q", AriaLabel: "Search"}))},
		{"InputGroup", string(components.InputGroup(components.InputGroupProps{
			Name: "amount", Label: "Amount", Start: "USD",
		}))},
		{"Item", string(components.Item(components.ItemProps{Title: "A row", Description: "with a line under it"}))},
		{"Kbd", string(components.Kbd(components.KbdProps{Keys: []string{"Ctrl", "K"}}))},
		{"Label", string(components.Label(components.LabelProps{For: "email", Text: "Email"}))},
		{"Popover", string(components.Popover(components.PopoverProps{
			ID: "p", Label: "More", Content: "Something else.",
		}))},
		{"Progress", string(components.Progress(components.ProgressProps{Label: "Upload", Value: 40}))},
		{"RadioGroup", string(components.RadioGroup(components.RadioGroupProps{
			Name: "plan", Label: "Plan", Options: []components.RadioOption{{Label: "Free", Value: "free"}},
		}))},
		{"RangeSlider", string(components.RangeSlider(components.RangeSliderProps{
			Name: "volume", Label: "Volume", Value: 30, ShowValue: true,
		}))},
		{"Select", string(components.Select(components.SelectProps{
			Name: "plan", Label: "Plan", Options: []components.SelectOption{{Label: "Free", Value: "free"}},
		}))},
		{"Separator", string(components.Separator(components.SeparatorProps{}))},
		{"Sidebar", string(components.Sidebar(components.SidebarProps{
			Groups: []components.SidebarGroup{{Label: "Main", Items: []components.SidebarItem{{Label: "Home", URL: "/"}}}},
		}))},
		{"Skeleton", string(components.Skeleton(components.SkeletonProps{Shape: "line"}))},
		{"StatCard", string(components.StatCard(components.StatCardProps{Title: "Connections"}))},
		{"Switch", string(components.Switch(components.SwitchProps{Name: "notify", Label: "Email me"}))},
		{"Table", string(components.Table(components.TableProps{
			Columns: []components.TableColumn{{Label: "Name"}},
			Rows:    []components.TableRow{{Cells: []components.TableCell{{Text: "Ada"}}}},
		}))},
		{"Tabs", string(components.Tabs(components.TabsProps{
			Tabs: []components.Tab{{ID: "one", Label: "One", Panel: "First."}},
		}))},
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
		"Combobox": string(components.Combobox(components.ComboboxProps{
			Name: "x", Label: payload, SearchURL: "/x",
			Options: []components.ComboboxOption{{Value: payload}},
		})),
		"Command": string(components.Command(components.CommandProps{
			ID: "x", Groups: []components.CommandGroup{{Heading: payload, Items: []components.CommandItem{
				{Label: payload, URL: "/x"},
			}}},
		})),
		"InputGroup": string(components.InputGroup(components.InputGroupProps{
			Name: "x", Label: payload, Start: payload,
		})),
		"RangeSlider": string(components.RangeSlider(components.RangeSliderProps{
			Name: "x", Label: payload,
		})),
		"ButtonGroup": string(components.ButtonGroup(components.ButtonGroupProps{
			Label: payload, Buttons: []components.ButtonProps{{Label: payload}},
		})),
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

// clientDirective matches the two shapes a client-side expression is written in:
// a directive named x-something, and the @ shorthand for the event form.
//
// The character before the name has to be one that does not continue a name, so
// hx-get and hx-on:click are not matched -- those are HTMX, they are read by a
// library that parses attributes rather than evaluating them, and they stay.
// Requiring an `=` after the @ form is what separates it from the view
// compiler's own @if, @for and @go, which take parentheses or nothing.
var clientDirective = regexp.MustCompile(`(?:^|[^\w:@-])(x-[a-z][\w.:-]*|@[\w.:-]+\s*=)`)

// TestNoComponentEvaluatesAnExpression is the guard, and it is the point of the
// rewrite rather than a tidying rule.
//
// A directive's value is a string the client library compiles with the Function
// constructor. This framework serves script-src 'self' with no unsafe-eval, so
// that compilation throws at the point of evaluation: a component carrying one
// is not a noisy component, it is a dead one, and it fails in the browser
// console rather than in any build.
//
// Nothing else in the toolchain would catch it. The markup is valid, the view
// compiles, every test above still renders a string with a class in it, and the
// component looks correct in review. So the source is read here.
//
// Both the sources and what they compile to are read. The source is where a
// directive would be written; the compiled file is what a browser is actually
// sent, and a stale one ships a directive that the source no longer has.
func TestNoComponentEvaluatesAnExpression(t *testing.T) {
	for _, pattern := range []string{"*.kyse.go", "*.go"} {
		files, err := filepath.Glob(filepath.Join("..", "..", "components", pattern))
		if err != nil {
			t.Fatalf("reading the component directory: %v", err)
		}
		if len(files) == 0 {
			// A glob that matches nothing would make this test pass by finding
			// nothing to read, which is the shape of failure it exists to prevent.
			t.Fatalf("no components/%s found; this test is checking nothing", pattern)
		}

		for _, file := range files {
			// *.go matches *.kyse.go too, and a line reported twice reads as two
			// directives.
			if pattern == "*.go" && strings.HasSuffix(file, ".kyse.go") {
				continue
			}
			body, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("reading %s: %v", file, err)
			}
			for at, line := range strings.Split(string(body), "\n") {
				found := clientDirective.FindStringSubmatch(asMarkup(line))
				if found == nil {
					continue
				}
				t.Errorf("%s:%d writes %q, which is an expression the browser is asked to compile:\n\t%s\n"+
					"\tBehaviour is delegation on data-* attributes. Put the value in a data-* attribute and let the script read it.",
					filepath.Base(file), at+1, strings.TrimSpace(found[1]), strings.TrimSpace(line))
			}
		}
	}
}

// asMarkup turns a line of a compiled view back into the markup it writes, far
// enough for the pattern above to read it.
//
// A compiled view holds its markup inside Go string literals, so a tab before an
// attribute arrives as the two characters `\` and `t`. The pattern asks what
// comes before a name, and a literal `t` there is a letter -- which is how the
// first version of this test read a compiled `\t\tx-on:input` as part of a
// longer word and passed over it. A source line has none of these sequences and
// goes through unchanged.
func asMarkup(line string) string {
	return strings.NewReplacer(`\t`, "\t", `\n`, "\n", `\"`, `"`).Replace(line)
}

// TestEveryComponentIsInTheTable reads the directory instead of the list.
//
// The table above is hand-kept, and a hand-kept list has one failure mode: a
// missing line looks exactly like a component that does not exist. That is not
// hypothetical here. The library went from twelve components to thirty-seven in
// one change and the table stayed green at sixteen, because nothing compared it
// to anything.
//
// So this compares. Every components/*.kyse.go is a component that ships, and
// the file name is the component name with the dashes taken out. A source with
// no row is a component nobody renders in a test, and the failure names it.
//
// It reads the source tree rather than the package because Go cannot enumerate
// a package's functions at run time without a registry, and a registry that
// exists only for a test is a second list to forget to update.
func TestEveryComponentIsInTheTable(t *testing.T) {
	sources, err := filepath.Glob(filepath.Join("..", "..", "components", "*.kyse.go"))
	if err != nil {
		t.Fatalf("reading the component directory: %v", err)
	}
	if len(sources) == 0 {
		// A glob that matches nothing would make this test pass by finding
		// nothing to check, which is the shape of failure it exists to prevent.
		t.Fatal("no component sources found; this test is checking nothing")
	}

	table, err := os.ReadFile("components_test.go")
	if err != nil {
		t.Fatalf("reading this file: %v", err)
	}
	body := string(table)

	for _, source := range sources {
		name := componentName(filepath.Base(source))
		if !strings.Contains(body, `{"`+name+`",`) {
			t.Errorf("components/%s ships %s, and TestEveryComponentRenders has no row for it",
				filepath.Base(source), name)
		}
	}
}

// componentName turns a source file name into the exported function it holds:
// stat-card.kyse.go is StatCard, and dropdown-menu.kyse.go is DropdownMenu.
func componentName(file string) string {
	base := strings.TrimSuffix(file, ".kyse.go")
	var out strings.Builder
	for _, part := range strings.Split(base, "-") {
		if part == "" {
			continue
		}
		out.WriteString(strings.ToUpper(part[:1]))
		out.WriteString(part[1:])
	}
	return out.String()
}
