package unit_test

import (
	"github.com/arandu-io/hesape/view"
	"net/http"
	"net/http/httptest"
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
// The component suite compares its names against the runtime linked by go.mod.
// Reading a neighboring checkout would test unrelated bytes or skip on a clean
// consumer. The runtime's own suite independently checks its implementation.
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

// servedScript reads the bytes linked into this build, not a sibling checkout.
// A module release has no neighboring source tree and must never skip this proof.
func servedScript(t *testing.T) string {
	t.Helper()
	path := view.Asset("ui.js")
	response := httptest.NewRecorder()
	view.Handler(response, httptest.NewRequest(http.MethodGet, path, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("linked ui.js status = %d", response.Code)
	}
	if !strings.Contains(path, view.AssetHash(response.Body.Bytes())) {
		t.Fatal("linked runtime bytes do not match their content-addressed URL")
	}
	return response.Body.String()
}
