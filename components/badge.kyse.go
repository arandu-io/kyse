//go:build kyse

package components

@go
// BadgeProps is a small piece of status beside something else: published or
// draft, open or closed, the count on a tab.
//
// It publishes one part, root, which is the span itself.
type BadgeProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label is the text.
	Label string
	// Variant is "primary", "secondary", "outline", "ghost", "destructive" or
	// "link". Empty is the default.
	Variant string
}

// PartNames are the parts this component publishes.
func (p BadgeProps) PartNames() []string { return []string{"root"} }
@endgo

<span
	data-part="root"
	class="{{ .RootClass("badge") }}"
	@attributes(.RootAttrs())
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
>{{ .Label }}</span>
