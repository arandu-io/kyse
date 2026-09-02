# Upgrade guide

What changed in a way that stops your code compiling, and what to write instead.

Additions are not listed. A new component or a new prop breaks nobody, and a
file that recorded every one of them would be a changelog nobody reads to find
the two lines that matter.

## Before v1.0.0

While the version starts with `v0.`, the API can break. That is what `v0.` means
in Go and it is deliberate — the alternative is freezing a shape before anyone
has drawn a screen with it. What is not deliberate is breaking it quietly, which
is what this file exists to stop.

This module is imported and never copied. A project draws a button by requiring
a version of it, so a prop that stops existing is a build that stops in somebody
else's repository, at the moment they upgrade and not before. Every release from
here is compared against the one before it, package by package, by `apidiff` in
CI, and an incompatible change with no entry here fails the build.

---

## Unreleased — Dialog submissions use the browser's native transport

`DialogProps.Method` may still be PUT, PATCH or DELETE, but an HTML form cannot
send any of those methods. The rendered form now sends POST and carries the
intended method in a hidden `_method` field. Applications must register
`hesape/http/middleware.OverrideMethod` **after** Framework's `CSRFProtect` and
before routing; the Arandu skeleton does this in its default pipeline.

The dialog's CSRF field is `_token`, and this module now requires Framework
v0.42.0, the release whose `CSRFProtect` reads that spelling. A public test
renders the component, simulates the browser submission, and crosses
`CSRFProtect → OverrideMethod → Router` for all three spoofed methods.

### v0.15.0 and v0.15.1 are retracted

Both releases selected Framework v0.41.0 together with Hesape v0.21.0. That
pair disagrees about the form field: Framework reads `_csrf`, while Hesape's
`@csrf` directive writes `_token`. v0.15.1 also changed Dialog itself to
`_token` without raising the Framework floor, so its own confirmation form was
refused with 419. In both releases, a Dialog configured for PUT, PATCH or DELETE
also wrote that method directly on the form, which browsers do not submit.

Use the next Kyse release rather than either retracted version. Retraction is
metadata published by that next release; it cannot remove an immutable tag, but
the Go command will warn before selecting one of them.

## v0.15.0 — every component takes a class, attributes and parts

### `ThemeToggle` takes props

```
- ./components.ThemeToggle: changed from func() html/template.HTML to func(ThemeToggleProps) html/template.HTML
```

```go
{!! components.ThemeToggle() !!}                              // before
{!! components.ThemeToggle(components.ThemeToggleProps{}) !!} // now
```

It was the one component that could not take a class, an attribute or a part,
and "every component except that one" is a rule nobody remembers. It carries no
data of its own and still needs somewhere for a caller to write.

### Twenty-two props types are no longer comparable

```
- ./components.AlertProps: old is comparable, new is not
- ./components.AvatarProps: old is comparable, new is not
- ./components.BadgeProps: old is comparable, new is not
- ./components.ButtonProps: old is comparable, new is not
- ./components.CardProps: old is comparable, new is not
- ./components.CheckboxProps: old is comparable, new is not
- ./components.CollapsibleProps: old is comparable, new is not
- ./components.DialogProps: old is comparable, new is not
- ./components.EmptyProps: old is comparable, new is not
- ./components.FieldProps: old is comparable, new is not
- ./components.InputGroupProps: old is comparable, new is not
- ./components.InputProps: old is comparable, new is not
- ./components.ItemProps: old is comparable, new is not
- ./components.LabelProps: old is comparable, new is not
- ./components.PopoverProps: old is comparable, new is not
- ./components.ProgressProps: old is comparable, new is not
- ./components.RangeSliderProps: old is comparable, new is not
- ./components.SeparatorProps: old is comparable, new is not
- ./components.SkeletonProps: old is comparable, new is not
- ./components.SwitchProps: old is comparable, new is not
- ./components.TextareaProps: old is comparable, new is not
- ./components.ToastProps: old is comparable, new is not
```

Every props type now embeds `ComponentProps`, which holds two maps, and a struct
containing a map cannot be compared with `==` or used as a map key.

Nothing in this repository or in the Arandu tree did either, so in practice this
breaks a call nobody makes. If you compared two props values, compare the fields
you meant instead.

### What you get for it

Every component takes `ComponentProps` as an embedded field:

```go
{!! components.ButtonGroup(components.ButtonGroupProps{
    Label: "Message actions",
    ComponentProps: components.ComponentProps{
        Class: "w-full",
        Parts: components.Parts{
            "root":   {Class: "gap-2"},
            "button": {Class: "rounded-xl"},
        },
        Attrs: components.Attrs{"data-testid": "message-actions"},
    },
    Buttons: []components.ButtonProps{{Label: "Archive"}},
}) !!}
```

Each component publishes the parts it draws — `PartNames()` on its props, and
the same names as `data-part` on the elements — so a caller addresses an inner
element by name rather than by a selector written against the structure.

Classes are **added** to the component's own, not substituted. Order inside the
attribute decides nothing; what settles a collision is the layer, and a Tailwind
utility beats a component class there. Utility against utility is the one case
that needs the importance marker: `p-8!`.

Attributes go through a check that refuses every name whose value a browser or a
library would act on — `on*`, `hx-on*`, `x-*`, `style`, `href` and the other
address attributes, and the `class`, `role` and `aria-*` the component owns. A
refused name is an error at render, naming the attribute and the view's line.

### A class written outside a view is not compiled into the stylesheet

Tailwind reads the source, so a class only exists if its name appears in a file
the stylesheet is compiled from. Those are the project's `resources/views/**`,
its compiled output, and this module.

A class written in a `.kyse.go` is therefore found. A class written in ordinary
application Go — a catalogue, a helper that returns props — **is not**, and the
element renders unstyled. So is anything assembled at run time, `"text-"+variant`
included. It is the same limit the progress bar's width table is written around.

### Fixed on the way past

The range slider wrote `class` twice on its `<input>` — `class="input"` and,
lower down, `class="{{ .FillClass() }}"`. HTML keeps the first occurrence of a
repeated attribute and drops the rest, so the fill class was in the markup and
never on the page. Both now go through one call.

## v0.13.0 — the range slider paints with a class

```
- ./components: RangeSliderProps.Fill: removed
```

`Fill()` returned how much of the track is coloured as a percentage string, meant
for an inline `style`. It is gone, and `FillClass()` returns the class that
paints the same thing:

```go
// before
<div style="width: {{ .Fill }}">

// after
<div class="{{ .FillClass }}">
```

The reason is the content security policy this stack serves under. An inline
`style` needs `style-src` to allow unsafe inline, and it does not -- so the
percentage was being computed, written into the markup, and refused by the
browser.

The class comes out of twenty-one fixed positions, so the painted half moves in
steps of five percent until the pointer takes over. What is announced to a screen
reader, and what is submitted with the form, is the exact value and not the
rounded one.

## v0.6.0 — an input asks the page instead of being handed the message

One change, and it is the only incompatible one this module has shipped. It cost
two projects a build each, which is what the rest of this file is for.

### `FieldProps.Error` and `TextareaProps.Error` are gone

| was | is | what to do |
|---|---|---|
| `FieldProps.Error string` | `FieldProps.Page Page` | Pass the screen instead of the message. The component looks the message up by the `Name` it was already given |
| `TextareaProps.Error string` | `TextareaProps.Page Page` | Same |

In a view, the change is one line per field:

```go
{!! components.Field(components.FieldProps{
	Name: "email", Label: "Email", Type: "email",
	Value: .Email, Error: .EmailError,
	Autocomplete: "email", Required: true,
}) !!}
```

becomes

```go
{!! components.Field(components.FieldProps{
	Name: "email", Label: "Email", Type: "email",
	Value: .Email, Page: .,
	Autocomplete: "email", Required: true,
}) !!}
```

`Page: .` passes the screen itself, which is what a screen embedding
`view.Page` can do. A screen that holds its page data in a field passes
`Page: .Page`.

**Why the field went rather than gaining a sibling.** `Error: .FieldError("email")`
writes the field's name twice in one call, and nothing checked that the two
agreed. When they did not, nothing said so: an input labelled "Confirm password"
drew the message belonging to `password`, or drew none, on the one screen where a
missing message is the entire complaint. Asked this way the name is written once.

### What `Page` accepts

Anything with the two methods:

```go
FieldError(name string) string
OldOr(name, fallback string) string
```

`github.com/arandu-io/framework/view.Page` has both, so a screen that embeds it
needs nothing else.

A screen that keeps its messages in **named fields** rather than in the map
`view.Page` carries has to write `FieldError` itself, mapping the form field name
onto the field a handler filled. Both projects that broke did exactly this, and
kept the typed fields as the source — a handler that sets
`PasswordConfirmatonError` does not build, where a map key spelt the same way is
simply never read:

```go
func (p AuthPage) FieldError(name string) string {
	switch name {
	case "name":
		return p.NameError
	case "email":
		return p.EmailError
	case "password":
		return p.PasswordError
	case "password_confirmation":
		return p.PasswordConfirmationError
	}
	return p.Page.FieldError(name)
}
```

The last line matters: a name with no field of its own falls through to
`view.Page`, which is where anything generic puts a message.

A `nil` `Page` is the ordinary case for a screen with no form behind it. It draws
no message and keeps whatever `Value` it was given, rather than panicking.

### One thing that still compiles and does not mean the same

`Value` is now the fallback rather than the answer. With a `Page` set, the input
is drawn with what was typed on the rejected attempt and falls back to `Value`
when there was none — so a rejected edit form comes back carrying the change
somebody was in the middle of making instead of reverting to the stored row.
That is the behaviour a form wants, and it arrives without a line changing, which
is why it is written down here.

### What it cost

Recorded because it is the argument for this file existing:

- `arandu-io/examples` — every form in the application stopped compiling.
  Thirty-one call sites across eleven views. No call there had a `Name` that
  disagreed with the message it passed, which was checked before they were
  changed.
- `arandu-io/ui` — `go run github.com/arandu-io/ui auth` installed a project
  that did not build. The nine published views still passed `Error`, and the
  tests compared published text against golden files without ever compiling it.

Neither found out until they upgraded.

---

## Everything else published so far

Measured, not remembered: `apidiff` between every pair of consecutive tags from
`v0.1.0` to `v0.9.0`, over every package of the module. Nothing else in that
range is incompatible.

Two results from that sweep are worth writing down, because both look like
entries and neither is one:

- **`v0.6.0` to `v0.7.0`, `fonts.Registered`.** `apidiff` reports it as changed
  from `func() []view.Asset` to `func() []view.Asset` — the same text on both
  sides. The signature did not move; the version of the framework underneath did,
  which gives the type a new identity. The CI step drops a line whose two halves
  are identical for this reason, and keeps every line where they differ.
- **`v0.4.0`, `fonts`.** It does not load at all: it calls `view.RegisterAsset`
  against a framework tag that predates it. So there is nothing to compare it
  against, and no measurement of that release exists here. `v0.4.1` is the fix,
  and it is one commit later.

## Deprecation

Nothing is deprecated right now, and that is a statement rather than an
omission: the entry above is what has already broken, and no symbol is currently
on its way out.

When one is, it goes through three steps and this file records each:

1. `// Deprecated:` on the symbol, naming what replaces it
2. an entry here, with the before and after
3. removal, never in the same release as step 1

## How this is checked

`.github/workflows/ci.yml` runs `apidiff` on every pull request, comparing the
working tree against the newest tag, package by package. Incompatible changes are
printed in the build log, and the build fails when there are some and this file
has no entry for them.

The point is not to prevent the break. It is to make the break something
somebody decided, in a diff a reviewer can see, instead of something a person
discovers when their build stops.

No view build runs in that step, and none is needed. The compiled components are
committed, and `go get` has no build step, so the `.go` in the tree is the
surface every consumer imports — on both sides of the comparison.
