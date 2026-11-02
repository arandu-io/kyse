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
	// "link". Empty is the default, and so is "default": a status enum that
	// says "default" would otherwise draw a badge no rule styles, which is the
	// bare text of the label with none of the pill around it.
	Variant string
}

// VariantName is the variant as written to data-variant: empty for the default.
func (p BadgeProps) VariantName() string {
	if p.Variant == "default" {
		return ""
	}
	return p.Variant
}

// PartNames are the parts this component publishes.
func (p BadgeProps) PartNames() []string { return []string{"root"} }
@endgo

<span
	data-part="root"
	class="{{ .RootClass("badge", "max-w-full") }}"
	@if(.VariantName() != "")
		data-variant="{{ .VariantName() }}"
	@endif
	@attributes(.RootAttrs())
>{{ .Label }}</span>
