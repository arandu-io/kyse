//go:build kyse

package components

import "strconv"

@go
// ActiveSearchProps is a search box that replaces the results as somebody
// types.
//
// The three things that go wrong with one are all decided here rather than
// left to whoever writes the page. It waits before asking, so a word typed at
// speed is one request and not eight. It cancels the request still in flight
// when a newer one starts, so an old answer cannot arrive after a new one and
// overwrite it. And the results region announces its own changes, so a
// keyboard reaches a list that has quietly become something else.
//
// Without script it is a form: pressing Enter submits it to the same endpoint
// and the whole page comes back with the results in it. That is not a
// fallback bolted on -- it is the same URL answering the same parameter, and
// it is what makes the box work before the script has loaded.
//
// It publishes root, label, input, indicator and results.
type ActiveSearchProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Name is the query parameter the typed text is sent under, and the id the
	// label points at.
	Name string
	// Label is the text above the box. It may be drawn or hidden, but it is
	// never absent: a search box whose only name is its placeholder has no
	// name once anything is typed into it.
	Label string
	// LabelHidden keeps the name and draws nothing, for a box whose position
	// on the page already says what it searches.
	LabelHidden bool
	// URL is the endpoint. It answers with the results fragment for a swap,
	// and with a page for a submit.
	URL string
	// ResultsID is the id of the region the answer is swapped into. Empty
	// means Name with "-results" after it.
	ResultsID string
	// Value is what is in the box, which for a page reached by submitting is
	// what was searched for.
	Value string
	// Placeholder is the grey text in the empty box.
	Placeholder string
	// Delay is how long typing pauses before the search runs, in
	// milliseconds. Zero waits 300.
	Delay int
	// MinLength is how many characters are needed before anything is asked
	// for. Zero asks from the first, which for a large set is a request that
	// returns everything.
	MinLength int
	// Hint is the sentence under the box.
	Hint string
	// EmptyMessage is what the results region says before anything has been
	// searched for.
	EmptyMessage string
}

// Region is the id of the results region.
func (p ActiveSearchProps) Region() string {
	if p.ResultsID != "" {
		return p.ResultsID
	}
	return p.Name + "-results"
}

// Trigger is when the search runs.
//
// It is three things at once, and each one is load-bearing. "input changed"
// runs only when the text actually differs, so an arrow key or a modifier does
// not search again. "delay" waits for typing to stop. "search" is HTMX's own
// event for the box being cleared with the native clear button, which fires no
// input event in some engines.
func (p ActiveSearchProps) Trigger() string {
	delay := p.Delay
	if delay <= 0 {
		delay = 300
	}
	return "input changed delay:" + strconv.Itoa(delay) + "ms, search"
}

// DescribedBy is the id of the hint, and empty when there is none.
func (p ActiveSearchProps) DescribedBy() string {
	if p.Hint != "" {
		return p.Name + "-hint"
	}
	return ""
}

// PartNames are the parts this component publishes.
func (p ActiveSearchProps) PartNames() []string {
	return []string{"root", "label", "group", "input", "indicator", "results", "hint"}
}
@endgo

{{-- A form, so Enter submits to the same endpoint and the page comes back with
     the results in it. HTMX takes the typing before that ever happens; what is
     left is the path for a browser that has not run the script yet. --}}
<form
	data-part="root"
	class="{{ .RootClass("field") }}"
	action="{{ .URL }}"
	method="get"
	role="search"
	@attributes(.RootAttrs())
>
	<label
		data-part="label"
		@if(.LabelHidden)
			class="{{ .PartClass("label", "sr-only") }}"
		@endif
		@if(!.LabelHidden)
			class="{{ .PartClass("label", "label") }}"
		@endif
		for="{{ .Name }}"
		@attributes(.PartAttrs("label"))
	>{{ .Label }}</label>

	<div
		data-part="group"
		class="{{ .PartClass("group", "input-group") }}"
		@attributes(.PartAttrs("group"))
	>
		{{-- hx-sync replace is what stops an old answer overwriting a new one:
		     the request in flight is abandoned the moment another starts, so
		     the last answer to arrive is always the last one asked for. --}}
		<input
			data-part="input"
			@if(.PartClass("input") != "")
				class="{{ .PartClass("input") }}"
			@endif
			type="search"
			id="{{ .Name }}"
			name="{{ .Name }}"
			value="{{ .Value }}"
			autocomplete="off"
			hx-get="{{ .URL }}"
			hx-trigger="{{ .Trigger() }}"
			hx-target="#{{ .Region() }}"
			hx-swap="innerHTML"
			hx-sync="this:replace"
			hx-indicator="#{{ .Name }}-indicator"
			@attributes(.PartAttrs("input"))
			@if(.Placeholder != "")
				placeholder="{{ .Placeholder }}"
			@endif
			@if(.MinLength > 0)
				minlength="{{ .MinLength }}"
			@endif
			@if(.DescribedBy() != "")
				aria-describedby="{{ .DescribedBy() }}"
			@endif
		>

		<span
			data-part="indicator"
			class="{{ .PartClass("indicator", "spinner") }}"
			id="{{ .Name }}-indicator"
			data-align="end"
			aria-hidden="true"
			@attributes(.PartAttrs("indicator"))
		></span>
	</div>

	@if(.Hint != "")
		<p
			data-part="hint"
			id="{{ .Name }}-hint"
			class="{{ .PartClass("hint", "text-muted-foreground text-sm") }}"
			@attributes(.PartAttrs("hint"))
		>{{ .Hint }}</p>
	@endif

	{{-- The region is on the page before any answer arrives, and it is polite
	     rather than assertive: results appearing are not an interruption. It is
	     announced as a region so a keyboard can reach it after typing rather
	     than tabbing forward through whatever is now inside it. --}}
	<div
		data-part="results"
		@if(.PartClass("results") != "")
			class="{{ .PartClass("results") }}"
		@endif
		id="{{ .Region() }}"
		role="region"
		aria-live="polite"
		aria-busy="false"
		aria-label="{{ .Label }}"
		@attributes(.PartAttrs("results"))
	>{{ .EmptyMessage }}</div>
</form>
