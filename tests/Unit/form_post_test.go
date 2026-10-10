package unit

import (
	"html"
	"regexp"
	"strings"
)

// The form a browser posts when no script runs, read from rendered markup.
//
// A component whose behaviour fills in the field that submits has to submit
// something useful before the behaviour runs, and the only way to know what it
// submits is to do what the browser does: find the inputs somebody can see,
// type into them, and send every named input that is not disabled.

// inputTag is the start tag of one input, and inputAttr one attribute in it,
// with or without a value.
var (
	inputTag  = regexp.MustCompile(`<input\s[^>]*>`)
	inputAttr = regexp.MustCompile(`([a-z][a-z0-9-]*)(?:="([^"]*)")?`)
)

// input is one drawn input, read the way a browser reads its attributes: the
// first of two with the same name wins.
type input map[string]string

func (in input) has(name string) bool { _, ok := in[name]; return ok }

// inputs reads every input start tag in html, in document order.
func inputs(markup string) []input {
	var out []input
	for _, tag := range inputTag.FindAllString(markup, -1) {
		in := input{}
		for _, m := range inputAttr.FindAllStringSubmatch(strings.TrimPrefix(tag, "<input"), -1) {
			if !in.has(m[1]) {
				in[m[1]] = html.UnescapeString(m[2])
			}
		}
		out = append(out, in)
	}
	return out
}

// postWithoutScript is what a browser with no script running sends for markup
// once somebody has typed code into it: they type into the inputs they can see
// -- the whole code into one box, or a character into each of several -- and
// the form posts every named input that is not disabled, under its name and
// with its value, in document order. Nothing joins anything.
func postWithoutScript(markup, code string) map[string][]string {
	drawn := inputs(markup)
	var visible []input
	for _, in := range drawn {
		if in["type"] != "hidden" && !in.has("hidden") && !in.has("disabled") {
			visible = append(visible, in)
		}
	}
	if len(visible) == 1 {
		visible[0]["value"] = code
	} else {
		for at, in := range visible {
			if at < len(code) {
				in["value"] = code[at : at+1]
			}
		}
	}
	posted := map[string][]string{}
	for _, in := range drawn {
		if in["name"] == "" || in.has("disabled") {
			continue
		}
		posted[in["name"]] = append(posted[in["name"]], in["value"])
	}
	return posted
}
