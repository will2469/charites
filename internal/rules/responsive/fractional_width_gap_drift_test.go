package responsive_test

import (
	"testing"

	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/rules/responsive"
)

func TestFractionalWidthGapDriftRule_AcceptanceMatrix(t *testing.T) {
	rule := responsive.NewFractionalWidthGapDriftRule()

	cases := []struct {
		name      string
		code      string
		wantDiags int
	}{
		{
			name:      "1. w-1/2 standalone is not inherently a defect",
			code:      `<div className="w-1/2">Content</div>`,
			wantDiags: 0,
		},
		{
			name:      "2. w-2/3 standalone is not inherently a defect",
			code:      `<div className="w-2/3">Content</div>`,
			wantDiags: 0,
		},
		{
			name:      "3. w-full md:w-2/3 compliant mobile-first progression",
			code:      `<div className="w-full md:w-2/3">Content</div>`,
			wantDiags: 0,
		},
		{
			name: "4. flex gap-4 w-1/2 + w-1/2 default nowrap and shrink-1 absorbs gap safely",
			code: `<div className="flex gap-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "5. flex flex-wrap gap-4 w-1/2 + w-1/2 wrap causes unexpected line breaking",
			code: `<div className="flex flex-wrap gap-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "6. flex gap-4 w-1/2 shrink-0 + w-1/2 shrink-0 rigid items cause container blowout",
			code: `<div className="flex gap-4">
				<div className="w-1/2 shrink-0">A</div>
				<div className="w-1/2 shrink-0">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "7. flex gap-4 w-1/2 shrink-0 + w-1/2 one rigid item leaves unresolvable gap surplus",
			code: `<div className="flex gap-4">
				<div className="w-1/2 shrink-0">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "8. flex gap-4 w-1/2 + w-1/2 shrink-0 second rigid item leaves unresolvable gap surplus",
			code: `<div className="flex gap-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2 shrink-0">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "9. flex gap-4 w-2/3 + w-1/3 default nowrap and shrink-1 absorbs gap safely",
			code: `<div className="flex gap-4">
				<div className="w-2/3">A</div>
				<div className="w-1/3">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "10. flex gap-0 w-1/2 + w-1/2 zero gap has zero spacing drift",
			code: `<div className="flex gap-0">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "11. flex gap-y-4 w-1/2 + w-1/2 gap-y is vertical only, row spacing is zero",
			code: `<div className="flex gap-y-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "12. flex flex-col gap-4 w-1/2 + w-1/2 items stack vertically in a column",
			code: `<div className="flex flex-col gap-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "13. flex-col md:flex-row md:flex-wrap md:gap-4 clean column on mobile, wraps at md tier",
			code: `<div className="flex flex-col md:flex-row md:flex-wrap md:gap-4">
				<div className="w-full md:w-2/3">A</div>
				<div className="w-full md:w-1/3">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "14. flex-col md:flex-row md:gap-4 with md:shrink-0 clean column on mobile, blows out at md tier",
			code: `<div className="flex flex-col md:flex-row md:gap-4">
				<div className="w-full md:w-2/3 md:shrink-0">A</div>
				<div className="w-full md:w-1/3">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "15. flex gap-4 flex-1 fluid distribution distributes space cleanly",
			code: `<div className="flex gap-4">
				<div className="flex-1">A</div>
				<div className="flex-1">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "16. flex gap-4 flex-1 basis-0 canonical safe flex pattern",
			code: `<div className="flex gap-4">
				<div className="flex-1 basis-0">A</div>
				<div className="flex-1 basis-0">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "17. carousel intent with overflow-x-auto and shrink-0 peek items",
			code: `<div className="flex overflow-x-auto gap-4">
				<div className="w-4/5 shrink-0">Slide 1</div>
				<div className="w-4/5 shrink-0">Slide 2</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "18. CSS Grid alternative automatically accounts for gap",
			code: `<div className="grid grid-cols-1 md:grid-cols-3 gap-4">
				<div className="md:col-span-2">A</div>
				<div>B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name:      "19. w-[60%] standalone arbitrary percentage is not a defect",
			code:      `<div className="w-[60%]">Content</div>`,
			wantDiags: 0,
		},
		{
			name: "20. flex flex-wrap gap-4 w-[60%] + w-[40%] arbitrary percentages wrap unexpectedly",
			code: `<div className="flex flex-wrap gap-4">
				<div className="w-[60%]">A</div>
				<div className="w-[40%]">B</div>
			</div>`,
			wantDiags: 1,
		},
		{
			name: "21. flex gap-4 w-1/4 + w-1/4 fractional sum 50% < 100% fits comfortably",
			code: `<div className="flex gap-4">
				<div className="w-1/4">A</div>
				<div className="w-1/4">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "22. flex gap-4 w-1/4 shrink-0 + w-1/4 fractional sum 50% < 100% fits even with shrink-0",
			code: `<div className="flex gap-4">
				<div className="w-1/4 shrink-0">A</div>
				<div className="w-1/4">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "23. flex gap-4 gap-x-0 w-1/2 + w-1/2 gap-x-0 overrides horizontal gap to zero",
			code: `<div className="flex gap-4 gap-x-0">
				<div className="w-1/2 shrink-0">A</div>
				<div className="w-1/2 shrink-0">B</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "24. Non-fractional child with shrink-0 does not trigger hasNoShrink on fractional allocation",
			code: `<div className="flex gap-4">
				<div className="w-1/2">A</div>
				<div className="w-1/2">B</div>
				<div className="w-auto shrink-0">C</div>
			</div>`,
			wantDiags: 0,
		},
		{
			name: "25. Display override md:block deactivates flex on md tier",
			code: `<div className="flex md:block md:gap-4">
				<div className="w-full md:w-1/2 md:shrink-0">A</div>
				<div className="w-full md:w-1/2 md:shrink-0">B</div>
			</div>`,
			wantDiags: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := `export function Layout() { return (` + tc.code + `); }`
			root, err := tsx.Extract([]byte(src))
			if err != nil {
				t.Fatalf("failed to parse TSX: %v", err)
			}

			var totalDiags int
			for node := range root.Walk() {
				totalDiags += len(rule.Evaluate(node))
			}

			if totalDiags != tc.wantDiags {
				t.Errorf("[%s] got %d diagnostics, want %d", tc.name, totalDiags, tc.wantDiags)
			}
		})
	}
}
