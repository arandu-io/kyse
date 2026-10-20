package components

import "github.com/arandu-io/hesape/pagination"

// Paginator is the page a query returned, as hesape/pagination hands it back:
// a *pagination.LengthAwarePage, a *pagination.Page or a
// *pagination.CursorPage. [Pagination] and [DataTable] draw from it directly,
// so a list passes on what its repository returned instead of copying three
// numbers out of it:
//
//	{!! components.Pagination(components.PaginationProps{Paginator: .Pager}) !!}
//
// where Pager is the *pagination.LengthAwarePage the controller got back from
// Paginate, beside the rows.
//
// The five methods are the ones all three pages answer, so all three fit. A
// page that also knows its total -- the length-aware one -- is drawn with its
// numbers, its window and the range of rows it holds; the other two know only
// whether there is a page either side of this one, and are drawn as the two
// steps.
//
// It is declared here, in the package that reads it, rather than in the one
// whose types satisfy it: the pagination package has no interface on purpose,
// and the methods a pager reads are the pager's question.
type Paginator interface {
	// HasPages is whether there is more than one page, which is when a pager is
	// drawn at all.
	HasPages() bool
	// OnFirstPage is whether this is the first page.
	OnFirstPage() bool
	// HasMorePages is whether a page follows this one.
	HasMorePages() bool
	// PreviousPageURL and NextPageURL are the two steps, each empty when there
	// is nowhere to step to.
	PreviousPageURL() string
	NextPageURL() string
}

// countedPage is the half of a Paginator that knows where it sits: the
// length-aware page.
type countedPage interface {
	Paginator
	CurrentPage() int
	LastPage() int
	Total() int
	FirstItem() int
	LastItem() int
	URL(page int) string
	Links() []pagination.Link
	GetPageName() string
}

// numberedPage is a page that knows its own number but not the last one: the
// simple page, which only knows whether one more row exists.
type numberedPage interface {
	Paginator
	CurrentPage() int
	GetPageName() string
}

// presentPaginator is the paginator, or nil when there is none -- including a
// nil pointer of one of the three pagination types, which a list whose query
// failed can hand on and whose methods would dereference it.
func presentPaginator(p Paginator) Paginator {
	switch page := p.(type) {
	case nil:
		return nil
	case *pagination.LengthAwarePage:
		if page == nil {
			return nil
		}
	case *pagination.Page:
		if page == nil {
			return nil
		}
	case *pagination.CursorPage:
		if page == nil {
			return nil
		}
	}
	return p
}

// counted is the paginator when it knows its total.
func counted(p Paginator) (countedPage, bool) {
	page, ok := presentPaginator(p).(countedPage)
	return page, ok
}

// steps is a Paginator that is only its two addresses. A data table draws a
// page that knows no total from it, with each address built from the table's
// own state rather than from whatever query the request arrived with.
type steps struct {
	previous, next string
}

func (s steps) HasPages() bool          { return s.previous != "" || s.next != "" }
func (s steps) OnFirstPage() bool       { return s.previous == "" }
func (s steps) HasMorePages() bool      { return s.next != "" }
func (s steps) PreviousPageURL() string { return s.previous }
func (s steps) NextPageURL() string     { return s.next }
