package unit_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Every behaviour a component emits has an answer in the script the runtime
// serves.
//
// This is the half of the claim that reads the components: it walks the
// constants in the component package and requires each one to be registered in
// hesape's ui.js. The other half lives there and reads the script.
//
// Two halves because neither repository can hold both: this one declares the
// names and cannot import the view runtime, and that one serves the script and
// must not depend on a component library. What connects them is the string, so
// the string is what each side checks.
//
// The password box shipped without an answer from the day it was written. Every
// project that drew a sign-up form implemented the behaviour by hand, and the
// console reported the behaviour missing on every page. Found by the Corujão.ai
// team, twice.

// emittedBehaviour is the constant a component declares its behaviour name in.
var emittedBehaviour = regexp.MustCompile(`(?m)^const ([A-Za-z]+)Behavior = "([a-z-]+)"`)

// registeredBehaviour is a name ui.js answers to.
var registeredBehaviour = regexp.MustCompile(`ui\.define\(\'([a-z-]+)\'`)

func TestEveryBehaviourTheComponentsEmitHasAnAnswer(t *testing.T) {
	t.Parallel()

	names := map[string]string{}
	entries, err := os.ReadDir(filepath.Join("..", "..", "components"))
	if err != nil {
		t.Fatalf("reading the components: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		body, err := os.ReadFile(filepath.Join("..", "..", "components", name))
		if err != nil {
			t.Fatalf("reading %s: %v", name, err)
		}
		for _, match := range emittedBehaviour.FindAllStringSubmatch(string(body), -1) {
			names[match[2]] = name
		}
	}
	if len(names) == 0 {
		t.Fatal("no component declares a behaviour name, so this test proved nothing")
	}

	script := servedScript(t)
	answered := map[string]bool{}
	for _, match := range registeredBehaviour.FindAllStringSubmatch(script, -1) {
		answered[match[1]] = true
	}

	for behaviour, file := range names {
		if !answered[behaviour] {
			t.Errorf("%s emits the behaviour %q and the served ui.js registers nothing under that name: "+
				"the markup renders, nothing mounts, and every project that draws this component writes the behaviour by hand",
				file, behaviour)
		}
	}
}

// servedScript reads ui.js out of the hesape checkout beside this one, so what
// is checked is the file this build would serve.
func servedScript(t *testing.T) string {
	t.Helper()

	beside := filepath.Join("..", "..", "..", "hesape", "view", "assets", "ui.js")
	body, err := os.ReadFile(beside)
	if err != nil {
		t.Skipf("no hesape checkout beside this one: the served script cannot be read (%v)", err)
	}
	return string(body)
}
