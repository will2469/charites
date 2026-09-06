# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html).

---

## [v1.0.0-beta.1] - 2026-09-06

### Beta Evaluation Rationale & Field-Testing Focus
* **Real-World Empirical Validation:** Over 50% of advanced rule combinations, nested component hierarchies, and bespoke Astro/React architectural patterns have only been tested internally. Testing on real-world projects is required to evaluate practical effectiveness, diagnostic precision, and false positive (FP) / false negative (FN) rates under production conditions.
* **Deferred Scopes:** Certain expansion domains-notably `seo.*` (`SPEC-EXP-12-SEO`)-have been intentionally deferred to avoid superficial overlap with generic linters and keep the compiler focused on design token integrity, accessibility, and Core Web Vitals.
* **Community & Agent Feedback:** Field feedback and real-world false positive reports are gathered via GitHub Issues and the built-in MCP two-phase HITL tool (`charites_report_issue`).

### Added
* **Ultra-Fast Zero-CGO Static Analysis Compiler:**
  - High-performance Go 1.26 AST parsing engine for `.astro`, `.tsx`, `.jsx`, and `.css` files with sub-millisecond per-file traversal without Node.js runtime or CGO overhead.
  - Unified Intermediate Representation (Leaf IR) streaming Astro and TSX/JSX ASTs into a normalized node graph.
  - SSOT Multi-Format Design Token Engine parsing `global.css`, `index.css`, `@theme`, and `tokens.json` (W3C DTCG format) with directed graph cycle detection (`ErrCycleDetected`) and recursion budget limits.
  - Evidence-based token verification ("The Banana Test") guaranteeing zero false positives for untokenized custom utility classes when no semantic token is declared.
* **90 Canonical Quality Rules across 8 Domains:**
  - `theme.*` (32 rules): Design token enforcement, slash opacity elimination (`bg-primary/10` $\rightarrow$ `bg-primary-light`), dark mode elevation preservation, CSS Cascade Layer boundaries.
  - `a11y.*` (16 rules): WCAG 2.2 AA standards, mathematical relative color contrast calculation, form control label bindings, modal keyboard trap prevention, iOS Safari auto-zoom prevention.
  - `responsive.*` (17 rules): Apple HIG/WCAG touch target ergonomics ($\ge 44 \times 44\text{px}$), container queries (`@container`), dynamic viewport units (`dvh`/`svh`), responsive table overflow wrapping, mobile keyboard safe areas.
  - `lcp.*`, `cls.*`, `inp.*`, `performance.*` (25 rules): Core Web Vitals optimization including hero image priority (`fetchpriority="high"`), preload links, explicit media aspect-ratio dimensions, layout shift prevention, long task unyielding detection, and Astro island deferred hydration.
* **Pure Stateless Model Context Protocol (MCP) Server (2026-07-28 Standard):**
  - `charites_scan`: Workspace component static analysis with rich structured diagnostics and online wiki links.
  - `charites_explain_rule`: Returns complete 8-Pillars architectural rationale, risk taxonomy, non-compliant examples, and remediation guidance.
  - `charites_list_rules`: Dynamic discovery of all 90 registered rules, domains, and default severities.
  - `charites_report_issue`: Two-Phase Human-in-the-Loop (HITL) reporting tool with SHA-256 draft signatures (Phase 1) and user-verified submission (Phase 2).
* **Rich Multi-Format CLI Reporters:**
  - Inline ANSI terminal reporter with colorized source snippets and actionable remediation hints.
  - Machine-readable JSON streaming format (`--format=json`) with rule metadata and online doc URLs.
  - Markdown audit reporter (`--format=markdown` / `-o report.md`) with executive scorecards, category violation breakdowns, and direct links to online wiki documentation.
* **1-SSOT Tri-Corpus Testing Harness:**
  - 17-pattern adversarial test matrix (P1-P5 positive, N1-N5 negative, A1-A7 adversarial) across all rules.
  - Continuous Go 1.26 fuzzing suite with 14,000+ synthetic mutations verifying zero crashes, memory leaks, or panic hazards.
* **Self-Management & Installation:**
  - Automated in-place self-update (`charites update`) and uninstaller (`charites uninstall`).
  - Cross-platform installation via `go install`, curl installer (`install.sh`), and PowerShell (`install.ps1`).

---
