package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/scanner"
)

// RunDrift mengeksekusi pipeline analisis Component-Scoped Style Drift repositori.
func RunDrift(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("drift", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var format string
	var outputFile string
	var configPath string
	var noColor bool
	var dominanceThreshold float64
	var outlierThreshold float64
	var minOccurrences int
	var reportFlag bool

	fs.StringVar(&format, "format", "inline", "Format output: inline, json, atau markdown (alias: md)")
	fs.StringVar(&format, "f", "inline", "Format output (shorthand)")
	fs.StringVar(&outputFile, "output", "", "Path berkas untuk menyimpan laporan drift (misal: drift.md)")
	fs.StringVar(&outputFile, "o", "", "Path berkas laporan (shorthand)")
	fs.StringVar(&configPath, "config", "charites.yaml", "Path ke berkas konfigurasi")
	fs.BoolVar(&noColor, "no-color", false, "Matikan pewarnaan ANSI di terminal")
	fs.BoolVar(&reportFlag, "report", true, "Cetak laporan ringkasan visual cluster")
	fs.Float64Var(&dominanceThreshold, "dominance", 0, "Ambang batas dominansi kanonikal (persen)")
	fs.Float64Var(&outlierThreshold, "outlier", 0, "Ambang batas outlier rogue (persen)")
	fs.IntVar(&minOccurrences, "min-occurrences", 0, "Batas minimum populasi cluster")

	reorderedArgs, positionalArgs := partitionArgs(args)
	if err := fs.Parse(reorderedArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprint(stdout, UsageString())
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v. Run 'charites --help' for usage.\n", err)
		return ExitOperational
	}

	target := "."
	if len(positionalArgs) == 1 {
		target = positionalArgs[0]
	}

	isExplicitConfig := isFlagPassed(args, "-config", "--config")
	cfg, ok := resolveScanConfig(target, configPath, isExplicitConfig, stderr)
	if !ok {
		return ExitOperational
	}

	if !isFlagPassed(args, "-f", "--format") && cfg != nil && cfg.Format != "" {
		format = cfg.Format
	}
	if !isFlagPassed(args, "-o", "--output") && cfg != nil && cfg.Output != "" {
		outputFile = cfg.Output
	}

	format = resolveReportFormat(format, outputFile)
	if !validateScanTargetAndFormat(target, format, positionalArgs, stderr) {
		return ExitOperational
	}

	occurrences, ok := scanOccurrences(target, cfg, stderr)
	if !ok {
		return ExitOperational
	}

	opts := resolveDriftOptions(cfg, dominanceThreshold, outlierThreshold, minOccurrences)
	driftAnalyzer := drift.NewAnalyzer(opts)
	report := driftAnalyzer.Analyze(occurrences)

	return writeDriftReport(report, outputFile, format, noColor, reportFlag, stdout, stderr)
}

func scanOccurrences(target string, cfg *config.Config, stderr io.Writer) ([]drift.StyleOccurrence, bool) {
	matcher, ok := buildScanMatcher(target, cfg, nil, stderr)
	if !ok {
		return nil, false
	}

	walker := scanner.NewWalker(matcher, scanner.DefaultExtensions)
	eng := analyzer.NewEngine(nil)
	pool := scanner.NewPool(0)
	ctx := context.Background()

	_, occurrences, err := pool.RunWithOccurrences(ctx, walker, target, eng)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: drift scan failed: %v\n", err)
		return nil, false
	}

	return occurrences, true
}

func resolveDriftOptions(cfg *config.Config, dominance, outlier float64, minOccurrences int) drift.Options {
	opts := drift.DefaultOptions()
	if cfg != nil {
		opts = cfg.Drift.ToDriftOptions()
	}
	if dominance > 0 {
		opts.DominanceThreshold = dominance
	}
	if outlier > 0 {
		opts.OutlierThreshold = outlier
	}
	if minOccurrences > 0 {
		opts.MinClusterOccurrences = minOccurrences
	}
	return opts
}

func writeDriftReport(report drift.Report, outputFile, format string, noColor, reportFlag bool, stdout, stderr io.Writer) int {
	if outputFile != "" {
		cleanOutput := filepath.Clean(outputFile)
		outDir := filepath.Dir(cleanOutput)
		if mErr := os.MkdirAll(outDir, 0o750); mErr != nil {
			_, _ = fmt.Fprintf(stderr, "charites: error: failed to create output directory: %v\n", mErr)
			return ExitOperational
		}
		f, cErr := os.Create(cleanOutput)
		if cErr != nil {
			_, _ = fmt.Fprintf(stderr, "charites: error: failed to create output file: %v\n", cErr)
			return ExitOperational
		}
		_ = drift.RenderReport(f, &report, format, noColor)
		_ = f.Close()
		_, _ = fmt.Fprintf(stdout, "Charites style drift report saved to %s\n", cleanOutput)
	} else if reportFlag {
		_ = drift.RenderReport(stdout, &report, format, noColor)
	}

	return ExitClean
}
