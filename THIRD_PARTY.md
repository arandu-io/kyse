# Third-party notices

This library is written here, and three projects are upstream of it. Two of them
have their notice beside what was taken; the third has it here, because what was
taken from it is a shape rather than a file.

## shadcn-htmx

The component set, the `data-variant` / `data-size` attribute API and the HTMX
patterns are adapted from [shadcn-htmx](https://github.com/productdevbook/shadcn-htmx).
No file of it is in this repository — the components are Go functions compiled
from `.kyse.go` templates, and the upstream is TypeScript over Go
`html/template`. What travelled is the vocabulary: which components exist, what
each is called, which attribute carries the variant, and how the ARIA is put
together.

That is why this notice is a file of its own rather than a licence beside a
directory. There is no directory. The debt is real all the same, and MIT asks
for the notice to travel with a substantial portion of the work.

```
MIT License

Copyright (c) 2026 productdevbook

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

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
