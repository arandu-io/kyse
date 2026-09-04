package unit

import (
	"regexp"
	"strings"
	"testing"

	"github.com/arandu-io/hesape/validation"
	"github.com/arandu-io/kyse/components"
)

// TestThePolicyReachesTheClientAsTheValidatorDeclaredIt is the property the
// whole component exists for.
//
// The lines a person reads and the rule the server rejects with have to be one
// declaration. So the props handed to the behaviour are compared against the
// policy's own AppliedRules rather than against a list written here: a second
// list in this file would be the very thing the component was built to avoid,
// and it would agree on the day it was written.
func TestThePolicyReachesTheClientAsTheValidatorDeclaredIt(t *testing.T) {
	policy := validation.NewPassword(12).Max(64).MixedCase().Letters().Numbers().Symbols()

	html := string(components.Password(components.PasswordProps{
		Name: "password", Label: "Password", Policy: policy,
	}))

	for key, value := range policy.AppliedRules() {
		if !strings.Contains(html, "&#34;"+key+"&#34;") {
			t.Errorf("the policy declares %q and the behaviour is not told about it:\n%s", key, html)
		}
		_ = value
	}
	if !strings.Contains(html, `data-kyse-behavior="password"`) {
		t.Errorf("the field names no behaviour, so nothing ever opens the panel:\n%s", html)
	}
}

// TestEveryRequirementLineNamesARuleThePolicyAsksFor holds the other direction.
//
// A line drawn for something the validator does not check is a promise the
// server never keeps, and it is the failure that looks like success: somebody
// satisfies the list and is rejected anyway, or satisfies it and is let through
// on a rule that was never enforced.
func TestEveryRequirementLineNamesARuleThePolicyAsksFor(t *testing.T) {
	policy := validation.NewPassword(10).Numbers().Symbols()
	rules := policy.AppliedRules()

	html := string(components.Password(components.PasswordProps{
		Name: "password", Label: "Password", Policy: policy,
	}))

	drawn := regexp.MustCompile(`data-requirement="([a-zA-Z]+)"`).FindAllStringSubmatch(html, -1)
	if len(drawn) == 0 {
		t.Fatalf("no requirement was drawn at all:\n%s", html)
	}

	for _, found := range drawn {
		key := found[1]
		if key == "custom" {
			// The merged rules are one line, and the policy holds them under a
			// name of its own rather than as a flag.
			if merged, _ := rules["customRules"].([]string); len(merged) == 0 {
				t.Errorf("a line was drawn for merged rules the policy has none of")
			}
			continue
		}
		switch value := rules[key].(type) {
		case int:
			if value <= 0 {
				t.Errorf("the line %q is drawn and the policy asks for nothing there", key)
			}
		case bool:
			if !value {
				t.Errorf("the line %q is drawn and the policy does not ask for it", key)
			}
		default:
			t.Errorf("the line %q names nothing the policy declares", key)
		}
	}

	// The ones this policy does not ask for stay off the list.
	for _, absent := range []string{"max", "mixedCase", "letters", "uncompromised"} {
		if strings.Contains(html, `data-requirement="`+absent+`"`) {
			t.Errorf("the policy does not ask for %q and a line was drawn for it:\n%s", absent, html)
		}
	}
}

// TestAFieldWithNoPolicyTakesTheRegisteredDefault is the case a sign-up form
// writes: no policy on the field, and the one the application registered.
func TestAFieldWithNoPolicyTakesTheRegisteredDefault(t *testing.T) {
	html := string(components.Password(components.PasswordProps{Name: "password", Label: "Password"}))

	if !strings.Contains(html, "At least 8 characters") {
		t.Errorf("a field given no policy did not fall back to the default:\n%s", html)
	}
	if !strings.Contains(html, "&#34;min&#34;:8") {
		t.Errorf("the default policy did not reach the behaviour:\n%s", html)
	}
}

// TestThePanelIsDrawnClosed is what keeps the component honest about who opens
// it.
//
// The panel is the client's to open, on the first keystroke, and a panel that
// arrived open would be a component that had decided somebody is typing before
// anybody has. Every line arrives unmet for the same reason: the box is empty.
func TestThePanelIsDrawnClosed(t *testing.T) {
	html := string(components.Password(components.PasswordProps{
		Name: "password", Label: "Password",
		Policy: validation.NewPassword(12).Letters().Numbers(),
	}))

	panel := regexp.MustCompile(`(?s)<section[^>]*data-part="panel".*?>`).FindString(html)
	if panel == "" {
		t.Fatalf("the panel was not drawn at all:\n%s", html)
	}
	if !strings.Contains(panel, "hidden") {
		t.Errorf("the panel arrived open:\n%s", panel)
	}
	if strings.Contains(html, `data-met="true"`) {
		t.Errorf("a requirement arrived met on an empty box:\n%s", html)
	}
	if !strings.Contains(html, `class="w-0"`) {
		t.Errorf("the strength bar arrived filled on an empty box:\n%s", html)
	}
}

// TestThePasswordBoxCarriesNoValue.
//
// The flash keeps the message for a password and drops what was typed, so there
// is nothing to redraw. A value here would be a secret written into the markup
// of a page that is cached and scrolled past over somebody's shoulder.
func TestThePasswordBoxCarriesNoValue(t *testing.T) {
	html := string(components.Password(components.PasswordProps{
		Name: "password", Label: "Password",
		Page: page{
			errs: map[string]string{"password": "is too short"},
			old:  map[string]string{"password": "hunter2"},
		},
	}))

	if strings.Contains(html, "hunter2") || strings.Contains(html, "value=") {
		t.Errorf("the password box came back carrying a value:\n%s", html)
	}
	if !strings.Contains(html, "is too short") {
		t.Errorf("the box does not say why it was rejected:\n%s", html)
	}
	if !strings.Contains(html, `aria-invalid="true"`) {
		t.Errorf("the rejected box is not marked invalid:\n%s", html)
	}
}

// TestTheBoxIsDescribedByItsChecklist.
//
// A list that is only drawn is a list the people who most need the rules never
// hear. The summary is deliberately not named: it is a live region, and naming
// it as well would read the count out on every focus and again on every
// keystroke that moved it.
func TestTheBoxIsDescribedByItsChecklist(t *testing.T) {
	for _, c := range []struct {
		name  string
		props components.PasswordProps
		want  string
	}{
		{"bare", components.PasswordProps{Name: "password"}, `aria-describedby="password-requirements"`},
		{"hint", components.PasswordProps{Name: "password", Hint: "Long is better."},
			`aria-describedby="password-hint password-requirements"`},
		{"message", components.PasswordProps{Name: "password",
			Page: page{errs: map[string]string{"password": "is required"}}},
			`aria-describedby="password-error password-requirements"`},
	} {
		t.Run(c.name, func(t *testing.T) {
			html := string(components.Password(c.props))
			if !strings.Contains(html, c.want) {
				t.Errorf("the box is not described by %s:\n%s", c.want, html)
			}
			described := regexp.MustCompile(`aria-describedby="([^"]*)"`).FindStringSubmatch(html)
			if described != nil && strings.Contains(described[1], "password-strength") {
				t.Errorf("the live summary is named as a description as well:\n%s", described[1])
			}
		})
	}
}

// TestTheRevealControlSaysWhichStateItIsIn.
//
// A control that swaps the box to text and back is two states, and the one it
// is in has to be readable without seeing the characters change. Both words
// travel with it so that whoever flips them does not have to carry a second
// copy of the wording.
func TestTheRevealControlSaysWhichStateItIsIn(t *testing.T) {
	html := string(components.Password(components.PasswordProps{Name: "password", Label: "Password"}))

	for _, want := range []string{
		`aria-pressed="false"`,
		`data-reveal-hidden="Show"`,
		`data-reveal-shown="Hide"`,
		`aria-controls="password"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the reveal control does not carry %s:\n%s", want, html)
		}
	}

	translated := string(components.Password(components.PasswordProps{
		Name: "password", ShowLabel: "Mostrar", HideLabel: "Ocultar",
	}))
	if !strings.Contains(translated, `data-reveal-shown="Ocultar"`) {
		t.Errorf("the pair of words is not translated together:\n%s", translated)
	}
}

// TestARevealedBoxIsStillNotProse.
//
// Revealing is a swap to type="text", and a text box is prose to a browser: it
// is offered what was saved for other text boxes, its content goes to a spell
// checker, and a phone capitalises it. These are written whatever the type is,
// because the component cannot know when the swap happens.
func TestARevealedBoxIsStillNotProse(t *testing.T) {
	html := string(components.Password(components.PasswordProps{Name: "password"}))

	for _, want := range []string{
		`type="password"`,
		`autocomplete="new-password"`,
		`spellcheck="false"`,
		`autocapitalize="none"`,
		`autocorrect="off"`,
	} {
		if !strings.Contains(html, want) {
			t.Errorf("the box does not carry %s:\n%s", want, html)
		}
	}

	signIn := string(components.Password(components.PasswordProps{Name: "password", Autocomplete: "current-password"}))
	if !strings.Contains(signIn, `autocomplete="current-password"`) {
		t.Errorf("a sign-in box could not name the password it recalls:\n%s", signIn)
	}
}

// TestACallerKeepsTheBehaviourTheyNamed.
//
// The default is a default and not a seizure. A field driven by an application's
// own behaviour would otherwise be impossible to write, because the component
// would overwrite the name on every render.
func TestACallerKeepsTheBehaviourTheyNamed(t *testing.T) {
	html := string(components.Password(components.PasswordProps{
		Name: "password",
		ComponentProps: components.ComponentProps{
			Behavior: components.Behavior{Name: "our-password", Props: map[string]any{"strict": true}},
		},
	}))

	if !strings.Contains(html, `data-kyse-behavior="our-password"`) {
		t.Errorf("the caller's behaviour was overwritten:\n%s", html)
	}
	if !strings.Contains(html, `data-kyse-props="{&#34;strict&#34;:true}"`) {
		t.Errorf("the caller's props were overwritten:\n%s", html)
	}
}
