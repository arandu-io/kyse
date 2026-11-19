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

## v0.54.0 — `FileUpload` draws the list of what was chosen

Requires hesape v0.110.0, which carries the behaviour that fills the list. No
field was removed or renamed and nothing needs to compile differently. What
changes is what a field that leaves its files to the form draws, and what a
screen that worked around the empty list can now delete.

| it was | it is |
| --- | --- |
| an empty `<ul data-part="list" aria-live="polite">` | the same list without `aria-live`, a `<p data-part="summary" role="status">` before it that is the polite announcement, and a `<template data-file-template>` after it with the row |
| the files chosen were not shown | one row per file, drawn by the `file-upload` behaviour: a thumbnail for an image, the name, the size, and a count in the status line |
| no limits | `MaxFiles` and `MaxBytes`, warned before the submit |

A field with an `UploadURL` is not changed.

- **Remove your own script and rules.** An application that listed the files
  with a script of its own (a `change` listener that fills `[data-part="list"]`)
  deletes it: left in place it redraws the list after the behaviour does, with
  its own sentences. The classes the old rules named (`file-upload-summary`,
  `file-upload-item`, `file-upload-name`, `file-upload-size`) are not drawn any
  more, so those rules do nothing.
- **Pass the limits where you validated them.** `MaxFiles: 6, MaxBytes: 5 << 20`
  is what the page warns about; keep the server's check, which stays the rule.
- **A project with its own `Behavior` on the props** is given no props by this
  component: pass `FileUploadProps{...}.BehaviorProps()` as `Behavior.Props` to
  keep the limits and the sentences.
- **Translate the sentences** in `lang/<locale>/kyse.json` under
  `file_upload`, or write them on the props. Rebuild your views with
  `aru view:build`.

## v0.53.1 — swapped content says where the focus goes

Requires hesape v0.109.0, which serves YPY v0.1.1. Since that version a swap
focuses the `[autofocus]` element of the content it brings in only when that
element carries `ypy-focus` (or the application turns `config.autofocus` on).
Nothing here turns it on for you, and nothing needs to compile differently.

The components that draw an `[autofocus]` element now write `ypy-focus="auto"`
on it, so the cursor lands where it did before, also when the markup arrives by
a swap or a boosted link: `EditInPlace` (the field and the textarea) and the
field of every component with an `Autofocus` prop (`Field`, `Input`,
`InputGroup`, `Password`, `Masked`, `Combobox`, `OneTimeCode`). The cancel button
of an alert `Dialog` or `FormDialog` is left alone: the browser focuses it when
the dialog opens as a modal.

If you draw an `[autofocus]` element of your own into content a swap brings in,
add `ypy-focus="auto"` to it. Rebuild your views with `aru view:build`.

## v0.53.0 — the request props and attributes are YPY's

This release stops compiling in every application that sets a request prop, and
changes what is drawn for every component that asks the server for something.
The view layer's request library is YPY now, and the components speak only its
native language: `ypy-*` attributes, `ypy:*` events, `YPY-*` headers, morph as the
default swap, and the native rule for a response of 400 or more.

### The fields

Every field of a props struct that carried the old library's prefix is now
spelled `Ypy...`. The sets are the same and the values mean what they meant: an
address, a selector, a swap style, a question, a trigger.

| it was | it is | on |
| --- | --- | --- |
| the old `...Post` | `YpyPost` | `ButtonProps`, `MenuItem`, `SegmentedControlProps`, `SelectProps`, `SplitButtonProps`, `ToolbarItem` |
| the old `...Get` | `YpyGet` | `ButtonProps`, `MenuItem`, `SegmentedControlProps`, `SelectProps`, `SplitButtonProps`, `ToolbarItem` |
| the old `...Target` | `YpyTarget` | `ButtonProps`, `MenuItem`, `PaginationProps`, `SegmentedControlProps`, `SelectProps`, `SplitButtonProps`, `TableProps`, `ToolbarItem` |
| the old `...Swap` | `YpySwap` | `ButtonProps`, `PaginationProps`, `SegmentedControlProps`, `SelectProps`, `SplitButtonProps`, `TableProps`, `ToolbarItem` |
| the old `...Confirm` | `YpyConfirm` | `ButtonProps`, `MenuItem` |
| the old `...Trigger` | `YpyTrigger` | `SegmentedControlProps`, `SelectProps` |
| the old `...Include` | `YpyInclude` | `SelectProps` |

The row actions of a `DataTable` are `MenuItem`s, so `row.Actions[n].YpyPost`,
`YpyGet`, `YpyTarget` and `YpyConfirm` follow the table above.

What an application does:

- **Rename the fields with `ypy migrate --write .`**, run over the views and the
  Go that builds the props. It renames the Go identifiers and the attributes in
  your own markup in the same pass, and `ypy check .` lists what is left for a
  person.
- **Load the YPY script in the layout** (the hesape and framework releases that
  embed `ypy.min.js`) and delete the old one. A page drawn by this release asks
  for `ypy-*` and an older script reads none of them. The test server of a
  project that renders components needs the same swap: `view.Asset("ypy.min.js")`.
- **Run `aru view:build` with the `aru` of the same move** before regenerating
  the compiled files. The older compiler does not know the new names, and loses
  the check on the address of a request.
- **A value that opens a swap, a header or a confirmation with a script prefix
  stops the build**: `YpyConfirm` that starts with `js:` or `javascript:` is
  refused at render with an error naming the line.

### What the components do now

- **A response of 400 or more swaps nothing.** YPY puts such a response into the
  target by default, and that would replace a row, a list or a button with an
  error page. Every element that asks for a fragment now writes
  `ypy-error-swap="none"`: the page stays as it was and `ypy:response:error`
  fires on the element, for a listener to show. It covers `ActiveSearch`,
  `Autocomplete`, `Button`, `Combobox`, `DataTable`, `DeleteRow`,
  `DropdownMenu`, `EditInPlace`, `LazyLoad`, `LoadMore`, `OptimisticToggle`,
  `Pagination`, `SegmentedControl`, `Select`, `SplitButton`, `Table` and
  `Toolbar`. The save of an `EditInPlace` keeps the native 422: an invalid save
  draws the form again with its message, and a failure of the server (5xx)
  leaves the form as it was.
- **The default swap is morph.** A `Ypy...Swap` left empty used to mean the old
  library's `innerHTML`. It now means YPY's default, which keeps what the reader
  is typing, selecting or scrolling. Write `YpySwap: "innerHTML"` where the raw
  replacement was the point. The components themselves always write their swap.
- **A 204 and a 304 swap nothing.** A `DeleteRow` removes its row when the
  endpoint answers 200 with an empty body; a 204 leaves the row where it is. The
  comment of the component said the opposite.
- **Nothing is inherited without being said.** The components write every
  attribute on the element that uses it, so none relies on a container's
  attribute reaching its descendants. A page of yours that did, writes
  `:inherited` on the attribute (`ypy-target:inherited="#out"`).
- **The layout writes `ypy-boost="true"` once, and a component that must not be
  boosted writes `ypy-boost="false"`.** `DialogTrigger`, `FormDialog` (`NoBoost`)
  and the links of a `DataTable` already do.
- **`Masked` knows `U` and `X`**, the two tokens of hesape's mask that write
  capitals: `U` a letter, `X` a letter or a digit, so a plate typed `abc1d23`
  reads `ABC-1D23` in the first frame as it does under the script. Needs hesape
  that carries them.
- **The script of the view layer follows YPY's order of events**: for the main
  swap `ypy:after:settle` fires before `ypy:after:swap`. A listener of your own
  that carried state from the first to the second needs to read it in this order.

## v0.50.0

Nothing stops compiling. One thing changes in what is drawn, and it undoes a
break of v0.48.0.

**A `DataTable` is cards by the width of the screen, not by the width of the block it
sits in.** v0.48.0 drew cards where the table's own region was narrower than 600px
(a container query), so a short list inside a card of a dashboard, which is under
600px on any desktop, became a list of cards with the sort bar for a head. From
v0.50.0 the rule is a media query on the screen: under 600px wide the rows are
cards, from 600px up the table is the table it was in v0.47.0, in any block. Ayra's
`DataTable` asks the width of the window in the same way. Rebuild the stylesheet
with `aru view:build`: the cards are `max-[600px]:` utilities now, and the region is
no longer a container (`@container/dt` is gone from the markup).

The rule lives in the component and not in each application: a `DataTable` is
described once, with typed columns (`Primary`, `CardHidden`, `Actions`), and each
layer decides how to draw it. Kyse draws a table on a computer and cards on a
phone's browser, from the width of the screen; Ayra draws cards on a phone and a
table on a wide window, from the width of the window. Nobody marks a table to get
this. The two props below are the explicit exception.

Two props give the choice back to the list:

| you want | write |
| --- | --- |
| the table on a phone too (a short list beside a video, in a card) | `Mode: components.DataTableModeTable` |
| the cards on a desktop too | `Mode: components.DataTableModeCards` |
| the default, cards under 600px | nothing, or `DataTableModeAuto` |

A script that adds a row as it arrives should not write the classes itself. Set
`RowTemplate: true`, clone `template[data-part="row-template"]` in the list and
write the cells: the row has the classes of the others in every mode.

---

## v0.48.0

Nothing stops compiling. Three things change in what is drawn, and the last one
needs your attention if a test of yours reads the footer.

**The English footer of a `DataTable` is no longer "{total} row(s)."** It says
"1 row." for one row and "{total} rows." for every other count. A project that
passes `Lines.Rows`, or whose catalogue carries `kyse.data_table.rows`, draws
exactly the sentence it drew, for every count. To say a singular, add the two
forms:

| it was | it can be |
| --- | --- |
| `Lines{Rows: "{total} linha(s)."}` | `Lines{RowsOne: "{total} linha.", RowsMany: "{total} linhas."}` |
| `kyse.data_table.rows` | `kyse.data_table.rows_one` and `kyse.data_table.rows_many` |

A form written on the props wins over the catalogue, and a form written beside
`Rows` wins for its own count. The sentence of a selected row ("{n} of {total}
row(s) selected.") is not split.

**A `DataTable` is drawn as cards where its region is narrower than 600px.** The
region is now a container (`@container/dt`) and every row, cell and header
carries classes under `@max-[600px]/dt:`. Tailwind compiles them from the
component sources, so rebuild the stylesheet (`aru view:build`) when you
upgrade; a stylesheet compiled before it has no card rules and draws the table
at every width, which is the old behaviour and not a break. Mark the card's
title with `TableColumn.Primary`, and leave a column out of it with
`CardHidden`; with none marked the first column the card shows is the title. A
`Parts{"row": {Replace: true, ...}}` or the same on `cell` drops the card
classes of that part, which is the way to keep a table a table.

A cell is the same element it was and its text is only its value: the name of the
column is the cell's `data-label`, which the card draws in front of it.

**The footer of a list the browser narrows** (`Complete: true`) says one row in
the singular only with `hesape` v0.93.0 or later, which carries `rowsOne` in the
table behaviour. With an older one the footer is right on the first paint and
says the plural ("1 rows.") after a search leaves exactly one row.

---

## v0.46.0

Nothing stops compiling. What changes is the name of three parts of `Calendar`,
and so of the `DateTimePicker` that draws one: they were called `previous`,
`next` and `footer`, which is what a `DataTable` calls the parts of its pager, so
a calendar inside a table's toolbar was found first by anything that looks the
pager up by part. A style, a test or a script that reached them by the old name
stops finding them.

| it was | it is |
| --- | --- |
| `Parts{"previous": {...}}`, `data-part="previous"` on a calendar | `"month-previous"` |
| `Parts{"next": {...}}`, `data-part="next"` on a calendar | `"month-next"` |
| `Parts{"footer": {...}}`, `data-part="footer"` on a calendar | `"calendar-footer"` |

`DateTimePicker` with `Calendar` also writes `popover` on its group by default.
The shipped popover behaviour binds to an element of that name and a group
without it drew a calendar nothing opened, so an application that added the class
to the group by hand (`Parts{"group": {Class: "popover"}}`) can stop.

`DateRangePicker`, a period chosen in one field, is new and breaks nobody. It
needs `hesape` v0.89.2 or later, which carries its behaviour, and the stylesheet
of the skeleton (`date-range-picker.css`).

---

## v0.45.0

A build stops here, and only a list drawn with `DataTable` is touched. The data
table is drawn as shadcn/ui's data table and works the way a server-rendered
list has to: one value says what the list is showing, every control is an
address that value builds, and one form makes one request. `Table`, `Button` and
every other component keep their props, and `Button` gains `Name` and `Value`.

The state moved out of five fields into one:

| it was | it is | write instead |
| --- | --- | --- |
| `DataTableProps{SearchName: "q", Query: "ada"}` | `Searchable` draws the box; the text travels under the state's search name (`q`) | `DataTableProps{Searchable: true, State: components.DataTableValues{Q: "ada"}}` |
| `SortKey: "total", SortDir: "desc"` | `DataTableProps.SortKey` and `SortDir` are gone; the order is the state's | `State: components.DataTableValues{Sort: "total", Dir: "desc"}` |
| `DataTableFilter{Selected: []string{"paid"}}` | `DataTableFilter.Selected` is gone; `DataTableState.Chosen(key)` says what is chosen | `State: components.DataTableValues{Facets: map[string][]string{"status": {"paid"}}}` |
| `Page: 3, Pages: 9, Total: 137` | `DataTableProps.Page`, `Pages` and `Total` are gone, as their deprecation said | `Paginator: page`, the `*pagination.LengthAwarePage` the query returned |
| `PageSize: 20` on a server list | read from the state first, then `PageSize`, then the page | `State: components.DataTableValues{Size: 20, DefaultSize: 10}` |

`State` is any `components.DataTableState`: `Params`, the values, and
`Values(change)`, the canonical query after one step. `DataTableValues` is the
plain implementation -- values written once, in a stable order, with the default
page and size left out -- and it is what a list that holds all its rows passes. A
nil `State` is a list in its first state.

A list read from a database is not written by hand at all. Declare it once with
`pagination.NewListing` (hesape v0.86.0), let the query builder apply it, and
hand the result to the table: `DataTableProps.Result` takes the
`*pagination.ListingResult` and fills the state, the page with its total, the
count beside every facet value, the facet's own query parameter and the page
sizes the listing allows, each of which a field written beside it still
overrides. `components.DataTableStateFrom(state)` is the same adapter for a state
you hold, and it answers every step with the listing state's own method, so
the rules are the listing's: a search, a facet or an order goes back to the
first page, a new size keeps the first row on screen, a column leaves the page
alone, and a step the listing did not declare changes nothing. The cursor of a
keyset page is not a part of a listing and is ignored there; a list paged by
cursor keeps `DataTableValues`.

```go
// The service: the listing is declared once, beside the model.
var invoiceListing = must(pagination.NewListing(pagination.ListingSpec{
	Search:      []string{"number", "customer"},
	Sorts:       map[string]string{"number": "number", "total": "total"},
	DefaultSort: "number",
	Facets:      []pagination.FacetSpec{{Key: "status", Column: "status", Values: []string{"paid", "open"}}},
	PageSizes:   []int{10, 20, 50},
	Hideable:    []string{"total"},
}))

func (s *Invoices) List(ctx context.Context, g auth.Grant, values url.Values) ([]query.Record, *pagination.ListingResult, error) {
	return s.db.Table("invoices").Listing(ctx, g, invoiceListing.Parse(values), "/invoices")
}

// The controller: authorize, then ask. The tenant is the Grant's.
rows, result, err := invoices.List(ctx, g, r.URL.Query())

// The view: the result goes to the table as it came back.
{!! components.DataTable(components.DataTableProps{
	ID: "invoices", Label: "Invoices", URL: "/invoices", Searchable: true, FullPager: true,
	Result:  .Result,
	Filters: []components.DataTableFilter{{Key: "status", Label: "Status"}},
	Columns: .Columns,
	Rows:    .Rows,
}) !!}
```

A `DataTableFilter` that names no `Options` lists the facet's own values, and one
whose facet travels under another `Param` writes it on its boxes.

Every address is the state's answer, so none is built from loose parameters any
more, and these exported methods of `DataTableProps` are gone with the parameters:
`Address`, `SortURL`, `PageURL` and `SearchAddress`, which built an address from
the fields above and accumulated a facet value on every control; `SearchTrigger`,
`ApplyText`, `Filtered`, `FacetChosen`, `FacetSummary`, `FacetName`, `Grid`,
`Pager`, `NextPage`, `PreviousPage`, `NextText`, `PreviousText`, `PageOfTemplate`,
`PerPageText`, `Ranged`, `Showing` and `SelectedText`. `Chosen` now answers the
values chosen for a facet and not whether one option is in force; `OptionChosen`
answers the second. The replacements are `Href(DataTableChange{...})`, `Active`,
`Chosen` and the other methods of the type.

| it was | it is | write instead |
| --- | --- | --- |
| `BulkAction` empty posted the selection to `URL`, which answered 204 and did nothing | the bar is drawn only when `BulkAction` is set | `BulkAction: "/invoices/bulk"` |
| every bulk button posted the same fields | `ButtonProps.Name` and `Value` are written on the button, so the server can tell which was pressed | `ButtonProps{Label: "Archive", Type: "submit", Name: "do", Value: "archive"}` |
| the rows ticked belonged to the page shown | the browser keeps the keys ticked per table across a search, an order, a facet and a page, and the bar posts all of them | nothing |
| a complete list's links could be boosted into a request | a complete list's links and form turn the old library's boost off and the behaviour takes the click | nothing |
| a sort, a page step or a search lost the hidden columns | the hidden columns are `State.Hidden()`, travel as `hide`, and survive a swap and a reload | `TableColumn{Hideable: true}`; `TableColumn.Hidden` is where a list with no `State` starts |
| `PageSizes` drew a chooser on a complete list only | the chooser, a menu of addresses, is drawn on both, and a new size keeps the first row on screen | `PageSizes: []int{10, 20, 50}` |

`DataTableLines` loses `Showing`, `PerPage` and `Page`, and gains `Filter`,
`Rows`, `SelectedCount`, `RowsPerPage`, `First`, `Last`, `SelectRow`, `Reset`,
`ClearFilters`, `NoFacetResults`, `FacetSelected`, `OpenMenu`, `Actions`,
`ToggleColumns`, `Asc`, `Desc` and `Hide`. The catalogue keys
`kyse.data_table.showing`, `kyse.data_table.per_page` and `kyse.data_table.page`
are no longer asked for, and the new ones are listed in the package doc. Three
English floors changed with the reference: `kyse.data_table.no_result` is "No
results.", `kyse.data_table.previous` and `next` are "Go to previous page" and
"Go to next page", and a data table without a catalogue names its sortable
headers "{column}, sorted ascending", "descending" or "not sorted" -- `Table`
does too, and draws the arrow for the direction.

The parts `indicator` and `showing` are no longer drawn. `wrapper`, `reset`,
`facet-search`, `header-menu`, `row-menu`, `first` and `last` are new, and the
rest keep their names and their elements.

The look is the shadcn/ui data table's, written as Tailwind utilities in the
markup: a project's stylesheet needs nothing of the data table beyond what
`aru view:build` already compiles from this module. The rules of
`data-table.css` and the `.table-columns`, `.table-count` and `.table-sort`
rules of `table.css` that Basecoat shipped for the old markup are not read by it
and can go. The behaviour that makes the menus, the selection and the swap work
is in the view layer's script from hesape v0.87.0; this release requires it.

## v0.42.0

Nothing stops a build, and a `Sidebar` with no `Mode` sends the markup it sent
before, byte for byte. The column can fold to a strip of icons on a wide screen
instead of disappearing, so a view can stop drawing its own.

| it was | it is | write instead |
| --- | --- | --- |
| a collapse button whose action hid the whole column on a wide screen | `SidebarProps.Mode`: `"full"` or `"slim"`, drawn by the server and written on the root as `data-sidebar-mode` | `SidebarProps{Mode: "full"}`, or `"slim"` |
| an icon pasted into `Label`, or a hand-built list to have one | `SidebarItem.Icon`, drawn before the label and all the strip shows | `SidebarItem{Label: "Plates", URL: "/plates", Icon: icons.Car(icons.Props{})}` |
| a tip, a separator and a narrower margin written in the project's stylesheet | the strip's look travels with the component as utilities; `--sidebar-slim-width` (3rem) sets its width | nothing |

In the strip each link is named by its label (`aria-label` and the hidden
text), the tip beside it appears on hover and on keyboard focus, a group's
heading is drawn as a line, and the current entry keeps `aria-current`. The list
does not scroll in the strip, so the tips are not clipped at its edge. On a
narrow screen the column is the overlay with its labels in either mode, and
every transition is shortened and taken away under reduced motion.

Switching after the first paint is setting `data-sidebar-mode` on the root to
the other value. From hesape v0.84.0 the view layer's script does it: give the
header button `data-sidebar-toggle="<sidebar id>"`. On a wide screen it folds the
sidebar between full and slim and remembers the choice per visitor
(`localStorage` key `arandu.sidebar.<id>`, read and written inside try/catch, so
a page without storage still renders the mode the server drew); on a narrow
screen, or on a sidebar with no mode, it opens and closes the overlay as before,
and it keeps `aria-expanded` on every button that points at the sidebar.

```html
<button type="button" data-sidebar-toggle="panel-sidebar" aria-controls="panel-sidebar">Menu</button>
```

A page that registered its own `toggle-sidebar` action for v0.42.0 deletes it
and the restore loop, and changes the button to `data-sidebar-toggle`.

## v0.41.0

Nothing stops a build, and every call that compiled draws what it drew. The
components take what a screen already holds instead of asking a view to copy
it out, and the ways that copied it are deprecated.

| it was | it is | write instead |
| --- | --- | --- |
| `PaginationProps{Page: p.CurrentPage(), Pages: p.LastPage(), URL: "/x?page={page}"}` | `Paginator` takes the `*pagination.LengthAwarePage`, `*pagination.Page` or `*pagination.CursorPage` the query returned; `Page`, `Pages`, `URL` and `Around` are deprecated | `PaginationProps{Paginator: page}` |
| `DataTableProps{Page: …, Pages: …, Total: …}` | `Paginator` as above; `Page`, `Pages` and `Total` are deprecated | `DataTableProps{Paginator: page}` |
| `Token: .Page.CSRFToken()` on `Dialog`, `FormDialog`, `Table`, `DataTable` | the form posts the page's token when `Token` is empty | `Page: .Page` (`PageData: .Page` on `DataTable`) |
| `ToastProps{Title: .Page.Flashed("status")}` behind an `@if` | the toast draws what the page flashed under `FlashKey` (`"status"`) when `Title` is empty, and nothing when nothing was | `ToastProps{Page: .Page}` |
| `CancelLabel: "Cancelar"` and the other default words, on every call | read from the catalogue group `kyse` when the page is a `TranslatorPage` | one `lang/<locale>/kyse.json`; the keys are in the package doc |

A pager drawn from a `Paginator` uses the page's own window (`OnEachSide`, three
either side by default) rather than one either side, and the page's addresses,
which keep the query the request arrived with. A data table keeps building its
addresses from its own fields, with the page's parameter name. A page that
knows no total -- the simple page and the cursor page -- draws the two steps and
neither the range nor "Page n of total".

`hesape/view.Page` answers `CSRFToken` and `Flashed`, so `Page: .Page` is all
the token and the toast need. It does not carry a translator: a screen's own
page data implements `TranslatorPage` by returning the translator it was wired
with and the request's locale. A page that does not draws the English, as
before.

## v0.39.0

Nothing stops a build, and one thing stops a render that used to come out
empty. An attribute the view layer refuses -- the request library's own, `on*`, `style`, an
address, `class`, `role`, `aria-*` -- written in `Attrs` or in
`Parts[...].Attrs` used to make the whole component render as nothing. It now
panics with a `*components.RefusedAttribute` naming the part, the attribute and
the reason, which a test that renders the page reports and the exception
handler turns into an error page.

| it was | it is | write instead |
| --- | --- | --- |
| an `Attrs` entry that turned the boost off on a `FormDialog` drew no dialog | a panic naming that attribute on `content` | `NoBoost: true` |
| any refused name in `Attrs` drew an empty component | a panic naming it | the component's own field for it, or nothing |

A value the escapers refuse because it arrived as data -- an address whose
scheme a browser acts on -- still draws nothing, as before.

## v0.36.0

Nothing stops a build. `mailui` can draw a message in a language other than
English: `LayoutProps.Lang` sets the `lang` attribute (empty, or anything that
is not a BCP 47 tag, stays `"en"`), and `FallbackWith(FallbackProps{Href,
Text})` writes the address under a button after the sentence it is given.
`Fallback(href)` still writes the English sentence, byte for byte.

## v0.35.0

Nothing stops a build. The data table is drawn as shadcn/ui's:

- each facet is a dashed button with what is chosen beside it and a menu of
  checkboxes, and "Apply filters" is drawn only for a browser with no script;
- the footer carries the selection count, the range, the page size, "Page n of
  total" and the two steps, and the table no longer draws its own count when a
  data table composes it;
- a sortable header carries the arrows icon, and every checkbox the class
  `input`.

A complete list filters by its facets in the browser when each row says its
value: `TableRow.Facets` (`{"status": "open"}`). `DataTableLines` gains
`PageOf`, `Back` and `Forward`. The stylesheet that draws all of it is the
starter's `resources/css/basecoat/components/data-table.css` and `table.css`;
copy both into an application created before this release. kyse now requires
Hesape v0.54.1, whose script works the facets and the pager.

## v0.34.2

Nothing to change. A data table's search, facets, sortable headers and pages
now take the region out of the answer (the select equal to the target), so
an endpoint that answers the whole page swaps one list for one list instead
of drawing the page inside the table.

## v0.34.1

Nothing below stops a build of this module's own symbols.

### kyse requires Framework v0.56.0 and Hesape v0.54.0

Requiring this version moves an application to both. Follow their upgrade
notes between the versions you had and these: Framework v0.56.0 removes
`config.Config.SessionTTL`, so `config.Load` no longer reads `SESSION_TTL`
(the session lasts `SESSION_LIFETIME` minutes), and puts `APP_NAME` on every
request; Hesape v0.54.0's `view.New` reads it into `Page.AppName`. A
controller that assigned the configured name by hand can delete the
assignment, and an application that never set `APP_NAME` now draws
`arandu-app` where its brand was blank.

## v0.34.0

Nothing below stops a build of this module's own symbols. Two components draw
different markup, and kyse requires Hesape v0.53.0, whose script mounts it.

### `OneTimeCode` posts the code without script

The `raw` part is a visible `type="text"` input carrying `id`, `name`,
`inputmode`, `autocomplete="one-time-code"`, `maxlength`, `pattern` and
`autofocus`. The squares are drawn `hidden`, without `id` or `autofocus`, and
the `one-time-code` behaviour shows them and turns the input into the hidden
field. A `Label` whose `For` names the field keeps working.

### `Masked` posts what was typed without script

The visible `input` part carries `Name`, and the `raw` part is drawn
`disabled`. The `mask` behaviour enables it and renames the box to
`DisplayName`, which the component hands over as `displayName`. Without script
the value arrives as typed: validate it with `mask.Pattern.Accepts` and store
`Unmask(value)` from the same pattern.

Styles or tests that targeted `[data-part=square]#<id>`,
`input[type=hidden][name=…]` or `name="<Name>_display"` in the drawn markup
follow it.

### kyse requires Hesape v0.53.0

Requiring this version moves an application to it. Hesape v0.53.0 makes
`mask.Pattern.Accepts` take the raw value as well as the formatted one.

## v0.33.0

Nothing below stops a build of this module's own symbols.

### kyse requires Framework v0.55.1 and Hesape v0.52.0

Requiring this version moves an application to both. Follow their upgrade
notes between the versions you had and these: Hesape v0.52.0 removed the
Php/Artisan/Symfony-named symbols and deprecated the SessionManager path;
Framework v0.54.0 refuses a malformed boolean in the environment, and v0.55
replaces `Configuration.Session` with `bootstrap.Session{Secure, Lifetime}`
and refuses at boot a `SESSION_*` variable it does not read.

## v0.32.0

Nothing below stops a build of this module's own symbols.

### kyse requires Framework v0.50.0 and Hesape v0.50.3

Hesape v0.50.3 does not build with Framework v0.49.0 or older, so the framework
floor moved with it. An application on an older framework follows that
framework's v0.50.0 upgrade notes: `KeyBySession` takes the session store, and
`CSRFProtect` binds each token to the session or to a signed guest cookie. A
test that issued a CSRF token with `Issue("")` now gets `ErrUnboundToken`: pass
what `CSRF.Binding` returns.

### `Masked` and `OneTimeCode` write the attributes of every part

The attributes written for the `input` and `raw` parts of `Masked`, and for the
`square` and `raw` parts of `OneTimeCode`, were dropped. They are now written,
last in the tag, like every other part's.

## v0.31.0

### `OneTimeCode` draws its rejection

The message `Page` holds for `Name` is now drawn inside the group as the
`message` part, under `ElementID()-error`, and named first in the group's
`aria-describedby`. A view that drew its own sentence beside the field and
pointed `DescribedBy` at it now shows it twice: remove the sentence and the
`DescribedBy`.

## v0.30.0

Nothing below stops a build. Each is a change in what a component draws, and
each closes a way the markup could be made to carry something its caller did
not write.

### `DataTable` no longer carries the query string of `URL`

`DataTableProps.URL` was documented as the endpoint without a query string,
but its query was copied onto every sort header, every page link and, as
hidden fields, into the search form. A table drawn with the request's own
address therefore resubmitted any parameter written into a link
(`/users?_method=DELETE&role=admin`) from every control on the page.

The query and fragment of `URL` are now dropped. Every address and hidden field
is built from the table's own state: the search, the order, the page and the
chosen `Filters`. A parameter that was riding on `URL` has to become one of
those — a facet is a `DataTableFilter` — or part of the path.

### `Dialog` and `Table` write the CSRF token only into a form that posts

A `Dialog` with `Method: "get"`, or a `Table` whose `BulkMethod` was `get` or a
method a browser cannot submit (which it then sends as GET), wrote `_token`
into that form, so the token ended up in the address, the history and the
access log. The field is now written only when the form's method is POST —
including PUT, PATCH and DELETE on `Dialog`, which travel as POST. Both props
gained `SendsToken()`, which says whether it will be written.

### `Attrs` can no longer replace an attribute the component writes

The caller's `Attrs` — on the root or on any part — are now written after
every attribute the component writes itself, on every element. A browser keeps
the first of two attributes with the same name, so where `Attrs` used to win a
collision on some elements (the `value` of `Checkbox` and `Switch`, the `name`
and `value` of `Combobox`'s input, the `id` of a tab, among others) the
component's own value now wins everywhere. An `Attrs` entry still adds any
attribute the component did not write. A call that relied on `Attrs` to change
a `name`, `value`, `type`, `form` or `id` sets the component's own field
instead.

## v0.15.2 — Dialog submissions use the browser's native transport

`DialogProps.Method` may still be PUT, PATCH or DELETE, but an HTML form cannot
send any of those methods. The rendered form now sends POST and carries the
intended method in a hidden `_method` field. Applications must register
`hesape/http/middleware.OverrideMethod` **after** Framework's `CSRFProtect` and
before routing; the Arandu skeleton does this in its default pipeline.

The dialog's CSRF field is `_token`, and this module now requires Framework
v0.42.1, whose `CSRFProtect` reads that spelling. A public test
renders the component, simulates the browser submission, and crosses
`CSRFProtect → OverrideMethod → Router` for all three spoofed methods.

### v0.15.0 and v0.15.1 are retracted

Both releases selected Framework v0.41.0 together with Hesape v0.21.0. That
pair disagrees about the form field: Framework reads `_csrf`, while Hesape's
`@csrf` directive writes `_token`. v0.15.1 also changed Dialog itself to
`_token` without raising the Framework floor, so its own confirmation form was
refused with 419. In both releases, a Dialog configured for PUT, PATCH or DELETE
also wrote that method directly on the form, which browsers do not submit.

Use v0.15.2 or newer rather than either retracted version. Retraction is
metadata published by v0.15.2; it cannot remove an immutable tag, but
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
library would act on — `on*`, the request library's `on` attributes, `x-*`, `style`, `href` and the other
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
