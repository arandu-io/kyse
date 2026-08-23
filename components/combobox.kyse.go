//go:build kyse

package components

import (
	"strconv"

	"github.com/arandu-io/kyse/icons"
)

@go
// ComboboxProps is a text box that chooses one value out of a list too long to
// put on the page.
//
// # Where the filtering happens
//
// On the server, on every keystroke, and nowhere else: a list short enough to
// send to the browser whole is a list a select already draws, so the component
// that exists for the other case has to keep working at the size that made it
// necessary.
//
// SearchURL is what carries that, and it is not optional. A combobox with no
// URL has no list it could ever show, so it draws nothing rather than a box
// that stays empty however long somebody types into it.
//
// The endpoint answers with this same component, rendered with Options set to
// what it matched. The request asks for the listbox alone and swaps it in
// place, so the endpoint writes no markup of its own and the ARIA below is
// written once.
//
// # What is submitted
//
// Two names. Name carries the chosen value, in a hidden input, and it is what
// the form is about. Query carries the text somebody typed, because the
// browser sends the box it fired the request from and a box with no name sends
// nothing. A handler that only wants the first ignores the second.
type ComboboxProps struct {
	// Name is the form field name, the id everything else is hung off, and
	// what Page is asked about.
	Name string
	// Label is the text above the box.
	Label string
	// SearchURL is where the list comes from. Empty draws nothing at all.
	SearchURL string
	// QueryParam is the parameter the typed text is sent under. Empty means
	// Name with "-query" after it, which is a name two comboboxes on one page
	// cannot share.
	QueryParam string
	// Options are what the server matched for the current query, in the order
	// it wants them read. Empty on the first paint, which is correct: nobody
	// has typed anything yet.
	Options []ComboboxOption
	// Value is the chosen value when nothing was rejected: the stored record
	// on an edit form, and empty on a create form.
	Value string
	// ValueLabel is what Value reads as in the box. Empty falls back to Value
	// itself, which is right when the value is already a word and wrong when
	// it is an identifier.
	ValueLabel string
	// Placeholder is the grey text inside an empty box.
	Placeholder string
	// Hint is the sentence under the box that is always there.
	Hint string
	// EmptyText is the sentence drawn in place of the list when the server
	// matched nothing. Empty takes the stylesheet's own wording.
	EmptyText string
	// Page is the screen's own view.Page, which is what this box asks for its
	// message and for what was chosen. Nil draws no message and keeps Value.
	Page Page
	// Delay is how long typing pauses before the server is asked, in
	// milliseconds. Zero means 200, which is long enough that a typed word
	// costs one request and short enough that the list feels attached to the
	// keyboard.
	Delay int
	// Required marks the box required, in the markup and to the browser.
	Required bool
	// Disabled draws the box as unavailable and stops it answering.
	Disabled bool
	// Autofocus puts the cursor here on load. At most one per screen.
	Autofocus bool
}

// ComboboxOption is one line of the list.
type ComboboxOption struct {
	// Value is what is submitted when this line is chosen.
	Value string
	// Label is what the line reads as. Empty falls back to Value.
	Label string
	// Disabled draws the line as a suggestion that cannot be chosen. It stays
	// visible, because a line that vanished would read as a list that is
	// missing it.
	Disabled bool
}

// Text is what the option reads as, falling back to the value it carries.
func (o ComboboxOption) Text() string {
	if o.Label == "" {
		return o.Value
	}
	return o.Label
}

// OptionID names one line of the list.
//
// It is numbered rather than named after the value it carries, because
// aria-activedescendant holds an id and an id with a space in it is read as
// two. A value is somebody's data and can hold anything; a position cannot.
//
// The id is what the Enter key is resolved through, so it has to be unique in
// the document: two boxes sharing a Name would number their lines the same and
// the earlier one would answer for both. Name is already required to be unique
// per page for what is submitted, and this is the second reason.
func (p ComboboxProps) OptionID(at int) string {
	return p.Name + "-option-" + strconv.Itoa(at)
}

// ListboxID names the list, which is what the box says it controls and where
// the server's answer is swapped in.
func (p ComboboxProps) ListboxID() string {
	return p.Name + "-listbox"
}

// ListboxTarget is the same list as a selector, for the request.
func (p ComboboxProps) ListboxTarget() string {
	return "#" + p.ListboxID()
}

// Query is the parameter the typed text is sent under.
func (p ComboboxProps) Query() string {
	if p.QueryParam != "" {
		return p.QueryParam
	}
	return p.Name + "-query"
}

// Trigger is when the server is asked: after typing stops, and only when what
// was typed changed.
func (p ComboboxProps) Trigger() string {
	delay := p.Delay
	if delay == 0 {
		delay = 200
	}
	return "input changed delay:" + strconv.Itoa(delay) + "ms"
}

// Message is what validation left for this box, or empty.
func (p ComboboxProps) Message() string {
	if p.Page == nil {
		return ""
	}
	return p.Page.FieldError(p.Name)
}

// Current is the value the box is drawn with: what was chosen on the rejected
// attempt, and Value when there was none.
func (p ComboboxProps) Current() string {
	if p.Page == nil {
		return p.Value
	}
	return p.Page.OldOr(p.Name, p.Value)
}

// CurrentLabel is what Current reads as in the box.
//
// A value that came back from a rejected attempt is a bare value with no label
// beside it, so it is shown as itself. ValueLabel is only right for the value
// it was given for, and using it for another one would put the wrong word in
// the box under a message complaining about the value.
func (p ComboboxProps) CurrentLabel() string {
	if p.Current() != p.Value {
		return p.Current()
	}
	if p.ValueLabel != "" {
		return p.ValueLabel
	}
	return p.Value
}

// DescribedBy names the element that explains this box, so the error and the
// hint are announced rather than only shown.
func (p ComboboxProps) DescribedBy() string {
	if p.Message() != "" {
		return p.Name + "-error"
	}
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}
@endgo

@if(.SearchURL != "")
	<div class="field">
		<label class="label" for="{{ .Name }}">{{ .Label }}</label>

		{{-- data-combobox is the scope the client script resolves everything from,
		     and it is the whole of what the client is told. Open is
		     aria-expanded on the box and aria-hidden on the popover, the active
		     line is aria-activedescendant, and the chosen one is aria-selected:
		     the attributes the markup has to carry anyway are the state, so
		     there is no second copy of it to fall out of step. --}}
		<div class="combobox" data-combobox>
			<input
				type="text"
				role="combobox"
				id="{{ .Name }}"
				name="{{ .Query() }}"
				value="{{ .CurrentLabel() }}"
				autocomplete="off"
				autocorrect="off"
				spellcheck="false"
				aria-autocomplete="list"
				aria-expanded="false"
				aria-controls="{{ .ListboxID() }}"
				hx-get="{{ .SearchURL }}"
				hx-trigger="{{ .Trigger() }}"
				hx-target="{{ .ListboxTarget() }}"
				hx-select="{{ .ListboxTarget() }}"
				hx-swap="outerHTML"
				hx-sync="this:replace"
				@if(.Placeholder != "")
					placeholder="{{ .Placeholder }}"
				@endif
				@if(.DescribedBy() != "")
					aria-describedby="{{ .DescribedBy() }}"
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
				@if(.Autofocus)
					autofocus
				@endif
			>

			{!! icons.CaretDown(icons.Props{}) !!}

			<div data-popover aria-hidden="true">
				<div
					role="listbox"
					id="{{ .ListboxID() }}"
					aria-orientation="vertical"
					@if(.EmptyText != "")
						data-empty="{{ .EmptyText }}"
					@endif
				>
					@for(at := 0; at < len(.Options); at++)
						<div
							role="option"
							id="{{ .OptionID(at) }}"
							data-value="{{ .Options[at].Value }}"
							data-label="{{ .Options[at].Text() }}"
							@if(.Options[at].Disabled)
								aria-disabled="true"
							@endif
							@if(.Options[at].Value == .Current())
								aria-selected="true"
							@endif
						>{{ .Options[at].Text() }}</div>
					@endfor
				</div>
			</div>

			<input type="hidden" name="{{ .Name }}" value="{{ .Current() }}" data-combobox-value>
		</div>

		@if(.Message() != "")
			<p id="{{ .Name }}-error" class="text-destructive text-sm">{{ .Message() }}</p>
		@endif
		@if(.Message() == "")
			@if(.Hint != "")
				<p id="{{ .Name }}-hint" class="text-muted-foreground text-sm">{{ .Hint }}</p>
			@endif
		@endif
	</div>
@endif
