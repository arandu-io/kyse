package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// The page every screen embeds is hesape's, and the contract says so: a doc
// that named a page nobody embeds any more sends a reader to the wrong package
// to find what a field asks.
var _ components.Page = view.Page{}

func TestThePageContractNamesThePageEveryScreenEmbeds(t *testing.T) {
	body, err := os.ReadFile(filepath.Join("..", "..", "components", "page.go"))
	if err != nil {
		t.Fatal(err)
	}
	doc := string(body)
	if !strings.Contains(doc, "hesape/view.Page implements") {
		t.Error("page.go does not name hesape/view.Page as the implementation of Page")
	}
	if strings.Contains(doc, "framework/view.Page") {
		t.Error("page.go still names framework/view.Page, which no screen embeds")
	}
}
