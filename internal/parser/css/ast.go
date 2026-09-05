package css

// Rule adalah interface penanda seluruh node aturan tingkat tinggi dalam stylesheet CSS.
type Rule interface {
	ruleNode()
	GetSpan() SourceSpan
}

// Declaration merepresentasikan pasangan nama properti dan nilai (misal: "--banana: #123456").
type Declaration struct {
	Property string
	Value    string
	RawValue string
	Span     SourceSpan
}

// AtRule merepresentasikan aturan berawalan @ (seperti @layer, @media, @supports, @theme).
type AtRule struct {
	Name         string        // e.g. "@layer", "@media", "@theme"
	Prelude      string        // e.g. "theme", "(prefers-color-scheme: dark)"
	Declarations []Declaration // Deklarasi langsung di dalam blok
	Rules        []Rule        // Aturan bersarang (nested rules)
	Span         SourceSpan
}

func (a *AtRule) ruleNode() {}

// GetSpan mengembalikan rentang posisi sumber (SourceSpan) untuk AtRule.
func (a *AtRule) GetSpan() SourceSpan { return a.Span }

// StyleRule merepresentasikan blok aturan dengan selektor CSS (seperti :root, .dark, .card).
type StyleRule struct {
	Selector     string
	Declarations []Declaration // Deklarasi properti CSS (--color-*, color-scheme, dsb.)
	Rules        []Rule        // Aturan bersarang (CSS nesting seperti &:hover, .child)
	Span         SourceSpan
}

func (s *StyleRule) ruleNode() {}

// GetSpan mengembalikan rentang posisi sumber (SourceSpan) untuk StyleRule.
func (s *StyleRule) GetSpan() SourceSpan { return s.Span }

// StyleSheet merepresentasikan berkas CSS lengkap yang telah diparsing.
type StyleSheet struct {
	Rules        []Rule
	Declarations []Declaration
	Span         SourceSpan
}
