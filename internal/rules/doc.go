package rules

import (
	"github.com/will2469/charites/internal/ir"
)

// CodeExample merepresentasikan contoh kode untuk dokumentasi rule.
type CodeExample = ir.CodeExample

// RiskItem merepresentasikan taksonomi risiko jika rule dilanggar.
type RiskItem = ir.RiskItem

// RuleDocumentation menyimpan metadata dokumentasi lengkap 8-Pillars untuk suatu rule.
type RuleDocumentation = ir.RuleDocumentation

// DocumentedRule adalah interface opsional yang diimplementasikan oleh rule yang menyediakan dokumentasi kaya.
type DocumentedRule interface {
	Rule
	Doc() ir.RuleDocumentation
}
