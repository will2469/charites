package browser

import (
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// DateInputFormatAssumptionRule mendeteksi asumsi format lokal saat parsing nilai input tanggal native.
type DateInputFormatAssumptionRule struct{}

// NewDateInputFormatAssumptionRule membuat instance baru dari DateInputFormatAssumptionRule.
func NewDateInputFormatAssumptionRule() *DateInputFormatAssumptionRule {
	return &DateInputFormatAssumptionRule{}
}

// ID mengembalikan identifier kanonikal rule.
func (r *DateInputFormatAssumptionRule) ID() string {
	return "browser.date-input-format-assumption"
}

// Description mengembalikan ringkasan aturan.
func (r *DateInputFormatAssumptionRule) Description() string {
	return "Prohibits localized string splitting assumptions on HTML5 date input values in favor of normative ISO 8601 parsing"
}

// Category mengembalikan nama kategori rule.
func (r *DateInputFormatAssumptionRule) Category() string {
	return "browser"
}

// DefaultSeverity mengembalikan tingkat keparahan bawaan (error).
func (r *DateInputFormatAssumptionRule) DefaultSeverity() ir.Severity {
	return ir.SeverityError
}

// Doc mengembalikan spesifikasi dokumentasi 8-Pillars lengkap untuk generator wiki.
func (r *DateInputFormatAssumptionRule) Doc() ir.RuleDocumentation {
	return ir.RuleDocumentation{
		TargetStandards: []string{
			"HTML Living Standard Section 4.10.5.1.7 (Date State - type=date)",
			"W3C RFC 3339 / ISO 8601 Normative Date Representation (YYYY-MM-DD)",
		},
		CoreInvariant: "Native <input type=\"date\"> values are guaranteed by W3C specification to be serialized strictly as ISO 8601 (YYYY-MM-DD). Code must not split values by localized delimiters ('/' or '.').",
		Grounding: "While the browser UI may render localized date pickers according to OS settings (e.g., DD/MM/YYYY in Indonesia/UK, MM/DD/YYYY in US), " +
			"the programmatic 'element.value' is ALWAYS serialized in ISO 8601 format ('YYYY-MM-DD').\n\n" +
			"Splitting by '/' causes catastrophic silent failures because the delimiter never exists in 'element.value'.",
		Risks: []ir.RiskItem{
			{
				Vector:   "Silent Date Parsing Failure",
				Severity: "HIGH",
				Impact:   "Splitting ISO 8601 date string by '/' returns an array of length 1, corrupting day, month, and year data sent to APIs.",
			},
			{
				Vector:   "Cross-Locale Form Submission Corruption",
				Severity: "HIGH",
				Impact:   "User birth dates, appointment bookings, or legal document dates are stored as NaN, undefined, or incorrect epochs.",
			},
		},
		BadExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Splitting native date input value by '/' based on UI display assumption",
				Code: `<input
  type="date"
  onChange={(e) => {
    // BUG: e.target.value is '2026-09-06'. Splitting by '/' fails!
    const [day, month, year] = e.target.value.split('/');
    saveDate(day, month, year);
  }}
/>`,
			},
		},
		GoodExamples: []ir.CodeExample{
			{
				Language: "tsx",
				Comment:  "Splitting native date input by normative ISO 8601 dash delimiter or using valueAsDate",
				Code: `<input
  type="date"
  onChange={(e) => {
    // Correct: ISO 8601 format YYYY-MM-DD
    const [year, month, day] = e.target.value.split('-');
    saveDate(day, month, year);
  }}
/>`,
			},
		},
	}
}

// Evaluate memeriksa apakah terdapat pemecahan nilai tanggal dengan delimiter lokal non-standar.
func (r *DateInputFormatAssumptionRule) Evaluate(node *ir.Node) []ir.Diagnostic {
	if node == nil {
		return nil
	}

	tag := strings.ToLower(node.Tag)
	if tag == "input" {
		return r.evalDateInputNode(node)
	}
	if tag == "script" {
		return r.evalDateScriptNode(node)
	}

	return nil
}

func (r *DateInputFormatAssumptionRule) evalDateInputNode(node *ir.Node) []ir.Diagnostic {
	inputType := strings.ToLower(node.Attributes["type"])
	if inputType != "date" {
		return nil
	}

	for attrName, attrVal := range node.Attributes {
		if isScriptAttribute(attrName) && hasNonStandardDateSplit(attrVal) {
			return []ir.Diagnostic{
				{
					Line:     node.Span.Line,
					Column:   node.Span.Column,
					Rule:     r.ID(),
					Severity: r.DefaultSeverity(),
					Message:  "Native <input type=\"date\"> value is split by localized delimiter ('/' or '.'). HTML5 guarantees '.value' is always formatted in ISO 8601 (YYYY-MM-DD).",
					Hint:     "Split using '-' (e.g. value.split('-')) or use 'input.valueAsDate' / 'new Date(value)' instead.",
				},
			}
		}
	}
	return nil
}

func (r *DateInputFormatAssumptionRule) evalDateScriptNode(node *ir.Node) []ir.Diagnostic {
	scriptText := getStyleNodeText(node)
	if !hasDateInputContext(scriptText) || !hasNonStandardDateSplit(scriptText) {
		return nil
	}

	line := node.Span.Line
	lines := strings.Split(scriptText, "\n")
	for idx, l := range lines {
		if hasNonStandardDateSplit(l) {
			line = node.Span.Line + idx
			break
		}
	}

	return []ir.Diagnostic{
		{
			Line:     line,
			Column:   node.Span.Column,
			Rule:     r.ID(),
			Severity: r.DefaultSeverity(),
			Message:  "Date input value is split by localized delimiter ('/' or '.'). HTML5 guarantees '.value' is always formatted in ISO 8601 (YYYY-MM-DD).",
			Hint:     "Split using '-' (e.g. value.split('-')) or use 'input.valueAsDate' / 'new Date(value)' instead.",
		},
	}
}

func hasNonStandardDateSplit(s string) bool {
	lower := strings.ToLower(s)
	if !strings.Contains(lower, "split(") {
		return false
	}

	// Cek pemisahan dengan '/' atau '.'
	return strings.Contains(lower, ".split('/')") ||
		strings.Contains(lower, ".split(\"/\")") ||
		strings.Contains(lower, ".split(`/`)") ||
		strings.Contains(lower, ".split('.')") ||
		strings.Contains(lower, ".split(\".\")") ||
		strings.Contains(lower, ".split(/\\//)")
}

func hasDateInputContext(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "type=\"date\"") ||
		strings.Contains(lower, "type='date'") ||
		strings.Contains(lower, "type=`date`") ||
		strings.Contains(lower, "dateinput") ||
		strings.Contains(lower, "date_input") ||
		strings.Contains(lower, "birth-date") ||
		strings.Contains(lower, "birthdate") ||
		strings.Contains(lower, "input[type=date]") ||
		strings.Contains(lower, "date")
}
