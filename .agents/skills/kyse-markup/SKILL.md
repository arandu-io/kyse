---
name: kyse-markup
description: Write the markup inside a kyse component — the part of a .kyse.go file after @endgo. Use when choosing a CSS class for a component, when tempted to add a stylesheet, a <style> block, an inline style, a <script>, an onclick, an x-data, an x-on:click or an @click, when a component needs to be interactive, when writing a loop or an empty state, or when deciding between {{ }} and {!! !!}. Also use when the request is to "style this component", "make it open a menu", "add the hover state", "port the CSS from shadcn", or when a component renders but looks unstyled. Covers the class names the shipped stylesheet knows, why behaviour is data-* delegation, the loop that wraps the wrong element, and the escaping rules.
license: MIT
---

# Writing the markup of a component

The part after `@endgo` is HTML with directives in it. Most of what a component
library usually reaches for to make a component look and behave right is absent
here, and each absence has its own reason.

## No CSS in this repository

There is no `.css` file in this tree and there is no place to put one. The
markup carries semantic class names — `btn`, `field`, `card`, `table` — and the
rules behind them are Basecoat, vendored into the project skeleton under
`resources/css/basecoat/`, one stylesheet per component.

So **a class you invent is a component with no appearance**. It renders, it
passes the smoke test, it looks correct in review, and it is unstyled in
somebody's project. Before writing a class name, check that the stylesheet knows
it — with the skeleton checked out beside this repository:

```sh
grep -rn '\.your-class' ../arandu/resources/css/basecoat/components/
```

The vocabulary is the shadcn one, adapted: `accordion`, `alert`, `alert-dialog`,
`avatar`, `avatar-group`, `badge`, `breadcrumb`, `btn`, `button-group`, `card`,
`combobox`, `command`, `command-dialog`, `dialog`, `drawer`, `dropdown-menu`,
`empty`, `field`, `field-separator`, `fieldset`, `input`, `input-group`, `item`,
`item-group`, `kbd`, `label`, `popover`, `progress`, `scrollbar`, `select`,
`sidebar`, `skeleton`, `table`, `table-container`, `tabs`, `textarea`, `toast`,
`toast-content`, `toaster`, `tooltip`.

Not everything with a stylesheet has a class of its own. A checkbox is
`.field > input[type='checkbox']` and a switch is the same input with
`role="switch"`, which is why `Checkbox` and `Switch` write `class="field"` and
`class="input"` and nothing else. Read the stylesheet before inventing a name
for a thing it already selects.

Variants and sizes are attributes rather than classes —
`data-variant="destructive"`, `data-size="sm"` — and the stylesheet reads them
off the element.

Tailwind utilities are available for one-off spacing and typography inside a
component (`mt-4`, `text-muted-foreground text-sm`, `tabular-nums`). They are
not how a component gets its identity: the same button written with utilities is
443 characters of markup and thirty as `class="btn"`.

If the appearance genuinely needs a new rule, the rule goes in the skeleton's
stylesheet, in that repository, and the component here uses the name it defines.

## No JavaScript in this repository

There is no `.js` file here either, and no component writes a `<script>`.
Behaviour comes from the view layer's own script, which the framework embeds and
serves: it binds once on `document` and dispatches by looking at `data-*`
attributes on the element an event came from. Markup HTMX swaps in is live the
moment it lands, because there is nothing to initialise.

**So a component makes itself interactive by carrying the attributes that script
already dispatches on**, and by carrying its state in the ARIA it has to write
anyway. `data-combobox` is the whole of what a combobox tells the client: open
is `aria-expanded` on the box and `aria-hidden` on the popover, the active line
is `aria-activedescendant`, the chosen one is `aria-selected`. The DOM is the
state and there is one copy of it.

A component that would need a script of its own is a component that does not
enter this repository. The behaviour belongs in that script, in the framework,
where it is one implementation for everybody.

## No directives that hold an expression

`x-data`, `x-on:click`, `x-bind:`, `:class`, `@click` — none of them work, and
they do not fail loudly. A directive's value is a string the client library
compiles with the `Function` constructor, and pages are served
`script-src 'self'` with no `unsafe-eval`, so the compilation throws at the
point of evaluation. The component is not noisy; it is dead, and it fails in a
browser console rather than in any build.

Nothing else in the toolchain would catch it: the markup is valid, the view
compiles, and the smoke test still finds a class in the output. So the sources
are read.

`TestNoComponentEvaluatesAnExpression` at `tests/Unit/components_test.go:354`
globs `components/*.kyse.go` and `components/*.go`, and fails on any `x-name` or
`@name=` it finds. Both halves are read: the source is where the directive would
be written, the compiled file is what a browser is actually sent. The failure
names the file, the line, the attribute, and what to do instead.

`hx-get`, `hx-post`, `hx-target`, `hx-on:click` are not matched and are correct
here. HTMX parses attributes rather than evaluating them.

## The loop wraps the repeated element and nothing else

This is the defect that no assertion about content can see. A loop written
around a whole table renders one complete table per row, each holding the full
row set — every label present, every header present, nothing missing and nothing
extra named. A test looking for substrings passes on five copies of the card as
happily as on one.

`@forelse` is where it comes from, because it puts both halves inside the loop.
So the empty state is a branch beside the loop, not a clause in it:

```go
@if(len(.Rows) > 0)
	<table class="table">
		<tbody>
			@foreach(.Rows as row)
				<tr>...</tr>
			@endforeach
		</tbody>
	</table>
@else
	{!! Empty(.Empty) !!}
@endif
```

`Table` and `StatCard` are both written this way. `TestStatCardDrawsOneTable` at
`tests/Unit/stat-card_test.go:60` counts rather than searches: one `<table`, one
`</table>`, one `<thead`, and one `<tr` per row plus the header.

A headed table with no body is not the empty state either. It reads as a table
that failed to load; draw the `Empty` component instead.

## `{{ }}` escapes, `{!! !!}` does not

`{{ }}` on everything that came from a caller. `{!! !!}` only on a call to a
named function that returns `template.HTML` — another component, or an icon:

```go
{!! Empty(.Empty) !!}
{!! icons.Trash(icons.Props{Label: "Delete this post"}) !!}
```

A component is entitled to skip escaping because everything it interpolated was
escaped by the view compiler when it was generated. **A value has been through
nothing.** `{!! .Body !!}` is stored cross-site scripting the first time one of
them comes from a person, and it runs for every reader of the page. In an
application `aru doctor` reports the shape as `raw-output-is-not-a-component`;
that check reads applications, not this library, so here the rule is yours to
keep.

`TestAPropIsEscaped` at `tests/Unit/components_test.go:131` renders twelve
components with `<script>alert(1)</script>` in a prop and fails if it survives.
Add your component to that map when its props carry text.

The one deliberate exception is `TableCell.HTML`, which exists so a cell can
hold a badge or a button. Its doc comment says what makes it safe — pass what
another component returned, never a string assembled around a value somebody
typed — and a new prop of that shape needs the same sentence.

`{{-- --}}` is a comment that is stripped and never reaches the page. The
existing components use it to record why the markup is the shape it is, next to
the markup.

## Directives

The set is closed: `@extends` · `@section` / `@endsection` · `@yield` ·
`@include` · `@if` / `@elseif` / `@else` / `@endif` · `@foreach` /
`@endforeach` · `@forelse` / `@empty` / `@endforelse` · `@for` / `@endfor` ·
`@while` / `@endwhile` · `@continue` · `@break` · `@go` / `@endgo` · `@csrf`.

A component uses `@if`, `@foreach`, `@for` and `@go`, and nothing else — this
returns nothing across all 37 of them:

```sh
grep -ohE '@(extends|section|yield|include|csrf|while|forelse)' components/*.kyse.go
```

`@extends`, `@section`, `@yield` and `@include` belong to pages and layouts, and
a component is neither: nothing renders it by name, no controller hands it data,
and it has no layout. `@forelse` is in the language and unused here, for the
reason in the section above.

A loop binding is an ordinary Go name. `@foreach(.Rows as row)` is fine; what
the compiler refuses is a Go predeclared identifier — `nil`, `len`, `string` —
and it says so with the file and the line.

## Accessibility is part of the markup, not a later pass

The components that pay for themselves do it here: a `Field` is a label, an
input, a conditional message and the two ARIA attributes that tie them
together, and the two attributes are the ones people forget. A message being
present does three things at once — draws the sentence, sets `aria-invalid`, and
points `aria-describedby` at the sentence.

An icon-only control needs a name: pass `Label` to the icon, or write an
`aria-label` on the button. An icon beside a word that already says it stays
decorative.

## A note on inline styles

`style-src 'self'` is served with no `unsafe-inline`, so a `style` attribute is
dropped the same way an inline script is. Use a class. `range-slider.kyse.go`
currently writes one to seed a custom property, and nothing in the build catches
it — that is a known inconsistency, not a precedent to copy.
