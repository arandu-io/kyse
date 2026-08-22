<p align="center">
  <img src=".github/logo.png" alt="Arandu" width="140" height="140">
</p>

<h1 align="center">arandu-io/kyse</h1>

<p align="center">The component library kyse views import to draw a screen.</p>

<p align="center">
<a href="https://github.com/arandu-io/kyse/actions/workflows/ci.yml"><img src="https://github.com/arandu-io/kyse/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
<a href="https://pkg.go.dev/github.com/arandu-io/kyse/components"><img src="https://pkg.go.dev/badge/github.com/arandu-io/kyse.svg" alt="Go Reference"></a>
<a href="https://github.com/arandu-io/kyse/tags"><img src="https://img.shields.io/github/v/tag/arandu-io/kyse?label=version" alt="Latest Version"></a>
<a href="LICENSE.md"><img src="https://img.shields.io/github/license/arandu-io/kyse" alt="License"></a>
</p>

## What this is

kyse names two things, and the name is overloaded on purpose. One is the
template language an Arandu view is written in: a `.kyse.go` file, compiled to
plain Go by the compiler in `aru/internal/kyse` — a build step, run by
`aru view:build`, never a runtime a request passes through. The other is this
repository: the library of components, written in that language, that a view
imports to draw a screen.

The components themselves are adapted from shadcn/ui by way of shadcn-htmx, and
importing the second is how a view writes the first.

## Install

```sh
go get github.com/arandu-io/kyse/components
```

That is the whole installation. Nothing is copied into your project, there is no
directory to keep in sync, and `aru view:build` has nothing extra to compile.

## Use

```go
//go:build kyse

package posts

import "github.com/arandu-io/kyse/components"

@extends('layouts.app')

@section('content')
	<form method="post" action="{{ .SubmitURL }}">
		@csrf
		{!! components.Field(components.FieldProps{
			Name:     "title",
			Label:    "Title",
			Value:    .Post.Title,
			Page:     .,
			Required: true,
		}) !!}

		{!! components.Button(components.ButtonProps{
			Label: "Publish",
			Type:  "submit",
		}) !!}
	</form>
@endsection
```

Every component is an ordinary exported Go function that returns
`template.HTML`. That is the point of it: a component that does not exist is
`undefined: components.Buton` and a prop that does not exist is `unknown field
Labl`, both at the line of the `.kyse.go` you wrote. A component library
resolved by string at run time would report the same two mistakes as a blank
space and a 500.

## What is here

| | |
|---|---|
| `Alert` | a message about what happened, or is about to |
| `Avatar` | who somebody is, in a circle, falling back to initials |
| `Badge` | a small piece of status beside something else |
| `Button` | six variants, five sizes, and the HTMX attributes |
| `Card` | one item in a list: a title, a sentence, where it goes |
| `Dialog` | a confirmation, in a `<dialog>` the browser manages |
| `Empty` | what a list draws when it has nothing in it |
| `Field` | a labelled input, with its error and its hint |
| `StatCard` | a table of numbers with a heading over it |
| `Textarea` | the same, for a textarea |
| `ThemeToggle` | light, dark and six accents, kept on the device |
| `Toast` | a flash message HTMX appends to the tray |

Twelve. What is here is what has been adapted, not what is planned.

Two more packages travel in the same module and neither draws a page. `mailui`
draws the messages an application sends, in tables and inline attributes,
because a stylesheet does not survive a mail client. `fonts` serves the faces a
project vendored and writes the preload links for them.

## Icons

[Phosphor](https://github.com/phosphor-icons/core) (MIT), regular weight, as
1512 exported functions in `github.com/arandu-io/kyse/icons`:

```go
{!! icons.Trash(icons.Props{Label: "Delete this post"}) !!}
```

One function per icon, rather than `go:embed` or a map keyed by name, because
the linker can prove a function is unreachable and cannot prove that about a map
lookup. Measured on a program that draws one icon: 176 bytes. The same icon
reached through a name-keyed map costs 1.24 MB, because the map keeps all 1512
alive. See [ADR 0033](https://github.com/arandu-io/docs).

No `Label` means the icon is decorative and is marked `aria-hidden`; a `Label`
makes it the accessible name, which is what an icon-only button needs. There is
no size and no colour prop: the svg is `currentColor` at `1em`, and both lose to
the stylesheet rule of whatever contains it.

The paths are regenerated from a pinned commit with `go generate ./...`, and
Phosphor's licence sits in `icons/LICENSE.md` beside them.

## The stylesheet

The markup carries semantic class names — `btn`, `field`, `card` — and the rules
behind them ship with the skeleton, under `resources/css/basecoat/`. A project
created with `aru new` already has them.

That split is deliberate. The same button written with utility classes is 443
characters of markup; written as `class="btn"` it is thirty. The rules exist once
in a stylesheet you own and can edit, rather than once per element in every view
that draws a button.

## Working on a component

The sources are the `.kyse.go` files under `components/`, and the `.go` beside
each one is compiled output that is committed — a module whose generated files
are missing is a module `go get` cannot use.

```sh
aru view:build
```

Run it from the root of this module. `aru view:build` compiles the views of
whatever module it is in: a project keeps them under `resources/views`, a library
keeps them at its root, and the command looks rather than being told. Commit what
it writes.

## Where this came from

The component set, the `data-variant` / `data-size` attribute API and the HTMX
patterns are adapted from [shadcn-htmx](https://github.com/productdevbook/shadcn-htmx)
(MIT). The stylesheet is [Basecoat](https://github.com/hunvreus/basecoat) (MIT).
The icons are [Phosphor](https://github.com/phosphor-icons/core) (MIT). Every
licence travels with what was taken: Basecoat's beside the vendored CSS in the
skeleton, Phosphor's in [icons/LICENSE.md](icons/LICENSE.md), and this project's
own in [LICENSE.md](LICENSE.md).

Nothing arrives through npm and nothing is fetched from a CDN. There is no
`package.json` in an Arandu project and the Content-Security-Policy is
`script-src 'self'`, so a script from another host would not run even if one were
referenced.

## Learning Arandu

The API reference is generated from the doc comments and lives on
[pkg.go.dev](https://pkg.go.dev/github.com/arandu-io/kyse). Every exported
symbol carries one, and that is deliberate: it is the documentation that cannot
drift from the code, because it sits in the same file.

The CLI documents itself. `aru help` lists every command, and each one explains
what it writes and what to do with it. `aru doctor` explains what it found and
what breaks, not which rule was violated.

A guide and a website do not exist yet, and that is a decision rather than a
gap: a guide written against an API that still moves is work done twice, and the
second time is worse — there is wrong documentation published. The site is the
next phase, and it will be an Arandu application.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security Vulnerabilities

Please review [our security policy](SECURITY.md) on how to report a
vulnerability. Never open a public issue for one.

## License

Open-sourced software licensed under the [MIT license](LICENSE.md).
