// Package kyse holds what a caller writes in a view and a component reads,
// where the two need a type and neither owns it.
//
// It is one type today. The components live under components/, and what is here
// is what they take rather than what they are.
package kyse

import "github.com/arandu-io/hesape/view"

// CSS is a block of stylesheet scoped to one component instance.
//
//	components.Card(components.CardProps{
//		ComponentProps: components.ComponentProps{
//			Style: kyse.CSS(`
//				& { gap: 6px; }
//				& [data-part="title"] { letter-spacing: -0.01em; }
//			`),
//		},
//	})
//
// # Where the block goes, and why not into the page
//
// It does not travel in the page. The policy is style-src 'self' with no
// unsafe-inline, so neither a style attribute nor a style element survives, and
// both fail by being dropped rather than by saying anything -- the markup reads
// correctly and the browser ignores it.
//
// What travels is a class. `aru view:build` reads these blocks out of the
// source, replaces & with the class, and writes the rules into the project's
// stylesheet, which is served from the origin like every other asset. The class
// is the hash of the block's own text, so the build and the render agree without
// a table between them: neither knows about the other, and both are looking at
// the same bytes.
//
// # The one condition
//
// The argument has to be a string literal. A block built at run time is a hash
// the render can compute and the build never saw, so the rule is never emitted
// and the element carries a class nothing styles -- which renders, and is
// silent. `aru view:build` refuses a call whose argument is anything else, and
// names the file and the line.
//
// It is the same limit a progress bar's width is written around, for the same
// reason: what the stylesheet contains is decided at build time, from the text
// in the source.
type CSS string

// Class is the class this block is compiled under, or empty when there is no
// block.
func (c CSS) Class() string {
	if c == "" {
		return ""
	}
	return view.StyleClass(string(c))
}

// Text is the block as written, which is what the build reads and hashes.
func (c CSS) Text() string { return string(c) }
