package unit

import (
	"github.com/arandu-io/hesape/validation"
	"github.com/arandu-io/kyse/components"
	"strings"
	"testing"
)

func TestPasswordLabelsDoNotChangeTheNativePolicy(t *testing.T) {
	policy := validation.PasswordMin(12).Max(128)
	props := components.PasswordProps{Name: "password", Policy: policy,
		RequirementLabels: map[string]string{"min": "Pelo menos {value} caracteres", "max": "Até {value} caracteres"},
		StrengthLabel:     "{met} de {total} requisitos atendidos", MetLabel: "Atendido:", UnmetLabel: "Pendente:"}
	if props.StrengthText() != "0 de 1 requisitos atendidos" {
		t.Fatal("initial summary does not use the label template")
	}
	requirements := props.Requirements()
	if len(requirements) != 2 || requirements[0].Text != "Pelo menos 12 caracteres" || requirements[1].Text != "Até 128 caracteres" {
		t.Fatal("localized requirements lost their policy values")
	}
	if props.AppliedPolicy() != policy || policy.Passes("password", "short") {
		t.Fatal("label overrides changed validation")
	}
	html := string(components.Password(props))
	for _, want := range []string{`data-password-summary="{met} de {total} requisitos atendidos"`, `data-password-met="Atendido:"`, `data-password-status`, `aria-controls="password-requirements-panel"`} {
		if !strings.Contains(html, want) {
			t.Errorf("missing native bridge %s", want)
		}
	}
}

func TestPasswordLabelOverridesAreEscaped(t *testing.T) {
	hostile := `"><script>alert(1)</script>`
	html := string(components.Password(components.PasswordProps{Name: "password", StrengthLabel: hostile, MetLabel: hostile, UnmetLabel: hostile, RequirementLabels: map[string]string{"min": hostile}}))
	if strings.Contains(html, "<script>alert(1)</script>") {
		t.Fatal("password labels introduced executable markup")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Fatal("label contents were dropped rather than escaped")
	}
}
