package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/will2469/charites/internal/config"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// mockRule mengimplementasikan rules.Rule untuk pengujian konfigurasi.
type mockRule struct {
	id          string
	category    string
	defaultSev  ir.Severity
	description string
}

func (m *mockRule) ID() string                        { return m.id }
func (m *mockRule) Category() string                  { return m.category }
func (m *mockRule) DefaultSeverity() ir.Severity      { return m.defaultSev }
func (m *mockRule) Description() string               { return m.description }
func (m *mockRule) Evaluate(*ir.Node) []ir.Diagnostic { return nil }

func setupRegistry() *rules.Registry {
	reg := rules.NewRegistry()
	_ = reg.Register(&mockRule{id: "theme.opacity", category: "theme", defaultSev: ir.SeverityError, description: "Opacity rule"})
	_ = reg.Register(&mockRule{id: "theme.color", category: "theme", defaultSev: ir.SeverityWarn, description: "Color rule"})
	_ = reg.Register(&mockRule{id: "a11y.alt", category: "a11y", defaultSev: ir.SeverityError, description: "Alt text rule"})
	_ = reg.Register(&mockRule{id: "perf.bundle", category: "perf", defaultSev: ir.SeverityInfo, description: "Bundle size rule"})
	return reg
}

func TestResolveActiveRules_DefaultYes(t *testing.T) {
	reg := setupRegistry()

	// 1. Config nil (Zero-Config Default: YES)
	var nilConfig *config.Config
	active := nilConfig.ResolveActiveRules(reg, "", "")

	if len(active) != 4 {
		t.Fatalf("expected all 4 rules active under nil config, got %d", len(active))
	}

	for _, a := range active {
		if a.EffectiveSeverity != a.Rule.DefaultSeverity() {
			t.Errorf("rule %s expected severity %v, got %v", a.Rule.ID(), a.Rule.DefaultSeverity(), a.EffectiveSeverity)
		}
	}

	// 2. Empty Config (Rules map empty)
	emptyConfig := &config.Config{}
	activeEmpty := emptyConfig.ResolveActiveRules(reg, "", "")
	if len(activeEmpty) != 4 {
		t.Fatalf("expected all 4 rules active under empty config, got %d", len(activeEmpty))
	}
}

func parseTestConfig(t *testing.T) *config.Config {
	t.Helper()
	yamlContent := `
format: inline
scan_path: ./src
theme: ./custom/theme.css
rules:
  theme.opacity: warn
  theme.color: "off"
  a11y.alt: disabled
  perf.bundle: error
ignore:
  - "dist/**"
  - 'build/**'
`
	cfg, err := config.Parse([]byte(yamlContent))
	if err != nil {
		t.Fatalf("config.Parse failed: %v", err)
	}
	return cfg
}

func TestResolveActiveRules_OverridesAndPrecedence(t *testing.T) {
	reg := setupRegistry()
	cfg := parseTestConfig(t)

	if cfg.Format != "inline" || cfg.ScanPath != "./src" || cfg.Theme != "./custom/theme.css" {
		t.Errorf("unexpected top-level values: format=%q, scan_path=%q, theme=%q", cfg.Format, cfg.ScanPath, cfg.Theme)
	}
	if len(cfg.Ignore) != 2 || cfg.Ignore[0] != "dist/**" || cfg.Ignore[1] != "build/**" {
		t.Errorf("unexpected ignore list: %+v", cfg.Ignore)
	}

	active := cfg.ResolveActiveRules(reg, "", "")
	if len(active) != 2 {
		t.Fatalf("expected 2 active rules, got %d", len(active))
	}

	ruleMap := make(map[string]ir.Severity)
	for _, a := range active {
		ruleMap[a.Rule.ID()] = a.EffectiveSeverity
	}

	if sev, ok := ruleMap["theme.opacity"]; !ok || sev != ir.SeverityWarn {
		t.Errorf("theme.opacity expected warn, got %v (ok=%v)", sev, ok)
	}
	if sev, ok := ruleMap["perf.bundle"]; !ok || sev != ir.SeverityError {
		t.Errorf("perf.bundle expected error, got %v (ok=%v)", sev, ok)
	}
}

func TestResolveActiveRules_PolicyOverridesCLI(t *testing.T) {
	reg := setupRegistry()
	cfg := parseTestConfig(t)

	activeForced := cfg.ResolveActiveRules(reg, "", "theme.color")
	if len(activeForced) != 0 {
		t.Errorf("expected theme.color to be disabled by policy despite CLI flag, got %d", len(activeForced))
	}
}

func TestResolveActiveRules_CategoryFilter(t *testing.T) {
	reg := setupRegistry()
	cfg := parseTestConfig(t)

	activeCategory := cfg.ResolveActiveRules(reg, "theme", "")
	if len(activeCategory) != 1 || activeCategory[0].Rule.ID() != "theme.opacity" {
		t.Fatalf("expected only theme.opacity under category=theme, got %d", len(activeCategory))
	}
}

func TestResolveActiveRules_NilRegistry(t *testing.T) {
	cfg := &config.Config{}
	active := cfg.ResolveActiveRules(nil, "", "")
	if len(active) != 0 {
		t.Errorf("expected 0 rules for nil registry, got %d", len(active))
	}
}

func TestConfig_Load(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "charites.yaml")

	content := `
rules:
  theme.opacity: info
`
	if err := os.WriteFile(cfgPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}

	// Load file yang ada
	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatalf("config.Load failed: %v", err)
	}
	if cfg == nil || cfg.Rules["theme.opacity"] != "info" {
		t.Errorf("unexpected loaded config: %+v", cfg)
	}

	// Load file eksplisit yang tidak ada (harus error)
	_, err = config.Load(filepath.Join(tmpDir, "non_existent.yaml"))
	if err == nil {
		t.Error("expected error for non-existent explicit config path, got nil")
	}

	// Load default ketika file tidak ada di current dir
	origWd, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origWd) }()

	_ = os.Remove(cfgPath)
	defCfg, err := config.Load("")
	if err != nil {
		t.Fatalf("expected nil error for missing default config, got %v", err)
	}
	if defCfg != nil {
		t.Errorf("expected nil config for missing default, got %+v", defCfg)
	}

	// Test charites.yml fallback
	ymlPath := filepath.Join(tmpDir, "charites.yml")
	if writeErr := os.WriteFile(ymlPath, []byte("scan_path: ./app\n"), 0o600); writeErr != nil {
		t.Fatalf("failed to write yml: %v", writeErr)
	}
	loadedYml, loadErr := config.Load("")
	if loadErr != nil || loadedYml == nil || loadedYml.ScanPath != "./app" {
		t.Errorf("expected fallback to charites.yml, got cfg=%+v, err=%v", loadedYml, loadErr)
	}
}

func TestResolveActiveRules_RuleFilterNotFound(t *testing.T) {
	reg := setupRegistry()
	cfg := &config.Config{}
	active := cfg.ResolveActiveRules(reg, "non-existent-category", "non-existent-rule")
	if len(active) != 0 {
		t.Errorf("expected 0 rules, got %d", len(active))
	}
}

func TestConfig_ParseConvention(t *testing.T) {
	// 1. Inline format
	inlineYAML := `
convention:
  prefixes: ["--custom-", "--"]
  opacity_mappings:
    "10": ["-light", "-subtle"]
    "20": ["-light"]
  fallbacks:
    secondary: ["muted"]
`
	cfgInline, err := config.Parse([]byte(inlineYAML))
	if err != nil {
		t.Fatalf("failed parsing inline convention: %v", err)
	}
	if len(cfgInline.Convention.Prefixes) != 2 || cfgInline.Convention.Prefixes[0] != "--custom-" {
		t.Errorf("unexpected prefixes: %+v", cfgInline.Convention.Prefixes)
	}
	if len(cfgInline.Convention.OpacityMappings["10"]) != 2 {
		t.Errorf("unexpected opacity mappings for 10: %+v", cfgInline.Convention.OpacityMappings["10"])
	}
	if len(cfgInline.Convention.Fallbacks["secondary"]) != 1 || cfgInline.Convention.Fallbacks["secondary"][0] != "muted" {
		t.Errorf("unexpected fallbacks for secondary: %+v", cfgInline.Convention.Fallbacks["secondary"])
	}

	// 2. Multi-line bullet format
	bulletYAML := `
convention:
  prefixes:
    - "--my-theme-"
    - "--brand-"
  opacity_mappings:
    "15":
      - "-tint"
      - "-soft"
  fallbacks:
    accent:
      - "highlight"
`
	cfgBullet, err := config.Parse([]byte(bulletYAML))
	if err != nil {
		t.Fatalf("failed parsing bullet convention: %v", err)
	}
	if len(cfgBullet.Convention.Prefixes) != 2 || cfgBullet.Convention.Prefixes[0] != "--my-theme-" {
		t.Errorf("unexpected prefixes: %+v", cfgBullet.Convention.Prefixes)
	}
	if len(cfgBullet.Convention.OpacityMappings["15"]) != 2 || cfgBullet.Convention.OpacityMappings["15"][0] != "-tint" {
		t.Errorf("unexpected opacity mappings for 15: %+v", cfgBullet.Convention.OpacityMappings["15"])
	}
	if len(cfgBullet.Convention.Fallbacks["accent"]) != 1 || cfgBullet.Convention.Fallbacks["accent"][0] != "highlight" {
		t.Errorf("unexpected fallbacks for accent: %+v", cfgBullet.Convention.Fallbacks["accent"])
	}
}
