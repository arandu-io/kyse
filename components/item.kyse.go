//go:build kyse

package components

@go
// ItemProps is one row of a list: a picture, a line, a sentence, and the
// control that acts on it.
//
// It is the shape half the screens in an application are made of -- a setting,
// a member of a team, a connected account, a file, a device. The reason it is a
// component rather than four elements is that the four parts have to keep the
// same order and the same spacing on every one of those screens, and a row
// written by hand keeps them until the day somebody adds a fifth.
//
// Several of them go inside <div class="item-group">, which draws the gap
// between rows and the rule across a divider.
//
// # It does not navigate on its own
//
// A whole row that is a link is <a class="item"> written where it is used. A
// component that chose between an anchor and an article would have to write out
// everything inside it twice, once under each element, and two copies of the
// same four parts is the drift this exists to prevent.
// It publishes root, icon, content, title, description and action.
type ItemProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Title is the line, and what the row is about.
	Title string
	// Description is the sentence under it. Empty draws none, which is the
	// dense row: a name and a control, nothing else.
	Description string
	// Icon is the picture on the left. Empty draws none, and the row closes up
	// around it.
	//
	// Nothing here escapes it, so what goes in is what the page gets, and the
	// caller is what makes that safe: pass what an icon function returned, and
	// never a string assembled around a value somebody typed.
	Icon template.HTML
	// Action is the control on the right: a button, a link, a badge, a
	// chevron. Empty draws none.
	//
	// It is markup for the same reason Icon is, and carries the same
	// obligation: pass what another component returned rather than a string
	// built here.
	Action template.HTML
	// Variant is "default", "outline" or "muted". Empty is the default, which
	// draws no border -- the row that belongs to a list rather than standing on
	// its own.
	Variant string
	// Size is "default", "sm" or "xs". Empty is the default.
	Size string
}

// PartNames are the parts this component publishes.
func (p ItemProps) PartNames() []string {
	return []string{"root", "icon", "content", "title", "description", "action"}
}
@endgo

<article
	data-part="root"
	class="{{ .RootClass("item") }}"
	@attributes(.RootAttrs())
	@if(.Variant != "")
		data-variant="{{ .Variant }}"
	@endif
	@if(.Size != "")
		data-size="{{ .Size }}"
	@endif
>
	@if(.Icon != "")
		<figure
			data-part="icon"
			@if(.PartClass("icon") != "")
				class="{{ .PartClass("icon") }}"
			@endif
			@attributes(.PartAttrs("icon"))
		>{!! .Icon !!}</figure>
	@endif
	<section
		data-part="content"
		@if(.PartClass("content") != "")
			class="{{ .PartClass("content") }}"
		@endif
		@attributes(.PartAttrs("content"))
	>
		<h3
			data-part="title"
			@if(.PartClass("title") != "")
				class="{{ .PartClass("title") }}"
			@endif
			@attributes(.PartAttrs("title"))
		>{{ .Title }}</h3>
		@if(.Description != "")
			<p
				data-part="description"
				@if(.PartClass("description") != "")
					class="{{ .PartClass("description") }}"
				@endif
				@attributes(.PartAttrs("description"))
			>{{ .Description }}</p>
		@endif
	</section>
	@if(.Action != "")
		<aside
			data-part="action"
			@if(.PartClass("action") != "")
				class="{{ .PartClass("action") }}"
			@endif
			@attributes(.PartAttrs("action"))
		>{!! .Action !!}</aside>
	@endif
</article>
