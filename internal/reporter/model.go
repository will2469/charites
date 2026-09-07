package reporter

// jsonViolationItem merepresentasikan satu pelanggaran aturan dalam konteks berkas tertentu.
type jsonViolationItem struct {
	Line        int    `json:"line"`
	Column      int    `json:"column"`
	Rule        string `json:"rule"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	Hint        string `json:"hint,omitempty"`
	DocURL      string `json:"doc_url"`
	Suppression string `json:"suppression,omitempty"`
}

// jsonFileGroup merepresentasikan kumpulan pelanggaran untuk satu berkas tertentu.
type jsonFileGroup struct {
	File         string              `json:"file"`
	ErrorCount   int                 `json:"error_count"`
	WarningCount int                 `json:"warning_count"`
	InfoCount    int                 `json:"info_count"`
	TotalIssues  int                 `json:"total_issues"`
	Violations   []jsonViolationItem `json:"violations"`
}

// jsonDiagnostic merepresentasikan satu titik diagnostik flat stream untuk backward compatibility.
type jsonDiagnostic struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Rule     string `json:"rule"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Hint     string `json:"hint,omitempty"`
	DocURL   string `json:"doc_url"`
}

// jsonDocument merepresentasikan struktur dokumen JSON lengkap hasil pemindaian.
type jsonDocument struct {
	Version     string           `json:"version"`
	Summary     ScanSummary      `json:"summary"`
	Files       []jsonFileGroup  `json:"files"`
	Diagnostics []jsonDiagnostic `json:"diagnostics"`
}
