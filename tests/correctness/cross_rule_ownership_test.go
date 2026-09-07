package correctness_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/ergonomy"
	"github.com/will2469/charites/internal/rules/ux"
)

// TestCrossRuleOwnership memverifikasi ortogonalitas dan hierarki yield antara empat aturan:
// 1. ux.number-input-identity-misuse
// 2. ux.number-input-missing-bounds
// 3. ergonomy.number-input-wheel-hazard
// 4. ux.multiline-input-misuse
func TestCrossRuleOwnership(t *testing.T) {
	wheelRule := ergonomy.NewNumberInputWheelHazardRule()
	identityRule := ux.NewNumberInputIdentityMisuseRule()
	boundsRule := ux.NewNumberInputMissingBoundsRule()
	multilineRule := ux.NewMultilineInputMisuseRule()

	cases := []struct {
		name                 string
		code                 string
		wantIdentityMisuse   int
		wantMissingBounds    int
		wantWheelHazard      int
		wantMultilineMisuse  int
		rationaleDescription string
	}{
		{
			name:                 "Postal Code type=number yields bounds to identity, wheel hazard remains independent",
			code:                 `<input type="number" name="postal_code" />`,
			wantIdentityMisuse:   1,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			wantMultilineMisuse:  0,
			rationaleDescription: "Identity field suppresses missing-bounds; wheel hazard is an orthogonal interaction hazard",
		},
		{
			name:                 "Civil NIK type=number yields bounds to identity, wheel hazard remains independent",
			code:                 `<input type="number" name="nik" />`,
			wantIdentityMisuse:   1,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			wantMultilineMisuse:  0,
			rationaleDescription: "Local identity token NIK yields bounds to identity-misuse",
		},
		{
			name:                 "Mathematical quantity type=number lacks both bounds and wheel protection",
			code:                 `<input type="number" name="quantity" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    1,
			wantWheelHazard:      1,
			wantMultilineMisuse:  0,
			rationaleDescription: "Legitimate numeric input lacks domain lower bound and wheel protection",
		},
		{
			name:                 "Mathematical quantity type=number with min=0 still lacks wheel protection",
			code:                 `<input type="number" name="quantity" min="0" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			wantMultilineMisuse:  0,
			rationaleDescription: "Lower bound satisfies bounds rule, but wheel hazard remains unhandled",
		},
		{
			name:                 "Fully compliant numeric input with bound and wheel blur",
			code:                 `<input type="number" name="quantity" min="0" onWheel={(e) => e.currentTarget.blur()} />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "All invariants satisfied, 0 diagnostics across all rules",
		},
		{
			name:                 "Compliant text numeric identity field",
			code:                 `<input type="text" name="postal_code" inputMode="numeric" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Using type=text with inputMode=numeric is the recommended pattern for identity fields",
		},
		{
			name:                 "Standard text field",
			code:                 `<input type="text" name="nik" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Plain text input triggers no number-specific rules",
		},
		{
			name:                 "Multiline input misuse with text notes",
			code:                 `<input type="text" name="notes" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  1,
			rationaleDescription: "Open-ended notes field on text input triggers ux.multiline-input-misuse",
		},
		{
			name:                 "Number input named notes without safeguards",
			code:                 `<input type="number" name="notes" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    1,
			wantWheelHazard:      1,
			wantMultilineMisuse:  0,
			rationaleDescription: "Type=number is InputTypeOther, so multiline rule never fires; missing bounds and wheel fire",
		},
		{
			name:                 "Numeric notes count input with safeguards",
			code:                 `<input type="number" name="notes_count" min="0" onWheel={(e) => e.currentTarget.blur()} />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Count qualifier clarifies quantity, and all numeric invariants are satisfied",
		},
		{
			name:                 "Text input for note title",
			code:                 `<input name="note_title" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Title qualifier marks single-line intent",
		},
		{
			name:                 "Proper textarea component",
			code:                 `<textarea name="notes" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Textarea element is the compliant control for open-ended notes",
		},
		{
			name:                 "Multi-channel decoupling test: Unknown identity name with multiline placeholder",
			code:                 `<input name="user_name" placeholder="Catatan tambahan" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  1,
			rationaleDescription: "Decoupled architecture: name evaluates to Unknown identity, placeholder resolves to Multiline intent",
		},
		{
			name:                 "Qualifier override: Multiline name overridden by single-line placeholder",
			code:                 `<input name="description" placeholder="Short description" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  0,
			rationaleDescription: "Explicit single-line qualifier 'short' wins over 'description'",
		},
		{
			name:                 "First-class aria-label channel multiline intent",
			code:                 `<input name="foo" aria-label="Catatan tambahan" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			wantMultilineMisuse:  1,
			rationaleDescription: "aria-label provides first-class content intent evidence when name is arbitrary",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `export function Form() { return (` + tc.code + `); }`
			root, err := tsx.Extract([]byte(src))
			if err != nil {
				t.Fatalf("failed to parse TSX snippet: %v", err)
			}

			var gotIdentity, gotBounds, gotWheel, gotMultiline int
			for node := range root.Walk() {
				gotIdentity += len(identityRule.Evaluate(node))
				gotBounds += len(boundsRule.Evaluate(node))
				gotWheel += len(wheelRule.Evaluate(node))
				gotMultiline += len(multilineRule.Evaluate(node))
			}

			if gotIdentity != tc.wantIdentityMisuse {
				t.Errorf("[%s] identity-misuse count = %d, want %d (%s)", tc.name, gotIdentity, tc.wantIdentityMisuse, tc.rationaleDescription)
			}
			if gotBounds != tc.wantMissingBounds {
				t.Errorf("[%s] missing-bounds count = %d, want %d (%s)", tc.name, gotBounds, tc.wantMissingBounds, tc.rationaleDescription)
			}
			if gotWheel != tc.wantWheelHazard {
				t.Errorf("[%s] wheel-hazard count = %d, want %d (%s)", tc.name, gotWheel, tc.wantWheelHazard, tc.rationaleDescription)
			}
			if gotMultiline != tc.wantMultilineMisuse {
				t.Errorf("[%s] multiline-misuse count = %d, want %d (%s)", tc.name, gotMultiline, tc.wantMultilineMisuse, tc.rationaleDescription)
			}
		})
	}
}
