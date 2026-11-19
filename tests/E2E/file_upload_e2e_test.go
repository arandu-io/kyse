//go:build e2e

package e2e

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/arandu-io/framework/test/browser"
	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse/components"
)

// uploadPicker is what a person does with the system's file dialog, which a
// page cannot open: the input is handed a DataTransfer and told it changed,
// which is what the browser does when the dialog closes. The images are real
// PNGs drawn on a canvas, because a thumbnail is made by decoding one.
const uploadPicker = `
window.__pick = async function (specs) {
  var dt = new DataTransfer();
  for (var i = 0; i < specs.length; i++) {
    var s = specs[i], blob;
    if (s.image) {
      var c = document.createElement('canvas');
      c.width = s.w || 64; c.height = s.h || 40;
      var g = c.getContext('2d'); g.fillStyle = '#c33'; g.fillRect(0, 0, c.width, c.height);
      blob = await new Promise(function (done) { c.toBlob(done, 'image/png'); });
    } else {
      blob = new Blob([new Uint8Array(s.size || 10)], { type: 'text/plain' });
    }
    dt.items.add(new File([blob], s.name, { type: s.image ? 'image/png' : 'text/plain' }));
  }
  var input = document.querySelector('input[type=file]');
  input.files = dt.files;
  input.dispatchEvent(new Event('change', { bubbles: true }));
  for (var n = 0; n < 100; n++) {
    var want = specs.filter(function (s) { return s.image; }).length;
    if (document.querySelectorAll('[data-part="list"] img:not([hidden])').length >= want) return true;
    await new Promise(function (done) { setTimeout(done, 50); });
  }
  return false;
};
window.__probe = function () {
  var root = document.querySelector('[data-kyse-behavior="file-upload"]');
  var summary = root.querySelector('[data-part="summary"]');
  var input = root.querySelector('input[type=file]');
  var colour = function (el) { return getComputedStyle(el).color; };
  var cls = function (name) {
    var probe = document.createElement('span');
    probe.className = name;
    document.body.appendChild(probe);
    var c = colour(probe);
    probe.remove();
    return c;
  };
  var rows = Array.prototype.map.call(root.querySelectorAll('[data-part="list"] > li'), function (li) {
    var img = li.querySelector('img'), name = li.querySelector('[data-part="file-name"]'), size = li.querySelector('[data-part="file-size"]');
    var box = img.getBoundingClientRect();
    return {
      name: name.textContent, size: size.textContent, over: li.getAttribute('data-over'),
      picture: img.hidden ? '' : img.getAttribute('src').slice(0, 22),
      pictureWidth: img.hidden ? 0 : Math.round(box.width), pictureHeight: img.hidden ? 0 : Math.round(box.height),
      nameColour: colour(name), sizeColour: colour(size), ellipsis: getComputedStyle(name).textOverflow,
      nameRight: Math.round(name.getBoundingClientRect().right), rowRight: Math.round(li.getBoundingClientRect().right),
    };
  });
  return JSON.stringify({
    summary: summary.textContent, over: summary.getAttribute('data-over'), summaryColour: colour(summary),
    summaryHidden: summary.getBoundingClientRect().width > 1,
    invalid: input.validationMessage, ariaInvalid: input.getAttribute('aria-invalid'),
    destructive: cls('text-destructive'), foreground: cls('text-foreground'), muted: cls('text-muted-foreground'),
    overflow: document.documentElement.scrollWidth > window.innerWidth,
    rows: rows,
  });
};`

type uploadRow struct {
	Name, Size, Over, Picture        string
	PictureWidth, PictureHeight      int
	NameColour, SizeColour, Ellipsis string
	NameRight, RowRight              int
}

type uploadProbe struct {
	Summary       string
	Over          string
	SummaryColour string
	SummaryHidden bool
	Invalid       string
	AriaInvalid   string
	Destructive   string
	Foreground    string
	Muted         string
	Overflow      bool
	Rows          []uploadRow
}

func uploadApp(t *testing.T) *app {
	t.Helper()
	css, err := stylesheet()
	if err != nil {
		t.Skipf("the stylesheet cannot be compiled here: %v", err)
	}
	a := &app{}
	mux := http.NewServeMux()
	mux.HandleFunc(view.AssetPath, view.Handler)
	mux.HandleFunc("GET /app.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
		w.Write(css)
	})
	mux.HandleFunc("GET /pick.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprint(w, uploadPicker)
	})
	field := func() string {
		return string(components.FileUpload(components.FileUploadProps{
			Name: "photos", Label: "Photos", Multiple: true, Accept: "image/*,.txt",
			Hint: "Up to 3 photos, 1 MB each.", MaxFiles: 3, MaxBytes: 1 << 20,
		}))
	}
	form := func(title string) string {
		return layout(title, `<form class="flex max-w-md flex-col gap-4">`+field()+`<button class="btn" type="submit">Send</button></form><script src="/pick.js"></script>`)
	}
	mux.HandleFunc("GET /files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, form("Photos"))
	})
	mux.HandleFunc("GET /start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		fmt.Fprint(w, layout("Start", `<a href="/files">Add photos</a><script src="/pick.js"></script>`))
	})
	a.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:")
		mux.ServeHTTP(w, r)
	}))
	t.Cleanup(a.Close)
	return a
}

func probe(t *testing.T, page *browser.Page) uploadProbe {
	t.Helper()
	var out uploadProbe
	if err := json.Unmarshal([]byte(eval(t, page, `window.__probe()`)), &out); err != nil {
		t.Fatalf("the field does not read: %v", err)
	}
	return out
}

const (
	photosAndANote = `[{"name":"front-of-the-house-with-a-very-long-name-that-has-to-be-cut-somewhere.png","image":true},{"name":"kitchen.png","image":true,"w":120,"h":80},{"name":"lease.txt","size":2048}]`
	tooMuch        = `[{"name":"a.png","image":true},{"name":"b.txt","size":10},{"name":"c.txt","size":10},{"name":"huge.txt","size":1572864}]`
)

func choose(t *testing.T, page *browser.Page, specs string) {
	t.Helper()
	page.Evaluate(`window.__pick(`+specs+`)`, nil)
}

// The component as an application draws it, with the stylesheet built from its
// sources and the behaviour the view runtime serves: a pick lists each file
// with a picture at its size, a name that is cut where the row ends and a size,
// all in the page's colours, under a policy that allows no inline script or
// style and no blob: picture.
func TestAPickedFileIsListedAsTheComponentDrawsIt(t *testing.T) {
	page := open(t, uploadApp(t), 1280)
	page.Visit("/files")
	choose(t, page, photosAndANote)
	got := probe(t, page)

	if got.Summary != "3 files chosen" || got.Over != "" || got.Invalid != "" {
		t.Errorf("summary %q over %q invalid %q", got.Summary, got.Over, got.Invalid)
	}
	if len(got.Rows) != 3 {
		t.Fatalf("%d rows: %+v", len(got.Rows), got.Rows)
	}
	for i, row := range got.Rows[:2] {
		if row.Picture != "data:image/jpeg;base64" || row.PictureWidth != 48 || row.PictureHeight != 48 {
			t.Errorf("row %d picture %q at %dx%d, want a data: URL at 48x48", i, row.Picture, row.PictureWidth, row.PictureHeight)
		}
	}
	if got.Rows[2].Picture != "" {
		t.Errorf("a text file has a picture: %+v", got.Rows[2])
	}
	if got.Rows[0].Ellipsis != "ellipsis" || got.Rows[0].NameRight > got.Rows[0].RowRight {
		t.Errorf("a long name is not cut where the row ends: %+v", got.Rows[0])
	}
	if got.Rows[0].NameColour != got.Foreground || got.Rows[0].SizeColour != got.Muted {
		t.Errorf("a row is not drawn in the page's colours: %+v, foreground %s, muted %s", got.Rows[0], got.Foreground, got.Muted)
	}
	if got.Overflow {
		t.Error("the page scrolls sideways")
	}
	shot(t, page, "file-upload-light")
}

// A pick over a limit is warned in the destructive colour, on the line and on
// the row that is over, and the form is held.
func TestALimitIsWarnedInTheDestructiveColour(t *testing.T) {
	page := open(t, uploadApp(t), 1280)
	page.Visit("/files")
	choose(t, page, tooMuch)
	got := probe(t, page)

	for _, want := range []string{"4 files chosen", "Send up to 3 files at a time", "One file is over 1 MB"} {
		if !strings.Contains(got.Summary, want) {
			t.Errorf("the summary %q does not say %q", got.Summary, want)
		}
	}
	if got.Over != "true" || got.SummaryColour != got.Destructive {
		t.Errorf("the summary is %q in %s, want over in %s", got.Over, got.SummaryColour, got.Destructive)
	}
	if got.Invalid == "" || got.AriaInvalid != "true" {
		t.Errorf("the input is not held: %q %q", got.Invalid, got.AriaInvalid)
	}
	last := got.Rows[3]
	if last.Over != "true" || last.NameColour != got.Destructive || last.SizeColour != got.Destructive {
		t.Errorf("the row over the size is not marked: %+v, destructive %s", last, got.Destructive)
	}
	if got.Rows[0].Over != "" || got.Rows[0].NameColour == got.Destructive {
		t.Errorf("a row under the size is marked: %+v", got.Rows[0])
	}
	shot(t, page, "file-upload-over-light")
}

// Both themes: the same pick in the dark theme is drawn in the dark theme's
// colours, and none of them is the light theme's.
func TestTheListIsDrawnInBothThemes(t *testing.T) {
	page := open(t, uploadApp(t), 1280)
	page.Visit("/files")
	choose(t, page, tooMuch)
	light := probe(t, page)

	page.Evaluate(`document.documentElement.classList.add('dark')`, nil)
	dark := probe(t, page)
	if dark.Foreground == light.Foreground || dark.Muted == light.Muted || dark.Destructive == light.Destructive {
		t.Errorf("the dark theme draws the light colours: light %s %s %s, dark %s %s %s",
			light.Foreground, light.Muted, light.Destructive, dark.Foreground, dark.Muted, dark.Destructive)
	}
	if dark.Rows[0].NameColour != dark.Foreground || dark.Rows[0].SizeColour != dark.Muted ||
		dark.Rows[3].NameColour != dark.Destructive || dark.SummaryColour != dark.Destructive {
		t.Errorf("a row is not drawn in the dark theme's colours: %+v", dark)
	}
	shot(t, page, "file-upload-over-dark")
}

// On a phone the rows fit and the page does not scroll sideways.
func TestTheListFitsAPhone(t *testing.T) {
	page := open(t, uploadApp(t), 390)
	page.Visit("/files")
	choose(t, page, photosAndANote)
	got := probe(t, page)
	if got.Overflow {
		t.Error("the page scrolls sideways at 390px")
	}
	if got.Rows[0].NameRight > got.Rows[0].RowRight {
		t.Errorf("a name runs out of its row: %+v", got.Rows[0])
	}
	shot(t, page, "file-upload-phone")
}

// The behaviour is mounted on a page a boosted link brought in.
func TestTheListWorksAfterABoostedNavigation(t *testing.T) {
	page := open(t, uploadApp(t), 1280)
	page.Visit("/start")
	page.ClickText("Add photos")
	page.ExpectURL("/files")
	choose(t, page, photosAndANote)
	if got := probe(t, page); got.Summary != "3 files chosen" || len(got.Rows) != 3 {
		t.Errorf("after a boosted navigation: %q, %d rows", got.Summary, len(got.Rows))
	}
}
