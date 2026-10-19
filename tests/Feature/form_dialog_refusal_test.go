package feature

import (
	"html"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"

	hhttp "github.com/arandu-io/hesape/http"
	"github.com/arandu-io/hesape/session"
	"github.com/arandu-io/hesape/validation"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

var formFieldPattern = regexp.MustCompile(`<input[^>]*\bname="([^"]+)"[^>]*\bvalue="([^"]*)"`)

// discountsPage is the list a "New discount" dialog sits on, drawn for a page
// the way a controller draws it.
func discountsPage(page components.Page) string {
	return string(components.FormDialog(components.FormDialogProps{
		ID:     "new-discount",
		Title:  "New discount",
		Action: "/discounts",
		Token:  "t",
		Page:   page,
		Body: []template.HTML{
			components.Field(components.FieldProps{Name: "name", Label: "Name", Page: page}),
		},
	}))
}

// TestARefusedFormDialogComesBackOpen runs the refusal the way a browser and
// the framework run it, with nothing added to either: the form is posted with
// what its markup carries, the handler refuses it with hesape's Reject, the
// browser follows the redirect with the flash cookie, and the page that lands
// is drawn from what the flash brought back.
//
// The dialog has to come back open, with the message and the typed value in
// it. Each half is owned by a different repository -- the field the form posts
// is kyse's, what the flash keeps is hesape's -- and a test of either alone
// would pass the day the other started dropping it.
func TestARefusedFormDialogComesBackOpen(t *testing.T) {
	flash := session.NewFlash([]byte("0123456789abcdef0123456789abcdef"), false)

	markup := discountsPage(nil)
	if strings.Contains(markup, "data-kyse-behavior") {
		t.Fatalf("the dialog opens on a page nobody was refused on:\n%s", markup)
	}

	values := url.Values{}
	for _, field := range formFieldPattern.FindAllStringSubmatch(markup, -1) {
		values.Set(html.UnescapeString(field[1]), html.UnescapeString(field[2]))
	}
	values.Set("name", "BLACKFRIDAY")

	post := httptest.NewRequest(http.MethodPost, "/discounts", strings.NewReader(values.Encode()))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Referer", "http://example.com/discounts")
	refused := httptest.NewRecorder()
	hhttp.Reject(refused, post, flash, validation.Errors{"name": {"The name is already taken."}})

	if location := refused.Header().Get("Location"); location != "/discounts" {
		t.Fatalf("the refusal sent the browser to %q, want the list the dialog sits on", location)
	}

	back := httptest.NewRequest(http.MethodGet, "/discounts", nil)
	back.Header.Set("Accept", "text/html")
	for _, cookie := range refused.Result().Cookies() {
		back.AddCookie(cookie)
	}
	got, ok := flash.Read(httptest.NewRecorder(), back)
	if !ok {
		t.Fatal("the page the refusal landed on found nothing in the flash")
	}

	landed := discountsPage(view.Page{Errors: got.Errors, Old: got.Old})
	for _, want := range []string{
		`data-kyse-behavior="form-dialog"`,
		`aria-invalid="true"`,
		"The name is already taken.",
		`value="BLACKFRIDAY"`,
	} {
		if !strings.Contains(landed, want) {
			t.Errorf("the dialog came back without %q:\n%s", want, landed)
		}
	}
}

// twoDialogs is a subscription page with two dialogs that each ask for a
// reason, written the way DialogPage says they are.
func twoDialogs(page components.Page) string {
	dialog := func(id, action string) string {
		scoped := components.DialogPage(page, id)
		return string(components.FormDialog(components.FormDialogProps{
			ID: id, Title: id, Action: action, Token: "t", Page: page,
			Body: []template.HTML{
				components.HiddenInput(components.HiddenInputProps{Name: "subscription", Value: "sub_1", Page: scoped}),
				components.Field(components.FieldProps{Name: "reason", ID: id + "-reason", Label: "Reason", Page: scoped}),
			},
		}))
	}
	return dialog("cancel-subscription", "/subscriptions/cancel") + dialog("pause-subscription", "/subscriptions/pause")
}

// TestOnlyTheRefusedDialogComesBackWithTheInput is the same round trip with
// two dialogs whose fields share a name. Before the fields were scoped, the
// message and the typed reason came back in both.
func TestOnlyTheRefusedDialogComesBackWithTheInput(t *testing.T) {
	flash := session.NewFlash([]byte("0123456789abcdef0123456789abcdef"), false)

	markup := twoDialogs(nil)
	if strings.Count(markup, `id="cancel-subscription-reason"`) != 1 || strings.Count(markup, `id="pause-subscription-reason"`) != 1 {
		t.Fatalf("the two reasons do not carry an id each:\n%s", markup)
	}
	cancel := markup[:strings.Index(markup, `id="pause-subscription"`)]
	values := url.Values{}
	for _, field := range formFieldPattern.FindAllStringSubmatch(cancel, -1) {
		values.Set(html.UnescapeString(field[1]), html.UnescapeString(field[2]))
	}
	values.Set("reason", "Too expensive")

	post := httptest.NewRequest(http.MethodPost, "/subscriptions/cancel", strings.NewReader(values.Encode()))
	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.Header.Set("Referer", "http://example.com/subscriptions")
	refused := httptest.NewRecorder()
	hhttp.Reject(refused, post, flash, validation.Errors{"reason": {"Say why in more words."}})

	back := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
	back.Header.Set("Accept", "text/html")
	for _, cookie := range refused.Result().Cookies() {
		back.AddCookie(cookie)
	}
	got, ok := flash.Read(httptest.NewRecorder(), back)
	if !ok {
		t.Fatal("the page the refusal landed on found nothing in the flash")
	}

	landed := twoDialogs(view.Page{Errors: got.Errors, Old: got.Old})
	at := strings.Index(landed, `id="pause-subscription"`)
	first, second := landed[:at], landed[at:]
	for _, want := range []string{"Say why in more words.", `value="Too expensive"`, `data-kyse-behavior="form-dialog"`} {
		if !strings.Contains(first, want) {
			t.Errorf("the refused dialog came back without %q:\n%s", want, first)
		}
	}
	for _, unwanted := range []string{"Say why in more words.", "Too expensive", "data-kyse-behavior", "aria-invalid"} {
		if strings.Contains(second, unwanted) {
			t.Errorf("the dialog nobody submitted came back with %q:\n%s", unwanted, second)
		}
	}
}
