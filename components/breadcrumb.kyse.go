//go:build kyse

package components

import "github.com/arandu-io/kyse/icons"

@go
// BreadcrumbProps is the trail from the front of a site down to where somebody
// is standing.
//
// The last crumb is the page being looked at, so it is drawn as text and not as
// a link, and it carries aria-current -- which is the part that is announced.
// A trail whose last entry links to the page it is already on is the mistake
// this shape cannot make: the current page is taken from the end of the list
// and its address is never read.
//
// The separators are drawn here rather than by the caller, and they are marked
// hidden. A screen reader that reads the arrows aloud reads the trail as
// "Home, greater than, Posts, greater than", which is four words of punctuation
// for three of content.
type BreadcrumbProps struct {
	// Items are the crumbs, from the outermost down to the page itself. Empty
	// draws an empty trail rather than a stray separator.
	Items []Crumb
	// Label names the landmark for a screen reader, so a page with two
	// navigations is not two of the same word. Empty means "Breadcrumb".
	Label string
}

// Crumb is one step of the trail.
type Crumb struct {
	// Label is the text.
	Label string
	// URL is where the step goes. Empty draws it as plain text, which is what a
	// grouping with no page behind it is. It is ignored on the last crumb.
	URL string
}

// Trail is every crumb above the current page.
func (p BreadcrumbProps) Trail() []Crumb {
	if len(p.Items) == 0 {
		return nil
	}
	return p.Items[:len(p.Items)-1]
}

// Current is the page being looked at, which is the last crumb.
func (p BreadcrumbProps) Current() Crumb {
	if len(p.Items) == 0 {
		return Crumb{}
	}
	return p.Items[len(p.Items)-1]
}

// Landmark is what the navigation is announced as.
func (p BreadcrumbProps) Landmark() string {
	if p.Label == "" {
		return "Breadcrumb"
	}
	return p.Label
}
@endgo

<nav class="breadcrumb" aria-label="{{ .Landmark() }}">
	<ol>
		@foreach(.Trail() as crumb)
			<li>
				@if(crumb.URL != "")
					<a href="{{ crumb.URL }}">{{ crumb.Label }}</a>
				@else
					<span>{{ crumb.Label }}</span>
				@endif
			</li>
			<li aria-hidden="true" data-rtl-flip>{!! icons.CaretRight(icons.Props{}) !!}</li>
		@endforeach
		@if(len(.Items) > 0)
			<li><span aria-current="page">{{ .Current().Label }}</span></li>
		@endif
	</ol>
</nav>
