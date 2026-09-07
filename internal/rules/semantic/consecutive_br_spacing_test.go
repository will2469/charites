package semantic_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/semantic"
)

func TestConsecutiveBRSpacingRule(t *testing.T) {
	rule := semantic.NewConsecutiveBRSpacingRule()

	cases := []struct {
		name      string
		code      string
		wantDiags int
	}{
		{
			name: "Single br is permitted",
			code: `<p>
  Line 1
  <br />
  Line 2
</p>`,
			wantDiags: 0,
		},
		{
			name: "Double br produces exactly 1 diagnostic",
			code: `<div>
  <p>Text 1</p>
  <br />
  <br />
  <p>Text 2</p>
</div>`,
			wantDiags: 1,
		},
		{
			name: "Triple br produces exactly 1 diagnostic",
			code: `<div>
  <p>Text 1</p>
  <br />
  <br />
  <br />
  <p>Text 2</p>
</div>`,
			wantDiags: 1,
		},
		{
			name: "Quadruple br produces exactly 1 diagnostic",
			code: `<div>
  <p>Text 1</p>
  <br />
  <br />
  <br />
  <br />
  <p>Text 2</p>
</div>`,
			wantDiags: 1,
		},
		{
			name: "Comment between br does not break run",
			code: `<div>
  <p>Text 1</p>
  <br />
  {/* spacing */}
  <br />
  <p>Text 2</p>
</div>`,
			wantDiags: 1,
		},
		{
			name: "Whitespace between br does not break run",
			code: `<div>
  <p>Text 1</p>
  <br />

  <br />
  <p>Text 2</p>
</div>`,
			wantDiags: 1,
		},
		{
			name: "Meaningful element between br resets run",
			code: `<div>
  <br />
  <p>Paragraph</p>
  <br />
</div>`,
			wantDiags: 0,
		},
		{
			name: "Non-flattening: br in different child parents does not trigger",
			code: `<div>
  <p><br /></p>
  <p><br /></p>
</div>`,
			wantDiags: 0,
		},
		{
			name: "Prefix containment bait: Breadcrumb is not br",
			code: `<div>
  <Breadcrumb />
  <Breadcrumb />
</div>`,
			wantDiags: 0,
		},
		{
			name: "Prefix containment bait: BrandedButton is not br",
			code: `<div>
  <BrandedButton />
  <BrandedButton />
</div>`,
			wantDiags: 0,
		},
		{
			name: "State machine run-reset: two distinct runs produce exactly 2 diagnostics",
			code: `<div>
  <br />
  <br />
  <p>Middle Element</p>
  <br />
  <br />
</div>`,
			wantDiags: 2,
		},
		{
			name: "Astrades dialog specimen inside Fragment",
			code: `<>
  Apakah Anda yakin?
  <br />
  <br />
  <strong>PERINGATAN</strong>
</>`,
			wantDiags: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `export function TestComponent() { return (` + tc.code + `); }`
			root, err := tsx.Extract([]byte(src))
			if err != nil {
				t.Fatalf("failed to parse TSX: %v", err)
			}

			var diagsCount int
			for node := range root.Walk() {
				findings := rule.Evaluate(node)
				diagsCount += len(findings)
			}

			if diagsCount != tc.wantDiags {
				t.Errorf("[%s] got %d diagnostics, want %d", tc.name, diagsCount, tc.wantDiags)
			}
		})
	}
}
