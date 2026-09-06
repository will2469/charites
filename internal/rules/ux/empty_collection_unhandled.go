package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// EmptyCollectionUnhandledRule menyarankan penyediaan penanganan cabang kondisi koleksi kosong
// pada pemetaan daftar dinamis ({items.map(...)}) untuk mencegah kebutaan state kosong (Zero-State Blindness).
type EmptyCollectionUnhandledRule struct{}

// NewEmptyCollectionUnhandledRule membuat instance baru dari EmptyCollectionUnhandledRule.
func NewEmptyCollectionUnhandledRule() *EmptyCollectionUnhandledRule {
	return &EmptyCollectionUnhandledRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *EmptyCollectionUnhandledRule) ID() string {
	return "ux.empty-collection-unhandled"
}

// Description mengembalikan ringkasan aturan.
func (r *EmptyCollectionUnhandledRule) Description() string {
	return "Advises handling empty collection state when mapping dynamic items to avoid zero-state blindness"
}

// Category mengembalikan nama kategori rule.
func (r *EmptyCollectionUnhandledRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (info/advisory).
func (r *EmptyCollectionUnhandledRule) DefaultSeverity() ir.Severity {
	return ir.SeverityInfo
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *EmptyCollectionUnhandledRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Zero-State Usability & Mental Model Continuity (Nielsen Norman Group)",
			"Feedforward Principle & Gulf of Evaluation (Don Norman)",
			"ISO 9241-110 Ergonomics of Human-System Interaction (Suitability for Learning)",
		},
		CoreInvariant: "Dynamic collection rendering expressions must handle empty collection states ('collection.length === 0') with informative fallback zero-state UI.",
		Grounding: "When dynamic lists, tables, or feed collections contain 0 records and render nothing, " +
			"users are stranded in an ambiguous vacuum: did the request fail, is it still loading, or are there genuinely zero records?\n\n" +
			"Zero-state blindness forces users to refresh repeatedly or assume the application is broken. " +
			"A dedicated empty state component (e.g. '<EmptyState />' with an illustration, clarifying text, and a call-to-action button) " +
			"confirms system status and proactively guides user next steps.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Zero-State Blindness & System Status Ambiguity",
				Severity: "LOW",
				Impact:   "Users perceive blank empty containers as silent application crashes or perpetual loading freezes.",
			},
			{
				Vector:   "Workflow Dead Ends",
				Severity: "LOW",
				Impact:   "Without an actionable empty state CTA (e.g. 'Create First Invoice'), users cannot self-discover how to populate the collection.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Dynamic list rendering items via .map() without handling empty array state",
				Code: `<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Explicit empty state fallback branch when collection has 0 items",
				Code: `<div className="space-y-3">
  <h2 className="text-lg font-bold">Daftar Tagihan</h2>
  {invoices.length === 0 ? (
    <EmptyState
      title="Belum Ada Tagihan"
      description="Buat tagihan pertama Anda untuk mulai menerima pembayaran."
      actionText="Buat Tagihan"
    />
  ) : (
    <List items={invoices.map(inv => <InvoiceRow key={inv.id} data={inv} />)} />
  )}
</div>`,
			},
		},
	}
}

// Evaluate memeriksa apakah perenderan koleksi dinamis memiliki penanganan state kosong.
func (r *EmptyCollectionUnhandledRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	// Kasus 1: Atribut yang memuat ekspresi mapping (.map(...) tanpa cek length/empty)
	for attrName, attrVal := range node.Attributes {
		if strings.Contains(attrVal, ".map(") {
			if hasEmptyStateInSubtree(node) {
				return nil
			}
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Dynamic collection mapped in attribute %q lacks explicit empty state handling ('length === 0' or <EmptyState /> fallback).",
						attrName,
					),
					Hint: "Provide an empty state fallback branch when the collection contains 0 items to prevent zero-state blindness.",
				},
			}
		}
	}

	// Kasus 2: Kontainer daftar dinamis dengan item berulang tanpa penanganan empty state
	if isDynamicListContainer(node) {
		if hasEmptyStateInSubtree(node) {
			return nil
		}
		if hasDynamicItemChildren(node) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Collection container <%s> renders dynamic item components without zero-state fallback.",
						node.Tag,
					),
					Hint: "Add a conditional branch (e.g. 'items.length === 0 ? <EmptyState /> : ...') to guide users when no items exist.",
				},
			}
		}
	}

	return nil
}

func isDynamicListContainer(node *ir.Node) bool {
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "list" || tagLower == "feed" || tagLower == "collection" ||
		tagLower == "tbody" || strings.HasSuffix(tagLower, "list") {
		return true
	}
	for _, cls := range node.Classes {
		base := strings.ToLower(StripVariantsOnlyBase(cls))
		if base == "item-list" || base == "data-list" || base == "invoice-list" || base == "card-grid" {
			return true
		}
	}
	return false
}

func hasDynamicItemChildren(node *ir.Node) bool {
	for _, child := range node.Children {
		if child.Type != ir.NodeElement {
			continue
		}
		tagLower := strings.ToLower(child.Tag)
		if strings.HasSuffix(tagLower, "row") || strings.HasSuffix(tagLower, "card") ||
			strings.HasSuffix(tagLower, "item") {
			return true
		}
	}
	return false
}
