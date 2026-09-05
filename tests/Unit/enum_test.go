package unit

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/enum"
	"github.com/arandu-io/kyse/components"
)

// The two fixtures below are what the generator writes, cut down to the four
// methods the contract names. They are here rather than imported because the
// property under test is that a form follows the type it was pointed at, and
// that can only be shown by moving a type this suite owns.
//
// The assertion is at package scope so that a type which stopped satisfying the
// contract fails to build rather than failing one test.
var (
	_ enum.Enum = invoiceStatusDraft
	_ enum.Enum = priorityLow
)

// invoiceStatus is a set stored as text, where the shown spelling and the
// stored one are the same string.
type invoiceStatus string

const (
	invoiceStatusDraft invoiceStatus = "draft"
	invoiceStatusSent  invoiceStatus = "sent"
	invoiceStatusPaid  invoiceStatus = "paid"
)

// invoiceStatusValues lists them in declaration order.
func invoiceStatusValues() []invoiceStatus {
	return []invoiceStatus{invoiceStatusDraft, invoiceStatusSent, invoiceStatusPaid}
}

func (v invoiceStatus) Valid() bool {
	switch v {
	case invoiceStatusDraft, invoiceStatusSent, invoiceStatusPaid:
		return true
	}
	return false
}

func (v invoiceStatus) String() string { return string(v) }

func (v invoiceStatus) Label() string {
	switch v {
	case invoiceStatusDraft:
		return "Draft"
	case invoiceStatusSent:
		return "Sent"
	case invoiceStatusPaid:
		return "Paid"
	}
	return v.String()
}

func (v invoiceStatus) Value() (driver.Value, error) {
	if !v.Valid() {
		return nil, fmt.Errorf("invoice status: refusing to write %q", v.String())
	}
	return string(v), nil
}

// priority is a set stored as a number, where the shown spelling is a word the
// column never holds.
type priority int

const (
	priorityLow    priority = 1
	priorityNormal priority = 2
	priorityHigh   priority = 3
)

// priorityValues lists them in declaration order.
func priorityValues() []priority {
	return []priority{priorityLow, priorityNormal, priorityHigh}
}

func (v priority) Valid() bool {
	switch v {
	case priorityLow, priorityNormal, priorityHigh:
		return true
	}
	return false
}

func (v priority) String() string {
	switch v {
	case priorityLow:
		return "low"
	case priorityNormal:
		return "normal"
	case priorityHigh:
		return "high"
	}
	return fmt.Sprintf("priority(%d)", int(v))
}

func (v priority) Label() string {
	switch v {
	case priorityLow:
		return "Low"
	case priorityNormal:
		return "Normal"
	case priorityHigh:
		return "High"
	}
	return v.String()
}

func (v priority) Value() (driver.Value, error) {
	if !v.Valid() {
		return nil, fmt.Errorf("priority: refusing to write %q", v.String())
	}
	return int64(v), nil
}

// drawnOption is one line as the rendered list carries it.
type drawnOption struct {
	value string
	text  string
}

// optionPattern reads the lines out of a rendered list.
//
// The attributes are written one per line, so the run before the value and the
// run after it are matched as "anything that is not the end of the tag" rather
// than with the dot-matches-newline flag, which would also let the two runs
// swallow the tag they are inside.
var optionPattern = regexp.MustCompile(`<option\b[^>]*\bvalue="([^"]*)"[^>]*>([^<]*)</option>`)

// drawnOptions is every line of a rendered list, in document order.
func drawnOptions(html string) []drawnOption {
	found := optionPattern.FindAllStringSubmatch(html, -1)
	lines := make([]drawnOption, 0, len(found))
	for _, one := range found {
		lines = append(lines, drawnOption{value: one[1], text: one[2]})
	}
	return lines
}

// TestTheListDrawnIsTheListTheTypeDeclares is the property the function exists
// for, and it is written so that nothing here can hold a case.
//
// What is expected is read off invoiceStatusValues, so a case added to the type
// changes both sides of the comparison and no line of this file names one. A
// test that spelled the three out would be the second copy the function was
// written to remove, and it would agree on the day it was written.
func TestTheListDrawnIsTheListTheTypeDeclares(t *testing.T) {
	html := string(components.Select(components.SelectProps{
		Name:    "status",
		Label:   "Status",
		Options: components.SelectOptionsFrom(invoiceStatusValues()...),
	}))

	drawn := drawnOptions(html)
	cases := invoiceStatusValues()
	if len(drawn) != len(cases) {
		t.Fatalf("the type declares %d cases and the list drew %d lines: %v\n%s",
			len(cases), len(drawn), drawn, html)
	}
	for at, one := range cases {
		if drawn[at].value != one.String() {
			t.Errorf("line %d submits %q, and case %d of the type is %q",
				at, drawn[at].value, at, one.String())
		}
		if drawn[at].text != one.Label() {
			t.Errorf("line %d reads %q, and case %d of the type is labelled %q",
				at, drawn[at].text, at, one.Label())
		}
	}
}

// TestALineSubmitsTheShownSpellingAndNotTheStoredOne is the trap a set stored
// as a number sets.
//
// The column holds 2 and the parser beside the type reads "normal", so a list
// that sent what the column holds would come back as a value nothing can read
// -- and it would look right in the markup, because a number is a plausible
// value for a form to send.
func TestALineSubmitsTheShownSpellingAndNotTheStoredOne(t *testing.T) {
	html := string(components.Select(components.SelectProps{
		Name:    "priority",
		Label:   "Priority",
		Options: components.SelectOptionsFrom(priorityValues()...),
	}))

	for _, one := range priorityValues() {
		if !strings.Contains(html, `value="`+one.String()+`"`) {
			t.Errorf("%q is a case and no line submits it:\n%s", one.String(), html)
		}
		stored, err := one.Value()
		if err != nil {
			t.Fatalf("%v is a case and Value refused it: %v", one, err)
		}
		if written := fmt.Sprint(stored); strings.Contains(html, `value="`+written+`"`) {
			t.Errorf("a line submits %q, which is what the column holds and not what the parser reads:\n%s",
				written, html)
		}
	}
}

// TestAFormMayDrawPartOfTheSet holds the other half of taking the cases: the
// form that shows two of three says so by naming two, and the two it names are
// still checked by the compiler.
func TestAFormMayDrawPartOfTheSet(t *testing.T) {
	html := string(components.Select(components.SelectProps{
		Name:    "status",
		Label:   "Status",
		Options: components.SelectOptionsFrom(invoiceStatusDraft, invoiceStatusSent),
	}))

	drawn := drawnOptions(html)
	if len(drawn) != 2 {
		t.Fatalf("two cases were named and the list drew %d lines: %v", len(drawn), drawn)
	}
	if strings.Contains(html, `value="`+invoiceStatusPaid.String()+`"`) {
		t.Errorf("a case nobody named was drawn:\n%s", html)
	}
}

// TestNoCaseCarriesAStateTheTypeCannotHold keeps the returned line to what the
// set knows about.
//
// A case is either declared or it is not; there is nothing on the type saying
// one of them is unavailable this afternoon. Deriving that would be inventing
// an answer, so the line comes back available and the screen that knows better
// says so itself.
func TestNoCaseCarriesAStateTheTypeCannotHold(t *testing.T) {
	for at, option := range components.SelectOptionsFrom(invoiceStatusValues()...) {
		if option.Disabled {
			t.Errorf("line %d came back disabled, and no case says it is", at)
		}
	}
}

// TestASetWithNoCasesNamedDrawsNoLines is the empty call, which a form reaches
// by filtering a set down to nothing.
func TestASetWithNoCasesNamedDrawsNoLines(t *testing.T) {
	if options := components.SelectOptionsFrom[invoiceStatus](); len(options) != 0 {
		t.Errorf("no case was named and %d lines came back: %v", len(options), options)
	}
}
