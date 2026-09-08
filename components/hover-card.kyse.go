//go:build kyse

package components

@go
// HoverCardProps is a preview of what is behind a link, shown when the pointer
// rests on it or the keyboard focuses it.
//
// It is a preview and never the only way to reach something. Everything in the
// card is also behind the link, because a card that appears on hover is a card
// that never appears on a touch screen -- there is no hover there, and there
// never will be. So the link works, and this is what a pointer gets in
// addition.
//
// It draws its own trigger, for the same reason a tooltip does: the card has
// to be part of the same hover target as the link, or it vanishes while the
// pointer is travelling to it.
//
// It publishes root, trigger, content, avatar, title, subtitle, description
// and meta.
type HoverCardProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// ID is what the trigger and the card use to point at each other.
	ID string
	// Label is the text of the link.
	Label string
	// URL is where the link goes. The card previews it; the link reaches it.
	URL string
	// Title is the heading inside the card. Empty uses the label.
	Title string
	// Subtitle is the smaller line under it: a handle, a role, a date.
	Subtitle string
	// Description is the sentence.
	Description string
	// ImageURL draws a picture at the top of the card: an avatar, a cover.
	ImageURL string
	// Alt is what the picture shows, for somebody who cannot see it.
	Alt string
	// Meta is the grey line at the bottom: a count, a joined date.
	Meta string
	// Side is where the card opens: "top", "bottom", "left" or "right". Empty
	// opens downward.
	Side string
}

// ContentID is the id of the card, which the trigger points at.
func (p HoverCardProps) ContentID() string { return p.ID + "-card" }

// Heading is the title inside the card, falling back to the link's own text --
// a card whose heading is missing is a card with no subject.
func (p HoverCardProps) Heading() string {
	if p.Title != "" {
		return p.Title
	}
	return p.Label
}

// PartNames are the parts this component publishes.
func (p HoverCardProps) PartNames() []string {
	return []string{"root", "trigger", "content", "avatar", "title", "subtitle", "description", "meta"}
}
@endgo

<span
	data-part="root"
	class="{{ .RootClass("hover-card") }}"
	@if(.Side != "")
		data-side="{{ .Side }}"
	@endif
	@attributes(.RootAttrs())
>
	<a
		data-part="trigger"
		class="{{ .PartClass("trigger", "link") }}"
		href="{{ .URL }}"
		aria-describedby="{{ .ContentID() }}"
		@attributes(.PartAttrs("trigger"))
	>{{ .Label }}</a>

	{{-- The card is not hidden from assistive technology, and it is not a live
	     region either. It is a description the trigger points at, which is read
	     when the trigger is reached and not when the card appears -- so it says
	     the same thing to a pointer and to a keyboard. --}}
	<span
		data-part="content"
		class="{{ .PartClass("content", "hover-card-content") }}"
		id="{{ .ContentID() }}"
		role="note"
		@if(.Side != "")
			data-side="{{ .Side }}"
		@endif
		@attributes(.PartAttrs("content"))
	>
		@if(.ImageURL != "")
			<img
				data-part="avatar"
				class="{{ .PartClass("avatar", "hover-card-avatar") }}"
				src="{{ .ImageURL }}"
				alt="{{ .Alt }}"
				@attributes(.PartAttrs("avatar"))
			>
		@endif
		<span
			data-part="title"
			class="{{ .PartClass("title", "font-semibold") }}"
			@attributes(.PartAttrs("title"))
		>{{ .Heading() }}</span>
		@if(.Subtitle != "")
			<span
				data-part="subtitle"
				class="{{ .PartClass("subtitle", "text-muted-foreground text-sm") }}"
				@attributes(.PartAttrs("subtitle"))
			>{{ .Subtitle }}</span>
		@endif
		@if(.Description != "")
			<span
				data-part="description"
				class="{{ .PartClass("description", "text-sm") }}"
				@attributes(.PartAttrs("description"))
			>{{ .Description }}</span>
		@endif
		@if(.Meta != "")
			<span
				data-part="meta"
				class="{{ .PartClass("meta", "text-muted-foreground text-xs") }}"
				@attributes(.PartAttrs("meta"))
			>{{ .Meta }}</span>
		@endif
	</span>
</span>
