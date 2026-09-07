package drift

import (
	"fmt"

	"github.com/will2469/charites/internal/ir"
)

// RuleID merepresentasikan Charites Rule ID kanonikal untuk style drift.
const RuleID = "design.component-style-drift"

// SynthesizeDiagnostics mengonversi seluruh temuan outlier dan chromatic drift ke slice ir.Diagnostic.
func SynthesizeDiagnostics(report Report) []ir.Diagnostic {
	diags := make([]ir.Diagnostic, 0, len(report.ContrastHazards))

	// 1. Diagnostik Geometris & Outlier
	for _, cluster := range report.GeometricClusters {
		if cluster.InsufficientData || len(cluster.Outliers) == 0 {
			continue
		}

		for _, outlier := range cluster.Outliers {
			for _, occ := range outlier.Occurrences {
				msg := fmt.Sprintf(
					"Rogue variation '%s' on <%s> (%d of %d occurrences, %.1f%%). Dominant standard is '%s' (%d of %d, %.1f%%).",
					outlier.Token,
					occ.Scope.Kind.String(),
					outlier.Count,
					cluster.Total,
					outlier.Percentage,
					cluster.CanonicalToken,
					cluster.CanonicalCount,
					cluster.Total,
					cluster.CanonicalDominance,
				)

				hint := fmt.Sprintf(
					"Consolidate to dominant standard '%s' or configure an exception in charites.yaml if intentional.",
					cluster.CanonicalToken,
				)

				diags = append(diags, ir.Diagnostic{
					File:     occ.FilePath,
					Line:     occ.Span.Line,
					Column:   occ.Span.Column,
					Rule:     RuleID,
					Severity: ir.SeverityWarn,
					Message:  msg,
					Hint:     hint,
				})
			}
		}
	}

	// 2. Diagnostik Chromatic Drift
	for _, chrom := range report.ChromaticClusters {
		for hue, count := range chrom.HueCounts {
			pct := float64(count) / float64(chrom.Total) * 100.0
			for _, occ := range chrom.HueOccurrences[hue] {
				msg := fmt.Sprintf(
					"Chromatic hue drift: '%s' on <%s> (%d of %d, %.1f%%) competes with other %s hues.",
					hue,
					chrom.Scope.Kind.String(),
					count,
					chrom.Total,
					pct,
					chrom.CandidateIntent.String(),
				)

				diags = append(diags, ir.Diagnostic{
					File:     occ.FilePath,
					Line:     occ.Span.Line,
					Column:   occ.Span.Column,
					Rule:     RuleID,
					Severity: ir.SeverityWarn,
					Message:  msg,
					Hint:     chrom.Recommendation,
				})
			}
		}
	}

	// 3. Diagnostik Contrast Hazard
	for _, hazard := range report.ContrastHazards {
		diags = append(diags, ir.Diagnostic{
			File:     hazard.FilePath,
			Line:     hazard.Span.Line,
			Column:   hazard.Span.Column,
			Rule:     RuleID,
			Severity: ir.SeverityWarn,
			Message:  hazard.Message,
			Hint:     hazard.Hint,
		})
	}

	return diags
}
