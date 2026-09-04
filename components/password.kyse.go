//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/hesape/validation"
)

@go
// PasswordBehavior is the name the client behaviour is registered under with
// arandu.ui.define, and the name this component writes when the caller named
// none.
//
// It is exported because the registration and the markup have to agree on one
// string, and a literal typed on both sides is a pair that drifts silently: a
// misspelling mounts nothing and the panel simply never opens.
const PasswordBehavior = "password"

// PasswordProps is a password box with the policy it will be judged by drawn
// under it: a strength summary, one line per requirement, and a control that
// reveals what was typed.
//
// # The policy comes from the validator, and there is only one of it
//
// Policy is the same *validation.Password the server checks the submitted value
// with. The lines under the box are drawn from its AppliedRules, so the list a
// person reads and the rule that rejects them are the same declaration. A
// component that took a minimum length of its own would be a second policy: it
// would agree on the first day, and on the day somebody raised the minimum in
// one place it would tell people a password is acceptable and then refuse it.
//
// A nil Policy takes validation.PasswordDefault, which is what the application
// registered with validation.PasswordDefaults, or a minimum of eight when it
// registered nothing. So a form that has said nothing gets a sensible rule, and
// a project that raised its minimum once gets the raised one everywhere without
// writing it on any field.
//
// # What the client does, and what it does not decide
//
// The panel is opened and the lines are ticked by a behaviour named
// PasswordBehavior, reached through the attributes ComponentProps writes:
// the name, and AppliedRules as its props. Nothing here is a script and nothing
// evaluates one -- the attribute carries a name that is looked up in a map,
// which is what lets the page keep a policy with no unsafe-eval.
//
// That behaviour is a convenience and never the check. It runs in a browser, on
// a value the browser owns, so a form posted around it arrives with whatever it
// arrives with. What decides is the server running the same *validation.Password
// this component was handed.
//
// # What the policy cannot say, and so this does not draw
//
// validation.Password asks for at least one letter, number or symbol -- a
// count of two is not something it can express, and neither is a set of
// permitted characters. Rules merged with Policy.Rules are the escape hatch for
// both, and a rule string is not a sentence: CustomRuleLabel is how the one
// line covering them is worded, because the component cannot turn a pattern
// into English and inventing typed props for what the validator cannot check
// would be a requirement nothing enforces.
//
// It publishes root, label, group, input, reveal, panel, meter, fill,
// strength, requirements, requirement, done, message and hint.
type PasswordProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the form field name, the id everything else is hung off, and what
	// Page is asked about.
	Name string
	// Label is the text above the box.
	Label string
	// Policy is the rule the submitted value is checked against, and what the
	// requirement lines are drawn from. Nil takes the application's registered
	// default.
	Policy *validation.Password
	// CustomRuleLabel is the sentence for the rules Policy.Rules merged in.
	// Empty draws a line saying they are checked on submit, because a policy
	// that asks for something the panel never mentions is a rejection nobody
	// was warned about.
	CustomRuleLabel string
	// Hint is the sentence under the box that is always there.
	Hint string
	// Page is the screen's own view.Page, which is what this box asks for its
	// message. Nil draws no message.
	Page Page
	// Autocomplete is the browser hint. Empty means "new-password", which is
	// the sign-up form this component is mostly drawn on; a sign-in box wants
	// "current-password".
	Autocomplete string
	// Placeholder is the grey text inside the empty box.
	Placeholder string
	// ShowLabel is the reveal control before it is pressed. Empty means "Show".
	ShowLabel string
	// HideLabel is what the behaviour writes into the reveal control once the
	// value is visible. Empty means "Hide". It is a prop of this component
	// rather than of the behaviour so that the pair is translated together.
	HideLabel string
	// DoneLabel is the control that closes the panel. Empty means "Done".
	DoneLabel string
	// Required marks the box required, in the markup and to the browser.
	Required bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
}

// PasswordRequirement is one line of the checklist: the key the policy declares
// it under, and the sentence drawn for it.
//
// The key is the one AppliedRules uses, and it is written onto the line, so the
// behaviour matches a line to the rule it was given by the same name rather
// than by its position or its wording.
type PasswordRequirement struct {
	// Key is the name this requirement has in the policy: min, max, letters,
	// mixedCase, numbers, symbols, uncompromised, custom.
	Key string
	// Text is the sentence the line reads.
	Text string
}

// AppliedPolicy is the policy in force for this field: the one it was given, or
// the application's registered default.
func (p PasswordProps) AppliedPolicy() *validation.Password {
	if p.Policy != nil {
		return p.Policy
	}
	return validation.PasswordDefault()
}

// Requirements are the lines under the box, in the order they are read.
//
// Length first, because it is the one every policy has and the one most
// passwords fail; then the kinds of character, in the order the validator
// checks them; then what only the server can answer.
func (p PasswordProps) Requirements() []PasswordRequirement {
	rules := p.AppliedPolicy().AppliedRules()

	number := func(key string) int {
		size, _ := rules[key].(int)
		return size
	}
	asked := func(key string) bool {
		on, _ := rules[key].(bool)
		return on
	}

	var out []PasswordRequirement
	if least := number("min"); least > 0 {
		out = append(out, PasswordRequirement{"min", "At least " + strconv.Itoa(least) + " characters"})
	}
	if most := number("max"); most > 0 {
		out = append(out, PasswordRequirement{"max", "At most " + strconv.Itoa(most) + " characters"})
	}
	if asked("letters") {
		out = append(out, PasswordRequirement{"letters", "At least one letter"})
	}
	if asked("mixedCase") {
		out = append(out, PasswordRequirement{"mixedCase", "At least one uppercase and one lowercase letter"})
	}
	if asked("numbers") {
		out = append(out, PasswordRequirement{"numbers", "At least one number"})
	}
	if asked("symbols") {
		out = append(out, PasswordRequirement{"symbols", "At least one symbol"})
	}
	if asked("uncompromised") {
		out = append(out, PasswordRequirement{"uncompromised", "Not found in a known data leak"})
	}
	if merged, _ := rules["customRules"].([]string); len(merged) > 0 {
		out = append(out, PasswordRequirement{"custom", p.customRuleText()})
	}
	return out
}

// customRuleText is the sentence for the rules the policy merged in.
//
// The rule strings themselves are not drawn. A pattern read out to somebody
// filling in a form is noise, and the component has no way to turn one into the
// sentence it means.
func (p PasswordProps) customRuleText() string {
	if p.CustomRuleLabel != "" {
		return p.CustomRuleLabel
	}
	return "Meets the remaining rules, checked when the form is sent"
}

// RequirementTotal is how many lines the checklist has, which is the scale the
// strength summary counts against.
func (p PasswordProps) RequirementTotal() int { return len(p.Requirements()) }

// StrengthText is what the summary says before anything has been typed.
//
// It is a sentence rather than a colour. A bar is the fastest way to read
// progress and the only way to read it is to see it, so the same count is
// written out here and this is the element pointed at as the live region.
func (p PasswordProps) StrengthText() string {
	return "0 of " + strconv.Itoa(p.RequirementTotal()) + " requirements met"
}

// InputAutocomplete is the browser hint, defaulting to a password being chosen
// rather than one being recalled.
func (p PasswordProps) InputAutocomplete() string {
	if p.Autocomplete != "" {
		return p.Autocomplete
	}
	return "new-password"
}

// RevealText is what the reveal control reads while the value is hidden.
func (p PasswordProps) RevealText() string {
	if p.ShowLabel != "" {
		return p.ShowLabel
	}
	return "Show"
}

// ConcealText is what the reveal control reads once the value is visible. The
// behaviour writes it in; it is declared here so both words come from the same
// call and are translated together.
func (p PasswordProps) ConcealText() string {
	if p.HideLabel != "" {
		return p.HideLabel
	}
	return "Hide"
}

// DoneText is the control that closes the panel.
func (p PasswordProps) DoneText() string {
	if p.DoneLabel != "" {
		return p.DoneLabel
	}
	return "Done"
}

// PanelID names the section that expands under the box.
func (p PasswordProps) PanelID() string { return p.Name + "-requirements-panel" }

// RequirementsID names the checklist, which is what the box is described by.
func (p PasswordProps) RequirementsID() string { return p.Name + "-requirements" }

// StrengthID names the summary, which is the live region.
func (p PasswordProps) StrengthID() string { return p.Name + "-strength" }

// Message is what validation left for this box, or empty.
func (p PasswordProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// DescribedBy names what explains this box, in the order it should be read: the
// complaint or the hint first, then the checklist.
//
// The summary is deliberately not in the list. It is a live region, so it is
// announced when it changes; naming it here as well would read the count out on
// every focus and again on every keystroke that moved it.
func (p PasswordProps) DescribedBy() string {
	switch {
	case p.Message() != "":
		return p.Name + "-error " + p.RequirementsID()
	case p.Hint != "":
		return p.Name + "-hint " + p.RequirementsID()
	}
	return p.RequirementsID()
}

// RootAttrs are the outermost element's attributes, with the client bridge
// filled in from the policy when the caller named no behaviour of their own.
//
// The props handed over are AppliedRules verbatim, so the behaviour ticking the
// lines is reading the same declaration the server rejects with. A caller who
// named a behaviour keeps it and owns its props: overriding theirs would make
// the field impossible to drive with anything but the shipped behaviour.
func (p PasswordProps) RootAttrs() map[string]string {
	if p.Behavior.Name == "" {
		p.Behavior = Behavior{Name: PasswordBehavior, Props: p.AppliedPolicy().AppliedRules()}
	}
	return p.ComponentProps.RootAttrs()
}

// PartNames are the parts this component publishes.
func (p PasswordProps) PartNames() []string {
	return []string{
		"root", "label", "group", "input", "reveal", "panel", "meter", "fill",
		"strength", "requirements", "requirement", "done", "message", "hint",
	}
}
@endgo

{{-- The box carries no value and never will. The flash keeps the message for a
     password and drops what was typed, so there is nothing to redraw, and a
     value attribute here would be a password written into the markup of a page
     that is cached, logged and scrolled past over somebody's shoulder.

     spellcheck, autocapitalize and autocorrect are written whatever the state
     is, and that is what makes revealing safe. Revealing is a swap to
     type="text", and a text box is prose to a browser: it is offered the
     values saved for other text boxes, its content is sent to a spell checker,
     and a phone capitalises the first letter of it. The alternatives were
     considered and are worse -- a second visible box means two controls holding
     one secret and a password manager filling the wrong one, and hiding the
     glyphs in CSS is not supported everywhere and leaves the value selectable
     and copyable in plain text anyway. So the swap stays, autocomplete keeps
     naming the field a password through it, and these three say the box is not
     prose whichever type it is wearing. --}}

<div
	data-part="root"
	class="{{ .RootClass("field") }}"
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		class="{{ .PartClass("label", "label") }}"
		for="{{ .Name }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	<div
		data-part="group"
		class="{{ .PartClass("group", "input-group") }}"
		@attributes(.PartAttrs("group"))
	>
		<input
			data-part="input"
			@if(.PartClass("input") != "")
				class="{{ .PartClass("input") }}"
			@endif
			type="password"
			id="{{ .Name }}"
			name="{{ .Name }}"
			autocomplete="{{ .InputAutocomplete() }}"
			spellcheck="false"
			autocapitalize="none"
			autocorrect="off"
			aria-describedby="{{ .DescribedBy() }}"
			@attributes(.PartAttrs("input"))
			@if(.Placeholder != "")
				placeholder="{{ .Placeholder }}"
			@endif
			@if(.Message() != "")
				aria-invalid="true"
			@endif
			@if(.Required)
				required
			@endif
			@if(.Autofocus)
				autofocus
			@endif
		>

		{{-- The word flips and aria-pressed carries the state, so the control
		     reads as a toggle rather than as two buttons that look alike. It is
		     a word and not a picture on purpose: an icon-only control needs a
		     name that says which of the two states it is in, and the name and
		     the picture then have to be kept in step by whoever swaps them. --}}
		<button
			data-part="reveal"
			class="{{ .PartClass("reveal", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			data-align="end"
			data-reveal-hidden="{{ .RevealText() }}"
			data-reveal-shown="{{ .ConcealText() }}"
			aria-pressed="false"
			aria-controls="{{ .Name }}"
			@attributes(.PartAttrs("reveal"))
		>{{ .RevealText() }}</button>
	</div>

	{{-- Closed, and closed by the hidden attribute rather than by a class: what
	     the panel is worth is decided by whether it is in the accessibility
	     tree, and a class that only stops it being painted leaves it read out
	     under a box nobody has typed in yet. The behaviour drops the attribute
	     on the first keystroke, which is also when aria-describedby starts
	     resolving to the checklist. --}}
	<section
		data-part="panel"
		class="{{ .PartClass("panel", "mt-2 flex flex-col gap-2") }}"
		id="{{ .PanelID() }}"
		hidden
		@attributes(.PartAttrs("panel"))
	>
		{{-- The bar restates in colour what the sentence under it says in
		     words, so it is marked decorative: announced as well, it would read
		     the same count twice on every keystroke. The sentence is the live
		     region, and it is what a reader who cannot see the fill is left
		     with. The fill starts empty because the box does. --}}
		<div
			data-part="meter"
			class="{{ .PartClass("meter", "progress") }}"
			aria-hidden="true"
			@attributes(.PartAttrs("meter"))
		>
			<span
				data-part="fill"
				class="{{ .PartClass("fill", "w-0") }}"
				@attributes(.PartAttrs("fill"))
			></span>
		</div>

		<p
			data-part="strength"
			class="{{ .PartClass("strength", "text-muted-foreground text-sm") }}"
			id="{{ .StrengthID() }}"
			aria-live="polite"
			aria-atomic="true"
			@attributes(.PartAttrs("strength"))
		>{{ .StrengthText() }}</p>

		{{-- Every line starts unmet, which is true rather than convenient: the
		     box is empty. data-met is the one copy of the state -- the
		     stylesheet paints from it and the behaviour writes it -- and the
		     word beside it is that state spoken, because a tick that is only a
		     colour or only a shape is a line whose state cannot be read out. --}}
		<ul
			data-part="requirements"
			class="{{ .PartClass("requirements", "flex flex-col gap-1 text-sm") }}"
			id="{{ .RequirementsID() }}"
			@attributes(.PartAttrs("requirements"))
		>
			@foreach(.Requirements() as requirement)
				<li
					data-part="requirement"
					class="{{ .PartClass("requirement", "flex items-center gap-2") }}"
					data-requirement="{{ requirement.Key }}"
					data-met="false"
					@attributes(.PartAttrs("requirement"))
				><span class="sr-only">Not met:</span>{{ requirement.Text }}</li>
			@endforeach
		</ul>

		<button
			data-part="done"
			class="{{ .PartClass("done", "btn") }}"
			type="button"
			data-variant="ghost"
			data-size="sm"
			aria-controls="{{ .PanelID() }}"
			@attributes(.PartAttrs("done"))
		>{{ .DoneText() }}</button>
	</section>

	@if(.Message() != "")
		<p
			data-part="message"
			id="{{ .Name }}-error"
			class="{{ .PartClass("message", "text-destructive text-sm") }}"
			@attributes(.PartAttrs("message"))
		>{{ .Message() }}</p>
	@endif
	@if(.Message() == "")
		@if(.Hint != "")
			<p
				data-part="hint"
				id="{{ .Name }}-hint"
				class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("hint"))
			>{{ .Hint }}</p>
		@endif
	@endif
</div>
