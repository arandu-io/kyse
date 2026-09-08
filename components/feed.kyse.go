//go:build kyse

package components

import "strconv"

@go
// FeedProps is a stream of articles that keeps growing as somebody reads.
//
// It is role=feed and not a list, and the difference is the promise it makes.
// A feed says: there are more of these, they arrive as you go, and each one is
// an article you can move between with a key rather than by tabbing through
// everything inside it. A screen reader that knows it is in a feed reads it
// that way; one that finds a list of divs reads it as the page.
//
// Each article carries its position and, when it is known, the total. When it
// is not known -- which is the usual case for a feed -- the total is left off
// rather than guessed, because a feed that claims 20 of 20 and then loads
// twenty more has lied about the one thing the attribute is for.
//
// It publishes root, article, title, meta, body and footer.
type FeedProps struct {
	// ComponentProps is the class, attributes and parts the caller adds.
	ComponentProps
	// Label names the feed. Without it it announces with no subject.
	Label string
	// Items are the articles, in the order they read.
	Items []FeedItem
	// Offset is how many articles came before these, for a page that is not
	// the first. The positions carry on from it.
	Offset int
	// Total is how many there are in all. Zero leaves it unstated, which is
	// honest for a stream nobody has counted.
	Total int
	// Busy marks the feed as loading more, which is what stops a screen reader
	// announcing a half-inserted article.
	Busy bool
}

// FeedItem is one article.
type FeedItem struct {
	// ID is the article's id, which the heading and the article point at.
	ID string
	// Title is the heading.
	Title string
	// URL makes the heading a link.
	URL string
	// Meta is the line under the heading: an author, a time, a source.
	Meta string
	// Body is the text of the article.
	Body string
	// ImageURL draws a picture in it.
	ImageURL string
	// Alt is what the picture shows.
	Alt string
	// Footer is the line at the bottom: a count, a tag.
	Footer string
}

// Position is where one article sits in the whole feed, counting from one.
func (p FeedProps) Position(at int) int { return p.Offset + at + 1 }

// ArticleID is the id of one article, falling back to its position when the
// caller gave none -- an article with no id cannot be pointed at by its own
// heading.
func (p FeedProps) ArticleID(item FeedItem, at int) string {
	if item.ID != "" {
		return item.ID
	}
	return "feed-article-" + strconv.Itoa(p.Position(at))
}

// State is what aria-busy says.
func (p FeedProps) State() string {
	if p.Busy {
		return "true"
	}
	return "false"
}

// PartNames are the parts this component publishes.
func (p FeedProps) PartNames() []string {
	return []string{"root", "article", "title", "meta", "media", "body", "footer"}
}
@endgo

<div
	data-part="root"
	class="{{ .RootClass("feed") }}"
	role="feed"
	aria-label="{{ .Label }}"
	aria-busy="{{ .State() }}"
	@attributes(.RootAttrs())
>
	@for(at := 0; at < len(.Items); at++)
		{{-- tabindex on the article, because moving between articles is the
		     interaction a feed promises and an article that cannot take focus
		     cannot be moved to. --}}
		<article
			data-part="article"
			class="{{ .PartClass("article", "feed-article") }}"
			id="{{ .ArticleID(.Items[at], at) }}"
			tabindex="0"
			aria-labelledby="{{ .ArticleID(.Items[at], at) }}-title"
			aria-posinset="{{ .Position(at) }}"
			@attributes(.PartAttrs("article"))
			@if(.Total > 0)
				aria-setsize="{{ .Total }}"
			@endif
			@if(.Total == 0)
				aria-setsize="-1"
			@endif
		>
			<h3
				data-part="title"
				class="{{ .PartClass("title", "font-semibold") }}"
				id="{{ .ArticleID(.Items[at], at) }}-title"
				@attributes(.PartAttrs("title"))
			>
				@if(.Items[at].URL != "")
					<a class="hover:underline" href="{{ .Items[at].URL }}">{{ .Items[at].Title }}</a>
				@endif
				@if(.Items[at].URL == "")
					{{ .Items[at].Title }}
				@endif
			</h3>

			@if(.Items[at].Meta != "")
				<p
					data-part="meta"
					class="{{ .PartClass("meta", "text-muted-foreground text-sm") }}"
					@attributes(.PartAttrs("meta"))
				>{{ .Items[at].Meta }}</p>
			@endif

			@if(.Items[at].ImageURL != "")
				<img
					data-part="media"
					class="{{ .PartClass("media", "feed-media") }}"
					src="{{ .Items[at].ImageURL }}"
					alt="{{ .Items[at].Alt }}"
					loading="lazy"
					decoding="async"
					@attributes(.PartAttrs("media"))
				>
			@endif

			@if(.Items[at].Body != "")
				<p
					data-part="body"
					@if(.PartClass("body") != "")
						class="{{ .PartClass("body") }}"
					@endif
					@attributes(.PartAttrs("body"))
				>{{ .Items[at].Body }}</p>
			@endif

			@if(.Items[at].Footer != "")
				<p
					data-part="footer"
					class="{{ .PartClass("footer", "text-muted-foreground text-xs") }}"
					@attributes(.PartAttrs("footer"))
				>{{ .Items[at].Footer }}</p>
			@endif
		</article>
	@endfor
</div>
