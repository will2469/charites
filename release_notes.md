# Release Notes - Charites v1.0.0-beta.1 (2026-09-06)

Welcome to **Charites v1.0.0-beta.1**, the initial public beta release of the ultra-fast, zero-CGO, zero-Node.js compile-time static analyzer and design token linter for **Astro**, **React TSX/JSX**, and **Tailwind CSS**.

---

## 1. Beta Evaluation Rationale & Field-Testing Focus

> [!NOTE]
> **Why is this release marked `v1.0.0-beta.1`?**
> Although Charites achieves **100% test passage** across its internal suite (including the canonical 1-SSOT Tri-Corpus, E2E CLI tests, and 14,000+ continuous Go 1.26 fuzzing mutations), this release is deliberately designated as **Beta** for the following reasons:
>
> 1. **Empirical Real-World Project Validation:** Over 50% of the advanced rule combinations, nested component structures, and bespoke Astro/React architectural patterns have not yet been stress-tested on diverse, production-scale monorepos outside this repository.
> 2. **False Positive & Precision Tuning:** Static analysis on modern frontend codebases can encounter unusual template string patterns, complex HOC wrappers, or proprietary UI wrappers. The primary objective of `v1.0.0-beta.1` is to measure real-world false positive (FP) and false negative (FN) rates directly in external repositories.
> 3. **Deferred Expansion Scopes:** Certain expansion domains have been intentionally deferred-such as the `seo.*` domain (`SPEC-EXP-12-SEO`)-to prevent superficial redundancy against established linters (e.g. HTMLHint/axe) and keep the core compiler laser-focused on design systems, accessibility, and Core Web Vitals.
>
> We encourage developers and teams to test `v1.0.0-beta.1` on their active Astro and React projects and report false-positives or edge cases via the built-in MCP tool `charites_report_issue` or directly on GitHub.

---

## 2. Highlighted Capabilities & Core Architecture

### Ultra-Fast Zero-CGO Static Analysis Compiler
- **Native Go 1.26 Architecture:** Sub-millisecond per-file AST traversal without Node.js runtime or CGO overhead.
- **Unified Intermediate Representation (Leaf IR):** Streams `.astro` and `.tsx`/`.jsx` files into a single normalized AST representation.
- **SSOT Multi-Format Design Token Engine:** Auto-discovers `global.css`, `index.css`, `@theme`, and `tokens.json` (W3C DTCG format). Constructs a directed token graph with cycle detection (`ErrCycleDetected`) and recursion budget limits.
- **The Banana Test:** Implements evidence-based token verification. Custom or untokenized utilities pass cleanly without false positives if no corresponding semantic token is defined in the token graph.

### 90 Canonical Rules across 8 Quality Domains
1. **Theme & Design Tokens (`theme.*` - 32 Rules):** Enforces design token consistency, eliminates hardcoded slash opacity modifiers (`bg-primary/10` $\rightarrow$ `bg-primary-light`), prevents dark mode elevation collapse, and safeguards CSS Cascade Layer integrity.
2. **Accessibility (`a11y.*` - 16 Rules):** Enforces WCAG 2.2 AA standards, calculates mathematical relative color contrast ratios without headless browsers, validates form control bindings, eliminates modal keyboard traps, and prevents iOS Safari auto-zoom hazards.
3. **Responsive Ergonomics (`responsive.*` - 17 Rules):** Enforces Apple HIG / WCAG touch target dimensions ($\ge 44 \times 44\text{px}$), modern container queries (`@container`), dynamic viewport units (`dvh`/`svh`), responsive table overflow wrapping, and mobile virtual keyboard clearance.
4. **Core Web Vitals & Runtime Performance (`lcp.*`, `cls.*`, `inp.*`, `performance.*` - 25 Rules):
   - **LCP:** Prioritizes hero assets (`fetchpriority="high"`), preloads delayed-discovery hero images, and prevents render-blocking head scripts.
   - **CLS:** Requires explicit dimensions and aspect ratios on media, reserves layout space for dynamic content, and eliminates layout-triggering animations.
   - **INP:** Detects render-blocking scripts, unyielded long tasks, cascading React Context domain coupling, and un-deferred Astro islands.

### Model Context Protocol (MCP 2026-07-28 Pure Stateless Standard)
Charites embeds a native, pure stateless MCP server for IDEs and AI coding agents (Claude Desktop, Cursor, Antigravity):
- `charites_scan`: Audits workspace components and provides structured diagnostics with online wiki links and remediation hints.
- `charites_explain_rule`: Returns complete 8-Pillars architectural documentation, risk taxonomy, non-compliant examples, and remediation advice.
- `charites_list_rules`: Discovers all available rules and metadata.
- `charites_report_issue`: **Cryptographic Two-Phase Human-in-the-Loop (HITL)** reporting tool that generates SHA-256 signed drafts (Phase 1) and submits verified issues via GitHub CLI or prefilled browser URLs upon explicit user confirmation (Phase 2).

### Multi-Format Reporters & CLI Tooling
- **Terminal Inline ANSI:** Clean, colorized diagnostic reports with source code snippets, column pointers, and actionable remediation tips.
- **Streaming JSON (`--format=json`):** Machine-readable diagnostic envelope including `doc_url` and rule metadata.
- **Markdown Audit Reports (`--format=markdown` / `-o report.md`):** Executive summary scorecards, category violation breakdowns, and direct links to online wiki documentation.
- **Self-Management:** In-place binary self-update (`charites update`) and clean uninstaller (`charites uninstall`).

### 1-SSOT Tri-Corpus Testing Harness
Every rule is rigorously validated against a 17-pattern matrix in `tests/correctness/<rule_id>/`:
- **Positive (P1-P5):** Obvious, indirect, helper-wrapped, nested, and aliased violations.
- **Negative (N1-N5):** Valid design tokens, explicit ignore directives, third-party libraries, semantic HTML, and untokenized custom values (Banana Test).
- **Adversarial (A1-A7):** Template literals, ternary expressions, spread properties, dynamic object classes, variable shadowing, and cyclic tokens.
- **Continuous Fuzzing:** Over 14,000 synthetic CSS and AST mutations tested with zero crashes or memory leaks.

---

## 3. Installation & Getting Started

### Via Go Toolchain (Recommended for Go Developers)
```bash
go install github.com/will2469/charites/cmd/charites@v1.0.0-beta.1
```

### Linux & macOS (Automated Script)
```bash
curl -fsSL https://raw.githubusercontent.com/will2469/charites/main/scripts/install.sh | bash
```

### Windows (PowerShell)
```powershell
irm https://raw.githubusercontent.com/will2469/charites/main/scripts/install.ps1 | iex
```

### Basic Usage
```bash
# Scan current workspace with default inline ANSI reporter
charites scan .

# Generate comprehensive Markdown audit report
charites scan -f markdown -o audit-report.md .

# Run MCP server for AI agent integration
charites mcp
```

---

## 4. Full Changelog

All notable commits and features leading up to this release are documented in [CHANGELOG.md](CHANGELOG.md).
