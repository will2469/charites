package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/config"
)

func TestIgnoreMatcher_BuiltinHardExclusionImmunity(t *testing.T) {
	// Buat matcher dengan aturan yang sengaja mencoba me-negasikan builtin exclusions
	lines := []string{
		"!node_modules/**",
		"!.git/**",
		"!dist/**",
		"!coverage/**",
	}
	matcher := config.NewIgnoreMatcher(lines)

	// Verifikasi direktori builtin tetap dipangkas 100%
	for _, b := range config.BuiltinExclusions() {
		if !matcher.ShouldIgnoreDir(b, b) {
			t.Errorf("builtin directory %q must be ignored despite negation rule", b)
		}
		if !matcher.ShouldIgnoreDir("sub", b+"/sub") {
			t.Errorf("sub-directory of builtin %q must be ignored", b)
		}
		if !matcher.ShouldIgnoreFile(b + "/index.tsx") {
			t.Errorf("file inside builtin %q must be ignored", b)
		}
	}
}

func TestIgnoreMatcher_HasBuiltinAncestor(t *testing.T) {
	matcher := config.NewIgnoreMatcher(nil)

	testCases := []struct {
		path     string
		expected bool
	}{
		{"node_modules/react/index.d.ts", true},
		{"./node_modules/pkg/Button.tsx", true},
		{".git/hooks/pre-commit", true},
		{"src/nested/dist/bundle.js", true},
		{"tests/coverage/lcov.info", true},
		{".turbo/cache/123", true},
		{".next/static/chunk.js", true},
		{"build/output/main.js", true},
		{"src/components/Button.tsx", false},
		{"apps/web/pages/index.astro", false},
		{"tests/fixtures/regression/test.tsx", false},
	}

	for _, tc := range testCases {
		got := matcher.HasBuiltinAncestor(tc.path)
		if got != tc.expected {
			t.Errorf("HasBuiltinAncestor(%q) = %v; expected %v", tc.path, got, tc.expected)
		}
	}
}

func TestIgnoreMatcher_SequentialEvaluation(t *testing.T) {
	lines := []string{
		"# Abaikan seluruh fixtures",
		"tests/fixtures/**",
		"# Kecualikan fixture khusus (negation)",
		"!tests/fixtures/keep.tsx",
		"# Abaikan direktori logs/",
		"logs/",
		"# Abaikan pola ekstensi",
		"*.tmp",
	}
	matcher := config.NewIgnoreMatcher(lines)

	// File yang cocok dengan ignore awal
	if !matcher.ShouldIgnoreFile("tests/fixtures/bad.tsx") {
		t.Error("expected tests/fixtures/bad.tsx to be ignored")
	}

	// File yang di-unignore oleh negasi
	if matcher.ShouldIgnoreFile("tests/fixtures/keep.tsx") {
		t.Error("expected tests/fixtures/keep.tsx to NOT be ignored due to negation")
	}

	// Dir-only rule
	if !matcher.ShouldIgnoreDir("logs", "logs") {
		t.Error("expected logs/ dir to be ignored")
	}
	if matcher.ShouldIgnoreFile("logs") {
		t.Error("expected logs file to NOT match dir-only rule")
	}

	// Basename pattern
	if !matcher.ShouldIgnoreFile("src/temp/cache.tmp") {
		t.Error("expected *.tmp file to be ignored")
	}
	if matcher.ShouldIgnoreFile("src/temp/cache.tsx") {
		t.Error("expected cache.tsx to NOT be ignored")
	}
}

func TestIgnoreMatcher_LoadAndAddPatterns(t *testing.T) {
	tmpDir := t.TempDir()
	ignoreFile := filepath.Join(tmpDir, ".charitesignore")

	content := `
# Custom ignores
custom_vendor/**
`
	if err := os.WriteFile(ignoreFile, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write ignore file: %v", err)
	}

	matcher, err := config.LoadIgnore(ignoreFile)
	if err != nil {
		t.Fatalf("LoadIgnore failed: %v", err)
	}

	if !matcher.ShouldIgnoreFile("custom_vendor/lib.tsx") {
		t.Error("expected custom_vendor/lib.tsx to be ignored")
	}

	// Add dynamic patterns from config
	matcher.AddPatterns([]string{"legacy/**"})
	if !matcher.ShouldIgnoreFile("legacy/old.tsx") {
		t.Error("expected dynamically added pattern legacy/** to be ignored")
	}

	// Missing file returns empty matcher with builtins intact
	missingMatcher, err := config.LoadIgnore(filepath.Join(tmpDir, "missing_ignore"))
	if err != nil {
		t.Fatalf("expected nil error for missing ignore file, got %v", err)
	}
	if !missingMatcher.HasBuiltinAncestor("node_modules/foo.tsx") {
		t.Error("builtin exclusions must be preserved in empty matcher")
	}

	// Read error when path is a directory
	_, err = config.LoadIgnore(tmpDir)
	if err == nil {
		t.Error("expected error when reading directory as ignore file")
	}
}

func TestIgnoreMatcher_GlobEdgeCases(t *testing.T) {
	lines := []string{
		"/anchored/file.tsx",
		"nested/**/target.tsx",
		"exact-match",
	}
	matcher := config.NewIgnoreMatcher(lines)

	if !matcher.ShouldIgnoreFile("anchored/file.tsx") {
		t.Error("expected anchored/file.tsx to be ignored")
	}
	if matcher.ShouldIgnoreFile("sub/anchored/file.tsx") {
		t.Error("expected sub/anchored/file.tsx NOT to match /anchored/file.tsx")
	}

	if !matcher.ShouldIgnoreFile("nested/a/b/c/target.tsx") {
		t.Error("expected nested/a/b/c/target.tsx to be ignored by **")
	}
	if !matcher.ShouldIgnoreFile("nested/target.tsx") {
		t.Error("expected nested/target.tsx to be ignored by ** (0 intermediate segments)")
	}
	if matcher.ShouldIgnoreFile("other/target.tsx") {
		t.Error("expected other/target.tsx NOT to match nested/**/target.tsx")
	}

	if !matcher.ShouldIgnoreFile("exact-match") {
		t.Error("expected exact-match to be ignored")
	}
}
