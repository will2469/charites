package drift_test

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/ir"
)

func TestResolveScope(t *testing.T) {
	tests := []struct {
		tag       string
		role      string
		wantKind  drift.ComponentKind
		wantConf  drift.ScopeConfidence
		wantFound bool
	}{
		{tag: "Button", wantKind: drift.KindButton, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "button", wantKind: drift.KindButton, wantConf: drift.ConfidenceNative, wantFound: true},
		{tag: "div", role: "button", wantKind: drift.KindButton, wantConf: drift.ConfidenceSemantic, wantFound: true},
		{tag: "Input", wantKind: drift.KindInput, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "textarea", wantKind: drift.KindInput, wantConf: drift.ConfidenceNative, wantFound: true},
		{tag: "Card", wantKind: drift.KindCard, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "DialogContent", wantKind: drift.KindDialog, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "SheetContent", wantKind: drift.KindSheet, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "Badge", wantKind: drift.KindBadge, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "div", role: "status", wantKind: drift.KindBadge, wantConf: drift.ConfidenceSemantic, wantFound: true},
		{tag: "PopoverContent", wantKind: drift.KindMenu, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "Alert", wantKind: drift.KindAlert, wantConf: drift.ConfidenceExact, wantFound: true},
		{tag: "div", wantFound: false},
	}

	for _, tt := range tests {
		node := &ir.Node{
			Type: ir.NodeElement,
			Tag:  tt.tag,
		}
		if tt.role != "" {
			node.Attributes = map[string]string{"role": tt.role}
		}

		gotScope, ok := drift.ResolveScope(node)
		if ok != tt.wantFound {
			t.Errorf("ResolveScope(%s, role=%s) found = %v, want %v", tt.tag, tt.role, ok, tt.wantFound)
			continue
		}
		if tt.wantFound {
			if gotScope.Kind != tt.wantKind || gotScope.Confidence != tt.wantConf {
				t.Errorf("ResolveScope(%s, role=%s) = %v, want Kind=%v, Conf=%v",
					tt.tag, tt.role, gotScope, tt.wantKind, tt.wantConf)
			}
		}
	}
}

func TestExtractOccurrences_FineGrained(t *testing.T) {
	node := &ir.Node{
		Type:       ir.NodeElement,
		Tag:        "Button",
		RawClasses: "rounded-md active:scale-95 disabled:opacity-50 focus-visible:ring-2 bg-primary",
		Classes: []string{
			"rounded-md",
			"active:scale-95",
			"disabled:opacity-50",
			"focus-visible:ring-2",
			"bg-primary",
		},
		Attributes: map[string]string{
			"size": "default",
		},
	}

	occs := drift.ExtractOccurrences("Test.tsx", node)
	if len(occs) != 5 {
		t.Fatalf("expected 5 occurrences, got %d", len(occs))
	}

	foundMap := make(map[drift.PropertyCategory]drift.StyleOccurrence)
	for _, o := range occs {
		foundMap[o.Category] = o
	}

	if o, ok := foundMap[drift.CatRounded]; !ok || o.Token != "rounded-md" {
		t.Errorf("missing or invalid CatRounded: %+v", o)
	}
	if o, ok := foundMap[drift.CatActiveScale]; !ok || o.Token != "scale-95" || o.Variant != "active" {
		t.Errorf("missing or invalid CatActiveScale: %+v", o)
	}
	if o, ok := foundMap[drift.CatDisabledOpacity]; !ok || o.Token != "opacity-50" || o.Variant != "disabled" {
		t.Errorf("missing or invalid CatDisabledOpacity: %+v", o)
	}
	if o, ok := foundMap[drift.CatFocusRingWidth]; !ok || o.Token != "ring-2" || o.Variant != "focus-visible" {
		t.Errorf("missing or invalid CatFocusRingWidth: %+v", o)
	}
	if o, ok := foundMap[drift.CatBackground]; !ok || o.Token != "bg-primary" {
		t.Errorf("missing or invalid CatBackground: %+v", o)
	}
}

func TestContrastRatio_WCAG(t *testing.T) {
	// 1. yellow-400 (#facc15) dengan white (#ffffff)
	ratio, failsAA, ok := drift.CheckContrastHazard("bg-yellow-400", "text-white")
	if !ok {
		t.Fatalf("expected CheckContrastHazard to resolve bg-yellow-400 and text-white")
	}
	if !failsAA {
		t.Errorf("expected bg-yellow-400 + text-white to fail WCAG AA (< 4.5:1)")
	}
	// Rasio kontras yellow-400 dan white adalah sekitar 1.53:1
	if math.Abs(ratio-1.53) > 0.05 {
		t.Errorf("expected contrast ratio ~1.53, got %.2f", ratio)
	}

	// 2. black dengan white
	ratioBW, failsBW, okBW := drift.CheckContrastHazard("bg-black", "text-white")
	if !okBW || failsBW || ratioBW < 20.0 {
		t.Errorf("expected bg-black + text-white to pass WCAG AA, got ratio %.2f, fails=%v", ratioBW, failsBW)
	}

	// 3. Token dinamis (var(--warning)) diabaikan secara anggun
	_, _, okVar := drift.CheckContrastHazard("bg-[var(--warning)]", "text-white")
	if okVar {
		t.Errorf("expected dynamic var(--warning) to return ok=false")
	}
}

func TestClusterOccurrences_StatisticalGates(t *testing.T) {
	scope := drift.ScopeKey{Kind: drift.KindButton, Confidence: drift.ConfidenceExact}

	t.Run("InsufficientDataGate", func(t *testing.T) {
		// 2x rounded-md, 1x rounded-lg (total 3 < min 10)
		var occs []drift.StyleOccurrence
		for i := 0; i < 2; i++ {
			occs = append(occs, drift.StyleOccurrence{
				Scope:    scope,
				Category: drift.CatRounded,
				Token:    "rounded-md",
			})
		}
		occs = append(occs, drift.StyleOccurrence{
			Scope:    scope,
			Category: drift.CatRounded,
			Token:    "rounded-lg",
		})

		clusters := drift.ClusterOccurrences(occs, drift.DefaultOptions())
		key := drift.ClusterKey{Scope: scope, Category: drift.CatRounded}
		c := clusters[key]
		if c == nil {
			t.Fatalf("expected cluster for key, got nil")
		}
		if !c.InsufficientData {
			t.Errorf("expected InsufficientData=true for total 3 < 10, got false")
		}
		if len(c.Outliers) > 0 {
			t.Errorf("expected 0 outliers for InsufficientData, got %d", len(c.Outliers))
		}
	})

	t.Run("ValidDriftDecision", func(t *testing.T) {
		// 10x rounded-md (90.9%), 1x rounded-2xl (9.1%)
		var occs []drift.StyleOccurrence
		for i := 0; i < 10; i++ {
			occs = append(occs, drift.StyleOccurrence{
				FilePath: "ButtonCanonical.tsx",
				Scope:    scope,
				Category: drift.CatRounded,
				Token:    "rounded-md",
			})
		}
		occs = append(occs, drift.StyleOccurrence{
			FilePath: "ButtonRogue.tsx",
			Scope:    scope,
			Category: drift.CatRounded,
			Token:    "rounded-2xl",
		})

		clusters := drift.ClusterOccurrences(occs, drift.DefaultOptions())
		key := drift.ClusterKey{Scope: scope, Category: drift.CatRounded}
		c := clusters[key]
		if c == nil {
			t.Fatalf("expected cluster for key, got nil")
		}
		if c.InsufficientData {
			t.Errorf("expected InsufficientData=false for total 11 >= 10")
		}
		if c.CanonicalToken != "rounded-md" {
			t.Errorf("expected canonical 'rounded-md', got %s", c.CanonicalToken)
		}
		if len(c.Outliers) != 1 {
			t.Fatalf("expected 1 outlier, got %d", len(c.Outliers))
		}
		if c.Outliers[0].Token != "rounded-2xl" {
			t.Errorf("expected outlier 'rounded-2xl', got %s", c.Outliers[0].Token)
		}
	})
}

func TestClusterOccurrences_ContextExceptions(t *testing.T) {
	scope := drift.ScopeKey{Kind: drift.KindButton, Confidence: drift.ConfidenceExact}

	// 10x rounded-md, 1x rounded-full dengan size="icon"
	var occs []drift.StyleOccurrence
	for i := 0; i < 10; i++ {
		occs = append(occs, drift.StyleOccurrence{
			FilePath: "Button.tsx",
			Scope:    scope,
			Category: drift.CatRounded,
			Token:    "rounded-md",
		})
	}
	occs = append(occs, drift.StyleOccurrence{
		FilePath: "IconButton.tsx",
		Scope:    scope,
		Category: drift.CatRounded,
		Token:    "rounded-full",
		Props:    map[string]string{"size": "icon"},
	})

	opts := drift.DefaultOptions()
	opts.Exceptions["Button"] = map[string][]drift.ContextException{
		"rounded": {
			{Token: "rounded-full", Props: map[string]string{"size": "icon"}},
		},
	}

	clusters := drift.ClusterOccurrences(occs, opts)
	key := drift.ClusterKey{Scope: scope, Category: drift.CatRounded}
	c := clusters[key]
	if c == nil {
		t.Fatalf("expected cluster, got nil")
	}
	if len(c.Outliers) != 0 {
		t.Errorf("expected rounded-full exception to be filtered out, got %d outliers", len(c.Outliers))
	}
}

func TestAnalyzeChromaticDrift(t *testing.T) {
	scope := drift.ScopeKey{Kind: drift.KindButton, Confidence: drift.ConfidenceExact}

	// 6x bg-amber-500, 3x bg-yellow-400, 1x bg-orange-500
	var occs []drift.StyleOccurrence
	addHues := func(token string, count int) {
		for i := 0; i < count; i++ {
			occs = append(occs, drift.StyleOccurrence{
				FilePath: "Form.tsx",
				Scope:    scope,
				Category: drift.CatBackground,
				Token:    token,
			})
		}
	}

	addHues("bg-amber-500", 6)
	addHues("bg-yellow-400", 3)
	addHues("bg-orange-500", 1)

	chromClusters := drift.AnalyzeChromaticDrift(occs)
	if len(chromClusters) != 1 {
		t.Fatalf("expected 1 chromatic drift cluster, got %d", len(chromClusters))
	}

	cluster := chromClusters[0]
	if cluster.Family != drift.FamilyWarmAmber {
		t.Errorf("expected FamilyWarmAmber, got %v", cluster.Family)
	}
	if cluster.CandidateIntent != drift.CandidateWarning {
		t.Errorf("expected CandidateWarning, got %v", cluster.CandidateIntent)
	}
	if len(cluster.HueCounts) != 3 {
		t.Errorf("expected 3 competing hues, got %d", len(cluster.HueCounts))
	}
	if !strings.Contains(cluster.Recommendation, "--warning") {
		t.Errorf("expected recommendation to mention --warning, got: %s", cluster.Recommendation)
	}
}

func TestRenderReport(t *testing.T) {
	scope := drift.ScopeKey{Kind: drift.KindButton, Confidence: drift.ConfidenceExact}
	var occs []drift.StyleOccurrence
	for i := 0; i < 10; i++ {
		occs = append(occs, drift.StyleOccurrence{
			FilePath:   "Button.tsx",
			Scope:      scope,
			Category:   drift.CatRounded,
			Token:      "rounded-md",
			RawClasses: "rounded-md bg-amber-500 text-white",
		})
	}
	occs = append(occs, drift.StyleOccurrence{
		FilePath:   "Modal.tsx",
		Scope:      scope,
		Category:   drift.CatRounded,
		Token:      "rounded-2xl",
		RawClasses: "rounded-2xl bg-yellow-400 text-white",
	})

	analyzer := drift.NewAnalyzer(drift.DefaultOptions())
	report := analyzer.Analyze(occs)

	// 1. Inline
	var bufInline bytes.Buffer
	if err := drift.RenderReport(&bufInline, &report, "inline", true); err != nil {
		t.Fatalf("RenderReport inline failed: %v", err)
	}
	inlineOut := bufInline.String()
	if !strings.Contains(inlineOut, "CANONICAL") || !strings.Contains(inlineOut, "rounded-2xl") {
		t.Errorf("inline report missing key content: %s", inlineOut)
	}

	// 2. Markdown
	var bufMD bytes.Buffer
	if err := drift.RenderReport(&bufMD, &report, "markdown", true); err != nil {
		t.Fatalf("RenderReport markdown failed: %v", err)
	}
	mdOut := bufMD.String()
	if !strings.Contains(mdOut, "# Charites Component-Scoped Style Drift Report") {
		t.Errorf("markdown report missing header: %s", mdOut)
	}

	// 3. JSON
	var bufJSON bytes.Buffer
	if err := drift.RenderReport(&bufJSON, &report, "json", true); err != nil {
		t.Fatalf("RenderReport json failed: %v", err)
	}
	jsonOut := bufJSON.String()
	if !strings.Contains(jsonOut, `"total_occurrences": 11`) {
		t.Errorf("json report missing total: %s", jsonOut)
	}

	// 4. SynthesizeDiagnostics
	diags := drift.SynthesizeDiagnostics(report)
	if len(diags) == 0 {
		t.Errorf("expected synthesized diagnostics > 0, got 0")
	}
	var foundRogueDiag bool
	for _, d := range diags {
		if strings.Contains(d.Message, "Rogue variation 'rounded-2xl'") {
			foundRogueDiag = true
			break
		}
	}
	if !foundRogueDiag {
		t.Errorf("expected diagnostic for rogue variation 'rounded-2xl'")
	}
}
