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
			return []string{string(components.Avatar(components.AvatarProps{ComponentProps: c, Name: "Ada Lovelace"}))}
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
