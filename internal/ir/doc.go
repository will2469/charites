package ir

// CodeExample merepresentasikan contoh kode untuk dokumentasi rule.
type CodeExample struct {
	Language string
	Code     string
	Comment  string
}

// RiskItem merepresentasikan taksonomi risiko jika rule dilanggar.
type RiskItem struct {
	Vector   string
	Severity string
	Impact   string
}

// RuleDocumentation menyimpan metadata dokumentasi lengkap 8-Pillars untuk suatu rule.
type RuleDocumentation struct {
	TargetStandards []string
	CoreInvariant   string
	Grounding       string
	BadExamples     []CodeExample
	GoodExamples    []CodeExample
	Risks           []RiskItem
}
