package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// ReactRedundantFunctionMemoizationRule mengaudit pembungkusan useCallback yang sia-sia karena hanya diteruskan ke elemen native DOM.
type ReactRedundantFunctionMemoizationRule struct{}

// NewReactRedundantFunctionMemoizationRule membuat instance baru dari ReactRedundantFunctionMemoizationRule.
func NewReactRedundantFunctionMemoizationRule() *ReactRedundantFunctionMemoizationRule {
	return &ReactRedundantFunctionMemoizationRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *ReactRedundantFunctionMemoizationRule) ID() string {
	return "performance.react-redundant-function-memoization"
}

// Category mengembalikan kategori aturan ('performance').
func (r *ReactRedundantFunctionMemoizationRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info/advisory).
func (r *ReactRedundantFunctionMemoizationRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *ReactRedundantFunctionMemoizationRule) Description() string {
	return "Mengaudit penggunaan useCallback pada callback yang hanya dikonsumsi oleh elemen native HTML tanpa konsumen peka identitas referensial."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *ReactRedundantFunctionMemoizationRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"React Official Documentation (When to use useCallback & Hook Overhead)",
			"React Compiler Architecture Specification (Automated Memoization Economy)",
			"Dan Abramov Architecture Notes ('A Complete Guide to useEffect & useCallback')",
		},
		CoreInvariant: "Functions passed exclusively to native HTML elements (<button>, <input>) must not be wrapped in 'useCallback'; native DOM elements do not perform shallow equality checks, making hook allocation a net negative overhead.",
		Grounding: "A common misconception among React developers is that wrapping every function in `useCallback` improves performance.\n\n" +
			"In reality, `useCallback` requires allocating an internal Hook cell, preserving a dependency array in memory, and executing array comparisons on every render cycle.\n\n" +
			"Native HTML elements (`<button onClick={...}>`) do not inspect prop referential equality; they simply attach or update event listeners. Unless a callback is passed to a `React.memo` component or included in another hook's dependency list, `useCallback` introduces pure overhead with zero performance benefit.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Hook Memory & GC Overhead",
				Severity: "LOW",
				Impact:   "Increases memory footprint and garbage collector pressure by retaining closures and dependency arrays across component lifecycles.",
			},
			{
				Vector:   "Codebase Complexity & Clutter",
				Severity: "LOW",
				Impact:   "Obscures real optimization sites and complicates eventual migration to automatic compiler memoization (React Compiler).",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Membungkus handler tombol native dengan useCallback adalah pemborosan hook",
				Code: `const handleClick = useCallback(() => {
  setOpen(true);
}, []);

return <button onClick={handleClick}>Buka Modal</button>;`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Gunakan deklarasi fungsi reguler untuk elemen DOM biasa",
				Code: `const handleClick = () => {
  setOpen(true);
};

return <button onClick={handleClick}>Buka Modal</button>;`,
			},
		},
	}
}

// Evaluate memeriksa apakah fungsi yang di-memoize hanya dikonsumsi oleh elemen HTML native.
func (r *ReactRedundantFunctionMemoizationRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || !isSourceRootOrScript(node) {
		return nil
	}

	fileSrc := getFileSourceContent(node)
	if len(fileSrc) == 0 {
		return nil
	}

	violations := findRedundantFunctionMemoizations(fileSrc)
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
			Message:  fmt.Sprintf("Handler '%s' is memoized with 'useCallback' but is only consumed by native DOM element '<%s>'. Wrapping callbacks passed only to native elements adds unnecessary hook memory and overhead without performance benefit.", v.HandlerName, v.ElementTag),
			Hint:     "Replace 'useCallback' with a plain arrow function or function declaration since native DOM elements do not participate in React referential equality checks.",
		})
	}

	return diags
}
