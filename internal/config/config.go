package config

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/will2469/charites/internal/ir"
	"github.com/will2469/charites/internal/rules"
)

// ActiveRule membungkus rule singleton dengan EffectiveSeverity hasil resolusi konfigurasi.
// Mempertahankan invariant immutability dari registry singleton.
type ActiveRule struct {
	Rule              rules.Rule
	EffectiveSeverity ir.Severity
}

// Config merepresentasikan konfigurasi proyek dari charites.yaml.
type Config struct {
	Format   string            `json:"format" yaml:"format"`
	ScanPath string            `json:"scan_path" yaml:"scan_path"`
	Rules    map[string]string `json:"rules" yaml:"rules"`   // "rule-id": "off" | "warn" | "error" | "info"
	Ignore   []string          `json:"ignore" yaml:"ignore"` // Pola path tambahan
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
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	currentSection := ""

	for scanner.Scan() {
		line := scanner.Text()

		// Buang komentar dan spasi di ujung
		if idx := strings.IndexByte(line, '#'); idx != -1 {
			line = line[:idx]
		}
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Hitung indentasi spasi
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if indent == 0 {
			parseTopLevel(trimmed, cfg, &currentSection)
		} else {
			parseIndentedSection(trimmed, currentSection, cfg)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
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
		case "format":
			cfg.Format = val
		case "scan_path":
			cfg.ScanPath = val
		}
	}
}

func parseIndentedSection(trimmed, currentSection string, cfg *Config) {
	switch currentSection {
	case "rules":
		if k, v, found := strings.Cut(trimmed, ":"); found {
			ruleID := strings.TrimSpace(k)
			val := cleanValue(v)
			if ruleID != "" {
				cfg.Rules[ruleID] = val
			}
		}
	case "ignore":
		if strings.HasPrefix(trimmed, "-") {
			val := cleanValue(strings.TrimPrefix(trimmed, "-"))
			if val != "" {
				cfg.Ignore = append(cfg.Ignore, val)
			}
		}
	}
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
