package fuzz

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/analyzer"
	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/parser/astro"
	"github.com/will2469/charites/internal/parser/tsx"
	"github.com/will2469/charites/internal/reporter"
	"github.com/will2469/charites/internal/rules/theme"
)

func FuzzAstroPipeline(f *testing.F) {
	// 1. Seed Corpus
	seedFiles := []string{
		"../fixtures/sample.astro",
		"../fixtures/astro/clean.astro",
		"../fixtures/astro/opacity_violations.astro",
		"../fixtures/astro/complex_frontmatter.astro",
		"../fixtures/astro/inline_ignore.astro",
	}

	for _, file := range seedFiles {
		if content, err := os.ReadFile(filepath.Clean(file)); err == nil {
			f.Add(content)
		}
	}

	seeds := [][]byte{
		[]byte(""),
		[]byte("---\nconst a = 1;\n---\n<div class=\"bg-primary/10\"></div>"),
		[]byte("<broken <button class='btn'>Click</button>"),
		[]byte("<!-- unclosed comment"),
		[]byte("<div class={`p-4 ${dynamic} text-sm`}></div>"),
		[]byte("<slot name='header' /><input type='text' disabled>"),
		[]byte("<!-- charites:ignore theme.hardcode-opacity-color -->\n<div class=\"bg-primary/20\"></div>"),
		[]byte("<<<<<<>>>>>>"),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	// 2. Fuzz Target
	rule := theme.NewHardcodeOpacityColorRule()
	activeRules := []config.ActiveRule{
		{Rule: rule, EffectiveSeverity: rule.DefaultSeverity()},
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unhandled panic in FuzzAstroPipeline with input %q: %v", string(data), r)
			}
		}()

		root, err := astro.Parse(data)
		if err != nil || root == nil {
			return
		}

		inlineIgnores := analyzer.ParseDirectives(data)
		engine := analyzer.NewEngine(activeRules)
		diags := engine.AnalyzeTree("virtual.astro", root, inlineIgnores)

		// Evaluasi format keluaran reporter tanpa panic
		repJSON := reporter.NewJSONReporter()
		var bufJSON bytes.Buffer
		_ = repJSON.Render(&bufJSON, &reporter.ScanResult{
			Version:     "1.0.0",
			Diagnostics: diags,
			Summary: reporter.ScanSummary{
				ScannedFiles: 1,
				Passed:       len(diags) == 0,
			},
		})

		repInline := reporter.NewInlineReporter(reporter.ColorNever)
		var bufInline bytes.Buffer
		_ = repInline.Render(&bufInline, &reporter.ScanResult{
			Version:     "1.0.0",
			Diagnostics: diags,
			Summary: reporter.ScanSummary{
				ScannedFiles: 1,
				Passed:       len(diags) == 0,
			},
		})
	})
}

func FuzzTSXPipeline(f *testing.F) {
	// 1. Seed Corpus
	seedFiles := []string{
		"../fixtures/sample.tsx",
		"../fixtures/tsx/clean.tsx",
		"../fixtures/tsx/opacity_violations.tsx",
		"../fixtures/tsx/template_literals.tsx",
		"../fixtures/tsx/inline_ignore.tsx",
	}

	for _, file := range seedFiles {
		if content, err := os.ReadFile(filepath.Clean(file)); err == nil {
			f.Add(content)
		}
	}

	seeds := [][]byte{
		[]byte(""),
		[]byte("export const Box = () => <div className='bg-primary/10'>Hello</div>;"),
		[]byte("<button {...props} className={`p-4 ${dyn} text-sm`}>Click</button>"),
		[]byte("{/* charites:ignore theme.hardcode-opacity-color */}\n<div className=\"bg-primary/20\" />"),
		[]byte("<broken <tag attr='unclosed>"),
		[]byte("<><span>Fragment</span></>"),
		[]byte("<<<<<///// >>>>>"),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	// 2. Fuzz Target
	rule := theme.NewHardcodeOpacityColorRule()
	activeRules := []config.ActiveRule{
		{Rule: rule, EffectiveSeverity: rule.DefaultSeverity()},
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("unhandled panic in FuzzTSXPipeline with input %q: %v", string(data), r)
			}
		}()

		root, err := tsx.Extract(data)
		if err != nil || root == nil {
			return
		}

		inlineIgnores := analyzer.ParseDirectives(data)
		engine := analyzer.NewEngine(activeRules)
		diags := engine.AnalyzeTree("virtual.tsx", root, inlineIgnores)

		// Evaluasi format keluaran reporter tanpa panic
		repJSON := reporter.NewJSONReporter()
		var bufJSON bytes.Buffer
		_ = repJSON.Render(&bufJSON, &reporter.ScanResult{
			Version:     "1.0.0",
			Diagnostics: diags,
			Summary: reporter.ScanSummary{
				ScannedFiles: 1,
				Passed:       len(diags) == 0,
			},
		})

		repInline := reporter.NewInlineReporter(reporter.ColorNever)
		var bufInline bytes.Buffer
		_ = repInline.Render(&bufInline, &reporter.ScanResult{
			Version:     "1.0.0",
			Diagnostics: diags,
			Summary: reporter.ScanSummary{
				ScannedFiles: 1,
				Passed:       len(diags) == 0,
			},
		})
	})
}
