//go:build kyse

package components

@go
// MaskedBehavior is the name the client behaviour is registered under with
// arandu.ui.define, and the name this component writes.
//
// Exported because the registration and the markup have to agree on one string,
// and a literal typed on both sides is a pair that drifts silently.
const MaskedBehavior = "mask"

// MaskedProps is a text box that formats what somebody types, and submits what
// a program stores.
//
// It draws two inputs. The visible one shows 123.456.789-00 and carries no form
// name; the hidden one beside it carries Name and holds 12345678900, kept in
// step on every keystroke. So the value that reaches the server is raw, and
// nothing there has to know a mask existed.
//
// # Why two inputs and not one
//
// The alternative is submitting the formatted value and unmasking it on the
// server. That works until one of the places that read the field forgets, and
// what a forgotten one writes is punctuation in a column somebody will later
// compare against a value without it. Two inputs make the raw value the only
// thing with a name, so forgetting is not available.
//
// The cost is stated rather than hidden: with scripting off, the hidden input
// stays empty and the field arrives blank. A form that must work without
// scripting takes Input and unmasks on the server; this component is for the
// ordinary case, and the ordinary case has scripting.
//
// # The pattern
//
// Pattern is the mask, in the tokens mask.Pattern documents: 0 a digit, 9 an
// optional digit, # a digit repeated to the end, A alphanumeric, S a letter.
// Everything else is written through.
//
//	components.Masked(components.MaskedProps{
//		Name:    "postcode",
//		Pattern: "00000-000",
//		Page:    page,
//	})
//
// Patterns is the same thing for a field that takes either of two shapes: the
// shortest one that still holds what has been typed is the one applied, so the
// field grows into the longer form at the character that needs it.
//
//	// A Brazilian document field, declared by the application that has one.
//	Patterns: []string{"000.000.000-00", "00.000.000/0000-00"},
//
// No pattern here belongs to a country. A national document is a rule before it
// is a shape -- eleven digits are not a CPF until the check digits agree -- and
// a framework that named the shape would be implying it checks the rule. The
// application names its own, which costs it one line.
//
// # It formats and does not decide
//
// A mask says what may be typed. Whether what was typed means anything is the
// server's answer: the value arrives raw, the validator reads it, and this
// component is not consulted. data-mask-complete on the visible input says
// whether the pattern is filled, for a stylesheet to show -- and a filled
// pattern is a shape, never a verdict.
type MaskedProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name. It goes on the hidden input, because the
	// hidden input is what submits.
	Name string
	// ID is the id the visible input carries, so a Label can point at it.
	// Empty uses Name.
	ID string
	// Pattern is the mask. Use Patterns instead for a field with two shapes.
	Pattern string
	// Patterns are the shapes this field takes, shortest first in effect: the
	// one applied is the shortest whose capacity still holds what was typed.
	Patterns []string
	// Reverse fills the pattern from the right, which is how money is typed:
	// the first digit entered is the last one of the value.
	Reverse bool
	// ClearIfNotMatch empties the field on the way out when the pattern was
	// never filled. A half-typed value left behind is one somebody submits
	// without looking.
	ClearIfNotMatch bool
	// Value is the raw value this starts with. It is formatted on mount, so a
	// stored 12345678900 is shown as the pattern writes it.
	Value string
	// Placeholder is the grey text inside an empty box.
	Placeholder string
	// Page is the screen's own view.Page, asked for what was typed and for
	// whether it was rejected.
	Page Page
	// AriaLabel is the accessible name for a box no label points at.
	AriaLabel string
	// DescribedBy is the id of the element that explains this box.
	DescribedBy string
	// Autocomplete is the browser hint.
	Autocomplete string
	// Required marks the box required.
	Required bool
	// Disabled takes it out of the form and out of the tab order.
	Disabled bool
	// Readonly shows the value and refuses edits.
	Readonly bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
}

// ElementID is the id the visible input carries.
func (p MaskedProps) ElementID() string {
	if p.ID != "" {
		return p.ID
	}
	return p.Name
}

// DisplayName is the name of the visible input.
//
// It has one, and it is not the field's: a box with no name at all is a box the
// browser will not remember or autofill, and one carrying the field's name
// would submit the formatted value beside the raw one.
func (p MaskedProps) DisplayName() string { return p.Name + "_display" }

// Applied are the patterns this field uses, whichever way they were given.
func (p MaskedProps) Applied() []string {
	if len(p.Patterns) > 0 {
		return p.Patterns
	}
	if p.Pattern != "" {
		return []string{p.Pattern}
	}
	return nil
}

// Current is the raw value the hidden input starts with: what was typed on a
// rejected attempt, and Value when there was none.
func (p MaskedProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// Formatted is Current as the pattern writes it, which is what the visible box
// starts with.
//
// Formatted here and not left to the behaviour, so the value is right in the
// first frame: a box that shows the raw value until a script runs is a box that
// flickers on every page load.
func (p MaskedProps) Formatted() string {
	patterns := p.Applied()
	if len(patterns) == 0 {
		return p.Current()
	}
	return maskApply(maskFor(patterns, p.Current()), p.Current())
}

// Message is the rejection for this field, or empty.
func (p MaskedProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in from the pattern when the caller named no behaviour of their own.
//
// The props handed over are the patterns themselves, so the behaviour formats
// with the same declaration this component rendered with. A caller who named a
// behaviour keeps it and owns its props.
func (p MaskedProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		props := map[string]any{"pattern": p.Applied()}
		if p.Reverse {
			props["reverse"] = true
		}
		if p.ClearIfNotMatch {
			props["clearIfNotMatch"] = true
		}
		p.Behavior = Behavior{Name: MaskedBehavior, Props: props}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
//
// raw is the hidden input, and it is published because a caller writing a form
// that reads its own fields needs to reach it.
func (p MaskedProps) PartNames() []string { return []string{"root", "input", "raw"} }

// maskAccepts reports what a pattern character takes, and whether it is a token.
//
// The dictionary is the one the client behaviour reads, spelled once on each
// side because the two cannot import from each other: this renders the first
// frame and that formats every keystroke after it, and a value formatted
// differently by the two would flicker on load. A test compares the two tables.
//
//	0  a digit, required
//	9  a digit, optional
//	#  a digit, repeated to the end
//	A  a letter or a digit
//	S  a letter
func maskAccepts(symbol, value rune) (ok, isToken, repeating bool) {
	digit := value >= '0' && value <= '9'
	letter := (value >= 'a' && value <= 'z') || (value >= 'A' && value <= 'Z')
	switch symbol {
	case '0', '9':
		return digit, true, false
	case '#':
		return digit, true, true
	case 'A':
		return digit || letter, true, false
	case 'S':
		return letter, true, false
	}
	return false, false, false
}

// maskCapacity is how many characters a pattern's tokens accept, and -1 for one
// that repeats and therefore has no end.
func maskCapacity(pattern string) int {
	var n int
	for _, symbol := range pattern {
		ok, isToken, repeating := maskAccepts(symbol, '0')
		_ = ok
		if !isToken {
			continue
		}
		if repeating {
			return -1
		}
		n++
	}
	return n
}

// maskFor picks the shortest pattern that still holds what is there.
func maskFor(patterns []string, value string) string {
	if len(patterns) == 0 {
		return ""
	}
	if len(patterns) == 1 {
		return patterns[0]
	}
	widest := patterns[0]
	for _, pattern := range patterns {
		if maskCapacity(pattern) > maskCapacity(widest) {
			widest = pattern
		}
	}
	typed := len([]rune(maskUnmask(widest, value)))

	best, found := "", false
	for _, pattern := range patterns {
		size := maskCapacity(pattern)
		if size < 0 || size < typed {
			continue
		}
		if !found || size < maskCapacity(best) {
			best, found = pattern, true
		}
	}
	if !found {
		return widest
	}
	return best
}

// maskApply formats a value, dropping what the pattern cannot accept.
func maskApply(pattern, value string) string {
	if pattern == "" {
		return value
	}
	var out []rune
	symbols := []rune(pattern)
	input := []rune(value)

	var at int
	for cursor := 0; cursor < len(symbols) && at < len(input); {
		symbol := symbols[cursor]
		takes, isToken, repeating := maskAccepts(symbol, input[at])
		if !isToken {
			out = append(out, symbol)
			if input[at] == symbol {
				at++
			}
			cursor++
			continue
		}
		if !takes {
			at++
			continue
		}
		out = append(out, input[at])
		at++
		if !repeating {
			cursor++
		}
	}
	for len(out) > 0 {
		last := out[len(out)-1]
		if _, isToken, _ := maskAccepts('A', last); isToken && (last >= '0' && last <= '9' ||
			last >= 'a' && last <= 'z' || last >= 'A' && last <= 'Z') {
			break
		}
		out = out[:len(out)-1]
	}
	return string(out)
}

// maskUnmask gives back what the hidden input carries.
func maskUnmask(pattern, value string) string {
	if pattern == "" {
		return value
	}
	var out []rune
	symbols := []rune(pattern)
	input := []rune(value)

	var at int
	for cursor := 0; cursor < len(symbols) && at < len(input); {
		symbol := symbols[cursor]
		takes, isToken, repeating := maskAccepts(symbol, input[at])
		if !isToken {
			if input[at] == symbol {
				at++
			}
			cursor++
			continue
		}
		if !takes {
			at++
			continue
		}
		out = append(out, input[at])
		at++
		if !repeating {
			cursor++
		}
	}
	return string(out)
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("masked") }}"
	@attributes(.RootAttrs())
>
	<input
		data-part="input"
		class="{{ .PartClass("input", "input") }}"
		type="text"
		id="{{ .ElementID() }}"
		name="{{ .DisplayName() }}"
		value="{{ .Formatted() }}"
		autocomplete="off"
		@if(.Placeholder != "")
			placeholder="{{ .Placeholder }}"
		@endif
		@if(.AriaLabel != "")
			aria-label="{{ .AriaLabel }}"
		@endif
		@if(.DescribedBy != "")
			aria-describedby="{{ .DescribedBy }}"
		@endif
		@if(.Message() != "")
			aria-invalid="true"
		@endif
		@if(.Required)
			required
		@endif
		@if(.Disabled)
			disabled
		@endif
		@if(.Readonly)
			readonly
		@endif
		@if(.Autofocus)
			autofocus
		@endif
	>
	<input
		data-part="raw"
		class="{{ .PartClass("raw", "") }}"
		type="hidden"
		name="{{ .Name }}"
		value="{{ .Current() }}"
		@if(.Autocomplete != "")
			autocomplete="{{ .Autocomplete }}"
		@endif
	>
</div>
