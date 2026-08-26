package unit

import (
	"strings"
	"testing"

	"github.com/arandu-io/kyse/fonts"
)

// A woff2 header, which is all Register needs: it reads the extension and hands
// the bytes to the asset table.
var body = []byte("wOF2 not really, but not empty either")

// TestRegisterKeepsTheFormatTheFileArrivedIn.
//
// The extension decides the content type. A .ttf served as font/woff2 is a font
// every browser fetches and then refuses, with an error that names neither the
// file nor the declaration -- and the page draws in the fallback with nothing
// saying why.
func TestRegisterKeepsTheFormatTheFileArrivedIn(t *testing.T) {
	fonts.Register("a-400.woff2", body)
	fonts.Register("b-400.ttf", body)
	fonts.Register("c-400.otf", body)

	want := map[string]string{
		"a-400.woff2": "font/woff2",
		"b-400.ttf":   "font/ttf",
		"c-400.otf":   "font/otf",
	}
	for _, a := range fonts.Registered() {
		if expected, ok := want[a.Name]; ok && a.ContentType != expected {
			t.Errorf("%s is served as %s, want %s", a.Name, a.ContentType, expected)
		}
	}
}

// TestSomethingThatIsNotAFaceIsRefusedAtBoot.
//
// It is a mistake in a generated file, so it belongs at boot rather than as a
// 404 on every page of a running application.
func TestSomethingThatIsNotAFaceIsRefusedAtBoot(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a .zip was accepted as a face")
		}
	}()
	fonts.Register("archive-400.zip", body)
}

// TestPreloadsCarryCrossoriginAndTheRightType.
//
// A font is fetched in CORS mode by specification, so a preload without
// crossorigin is a SECOND request rather than a warm cache -- the most common
// way a preload makes a page slower. A type that does not match the response is
// the other way.
func TestPreloadsCarryCrossoriginAndTheRightType(t *testing.T) {
	fonts.Register("preload-400.woff2", body)

	markup := string(fonts.Preloads())
	if !strings.Contains(markup, "preload-400.woff2") {
		t.Fatalf("the registered face is not preloaded:\n%s", markup)
	}
	for _, want := range []string{`rel="preload"`, `as="font"`, `type="font/woff2"`, "crossorigin"} {
		if !strings.Contains(markup, want) {
			t.Errorf("the link is missing %s:\n%s", want, markup)
		}
	}
	// The stylesheet and the scripts are registered too, and preloading a
	// stylesheet as a font is a request the browser makes and discards.
	if strings.Contains(markup, ".css") || strings.Contains(markup, ".js") {
		t.Error("something that is not a face was preloaded as one")
	}
}

// TestAFontNameCannotAddAPreloadAttribute covers the application-controlled
// part of the preload URL. A quote is a valid file-name byte, but it must stay
// inside href rather than becoming an event-handler attribute.
func TestAFontNameCannotAddAPreloadAttribute(t *testing.T) {
	fonts.Register(`hostile" onload="alert(1).woff2`, body)

	markup := string(fonts.Preloads())
	if strings.Contains(markup, `" onload="`) {
		t.Fatalf("a font name added an event-handler attribute:\n%s", markup)
	}
	if !strings.Contains(markup, `hostile&#34; onload=&#34;alert(1).woff2`) {
		t.Errorf("the hostile name was not kept as escaped href data:\n%s", markup)
	}
}

// TestTheOrderIsStable: a map has no order, and markup that differs by line
// order between two runs of one binary is markup no diff and no cache can
// compare.
func TestTheOrderIsStable(t *testing.T) {
	fonts.Register("zzz-400.woff2", body)
	fonts.Register("aaa-400.woff2", body)

	first := string(fonts.Preloads())
	for range 5 {
		if got := string(fonts.Preloads()); got != first {
			t.Fatal("two calls produced different markup")
		}
	}
	if strings.Index(first, "aaa-400") > strings.Index(first, "zzz-400") {
		t.Error("the order is not by name")
	}
}
