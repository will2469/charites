package correctness_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/ergonomy"
	"github.com/will2469/charites/internal/rules/ux"
)

// TestCrossRuleOwnership memverifikasi ortogonalitas dan hierarki yield antara tiga aturan:
// 1. ux.number-input-identity-misuse
// 2. ux.number-input-missing-bounds
// 3. ergonomy.number-input-wheel-hazard
func TestCrossRuleOwnership(t *testing.T) {
	wheelRule := ergonomy.NewNumberInputWheelHazardRule()
	identityRule := ux.NewNumberInputIdentityMisuseRule()
	boundsRule := ux.NewNumberInputMissingBoundsRule()

	cases := []struct {
		name                 string
		code                 string
		wantIdentityMisuse   int
		wantMissingBounds    int
		wantWheelHazard      int
		rationaleDescription string
	}{
		{
			name:                 "Postal Code type=number yields bounds to identity, wheel hazard remains independent",
			code:                 `<input type="number" name="postal_code" />`,
			wantIdentityMisuse:   1,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			rationaleDescription: "Identity field suppresses missing-bounds; wheel hazard is an orthogonal interaction hazard",
		},
		{
			name:                 "Civil NIK type=number yields bounds to identity, wheel hazard remains independent",
			code:                 `<input type="number" name="nik" />`,
			wantIdentityMisuse:   1,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			rationaleDescription: "Local identity token NIK yields bounds to identity-misuse",
		},
		{
			name:                 "Mathematical quantity type=number lacks both bounds and wheel protection",
			code:                 `<input type="number" name="quantity" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    1,
			wantWheelHazard:      1,
			rationaleDescription: "Legitimate numeric input lacks domain lower bound and wheel protection",
		},
		{
			name:                 "Mathematical quantity type=number with min=0 still lacks wheel protection",
			code:                 `<input type="number" name="quantity" min="0" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      1,
			rationaleDescription: "Lower bound satisfies bounds rule, but wheel hazard remains unhandled",
		},
		{
			name:                 "Fully compliant numeric input with bound and wheel blur",
			code:                 `<input type="number" name="quantity" min="0" onWheel={(e) => e.currentTarget.blur()} />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			rationaleDescription: "All invariants satisfied, 0 diagnostics across all three rules",
		},
		{
			name:                 "Compliant text numeric identity field",
			code:                 `<input type="text" name="postal_code" inputMode="numeric" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			rationaleDescription: "Using type=text with inputMode=numeric is the recommended pattern for identity fields",
		},
		{
			name:                 "Standard text field",
			code:                 `<input type="text" name="nik" />`,
			wantIdentityMisuse:   0,
			wantMissingBounds:    0,
			wantWheelHazard:      0,
			rationaleDescription: "Plain text input triggers no number-specific rules",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `export function Form() { return (` + tc.code + `); }`
			root, err := tsx.Extract([]byte(src))
			if err != nil {
				t.Fatalf("failed to parse TSX snippet: %v", err)
			}

			var gotIdentity, gotBounds, gotWheel int
			for node := range root.Walk() {
				gotIdentity += len(identityRule.Evaluate(node))
				gotBounds += len(boundsRule.Evaluate(node))
				gotWheel += len(wheelRule.Evaluate(node))
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
		})
	}
}
