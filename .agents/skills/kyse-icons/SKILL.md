---
name: kyse-icons
description: Draw, find, add or regenerate an icon in the kyse icons package (Go). Use when the request is to "add an icon", "put a trash icon on the button", "which icon is that", "the icon is the wrong size", "make the icon red", "add a new icon to the set", "update Phosphor", or when a build says undefined icons.Something. Also use when tempted to paste an SVG into a component, to write a Size or Color prop, to embed the icon set with go:embed, or to reach an icon by name from a map. Covers the 1512 exported functions, how to find the one you want, the single Label prop and what it does to accessibility, and the pinned generator that is the only way the package changes.
license: MIT
---

# Icons

Phosphor, regular weight, 1512 of them, one exported function each:

```go
{!! icons.Trash(icons.Props{Label: "Delete this post"}) !!}
```

Count them with `grep -hcE '^func [A-Z]' icons/icons_*.go | paste -sd+ - | bc`.

## Why one function each

An application uses thirty icons out of six hundred kilobytes of path data.
`go:embed`, or a map keyed by name, puts all of it in every binary that draws
one arrow — 1.2 MB, measured — because the linker cannot prove a map lookup will
not ask for the rest. It can prove that about a function: an icon nobody calls
is dead code, and dead code and its string data are dropped. A program that
draws one icon carries 176 bytes of the package.

The second reason is the same one the components package makes. `icons.Trahs`
is undefined at the line of the file that wrote it. An icon named by string is a
blank square in production.

So: **never reach an icon through a variable, a map or a string.** If you find
yourself wanting to, the shape is wrong — pass the rendered `template.HTML`, or
branch on the call.

## Finding the one you want

The upstream synonyms are in each doc comment, so the name is greppable by what
you would have called it. The comment wraps onto a second line, so search with
context and keep the line that names the function:

```sh
grep -rn -i -B2 'zoom' icons/icons_*.go | grep 'is the Phosphor'
```

It lists candidates rather than an answer — `zoom` returns
`MagnifyingGlassPlus`, `MagnifyingGlassMinus` and `VideoConference` — and you
pick. The same comments are the search index on pkg.go.dev.

Names are the Phosphor names in PascalCase: `arrow-right` is `ArrowRight`,
`magnifying-glass` is `MagnifyingGlass`.

## The one prop

```go
type Props struct {
	Label string
}
```

That is the whole API, and the field is the one an icon gets wrong.

- **Empty means decorative.** The icon is marked `aria-hidden="true"` and gets
  no role and no label. This is the common case: an icon beside the word
  "Delete" is drawn for the eye, and announcing "trash, Delete" reads the button
  twice.
- **Set means the icon is the name.** It gets `role="img"` and the label, and no
  `aria-hidden` — an `aria-hidden` element with a label has no accessible name
  at all. This is the icon-only button, which without a label is announced as
  "button" with nothing on the screen saying which one.

**There is no Size and no Color, and adding one is not a small convenience.**
The svg is `fill="currentColor"` at `width="1em" height="1em"`. Both are
presentation attributes, so both lose to any CSS rule, which is what makes the
stylesheet the one place an icon's size is decided: Basecoat sizes svg
descendants of a button, an alert or a badge. A prop here would be a second
place, and the two would disagree.

There is no `Class` either, for the same reason.
`TestSizeComesFromTheLineOrTheStylesheet` at `tests/Unit/icons_test.go:153`
fails if the renderer stops writing `1em`, or starts writing a class.

To make an icon larger or a different colour, style the element that contains
it — that is what `class="block size-4"` around an icon in `ThemeToggle` is
doing.

## Changing the set

`icons/*.go` is generated and the pin is in the source:
`internal/icongen/main.go` names one commit of `phosphor-icons/core`. Two runs a
year apart produce the same bytes, and moving to a newer Phosphor is an edit
somebody makes on purpose and reviews as a diff.

```sh
export GOWORK=off
go generate ./...
```

It prints `icongen: wrote 1512 icons to …/icons` and leaves the tree unchanged
when nothing moved. It also rewrites `tests/Unit/icons_all_test.go`, the
registry of every icon that the tests walk — that file is generated too, and it
lives under `tests/` rather than in the package because a map of every icon in
the package proper would pin all of them into every binary that draws one.

**Never hand-edit a file under `icons/`.** That edit compiles, and the next
regeneration silently undoes it. Two things catch it:
`TestTheGeneratedIconsAreTheOnesInTheTree` at `tests/Unit/icons_test.go:231`
fails on a generated file that no longer names the pinned commit, on an exported
function missing from the registry, and on a set whose size is not 1512; and CI
regenerates and fails if the tree moved.

To add an icon that Phosphor has and this package does not, move the pin. To add
one Phosphor does not have, do not — the set is a design decision the library
makes once, like the stylesheet, and a hand-written icon is the first of a
second set.

Only the regular weight ships. Phosphor draws six, and shipping two would mean
`Trash` and `TrashBold`: two names for one idea, and a question on every call
site that has no good answer.

## What the generator refuses, and why it matters

`icongen` carries out the value of one `d` attribute and nothing else. Every
regular Phosphor SVG is one `<path>` inside a wrapper the program knows by
heart, and it refuses anything else — another element, another attribute, a
different wrapper. A vendoring step that copied whole SVG files would one day
copy a `<script>` or an `onload=` with them.

That is also what makes it safe for a view to interpolate an icon with `{!! !!}`.
Three tests in `tests/Unit/icons_test.go` hold every one of the 1512 to it on
every run:

| test | line | what it holds |
| --- | --- | --- |
| `TestEveryIconIsOneWellFormedSVG` | 72 | one `<svg>`, one `<path>`, the Phosphor viewBox, non-empty path data |
| `TestNoIconCarriesMarkup` | 101 | no angle bracket, quote, `script`, `javascript:` or `url(` in the path data, and no `on*` attribute anywhere |
| `TestEveryIconTakesTheColourOfItsText` | 130 | `fill="currentColor"` on the svg and on nothing under it |

They walk the whole set rather than a sample, because a sample would pass on the
day one file is different, which is the only day any of this matters.

`TestALabelIsEscaped` at line 209 reads one icon, and that is enough: the label
is escaped in the shared renderer that all 1512 call, so one of them proves it
for the set.

## The licence travels with the artwork

The paths are somebody else's work under MIT, and the one condition is that the
notice travels with them. It is in `icons/LICENSE.md`, beside the files, and
`TestTheUpstreamLicenceIsHere` at `tests/Unit/icons_test.go:295` fails if it
stops carrying the copyright line. A directory of vendored artwork that lost its
licence file is a violation nothing else in the build would notice.
