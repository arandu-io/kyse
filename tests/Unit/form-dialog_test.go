package unit

import (
	"html/template"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/kyse/components"
)

// refusedPage is the page a refused form is sent back to: the messages by
// field and what was typed, the way view.Page answers both.
type refusedPage struct {
	errors map[string]string
	old    map[string]string
}

func (p refusedPage) FieldError(name string) string { return p.errors[name] }

func (p refusedPage) OldOr(name, fallback string) string {
	if value, typed := p.old[name]; typed {
		return value
	}
	return fallback
}

func discountDialog(page components.Page) components.FormDialogProps {
	return components.FormDialogProps{
		ID:          "new-discount",
		Title:       "New discount",
		Description: "It applies at checkout.",
		Action:      "/discounts",
		Token:       "t",
		Page:        page,
		Body: []template.HTML{
			components.Field(components.FieldProps{Name: "name", Label: "Name", Page: page}),
			components.Field(components.FieldProps{Name: "amount", Label: "Amount", Page: page}),
		},
	}
}

// The page a refused form lands on draws every dialog it has, and only the one
// whose form was refused opens. That one carries its messages and what was
// typed, because they are drawn by the fields inside it, which asked the page
// themselves.
func TestTheRefusedFormDialogOpensItselfAndNoOther(t *testing.T) {
	page := refusedPage{
		errors: map[string]string{"name": "The name is already taken."},
		old:    map[string]string{components.FormDialogField: "new-discount", "name": "BLACKFRIDAY"},
	}

	refused := string(components.FormDialog(discountDialog(page)))
	other := discountDialog(page)
	other.ID = "new-benefit"
	untouched := string(components.FormDialog(other))

	for _, want := range []string{
		`data-kyse-behavior="form-dialog"`,
		`data-kyse-props="{&#34;open&#34;:true}"`,
		"\n\t\topen\n",
		`aria-invalid="true"`,
		"The name is already taken.",
		`value="BLACKFRIDAY"`,
	} {
		if !strings.Contains(refused, want) {
			t.Errorf("the refused dialog does not carry %q:\n%s", want, refused)
		}
	}
	for _, unwanted := range []string{"data-kyse-behavior", "\n\t\topen\n"} {
		if strings.Contains(untouched, unwanted) {
			t.Errorf("a dialog whose form was not the one refused carries %q:\n%s", unwanted, untouched)
		}
	}
	if (components.FormDialogProps{ID: "x"}).Opened() {
		t.Error("a dialog with no page opened by itself")
	}
	if !(components.FormDialogProps{ID: "x", Open: true}).Opened() {
		t.Error("Open did not open it")
	}
}

// The ID the page is asked about is the one the form posts, in the field the
// page is asked about. If the two drifted apart every refusal would land on a
// closed dialog, and the page would look as if nothing had been sent.
func TestTheFormDialogPostsItsOwnID(t *testing.T) {
	html := string(components.FormDialog(discountDialog(nil)))
	want := `<input type="hidden" name="` + components.FormDialogField + `" value="new-discount">`
	if !strings.Contains(html, want) {
		t.Fatalf("the form does not post its ID as %s:\n%s", want, html)
	}
	if strings.Contains(html, "data-kyse-behavior") {
		t.Errorf("a dialog waiting for its trigger writes a behaviour it does not need:\n%s", html)
	}
}

var submitButtons = regexp.MustCompile(`(?s)<button\s[^>]*type="submit"[^>]*>`)

// Enter in a field presses the first submit button the form owns. Cancel and
// close are submit buttons too -- of the dismiss form -- so the one Enter
// reaches has to be the one that sends the form, and the two ways out have to
// belong to the other form.
func TestEnterSendsTheFormAndTheWaysOutCloseIt(t *testing.T) {
	html := string(components.FormDialog(discountDialog(nil)))

	buttons := submitButtons.FindAllString(html, -1)
	if len(buttons) != 3 {
		t.Fatalf("want 3 submit buttons (submit, cancel, close), got %d:\n%s", len(buttons), html)
	}
	if !strings.Contains(buttons[0], `data-part="submit"`) || strings.Contains(buttons[0], "form=") {
		t.Errorf("the first submit button is not the form's own submit: %s", buttons[0])
	}
	for _, button := range buttons[1:] {
		if !strings.Contains(button, `form="new-discount-dismiss"`) {
			t.Errorf("a way out is owned by the form it should leave: %s", button)
		}
	}
	if !strings.Contains(html, `<form method="dialog" id="new-discount-dismiss" hidden></form>`) {
		t.Errorf("there is no dismiss form for the ways out to submit:\n%s", html)
	}
	if strings.Count(html, "<form") != 2 || strings.Count(html, "</form>") != 2 {
		t.Errorf("want the form and the dismiss form, side by side and not nested:\n%s", html)
	}
}

func TestTheFormDialogSizes(t *testing.T) {
	for size, want := range map[string]string{
		"":   "sm:max-w-lg",
		"md": "sm:max-w-lg",
		"sm": "sm:max-w-sm",
		"lg": "sm:max-w-2xl",
		"xl": "sm:max-w-lg",
	} {
		props := discountDialog(nil)
		props.Size = size
		html := string(components.FormDialog(props))
		if !strings.Contains(html, want) {
			t.Errorf("Size %q: want %s:\n%s", size, want, html)
		}
		if !strings.Contains(html, "max-sm:h-dvh") {
			t.Errorf("Size %q: the dialog does not take the whole screen of a phone", size)
		}
	}
}

// An alert is a question, and stays a card on a phone; the role is what tells
// assistive technology it has to be answered, and the way out is where the
// cursor starts.
func TestTheDestructiveFormDialog(t *testing.T) {
	props := discountDialog(nil)
	props.Alert = true
	props.SubmitVariant = "destructive"
	props.SubmitLabel = "Cancel subscription"
	props.CancelLabel = "Keep it"
	html := string(components.FormDialog(props))

	for _, want := range []string{
		`class="alert-dialog"`,
		`role="alertdialog"`,
		`data-variant="destructive"`,
		"autofocus",
		"Cancel subscription",
		"Keep it",
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the alert does not carry %q:\n%s", want, html)
		}
	}
	for _, unwanted := range []string{"max-sm:h-dvh", "data-size=", `data-part="close"`} {
		if strings.Contains(html, unwanted) {
			t.Errorf("the alert carries %q:\n%s", unwanted, html)
		}
	}
}

func TestTheFormDialogNamesWhatItDescribes(t *testing.T) {
	html := string(components.FormDialog(discountDialog(nil)))
	for _, want := range []string{
		`aria-labelledby="new-discount-title"`,
		`id="new-discount-title"`,
		`aria-describedby="new-discount-description"`,
		`id="new-discount-description"`,
		`aria-label="Close"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("want %q:\n%s", want, html)
		}
	}

	props := discountDialog(nil)
	props.Description = ""
	if html := string(components.FormDialog(props)); strings.Contains(html, "aria-describedby") {
		t.Errorf("a dialog with no description points at one:\n%s", html)
	}
}

func TestTheFormDialogSendsTheTokenOnlyWhenItPosts(t *testing.T) {
	for _, tc := range []struct {
		method   string
		posts    bool
		override string
	}{
		{"", true, ""},
		{"post", true, ""},
		{"patch", true, "PATCH"},
		{"PUT", true, "PUT"},
		{"delete", true, "DELETE"},
		{"get", false, ""},
	} {
		props := discountDialog(nil)
		props.Method = tc.method
		props.Token = "secret-token"
		html := string(components.FormDialog(props))
		if got := strings.Contains(html, "secret-token"); got != tc.posts {
			t.Errorf("Method %q: token written = %v, want %v", tc.method, got, tc.posts)
		}
		hasOverride := strings.Contains(html, `name="_method" value="`+tc.override+`"`)
		if (tc.override != "") != hasOverride {
			t.Errorf("Method %q: _method %q written = %v:\n%s", tc.method, tc.override, hasOverride, html)
		}
	}

	props := discountDialog(nil)
	props.Multipart = true
	if html := string(components.FormDialog(props)); !strings.Contains(html, `enctype="multipart/form-data"`) {
		t.Errorf("a multipart form dialog is not encoded as one:\n%s", html)
	}
}

// Without the script the trigger is a link to the page with the form on it;
// without a page behind it, it is a button.
func TestTheDialogTriggerFallsBackToThePage(t *testing.T) {
	link := string(components.DialogTrigger(components.DialogTriggerProps{
		Dialog: "new-discount", Label: "New discount", Href: "/discounts/new",
	}))
	for _, want := range []string{`<a`, `href="/discounts/new"`, `data-dialog-open="new-discount"`, `aria-haspopup="dialog"`, `class="btn"`} {
		if !strings.Contains(link, want) {
			t.Errorf("the link trigger does not carry %q:\n%s", want, link)
		}
	}

	button := string(components.DialogTrigger(components.DialogTriggerProps{
		Dialog: "confirm", Label: "Delete", Variant: "destructive",
	}))
	for _, want := range []string{`<button`, `type="button"`, `data-dialog-open="confirm"`, `data-variant="destructive"`} {
		if !strings.Contains(button, want) {
			t.Errorf("the button trigger does not carry %q:\n%s", want, button)
		}
	}
	if strings.Contains(button, "href=") {
		t.Errorf("a trigger with no page behind it links somewhere:\n%s", button)
	}

	iconOnly := string(components.DialogTrigger(components.DialogTriggerProps{
		Dialog: "new-discount", Label: "New discount", Href: "/discounts/new",
		Icon: template.HTML("<svg></svg>"), IconOnly: true,
	}))
	if !strings.Contains(iconOnly, `aria-label="New discount"`) {
		t.Errorf("an icon-only trigger has no name:\n%s", iconOnly)
	}
}

// A form whose answer shows a secret once has to be posted by the browser and
// answered as a document, never swapped in by a boosted layout. The caller
// cannot write hx-boost itself, so the dialog writes it on the form.
func TestAFormDialogCanKeepItsFormOutOfABoostedLayout(t *testing.T) {
	plain := string(components.FormDialog(discountDialog(nil)))
	if strings.Contains(plain, "hx-boost") {
		t.Errorf("a form dialog opted out of boosting without being asked:\n%s", plain)
	}

	props := discountDialog(nil)
	props.NoBoost = true
	html := string(components.FormDialog(props))
	form := html[strings.Index(html, `<form`):strings.Index(html, `<input type="hidden"`)]
	if !strings.Contains(form, `hx-boost="false"`) {
		t.Errorf("NoBoost did not reach the form:\n%s", form)
	}
	if strings.Count(html, "hx-boost") != 1 {
		t.Errorf("hx-boost written somewhere other than the form:\n%s", html)
	}
}

// The case that drew an empty space where a dialog should have been: the view
// tried to write hx-boost through Attrs. It now stops the render and says so.
func TestAFormDialogGivenARefusedAttributeSaysSo(t *testing.T) {
	props := discountDialog(nil)
	props.Parts = components.Parts{"content": {Attrs: components.Attrs{"hx-boost": "false"}}}

	defer func() {
		refused, ok := recover().(*components.RefusedAttribute)
		if !ok {
			t.Fatal("the dialog rendered instead of refusing hx-boost")
		}
		if refused.Part != "content" || refused.Name != "hx-boost" {
			t.Errorf("refused %q on %q, want hx-boost on content", refused.Name, refused.Part)
		}
		if !strings.Contains(refused.Error(), `"hx-boost"`) || !strings.Contains(refused.Error(), "drawn nothing") {
			t.Errorf("the message does not say what was refused and what it cost: %s", refused.Error())
		}
	}()
	components.FormDialog(props)
}
