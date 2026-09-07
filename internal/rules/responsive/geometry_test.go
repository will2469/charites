package responsive_test

import (
	"testing"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules/responsive"
)

func TestGeometry_ParseFractionalWidth(t *testing.T) {
	cases := []struct {
		input string
		want  float64
		ok    bool
	}{
		{"w-1/2", 0.5, true},
		{"w-1/3", 1.0 / 3.0, true},
		{"w-2/3", 2.0 / 3.0, true},
		{"w-1/4", 0.25, true},
		{"w-3/4", 0.75, true},
		{"w-1/5", 0.2, true},
		{"w-4/5", 0.8, true},
		{"w-1/6", 1.0 / 6.0, true},
		{"w-5/6", 5.0 / 6.0, true},
		{"w-1/12", 1.0 / 12.0, true},
		{"w-11/12", 11.0 / 12.0, true},
		{"w-[60%]", 0.60, true},
		{"w-[33.33%]", 0.3333, true},
		{"w-[100%]", 0, false}, // 100% is full width, not fractional
		{"w-full", 0, false},
		{"w-auto", 0, false},
		{"w-64", 0, false},
		{"w-2/2", 0, false}, // equal to 1, not fraction
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			node := &ir.Node{
				Type:    ir.NodeElement,
				Classes: []string{tc.input},
			}
			geom := responsive.ResolveGeometry(&ir.Node{
				Type:     ir.NodeElement,
				Classes:  []string{"flex"},
				Children: []*ir.Node{node},
			})
			if len(geom.Children) != 1 {
				t.Fatalf("expected 1 child geometry, got %d", len(geom.Children))
			}
			cw := geom.Children[0].Width[responsive.TierBaseline]
			if cw.IsFractional != tc.ok {
				t.Fatalf("[%s] IsFractional = %v, want %v", tc.input, cw.IsFractional, tc.ok)
			}
			if tc.ok {
				diff := cw.Ratio - tc.want
				if diff < -0.001 || diff > 0.001 {
					t.Fatalf("[%s] Ratio = %f, want %f", tc.input, cw.Ratio, tc.want)
				}
			}
		})
	}
}

func TestGeometry_GapPrecedence_OrderInsensitive(t *testing.T) {
	// 1. gap-4 gap-x-0: gap-x-0 wins over shorthand gap-4
	node1 := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-4", "gap-x-0"},
	}
	geom1 := responsive.ResolveGeometry(node1)
	if geom1.Container[responsive.TierBaseline].Gap.HasHorizontal {
		t.Errorf("expected gap-x-0 to override horizontal gap to false")
	}

	// 2. gap-x-0 gap-4: gap-x-0 must STILL win despite being earlier in markup
	node2 := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-x-0", "gap-4"},
	}
	geom2 := responsive.ResolveGeometry(node2)
	if geom2.Container[responsive.TierBaseline].Gap.HasHorizontal {
		t.Errorf("expected gap-x-0 to override horizontal gap to false regardless of class order")
	}

	// 3. gap-y-4: horizontal is false, vertical is true
	node3 := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-y-4"},
	}
	geom3 := responsive.ResolveGeometry(node3)
	if geom3.Container[responsive.TierBaseline].Gap.HasHorizontal {
		t.Errorf("expected gap-y-4 to not declare horizontal gap")
	}
	if !geom3.Container[responsive.TierBaseline].Gap.HasVertical {
		t.Errorf("expected gap-y-4 to declare vertical gap")
	}
}

func TestGeometry_DisplayOverride_DeactivatesFlex(t *testing.T) {
	node := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "md:block"},
	}
	geom := responsive.ResolveGeometry(node)

	if !geom.Container[responsive.TierBaseline].Display.IsFlex() {
		t.Errorf("expected baseline to be flex")
	}
	if geom.Container[responsive.TierMd].Display.IsFlex() {
		t.Errorf("expected md:block to override and deactivate flex on md tier")
	}
}

func TestGeometry_Classifier_NativeShrinkVsWrapVsNoShrink(t *testing.T) {
	// 1. Default nowrap + default shrink-1: SAFE (nil)
	safeNode := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
		},
	}
	geomSafe := responsive.ResolveGeometry(safeNode)
	if finding := responsive.ClassifyFlexGapDrift(geomSafe); finding != nil {
		t.Fatalf("expected nil for default nowrap + shrink-1, got %+v", finding)
	}

	// 2. Wrap active: DRIFT
	wrapNode := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "flex-wrap", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
		},
	}
	geomWrap := responsive.ResolveGeometry(wrapNode)
	if finding := responsive.ClassifyFlexGapDrift(geomWrap); finding == nil {
		t.Fatalf("expected drift finding for flex-wrap, got nil")
	}

	// 3. One fractional child has shrink-0: DRIFT
	shrink0Node := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-1/2", "shrink-0"}},
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
		},
	}
	geomShrink0 := responsive.ResolveGeometry(shrink0Node)
	if finding := responsive.ClassifyFlexGapDrift(geomShrink0); finding == nil {
		t.Fatalf("expected drift finding when fractional child has shrink-0, got nil")
	}

	// 4. Non-fractional child has shrink-0: SAFE (nil)
	nonFracShrink0Node := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
			{Type: ir.NodeElement, Classes: []string{"w-1/2"}},
			{Type: ir.NodeElement, Classes: []string{"w-auto", "shrink-0"}},
		},
	}
	geomNonFrac := responsive.ResolveGeometry(nonFracShrink0Node)
	if finding := responsive.ClassifyFlexGapDrift(geomNonFrac); finding != nil {
		t.Fatalf("expected nil when non-fractional child has shrink-0, got %+v", finding)
	}

	// 5. Fraction sum under 100%: SAFE (nil)
	sumUnder100 := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-1/4", "shrink-0"}},
			{Type: ir.NodeElement, Classes: []string{"w-1/4"}},
		},
	}
	geomSumUnder := responsive.ResolveGeometry(sumUnder100)
	if finding := responsive.ClassifyFlexGapDrift(geomSumUnder); finding != nil {
		t.Fatalf("expected nil when fraction sum < 100%%, got %+v", finding)
	}

	// 6. Carousel Intent Heuristic: SAFE (nil)
	carouselNode := &ir.Node{
		Type:    ir.NodeElement,
		Classes: []string{"flex", "overflow-x-auto", "gap-4"},
		Children: []*ir.Node{
			{Type: ir.NodeElement, Classes: []string{"w-4/5", "shrink-0"}},
			{Type: ir.NodeElement, Classes: []string{"w-4/5", "shrink-0"}},
		},
	}
	geomCarousel := responsive.ResolveGeometry(carouselNode)
	if finding := responsive.ClassifyFlexGapDrift(geomCarousel); finding != nil {
		t.Fatalf("expected nil for carousel intent, got %+v", finding)
	}
}
