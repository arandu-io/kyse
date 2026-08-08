<h1 align="center">arandu-io/kyse</h1>

<p align="center">The component library for Arandu. Import it and draw a screen.</p>

<p align="center">
<a href="https://github.com/arandu-io/kyse/actions/workflows/ci.yml"><img src="https://github.com/arandu-io/kyse/actions/workflows/ci.yml/badge.svg" alt="Build Status"></a>
<a href="https://pkg.go.dev/github.com/arandu-io/kyse/components"><img src="https://pkg.go.dev/badge/github.com/arandu-io/kyse.svg" alt="Go Reference"></a>
<a href="https://github.com/arandu-io/kyse/tags"><img src="https://img.shields.io/github/v/tag/arandu-io/kyse?label=version" alt="Latest Version"></a>
<a href="LICENSE.md"><img src="https://img.shields.io/github/license/arandu-io/kyse" alt="License"></a>
</p>

## What this is

The components a screen is made of, written in kyse and adapted from
shadcn/ui by way of shadcn-htmx.

The name is the view language they are written in: a `.kyse.go` is what an
Arandu view is, and this is the library of them. The compiler that turns one
into Go is part of `aru` and stays there — it is a build step, not a
dependency, and a project already has it.

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
			Value:    .Form.Title,
			Error:    .Errors.First("title"),
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
| `Textarea` | the same, for a textarea |
| `ThemeToggle` | light, dark and six accents, kept on the device |
| `Toast` | a flash message HTMX appends to the tray |

Eleven, and the catalogue upstream has eighty-two. What is here is what has been
adapted, not what is planned.

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
Both licences travel with what was taken: Basecoat's beside the vendored CSS in
the skeleton, and this project's own in [LICENSE.md](LICENSE.md).

Nothing arrives through npm and nothing is fetched from a CDN. There is no
`package.json` in an Arandu project and the Content-Security-Policy is
`script-src 'self'`, so a script from another host would not run even if one were
referenced.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## Security Vulnerabilities

Please review [our security policy](SECURITY.md) on how to report a
vulnerability. Never open a public issue for one.

## License

Open-sourced software licensed under the [MIT license](LICENSE.md).
