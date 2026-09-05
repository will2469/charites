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
	var category string
	var rule string
	var configPath string
	var noColor bool
	var failOnWarn bool
	var extFlags stringSliceFlag
	var ignoreFlags stringSliceFlag

	fs.StringVar(&format, "format", "inline", "Output format: inline or json")
	fs.StringVar(&format, "f", "inline", "Output format (shorthand)")

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

	// 1. Reorder argumen untuk fleksibilitas POSIX/GNU (flag dapat diletakkan di sebelum atau sesudah path)
	reorderedArgs, positionalArgs := partitionArgs(args)

	if err := fs.Parse(reorderedArgs); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, _ = fmt.Fprint(stdout, UsageString())
			return ExitClean
		}
		_, _ = fmt.Fprintf(stderr, "charites: error: %v. Run 'charites --help' for usage.\n", err)
		return ExitOperational
	}

	// 2. Validasi Batasan Target Path
	if len(positionalArgs) > 1 {
		_, _ = fmt.Fprintf(stderr, "charites: error: multiple scan targets not supported. Specify a single path.\n")
		return ExitOperational
	}

	target := "."
	if len(positionalArgs) == 1 {
		target = positionalArgs[0]
	}

	// 3. Validasi Keberadaan Target Path
	if _, err := os.Stat(target); err != nil {
		_, _ = fmt.Fprintf(stderr, "charites: error: scan target %q does not exist.\n", target)
		return ExitOperational
	}

	// 4. Validasi Format Output
	formatLower := strings.ToLower(strings.TrimSpace(format))
	if formatLower != "inline" && formatLower != "json" {
		_, _ = fmt.Fprintf(stderr, "charites: error: unsupported format %q. Supported formats: inline, json.\n", format)
		return ExitOperational
	}

	// 5. Validasi & Normalisasi Flag --ext
	normalizedExts, extErr := normalizeExtensions(extFlags)
	if extErr != nil {
		_, _ = fmt.Fprintln(stderr, extErr.Error())
		return ExitOperational
	}

	// 6. Validasi Keberadaan & Konflik --category dan --rule
	reg := rules.DefaultRegistry()

	if category != "" {
		if len(reg.ByCategory(category)) == 0 {
			_, _ = fmt.Fprintf(stderr, "charites: error: unknown category %q.\n", category)
			return ExitOperational
		}
	}

	if rule != "" {
		r, exists := reg.Get(rule)
		if !exists {
			_, _ = fmt.Fprintf(stderr, "charites: error: unknown rule %q.\n", rule)
			return ExitOperational
		}
		if category != "" && r.Category() != category {
			_, _ = fmt.Fprintf(stderr, "charites: error: rule %q does not belong to category %q.\n", rule, category)
			return ExitOperational
		}
	}

	// 7. Resolusi Konfigurasi charites.yaml
	var cfg *config.Config
	isExplicitConfig := isFlagPassed(args, "-config", "--config")

	if isExplicitConfig {
		var err error
		cfg, err = config.Load(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				_, _ = fmt.Fprintf(stderr, "charites: error: config file not found: %q.\n", configPath)
				return ExitOperational
			}
			_, _ = fmt.Fprintf(stderr, "charites: error: failed to parse config %q: %v\n", configPath, err)
			return ExitOperational
		}
	} else {
		// Zero-config default: periksa target terlebih dahulu jika berupa direktori
		targetConfig := filepath.Join(target, "charites.yaml")
		if _, err := os.Stat(targetConfig); err == nil {
			cfg, _ = config.Load(targetConfig)
		} else {
			targetConfigYml := filepath.Join(target, "charites.yml")
			if _, err := os.Stat(targetConfigYml); err == nil {
				cfg, _ = config.Load(targetConfigYml)
			} else {
				cfg, _ = config.Load("")
			}
		}
	}

	// 8. Resolusi Active Rules (3-Tier Precedence: Registry -> CLI Scope -> Config Policy)
	activeRules := cfg.ResolveActiveRules(reg, category, rule)

	// 9. Penyiapan Ignore Matcher & Direct-Target Safety
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

	// Direct-Target Safety Check (kekebalan builtin hard exclusions)
	if matcher.HasBuiltinAncestor(target) {
		_, _ = fmt.Fprintf(stderr, "charites: error: scan target %q is within excluded directory (builtin hard exclusion).\n", target)
		return ExitOperational
	}

	// 10. Inisialisasi Scanner Walker & Analyzer Engine
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

	// 11. Agregasi Metrik Temuan
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

	result := &reporter.ScanResult{
		Version: Version,
		Summary: reporter.ScanSummary{
			ScannedFiles: int(ca.count.Load()),
			DurationMS:   durationMS,
			ErrorCount:   errCount,
			WarningCount: warnCount,
			InfoCount:    infoCount,
			Passed:       passed,
		},
		Diagnostics: diags,
	}

	// 12. Presentasi Dokumen Laporan
	if formatLower == "json" {
		rep := reporter.NewJSONReporter()
		_ = rep.Render(stdout, result)
	} else {
		colorMode := reporter.ResolveColorMode(noColor, stdout)
		rep := reporter.NewInlineReporter(colorMode)
		_ = rep.Render(stdout, result)
	}

	// 13. Resolusi Kode Keluar
	return ResolveExitCode(&result.Summary, failOnWarn)
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
