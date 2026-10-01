package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// An HTML parser keeps the first of two attributes with the same name and
// drops the second. So whichever of the component and the caller writes a name
// first owns it, and the rule is that the component does: every tag writes the
// caller's Attrs last. A caller's name, value, type, form or id can then fill a
// slot the component left empty, and can never replace what the component
// wrote -- the field a form submits under, the value it submits, the id a
// label and an aria-describedby point at.

// TestTheCallersAttributesAreTheLastOfEveryTag reads the sources: after an
// @attributes line, the next line closes the tag.
func TestTheCallersAttributesAreTheLastOfEveryTag(t *testing.T) {
	sources, err := filepath.Glob(filepath.Join("..", "..", "components", "*.kyse.go"))
	if err != nil || len(sources) == 0 {
		t.Fatalf("no component sources found: %v", err)
	}
	for _, source := range sources {
		data, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(string(data), "\n")
		for i, line := range lines {
			if !strings.HasPrefix(strings.TrimSpace(line), "@attributes(") {
				continue
			}
			next := ""
			for _, after := range lines[i+1:] {
				if next = strings.TrimSpace(after); next != "" {
					break
				}
			}
			if !strings.HasPrefix(next, ">") && !strings.HasPrefix(next, "/>") {
				t.Errorf("%s:%d: %q is followed by %q. The caller's attributes are the last of the tag, "+
					"so that on a name both write the component's comes first and is the one a browser keeps",
					filepath.Base(source), i+1, strings.TrimSpace(line), next)
			}
		}
	}
}

// TestTheComponentKeepsTheAttributesItWrites renders components with Attrs
// that collide with what they write, and checks the component's value comes
// first in the tag.
func TestTheComponentKeepsTheAttributesItWrites(t *testing.T) {
	forged := components.Attrs{"name": "forged", "value": "forged", "type": "hidden", "id": "forged"}
	cases := []struct {
		name   string
		html   string
		marker string
		owned  []string
	}{
		{
			name: "Input",
			html: string(components.Input(components.InputProps{
				Name: "email", Value: "ada@example.test",
				ComponentProps: components.ComponentProps{Attrs: forged},
			})),
			marker: `<input`,
			owned:  []string{`name="email"`, `value="ada@example.test"`, `type="text"`, `id="email"`},
		},
		{
			name: "Textarea",
			html: string(components.Textarea(components.TextareaProps{
				Name: "bio", Label: "Bio",
				ComponentProps: components.ComponentProps{Parts: components.Parts{"input": {Attrs: forged}}},
			})),
			marker: `<textarea`,
			owned:  []string{`name="bio"`, `id="bio"`},
		},
		{
			name: "Checkbox",
			html: string(components.Checkbox(components.CheckboxProps{
				Name: "terms", Label: "Terms", Value: "accepted",
				ComponentProps: components.ComponentProps{Parts: components.Parts{"input": {Attrs: forged}}},
			})),
			marker: `type="checkbox"`,
			owned:  []string{`name="terms"`, `value="accepted"`, `type="checkbox"`, `id="terms"`},
		},
		{
			name: "Switch",
			html: string(components.Switch(components.SwitchProps{
				Name: "notify", Label: "Notify", Value: "on",
				ComponentProps: components.ComponentProps{Parts: components.Parts{"input": {Attrs: forged}}},
			})),
			marker: `role="switch"`,
			owned:  []string{`name="notify"`, `value="on"`, `type="checkbox"`},
		},
		{
			name: "Combobox",
			html: string(components.Combobox(components.ComboboxProps{
				Name: "country", Label: "Country", SearchURL: "/countries",
				ComponentProps: components.ComponentProps{Parts: components.Parts{"input": {Attrs: forged}}},
			})),
			marker: `role="combobox"`,
			owned:  []string{`name="country-query"`, `type="text"`, `id="country"`},
		},
		{
			name: "Tabs",
			html: string(components.Tabs(components.TabsProps{
				ID: "settings", Active: "general",
				Tabs:           []components.Tab{{ID: "general", Label: "General"}},
				ComponentProps: components.ComponentProps{Parts: components.Parts{"trigger": {Attrs: forged}}},
			})),
			marker: `role="tab"`,
			owned:  []string{`type="button"`, `id="settings-tab-general"`},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tag := tagContaining(tc.html, tc.marker)
			if tag == "" {
				t.Fatalf("no tag contains %q:\n%s", tc.marker, tc.html)
			}
			for _, owned := range tc.owned {
				attr := owned[:strings.Index(owned, "=")]
				at := strings.Index(tag, " "+owned)
				if at < 0 {
					at = strings.Index(tag, "\t"+owned)
				}
				if at < 0 {
					t.Errorf("the component did not write %s in:\n%s", owned, tag)
					continue
				}
				if caller := strings.Index(tag, attr+`="forged"`); caller >= 0 && caller < at {
					t.Errorf("the caller's %s=\"forged\" comes before the component's %s, so a browser keeps the caller's:\n%s",
						attr, owned, tag)
				}
			}
		})
	}
}

// tagContaining is the start tag that holds marker, from its < to its >.
func tagContaining(html, marker string) string {
	at := strings.Index(html, marker)
	if at < 0 {
		return ""
	}
	start := strings.LastIndex(html[:at+1], "<")
	end := strings.Index(html[at:], ">")
	if start < 0 || end < 0 {
		return ""
	}
	return html[start : at+end+1]
}
