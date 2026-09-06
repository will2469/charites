package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactUnstableHookReferenceRule mengaudit custom hook yang mengembalikan fungsi tidak stabil tanpa useCallback.
type ReactUnstableHookReferenceRule struct{}

// NewReactUnstableHookReferenceRule membuat instance baru dari ReactUnstableHookReferenceRule.
func NewReactUnstableHookReferenceRule() *ReactUnstableHookReferenceRule {
	return &ReactUnstableHookReferenceRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *ReactUnstableHookReferenceRule) ID() string {
	return "performance.react-unstable-hook-reference"
}

// Category mengembalikan kategori aturan ('performance').
func (r *ReactUnstableHookReferenceRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warning).
func (r *ReactUnstableHookReferenceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *ReactUnstableHookReferenceRule) Description() string {
	return "Mengaudit custom hook yang mengembalikan referensi fungsi tidak stabil tanpa dibungkus useCallback, yang memicu re-render loop pada komponen konsumen."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactUnstableHookReferenceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Documentation (Building Your Own Custom Hooks)",
			"React Hooks Referential Integrity & Stable Function Contracts",
			"React Hooks Exhaustive Dependencies Safety Guidelines",
		},
		CoreInvariant: "Custom React hooks exposing helper functions must stabilize them with 'useCallback'; returning fresh function instances causes downstream consumers using them in effect dependencies to trigger infinite render loops.",
		Grounding: "Custom hooks frequently return an object containing state and mutation functions (e.g. `{ data, refetch, reset }`).\n\n" +
			"If these functions are defined as regular arrow functions without `useCallback`, a brand-new function reference is created in memory on every render pass of the consuming component.\n\n" +
			"When the consuming component passes this function into the dependency array of `useEffect` or `useMemo`, or passes it down to a memoized child component, the newly allocated reference violates referential equality, defeating memoization and in many cases causing uncontrollable infinite re-render loops.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Downstream Infinite Render Loops",
				Severity: "HIGH",
				Impact:   "Triggers continuous re-execution of downstream useEffect hooks that list the unmemoized helper function in their dependency arrays.",
			},
			{
				Vector:   "Bypassed Child Memoization",
				Severity: "MEDIUM",
				Impact:   "Breaks shallow prop comparison (React.memo) across all child components consuming functions returned from the custom hook.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "refetch dialokasikan sebagai fungsi baru di setiap pemanggilan hook",
				Code: `export function useProfile(userId: string) {
  const [data, setData] = useState(null);
  const refetch = () => { fetchProfile(userId).then(setData); };
  return { data, refetch };
}`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Menstabilkan referensi fungsi dengan useCallback",
				Code: `export function useProfile(userId: string) {
  const [data, setData] = useState(null);
  const refetch = useCallback(() => {
    fetchProfile(userId).then(setData);
  }, [userId]);
  return { data, refetch };
}`,
			},
		},
	}
}

// Evaluate memeriksa apakah custom hook mengembalikan referensi fungsi tidak stabil.
func (r *ReactUnstableHookReferenceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !isSourceRootOrScript(node) {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	if len(fileSrc) == 0 {
		return nil
	}

	violations := findUnstableHookReferences(fileSrc)
	if len(violations) == 0 {
		return nil
	}

	diags := make([]ir.Diagnostic, 0, len(violations))
	for _, v := range violations {
		diags = append(diags, ir.Diagnostic{
			Line:     v.Line,
			Column:   1,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Custom hook '%s' exports unstable function reference '%s' without memoizing it via 'useCallback'. Consuming components passing this function into effect dependency arrays or memoized props will trigger infinite re-renders or defeated memoization.", v.HookName, v.FunctionName),
			Hint:     fmt.Sprintf("Wrap '%s' in 'useCallback' with proper dependencies before returning it from custom hook '%s'.", v.FunctionName, v.HookName),
		})
	}

	return diags
}
