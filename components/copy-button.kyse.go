//go:build kyse

package components

@go
// CopyButtonBehavior is the name the client behaviour is registered under with
// arandu.ui.define.
const CopyButtonBehavior = "copy"

// CopyButtonProps is the button beside a value that puts it on the clipboard.
//
// It is for the values people otherwise select by dragging: an API key, a
// command, a hash, an address. Dragging over a value that wraps takes the
// newline with it, and over one that scrolls takes half of it.
//
// It draws as a button and not as an icon on its own, because the confirmation
// has to be announced and not only shown. A tick that appears is a tick that
// somebody using a screen reader is never told about, so the button's own
// label changes and the change is announced.
//
// Without script it is a button that does nothing, so the value it copies is
// always drawn on the page as well -- this is beside a value, never instead of
// one.
//
// It publishes root and feedback.
type CopyButtonProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Value is what goes on the clipboard. It is the value itself, not an id,
	// for the case where what is copied is not what is shown -- a key shown
	// truncated, a command shown wrapped.
	Value string
	// SourceID is the id of the element to copy the text of, for when what is
	// copied is exactly what is drawn and repeating it here would be two
	// copies of one string. Value wins when both are set.
	SourceID string
	// Label is the text on the button. Empty says "Copy".
	Label string
	// CopiedLabel is what it says after. Empty says "Copied".
	CopiedLabel string
	// Variant and Size style it. See ButtonProps.
	Variant string
	Size    string
	// IconOnly draws the label as the accessible name rather than as text,
	// for a button that sits inside a code block and has no room. The name is
	// still there; only the drawing changes.
	IconOnly bool
}

// Text is the label, or the word used when none was given.
func (p CopyButtonProps) Text() string {
	if p.Label != "" {
		return p.Label
	}
	return "Copy"
}

// Confirmation is what it says after the copy.
func (p CopyButtonProps) Confirmation() string {
	if p.CopiedLabel != "" {
		return p.CopiedLabel
	}
	return "Copied"
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in when the caller named no behaviour of their own.
//
// The value and the confirmation are handed over as props rather than read off
// the element, so the behaviour does not have to know which attribute the
// markup happened to use.
func (p CopyButtonProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: CopyButtonBehavior, Props: map[string]any{
			"value":  p.Value,
			"source": p.SourceID,
			"copied": p.Confirmation(),
			"label":  p.Text(),
		}}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p CopyButtonProps) PartNames() []string { return []string{"root", "feedback"} }
@endgo

<button
	data-part="root"
	class="{{ .RootClass("btn") }}"
	type="button"
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
	@if(.IconOnly)
		aria-label="{{ .Text() }}"
	@endif
	@attributes(.RootAttrs())
>
	@if(!.IconOnly)
		{{ .Text() }}
	@endif
	<span
		data-part="feedback"
		class="{{ .PartClass("feedback", "sr-only") }}"
		role="status"
		aria-live="polite"
		@attributes(.PartAttrs("feedback"))
	></span>
</button>
