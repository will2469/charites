package config

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/will2469/charites/internal/drift"
	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
	"github.com/will2469/charites/internal/token"
)

// ActiveRule membungkus rule singleton dengan EffectiveSeverity hasil resolusi konfigurasi.
// Mempertahankan invariant immutability dari registry singleton.
type ActiveRule struct {
	Rule              rules.Rule
	EffectiveSeverity ir.Severity
}

// ConventionConfig merepresentasikan konfigurasi inferensi semantik token di charites.yaml.
type ConventionConfig = token.ConventionConfig

// DriftConfig merepresentasikan konfigurasi kebijakan statistik analisis style drift.
type DriftConfig struct {
	DominanceThreshold    float64                                        `json:"dominance_threshold" yaml:"dominance_threshold"`
	OutlierThreshold      float64                                        `json:"outlier_threshold" yaml:"outlier_threshold"`
	MinClusterOccurrences int                                            `json:"min_cluster_occurrences" yaml:"min_cluster_occurrences"`
	Exceptions            map[string]map[string][]drift.ContextException `json:"exceptions" yaml:"exceptions"`
}

// ToDriftOptions mengonversi DriftConfig menjadi drift.Options untuk engine clusterer.
func (dc DriftConfig) ToDriftOptions() drift.Options {
	opts := drift.DefaultOptions()
	if dc.DominanceThreshold > 0 {
		opts.DominanceThreshold = dc.DominanceThreshold
	}
	if dc.OutlierThreshold > 0 {
		opts.OutlierThreshold = dc.OutlierThreshold
	}
	if dc.MinClusterOccurrences > 0 {
		opts.MinClusterOccurrences = dc.MinClusterOccurrences
	}
	if len(dc.Exceptions) > 0 {
		opts.Exceptions = dc.Exceptions
	}
	return opts
}

// Config merepresentasikan konfigurasi proyek dari charites.yaml.
type Config struct {
	Format     string            `json:"format" yaml:"format"`
	Output     string            `json:"output" yaml:"output"`                           // Path berkas output laporan (misal: report.md)
	Telemetry  *bool             `json:"telemetry,omitempty" yaml:"telemetry,omitempty"` // Status izin pelaporan telemetri / issue (default: true)
	ScanPath   string            `json:"scan_path" yaml:"scan_path"`
	Theme      string            `json:"theme" yaml:"theme"`           // Custom path ke SSOT tema (CSS/JSON) jika di luar path standar
	Convention ConventionConfig  `json:"convention" yaml:"convention"` // Konfigurasi konvensi semantik token
	Drift      DriftConfig       `json:"drift" yaml:"drift"`           // Konfigurasi ambang batas dan exception style drift
	Rules      map[string]string `json:"rules" yaml:"rules"`           // "rule-id": "off" | "warn" | "error" | "info"
	Ignore     []string          `json:"ignore" yaml:"ignore"`         // Pola path tambahan
}

// Load membaca dan mem-parse berkas konfigurasi dari path yang diberikan.
// Jika path kosong, mencoba membaca "charites.yaml" atau "charites.yml" di direktori kerja.
// Jika berkas default tidak ditemukan, mengembalikan nil, nil (memenuhi invarian Zero-Config Default: YES).
// Jika path eksplisit diberikan namun tidak ditemukan, mengembalikan error.
func Load(path string) (*Config, error) {
	if path == "" {
		for _, defaultName := range []string{"charites.yaml", "charites.yml"} {
			if _, err := os.Stat(defaultName); err == nil {
				path = defaultName
				break
			}
		}
		if path == "" {
			// Zero-config operability: Default: YES
			return nil, nil
		}
	}

	data, err := os.ReadFile(filepath.Clean(path)) //nolint:gosec // controlled configuration path
	if err != nil {
		return nil, err
	}

	return Parse(data)
}

// Parse mem-parse data konfigurasi YAML secara deterministik tanpa external dependency.
func Parse(data []byte) (*Config, error) {
	cfg := &Config{
		Rules:  make(map[string]string),
		Ignore: make([]string, 0),
		Convention: ConventionConfig{
			OpacityMappings: make(map[string][]string),
			Fallbacks:       make(map[string][]string),
			Prefixes:        make([]string, 0),
		},
		Drift: DriftConfig{
			DominanceThreshold:    80.0,
			OutlierThreshold:      10.0,
			MinClusterOccurrences: 10,
			Exceptions:            make(map[string]map[string][]drift.ContextException),
		},
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.IndexByte(line, '#'); idx != -1 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		indent := len(line) - len(strings.TrimLeft(line, " "))
		if indent == 0 {
			parseTopLevel(trimmed, cfg, &currentSection)
			continue
		}

		if nextSec, ok := resolveSubSection(trimmed, currentSection, indent); ok {
			currentSection = nextSec
			continue
		}

		parseIndentedSection(trimmed, currentSection, cfg)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func resolveSubSection(trimmed, currentSection string, indent int) (string, bool) {
	if !strings.HasSuffix(trimmed, ":") {
		return currentSection, false
	}
	sub := cleanValue(strings.TrimSuffix(trimmed, ":"))
	if strings.HasPrefix(currentSection, "convention") {
		return resolveConventionSubSection(sub, currentSection, indent), true
	}
	if strings.HasPrefix(currentSection, "drift") {
		return resolveDriftSubSection(sub, currentSection, indent), true
	}
	return currentSection, false
}

func resolveConventionSubSection(sub, currentSection string, indent int) string {
	if indent <= 2 {
		return "convention." + sub
	}
	switch {
	case strings.HasPrefix(currentSection, "convention.opacity_mappings"):
		return "convention.opacity_mappings." + sub
	case strings.HasPrefix(currentSection, "convention.fallbacks"):
		return "convention.fallbacks." + sub
	default:
		return "convention." + sub
	}
}

func resolveDriftSubSection(sub, currentSection string, indent int) string {
	switch {
	case indent <= 2:
		return "drift." + sub
	case indent <= 4 && strings.HasPrefix(currentSection, "drift.exceptions"):
		return "drift.exceptions." + sub
	case indent <= 6 && strings.HasPrefix(currentSection, "drift.exceptions."):
		parts := strings.Split(currentSection, ".")
		if len(parts) >= 3 {
			return "drift.exceptions." + parts[2] + "." + sub
		}
	}
	return currentSection
}

func parseTopLevel(trimmed string, cfg *Config, currentSection *string) {
	if strings.HasSuffix(trimmed, ":") {
		*currentSection = strings.TrimSuffix(trimmed, ":")
		return
	}

	if k, v, found := strings.Cut(trimmed, ":"); found {
		*currentSection = ""
		key := strings.TrimSpace(k)
		val := cleanValue(v)
		switch key {
		case "format", "report_format":
			cfg.Format = val
		case "output", "output_file", "report_file":
			cfg.Output = val
		case "scan_path":
			cfg.ScanPath = val
		case "theme":
			cfg.Theme = val
		case "telemetry":
			b := parseBool(val)
			cfg.Telemetry = &b
		}
	}
}

func parseBool(val string) bool {
	norm := strings.ToLower(strings.TrimSpace(val))
	return norm == "true" || norm == "1" || norm == "on" || norm == "yes"
}

func parseIndentedSection(trimmed, currentSection string, cfg *Config) {
	switch {
	case currentSection == "rules":
		if k, v, found := strings.Cut(trimmed, ":"); found {
			ruleID := strings.TrimSpace(k)
			val := cleanValue(v)
			if ruleID != "" {
				cfg.Rules[ruleID] = val
			}
		}
	case currentSection == "ignore":
		if strings.HasPrefix(trimmed, "-") {
			val := cleanValue(strings.TrimPrefix(trimmed, "-"))
			if val != "" {
				cfg.Ignore = append(cfg.Ignore, val)
			}
		}
	case strings.HasPrefix(currentSection, "convention"):
		parseConventionSection(trimmed, currentSection, cfg)
	case strings.HasPrefix(currentSection, "drift"):
		parseDriftSection(trimmed, currentSection, cfg)
	}
}

func parseConventionSection(trimmed, currentSection string, cfg *Config) {
	if cfg.Convention.OpacityMappings == nil {
		cfg.Convention.OpacityMappings = make(map[string][]string)
	}
	if cfg.Convention.Fallbacks == nil {
		cfg.Convention.Fallbacks = make(map[string][]string)
	}

	switch {
	case currentSection == "convention":
		if k, v, found := strings.Cut(trimmed, ":"); found {
			if strings.TrimSpace(k) == "prefixes" {
				cfg.Convention.Prefixes = parseStringList(cleanValue(v))
			}
		}
	case strings.HasPrefix(currentSection, "convention.opacity_mappings"):
		parseConventionOpacity(trimmed, currentSection, cfg)
	case strings.HasPrefix(currentSection, "convention.fallbacks"):
		parseConventionFallbacks(trimmed, currentSection, cfg)
	case currentSection == "convention.prefixes":
		parseConventionPrefixes(trimmed, cfg)
	}
}

func parseConventionOpacity(trimmed, currentSection string, cfg *Config) {
	if currentSection == "convention.opacity_mappings" {
		if k, v, found := strings.Cut(trimmed, ":"); found {
			op := cleanValue(k)
			val := strings.TrimSpace(v)
			if val != "" {
				cfg.Convention.OpacityMappings[op] = parseStringList(val)
			}
		}
		return
	}
	op := strings.TrimPrefix(currentSection, "convention.opacity_mappings.")
	if strings.HasPrefix(trimmed, "-") {
		val := cleanValue(strings.TrimPrefix(trimmed, "-"))
		if val != "" {
			cfg.Convention.OpacityMappings[op] = append(cfg.Convention.OpacityMappings[op], val)
		}
	}
}

func parseConventionFallbacks(trimmed, currentSection string, cfg *Config) {
	if currentSection == "convention.fallbacks" {
		if k, v, found := strings.Cut(trimmed, ":"); found {
			base := cleanValue(k)
			val := strings.TrimSpace(v)
			if val != "" {
				cfg.Convention.Fallbacks[base] = parseStringList(val)
			}
		}
		return
	}
	base := strings.TrimPrefix(currentSection, "convention.fallbacks.")
	if strings.HasPrefix(trimmed, "-") {
		val := cleanValue(strings.TrimPrefix(trimmed, "-"))
		if val != "" {
			cfg.Convention.Fallbacks[base] = append(cfg.Convention.Fallbacks[base], val)
		}
	}
}

func parseConventionPrefixes(trimmed string, cfg *Config) {
	if strings.HasPrefix(trimmed, "-") {
		val := cleanValue(strings.TrimPrefix(trimmed, "-"))
		if val != "" {
			cfg.Convention.Prefixes = append(cfg.Convention.Prefixes, val)
		}
		return
	}
	cfg.Convention.Prefixes = parseStringList(trimmed)
}

func parseDriftSection(trimmed, currentSection string, cfg *Config) {
	if cfg.Drift.Exceptions == nil {
		cfg.Drift.Exceptions = make(map[string]map[string][]drift.ContextException)
	}

	if currentSection == "drift" {
		parseDriftRoot(trimmed, cfg)
		return
	}

	if strings.HasPrefix(currentSection, "drift.exceptions.") {
		parseDriftExceptions(trimmed, currentSection, cfg)
	}
}

func parseDriftRoot(trimmed string, cfg *Config) {
	k, v, found := strings.Cut(trimmed, ":")
	if !found {
		return
	}
	key := strings.TrimSpace(k)
	val := cleanValue(v)
	switch key {
	case "dominance_threshold":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.Drift.DominanceThreshold = f
		}
	case "outlier_threshold", "threshold_percentage":
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			cfg.Drift.OutlierThreshold = f
		}
	case "min_cluster_occurrences", "min_occurrences":
		if n, err := strconv.Atoi(val); err == nil {
			cfg.Drift.MinClusterOccurrences = n
		}
	}
}

func parseDriftExceptions(trimmed, currentSection string, cfg *Config) {
	rem := strings.TrimPrefix(currentSection, "drift.exceptions.")
	parts := strings.Split(rem, ".")
	if len(parts) == 1 {
		parseDriftExceptionCategoryList(parts[0], trimmed, cfg)
		return
	}
	if len(parts) == 2 && strings.HasPrefix(trimmed, "-") {
		comp := parts[0]
		cat := parts[1]
		val := cleanValue(strings.TrimPrefix(trimmed, "-"))
		if val != "" {
			addDriftException(cfg, comp, cat, drift.ContextException{Token: val})
		}
	}
}

func parseDriftExceptionCategoryList(comp, trimmed string, cfg *Config) {
	k, v, found := strings.Cut(trimmed, ":")
	if !found {
		return
	}
	cat := cleanValue(k)
	val := strings.TrimSpace(v)
	if val == "" {
		return
	}
	tokens := parseStringList(val)
	for _, tok := range tokens {
		addDriftException(cfg, comp, cat, drift.ContextException{Token: tok})
	}
}

func addDriftException(cfg *Config, comp, cat string, ex drift.ContextException) {
	if cfg.Drift.Exceptions[comp] == nil {
		cfg.Drift.Exceptions[comp] = make(map[string][]drift.ContextException)
	}
	cfg.Drift.Exceptions[comp][cat] = append(cfg.Drift.Exceptions[comp][cat], ex)
}

func parseStringList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		raw = raw[1 : len(raw)-1]
	}
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		cleaned := cleanValue(p)
		if cleaned != "" {
			result = append(result, cleaned)
		}
	}
	return result
}

// cleanValue membersihkan tanda kutip (single/double) dan spasi dari nilai konfigurasi.
func cleanValue(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			s = s[1 : len(s)-1]
		}
	}
	return strings.TrimSpace(s)
}

// ResolveActiveRules menerapkan 3-tier precedence:
// Registry (Base) -> CLI Scope (--rule/--category) -> Config Policy (charites.yaml).
func (c *Config) ResolveActiveRules(reg *rules.Registry, categoryFilter, ruleFilter string) []ActiveRule {
	if reg == nil {
		return nil
	}

	active := make([]ActiveRule, 0, reg.Count())

	// Iterasi deterministik (All() terurut leksikografis berdasarkan Rule.ID())
	for _, rule := range reg.All() {
		id := rule.ID()

		// 1. CLI Candidate Scope Filter
		if ruleFilter != "" && id != ruleFilter {
			continue
		}
		if categoryFilter != "" && rule.Category() != categoryFilter {
			continue
		}

		// 2. Config Policy Resolution
		effectiveSev, enabled := c.resolveRuleSeverity(rule)
		if !enabled {
			continue // Policy menonaktifkan rule (Policy mengalahkan CLI filter)
		}

		active = append(active, ActiveRule{
			Rule:              rule,
			EffectiveSeverity: effectiveSev,
		})
	}

	return active
}

func (c *Config) resolveRuleSeverity(rule rules.Rule) (ir.Severity, bool) {
	effectiveSev := rule.DefaultSeverity()
	if c == nil || c.Rules == nil {
		return effectiveSev, true
	}

	override, exists := c.Rules[rule.ID()]
	if !exists {
		return effectiveSev, true
	}

	val := strings.ToLower(strings.TrimSpace(override))
	if val == "off" || val == "false" || val == "disable" || val == "disabled" {
		return effectiveSev, false
	}

	switch val {
	case "error":
		effectiveSev = ir.SeverityError
	case "warn", "warning":
		effectiveSev = ir.SeverityWarn
	case "info":
		effectiveSev = ir.SeverityInfo
	}

	return effectiveSev, true
}

// IsTelemetryEnabled mengembalikan apakah pelaporan issue dan telemetri diizinkan.
// Variabel lingkungan CHARITES_TELEMETRY memiliki presedensi tertinggi atas berkas konfigurasi.
// Default bernilai true (opt-out).
func (c *Config) IsTelemetryEnabled() bool {
	if env := os.Getenv("CHARITES_TELEMETRY"); env != "" {
		norm := strings.ToLower(strings.TrimSpace(env))
		if norm == "false" || norm == "0" || norm == "off" || norm == "no" {
			return false
		}
		if norm == "true" || norm == "1" || norm == "on" || norm == "yes" {
			return true
		}
	}
	if c != nil && c.Telemetry != nil {
		return *c.Telemetry
	}
	return true
}
