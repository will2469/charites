# Contributing to Charites

Thank you for your interest in improving Charites! We welcome contributions from developers passionate about frontend performance, design token integrity, accessibility (a11y), and static analysis compiler engineering.

---

## Code of Conduct

We are committed to providing a welcoming, inclusive, and harassment-free environment for everyone. Please read our [Code of Conduct](CODE_OF_CONDUCT.md) before contributing.

---

## Development Setup

### Prerequisites

- **Go 1.26 or higher**
- **Git**
- *(Optional)* **Node.js 20+** (if testing web framework integrations)

### Getting Started

```bash
# Clone the repository
git clone https://github.com/will2469/charites.git
cd charites

# Install local git hooks
make setup-hooks

# Run tests
make test

# Run full test suite with Go Race Detector
make test-race

# Run code linters & hygiene checks
make lint
```

---

## How to Add a New Analyzer Rule

Adding a new rule to Charites is automated through the rule scaffolding harness:

1. **Scaffold the Rule:**
   Run the automated generator from the repository root:
   ```bash
   ./.agents/skills/charites-rule-scaffold/scripts/scaffold_rule.sh theme hardcode-opacity-color HIGH "Detects hardcoded opacity color slash modifiers"
   ```
   This automatically creates:
   - Rule implementation: `internal/rules/<category>/<snake_slug>.go`
   - 1-SSOT Tri-Corpus test fixtures: `tests/correctness/<category>/<slug>/`

2. **Register the Rule:**
   Register the new rule constructor in `internal/rules/builtin.go`:
   ```go
   func init() {
       // ...
       MustRegister(theme.NewHardcodeOpacityColorRule())
   }
   ```

3. **Implement the AST Inspection & Documentation:**
   - Implement node inspection logic in `Evaluate(node *ir.Node) []ir.Diagnostic`.
   - Provide complete 8-Pillars documentation in `Doc() ir.RuleDocumentation`.

4. **Populate Tri-Corpus Fixtures:**
   - `positive/`: Non-compliant test files with true violations.
   - `negative/`: Compliant test files, valid tokens, and ignored lines (`charites:ignore <category>.<slug>`).
   - `adversarial/`: Stress testing edge cases (fractions, line heights, arbitrary values).

5. **Verify Tests:**
   ```bash
   go test -v ./tests/correctness/<category>/<slug>/...
   ```

6. **Compile Wiki Documentation (SSOT):**
   **Never edit markdown files manually.** Run the automated wiki generator:
   ```bash
   make wiki
   ```
   This deterministically compiles `wiki/Home.md`, `wiki/<category>.md`, and `wiki/<category>/<slug>.md` directly from your Go rule definition.

---

## Code Style & Architectural Invariants

- **Semgrep Canonical Identifiers:** Rule IDs must strictly follow `<category>.<slug>` format (e.g. `theme.hardcode-color`, `a11y.alt-text`). Arbitrary abbreviations like `txx` or `axx` are strictly prohibited.
- **Suppression Directives:** Rules must honor inline directives:
  - Astro/HTML: `<!-- charites:ignore <category>.<slug> [reason] -->`
  - TSX/JSX/JS: `// charites:ignore <category>.<slug> [reason]`
  - CSS: `/* charites:ignore <category>.<slug> [reason] */`
- **Anti-Fat Code (~250 Lines/File):** Keep rule analyzers modular, cohesive, and concise.
- **Zero-False-Positive Target:** Always evaluate real-world frontend idioms and design tokens to prevent developer alert fatigue.
- **Conventional Commits:** All commits must follow the [Conventional Commits](https://www.conventionalcommits.org/) standard:
  - `feat(theme): implement theme.hardcode-opacity-color rule`
  - `fix(parser): correct line span for Astro frontmatter expressions`
  - `docs(wiki): update 8-pillars documentation for a11y.alt-text`
  - `perf(ir): optimize sync.Pool allocator for AST traversal`
  - `test(corpus): add adversarial ternary tests for responsive rules`

---

## Submitting Pull Requests

1. Fork the repository and create your branch from `main`:
   ```bash
   git checkout -b feat/my-new-rule
   ```
2. Ensure all tests and linters pass:
   ```bash
   make test-race
   make lint
   ```
3. Push your branch and open a Pull Request using the standard PR template.
