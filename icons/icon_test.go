package icons_test

import (
	"encoding/xml"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/icons"
)

// What is checked here is what vendoring 1512 pictures from another project can
// get wrong without anybody looking at the result: markup that arrived with the
// artwork, an icon that ignores the colour of the text it sits in, an icon-only
// button that a screen reader has no name for, and a generated file somebody
// edited by hand.
//
// Every test walks the whole set. A sample would pass on the day one file is
// different, which is the only day any of this matters.

// elements returns the elements of a rendered icon, in document order, and
// fails if it is not well-formed XML.
func elements(t *testing.T, name, markup string) []xml.StartElement {
	t.Helper()

	var out []xml.StartElement
	d := xml.NewDecoder(strings.NewReader(markup))
	for {
		tok, err := d.Token()
		if err == io.EOF {
			return out
		}
		if err != nil {
			t.Fatalf("%s does not parse as XML: %v\n%s", name, err, markup)
		}
		if start, ok := tok.(xml.StartElement); ok {
			out = append(out, start)
		}
	}
}

// attr reads one attribute, and reports whether it was written at all. Absent
// and empty are different things: aria-hidden="false" is not the absence of
// aria-hidden, and neither is aria-label="".
func attr(e xml.StartElement, name string) (string, bool) {
	for _, a := range e.Attr {
		if a.Name.Local == name {
			return a.Value, true
		}
	}
	return "", false
}

// TestEveryIconIsOneWellFormedSVG.
//
// An icon is interpolated with {!! !!}, which does not escape, so a malformed
// one is not a malformed icon -- it is the rest of the page reparented into
// whatever tag was left open. The browser recovers silently and the layout is
// wrong somewhere else.
func TestEveryIconIsOneWellFormedSVG(t *testing.T) {
	for name, icon := range all {
		markup := string(icon(icons.Props{}))
		els := elements(t, name, markup)

		if len(els) != 2 {
			t.Errorf("%s drew %d elements, want an svg and a path:\n%s", name, len(els), markup)
			continue
		}
		if els[0].Name.Local != "svg" || els[1].Name.Local != "path" {
			t.Errorf("%s drew <%s><%s>, want <svg><path>", name, els[0].Name.Local, els[1].Name.Local)
			continue
		}
		if box, _ := attr(els[0], "viewBox"); box != "0 0 256 256" {
			t.Errorf("%s has viewBox %q, so it is not on the Phosphor grid", name, box)
		}
		if d, ok := attr(els[1], "d"); !ok || d == "" {
			t.Errorf("%s has no path data", name)
		}
	}
}

// TestNoIconCarriesMarkup.
//
// The generator writes the value of one d attribute and nothing else, so this
// should be true by construction. It is checked anyway because the property is
// the whole reason vendoring somebody else's artwork into this binary is safe:
// the day a Phosphor SVG arrives with a <script>, an onload= or an
// xlink:href=, this is what says so.
func TestNoIconCarriesMarkup(t *testing.T) {
	for name, icon := range all {
		els := elements(t, name, string(icon(icons.Props{})))
		if len(els) < 2 {
			continue // already reported by the well-formed test
		}
		d, _ := attr(els[1], "d")

		for _, forbidden := range []string{"<", ">", "&", "\"", "'", "script", "javascript:", "url("} {
			if strings.Contains(strings.ToLower(d), forbidden) {
				t.Errorf("%s has %q in its path data:\n%s", name, forbidden, d)
			}
		}
		for _, e := range els {
			for _, a := range e.Attr {
				if strings.HasPrefix(strings.ToLower(a.Name.Local), "on") {
					t.Errorf("%s wrote the event handler %s", name, a.Name.Local)
				}
			}
		}
	}
}

// TestEveryIconTakesTheColourOfItsText.
//
// fill="currentColor" is what makes an icon work in a destructive button, in a
// dark theme and in a muted caption without anybody passing it a colour. A
// hard-coded fill is invisible in review -- it looks right in the one place it
// was tried -- and wrong everywhere else.
func TestEveryIconTakesTheColourOfItsText(t *testing.T) {
	for name, icon := range all {
		els := elements(t, name, string(icon(icons.Props{})))
		if len(els) == 0 {
			continue
		}
		if fill, _ := attr(els[0], "fill"); fill != "currentColor" {
			t.Errorf("%s has fill=%q, so it will not follow the text colour", name, fill)
		}
		for _, e := range els[1:] {
			if fill, written := attr(e, "fill"); written {
				t.Errorf("%s overrides the colour on its %s with %q", name, e.Name.Local, fill)
			}
		}
	}
}

// TestSizeComesFromTheLineOrTheStylesheet.
//
// An svg with a viewBox and no size is 300 by 150 pixels, which is the failure
// people meet as "the icon took the whole page". 1em is a presentation
// attribute, so it is the fallback and loses to the stylesheet rule of whatever
// contains it.
func TestSizeComesFromTheLineOrTheStylesheet(t *testing.T) {
	els := elements(t, "House", string(icons.House(icons.Props{})))

	for attribute, want := range map[string]string{"width": "1em", "height": "1em"} {
		if got, _ := attr(els[0], attribute); got != want {
			t.Errorf("%s is %q, want %q", attribute, got, want)
		}
	}
	// A class would be a second place the size is decided, and Basecoat's rules
	// are written to size an svg that has none.
	if class, written := attr(els[0], "class"); written {
		t.Errorf("an icon wrote class=%q; size and colour come from the container", class)
	}
}

// TestAnIconWithNoLabelIsHiddenFromScreenReaders: it sits beside the word that
// already says what it means, and announcing both reads the button twice.
func TestAnIconWithNoLabelIsHiddenFromScreenReaders(t *testing.T) {
	els := elements(t, "Trash", string(icons.Trash(icons.Props{})))

	if hidden, _ := attr(els[0], "aria-hidden"); hidden != "true" {
		t.Errorf("a decorative icon is not aria-hidden: %q", hidden)
	}
	for _, unwanted := range []string{"role", "aria-label"} {
		if v, written := attr(els[0], unwanted); written {
			t.Errorf("a decorative icon wrote %s=%q", unwanted, v)
		}
	}
}

// TestAnIconWithALabelIsAnnounced.
//
// This is the icon-only button. Without a name it is announced as "button", and
// the page offers no way to find out which one.
func TestAnIconWithALabelIsAnnounced(t *testing.T) {
	els := elements(t, "Trash", string(icons.Trash(icons.Props{Label: "Delete this post"})))

	if role, _ := attr(els[0], "role"); role != "img" {
		t.Errorf("a labelled icon has role %q, want img", role)
	}
	if label, _ := attr(els[0], "aria-label"); label != "Delete this post" {
		t.Errorf("aria-label is %q", label)
	}
	// aria-hidden wins over everything under it: an aria-hidden element with a
	// label is an element with no accessible name at all.
	if v, written := attr(els[0], "aria-hidden"); written {
		t.Errorf("a labelled icon is also aria-hidden=%q, so it has no name", v)
	}
}

// TestALabelIsEscaped.
//
// The label is the one thing an icon takes from the caller, and a view
// interpolates the result with {!! !!}. A label that reaches the page unescaped
// makes every icon an injection point, and it does it silently: the icon still
// draws.
func TestALabelIsEscaped(t *testing.T) {
	const payload = `"><script>alert(1)</script>`

	markup := string(icons.Trash(icons.Props{Label: payload}))
	if strings.Contains(markup, "<script>") {
		t.Errorf("a label wrote a script tag through:\n%s", markup)
	}
	els := elements(t, "Trash", markup)
	if label, _ := attr(els[0], "aria-label"); label != payload {
		t.Errorf("the label came back as %q, want it escaped and intact", label)
	}
}

// TestTheGeneratedIconsAreTheOnesInTheTree.
//
// The set is generated from a pinned commit, so the only way a file here can
// stop matching upstream is somebody editing it. That edit compiles.
//
// Two halves. Every generated file has to still name the commit it came from --
// a file with the header removed is a file nobody can regenerate. And the
// exported functions in the package have to be exactly the registry the other
// tests walk, so an icon added by hand is not simply untested.
func TestTheGeneratedIconsAreTheOnesInTheTree(t *testing.T) {
	const pin = "2b75f3ad12b420c9504ef05df8d2564a28f8500e"

	matches, err := filepath.Glob("icons_*.go")
	if err != nil {
		t.Fatal(err)
	}
	var files []string
	for _, file := range matches {
		if !strings.HasSuffix(file, "_test.go") {
			files = append(files, file)
		}
	}
	if len(files) == 0 {
		t.Fatal("no generated files: run `go generate ./...`")
	}

	declared := map[string]bool{}
	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), pin) {
			t.Errorf("%s does not name the Phosphor commit it was generated from", file)
		}

		parsed, err := parser.ParseFile(token.NewFileSet(), file, body, parser.ParseComments|parser.SkipObjectResolution)
		if err != nil {
			t.Fatal(err)
		}
		for _, decl := range parsed.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !fn.Name.IsExported() {
				continue
			}
			declared[fn.Name.Name] = true
			if fn.Doc == nil {
				t.Errorf("%s has no doc comment", fn.Name.Name)
			}
		}
	}

	for name := range declared {
		if _, listed := all[name]; !listed {
			t.Errorf("%s is an icon no test covers: it is missing from all_test.go", name)
		}
	}
	for name := range all {
		if !declared[name] {
			t.Errorf("all_test.go lists %s and no file declares it", name)
		}
	}
	if len(all) != 1512 {
		t.Errorf("the set is %d icons; the pinned Phosphor is 1512", len(all))
	}
}

// TestTheUpstreamLicenceIsHere.
//
// The paths in this directory are somebody else's work under the MIT licence,
// and the one condition it puts on using them is that the notice travels with
// them. A directory of vendored artwork that lost its licence file is a licence
// violation that nothing else in the build would notice.
func TestTheUpstreamLicenceIsHere(t *testing.T) {
	body, err := os.ReadFile("LICENSE.md")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"MIT License", "Copyright (c) 2023 Phosphor Icons"} {
		if !strings.Contains(string(body), want) {
			t.Errorf("LICENSE.md no longer carries %q", want)
		}
	}
}
