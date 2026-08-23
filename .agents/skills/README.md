# Skills

Procedures an assistant follows when working on this component library.

They live in `.agents/skills/<name>/SKILL.md`, which is the path the coding
assistants read from — Cursor, Codex, Cline, Copilot, Gemini CLI, Amp, OpenCode,
Warp, Zed and the rest all look there. It is one directory rather than a file
per vendor, so a skill written once is read by whatever this project is being
written with.

Each file opens with frontmatter carrying a `name` and a `description`. The
description is what a tool reads to decide whether the skill is relevant, so it
names the situation you are in rather than the topic it covers.

| skill | when it fires |
| --- | --- |
| `kyse-component` | adding a component, changing a prop, or a component that does not compile |
| `kyse-markup` | writing the markup inside a component: classes, attributes, escaping, behaviour |
| `kyse-icons` | drawing an icon, adding one, or moving the icon set |

## Why these exist

The audience of this repository is somebody changing the library, not somebody
using it. The questions are different ones: how a component is added, what a
component is allowed to do, and what fails if it does something else.

A model asked to write a component here fills the gap with the component
libraries it does know and produces a stylesheet, a `<script>`, an `x-on:click`
and a `class` invented on the spot. None of those work: there is no CSS in this
repository, no JavaScript, a test that reads the sources and fails on a
directive, and a stylesheet that ships elsewhere and knows only the class names
it was written for. `AGENTS.md` at the root lists what each of those maps to.

The rest of the answer is that the library is built to be checked rather than
trusted. The smoke table is compared against the directory, the compiled files
are compared against their sources, every prop is rendered with a script tag in
it, and the whole icon set is parsed as XML on every run. An assistant that runs
`go test -race ./...` is not guessing.

## Adding your own

A skill in this directory is yours and travels with the repository. Keep it a
procedure rather than a description: a file that says "read the documentation"
never changes what anybody does.
