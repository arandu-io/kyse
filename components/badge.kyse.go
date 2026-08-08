//go:build kyse

package components

@go
// BadgeProps is a small piece of status beside something else: published or
// draft, open or closed, the count on a tab.
type BadgeProps struct {
	// Label is the text.
	Label string
	// Variant is "primary", "secondary", "outline", "ghost", "destructive" or
	// "link". Empty is the default.
	Variant string
}
@endgo

<span
	class="badge"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
>{{ .Label }}</span>
