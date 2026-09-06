package ux

import (
	"fmt"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// RadioOverchoiceRule mendeteksi kumpulan opsi radio datar yang berlebihan (> 7 opsi)
// tanpa sarana filter/pencarian atau abstraksi combobox/select, melanggar Hukum Hick-Hyman.
type RadioOverchoiceRule struct{}

// NewRadioOverchoiceRule membuat instance baru dari RadioOverchoiceRule.
func NewRadioOverchoiceRule() *RadioOverchoiceRule {
	return &RadioOverchoiceRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *RadioOverchoiceRule) ID() string {
	return "ux.radio-overchoice"
}

// Description mengembalikan ringkasan aturan.
func (r *RadioOverchoiceRule) Description() string {
	return "Warns when radio groups present excessive flat options (> 7) without filtering or combobox grouping, violating Hick-Hyman Law"
}

// Category mengembalikan nama kategori rule.
func (r *RadioOverchoiceRule) Category() string {
	return "ux"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (warn).
func (r *RadioOverchoiceRule) DefaultSeverity() ir.Severity {
	return ir.SeverityWarn
}

// Doc mengembalikan dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *RadioOverchoiceRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"Hick-Hyman Law of Decision Latency (Reaction Time T = b * log2(n + 1))",
			"W3C WAI-ARIA Authoring Practices Guide 1.2 (Radio Group Design Pattern)",
			"Nielsen Norman Group Guidelines on Selection Controls (Radio Buttons vs Dropdown Menus)",
		},
		CoreInvariant: "Radio groups sharing the same name or contained within a '<RadioGroup>' must not present more than 7 flat unsearchable options without filter mechanisms or combobox grouping.",
		Grounding: "The Hick-Hyman Law mathematically models cognitive choice reaction time as a logarithmic function of the number of options presented. " +
			"Radio buttons are optimized for rapid, mutually exclusive scanning when choices are few (2 to 4 options).\n\n" +
			"When developers present 8 or more flat radio buttons (e.g. selecting from 34 provinces or 50 states), " +
			"users must visually inspect every single choice sequentially. This drastically inflates interaction latency and induces decision paralysis. " +
			"For option counts exceeding 7, a searchable '<Combobox>' or grouped dropdown '<Select>' is strongly mandated.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Decision Paralysis & Extended Scan Latency",
				Severity: "MEDIUM",
				Impact:   "Users experience substantial friction locating their desired choice among dozens of unstructured radio options.",
			},
			{
				Vector:   "Excessive Vertical Viewport Consumption",
				Severity: "MEDIUM",
				Impact:   "Long lists of vertical radio items force extensive scrolling on mobile screens, pushing submit buttons out of view.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Ten flat radio buttons in a group without filtering or select abstraction",
				Code: `<div className="space-y-2">
  <label className="text-sm font-semibold">Pilih Wilayah Kerja</label>
  <input type="radio" name="region" value="reg-1" /> Wilayah 1
  <input type="radio" name="region" value="reg-2" /> Wilayah 2
  <input type="radio" name="region" value="reg-3" /> Wilayah 3
  <input type="radio" name="region" value="reg-4" /> Wilayah 4
  <input type="radio" name="region" value="reg-5" /> Wilayah 5
  <input type="radio" name="region" value="reg-6" /> Wilayah 6
  <input type="radio" name="region" value="reg-7" /> Wilayah 7
  <input type="radio" name="region" value="reg-8" /> Wilayah 8
  <input type="radio" name="region" value="reg-9" /> Wilayah 9
  <input type="radio" name="region" value="reg-10" /> Wilayah 10
</div>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Searchable Combobox for large dataset, keeping cognitive load low",
				Code: `<div className="space-y-2">
  <label className="text-sm font-semibold">Pilih Wilayah Kerja</label>
  <Combobox
    options={regionOptions}
    placeholder="Cari atau pilih wilayah..."
    searchable
  />
</div>`,
			},
			{
				Language: "tsx",
				Comment:  "Compact radio group with 3 clear, distinct choices adhering to Hick-Hyman Law",
				Code: `<RadioGroup name="billing_cycle" className="flex gap-4">
  <RadioGroupItem value="monthly" label="Bulanan" />
  <RadioGroupItem value="annual" label="Tahunan (Hemat 20%)" />
  <RadioGroupItem value="lifetime" label="Seumur Hidup" />
</RadioGroup>`,
			},
		},
	}
}

// Evaluate memeriksa apakah grup radio memiliki opsi berlebih.
func (r *RadioOverchoiceRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil || node.Type != ir.NodeElement {
		return nil
	}
	if isRadioGroupContainer(node) {
		return r.evaluateRadioGroup(node)
	}
	if isRadioInput(node) {
		return r.evaluateRadioInput(node)
	}
	return nil
}

func (r *RadioOverchoiceRule) evaluateRadioGroup(node *ir.Node) []ir.Diagnostic {
	count := countRadioGroupOptions(node)
	if count <= 4 || hasFilterOrSearchInput(node) {
		return nil
	}

	name, _ := getAttrCaseInsensitive(node, "name")
	nameClean := cleanAttrValue(name)
	if nameClean == "" {
		nameClean = "options"
	}

	return r.buildOverchoiceDiagnostic(node, nameClean, count)
}

func (r *RadioOverchoiceRule) evaluateRadioInput(node *ir.Node) []ir.Diagnostic {
	if isInsideRadioGroupContainer(node) {
		return nil
	}

	nameVal, ok := getAttrCaseInsensitive(node, "name")
	nameClean := cleanAttrValue(nameVal)
	if !ok || nameClean == "" {
		return nil
	}

	scope := findRadioScope(node)
	if scope == nil {
		return nil
	}

	isFirst, totalCount := analyzeRadioInScope(scope, nameClean, node)
	if !isFirst || totalCount <= 4 || hasFilterOrSearchInput(scope) {
		return nil
	}

	return r.buildOverchoiceDiagnostic(node, nameClean, totalCount)
}

func (r *RadioOverchoiceRule) buildOverchoiceDiagnostic(node *ir.Node, name string, count int) []ir.Diagnostic {
	if count > 7 {
		return []ir.Diagnostic{
			{
				Line:     node.Span.Line,
				Column:   node.Span.Column,
				Rule:     r.ID(),
				Severity: r.DefaultSeverity(),
				Message: fmt.Sprintf(
					"Radio group %q presents %d flat options exceeding Hick-Hyman threshold (maximum 7). Users experience elevated decision latency.",
					name, count,
				),
				Hint: "Consider replacing with a searchable Combobox or Select dropdown to reduce cognitive scan overhead.",
			},
		}
	}

	// 5 <= count <= 7: Advisory
	return []ir.Diagnostic{
		{
			Line:     node.Span.Line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: ir.SeverityInfo,
			Message: fmt.Sprintf(
				"Radio group %q presents %d flat options. Hick-Hyman Law recommends considering a Select dropdown when choices exceed 4.",
				name, count,
			),
			Hint: "Consider using a Select dropdown if options are not immediately compared side-by-side.",
		},
	}
}

func isRadioGroupContainer(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	tagLower := strings.ToLower(node.Tag)
	if tagLower == "radiogroup" || strings.HasSuffix(tagLower, "radiogroup") {
		return true
	}
	if role, ok := getAttrCaseInsensitive(node, "role"); ok && strings.ToLower(role) == "radiogroup" {
		return true
	}
	return false
}

func isInsideRadioGroupContainer(node *ir.Node) bool {
	curr := node.Parent
	for curr != nil {
		if isRadioGroupContainer(curr) {
			return true
		}
		curr = curr.Parent
	}
	return false
}

func isRadioInput(node *ir.Node) bool {
	if node == nil || node.Type != ir.NodeElement {
		return false
	}
	tagLower := strings.ToLower(node.Tag)
	if tagLower != "input" && node.Tag != "Input" {
		return false
	}
	typeVal, ok := getAttrCaseInsensitive(node, "type")
	return ok && cleanAttrValue(typeVal) == "radio"
}

func countRadioGroupOptions(node *ir.Node) int {
	count := 0
	for child := range node.Walk() {
		if child == node || child.Type != ir.NodeElement {
			continue
		}
		tagLower := strings.ToLower(child.Tag)
		if tagLower == "radiogroupitem" || strings.HasSuffix(tagLower, "radiogroupitem") ||
			tagLower == "radio" || child.Tag == "Radio" || isRadioInput(child) {
			count++
		}
	}
	return count
}

func findRadioScope(node *ir.Node) *ir.Node {
	curr := node.Parent
	lastValid := node.Parent
	for curr != nil {
		tagLower := strings.ToLower(curr.Tag)
		if tagLower == "form" || tagLower == "fieldset" {
			return curr
		}
		lastValid = curr
		curr = curr.Parent
	}
	return lastValid
}

func analyzeRadioInScope(scope *ir.Node, name string, target *ir.Node) (isFirst bool, count int) {
	var firstRadio *ir.Node
	for n := range scope.Walk() {
		if !isRadioInput(n) {
			continue
		}
		nameVal, ok := getAttrCaseInsensitive(n, "name")
		if ok && cleanAttrValue(nameVal) == name {
			if firstRadio == nil {
				firstRadio = n
			}
			count++
		}
	}
	return firstRadio == target, count
}

func hasFilterOrSearchInput(node *ir.Node) bool {
	if node == nil {
		return false
	}
	if hasSearchableAttribute(node) {
		return true
	}
	for child := range node.Walk() {
		if child.Type == ir.NodeElement && isSearchOrFilterNode(child) {
			return true
		}
	}
	return false
}

func hasSearchableAttribute(node *ir.Node) bool {
	_, ok := getAttrCaseInsensitive(node, "searchable", "filterable", "filter")
	return ok
}

func isSearchOrFilterNode(child *ir.Node) bool {
	tagLower := strings.ToLower(child.Tag)
	if strings.Contains(tagLower, "search") || strings.Contains(tagLower, "filter") {
		return true
	}
	if tagLower == "input" || child.Tag == "Input" {
		typeVal, _ := getAttrCaseInsensitive(child, "type")
		if cleanAttrValue(typeVal) == "search" {
			return true
		}
		for _, key := range [...]string{"placeholder", "aria-label", "name", "id"} {
			if val, ok := getAttrCaseInsensitive(child, key); ok {
				valLower := strings.ToLower(val)
				if strings.Contains(valLower, "search") || strings.Contains(valLower, "cari") || strings.Contains(valLower, "filter") {
					return true
				}
			}
		}
	}
	return false
}
