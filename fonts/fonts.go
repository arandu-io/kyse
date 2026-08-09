// Package fonts serves the faces a project vendored.
//
// It is the runtime half of `aru font:add`. That command downloads a family, or
// takes a file of your own, writes it under assets/fonts/ and generates the Go
// that calls Register from init(). Everything here is what happens afterwards:
// the bytes reach the served assets, and the layout gets its preload links.
//
// It is in kyse and not in the framework, and that is the point. Which formats
// ship, what a woff2 costs against a ttf, and what a <link rel=preload> has to
// say are typographic decisions -- they belong to the stack that draws the page.
// The framework serves bytes at a content-addressed URL and has no opinion about
// what they are.
//
// # Nothing is fetched, ever
//
// The files are committed to the repository and embedded in the binary. The CSP
// is default-src 'self' with font-src 'self' spelled out beside it, and a font
// from a CDN would mean loosening it in order to look different. There is no
// package.json anywhere near this (RULE 13).
package fonts

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/arandu-io/framework/view"
)

// The formats a vendored face may be in, and what each is served as.
//
// WOFF2 is what ships, and what `aru font:add` produces from the catalogue:
// every browser released since 2016 reads it, and it is Brotli-compressed
// inside the container, so it is already smaller than anything a server could do
// to it.
//
// TrueType and OpenType are accepted for one case, and the CLI says so every
// time it vendors one: a font being DRAWN is a .ttf long before it is a .woff2,
// and refusing it would make the command unusable during exactly the work
// somebody asked it for. They cost two to four times the bytes.
var formats = map[string]string{
	".woff2": "font/woff2",
	".ttf":   "font/ttf",
	".otf":   "font/otf",
}

// Register adds one vendored face to the served assets.
//
//	//go:embed fonts/young-serif-400-latin.woff2
//	var display0 []byte
//
//	func init() { fonts.Register("young-serif-400-latin.woff2", display0) }
//
// That file is generated. Editing it by hand is editing something `aru font:add`
// rewrites, and the command is what keeps it in step with the stylesheet beside
// it -- a face registered and not referenced is bytes in the binary nobody
// draws, and one referenced and not registered is a 404 on every page.
//
// The extension decides the content type, so it is kept: a .ttf renamed .woff2
// is a font every browser fetches and then refuses, with an error that names
// neither the file nor the declaration.
func Register(name string, body []byte) {
	contentType, ok := formats[extension(name)]
	if !ok {
		panic("fonts: " + name + " is not a face this serves -- want .woff2, .ttf or .otf")
	}
	view.RegisterAsset(name, contentType, body)
}

// Preloads is the <link rel=preload> for every registered face.
//
// Without it a font is discovered two round trips deep: the browser fetches the
// page, parses it, fetches the stylesheet, parses THAT, and only then learns
// there is a font -- by which time the heading has already painted in the
// fallback. The preload starts the fetch with the page.
//
// It returns every registered face rather than taking names, so a layout writes
// one line and never touches it again:
//
//	{!! fonts.Preloads() !!}
//
// Swapping the family is `aru font:add` running, not a template being edited.
//
// crossorigin is required even same-origin. A font is fetched in CORS mode by
// specification, and a preload without it is a SECOND request rather than a warm
// cache -- the most common way a preload makes a page slower instead of faster.
//
// The type is what the response will actually carry. A preload whose type does
// not match is the other way to make it a second request.
func Preloads() template.HTML {
	var b strings.Builder
	for _, a := range view.Assets() {
		if _, ok := formats[extension(a.Name)]; !ok {
			continue
		}
		fmt.Fprintf(&b, `<link rel="preload" href="%s" as="font" type="%s" crossorigin>`,
			a.URL, a.ContentType)
	}
	return template.HTML(b.String())
}

// Registered is the vendored faces, for a health check or a debug page.
//
// It answers the question `aru font:list` answers from the manifest, but from
// the running binary -- which is the one that matters when a page is drawing in
// the fallback and nobody knows why.
func Registered() []view.Asset {
	var out []view.Asset
	for _, a := range view.Assets() {
		if _, ok := formats[extension(a.Name)]; ok {
			out = append(out, a)
		}
	}
	return out
}

// extension is the lowercased suffix from the last dot.
//
// filepath.Ext is not used: this is a URL path segment and not a file path, and
// on Windows filepath would treat a backslash as a separator in something that
// never has one.
func extension(name string) string {
	at := strings.LastIndexByte(name, '.')
	if at < 0 {
		return ""
	}
	return strings.ToLower(name[at:])
}
