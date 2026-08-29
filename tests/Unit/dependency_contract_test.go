package unit

import (
	"encoding/json"
	"os/exec"
	"path/filepath"
	"slices"
	"testing"
)

// TestTheNativeViewDependencyBudget keeps ownership visible in the module
// graph. Generated components compile against hesape/view; Framework remains
// direct for the application adapter and asset registry. A third direct
// dependency is a new shipping boundary and must fail here rather than enter
// every application that imports a component.
func TestTheNativeViewDependencyBudget(t *testing.T) {
	t.Parallel()

	root, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("finding the module root: %v", err)
	}
	command := exec.Command("go", "mod", "edit", "-json")
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		t.Fatalf("reading go.mod through the Go tool: %v", err)
	}

	var module struct {
		Require []struct {
			Path     string
			Indirect bool
		}
	}
	if err := json.Unmarshal(output, &module); err != nil {
		t.Fatalf("decoding go mod edit -json: %v", err)
	}

	direct := make([]string, 0, len(module.Require))
	for _, requirement := range module.Require {
		if !requirement.Indirect {
			direct = append(direct, requirement.Path)
		}
	}
	slices.Sort(direct)
	want := []string{
		"github.com/arandu-io/framework",
		"github.com/arandu-io/hesape",
	}
	if !slices.Equal(direct, want) {
		t.Fatalf("direct dependencies = %v, want exactly %v", direct, want)
	}
}
