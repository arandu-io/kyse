package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

func controlled() components.DataTableProps {
	p := listing()
	p.State = components.DataTableValues{
		Q: "ada", Sort: "total", Dir: "desc", Size: 20, DefaultSize: 10,
		Facets: map[string][]string{
			"status": {"paid"}, "kind": {"b"}, "code": {"AB-1"}, "de": {"2026-10-01"}, "ate": {"2026-10-31"},
		},
	}
	p.Controls = []components.DataTableControl{
		{Key: "code", Label: "Código", Placeholder: "AB-1"},
		{Key: "kind", Kind: components.DataTableSelect, Label: "Tipo", Placeholder: "Todos",
			Options: []components.SelectOption{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}}},
		{Key: "de", ToKey: "ate", Kind: components.DataTablePeriod, Label: "Período",
			Period: components.DateRangePickerProps{DateFormat: "02/01/2006", Placeholder: "Escolher período"}},
	}
	return p
}

// The controls sit in the one form that owns the search, so a request carries
// all of them; the value each holds is the state's.
func TestTheControlsAreFieldsOfTheQueryFormHoldingTheStatesValues(t *testing.T) {
	html := string(components.DataTable(controlled()))
	form := html[strings.Index(html, `<form`):strings.Index(html, `</form>`)]
	form = strings.ReplaceAll(strings.Join(strings.Fields(form), " "), " >", ">")

	for _, want := range []string{
		`data-part="control" data-control="text"`, `data-part="control" data-control="select"`, `data-part="control" data-control="period"`,
		`name="code"`, `value="AB-1"`, `placeholder="AB-1"`,
		`<option value="b" selected>B</option>`, `<option value="">Todos</option>`,
		`name="de"`, `value="2026-10-01"`, `name="ate"`, `value="2026-10-31"`, `01/10/2026 – 31/10/2026`,
		`for="invoices-control-code"`, `id="invoices-control-kind"`,
	} {
		if !strings.Contains(form, want) {
			t.Errorf("the query form lacks %s", want)
		}
	}
	for _, key := range []string{"code", "kind", "de", "ate"} {
		if n := strings.Count(form, `name="`+key+`"`); n != 1 {
			t.Errorf("the field %q is written %d times in the form: a control's value is its own field and never also a hidden copy", key, n)
		}
	}
}

// The state owns the keys: a new order, a page and a size keep what the controls
// hold, and a control's value is never doubled into an address.
func TestEveryAddressOfTheTableCarriesTheControls(t *testing.T) {
	p := controlled()
	for name, href := range map[string]string{
		"sort": p.SortHref(p.Columns[1]), "page": p.PageHref(3), "size": p.SizeHref(50), "column": p.ColumnHref(p.Columns[0]),
	} {
		query := parsed(t, href)
		for key, want := range map[string]string{"code": "AB-1", "kind": "b", "de": "2026-10-01", "ate": "2026-10-31"} {
			if got := query[key]; len(got) != 1 || got[0] != want {
				t.Errorf("%s: %s=%v in %s, want once %s", name, key, got, href, want)
			}
		}
	}
	// Reset drops them with the search.
	if reset := parsed(t, p.ResetHref()); len(reset["code"])+len(reset["kind"])+len(reset["de"])+len(reset["q"]) != 0 {
		t.Errorf("Reset keeps %v", reset)
	}
}

// A narrowed list is a list with Reset, and the hidden fields of the form do not
// repeat what a control already sends.
func TestAControlNarrowsTheListAndIsNotAHiddenField(t *testing.T) {
	p := controlled()
	p.State = components.DataTableValues{Facets: map[string][]string{"kind": {"a"}}}
	if !p.Active() {
		t.Error("a chosen control does not narrow the list")
	}
	if !strings.Contains(string(components.DataTable(p)), `data-part="reset"`) {
		t.Error("a list narrowed only by a control has no Reset")
	}
	p.State = components.DataTableValues{}
	if p.Active() {
		t.Error("a list with no control chosen is narrowed")
	}
	for _, field := range controlled().QueryFields() {
		switch field.Name {
		case "code", "kind", "de", "ate":
			t.Errorf("the form also carries %s as a hidden field", field.Name)
		}
	}
}

// A list with no search and no facet but a control still has the form, its
// request and the fallback button for a page without script.
func TestAListWithOnlyControlsIsAQueryForm(t *testing.T) {
	p := controlled()
	p.Searchable, p.Filters = false, nil
	html := string(components.DataTable(p))
	for _, want := range []string{`hx-get="/invoices"`, `hx-trigger="submit, change"`, `<noscript>`, `data-part="query"`} {
		if !strings.Contains(html, want) {
			t.Errorf("a list with only controls lacks %s", want)
		}
	}
}

// A control's value falls back to the one it was given when the state has none,
// and a period shows its shortcuts as the picker does.
func TestAControlHoldsItsOwnValueWhenTheStateHasNone(t *testing.T) {
	p := controlled()
	p.State = components.DataTableValues{}
	p.Controls[0].Value = "ZZ-9"
	p.Controls[2].Value, p.Controls[2].ToValue = "2026-09-01", "2026-09-30"
	p.Controls[2].Period.Presets = []components.DateRangePreset{{Label: "Ontem", Days: 1, EndsAgo: 1}}
	html := string(components.DataTable(p))
	for _, want := range []string{`value="ZZ-9"`, `value="2026-09-01"`, `01/09/2026 – 30/09/2026`, `data-range-ends-ago="1"`} {
		if !strings.Contains(html, want) {
			t.Errorf("the table lacks %s", want)
		}
	}
}

// On a complete list the form is the browser's: no request, and the rows carry
// the values the controls narrow them by.
func TestTheControlsOfACompleteListAreNotRequests(t *testing.T) {
	p := controlled()
	p.Complete, p.Paginator = true, nil
	p.Rows = []components.TableRow{{Key: "1", Facets: map[string]string{"kind": "b", "code": "AB-1", "de": "2026-10-05"},
		Cells: []components.TableCell{{Text: "1"}, {Text: "2"}}}}
	html := string(components.DataTable(p))
	if strings.Contains(html, `id="invoices-query"`) && strings.Contains(html[strings.Index(html, `id="invoices-query"`):strings.Index(html, `</form>`)], "hx-get") {
		t.Error("the form of a complete list asks the server")
	}
	for _, want := range []string{`data-facet-kind="b"`, `data-facet-code="AB-1"`, `data-facet-de="2026-10-05"`, `data-complete="true"`} {
		if !strings.Contains(html, want) {
			t.Errorf("a complete list lacks %s", want)
		}
	}
}
