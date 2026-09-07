package drift

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"

	"github.com/will2469/charites/internal/ir"
)

// ContrastHazard menyimpan temuan pelanggaran kontras aksesibilitas WCAG 1.4.3 pada elemen.
type ContrastHazard struct {
	FilePath  string
	Span      ir.Span
	Scope     ScopeKey
	BgToken   string
	TextToken string
	Ratio     float64
	Message   string
	Hint      string
}

// Report menyimpan seluruh ringkasan audit style drift geometris dan kromatik.
type Report struct {
	TotalOccurrences  int
	GeometricClusters map[ClusterKey]*Cluster
	ChromaticClusters []ChromaticCluster
	ContrastHazards   []ContrastHazard
	Options           Options
}

// Analyzer mengoordinasikan pipeline analisis style drift repositori.
type Analyzer struct {
	Options Options
}

// NewAnalyzer menginisialisasi Analyzer baru dengan opsi yang diberikan.
func NewAnalyzer(opts Options) *Analyzer {
	return &Analyzer{
		Options: opts,
	}
}

// Analyze mengevaluasi daftar kemunculan gaya repositori dan menghasilkan Report lengkap.
func (a *Analyzer) Analyze(occurrences []StyleOccurrence) Report {
	geomClusters := ClusterOccurrences(occurrences, a.Options)
	chromClusters := AnalyzeChromaticDrift(occurrences)
	hazards := a.detectContrastHazards(occurrences)

	return Report{
		TotalOccurrences:  len(occurrences),
		GeometricClusters: geomClusters,
		ChromaticClusters: chromClusters,
		ContrastHazards:   hazards,
		Options:           a.Options,
	}
}

// detectContrastHazards memindai elemen yang menggunakan kombinasi background dan text primitif untuk memeriksa rasio kontras.
func (a *Analyzer) detectContrastHazards(occurrences []StyleOccurrence) []ContrastHazard {
	type elementLoc struct {
		filePath string
		line     int
		column   int
	}

	seen := make(map[elementLoc]bool)
	var hazards []ContrastHazard

	for _, occ := range occurrences {
		loc := elementLoc{filePath: occ.FilePath, line: occ.Span.Line, column: occ.Span.Column}
		if seen[loc] {
			continue
		}

		bgClass, textClass := extractBgAndTextClasses(occ.RawClasses)
		if bgClass == "" || textClass == "" {
			continue
		}

		ratio, failsAA, ok := CheckContrastHazard(bgClass, textClass)
		if ok && failsAA {
			seen[loc] = true
			hazards = append(hazards, ContrastHazard{
				FilePath:  occ.FilePath,
				Span:      occ.Span,
				Scope:     occ.Scope,
				BgToken:   bgClass,
				TextToken: textClass,
				Ratio:     ratio,
				Message: fmt.Sprintf(
					"WCAG 1.4.3 Contrast Hazard: '%s' with '%s' yields contrast ratio %.2f:1 (fails minimum 4.5:1).",
					bgClass,
					textClass,
					ratio,
				),
				Hint: "Use a dark text foreground token (e.g. text-black) or define a high-contrast semantic token in global.css.",
			})
		}
	}

	return hazards
}

func extractBgAndTextClasses(rawClasses string) (string, string) {
	var bgClass, textClass string
	for _, c := range strings.Fields(rawClasses) {
		if strings.HasPrefix(c, "bg-") && !strings.Contains(c, ":") {
			bgClass = c
		}
		if strings.HasPrefix(c, "text-") && !strings.Contains(c, ":") {
			textClass = c
		}
	}
	return bgClass, textClass
}

// RenderReport merender Report ke io.Writer sesuai format yang dipilih (inline ANSI, json, atau markdown).
func RenderReport(w io.Writer, report *Report, format string, noColor bool) error {
	formatLower := strings.ToLower(strings.TrimSpace(format))

	switch formatLower {
	case "json":
		return renderJSONReport(w, report)
	case "markdown", "md":
		return renderMarkdownReport(w, report)
	default:
		return renderInlineReport(w, report, noColor)
	}
}

func renderJSONReport(w io.Writer, report *Report) error {
	type jsonOutput struct {
		TotalOccurrences int             `json:"total_occurrences"`
		ClustersCount    int             `json:"clusters_count"`
		ChromaticCount   int             `json:"chromatic_drift_count"`
		HazardsCount     int             `json:"contrast_hazards_count"`
		Diagnostics      []ir.Diagnostic `json:"diagnostics"`
	}

	out := jsonOutput{
		TotalOccurrences: report.TotalOccurrences,
		ClustersCount:    len(report.GeometricClusters),
		ChromaticCount:   len(report.ChromaticClusters),
		HazardsCount:     len(report.ContrastHazards),
		Diagnostics:      SynthesizeDiagnostics(*report),
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func renderMarkdownReport(w io.Writer, report *Report) error {
	var sb strings.Builder
	sb.WriteString("# Charites Component-Scoped Style Drift Report\n\n")
	sb.WriteString(fmt.Sprintf("**Total Style Occurrences:** %d\n\n", report.TotalOccurrences))

	renderMarkdownGeometric(&sb, report)
	renderMarkdownChromatic(&sb, report)
	renderMarkdownHazards(&sb, report)

	_, err := io.WriteString(w, sb.String())
	return err
}

func renderMarkdownGeometric(sb *strings.Builder, report *Report) {
	sb.WriteString("## 1. Geometric & Micro-Interaction Drift Clusters\n\n")
	sortedKeys := getSortedClusterKeys(report.GeometricClusters)

	for _, k := range sortedKeys {
		c := report.GeometricClusters[k]
		if c.InsufficientData {
			continue
		}

		sb.WriteString(fmt.Sprintf("### <%s> - %s (%d occurrences)\n\n", k.Scope.String(), k.Category.String(), c.Total))
		sb.WriteString("| Token | Occurrences | Percentage | Status |\n")
		sb.WriteString("| :--- | :--- | :--- | :--- |\n")

		for _, stat := range c.TokenStatsList() {
			status := "Standard"
			if stat.Token == c.CanonicalToken {
				status = " **Canonical**"
			} else {
				for _, out := range c.Outliers {
					if out.Token == stat.Token {
						status = " **Rogue Outlier**"
						break
					}
				}
			}
			sb.WriteString(fmt.Sprintf("| `%s` | %d | %.1f%% | %s |\n", stat.Token, stat.Count, stat.Percentage, status))
		}
		sb.WriteString("\n")
	}
}

func renderMarkdownChromatic(sb *strings.Builder, report *Report) {
	if len(report.ChromaticClusters) == 0 {
		return
	}
	sb.WriteString("## 2. Chromatic Family & Variant Drift\n\n")
	for _, chrom := range report.ChromaticClusters {
		sb.WriteString(fmt.Sprintf("### <%s> -> %s (Candidate Intent: %s)\n\n",
			chrom.Scope.String(), chrom.Family.String(), chrom.CandidateIntent.String()))
		sb.WriteString(fmt.Sprintf(">  **Recommendation:** %s\n\n", chrom.Recommendation))
		sb.WriteString("| Hue | Occurrences | Percentage |\n| :--- | :--- | :--- |\n")
		for hue, count := range chrom.HueCounts {
			pct := float64(count) / float64(chrom.Total) * 100
			sb.WriteString(fmt.Sprintf("| `%s` | %d | %.1f%% |\n", hue, count, pct))
		}
		sb.WriteString("\n")
	}
}

func renderMarkdownHazards(sb *strings.Builder, report *Report) {
	if len(report.ContrastHazards) == 0 {
		return
	}
	sb.WriteString("## 3. WCAG 1.4.3 Contrast Hazards\n\n")
	sb.WriteString("| Location | Component | Background | Text | Contrast Ratio |\n")
	sb.WriteString("| :--- | :--- | :--- | :--- | :--- |\n")
	for _, h := range report.ContrastHazards {
		sb.WriteString(fmt.Sprintf("| `%s:%d` | `<%s>` | `%s` | `%s` | **%.2f:1** (Fails AA) |\n",
			h.FilePath, h.Span.Line, h.Scope.Kind.String(), h.BgToken, h.TextToken, h.Ratio))
	}
	sb.WriteString("\n")
}

type reportWriter struct {
	w   io.Writer
	err error
}

func (rw *reportWriter) printf(format string, a ...any) {
	if rw.err != nil {
		return
	}
	_, rw.err = fmt.Fprintf(rw.w, format, a...)
}

func (rw *reportWriter) println(a ...any) {
	if rw.err != nil {
		return
	}
	_, rw.err = fmt.Fprintln(rw.w, a...)
}

type colorPalette struct {
	bold   string
	dim    string
	cyan   string
	green  string
	yellow string
	red    string
	reset  string
}

func newColorPalette(noColor bool) colorPalette {
	if noColor {
		return colorPalette{}
	}
	return colorPalette{
		bold:   "\033[1m",
		dim:    "\033[2m",
		cyan:   "\033[36m",
		green:  "\033[32m",
		yellow: "\033[33m",
		red:    "\033[31m",
		reset:  "\033[0m",
	}
}

func renderInlineReport(w io.Writer, report *Report, noColor bool) error {
	rw := &reportWriter{w: w}
	pal := newColorPalette(noColor)

	rw.printf("\n%s%sCharites Component-Scoped Style Drift Report%s\n", pal.bold, pal.cyan, pal.reset)
	rw.printf("%sTotal Analyzed Occurrences: %d%s\n\n", pal.dim, report.TotalOccurrences, pal.reset)

	renderInlineGeometric(rw, report, pal)
	renderInlineChromatic(rw, report, pal)
	renderInlineHazards(rw, report, pal)

	return rw.err
}

func renderInlineGeometric(rw *reportWriter, report *Report, pal colorPalette) {
	sortedKeys := getSortedClusterKeys(report.GeometricClusters)
	hasClusters := false

	for _, k := range sortedKeys {
		c := report.GeometricClusters[k]
		if c.InsufficientData {
			continue
		}
		hasClusters = true
		rw.printf("%s[%s -> %s]%s (Total: %d)\n", pal.bold, k.Scope.String(), k.Category.String(), pal.reset, c.Total)
		renderClusterTokens(rw, c, pal)
		rw.println()
	}

	if !hasClusters {
		rw.printf("%sNo style drift detected across components.%s\n\n", pal.green, pal.reset)
	}
}

func renderClusterTokens(rw *reportWriter, c *Cluster, pal colorPalette) {
	for _, stat := range c.TokenStatsList() {
		bar := renderBar(stat.Percentage)
		if stat.Token == c.CanonicalToken {
			rw.printf("  %s%-22s%s %4d (%5.1f%%) %s%s%s %s[CANONICAL]%s\n",
				pal.green, stat.Token, pal.reset, stat.Count, stat.Percentage, pal.green, bar, pal.reset, pal.bold, pal.reset)
			continue
		}

		isOutlier := false
		for _, out := range c.Outliers {
			if out.Token == stat.Token {
				isOutlier = true
				break
			}
		}

		if isOutlier {
			rw.printf("  %s%-22s%s %4d (%5.1f%%) %s%s%s %s[ ROGUE DRIFT]%s\n",
				pal.yellow, stat.Token, pal.reset, stat.Count, stat.Percentage, pal.yellow, bar, pal.reset, pal.red, pal.reset)
		} else {
			rw.printf("  %-22s %4d (%5.1f%%) %s%s%s\n",
				stat.Token, stat.Count, stat.Percentage, pal.dim, bar, pal.reset)
		}
	}
}

func renderInlineChromatic(rw *reportWriter, report *Report, pal colorPalette) {
	if len(report.ChromaticClusters) == 0 {
		return
	}
	rw.printf("%s%s=== Chromatic Family & Variant Drift ===%s\n\n", pal.bold, pal.yellow, pal.reset)
	for _, chrom := range report.ChromaticClusters {
		rw.printf("%s[CLUSTER: <%s> -> %s (Candidate: %s)]%s\n",
			pal.bold, chrom.Scope.String(), chrom.Family.String(), chrom.CandidateIntent.String(), pal.reset)
		rw.printf("    %sChromatic Drift Detected: %d competing hues for %s intent!%s\n",
			pal.yellow, len(chrom.HueCounts), chrom.CandidateIntent.String(), pal.reset)

		for hue, count := range chrom.HueCounts {
			pct := float64(count) / float64(chrom.Total) * 100
			rw.printf("      - %-16s ──► %2dx (%5.1f%%)\n", hue, count, pct)
		}
		rw.printf("\n   %sActionable Architectural Recommendation:%s\n      %s\n\n",
			pal.cyan, pal.reset, chrom.Recommendation)
	}
}

func renderInlineHazards(rw *reportWriter, report *Report, pal colorPalette) {
	if len(report.ContrastHazards) == 0 {
		return
	}
	rw.printf("%s%s=== WCAG 1.4.3 Contrast Hazards ===%s\n\n", pal.bold, pal.red, pal.reset)
	for _, h := range report.ContrastHazards {
		rw.printf("  %s  %s:%d%s <%s> %s + %s yields %.2f:1 (fails AA)\n",
			pal.red, h.FilePath, h.Span.Line, pal.reset, h.Scope.Kind.String(), h.BgToken, h.TextToken, h.Ratio)
	}
	rw.println()
}

func renderBar(percentage float64) string {
	const barWidth = 20
	fillLen := int(math.Round((percentage / 100.0) * float64(barWidth)))
	if fillLen < 0 {
		fillLen = 0
	}
	if fillLen > barWidth {
		fillLen = barWidth
	}
	return "[" + strings.Repeat("█", fillLen) + strings.Repeat("░", barWidth-fillLen) + "]"
}

func getSortedClusterKeys(clusters map[ClusterKey]*Cluster) []ClusterKey {
	keys := make([]ClusterKey, 0, len(clusters))
	for k := range clusters {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(a, b ClusterKey) int {
		if a.Scope.Kind != b.Scope.Kind {
			return int(a.Scope.Kind) - int(b.Scope.Kind)
		}
		if a.Scope.Confidence != b.Scope.Confidence {
			return int(a.Scope.Confidence) - int(b.Scope.Confidence)
		}
		return int(a.Category) - int(b.Category)
	})

	return keys
}
