//go:build kyse

package components

@go
// RelativeTimeBehavior is the name the client behaviour is registered under
// with arandu.ui.define.
const RelativeTimeBehavior = "relative-time"

// RelativeTimeProps is a moment, written so it is readable by a person and by
// a machine, and so the person reads it in their own time zone.
//
// The element is <time>, and the two halves are deliberate. The datetime
// attribute is the instant, in UTC, unambiguous, and it is what a browser, a
// crawler and a calendar read. The text is what a person sees, and it is
// rendered here, on the server, so a page with no script working still says
// when.
//
// What the script changes is the text and nothing else: it rewrites it in the
// visitor's locale and offset. A server has no way to know either -- the
// header that would say so does not exist -- so the choice is a wrong time
// zone for everyone or a right one for everyone once the page is live, and the
// server's answer stays as the floor under it.
//
// It publishes root.
type RelativeTimeProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// DateTime is the instant, in RFC 3339, in UTC: "2026-09-07T14:30:00Z".
	// It is a string and not a time.Time because what a view has is what a
	// query returned, already formatted, and because an unset time.Time is the
	// year one rather than nothing.
	DateTime string
	// Label is what the element reads before any script has run, and after one
	// has if it cannot: "7 September 2026". A page that renders this empty is
	// a page whose timestamps are blank without JavaScript.
	Label string
	// Style is how the script rewrites it: "relative" for "3 hours ago",
	// "date" for the day, "datetime" for the day and the time, "time" for the
	// time alone. Empty leaves the server's label as it is, which is the right
	// answer for a timestamp that is already exactly what it should say.
	Style string
	// Title puts the full instant in the tooltip, which is how a relative
	// label stays answerable: "3 hours ago" is useless in a support ticket
	// read tomorrow.
	Title string
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when a Style was asked for and the caller named no behaviour of
// their own.
//
// A caller who named one keeps it and owns its props: overriding theirs would
// make the element impossible to drive with anything but the shipped
// behaviour.
func (p RelativeTimeProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" && p.Style != "" {
		p.Behavior = Behavior{Name: RelativeTimeBehavior, Props: map[string]any{"style": p.Style}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p RelativeTimeProps) PartNames() []string { return []string{"root"} }
@endgo

@if(.DateTime != "")
	<time
		data-part="root"
		class="{{ .RootClass("relative-time") }}"
		datetime="{{ .DateTime }}"
		@if(.Title != "")
			title="{{ .Title }}"
		@endif
		@attributes(.RootAttrs())
	>{{ .Label }}</time>
@endif
