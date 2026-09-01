//go:build kyse

package components

@go
// AlertProps is a message about what just happened, or about what is about to.
//
// role is set from Variant rather than passed in: a destructive alert is an
// assertive live region and everything else is polite, and getting that pair
// wrong is the difference between a screen reader interrupting somebody and
// never mentioning the error at all.
// It publishes root, title and message.
type AlertProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
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

// PartNames are the parts this component publishes.
func (p AlertProps) PartNames() []string { return []string{"root", "title", "message"} }
@endgo

<div
	data-part="root"
	class="{{ .RootClass("alert") }}"
	role="{{ .Role() }}"
	@attributes(.RootAttrs())
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
>
	<h2
		data-part="title"
		@if(.PartClass("title") != "")
			class="{{ .PartClass("title") }}"
		@endif
		@attributes(.PartAttrs("title"))
	>{{ .Title }}</h2>
	@if(.Message != "")
		<section
			data-part="message"
			@if(.PartClass("message") != "")
				class="{{ .PartClass("message") }}"
			@endif
			@attributes(.PartAttrs("message"))
		>{{ .Message }}</section>
	@endif
</div>
