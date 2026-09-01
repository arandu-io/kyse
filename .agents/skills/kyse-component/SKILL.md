---
name: kyse-component
description: Add, change or remove a component in the kyse component library (Go). Use when the request is to "add a component", "make a button", "write a new card", "add a prop", "rename a field", "port a shadcn component", "the component does not compile", "undefined components.X", "unknown field in struct literal", or when a pull request against this repository touches anything under components/. Covers the .kyse.go source, the Props struct the compiler turns into the function parameter, aru view:build and the compiled .go that is committed, the row every component needs in the smoke table, and what happens when a prop that already shipped stops existing.
license: MIT
---

# Adding a component

A component is one file you write and one file the compiler writes. Both are
committed, and the suite compares them against each other and against the
directory they sit in.

## What you are producing

`components/stat-card.kyse.go` is the source. `aru view:build` turns it into
`components/stat-card.go`, which declares:

```go
func StatCard(p StatCardProps) template.HTML
```

Two things decide that signature and neither is the struct name:

- **The function name comes from the file name**, with the dashes taken out.
  `stat-card.kyse.go` is `StatCard`, `dropdown-menu.kyse.go` is `DropdownMenu`.
- **The parameter type is the first struct declared in the `@go` block.** Name
  it `<Component>Props` — every component here does, and the smoke table reads
  as nonsense otherwise. A component whose `@go` block declares no struct takes
  no arguments at all; `ThemeToggle()` is the one that does this, because the
  theme is client state and there is nothing for a caller to pass.

## The procedure

**1. Write the source.** The build tag is what keeps the compiler out of the
file: everything after it is Go the compiler would reject, and the tag is why it
never sees it.

```go
//go:build kyse

package components

@go
// StatCardProps is a table of numbers with a heading over it.
type StatCardProps struct {
	// Title is the heading over the table.
	Title string
	// Rows are the lines. Empty draws Empty rather than a headed table with
	// nothing under it.
	Rows []StatRow
	// Empty is what stands in when there are no rows.
	Empty EmptyProps
}
@endgo

<article class="card p-5">
	<h3 class="font-semibold tracking-tight">{{ .Title }}</h3>
	...
</article>
```

Read `kyse-markup` before writing the part after `@endgo`. What may appear
there is a short list, and most of what a component library usually reaches for
is not on it.

**2. Compile it, from the root of this module.**

```sh
export GOWORK=off
aru view:build
```

It reports how many views it compiled — `kyse: 37 view(s) compiled` before you
added yours — and writes the `.go` beside each source. **Commit what it
writes.** `go get` has no build step to run, so the compiled
files in the repository are the ones every consumer imports; a missing one is a
module that does not build for them, and a stale one ships a component whose
source says something else. `TestTheGeneratedFilesMatchTheirSources` in
`tests/Feature/components_test.go:32` fails on a `.go` with no source, and CI
rebuilds everything and fails if the tree moved.

**3. Add the row to the smoke table.** `TestEveryComponentRenders` in
`tests/Unit/components_test.go:30` is a table, one line per component, each one
calling it with the smallest props that make sense:

```go
{"StatCard", string(components.StatCard(components.StatCardProps{Title: "Connections"}))},
```

This is not optional and it is not on your honour.
`TestEveryComponentIsInTheTable` at `tests/Unit/components_test.go:417` globs
`components/*.kyse.go`, turns each file name into the function name, and fails
when the table has no row starting `{"StatCard",`. The library went from twelve
components to thirty-seven in one change and the table stayed green at sixteen,
because nothing compared it to anything. Now something does.

The row also asserts three things about what came back: it is not empty, it
contains a `<`, and it contains a `class="`. The last one is why a component
with no class fails here rather than looking unstyled in somebody's project.

**4. Run the gates.**

```sh
export GOWORK=off
aru view:build && gofmt -l $(find . -name '*.go' -not -path '*/testdata/*' -not -name '*.kyse.go') \
  && go build ./... && go vet ./... && go test -race ./... && bash tests/test-layout-guard.sh
```

## Where a test goes

Under `tests/`, in a capitalized category directory declaring a lowercase
package: `tests/Unit` holds `package unit`. A test that genuinely needs
something the package does not export goes beside the code as
`*_internal_test.go`, and the suffix is how it says so.
`tests/test-layout-guard.sh` runs in CI and rejects anything else — including a
file named `ButtonTest.go`, which compiles as ordinary code and runs no tests at
all.

Prefer the first. A test under `tests/` sees what a caller sees, which for a
component library is the whole point.

## Props conventions the existing components keep

- **Every field carries a doc comment.** The reference on pkg.go.dev is
  generated from them, and a prop with no comment is a prop nobody can use
  without reading the markup.
- **The zero value is the sensible default, and the component says so.**
  `Variant` empty is the default variant; `Type` empty is `"button"`, because a
  button inside a form that does not say otherwise submits it.
- **An optional attribute is written only when it carries something.**
  `hx-post=""` is not the absence of `hx-post` — HTMX acts on it and posts to
  the current URL. `TestAnEmptyOptionalAttributeIsNotWritten` at
  `tests/Unit/components_test.go:174` holds `Button` to this.
- **An input asks the page rather than being handed its message.** Take
  `Page Page` and look the message and the typed value up by the `Name` the
  component already has. The alternative writes the field's name three times in
  one call and nothing checks that the three agree.
- **Derived values are methods on the props struct**, in the `@go` block:
  `InputType()`, `DescribedBy()`, `AlignClass()`. The markup stays readable and
  the logic is testable without rendering.
- **`ComponentProps` is the first field, and `PartNames()` is a method beside
  the others.** They are what lets a caller add a class, an attribute or a part
  without forking the component. The doc comment on the struct lists the parts
  it publishes, because that list is the component's public promise and
  pkg.go.dev is where somebody reads it.
- **A row in `tests/Unit/component-parts_test.go`.** Not optional:
  `TestEveryExtensibleComponentIsInThisTable` reads the directory and fails on a
  component that embeds the type and has no row. If the component draws a part
  only in some states — a message or a hint, never both — the row answers one
  rendering per state.
- **Nothing is a part unless a caller would reach for it.** A group whose
  members are themselves components does not publish them: the button inside a
  `ButtonGroup` is a `ButtonProps` and carries its own class, so a part here
  would be a second way to the same element.

## Removing or renaming a prop

This module is imported and never copied, so a field that stops existing is a
build that stops in somebody else's repository, weeks later, when they upgrade.

Breaking is allowed while the version starts with `v0.` — freezing a shape this
early is worse. Breaking *quietly* is not. CI compares every package against the
last release with `apidiff`, and an incompatible change with no entry in
`UPGRADE.md` fails the build. Write the entry as a table: what it was, what it
is, and what somebody has to write instead. Additions are not listed there; a
new component or a new prop breaks nobody.

## What a component may not take on

`arandu.mod.toml` declares this library's permissions, and all four are
`false`: network, filesystem, exec, migrations. A component renders markup. If
the thing you are building needs to fetch, read a file, run a program or own a
table, it is not a component, and the answer is not to widen the library.
