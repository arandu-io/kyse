//go:build kyse

package components

@go
// AlertProps is a message about what just happened, or about what is about to.
//
// role is set from Variant rather than passed in: a destructive alert is an
// assertive live region and everything else is polite, and getting that pair
// wrong is the difference between a screen reader interrupting somebody and
// never mentioning the error at all.
type AlertProps struct {
	// Variant is "default" or "destructive". Empty is the default.
	Variant string
	// Title is the line in bold. It can stand alone.
	Title string
	// Message is the sentence under it, and may be empty.
	Message string
}

// Role is what assistive technology treats this as.
func (p AlertProps) Role() string {
	if p.Variant == "destructive" {
		return "alert"
	}
	return "status"
}
@endgo

<div
	class="alert"
	role="{{ .Role() }}"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
>
	<h2>{{ .Title }}</h2>
	@if(.Message != "")
		<section>{{ .Message }}</section>
	@endif
</div>
