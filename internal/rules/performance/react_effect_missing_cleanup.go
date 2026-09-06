package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactEffectMissingCleanupRule mendeteksi penggunaan useEffect atau useLayoutEffect
// yang mengakuisisi resource persisten (listener, timer, observer) tanpa fungsi cleanup simetris.
type ReactEffectMissingCleanupRule struct{}

// NewReactEffectMissingCleanupRule membuat instance baru dari ReactEffectMissingCleanupRule.
func NewReactEffectMissingCleanupRule() *ReactEffectMissingCleanupRule {
	return &ReactEffectMissingCleanupRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *ReactEffectMissingCleanupRule) ID() string {
	return "performance.react-effect-missing-cleanup"
}

// Description mengembalikan ringkasan aturan.
func (r *ReactEffectMissingCleanupRule) Description() string {
	return "Effect hook acquiring persistent resource (listener, interval, observer) lacks a symmetrical cleanup return function, causing memory leaks"
}

// Category mengembalikan nama kategori rule.
func (r *ReactEffectMissingCleanupRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *ReactEffectMissingCleanupRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactEffectMissingCleanupRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Documentation (Synchronizing with Effects & Effect Cleanup Invariants)",
			"W3C EventTarget and Observer Lifecycle Specifications",
			"Google Chrome Memory Profiling Guidelines (Retained DOM Detached Node Prevention)",
		},
		CoreInvariant: "React effect hooks ('useEffect', 'useLayoutEffect') that acquire persistent resources (event listeners, intervals, observers, WebSockets) must return a symmetrical cleanup function to release references upon unmount or dependency changes.",
		Grounding: "When an effect registers an external subscription (such as `window.addEventListener`, `setInterval`, or an `IntersectionObserver`) without returning a cleanup function, that subscription remains active in the browser memory even after the component is unmounted.\n\n" +
			"The orphaned subscription retains references to component state, props, and closures, preventing the JavaScript garbage collector from reclaiming the component tree's memory.\n\n" +
			"Furthermore, triggered callbacks continue attempting to execute against unmounted components, causing unhandled errors, stale state updates, and compounding memory leaks during client-side navigation.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Persistent Memory Leaks",
				Severity: "CRITICAL",
				Impact:   "Orphaned event listeners and observers retain unmounted component closures, leading to runaway heap memory growth in single-page applications.",
			},
			{
				Vector:   "Zombie Handler Execution",
				Severity: "HIGH",
				Impact:   "Callbacks trigger state updates on unmounted components, causing React warnings and erratic background behavior.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Window event listener registered in useEffect without cleanup return function",
				Code: `useEffect(() => {
  const onResize = () => setWidth(window.innerWidth);
  window.addEventListener('resize', onResize);
}, []);`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Symmetrical cleanup function returned to remove listener on unmount",
				Code: `useEffect(() => {
  const onResize = () => setWidth(window.innerWidth);
  window.addEventListener('resize', onResize);
  return () => window.removeEventListener('resize', onResize);
}, []);`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat pemanggilan effect hook yang kehilangan fungsi cleanup.
func (r *ReactEffectMissingCleanupRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if !isSourceRootOrScript(node) {
		return nil
	}

	src := getFileSourceContent(node)
	if src == "" {
		return nil
	}

	violations := findMissingEffectCleanups(src)
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
			Message:  fmt.Sprintf("Hook '%s' acquires persistent resource '%s' without returning a symmetrical cleanup function. Failing to remove subscriptions or timers triggers permanent memory leaks and zombie callback execution after component unmount.", v.EffectName, v.Resource),
			Hint:     "Return a cleanup function (e.g. 'return () => { window.removeEventListener(...); }' or 'clearInterval') at the end of the effect callback.",
		})
	}

	return diags
}
