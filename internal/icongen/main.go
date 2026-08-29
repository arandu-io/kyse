// Command icongen vendors the Phosphor icon set into Go source.
//
// It fetches one pinned commit of github.com/phosphor-icons/core, reads the
// duotone weight, and writes one exported function per icon into ../icons,
// along with the registry of them all as a test file under ../tests/Unit.
// Run it with `go generate ./...` from anywhere in this module.
//
// Two properties are worth the program existing rather than a shell script.
//
// The first is that nothing but path geometry crosses the boundary. Every
// duotone Phosphor SVG is one or two <path> elements inside a wrapper this
// program knows by heart. The shade may carry the exact opacity this weight
// declares; every other value is a d attribute. It refuses anything else --
// another element, another attribute, a different wrapper. A vendoring step
// that copied whole SVG files would happily copy a <script> or an onload= along
// with them one day; this one cannot, because the only values it carries out
// are path data.
//
// The second is that the pin is in the source. The commit below is what
// regenerates, so two runs a year apart produce the same bytes, and moving to a
// newer Phosphor is an edit somebody makes on purpose and reviews as a diff.
package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// The pin. Phosphor tags releases sporadically -- the tree here calls itself
// 2.1.1 in package.json and the newest tag in the repository is v2.0.8 -- so
// the commit is the version, and the name is only for people.
const (
	phosphorRepo    = "https://github.com/phosphor-icons/core"
	phosphorCommit  = "2b75f3ad12b420c9504ef05df8d2564a28f8500e"
	phosphorVersion = "2.1.1"
)

// weight is the one weight this library ships. Six weights would be six names
// for the same idea, and the icon set is a design decision the library makes
// once, like the stylesheet.
//
// Duotone, because that is the weight the identity draws in. It is still one
// weight and still one name per icon: Heart is duotone, and there is no
// HeartDuotone beside it.
const weight = "duotone"

// suffix is what upstream appends to a file name of this weight: the regular
// assets are heart.svg and these are heart-duotone.svg. It is trimmed before
// the name is looked up in the catalogue, which lists the icon and not the
// file.
const suffix = "-" + weight

// wrapper is every duotone Phosphor SVG, exactly. The generator asserts it
// rather than parsing: an upstream file that no longer looks like this is a
// change worth stopping for, not one worth accommodating.
//
// Two paths, and always in this order: the shade at one fifth opacity, then the
// line over it. That is what makes the weight duotone -- one colour, two
// strengths of it -- and it is why an icon of this weight carries two strings
// rather than one.
const (
	wrapperOpen  = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 256 256" fill="currentColor">`
	wrapperClose = `</svg>`
	pathOpen     = `<path d="`
	shadeClose   = `" opacity="0.2"/>`
	pathClose    = `"/>`
)

// pathData is what an SVG path is allowed to contain: the command letters, the
// digits, and the three separators. Anything else -- a quote, a backtick, an
// angle bracket -- means the file is not what this program thinks it is.
var pathData = regexp.MustCompile(`^[MmLlHhVvCcSsQqTtAaZz0-9., \-]+$`)

// icon is one entry of the set.
type icon struct {
	// Name is the Phosphor name, kebab-case: "arrow-right".
	Name string
	// GoName is the exported Go identifier: "ArrowRight".
	GoName string
	// Shade is the d attribute of the first path, the one drawn at a fifth of
	// the colour. Data is the second, the line over it. Two fields and not one
	// string of markup, for the reason pathsOf gives.
	Shade string
	// Data is the d attribute of the line path, and nothing else.
	Data string
	// Tags are the upstream synonyms, which end up in the doc comment so that
	// searching pkg.go.dev for "delete" finds Trash.
	Tags []string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("icongen: ")

	root, err := moduleRoot()
	if err != nil {
		log.Fatal(err)
	}

	src, err := os.MkdirTemp("", "icongen-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(src)

	log.Printf("fetching %s at %s", phosphorRepo, phosphorCommit[:12])
	if err := fetch(src); err != nil {
		log.Fatal(err)
	}

	icons, err := read(src)
	if err != nil {
		log.Fatal(err)
	}

	out := filepath.Join(root, "icons")
	// The registry is an external test file: it reaches nothing unexported, so
	// it belongs with the other external tests rather than beside the code.
	registryFile := filepath.Join(root, "tests", "Unit", "icons_all_test.go")
	if err := write(out, registryFile, src, icons); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %d icons to %s", len(icons), out)
}

// moduleRoot walks up for this module's go.mod.
//
// It does not trust the working directory. `go generate` runs the command from
// the package that carries the directive, but somebody running it by hand from
// the module root would otherwise write 1512 files there.
func moduleRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		body, err := os.ReadFile(filepath.Join(dir, "go.mod"))
		if err == nil {
			if strings.HasPrefix(string(body), "module github.com/arandu-io/kyse\n") {
				return dir, nil
			}
			return "", fmt.Errorf("%s is not the kyse module", dir)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no go.mod above the working directory")
		}
		dir = parent
	}
}

// fetch checks out the pinned commit into dir.
//
// Blobless and sparse: the repository is eighty megabytes and all of it but the
// selected weight, the catalogue and the licence is waste. This way it is nine
// megabytes and about two seconds, which is what makes the CI step that
// regenerates and compares affordable on every push.
func fetch(dir string) error {
	for _, args := range [][]string{
		{"init", "--quiet"},
		{"remote", "add", "origin", phosphorRepo},
		{"sparse-checkout", "set", "--no-cone", "/LICENSE", "/src/icons.ts", "/assets/" + weight + "/*"},
		{"fetch", "--quiet", "--depth", "1", "--filter=blob:none", "origin", phosphorCommit},
		{"checkout", "--quiet", "FETCH_HEAD"},
	} {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
		}
	}
	return nil
}

// read turns the checkout into the icon set, refusing anything unexpected.
func read(src string) ([]icon, error) {
	catalog, err := readCatalog(filepath.Join(src, "src", "icons.ts"))
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(src, "assets", weight)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var icons []icon
	seen := map[string]string{}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".svg") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".svg")
		trimmed, ok := strings.CutSuffix(name, suffix)
		if !ok {
			return nil, fmt.Errorf("%s: a file of the %s weight is expected to end in %s", e.Name(), weight, suffix)
		}
		name = trimmed

		body, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, err
		}
		shade, data, err := pathsOf(strings.TrimSpace(string(body)))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}

		goName, err := identifier(name)
		if err != nil {
			return nil, err
		}
		if other, taken := seen[goName]; taken {
			return nil, fmt.Errorf("%s and %s both want to be called %s", other, name, goName)
		}
		seen[goName] = name

		entry, known := catalog[name]
		if !known {
			return nil, fmt.Errorf("%s.svg is not in src/icons.ts", name)
		}
		// Upstream publishes the PascalCase name it intends. Deriving it here
		// and checking rather than reading it means the day Phosphor decides
		// that some icon is spelled differently, this stops instead of renaming
		// somebody's call site quietly.
		if entry.pascal != goName {
			return nil, fmt.Errorf("%s: upstream calls it %s, this derives %s", name, entry.pascal, goName)
		}
		delete(catalog, name)

		icons = append(icons, icon{Name: name, GoName: goName, Shade: shade, Data: data, Tags: entry.tags})
	}

	for name := range catalog {
		return nil, fmt.Errorf("src/icons.ts lists %s and assets/%s has no file for it", name, weight)
	}
	if len(icons) == 0 {
		return nil, fmt.Errorf("no icons in %s", dir)
	}

	sort.Slice(icons, func(i, j int) bool { return icons[i].GoName < icons[j].GoName })
	return icons, nil
}

// pathsOf pulls the two d attributes out of a duotone Phosphor SVG, and refuses
// everything that is not them.
//
// This is the whole security argument of the package. What is written into Go
// source is two strings that have been proved to hold nothing but path
// commands, so no icon can carry a <script>, an onclick, an xlink:href or a
// third element, whatever arrives upstream.
//
// The opacity is asserted rather than read. It is the same value on all of
// them, it is what the weight means, and reading it would make it a number this
// package carries fifteen hundred copies of.
func pathsOf(svg string) (shade, line string, err error) {
	body, ok := strings.CutPrefix(svg, wrapperOpen)
	if !ok {
		return "", "", fmt.Errorf("does not open with the wrapper this generator knows")
	}
	body, ok = strings.CutSuffix(body, wrapperClose)
	if !ok {
		return "", "", fmt.Errorf("does not close with %s", wrapperClose)
	}

	rest, ok := strings.CutPrefix(body, pathOpen)
	if !ok {
		return "", "", fmt.Errorf("does not open with a path: %.60q", body)
	}

	// A shade is usual and not universal. Upstream draws none where there is
	// nothing to fill -- cell-signal-none is an empty signal, and a fifth of a
	// colour over nothing is nothing. Those files carry the line alone, and an
	// icon whose shade is empty is a shape of this weight rather than a file
	// that stopped being one.
	if before, after, found := strings.Cut(rest, shadeClose); found {
		shade, rest = before, after
		if rest, ok = strings.CutPrefix(rest, pathOpen); !ok {
			return "", "", fmt.Errorf("holds a shade and no line after it: %.60q", body)
		}
	}

	line, ok = strings.CutSuffix(rest, pathClose)
	if !ok {
		return "", "", fmt.Errorf("holds more than the paths of this weight: %.60q", rest)
	}

	if shade != "" && !pathData.MatchString(shade) {
		return "", "", fmt.Errorf("the shade path holds something that is not path data: %.60q", shade)
	}
	if !pathData.MatchString(line) {
		return "", "", fmt.Errorf("the line path holds something that is not path data: %.60q", line)
	}
	return shade, line, nil
}

// identifier is the Phosphor name as an exported Go one: arrow-right becomes
// ArrowRight.
//
// The mapping is mechanical on purpose. A table of special cases would be a
// second thing to maintain, and a name somebody has to look up is a name they
// get wrong.
func identifier(name string) (string, error) {
	var b strings.Builder
	for _, word := range strings.Split(name, "-") {
		if word == "" {
			return "", fmt.Errorf("%q has an empty word in it", name)
		}
		b.WriteString(strings.ToUpper(word[:1]))
		b.WriteString(word[1:])
	}
	out := b.String()
	if c := out[0]; c < 'A' || c > 'Z' {
		return "", fmt.Errorf("%q does not start with a letter, so it cannot be an exported function", name)
	}
	return out, nil
}

// entry is what the upstream catalogue says about one icon.
type entry struct {
	pascal string
	tags   []string
}

var (
	reName    = regexp.MustCompile(`^\s*name:\s*"([^"]+)",`)
	rePascal  = regexp.MustCompile(`^\s*pascal_name:\s*"([^"]+)",`)
	reTagsOne = regexp.MustCompile(`^\s*tags:\s*\[(.*)],\s*$`)
	reTagsAll = regexp.MustCompile(`^\s*tags:\s*\[\s*$`)
	reQuoted  = regexp.MustCompile(`"([^"]*)"`)
)

// readCatalog reads src/icons.ts for the names and the tags.
//
// It scans lines rather than parsing TypeScript, which is enough because the
// file is machine-formatted and, more to the point, pinned: it cannot change
// under this program without somebody moving the commit above.
func readCatalog(path string) (map[string]entry, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	out := map[string]entry{}
	var name string
	var current entry
	var collecting bool

	flush := func() {
		if name != "" {
			out[name] = current
		}
		name, current = "", entry{}
	}

	for _, line := range strings.Split(string(body), "\n") {
		switch {
		case collecting:
			if strings.HasPrefix(strings.TrimSpace(line), "]") {
				collecting = false
				continue
			}
			current.tags = append(current.tags, quoted(line)...)
		case reName.MatchString(line):
			flush()
			name = reName.FindStringSubmatch(line)[1]
		case rePascal.MatchString(line):
			current.pascal = rePascal.FindStringSubmatch(line)[1]
		case reTagsOne.MatchString(line):
			current.tags = quoted(reTagsOne.FindStringSubmatch(line)[1])
		case reTagsAll.MatchString(line):
			collecting = true
		}
	}
	flush()

	if len(out) == 0 {
		return nil, fmt.Errorf("%s: read no icons, so its format changed", path)
	}
	return out, nil
}

// quoted pulls the double-quoted strings out of a line, dropping Phosphor's
// "*new*" marker, which is a release note rather than a synonym.
func quoted(line string) []string {
	var out []string
	for _, m := range reQuoted.FindAllStringSubmatch(line, -1) {
		if tag := m[1]; tag != "" && tag != "*new*" {
			out = append(out, tag)
		}
	}
	return out
}

// write emits the generated Go and the upstream licence. out takes the icon
// sources and the licence; registryFile is the one path outside it, the map the
// tests walk.
//
// The icons are split by initial letter. One file of a megabyte is a file no
// editor opens and no diff is readable in; twenty-six mean that Phosphor adding
// an icon shows up as a change to one of them.
func write(out, registryFile, src string, icons []icon) error {
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}

	// The licence travels with what was taken, unaltered. Copied rather than
	// retyped so that a change upstream is a diff here.
	licence, err := os.ReadFile(filepath.Join(src, "LICENSE"))
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(out, "LICENSE.md"), licence, 0o644); err != nil {
		return err
	}

	// A run that generates fewer letters than the tree has leaves the missing
	// ones behind, and a stale icons_q.go compiles perfectly.
	//
	// The sweep skips _test.go, and that is not caution: icons_*.go once matched
	// a test file sitting in this directory, and the first run of this generator
	// deleted the tests. The package still compiled, still passed, and reported
	// "no tests to run" in a line nobody reads.
	stale, err := filepath.Glob(filepath.Join(out, "icons_*.go"))
	if err != nil {
		return err
	}
	for _, f := range stale {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		if err := os.Remove(f); err != nil {
			return err
		}
	}

	groups := map[byte][]icon{}
	for _, i := range icons {
		groups[i.GoName[0]] = append(groups[i.GoName[0]], i)
	}
	for letter, group := range groups {
		name := fmt.Sprintf("icons_%c.go", letter|0x20)
		if err := emit(filepath.Join(out, name), source(group)); err != nil {
			return err
		}
	}
	if err := os.MkdirAll(filepath.Dir(registryFile), 0o755); err != nil {
		return err
	}
	return emit(registryFile, registry(icons))
}

// emit gofmts and writes one file.
func emit(path, body string) error {
	formatted, err := format.Source([]byte(body))
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	return os.WriteFile(path, formatted, 0o644)
}

// header is on every generated file, so that a file lifted out of the tree
// still says where its contents came from and under what licence. pkg is the
// package clause, which is not the same in the registry: that one is a test
// file and belongs to the package of the tests that walk it.
func header(pkg string) string {
	return fmt.Sprintf(`// Code generated by internal/icongen. DO NOT EDIT.
//
// Phosphor Icons %s, %s weight, from
// %s
// at commit %s.
//
// The paths are MIT, Copyright (c) 2023 Phosphor Icons. The licence is in
// LICENSE.md beside this file.

package %s
`, phosphorVersion, weight, phosphorRepo, phosphorCommit, pkg)
}

// source is one letter's worth of icons.
func source(group []icon) string {
	var b strings.Builder
	b.WriteString(header("icons"))
	b.WriteString("\nimport \"html/template\"\n")
	for _, i := range group {
		b.WriteString("\n")
		b.WriteString(comment(i))
		fmt.Fprintf(&b, "func %s(p Props) template.HTML {\n\treturn draw(p, `%s`, `%s`)\n}\n", i.GoName, i.Shade, i.Data)
	}
	return b.String()
}

// comment is the doc comment for one icon: the name it has upstream, and the
// synonyms, which is what makes the set findable. Nobody looking for a delete
// button searches for "trash".
func comment(i icon) string {
	line := fmt.Sprintf("%s is the Phosphor %q icon.", i.GoName, i.Name)
	if len(i.Tags) > 0 {
		line += " Tags: " + strings.Join(i.Tags, ", ") + "."
	}
	return wrap(line, 74)
}

// wrap breaks a doc comment at width columns, counting the "// ".
func wrap(text string, width int) string {
	var b strings.Builder
	column := 0
	for _, word := range strings.Fields(text) {
		switch {
		case column == 0:
			b.WriteString("// ")
			column = 3
		case column+1+len(word) > width:
			b.WriteString("\n// ")
			column = 3
		default:
			b.WriteString(" ")
			column++
		}
		b.WriteString(word)
		column += len(word)
	}
	b.WriteString("\n")
	return b.String()
}

// registry is the map the tests walk, and it is generated into a _test.go on
// purpose.
//
// A package-level map naming every icon would be reachable from any program
// that imports this package, and the linker would then keep all of them in
// every binary -- the exact cost the one-function-per-icon shape exists to
// avoid. Test files are compiled only into the test binary, so the tests can be
// exhaustive without anybody paying for it.
func registry(icons []icon) string {
	var b strings.Builder
	b.WriteString(header("unit"))
	b.WriteString(`
import (
	"html/template"

	"github.com/arandu-io/kyse/icons"
)

// all is every icon in the package, so that a test can hold the whole set to
// the same rule rather than a sample of it.
//
// It lives in a _test.go because a map of every icon is reachable from
// everything: in the package proper it would pin all of them into every binary
// that draws one, which is what this package is shaped to avoid.
var all = map[string]func(icons.Props) template.HTML{
`)
	for _, i := range icons {
		fmt.Fprintf(&b, "\t%q: icons.%s,\n", i.GoName, i.GoName)
	}
	b.WriteString("}\n")
	return b.String()
}
