# Release Notes

## [Unreleased](https://github.com/arandu-io/kyse/compare/v0.48.1...HEAD)

- `DataTable` cards: a cell with nothing in it is left out of the card, so the name of its column is no longer drawn in front of nothing. A row with no action used to show "Actions" alone in its card; it now shows neither the label nor a rule. The title of the card is drawn even when it is empty. `TableCell.Empty()` and `DataTableProps.CardCellFor(i, cell)` say so; `CardCell(i)` is unchanged. Above 600px nothing changes.
- `TableColumn.Actions` marks the column that holds a row's own controls (icon buttons, a link, a form). A card draws its cell as the foot of the card, at the end and under a rule, with no name in front of the controls, which stops the label from crowding the first icon; above 600px it is a column like any other, header included. It is how Ayra's `DataTable` draws a row's actions. Mark the column with `Actions: true` and rebuild the stylesheet with `aru view:build`.

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
