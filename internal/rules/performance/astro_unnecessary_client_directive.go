package performance

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// AstroUnnecessaryClientDirectiveRule mengaudit komponen statis yang dipaksa terhidrasi menggunakan direktif client.
type AstroUnnecessaryClientDirectiveRule struct{}

// NewAstroUnnecessaryClientDirectiveRule membuat instance baru dari AstroUnnecessaryClientDirectiveRule.
func NewAstroUnnecessaryClientDirectiveRule() *AstroUnnecessaryClientDirectiveRule {
	return &AstroUnnecessaryClientDirectiveRule{}
}

// ID mengembalikan identifier unik kanonikal aturan.
func (r *AstroUnnecessaryClientDirectiveRule) ID() string {
	return "performance.astro-unnecessary-client-directive"
}

// Category mengembalikan kategori aturan ('performance').
func (r *AstroUnnecessaryClientDirectiveRule) Category() string {
	return "performance"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *AstroUnnecessaryClientDirectiveRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Description mengembalikan deskripsi tujuan dan konteks aturan.
func (r *AstroUnnecessaryClientDirectiveRule) Description() string {
	return "Menegakkan prinsip Zero-JS Astro dengan melarang penambahan direktif hidrasi (client:*) pada komponen antarmuka yang murni statis."
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *AstroUnnecessaryClientDirectiveRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Astro Islands Architecture Specification (Zero-JS Baseline Principle)",
			"W3C Web Performance Client-Side Script Minimization Invariants",
			"Astro Official Documentation ('Template Directives: client:*')",
		},
		CoreInvariant: "Static UI components must not include 'client:*' hydration directives; adding hydration directives to non-interactive components forces the framework runtime and component bundle to be downloaded, violating Astro's Zero-JS guarantee.",
		Grounding: "Astro by default renders all components to pure, static HTML at build time with zero client-side JavaScript overhead.\n\n" +
			"When a developer unnecessarily adds a `client:*` directive (`client:load`, `client:idle`, `client:visible`) to a purely presentational component, Astro treats it as an interactive island.\n\n" +
			"This forces the bundler to extract the component into a separate client bundle and ship the framework runtime (such as React or Vue, weighing 30-50KB+) to the browser, needlessly delaying page interactivity and squandering network bandwidth.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Zero-JS Guarantee Violation",
				Severity: "HIGH",
				Impact:   "Transmits unnecessary framework runtimes and component code to the client, increasing page weight and parse time.",
			},
			{
				Vector:   "Main Thread Hydration Lag",
				Severity: "MEDIUM",
				Impact:   "Wastes browser CPU cycles hydrating static DOM trees that have no event listeners or interactive state.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Header statis dipaksa terhidrasi ke peramban",
				Code: `---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic client:load title="Selamat Datang" />`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "astro",
				Comment:  "Dirender sebagai pure static HTML tanpa JavaScript",
				Code: `---
import HeaderStatic from '../components/HeaderStatic.tsx';
---
<HeaderStatic title="Selamat Datang" />`,
			},
		},
	}
}

// Evaluate memeriksa apakah komponen statis menyertakan direktif client:*.
func (r *AstroUnnecessaryClientDirectiveRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement || len(node.Attributes) == 0 {
		return nil
	}

	dir, isClient := hasClientDirective(node)
	if !isClient {
		return nil
	}

	if !isStaticComponentTag(node.Tag) {
		return nil
	}

	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  fmt.Sprintf("Static component '<%s>' is hydrated with '%s', sending unnecessary JavaScript runtime to the client and violating Astro Zero-JS architecture.", node.Tag, dir),
			Hint:     fmt.Sprintf("Remove the '%s' directive so '<%s>' is rendered as pure zero-JS static HTML.", dir, node.Tag),
		},
	}
}
