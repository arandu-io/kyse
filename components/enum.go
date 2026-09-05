package components

import "github.com/arandu-io/hesape/enum"

// SelectOptionsFrom is the cases of a closed set as the lines of a [Select], in
// the order they were given.
//
//	components.Select(components.SelectProps{
//		Name:    "status",
//		Label:   "Status",
//		Options: components.SelectOptionsFrom(enums.InvoiceStatusValues()...),
//	})
//
// # Why it takes the cases and not a list of options
//
// A list of options is something anybody can write, and one written by hand is
// a second copy of a set the type already declares: it agrees on the day it is
// written and disagrees on the day a case is added, with nothing comparing the
// two. Cases cannot be written that way -- a constant that does not exist is a
// build error at the line of the view -- so the list a form draws and the list
// the type declares are one list.
//
// Naming fewer than all of them is how a form draws part of a set, and each one
// named is still checked.
//
// # What each line carries
//
// Value is the shown spelling, which is what a submitted form sends back and
// what the parser beside the type reads. For a set stored as an integer the
// shown spelling and the stored one differ, and a control that sent the number
// would come back as a value nothing can read.
//
// Label is what the line reads, kept apart from Value so that rewording a
// choice cannot change what choosing it submits.
//
// Disabled is false on every line. Whether a case can be picked today is the
// screen's answer and not the type's, so a caller with one to grey out sets it
// on the line this returns.
func SelectOptionsFrom[E enum.Enum](values ...E) []SelectOption {
	cases := enum.Options(values...)
	options := make([]SelectOption, 0, len(cases))
	for _, one := range cases {
		options = append(options, SelectOption{Label: one.Label, Value: one.Value})
	}
	return options
}
