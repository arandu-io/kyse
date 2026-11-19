package unit

import (
	"encoding/json"
	"html"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/translation"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// A field that leaves its files to the form is drawn with everything the
// behaviour needs to say what was chosen, and a field that posts each file as
// it is chosen is drawn without any of it, because the server's answer fills
// its list.

var fileBehaviourProps = regexp.MustCompile(`data-kyse-props="([^"]*)"`)

// fileProps reads the props the behaviour is given back out of the markup, the
// way the browser does: the attribute decoded, then parsed as JSON.
func fileProps(t *testing.T, markup string) map[string]any {
	t.Helper()
	match := fileBehaviourProps.FindStringSubmatch(markup)
	if match == nil {
		t.Fatalf("the root carries no props for the behaviour:\n%s", markup)
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(html.UnescapeString(match[1])), &props); err != nil {
		t.Fatalf("the props are not JSON: %v\n%s", err, match[1])
	}
	return props
}

func TestAFileFieldDrawsWhatTheBehaviourListsWith(t *testing.T) {
	t.Parallel()
	markup := string(components.FileUpload(components.FileUploadProps{
		Name: "photos", Label: "Photos", Multiple: true, MaxFiles: 6, MaxBytes: 5 << 20,
	}))

	// The count is a polite status line that is on the page before anything is
	// chosen, and empty, so a screen reader has it in the tree when the words
	// arrive.
	summary := regexp.MustCompile(`(?s)<p\s[^>]*data-part="summary"[^>]*></p>`).FindString(markup)
	if summary == "" {
		t.Fatalf("no empty status line for the count:\n%s", markup)
	}
	for _, want := range []string{`role="status"`, `aria-live="polite"`, `aria-atomic="true"`, `empty:sr-only`, `data-[over=true]:text-destructive`} {
		if !strings.Contains(summary, want) {
			t.Errorf("the status line lacks %s: %s", want, summary)
		}
	}

	// The list the page draws is not itself a live region: the status line is
	// what is announced, and a list that re-read every row on every pick would
	// be noise.
	list := regexp.MustCompile(`(?s)<ul\s[^>]*data-part="list"[^>]*>`).FindString(markup)
	if list == "" || strings.Contains(list, "aria-live") {
		t.Errorf("the list is missing or is a live region: %s", list)
	}

	// The row is a template outside the list, so a redraw of the list cannot
	// take it away, with every class on it so the stylesheet built from these
	// sources has them in both themes.
	template := regexp.MustCompile(`(?s)<template data-file-template>.*?</template>`).FindString(markup)
	if template == "" {
		t.Fatalf("no template for the row:\n%s", markup)
	}
	if strings.Index(markup, "</ul>") > strings.Index(markup, "<template") {
		t.Error("the template is inside the list the behaviour empties on every pick")
	}
	for _, want := range []string{`data-part="file"`, `data-part="thumbnail"`, `data-part="file-name"`, `data-part="file-size"`,
		`bg-muted`, `text-muted-foreground`, `group-data-[over=true]:text-destructive`, `alt=""`, `hidden`} {
		if !strings.Contains(template, want) {
			t.Errorf("the row lacks %s: %s", want, template)
		}
	}
	if strings.Contains(template, "<script") || strings.Contains(template, " style=") || strings.Contains(template, " on") {
		t.Errorf("the row carries script or style: %s", template)
	}

	props := fileProps(t, markup)
	for key, want := range map[string]any{
		"maxFiles": float64(6), "maxBytes": float64(5 << 20),
		"chosenOne":  "1 file chosen",
		"chosenMany": "{n} files chosen",
		"tooMany":    "Send up to {max} files at a time: choose again.",
		"tooBigOne":  "One file is over {size}: replace it before sending.",
		"tooBigMany": "{n} files are over {size}: replace them before sending.",
	} {
		if props[key] != want {
			t.Errorf("props[%q] = %v, want %v", key, props[key], want)
		}
	}
}

func TestAFileFieldWithNoLimitSendsNone(t *testing.T) {
	t.Parallel()
	props := fileProps(t, string(components.FileUpload(components.FileUploadProps{Name: "doc", Label: "Document"})))
	for _, key := range []string{"maxFiles", "maxBytes"} {
		if _, there := props[key]; there {
			t.Errorf("a field with no limit sends %s", key)
		}
	}
}

func TestAFieldThatPostsOnChangeLeavesTheListToTheServer(t *testing.T) {
	t.Parallel()
	markup := string(components.FileUpload(components.FileUploadProps{
		Name: "doc", Label: "Document", UploadURL: "/upload", MaxFiles: 3,
	}))
	for _, absent := range []string{`data-part="summary"`, "<template", `data-kyse-props`} {
		if strings.Contains(markup, absent) {
			t.Errorf("a field that posts each file draws %s:\n%s", absent, markup)
		}
	}
	list := regexp.MustCompile(`(?s)<ul\s[^>]*data-part="list"[^>]*>`).FindString(markup)
	if !strings.Contains(list, `aria-live="polite"`) {
		t.Errorf("the list the server fills is not polite: %s", list)
	}
}

func TestTheSentencesOfAFileFieldAreThePagesOwn(t *testing.T) {
	t.Parallel()
	loader := translation.NewArrayLoader().AddMessages("pt-BR", "kyse", translation.Lines{
		"file_upload.chosen_one":   "1 arquivo escolhido",
		"file_upload.chosen_many":  "{n} arquivos escolhidos",
		"file_upload.too_many":     "Envie até {max} por vez.",
		"file_upload.too_big_one":  "Um passa de {size}.",
		"file_upload.too_big_many": "{n} passam de {size}.",
	}, "")
	page := translatedPage{Page: view.Page{}, translator: translation.New(loader, "en", "en"), locale: "pt-BR"}

	props := fileProps(t, string(components.FileUpload(components.FileUploadProps{Name: "fotos", Label: "Fotos", Page: page})))
	for key, want := range map[string]string{
		"chosenOne": "1 arquivo escolhido", "chosenMany": "{n} arquivos escolhidos",
		"tooMany": "Envie até {max} por vez.", "tooBigOne": "Um passa de {size}.", "tooBigMany": "{n} passam de {size}.",
	} {
		if props[key] != want {
			t.Errorf("props[%q] = %v, want %q", key, props[key], want)
		}
	}

	// A sentence written on the props wins over the catalogue, and the ones it
	// does not name keep speaking for the page.
	props = fileProps(t, string(components.FileUpload(components.FileUploadProps{
		Name: "fotos", Label: "Fotos", Page: page, ChosenOne: "Uma foto",
	})))
	if props["chosenOne"] != "Uma foto" || props["chosenMany"] != "{n} arquivos escolhidos" {
		t.Errorf("a written sentence does not win: %v", props)
	}
}

func TestAFileFieldEscapesWhatACallerWrote(t *testing.T) {
	t.Parallel()
	const payload = `<script>alert(1)</script>`
	markup := string(components.FileUpload(components.FileUploadProps{
		Name: "x", Label: payload, ChosenOne: payload, ChosenMany: payload, TooMany: payload, TooBigOne: payload, TooBigMany: payload,
	}))
	if strings.Contains(markup, payload) {
		t.Errorf("a prop survived unescaped:\n%s", markup)
	}
}
