package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"time"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/reporter"
	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/scanner"
)

type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(val string) error {
	*s = append(*s, val)
	return nil
}

type countingAnalyzer struct {
	inner scanner.FileAnalyzer
	count atomic.Int64
}

func (c *countingAnalyzer) AnalyzeFile(path string) ([]ir.Diagnostic, error) {
	c.count.Add(1)
	return c.inner.AnalyzeFile(path)
}

// RunScan mengorkestrasi pipeline pemindaian kode frontend sesuai kontrak SPEC-05-CLI.
func RunScan(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("scan", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.Usage = func() {}

	var format string
	var outputFile string
	var category string
	var rule string
	var configPath string
	var noColor bool
	var failOnWarn bool
	var extFlags stringSliceFlag
	var ignoreFlags stringSliceFlag

	fs.StringVar(&format, "format", "inline", "Output format: inline, json, or markdown (alias: md)")
	fs.StringVar(&format, "f", "inline", "Output format (shorthand)")

	fs.StringVar(&outputFile, "output", "", "Path to output report file (e.g. report.md)")
	fs.StringVar(&outputFile, "o", "", "Path to output report file (shorthand)")

	fs.Var(&extFlags, "ext", "Filter extensions: astro, tsx, jsx")
	fs.Var(&extFlags, "e", "Filter extensions (shorthand)")

	fs.StringVar(&category, "category", "", "Filter by rule category")
	fs.StringVar(&category, "c", "", "Filter by rule category (shorthand)")

	fs.StringVar(&rule, "rule", "", "Filter by canonical rule ID")
	fs.StringVar(&rule, "r", "", "Filter by canonical rule ID (shorthand)")

	fs.StringVar(&configPath, "config", "charites.yaml", "Path to config file")

	fs.Var(&ignoreFlags, "ignore", "Additional custom ignore patterns")

	fs.BoolVar(&noColor, "no-color", false, "Disable ANSI color formatting")
	fs.BoolVar(&failOnWarn, "fail-on-warn", false, "Exit with code 1 on warnings")

	// 1. Reorder argumen untuk fleksibilitas POSIX/GNU
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

	normalizedExts, extErr := normalizeExtensions(extFlags)
	if extErr != nil {
		_, _ = fmt.Fprintln(stderr, extErr.Error())
		return ExitOperational
	}

	reg := rules.DefaultRegistry()
	if !validateCategoryAndRule(reg, category, rule, stderr) {
		return ExitOperational
	}

	activeRules := cfg.ResolveActiveRules(reg, category, rule)
	matcher, ok := buildScanMatcher(target, cfg, ignoreFlags, stderr)
	if !ok {
		return ExitOperational
	}

	walker := scanner.NewWalker(matcher, normalizedExts)
	eng := analyzer.NewEngine(activeRules)
	ca := &countingAnalyzer{inner: eng}

	pool := scanner.NewPool(0)
	ctx := context.Background()

	startTime := time.Now()
	diags, err := pool.Run(ctx, walker, target, ca)
	durationMS := time.Since(startTime).Milliseconds()

	if ctx.Err() != nil {
		return ExitInterrupted
	}
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: scan execution failed: %v\n", err)
		return ExitOperational
	}

	result := buildScanResult(ca, diags, durationMS, failOnWarn, activeRules, target, startTime)
	if outputFile != "" {
		if code := saveReportToFile(outputFile, result, format, noColor, stdout, stderr); code != -1 {
			return code
		}
	} else {
		renderScanResult(stdout, result, format, noColor)
	}

	return ResolveExitCode(&result.Summary, failOnWarn)
}

func resolveReportFormat(format, outputFile string) string {
	if outputFile != "" && format == "inline" {
		lowerOut := strings.ToLower(outputFile)
		if strings.HasSuffix(lowerOut, ".md") || strings.HasSuffix(lowerOut, ".markdown") {
			return "markdown"
		}
		if strings.HasSuffix(lowerOut, ".json") {
			return "json"
		}
	}
	return format
}

func saveReportToFile(outputFile string, result *reporter.ScanResult, format string, noColor bool, stdout, stderr io.Writer) int {
	cleanOutput := filepath.Clean(outputFile)
	outDir := filepath.Dir(cleanOutput)
	if err := os.MkdirAll(outDir, 0o750); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to create output directory: %v\n", err)
		return ExitOperational
	}
	f, err := os.Create(filepath.Clean(cleanOutput))
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: failed to create output report file: %v\n", err)
		return ExitOperational
	}
	renderScanResult(f, result, format, noColor)
	_ = f.Close()
	_, _ = fmt.Fprintf(stdout, "Charites audit report saved to %s\n", cleanOutput)
	return -1
}

func validateScanTargetAndFormat(target, format string, positionalArgs []string, stderr io.Writer) bool {
	if len(positionalArgs) > 1 {
		_, _ = fmt.Fprintf(stderr, "charites: error: multiple scan targets not supported. Specify a single path.\n")
		return false
	}
	if _, err := os.Stat(target); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: scan target %q does not exist.\n", target)
		return false
	}
	formatLower := strings.ToLower(strings.TrimSpace(format))
	if formatLower != "inline" && formatLower != "json" && formatLower != "markdown" && formatLower != "md" {
		_, _ = fmt.Fprintf(stderr, "charites: error: unsupported format %q. Supported formats: inline, json, markdown, md.\n", format)
		return false
	}
	return true
}

func validateCategoryAndRule(reg *rules.Registry, category, rule string, stderr io.Writer) bool {
	if category != "" && len(reg.ByCategory(category)) == 0 {
		_, _ = fmt.Fprintf(stderr, "charites: error: unknown category %q.\n", category)
		return false
	}
	if rule != "" {
		r, exists := reg.Get(rule)
		if !exists {
			_, _ = fmt.Fprintf(stderr, "charites: error: unknown rule %q.\n", rule)
			return false
		}
		if category != "" && r.Category() != category {
			_, _ = fmt.Fprintf(stderr, "charites: error: rule %q does not belong to category %q.\n", rule, category)
			return false
		}
	}
	return true
}

func resolveScanConfig(target, configPath string, isExplicit bool, stderr io.Writer) (*config.Config, bool) {
	if isExplicit {
		cfg, err := config.Load(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				_, _ = fmt.Fprintf(stderr, "charites: error: config file not found: %q.\n", configPath)
			} else {
				_, _ = fmt.Fprintf(stderr, "charites: error: failed to parse config %q: %v\n", configPath, err)
			}
			return nil, false
		}
		return cfg, true
	}

	for _, cand := range []string{"charites.yaml", "charites.yml"} {
		p := filepath.Join(target, cand)
		if _, err := os.Stat(p); err == nil {
			cfg, _ := config.Load(p)
			return cfg, true
		}
	}
	cfg, _ := config.Load("")
	return cfg, true
}

func buildScanMatcher(target string, cfg *config.Config, ignoreFlags []string, stderr io.Writer) (*config.IgnoreMatcher, bool) {
	var matcher *config.IgnoreMatcher
	targetIgnore := filepath.Join(target, ".charitesignore")
	if _, err := os.Stat(targetIgnore); err == nil {
		matcher, _ = config.LoadIgnore(targetIgnore)
	} else {
		matcher, _ = config.LoadIgnore(".charitesignore")
	}
	if matcher == nil {
		matcher = config.NewIgnoreMatcher(nil)
	}

	if cfg != nil && len(cfg.Ignore) > 0 {
		matcher.AddPatterns(cfg.Ignore)
	}
	if len(ignoreFlags) > 0 {
		matcher.AddPatterns(ignoreFlags)
	}

	if matcher.HasBuiltinAncestor(target) {
		_, _ = fmt.Fprintf(stderr, "charites: error: scan target %q is within excluded directory (builtin hard exclusion).\n", target)
		return nil, false
	}
	return matcher, true
}

func buildScanResult(ca *countingAnalyzer, diags []ir.Diagnostic, durationMS int64, failOnWarn bool, activeRules []config.ActiveRule, target string, startTime time.Time) *reporter.ScanResult {
	var errCount, warnCount, infoCount int
	for _, d := range diags {
		switch d.Severity {
		case ir.SeverityError:
			errCount++
		case ir.SeverityWarn:
			warnCount++
		case ir.SeverityInfo:
			infoCount++
		}
	}

	passed := errCount == 0
	if failOnWarn && warnCount > 0 {
		passed = false
	}

	counts := make(map[string]int)
	for _, d := range diags {
		counts[d.Rule]++
	}

	attached := make([]reporter.RuleAuditInfo, 0, len(activeRules))
	for _, ar := range activeRules {
		r := ar.Rule
		c := counts[r.ID()]
		st := "PASS"
		if c > 0 {
			st = "FAILED"
		}
		sev := string(ar.EffectiveSeverity)
		if sev == "" {
			sev = string(r.DefaultSeverity())
		}
		attached = append(attached, reporter.RuleAuditInfo{
			ID:          r.ID(),
			Category:    r.Category(),
			Description: r.Description(),
			Severity:    sev,
			IssuesFound: c,
			Status:      st,
		})
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = target
	}
	absRoot, err := filepath.Abs(cwd)
	if err != nil {
		absRoot = cwd
	}

	return &reporter.ScanResult{
		Version:   Version,
		Timestamp: startTime,
		RootDir:   absRoot,
		Summary: reporter.ScanSummary{
			ScannedFiles: int(ca.count.Load()),
			DurationMS:   durationMS,
			ErrorCount:   errCount,
			WarningCount: warnCount,
			InfoCount:    infoCount,
			Passed:       passed,
		},
		Diagnostics:   diags,
		AttachedRules: attached,
	}
}

func renderScanResult(w io.Writer, result *reporter.ScanResult, format string, noColor bool) {
	fmtLower := strings.ToLower(strings.TrimSpace(format))
	switch fmtLower {
	case "json":
		rep := reporter.NewJSONReporter()
		_ = rep.Render(w, result)
	case "markdown", "md":
		rep := reporter.NewMarkdownReporter(
			reporter.WithRootDir(result.RootDir),
			reporter.WithTimestamp(result.Timestamp),
		)
		_ = rep.Render(w, result)
	default:
		colorMode := reporter.ResolveColorMode(noColor, w)
		rep := reporter.NewInlineReporter(colorMode)
		_ = rep.Render(w, result)
	}
}

func partitionArgs(args []string) ([]string, []string) {
	var flagArgs []string
	var posArgs []string

	takesValue := map[string]bool{
		"-f": true, "--format": true,
		"-e": true, "--ext": true,
		"-c": true, "--category": true,
		"-r": true, "--rule": true,
		"-config": true, "--config": true,
		"-ignore": true, "--ignore": true,
		"-o": true, "--output": true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flagArgs = append(flagArgs, arg)
			// Jika flag membutuhkan nilai dan tidak menggunakan '=', ambil argumen berikutnya
			if !strings.Contains(arg, "=") {
				k := arg
				if takesValue[k] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					i++
					flagArgs = append(flagArgs, args[i])
				}
			}
		} else {
			posArgs = append(posArgs, arg)
		}
	}

	return flagArgs, posArgs
}

func isFlagPassed(args []string, names ...string) bool {
	for _, arg := range args {
		for _, name := range names {
			if arg == name || strings.HasPrefix(arg, name+"=") {
				return true
			}
		}
	}
	return false
}

func normalizeExtensions(extFlags []string) ([]string, error) {
	if len(extFlags) == 0 {
		return scanner.DefaultExtensions, nil
	}

	var rawTokens []string
	for _, val := range extFlags {
		trimmed := strings.TrimSpace(val)
		if trimmed == "" {
			return nil, fmt.Errorf("charites: error: empty extension flag.")
		}
		parts := strings.Split(trimmed, ",")
		for _, p := range parts {
			pt := strings.TrimSpace(p)
			if pt == "" {
				return nil, fmt.Errorf("charites: error: empty extension flag.")
			}
			rawTokens = append(rawTokens, pt)
		}
	}

	validSet := map[string]bool{
		".astro": true,
		".tsx":   true,
		".jsx":   true,
	}

	var result []string
	for _, tok := range rawTokens {
		low := strings.ToLower(tok)
		if !strings.HasPrefix(low, ".") {
			low = "." + low
		}
		if !validSet[low] {
			return nil, fmt.Errorf("charites: error: unsupported extension %q. Supported extensions: .astro, .tsx, .jsx.", tok)
		}
		if !slices.Contains(result, low) {
			result = append(result, low)
		}
	}

	return result, nil
}
