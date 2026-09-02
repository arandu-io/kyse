package feature

import (
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	fhttp "github.com/arandu-io/framework/http"
	frameworkmiddleware "github.com/arandu-io/framework/http/middleware"
	"github.com/arandu-io/framework/security"
	"github.com/arandu-io/hesape/http/middleware"
	"github.com/arandu-io/kyse/components"
)

var (
	confirmFormPattern = regexp.MustCompile(`(?s)<form method="([^"]+)" action="([^"]+)">(.*?)<button\s+data-part="confirm"`)
	hiddenInputPattern = regexp.MustCompile(`<input type="hidden" name="([^"]+)" value="([^"]*)">`)
)

// TestDialogNativeSubmissionsReachProtectedRoutes exercises the browser's
// public form contract rather than checking either side's field names. The
// rendered form is submitted as a browser submits it, then the real CSRF and
// method-override middleware decide whether the registered route is reached.
func TestDialogNativeSubmissionsReachProtectedRoutes(t *testing.T) {
	const path = "/invoices/7"

	csrf := security.NewCSRF([]byte("0123456789abcdef0123456789abcdef"), time.Hour)
	token, err := csrf.Issue("")
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{http.MethodPut, http.MethodPatch, http.MethodDelete} {
		t.Run(want, func(t *testing.T) {
			reached := ""
			router := fhttp.NewRouter()
			handle := func(w http.ResponseWriter, r *http.Request) {
				reached = r.Method
				w.WriteHeader(http.StatusNoContent)
			}
			switch want {
			case http.MethodPut:
				router.Put(path, handle)
			case http.MethodPatch:
				router.Patch(path, handle)
			case http.MethodDelete:
				router.Delete(path, handle)
			}

			handler := frameworkmiddleware.CSRFProtect(csrf, func(*http.Request) string { return "" })(
				middleware.OverrideMethod()(router),
			)
			markup := components.Dialog(components.DialogProps{
				ID:     strings.ToLower(want) + "-invoice",
				Title:  "Change invoice?",
				Action: path,
				Method: strings.ToLower(want),
				Token:  token,
			})

			response := submitConfirmForm(t, string(markup), handler)
			if reached != want {
				t.Fatalf("the native form reached the router as %q with status %d, want %s: %s", reached, response.Code, want, response.Body.String())
			}
			if response.Code != http.StatusNoContent {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
			}
		})
	}
}

// submitConfirmForm applies native form submission semantics to the rendered
// confirm form. Browsers submit only GET and POST; an unknown method defaults
// to GET, which is why a PUT, PATCH or DELETE attribute cannot reach that route.
func submitConfirmForm(t *testing.T, markup string, handler http.Handler) *httptest.ResponseRecorder {
	t.Helper()

	form := confirmFormPattern.FindStringSubmatch(markup)
	if form == nil {
		t.Fatalf("dialog has no confirm form:\n%s", markup)
	}

	values := url.Values{}
	for _, field := range hiddenInputPattern.FindAllStringSubmatch(form[3], -1) {
		values.Add(html.UnescapeString(field[1]), html.UnescapeString(field[2]))
	}

	method := http.MethodGet
	target := html.UnescapeString(form[2])
	var body io.Reader
	switch strings.ToLower(html.UnescapeString(form[1])) {
	case "post":
		method = http.MethodPost
		body = strings.NewReader(values.Encode())
	case "get":
		query := values.Encode()
		if query != "" {
			target += "?" + query
		}
	}

	request := httptest.NewRequest(method, target, body)
	if method == http.MethodPost {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}
