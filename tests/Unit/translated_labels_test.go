package unit

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/arandu-io/hesape/pagination"
	"github.com/arandu-io/hesape/translation"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// A component draws a few words nobody wrote -- "Cancel", "Next page",
// "Close" -- and drew them in English whatever the page was served in, so an
// application in another language wrote every one of them on every call. These
// hold that a page carrying a translator is asked first, under the keys the
// package doc lists, and that a page without one draws exactly what it drew.

// translatedPage is a screen's page data that carries the translator it was
// wired with and the locale of the request.
type translatedPage struct {
	view.Page
	translator *translation.Translator
	locale     string
}

func (p translatedPage) Translator() *translation.Translator { return p.translator }
func (p translatedPage) Locale() string                      { return p.locale }

var _ components.TranslatorPage = translatedPage{}

// documentedKeys are the keys the package doc lists, read from it.
func documentedKeys(t *testing.T) []string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("..", "..", "components", "doc.go"))
	if err != nil {
		t.Fatal(err)
	}
	var keys []string
	for _, match := range regexp.MustCompile(`(?m)^//\tkyse\.([a-z_.]+)\s`).FindAllStringSubmatch(string(body), -1) {
		keys = append(keys, match[1])
	}
	slices.Sort(keys)
	return keys
}

// askedKeys are the keys the components ask for, read from their sources.
func askedKeys(t *testing.T) []string {
	t.Helper()
	sources, err := filepath.Glob(filepath.Join("..", "..", "components", "*.kyse.go"))
	if err != nil {
		t.Fatal(err)
	}
	asked := map[string]bool{}
	pattern := regexp.MustCompile(`(?:say\([^,]+, |\{&lines\.\w+, )"([a-z_.]+)"`)
	for _, source := range sources {
		body, err := os.ReadFile(source)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range pattern.FindAllStringSubmatch(string(body), -1) {
			asked[match[1]] = true
		}
	}
	// SortName asks with a key chosen by the order, not a literal.
	for _, key := range []string{"table.ascending", "table.descending", "table.unsorted"} {
		asked[key] = true
	}
	keys := make([]string, 0, len(asked))
	for key := range asked {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func TestEveryKeyAComponentAsksForIsDocumented(t *testing.T) {
	documented, asked := documentedKeys(t), askedKeys(t)
	if len(documented) == 0 {
		t.Fatal("the package doc lists no key: this test proves nothing")
	}
	if !slices.Equal(documented, asked) {
		t.Errorf("the documented keys and the asked keys differ:\ndocumented %v\nasked      %v", documented, asked)
	}
}

// marked is a catalogue answering every documented key with a marker that
// keeps the English placeholders, so a rendering shows which key it drew.
func marked(t *testing.T) translation.Lines {
	lines := translation.Lines{}
	for _, key := range documentedKeys(t) {
		lines[key] = "PT[" + key + "]"
	}
	for key, placeholders := range map[string]string{
		"pagination.page":      " {n}",
		"table.selected":       " {n}/{total}",
		"table.ascending":      " {column}",
		"table.descending":     " {column}",
		"table.unsorted":       " {column}",
		"data_table.showing":   " {from}-{to}/{total}",
		"data_table.page":      " {n}",
		"data_table.page_of":   " {n}/{total}",
		"data_table.search":    " {label}",
		"data_table.no_result": "",
	} {
		lines[key] += placeholders
	}
	return lines
}

func ptPage(lines translation.Lines) translatedPage {
	loader := translation.NewArrayLoader().AddMessages("pt-BR", "kyse", lines, "")
	return translatedPage{translator: translation.New(loader, "en", "en"), locale: "pt-BR"}
}

// everyLabel renders every component that draws one of its own words, in the
// states that draw all of them, with no label written on the props.
func everyLabel(page components.Page) string {
	columns := []components.TableColumn{{Label: "Number", Key: "number", Sortable: true, Hideable: true}, {Label: "Status", Key: "status", Sortable: true}}
	rows := []components.TableRow{{Key: "1", Cells: []components.TableCell{{Text: "2026-114"}, {Text: "open"}}}}
	bulk := []components.ButtonProps{{Label: "Archive", Type: "submit"}}
	paginator := pagination.NewLengthAwarePage(1, 137, 20, 3, pagination.Options{Path: "/invoices"})
	return strings.Join([]string{
		string(components.Pagination(components.PaginationProps{Paginator: paginator, PageData: page})),
		string(components.Dialog(components.DialogProps{ID: "delete", Title: "Delete?", Action: "/x", Page: page})),
		string(components.FormDialog(components.FormDialogProps{ID: "rename", Title: "Rename", Action: "/x", Page: page})),
		string(components.Toast(components.ToastProps{Title: "Saved.", Dismissible: true, Page: page})),
		string(components.ThemeToggle(components.ThemeToggleProps{Page: page})),
		string(components.Table(components.TableProps{
			ID: "plain", Columns: columns, Rows: rows, SelectName: "ids", BulkActions: bulk,
			SortKey: "number", SortDir: "asc", SortURL: "/x?sort={sort}&dir={dir}", Page: page,
		})),
		string(components.DataTable(components.DataTableProps{
			ID: "invoices", Label: "Invoices", URL: "/invoices", Columns: columns, Rows: rows,
			SearchName: "q", SelectName: "ids", BulkActions: bulk, SortKey: "number", SortDir: "desc",
			Filters:   []components.DataTableFilter{{Key: "status", Label: "Status", Options: []components.SelectOption{{Value: "open", Label: "Open"}}}},
			Paginator: paginator, PageData: page,
		})),
	}, "\n")
}

func TestEveryDocumentedKeyIsDrawnFromTheCatalogue(t *testing.T) {
	html := everyLabel(ptPage(marked(t)))
	for _, key := range documentedKeys(t) {
		if !strings.Contains(html, "PT["+key+"]") {
			t.Errorf("kyse.%s was not drawn from the catalogue", key)
		}
	}
	for _, english := range []string{">Confirm<", ">Cancel<", ">Save<", `aria-label="Close"`, `aria-label="Pagination"`, "Select all rows", "Bulk actions", `aria-label="Change the theme"`, `aria-label="Theme"`, "<span>Auto</span>", "<span>Light</span>", "<span>Dark</span>"} {
		if strings.Contains(html, english) {
			t.Errorf("drew the English %s with a catalogue that carries it", english)
		}
	}
}

func TestACataloguePlaceholderIsFilled(t *testing.T) {
	html := everyLabel(ptPage(marked(t)))
	for _, want := range []string{"PT[pagination.page] 3", "PT[data_table.page_of] 3/7", "PT[data_table.showing] 41-41/137", "PT[data_table.search] Invoices"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q", want)
		}
	}
}

// A page that is not a TranslatorPage -- the view.Page every screen embeds, or
// one whose translator is nil -- draws byte for byte what no page drew.
func TestAPageWithNoTranslatorDrawsTheEnglish(t *testing.T) {
	before := everyLabel(nil)
	for name, page := range map[string]components.Page{
		"view.Page":      view.Page{},
		"nil translator": translatedPage{locale: "pt-BR"},
		"empty catalogue": translatedPage{
			translator: translation.New(translation.NewArrayLoader(), "en", "en"), locale: "pt-BR",
		},
	} {
		if after := everyLabel(page); after != before {
			t.Errorf("%s drew differently from no page at all", name)
		}
	}
}

func TestALabelWrittenOnThePropsWinsOverTheCatalogue(t *testing.T) {
	page := ptPage(marked(t))
	html := string(components.Dialog(components.DialogProps{
		ID: "delete", Title: "Delete?", Action: "/x", Page: page, ConfirmLabel: "Excluir", CancelLabel: "Voltar",
	}))
	if !strings.Contains(html, ">Excluir<") || !strings.Contains(html, ">Voltar<") || strings.Contains(html, "PT[") {
		t.Errorf("the catalogue replaced labels the caller wrote:\n%s", html)
	}
	lines := string(components.DataTable(components.DataTableProps{
		ID: "invoices", URL: "/invoices", PageData: page, Lines: components.DataTableLines{PerPage: "por página"},
		Complete: true, PageSize: 10, PageSizes: []int{10, 20},
	}))
	if !strings.Contains(lines, "por página") || strings.Contains(lines, "PT[data_table.per_page]") {
		t.Errorf("the catalogue replaced a line the caller wrote:\n%s", lines)
	}
}

// A key the catalogue does not carry says the English, beside the ones it does.
func TestAKeyTheCatalogueLacksSaysTheEnglish(t *testing.T) {
	page := ptPage(translation.Lines{"dialog.cancel": "Cancelar"})
	html := string(components.Dialog(components.DialogProps{ID: "delete", Title: "Delete?", Action: "/x", Page: page}))
	if !strings.Contains(html, ">Cancelar<") || !strings.Contains(html, ">Confirm<") {
		t.Errorf("want the translated way out and the English way through:\n%s", html)
	}
}

// The catalogue is the application's own file, nested as the package doc says,
// and the locale is the request's.
func TestTheCatalogueIsTheApplicationsKyseFile(t *testing.T) {
	files := fstest.MapFS{
		"lang/pt-BR/kyse.json": {Data: []byte(`{"pagination": {"next": "Próxima", "page": "Página {n}"}}`)},
		"lang/es/kyse.json":    {Data: []byte(`{"pagination": {"next": "Siguiente"}}`)},
	}
	loader, err := translation.NewFileLoader(files, "lang")
	if err != nil {
		t.Fatal(err)
	}
	translator := translation.New(loader, "en", "en")
	paginator := pagination.NewLengthAwarePage(20, 137, 20, 3, pagination.Options{Path: "/invoices"})

	pt := string(components.Pagination(components.PaginationProps{Paginator: paginator, PageData: translatedPage{translator: translator, locale: "pt-BR"}}))
	if !strings.Contains(pt, ">Próxima</a>") || !strings.Contains(pt, `aria-label="Página 4"`) {
		t.Errorf("the pt-BR file was not read:\n%s", pt)
	}
	es := string(components.Pagination(components.PaginationProps{Paginator: paginator, PageData: translatedPage{translator: translator, locale: "es"}}))
	if !strings.Contains(es, ">Siguiente</a>") || !strings.Contains(es, `aria-label="Page 4"`) {
		t.Errorf("the request's locale did not choose the file, or a missing line was not English:\n%s", es)
	}
}
