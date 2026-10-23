package unit

import (
	stdhtml "html"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
	"github.com/arandu-io/kyse/icons"
)

// The sidebar is drawn in three shapes a reader can tell apart -- the full
// column, the slim strip, and the overlay on a narrow screen -- and a fourth
// that must not move at all: the column of a page that never asked for a mode.
//
// The goldens are the whole markup, compared byte for byte, because what the
// slim strip looks like is decided by classes, and a class that went missing
// renders, passes every assertion about content, and draws a column that does
// not fold. The assertions beside them say which bytes matter and why.
//
// KYSE_UPDATE_GOLDEN=1 rewrites the goldens of the modes. It never rewrites
// sidebar-plain.golden.html: that one is the markup the release before modes
// existed sent, and a change to it is a change for every page that upgrades.

// sidebarGroups are the entries every golden draws: one current, one
// disabled, one with no icon, and a group with no heading.
func sidebarGroups() []components.SidebarGroup {
	return []components.SidebarGroup{
		{Label: "Manage", Items: []components.SidebarItem{
			{Label: "Dashboard", URL: "/", Icon: icons.House(icons.Props{}), Current: true},
			{Label: "Plates", URL: "/plates", Icon: icons.Car(icons.Props{}), Disabled: true},
			{Label: "reports", URL: "/reports"},
		}},
		{Items: []components.SidebarItem{
			{Label: "Account", URL: "/account", Icon: icons.User(icons.Props{})},
		}},
	}
}

func sidebarWith(mode string, mobileOpen bool) string {
	return string(components.Sidebar(components.SidebarProps{
		ID:         "main",
		Label:      "Main menu",
		Header:     template.HTML("<b>Acme</b>"),
		Footer:     template.HTML("<b>v1</b>"),
		Groups:     sidebarGroups(),
		Mode:       mode,
		MobileOpen: mobileOpen,
	}))
}

// golden compares html against testdata/<name>, and rewrites it first when
// asked to and allowed to.
func golden(t *testing.T, name, html string, rewritable bool) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if rewritable && os.Getenv("KYSE_UPDATE_GOLDEN") == "1" {
		if err := os.WriteFile(path, []byte(html), 0o644); err != nil {
			t.Fatalf("writing %s: %v", path, err)
		}
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading %s: %v (KYSE_UPDATE_GOLDEN=1 writes it)", path, err)
	}
	if html != string(want) {
		t.Errorf("%s differs from what was rendered.\n--- rendered\n%s\n--- golden\n%s", path, html, want)
	}
}

// TestASidebarWithoutAModeSendsTheMarkupItSentBefore holds the compatibility
// promise to the byte: the golden was rendered by v0.41.0 from these props.
func TestASidebarWithoutAModeSendsTheMarkupItSentBefore(t *testing.T) {
	t.Parallel()

	html := string(components.Sidebar(components.SidebarProps{
		ID: "main", Header: template.HTML("<b>Acme</b>"), Footer: template.HTML("<b>v1</b>"),
		Groups: []components.SidebarGroup{
			{Label: "Manage", Items: []components.SidebarItem{
				{Label: "Posts", URL: "/posts", Current: true},
				{Label: "Users", URL: "/u", Disabled: true},
			}},
			{Items: []components.SidebarItem{{Label: "Home", URL: "/"}}},
		},
	}))
	golden(t, "sidebar-plain.golden.html", html, false)

	for _, absent := range []string{"sidebar-mode", `data-part="tooltip"`, "aria-label=\"Posts\"", "peer"} {
		if strings.Contains(html, absent) {
			t.Errorf("a sidebar with no mode carries %q, which only a mode draws", absent)
		}
	}
}

func TestTheFullSidebarGolden(t *testing.T) {
	t.Parallel()
	golden(t, "sidebar-full.golden.html", sidebarWith("full", false), true)
}

func TestTheSlimSidebarGolden(t *testing.T) {
	t.Parallel()
	golden(t, "sidebar-slim.golden.html", sidebarWith("slim", false), true)
}

func TestTheMobileDrawerGolden(t *testing.T) {
	t.Parallel()
	golden(t, "sidebar-mobile-drawer.golden.html", sidebarWith("slim", true), true)
}

// TestTheModeIsWrittenOnTheRoot covers the three answers: none, full, slim --
// and a value nobody defined, which is full rather than a strip with no rules
// behind it.
func TestTheModeIsWrittenOnTheRoot(t *testing.T) {
	t.Parallel()

	for mode, want := range map[string]string{
		"full":    `data-sidebar-mode="full"`,
		"slim":    `data-sidebar-mode="slim"`,
		"compact": `data-sidebar-mode="full"`,
	} {
		html := sidebarWith(mode, false)
		if !strings.Contains(html, want) {
			t.Errorf("Mode %q: the root does not carry %s", mode, want)
		}
		if strings.Count(html, "data-sidebar-mode=") != 1 {
			t.Errorf("Mode %q: the mode is written %d times, want once on the root", mode, strings.Count(html, "data-sidebar-mode="))
		}
	}
	if strings.Contains(sidebarWith("", false), "data-sidebar-mode") {
		t.Error("an empty Mode wrote a mode, so the page would start answering to a toggle it never asked for")
	}
}

// TestTheIconIsDrawnBeforeTheLabel in every mode, and in the plain column too:
// an icon is not something only the strip may have.
func TestTheIconIsDrawnBeforeTheLabel(t *testing.T) {
	t.Parallel()

	icon := string(icons.House(icons.Props{}))
	for _, mode := range []string{"", "full", "slim"} {
		html := string(components.Sidebar(components.SidebarProps{Mode: mode, Groups: []components.SidebarGroup{{
			Items: []components.SidebarItem{{Label: "Dashboard", URL: "/", Icon: template.HTML(icon)}},
		}}}))
		at, label := strings.Index(html, icon), strings.Index(html, ">Dashboard</span>")
		if at < 0 || label < 0 || at > label {
			t.Errorf("Mode %q: the icon is not drawn before the label (icon at %d, label at %d)", mode, at, label)
		}
		if !strings.Contains(html, `aria-hidden="true"`) {
			t.Errorf("Mode %q: an icon passed with no Label of its own lost its aria-hidden", mode)
		}
	}
}

var anchorTag = regexp.MustCompile(`<a\s[^>]*data-part="link"[^>]*>`)

// TestTheSlimStripNamesEveryLinkByItsLabel: the label leaves the eye and stays
// with the link, as its aria-label and as the hidden text, and the tip beside
// it repeats it for the eye and for nobody else.
func TestTheSlimStripNamesEveryLinkByItsLabel(t *testing.T) {
	t.Parallel()

	html := flat(sidebarWith("slim", false))
	links := anchorTag.FindAllString(html, -1)
	if len(links) != 4 {
		t.Fatalf("drew %d links, want 4", len(links))
	}
	for i, label := range []string{"Dashboard", "Plates", "reports", "Account"} {
		if !strings.Contains(links[i], `aria-label="`+label+`"`) {
			t.Errorf("the link for %q is not named by its label in the strip: %s", label, links[i])
		}
		tip := `aria-hidden="true">` + label + `</span>`
		if !strings.Contains(html, tip) {
			t.Errorf("no tip repeats %q, aria-hidden, beside its link", label)
		}
		hidden := `class="md:in-data-[sidebar-mode=slim]:sr-only">` + label + `</span>`
		if !strings.Contains(html, hidden) {
			t.Errorf("the text of %q is not kept as the link's name when the strip hides it", label)
		}
	}
	if got := strings.Count(html, `data-part="tooltip"`); got != 4 {
		t.Errorf("drew %d tips, want one per entry", got)
	}
	if strings.Contains(html, `role="tooltip"`) {
		t.Error("the tip is a tooltip role, and it would be announced as a second name")
	}

	full := sidebarWith("full", false)
	if strings.Contains(full, `aria-label="Dashboard"`) {
		t.Error("the full column names its links twice: the label is on the screen")
	}
}

// TestTheTipAppearsOnHoverAndOnKeyboardFocus, and only in the strip and only
// on a wide screen. The tip's look is the tooltip's own class.
func TestTheTipAppearsOnHoverAndOnKeyboardFocus(t *testing.T) {
	t.Parallel()

	for _, side := range []string{"", "right"} {
		html := string(components.Sidebar(components.SidebarProps{Mode: "slim", Side: side, Groups: sidebarGroups()}))
		for _, class := range []string{
			"tooltip-content", "hidden", "md:in-data-[sidebar-mode=slim]:block",
			"peer-hover:opacity-100", "peer-focus-visible:opacity-100",
		} {
			if !hasClassOnPart(html, "tooltip", class) {
				t.Errorf("Side %q: the tip does not carry %q", side, class)
			}
		}
		if !hasClassOnPart(html, "link", "peer") {
			t.Errorf("Side %q: the link is not the peer the tip watches", side)
		}
	}
	right := string(components.Sidebar(components.SidebarProps{Mode: "slim", Side: "right", Groups: sidebarGroups()}))
	if !hasClassOnPart(right, "tooltip", "right-[calc(100%+0.75rem)]") {
		t.Error("on the right the tip is not drawn on the side away from the edge")
	}
}

// TestTheCurrentEntryStaysMarkedInTheStrip: aria-current is what both the
// stylesheet and a screen reader read, in every mode.
func TestTheCurrentEntryStaysMarkedInTheStrip(t *testing.T) {
	t.Parallel()

	html := sidebarWith("slim", false)
	links := anchorTag.FindAllString(html, -1)
	if len(links) == 0 || !strings.Contains(links[0], `aria-current="page"`) {
		t.Fatalf("the current entry is not marked in the strip: %v", links)
	}
	if strings.Count(html, `aria-current="page"`) != 1 {
		t.Errorf("aria-current is written %d times, want once", strings.Count(html, `aria-current="page"`))
	}
	if !strings.Contains(links[1], `aria-disabled="true"`) {
		t.Error("the disabled entry stopped saying so in the strip")
	}
}

// TestAGroupHeadingBecomesALineInTheStrip and keeps its words: the heading
// keeps the text, clipped, and the group keeps its name.
func TestAGroupHeadingBecomesALineInTheStrip(t *testing.T) {
	t.Parallel()

	html := sidebarWith("slim", false)
	for _, class := range []string{
		"md:in-data-[sidebar-mode=slim]:h-0", "md:in-data-[sidebar-mode=slim]:border-t",
		"md:in-data-[sidebar-mode=slim]:overflow-hidden",
	} {
		if !hasClassOnPart(html, "group-label", class) {
			t.Errorf("the group heading does not carry %q in the strip", class)
		}
	}
	if !strings.Contains(html, `>Manage</h3>`) || !strings.Contains(html, `aria-label="Manage"`) {
		t.Error("the group lost its name in the strip")
	}
}

// TestAnEntryWithNoIconDrawsItsInitialInTheStrip rather than an empty square.
func TestAnEntryWithNoIconDrawsItsInitialInTheStrip(t *testing.T) {
	t.Parallel()

	html := sidebarWith("slim", false)
	if !strings.Contains(html, `aria-hidden="true">R</span>`) {
		t.Error("the entry with no icon draws no letter in the strip")
	}
	if got := (components.SidebarItem{Label: "  émile"}).Initial(); got != "É" {
		t.Errorf("Initial() = %q, want %q", got, "É")
	}
	if got := (components.SidebarItem{}).Initial(); got != "" {
		t.Errorf("Initial() of an empty label = %q, want empty", got)
	}
}

// TestTheNarrowScreenIsTheOverlayInEitherMode: every rule the strip adds waits
// on a wide screen, so the overlay is drawn with its labels whatever the mode.
func TestTheNarrowScreenIsTheOverlayInEitherMode(t *testing.T) {
	t.Parallel()

	html := sidebarWith("slim", true)
	if !strings.Contains(html, `data-initial-mobile-open="true"`) {
		t.Fatal("MobileOpen no longer opens the overlay")
	}
	for _, class := range classesIn(html) {
		if strings.Contains(class, "sidebar-mode") && !strings.HasPrefix(class, "md:") {
			t.Errorf("%q acts on a narrow screen, where the column is the overlay with its labels", class)
		}
	}
}

// TestTheStripMovesWithAShortTransitionAndStillUnderReducedMotion.
func TestTheStripMovesWithAShortTransitionAndStillUnderReducedMotion(t *testing.T) {
	t.Parallel()

	for _, side := range []string{"", "right"} {
		html := string(components.Sidebar(components.SidebarProps{Mode: "full", Side: side, Groups: sidebarGroups()}))
		for part, classes := range map[string][]string{
			"nav":  {"transition-[transform,translate,width]", "duration-200", "motion-reduce:transition-none"},
			"root": {"[&+*]:duration-200", "motion-reduce:[&+*]:transition-none"},
		} {
			for _, class := range classes {
				if !hasClassOnPart(html, part, class) {
					t.Errorf("Side %q: %s does not carry %q", side, part, class)
				}
			}
		}
	}
	left := sidebarWith("slim", false)
	if !hasClassOnPart(left, "root", "md:data-[sidebar-mode=slim]:[&:not([aria-hidden=true])+*]:ml-[var(--sidebar-slim-width,3rem)]") {
		t.Error("the page beside a left strip is not given the strip's width")
	}
	right := string(components.Sidebar(components.SidebarProps{Mode: "slim", Side: "right", Groups: sidebarGroups()}))
	if !hasClassOnPart(right, "root", "md:data-[sidebar-mode=slim]:[&:not([aria-hidden=true])+*]:mr-[var(--sidebar-slim-width,3rem)]") {
		t.Error("the page beside a right strip is not given the strip's width")
	}
	if strings.Contains(right, "ml-[var(--sidebar-slim-width") {
		t.Error("a right strip pushes the page from the left as well")
	}
}

// TestTheSlimStripCarriesNoScriptAndNoStyle: everything it does is a class
// and an attribute, under script-src 'self' and style-src 'self'.
func TestTheSlimStripCarriesNoScriptAndNoStyle(t *testing.T) {
	t.Parallel()

	html := sidebarWith("slim", false)
	for _, forbidden := range []string{"<script", " style=", " onclick", " x-", "@click", "hx-on", "localStorage"} {
		if strings.Contains(html, forbidden) {
			t.Errorf("the strip carries %q", forbidden)
		}
	}
}

var spaces = regexp.MustCompile(`\s+`)

// flat is html with every run of whitespace made one space and none before
// the end of a tag, so an assertion can read a tag whose attributes the
// component writes one per line.
func flat(html string) string {
	return strings.ReplaceAll(spaces.ReplaceAllString(html, " "), " >", ">")
}

var partTag = regexp.MustCompile(`<[a-z0-9]+\s[^>]*>`)
var classAttr = regexp.MustCompile(`class="([^"]*)"`)

// hasClassOnPart is whether some element with data-part=part carries class.
func hasClassOnPart(html, part, class string) bool {
	for _, tag := range partTag.FindAllString(html, -1) {
		if !strings.Contains(tag, `data-part="`+part+`"`) {
			continue
		}
		m := classAttr.FindStringSubmatch(tag)
		if m == nil {
			continue
		}
		for _, c := range strings.Fields(stdhtml.UnescapeString(m[1])) {
			if c == class {
				return true
			}
		}
	}
	return false
}

// classesIn is every class written anywhere in html.
func classesIn(html string) []string {
	var out []string
	for _, m := range classAttr.FindAllStringSubmatch(html, -1) {
		out = append(out, strings.Fields(stdhtml.UnescapeString(m[1]))...)
	}
	return out
}
