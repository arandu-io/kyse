package unit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The shade of a duotone icon is a path with an opacity attribute. A component
// that hides it with a selector on that path (`[&_path[opacity]]:hidden`)
// draws a regular icon from a duotone set, and the family draws every icon in
// both layers. The sources are read rather than a render, so a component that
// only shows the icon in a state a test does not reach is held to it too.
func TestNoComponentHidesTheShadeOfAnIcon(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "components", "*.go"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no component sources found: %v", err)
	}
	for _, file := range files {
		body, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		for _, banned := range []string{"path[opacity]", "[&_path", "[&>path"} {
			if strings.Contains(string(body), banned) {
				t.Errorf("%s: %q hides or restyles the shade of a duotone icon", filepath.Base(file), banned)
			}
		}
	}
}
