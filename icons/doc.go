// Package icons draws the Phosphor icon set, one exported function per icon.
//
//	{!! icons.Trash(icons.Props{Label: "Delete this post"}) !!}
//
// The shape is the reason the package can exist at all. Phosphor is 1512 icons,
// six hundred kilobytes of path data, and an application uses thirty of them.
// Embedding the set with go:embed, or reaching it through a map keyed by name,
// would retain the whole catalogue because the linker cannot prove a map lookup
// will not ask for the rest. One function per icon it can prove: an icon nobody
// calls is dead code, and dead code and its string data are dropped.
//
// It is also the same argument the components package makes. icons.Trahs is
// undefined at the line of the .kyse.go that wrote it, where an icon named by
// string is a blank square in production.
//
// # What an icon renders as
//
// An svg with fill="currentColor", so it takes the colour of the text around
// it, and width and height of 1em, so it takes the size. Both are presentation
// attributes and lose to any CSS rule, which is what makes the stylesheet the
// one place size is decided: Basecoat sizes svg descendants of a button, an
// alert or a badge, and an icon dropped into a paragraph with no rule to match
// falls back to the height of the line rather than to the 300 by 150 pixels an
// svg with no size defaults to.
//
// There is no Size and no Class. Either would be a second place the size of an
// icon is decided, and the first one is the stylesheet the project owns.
//
// # The accessible name
//
// An icon with no Label is decorative and is marked aria-hidden, because it is
// beside a word that already says what it means and a screen reader announcing
// both reads everything twice. An icon with a Label is the label: that is the
// icon-only button, which without one is announced as "button".
//
// # Weight
//
// The duotone weight, and only it. Phosphor draws six, and shipping two would
// mean Trash and TrashDuotone: two names for one idea, and a question on every
// call site that has no good answer. The icon set is a decision the library
// makes once, like the stylesheet.
//
// Duotone because that is the weight the identity draws in, and one weight is
// still one name: Trash is duotone, and nothing is called anything else.
//
// It is two paths of one colour rather than two colours. The shade is
// currentColor at a fifth of its strength, so an icon still takes the colour of
// the text around it and still needs no palette of its own.
package icons

//go:generate go run github.com/arandu-io/kyse/internal/icongen
