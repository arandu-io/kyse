package unit

import (
	"html/template"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// A response of 400 or more is swapped by YPY into the target unless the
// element says otherwise, so a refused delete would replace the row with an
// error page and a failed fetch would replace a list with one. Every element a
// component draws to ask the server for a fragment says that a refusal swaps
// nothing; the form that wants its 422 back keeps the native rule and refuses
// only a failure of the server.
func TestEveryElementThatAsksForAFragmentSaysWhatARefusalSwaps(t *testing.T) {
	t.Parallel()

	verbs := regexp.MustCompile(`<[^<>]*\sypy-(get|post|put|delete)=[^<>]*>`)
	for name, html := range map[string]template.HTML{
		"Button": components.Button(components.ButtonProps{Label: "Go", YpyPost: "/go", YpyTarget: "#out"}),
		"SplitButton": components.SplitButton(components.SplitButtonProps{
			Label: "Go", YpyPost: "/go", YpyTarget: "#out",
		}),
		"DeleteRow":        components.DeleteRow(components.DeleteRowProps{URL: "/rows/1"}),
		"LoadMore":         components.LoadMore(components.LoadMoreProps{URL: "/more"}),
		"LazyLoad":         components.LazyLoad(components.LazyLoadProps{URL: "/lazy", Label: "Panel"}),
		"OptimisticToggle": components.OptimisticToggle(components.OptimisticToggleProps{URL: "/toggle", Label: "Star"}),
		"ActiveSearch":     components.ActiveSearch(components.ActiveSearchProps{Name: "q", Label: "Search", URL: "/search"}),
		"Select": components.Select(components.SelectProps{
			Name: "kind", Label: "Kind", YpyGet: "/kinds", YpyTarget: "#out",
			Options: []components.SelectOption{{Label: "A", Value: "a"}},
		}),
		"SegmentedControl": components.SegmentedControl(components.SegmentedControlProps{
			Name: "view", Label: "View", YpyGet: "/view", YpyTarget: "#out",
			Options: []components.SegmentedOption{{Label: "A", Value: "a"}, {Label: "B", Value: "b"}},
		}),
		"EditInPlace": components.EditInPlace(components.EditInPlaceProps{
			ID: "t", Name: "t", Label: "T", Value: "v", Editing: true, SaveURL: "/t", CancelURL: "/t/cancel",
		}),
	} {
		found := verbs.FindAllString(string(html), -1)
		if len(found) == 0 {
			t.Errorf("%s draws no element that asks the server for anything:\n%s", name, html)
			continue
		}
		for _, tag := range found {
			if !strings.Contains(tag, `ypy-error-swap="none"`) && !strings.Contains(tag, `ypy-status:5xx="swap:none"`) {
				t.Errorf("%s asks the server and does not say what a refusal swaps:\n%s", name, tag)
			}
		}
	}
}
