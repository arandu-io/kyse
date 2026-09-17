package unit

import (
	"github.com/arandu-io/hesape/validation"
	"github.com/arandu-io/kyse/components"
	"regexp"
	"strings"
	"testing"
)

func TestPasswordUpperBoundIsAConstraintRatherThanStrength(t *testing.T) {
	policy := validation.PasswordMin(12).Max(128).Letters().MixedCase().Numbers().Symbols()
	props := components.PasswordProps{Name: "password", Policy: policy}
	if props.RequirementTotal() != 5 {
		t.Fatal("maximum length awarded a strength goal")
	}
	html := string(components.Password(props))
	row := regexp.MustCompile(`(?s)<li\b[^>]*data-requirement="max"[^>]*>`).FindString(html)
	if row == "" || !strings.Contains(row, "hidden") {
		t.Fatal("technical limit was not hidden initially")
	}
	if !strings.Contains(html, `aria-valuemax="5"`) {
		t.Fatal("bar scale drifted from the policy")
	}
	if props.AppliedPolicy() != policy || props.AppliedPolicy().Passes("password", strings.Repeat("a", 128)+"Z1!") {
		t.Fatal("hiding the goal disabled its validation")
	}
}

func TestPasswordPanelUsesTheNativeCompactLayout(t *testing.T) {
	html := string(components.Password(components.PasswordProps{Name: "password"}))
	for _, marker := range []string{"bottom-full!", "top-auto!", "h-1.5!", "flex-wrap", "border-dashed", "self-start", "sr-only"} {
		if !strings.Contains(html, marker) {
			t.Fatalf("missing compact native part %s", marker)
		}
	}
	if strings.Contains(html, "style=") || strings.Contains(html, "<script") {
		t.Fatal("presentation introduced inline execution")
	}
}
