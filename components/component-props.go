package components

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/arandu-io/hesape/view"
	"github.com/arandu-io/kyse"
)

// ComponentProps is what every component takes besides its own data: the class,
// the attributes and the per-part overrides a caller adds to this one instance.
//
// It is embedded, so it is written as a field of the props it belongs to:
//
//	components.ButtonGroup(components.ButtonGroupProps{
//		Label: "Message actions",
//		ComponentProps: components.ComponentProps{
//			Class: "w-full",
//			Parts: components.Parts{
//				"label":  {Class: "text-xs"},
//				"button": {Class: "rounded-xl"},
//			},
//		},
//		Buttons: []components.ButtonProps{{Label: "Archive"}},
//	})
//
// # Why parts exist, rather than a class on the root and nothing else
//
// A component draws more than one element, and the one a caller wants to reach
// is usually not the outermost. Without a name for the inner ones the only way
// to reach them is a selector written against the structure -- `div > div:nth-child(2) > button`
// -- which is a promise the component never made and breaks the first time an
// element moves.
//
// So each component publishes the names it draws, and each named element
// carries data-part. What a caller writes is the name, and what the component
// owes in return is that the name keeps meaning the same element. PartNames is
// where that promise is written down, and a test compares it against the
// data-part attributes the component actually renders, in both directions: a
// name that is published and not drawn fails, and one that is drawn and not
// published fails too.
//
// "root" is the outermost element and every component has it. Class and Attrs
// are the same thing as Parts["root"], spelled shorter because it is what most
// calls want.
//
// # What this does not decide
//
// It does not replace the component's own classes, it adds to them -- see
// PartClass. And a class written here has to be a literal in a file the
// stylesheet is compiled from, because that is how the class ends up in the
// stylesheet at all; a class assembled at run time is a name no rule was
// emitted for, and the element renders unstyled. The same limit the width of a
// progress bar is written around.
type ComponentProps struct {
	// Class is added to the outermost element, after the component's own.
	// Shorthand for Parts["root"].Class.
	Class string
	// Replace writes Class instead of adding to it, so the outermost element
	// carries only what this call wrote. Shorthand for Parts["root"].Replace,
	// and documented there.
	Replace bool
	// Attrs are written on the outermost element. Shorthand for
	// Parts["root"].Attrs.
	Attrs Attrs
	// Parts is what to add to each named element, by the names the component
	// publishes in PartNames.
	//
	// A name no component publishes is a mistake nothing can catch here: the
	// map takes any string. It is caught by writing the name and seeing nothing
	// change, which is why PartNames is on every component and is part of what
	// pkg.go.dev shows.
	Parts Parts
	// Theme names a variation the stylesheet defines, written on the outermost
	// element as data-theme. Empty takes the one Configure set, and the theme
	// on an instance beats the global.
	Theme string
	// Style is a block of CSS scoped to this component, compiled into the
	// project's stylesheet at build time and reached by a class. Zero writes
	// nothing.
	Style kyse.CSS
	// Behavior names a client behaviour registered with arandu.ui.define, and
	// the props handed to it. Zero writes nothing.
	Behavior Behavior
	// Events maps a DOM event to the name of an action registered with
	// arandu.ui.action. Zero writes nothing.
	Events Events
}

// Attrs are attributes written on an element as they are given.
//
// Names are checked when they are rendered, not here, and what is refused is
// every name whose value a browser or a library would act on: an event handler,
// an Alpine directive, a style, an address, and the class, role and aria- a
// component owns. Those have somewhere else to go -- a field of the component,
// or the class field beside this one.
//
// They are written after every attribute the component writes itself, on
// every element of every component. An HTML parser keeps the first of two
// attributes with the same name and drops the second, so on a collision the
// component's value is the one the browser sees: an Attrs entry can add an
// attribute the component left out, and can never replace one it wrote --
// the name, value, type and form a control submits with, or the id that a
// label and an aria-describedby point at. A value for one of those goes in the
// component's own field.
type Attrs map[string]string

// Parts is what a caller adds to each named element of a component.
type Parts map[string]PartProps

// PartProps is the class and the attributes for one named element.
type PartProps struct {
	// Class is added to that element's own classes, after them.
	Class string
	// Replace writes Class instead of adding to it, so the element carries
	// only what the caller wrote.
	//
	// # Why adding is not enough on its own
	//
	// Adding covers the common case and cannot cover all of it. A caller who
	// wants one more unit of padding adds a class and wins on the cascade; a
	// caller who wants the border gone, or the background gone, or the whole
	// look of a control replaced, cannot -- there is no utility that unsets
	// what another rule set, and a component class that draws a border will
	// keep drawing it however many classes are written after it. The answer
	// used to be to fight the specificity, which works until the stylesheet
	// changes.
	//
	// So the caller says which they meant. Adding stays the default, because
	// it is what almost every call wants and it costs the caller nothing.
	//
	// # What it does not take away
	//
	// The element, its data-part, its role and its ARIA. Those are what the
	// component is; the class is what it looks like. A part that stopped being
	// reachable, or a control that stopped announcing itself, would be a
	// different component rather than the same one restyled.
	Replace bool
	// Attrs are written on that element.
	Attrs Attrs
}

// PartClass is the class attribute of one part: what the component draws with,
// then what the caller added.
//
// # Why the caller's classes are added and not substituted
//
// Substituting would mean a caller who wants one more unit of padding has to
// know, and repeat, every class the component draws itself with -- and repeat
// them again whenever the component changes. Adding costs nothing to the caller
// who wants nothing.
//
// # Why the order in the attribute is not the point
//
// class="btn p-8" and class="p-8 btn" are the same document to a browser: the
// attribute is a set, and nothing in CSS reads its order. What decides is the
// layer, and the two live in different ones -- the component classes are
// declared in @layer components and a Tailwind utility in @layer utilities,
// and utilities win, above specificity. So a utility written here takes effect
// against a component class without an !important and without depending on
// where it sits in the string.
//
// The collision the layer does not settle is a utility against a utility: a
// component that draws itself with p-5 and a caller who writes p-8 put two
// rules in the same layer, and the order between them is Tailwind's canonical
// one rather than the one either of them wrote. That is what the importance
// marker is for -- p-8! -- and it is the only case that needs it.
//
// Empty pieces and repeats are dropped, so a caller who writes a class the
// component already draws with does not double it.
// own is variadic because a component's own classes are not always one string:
// a skeleton's shape and its width are decided separately, and a component that
// had to concatenate them first would be building a class list with a + before
// handing it to the thing whose job is building class lists.
func (c ComponentProps) PartClass(part string, own ...string) string {
	props := c.partProps(part)
	if props.Replace {
		return joinClasses(props.Class)
	}
	return joinClasses(append(own, props.Class)...)
}

// PartAttrs is the attributes to write on one part.
//
// The result goes through view.Attributes, which is what decides whether a name
// may be written at all; nothing here refuses one. What this decides is only
// which part a caller was addressing.
//
// A name the view layer refuses -- an hx- or on- attribute, a style, an
// address, class, role or aria- -- panics with a *RefusedAttribute; see there
// for why it is not left to the render.
func (c ComponentProps) PartAttrs(part string) map[string]string {
	return checkedAttrs(part, c.partProps(part).Attrs)
}

// RefusedAttribute is what a component panics with when its caller handed it
// an attribute the view layer will not write.
//
// # Why a panic and not the empty component
//
// The render refuses the same name, and a component that refused anything
// returns nothing at all: the markup built up to the refusal ends inside an
// attribute with its quote still open, and handing that back would give the
// rest of the page to an element the browser is still reading. That is right
// for a value that arrived as data, and it stays the answer for one -- an
// address whose scheme a browser acts on draws nothing.
//
// It is wrong for a name the caller wrote. Attrs{"hx-boost": "false"} on a
// FormDialog drew an empty space where the dialog should be, with nothing in a
// log, a test or the console to say why: the screen looked as if the dialog had
// never been written. A name in a view is the view's mistake, it is the same on
// every render, and it is found the first time the page is drawn -- by a test
// that renders it, or by the exception handler that turns the panic into an
// error page in development. So it stops the render where the mistake is,
// naming the part, the attribute and the reason.
type RefusedAttribute struct {
	// Part is the part the attribute was written for; "root" for Attrs.
	Part string
	// Name is the attribute.
	Name string
	// Err is the view layer's reason.
	Err error
}

// Error names the part, the attribute and the reason.
func (r *RefusedAttribute) Error() string {
	return fmt.Sprintf("kyse: the attribute %q written for the part %q is refused, "+
		"and the component would have drawn nothing: %v", r.Name, r.Part, r.Err)
}

// Unwrap is the view layer's reason.
func (r *RefusedAttribute) Unwrap() error { return r.Err }

// checkedAttrs answers attrs, or panics with the first name the view layer
// refuses, in the order the render would reach them.
func checkedAttrs(part string, attrs map[string]string) map[string]string {
	if _, err := view.Attributes(attrs); err == nil {
		return attrs
	}
	names := make([]string, 0, len(attrs))
	for name := range attrs {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		if _, err := view.Attributes(map[string]string{name: attrs[name]}); err != nil {
			panic(&RefusedAttribute{Part: part, Name: name, Err: err})
		}
	}
	return attrs
}

// RootAttrs is everything written on the outermost element besides its class:
// what the caller handed for the root, plus the theme, the behaviour and the
// events, which are all attributes too.
//
// They are one bag rather than six conditionals in every component's markup.
// The alternative is worth naming, because it looks simpler until it is
// written: an event maps to data-kyse-on-<event>, and a name assembled from
// data is exactly what markup cannot write -- the compiler reads the name off
// the source to decide the escape, and there is no source to read. Putting them
// through the bag puts them through the one check that reads names.
//
// The component's own values win a collision. A caller who writes
// data-kyse-behavior by hand is asking for the field beside it.
//
// A name the view layer refuses panics with a *RefusedAttribute, as PartAttrs
// does.
func (c ComponentProps) RootAttrs() map[string]string {
	out := map[string]string{}
	for name, value := range c.partProps("root").Attrs {
		out[name] = value
	}
	if theme := c.ThemeName(); theme != "" {
		out["data-theme"] = theme
	}
	if c.Behavior.Name != "" {
		out["data-kyse-behavior"] = c.Behavior.Name
		if props := c.Behavior.PropsJSON(); props != "" {
			out["data-kyse-props"] = props
		}
	}
	for _, event := range c.Events.EventNames() {
		out["data-kyse-on-"+event] = c.Events[event]
	}
	return checkedAttrs("root", out)
}

// ThemeName is the theme in force for this instance: the one it was given, or
// the one Configure set.
//
// Empty means no data-theme is written, which is the state every rule written
// for `:not([data-theme])` applies in.
func (c ComponentProps) ThemeName() string {
	if c.Theme != "" {
		return c.Theme
	}
	return configured().Theme
}

// StyleClass is the class the scoped block is compiled under, or empty when
// there is no block.
func (c ComponentProps) StyleClass() string { return c.Style.Class() }

// RootClass is the class attribute of the outermost element: the component's
// own, the caller's, and the scoped block's.
//
// It is a separate method rather than PartClass("root", own) because the root
// is the one element the scoped block attaches to, and a component should not
// have to remember to append it.
func (c ComponentProps) RootClass(own ...string) string {
	props := c.partProps("root")
	if props.Replace {
		// The scoped block's class stays: it names CSS this call itself asked
		// to be compiled, so dropping it would silently discard the Style the
		// same caller wrote.
		return joinClasses(props.Class, c.Style.Class())
	}
	return joinClasses(append(own, props.Class, c.Style.Class())...)
}

// partProps is what the caller wrote for one part, with root reading the
// shorthand fields as well.
//
// The shorthand loses to the explicit form on a collision, and the collision is
// a caller writing both Class and Parts["root"].Class in one call. Neither
// answer is obviously right, so the one that reads as more deliberate wins.
func (c ComponentProps) partProps(part string) PartProps {
	p := c.Parts[part]
	if part == "root" && c.Replace {
		p.Replace = true
	}
	if part != "root" {
		return p
	}
	if p.Class == "" {
		p.Class = c.Class
	}
	if len(p.Attrs) == 0 {
		p.Attrs = c.Attrs
	}
	return p
}

// joinClasses joins class lists with a space, dropping empty pieces and any
// repeat of a piece already seen.
//
// Deduplicating rather than concatenating, because a caller repeating a class
// the component already draws with is writing the same rule twice, and an
// attribute that says "btn btn" is a diff nobody can read.
func joinClasses(lists ...string) string {
	seen := map[string]bool{}
	kept := make([]string, 0, len(lists))
	for _, list := range lists {
		for _, class := range strings.Fields(list) {
			if seen[class] {
				continue
			}
			seen[class] = true
			kept = append(kept, class)
		}
	}
	return strings.Join(kept, " ")
}

// Behavior names a client behaviour and the props it is given.
//
// The name is looked up in the registry a project fills with
// arandu.ui.define(name, {mounted, updated, destroyed}); nothing here is a
// script and nothing anywhere evaluates one. That is the difference between
// this and an inline handler, and it is why the policy can stay
// script-src 'self' with no unsafe-eval.
type Behavior struct {
	// Name is the registered behaviour. Empty writes nothing.
	Name string
	// Props are handed to the behaviour as its ctx.props, encoded as JSON into
	// an attribute. A value that will not encode is dropped along with the
	// whole set, and the behaviour still mounts with no props -- see
	// BehaviorProps.
	Props map[string]any
}

// PropsJSON is Props as the JSON the behaviour reads from ctx.props, or empty
// when there are none.
//
// The attribute it lands in is escaped as HTML, so the quotes arrive as &#34;
// and the browser decodes them before any script reads the value -- which is
// what makes JSON in an attribute work, and is why the encoding is done here
// rather than by whoever writes the markup.
//
// A value that will not encode -- a channel, a cycle -- drops the whole set
// rather than half of it. The behaviour then mounts with no props, which is
// visible the first time it runs and is better than mounting with props that
// silently lost a field.
func (b Behavior) PropsJSON() string {
	if b.Name == "" || len(b.Props) == 0 {
		return ""
	}
	encoded, err := json.Marshal(b.Props)
	if err != nil {
		return ""
	}
	return string(encoded)
}

// Events maps a DOM event name to the action that answers it.
//
//	Events{"click": "archive-message"}
//
// The action is a name registered with arandu.ui.action(name, fn), and the
// event is one of the ones ui.js listens for on the document: click, input,
// change, keydown and submit. An event outside that set writes an attribute
// nothing reads, which is visible on the page and therefore fixable.
type Events map[string]string

// EventNames are the events in this set, sorted, so a component draws them in
// one order.
func (e Events) EventNames() []string {
	if len(e) == 0 {
		return nil
	}
	out := make([]string, 0, len(e))
	for name := range e {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
