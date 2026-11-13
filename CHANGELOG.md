# Release Notes

## [v0.52.1](https://github.com/arandu-io/kyse/compare/v0.52.0...v0.52.1) - 2026-11-13

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.52.0...v0.52.1

## v0.52.1

- `Password`: the strength bar draws. Its fill is a span, a span is inline, and the width the behaviour writes on every keystroke does not apply to an inline box, so in an application whose stylesheet does not make the track a flex container (Basecoat's `progress.css` does; an application may leave it out) the bar stayed empty at every count while the lines under it ticked. The fill is now `block h-full`, so the markup carries what it needs. Rebuild the stylesheet with `aru view:build`.
- Tests: the checklist is driven in a real browser, key by key, by a paste and inside a fragment an htmx swap brought in, and held to the policy's own verdict for every line and to the width of the bar; a unit test holds the attributes the markup writes to the ones the behaviour in hesape reads.

## [v0.52.0](https://github.com/arandu-io/kyse/compare/v0.51.0...v0.52.0) - 2026-11-03

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.51.0...v0.52.0

## v0.52.0

- `MenuItem` gains `Dialog`, the id of a `Dialog` drawn beside the menu: the entry opens it, centred, and sends nothing itself, as `DialogTriggerProps.Dialog` does for a button. A menu entry had no way to ask "are you sure" except `HxConfirm`, which htmx turns into the browser's own confirm box. A link entry that names a dialog is left out of a boosted layout (`hx-boost="false"`), as a `DialogTrigger` link is.
- `HxConfirm` of `MenuItem` and `ButtonProps` is documented for what hesape v0.103.0 does with it: the view layer's script takes the question before htmx does and draws it as the confirmation dialog, with `data-confirm-title`, `data-confirm-label` and `data-confirm-cancel` to name its words, and the control's `destructive` variant on the button that goes through with it. Nothing in this library can reach the browser's own box any more. Use `Dialog` where the words matter.

## [v0.51.0](https://github.com/arandu-io/kyse/compare/v0.50.0...v0.51.0) - 2026-11-02

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.50.0...v0.51.0

## v0.51.0

- `FormDialog` gains `Forms`, for a dialog whose `Body` holds forms of its own: the panels of a `Tabs`, one per kind of change, each posting to its own address. The dialog is then a `div` panel and not a form (a form inside a form is not HTML, and the browser drops the inner one), so it draws no `Action`, no hidden fields and no footer; the title, the description, the scrolling body and the close button stay, and a page sent back with the dialog's name in `_dialog` still opens it. Without `Forms` the dialog is the one form it always was.

## [v0.50.0](https://github.com/arandu-io/kyse/compare/v0.49.1...v0.50.0) - 2026-11-02

- `DataTable` cards are chosen by the width of the screen and not of the block: under 600px wide the rows are cards (a media query, `max-[600px]:`), from 600px up the table is the table it was before v0.48.0, however narrow the block it sits in. v0.48.0 asked the width of the region with a container query, so a short list inside a card of a dashboard was drawn as cards on a desktop. The region is no longer a container (`@container/dt` is gone). It is how Ayra decides. Rebuild the stylesheet with `aru view:build`.
- `DataTable` gains `Mode`, with the names of Ayra's `DataTable`: `DataTableModeAuto` (the zero value, the rule above), `DataTableModeTable` (the table at every width, a short list that stays a list of lines on a phone) and `DataTableModeCards` (the cards at every width). The table mode draws no card class at all, so it is the markup and the pixels of the table before v0.48.0, with `data-label` on the cells as the only addition. The cards mode writes the card classes under `min-[0px]:`, a query that is always true, so they still sit after the table's own utilities in the stylesheet.
- `DataTable` gains `RowTemplate`, which draws an inert `<template data-part="row-template">` after the table: one row with the classes of the row and of each cell, the `data-column` and the `data-label`, in the mode the table is in. A script that puts a row into the list as it arrives, a read from a socket, clones it and writes the cells, so the row it adds is the row the server draws and no class is written in a script, where the stylesheet was never compiled with it.

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.49.1...v0.50.0

## [v0.49.1](https://github.com/arandu-io/kyse/compare/v0.49.0...v0.49.1) - 2026-11-02

- `DataTable` column menu: a column's label is written as the caller wrote it; the menu no longer capitalizes every word ("Data de renovação", not "Data De Renovação").
- `Progress`: the fill is a block, so its width draws; it used to be an empty track.
- `Badge`: the "default" variant draws as a pill, and a long label wraps inside the screen instead of widening the page.
- `Alert`: a long e-mail address or URL in the title or the message wraps inside the card.
- `Avatar`: a picture that fails to load no longer draws its alt text and the broken-image icon outside the circle.
- `StatCard`: the first column header is translated through the catalogue instead of the English "Name".
- `Masked`: a mask that takes digits only asks for the number pad (`inputmode="numeric"`); a mask with a letter in it (the alphanumeric CNPJ) keeps the full keyboard.

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.49.0...v0.49.1

## [v0.49.0](https://github.com/arandu-io/kyse/compare/v0.48.1...v0.49.0) - 2026-11-02

- `DataTable` cards: a cell with nothing in it is left out of the card, so the name of its column is no longer drawn in front of nothing. A row with no action used to show "Actions" alone in its card; it now shows neither the label nor a rule. The title of the card is drawn even when it is empty. `TableCell.Empty()` and `DataTableProps.CardCellFor(i, cell)` say so; `CardCell(i)` is unchanged. Above 600px nothing changes.
- `TableColumn.Actions` marks the column that holds a row's own controls (icon buttons, a link, a form). A card draws its cell as the foot of the card, at the end and under a rule, with no name in front of the controls, which stops the label from crowding the first icon; above 600px it is a column like any other, header included. It is how Ayra's `DataTable` draws a row's actions. Mark the column with `Actions: true` and rebuild the stylesheet with `aru view:build`.

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.48.1...v0.49.0

## [v0.48.1](https://github.com/arandu-io/kyse/compare/v0.48.0...v0.48.1) - 2026-11-01

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.48.0...v0.48.1

## [v0.48.0](https://github.com/arandu-io/kyse/compare/v0.47.0...v0.48.0) - 2026-11-01

- `DataTable` footer with a singular and a plural: `DataTableLines` gains `RowsOne` ("{total} row.") and `RowsMany` ("{total} rows."), `DataTableLines.RowsFor(total)` chooses between them, and the catalogue gains `kyse.data_table.rows_one` and `kyse.data_table.rows_many`. `Rows` and `kyse.data_table.rows` keep saying one sentence for both counts, so a project that passes them draws what it drew.
- `DataTable` as cards on a narrow screen: under 600px (a container query on the region, the breakpoint Ayra uses) every row is a card and nothing scrolls sideways. `TableColumn` gains `Primary` (the card's title; with none marked, the first column the card shows) and `CardHidden` (the column is left out of the card), the same names as Ayra's `DataColumn`. Each cell carries a `data-label` with the name of its column, the head stays as a bar of the controls that order the list, and a row's menu sits along the foot of its card. At 600px and above the table is unchanged.
- The behaviour props of a `DataTable` carry `rowsOne` beside `rows`; the footer of a list the browser narrows says one row in the singular with `hesape` carrying it, and says the plural as before with an older one.

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.47.0...v0.48.0

## [v0.47.0](https://github.com/arandu-io/kyse/compare/v0.46.0...v0.47.0) - 2026-11-01

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.46.0...v0.47.0

## [v0.46.0](https://github.com/arandu-io/kyse/compare/v0.45.0...v0.46.0) - 2026-10-29

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.45.0...v0.46.0

## [v0.45.0](https://github.com/arandu-io/kyse/compare/v0.44.0...v0.45.0) - 2026-10-25

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.44.0...v0.45.0

## [v0.44.0](https://github.com/arandu-io/kyse/compare/v0.43.1...v0.44.0) - 2026-10-25

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.43.1...v0.44.0

## [v0.43.1](https://github.com/arandu-io/kyse/compare/v0.43.0...v0.43.1) - 2026-10-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.43.0...v0.43.1

## [v0.43.0](https://github.com/arandu-io/kyse/compare/v0.42.1...v0.43.0) - 2026-10-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.42.1...v0.43.0

## [v0.42.1](https://github.com/arandu-io/kyse/compare/v0.42.0...v0.42.1) - 2026-10-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.42.0...v0.42.1

## [v0.42.0](https://github.com/arandu-io/kyse/compare/v0.41.0...v0.42.0) - 2026-10-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.41.0...v0.42.0

## [v0.41.0](https://github.com/arandu-io/kyse/compare/v0.40.0...v0.41.0) - 2026-10-20

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.40.0...v0.41.0

## [v0.40.0](https://github.com/arandu-io/kyse/compare/v0.39.0...v0.40.0) - 2026-10-20

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.39.0...v0.40.0

## [v0.39.0](https://github.com/arandu-io/kyse/compare/v0.38.0...v0.39.0) - 2026-10-20

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.38.0...v0.39.0

## [v0.38.0](https://github.com/arandu-io/kyse/compare/v0.37.0...v0.38.0) - 2026-10-19

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.37.0...v0.38.0

## [v0.37.0](https://github.com/arandu-io/kyse/compare/v0.36.0...v0.37.0) - 2026-10-19

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.36.0...v0.37.0

## [v0.36.0](https://github.com/arandu-io/kyse/compare/v0.35.0...v0.36.0) - 2026-10-19

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.35.0...v0.36.0

## [v0.34.1](https://github.com/arandu-io/kyse/compare/v0.34.0...v0.34.1) - 2026-10-10

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.34.0...v0.34.1

## [v0.34.0](https://github.com/arandu-io/kyse/compare/v0.33.0...v0.34.0) - 2026-10-10

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.33.0...v0.34.0

## [v0.33.0](https://github.com/arandu-io/kyse/compare/v0.32.0...v0.33.0) - 2026-10-09

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.32.0...v0.33.0

## [v0.32.0](https://github.com/arandu-io/kyse/compare/v0.31.0...v0.32.0) - 2026-10-09

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.31.0...v0.32.0

## [v0.31.0](https://github.com/arandu-io/kyse/compare/v0.30.0...v0.31.0) - 2026-10-09

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.30.0...v0.31.0

## [v0.30.0](https://github.com/arandu-io/kyse/compare/v0.29.4...v0.30.0) - 2026-10-01

## What's Changed
* Bump the gomod group with 2 updates by @dependabot[bot] in https://github.com/arandu-io/kyse/pull/2

## New Contributors
* @dependabot[bot] made their first contribution in https://github.com/arandu-io/kyse/pull/2

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.29.4...v0.30.0

## [v0.29.4](https://github.com/arandu-io/kyse/compare/v0.29.3...v0.29.4) - 2026-09-18

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.29.3...v0.29.4

## [v0.29.3](https://github.com/arandu-io/kyse/compare/v0.29.2...v0.29.3) - 2026-09-17

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.29.2...v0.29.3

## [v0.29.2](https://github.com/arandu-io/kyse/compare/v0.29.1...v0.29.2) - 2026-09-17

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.29.1...v0.29.2

## [v0.29.1](https://github.com/arandu-io/kyse/compare/v0.29.0...v0.29.1) - 2026-09-17

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.29.0...v0.29.1

## [v0.29.0](https://github.com/arandu-io/kyse/compare/v0.28.0...v0.29.0) - 2026-09-09

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.28.0...v0.29.0

## [v0.28.0](https://github.com/arandu-io/kyse/compare/v0.27.0...v0.28.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.27.0...v0.28.0

## [v0.27.0](https://github.com/arandu-io/kyse/compare/v0.26.0...v0.27.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.26.0...v0.27.0

## [v0.26.0](https://github.com/arandu-io/kyse/compare/v0.25.0...v0.26.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.25.0...v0.26.0

## [v0.25.0](https://github.com/arandu-io/kyse/compare/v0.24.0...v0.25.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.24.0...v0.25.0

## [v0.24.0](https://github.com/arandu-io/kyse/compare/v0.23.1...v0.24.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.23.1...v0.24.0

## [v0.23.1](https://github.com/arandu-io/kyse/compare/v0.23.0...v0.23.1) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.23.0...v0.23.1

## [v0.23.0](https://github.com/arandu-io/kyse/compare/v0.22.0...v0.23.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.22.0...v0.23.0

## [v0.22.0](https://github.com/arandu-io/kyse/compare/v0.21.0...v0.22.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.21.0...v0.22.0

## [v0.21.0](https://github.com/arandu-io/kyse/compare/v0.20.0...v0.21.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.20.0...v0.21.0

## [v0.20.0](https://github.com/arandu-io/kyse/compare/v0.19.0...v0.20.0) - 2026-09-08

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.19.0...v0.20.0

## [v0.19.0](https://github.com/arandu-io/kyse/compare/v0.18.0...v0.19.0) - 2026-09-07

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.18.0...v0.19.0

## [v0.18.0](https://github.com/arandu-io/kyse/compare/v0.17.0...v0.18.0) - 2026-09-05

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.17.0...v0.18.0

## [v0.17.0](https://github.com/arandu-io/kyse/compare/v0.15.2...v0.17.0) - 2026-09-04

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.15.2...v0.17.0

## [v0.17.0](https://github.com/arandu-io/kyse/compare/v0.16.0...v0.17.0) - 2026-09-03

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.16.0...v0.17.0

## [v0.16.0](https://github.com/arandu-io/kyse/compare/v0.15.2...v0.16.0) - 2026-09-03

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.15.2...v0.16.0

## [v0.15.2](https://github.com/arandu-io/kyse/compare/v0.15.1...v0.15.2) - 2026-09-02

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.15.1...v0.15.2

## [v0.15.1](https://github.com/arandu-io/kyse/compare/v0.15.0...v0.15.1) - 2026-09-01

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.15.0...v0.15.1

## [v0.15.0](https://github.com/arandu-io/kyse/compare/v0.14.0...v0.15.0) - 2026-09-01

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.14.0...v0.15.0

## [v0.14.0](https://github.com/arandu-io/kyse/compare/v0.13.0...v0.14.0) - 2026-08-29

The icon set is drawn in the duotone weight, which is the one the identity uses.
Still one weight and one name per icon: `Heart` is duotone, and there is no
`HeartDuotone` beside it. An icon is two paths of one colour -- the shade at a
fifth of `currentColor`, then the line over it -- so it still takes the colour
of the text around it and still needs no palette.

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.13.0...v0.14.0

## [v0.13.0](https://github.com/arandu-io/kyse/compare/v0.12.0...v0.13.0) - 2026-08-29

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.12.0...v0.13.0

## [v0.12.1](https://github.com/arandu-io/kyse/compare/v0.12.0...v0.12.1) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.12.0...v0.12.1

## [v0.12.0](https://github.com/arandu-io/kyse/compare/v0.11.0...v0.12.0) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.11.0...v0.12.0

## [v0.11.2](https://github.com/arandu-io/kyse/compare/v0.11.1...v0.11.2) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.11.1...v0.11.2

## [v0.11.1](https://github.com/arandu-io/kyse/compare/v0.11.0...v0.11.1) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.11.0...v0.11.1

## [v0.11.0](https://github.com/arandu-io/kyse/compare/v0.10.0...v0.11.0) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.10.0...v0.11.0

## [v0.10.0](https://github.com/arandu-io/kyse/compare/v0.9.0...v0.10.0) - 2026-08-23

**Full Changelog**: https://github.com/arandu-io/kyse/compare/v0.9.0...v0.10.0

Written by the release, not by hand. `.github/workflows/release.yml` publishes
the release a tag names and prepends what it said, so the tag, the release and
this file cannot disagree.

Tags cut before that workflow existed have no release and are not listed here.
`git tag` is where they are.
