//go:build kyse

package components

@go
// HiddenInputProps is a value a form carries that nobody types: the record an
// edit form is about, the subscription a dialog cancels, the page to go back
// to.
//
// It is a component rather than a line of markup in the view for the reason
// every input here is: the value is escaped by construction, and the value of
// a refused attempt comes back from Page by Name -- so a dialog that was
// refused for one subscription is drawn again carrying that subscription, not
// whichever one the page happens to draw first.
//
// It has no appearance, so it is the one component that writes no class of
// its own; Class still reaches it.
//
// It publishes root.
type HiddenInputProps struct {
	// ComponentProps is the class and attributes the caller adds.
	ComponentProps
	// Name is the form field name, and what Page is asked about.
	Name string
	// ID is the id it carries. Empty writes none: nothing points at a hidden
	// input, and an id derived from Name would repeat on a page with two forms
	// carrying the same field.
	ID string
	// Value is what it submits when nothing was refused.
	Value string
	// Page is where the value of a refused attempt is looked up. Nil keeps
	// Value.
	Page Page
}

// Current is the value it submits: what the refused attempt carried, or Value.
func (p HiddenInputProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// PartNames are the parts this component publishes.
func (p HiddenInputProps) PartNames() []string { return []string{"root"} }
@endgo

<input
	data-part="root"
	type="hidden"
	name="{{ .Name }}"
	value="{{ .Current() }}"
	@if(.ID != "")
		id="{{ .ID }}"
	@endif
	@if(.RootClass() != "")
		class="{{ .RootClass() }}"
	@endif
	@attributes(.RootAttrs())
>
