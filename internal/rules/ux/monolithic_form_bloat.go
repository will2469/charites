package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// MonolithicFormBloatRule mendeteksi form masif tanpa pemisahan batas kognitif (chunking boundary)
// seperti <fieldset>, Stepper, Tabs, atau Wizard yang memuat > 9 field total atau > 7 field per-chunk.
type MonolithicFormBloatRule struct{}

// NewMonolithicFormBloatRule membuat instance baru dari MonolithicFormBloatRule.
func NewMonolithicFormBloatRule() *MonolithicFormBloatRule {
	return &MonolithicFormBloatRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *MonolithicFormBloatRule) ID() string {
	return "ux.monolithic-form-bloat"
}

// Description mengembalikan ringkasan aturan.
func (r *MonolithicFormBloatRule) Description() string {
	return "Warns when a monolithic form contains excessive unchunked inputs (> 9 total or > 7 per chunk), violating Cognitive Load Theory"
}

// Category mengembalikan nama kategori rule.
func (r *MonolithicFormBloatRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *MonolithicFormBloatRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *MonolithicFormBloatRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Sweller's Cognitive Load Theory (Intrinsic & Germane Cognitive Load Management)",
			"Progressive Task Completion & Form Usability Principles (Wroblewski & Nielsen Norman Group)",
			"W3C WAI-ARIA Authoring Practices Guide (Form Landmark & Fieldset Segmentation)",
		},
		CoreInvariant: "Forms containing more than 9 total interactive fields must segment fields into chunks ('<fieldset>', Stepper, or Tabs), with no single chunk exceeding 7 fields.",
		Grounding: "Long, monolithic forms containing 10 or more unorganized inputs overload user working memory and generate psychological intimidation. " +
			"According to Cognitive Load Theory, breaking complex information into manageable, cohesive chunks reduces task completion friction and error rates.\n\n" +
			"Forms must segment large input groups into semantic '<fieldset>' elements with explanatory '<legend>' titles, " +
			"or utilize multi-step wizards ('<Stepper>', progressive tabs). Furthermore, each individual chunk must not exceed 7 active inputs to preserve perceptual focus.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Severe Form Abandonment",
				Severity: "MEDIUM",
				Impact:   "Users faced with tall, endless walls of inputs drop off significantly before completing registration or checkout flows.",
			},
			{
				Vector:   "Field Omission & Data Entry Fatigue",
				Severity: "MEDIUM",
				Impact:   "Dense unchunked inputs cause users to overlook required fields, leading to repeated validation errors upon submission.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Monolithic form with 10 flat interactive inputs without fieldset or step chunking",
				Code: `<form onSubmit={handleSubmit} className="space-y-4 max-w-md">
  <input name="f1" placeholder="First Name" />
  <input name="f2" placeholder="Last Name" />
  <input name="f3" placeholder="Email" />
  <input name="f4" placeholder="Phone" />
  <input name="f5" placeholder="Address" />
  <input name="f6" placeholder="City" />
  <input name="f7" placeholder="State" />
  <input name="f8" placeholder="Zip Code" />
  <input name="f9" placeholder="Company" />
  <input name="f10" placeholder="Job Title" />
  <button type="submit">Kirim</button>
</form>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Segmented into two semantic fieldsets with legends, keeping fields per chunk under 7",
				Code: `<form onSubmit={handleSubmit} className="space-y-6 max-w-md">
  <fieldset className="space-y-4 border p-4 rounded-lg">
    <legend className="font-semibold text-sm">Informasi Pribadi</legend>
    <input name="f1" placeholder="First Name" />
    <input name="f2" placeholder="Last Name" />
    <input name="f3" placeholder="Email" />
    <input name="f4" placeholder="Phone" />
  </fieldset>

  <fieldset className="space-y-4 border p-4 rounded-lg">
    <legend className="font-semibold text-sm">Alamat Pengiriman</legend>
    <input name="f5" placeholder="Address" />
    <input name="f6" placeholder="City" />
    <input name="f7" placeholder="State" />
    <input name="f8" placeholder="Zip Code" />
  </fieldset>

  <button type="submit">Kirim</button>
</form>`,
			},
		},
	}
}

// Evaluate memeriksa apakah form memuat terlalu banyak field tanpa pemisahan chunking.
func (r *MonolithicFormBloatRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}

	tagLower := strings.ToLower(node.Tag)
	isForm := tagLower == "form"
	if !isForm {
		if role, ok := getAttrCaseInsensitive(node, "role"); ok && strings.ToLower(role) == "form" {
			isForm = true
		}
	}
	if !isForm {
		return nil
	}

	chunks := findTopLevelFormChunks(node)

	// Kasus 1: Tidak ada mekanisme chunking sama sekali
	if len(chunks) == 0 {
		totalFields := countInteractiveFieldsInSubtree(node, nil)
		if totalFields > 9 {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message: fmt.Sprintf(
						"Monolithic form contains %d interactive fields without chunking boundaries (<fieldset>, Stepper, or multi-step tabs), violating Cognitive Load Theory (maximum 9 unstructured fields).",
						totalFields,
					),
					Hint: "Segment fields into semantic <fieldset> groups with <legend> titles, or split into a multi-step wizard.",
				},
			}
		}
		return nil
	}

	// Kasus 2: Ada chunking, periksa kapasitas per-chunk (> 7 field)
	var diags []ir.Diagnostic
	for _, chunk := range chunks {
		chunkFields := countInteractiveFieldsInSubtree(chunk, nil)
		if chunkFields > 7 {
			diags = append(diags, ir.Diagnostic{
				Line:     chunk.Span.Line,
				Column:   chunk.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"Form section/chunk contains %d interactive fields exceeding recommended cognitive chunk capacity (maximum 7 fields per chunk).",
					chunkFields,
				),
				Hint: "Split this section into smaller logical fieldsets or progressive disclosure groups.",
			})
		}
	}

	// Periksa juga field yang berada di luar seluruh chunk
	unchunkedFields := countFieldsOutsideChunks(node, chunks)
	if unchunkedFields > 7 {
		diags = append(diags, ir.Diagnostic{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message: fmt.Sprintf(
				"Form contains %d unchunked fields outside of section containers (maximum 7 unchunked fields allowed).",
				unchunkedFields,
			),
			Hint: "Enclose remaining free fields within designated fieldsets or move to an adjacent step.",
		})
	}

	return diags
}

func findTopLevelFormChunks(form *ir.Node) []*ir.Node {
	var chunks []*ir.Node

	var collect func(curr *ir.Node)
	collect = func(curr *ir.Node) {
		if curr == nil {
			return
		}
		for _, child := range curr.Children {
			if child.Type != ir.NodeElement {
				continue
			}
			if isFormChunkingContainer(child) {
				chunks = append(chunks, child)
				// Jangan masuk ke anak chunk ini agar tidak menduplikasi top-level chunk
			} else {
				collect(child)
			}
		}
	}

	collect(form)
	return chunks
}

func countInteractiveFieldsInSubtree(root *ir.Node, excludeSubtrees []*ir.Node) int {
	if root == nil {
		return 0
	}

	radioGroups := make(map[string]bool)
	count := 0

	for n := range root.Walk() {
		if n == root {
			continue
		}
		if isNodeInSubtrees(n, excludeSubtrees) {
			continue
		}

		isField, groupKey := isInteractiveFormField(n)
		if isField {
			if groupKey != "" {
				if !radioGroups[groupKey] {
					radioGroups[groupKey] = true
					count++
				}
			} else {
				count++
			}
		}
	}

	return count
}

func isNodeInSubtrees(node *ir.Node, subtrees []*ir.Node) bool {
	if len(subtrees) == 0 || node == nil {
		return false
	}
	curr := node
	for curr != nil {
		for _, sub := range subtrees {
			if curr == sub {
				return true
			}
		}
		curr = curr.Parent
	}
	return false
}

func countFieldsOutsideChunks(form *ir.Node, chunks []*ir.Node) int {
	return countInteractiveFieldsInSubtree(form, chunks)
}
