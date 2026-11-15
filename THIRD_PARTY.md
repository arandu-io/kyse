# Third-party notices

This library is written here, and two projects are upstream of it. Each has its
notice beside what was taken.

## Basecoat

The stylesheet is [Basecoat](https://github.com/hunvreus/basecoat), MIT, and it
is vendored rather than adapted: the CSS files are in the project skeleton under
`resources/css/basecoat/`, with the upstream licence beside them in
`resources/css/basecoat/LICENSE.md`. It arrives as a directory of the project
and not as a dependency, so changing how a button looks is editing a file you
can see.

## Phosphor

The icons are [Phosphor](https://github.com/phosphor-icons/core), MIT. The path
geometry is generated into Go functions, and the licence is in
[icons/LICENSE.md](icons/LICENSE.md) beside them.

## What does not arrive

Nothing here comes through npm and nothing is fetched from a CDN. There is no
`package.json` in this repository or in a project that uses it, and the
Content-Security-Policy a page is served under is `script-src 'self'`, so a
script from another host would not run even if one were referenced.
