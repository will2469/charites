package ux_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/ux"
)

func TestSpacingRhythmDriftRule(t *testing.T) {
	rule := ux.NewSpacingRhythmDriftRule()

	cases := []struct {
		name      string
		code      string
		wantDiags int
	}{
		{
			name: "Case 1: Uniform rhythm (4, 4, 4, 4)",
			code: `
export function UniformRhythm() {
  return (
    <div>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-4" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 2: Isolated outlier (4, 4, 7, 4)",
			code: `
export function MarginOutlier() {
  return (
    <div>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-7" />
      <Field className="mb-4" />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 3: Homogeneous group outlier (3, 3, 3, 6)",
			code: `
export function RepeatedGroupOutlier() {
  return (
    <div className="flex flex-col">
      <Field className="mb-3" />
      <Field className="mb-3" />
      <Field className="mb-3" />
      <Field className="mb-6" />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 4: Majority three items (4, 4, 5)",
			code: `
export function MajorityThreeItems() {
  return (
    <div>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-5" />
    </div>
  );
}`,
			wantDiags: 1,
		},
		{
			name: "Case 5: Insufficient sample size (3, 6)",
			code: `
export function InsufficientSample() {
  return (
    <div>
      <Field className="mb-3" />
      <Field className="mb-6" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 6: Section with field groups (RoleGroup classification)",
			code: `
export function SectionWithFieldGroups() {
  return (
    <section className="flex flex-col gap-4">
      <Field />
      <Field />
      <Field />
    </section>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 7: Distinct semantic roles (SectionHeader vs FieldGroup vs ActionFooter)",
			code: `
export function DistinctSemanticRoles() {
  return (
    <div className="flex flex-col">
      <SectionHeader className="mb-3" />
      <FieldGroup className="mb-3" />
      <ActionFooter className="mb-6" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 8: Responsive spacing isolation (gap-4 vs md:gap-6)",
			code: `
export function ResponsiveSpacing() {
  return (
    <div className="flex flex-col gap-4 md:gap-6">
      <Field />
      <Field />
      <Field />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 9: Dynamic spacing immunity (gap-${spacing})",
			code: `
export function DynamicSpacing({ spacing }: { spacing: number }) {
  return (
    <div className="flex flex-col">
      <Field />
      <Field />
      <Field />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 10: Different layout groups (Card vs Modal)",
			code: `
export function DifferentLayoutGroups() {
  return (
    <div>
      <Card className="flex flex-col gap-4">
        <Field />
        <Field />
      </Card>
      <Modal className="flex flex-col gap-8">
        <Field />
        <Field />
      </Modal>
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 11: Different axes (gap-y-4 vs gap-x-6)",
			code: `
export function DifferentAxes() {
  return (
    <div className="flex gap-x-6 gap-y-4">
      <Field />
      <Field />
      <Field />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 12: Two equally common values (4, 4, 6, 6 tie)",
			code: `
export function TieSpacing() {
  return (
    <div>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-6" />
      <Field className="mb-6" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 13: Multiple outliers (4, 4, 7, 9, 4)",
			code: `
export function MultipleOutliers() {
  return (
    <div>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-7" />
      <Field className="mb-9" />
      <Field className="mb-4" />
    </div>
  );
}`,
			wantDiags: 2,
		},
		{
			name: "Case 14: Margin Guard Proof 1 (Competing parent gap-4)",
			code: `
export function MarginCompetingParentGap() {
  return (
    <div className="flex flex-col gap-4">
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-7" />
      <Field className="mb-4" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 15: Margin Guard Proof 2 (Competing parent space-y-4)",
			code: `
export function MarginCompetingParentSpace() {
  return (
    <div className="space-y-4">
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-7" />
      <Field className="mb-4" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 16: Margin Guard Proof 3 (Non-adjacent margin)",
			code: `
export function MarginNonAdjacent() {
  return (
    <div>
      <Field className="mb-4" />
      <div className="border-b" />
      <Field className="mb-4" />
      <Field className="mb-7" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 17: Margin Guard Proof 4 (Mixed semantic siblings)",
			code: `
export function MarginMixedSemantics() {
  return (
    <div>
      <Heading className="mb-4" />
      <Field className="mb-4" />
      <Button className="mb-7" />
    </div>
  );
}`,
			wantDiags: 0,
		},
		{
			name: "Case 18: Margin Guard Proof 5 (No layout proof - span wrapper)",
			code: `
export function MarginNoLayoutProof() {
  return (
    <span>
      <Field className="mb-4" />
      <Field className="mb-4" />
      <Field className="mb-7" />
      <Field className="mb-4" />
    </span>
  );
}`,
			wantDiags: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root, err := tsx.Extract([]byte(tc.code))
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
