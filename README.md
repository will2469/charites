# Charites: Compiler-Grade Static Analyzer & AST Linter for Modern Web

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-blue.svg)](https://golang.org/)
[![Target Ecosystem](https://img.shields.io/badge/Targets-Astro%20%7C%20React%20TSX%20%7C%20CSS-6366f1.svg)](https://github.com/will2469/charites)
[![Documentation](https://img.shields.io/badge/docs-Architecture%20Portal-green.svg)](docs/README.md)

**Charites** (`Χάριτες`) is an ultra-fast, compiler-grade static analyzer and AST linter built with pure Go 1.26+ for the modern web ecosystem (**Astro**, **React/TSX**, and **CSS / Tailwind CSS design tokens**).

Designed to eliminate aesthetic fatigue, accessibility failures, responsive regressions, and hardcoded styling anti-patterns, Charites bridges the gap between design tokens and production codebases. It enforces strict web invariants (*WCAG 2.2, OKLCH semantic palettes, 44x44px Fitts's Law touch ergonomics, and Zero-JS SSR defaults*) with zero runtime dependencies and sub-100ms repository-wide scanning latency.

---

## Key Architectural Principles

1. **Zero Runtime Dependency (Single Static Binary):**
   Compiled to a standalone static binary with `CGO_ENABLED=0`. Requires **no Node.js, no Go runtime, no Python**, and runs anywhere instantly.
2. **Sub-100ms Latency:**
   High-throughput *Leaf IR* AST traversal parses and evaluates thousands of frontend files in under 100 milliseconds.
3. **1-SSOT Tri-Corpus Adversarial Quality:**
   Every analyzer rule is verified against a 17-pattern matrix: *Positive (P1-P5)*, *Negative (N1-N5)*, and *Adversarial (A1-A7)* to guarantee zero false positives and evasion resistance.
4. **Semgrep Canonical Identifiers:**
   Standardized `<category>.<slug>` rule identifiers (e.g. `theme.hardcode-color`, `a11y.alt-text`) with seamless `charites:ignore` inline directive suppression.
5. **Native Model Context Protocol (MCP):**
   Built-in pure stateless MCP server (protocol revision 2026-07-28) providing instant static analysis intelligence (`charites_scan`, `charites_explain_rule`, `charites_list_rules`) directly to AI coding assistants and IDEs.

---

## Quickstart

### Installation

#### Linux & macOS (One-Line Automated Installer)

```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash
```

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/will2469/charites/main/scripts/install.ps1 | iex
```

#### Via Go Toolchain (Optional)

```bash
go install github.com/will2469/charites/cmd/charites@latest
```

_Or download standalone pre-compiled binaries directly from [GitHub Releases](https://github.com/will2469/charites/releases)._

---

### Basic Usage

```bash
# 1. Scan current repository for design system and accessibility violations
charites scan .

# 2. Scan specific frontend file directly
charites scan src/components/Header.astro

# 3. Filter scan by specific rule category (theme, a11y, responsive, perf, tailwind)
charites scan . --category=theme

# 4. Filter scan by specific file extension
charites scan . --ext=astro,tsx

# 5. Output machine-readable JSON stream for CI/CD pipelines
charites scan . --format=json

# 6. Launch native Model Context Protocol (MCP) server for AI assistants
charites mcp
```

---

## In-Code Suppression Directives

To legitimately suppress a warning on a specific line or block:

```astro
<!-- charites:ignore theme.hardcode-color Brand exception -->
<div style="color: #ff5722;">Brand Logo</div>
```

```tsx
// charites:ignore theme.hardcode-color Third-party widget
<span style={{ backgroundColor: "#123456" }}>Widget</span>
```

```css
/* charites:ignore theme.hardcode-color Legacy vendor token */
.custom-banner {
  background-color: #00bcd4;
}
```

---

## Documentation & Architecture Portal

The comprehensive architectural specifications are available in the [docs/](docs/README.md) portal:

* **[Specification (docs/01-spec/)](docs/01-spec/)**: AST Leaf IR, rule specifications, configuration engine, and CLI contracts.
* **[Architecture (docs/02-architecture/)](docs/02-architecture/)**: AST parsing pipeline, token evaluator, and memory pool optimization.
* **[Testing Strategy (docs/03-testing/)](docs/03-testing/)**: 1-SSOT Tri-Corpus golden fixtures and adoption matrices.
* **[Quality Gates (docs/04-quality/)](docs/04-quality/)**: Benchmarks, race condition invariants, and security guidelines.
* **[Release Management (docs/05-release/)](docs/05-release/)**: SemVer 2.0.0 criteria, changelogs, and packaging.
* **[Roadmap & Phases (docs/06-roadmap/)](docs/06-roadmap/)**: Detailed implementation roadmap from Phase 0 to Phase 8.

---

## Contributing & Community

* **Contributing Guidelines:** [CONTRIBUTING.md](CONTRIBUTING.md)
* **Code of Conduct:** [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)
* **Security Policy:** [SECURITY.md](SECURITY.md)

---

## License

Charites is licensed under the [MIT License](LICENSE).
