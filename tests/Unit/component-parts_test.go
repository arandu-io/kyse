package unit

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// extensible is one row per component that embeds ComponentProps: how to render
// it with a caller's props, and the parts it says it publishes.
//
// Both halves come from the real type -- the render calls the real function and
// partNames calls the real method -- so neither can drift from what ships
// without this file failing to compile. What is hand-kept is only the
// membership of the list, and TestEveryExtensibleComponentIsInThisTable reads
// the directory to hold that.
//
// render answers one string per state, because some parts cannot coexist. A
// field draws a message or a hint and never both, so no single rendering shows
// everything it publishes, and a test that demanded one would be asking the
// component to do something it must not do.
var extensible = []struct {
	name      string
	render    func(components.ComponentProps) []string
	partNames func() []string
}{
	{
		"Alert",
		func(c components.ComponentProps) []string {
			return []string{string(components.Alert(components.AlertProps{
				ComponentProps: c, Title: "Saved", Message: "The post is live.",
			}))}
		},
		components.AlertProps{}.PartNames,
	},
	{
		"Avatar",
		func(c components.ComponentProps) []string {
			props := components.AvatarProps{ComponentProps: c, Name: "Ada Lovelace"}
			withInitials := string(components.Avatar(props))
			props.ImageURL = "/ada.png"
			return []string{withInitials, string(components.Avatar(props))}
		},
		components.AvatarProps{}.PartNames,
	},
	{
		"Badge",
		func(c components.ComponentProps) []string {
			return []string{string(components.Badge(components.BadgeProps{ComponentProps: c, Label: "draft"}))}
		},
		components.BadgeProps{}.PartNames,
	},
	{
		"Button",
		func(c components.ComponentProps) []string {
			return []string{string(components.Button(components.ButtonProps{ComponentProps: c, Label: "Save"}))}
		},
		components.ButtonProps{}.PartNames,
	},
	{
		"Breadcrumb",
		func(c components.ComponentProps) []string {
			return []string{string(components.Breadcrumb(components.BreadcrumbProps{
				ComponentProps: c,
				Items: []components.Crumb{
					{Label: "Home", URL: "/"},
					{Label: "Posts", URL: "/posts"},
					{Label: "A post"},
				},
			}))}
		},
		components.BreadcrumbProps{}.PartNames,
	},
	{
		"StatCard",
		func(c components.ComponentProps) []string {
			return []string{string(components.StatCard(components.StatCardProps{
				ComponentProps: c,
				Title:          "Connections",
				Meta:           "read at 21:44",
				Columns:        []string{"Open"},
				Rows:           []components.StatRow{{Label: "acme", Values: []string{"12"}}},
			}))}
		},
		components.StatCardProps{}.PartNames,
	},
	{
		"Table",
		func(c components.ComponentProps) []string {
			return []string{string(components.Table(components.TableProps{
				ComponentProps: c,
				Caption:        "Invoices",
				Columns:        []components.TableColumn{{Label: "Number"}},
				Rows:           []components.TableRow{{Cells: []components.TableCell{{Text: "0001"}}}},
			}))}
		},
		components.TableProps{}.PartNames,
	},
	{
		"Tabs",
		func(c components.ComponentProps) []string {
			return []string{string(components.Tabs(components.TabsProps{
				ComponentProps: c,
				ID:             "settings",
				Tabs:           []components.Tab{{Label: "General", Panel: template.HTML("<p>x</p>")}},
			}))}
		},
		components.TabsProps{}.PartNames,
	},
	{
		"Toast",
		func(c components.ComponentProps) []string {
			return []string{string(components.Toast(components.ToastProps{
				ComponentProps: c,
				Title:          "Saved",
				Message:        "The post is live.",
				ActionLabel:    "Undo",
				ActionURL:      "/undo",
			}))}
		},
		components.ToastProps{}.PartNames,
	},
	{
		"Accordion",
		func(c components.ComponentProps) []string {
			return []string{string(components.Accordion(components.AccordionProps{
				ComponentProps: c,
				Items: []components.AccordionItem{
					{Label: "What is it", Content: template.HTML("<p>A framework.</p>")},
				},
			}))}
		},
		components.AccordionProps{}.PartNames,
	},
	{
		"Collapsible",
		func(c components.ComponentProps) []string {
			return []string{string(components.Collapsible(components.CollapsibleProps{
				ComponentProps: c, Label: "Details", Content: "More.",
			}))}
		},
		components.CollapsibleProps{}.PartNames,
	},
	{
		"Dialog",
		func(c components.ComponentProps) []string {
			return []string{string(components.Dialog(components.DialogProps{
				ComponentProps: c,
				ID:             "confirm",
				Title:          "Delete the post?",
				Message:        "This cannot be undone.",
				Action:         "/posts/1",
				Token:          "t",
			}))}
		},
		components.DialogProps{}.PartNames,
	},
	{
		"Drawer",
		func(c components.ComponentProps) []string {
			return []string{string(components.Drawer(components.DrawerProps{
				ComponentProps: c,
				ID:             "menu",
				Title:          "Menu",
				Description:    "Where to go.",
				Links:          []components.DrawerLink{{Label: "Home", Href: "/"}},
			}))}
		},
		components.DrawerProps{}.PartNames,
	},
	{
		"Popover",
		func(c components.ComponentProps) []string {
			return []string{string(components.Popover(components.PopoverProps{
				ComponentProps: c,
				ID:             "help",
				Label:          "Help",
				Title:          "About this",
				Description:    "What it does.",
				Content:        template.HTML("<p>x</p>"),
			}))}
		},
		components.PopoverProps{}.PartNames,
	},
	{
		"InputGroup",
		func(c components.ComponentProps) []string {
			props := components.InputGroupProps{
				ComponentProps: c, Name: "search", Label: "Search",
				Start: "/", End: "go", Hint: "Press enter.",
			}
			withHint := string(components.InputGroup(props))
			props.Page = page{errs: map[string]string{"search": "Required."}}
			return []string{withHint, string(components.InputGroup(props))}
		},
		components.InputGroupProps{}.PartNames,
	},
	{
		"RadioGroup",
		func(c components.ComponentProps) []string {
			props := components.RadioGroupProps{
				ComponentProps: c, Name: "plan", Label: "Plan", Hint: "Change any time.",
				Options: []components.RadioOption{{Label: "Free", Value: "free"}},
			}
			withHint := string(components.RadioGroup(props))
			props.Page = page{errs: map[string]string{"plan": "Required."}}
			return []string{withHint, string(components.RadioGroup(props))}
		},
		components.RadioGroupProps{}.PartNames,
	},
	{
		"RangeSlider",
		func(c components.ComponentProps) []string {
			props := components.RangeSliderProps{
				ComponentProps: c, Name: "volume", Label: "Volume",
				ShowValue: true, Hint: "Louder to the right.",
			}
			withHint := string(components.RangeSlider(props))
			props.Page = page{errs: map[string]string{"volume": "Required."}}
			return []string{withHint, string(components.RangeSlider(props))}
		},
		components.RangeSliderProps{}.PartNames,
	},
	{
		"Select",
		func(c components.ComponentProps) []string {
			props := components.SelectProps{
				ComponentProps: c, Name: "plan", Label: "Plan", Hint: "Change any time.",
				Options: []components.SelectOption{{Label: "Free", Value: "free"}},
			}
			withHint := string(components.Select(props))
			props.Page = page{errs: map[string]string{"plan": "Required."}}
			return []string{withHint, string(components.Select(props))}
		},
		components.SelectProps{}.PartNames,
	},
	{
		"ThemeToggle",
		func(c components.ComponentProps) []string {
			return []string{string(components.ThemeToggle(components.ThemeToggleProps{ComponentProps: c}))}
		},
		components.ThemeToggleProps{}.PartNames,
	},
	{
		"Combobox",
		func(c components.ComponentProps) []string {
			props := components.ComboboxProps{
				ComponentProps: c, Name: "city", Label: "City", SearchURL: "/cities",
				Hint: "Start typing.", Options: []components.ComboboxOption{{Label: "Lisbon", Value: "lis"}},
			}
			withHint := string(components.Combobox(props))
			props.Page = page{errs: map[string]string{"city": "Required."}}
			return []string{withHint, string(components.Combobox(props))}
		},
		components.ComboboxProps{}.PartNames,
	},
	{
		"Command",
		func(c components.ComponentProps) []string {
			return []string{string(components.Command(components.CommandProps{
				ComponentProps: c,
				ID:             "palette",
				Label:          "Commands",
				Groups: []components.CommandGroup{{
					Heading: "Posts",
					Items:   []components.CommandItem{{Label: "New post", URL: "/posts/new"}},
				}},
			}))}
		},
		components.CommandProps{}.PartNames,
	},
	{
		"DropdownMenu",
		func(c components.ComponentProps) []string {
			return []string{string(components.DropdownMenu(components.DropdownMenuProps{
				ComponentProps: c,
				ID:             "actions",
				Label:          "Actions",
				Items:          []components.MenuItem{{Label: "Edit", URL: "/edit"}},
			}))}
		},
		components.DropdownMenuProps{}.PartNames,
	},
	{
		"Sidebar",
		func(c components.ComponentProps) []string {
			return []string{string(components.Sidebar(components.SidebarProps{
				ComponentProps: c,
				ID:             "main",
				Header:         template.HTML("<b>Acme</b>"),
				Footer:         template.HTML("<b>v1</b>"),
				Groups: []components.SidebarGroup{{
					Label: "Manage",
					Items: []components.SidebarItem{{Label: "Posts", URL: "/posts"}},
				}},
			}))}
		},
		components.SidebarProps{}.PartNames,
	},
	{
		"Card",
		func(c components.ComponentProps) []string {
			return []string{string(components.Card(components.CardProps{
				ComponentProps: c,
				Title:          "A title",
				Description:    "A sentence.",
				Meta:           "yesterday",
			}))}
		},
		components.CardProps{}.PartNames,
	},
	{
		"ButtonGroup",
		func(c components.ComponentProps) []string {
			return []string{string(components.ButtonGroup(components.ButtonGroupProps{
				ComponentProps: c,
				Label:          "Message actions",
				Buttons:        []components.ButtonProps{{Label: "Archive"}, {Label: "Report"}},
			}))}
		},
		components.ButtonGroupProps{}.PartNames,
	},
	{
		"Empty",
		func(c components.ComponentProps) []string {
			return []string{string(components.Empty(components.EmptyProps{
				ComponentProps: c,
				Title:          "No posts",
				Message:        "Nothing has been published.",
				ActionLabel:    "Write one",
				ActionURL:      "/posts/new",
			}))}
		},
		components.EmptyProps{}.PartNames,
	},
	{
		"Checkbox",
		func(c components.ComponentProps) []string {
			props := components.CheckboxProps{ComponentProps: c, Name: "terms", Label: "I agree", Hint: "You can change this later."}
			withHint := string(components.Checkbox(props))
			props.Page = page{errs: map[string]string{"terms": "Required."}}
			return []string{withHint, string(components.Checkbox(props))}
		},
		components.CheckboxProps{}.PartNames,
	},
	{
		"Field",
		func(c components.ComponentProps) []string {
			props := components.FieldProps{ComponentProps: c, Name: "email", Label: "Email", Hint: "We never share it."}
			withHint := string(components.Field(props))
			props.Page = page{errs: map[string]string{"email": "Required."}}
			return []string{withHint, string(components.Field(props))}
		},
		components.FieldProps{}.PartNames,
	},
	{
		"Input",
		func(c components.ComponentProps) []string {
			return []string{string(components.Input(components.InputProps{ComponentProps: c, Name: "email"}))}
		},
		components.InputProps{}.PartNames,
	},
	{
		"Item",
		func(c components.ComponentProps) []string {
			return []string{string(components.Item(components.ItemProps{
				ComponentProps: c,
				Title:          "Billing",
				Description:    "Cards and invoices.",
				Icon:           template.HTML(`<svg></svg>`),
				Action:         template.HTML(`<button></button>`),
			}))}
		},
		components.ItemProps{}.PartNames,
	},
	{
		"Kbd",
		func(c components.ComponentProps) []string {
			return []string{string(components.Kbd(components.KbdProps{ComponentProps: c, Keys: []string{"⌘", "K"}}))}
		},
		components.KbdProps{}.PartNames,
	},
	{
		"Label",
		func(c components.ComponentProps) []string {
			return []string{string(components.Label(components.LabelProps{
				ComponentProps: c, For: "email", Text: "Email", Required: true,
			}))}
		},
		components.LabelProps{}.PartNames,
	},
	{
		"Progress",
		func(c components.ComponentProps) []string {
			return []string{string(components.Progress(components.ProgressProps{
				ComponentProps: c, Label: "Upload", Value: 40,
			}))}
		},
		components.ProgressProps{}.PartNames,
	},
	{
		"Separator",
		func(c components.ComponentProps) []string {
			return []string{string(components.Separator(components.SeparatorProps{ComponentProps: c}))}
		},
		components.SeparatorProps{}.PartNames,
	},
	{
		"Skeleton",
		func(c components.ComponentProps) []string {
			return []string{string(components.Skeleton(components.SkeletonProps{ComponentProps: c}))}
		},
		components.SkeletonProps{}.PartNames,
	},
	{
		"Switch",
		func(c components.ComponentProps) []string {
			props := components.SwitchProps{ComponentProps: c, Name: "alerts", Label: "Email alerts", Hint: "Sent once a day."}
			withHint := string(components.Switch(props))
			props.Page = page{errs: map[string]string{"alerts": "Required."}}
			return []string{withHint, string(components.Switch(props))}
		},
		components.SwitchProps{}.PartNames,
	},
	{
		"Textarea",
		func(c components.ComponentProps) []string {
			props := components.TextareaProps{ComponentProps: c, Name: "body", Label: "Body", Hint: "Markdown is allowed."}
			withHint := string(components.Textarea(props))
			props.Page = page{errs: map[string]string{"body": "Required."}}
			return []string{withHint, string(components.Textarea(props))}
		},
		components.TextareaProps{}.PartNames,
	},
}

// TestEveryComponentTakesAClass is the first half of the promise: a caller can
// add a class to any component, and it reaches the element.
//
// A component migrated by embedding ComponentProps and then left drawing
// class="btn" compiles, renders, and silently ignores everything a caller
// writes. This is what catches that, here rather than on somebody's page.
func TestEveryComponentTakesAClass(t *testing.T) {
	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			for _, got := range c.render(components.ComponentProps{Class: "kyse-probe-root"}) {
				if !strings.Contains(got, "kyse-probe-root") {
					t.Fatalf("the class a caller added is not in the output:\n%s", got)
				}
			}
		})
	}
}

// TestEveryPublishedPartIsReachable is the other half, and it asserts in both
// directions.
//
// A name is only a promise if writing it changes something, and the published
// list is only true if it names what is actually drawn. So: every name in
// PartNames has to reach an element, and every data-part rendered has to be a
// name PartNames publishes. A part that is drawn and not published is a handle
// nobody can find; one that is published and not drawn is a handle that does
// nothing.
func TestEveryPublishedPartIsReachable(t *testing.T) {
	drawn := regexp.MustCompile(`data-part="([a-z-]+)"`)

	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			published := c.partNames()
			if len(published) == 0 {
				t.Fatal("the component publishes no parts, and every component has a root")
			}

			for _, part := range published {
				probe := "kyse-probe-" + part
				states := c.render(components.ComponentProps{
					Parts: components.Parts{part: {Class: probe}},
				})
				reached := false
				for _, got := range states {
					reached = reached || strings.Contains(got, probe)
				}
				if !reached {
					t.Errorf("the part %q is published and a class written for it appears in none of its %d states:\n%s",
						part, len(states), strings.Join(states, "\n---\n"))
				}
			}

			var found []string
			for _, got := range c.render(components.ComponentProps{}) {
				for _, m := range drawn.FindAllStringSubmatch(got, -1) {
					found = append(found, m[1])
				}
			}
			if diff := missing(found, published); len(diff) > 0 {
				t.Errorf("these parts are drawn and not published: %v", diff)
			}
			if diff := missing(published, found); len(diff) > 0 {
				t.Errorf("these parts are published and not drawn: %v", diff)
			}
		})
	}
}

// TestAnAttrCannotCarryAScript is what stands between a bag of attributes and
// the policy this framework serves under.
//
// The refusal is view.Attributes's, and this is the assertion that a component
// actually goes through it. A component writing the bag some other way would
// pass every other test in this package and put an onclick on the page.
func TestAnAttrCannotCarryAScript(t *testing.T) {
	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			for _, attrs := range []components.Attrs{
				{"onclick": "alert(1)"},
				{"hx-on:click": "alert(1)"},
				{"x-data": "{}"},
				{"style": "display:none"},
				{"href": "javascript:alert(1)"},
				{`x" onerror="alert(1)`: "1"},
			} {
				for _, got := range c.render(components.ComponentProps{Attrs: attrs}) {
					for name := range attrs {
						if strings.Contains(got, name) {
							t.Errorf("the attribute %q was written:\n%s", name, got)
						}
					}
				}
			}
		})
	}
}

// TestAnInertAttrIsWritten is the other side of the refusals: the attributes a
// caller reaches for have to actually arrive, or the field is a list of things
// that do not work.
func TestAnInertAttrIsWritten(t *testing.T) {
	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			for _, got := range c.render(components.ComponentProps{
				Attrs: components.Attrs{"data-testid": "probe"},
			}) {
				if !strings.Contains(got, `data-testid="probe"`) {
					t.Fatalf("an inert attribute did not reach the element:\n%s", got)
				}
			}
		})
	}
}

// bridgeEvents are the events the client script listens for on the document,
// which is the whole set an action can answer.
//
// They are listed rather than sampled because the attribute name is assembled
// from the event -- data-kyse-on-<event> -- and a component that wrote one of
// them and dropped the rest would pass a test that only asked about click.
var bridgeEvents = map[string]string{
	"change":  "kyse-probe-change",
	"click":   "kyse-probe-click",
	"input":   "kyse-probe-input",
	"keydown": "kyse-probe-keydown",
	"submit":  "kyse-probe-submit",
}

// TestEveryComponentCarriesTheClientBridge is what a behaviour needs to be
// reachable at all: the name, the props and the events have to arrive on the
// element, or an application registers a behaviour that nothing ever mounts.
//
// It asserts on the outermost tag rather than on the whole rendering, because
// where the attributes land decides what happens. A behaviour is mounted on the
// element carrying data-kyse-behavior and torn down with it, so the same
// attribute on an inner element gives the behaviour a shorter life than the
// component it was written for.
//
// The props are checked in their escaped spelling, which is the one the browser
// is sent. The value is JSON inside an HTML attribute, so the quotes travel as
// &#34; and are decoded by the parser before any script reads them; a component
// that wrote the raw text would end the attribute at the first quote and hand
// the behaviour a fragment.
//
// Every component takes these fields because every component embeds
// ComponentProps. So a component that takes them and writes nothing compiles,
// renders, and is silent -- which is what this exists to make loud, and why it
// runs over the whole table rather than over a component somebody remembered.
func TestEveryComponentCarriesTheClientBridge(t *testing.T) {
	props := components.ComponentProps{
		Behavior: components.Behavior{
			Name:  "kyse-probe-behaviour",
			Props: map[string]any{"confirm": true},
		},
		Events: components.Events{},
	}
	for event, action := range bridgeEvents {
		props.Events[event] = action
	}

	want := []string{
		`data-kyse-behavior="kyse-probe-behaviour"`,
		`data-kyse-props="{&#34;confirm&#34;:true}"`,
	}
	for event, action := range bridgeEvents {
		want = append(want, `data-kyse-on-`+event+`="`+action+`"`)
	}
	sort.Strings(want)

	for _, c := range extensible {
		t.Run(c.name, func(t *testing.T) {
			for _, got := range c.render(props) {
				root := rootTag(t, got)
				for _, attribute := range want {
					if !strings.Contains(root, attribute) {
						t.Errorf("the outermost element does not carry %s:\n%s", attribute, root)
					}
				}
			}
		})
	}
}

// TestTheBridgeIsWrittenInOnePlace holds the emission down to ComponentProps.
//
// The attributes are one bag, RootAttrs, and every component writes that bag.
// A component that spelled an attribute into its markup instead would work,
// which is the problem: it would work for that component, at that name, and
// leave the next one to be remembered rather than compiled. Thirty-seven copies
// of one emission is thirty-seven places for the thirty-eighth to be forgotten.
//
// Both the sources and what they compile to are read, for the reason
// TestNoComponentEvaluatesAnExpression reads both: the source is where the line
// would be written, and the compiled file is what a browser is actually sent.
func TestTheBridgeIsWrittenInOnePlace(t *testing.T) {
	for _, pattern := range []string{"*.kyse.go", "*.go"} {
		files, err := filepath.Glob(filepath.Join("..", "..", "components", pattern))
		if err != nil {
			t.Fatalf("reading the component directory: %v", err)
		}
		if len(files) == 0 {
			t.Fatalf("no components/%s found; this test is checking nothing", pattern)
		}

		for _, file := range files {
			base := filepath.Base(file)
			// *.go matches *.kyse.go too, and component-props.go is where the
			// attribute names are supposed to be written.
			if (pattern == "*.go" && strings.HasSuffix(file, ".kyse.go")) || base == "component-props.go" {
				continue
			}
			body, err := os.ReadFile(file)
			if err != nil {
				t.Fatalf("reading %s: %v", file, err)
			}
			for at, line := range strings.Split(string(body), "\n") {
				if !strings.Contains(line, "data-kyse-") {
					continue
				}
				t.Errorf("%s:%d writes a bridge attribute of its own:\n\t%s\n"+
					"\tThe behaviour, the props and the events are the ComponentProps fields, and RootAttrs is what writes them.",
					base, at+1, strings.TrimSpace(line))
			}
		}
	}
}

// rootTag is the outermost tag of a rendering, which is the element every
// component publishes as its root.
func rootTag(t *testing.T, markup string) string {
	t.Helper()

	tag := firstTag.FindString(markup)
	if tag == "" {
		t.Fatalf("the rendering opens no tag:\n%s", markup)
	}
	return tag
}

var firstTag = regexp.MustCompile(`(?s)<[a-zA-Z][^>]*>`)

// TestEveryExtensibleComponentIsInThisTable reads the component directory, the
// way TestEveryComponentIsInTheTable does, so that migrating a component and
// forgetting this file fails here rather than leaving the component untested.
//
// A source embedding ComponentProps is the definition of migrated, and it is
// the same string this test greps for -- so the answer comes from the tree
// rather than from anybody's memory of it.
func TestEveryExtensibleComponentIsInThisTable(t *testing.T) {
	sources, err := filepath.Glob(filepath.Join("..", "..", "components", "*.kyse.go"))
	if err != nil {
		t.Fatalf("reading the component directory: %v", err)
	}
	if len(sources) == 0 {
		t.Fatal("no component sources found; this test is checking nothing")
	}

	listed := map[string]bool{}
	for _, c := range extensible {
		listed[c.name] = true
	}

	for _, source := range sources {
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatalf("reading %s: %v", source, err)
		}
		// The field, on a line of its own, and not the doc comment above it.
		if !regexp.MustCompile(`(?m)^\tComponentProps$`).Match(body) {
			continue
		}
		name := componentName(filepath.Base(source))
		if !listed[name] {
			t.Errorf("components/%s embeds ComponentProps and has no row in extensible",
				filepath.Base(source))
		}
	}
}

// missing returns the members of a that b does not have, sorted.
func missing(a, b []string) []string {
	have := map[string]bool{}
	for _, s := range b {
		have[s] = true
	}
	seen := map[string]bool{}
	var out []string
	for _, s := range a {
		if have[s] || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}
